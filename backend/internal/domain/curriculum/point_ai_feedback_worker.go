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

type pointFeedbackGenerateProcessor struct {
	h *Handler
}

type pointFeedbackGenerateJobRef struct {
	PlanetID   string `json:"planet_id"`
	PointID    string `json:"point_id"`
	QuestionID string `json:"question_id"`
	Language   string `json:"learning_language"`
}

type pointFeedbackGenerateJobResult struct {
	PlanetID         string `json:"planet_id"`
	PointID          string `json:"point_id"`
	QuestionID       string `json:"question_id"`
	Feedback         string `json:"feedback"`
	Provider         string `json:"provider,omitempty"`
	Model            string `json:"model,omitempty"`
	BillingStatus    string `json:"billing_status,omitempty"`
	LearningLanguage string `json:"learning_language,omitempty"`
}

func NewPointFeedbackGenerateProcessor(handler *Handler) llmjobs.Processor {
	return &pointFeedbackGenerateProcessor{h: handler}
}

func CreatePointFeedbackGenerateJobInput(userID, planetID, pointID, questionID uuid.UUID, learningLanguage string) llmjobs.CreateJobInput {
	providerMode := "pending"
	return llmjobs.CreateJobInput{
		UserID:  &userID,
		Feature: llmjobs.FeaturePointFeedbackGenerate,
		RequestRef: map[string]any{
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"question_id":       questionID.String(),
			"learning_language": normalizeLearningLanguage(learningLanguage),
		},
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *pointFeedbackGenerateProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_feedback_generate_failed", "invalid point feedback job", false)
	}
	var ref pointFeedbackGenerateJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_feedback_generate_failed", "invalid point feedback payload", false)
	}
	planetID, err := uuid.Parse(strings.TrimSpace(ref.PlanetID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_feedback_generate_failed", "invalid planet reference", false)
	}
	pointID, err := uuid.Parse(strings.TrimSpace(ref.PointID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_feedback_generate_failed", "invalid point reference", false)
	}
	questionID, err := uuid.Parse(strings.TrimSpace(ref.QuestionID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_feedback_generate_failed", "invalid question reference", false)
	}
	result, err := p.h.executePointFeedbackGenerate(ctx, *job.UserID, planetID, pointID, questionID)
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapPointFeedbackGenerateProcessorError(err)
	}
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"question_id":       questionID.String(),
			"feedback":          result.Feedback,
			"provider":          result.Provider,
			"model":             result.Model,
			"billing_status":    result.BillingStatus,
			"learning_language": result.LearningLanguage,
		},
		LLMGenerationElapsedMS: &result.LLMElapsedMS,
	}, nil
}

type pointFeedbackGenerateExecutionResult struct {
	Feedback         string
	Provider         string
	Model            string
	BillingStatus    string
	LearningLanguage string
	LLMElapsedMS     int
}

func (h *Handler) executePointFeedbackGenerate(ctx context.Context, userID, planetID, pointID, questionID uuid.UUID) (*pointFeedbackGenerateExecutionResult, error) {
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
	var targetQuestion *CoursePointQuestion
	for idx := range resolvedPoint.Questions {
		if resolvedPoint.Questions[idx].ID == questionID {
			targetQuestion = &resolvedPoint.Questions[idx]
			break
		}
	}
	if targetQuestion == nil {
		return nil, errLearningPointQuestionNotFound
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
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_question_feedback", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
		"question_id":       questionID.String(),
	})

	llmCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	llmStarted := time.Now()
	raw, llmErr := client.Complete(llmCtx, buildPointAIFeedbackPrompt(planet.GoalContext, resolvedPoint, lessonTitle, *targetQuestion, learningLanguage))
	llmElapsed := int(time.Since(llmStarted).Milliseconds())
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		return nil, errPointAIFeedbackGenerateFailed
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	feedback := normalizePointAIFeedback(raw, learningLanguage)
	h.observePointFeedbackAIOutput(ctx, userID, planetID, pointID, questionID, feedback, learningLanguage, "point_feedback_generate_worker")
	return &pointFeedbackGenerateExecutionResult{
		Feedback:         feedback,
		Provider:         strings.TrimSpace(exec.Provider),
		BillingStatus:    strings.TrimSpace(exec.BillingStatus),
		LearningLanguage: learningLanguage,
		LLMElapsedMS:     llmElapsed,
	}, nil
}

var (
	errPointAIFeedbackGenerateFailed = errors.New("point ai feedback generate failed")
	errLearningPointQuestionNotFound = errors.New("learning point question not found")
)

func mapPointFeedbackGenerateProcessorError(err error) error {
	switch {
	case errors.Is(err, errNoPointAILLMAvailable):
		return llmjobs.NewProcessorError("no_llm_available", err.Error(), true)
	case errors.Is(err, errPointAIFeedbackGenerateFailed):
		return llmjobs.NewProcessorError("ai_feedback_failed", err.Error(), true)
	case errors.Is(err, errLearningPointNotFound):
		return llmjobs.NewProcessorError("learning_point_not_found", err.Error(), false)
	case errors.Is(err, errLearningPointQuestionNotFound):
		return llmjobs.NewProcessorError("learning_point_question_not_found", err.Error(), false)
	default:
		return llmjobs.NewProcessorError("point_feedback_generate_failed", err.Error(), false)
	}
}

func (h *Handler) respondWithPointFeedbackGenerateJobResult(c *gin.Context, job *llmjobs.Job) {
	if job == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_feedback_generate_result_invalid", "error_code": "point_feedback_generate_result_invalid"})
		return
	}
	var result pointFeedbackGenerateJobResult
	if err := json.Unmarshal(job.ResultRef, &result); err != nil || strings.TrimSpace(result.Feedback) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_feedback_generate_result_invalid", "error_code": "point_feedback_generate_result_invalid"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"feedback": result.Feedback, "learning_language": result.LearningLanguage})
}
