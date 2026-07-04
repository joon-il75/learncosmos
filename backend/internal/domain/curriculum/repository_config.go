package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	secretmanager "github.com/learnweaver/backend/internal/pkg/secretmanager"
)

type LLMSetting struct {
	Feature     string
	Provider    string
	Model       string
	EndpointURL *string
}

type PointTransactionDetail struct {
	Feature       string
	ReferenceType string
	ReferenceID   *uuid.UUID
	Description   string
	Metadata      map[string]any
}

type UserRuntimeAIConfig struct {
	Mode        string
	Provider    string
	APIKey      string
	EndpointURL *string
}

type AIUsageEventInput struct {
	UserID           uuid.UUID
	Source           string
	Provider         string
	Model            string
	Feature          string
	BillingStatus    string
	InputTokens      *int
	OutputTokens     *int
	EstimatedCostUSD *float64
	Success          bool
	ErrorCode        string
	Metadata         map[string]any
}

func (r *Repository) GetPointSetting(ctx context.Context, key string) (int, error) {
	var value int
	err := r.pool.QueryRow(ctx, `
		SELECT value
		FROM point_settings
		WHERE key = $1
	`, key).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("get point setting %s: %w", key, err)
	}
	return value, nil
}

func (r *Repository) RecordAIUsageEvent(ctx context.Context, input AIUsageEventInput) error {
	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = "byok"
	}
	metadata := input.Metadata
	if metadata == nil {
		metadata = map[string]any{}
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO ai_usage_events (
			id, user_id, source, provider, model, feature, billing_status,
			input_tokens, output_tokens, estimated_cost_usd, success, error_code, metadata, created_at
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, NULLIF($7, ''),
		        $8, $9, $10, $11, NULLIF($12, ''), $13::jsonb, NOW())
	`, uuid.New(), input.UserID, source, strings.TrimSpace(input.Provider), strings.TrimSpace(input.Model),
		strings.TrimSpace(input.Feature), strings.TrimSpace(input.BillingStatus), input.InputTokens, input.OutputTokens,
		input.EstimatedCostUSD, input.Success, strings.TrimSpace(input.ErrorCode), string(metadataJSON))
	return err
}

func (r *Repository) UserHasPremiumAccess(ctx context.Context, userID uuid.UUID) (bool, error) {
	var premiumAccess bool
	err := r.pool.QueryRow(ctx, `
		SELECT premium_access
		FROM users
		WHERE id = $1
	`, userID).Scan(&premiumAccess)
	if err != nil {
		return false, fmt.Errorf("get user premium access: %w", err)
	}
	return premiumAccess, nil
}

func (r *Repository) GetLLMSetting(ctx context.Context, feature string) (*LLMSetting, error) {
	var setting LLMSetting
	err := r.pool.QueryRow(ctx, `
		SELECT feature, provider, model, endpoint_url
		FROM system_llm_settings
		WHERE feature = $1
	`, feature).Scan(&setting.Feature, &setting.Provider, &setting.Model, &setting.EndpointURL)
	if err != nil {
		return nil, fmt.Errorf("get llm setting %s: %w", feature, err)
	}
	return &setting, nil
}

type decryptUserAIKeyFunc func(string) (string, error)

func buildUserRuntimeAIConfig(provider string, apiKeyEncrypted *string, endpointURL *string, isEnabled bool, decrypt decryptUserAIKeyFunc) (*UserRuntimeAIConfig, error) {
	config := &UserRuntimeAIConfig{
		Mode:        "managed_credit",
		Provider:    strings.TrimSpace(provider),
		EndpointURL: endpointURL,
	}
	if isEnabled && apiKeyEncrypted != nil && strings.TrimSpace(*apiKeyEncrypted) != "" {
		decrypted, err := decrypt(*apiKeyEncrypted)
		if err != nil {
			return nil, fmt.Errorf("decrypt user ai key: %w", err)
		}
		config.Mode = "byok"
		config.APIKey = strings.TrimSpace(decrypted)
	}
	return config, nil
}

func (r *Repository) GetUserRuntimeAIConfig(ctx context.Context, userID uuid.UUID) (*UserRuntimeAIConfig, error) {
	var provider string
	var apiKeyEncrypted *string
	var endpointURL *string
	var isEnabled bool

	err := r.pool.QueryRow(ctx, `
		SELECT provider, api_key_encrypted, endpoint_url, COALESCE(is_enabled, true)
		FROM user_api_keys
		WHERE user_id = $1
	`, userID).Scan(&provider, &apiKeyEncrypted, &endpointURL, &isEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user ai config: %w", err)
	}

	return buildUserRuntimeAIConfig(provider, apiKeyEncrypted, endpointURL, isEnabled, secretmanager.Decrypt)
}

func (r *Repository) GetPointBalances(ctx context.Context, userID uuid.UUID) (int, int, error) {
	var freeBalance, paidBalance int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(free_balance, 0), COALESCE(paid_balance, 0)
		FROM ai_point_wallets
		WHERE user_id = $1
	`, userID).Scan(&freeBalance, &paidBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("get point balances: %w", err)
	}
	return freeBalance, paidBalance, nil
}

func (r *Repository) ChargePoints(ctx context.Context, userID uuid.UUID, amount int, transactionType string) error {
	return r.ChargePointsWithDetail(ctx, userID, amount, transactionType, nil)
}

func (r *Repository) ChargePointsWithDetail(ctx context.Context, userID uuid.UUID, amount int, transactionType string, detail *PointTransactionDetail) error {
	if amount <= 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin point charge tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var freeBalance, paidBalance int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(free_balance, 0), COALESCE(paid_balance, 0)
		FROM ai_point_wallets
		WHERE user_id = $1
	`, userID).Scan(&freeBalance, &paidBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insufficient_points")
		}
		return fmt.Errorf("load point balance: %w", err)
	}

	if freeBalance+paidBalance < amount {
		return fmt.Errorf("insufficient_points")
	}

	remaining := amount
	newFree := freeBalance
	newPaid := paidBalance

	if newFree >= remaining {
		newFree -= remaining
		remaining = 0
	} else {
		remaining -= newFree
		newFree = 0
	}
	if remaining > 0 {
		newPaid -= remaining
	}

	_, err = tx.Exec(ctx, `
		UPDATE ai_point_wallets
		SET free_balance = $1, paid_balance = $2, updated_at = NOW()
		WHERE user_id = $3
	`, newFree, newPaid, userID)
	if err != nil {
		return fmt.Errorf("update point balance: %w", err)
	}

	metadataJSON, err := marshalPointTransactionMetadata(detail)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_transactions (id, user_id, type, amount, feature, reference_type, reference_id, description, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, NOW())
	`, uuid.New().String(), userID, transactionType, -amount, pointTransactionFeature(detail), pointTransactionReferenceType(detail), pointTransactionReferenceID(detail), pointTransactionDescription(detail), metadataJSON)
	if err != nil {
		return fmt.Errorf("insert point transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit point charge: %w", err)
	}
	return nil
}

func chargePointsTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID, amount int, transactionType string) error {
	return chargePointsTxWithDetail(ctx, tx, userID, amount, transactionType, nil)
}

func chargePointsTxWithDetail(ctx context.Context, tx pgx.Tx, userID uuid.UUID, amount int, transactionType string, detail *PointTransactionDetail) error {
	if amount <= 0 {
		return nil
	}

	var freeBalance, paidBalance int
	err := tx.QueryRow(ctx, `
		SELECT COALESCE(free_balance, 0), COALESCE(paid_balance, 0)
		FROM ai_point_wallets
		WHERE user_id = $1
	`, userID).Scan(&freeBalance, &paidBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insufficient_points")
		}
		return fmt.Errorf("load point balance: %w", err)
	}

	if freeBalance+paidBalance < amount {
		return fmt.Errorf("insufficient_points")
	}

	remaining := amount
	newFree := freeBalance
	newPaid := paidBalance

	if newFree >= remaining {
		newFree -= remaining
		remaining = 0
	} else {
		remaining -= newFree
		newFree = 0
	}
	if remaining > 0 {
		newPaid -= remaining
	}

	_, err = tx.Exec(ctx, `
		UPDATE ai_point_wallets
		SET free_balance = $1, paid_balance = $2, updated_at = NOW()
		WHERE user_id = $3
	`, newFree, newPaid, userID)
	if err != nil {
		return fmt.Errorf("update point balance: %w", err)
	}

	metadataJSON, err := marshalPointTransactionMetadata(detail)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_transactions (id, user_id, type, amount, feature, reference_type, reference_id, description, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, NOW())
	`, uuid.New().String(), userID, transactionType, -amount, pointTransactionFeature(detail), pointTransactionReferenceType(detail), pointTransactionReferenceID(detail), pointTransactionDescription(detail), metadataJSON)
	if err != nil {
		return fmt.Errorf("insert point transaction: %w", err)
	}

	return nil
}

func pointTransactionFeature(detail *PointTransactionDetail) any {
	if detail == nil || strings.TrimSpace(detail.Feature) == "" {
		return nil
	}
	return strings.TrimSpace(detail.Feature)
}

func pointTransactionReferenceType(detail *PointTransactionDetail) any {
	if detail == nil || strings.TrimSpace(detail.ReferenceType) == "" {
		return nil
	}
	return strings.TrimSpace(detail.ReferenceType)
}

func pointTransactionReferenceID(detail *PointTransactionDetail) any {
	if detail == nil || detail.ReferenceID == nil {
		return nil
	}
	return *detail.ReferenceID
}

func pointTransactionDescription(detail *PointTransactionDetail) any {
	if detail == nil || strings.TrimSpace(detail.Description) == "" {
		return nil
	}
	return strings.TrimSpace(detail.Description)
}

func marshalPointTransactionMetadata(detail *PointTransactionDetail) (string, error) {
	if detail == nil || len(detail.Metadata) == 0 {
		return "{}", nil
	}
	payload, err := json.Marshal(detail.Metadata)
	if err != nil {
		return "", fmt.Errorf("marshal point transaction metadata: %w", err)
	}
	return string(payload), nil
}

func (r *Repository) RecordBYOKValidation(ctx context.Context, userID uuid.UUID, provider string, valid bool, message string) error {
	status := "valid"
	var errorText *string
	var lastFailedAt *time.Time
	var nextRetryAt *time.Time
	now := time.Now()

	if !valid {
		status = "invalid"
		trimmed := logsafe.ProviderErrorCode(message)
		errorText = &trimmed
		lastFailedAt = &now
		retry := now.Add(15 * time.Minute)
		nextRetryAt = &retry
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_api_keys (id, user_id, provider, last_validation_status, last_validated_at, last_validation_error, last_failed_at, next_retry_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET last_validation_status = EXCLUDED.last_validation_status,
		    last_validated_at = NOW(),
		    last_validation_error = EXCLUDED.last_validation_error,
		    last_failed_at = EXCLUDED.last_failed_at,
		    next_retry_at = EXCLUDED.next_retry_at,
		    updated_at = NOW()
	`, uuid.New(), userID, provider, status, errorText, lastFailedAt, nextRetryAt)
	return err
}
