package demoaccount

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidCode        = errors.New("invalid demo code")
	ErrExpiredAccount     = errors.New("expired demo account")
	ErrDisabledAccount    = errors.New("disabled demo account")
	ErrUnavailableAccount = errors.New("demo account unavailable")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateBatch(ctx context.Context, req CreateBatchRequest, loginSecret string) ([]CreatedAccount, error) {
	normalized := normalizeCreateRequest(req)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	created := make([]CreatedAccount, 0, normalized.Count)
	now := time.Now()
	for i := 0; i < normalized.Count; i++ {
		accountNumber := normalized.StartNumber + i
		numberText := fmt.Sprintf("%02d", accountNumber)
		userID := uuid.New().String()
		demoAccountID := uuid.New().String()
		walletID := uuid.New().String()
		transactionID := uuid.New().String()
		label := fmt.Sprintf("%s%s", normalized.LabelPrefix, numberText)
		email := fmt.Sprintf("%s%s@%s", normalized.EmailPrefix, numberText, normalized.EmailDomain)
		nickname := fmt.Sprintf("Demo %s", numberText)

		loginCode, err := GenerateLoginCode(accountNumber)
		if err != nil {
			return nil, err
		}
		loginCodeHash, err := HashLoginCode(loginCode, loginSecret)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO users (
				id, email, nickname, role, premium_access, avatar_url,
				status, ui_locale, learning_language, language_setup_completed_at,
				alpha_access_granted_at, created_at, updated_at
			)
			VALUES (
				$1, $2, $3, 'learner', false, $4,
				'active', $5, $6, $7,
				$7, $7, $7
			)
		`, userID, email, nickname, DefaultAvatarURL, normalized.UILocale, normalized.LearningLanguage, now)
		if err != nil {
			return nil, fmt.Errorf("insert demo user %s: %w", label, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO demo_accounts (
				id, user_id, label, email, login_code_hash, login_code,
				assigned_to, assignment_note, expires_at,
				created_by, created_by_actor, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), NULLIF($8, ''), $9, $10, $11, $12, $12)
		`, demoAccountID, userID, label, email, loginCodeHash, NormalizeLoginCode(loginCode), normalized.AssignedTo, normalized.AssignmentNote, normalized.ExpiresAt, normalized.CreatedBy, normalized.CreatedByActor, now)
		if err != nil {
			return nil, fmt.Errorf("insert demo account %s: %w", label, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO ai_point_wallets(id, user_id, free_balance, paid_balance, updated_at)
			VALUES($1, $2, $3, 0, $4)
		`, walletID, userID, normalized.InitialPoints, now)
		if err != nil {
			return nil, fmt.Errorf("insert demo wallet %s: %w", label, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO ai_point_transactions(id, user_id, type, amount, created_at)
			VALUES($1, $2, 'grant_free', $3, $4)
		`, transactionID, userID, normalized.InitialPoints, now)
		if err != nil {
			return nil, fmt.Errorf("insert demo point transaction %s: %w", label, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_policy_consents (
				id, user_id, policy_document_id, set_id, agreed_locale, requested_locale,
				agreed_at, ip_address, user_agent, created_at
			)
			SELECT DISTINCT ON (pds.id)
			       gen_random_uuid(), $1, pd.id, pds.id, pd.locale, $2,
			       $3, NULL, NULL, $3
			FROM policy_document_sets pds
			JOIN policy_documents pd
			  ON pd.set_id = pds.id
			 AND pd.is_active = true
			 AND pd.locale IN ($2::text, 'ko')
			WHERE pds.active = true
			  AND pds.required = true
			ORDER BY pds.id, (pd.locale = $2::text) DESC, (pd.locale = 'ko') DESC
		`, userID, normalized.UILocale, now)
		if err != nil {
			return nil, fmt.Errorf("insert demo consents %s: %w", label, err)
		}

		created = append(created, CreatedAccount{
			ID:        demoAccountID,
			UserID:    userID,
			Label:     label,
			Email:     email,
			Nickname:  nickname,
			LoginCode: loginCode,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *Repository) AuthenticateByCode(ctx context.Context, code string, loginSecret string) (*AuthenticatedUser, error) {
	loginCodeHash, err := HashLoginCode(code, loginSecret)
	if err != nil {
		return nil, ErrInvalidCode
	}

	row := r.pool.QueryRow(ctx, `
		SELECT
			da.id,
			u.id, u.email, u.nickname, u.role, u.premium_access, u.avatar_url, u.display_id,
			COALESCE(u.status, 'active') AS status,
			COALESCE(u.ui_locale, '') AS ui_locale,
			COALESCE(u.learning_language, '') AS learning_language,
			u.language_setup_completed_at,
			u.last_login_at, u.withdrawn_at, u.reactivated_at, u.created_at,
			da.disabled_at, da.expires_at
		FROM demo_accounts da
		JOIN users u ON u.id = da.user_id
		WHERE da.login_code_hash = $1
	`, loginCodeHash)

	var user AuthenticatedUser
	var disabledAt sql.NullTime
	var expiresAt sql.NullTime
	if err := row.Scan(
		&user.DemoAccountID,
		&user.ID, &user.Email, &user.Nickname, &user.Role, &user.PremiumAccess, &user.AvatarURL, &user.DisplayID,
		&user.Status, &user.UILocale, &user.LearningLanguage, &user.LanguageSetupCompletedAt,
		&user.LastLoginAt, &user.WithdrawnAt, &user.ReactivatedAt, &user.CreatedAt,
		&disabledAt, &expiresAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCode
		}
		return nil, err
	}

	now := time.Now()
	if disabledAt.Valid {
		return nil, ErrDisabledAccount
	}
	if expiresAt.Valid && expiresAt.Time.Before(now) {
		return nil, ErrExpiredAccount
	}
	if user.Status != "active" {
		return nil, ErrUnavailableAccount
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE demo_accounts
		SET activated_at = COALESCE(activated_at, $1),
		    last_used_at = $1,
		    updated_at = $1
		WHERE id = $2
	`, now, user.DemoAccountID)
	if err != nil {
		return nil, err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE users
		SET last_login_at = $1, updated_at = $1
		WHERE id = $2
	`, now, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) List(ctx context.Context, opts ListOptions) (ListResult, error) {
	opts = normalizeListOptions(opts)
	q := "%" + strings.ToLower(strings.TrimSpace(opts.Query)) + "%"
	offset := (opts.Page - 1) * opts.Limit

	statusFilter := `
		AND (
			$2 = 'all'
			OR ($2 = 'active' AND da.disabled_at IS NULL AND (da.expires_at IS NULL OR da.expires_at > NOW()) AND u.status = 'active')
			OR ($2 = 'disabled' AND da.disabled_at IS NOT NULL)
			OR ($2 = 'expired' AND da.disabled_at IS NULL AND da.expires_at IS NOT NULL AND da.expires_at <= NOW())
		)
	`

	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM demo_accounts da
		JOIN users u ON u.id = da.user_id
		WHERE (
			LOWER(da.label) LIKE $1
			OR LOWER(da.email) LIKE $1
			OR LOWER(u.nickname) LIKE $1
			OR LOWER(COALESCE(da.assigned_to, '')) LIKE $1
		)
	`+statusFilter, q, opts.Status).Scan(&total)
	if err != nil {
		return ListResult{}, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			da.id, da.user_id, da.label, da.email, u.nickname,
			COALESCE(da.assigned_to, ''), COALESCE(da.assignment_note, ''),
			da.login_code,
			CASE
				WHEN da.disabled_at IS NOT NULL THEN 'disabled'
				WHEN da.expires_at IS NOT NULL AND da.expires_at <= NOW() THEN 'expired'
				WHEN u.status <> 'active' THEN 'unavailable'
				ELSE 'active'
			END AS status,
			da.expires_at, da.activated_at, da.last_used_at, da.disabled_at,
			COALESCE(w.free_balance, 0), COALESCE(w.paid_balance, 0),
			da.created_at, da.updated_at
		FROM demo_accounts da
		JOIN users u ON u.id = da.user_id
		LEFT JOIN ai_point_wallets w ON w.user_id = u.id
		WHERE (
			LOWER(da.label) LIKE $1
			OR LOWER(da.email) LIKE $1
			OR LOWER(u.nickname) LIKE $1
			OR LOWER(COALESCE(da.assigned_to, '')) LIKE $1
		)
	`+statusFilter+`
		ORDER BY da.created_at DESC
		LIMIT $3 OFFSET $4
	`, q, opts.Status, opts.Limit, offset)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()

	items := make([]AccountSummary, 0)
	for rows.Next() {
		var item AccountSummary
		var expiresAt sql.NullTime
		var activatedAt sql.NullTime
		var lastUsedAt sql.NullTime
		var disabledAt sql.NullTime
		var loginCode sql.NullString
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Label, &item.Email, &item.Nickname,
			&item.AssignedTo, &item.AssignmentNote, &loginCode, &item.Status,
			&expiresAt, &activatedAt, &lastUsedAt, &disabledAt,
			&item.FreePoints, &item.PaidPoints,
			&createdAt, &updatedAt,
		); err != nil {
			return ListResult{}, err
		}
		item.ExpiresAt = formatNullableTime(expiresAt)
		if loginCode.Valid && strings.TrimSpace(loginCode.String) != "" {
			code := loginCode.String
			item.LoginCode = &code
		}
		item.ActivatedAt = formatNullableTime(activatedAt)
		item.LastUsedAt = formatNullableTime(lastUsedAt)
		item.DisabledAt = formatNullableTime(disabledAt)
		item.CreatedAt = formatTime(createdAt)
		item.UpdatedAt = formatTime(updatedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}

	return ListResult{Items: items, Total: total, Page: opts.Page, Limit: opts.Limit}, nil
}

func (r *Repository) RotateCode(ctx context.Context, id string, loginSecret string) (string, error) {
	var label string
	err := r.pool.QueryRow(ctx, `
		SELECT label
		FROM demo_accounts
		WHERE id = $1
	`, id).Scan(&label)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCode
		}
		return "", err
	}

	code, err := GenerateLoginCodeForLabel(label)
	if err != nil {
		return "", err
	}
	hash, err := HashLoginCode(code, loginSecret)
	if err != nil {
		return "", err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE demo_accounts
		SET login_code_hash = $2, login_code = $3, updated_at = NOW()
		WHERE id = $1
	`, id, hash, NormalizeLoginCode(code))
	if err != nil {
		return "", err
	}
	return code, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateRequest) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE demo_accounts
		SET
			assigned_to = COALESCE(NULLIF($2, ''), assigned_to),
			assignment_note = COALESCE(NULLIF($3, ''), assignment_note),
			expires_at = CASE WHEN $4 THEN $5 ELSE expires_at END,
			disabled_at = CASE
				WHEN $6::boolean IS NULL THEN disabled_at
				WHEN $6 THEN COALESCE(disabled_at, NOW())
				ELSE NULL
			END,
			disabled_reason = CASE
				WHEN $6::boolean IS NULL THEN disabled_reason
				WHEN $6 THEN NULLIF($7, '')
				ELSE NULL
			END,
			updated_at = NOW()
		WHERE id = $1
	`, id, stringPtrValue(req.AssignedTo), stringPtrValue(req.AssignmentNote), req.ExpiresAt != nil, nullableTimePtr(req.ExpiresAt), req.Disabled, stringPtrValue(req.DisabledReason))
	return err
}

func (r *Repository) PurgePreview(ctx context.Context, id string) (PurgePreview, error) {
	var preview PurgePreview
	err := r.pool.QueryRow(ctx, `
		SELECT da.id, da.user_id, da.label
		FROM demo_accounts da
		WHERE da.id = $1
	`, id).Scan(&preview.DemoAccountID, &preview.UserID, &preview.Label)
	if err != nil {
		return PurgePreview{}, err
	}

	countQueries := []struct {
		target *int
		query  string
	}{
		{&preview.CourseDrafts, `SELECT COUNT(*) FROM course_drafts WHERE user_id = $1`},
		{&preview.Courses, `SELECT COUNT(*) FROM courses WHERE user_id = $1`},
		{&preview.GoalProfiles, `SELECT COUNT(*) FROM course_goal_profiles WHERE user_id = $1`},
		{&preview.GoalRevisionLogs, `
			SELECT COUNT(*)
			FROM course_goal_revision_logs gr
			WHERE gr.course_draft_id IN (SELECT id FROM course_drafts WHERE user_id = $1)
		`},
		{&preview.CoursePoints, `
			SELECT COUNT(*)
			FROM course_points cp
			JOIN courses c ON c.id = cp.course_id
			WHERE c.user_id = $1
		`},
		{&preview.CourseDraftPoints, `
			SELECT COUNT(*)
			FROM course_draft_points cdp
			JOIN course_drafts cd ON cd.id = cdp.course_draft_id
			WHERE cd.user_id = $1
		`},
		{&preview.Attachments, `
			SELECT COUNT(*)
			FROM course_point_attachments cpa
			WHERE cpa.user_id = $1
			   OR cpa.course_point_id IN (
			     SELECT cp.id
			     FROM course_points cp
			     JOIN courses c ON c.id = cp.course_id
			     WHERE c.user_id = $1
			   )
		`},
		{&preview.AIUsageEvents, `SELECT COUNT(*) FROM ai_usage_events WHERE user_id = $1`},
		{&preview.RecommendationEvents, `SELECT COUNT(*) FROM recommendation_events WHERE user_id = $1`},
		{&preview.RecommendationLearnerEvents, `SELECT COUNT(*) FROM recommendation_rollout_learner_events WHERE user_id = $1`},
		{&preview.Contents, `SELECT COUNT(*) FROM contents WHERE user_id = $1`},
		{&preview.SafeDeletableContents, safeDeletableContentsCountQuery},
	}
	for _, item := range countQueries {
		if err := r.pool.QueryRow(ctx, item.query, preview.UserID).Scan(item.target); err != nil {
			return PurgePreview{}, err
		}
	}

	objectKeys, err := r.collectObjectKeys(ctx, preview.UserID)
	if err != nil {
		return PurgePreview{}, err
	}
	preview.ObjectKeys = objectKeys
	preview.ObjectStorageFiles = len(objectKeys)
	return preview, nil
}

func (r *Repository) PurgeLearningData(ctx context.Context, id string) (PurgePreview, error) {
	preview, err := r.PurgePreview(ctx, id)
	if err != nil {
		return PurgePreview{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return PurgePreview{}, err
	}
	defer tx.Rollback(ctx)

	userID := preview.UserID
	deleteStatements := []string{
		`DELETE FROM course_goal_revision_logs WHERE course_draft_id IN (SELECT id FROM course_drafts WHERE user_id = $1)`,
		`DELETE FROM recommendation_rollout_learner_events WHERE user_id = $1`,
		`DELETE FROM recommendation_events WHERE user_id = $1`,
		`DELETE FROM ai_usage_events WHERE user_id = $1`,
		`DELETE FROM course_point_attachments WHERE user_id = $1 OR course_point_id IN (
			SELECT cp.id FROM course_points cp JOIN courses c ON c.id = cp.course_id WHERE c.user_id = $1
		)`,
		`DELETE FROM courses WHERE user_id = $1`,
		`DELETE FROM course_drafts WHERE user_id = $1`,
		`DELETE FROM contents c WHERE c.user_id = $1 AND ` + safeDeletableContentsWhereClause,
	}
	for _, stmt := range deleteStatements {
		if _, err := tx.Exec(ctx, stmt, userID); err != nil {
			return PurgePreview{}, err
		}
	}

	now := time.Now()
	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_wallets(id, user_id, free_balance, paid_balance, updated_at)
		VALUES(gen_random_uuid(), $1, $2, 0, $3)
		ON CONFLICT(user_id) DO UPDATE
		SET free_balance = EXCLUDED.free_balance,
		    paid_balance = 0,
		    updated_at = EXCLUDED.updated_at
	`, userID, DefaultPoints, now)
	if err != nil {
		return PurgePreview{}, err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO ai_point_transactions(id, user_id, type, amount, created_at)
		VALUES(gen_random_uuid(), $1, 'grant_free', $2, $3)
	`, userID, DefaultPoints, now)
	if err != nil {
		return PurgePreview{}, err
	}
	_, err = tx.Exec(ctx, `
		UPDATE demo_accounts
		SET updated_at = $2
		WHERE id = $1
	`, id, now)
	if err != nil {
		return PurgePreview{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return PurgePreview{}, err
	}
	return preview, nil
}

func (r *Repository) collectObjectKeys(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT file_path
		FROM course_point_attachments cpa
		WHERE COALESCE(file_path, '') <> ''
		  AND (
			cpa.user_id = $1
			OR cpa.course_point_id IN (
				SELECT cp.id
				FROM course_points cp
				JOIN courses c ON c.id = cp.course_id
				WHERE c.user_id = $1
			)
		  )
		ORDER BY file_path
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func normalizeCreateRequest(req CreateBatchRequest) CreateBatchRequest {
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Count > 100 {
		req.Count = 100
	}
	if req.StartNumber <= 0 {
		req.StartNumber = 1
	}
	req.EmailPrefix = normalizeToken(req.EmailPrefix, DefaultEmailPrefix)
	req.EmailDomain = normalizeDomain(req.EmailDomain, DefaultEmailDomain)
	req.LabelPrefix = normalizeToken(req.LabelPrefix, DefaultLabelPrefix)
	req.UILocale = normalizeLocale(req.UILocale)
	req.LearningLanguage = normalizeLocale(req.LearningLanguage)
	if req.InitialPoints <= 0 {
		req.InitialPoints = DefaultPoints
	}
	if strings.TrimSpace(req.CreatedByActor) == "" {
		req.CreatedByActor = "super_admin"
	}
	req.AssignedTo = strings.TrimSpace(req.AssignedTo)
	req.AssignmentNote = strings.TrimSpace(req.AssignmentNote)
	return req
}

func normalizeListOptions(opts ListOptions) ListOptions {
	opts.Status = strings.TrimSpace(strings.ToLower(opts.Status))
	switch opts.Status {
	case "active", "disabled", "expired", "all":
	default:
		opts.Status = "all"
	}
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 100 {
		opts.Limit = 20
	}
	return opts
}

func formatNullableTime(value sql.NullTime) *string {
	if !value.Valid {
		return nil
	}
	text := formatTime(value.Time)
	return &text
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func nullableTimePtr(value **time.Time) *time.Time {
	if value == nil {
		return nil
	}
	return *value
}

const safeDeletableContentsWhereClause = `
	NOT EXISTS (SELECT 1 FROM course_draft_points x WHERE x.content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM course_points x WHERE x.content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM recommendation_events x WHERE x.content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM explorer_nodes x WHERE x.content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM content_recommendation_blocks x WHERE x.content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM course_point_material_reports x WHERE x.target_content_id = c.id OR x.replacement_content_id = c.id)
	AND NOT EXISTS (SELECT 1 FROM recommendation_rollout_learner_events x WHERE x.selected_content_id = c.id)
`

const safeDeletableContentsCountQuery = `
	SELECT COUNT(*)
	FROM contents c
	WHERE c.user_id = $1
	  AND ` + safeDeletableContentsWhereClause

func normalizeLocale(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func normalizeToken(value string, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "")
	if value == "" {
		return fallback
	}
	return value
}

func normalizeDomain(value string, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || !strings.Contains(value, ".") {
		return fallback
	}
	return value
}
