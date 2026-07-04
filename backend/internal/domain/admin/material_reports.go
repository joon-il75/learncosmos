package admin

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/auth"
)

type materialReportRow struct {
	ID                    string  `json:"id"`
	CoursePointID         string  `json:"course_point_id"`
	UserID                string  `json:"user_id"`
	ReporterDisplay       string  `json:"reporter_display"`
	ReporterEmail         string  `json:"reporter_email"`
	AttachmentID          *string `json:"attachment_id,omitempty"`
	TargetType            string  `json:"target_type"`
	ReportType            string  `json:"report_type"`
	Message               string  `json:"message"`
	Status                string  `json:"status"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	AdminNote             string  `json:"admin_note"`
	ReviewedBy            *string `json:"reviewed_by,omitempty"`
	ReviewerDisplay       string  `json:"reviewer_display"`
	ReviewedAt            *string `json:"reviewed_at,omitempty"`
	CourseID              string  `json:"course_id"`
	CourseTitle           string  `json:"course_title"`
	PointContentID        *string `json:"point_content_id,omitempty"`
	LevelTitle            string  `json:"level_title"`
	LessonTitle           string  `json:"lesson_title"`
	PointTitle            string  `json:"point_title"`
	PointType             string  `json:"point_type"`
	PointExternalURL      string  `json:"point_external_url"`
	AttachmentTitle       string  `json:"attachment_title"`
	AttachmentURL         string  `json:"attachment_url"`
	AttachmentFile        string  `json:"attachment_file_path"`
	AttachmentSize        *int64  `json:"attachment_file_size,omitempty"`
	AttachmentMime        string  `json:"attachment_mime_type"`
	AttachmentContext     string  `json:"attachment_source_context"`
	AttachmentType        string  `json:"attachment_type"`
	RecommendationBlocked bool    `json:"recommendation_blocked"`
}

type updateMaterialReportRequest struct {
	Status    string `json:"status"`
	AdminNote string `json:"admin_note"`
}

var materialReportStatuses = map[string]bool{
	"open":      true,
	"reviewing": true,
	"resolved":  true,
	"dismissed": true,
	"cancelled": true,
}

var materialReportTargetTypes = map[string]bool{
	"source":     true,
	"ai_summary": true,
	"attachment": true,
	"other":      true,
}

var materialReportTypes = map[string]bool{
	"broken_link":    true,
	"wrong_content":  true,
	"unsafe_content": true,
	"copyright":      true,
	"low_quality":    true,
	"other":          true,
}

// GET /api/v1/super-admin/material-reports
func (h *AdminHandler) listMaterialReportsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 5
	}
	offset := (page - 1) * limit

	where := []string{"1=1"}
	args := []any{}
	addFilter := func(condition string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(condition, len(args)))
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" && status != "all" {
		if !materialReportStatuses[status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		addFilter("cpmr.status = $%d", status)
	}
	if targetType := strings.TrimSpace(c.Query("target_type")); targetType != "" && targetType != "all" {
		if !materialReportTargetTypes[targetType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_type"})
			return
		}
		addFilter("cpmr.target_type = $%d", targetType)
	}
	if reportType := strings.TrimSpace(c.Query("report_type")); reportType != "" && reportType != "all" {
		if !materialReportTypes[reportType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_type"})
			return
		}
		addFilter("cpmr.report_type = $%d", reportType)
	}
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		args = append(args, "%"+q+"%")
		placeholder := fmt.Sprintf("$%d", len(args))
		where = append(where, fmt.Sprintf(`(
			cpmr.message ILIKE %[1]s
			OR cp.title ILIKE %[1]s
			OR c.title ILIKE %[1]s
			OR COALESCE(u.email, '') ILIKE %[1]s
			OR COALESCE(u.nickname, '') ILIKE %[1]s
			OR COALESCE(cpa.title, '') ILIKE %[1]s
		)`, placeholder))
	}

	whereSQL := strings.Join(where, " AND ")
	ctx := c.Request.Context()

	var total int
	countSQL := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM course_point_material_reports cpmr
		JOIN course_points cp ON cp.id = cpmr.course_point_id
		JOIN courses c ON c.id = cp.course_id
		LEFT JOIN users u ON u.id = cpmr.user_id
		LEFT JOIN course_point_attachments cpa ON cpa.id = cpmr.attachment_id
		WHERE %s
	`, whereSQL)
	if err := h.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
		return
	}

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, limit, offset)
	limitPos := len(queryArgs) - 1
	offsetPos := len(queryArgs)
	rowsSQL := fmt.Sprintf(`
		SELECT
			cpmr.id::text,
			cpmr.course_point_id::text,
			cpmr.user_id::text,
			COALESCE(NULLIF(u.nickname, ''), NULLIF(u.email, ''), cpmr.user_id::text) AS reporter_display,
			COALESCE(u.email, '') AS reporter_email,
			cpmr.attachment_id::text,
			cpmr.target_type,
			cpmr.report_type,
			cpmr.message,
			cpmr.status,
			cpmr.created_at,
			cpmr.updated_at,
			COALESCE(cpmr.admin_note, '') AS admin_note,
			cpmr.reviewed_by::text,
			COALESCE(cpmr.reviewed_by, '') AS reviewer_display,
			cpmr.reviewed_at,
			c.id::text AS course_id,
			c.title AS course_title,
			cp.content_id::text AS point_content_id,
			COALESCE(parent_cl.title, '') AS level_title,
			COALESCE(cl.title, '') AS lesson_title,
			cp.title AS point_title,
			cp.point_type,
			COALESCE(cp.external_url, '') AS point_external_url,
			COALESCE(cpa.title, '') AS attachment_title,
			COALESCE(cpa.url, '') AS attachment_url,
			COALESCE(cpa.file_path, '') AS attachment_file_path,
			cpa.file_size,
			COALESCE(cpa.mime_type, '') AS attachment_mime_type,
			COALESCE(cpa.source_context, '') AS attachment_source_context,
			COALESCE(cpa.attachment_type, '') AS attachment_type,
			EXISTS (
				SELECT 1
				FROM content_recommendation_blocks crb
				WHERE crb.active = true
				  AND (
				    (cp.content_id IS NOT NULL AND crb.content_id = cp.content_id)
				    OR (COALESCE(crb.url_key, '') <> '' AND crb.url_key IN (
				      lower(trim(COALESCE(cp.external_url, ''))),
				      lower(trim(COALESCE(cpa.url, '')))
				    ))
				  )
			) AS recommendation_blocked
		FROM course_point_material_reports cpmr
		JOIN course_points cp ON cp.id = cpmr.course_point_id
		JOIN courses c ON c.id = cp.course_id
		LEFT JOIN course_lessons cl ON cl.id = cp.course_lesson_id
		LEFT JOIN course_lessons parent_cl ON parent_cl.id = cl.parent_lesson_id
		LEFT JOIN users u ON u.id = cpmr.user_id
		LEFT JOIN course_point_attachments cpa ON cpa.id = cpmr.attachment_id
		WHERE %s
		ORDER BY cpmr.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, limitPos, offsetPos)

	rows, err := h.db.Query(ctx, rowsSQL, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
		return
	}
	defer rows.Close()

	reports := []materialReportRow{}
	for rows.Next() {
		row, err := scanMaterialReportRow(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
			return
		}
		reports = append(reports, row)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// PATCH /api/v1/super-admin/material-reports/:id
func (h *AdminHandler) updateMaterialReportHandler(c *gin.Context) {
	reportID := strings.TrimSpace(c.Param("id"))
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	var req updateMaterialReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Status = strings.TrimSpace(req.Status)
	req.AdminNote = strings.TrimSpace(req.AdminNote)
	if !materialReportStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if len([]rune(req.AdminNote)) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin_note too long"})
		return
	}

	adminID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query := `
		UPDATE course_point_material_reports
		SET status = $1,
		    admin_note = $2,
		    reviewed_by = CASE WHEN $1 = 'open' THEN NULL ELSE $3 END,
		    reviewed_at = CASE WHEN $1 IN ('resolved', 'dismissed') THEN NOW() ELSE NULL END,
		    updated_at = NOW()
		WHERE id = $4
		RETURNING id
	`

	var updatedID string
	tx, err := h.db.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "처리 실패"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	if err := tx.QueryRow(c.Request.Context(), query, req.Status, req.AdminNote, adminID, reportID).Scan(&updatedID); err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "신고 내역을 찾을 수 없습니다."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "처리 실패"})
		return
	}

	if req.Status == "resolved" {
		if err := h.blockRecommendationForMaterialReport(c.Request.Context(), tx, reportID, adminID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "추천 제외 처리 실패"})
			return
		}
	} else {
		if err := h.unblockRecommendationForMaterialReport(c.Request.Context(), tx, reportID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "추천 제외 해제 실패"})
			return
		}
	}

	if err := tx.Commit(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "처리 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedID, "status": req.Status, "message": "updated"})
}

func (h *AdminHandler) blockRecommendationForMaterialReport(ctx context.Context, tx pgx.Tx, reportID string, adminID string) error {
	var reportType string
	var targetType string
	var contentID sql.NullString
	var pointURL string
	var attachmentURL string
	if err := tx.QueryRow(ctx, `
		SELECT cpmr.report_type,
		       cpmr.target_type,
		       COALESCE(cpmr.target_content_id::text, cp.content_id::text),
		       COALESCE(cpmr.target_url, cp.external_url, ''),
		       COALESCE(cpa.url, '')
		FROM course_point_material_reports cpmr
		JOIN course_points cp ON cp.id = cpmr.course_point_id
		LEFT JOIN course_point_attachments cpa ON cpa.id = cpmr.attachment_id
		WHERE cpmr.id = $1
	`, reportID).Scan(&reportType, &targetType, &contentID, &pointURL, &attachmentURL); err != nil {
		return err
	}

	if reportType != "broken_link" {
		return nil
	}

	targetURL := pointURL
	if targetType == "attachment" && strings.TrimSpace(attachmentURL) != "" {
		targetURL = attachmentURL
	}
	urlKey := materialReportRecommendationURLKey(targetURL)
	if !contentID.Valid && urlKey == "" {
		return nil
	}

	var contentIDValue any
	if contentID.Valid {
		contentIDValue = contentID.String
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO content_recommendation_blocks (
			content_id, url, url_key, source_report_id, reason, created_by, active
		) VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT DO NOTHING
	`, contentIDValue, strings.TrimSpace(targetURL), urlKey, reportID, "broken_link_report_resolved", adminID)
	return err
}

func (h *AdminHandler) unblockRecommendationForMaterialReport(ctx context.Context, tx pgx.Tx, reportID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE content_recommendation_blocks
		SET active = false,
		    updated_at = NOW()
		WHERE source_report_id = $1
		  AND active = true
	`, reportID)
	return err
}

type materialReportScanner interface {
	Scan(dest ...any) error
}

func scanMaterialReportRow(scanner materialReportScanner) (materialReportRow, error) {
	var row materialReportRow
	var attachmentID sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullTime
	var pointContentID sql.NullString
	var createdAt time.Time
	var updatedAt time.Time
	var attachmentSize sql.NullInt64

	if err := scanner.Scan(
		&row.ID,
		&row.CoursePointID,
		&row.UserID,
		&row.ReporterDisplay,
		&row.ReporterEmail,
		&attachmentID,
		&row.TargetType,
		&row.ReportType,
		&row.Message,
		&row.Status,
		&createdAt,
		&updatedAt,
		&row.AdminNote,
		&reviewedBy,
		&row.ReviewerDisplay,
		&reviewedAt,
		&row.CourseID,
		&row.CourseTitle,
		&pointContentID,
		&row.LevelTitle,
		&row.LessonTitle,
		&row.PointTitle,
		&row.PointType,
		&row.PointExternalURL,
		&row.AttachmentTitle,
		&row.AttachmentURL,
		&row.AttachmentFile,
		&attachmentSize,
		&row.AttachmentMime,
		&row.AttachmentContext,
		&row.AttachmentType,
		&row.RecommendationBlocked,
	); err != nil {
		return materialReportRow{}, err
	}

	if attachmentID.Valid {
		row.AttachmentID = &attachmentID.String
	}
	if reviewedBy.Valid {
		row.ReviewedBy = &reviewedBy.String
	}
	if pointContentID.Valid {
		row.PointContentID = &pointContentID.String
	}
	if reviewedAt.Valid {
		formatted := reviewedAt.Time.Format(time.RFC3339)
		row.ReviewedAt = &formatted
	}
	if attachmentSize.Valid {
		row.AttachmentSize = &attachmentSize.Int64
	}
	row.CreatedAt = createdAt.Format(time.RFC3339)
	row.UpdatedAt = updatedAt.Format(time.RFC3339)
	return row, nil
}

func materialReportRecommendationURLKey(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Host != "" {
		parsed.Fragment = ""
		parsed.Scheme = strings.ToLower(parsed.Scheme)
		parsed.Host = strings.ToLower(parsed.Host)
		parsed.Path = strings.TrimRight(parsed.EscapedPath(), "/")
		if parsed.Path == "" {
			parsed.Path = "/"
		}
		return strings.TrimRight(parsed.String(), "/")
	}
	return strings.TrimRight(strings.ToLower(trimmed), "/")
}
