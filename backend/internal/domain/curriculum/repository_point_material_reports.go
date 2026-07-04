package curriculum

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) CreateLearningPointMaterialReport(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointMaterialReportRequest) (uuid.UUID, CoursePointMaterialReport, error) {
	targetType, reportType, message, err := normalizePointMaterialReportInput(req.TargetType, req.ReportType, req.Message)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("begin create point material report tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}
	if req.AttachmentID != nil {
		var exists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM course_point_attachments
				WHERE id = $1
				  AND course_point_id = $2
			)
		`, *req.AttachmentID, pointID).Scan(&exists); err != nil {
			return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("check point report attachment: %w", err)
		}
		if !exists {
			return uuid.Nil, CoursePointMaterialReport{}, errLearningPointNotFound
		}
	}

	if targetType == "source" {
		var hasPendingSourceReport bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM course_point_material_reports
				WHERE course_point_id = $1
				  AND user_id = $2
				  AND target_type = 'source'
				  AND status IN ('open', 'reviewing')
				  AND replaced_at IS NULL
			)
		`, pointID, userID).Scan(&hasPendingSourceReport); err != nil {
			return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("check pending source material report: %w", err)
		}
		if hasPendingSourceReport {
			return uuid.Nil, CoursePointMaterialReport{}, errLearningPointArtifactInvalidInput
		}
	}

	var targetContentID *uuid.UUID
	var targetURL string
	var targetTitle string
	if targetType == "source" {
		if err := tx.QueryRow(ctx, `
			SELECT content_id, COALESCE(external_url, ''), title
			FROM course_points
			WHERE id = $1
		`, pointID).Scan(&targetContentID, &targetURL, &targetTitle); err != nil {
			return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("snapshot point report source: %w", err)
		}
	} else if req.AttachmentID != nil {
		if err := tx.QueryRow(ctx, `
			SELECT NULL::uuid, COALESCE(url, ''), title
			FROM course_point_attachments
			WHERE id = $1
			  AND course_point_id = $2
		`, *req.AttachmentID, pointID).Scan(&targetContentID, &targetURL, &targetTitle); err != nil {
			return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("snapshot point report attachment: %w", err)
		}
	}

	var report CoursePointMaterialReport
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_material_reports (
			course_point_id, user_id, attachment_id, target_type, report_type, message,
			target_content_id, target_url, target_title
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), NULLIF($9, ''))
		RETURNING id, course_point_id, user_id, attachment_id, target_type, report_type, message, status,
		          target_content_id, COALESCE(target_url, ''), COALESCE(target_title, ''),
		          replacement_content_id, COALESCE(replacement_url, ''), COALESCE(replacement_title, ''), replaced_at,
		          created_at, updated_at
	`, pointID, userID, req.AttachmentID, targetType, reportType, message, targetContentID, targetURL, targetTitle).Scan(
		&report.ID,
		&report.CoursePointID,
		&report.UserID,
		&report.AttachmentID,
		&report.TargetType,
		&report.ReportType,
		&report.Message,
		&report.Status,
		&report.TargetContentID,
		&report.TargetURL,
		&report.TargetTitle,
		&report.ReplacementContentID,
		&report.ReplacementURL,
		&report.ReplacementTitle,
		&report.ReplacedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("insert learning point material report: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "material_reported", map[string]any{
		"report_id":   report.ID,
		"target_type": targetType,
		"report_type": reportType,
	}); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}
	if err := r.recordRecommendationRolloutLearnerEventTx(ctx, tx, userID, pointID, "material_reported", report); err != nil {
		log.Printf("[recommendation/rollout] material report telemetry skipped: %v", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning draft after point material report: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning course after point material report: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("commit create point material report tx: %w", err)
	}
	return courseID, report, nil
}

func (r *Repository) ReplaceLearningPointSourceMaterial(ctx context.Context, userID, planetID, pointID uuid.UUID, req ReplaceLearningPointMaterialRequest) (uuid.UUID, CoursePointMaterialReport, error) {
	if req.ReportID == nil {
		return uuid.Nil, CoursePointMaterialReport{}, errLearningPointArtifactInvalidInput
	}
	title := ""
	if req.Title != nil {
		title = strings.TrimSpace(*req.Title)
	}
	replacementURL := ""
	if req.URL != nil {
		replacementURL = strings.TrimSpace(*req.URL)
	}
	thumbnailURL := ""
	if req.ThumbnailURL != nil {
		thumbnailURL = strings.TrimSpace(*req.ThumbnailURL)
	}
	if replacementURL == "" && req.ContentID == nil {
		return uuid.Nil, CoursePointMaterialReport{}, errLearningPointArtifactInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("begin replace point material tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}

	var report CoursePointMaterialReport
	if err := tx.QueryRow(ctx, `
		SELECT id, course_point_id, user_id, attachment_id, target_type, report_type, message, status,
		       target_content_id, COALESCE(target_url, ''), COALESCE(target_title, ''),
		       replacement_content_id, COALESCE(replacement_url, ''), COALESCE(replacement_title, ''), replaced_at,
		       created_at, updated_at
		FROM course_point_material_reports
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND target_type = 'source'
		  AND status IN ('open', 'reviewing')
		  AND replaced_at IS NULL
	`, *req.ReportID, pointID, userID).Scan(
		&report.ID,
		&report.CoursePointID,
		&report.UserID,
		&report.AttachmentID,
		&report.TargetType,
		&report.ReportType,
		&report.Message,
		&report.Status,
		&report.TargetContentID,
		&report.TargetURL,
		&report.TargetTitle,
		&report.ReplacementContentID,
		&report.ReplacementURL,
		&report.ReplacementTitle,
		&report.ReplacedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, errLearningPointNotFound
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_points
		SET content_id = $1,
		    external_url = NULLIF($2, ''),
		    thumbnail_url = COALESCE(NULLIF($3, ''), thumbnail_url),
		    updated_at = NOW()
		WHERE id = $4
	`, req.ContentID, replacementURL, thumbnailURL, pointID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("replace point source material: %w", err)
	}

	if err := tx.QueryRow(ctx, `
		UPDATE course_point_material_reports
		SET replacement_content_id = $1,
		    replacement_url = NULLIF($2, ''),
		    replacement_title = NULLIF($3, ''),
		    replaced_at = NOW(),
		    updated_at = NOW()
		WHERE id = $4
		RETURNING id, course_point_id, user_id, attachment_id, target_type, report_type, message, status,
		          target_content_id, COALESCE(target_url, ''), COALESCE(target_title, ''),
		          replacement_content_id, COALESCE(replacement_url, ''), COALESCE(replacement_title, ''), replaced_at,
		          created_at, updated_at
	`, req.ContentID, replacementURL, title, report.ID).Scan(
		&report.ID,
		&report.CoursePointID,
		&report.UserID,
		&report.AttachmentID,
		&report.TargetType,
		&report.ReportType,
		&report.Message,
		&report.Status,
		&report.TargetContentID,
		&report.TargetURL,
		&report.TargetTitle,
		&report.ReplacementContentID,
		&report.ReplacementURL,
		&report.ReplacementTitle,
		&report.ReplacedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("save point material replacement: %w", err)
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "material_replaced", map[string]any{
		"report_id":              report.ID,
		"target_content_id":      report.TargetContentID,
		"target_url":             report.TargetURL,
		"replacement_content_id": report.ReplacementContentID,
		"replacement_url":        report.ReplacementURL,
		"replacement_title":      report.ReplacementTitle,
	}); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}
	if err := r.recordRecommendationRolloutLearnerEventTx(ctx, tx, userID, pointID, "material_replaced", report); err != nil {
		log.Printf("[recommendation/rollout] material replacement telemetry skipped: %v", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning draft after point material replacement: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning course after point material replacement: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("commit replace point material tx: %w", err)
	}
	return courseID, report, nil
}

func (r *Repository) CancelLearningPointSourceMaterialReport(ctx context.Context, userID, planetID, pointID, reportID uuid.UUID) (uuid.UUID, CoursePointMaterialReport, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("begin cancel point material report tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}

	var report CoursePointMaterialReport
	if err := tx.QueryRow(ctx, `
		UPDATE course_point_material_reports
		SET status = 'cancelled',
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND target_type = 'source'
		  AND status IN ('open', 'reviewing')
		  AND replaced_at IS NULL
		RETURNING id, course_point_id, user_id, attachment_id, target_type, report_type, message, status,
		          target_content_id, COALESCE(target_url, ''), COALESCE(target_title, ''),
		          replacement_content_id, COALESCE(replacement_url, ''), COALESCE(replacement_title, ''), replaced_at,
		          created_at, updated_at
	`, reportID, pointID, userID).Scan(
		&report.ID,
		&report.CoursePointID,
		&report.UserID,
		&report.AttachmentID,
		&report.TargetType,
		&report.ReportType,
		&report.Message,
		&report.Status,
		&report.TargetContentID,
		&report.TargetURL,
		&report.TargetTitle,
		&report.ReplacementContentID,
		&report.ReplacementURL,
		&report.ReplacementTitle,
		&report.ReplacedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, errLearningPointNotFound
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "material_report_cancelled", map[string]any{
		"report_id":   report.ID,
		"target_type": report.TargetType,
		"report_type": report.ReportType,
	}); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, err
	}
	if err := r.recordRecommendationRolloutLearnerEventTx(ctx, tx, userID, pointID, "material_report_cancelled", report); err != nil {
		log.Printf("[recommendation/rollout] material report cancel telemetry skipped: %v", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning draft after point material report cancel: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("touch learning course after point material report cancel: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointMaterialReport{}, fmt.Errorf("commit cancel point material report tx: %w", err)
	}
	return courseID, report, nil
}
