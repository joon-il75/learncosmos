package admin

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/auth"
)

type platformNoticeRow struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Locale      string  `json:"locale"`
	Title       string  `json:"title"`
	Summary     string  `json:"summary"`
	Body        string  `json:"body"`
	Status      string  `json:"status"`
	Pinned      bool    `json:"pinned"`
	PublishedAt *string `json:"published_at,omitempty"`
	CreatedBy   string  `json:"created_by"`
	UpdatedBy   string  `json:"updated_by"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type platformNoticeRequest struct {
	Slug    string `json:"slug"`
	Locale  string `json:"locale"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
	Status  string `json:"status"`
	Pinned  bool   `json:"pinned"`
}

var platformNoticeSlugPattern = regexp.MustCompile(`[^a-z0-9-]+`)
var platformNoticeStatuses = map[string]bool{"draft": true, "published": true, "archived": true}
var platformNoticeLocales = map[string]bool{"ko": true, "en": true}

func normalizePlatformNoticeLocale(raw string) string {
	locale := strings.ToLower(strings.TrimSpace(raw))
	if platformNoticeLocales[locale] {
		return locale
	}
	return "ko"
}

func normalizePlatformNoticeStatus(raw string) string {
	status := strings.ToLower(strings.TrimSpace(raw))
	if platformNoticeStatuses[status] {
		return status
	}
	return "draft"
}

func sanitizePlatformNoticeSlug(raw string) string {
	slug := strings.ToLower(strings.TrimSpace(raw))
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = platformNoticeSlugPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return fmt.Sprintf("notice-%d", time.Now().Unix())
	}
	if len(slug) > 96 {
		slug = strings.Trim(slug[:96], "-")
	}
	if slug == "" {
		return fmt.Sprintf("notice-%d", time.Now().Unix())
	}
	return slug
}

func validatePlatformNoticeRequest(req *platformNoticeRequest) (platformNoticeRequest, string) {
	normalized := platformNoticeRequest{
		Slug:    sanitizePlatformNoticeSlug(req.Slug),
		Locale:  normalizePlatformNoticeLocale(req.Locale),
		Title:   strings.TrimSpace(req.Title),
		Summary: strings.TrimSpace(req.Summary),
		Body:    strings.TrimSpace(req.Body),
		Status:  normalizePlatformNoticeStatus(req.Status),
		Pinned:  req.Pinned,
	}
	if normalized.Title == "" || len([]rune(normalized.Title)) > 160 {
		return normalized, "invalid title"
	}
	if len([]rune(normalized.Summary)) > 500 {
		return normalized, "summary too long"
	}
	if len([]rune(normalized.Body)) > 20000 {
		return normalized, "body too long"
	}
	return normalized, ""
}

func scanPlatformNoticeRow(scanner interface{ Scan(dest ...any) error }) (platformNoticeRow, error) {
	var row platformNoticeRow
	var publishedAt *time.Time
	var createdAt time.Time
	var updatedAt time.Time
	if err := scanner.Scan(
		&row.ID,
		&row.Slug,
		&row.Locale,
		&row.Title,
		&row.Summary,
		&row.Body,
		&row.Status,
		&row.Pinned,
		&publishedAt,
		&row.CreatedBy,
		&row.UpdatedBy,
		&createdAt,
		&updatedAt,
	); err != nil {
		return row, err
	}
	if publishedAt != nil {
		formatted := publishedAt.UTC().Format(time.RFC3339)
		row.PublishedAt = &formatted
	}
	row.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	row.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
	return row, nil
}

const platformNoticeSelectColumns = `
	id::text,
	slug,
	locale,
	title,
	COALESCE(summary, '') AS summary,
	COALESCE(body, '') AS body,
	status,
	pinned,
	published_at,
	COALESCE(created_by, '') AS created_by,
	COALESCE(updated_by, '') AS updated_by,
	created_at,
	updated_at
`

// GET /api/v1/public/platform-notices
func (h *AdminHandler) publicListPlatformNoticesHandler(c *gin.Context) {
	locale := normalizePlatformNoticeLocale(c.Query("locale"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}
	ctx := c.Request.Context()

	var total int
	if err := h.db.QueryRow(ctx, `SELECT COUNT(*) FROM platform_notices WHERE status = 'published' AND locale = $1`, locale).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT `+platformNoticeSelectColumns+`
		FROM platform_notices
		WHERE status = 'published' AND locale = $1
		ORDER BY pinned DESC, COALESCE(published_at, created_at) DESC, created_at DESC
		LIMIT $2
	`, locale, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}
	defer rows.Close()

	notices := []platformNoticeRow{}
	for rows.Next() {
		row, err := scanPlatformNoticeRow(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
			return
		}
		notices = append(notices, row)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notices": notices, "total": total, "limit": limit, "locale": locale})
}

// GET /api/v1/public/platform-notices/:slug
func (h *AdminHandler) publicGetPlatformNoticeHandler(c *gin.Context) {
	slug := sanitizePlatformNoticeSlug(c.Param("slug"))
	ctx := c.Request.Context()
	row, err := scanPlatformNoticeRow(h.db.QueryRow(ctx, `
		SELECT `+platformNoticeSelectColumns+`
		FROM platform_notices
		WHERE status = 'published' AND slug = $1
	`, slug))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "notice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notice": row})
}

// GET /api/v1/super-admin/platform-notices
func (h *AdminHandler) listPlatformNoticesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit
	locale := strings.ToLower(strings.TrimSpace(c.Query("locale")))
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))

	where := []string{"1=1"}
	args := []any{}
	if platformNoticeLocales[locale] {
		args = append(args, locale)
		where = append(where, fmt.Sprintf("locale = $%d", len(args)))
	}
	if platformNoticeStatuses[status] {
		args = append(args, status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	whereSQL := strings.Join(where, " AND ")
	ctx := c.Request.Context()

	var total int
	if err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM platform_notices WHERE "+whereSQL, args...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, limit, offset)
	limitPos := len(queryArgs) - 1
	offsetPos := len(queryArgs)
	rows, err := h.db.Query(ctx, fmt.Sprintf(`
		SELECT %s
		FROM platform_notices
		WHERE %s
		ORDER BY pinned DESC, updated_at DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, platformNoticeSelectColumns, whereSQL, limitPos, offsetPos), queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}
	defer rows.Close()

	notices := []platformNoticeRow{}
	for rows.Next() {
		row, err := scanPlatformNoticeRow(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
			return
		}
		notices = append(notices, row)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice list unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notices": notices, "total": total, "page": page, "limit": limit})
}

// POST /api/v1/super-admin/platform-notices
func (h *AdminHandler) createPlatformNoticeHandler(c *gin.Context) {
	var req platformNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	normalized, validationErr := validatePlatformNoticeRequest(&req)
	if validationErr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		return
	}
	adminID, _ := auth.GetCurrentUserID(c)
	row, err := scanPlatformNoticeRow(h.db.QueryRow(c.Request.Context(), `
		INSERT INTO platform_notices(slug, locale, title, summary, body, status, pinned, published_at, created_by, updated_by)
		VALUES($1, $2, $3, $4, $5, $6, $7, CASE WHEN $6 = 'published' THEN NOW() ELSE NULL END, $8, $8)
		RETURNING `+platformNoticeSelectColumns+`
	`, normalized.Slug, normalized.Locale, normalized.Title, normalized.Summary, normalized.Body, normalized.Status, normalized.Pinned, adminID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice create failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"notice": row})
}

// PATCH /api/v1/super-admin/platform-notices/:id
func (h *AdminHandler) updatePlatformNoticeHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notice id"})
		return
	}
	var req platformNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	normalized, validationErr := validatePlatformNoticeRequest(&req)
	if validationErr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		return
	}
	adminID, _ := auth.GetCurrentUserID(c)
	row, err := scanPlatformNoticeRow(h.db.QueryRow(c.Request.Context(), `
		UPDATE platform_notices
		SET slug = $1,
		    locale = $2,
		    title = $3,
		    summary = $4,
		    body = $5,
		    status = $6,
		    pinned = $7,
		    published_at = CASE
		      WHEN $6 = 'published' AND published_at IS NULL THEN NOW()
		      WHEN $6 <> 'published' THEN NULL
		      ELSE published_at
		    END,
		    updated_by = $8,
		    updated_at = NOW()
		WHERE id = $9
		RETURNING `+platformNoticeSelectColumns+`
	`, normalized.Slug, normalized.Locale, normalized.Title, normalized.Summary, normalized.Body, normalized.Status, normalized.Pinned, adminID, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "notice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notice": row})
}

// DELETE /api/v1/super-admin/platform-notices/:id
func (h *AdminHandler) deletePlatformNoticeHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notice id"})
		return
	}
	ct, err := h.db.Exec(c.Request.Context(), `DELETE FROM platform_notices WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notice delete failed"})
		return
	}
	if ct.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "notice not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
