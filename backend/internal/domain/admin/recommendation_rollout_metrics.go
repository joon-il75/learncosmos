package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AdminHandler) CreateRecommendationRolloutMetricSnapshot(c *gin.Context) {
	var req recommendationRolloutMetricSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.WindowHours <= 0 {
		req.WindowHours = 24
	}
	if req.WindowHours > 24*30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "window_hours too large"})
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
	windowEndedAt := time.Now()
	windowStartedAt := windowEndedAt.Add(-time.Duration(req.WindowHours) * time.Hour)

	var recommendationCount int
	var p95LatencyVal sql.NullFloat64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			COUNT(*),
			percentile_cont(0.95) WITHIN GROUP (ORDER BY latency_ms)::float8
		FROM recommendation_debug_runs
		WHERE run_type = 'recommendation_compare'
		  AND created_at >= $1
		  AND created_at < $2
		  AND latency_ms IS NOT NULL
	`, windowStartedAt, windowEndedAt).Scan(&recommendationCount, &p95LatencyVal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate recommendation metrics"})
		return
	}

	var shadowAttempted int64
	var shadowSucceeded int64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			COALESCE(SUM(CASE WHEN metrics_snapshot ? 'shadow_attempted' THEN (metrics_snapshot->>'shadow_attempted')::integer ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN metrics_snapshot ? 'shadow_succeeded' THEN (metrics_snapshot->>'shadow_succeeded')::integer ELSE 0 END), 0)
		FROM recommendation_debug_runs
		WHERE run_type = 'recommendation_compare'
		  AND created_at >= $1
		  AND created_at < $2
	`, windowStartedAt, windowEndedAt).Scan(&shadowAttempted, &shadowSucceeded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate shadow metrics"})
		return
	}

	var top5OverlapVal sql.NullFloat64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT AVG(NULLIF(lesson->'metrics'->>'shadow_overlap', '')::float8)
		FROM recommendation_debug_runs r
		CROSS JOIN LATERAL jsonb_array_elements(COALESCE(r.baseline_result_snapshot->'lessons', '[]'::jsonb)) lesson
		WHERE r.run_type = 'recommendation_compare'
		  AND r.created_at >= $1
		  AND r.created_at < $2
		  AND lesson->'metrics' ? 'shadow_overlap'
	`, windowStartedAt, windowEndedAt).Scan(&top5OverlapVal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate overlap metrics"})
		return
	}

	var labelTotal int64
	var goodFit int64
	var irrelevant int64
	var duplicate int64
	var brokenLink int64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE label = 'good_fit'),
			COUNT(*) FILTER (WHERE label = 'irrelevant'),
			COUNT(*) FILTER (WHERE label = 'duplicate'),
			COUNT(*) FILTER (WHERE label = 'broken_link')
		FROM recommendation_debug_labels
		WHERE updated_at >= $1
		  AND updated_at < $2
	`, windowStartedAt, windowEndedAt).Scan(&labelTotal, &goodFit, &irrelevant, &duplicate, &brokenLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate label metrics"})
		return
	}

	var learnerRecommendationCount int64
	var learnerReplacementCount int64
	var learnerBrokenLinkCount int64
	var learnerWrongContentCount int64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			COUNT(*) FILTER (WHERE event_type = 'recommendation_exposed'),
			COUNT(*) FILTER (WHERE event_type = 'material_replaced'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'broken_link'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'wrong_content')
		FROM recommendation_rollout_learner_events
		WHERE rollout_state_id = $1
		  AND created_at >= $2
		  AND created_at < $3
	`, stateID, windowStartedAt, windowEndedAt).Scan(&learnerRecommendationCount, &learnerReplacementCount, &learnerBrokenLinkCount, &learnerWrongContentCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate learner rollout metrics"})
		return
	}

	var fallbackRate any
	if shadowAttempted > 0 {
		fallbackRate = safeRate(shadowAttempted-shadowSucceeded, shadowAttempted)
	}
	var p95LatencyMS any
	if p95LatencyVal.Valid {
		p95LatencyMS = int(p95LatencyVal.Float64 + 0.5)
	}
	var top5OverlapValue any
	if top5OverlapVal.Valid {
		top5OverlapValue = top5OverlapVal.Float64
	}
	var goodFitRate any
	var irrelevantRate any
	var duplicateRate any
	var brokenLinkRate any
	if labelTotal > 0 {
		goodFitRate = safeRate(goodFit, labelTotal)
		irrelevantRate = safeRate(irrelevant, labelTotal)
		duplicateRate = safeRate(duplicate, labelTotal)
		brokenLinkRate = safeRate(brokenLink, labelTotal)
	}
	var selectionRate any
	var wrongContentRate any
	snapshotRecommendationCount := recommendationCount
	if learnerRecommendationCount > 0 {
		snapshotRecommendationCount = int(learnerRecommendationCount)
		selectionRate = safeRate(learnerReplacementCount, learnerRecommendationCount)
		brokenLinkRate = safeRate(learnerBrokenLinkCount, learnerRecommendationCount)
		wrongContentRate = safeRate(learnerWrongContentCount, learnerRecommendationCount)
	}
	metricsSnapshot, _ := json.Marshal(gin.H{
		"source":                         "manual_super_admin_snapshot",
		"window_hours":                   req.WindowHours,
		"shadow_attempted":               shadowAttempted,
		"shadow_succeeded":               shadowSucceeded,
		"label_total":                    labelTotal,
		"admin_recommendation_runs":      recommendationCount,
		"learner_recommendation_exposed": learnerRecommendationCount,
		"learner_material_replaced":      learnerReplacementCount,
		"learner_broken_link_reports":    learnerBrokenLinkCount,
		"learner_wrong_content_reports":  learnerWrongContentCount,
	})

	var saved gin.H
	var savedID string
	var createdAt time.Time
	if err := h.db.QueryRow(c.Request.Context(), `
		INSERT INTO recommendation_rollout_metrics (
			rollout_state_id,
			window_started_at,
			window_ended_at,
			recommendation_count,
			selection_rate,
			broken_link_rate,
			wrong_content_rate,
			fallback_rate,
			p95_latency_ms,
			top5_overlap,
			good_fit_rate,
			irrelevant_rate,
			duplicate_rate,
			metrics_snapshot
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb)
		RETURNING id::text, created_at
	`, stateID, windowStartedAt, windowEndedAt, snapshotRecommendationCount, selectionRate, brokenLinkRate, wrongContentRate, fallbackRate, p95LatencyMS, top5OverlapValue, goodFitRate, irrelevantRate, duplicateRate, string(metricsSnapshot)).Scan(&savedID, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rollout metric snapshot"})
		return
	}
	saved = gin.H{
		"id":                   savedID,
		"rollout_state_id":     state.ID,
		"window_started_at":    windowStartedAt,
		"window_ended_at":      windowEndedAt,
		"recommendation_count": snapshotRecommendationCount,
		"selection_rate":       selectionRate,
		"broken_link_rate":     brokenLinkRate,
		"wrong_content_rate":   wrongContentRate,
		"fallback_rate":        fallbackRate,
		"p95_latency_ms":       p95LatencyMS,
		"top5_overlap":         top5OverlapValue,
		"good_fit_rate":        goodFitRate,
		"irrelevant_rate":      irrelevantRate,
		"duplicate_rate":       duplicateRate,
		"metrics_snapshot":     rawJSONOrEmpty(string(metricsSnapshot)),
		"created_at":           createdAt,
	}
	c.JSON(http.StatusOK, gin.H{"metric": saved})
}
