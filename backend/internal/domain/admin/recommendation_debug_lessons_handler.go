package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

func (h *AdminHandler) GenerateRecommendationDebugLessons(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}
	scenarioUUID, err := uuid.Parse(scenario.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return
	}

	if h.recommendationDebugLessonsWorkerEnabled && h.llmGateway != nil {
		job, err := h.llmGateway.SubmitAndWait(
			c.Request.Context(),
			CreateRecommendationDebugLessonsGenerateJobInput(adminID, scenarioUUID),
			h.recommendationDebugLessonsWorkerSyncWait,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "lesson generation worker unavailable", "error_code": "admin_recommendation_debug_lessons_worker_unavailable"})
			return
		}
		h.respondWithRecommendationDebugLessonsJob(c, job)
		return
	}

	result, err := h.executeRecommendationDebugLessons(c.Request.Context(), scenarioUUID, adminID)
	if err != nil {
		h.respondWithRecommendationDebugLessonsError(c, scenario.ID, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"snapshot": result.Snapshot,
		"scenario": result.Scenario,
	})
}

func (h *AdminHandler) respondWithRecommendationDebugLessonsJob(c *gin.Context, job *llmjobs.Job) {
	if job == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "lesson generation worker unavailable", "error_code": "admin_recommendation_debug_lessons_worker_unavailable"})
		return
	}
	switch job.Status {
	case llmjobs.StatusSucceeded:
		var result recommendationDebugLessonsGenerateResult
		if err := json.Unmarshal(job.ResultRef, &result); err != nil || result.Snapshot == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lesson generation result unavailable", "error_code": "admin_recommendation_debug_lessons_result_unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"snapshot": result.Snapshot,
			"scenario": result.Scenario,
		})
	case llmjobs.StatusFailed, llmjobs.StatusCanceled, llmjobs.StatusExpired:
		code := "admin_recommendation_debug_lessons_failed"
		if job.ErrorCode != nil && strings.TrimSpace(*job.ErrorCode) != "" {
			code = strings.TrimSpace(*job.ErrorCode)
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": code, "error_code": code})
	default:
		c.JSON(http.StatusAccepted, gin.H{
			"error":      "admin_recommendation_debug_lessons_pending",
			"error_code": "admin_recommendation_debug_lessons_pending",
			"job_id":     job.ID.String(),
			"feature":    job.Feature,
			"status":     job.Status,
			"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
		})
	}
}

func (h *AdminHandler) executeRecommendationDebugLessons(ctx context.Context, scenarioID uuid.UUID, adminID uuid.UUID) (*recommendationDebugLessonsGenerateResult, error) {
	scenario, err := h.loadRecommendationDebugScenarioByID(ctx, scenarioID)
	if err != nil {
		return nil, err
	}
	profile, err := recommendationDebugGoalProfileFromSnapshot(scenario.GoalProfileSnapshot, adminID, scenario.InitialUserIntent)
	if err != nil || profile.InterviewState != goal.StateConfirmed || profile.ConfirmedGoal == nil || strings.TrimSpace(*profile.ConfirmedGoal) == "" {
		return nil, errRecommendationDebugLessonsInvalidState
	}

	confirmedGoal := strings.TrimSpace(*profile.ConfirmedGoal)
	req := curriculum.CreateCourseDraftRequest{
		SourceQuery:         strings.TrimSpace(scenario.CourseTitle),
		LearningGoal:        &confirmedGoal,
		GoalUserIntent:      &profile.UserIntent,
		GoalMotivation:      profile.Motivation,
		GoalUsageContext:    profile.UsageContext,
		GoalDifficultyLevel: profile.DifficultyLevel,
		GoalTimeHorizon:     profile.TimeHorizon,
		GoalOutputType:      profile.OutputType,
		GoalType:            profile.GoalType,
	}
	if strings.TrimSpace(req.SourceQuery) == "" {
		req.SourceQuery = strings.TrimSpace(scenario.InitialUserIntent)
	}

	draftUserID := adminID
	if draftUserID == uuid.Nil {
		draftUserID = scenarioID
	}

	provider, model, settingErr := goal.NewRepository(h.db).GetSystemLLMSetting(ctx, "default")
	if settingErr != nil {
		return nil, errRecommendationDebugLessonsUnavailable
	}
	apiKey, err := h.settingsStore.ResolveAPIKey(ctx, provider, os.Getenv("OPENAI_API_KEY"))
	if err != nil || strings.TrimSpace(apiKey) == "" {
		return nil, errRecommendationDebugLessonsUnavailable
	}
	client, err := llm.NewClient(strings.TrimSpace(strings.ToLower(provider)), apiKey, strings.TrimSpace(model))
	if err != nil {
		return nil, errRecommendationDebugLessonsUnavailable
	}

	startedAt := time.Now()
	service := curriculum.NewService()
	llmCtx, cancel := context.WithTimeout(ctx, 75*time.Second)
	defer cancel()
	draft, err := service.BuildDraftWithLLM(llmCtx, client, draftUserID, req, curriculum.DraftBuildOptions{
		SkipReview:        true,
		GenerationTimeout: 60 * time.Second,
		ReviewTimeout:     0,
	})
	if err != nil {
		log.Printf("[recommendation-debug] lesson generation failed scenario_id=%s provider=%s model=%s err=%v", scenario.ID, provider, model, err)
		return nil, errRecommendationDebugLessonsFailed
	}
	latencyMS := int(time.Since(startedAt).Milliseconds())

	snapshot := buildRecommendationDebugLessonsSnapshot(scenario, profile, draft, provider, model, latencyMS)
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}
	requestSnapshot, _ := json.Marshal(gin.H{
		"course_title":         scenario.CourseTitle,
		"initial_user_intent":  scenario.InitialUserIntent,
		"confirmed_goal":       confirmedGoal,
		"goal_profile_version": profile.Version,
	})
	providerSnapshot, _ := json.Marshal(gin.H{
		"provider": provider,
		"model":    model,
		"mode":     "super_admin_debug",
	})

	persistStarted := time.Now()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE recommendation_debug_scenarios
		SET generated_lessons_snapshot = $2::jsonb,
		    status = 'lessons_generated',
		    updated_at = NOW()
		WHERE id = $1
	`, scenarioID, string(snapshotJSON)); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_debug_runs (
			scenario_id,
			created_by,
			run_type,
			request_snapshot,
			baseline_result_snapshot,
			provider_snapshot,
			latency_ms
		)
		VALUES ($1, $2, 'lesson_generation', $3::jsonb, $4::jsonb, $5::jsonb, $6)
	`, scenarioID, recommendationDebugNullableAdminID(adminID), string(requestSnapshot), string(snapshotJSON), string(providerSnapshot), latencyMS); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	persistMS := int(time.Since(persistStarted).Milliseconds())

	return &recommendationDebugLessonsGenerateResult{
		Snapshot:  snapshot,
		Scenario:  map[string]any{"id": scenario.ID, "status": "lessons_generated", "generated_lessons_snapshot": string(snapshotJSON)},
		Provider:  provider,
		Model:     model,
		LatencyMS: latencyMS,
		PersistMS: persistMS,
	}, nil
}

func (h *AdminHandler) respondWithRecommendationDebugLessonsError(c *gin.Context, scenarioID string, err error) {
	switch {
	case errors.Is(err, errRecommendationDebugLessonsInvalidState):
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed goal is required"})
	case errors.Is(err, errRecommendationDebugLessonsUnavailable):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "system llm unavailable"})
	case errors.Is(err, errRecommendationDebugLessonsFailed):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "lesson generation failed"})
	default:
		log.Printf("[recommendation-debug] lesson generation failed scenario_id=%s err=%v", scenarioID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate lessons"})
	}
}
