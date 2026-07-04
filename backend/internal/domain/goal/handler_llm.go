package goal

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

type goalLLMExecution struct {
	Client        llm.Client
	Provider      string
	Model         string
	BillingStatus string
}

type goalInterviewCountingClient struct {
	base  llm.Client
	calls *int
}

func (c *goalInterviewCountingClient) Complete(ctx context.Context, prompt string) (string, error) {
	if c != nil && c.calls != nil {
		(*c.calls)++
	}
	if c == nil || c.base == nil {
		return "", errors.New("llm client is nil")
	}
	return c.base.Complete(ctx, prompt)
}

func (h *Handler) resolveGoalLLMExecution(ctx context.Context, userID uuid.UUID) (*goalLLMExecution, error) {
	provider, model, llmErr := h.repo.GetSystemLLMSetting(ctx, "default")
	if strings.TrimSpace(model) == "" {
		model = "gpt-4o-mini"
	}

	userConfig, cfgErr := h.repo.GetUserRuntimeAIConfig(ctx, userID)
	if cfgErr == nil && userConfig != nil && userConfig.Mode == "byok" && strings.EqualFold(userConfig.Provider, "openai") && strings.TrimSpace(userConfig.APIKey) != "" {
		byokProvider := strings.TrimSpace(strings.ToLower(userConfig.Provider))
		client, err := llm.NewClient(byokProvider, strings.TrimSpace(userConfig.APIKey), model)
		if err != nil {
			return nil, err
		}
		return &goalLLMExecution{
			Client:        h.trackBYOKLLMUsage(client, userID, byokProvider, model, "goal_interview", "byok_no_charge"),
			Provider:      byokProvider,
			Model:         model,
			BillingStatus: "byok_no_charge",
		}, nil
	}

	if llmErr != nil {
		return nil, llmErr
	}
	apiKey, _ := h.repo.GetSystemAPIKey(ctx, provider)
	client, err := llm.NewClient(strings.TrimSpace(strings.ToLower(provider)), apiKey, model)
	if err != nil {
		return nil, err
	}
	return &goalLLMExecution{
		Client:        client,
		Provider:      provider,
		Model:         model,
		BillingStatus: "charged",
	}, nil
}

func (h *Handler) resolveInterviewResponse(ctx context.Context, userID uuid.UUID, profile *GoalProfile, userMessage string) *InterviewAIResponse {
	startedAt := time.Now()
	if profile == nil {
		resp := h.service.BuildFallbackInterviewResponse(&GoalProfile{}, userMessage)
		h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "nil_profile")
		h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, "")
		return resp
	}
	profile.Language = normalizeLearningLanguage(profile.Language)
	if h.goalChatFastPathEnabled {
		if resp, ok := h.service.TryBuildGoalInterviewFastPath(profile, userMessage); ok {
			h.annotateGoalInterviewResponse(resp, "fast_path", 0, "")
			h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, profile.Language)
			return resp
		}
	}

	exec, llmErr := h.resolveGoalLLMExecution(ctx, userID)
	if llmErr == nil && exec != nil && exec.Client != nil {
		resp := h.resolveInterviewResponseWithExecution(ctx, userID, profile, userMessage, exec)
		h.logGoalInterviewTelemetry(startedAt, userID, resp, exec, profile.Language)
		return resp
	}

	resp := h.service.BuildFallbackInterviewResponse(profile, userMessage)
	h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "llm_unavailable")
	h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, profile.Language)
	return resp
}

func (h *Handler) resolveInterviewResponseWithExecution(ctx context.Context, userID uuid.UUID, profile *GoalProfile, userMessage string, exec *goalLLMExecution) *InterviewAIResponse {
	if profile == nil {
		resp := h.service.BuildFallbackInterviewResponse(&GoalProfile{}, userMessage)
		h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "nil_profile")
		return resp
	}
	if h.goalChatFastPathEnabled {
		if resp, ok := h.service.TryBuildGoalInterviewFastPath(profile, userMessage); ok {
			h.annotateGoalInterviewResponse(resp, "fast_path", 0, "")
			return resp
		}
	}
	if exec != nil && exec.Client != nil {
		llmCallCount := 0
		countingClient := &goalInterviewCountingClient{base: exec.Client, calls: &llmCallCount}
		// 2-step LLM 우선 시도
		var resp *InterviewAIResponse
		var err error
		if h.goalChatMode == "analysis_template" {
			resp, err = h.service.RunAnalysisTemplateReplyWithClient(ctx, profile, userMessage, countingClient)
		} else {
			resp, err = h.service.RunTwoStepLLMWithClient(ctx, profile, userMessage, countingClient)
		}
		if err == nil {
			if h.repo != nil && exec.BillingStatus == "byok_no_charge" {
				go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
			}
			h.annotateGoalInterviewResponse(resp, h.goalChatRouteName(), llmCallCount, "")
			return resp
		} else if h.repo != nil && exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, err.Error())
		}
		if !h.singlePromptFallbackLLMEnabled {
			resp := h.service.BuildFallbackInterviewResponse(profile, userMessage)
			h.annotateGoalInterviewResponse(resp, "local_fallback", llmCallCount, goalInterviewFallbackReason(err))
			return resp
		}
		// 2-step 실패 시 단일 프롬프트 fallback
		prompt := h.service.BuildInterviewPrompt(profile, userMessage)
		if raw, callErr := h.service.CallLLMWithClient(ctx, countingClient, prompt); callErr == nil {
			if aiResp, parseErr := h.service.ParseAIResponse(raw); parseErr == nil && !h.service.ShouldOverrideAIResponse(profile, userMessage, aiResp) {
				if h.repo != nil && exec.BillingStatus == "byok_no_charge" {
					go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
				}
				h.annotateGoalInterviewResponse(aiResp, "single_prompt_fallback", llmCallCount, goalInterviewFallbackReason(err))
				return aiResp
			}
		} else if h.repo != nil && exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, callErr.Error())
		}
		resp = h.service.BuildFallbackInterviewResponse(profile, userMessage)
		h.annotateGoalInterviewResponse(resp, "local_fallback", llmCallCount, goalInterviewFallbackReason(err))
		return resp
	}

	resp := h.service.BuildFallbackInterviewResponse(profile, userMessage)
	h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "llm_unavailable")
	return resp
}
