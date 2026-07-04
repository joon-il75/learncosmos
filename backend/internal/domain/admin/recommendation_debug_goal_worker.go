package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

type recommendationDebugGoalTurnProcessor struct {
	h *AdminHandler
}

type recommendationDebugGoalTurnJobRef struct {
	ScenarioID string `json:"scenario_id"`
	AdminID    string `json:"admin_id,omitempty"`
	EventType  string `json:"event_type"`
}

type recommendationDebugGoalTurnResult struct {
	ScenarioID         string `json:"scenario_id"`
	Status             string `json:"status"`
	GoalProfileVersion int    `json:"goal_profile_version"`
	MessageCount       int    `json:"message_count"`
	LatencyMS          int    `json:"latency_ms,omitempty"`
	PersistMS          int    `json:"persist_ms,omitempty"`
}

func NewRecommendationDebugGoalTurnProcessor(handler *AdminHandler) llmjobs.Processor {
	return &recommendationDebugGoalTurnProcessor{h: handler}
}

func CreateRecommendationDebugGoalTurnJobInput(adminID, scenarioID uuid.UUID, eventType string) llmjobs.CreateJobInput {
	providerMode := "super_admin_debug"
	requestRef := map[string]any{
		"scenario_id": scenarioID.String(),
		"event_type":  strings.TrimSpace(eventType),
	}
	if adminID != uuid.Nil {
		requestRef["admin_id"] = adminID.String()
	}
	return llmjobs.CreateJobInput{
		UserID:         recommendationDebugJobUserID(adminID),
		Feature:        llmjobs.FeatureAdminRecommendationDebugGoalTurn,
		RequestRef:     requestRef,
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
	}
}

func (p *recommendationDebugGoalTurnProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_goal_turn_failed", "invalid recommendation debug goal job", false)
	}
	var ref recommendationDebugGoalTurnJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_goal_turn_failed", "invalid recommendation debug goal payload", false)
	}
	scenarioID, err := uuid.Parse(strings.TrimSpace(ref.ScenarioID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("admin_recommendation_debug_goal_turn_failed", "invalid scenario reference", false)
	}
	adminID := uuid.Nil
	if strings.TrimSpace(ref.AdminID) != "" {
		if parsed, parseErr := uuid.Parse(strings.TrimSpace(ref.AdminID)); parseErr == nil {
			adminID = parsed
		}
	} else if job.UserID != nil {
		adminID = *job.UserID
	}

	result, err := p.h.executeRecommendationDebugGoalTurn(ctx, scenarioID, adminID, strings.TrimSpace(ref.EventType))
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapRecommendationDebugGoalTurnProcessorError(err)
	}
	return llmjobs.CompleteJobInput{
		ResultRef: map[string]any{
			"scenario_id":          result.ScenarioID,
			"status":               result.Status,
			"goal_profile_version": result.GoalProfileVersion,
			"message_count":        result.MessageCount,
		},
		LLMGenerationElapsedMS: &result.LatencyMS,
		PersistElapsedMS:       &result.PersistMS,
	}, nil
}

func (h *AdminHandler) executeRecommendationDebugGoalTurn(ctx context.Context, scenarioID uuid.UUID, adminID uuid.UUID, eventType string) (*recommendationDebugGoalTurnResult, error) {
	scenario, err := h.loadRecommendationDebugScenarioByID(ctx, scenarioID)
	if err != nil {
		return nil, err
	}
	profile, err := recommendationDebugGoalProfileFromSnapshot(scenario.GoalProfileSnapshot, adminID, scenario.InitialUserIntent)
	if err != nil {
		return nil, errRecommendationDebugGoalInvalidState
	}
	if profile.InterviewState == goal.StateConfirmed {
		return nil, errRecommendationDebugGoalInvalidState
	}
	userMessage := lastRecommendationDebugUserMessage(profile)
	if strings.TrimSpace(userMessage) == "" {
		return nil, errRecommendationDebugGoalInvalidState
	}
	if eventType == "" {
		eventType = "goal_chat_message"
	}

	startedAt := time.Now()
	h.applyRecommendationDebugGoalResponse(ctx, adminID, profile, userMessage)
	latencyMS := int(time.Since(startedAt).Milliseconds())

	persistStarted := time.Now()
	if err := h.saveRecommendationDebugGoalSnapshotContext(ctx, scenarioID, adminID, profile, eventType, userMessage); err != nil {
		return nil, err
	}
	persistMS := int(time.Since(persistStarted).Milliseconds())

	status := "draft"
	if profile.InterviewState == goal.StateConfirmed {
		status = "goal_confirmed"
	}
	return &recommendationDebugGoalTurnResult{
		ScenarioID:         scenario.ID,
		Status:             status,
		GoalProfileVersion: profileVersion(profile),
		MessageCount:       len(profile.Messages),
		LatencyMS:          latencyMS,
		PersistMS:          persistMS,
	}, nil
}

func lastRecommendationDebugUserMessage(profile *goal.GoalProfile) string {
	if profile == nil {
		return ""
	}
	for i := len(profile.Messages) - 1; i >= 0; i-- {
		if profile.Messages[i].Role == "user" {
			return strings.TrimSpace(profile.Messages[i].Content)
		}
	}
	return ""
}

var errRecommendationDebugGoalInvalidState = errors.New("goal turn state is invalid")

func mapRecommendationDebugGoalTurnProcessorError(err error) error {
	switch {
	case errors.Is(err, errRecommendationDebugGoalInvalidState):
		return llmjobs.NewProcessorError("admin_goal_turn_invalid_state", err.Error(), false)
	default:
		return llmjobs.NewProcessorError("admin_recommendation_debug_goal_turn_failed", err.Error(), false)
	}
}
