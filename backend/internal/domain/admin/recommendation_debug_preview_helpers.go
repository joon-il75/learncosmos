package admin

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
)

func (h *AdminHandler) listRecommendationDebugRuns(ctx context.Context, scenarioID uuid.UUID) ([]recommendationDebugRunRow, error) {
	rows, err := h.db.Query(ctx, `
		SELECT
			id::text,
			scenario_id::text,
			COALESCE(created_by::text, ''),
			run_type,
			request_snapshot::text,
			baseline_result_snapshot::text,
			shadow_result_snapshot::text,
			feature_snapshot::text,
			metrics_snapshot::text,
			provider_snapshot::text,
			latency_ms,
			created_at
		FROM recommendation_debug_runs
		WHERE scenario_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]recommendationDebugRunRow, 0)
	for rows.Next() {
		var item recommendationDebugRunRow
		if err := rows.Scan(
			&item.ID,
			&item.ScenarioID,
			&item.CreatedBy,
			&item.RunType,
			&item.RequestSnapshot,
			&item.BaselineResultSnapshot,
			&item.ShadowResultSnapshot,
			&item.FeatureSnapshot,
			&item.MetricsSnapshot,
			&item.ProviderSnapshot,
			&item.LatencyMS,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (h *AdminHandler) listRecommendationDebugLabels(ctx context.Context, scenarioID uuid.UUID) ([]recommendationDebugLabelRow, error) {
	rows, err := h.db.Query(ctx, `
		SELECT
			id::text,
			scenario_id::text,
			COALESCE(run_id::text, ''),
			COALESCE(created_by::text, ''),
			candidate_key,
			COALESCE(content_id::text, ''),
			url,
			label,
			note,
			feature_snapshot::text,
			baseline_rank,
			shadow_rank,
			created_at,
			updated_at
		FROM recommendation_debug_labels
		WHERE scenario_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]recommendationDebugLabelRow, 0)
	for rows.Next() {
		var item recommendationDebugLabelRow
		if err := rows.Scan(
			&item.ID,
			&item.ScenarioID,
			&item.RunID,
			&item.CreatedBy,
			&item.CandidateKey,
			&item.ContentID,
			&item.URL,
			&item.Label,
			&item.Note,
			&item.FeatureSnapshot,
			&item.BaselineRank,
			&item.ShadowRank,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (h *AdminHandler) ListRecommendationDebugDrafts(c *gin.Context) {
	repo := curriculum.NewRepository(h.db)
	drafts, err := repo.ListDraftsForAdmin(c.Request.Context(), 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list drafts"})
		return
	}

	items := make([]gin.H, 0, len(drafts))
	for _, draft := range drafts {
		items = append(items, gin.H{
			"id":           draft.ID,
			"user_id":      draft.UserID,
			"title":        draft.Title,
			"source_query": draft.SourceQuery,
			"status":       draft.Status,
			"updated_at":   draft.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"drafts": items})
}

func (h *AdminHandler) GetRecommendationDebugDraft(c *gin.Context) {
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}

	repo := curriculum.NewRepository(h.db)
	draft, err := repo.GetDraftByIDForAdmin(c.Request.Context(), draftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	lessons := make([]gin.H, 0)
	for _, mainTree := range draft.Lessons {
		for _, subTree := range mainTree.SubLessons {
			queryBundle, _ := h.buildRecommendationDebugLessonQuery(c.Request.Context(), draft.Draft, mainTree.Lesson, subTree.Lesson)
			candidateCount := 0
			for _, pt := range subTree.Points {
				if pt.Point.PointType == curriculum.PointTypeExploration &&
					pt.Point.SelectionState != nil &&
					*pt.Point.SelectionState == curriculum.ResourceSelectionCandidate {
					candidateCount++
				}
			}
			lessons = append(lessons, gin.H{
				"lesson_id":                subTree.Lesson.ID,
				"lesson_title":             subTree.Lesson.Title,
				"lesson_source_type":       subTree.Lesson.SourceType,
				"group_title":              mainTree.Lesson.Title,
				"generated_query":          queryBundle.RawQuery,
				"existing_candidate_count": candidateCount,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"draft": gin.H{
			"id":                   draft.Draft.ID,
			"user_id":              draft.Draft.UserID,
			"title":                draft.Draft.Title,
			"source_query":         draft.Draft.SourceQuery,
			"learning_goal":        derefString(draft.Draft.LearningGoal),
			"goal_profile_id":      draft.Draft.GoalProfileID,
			"goal_profile_version": draft.Draft.GoalProfileVersion,
			"status":               draft.Draft.Status,
			"group_count":          len(draft.Lessons),
		},
		"lessons": lessons,
	})
}

func (h *AdminHandler) PreviewRecommendationDebug(c *gin.Context) {
	var req recommendationDebugPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draftID, err := uuid.Parse(strings.TrimSpace(req.DraftID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft_id"})
		return
	}
	lessonID, err := uuid.Parse(strings.TrimSpace(req.LessonID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson_id"})
		return
	}
	if req.MaxPerLesson <= 0 || req.MaxPerLesson > 10 {
		req.MaxPerLesson = 5
	}

	repo := curriculum.NewRepository(h.db)
	draft, err := repo.GetDraftByIDForAdmin(c.Request.Context(), draftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	mainLessonTree, subLessonTree, ok := findRecommendationDebugLesson(draft, lessonID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		return
	}

	queryBundle, goalMeta := h.buildRecommendationDebugLessonQuery(c.Request.Context(), draft.Draft, mainLessonTree.Lesson, subLessonTree.Lesson)
	embeddingInfo, embeddingText := h.buildRecommendationDebugEmbedding(c.Request.Context(), queryBundle.DenseQuery)
	lexicalQuery := queryBundle.LexicalQueryExpanded
	if strings.TrimSpace(lexicalQuery) == "" {
		lexicalQuery = queryBundle.LexicalQueryKO
	}
	diagnostics, err := repo.SearchLessonCandidatesDiagnostics(c.Request.Context(), draft.Draft.UserID, lexicalQuery, embeddingText, req.MaxPerLesson, curriculum.LessonCandidateSearchOptions{
		PreferredFormat:   draft.Draft.PreferredFormat,
		PreferredLanguage: "ko",
		Now:               time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search lesson candidates"})
		return
	}
	externalStage := h.buildRecommendationDebugExternalStage(c.Request.Context(), draft.Draft.UserID, req.MaxPerLesson, queryBundle, diagnostics.MergedCandidates)
	candidateSummaries := buildRecommendationDebugCandidateSummaries(
		diagnostics,
		draft.Draft.PreferredFormat,
		"ko",
		time.Now(),
	)

	c.JSON(http.StatusOK, gin.H{
		"draft": gin.H{
			"id":           draft.Draft.ID,
			"user_id":      draft.Draft.UserID,
			"title":        draft.Draft.Title,
			"source_query": draft.Draft.SourceQuery,
		},
		"lesson": gin.H{
			"lesson_id":          subLessonTree.Lesson.ID,
			"lesson_title":       subLessonTree.Lesson.Title,
			"lesson_source_type": subLessonTree.Lesson.SourceType,
			"group_title":        mainLessonTree.Lesson.Title,
		},
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
				"goal_pattern_key":       goalMeta.PatternKey,
				"goal_subpattern_key":    goalMeta.SubpatternKey,
				"goal_query_hint":        goalMeta.QueryHint,
				"goal_query_terms":       goalMeta.QueryTerms,
				"dictionary_version":     recommendationDictionaryVersion,
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
		},
		"candidates": candidateSummaries,
	})
}

func (h *AdminHandler) PreviewRecommendationDebugDirect(c *gin.Context) {
	var req recommendationDebugDirectPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.SourceQuery) == "" || strings.TrimSpace(req.GroupTitle) == "" || strings.TrimSpace(req.LessonTitle) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_query, group_title, lesson_title are required"})
		return
	}
	if req.MaxPerLesson <= 0 || req.MaxPerLesson > 10 {
		req.MaxPerLesson = 5
	}

	ownerUserID := uuid.Nil
	if strings.TrimSpace(req.OwnerUserID) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(req.OwnerUserID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner_user_id"})
			return
		}
		ownerUserID = parsed
	}

	queryBundle := buildRecommendationDebugQueryFromInput(
		req.SourceQuery,
		req.GroupTitle,
		req.GroupObjective,
		req.LessonTitle,
		req.LessonObjective,
		req.LessonSummary,
	)
	embeddingInfo, embeddingText := h.buildRecommendationDebugEmbedding(c.Request.Context(), queryBundle.DenseQuery)
	repo := curriculum.NewRepository(h.db)
	lexicalQuery := queryBundle.LexicalQueryExpanded
	if strings.TrimSpace(lexicalQuery) == "" {
		lexicalQuery = queryBundle.LexicalQueryKO
	}
	diagnostics, err := repo.SearchLessonCandidatesDiagnostics(c.Request.Context(), ownerUserID, lexicalQuery, embeddingText, req.MaxPerLesson, curriculum.LessonCandidateSearchOptions{
		PreferredLanguage: "ko",
		Now:               time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search lesson candidates"})
		return
	}
	externalStage := h.buildRecommendationDebugExternalStage(c.Request.Context(), ownerUserID, req.MaxPerLesson, queryBundle, diagnostics.MergedCandidates)
	candidateSummaries := buildRecommendationDebugCandidateSummaries(
		diagnostics,
		nil,
		"ko",
		time.Now(),
	)

	c.JSON(http.StatusOK, gin.H{
		"input": gin.H{
			"owner_user_id":    ownerUserID,
			"source_query":     strings.TrimSpace(req.SourceQuery),
			"group_title":      strings.TrimSpace(req.GroupTitle),
			"group_objective":  strings.TrimSpace(req.GroupObjective),
			"lesson_title":     strings.TrimSpace(req.LessonTitle),
			"lesson_objective": strings.TrimSpace(req.LessonObjective),
			"lesson_summary":   strings.TrimSpace(req.LessonSummary),
		},
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
				"goal_pattern_key":       "",
				"goal_subpattern_key":    "",
				"goal_query_hint":        "",
				"goal_query_terms":       []string{},
				"dictionary_version":     recommendationDictionaryVersion,
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
		},
		"candidates": candidateSummaries,
	})
}
