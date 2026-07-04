package content

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
	"github.com/learnweaver/backend/internal/pkg/urlsafe"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// scanContent — 공통 컬럼 스캔 헬퍼
func scanContent(row interface {
	Scan(dest ...any) error
}) (*Content, error) {
	var c Content
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.ContentType,
		&c.ExternalSource,
		&c.ExternalContentID,
		&c.CanonicalURL,
		&c.URL,
		&c.Title,
		&c.Description,
		&c.ThumbnailURL,
		&c.DurationSeconds,
		&c.Author,
		&c.Language,
		&c.IsPublic,
		&c.QualityScore,
		&c.ContentStatus,
		&c.HealthScore,
		&c.LastCheckedAt,
		&c.HTTPStatus,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

const selectCols = `
	id, user_id, content_type, external_source, external_content_id, canonical_url, url, title, description,
	thumbnail_url, duration_seconds, author, language,
	is_public, quality_score, content_status, health_score, last_checked_at, http_status, created_at, updated_at`

// Create — 새 콘텐츠 등록
func (r *Repository) Create(ctx context.Context, userID uuid.UUID, req CreateContentRequest) (*Content, error) {
	contentURL, err := normalizeStoredContentURL(req.URL)
	if err != nil {
		return nil, err
	}
	req.URL = contentURL

	lang := req.Language
	if lang == "" {
		lang = "ko"
	}

	query := `
		INSERT INTO contents
			(user_id, content_type, external_source, external_content_id, canonical_url, url, title, description, thumbnail_url,
			 duration_seconds, author, language, is_public, search_text_ko)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false, $13)
		RETURNING` + selectCols

	canonicalURL := buildCanonicalURL(req.URL)
	searchTextKO := normalizer.BuildSearchTextKO(
		req.Title,
		valueOrEmpty(req.Description),
		valueOrEmpty(req.Author),
		string(req.ContentType),
		lang,
	)
	row := r.pool.QueryRow(ctx, query,
		userID,
		req.ContentType,
		normalizeOptionalString(req.ExternalSource),
		normalizeOptionalString(req.ExternalContentID),
		canonicalURL,
		req.URL,
		req.Title,
		req.Description,
		req.ThumbnailURL,
		req.DurationSeconds,
		req.Author,
		lang,
		searchTextKO,
	)
	c, err := scanContent(row)
	if err != nil {
		return nil, fmt.Errorf("content create: %w", err)
	}
	return c, nil
}

// GetByID — 단건 조회 (소유자 확인 포함)
func (r *Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (*Content, error) {
	query := `SELECT` + selectCols + `
		FROM contents
		WHERE id = $1 AND user_id = $2`

	row := r.pool.QueryRow(ctx, query, id, userID)
	c, err := scanContent(row)
	if err != nil {
		return nil, fmt.Errorf("content get: %w", err)
	}
	return c, nil
}

// List — 목록 조회 (본인 콘텐츠만 · 페이지네이션)
func (r *Repository) List(ctx context.Context, userID uuid.UUID, q ListContentsQuery) (*ListContentsResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	offset := (q.Page - 1) * q.Limit

	countQuery := `SELECT COUNT(*) FROM contents WHERE user_id = $1`
	args := []any{userID}
	argIdx := 2

	if q.ContentType != "" {
		countQuery += fmt.Sprintf(" AND content_type = $%d", argIdx)
		args = append(args, q.ContentType)
		argIdx++
	}

	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("content count: %w", err)
	}

	listQuery := `SELECT` + selectCols + `
		FROM contents
		WHERE user_id = $1`
	listArgs := []any{userID}
	listArgIdx := 2

	if q.ContentType != "" {
		listQuery += fmt.Sprintf(" AND content_type = $%d", listArgIdx)
		listArgs = append(listArgs, q.ContentType)
		listArgIdx++
	}

	listQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", listArgIdx, listArgIdx+1)
	listArgs = append(listArgs, q.Limit, offset)

	rows, err := r.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, fmt.Errorf("content list: %w", err)
	}
	defer rows.Close()

	contents := []Content{}
	for rows.Next() {
		c, err := scanContent(rows)
		if err != nil {
			return nil, fmt.Errorf("content scan: %w", err)
		}
		contents = append(contents, *c)
	}

	return &ListContentsResponse{
		Contents: contents,
		Total:    total,
		Page:     q.Page,
		Limit:    q.Limit,
	}, nil
}

// Update — 부분 수정 (소유자만)
func (r *Repository) Update(ctx context.Context, id, userID uuid.UUID, req UpdateContentRequest) (*Content, error) {
	existing, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	title := existing.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := existing.Description
	if req.Description != nil {
		description = req.Description
	}
	thumbnailURL := existing.ThumbnailURL
	if req.ThumbnailURL != nil {
		thumbnailURL = req.ThumbnailURL
	}
	author := existing.Author
	if req.Author != nil {
		author = req.Author
	}
	contentURL := existing.URL
	if req.URL != nil {
		normalizedURL, err := normalizeStoredContentURL(req.URL)
		if err != nil {
			return nil, err
		}
		contentURL = normalizedURL
	}
	language := existing.Language
	if req.Language != nil {
		language = *req.Language
	}

	query := `
		UPDATE contents
		SET title = $1, description = $2, thumbnail_url = $3,
		    author = $4, language = $5, search_text_ko = $6, canonical_url = $7, url = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING` + selectCols

	canonicalURL := buildCanonicalURL(contentURL)
	searchTextKO := normalizer.BuildSearchTextKO(
		title,
		valueOrEmpty(description),
		valueOrEmpty(author),
		string(existing.ContentType),
		language,
	)
	row := r.pool.QueryRow(ctx, query,
		title, description, thumbnailURL, author, language, searchTextKO, canonicalURL, contentURL, id, userID,
	)
	c, err := scanContent(row)
	if err != nil {
		return nil, fmt.Errorf("content update: %w", err)
	}
	return c, nil
}

func (r *Repository) SaveEmbeddingReady(ctx context.Context, id uuid.UUID, provider, model, sourceText string, vector []float32) error {
	if len(vector) == 0 {
		return fmt.Errorf("embedding vector is empty")
	}
	provider = normalizeEmbeddingProvider(provider)
	model = normalizeEmbeddingModel(model)
	pgvecStr := float32VectorToPGVector(vector)
	sourceHash := sourceTextHash(sourceText)
	preview := truncateString(strings.TrimSpace(sourceText), 240)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO content_embeddings (
			content_id, provider, model, dimension, embedding,
			source_text_hash, source_text_preview, status, error_message, embedded_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5::vector, $6, $7, 'ready', '', NOW(), NOW(), NOW())
		ON CONFLICT (content_id, provider, model, dimension)
		DO UPDATE SET
			embedding = EXCLUDED.embedding,
			source_text_hash = EXCLUDED.source_text_hash,
			source_text_preview = EXCLUDED.source_text_preview,
			status = 'ready',
			error_message = '',
			embedded_at = NOW(),
			updated_at = NOW()
	`, id, provider, model, len(vector), pgvecStr, sourceHash, preview)
	return err
}

func (r *Repository) SaveEmbeddingFailure(ctx context.Context, id uuid.UUID, provider, model string, dimension int, sourceText string, embedErr error) error {
	provider = normalizeEmbeddingProvider(provider)
	model = normalizeEmbeddingModel(model)
	if dimension <= 0 {
		dimension = 768
	}
	errMessage := ""
	if embedErr != nil {
		errMessage = logsafe.PersistedError(embedErr.Error())
	}
	sourceHash := sourceTextHash(sourceText)
	preview := truncateString(strings.TrimSpace(sourceText), 240)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO content_embeddings (
			content_id, provider, model, dimension,
			source_text_hash, source_text_preview, status, error_message,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'failed', $7, NOW(), NOW())
		ON CONFLICT (content_id, provider, model, dimension)
		DO UPDATE SET
			source_text_hash = EXCLUDED.source_text_hash,
			source_text_preview = EXCLUDED.source_text_preview,
			status = 'failed',
			error_message = EXCLUDED.error_message,
			updated_at = NOW()
	`, id, provider, model, dimension, sourceHash, preview, errMessage)
	return err
}

func float32VectorToPGVector(vector []float32) string {
	parts := make([]string, len(vector))
	for i, v := range vector {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func sourceTextHash(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}

func normalizeEmbeddingProvider(provider string) string {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return "embedding_gemma"
	}
	return provider
}

func normalizeEmbeddingModel(model string) string {
	model = strings.TrimSpace(model)
	if model == "" {
		return "embedding-gemma"
	}
	return model
}

func truncateString(value string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	return string(runes[:maxRunes])
}

func (r *Repository) FindByExternalSource(ctx context.Context, externalSource, externalContentID string) (*Content, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+selectCols+`
		FROM contents
		WHERE external_source = $1 AND external_content_id = $2
		LIMIT 1
	`, strings.TrimSpace(externalSource), strings.TrimSpace(externalContentID))

	content, err := scanContent(row)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *Repository) FindByCanonicalURL(ctx context.Context, canonicalURL string) (*Content, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+selectCols+`
		FROM contents
		WHERE canonical_url = $1
		LIMIT 1
	`, strings.TrimSpace(canonicalURL))

	content, err := scanContent(row)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *Repository) UpsertExternalContent(ctx context.Context, userID uuid.UUID, req CreateContentRequest) (*Content, error) {
	contentURL, err := normalizeStoredContentURL(req.URL)
	if err != nil {
		return nil, err
	}
	req.URL = contentURL
	canonicalURL := buildCanonicalURL(req.URL)
	if externalSource := normalizeOptionalString(req.ExternalSource); externalSource != nil {
		if externalContentID := normalizeOptionalString(req.ExternalContentID); externalContentID != nil {
			existing, err := r.FindByExternalSource(ctx, *externalSource, *externalContentID)
			if err == nil && existing != nil {
				return existing, nil
			}
		}
	}
	if canonicalURL != nil {
		existing, err := r.FindByCanonicalURL(ctx, *canonicalURL)
		if err == nil && existing != nil {
			return existing, nil
		}
	}
	return r.Create(ctx, userID, req)
}

type HealthCheckTarget struct {
	ID  uuid.UUID
	URL string
}

func (r *Repository) ListHealthCheckTargets(ctx context.Context, limit int) ([]HealthCheckTarget, error) {
	if limit <= 0 {
		limit = 200
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, url
		FROM contents
		WHERE content_type <> 'internal'
		  AND url IS NOT NULL
		  AND (
		    last_checked_at IS NULL
		    OR last_checked_at < NOW() - interval '24 hours'
		  )
		ORDER BY last_checked_at ASC NULLS FIRST, created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list health check targets: %w", err)
	}
	defer rows.Close()

	targets := make([]HealthCheckTarget, 0, limit)
	for rows.Next() {
		var target HealthCheckTarget
		if err := rows.Scan(&target.ID, &target.URL); err != nil {
			return nil, fmt.Errorf("scan health check target: %w", err)
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (r *Repository) UpdateContentHealth(ctx context.Context, id uuid.UUID, status ContentStatus, healthScore float64, httpStatus *int, checkedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE contents
		SET content_status = $1,
		    health_score = $2,
		    http_status = $3,
		    last_checked_at = $4,
		    updated_at = NOW()
		WHERE id = $5
	`, status, healthScore, httpStatus, checkedAt, id)
	if err != nil {
		return fmt.Errorf("update content health: %w", err)
	}
	return nil
}

// Delete — 삭제 (소유자만)
func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM contents WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("content delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func buildCanonicalURL(rawURL *string) *string {
	if rawURL == nil {
		return nil
	}
	normalizedURL, ok := urlsafe.NormalizeHTTPURL(*rawURL)
	if !ok {
		return nil
	}

	parsed, err := url.Parse(normalizedURL)
	if err != nil || parsed == nil {
		return nil
	}
	parsed.Fragment = ""
	normalized := strings.TrimRight(parsed.String(), "/")
	if normalized == "" {
		return nil
	}
	return &normalized
}
