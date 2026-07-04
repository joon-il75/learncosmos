package curriculum

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func resolveExplorerRecommendationLanguage(userLanguage string, draft *CourseDraft) string {
	if draft != nil && strings.TrimSpace(draft.GenerationLanguage) != "" {
		return normalizeLearningLanguage(draft.GenerationLanguage)
	}
	return normalizeLearningLanguage(userLanguage)
}

func (h *Handler) RecommendExplorerContent(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req RecommendExplorerContentRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}
	if req.MaxResults <= 0 || req.MaxResults > 6 {
		req.MaxResults = 6
	}
	exclusion := buildExplorerRecommendationExclusion(req)
	learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	recommendationLanguage := learningLanguage

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

	contextSourceQuery := strings.TrimSpace(req.Query)
	contextLevelTitle := strings.Join(compactExplorerRecommendationParts(
		req.RegionTitle,
		req.SubRegionTitle,
	), " ")
	contextLevelObjective := strings.Join(compactExplorerRecommendationParts(
		req.RegionDescription,
		req.SubRegionDescription,
	), " ")
	contextLessonTitle := strings.TrimSpace(req.NodeTitle)
	contextLessonSummary := strings.TrimSpace(req.NodeSummary)
	recommendationGoalMeta := RecommendationGoalMetadata{}
	if contextLessonTitle != "" && contextLessonTitle == strings.TrimSpace(req.Query) {
		// New point recommendation sometimes sends the visible search query as node_title.
		// Treat that as synthetic text and rely on region/course context instead.
		contextLessonTitle = ""
		contextLessonSummary = ""
	}

	if draftIDText := strings.TrimSpace(req.CourseDraftID); draftIDText != "" {
		draftID, parseErr := uuid.Parse(draftIDText)
		if parseErr == nil {
			draft, draftErr := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
			if draftErr == nil {
				recommendationLanguage = resolveExplorerRecommendationLanguage(learningLanguage, &draft.Draft)
				recommendationReq := CreateCourseDraftRequest{
					SourceQuery:      draft.Draft.SourceQuery,
					LearningGoal:     draft.Draft.LearningGoal,
					GoalProfileID:    draft.Draft.GoalProfileID,
					CurrentLevel:     draft.Draft.CurrentLevel,
					PreferredFormat:  draft.Draft.PreferredFormat,
					LearningLanguage: recommendationLanguage,
				}
				activeGoal, goalErr := h.goalRepo.GetActiveGoal(c.Request.Context(), draftID)
				if goalErr != nil {
					activeGoal, goalErr = h.goalRepo.GetActiveGoalByUser(c.Request.Context(), userID)
				}
				if goalErr == nil && activeGoal != nil {
					recommendationReq.GoalMotivation = activeGoal.Motivation
					recommendationReq.GoalUsageContext = activeGoal.UsageContext
					recommendationReq.GoalDifficultyLevel = activeGoal.DifficultyLevel
					recommendationReq.GoalTimeHorizon = activeGoal.TimeHorizon
					recommendationReq.GoalOutputType = activeGoal.OutputType
					recommendationReq.GoalType = activeGoal.GoalType
					if activeGoal.ConfirmedGoal != nil && strings.TrimSpace(derefString(recommendationReq.LearningGoal)) == "" {
						recommendationReq.LearningGoal = activeGoal.ConfirmedGoal
					}
				} else {
					log.Printf("[explorer/recommend] no_active_goal fallback: draft_id=%s", draftID.String())
				}
				recommendationGoalMeta = BuildRecommendationGoalMetadata(recommendationReq)
				contextSourceQuery = strings.Join(compactExplorerRecommendationParts(
					draft.Draft.SourceQuery,
					derefString(recommendationReq.LearningGoal),
					recommendationGoalMeta.QueryHint,
				), " ")
			}
		}
	}
	lessonSearchSpec := LessonRecommendationSearchSpec{}
	resolvedLessonIDText := ""
	resolutionSource := ""
	if pointIDText := strings.TrimSpace(req.PointID); pointIDText != "" {
		if pointID, parseErr := uuid.Parse(pointIDText); parseErr == nil {
			searchContext, specErr := h.repo.GetLessonRecommendationSearchContextByPointID(c.Request.Context(), userID, pointID)
			if specErr == nil {
				lessonSearchSpec = searchContext.Spec
				if strings.TrimSpace(lessonSearchSpec.PrimaryQuery) == "" {
					lessonSearchSpec = searchContext.FallbackSpec()
					if saveErr := h.repo.SaveFallbackLessonRecommendationSearchSpecByPointID(c.Request.Context(), userID, pointID, lessonSearchSpec); saveErr != nil {
						log.Printf("[explorer/recommend] lesson_search_spec_fallback_save_failed point_id=%s err=%v", pointID.String(), saveErr)
					}
				}
			} else {
				log.Printf("[explorer/recommend] lesson_search_spec_unavailable point_id=%s err=%v", pointID.String(), specErr)
			}
		}
	}
	if strings.TrimSpace(lessonSearchSpec.PrimaryQuery) == "" {
		if lessonIDText := strings.TrimSpace(req.LessonID); lessonIDText != "" {
			if lessonID, parseErr := uuid.Parse(lessonIDText); parseErr == nil {
				searchContext, specErr := h.repo.GetLessonRecommendationSearchContextByLessonID(c.Request.Context(), userID, lessonID)
				resolvedLessonID := lessonID
				resolutionSource = "direct_lesson"
				if specErr != nil {
					linkedLessonID, resolveErr := h.repo.ResolveDraftLessonIDByExplorerSubRegionID(c.Request.Context(), userID, lessonID)
					if resolveErr == nil {
						resolvedLessonID = linkedLessonID
						resolutionSource = "explorer_subregion"
						searchContext, specErr = h.repo.GetLessonRecommendationSearchContextByLessonID(c.Request.Context(), userID, resolvedLessonID)
						log.Printf("[explorer/recommend] lesson_search_spec_resolved_subregion lesson_id=%s resolved_lesson_id=%s", lessonID.String(), resolvedLessonID.String())
					} else {
						log.Printf("[explorer/recommend] lesson_search_spec_subregion_resolve_unavailable lesson_id=%s err=%v", lessonID.String(), resolveErr)
					}
				}
				if specErr != nil {
					linkedLessonID, resolveErr := h.repo.ResolveDraftLessonIDByExplorerRegionID(c.Request.Context(), userID, lessonID)
					if resolveErr == nil {
						resolvedLessonID = linkedLessonID
						resolutionSource = "explorer_region"
						searchContext, specErr = h.repo.GetLessonRecommendationSearchContextByLessonID(c.Request.Context(), userID, resolvedLessonID)
						log.Printf("[explorer/recommend] lesson_search_spec_resolved_region lesson_id=%s resolved_lesson_id=%s", lessonID.String(), resolvedLessonID.String())
					} else {
						log.Printf("[explorer/recommend] lesson_search_spec_region_resolve_unavailable lesson_id=%s err=%v", lessonID.String(), resolveErr)
					}
				}
				if specErr == nil {
					resolvedLessonIDText = resolvedLessonID.String()
					lessonSearchSpec = searchContext.Spec
					if strings.TrimSpace(lessonSearchSpec.PrimaryQuery) == "" {
						lessonSearchSpec = searchContext.FallbackSpec()
						if saveErr := h.repo.SaveFallbackLessonRecommendationSearchSpecByLessonID(c.Request.Context(), userID, resolvedLessonID, lessonSearchSpec); saveErr != nil {
							log.Printf("[explorer/recommend] lesson_search_spec_lesson_fallback_save_failed lesson_id=%s resolved_lesson_id=%s err=%v", lessonID.String(), resolvedLessonID.String(), saveErr)
						}
					}
				} else {
					log.Printf("[explorer/recommend] lesson_search_spec_unavailable lesson_id=%s err=%v", lessonID.String(), specErr)
				}
			}
		}
	}
	effectiveRecommendationQuery := strings.TrimSpace(req.Query)
	if strings.TrimSpace(lessonSearchSpec.PrimaryQuery) != "" {
		effectiveRecommendationQuery = strings.TrimSpace(lessonSearchSpec.PrimaryQuery)
	}
	req.QueryLanguage = recommendationLanguage

	var embedClient embedder.Embedder
	client, embedErr := newPrimaryEmbeddingClient()
	if embedErr == nil {
		embedClient = client
	} else {
		log.Printf("[explorer/recommend] embedding client unavailable provider=%s model=%s err=%v", primaryEmbeddingProvider(), primaryEmbeddingModel(), embedErr)
	}

	searchStages := buildExplorerSearchStages(
		normalizer.BuildLessonSearchQuery("", "", "", contextLessonTitle, effectiveRecommendationQuery, contextLessonSummary),
		normalizer.BuildLessonSearchQuery("", contextLevelTitle, contextLevelObjective, "", "", ""),
		normalizer.BuildLessonSearchQuery(contextSourceQuery, "", "", "", "", ""),
	)
	log.Printf(
		"[explorer/recommend] query=%q effective_query=%q query_language=%q draft=%q region=%q subregion=%q node=%q source=%q search_spec_source=%q pattern=%q subpattern=%q goal_hint=%q stage_count=%d",
		strings.TrimSpace(req.Query),
		effectiveRecommendationQuery,
		recommendationLanguage,
		strings.TrimSpace(req.CourseDraftID),
		contextLevelTitle,
		strings.TrimSpace(req.SubRegionTitle),
		contextLessonTitle,
		contextSourceQuery,
		strings.TrimSpace(lessonSearchSpec.Source),
		recommendationGoalMeta.PatternKey,
		recommendationGoalMeta.SubpatternKey,
		recommendationGoalMeta.QueryHint,
		len(searchStages),
	)
	if !exclusion.empty() {
		log.Printf(
			"[explorer/recommend] exclusion content_ids=%d urls=%d",
			len(exclusion.contentIDs),
			len(exclusion.urlKeys),
		)
	}
	for _, stage := range searchStages {
		log.Printf(
			"[explorer/recommend] stage=%s lexical=%q tokens=%v strong=%v min_strong=%d",
			stage.name,
			stage.bundle.LexicalQueryExpanded,
			stage.bundle.Tokens,
			stage.filter.strongTokens,
			stage.filter.minStrongMatch,
		)
	}

	candidates, err := h.searchExplorerInternalCandidatesByStage(
		c.Request.Context(),
		userID,
		searchStages,
		req.MaxResults,
		LessonCandidateSearchOptions{PreferredLanguage: recommendationLanguage, Now: time.Now()},
		embedClient,
		execution.Provider,
		execution.BillingStatus,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search candidates"})
		return
	}
	candidates = filterExcludedExplorerCandidates(candidates, exclusion)
	log.Printf("[explorer/recommend] internal_candidates=%d", len(candidates))

	externalQueryText := buildExplorerExternalSearchQueryWithPrimaryQuery(searchStages, effectiveRecommendationQuery, contextSourceQuery, strings.TrimSpace(req.Query))
	externalBundle := normalizer.BuildLessonSearchQuery(externalQueryText, "", "", "", "", "")
	log.Printf("[explorer/recommend] external_query=%q", externalQueryText)

	externalCandidates, extErr := h.searchExternalLessonCandidatesWithQueryFiltered(
		c.Request.Context(),
		20,
		externalQueryText,
		externalBundle,
		nil,
		recommendationLanguage,
		exclusion.blockedExternalResult,
	)
	if extErr == nil && len(externalCandidates) > 0 {
		log.Printf("[explorer/recommend] external_candidates=%d", len(externalCandidates))
		candidates = mergeExternalCandidates(candidates, externalCandidates, max(req.MaxResults*4, req.MaxResults+10))
	}
	candidates = filterExcludedExplorerCandidates(candidates, exclusion)
	candidates = rerankExplorerCandidatesWithPrimaryQuery(candidates, searchStages, effectiveRecommendationQuery, req.MaxResults)
	candidates = rerankExplorerCandidatesWithSearchSpec(candidates, lessonSearchSpec, req.MaxResults)
	candidates = filterExcludedExplorerCandidates(candidates, exclusion)
	for idx, candidate := range candidates {
		log.Printf(
			"[explorer/recommend] result[%d]=%q score=%.4f url=%q",
			idx,
			candidate.Title,
			candidate.RankScore,
			derefString(candidate.ExternalURL),
		)
	}

	if execution.ShouldCharge {
		var referenceID *uuid.UUID
		if draftIDText := strings.TrimSpace(req.CourseDraftID); draftIDText != "" {
			if parsed, parseErr := uuid.Parse(draftIDText); parseErr == nil {
				referenceID = &parsed
			}
		}
		if err := h.repo.ChargePointsWithDetail(c.Request.Context(), userID, execution.EffectiveCost, "use_lesson_rec", &PointTransactionDetail{
			Feature:       "explorer_recommendation_search",
			ReferenceType: "course_draft",
			ReferenceID:   referenceID,
			Description:   strings.TrimSpace(req.Query),
			Metadata: map[string]any{
				"query_language":     recommendationLanguage,
				"point_id":           strings.TrimSpace(req.PointID),
				"lesson_id":          strings.TrimSpace(req.LessonID),
				"resolved_lesson_id": resolvedLessonIDText,
				"resolution_source":  resolutionSource,
				"effective_query":    effectiveRecommendationQuery,
				"search_spec_source": strings.TrimSpace(lessonSearchSpec.Source),
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

	items := make([]RecommendExplorerContentItem, 0, len(candidates))
	for _, cand := range candidates {
		item := RecommendExplorerContentItem{
			Title:        cand.Title,
			Description:  cand.Description,
			ThumbnailURL: cand.ThumbnailURL,
			URL:          cand.ExternalURL,
			ContentType:  cand.ContentType,
			RankScore:    cand.RankScore,
		}
		if cand.ContentID != nil {
			id := cand.ContentID.String()
			item.ContentID = &id
		}
		items = append(items, item)
	}

	rolloutDebug := h.buildRecommendationRolloutReadOnlyDebug(c.Request.Context(), userID)
	h.recordRecommendationRolloutExposure(c.Request.Context(), userID, req, items, rolloutDebug)

	c.JSON(http.StatusOK, gin.H{
		"candidates":     items,
		"billing_status": execution.BillingStatus,
		"query_language": recommendationLanguage,
		"rollout_debug":  rolloutDebug,
		"point_preview": CourseDraftPointPreview{
			Cost:         execution.EffectiveCost,
			FreeBalance:  freeBalance,
			PaidBalance:  paidBalance,
			TotalBalance: freeBalance + paidBalance,
		},
	})
}
