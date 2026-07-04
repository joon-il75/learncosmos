package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *AdminHandler) GetRecommendationRolloutState(c *gin.Context) {
	row := h.db.QueryRow(c.Request.Context(), `
		SELECT
			id::text,
			mode,
			traffic_percent,
			embedding_provider,
			embedding_model,
			embedding_dimension,
			ranker_model_version,
			feature_schema_version,
			quality_gate_status,
			quality_gate_snapshot::text,
			rollback_policy_snapshot::text,
			COALESCE(approved_by::text, ''),
			approved_at,
			activated_at,
			rolled_back_at,
			rollback_reason,
			active,
			created_at,
			updated_at
		FROM recommendation_rollout_states
		WHERE active = true
		ORDER BY updated_at DESC
		LIMIT 1
	`)

	var item recommendationRolloutStateRow
	if err := row.Scan(
		&item.ID,
		&item.Mode,
		&item.TrafficPercent,
		&item.EmbeddingProvider,
		&item.EmbeddingModel,
		&item.EmbeddingDimension,
		&item.RankerModelVersion,
		&item.FeatureSchemaVersion,
		&item.QualityGateStatus,
		&item.QualityGateSnapshot,
		&item.RollbackPolicySnapshot,
		&item.ApprovedBy,
		&item.ApprovedAt,
		&item.ActivatedAt,
		&item.RolledBackAt,
		&item.RollbackReason,
		&item.Active,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"state": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": item})
}

func (h *AdminHandler) ListRecommendationRolloutEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			id::text,
			COALESCE(rollout_state_id::text, ''),
			COALESCE(created_by::text, ''),
			event_type,
			from_mode,
			to_mode,
			payload::text,
			created_at
		FROM recommendation_rollout_events
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout events"})
		return
	}
	defer rows.Close()

	items := make([]recommendationRolloutEventRow, 0)
	for rows.Next() {
		var item recommendationRolloutEventRow
		if err := rows.Scan(
			&item.ID,
			&item.RolloutStateID,
			&item.CreatedBy,
			&item.EventType,
			&item.FromMode,
			&item.ToMode,
			&item.Payload,
			&item.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan rollout event"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":  limit,
		"events": items,
	})
}

func (h *AdminHandler) ListRecommendationRolloutMetrics(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if limit <= 0 || limit > 100 {
		limit = 12
	}
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			id::text,
			COALESCE(rollout_state_id::text, ''),
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
			metrics_snapshot::text,
			created_at
		FROM recommendation_rollout_metrics
		ORDER BY window_ended_at DESC, created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout metrics"})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var item recommendationRolloutMetricRow
		if err := rows.Scan(
			&item.ID,
			&item.RolloutStateID,
			&item.WindowStartedAt,
			&item.WindowEndedAt,
			&item.RecommendationCount,
			&item.SelectionRate,
			&item.BrokenLinkRate,
			&item.WrongContentRate,
			&item.FallbackRate,
			&item.P95LatencyMS,
			&item.Top5Overlap,
			&item.GoodFitRate,
			&item.IrrelevantRate,
			&item.DuplicateRate,
			&item.MetricsSnapshot,
			&item.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan rollout metric"})
			return
		}
		items = append(items, gin.H{
			"id":                   item.ID,
			"rollout_state_id":     item.RolloutStateID,
			"window_started_at":    item.WindowStartedAt,
			"window_ended_at":      item.WindowEndedAt,
			"recommendation_count": item.RecommendationCount,
			"selection_rate":       nullFloat64Value(item.SelectionRate),
			"broken_link_rate":     nullFloat64Value(item.BrokenLinkRate),
			"wrong_content_rate":   nullFloat64Value(item.WrongContentRate),
			"fallback_rate":        nullFloat64Value(item.FallbackRate),
			"p95_latency_ms":       nullInt32Value(item.P95LatencyMS),
			"top5_overlap":         nullFloat64Value(item.Top5Overlap),
			"good_fit_rate":        nullFloat64Value(item.GoodFitRate),
			"irrelevant_rate":      nullFloat64Value(item.IrrelevantRate),
			"duplicate_rate":       nullFloat64Value(item.DuplicateRate),
			"metrics_snapshot":     rawJSONOrEmpty(item.MetricsSnapshot),
			"created_at":           item.CreatedAt,
		})
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":   limit,
		"metrics": items,
	})
}
