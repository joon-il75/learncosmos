package safety

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListActiveRules(ctx context.Context, locale string) ([]Rule, error) {
	if r == nil || r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, rule_type, pattern, risk_type, action, locale, description
		FROM moderation_rules
		WHERE is_active = TRUE
		ORDER BY
		  CASE action WHEN 'block' THEN 0 ELSE 1 END,
		  created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list active moderation rules: %w", err)
	}
	defer rows.Close()

	rules := make([]Rule, 0)
	for rows.Next() {
		var rule Rule
		var id uuid.UUID
		var action string
		if err := rows.Scan(&id, &rule.RuleType, &rule.Pattern, &rule.RiskType, &action, &rule.Locale, &rule.Description); err != nil {
			return nil, fmt.Errorf("scan moderation rule: %w", err)
		}
		rule.ID = &id
		rule.Action = Action(action)
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate moderation rules: %w", err)
	}
	return rules, nil
}

func (r *Repository) InsertLog(ctx context.Context, input ModerateInput, result *ModerationResult, hash string) error {
	if r == nil || r.pool == nil || result == nil {
		return nil
	}
	metadata := sanitizeSafetyMetadata(input.Metadata)
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal moderation metadata: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO moderation_logs (
			user_id, target_type, target_id, input_text_hash, risk_type, action,
			matched_rule_id, route, locale, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULLIF($9, ''), $10::jsonb)
	`, input.UserID, input.TargetType, input.TargetID, hash, result.RiskType, string(result.Action), result.MatchedRuleID, input.Route, strings.TrimSpace(input.Locale), string(metadataBytes))
	if err != nil {
		return fmt.Errorf("insert moderation log: %w", err)
	}
	return nil
}

func (r *Repository) ListLogs(ctx context.Context, input ListLogsInput) ([]LogEntry, error) {
	if r == nil || r.pool == nil {
		return []LogEntry{}, nil
	}
	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			ml.id, ml.user_id, u.email, ml.target_type, ml.target_id, ml.input_text_hash,
			ml.risk_type, ml.action, ml.matched_rule_id, ml.route, ml.locale, ml.metadata,
			ml.created_at, mr.pattern, mr.rule_type
		FROM moderation_logs ml
		LEFT JOIN users u ON u.id = ml.user_id
		LEFT JOIN moderation_rules mr ON mr.id = ml.matched_rule_id
		WHERE ($1 = '' OR ml.action = $1)
		  AND ($2 = '' OR ml.risk_type = $2)
		  AND ($3 = '' OR ml.target_type = $3)
		ORDER BY ml.created_at DESC
		LIMIT $4 OFFSET $5
	`, strings.TrimSpace(input.Action), strings.TrimSpace(input.RiskType), strings.TrimSpace(input.TargetType), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list moderation logs: %w", err)
	}
	defer rows.Close()

	items := make([]LogEntry, 0)
	for rows.Next() {
		var item LogEntry
		var action string
		var metadataBytes []byte
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserEmail, &item.TargetType, &item.TargetID, &item.InputTextHash,
			&item.RiskType, &action, &item.MatchedRuleID, &item.Route, &item.Locale, &metadataBytes,
			&item.CreatedAt, &item.MatchedPattern, &item.RuleType,
		); err != nil {
			return nil, fmt.Errorf("scan moderation log: %w", err)
		}
		item.Action = Action(action)
		if len(metadataBytes) > 0 {
			if err := json.Unmarshal(metadataBytes, &item.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal moderation metadata: %w", err)
			}
		}
		item.Metadata = sanitizeSafetyMetadata(item.Metadata)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("iterate moderation logs: %w", err)
	}
	return items, nil
}

func (r *Repository) ListAIOutputReviewDecisions(ctx context.Context, candidateKeys []string) (map[string]AIOutputReviewDecision, error) {
	decisions := map[string]AIOutputReviewDecision{}
	if r == nil || r.pool == nil || len(candidateKeys) == 0 {
		return decisions, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT candidate_key, status, note, reviewed_by, reviewed_at, metadata, created_at, updated_at
		FROM moderation_ai_output_review_decisions
		WHERE candidate_key = ANY($1::text[])
	`, candidateKeys)
	if err != nil {
		return nil, fmt.Errorf("list ai output review decisions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var decision AIOutputReviewDecision
		var status string
		var metadataBytes []byte
		if err := rows.Scan(&decision.CandidateKey, &status, &decision.Note, &decision.ReviewedBy, &decision.ReviewedAt, &metadataBytes, &decision.CreatedAt, &decision.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan ai output review decision: %w", err)
		}
		decision.Status = AIOutputReviewStatus(status)
		if len(metadataBytes) > 0 {
			if err := json.Unmarshal(metadataBytes, &decision.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal ai output review metadata: %w", err)
			}
		}
		decision.Metadata = sanitizeSafetyMetadata(decision.Metadata)
		decisions[decision.CandidateKey] = decision
	}
	if err := rows.Err(); err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("iterate ai output review decisions: %w", err)
	}
	return decisions, nil
}

func (r *Repository) UpsertAIOutputReviewDecision(ctx context.Context, input AIOutputReviewDecisionInput) (AIOutputReviewDecision, error) {
	if r == nil || r.pool == nil {
		return AIOutputReviewDecision{
			CandidateKey: strings.TrimSpace(input.CandidateKey),
			Status:       AIOutputReviewStatus(strings.TrimSpace(input.Status)),
			Note:         strings.TrimSpace(input.Note),
			ReviewedBy:   strings.TrimSpace(input.ReviewedBy),
			Metadata:     sanitizeSafetyMetadata(input.Metadata),
		}, nil
	}
	metadata := sanitizeSafetyMetadata(input.Metadata)
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return AIOutputReviewDecision{}, fmt.Errorf("marshal ai output review metadata: %w", err)
	}
	var decision AIOutputReviewDecision
	var status string
	var storedMetadataBytes []byte
	err = r.pool.QueryRow(ctx, `
		INSERT INTO moderation_ai_output_review_decisions (candidate_key, status, note, reviewed_by, reviewed_at, metadata)
		VALUES ($1, $2, $3, $4, NOW(), $5::jsonb)
		ON CONFLICT (candidate_key) DO UPDATE SET
		  status = EXCLUDED.status,
		  note = EXCLUDED.note,
		  reviewed_by = EXCLUDED.reviewed_by,
		  reviewed_at = NOW(),
		  metadata = EXCLUDED.metadata,
		  updated_at = NOW()
		RETURNING candidate_key, status, note, reviewed_by, reviewed_at, metadata, created_at, updated_at
	`, strings.TrimSpace(input.CandidateKey), strings.TrimSpace(input.Status), strings.TrimSpace(input.Note), strings.TrimSpace(input.ReviewedBy), string(metadataBytes)).Scan(
		&decision.CandidateKey, &status, &decision.Note, &decision.ReviewedBy, &decision.ReviewedAt, &storedMetadataBytes, &decision.CreatedAt, &decision.UpdatedAt,
	)
	if err != nil {
		return AIOutputReviewDecision{}, fmt.Errorf("upsert ai output review decision: %w", err)
	}
	decision.Status = AIOutputReviewStatus(status)
	if len(storedMetadataBytes) > 0 {
		if err := json.Unmarshal(storedMetadataBytes, &decision.Metadata); err != nil {
			return AIOutputReviewDecision{}, fmt.Errorf("unmarshal stored ai output review metadata: %w", err)
		}
	}
	decision.Metadata = sanitizeSafetyMetadata(decision.Metadata)
	return decision, nil
}

func (r *Repository) ListAIOutputObservations(ctx context.Context, maxHours int) ([]AIOutputLogObservation, error) {
	if r == nil || r.pool == nil {
		return []AIOutputLogObservation{}, nil
	}
	if maxHours <= 0 || maxHours > 24*30 {
		maxHours = 24 * 7
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			ml.target_type,
			ml.target_id,
			ml.action,
			ml.risk_type,
			COALESCE(ml.metadata->>'source_feature', ''),
			COALESCE(ml.matched_rule_id::text, ''),
			COALESCE(mr.pattern, ''),
			ml.created_at
		FROM moderation_logs ml
		LEFT JOIN moderation_rules mr ON mr.id = ml.matched_rule_id
		WHERE ml.metadata->>'direction' = 'ai_output'
		  AND ml.created_at >= NOW() - ($1::int * interval '1 hour')
		ORDER BY ml.created_at DESC
	`, maxHours)
	if err != nil {
		return nil, fmt.Errorf("list ai output observations: %w", err)
	}
	defer rows.Close()

	items := make([]AIOutputLogObservation, 0)
	for rows.Next() {
		var item AIOutputLogObservation
		var action string
		if err := rows.Scan(&item.TargetType, &item.TargetID, &action, &item.RiskType, &item.SourceFeature, &item.MatchedRuleID, &item.MatchedPattern, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ai output observation: %w", err)
		}
		item.Action = Action(action)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("iterate ai output observations: %w", err)
	}
	return items, nil
}
