package goal

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) goalChatRouteName() string {
	if h.goalChatMode == "analysis_template" {
		return "analysis_template"
	}
	return "two_step"
}

func (h *Handler) annotateGoalInterviewResponse(resp *InterviewAIResponse, route string, llmCallCount int, fallbackReason string) {
	if resp == nil {
		return
	}
	resp.Route = route
	resp.LLMCallCount = llmCallCount
	resp.FallbackReason = fallbackReason
}

func goalInterviewFallbackReason(err error) string {
	if err == nil {
		return ""
	}
	value := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(value, "step a parse"):
		return "step_a_parse_error"
	case strings.Contains(value, "step a llm"):
		return "step_a_error"
	case strings.Contains(value, "step b parse"):
		return "step_b_parse_error"
	case strings.Contains(value, "step b llm"):
		return "step_b_error"
	default:
		return "llm_error"
	}
}

func (h *Handler) logGoalInterviewTelemetry(startedAt time.Time, userID uuid.UUID, resp *InterviewAIResponse, exec *goalLLMExecution, language string) {
	if resp == nil {
		return
	}
	provider := ""
	billingStatus := ""
	if exec != nil {
		provider = exec.Provider
		billingStatus = exec.BillingStatus
	}
	log.Printf(
		"[goal/interview] user=%s route=%s mode=%s llm_call_count=%d provider=%s billing_status=%s language=%s fallback_reason=%s duration_ms=%d",
		userID,
		resp.Route,
		h.goalChatMode,
		resp.LLMCallCount,
		strings.TrimSpace(provider),
		strings.TrimSpace(billingStatus),
		normalizeLearningLanguage(language),
		strings.TrimSpace(resp.FallbackReason),
		time.Since(startedAt).Milliseconds(),
	)
}
