package admin

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (h *AdminHandler) loadActiveRecommendationRolloutState(ctx context.Context) (*recommendationRolloutStateRow, error) {
	row := h.db.QueryRow(ctx, `
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
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
