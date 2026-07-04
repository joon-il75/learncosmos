package goal

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/safety"
)

func (h *Handler) enforceGoalSafety(c *gin.Context, userID uuid.UUID, targetType string, targetID *uuid.UUID, text string, locale string, metadata map[string]any) bool {
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

func (h *Handler) observeGoalInterviewAIOutput(ctx context.Context, userID uuid.UUID, profile *GoalProfile, resp *InterviewAIResponse, sourceFeature string) {
	if h == nil || h.safetyService == nil || resp == nil {
		return
	}
	text := goalInterviewAIOutputText(resp)
	if text == "" {
		return
	}
	metadata := map[string]any{
		"source_feature": sourceFeature,
		"output_kind":    "goal_interview_response",
		"next_state":     string(resp.NextState),
		"response_route": strings.TrimSpace(resp.Route),
		"llm_call_count": resp.LLMCallCount,
	}
	if reason := strings.TrimSpace(resp.FallbackReason); reason != "" {
		metadata["fallback_reason"] = reason
	}
	var targetID *uuid.UUID
	if profile != nil {
		if profile.ID != uuid.Nil {
			id := profile.ID
			targetID = &id
			metadata["goal_profile_id"] = id.String()
		}
		if profile.CourseDraftID != nil {
			metadata["course_draft_id"] = profile.CourseDraftID.String()
		}
	}
	h.safetyService.ObserveAIOutput(ctx, safety.ModerateInput{
		UserID:     &userID,
		TargetType: safety.TargetAIGoalInterviewOutput,
		TargetID:   targetID,
		Text:       text,
		Locale:     profileLanguage(profile),
		Route:      "goal_interview",
		Metadata:   metadata,
	})
}

func goalInterviewAIOutputText(resp *InterviewAIResponse) string {
	if resp == nil {
		return ""
	}
	parts := []string{strings.TrimSpace(resp.Message)}
	if value := strings.TrimSpace(derefString(resp.ProposedGoal)); value != "" {
		parts = append(parts, value)
	}
	if value := strings.TrimSpace(derefString(resp.GoalCandidate)); value != "" {
		parts = append(parts, value)
	}
	for _, item := range resp.MissingInfo {
		if value := strings.TrimSpace(item); value != "" {
			parts = append(parts, value)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func profileLanguage(profile *GoalProfile) string {
	if profile == nil {
		return ""
	}
	return profile.Language
}
