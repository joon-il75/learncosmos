package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

func buildRecommendationDebugLessonsSnapshot(scenario recommendationDebugScenarioDetailRow, profile *goal.GoalProfile, draft *curriculum.DraftAggregate, provider, model string, latencyMS int) gin.H {
	lessons := make([]gin.H, 0)
	completionCriteria := []string{}
	title := ""
	description := ""
	if draft != nil {
		title = draft.Draft.Title
		description = derefString(draft.Draft.Description)
		completionCriteria = append(completionCriteria, draft.Draft.CompletionCriteria...)
		for _, lesson := range draft.Lessons {
			lessons = append(lessons, gin.H{
				"lesson_id":   lesson.Lesson.ID,
				"title":       lesson.Lesson.Title,
				"objective":   derefString(lesson.Lesson.Objective),
				"order_index": lesson.Lesson.OrderIndex,
				"source_type": lesson.Lesson.SourceType,
				"lesson_role": lesson.Lesson.LessonRole,
				"sub_lessons": len(lesson.SubLessons),
				"point_count": len(lesson.Points),
			})
		}
	}
	confirmedGoal := ""
	if profile != nil && profile.ConfirmedGoal != nil {
		confirmedGoal = strings.TrimSpace(*profile.ConfirmedGoal)
	}
	return gin.H{
		"course_title":        scenario.CourseTitle,
		"initial_user_intent": scenario.InitialUserIntent,
		"confirmed_goal":      confirmedGoal,
		"goal_profile": gin.H{
			"user_intent":      derefGoalString(profile, "user_intent"),
			"motivation":       derefGoalString(profile, "motivation"),
			"usage_context":    derefGoalString(profile, "usage_context"),
			"difficulty_level": derefGoalString(profile, "difficulty_level"),
			"time_horizon":     derefGoalString(profile, "time_horizon"),
			"goal_type":        derefGoalString(profile, "goal_type"),
			"output_type":      derefGoalString(profile, "output_type"),
			"version":          profileVersion(profile),
		},
		"draft_title":         title,
		"draft_description":   description,
		"completion_criteria": completionCriteria,
		"lessons":             lessons,
		"lesson_count":        len(lessons),
		"provider":            provider,
		"model":               model,
		"latency_ms":          latencyMS,
		"generated_at":        time.Now(),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func derefGoalString(profile *goal.GoalProfile, field string) string {
	if profile == nil {
		return ""
	}
	switch field {
	case "user_intent":
		return strings.TrimSpace(profile.UserIntent)
	case "motivation":
		return derefString(profile.Motivation)
	case "usage_context":
		return derefString(profile.UsageContext)
	case "difficulty_level":
		return derefString(profile.DifficultyLevel)
	case "time_horizon":
		return derefString(profile.TimeHorizon)
	case "goal_type":
		return derefString(profile.GoalType)
	case "output_type":
		return derefString(profile.OutputType)
	default:
		return ""
	}
}

func profileVersion(profile *goal.GoalProfile) int {
	if profile == nil || profile.Version <= 0 {
		return 1
	}
	return profile.Version
}

func (h *AdminHandler) loadRecommendationDebugScenarioForGoal(c *gin.Context) (recommendationDebugScenarioDetailRow, uuid.UUID, bool) {
	var empty recommendationDebugScenarioDetailRow
	scenarioID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return empty, uuid.Nil, false
	}
	adminIDRaw, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return empty, uuid.Nil, false
	}
	adminID, err := uuid.Parse(adminIDRaw)
	if err != nil {
		adminID = uuid.Nil
	}

	item, err := h.loadRecommendationDebugScenarioByID(c.Request.Context(), scenarioID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "scenario not found"})
			return empty, uuid.Nil, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load scenario"})
		return empty, uuid.Nil, false
	}

	return item, adminID, true
}

func (h *AdminHandler) loadRecommendationDebugScenarioByID(ctx context.Context, scenarioID uuid.UUID) (recommendationDebugScenarioDetailRow, error) {
	var item recommendationDebugScenarioDetailRow
	err := h.db.QueryRow(ctx, `
		SELECT
			id::text,
			COALESCE(created_by::text, ''),
			course_title,
			initial_user_intent,
			status,
			COALESCE(goal_profile_snapshot::text, ''),
			COALESCE(generated_lessons_snapshot::text, ''),
			notes,
			created_at,
			updated_at
		FROM recommendation_debug_scenarios
		WHERE id = $1
		  AND status <> 'archived'
	`, scenarioID).Scan(
		&item.ID,
		&item.CreatedBy,
		&item.CourseTitle,
		&item.InitialUserIntent,
		&item.Status,
		&item.GoalProfileSnapshot,
		&item.GeneratedLessonsSnapshot,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

func (h *AdminHandler) applyRecommendationDebugGoalResponse(ctx context.Context, adminID uuid.UUID, profile *goal.GoalProfile, userMessage string) {
	service := goal.NewService()
	resp := h.resolveRecommendationDebugGoalResponse(ctx, adminID, profile, userMessage)
	if profile == nil || resp == nil {
		return
	}
	profile.InterviewState = resp.NextState
	service.ApplyExtracted(profile, resp.Extracted)
	proposedGoal := strings.TrimSpace(derefString(resp.ProposedGoal))
	if extracted := service.ExtractProposedGoalFromReply(resp.Message); resp.NextState == goal.StateProposingGoal && extracted != "" {
		proposedGoal = extracted
	}
	if proposedGoal != "" {
		profile.ConfirmedGoal = &proposedGoal
	}
	profile.Messages = append(profile.Messages, goal.InterviewMessage{
		Role:    "lumi",
		Content: resp.Message,
	})
}

func (h *AdminHandler) resolveRecommendationDebugGoalResponse(ctx context.Context, adminID uuid.UUID, profile *goal.GoalProfile, userMessage string) *goal.InterviewAIResponse {
	service := goal.NewService()
	repo := goal.NewRepository(h.db)
	provider, model, settingErr := repo.GetSystemLLMSetting(ctx, "default")
	if strings.TrimSpace(model) == "" {
		model = "gpt-4o-mini"
	}
	apiKey, _ := repo.GetSystemAPIKey(ctx, provider)
	client, clientErr := llm.NewClient(strings.TrimSpace(strings.ToLower(provider)), apiKey, model)
	if settingErr == nil && clientErr == nil && client != nil {
		llmCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if resp, err := service.RunTwoStepLLMWithClient(llmCtx, profile, userMessage, client); err == nil {
			return resp
		}
		prompt := service.BuildInterviewPrompt(profile, userMessage)
		if raw, err := service.CallLLMWithClient(llmCtx, client, prompt); err == nil {
			if resp, parseErr := service.ParseAIResponse(raw); parseErr == nil && !service.ShouldOverrideAIResponse(profile, userMessage, resp) {
				return resp
			}
		}
	}
	return service.BuildFallbackInterviewResponse(profile, userMessage)
}

func (h *AdminHandler) saveRecommendationDebugGoalSnapshot(c *gin.Context, scenarioID string, adminID uuid.UUID, profile *goal.GoalProfile, eventType string, userMessage string) bool {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return false
	}
	if err := h.saveRecommendationDebugGoalSnapshotContext(c.Request.Context(), scenarioUUID, adminID, profile, eventType, userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save goal snapshot"})
		return false
	}
	return true
}

func (h *AdminHandler) saveRecommendationDebugGoalSnapshotContext(ctx context.Context, scenarioID uuid.UUID, adminID uuid.UUID, profile *goal.GoalProfile, eventType string, userMessage string) error {
	snapshot, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	requestSnapshot, _ := json.Marshal(gin.H{
		"event_type":   eventType,
		"user_message": userMessage,
	})

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	status := "draft"
	if profile.InterviewState == goal.StateConfirmed {
		status = "goal_confirmed"
	}
	if _, err := tx.Exec(ctx, `
		UPDATE recommendation_debug_scenarios
		SET goal_profile_snapshot = $2::jsonb,
		    status = $3,
		    updated_at = NOW()
		WHERE id = $1
	`, scenarioID, string(snapshot), status); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_debug_runs (
			scenario_id,
			created_by,
			run_type,
			request_snapshot,
			baseline_result_snapshot,
			provider_snapshot
		)
		VALUES ($1, $2, 'goal_chat', $3::jsonb, $4::jsonb, $5::jsonb)
	`, scenarioID, recommendationDebugNullableAdminID(adminID), string(requestSnapshot), string(snapshot), `{"mode":"super_admin_debug","source":"goal_service"}`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (h *AdminHandler) saveRecommendationDebugGoalProfileOnly(ctx context.Context, scenarioID uuid.UUID, profile *goal.GoalProfile) error {
	snapshot, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	status := "draft"
	if profile.InterviewState == goal.StateConfirmed {
		status = "goal_confirmed"
	}
	_, err = h.db.Exec(ctx, `
		UPDATE recommendation_debug_scenarios
		SET goal_profile_snapshot = $2::jsonb,
		    status = $3,
		    updated_at = NOW()
		WHERE id = $1
	`, scenarioID, string(snapshot), status)
	return err
}

func recommendationDebugNullableAdminID(adminID uuid.UUID) any {
	if adminID == uuid.Nil {
		return nil
	}
	return adminID
}

func recommendationDebugGoalProfileFromSnapshot(raw string, adminID uuid.UUID, userIntent string) (*goal.GoalProfile, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "{}" {
		return nil, fmt.Errorf("empty goal snapshot")
	}
	var profile goal.GoalProfile
	if err := json.Unmarshal([]byte(trimmed), &profile); err != nil {
		return nil, err
	}
	if profile.UserID == uuid.Nil {
		profile.UserID = adminID
	}
	if strings.TrimSpace(profile.UserIntent) == "" {
		profile.UserIntent = userIntent
	}
	if len(profile.Messages) == 0 {
		profile.Messages = []goal.InterviewMessage{{Role: "user", Content: userIntent}}
	}
	if profile.Version <= 0 {
		profile.Version = 1
	}
	return &profile, nil
}
