package admin

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/auth"
)

func (h *AdminHandler) PreviewRecommendationRolloutRouting(c *gin.Context) {
	userID := strings.TrimSpace(c.Query("user_id"))
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if len([]rune(userID)) > 160 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id too long"})
		return
	}
	state, err := h.loadActiveRecommendationRolloutState(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}
	if state == nil {
		c.JSON(http.StatusOK, gin.H{
			"user_id":         userID,
			"active":          false,
			"bucket":          0,
			"applied":         false,
			"mode":            "",
			"reason":          "active rollout state not found",
			"rollout_id":      "",
			"traffic_percent": 0,
		})
		return
	}
	bucket := recommendationRolloutBucket(userID, state.ID)
	applied := bucket < state.TrafficPercent && strings.HasPrefix(state.Mode, "learner_")
	reason := "outside rollout traffic"
	if applied {
		reason = "inside rollout traffic"
	}
	if state.Mode == "shadow_only" || state.Mode == "admin_preview" {
		reason = "mode is admin/shadow only"
	}
	if state.Mode == "rolled_back" {
		reason = "rollout is rolled back"
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":         userID,
		"active":          state.Active,
		"rollout_id":      state.ID,
		"mode":            state.Mode,
		"traffic_percent": state.TrafficPercent,
		"bucket":          bucket,
		"applied":         applied,
		"reason":          reason,
	})
}

func (h *AdminHandler) UpdateRecommendationRolloutRanker(c *gin.Context) {
	var req recommendationRolloutRankerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.RankerModelVersion = strings.TrimSpace(req.RankerModelVersion)
	req.FeatureSchemaVersion = strings.TrimSpace(req.FeatureSchemaVersion)
	if req.FeatureSchemaVersion == "" {
		req.FeatureSchemaVersion = "ranker-feature-v1"
	}
	if len([]rune(req.RankerModelVersion)) > 120 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ranker_model_version too long"})
		return
	}
	if len([]rune(req.FeatureSchemaVersion)) > 80 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feature_schema_version too long"})
		return
	}
	if req.RankerModelVersion != "" && !isSafeRecommendationVersionText(req.RankerModelVersion) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ranker_model_version contains unsupported characters"})
		return
	}
	if !isSafeRecommendationVersionText(req.FeatureSchemaVersion) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feature_schema_version contains unsupported characters"})
		return
	}
	adminID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	state, err := h.loadActiveRecommendationRolloutState(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}
	if state == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "active rollout state not found"})
		return
	}
	stateID, _ := uuid.Parse(state.ID)
	if _, err := h.db.Exec(c.Request.Context(), `
		UPDATE recommendation_rollout_states
		SET ranker_model_version = $2,
		    feature_schema_version = $3,
		    updated_at = NOW()
		WHERE id = $1
	`, stateID, req.RankerModelVersion, req.FeatureSchemaVersion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rollout ranker"})
		return
	}
	if _, err := h.db.Exec(c.Request.Context(), `
		INSERT INTO recommendation_rollout_events (
			rollout_state_id,
			created_by,
			event_type,
			from_mode,
			to_mode,
			payload
		)
		VALUES ($1, $2, 'ranker_config_updated', $3, $3, $4::jsonb)
	`, stateID, adminID, state.Mode, mustJSON(gin.H{
		"previous_ranker_model_version":   state.RankerModelVersion,
		"ranker_model_version":            req.RankerModelVersion,
		"previous_feature_schema_version": state.FeatureSchemaVersion,
		"feature_schema_version":          req.FeatureSchemaVersion,
	})); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record rollout ranker event"})
		return
	}
	nextState, err := h.loadActiveRecommendationRolloutState(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload rollout state"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"state": nextState})
}

func (h *AdminHandler) GetRecommendationRankerArtifactStatus(c *gin.Context) {
	state, err := h.loadActiveRecommendationRolloutState(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"artifact": h.buildRecommendationRankerArtifactStatus(c.Request.Context(), state),
	})
}

func (h *AdminHandler) TransitionRecommendationRollout(c *gin.Context) {
	var req recommendationRolloutTransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Action = strings.TrimSpace(strings.ToLower(req.Action))
	req.Reason = strings.TrimSpace(req.Reason)
	if req.Action != "promote" && req.Action != "rollback" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action must be promote or rollback"})
		return
	}
	if len([]rune(req.Reason)) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason too long"})
		return
	}
	adminID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	state, err := h.loadActiveRecommendationRolloutState(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}
	if state == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "active rollout state not found"})
		return
	}
	fromMode := state.Mode
	toMode := ""
	eventType := ""
	qualityGateStatus := state.QualityGateStatus
	qualityGateSnapshot := strings.TrimSpace(state.QualityGateSnapshot)
	rollbackReason := ""
	approvedAtSQL := "approved_at"
	rolledBackAtSQL := "rolled_back_at"
	trafficPercent := state.TrafficPercent

	readiness, err := h.evaluateRecommendationRolloutReadiness(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to evaluate rollout readiness"})
		return
	}
	readinessJSON, _ := json.Marshal(readiness["snapshot"])
	if req.Action == "promote" {
		canPromote, _ := readiness["can_promote"].(bool)
		if !canPromote {
			c.JSON(http.StatusConflict, gin.H{"error": "rollout checkpoints are not passed"})
			return
		}
		toMode = nextRecommendationRolloutMode(fromMode)
		if toMode == "" || toMode == fromMode || toMode == "shadow_only" {
			c.JSON(http.StatusConflict, gin.H{"error": "rollout cannot be promoted from current mode"})
			return
		}
		eventType = "manual_promote"
		qualityGateStatus = "passed"
		qualityGateSnapshot = string(readinessJSON)
		approvedAtSQL = "NOW()"
		switch toMode {
		case "admin_preview", "shadow_only":
			trafficPercent = 0
		case "learner_10_percent":
			trafficPercent = 10
		case "learner_50_percent":
			trafficPercent = 50
		case "learner_100_percent":
			trafficPercent = 100
		}
	} else {
		toMode = "rolled_back"
		eventType = "manual_rollback"
		qualityGateStatus = "failed"
		qualityGateSnapshot = string(readinessJSON)
		rollbackReason = firstNonEmpty(req.Reason, "manual rollback")
		rolledBackAtSQL = "NOW()"
		trafficPercent = 0
	}

	stateID, _ := uuid.Parse(state.ID)
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to transition rollout"})
		return
	}
	defer tx.Rollback(ctx)

	updateSQL := fmt.Sprintf(`
		UPDATE recommendation_rollout_states
		SET mode = $2,
		    traffic_percent = $3,
		    quality_gate_status = $4,
		    quality_gate_snapshot = $5::jsonb,
		    approved_by = CASE WHEN $6 = 'manual_promote' THEN $7 ELSE approved_by END,
		    approved_at = %s,
		    rolled_back_at = %s,
		    rollback_reason = $8,
		    updated_at = NOW()
		WHERE id = $1
		  AND active = true
	`, approvedAtSQL, rolledBackAtSQL)
	if _, err := tx.Exec(ctx, updateSQL, stateID, toMode, trafficPercent, qualityGateStatus, qualityGateSnapshot, eventType, adminID, rollbackReason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rollout state"})
		return
	}

	payload, _ := json.Marshal(gin.H{
		"reason":    req.Reason,
		"readiness": readiness,
	})
	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_rollout_events (
			rollout_state_id,
			created_by,
			event_type,
			from_mode,
			to_mode,
			payload
		)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, stateID, adminID, eventType, fromMode, toMode, string(payload)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record rollout event"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to transition rollout"})
		return
	}

	nextState, err := h.loadActiveRecommendationRolloutState(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload rollout state"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state":      nextState,
		"from_mode":  fromMode,
		"to_mode":    toMode,
		"event_type": eventType,
	})
}

func nextRecommendationRolloutMode(current string) string {
	switch current {
	case "shadow_only":
		return "admin_preview"
	case "admin_preview":
		return "learner_10_percent"
	case "learner_10_percent":
		return "learner_50_percent"
	case "learner_50_percent":
		return "learner_100_percent"
	case "learner_100_percent":
		return "learner_100_percent"
	case "rolled_back":
		return "shadow_only"
	default:
		return "shadow_only"
	}
}

func recommendationRolloutBucket(userID, rolloutID string) int {
	sum := sha256.Sum256([]byte(strings.TrimSpace(userID) + ":" + strings.TrimSpace(rolloutID)))
	return int(binary.BigEndian.Uint32(sum[:4]) % 100)
}
