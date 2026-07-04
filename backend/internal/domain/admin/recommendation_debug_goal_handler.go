package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

func (h *AdminHandler) StartRecommendationDebugGoalInterview(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugGoalMessageRequest
	_ = c.ShouldBindJSON(&req)
	req.Message = strings.TrimSpace(req.Message)
	userIntent := req.Message
	if userIntent == "" {
		userIntent = strings.TrimSpace(scenario.InitialUserIntent)
	}
	if userIntent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "initial goal request is required"})
		return
	}
	if len([]rune(userIntent)) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "initial goal request too long"})
		return
	}
	if userIntent != strings.TrimSpace(scenario.InitialUserIntent) {
		scenarioUUID, _ := uuid.Parse(scenario.ID)
		if _, err := h.db.Exec(c.Request.Context(), `
			UPDATE recommendation_debug_scenarios
			SET initial_user_intent = $2,
			    updated_at = NOW()
			WHERE id = $1
		`, scenarioUUID, userIntent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save initial goal request"})
			return
		}
		scenario.InitialUserIntent = userIntent
	}

	profile := &goal.GoalProfile{
		UserID:         adminID,
		UserIntent:     userIntent,
		InterviewState: goal.StateClarifying,
		Messages: []goal.InterviewMessage{
			{Role: "user", Content: userIntent},
		},
		Version:  1,
		IsActive: true,
	}
	if h.recommendationDebugGoalWorkerEnabled && h.llmGateway != nil {
		scenarioUUID, _ := uuid.Parse(scenario.ID)
		if err := h.saveRecommendationDebugGoalProfileOnly(c.Request.Context(), scenarioUUID, profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save goal snapshot"})
			return
		}
		job, err := h.llmGateway.SubmitAndWait(
			c.Request.Context(),
			CreateRecommendationDebugGoalTurnJobInput(adminID, scenarioUUID, "goal_chat_start"),
			h.recommendationDebugGoalWorkerSyncWait,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "goal chat worker unavailable", "error_code": "admin_recommendation_debug_goal_worker_unavailable"})
			return
		}
		h.respondWithRecommendationDebugGoalJob(c, job, scenarioUUID, adminID)
		return
	}

	h.applyRecommendationDebugGoalResponse(c.Request.Context(), adminID, profile, userIntent)
	if !h.saveRecommendationDebugGoalSnapshot(c, scenario.ID, adminID, profile, "goal_chat_start", userIntent) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

func (h *AdminHandler) SendRecommendationDebugGoalMessage(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugGoalMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}
	if len([]rune(req.Message)) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message too long"})
		return
	}

	profile, err := recommendationDebugGoalProfileFromSnapshot(scenario.GoalProfileSnapshot, adminID, scenario.InitialUserIntent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal interview not started"})
		return
	}
	if profile.InterviewState == goal.StateConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal already confirmed"})
		return
	}

	profile.Messages = append(profile.Messages, goal.InterviewMessage{Role: "user", Content: req.Message})
	if h.recommendationDebugGoalWorkerEnabled && h.llmGateway != nil {
		scenarioUUID, _ := uuid.Parse(scenario.ID)
		if err := h.saveRecommendationDebugGoalProfileOnly(c.Request.Context(), scenarioUUID, profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save goal snapshot"})
			return
		}
		job, err := h.llmGateway.SubmitAndWait(
			c.Request.Context(),
			CreateRecommendationDebugGoalTurnJobInput(adminID, scenarioUUID, "goal_chat_message"),
			h.recommendationDebugGoalWorkerSyncWait,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "goal chat worker unavailable", "error_code": "admin_recommendation_debug_goal_worker_unavailable"})
			return
		}
		h.respondWithRecommendationDebugGoalJob(c, job, scenarioUUID, adminID)
		return
	}

	h.applyRecommendationDebugGoalResponse(c.Request.Context(), adminID, profile, req.Message)
	if !h.saveRecommendationDebugGoalSnapshot(c, scenario.ID, adminID, profile, "goal_chat_message", req.Message) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

func (h *AdminHandler) respondWithRecommendationDebugGoalJob(c *gin.Context, job *llmjobs.Job, scenarioID uuid.UUID, adminID uuid.UUID) {
	if job == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "goal chat worker unavailable", "error_code": "admin_recommendation_debug_goal_worker_unavailable"})
		return
	}
	switch job.Status {
	case llmjobs.StatusSucceeded:
		scenario, err := h.loadRecommendationDebugScenarioByID(c.Request.Context(), scenarioID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal snapshot", "error_code": "admin_recommendation_debug_goal_result_unavailable"})
			return
		}
		profile, err := recommendationDebugGoalProfileFromSnapshot(scenario.GoalProfileSnapshot, adminID, scenario.InitialUserIntent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode goal snapshot", "error_code": "admin_recommendation_debug_goal_result_unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"goal": profile})
	case llmjobs.StatusFailed, llmjobs.StatusCanceled, llmjobs.StatusExpired:
		code := "admin_recommendation_debug_goal_turn_failed"
		if job.ErrorCode != nil && strings.TrimSpace(*job.ErrorCode) != "" {
			code = strings.TrimSpace(*job.ErrorCode)
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": code, "error_code": code})
	default:
		c.JSON(http.StatusAccepted, gin.H{
			"error":      "admin_recommendation_debug_goal_pending",
			"error_code": "admin_recommendation_debug_goal_pending",
			"job_id":     job.ID.String(),
			"feature":    job.Feature,
			"status":     job.Status,
			"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
		})
	}
}

func (h *AdminHandler) ConfirmRecommendationDebugGoal(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugGoalConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.ConfirmedGoal = strings.TrimSpace(req.ConfirmedGoal)
	if req.ConfirmedGoal == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed_goal is required"})
		return
	}
	if len([]rune(req.ConfirmedGoal)) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed_goal too long"})
		return
	}

	profile, err := recommendationDebugGoalProfileFromSnapshot(scenario.GoalProfileSnapshot, adminID, scenario.InitialUserIntent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal interview not started"})
		return
	}
	profile.ConfirmedGoal = &req.ConfirmedGoal
	profile.InterviewState = goal.StateConfirmed
	profile.Messages = append(profile.Messages, goal.InterviewMessage{
		Role:    "lumi",
		Content: "목표가 확정됐어요! 이제 탐험계획을 만들어볼까요? 🎯",
	})
	if !h.saveRecommendationDebugGoalSnapshot(c, scenario.ID, adminID, profile, "goal_chat_confirm", req.ConfirmedGoal) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"goal": profile})
}
