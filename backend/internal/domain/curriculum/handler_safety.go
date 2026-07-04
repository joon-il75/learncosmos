package curriculum

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/safety"
)

func (h *Handler) enforceSafety(c *gin.Context, userID uuid.UUID, targetType string, targetID *uuid.UUID, text string, locale string, metadata map[string]any) bool {
	if h == nil || h.safetyService == nil {
		return false
	}
	result, err := h.safetyService.Enforce(c.Request.Context(), safety.ModerateInput{
		UserID:     &userID,
		TargetType: targetType,
		TargetID:   targetID,
		Text:       text,
		Locale:     locale,
		Route:      c.FullPath(),
		Metadata:   metadata,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to moderate input", "error_code": "safety_moderation_failed"})
		return true
	}
	if safety.IsBlocked(result) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "input blocked by safety policy",
			"error_code": safety.ErrorCodeBlocked,
			"safety":     result,
		})
		return true
	}
	return false
}

func (h *Handler) observePointQuestionAIOutput(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, pointID uuid.UUID, questionID uuid.UUID, questionText string, questionType PointQuestionType, locale string, sourceFeature string) {
	if h == nil || h.safetyService == nil {
		return
	}
	text := strings.TrimSpace(questionText)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "point_question",
		"planet_id":      planetID.String(),
		"point_id":       pointID.String(),
		"question_id":    questionID.String(),
		"question_type":  string(questionType),
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAIPointQuestionOutput,
		TargetID:   &questionID,
		Text:       text,
		Locale:     locale,
		Route:      "point_question_generate",
		Metadata:   metadata,
	})
}

func (h *Handler) observePointFeedbackAIOutput(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, pointID uuid.UUID, questionID uuid.UUID, feedbackText string, locale string, sourceFeature string) {
	if h == nil || h.safetyService == nil {
		return
	}
	text := strings.TrimSpace(feedbackText)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "point_feedback",
		"planet_id":      planetID.String(),
		"point_id":       pointID.String(),
		"question_id":    questionID.String(),
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAIPointFeedbackOutput,
		TargetID:   &questionID,
		Text:       text,
		Locale:     locale,
		Route:      "point_feedback_generate",
		Metadata:   metadata,
	})
}

func (h *Handler) observePointSelfEvaluationAIOutput(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, pointID uuid.UUID, draft LearningPointSelfEvaluationAIDraft, locale string, sourceFeature string) {
	if h == nil || h.safetyService == nil {
		return
	}
	text := pointSelfEvaluationDraftText(draft)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "point_self_evaluation_draft",
		"planet_id":      planetID.String(),
		"point_id":       pointID.String(),
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAIPointSelfEvalOutput,
		TargetID:   &pointID,
		Text:       text,
		Locale:     locale,
		Route:      "point_self_evaluation_draft",
		Metadata:   metadata,
	})
}

func pointSelfEvaluationDraftText(draft LearningPointSelfEvaluationAIDraft) string {
	parts := make([]string, 0, 12+len(draft.ApplicationQuestions)*2)
	for _, item := range draft.ApplicationQuestions {
		if value := strings.TrimSpace(item.Question); value != "" {
			parts = append(parts, value)
		}
		if value := strings.TrimSpace(item.Intent); value != "" {
			parts = append(parts, value)
		}
	}
	for _, value := range []string{
		draft.UnderstandingReason,
		draft.ApplicationReason,
		draft.ProficiencyReason,
		draft.ProblemSolvingReason,
		draft.ExpressionReason,
		draft.GoalAlignmentNote,
	} {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func (h *Handler) observePointSummaryAIOutput(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, pointID uuid.UUID, summaryText string, locale string, sourceFeature string) {
	if h == nil || h.safetyService == nil {
		return
	}
	text := strings.TrimSpace(summaryText)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "point_summary",
		"planet_id":      planetID.String(),
		"point_id":       pointID.String(),
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAIPointSummaryOutput,
		TargetID:   &pointID,
		Text:       text,
		Locale:     locale,
		Route:      "point_ai_summary",
		Metadata:   metadata,
	})
}

func (h *Handler) observeCourseDraftAIOutput(ctx context.Context, userID uuid.UUID, draft *DraftAggregate, courseID uuid.UUID, sourceFeature string, extraMetadata map[string]any) {
	if h == nil || h.safetyService == nil || draft == nil {
		return
	}
	text := courseDraftAIOutputText(draft)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "course_draft",
		"draft_id":       draft.Draft.ID.String(),
		"course_id":      courseID.String(),
		"lesson_count":   len(draft.Lessons),
	}
	for key, value := range extraMetadata {
		metadata[key] = value
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAICourseDraftOutput,
		TargetID:   &draft.Draft.ID,
		Text:       text,
		Locale:     normalizeLearningLanguage(draft.Draft.GenerationLanguage),
		Route:      "course_draft_create",
		Metadata:   metadata,
	})
}

func courseDraftAIOutputText(draft *DraftAggregate) string {
	if draft == nil {
		return ""
	}
	parts := []string{
		strings.TrimSpace(draft.Draft.Title),
		strings.TrimSpace(draft.Draft.SourceQuery),
		strings.TrimSpace(derefStr(draft.Draft.LearningGoal)),
	}
	for _, lesson := range draft.Lessons {
		appendCourseDraftLessonText(&parts, lesson)
	}
	return strings.TrimSpace(strings.Join(compactStrings(parts), "\n"))
}

func appendCourseDraftLessonText(parts *[]string, lesson DraftLessonTree) {
	*parts = append(*parts, strings.TrimSpace(lesson.Lesson.Title), strings.TrimSpace(derefStr(lesson.Lesson.Objective)))
	for _, point := range lesson.Points {
		*parts = append(*parts, strings.TrimSpace(point.Point.Title), strings.TrimSpace(derefStr(point.Point.Description)))
	}
	for _, subLesson := range lesson.SubLessons {
		appendCourseDraftLessonText(parts, subLesson)
	}
}

func compactStrings(values []string) []string {
	compact := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			compact = append(compact, trimmed)
		}
	}
	return compact
}
