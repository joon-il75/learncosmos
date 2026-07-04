package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

type recommendationDebugLessonsGenerateProcessor struct {
	h *AdminHandler
}

type recommendationDebugLessonsGenerateJobRef struct {
	ScenarioID string `json:"scenario_id"`
	AdminID    string `json:"admin_id,omitempty"`
}

type recommendationDebugLessonsGenerateResult struct {
	Snapshot  map[string]any `json:"snapshot"`
	Scenario  map[string]any `json:"scenario"`
	Provider  string         `json:"provider,omitempty"`
	Model     string         `json:"model,omitempty"`
	LatencyMS int            `json:"latency_ms,omitempty"`
	PersistMS int            `json:"persist_ms,omitempty"`
}

func NewRecommendationDebugLessonsGenerateProcessor(handler *AdminHandler) llmjobs.Processor {
	return &recommendationDebugLessonsGenerateProcessor{h: handler}
}

func CreateRecommendationDebugLessonsGenerateJobInput(adminID, scenarioID uuid.UUID) llmjobs.CreateJobInput {
	providerMode := "super_admin_debug"
	requestRef := map[string]any{
		"scenario_id": scenarioID.String(),
	}
	if adminID != uuid.Nil {
		requestRef["admin_id"] = adminID.String()
	}
	return llmjobs.CreateJobInput{
		UserID:         recommendationDebugJobUserID(adminID),
		Feature:        llmjobs.FeatureAdminRecommendationDebugLessonsGenerate,
		RequestRef:     requestRef,
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *recommendationDebugLessonsGenerateProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_lessons_failed", "invalid recommendation debug lessons job", false)
	}
	var ref recommendationDebugLessonsGenerateJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_lessons_failed", "invalid recommendation debug lessons payload", false)
	}
	scenarioID, err := uuid.Parse(strings.TrimSpace(ref.ScenarioID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_lessons_failed", "invalid scenario reference", false)
	}
	adminID := uuid.Nil
	if strings.TrimSpace(ref.AdminID) != "" {
		if parsed, parseErr := uuid.Parse(strings.TrimSpace(ref.AdminID)); parseErr == nil {
			adminID = parsed
		}
	} else if job.UserID != nil {
		adminID = *job.UserID
	}

	started := time.Now()
	result, err := p.h.executeRecommendationDebugLessons(ctx, scenarioID, adminID)
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapRecommendationDebugLessonsProcessorError(err)
	}
	totalMS := int(time.Since(started).Milliseconds())
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"snapshot": result.Snapshot,
			"scenario": result.Scenario,
		},
		LLMGenerationElapsedMS: &result.LatencyMS,
		PersistElapsedMS:       &result.PersistMS,
		TotalJobElapsedMS:      &totalMS,
	}, nil
}

func recommendationDebugJobUserID(adminID uuid.UUID) *uuid.UUID {
	if adminID == uuid.Nil {
		return nil
	}
	return &adminID
}

var (
	errRecommendationDebugLessonsInvalidState = errors.New("confirmed goal is required")
	errRecommendationDebugLessonsUnavailable  = errors.New("system llm unavailable")
	errRecommendationDebugLessonsFailed       = errors.New("lesson generation failed")
)

func mapRecommendationDebugLessonsProcessorError(err error) error {
	switch {
	case errors.Is(err, errRecommendationDebugLessonsInvalidState):
		return llmjobs.NewProcessorError("confirmed_goal_required", err.Error(), false)
	case errors.Is(err, errRecommendationDebugLessonsUnavailable):
		return llmjobs.NewProcessorError("system_llm_unavailable", err.Error(), true)
	case errors.Is(err, errRecommendationDebugLessonsFailed):
		return llmjobs.NewProcessorError("lesson_generation_failed", err.Error(), true)
	default:
		return llmjobs.NewProcessorError("admin_recommendation_debug_lessons_failed", err.Error(), false)
	}
}
