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

type pointQuestionGenerateProcessor struct {
	h *Handler
}

type pointQuestionGenerateJobRef struct {
	PlanetID string `json:"planet_id"`
	PointID  string `json:"point_id"`
	Language string `json:"learning_language"`
}

type pointQuestionGenerateJobResult struct {
	CourseID         string `json:"course_id"`
	PlanetID         string `json:"planet_id"`
	PointID          string `json:"point_id"`
	QuestionID       string `json:"question_id"`
	Provider         string `json:"provider,omitempty"`
	Model            string `json:"model,omitempty"`
	BillingStatus    string `json:"billing_status,omitempty"`
	LearningLanguage string `json:"learning_language,omitempty"`
}

func NewPointQuestionGenerateProcessor(handler *Handler) llmjobs.Processor {
	return &pointQuestionGenerateProcessor{h: handler}
}

func CreatePointQuestionGenerateJobInput(userID, planetID, pointID uuid.UUID, learningLanguage string) llmjobs.CreateJobInput {
	providerMode := "pending"
	return llmjobs.CreateJobInput{
		UserID:  &userID,
		Feature: llmjobs.FeaturePointQuestionGenerate,
		RequestRef: map[string]any{
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"learning_language": normalizeLearningLanguage(learningLanguage),
		},
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *pointQuestionGenerateProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_question_generate_failed", "invalid point question job", false)
	}
	var ref pointQuestionGenerateJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_question_generate_failed", "invalid point question payload", false)
	}
	planetID, err := uuid.Parse(strings.TrimSpace(ref.PlanetID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_question_generate_failed", "invalid planet reference", false)
	}
	pointID, err := uuid.Parse(strings.TrimSpace(ref.PointID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_question_generate_failed", "invalid point reference", false)
	}
	result, err := p.h.executePointQuestionGenerate(ctx, *job.UserID, planetID, pointID)
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapPointQuestionGenerateProcessorError(err)
	}
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"course_id":         result.CourseID.String(),
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"question_id":       result.QuestionID.String(),
			"provider":          result.Provider,
			"model":             result.Model,
			"billing_status":    result.BillingStatus,
			"learning_language": result.LearningLanguage,
		},
		LLMGenerationElapsedMS: &result.LLMElapsedMS,
		PersistElapsedMS:       &result.PersistElapsedMS,
	}, nil
}

type pointQuestionGenerateExecutionResult struct {
	CourseID         uuid.UUID
	QuestionID       uuid.UUID
	Provider         string
	Model            string
	BillingStatus    string
	LearningLanguage string
	LLMElapsedMS     int
	PersistElapsedMS int
}

func (h *Handler) executePointQuestionGenerate(ctx context.Context, userID, planetID, pointID uuid.UUID) (*pointQuestionGenerateExecutionResult, error) {
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
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_question_generate", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
	})

	llmCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	llmStarted := time.Now()
	raw, llmErr := client.Complete(llmCtx, buildPointAIQuestionPrompt(planet.GoalContext, resolvedPoint, lessonTitle, learningLanguage))
	llmElapsed := int(time.Since(llmStarted).Milliseconds())
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		return nil, errPointAIQuestionGenerateFailed
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	suggestion := parsePointAIQuestionSuggestion(raw, resolvedPoint.Point.PointType, learningLanguage)
	persistStarted := time.Now()
	courseID, question, err := h.repo.createLearningPointQuestionWithAuthor(ctx, userID, planetID, pointID, nil, suggestion.Question, PointQuestionType(suggestion.QuestionType), nil, "ai")
	persistElapsed := int(time.Since(persistStarted).Milliseconds())
	if err != nil {
		return nil, err
	}
	h.observePointQuestionAIOutput(ctx, userID, planetID, pointID, question.ID, question.Question, question.QuestionType, learningLanguage, "point_question_generate_worker")
	return &pointQuestionGenerateExecutionResult{
		CourseID:         courseID,
		QuestionID:       question.ID,
		Provider:         strings.TrimSpace(exec.Provider),
		BillingStatus:    strings.TrimSpace(exec.BillingStatus),
		LearningLanguage: learningLanguage,
		LLMElapsedMS:     llmElapsed,
		PersistElapsedMS: persistElapsed,
	}, nil
}

var errPointAIQuestionGenerateFailed = errors.New("point ai question generate failed")

func mapPointQuestionGenerateProcessorError(err error) error {
	switch {
	case errors.Is(err, errNoPointAILLMAvailable):
		return llmjobs.NewProcessorError("no_llm_available", err.Error(), true)
	case errors.Is(err, errPointAIQuestionGenerateFailed):
		return llmjobs.NewProcessorError("ai_question_failed", err.Error(), true)
	case errors.Is(err, errLearningPointQuestionInvalidInput):
		return llmjobs.NewProcessorError("point_question_invalid", err.Error(), false)
	case errors.Is(err, errLearningPointNotFound):
		return llmjobs.NewProcessorError("learning_point_not_found", err.Error(), false)
	default:
		return llmjobs.NewProcessorError("point_question_generate_failed", err.Error(), false)
	}
}

func (h *Handler) respondWithPointQuestionGenerateJobResult(c *gin.Context, userID, fallbackPointID uuid.UUID, job *llmjobs.Job) {
	if job == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_question_generate_result_invalid", "error_code": "point_question_generate_result_invalid"})
		return
	}
	var result pointQuestionGenerateJobResult
	if err := json.Unmarshal(job.ResultRef, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_question_generate_result_invalid", "error_code": "point_question_generate_result_invalid"})
		return
	}
	courseID, err := uuid.Parse(strings.TrimSpace(result.CourseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_question_generate_result_invalid", "error_code": "point_question_generate_result_invalid"})
		return
	}
	pointID := fallbackPointID
	if parsedPointID, parseErr := uuid.Parse(strings.TrimSpace(result.PointID)); parseErr == nil {
		pointID = parsedPointID
	}
	questionID, err := uuid.Parse(strings.TrimSpace(result.QuestionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_question_generate_result_invalid", "error_code": "point_question_generate_result_invalid"})
		return
	}

	point, err := h.repo.GetPlanetPointDetailByIDAndDraftStatuses(c.Request.Context(), userID, courseID, pointID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ai question saved but failed to reload", "error_code": "point_question_generate_reload_failed"})
		return
	}
	response := gin.H{
		"planet": point.Planet,
		"point":  point,
	}
	for idx := range point.Point.Questions {
		if point.Point.Questions[idx].ID == questionID {
			response["question"] = point.Point.Questions[idx]
			break
		}
	}
	if _, ok := response["question"]; !ok {
		response["question_id"] = questionID.String()
	}
	c.JSON(http.StatusCreated, response)
}
