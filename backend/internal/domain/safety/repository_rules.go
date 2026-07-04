package safety

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrRuleDuplicate = errors.New("moderation rule already exists")
var ErrRuleNotFound = errors.New("moderation rule not found")

func (r *Repository) ListRules(ctx context.Context, input ListRulesInput) ([]ModerationRule, error) {
	if r == nil || r.pool == nil {
		return []ModerationRule{}, nil
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
		SELECT id, rule_type, pattern, risk_type, action, locale, description, is_active, created_at, updated_at
		FROM moderation_rules
		WHERE ($1::boolean IS NULL OR is_active = $1)
		  AND ($2 = '' OR action = $2)
		  AND ($3 = '' OR risk_type = $3)
		  AND ($4 = '' OR rule_type = $4)
		  AND ($5 = '' OR COALESCE(locale, '') = $5)
		  AND ($6 = '' OR pattern ILIKE '%' || $6 || '%' OR description ILIKE '%' || $6 || '%')
		ORDER BY is_active DESC, updated_at DESC, created_at DESC
		LIMIT $7 OFFSET $8
	`, input.IsActive, strings.TrimSpace(input.Action), strings.TrimSpace(input.RiskType), strings.TrimSpace(input.RuleType), strings.TrimSpace(input.Locale), strings.TrimSpace(input.Query), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list moderation rules: %w", err)
	}
	defer rows.Close()
	items := make([]ModerationRule, 0)
	for rows.Next() {
		var item ModerationRule
		var action string
		if err := rows.Scan(&item.ID, &item.RuleType, &item.Pattern, &item.RiskType, &action, &item.Locale, &item.Description, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan moderation rule: %w", err)
		}
		item.Action = Action(action)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate moderation rules: %w", err)
	}
	return items, nil
}

func (r *Repository) CreateRule(ctx context.Context, input RuleMutationInput, adminUserID string) (ModerationRule, error) {
	if r == nil || r.pool == nil {
		return ModerationRule{}, nil
	}
	row := r.pool.QueryRow(ctx, `
		INSERT INTO moderation_rules (rule_type, pattern, risk_type, action, locale, description)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), COALESCE($6, ''))
		RETURNING id, rule_type, pattern, risk_type, action, locale, description, is_active, created_at, updated_at
	`, value(input.RuleType), value(input.Pattern), value(input.RiskType), value(input.Action), value(input.Locale), value(input.Description))
	rule, err := scanRule(row)
	if err != nil {
		if isUniqueViolation(err) {
			return ModerationRule{}, ErrRuleDuplicate
		}
		return ModerationRule{}, fmt.Errorf("create moderation rule: %w", err)
	}
	if err := r.insertRuleEvent(ctx, rule.ID, adminUserID, "create", nil, ruleSnapshot(rule), map[string]any{}); err != nil {
		return ModerationRule{}, err
	}
	return rule, nil
}

func (r *Repository) UpdateRule(ctx context.Context, id uuid.UUID, input RuleMutationInput, adminUserID string) (ModerationRule, error) {
	if r == nil || r.pool == nil {
		return ModerationRule{}, nil
	}
	before, err := r.GetRule(ctx, id)
	if err != nil {
		return ModerationRule{}, err
	}
	row := r.pool.QueryRow(ctx, `
		UPDATE moderation_rules
		SET rule_type = COALESCE($2, rule_type),
		    pattern = COALESCE($3, pattern),
		    risk_type = COALESCE($4, risk_type),
		    action = COALESCE($5, action),
		    locale = CASE WHEN $6::text IS NULL THEN locale ELSE NULLIF($6, '') END,
		    description = COALESCE($7, description),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, rule_type, pattern, risk_type, action, locale, description, is_active, created_at, updated_at
	`, id, input.RuleType, input.Pattern, input.RiskType, input.Action, input.Locale, input.Description)
	after, err := scanRule(row)
	if err != nil {
		if isUniqueViolation(err) {
			return ModerationRule{}, ErrRuleDuplicate
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return ModerationRule{}, ErrRuleNotFound
		}
		return ModerationRule{}, fmt.Errorf("update moderation rule: %w", err)
	}
	if err := r.insertRuleEvent(ctx, id, adminUserID, "update", ruleSnapshot(before), ruleSnapshot(after), map[string]any{}); err != nil {
		return ModerationRule{}, err
	}
	return after, nil
}

func (r *Repository) SetRuleActive(ctx context.Context, id uuid.UUID, active bool, adminUserID string) (ModerationRule, error) {
	if r == nil || r.pool == nil {
		return ModerationRule{}, nil
	}
	before, err := r.GetRule(ctx, id)
	if err != nil {
		return ModerationRule{}, err
	}
	row := r.pool.QueryRow(ctx, `
		UPDATE moderation_rules
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, rule_type, pattern, risk_type, action, locale, description, is_active, created_at, updated_at
	`, id, active)
	after, err := scanRule(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ModerationRule{}, ErrRuleNotFound
		}
		return ModerationRule{}, fmt.Errorf("set moderation rule active: %w", err)
	}
	action := "deactivate"
	if active {
		action = "reactivate"
	}
	if err := r.insertRuleEvent(ctx, id, adminUserID, action, ruleSnapshot(before), ruleSnapshot(after), map[string]any{}); err != nil {
		return ModerationRule{}, err
	}
	return after, nil
}

func (r *Repository) GetRule(ctx context.Context, id uuid.UUID) (ModerationRule, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, rule_type, pattern, risk_type, action, locale, description, is_active, created_at, updated_at
		FROM moderation_rules
		WHERE id = $1
	`, id)
	rule, err := scanRule(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ModerationRule{}, ErrRuleNotFound
		}
		return ModerationRule{}, fmt.Errorf("get moderation rule: %w", err)
	}
	return rule, nil
}

func (r *Repository) insertRuleEvent(ctx context.Context, ruleID uuid.UUID, adminUserID string, action string, before map[string]any, after map[string]any, metadata map[string]any) error {
	if before == nil {
		before = map[string]any{}
	}
	if after == nil {
		after = map[string]any{}
	}
	if metadata == nil {
		metadata = map[string]any{}
	}
	beforeBytes, err := json.Marshal(before)
	if err != nil {
		return fmt.Errorf("marshal rule event before: %w", err)
	}
	afterBytes, err := json.Marshal(after)
	if err != nil {
		return fmt.Errorf("marshal rule event after: %w", err)
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal rule event metadata: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO moderation_rule_events (rule_id, admin_user_id, action, before_snapshot, after_snapshot, metadata)
		VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6::jsonb)
	`, ruleID, strings.TrimSpace(adminUserID), action, string(beforeBytes), string(afterBytes), string(metadataBytes))
	if err != nil {
		return fmt.Errorf("insert moderation rule event: %w", err)
	}
	return nil
}

type ruleScanner interface{ Scan(dest ...any) error }

func scanRule(row ruleScanner) (ModerationRule, error) {
	var rule ModerationRule
	var action string
	if err := row.Scan(&rule.ID, &rule.RuleType, &rule.Pattern, &rule.RiskType, &action, &rule.Locale, &rule.Description, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
		return ModerationRule{}, err
	}
	rule.Action = Action(action)
	return rule, nil
}

func ruleSnapshot(rule ModerationRule) map[string]any {
	locale := ""
	if rule.Locale != nil {
		locale = *rule.Locale
	}
	return map[string]any{
		"id":          rule.ID.String(),
		"rule_type":   rule.RuleType,
		"pattern":     rule.Pattern,
		"risk_type":   rule.RiskType,
		"action":      string(rule.Action),
		"locale":      locale,
		"description": rule.Description,
		"is_active":   rule.IsActive,
	}
}

func value(ptr *string) any {
	if ptr == nil {
		return nil
	}
	return *ptr
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
