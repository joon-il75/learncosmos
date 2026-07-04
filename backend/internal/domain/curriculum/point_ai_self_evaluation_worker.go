package curriculum

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
	"github.com/learnweaver/backend/internal/pkg/llm"
)

type pointSelfEvaluationDraftProcessor struct {
	h *Handler
}

type pointSelfEvaluationDraftJobRef struct {
	PlanetID         string `json:"planet_id"`
	PointID          string `json:"point_id"`
	Language         string `json:"learning_language"`
	UseStoredAnswers bool   `json:"use_stored_answers"`
}

type pointSelfEvaluationDraftJobResult struct {
	PlanetID         string                             `json:"planet_id"`
	PointID          string                             `json:"point_id"`
	Draft            LearningPointSelfEvaluationAIDraft `json:"draft"`
	Provider         string                             `json:"provider,omitempty"`
	Model            string                             `json:"model,omitempty"`
	BillingStatus    string                             `json:"billing_status,omitempty"`
	LearningLanguage string                             `json:"learning_language,omitempty"`
}

func NewPointSelfEvaluationDraftProcessor(handler *Handler) llmjobs.Processor {
	return &pointSelfEvaluationDraftProcessor{h: handler}
}

func CreatePointSelfEvaluationDraftJobInput(userID, planetID, pointID uuid.UUID, learningLanguage string, useStoredAnswers bool) llmjobs.CreateJobInput {
	providerMode := "pending"
	return llmjobs.CreateJobInput{
		UserID:  &userID,
		Feature: llmjobs.FeaturePointSelfEvaluationDraft,
		RequestRef: map[string]any{
			"planet_id":          planetID.String(),
			"point_id":           pointID.String(),
			"learning_language":  normalizeLearningLanguage(learningLanguage),
			"use_stored_answers": useStoredAnswers,
		},
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *pointSelfEvaluationDraftProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_self_evaluation_draft_failed", "invalid point self evaluation draft job", false)
	}
	var ref pointSelfEvaluationDraftJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_self_evaluation_draft_failed", "invalid point self evaluation draft payload", false)
	}
	planetID, err := uuid.Parse(strings.TrimSpace(ref.PlanetID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_self_evaluation_draft_failed", "invalid planet reference", false)
	}
	pointID, err := uuid.Parse(strings.TrimSpace(ref.PointID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_self_evaluation_draft_failed", "invalid point reference", false)
	}
	var answers []LearningPointSelfEvaluationApplicationAnswer
	if ref.UseStoredAnswers {
		answers, err = p.h.repo.GetLearningPointSelfEvaluationApplicationAnswers(ctx, *job.UserID, planetID, pointID)
		if err != nil {
			return llmjobs.CompleteJobInput{}, mapPointSelfEvaluationDraftProcessorError(err)
		}
		if len(answers) == 0 {
			return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_self_evaluation_answers_missing", "self evaluation application answers missing", false)
		}
	}
	result, err := p.h.executePointSelfEvaluationDraft(ctx, *job.UserID, planetID, pointID, answers)
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapPointSelfEvaluationDraftProcessorError(err)
	}
	p.h.observePointSelfEvaluationAIOutput(ctx, *job.UserID, planetID, pointID, result.Draft, result.LearningLanguage, "point_self_evaluation_draft_worker")
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"draft":             result.Draft,
			"provider":          result.Provider,
			"model":             result.Model,
			"billing_status":    result.BillingStatus,
			"learning_language": result.LearningLanguage,
		},
		LLMGenerationElapsedMS: &result.LLMElapsedMS,
	}, nil
}

type pointSelfEvaluationDraftExecutionResult struct {
	Draft            LearningPointSelfEvaluationAIDraft
	Provider         string
	Model            string
	BillingStatus    string
	LearningLanguage string
	LLMElapsedMS     int
}

func (h *Handler) executePointSelfEvaluationDraft(ctx context.Context, userID, planetID, pointID uuid.UUID, applicationAnswers []LearningPointSelfEvaluationApplicationAnswer) (*pointSelfEvaluationDraftExecutionResult, error) {
	tx, err := h.repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	_, courseID, err := h.repo.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, err
	}
	_ = tx.Rollback(ctx)

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(ctx, userID, courseID, learningPointStatuses())
	if err != nil {
		return nil, err
	}
	resolvedPoint, lessonTitle, found := findPointInPlanetLessons(planet.Lessons, pointID)
	if !found {
		return nil, errLearningPointNotFound
	}

	exec, err := h.resolveEvaluationExecution(ctx, userID, 0)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(exec.APIKey) == "" {
		return nil, errNoPointAILLMAvailable
	}
	learningLanguage := h.repo.GetUserLearningLanguage(ctx, userID)
	client, clientErr := llm.NewClient(exec.Provider, exec.APIKey, "")
	if clientErr != nil {
		return nil, clientErr
	}
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_self_evaluation_draft", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
	})

	llmCtx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()
	llmStarted := time.Now()
	raw, llmErr := client.Complete(llmCtx, buildPointSelfEvaluationDraftPrompt(planet.GoalContext, resolvedPoint, lessonTitle, applicationAnswers, learningLanguage))
	llmElapsed := int(time.Since(llmStarted).Milliseconds())
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		return nil, errPointSelfEvaluationDraftFailed
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	return &pointSelfEvaluationDraftExecutionResult{
		Draft:            parsePointSelfEvaluationDraft(raw, resolvedPoint, learningLanguage),
		Provider:         strings.TrimSpace(exec.Provider),
		BillingStatus:    strings.TrimSpace(exec.BillingStatus),
		LearningLanguage: learningLanguage,
		LLMElapsedMS:     llmElapsed,
	}, nil
}

var errPointSelfEvaluationDraftFailed = errors.New("point self evaluation draft failed")

func mapPointSelfEvaluationDraftProcessorError(err error) error {
	switch {
	case errors.Is(err, errNoPointAILLMAvailable):
		return llmjobs.NewProcessorError("no_llm_available", err.Error(), true)
	case errors.Is(err, errPointSelfEvaluationDraftFailed):
		return llmjobs.NewProcessorError("ai_self_evaluation_draft_failed", err.Error(), true)
	case errors.Is(err, errLearningPointNotFound):
		return llmjobs.NewProcessorError("learning_point_not_found", err.Error(), false)
	default:
		return llmjobs.NewProcessorError("point_self_evaluation_draft_failed", err.Error(), false)
	}
}

func (h *Handler) respondWithPointSelfEvaluationDraftJobResult(c *gin.Context, job *llmjobs.Job) {
	if job == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_self_evaluation_draft_result_invalid", "error_code": "point_self_evaluation_draft_result_invalid"})
		return
	}
	var result pointSelfEvaluationDraftJobResult
	if err := json.Unmarshal(job.ResultRef, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_self_evaluation_draft_result_invalid", "error_code": "point_self_evaluation_draft_result_invalid"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"draft": result.Draft, "learning_language": result.LearningLanguage})
}
