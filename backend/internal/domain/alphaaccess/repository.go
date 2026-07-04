package alphaaccess

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCodeNotFound = errors.New("alpha invite code not found")
	ErrCodeRevoked  = errors.New("alpha invite code revoked")
	ErrCodeExpired  = errors.New("alpha invite code expired")
	ErrCodeUsedUp   = errors.New("alpha invite code used up")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func NormalizeCode(code string) string {
	code = strings.TrimSpace(strings.ToUpper(code))
	code = strings.ReplaceAll(code, " ", "")
	return code
}

func GenerateCode() (string, error) {
	var raw [10]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	encoded := strings.TrimRight(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw[:]), "=")
	if len(encoded) < 12 {
		return "", errors.New("generated invite code too short")
	}
	return "LW-" + encoded[:4] + "-" + encoded[4:8] + "-" + encoded[8:12], nil
}

func (r *Repository) HasAlphaAccess(ctx context.Context, userID string) (bool, *time.Time, error) {
	var grantedAt *time.Time
	err := r.db.QueryRow(ctx, `
		SELECT alpha_access_granted_at
		FROM users
		WHERE id = $1
	`, userID).Scan(&grantedAt)
	if err != nil {
		return false, nil, err
	}
	return grantedAt != nil, grantedAt, nil
}

func (r *Repository) Redeem(ctx context.Context, userID, code string) (*time.Time, error) {
	code = NormalizeCode(code)
	if code == "" {
		return nil, ErrCodeNotFound
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var existingGrantedAt *time.Time
	if err := tx.QueryRow(ctx, `
		SELECT alpha_access_granted_at
		FROM users
		WHERE id = $1
	`, userID).Scan(&existingGrantedAt); err != nil {
		return nil, err
	}
	if existingGrantedAt != nil {
		return existingGrantedAt, tx.Commit(ctx)
	}

	var codeID string
	var status string
	var maxUses int
	var usedCount int
	var expiresAt time.Time
	err = tx.QueryRow(ctx, `
		SELECT id, status, max_uses, used_count, expires_at
		FROM alpha_invite_codes
		WHERE code = $1
		FOR UPDATE
	`, code).Scan(&codeID, &status, &maxUses, &usedCount, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCodeNotFound
		}
		return nil, err
	}
	if status != "active" {
		return nil, ErrCodeRevoked
	}
	if !expiresAt.After(time.Now()) {
		return nil, ErrCodeExpired
	}
	if usedCount >= maxUses {
		return nil, ErrCodeUsedUp
	}

	var provider sql.NullString
	var providerID sql.NullString
	err = tx.QueryRow(ctx, `
		SELECT provider, provider_id
		FROM social_accounts
		WHERE user_id = $1
		ORDER BY created_at ASC
		LIMIT 1
	`, userID).Scan(&provider, &providerID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	grantedAt := time.Now().UTC()
	_, err = tx.Exec(ctx, `
		INSERT INTO alpha_invite_code_uses(id, code_id, user_id, provider, provider_id, used_at)
		VALUES($1, $2, $3, $4, $5, $6)
	`, uuid.New().String(), codeID, userID, nullableString(provider), nullableString(providerID), grantedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE alpha_invite_codes
		SET used_count = used_count + 1,
		    updated_at = now()
		WHERE id = $1
	`, codeID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE users
		SET alpha_access_granted_at = $2
		WHERE id = $1
		  AND alpha_access_granted_at IS NULL
	`, userID, grantedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &grantedAt, nil
}

func (r *Repository) CreateInviteCode(ctx context.Context, createdBy *string, expiresAt time.Time, maxUses int, sentToNote, adminNote string) (InviteCodeResponse, error) {
	if maxUses < 1 {
		maxUses = 1
	}
	sentToNote = strings.TrimSpace(sentToNote)
	adminNote = strings.TrimSpace(adminNote)

	for attempt := 0; attempt < 5; attempt++ {
		code, err := GenerateCode()
		if err != nil {
			return InviteCodeResponse{}, err
		}

		var item InviteCodeResponse
		var createdAt time.Time
		var updatedAt time.Time
		err = r.db.QueryRow(ctx, `
			INSERT INTO alpha_invite_codes(code, max_uses, expires_at, sent_to_note, admin_note, created_by)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id, code, status, max_uses, used_count, expires_at, sent_to_note, admin_note, created_by, created_at, updated_at
		`, code, maxUses, expiresAt, sentToNote, adminNote, createdBy).Scan(
			&item.ID, &item.Code, &item.Status, &item.MaxUses, &item.UsedCount, &expiresAt,
			&item.SentToNote, &item.AdminNote, &item.CreatedBy, &createdAt, &updatedAt,
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			return InviteCodeResponse{}, err
		}
		item.ExpiresAt = formatTime(expiresAt)
		item.CreatedAt = formatTime(createdAt)
		item.UpdatedAt = formatTime(updatedAt)
		item.State = deriveState(item.Status, item.UsedCount, item.MaxUses, expiresAt)
		item.Uses = []InviteCodeUseRecord{}
		return item, nil
	}
	return InviteCodeResponse{}, errors.New("failed to generate unique invite code")
}

func (r *Repository) ListInviteCodes(ctx context.Context) ([]InviteCodeResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code, status, max_uses, used_count, expires_at, sent_to_note, admin_note, created_by, created_at, updated_at
		FROM alpha_invite_codes
		ORDER BY created_at DESC
		LIMIT 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []InviteCodeResponse{}
	for rows.Next() {
		var item InviteCodeResponse
		var expiresAt time.Time
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Status, &item.MaxUses, &item.UsedCount, &expiresAt,
			&item.SentToNote, &item.AdminNote, &item.CreatedBy, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		item.ExpiresAt = formatTime(expiresAt)
		item.CreatedAt = formatTime(createdAt)
		item.UpdatedAt = formatTime(updatedAt)
		item.State = deriveState(item.Status, item.UsedCount, item.MaxUses, expiresAt)
		item.Uses = []InviteCodeUseRecord{}
		items = append(items, item)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(items) == 0 {
		return items, nil
	}

	ids := make([]uuid.UUID, 0, len(items))
	indexByID := make(map[string]int, len(items))
	for i, item := range items {
		parsedID, err := uuid.Parse(item.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, parsedID)
		indexByID[item.ID] = i
	}

	useRows, err := r.db.Query(ctx, `
		SELECT
			icu.id,
			icu.code_id,
			icu.user_id,
			COALESCE(u.email, '') AS email,
			COALESCE(u.nickname, '') AS nickname,
			u.display_id,
			COALESCE(icu.provider, '') AS provider,
			COALESCE(icu.provider_id, '') AS provider_id,
			icu.used_at
		FROM alpha_invite_code_uses icu
		JOIN users u ON u.id = icu.user_id
		WHERE icu.code_id = ANY($1::uuid[])
		ORDER BY icu.used_at DESC
	`, ids)
	if err != nil {
		return nil, err
	}
	defer useRows.Close()

	for useRows.Next() {
		var codeID string
		var usedAt time.Time
		var rec InviteCodeUseRecord
		if err := useRows.Scan(&rec.ID, &codeID, &rec.UserID, &rec.Email, &rec.Nickname, &rec.DisplayID, &rec.Provider, &rec.ProviderID, &usedAt); err != nil {
			return nil, err
		}
		rec.UsedAt = formatTime(usedAt)
		if idx, ok := indexByID[codeID]; ok {
			items[idx].Uses = append(items[idx].Uses, rec)
		}
	}
	if useRows.Err() != nil {
		return nil, useRows.Err()
	}
	return items, nil
}

func (r *Repository) RevokeInviteCode(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE alpha_invite_codes
		SET status = 'revoked',
		    updated_at = now()
		WHERE id = $1
	`, id)
	return err
}

func deriveState(status string, usedCount, maxUses int, expiresAt time.Time) string {
	if status == "revoked" {
		return "revoked"
	}
	if !expiresAt.After(time.Now()) {
		return "expired"
	}
	if usedCount >= maxUses {
		return "used_up"
	}
	return "available"
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func nullableString(v sql.NullString) any {
	if !v.Valid {
		return nil
	}
	return v.String
}
