package curriculum

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) SearchDraftLessons(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req SearchCourseDraftLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.MaxPerLesson <= 0 || req.MaxPerLesson > 10 {
		req.MaxPerLesson = 5
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	cost, err := h.repo.GetPointSetting(c.Request.Context(), "lesson_rec_cost")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point policy"})
		return
	}

	freeBalance, paidBalance, err := h.repo.GetPointBalances(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point balance"})
		return
	}

	execution, err := h.resolveLessonSearchExecution(c.Request.Context(), userID, cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve search execution"})
		return
	}

	if execution.ShouldCharge && freeBalance+paidBalance < execution.EffectiveCost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
		return
	}

	for _, mainLesson := range draft.Lessons {
		for _, subLesson := range mainLesson.SubLessons {
			if req.LessonID != nil && subLesson.Lesson.ID != *req.LessonID {
				continue
			}

			queryBundle := buildLessonSearchQuery(draft.Draft, mainLesson.Lesson, subLesson.Lesson)
			var embeddingText *string
			embedClient, embedErr := newPrimaryEmbeddingClient()
			if embedErr == nil {
				embedCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
				vector, embErr := embedClient.Embed(embedCtx, queryBundle.DenseQuery)
				cancel()
				if embErr == nil && len(vector) > 0 {
					pgvec := float32VectorToPGVector(vector)
					embeddingText = &pgvec
				} else if embErr != nil {
					log.Printf("[recommendation.embedding] provider=%s model=%s failed err=%v", primaryEmbeddingProvider(), primaryEmbeddingModel(), embErr)
				}
			} else {
				log.Printf("[recommendation.embedding] client unavailable provider=%s model=%s err=%v", primaryEmbeddingProvider(), primaryEmbeddingModel(), embedErr)
			}

			lexicalQuery := queryBundle.LexicalQueryExpanded
			if strings.TrimSpace(lexicalQuery) == "" {
				lexicalQuery = queryBundle.LexicalQueryKO
			}
			candidates, err := h.repo.SearchLessonCandidates(c.Request.Context(), userID, lexicalQuery, embeddingText, req.MaxPerLesson, LessonCandidateSearchOptions{
				PreferredFormat:   draft.Draft.PreferredFormat,
				PreferredLanguage: "ko",
				Now:               time.Now(),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search lesson candidates"})
				return
			}
			if len(candidates) < req.MaxPerLesson {
				externalCandidates, extErr := h.searchExternalLessonCandidates(
					c.Request.Context(),
					userID,
					req.MaxPerLesson-len(candidates),
					queryBundle,
					draft.Draft.PreferredFormat,
				)
				if extErr == nil && len(externalCandidates) > 0 {
					candidates = mergeExternalCandidates(candidates, externalCandidates, req.MaxPerLesson)
				}
			}
			if err := h.repo.ReplaceLessonCandidateResources(c.Request.Context(), draft.Draft.ID, subLesson.Lesson.ID, candidates, queryBundle.RawQuery); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save lesson candidates"})
				return
			}
		}
	}

	if execution.ShouldCharge {
		if err := h.repo.ChargePointsWithDetail(c.Request.Context(), userID, execution.EffectiveCost, "use_lesson_rec", &PointTransactionDetail{
			Feature:       "lesson_candidate_search",
			ReferenceType: "course_draft",
			ReferenceID:   &draft.Draft.ID,
			Description:   strings.TrimSpace(draft.Draft.Title),
			Metadata: map[string]any{
				"query": strings.TrimSpace(draft.Draft.SourceQuery),
			},
		}); err != nil {
			if err.Error() == "insufficient_points" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to charge points"})
			return
		}
	}

	updatedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draft": updatedDraft,
		"point_preview": CourseDraftPointPreview{
			Cost:         execution.EffectiveCost,
			FreeBalance:  freeBalance,
			PaidBalance:  paidBalance,
			TotalBalance: freeBalance + paidBalance,
		},
		"billing_status": execution.BillingStatus,
		"policy_cost":    cost,
	})
}
