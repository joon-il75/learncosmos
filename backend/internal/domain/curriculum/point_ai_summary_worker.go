package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/pkg/llm"
	"github.com/learnweaver/backend/internal/pkg/metaparser"
)

type pointAISummaryProcessor struct {
	h *Handler
}

type pointAISummaryJobRef struct {
	PlanetID string `json:"planet_id"`
	PointID  string `json:"point_id"`
	Language string `json:"learning_language"`
}

type pointAISummaryJobResult struct {
	CourseID         string `json:"course_id"`
	PlanetID         string `json:"planet_id"`
	PointID          string `json:"point_id"`
	AISummaryID      string `json:"ai_summary_id,omitempty"`
	Provider         string `json:"provider,omitempty"`
	Model            string `json:"model,omitempty"`
	BillingStatus    string `json:"billing_status,omitempty"`
	LearningLanguage string `json:"learning_language,omitempty"`
}

func NewPointAISummaryProcessor(handler *Handler) llmjobs.Processor {
	return &pointAISummaryProcessor{h: handler}
}

func CreatePointAISummaryJobInput(userID, planetID, pointID uuid.UUID, learningLanguage string) llmjobs.CreateJobInput {
	providerMode := "pending"
	return llmjobs.CreateJobInput{
		UserID:  &userID,
		Feature: llmjobs.FeaturePointAISummary,
		RequestRef: map[string]any{
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"learning_language": normalizeLearningLanguage(learningLanguage),
		},
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *pointAISummaryProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_ai_summary_failed", "invalid point ai summary job", false)
	}
	var ref pointAISummaryJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_ai_summary_failed", "invalid point ai summary payload", false)
	}
	planetID, err := uuid.Parse(strings.TrimSpace(ref.PlanetID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_ai_summary_failed", "invalid planet reference", false)
	}
	pointID, err := uuid.Parse(strings.TrimSpace(ref.PointID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("point_ai_summary_failed", "invalid point reference", false)
	}
	result, err := p.h.executePointAISummary(ctx, *job.UserID, planetID, pointID, "point_ai_summary_worker")
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapPointAISummaryProcessorError(err)
	}
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"course_id":         result.CourseID.String(),
			"planet_id":         planetID.String(),
			"point_id":          pointID.String(),
			"provider":          result.Provider,
			"model":             result.Model,
			"billing_status":    result.BillingStatus,
			"learning_language": result.LearningLanguage,
		},
		LLMGenerationElapsedMS: &result.LLMElapsedMS,
		PersistElapsedMS:       &result.PersistElapsedMS,
	}, nil
}

type pointAISummaryExecutionResult struct {
	CourseID         uuid.UUID
	Provider         string
	Model            string
	BillingStatus    string
	LearningLanguage string
	LLMElapsedMS     int
	PersistElapsedMS int
}

func (h *Handler) executePointAISummary(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceFeature string) (*pointAISummaryExecutionResult, error) {
	tx, err := h.repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, courseID, point, err := h.repo.getLearningExplorationPointSourceTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return nil, err
	}
	if err := tx.Rollback(ctx); err != nil {
		log.Printf("[point-summary] rollback after source load failed point=%s err=%v", pointID, err)
	}

	sourceTitle := strings.TrimSpace(point.Title)
	sourceDescription := strings.TrimSpace(derefStr(point.Description))
	if point.ExternalURL != nil && strings.TrimSpace(*point.ExternalURL) != "" {
		metaCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
		meta, metaErr := metaparser.Parse(metaCtx, strings.TrimSpace(*point.ExternalURL), h.youtubeAPIKey)
		cancel()
		if metaErr == nil && meta != nil {
			if strings.TrimSpace(meta.Title) != "" {
				sourceTitle = strings.TrimSpace(meta.Title)
			}
			if strings.TrimSpace(meta.Description) != "" {
				sourceDescription = strings.TrimSpace(meta.Description)
			}
		}
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
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_ai_summary", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
	})

	llmCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	llmStarted := time.Now()
	raw, llmErr := client.Complete(llmCtx, buildPointAISummaryPrompt(point.Title, derefStr(point.Description), sourceTitle, sourceDescription, learningLanguage))
	llmElapsed := int(time.Since(llmStarted).Milliseconds())
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		return nil, errPointAISummaryFailed
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	summaryText := normalizePointAISummaryText(raw)
	persistStarted := time.Now()
	courseID, _, err = h.repo.UpsertLearningPointAISummary(ctx, userID, planetID, pointID, sourceTitle, sourceDescription, summaryText, learningLanguage)
	persistElapsed := int(time.Since(persistStarted).Milliseconds())
	if err != nil {
		return nil, err
	}
	h.observePointSummaryAIOutput(ctx, userID, planetID, pointID, summaryText, learningLanguage, sourceFeature)
	return &pointAISummaryExecutionResult{
		CourseID:         courseID,
		Provider:         strings.TrimSpace(exec.Provider),
		BillingStatus:    strings.TrimSpace(exec.BillingStatus),
		LearningLanguage: learningLanguage,
		LLMElapsedMS:     llmElapsed,
		PersistElapsedMS: persistElapsed,
	}, nil
}

var (
	errNoPointAILLMAvailable = errors.New("no point ai llm available")
	errPointAISummaryFailed  = errors.New("point ai summary failed")
)

func mapPointAISummaryProcessorError(err error) error {
	switch {
	case errors.Is(err, errNoPointAILLMAvailable):
		return llmjobs.NewProcessorError("no_llm_available", err.Error(), true)
	case errors.Is(err, errPointAISummaryFailed):
		return llmjobs.NewProcessorError("ai_summary_failed", err.Error(), true)
	case errors.Is(err, errLearningPointAISummaryInvalidInput):
		return llmjobs.NewProcessorError("point_ai_summary_empty", err.Error(), false)
	case errors.Is(err, errLearningPointNotFound):
		return llmjobs.NewProcessorError("learning_point_not_found", err.Error(), false)
	default:
		return llmjobs.NewProcessorError("point_ai_summary_failed", err.Error(), false)
	}
}
