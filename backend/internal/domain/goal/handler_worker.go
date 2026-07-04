package goal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

type goalInterviewWorkerPendingError struct {
	job *llmjobs.Job
}

func (e *goalInterviewWorkerPendingError) Error() string {
	return "goal interview worker job is pending"
}

func (h *Handler) applyGoalInterviewTurn(ctx context.Context, userID uuid.UUID, profile *GoalProfile, userMessage string) error {
	if h == nil || profile == nil {
		return nil
	}
	if h.goalChatFastPathEnabled {
		if resp, ok := h.service.TryBuildGoalInterviewFastPath(profile, userMessage); ok {
			h.annotateGoalInterviewResponse(resp, "fast_path", 0, "")
			h.observeGoalInterviewAIOutput(ctx, userID, profile, resp, "goal_interview_fast_path")
			h.applyInterviewResponse(profile, resp)
			return nil
		}
	}
	if h.shouldUseGoalInterviewWorker() {
		return h.applyGoalInterviewTurnWithWorker(ctx, userID, profile)
	}
	llmCtx, llmCancel := context.WithTimeout(ctx, 30*time.Second)
	defer llmCancel()
	resp := h.resolveInterviewResponse(llmCtx, userID, profile, userMessage)
	h.observeGoalInterviewAIOutput(ctx, userID, profile, resp, "goal_interview_sync")
	h.applyInterviewResponse(profile, resp)
	return nil
}

func (h *Handler) shouldUseGoalInterviewWorker() bool {
	return h != nil && h.goalChatWorkerEnabled && h.llmGateway != nil
}

func (h *Handler) applyGoalInterviewTurnWithWorker(ctx context.Context, userID uuid.UUID, profile *GoalProfile) error {
	if h == nil || profile == nil {
		return nil
	}
	startedAt := time.Now()
	if err := h.repo.UpdateGoalProfile(ctx, profile); err != nil {
		return err
	}
	input := CreateGoalInterviewTurnJobInput(userID, profile, h.goalChatMode)
	job, err := h.llmGateway.SubmitAndWait(ctx, input, h.goalChatWorkerSyncWait)
	if err != nil {
		resp := h.service.BuildFallbackInterviewResponse(profile, latestUserMessage(profile))
		h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "worker_submit_failed")
		h.observeGoalInterviewAIOutput(ctx, userID, profile, resp, "goal_interview_worker_submit_fallback")
		h.applyInterviewResponse(profile, resp)
		h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, profile.Language)
		return nil
	}
	if job == nil || job.Status != llmjobs.StatusSucceeded {
		if job != nil && !terminalGoalInterviewJobStatus(job.Status) {
			return &goalInterviewWorkerPendingError{job: job}
		}
		reason := "worker_failed"
		if job != nil {
			reason = strings.TrimSpace(derefString(job.ErrorCode))
		}
		if reason == "" {
			reason = "worker_failed"
		}
		resp := h.service.BuildFallbackInterviewResponse(profile, latestUserMessage(profile))
		h.annotateGoalInterviewResponse(resp, "local_fallback", 0, reason)
		h.observeGoalInterviewAIOutput(ctx, userID, profile, resp, "goal_interview_worker_fallback")
		h.applyInterviewResponse(profile, resp)
		h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, profile.Language)
		return nil
	}
	result, err := parseGoalInterviewTurnJobResult(job)
	if err != nil || result.Response == nil {
		resp := h.service.BuildFallbackInterviewResponse(profile, latestUserMessage(profile))
		h.annotateGoalInterviewResponse(resp, "local_fallback", 0, "worker_result_parse_failed")
		h.observeGoalInterviewAIOutput(ctx, userID, profile, resp, "goal_interview_worker_parse_fallback")
		h.applyInterviewResponse(profile, resp)
		h.logGoalInterviewTelemetry(startedAt, userID, resp, nil, profile.Language)
		return nil
	}
	result.Response.Route = strings.TrimSpace(result.Route)
	result.Response.LLMCallCount = result.LLMCallCount
	result.Response.FallbackReason = strings.TrimSpace(result.FallbackReason)
	exec := &goalLLMExecution{Provider: result.Provider, Model: result.Model, BillingStatus: result.BillingStatus}
	h.logGoalInterviewTelemetry(startedAt, userID, result.Response, exec, profile.Language)
	updated, err := h.repo.GetGoalByID(ctx, profile.ID)
	if err != nil {
		return err
	}
	*profile = *updated
	return nil
}

func terminalGoalInterviewJobStatus(status llmjobs.Status) bool {
	return status == llmjobs.StatusFailed || status == llmjobs.StatusCanceled || status == llmjobs.StatusExpired
}

func parseGoalInterviewTurnJobResult(job *llmjobs.Job) (*goalInterviewTurnJobResult, error) {
	if job == nil {
		return nil, errors.New("nil goal interview job")
	}
	var result goalInterviewTurnJobResult
	if err := json.Unmarshal(job.ResultRef, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func writeGoalInterviewPendingResponse(c *gin.Context, profile *GoalProfile, job *llmjobs.Job) {
	if job == nil {
		c.JSON(http.StatusAccepted, gin.H{
			"error":      "goal_interview_pending",
			"error_code": "goal_interview_pending",
			"goal":       profile,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"error":      "goal_interview_pending",
		"error_code": "goal_interview_pending",
		"goal":       profile,
		"job_id":     job.ID,
		"feature":    job.Feature,
		"status":     job.Status,
		"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
	})
}

func handleGoalInterviewTurnError(c *gin.Context, profile *GoalProfile, err error, fallbackMessage string) bool {
	if err == nil {
		return false
	}
	var pending *goalInterviewWorkerPendingError
	if errors.As(err, &pending) {
		writeGoalInterviewPendingResponse(c, profile, pending.job)
		return true
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMessage})
	return true
}
