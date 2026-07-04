package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/ranker"
)

func (h *AdminHandler) CompareRecommendationDebugScenarioRecommendations(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugScenarioCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.MaxPerLesson <= 0 || req.MaxPerLesson > 12 {
		req.MaxPerLesson = 5
	}
	req.RecommendationQuery = strings.TrimSpace(req.RecommendationQuery)
	if len([]rune(req.RecommendationQuery)) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recommendation_query too long"})
		return
	}

	var lessonSnapshot recommendationDebugGeneratedLessonsSnapshot
	if err := json.Unmarshal([]byte(strings.TrimSpace(scenario.GeneratedLessonsSnapshot)), &lessonSnapshot); err != nil || len(lessonSnapshot.Lessons) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated lessons are required"})
		return
	}
	lessonsToCompare := lessonSnapshot.Lessons
	if strings.TrimSpace(req.LessonID) != "" || req.OrderIndex != nil {
		lessonsToCompare = nil
		for _, lesson := range lessonSnapshot.Lessons {
			if strings.TrimSpace(req.LessonID) != "" && lesson.LessonID == strings.TrimSpace(req.LessonID) {
				lessonsToCompare = append(lessonsToCompare, lesson)
				break
			}
			if req.OrderIndex != nil && lesson.OrderIndex == *req.OrderIndex {
				lessonsToCompare = append(lessonsToCompare, lesson)
				break
			}
		}
		if len(lessonsToCompare) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found in generated snapshot"})
			return
		}
	}

	startedAt := time.Now()
	rolloutState, err := h.loadActiveRecommendationRolloutState(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}
	rankerModelVersion := ""
	featureSchemaVersion := "ranker-feature-v1"
	if rolloutState != nil {
		rankerModelVersion = strings.TrimSpace(rolloutState.RankerModelVersion)
		if strings.TrimSpace(rolloutState.FeatureSchemaVersion) != "" {
			featureSchemaVersion = strings.TrimSpace(rolloutState.FeatureSchemaVersion)
		}
	}
	rankerScorer := ranker.NewNoOpScorer(rankerModelVersion, featureSchemaVersion)
	repo := curriculum.NewRepository(h.db)
	ownerUserID, err := h.resolveRecommendationDebugContentOwner(c.Request.Context(), adminID, scenario)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content owner user unavailable"})
		return
	}
	lessonResults := make([]gin.H, 0, len(lessonsToCompare))
	totalCandidates := 0
	totalLexicalHits := 0
	totalVectorHits := 0
	shadowAttempted := 0
	shadowSucceeded := 0
	shadowRouteAttempted := 0
	shadowRouteSucceeded := 0

	for _, lesson := range lessonsToCompare {
		querySource := firstNonEmpty(
			req.RecommendationQuery,
			lessonSnapshot.CourseTitle,
			scenario.CourseTitle,
			lessonSnapshot.InitialUserIntent,
		)
		queryBundle := buildRecommendationDebugQueryFromInput(
			querySource,
			firstNonEmpty(lessonSnapshot.DraftTitle, scenario.CourseTitle),
			lessonSnapshot.ConfirmedGoal,
			lesson.Title,
			lesson.Objective,
			"",
		)
		embeddingInfo, embeddingText := h.buildRecommendationDebugEmbedding(c.Request.Context(), queryBundle.DenseQuery)
		lexicalQuery := queryBundle.LexicalQueryExpanded
		if strings.TrimSpace(lexicalQuery) == "" {
			lexicalQuery = queryBundle.LexicalQueryKO
		}
		diagnostics, err := repo.SearchLessonCandidatesDiagnostics(c.Request.Context(), ownerUserID, lexicalQuery, embeddingText, req.MaxPerLesson, curriculum.LessonCandidateSearchOptions{
			PreferredLanguage: "ko",
			Now:               time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search recommendation candidates"})
			return
		}
		shadowStage, shadowEmbeddingText := h.buildRecommendationDebugShadowEmbedding(c.Request.Context(), queryBundle.DenseQuery)
		if attempted, _ := shadowStage["attempted"].(bool); attempted {
			shadowAttempted++
		}
		if used, _ := shadowStage["used"].(bool); used {
			shadowSucceeded++
		}
		explorerContext := curriculum.ExplorerRecommendationContext{
			SourceQuery:     firstNonEmpty(lessonSnapshot.CourseTitle, scenario.CourseTitle, lessonSnapshot.InitialUserIntent),
			LearningGoal:    lessonSnapshot.ConfirmedGoal,
			RegionTitle:     firstNonEmpty(lessonSnapshot.DraftTitle, scenario.CourseTitle),
			RegionObjective: lessonSnapshot.ConfirmedGoal,
			NodeTitle:       lesson.Title,
			NodeSummary:     lesson.Objective,
			FallbackQuery:   firstNonEmpty(scenario.CourseTitle, lessonSnapshot.CourseTitle, lessonSnapshot.InitialUserIntent),
		}
		externalQuery := curriculum.BuildExplorerRecommendationExternalSearchQueryWithPrimaryQuery(explorerContext, req.RecommendationQuery)
		if strings.TrimSpace(externalQuery) == "" {
			externalQuery = curriculum.BuildExternalSearchQuery(queryBundle)
		}
		externalStage := h.buildRecommendationDebugExternalStageWithQuery(c.Request.Context(), ownerUserID, req.MaxPerLesson, externalQuery, diagnostics.MergedCandidates)
		candidates := buildRecommendationDebugCandidateSummaries(diagnostics, nil, "ko", time.Now())
		externalSummaries := buildRecommendationDebugExternalCandidateSummariesWithContext(externalStage, req.MaxPerLesson, &explorerContext, req.RecommendationQuery)
		candidates = mergeRecommendationDebugCandidateSummaries(candidates, externalSummaries, req.MaxPerLesson)
		rankerBaseline := applyRecommendationDebugRankerScores(c.Request.Context(), candidates, rankerScorer)
		totalCandidates += len(candidates)
		totalLexicalHits += len(diagnostics.LexicalCandidates)
		totalVectorHits += len(diagnostics.VectorCandidates)
		shadowCandidates := []gin.H{}
		rankerShadow := applyRecommendationDebugRankerScores(c.Request.Context(), shadowCandidates, rankerScorer)
		shadowSearch := gin.H{
			"used":            false,
			"candidate_count": 0,
			"route":           "shadow_vector",
			"reason":          "primary_embedding_unavailable",
			"fallback_reason": "content embedding unavailable; primary route keeps current ordering",
			"ranker":          rankerShadow,
		}
		if shadowEmbeddingText != nil {
			shadowRouteAttempted++
			shadowProvider := strings.TrimSpace(fmt.Sprint(shadowStage["provider"]))
			shadowModel := strings.TrimSpace(fmt.Sprint(shadowStage["model"]))
			shadowDimension, _ := shadowStage["dimension"].(int)
			shadowVectorCandidates, err := repo.SearchShadowVectorLessonCandidates(c.Request.Context(), ownerUserID, shadowProvider, shadowModel, shadowDimension, shadowEmbeddingText, req.MaxPerLesson)
			if err != nil {
				shadowSearch = gin.H{
					"used":            false,
					"candidate_count": 0,
					"route":           "shadow_vector",
					"reason":          err.Error(),
					"fallback_reason": "primary vector route failed; primary route keeps current ordering",
				}
			} else {
				shadowCandidates = buildRecommendationDebugShadowCandidateSummaries(shadowVectorCandidates)
				rankerShadow = applyRecommendationDebugRankerScores(c.Request.Context(), shadowCandidates, rankerScorer)
				fallbackReason := ""
				if len(shadowCandidates) == 0 {
					fallbackReason = "ready primary vector candidates not found; primary route keeps current ordering"
				} else {
					shadowRouteSucceeded++
				}
				shadowSearch = gin.H{
					"used":            true,
					"candidate_count": len(shadowCandidates),
					"route":           "shadow_vector",
					"reason":          "content_embeddings queried",
					"fallback_reason": fallbackReason,
					"ranker":          rankerShadow,
				}
			}
		}
		shadowOverlap := recommendationDebugTopKOverlap(candidates, shadowCandidates, req.MaxPerLesson)

		lessonResults = append(lessonResults, gin.H{
			"lesson": gin.H{
				"lesson_id":            lesson.LessonID,
				"title":                lesson.Title,
				"objective":            lesson.Objective,
				"order_index":          lesson.OrderIndex,
				"recommendation_query": req.RecommendationQuery,
			},
			"baseline": gin.H{
				"stages": gin.H{
					"query": gin.H{
						"detected_intent":        curriculum.DetectExternalSearchIntent(queryBundle),
						"raw_query":              queryBundle.RawQuery,
						"dense_query":            queryBundle.DenseQuery,
						"lexical_query":          lexicalQuery,
						"lexical_query_ko":       queryBundle.LexicalQueryKO,
						"lexical_query_expanded": queryBundle.LexicalQueryExpanded,
						"tokens":                 queryBundle.Tokens,
						"expanded_terms":         queryBundle.ExpandedTokens,
						"dictionary_version":     recommendationDictionaryVersion,
						"recommendation_query":   req.RecommendationQuery,
					},
					"embedding": embeddingInfo,
					"search": gin.H{
						"max_per_lesson":    req.MaxPerLesson,
						"retrieval_limit":   diagnostics.RetrievalLimit,
						"lexical_hit_count": len(diagnostics.LexicalCandidates),
						"vector_hit_count":  len(diagnostics.VectorCandidates),
						"merged_hit_count":  len(diagnostics.MergedCandidates),
						"failure_reason": recommendationDebugFailureReason(
							queryBundle,
							embeddingInfo,
							len(diagnostics.LexicalCandidates),
							len(diagnostics.VectorCandidates),
							len(diagnostics.MergedCandidates),
						),
					},
					"external": externalStage,
					"ranker":   rankerBaseline,
				},
				"candidates": candidates,
			},
			"shadow": gin.H{
				"enabled":    shadowStage["used"],
				"reason":     shadowStage["reason"],
				"embedding":  shadowStage,
				"search":     shadowSearch,
				"candidates": shadowCandidates,
				"ranker":     rankerShadow,
			},
			"metrics": gin.H{
				"candidate_count": len(candidates),
				"shadow_overlap":  shadowOverlap,
			},
		})
	}

	latencyMS := int(time.Since(startedAt).Milliseconds())
	compareSnapshot := gin.H{
		"scenario_id":            scenario.ID,
		"course_title":           scenario.CourseTitle,
		"confirmed_goal":         lessonSnapshot.ConfirmedGoal,
		"max_per_lesson":         req.MaxPerLesson,
		"lesson_count":           len(lessonResults),
		"total_candidates":       totalCandidates,
		"total_lexical_hits":     totalLexicalHits,
		"total_vector_hits":      totalVectorHits,
		"baseline_provider":      "embedding_gemma",
		"ranker_provider":        rankerScorer.Provider(),
		"ranker_noop":            rankerScorer.Provider() == "noop",
		"ranker_model_version":   rankerModelVersion,
		"feature_schema_version": featureSchemaVersion,
		"shadow_enabled":         shadowSucceeded > 0,
		"shadow_attempted":       shadowAttempted,
		"shadow_succeeded":       shadowSucceeded,
		"shadow_route":           "shadow_vector",
		"shadow_route_attempted": shadowRouteAttempted,
		"shadow_route_succeeded": shadowRouteSucceeded,
		"lessons":                lessonResults,
		"latency_ms":             latencyMS,
		"created_at":             time.Now(),
	}
	compareJSON, err := json.Marshal(compareSnapshot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode recommendation comparison"})
		return
	}
	requestSnapshot, _ := json.Marshal(gin.H{
		"max_per_lesson":       req.MaxPerLesson,
		"lesson_count":         len(lessonsToCompare),
		"lesson_id":            strings.TrimSpace(req.LessonID),
		"order_index":          req.OrderIndex,
		"recommendation_query": req.RecommendationQuery,
		"owner_user_id":        ownerUserID.String(),
	})
	metricsSnapshot, _ := json.Marshal(gin.H{
		"lesson_count":           len(lessonResults),
		"total_candidates":       totalCandidates,
		"total_lexical_hits":     totalLexicalHits,
		"total_vector_hits":      totalVectorHits,
		"ranker_provider":        rankerScorer.Provider(),
		"ranker_noop":            rankerScorer.Provider() == "noop",
		"ranker_model_version":   rankerModelVersion,
		"feature_schema_version": featureSchemaVersion,
		"shadow_enabled":         shadowSucceeded > 0,
		"shadow_attempted":       shadowAttempted,
		"shadow_succeeded":       shadowSucceeded,
		"shadow_route":           "shadow_vector",
		"shadow_route_attempted": shadowRouteAttempted,
		"shadow_route_succeeded": shadowRouteSucceeded,
	})
	providerSnapshot, _ := json.Marshal(gin.H{
		"baseline": gin.H{
			"provider":  "embedding_gemma",
			"dimension": 768,
		},
		"shadow": gin.H{
			"provider": "embedding_gemma",
			"route":    "shadow_vector",
		},
		"ranker": gin.H{
			"provider":               rankerScorer.Provider(),
			"model_version":          rankerScorer.ModelVersion(),
			"feature_schema_version": rankerScorer.FeatureSchemaVersion(),
			"enabled":                true,
			"noop":                   rankerScorer.Provider() == "noop",
			"reason":                 "no-op ranker scoring fallback is active; route order is preserved",
		},
	})

	scenarioUUID, _ := uuid.Parse(scenario.ID)
	tx, err := h.db.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recommendation comparison"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	if _, err := tx.Exec(c.Request.Context(), `
		UPDATE recommendation_debug_scenarios
		SET status = 'recommendation_tested',
		    updated_at = NOW()
		WHERE id = $1
	`, scenarioUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recommendation comparison"})
		return
	}
	var runID string
	if err := tx.QueryRow(c.Request.Context(), `
		INSERT INTO recommendation_debug_runs (
			scenario_id,
			created_by,
			run_type,
			request_snapshot,
			baseline_result_snapshot,
			shadow_result_snapshot,
			metrics_snapshot,
			provider_snapshot,
			latency_ms
		)
		VALUES ($1, $2, 'recommendation_compare', $3::jsonb, $4::jsonb, '{}'::jsonb, $5::jsonb, $6::jsonb, $7)
		RETURNING id::text
	`, scenarioUUID, recommendationDebugNullableAdminID(adminID), string(requestSnapshot), string(compareJSON), string(metricsSnapshot), string(providerSnapshot), latencyMS).Scan(&runID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recommendation comparison run"})
		return
	}
	if err := tx.Commit(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recommendation comparison"})
		return
	}

	compareSnapshot["run_id"] = runID
	c.JSON(http.StatusOK, gin.H{"comparison": compareSnapshot})
}
