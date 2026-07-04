package goal

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

type goalInterviewTurnProcessor struct {
	h *Handler
}

type goalInterviewTurnJobRef struct {
	GoalProfileID      string `json:"goal_profile_id"`
	GoalProfileVersion int    `json:"goal_profile_version"`
	MessageCount       int    `json:"message_count"`
	Mode               string `json:"mode"`
	Language           string `json:"language"`
}

type goalInterviewTurnJobResult struct {
	Response          *InterviewAIResponse `json:"response"`
	Route             string               `json:"route,omitempty"`
	LLMCallCount      int                  `json:"llm_call_count,omitempty"`
	FallbackReason    string               `json:"fallback_reason,omitempty"`
	SummarizedContext string               `json:"summarized_context,omitempty"`
	Provider          string               `json:"provider,omitempty"`
	Model             string               `json:"model,omitempty"`
	BillingStatus     string               `json:"billing_status,omitempty"`
}

func NewGoalInterviewTurnProcessor(handler *Handler) llmjobs.Processor {
	return &goalInterviewTurnProcessor{h: handler}
}

func CreateGoalInterviewTurnJobInput(userID uuid.UUID, profile *GoalProfile, mode string) llmjobs.CreateJobInput {
	feature := llmjobs.FeatureGoalInterviewTurn
	providerMode := "pending"
	ref := map[string]any{
		"mode":     normalizeGoalChatMode(mode),
		"language": normalizeLearningLanguage(""),
	}
	if profile != nil {
		ref["goal_profile_id"] = profile.ID.String()
		ref["goal_profile_version"] = profile.Version
		ref["message_count"] = len(profile.Messages)
		ref["language"] = normalizeLearningLanguage(profile.Language)
	}
	return llmjobs.CreateJobInput{
		UserID:         &userID,
		Feature:        feature,
		RequestRef:     ref,
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *goalInterviewTurnProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || p.h.repo == nil || p.h.service == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_failed", "invalid goal interview job", false)
	}
	var ref goalInterviewTurnJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_failed", "invalid goal interview job payload", false)
	}
	goalID, err := uuid.Parse(strings.TrimSpace(ref.GoalProfileID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_failed", "invalid goal profile reference", false)
	}
	profile, err := p.h.repo.GetGoalByID(ctx, goalID)
	if err != nil {
		code := "goal_interview_failed"
		if errors.Is(err, ErrNotFound) {
			code = "goal_not_found"
		}
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError(code, err.Error(), false)
	}
	if profile.UserID != *job.UserID || !profile.IsActive {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_not_found", "goal profile is not available", false)
	}
	if ref.GoalProfileVersion > 0 && profile.Version != ref.GoalProfileVersion {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_profile_version_changed", "goal profile version changed", false)
	}
	if ref.MessageCount > 0 && len(profile.Messages) != ref.MessageCount {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_profile_message_changed", "goal profile message changed", true)
	}
	userMessage := latestUserMessage(profile)
	if userMessage == "" {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_failed", "latest user message is empty", false)
	}
	profile.Language = normalizeLearningLanguage(profile.Language)
	exec, llmErr := p.h.resolveGoalLLMExecution(ctx, *job.UserID)
	var resp *InterviewAIResponse
	provider := ""
	model := ""
	billingStatus := ""
	if llmErr == nil && exec != nil && exec.Client != nil {
		provider = exec.Provider
		model = exec.Model
		billingStatus = exec.BillingStatus
		resp = p.h.resolveInterviewResponseWithExecution(ctx, *job.UserID, profile, userMessage, exec)
	} else {
		resp = p.h.service.BuildFallbackInterviewResponse(profile, userMessage)
		p.h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "llm_unavailable")
	}
	if resp == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_failed", "goal interview response is empty", false)
	}
	p.h.observeGoalInterviewAIOutput(ctx, *job.UserID, profile, resp, "goal_interview_worker")
	persistStarted := time.Now()
	p.h.applyInterviewResponse(profile, resp)
	if err := p.h.repo.UpdateGoalProfile(ctx, profile); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_interview_persist_failed", err.Error(), true)
	}
	persistElapsed := int(time.Since(persistStarted).Milliseconds())
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"response":             resp,
			"route":                strings.TrimSpace(resp.Route),
			"llm_call_count":       resp.LLMCallCount,
			"fallback_reason":      strings.TrimSpace(resp.FallbackReason),
			"summarized_context":   strings.TrimSpace(profile.SummarizedContext),
			"provider":             strings.TrimSpace(provider),
			"model":                strings.TrimSpace(model),
			"billing_status":       strings.TrimSpace(billingStatus),
			"goal_profile_id":      profile.ID.String(),
			"goal_profile_version": profile.Version,
		},
		PersistElapsedMS: &persistElapsed,
	}, nil
}

func latestUserMessage(profile *GoalProfile) string {
	if profile == nil {
		return ""
	}
	for idx := len(profile.Messages) - 1; idx >= 0; idx-- {
		if profile.Messages[idx].Role == "user" {
			return strings.TrimSpace(profile.Messages[idx].Content)
		}
	}
	return ""
}

func normalizeGoalChatMode(mode string) string {
	mode = strings.TrimSpace(mode)
	if mode == "analysis_template" {
		return mode
	}
	return "two_step"
}
