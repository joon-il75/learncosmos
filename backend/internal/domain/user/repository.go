package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/pkg/apperr"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	secretmanager "github.com/learnweaver/backend/internal/pkg/secretmanager"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindBySocialAccount(ctx context.Context, provider auth.Provider, providerID string) (*auth.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT u.id, COALESCE(u.email, '') AS email, u.nickname, u.role, u.premium_access, u.avatar_url, u.display_id,
		       COALESCE(u.status, 'active') AS status,
		       COALESCE(u.ui_locale, '') AS ui_locale,
		       COALESCE(u.learning_language, '') AS learning_language,
		       u.language_setup_completed_at,
		       u.last_login_at, u.withdrawn_at, u.reactivated_at, u.created_at
		FROM users u
		JOIN social_accounts sa ON sa.user_id = u.id
		WHERE sa.provider = $1 AND sa.provider_id = $2
	`, provider, providerID)

	u := &auth.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Nickname, &u.Role, &u.PremiumAccess, &u.AvatarURL, &u.DisplayID, &u.Status, &u.UILocale, &u.LearningLanguage, &u.LanguageSetupCompletedAt, &u.LastLoginAt, &u.WithdrawnAt, &u.ReactivatedAt, &u.CreatedAt)
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	u.Provider = string(provider)
	return u, nil
}

func (r *PostgresRepository) Create(ctx context.Context, u *auth.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users(id, email, nickname, role, premium_access, avatar_url, display_id, created_at)
		VALUES($1, NULLIF(btrim($2), ''), $3, $4, $5, $6, $7, $8)
	`, u.ID, u.Email, u.Nickname, u.Role, u.PremiumAccess, u.AvatarURL, u.DisplayID, u.CreatedAt)
	return err
}

func (r *PostgresRepository) CreateSocialAccount(ctx context.Context, sa *auth.SocialAccount) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO social_accounts(id, user_id, provider, provider_id, created_at)
		VALUES($1, $2, $3, $4, $5)
	`, sa.ID, sa.UserID, sa.Provider, sa.ProviderID, sa.CreatedAt)
	return err
}

func (r *PostgresRepository) GrantFreePoints(ctx context.Context, userID string, amount int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_wallets(id, user_id, free_balance, updated_at)
		VALUES($1, $2, $3, $4)
		ON CONFLICT(user_id) DO UPDATE
		SET free_balance = ai_point_wallets.free_balance + $3,
		    updated_at = $4
	`, uuid.New().String(), userID, amount, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_transactions(id, user_id, type, amount, feature, description, metadata, created_at)
		VALUES($1, $2, $3, $4, 'welcome_points', '신규 가입 웰컴 포인트', '{}'::jsonb, $5)
	`, uuid.New().String(), userID, "grant_free", amount, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT u.id, COALESCE(u.email, '') AS email, u.nickname, u.role, u.premium_access, u.avatar_url, u.display_id,
		       COALESCE(u.status, 'active') AS status,
		       COALESCE(u.ui_locale, '') AS ui_locale,
		       COALESCE(u.learning_language, '') AS learning_language,
		       u.language_setup_completed_at,
		       u.last_login_at, u.withdrawn_at, u.reactivated_at, u.created_at,
		       COALESCE(sa.provider, '') AS provider
		FROM users u
		LEFT JOIN social_accounts sa ON sa.user_id = u.id
		WHERE u.id = $1
		LIMIT 1
	`, id)

	u := &auth.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Nickname, &u.Role, &u.PremiumAccess, &u.AvatarURL, &u.DisplayID, &u.Status, &u.UILocale, &u.LearningLanguage, &u.LanguageSetupCompletedAt, &u.LastLoginAt, &u.WithdrawnAt, &u.ReactivatedAt, &u.CreatedAt, &u.Provider)
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return u, nil
}

func (r *PostgresRepository) GetPointBalances(ctx context.Context, userID string) (int, int, error) {
	var freeBalance, paidBalance int
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(free_balance, 0), COALESCE(paid_balance, 0)
		FROM ai_point_wallets
		WHERE user_id = $1
	`, userID).Scan(&freeBalance, &paidBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	return freeBalance, paidBalance, nil
}

func (r *PostgresRepository) GetPointSetting(ctx context.Context, key string) (int, error) {
	var value int
	err := r.pool.QueryRow(ctx, `
		SELECT value
		FROM point_settings
		WHERE key = $1
	`, key).Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func normalizePolicyLocale(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func (r *PostgresRepository) GetConsentStatus(ctx context.Context, userID string, locale string) (*auth.ConsentStatus, error) {
	requestedLocale := normalizePolicyLocale(locale)
	rows, err := r.pool.Query(ctx, `
		WITH selected_documents AS (
			SELECT DISTINCT ON (pds.id)
				pds.id AS set_id,
				pd.id,
				pd.type,
				pd.title,
				pd.version,
				pd.locale,
				pd.translation_status,
				($2::text) AS requested_locale,
				(pd.locale <> $2::text) AS fallback_used,
				pd.effective_at,
				pd.updated_at
			FROM policy_document_sets pds
			JOIN policy_documents pd
			  ON pd.set_id = pds.id
			 AND pd.is_active = true
			 AND pd.locale IN ($2::text, 'ko')
			WHERE pds.active = true
			  AND pds.required = true
			ORDER BY pds.id, (pd.locale = $2::text) DESC, (pd.locale = 'ko') DESC
		)
		SELECT pd.id, pd.set_id, pd.type, pd.title, pd.version, pd.locale, pd.translation_status,
		       pd.requested_locale, pd.fallback_used,
		       pd.effective_at, pd.updated_at, upc.agreed_at
		FROM selected_documents pd
		LEFT JOIN user_policy_consents upc
		  ON upc.set_id = pd.set_id
		 AND upc.user_id = $1
		ORDER BY pd.type
	`, userID, requestedLocale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formatTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.UTC().Format("2006-01-02T15:04:05Z")
	}

	status := &auth.ConsentStatus{
		RequiredDocuments: []auth.PolicyDocumentSummary{},
	}
	requiredCount := 0
	agreedCount := 0

	for rows.Next() {
		var docID, setID, docType, title, actualLocale, translationStatus, requestedLocale string
		var version int
		var fallbackUsed bool
		var effectiveAt time.Time
		var updatedAt time.Time
		var agreedAt *time.Time
		if err := rows.Scan(&docID, &setID, &docType, &title, &version, &actualLocale, &translationStatus, &requestedLocale, &fallbackUsed, &effectiveAt, &updatedAt, &agreedAt); err != nil {
			return nil, err
		}

		requiredCount++
		status.RequiredDocuments = append(status.RequiredDocuments, auth.PolicyDocumentSummary{
			ID:                docID,
			SetID:             setID,
			Type:              docType,
			Title:             title,
			Version:           version,
			Locale:            actualLocale,
			TranslationStatus: translationStatus,
			RequestedLocale:   requestedLocale,
			FallbackUsed:      fallbackUsed,
			EffectiveAt:       effectiveAt.UTC().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:         updatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})

		if agreedAt != nil {
			agreedCount++
		}

		switch docType {
		case "terms":
			status.TermsAgreed = agreedAt != nil
			status.TermsAgreedAt = formatTime(agreedAt)
		case "privacy":
			status.PrivacyAgreed = agreedAt != nil
			status.PrivacyAgreedAt = formatTime(agreedAt)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	status.RequiredConsentPending = requiredCount > agreedCount
	return status, nil
}

func (r *PostgresRepository) RecordRequiredConsents(ctx context.Context, userID, locale, ipAddress, userAgent string) error {
	requestedLocale := normalizePolicyLocale(locale)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_policy_consents (
			id, user_id, policy_document_id, set_id, agreed_locale, requested_locale, agreed_at, ip_address, user_agent, created_at
		)
		SELECT DISTINCT ON (pds.id)
		       gen_random_uuid(), $1, pd.id, pds.id, pd.locale, $2, NOW(), NULLIF($3, ''), NULLIF($4, ''), NOW()
		FROM policy_document_sets pds
		JOIN policy_documents pd
		  ON pd.set_id = pds.id
		 AND pd.is_active = true
		 AND pd.locale IN ($2::text, 'ko')
		WHERE pds.active = true
		  AND pds.required = true
		  AND NOT EXISTS (
			SELECT 1
			FROM user_policy_consents upc
			WHERE upc.user_id = $1
			  AND upc.set_id = pds.id
		  )
		ORDER BY pds.id, (pd.locale = $2::text) DESC, (pd.locale = 'ko') DESC
	`, userID, requestedLocale, strings.TrimSpace(ipAddress), strings.TrimSpace(userAgent))
	return err
}

func (r *PostgresRepository) GetUserAISettings(ctx context.Context, userID string) (*auth.UserAISettings, error) {
	var provider string
	var apiKeyEncrypted *string
	var endpointURL *string
	var updatedAt time.Time
	var lastValidationStatus *string
	var lastValidatedAt *time.Time
	var lastValidationError *string
	var lastFailedAt *time.Time
	var nextRetryAt *time.Time
	var isEnabled bool

	err := r.pool.QueryRow(ctx, `
		SELECT provider, api_key_encrypted, endpoint_url, updated_at,
		       last_validation_status, last_validated_at, last_validation_error, last_failed_at, next_retry_at,
		       COALESCE(is_enabled, true)
		FROM user_api_keys
		WHERE user_id = $1
	`, userID).Scan(&provider, &apiKeyEncrypted, &endpointURL, &updatedAt, &lastValidationStatus, &lastValidatedAt, &lastValidationError, &lastFailedAt, &nextRetryAt, &isEnabled)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	hasAPIKey := false
	if apiKeyEncrypted != nil && strings.TrimSpace(*apiKeyEncrypted) != "" {
		hasAPIKey = true
	}

	formatTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.UTC().Format("2006-01-02T15:04:05Z")
	}

	mode := "managed_credit"
	if hasAPIKey && isEnabled {
		mode = "byok"
	}

	return &auth.UserAISettings{
		Mode:                 mode,
		Provider:             provider,
		HasAPIKey:            hasAPIKey,
		IsEnabled:            hasAPIKey && isEnabled,
		EndpointURL:          endpointURL,
		UpdatedAt:            updatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		LastValidationStatus: derefString(lastValidationStatus),
		LastValidatedAt:      formatTime(lastValidatedAt),
		LastValidationError:  derefString(lastValidationError),
		LastFailedAt:         formatTime(lastFailedAt),
		NextRetryAt:          formatTime(nextRetryAt),
	}, nil
}

func (r *PostgresRepository) ListUserAIUsageEvents(ctx context.Context, userID string, startAt, endAt time.Time, limit, offset int) ([]auth.UserAIUsageEvent, int, auth.UserAIUsageSummary, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	var summary auth.UserAIUsageSummary
	if err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(COALESCE(input_tokens, 0) + COALESCE(output_tokens, 0)), 0),
			COALESCE(SUM(estimated_cost_usd), 0)::float8
		FROM ai_usage_events
		WHERE user_id = $1
		  AND source = 'byok'
		  AND created_at >= $2
		  AND created_at < $3
	`, userID, startAt, endAt).Scan(&total, &summary.TotalTokens, &summary.TotalEstimatedCostUSD); err != nil {
		return nil, 0, auth.UserAIUsageSummary{}, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, source, provider, COALESCE(model, ''), feature, COALESCE(billing_status, ''),
		       input_tokens, output_tokens, estimated_cost_usd::float8, success, COALESCE(error_code, ''), created_at
		FROM ai_usage_events
		WHERE user_id = $1
		  AND source = 'byok'
		  AND created_at >= $2
		  AND created_at < $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`, userID, startAt, endAt, limit, offset)
	if err != nil {
		return nil, 0, auth.UserAIUsageSummary{}, err
	}
	defer rows.Close()

	events := make([]auth.UserAIUsageEvent, 0)
	for rows.Next() {
		var event auth.UserAIUsageEvent
		var inputTokens sql.NullInt64
		var outputTokens sql.NullInt64
		var estimatedCost sql.NullFloat64
		var createdAt time.Time
		if err := rows.Scan(
			&event.ID,
			&event.Source,
			&event.Provider,
			&event.Model,
			&event.Feature,
			&event.BillingStatus,
			&inputTokens,
			&outputTokens,
			&estimatedCost,
			&event.Success,
			&event.ErrorCode,
			&createdAt,
		); err != nil {
			return nil, 0, auth.UserAIUsageSummary{}, err
		}
		if inputTokens.Valid {
			value := int(inputTokens.Int64)
			event.InputTokens = &value
		}
		if outputTokens.Valid {
			value := int(outputTokens.Int64)
			event.OutputTokens = &value
		}
		if estimatedCost.Valid {
			value := estimatedCost.Float64
			event.EstimatedCostUSD = &value
		}
		event.CreatedAt = createdAt.UTC().Format("2006-01-02T15:04:05Z")
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, auth.UserAIUsageSummary{}, err
	}
	return events, total, summary, nil
}

func (r *PostgresRepository) ListUserPointTransactions(ctx context.Context, userID string, startAt, endAt time.Time, limit, offset int) ([]auth.UserPointTransaction, int, auth.UserPointUsageSummary, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	var summary auth.UserPointUsageSummary
	if err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN type = 'grant_free' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'purchase' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type IN ('use_course_gen', 'use_lesson_rec') THEN ABS(amount) ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'refund' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(amount), 0)
		FROM ai_point_transactions
		WHERE user_id = $1
		  AND created_at >= $2
		  AND created_at < $3
	`, userID, startAt, endAt).Scan(
		&total,
		&summary.GrantedPoints,
		&summary.PurchasedPoints,
		&summary.UsedPoints,
		&summary.RefundedPoints,
		&summary.NetChange,
	); err != nil {
		return nil, 0, auth.UserPointUsageSummary{}, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, type, amount, COALESCE(feature, ''), COALESCE(reference_type, ''), reference_id::text, COALESCE(description, ''), metadata, created_at
		FROM ai_point_transactions
		WHERE user_id = $1
		  AND created_at >= $2
		  AND created_at < $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`, userID, startAt, endAt, limit, offset)
	if err != nil {
		return nil, 0, auth.UserPointUsageSummary{}, err
	}
	defer rows.Close()

	transactions := make([]auth.UserPointTransaction, 0)
	for rows.Next() {
		var transaction auth.UserPointTransaction
		var referenceID sql.NullString
		var metadataBytes []byte
		var createdAt time.Time
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Feature,
			&transaction.ReferenceType,
			&referenceID,
			&transaction.Description,
			&metadataBytes,
			&createdAt,
		); err != nil {
			return nil, 0, auth.UserPointUsageSummary{}, err
		}
		if referenceID.Valid {
			value := referenceID.String
			transaction.ReferenceID = &value
		}
		if len(metadataBytes) > 0 && string(metadataBytes) != "{}" {
			var metadata map[string]any
			if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
				transaction.Metadata = metadata
			}
		}
		transaction.CreatedAt = createdAt.UTC().Format("2006-01-02T15:04:05Z")
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, auth.UserPointUsageSummary{}, err
	}
	return transactions, total, summary, nil
}

func (r *PostgresRepository) UpsertUserAISettings(ctx context.Context, userID, provider, apiKey string, endpointURL *string, isEnabled bool) error {
	encrypted, err := secretmanager.Encrypt(strings.TrimSpace(apiKey))
	if err != nil {
		return err
	}

	var normalizedEndpoint *string
	if endpointURL != nil {
		trimmed := strings.TrimSpace(*endpointURL)
		if trimmed != "" {
			normalizedEndpoint = &trimmed
		}
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO user_api_keys (id, user_id, provider, api_key_encrypted, endpoint_url, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET provider = EXCLUDED.provider,
		    api_key_encrypted = EXCLUDED.api_key_encrypted,
		    endpoint_url = EXCLUDED.endpoint_url,
		    is_enabled = EXCLUDED.is_enabled,
		    updated_at = NOW()
	`, uuid.New().String(), userID, provider, encrypted, normalizedEndpoint, isEnabled)
	return err
}

func (r *PostgresRepository) UpdateUserAISettingsEnabled(ctx context.Context, userID string, isEnabled bool) error {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE user_api_keys
		SET is_enabled = $2,
		    updated_at = NOW()
		WHERE user_id = $1
		  AND api_key_encrypted IS NOT NULL
		  AND api_key_encrypted <> ''
	`, userID, isEnabled)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *PostgresRepository) DeleteUserAISettings(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM user_api_keys
		WHERE user_id = $1
	`, userID)
	return err
}

func (r *PostgresRepository) RecordUserAIValidation(ctx context.Context, userID, provider string, valid bool, message string) error {
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
		SET provider = EXCLUDED.provider,
		    last_validation_status = EXCLUDED.last_validation_status,
		    last_validated_at = NOW(),
		    last_validation_error = EXCLUDED.last_validation_error,
		    last_failed_at = EXCLUDED.last_failed_at,
		    next_retry_at = EXCLUDED.next_retry_at,
		    updated_at = NOW()
	`, uuid.New().String(), userID, provider, status, errorText, lastFailedAt, nextRetryAt)
	return err
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (r *PostgresRepository) DisplayIDExists(ctx context.Context, displayID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE display_id = $1)`, displayID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) UpdateNickname(ctx context.Context, userID, nickname string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET nickname = $1, updated_at = NOW() WHERE id = $2`, nickname, userID)
	return err
}

func (r *PostgresRepository) UpdateEmail(ctx context.Context, userID, email string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET email = NULLIF(btrim($1), ''), updated_at = NOW() WHERE id = $2`, email, userID)
	return err
}

func (r *PostgresRepository) UpdateLanguagePreferences(ctx context.Context, userID, uiLocale, learningLanguage string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users
		SET ui_locale = $1,
		    learning_language = $2,
		    language_setup_completed_at = COALESCE(language_setup_completed_at, NOW()),
		    updated_at = NOW()
		WHERE id = $3
	`, normalizePolicyLocale(uiLocale), normalizePolicyLocale(learningLanguage), userID)
	return err
}

func (r *PostgresRepository) EmailExists(ctx context.Context, email, excludeUserID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email IS NOT NULL
			  AND btrim(email) <> ''
			  AND lower(email) = lower(btrim($1))
			  AND id != $2
		)`, email, excludeUserID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) UpdateAvatarURL(ctx context.Context, userID, avatarURL string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`, avatarURL, userID)
	return err
}

func (r *PostgresRepository) UpdateOnLogin(ctx context.Context, userID, email, nickname, avatarURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users
		SET
		  email      = CASE WHEN email IS NULL OR email = '' THEN NULLIF(btrim($1), '') ELSE email END,
		  nickname   = CASE WHEN nickname IS NULL OR nickname = '' THEN $2 ELSE nickname END,
		  avatar_url = CASE
		    WHEN avatar_url IS NULL OR avatar_url = '' OR avatar_url NOT LIKE '/api/v1/users/avatars/%' THEN $3
		    ELSE avatar_url
		  END,
		  updated_at = NOW()
		WHERE id = $4
	`, email, nickname, avatarURL, userID)
	return err
}

func (r *PostgresRepository) ReactivateUser(ctx context.Context, userID, email, nickname, avatarURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users
		SET
		  status         = 'active',
		  withdrawn_at   = NULL,
		  reactivated_at = NOW(),
		  email          = CASE WHEN email IS NULL OR email = '' THEN NULLIF(btrim($1), '') ELSE email END,
		  nickname       = CASE WHEN nickname IS NULL OR nickname = '' OR nickname = '탈퇴한 학습자' THEN $2 ELSE nickname END,
		  avatar_url     = CASE
		    WHEN avatar_url IS NULL OR avatar_url = '' OR avatar_url NOT LIKE '/api/v1/users/avatars/%' THEN $3
		    ELSE avatar_url
		  END,
		  updated_at     = NOW()
		WHERE id = $4
	`, email, nickname, avatarURL, userID)
	return err
}

func (r *PostgresRepository) WithdrawUser(ctx context.Context, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE users
		SET status = 'withdrawn',
		    withdrawn_at = NOW(),
		    nickname = '탈퇴한 학습자',
		    avatar_url = '',
		    updated_at = NOW()
		WHERE id = $1
	`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM user_api_keys
		WHERE user_id = $1
	`, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
