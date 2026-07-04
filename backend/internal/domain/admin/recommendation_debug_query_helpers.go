package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func findRecommendationDebugLesson(draft *curriculum.DraftAggregate, lessonID uuid.UUID) (curriculum.DraftLessonTree, curriculum.DraftLessonTree, bool) {
	if draft == nil {
		return curriculum.DraftLessonTree{}, curriculum.DraftLessonTree{}, false
	}

	for _, mainTree := range draft.Lessons {
		for _, subTree := range mainTree.SubLessons {
			if subTree.Lesson.ID == lessonID {
				return mainTree, subTree, true
			}
		}
	}

	return curriculum.DraftLessonTree{}, curriculum.DraftLessonTree{}, false
}

func (h *AdminHandler) buildRecommendationDebugLessonQuery(ctx context.Context, draft curriculum.CourseDraft, mainLesson curriculum.CourseDraftLesson, subLesson curriculum.CourseDraftLesson) (normalizer.LessonSearchQuery, curriculum.RecommendationGoalMetadata) {
	mainObjective := ""
	if mainLesson.Objective != nil {
		mainObjective = *mainLesson.Objective
	}
	lessonObjective := ""
	if subLesson.Objective != nil {
		lessonObjective = *subLesson.Objective
	}
	lessonSummary := ""
	if subLesson.Summary != nil {
		lessonSummary = *subLesson.Summary
	}
	recommendationReq := curriculum.CreateCourseDraftRequest{
		SourceQuery:     draft.SourceQuery,
		LearningGoal:    draft.LearningGoal,
		GoalProfileID:   draft.GoalProfileID,
		CurrentLevel:    draft.CurrentLevel,
		PreferredFormat: draft.PreferredFormat,
	}
	goalRepo := goal.NewRepository(h.db)
	activeGoal, err := goalRepo.GetActiveGoal(ctx, draft.ID)
	if err != nil {
		activeGoal, err = goalRepo.GetActiveGoalByUser(ctx, draft.UserID)
	}
	if err == nil && activeGoal != nil {
		recommendationReq.GoalMotivation = activeGoal.Motivation
		recommendationReq.GoalUsageContext = activeGoal.UsageContext
		recommendationReq.GoalDifficultyLevel = activeGoal.DifficultyLevel
		recommendationReq.GoalTimeHorizon = activeGoal.TimeHorizon
		recommendationReq.GoalOutputType = activeGoal.OutputType
		recommendationReq.GoalType = activeGoal.GoalType
		if activeGoal.ConfirmedGoal != nil && strings.TrimSpace(derefString(recommendationReq.LearningGoal)) == "" {
			recommendationReq.LearningGoal = activeGoal.ConfirmedGoal
		}
	}
	goalMeta := curriculum.BuildRecommendationGoalMetadata(recommendationReq)
	return buildRecommendationDebugQueryFromInput(
		strings.Join(compactRecommendationParts(
			draft.SourceQuery,
			derefString(recommendationReq.LearningGoal),
			goalMeta.QueryHint,
		), " "),
		mainLesson.Title,
		mainObjective,
		subLesson.Title,
		lessonObjective,
		lessonSummary,
	), goalMeta
}

func buildRecommendationDebugQueryFromInput(sourceQuery, levelTitle, levelObjective, lessonTitle, lessonObjective, lessonSummary string) normalizer.LessonSearchQuery {
	return normalizer.BuildLessonSearchQuery(
		sourceQuery,
		levelTitle,
		levelObjective,
		lessonTitle,
		lessonObjective,
		lessonSummary,
	)
}

func compactRecommendationParts(parts ...string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func recommendationDebugVectorToPGVector(vector []float32) string {
	parts := make([]string, len(vector))
	for i, v := range vector {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func recommendationDebugEmbeddingPreview(vector []float32, limit int) []float32 {
	if limit <= 0 || len(vector) <= limit {
		return vector
	}
	return vector[:limit]
}

func recommendationDebugFailureReason(query normalizer.LessonSearchQuery, embeddingInfo gin.H, lexicalHits, vectorHits, mergedHits int) string {
	if len(query.Tokens) == 0 {
		return "NO_QUERY_TOKENS"
	}
	if mergedHits > 0 {
		return ""
	}
	if lexicalHits == 0 && vectorHits == 0 {
		if used, _ := embeddingInfo["used"].(bool); !used {
			return "EMBEDDING_FAILED"
		}
		return "LEXICAL_ZERO_VECTOR_ZERO"
	}
	if lexicalHits == 0 && vectorHits > 0 {
		return "LEXICAL_ZERO_VECTOR_LOW"
	}
	return "SEED_CONTENT_SHORTAGE"
}
