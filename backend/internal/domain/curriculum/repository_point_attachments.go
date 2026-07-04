package curriculum

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ListLearningPointAttachments(ctx context.Context, userID, planetID, pointID uuid.UUID) ([]CoursePointAttachment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin list point attachments tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT cpa.id, cpa.course_point_id, cpa.artifact_id, cpa.user_id, cpa.provider, cpa.attachment_type,
		       COALESCE(cpa.source_context, 'work_attachment'),
		       COALESCE(cpa.title, ''), COALESCE(cpa.url, ''), COALESCE(cpa.file_path, ''),
		       cpa.file_size, COALESCE(cpa.mime_type, ''), cpa.created_at, cpa.updated_at
		FROM course_point_attachments cpa
		JOIN course_points cp ON cp.id = cpa.course_point_id
		WHERE cpa.course_point_id = $1
		  AND cp.course_id = $2
		ORDER BY cpa.created_at ASC
	`, pointID, courseID)
	if err != nil {
		return nil, fmt.Errorf("list learning point attachments: %w", err)
	}
	defer rows.Close()

	attachments := make([]CoursePointAttachment, 0)
	for rows.Next() {
		attachment, err := scanCoursePointAttachment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan learning point attachment: %w", err)
		}
		attachments = append(attachments, attachment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning point attachments: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit list point attachments tx: %w", err)
	}
	return attachments, nil
}

func (r *Repository) ResolveLearningPointCourseID(ctx context.Context, userID, planetID, pointID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin resolve learning point course tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit resolve learning point course tx: %w", err)
	}
	return courseID, nil
}

func (r *Repository) CountLearningPointAttachments(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceContext string) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point attachments tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	}
	count, err := countLearningPointAttachmentsTx(ctx, tx, pointID, sourceContext)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point attachments tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CountLearningPointAttachmentsByType(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceContext, attachmentType string) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point attachments by type tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	}
	count, err := countLearningPointAttachmentsByTypeTx(ctx, tx, pointID, sourceContext, attachmentType)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point attachments by type tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CountLearningPointAttachmentsExcludingType(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceContext, attachmentType string) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point attachments excluding type tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	}
	count, err := countLearningPointAttachmentsExcludingTypeTx(ctx, tx, pointID, sourceContext, attachmentType)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point attachments excluding type tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CountLearningPointMaterialFileAttachments(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceContext string) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point material file attachments tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	}
	count, err := countLearningPointMaterialFileAttachmentsTx(ctx, tx, pointID, sourceContext)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point material file attachments tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CountLearningPointArtifactAttachmentsByType(ctx context.Context, userID, planetID, pointID, artifactID uuid.UUID, attachmentType string) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point artifact attachments by type tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	} else if err := ensureLearningPointArtifactTx(ctx, tx, pointID, courseID, artifactID); err != nil {
		return 0, err
	}
	count, err := countLearningPointArtifactAttachmentsByTypeTx(ctx, tx, pointID, artifactID, attachmentType)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point artifact attachments by type tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CountLearningPointArtifactFileAttachments(ctx context.Context, userID, planetID, pointID, artifactID uuid.UUID) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin count point artifact file attachments tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return 0, err
	} else if err := ensureLearningPointArtifactTx(ctx, tx, pointID, courseID, artifactID); err != nil {
		return 0, err
	}
	count, err := countLearningPointArtifactFileAttachmentsTx(ctx, tx, pointID, artifactID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit count point artifact file attachments tx: %w", err)
	}
	return count, nil
}

func (r *Repository) CreateLearningPointAttachment(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointAttachmentRequest) (uuid.UUID, CoursePointAttachment, error) {
	provider, attachmentType, sourceContext, title, url, filePath, fileSize, mimeType, err := normalizePointAttachmentInput(req.Provider, req.AttachmentType, req.SourceContext, req.Title, req.URL, req.FilePath, req.MimeType, req.FileSize)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("begin create point attachment tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if sourceContext == "research_material" {
		confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
		if err != nil {
			return uuid.Nil, CoursePointAttachment{}, err
		}
		if confirmed {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointMaterialLocked
		}

		if attachmentType == "video" || attachmentType == "subtitle" || attachmentType == "thumbnail" {
			videoCount, err := countLearningPointAttachmentsByTypeTx(ctx, tx, pointID, sourceContext, attachmentType)
			if err != nil {
				return uuid.Nil, CoursePointAttachment{}, err
			}
			if videoCount >= 1 {
				return uuid.Nil, CoursePointAttachment{}, errLearningPointAttachmentLimitExceeded
			}
		} else {
			count, err := countLearningPointMaterialFileAttachmentsTx(ctx, tx, pointID, sourceContext)
			if err != nil {
				return uuid.Nil, CoursePointAttachment{}, err
			}
			if count >= maxLearningResearchMaterialAttachmentCount {
				return uuid.Nil, CoursePointAttachment{}, errLearningPointAttachmentLimitExceeded
			}
		}
	}
	var artifactID *uuid.UUID
	if sourceContext == "artifact" {
		if req.ArtifactID == nil {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointArtifactInvalidInput
		}
		if err := ensureLearningPointArtifactTx(ctx, tx, pointID, courseID, *req.ArtifactID); err != nil {
			return uuid.Nil, CoursePointAttachment{}, err
		}
		artifactID = req.ArtifactID
		if attachmentType == "video" || attachmentType == "subtitle" || attachmentType == "thumbnail" {
			count, err := countLearningPointArtifactAttachmentsByTypeTx(ctx, tx, pointID, *artifactID, attachmentType)
			if err != nil {
				return uuid.Nil, CoursePointAttachment{}, err
			}
			if count >= 1 {
				return uuid.Nil, CoursePointAttachment{}, errLearningPointAttachmentLimitExceeded
			}
		} else {
			count, err := countLearningPointArtifactFileAttachmentsTx(ctx, tx, pointID, *artifactID)
			if err != nil {
				return uuid.Nil, CoursePointAttachment{}, err
			}
			if count >= maxLearningResearchMaterialAttachmentCount {
				return uuid.Nil, CoursePointAttachment{}, errLearningPointAttachmentLimitExceeded
			}
		}
	}

	attachment, err := scanCoursePointAttachment(tx.QueryRow(ctx, `
		INSERT INTO course_point_attachments (
			course_point_id, artifact_id, user_id, provider, attachment_type, source_context, title, url, file_path, file_size, mime_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, course_point_id, artifact_id, user_id, provider, attachment_type,
		          COALESCE(source_context, 'work_attachment'),
		          COALESCE(title, ''), COALESCE(url, ''), COALESCE(file_path, ''),
		          file_size, COALESCE(mime_type, ''), created_at, updated_at
	`, pointID, artifactID, userID, provider, attachmentType, sourceContext, title, url, filePath, fileSize, mimeType))
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("insert learning point attachment: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "attachment_added", map[string]any{"attachment_id": attachment.ID, "source_context": attachment.SourceContext, "title": attachment.Title}); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := touchLearningCourseAfterPointAttachment(ctx, tx, draftID, courseID, userID); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("commit create point attachment tx: %w", err)
	}
	return courseID, attachment, nil
}

func (r *Repository) UpdateLearningPointAttachment(ctx context.Context, userID, planetID, pointID, attachmentID uuid.UUID, req UpdateLearningPointAttachmentRequest) (uuid.UUID, CoursePointAttachment, error) {
	provider, attachmentType, sourceContext, title, url, filePath, fileSize, mimeType, err := normalizePointAttachmentInput(req.Provider, req.AttachmentType, req.SourceContext, req.Title, req.URL, req.FilePath, req.MimeType, req.FileSize)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("begin update point attachment tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}

	var existingSourceContext string
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(source_context, 'work_attachment')
		FROM course_point_attachments
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
	`, attachmentID, pointID, courseID).Scan(&existingSourceContext); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("get learning point attachment source context: %w", err)
	}
	if sourceContext == "research_material" || existingSourceContext == "research_material" {
		confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
		if err != nil {
			return uuid.Nil, CoursePointAttachment{}, err
		}
		if confirmed {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointMaterialLocked
		}
	}
	var artifactID *uuid.UUID
	if sourceContext == "artifact" {
		if req.ArtifactID == nil {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointArtifactInvalidInput
		}
		if err := ensureLearningPointArtifactTx(ctx, tx, pointID, courseID, *req.ArtifactID); err != nil {
			return uuid.Nil, CoursePointAttachment{}, err
		}
		artifactID = req.ArtifactID
	}

	attachment, err := scanCoursePointAttachment(tx.QueryRow(ctx, `
		UPDATE course_point_attachments
		SET provider = $4,
		    attachment_type = $5,
		    source_context = $6,
		    artifact_id = $7,
		    title = $8,
		    url = $9,
		    file_path = $10,
		    file_size = $11,
		    mime_type = $12,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
		RETURNING id, course_point_id, artifact_id, user_id, provider, attachment_type,
		          COALESCE(source_context, 'work_attachment'),
		          COALESCE(title, ''), COALESCE(url, ''), COALESCE(file_path, ''),
		          file_size, COALESCE(mime_type, ''), created_at, updated_at
	`, attachmentID, pointID, courseID, provider, attachmentType, sourceContext, artifactID, title, url, filePath, fileSize, mimeType))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("update learning point attachment: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "attachment_added", map[string]any{"attachment_id": attachment.ID, "action": "updated", "source_context": attachment.SourceContext, "title": attachment.Title}); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := touchLearningCourseAfterPointAttachment(ctx, tx, draftID, courseID, userID); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("commit update point attachment tx: %w", err)
	}
	return courseID, attachment, nil
}

func (r *Repository) DeleteLearningPointAttachment(ctx context.Context, userID, planetID, pointID, attachmentID uuid.UUID) (uuid.UUID, CoursePointAttachment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("begin delete point attachment tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}

	var sourceContext string
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(source_context, 'work_attachment')
		FROM course_point_attachments
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
	`, attachmentID, pointID, courseID).Scan(&sourceContext); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("get learning point attachment source context: %w", err)
	}
	if sourceContext == "research_material" {
		confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
		if err != nil {
			return uuid.Nil, CoursePointAttachment{}, err
		}
		if confirmed {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointMaterialLocked
		}
	}

	attachment, err := scanCoursePointAttachment(tx.QueryRow(ctx, `
		DELETE FROM course_point_attachments
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
		RETURNING id, course_point_id, artifact_id, user_id, provider, attachment_type,
		          COALESCE(source_context, 'work_attachment'),
		          COALESCE(title, ''), COALESCE(url, ''), COALESCE(file_path, ''),
		          file_size, COALESCE(mime_type, ''), created_at, updated_at
	`, attachmentID, pointID, courseID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("delete learning point attachment: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "attachment_added", map[string]any{"attachment_id": attachmentID, "action": "deleted", "source_context": attachment.SourceContext, "title": attachment.Title}); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := touchLearningCourseAfterPointAttachment(ctx, tx, draftID, courseID, userID); err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("commit delete point attachment tx: %w", err)
	}
	return courseID, attachment, nil
}

func (r *Repository) GetLearningPointAttachment(ctx context.Context, userID, planetID, pointID, attachmentID uuid.UUID) (uuid.UUID, CoursePointAttachment, error) {
	return r.getPointAttachment(ctx, userID, planetID, pointID, attachmentID, false)
}

func (r *Repository) GetReadablePointAttachment(ctx context.Context, userID, planetID, pointID, attachmentID uuid.UUID) (uuid.UUID, CoursePointAttachment, error) {
	return r.getPointAttachment(ctx, userID, planetID, pointID, attachmentID, true)
}

func (r *Repository) getPointAttachment(ctx context.Context, userID, planetID, pointID, attachmentID uuid.UUID, allowArchived bool) (uuid.UUID, CoursePointAttachment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("begin get point attachment tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var courseID uuid.UUID
	if allowArchived {
		_, courseID, err = r.getReadablePointContextTx(ctx, tx, userID, planetID, pointID)
	} else {
		_, courseID, err = r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	}
	if err != nil {
		return uuid.Nil, CoursePointAttachment{}, err
	}

	attachment, err := scanCoursePointAttachment(tx.QueryRow(ctx, `
		SELECT cpa.id, cpa.course_point_id, cpa.artifact_id, cpa.user_id, cpa.provider, cpa.attachment_type,
		       COALESCE(cpa.source_context, 'work_attachment'),
		       COALESCE(cpa.title, ''), COALESCE(cpa.url, ''), COALESCE(cpa.file_path, ''),
		       cpa.file_size, COALESCE(cpa.mime_type, ''), cpa.created_at, cpa.updated_at
		FROM course_point_attachments cpa
		JOIN course_points cp ON cp.id = cpa.course_point_id
		WHERE cpa.id = $1
		  AND cpa.course_point_id = $2
		  AND cp.course_id = $3
	`, attachmentID, pointID, courseID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointAttachment{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("get learning point attachment: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointAttachment{}, fmt.Errorf("commit get point attachment tx: %w", err)
	}
	return courseID, attachment, nil
}

type coursePointAttachmentScanner interface {
	Scan(dest ...any) error
}

func scanCoursePointAttachment(scanner coursePointAttachmentScanner) (CoursePointAttachment, error) {
	var attachment CoursePointAttachment
	if err := scanner.Scan(
		&attachment.ID,
		&attachment.CoursePointID,
		&attachment.ArtifactID,
		&attachment.UserID,
		&attachment.Provider,
		&attachment.AttachmentType,
		&attachment.SourceContext,
		&attachment.Title,
		&attachment.URL,
		&attachment.FilePath,
		&attachment.FileSize,
		&attachment.MimeType,
		&attachment.CreatedAt,
		&attachment.UpdatedAt,
	); err != nil {
		return CoursePointAttachment{}, err
	}
	return attachment, nil
}

func touchLearningCourseAfterPointAttachment(ctx context.Context, tx pgx.Tx, draftID, courseID, userID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return fmt.Errorf("touch learning draft after point attachment change: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return fmt.Errorf("touch learning course after point attachment change: %w", err)
	}
	return nil
}

func countLearningPointAttachmentsTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID, sourceContext string) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND COALESCE(source_context, 'work_attachment') = $2
	`, pointID, sourceContext).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point attachments: %w", err)
	}
	return count, nil
}

func countLearningPointAttachmentsByTypeTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID, sourceContext, attachmentType string) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND COALESCE(source_context, 'work_attachment') = $2
		  AND attachment_type = $3
	`, pointID, sourceContext, attachmentType).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point attachments by type: %w", err)
	}
	return count, nil
}

func countLearningPointAttachmentsExcludingTypeTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID, sourceContext, attachmentType string) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND COALESCE(source_context, 'work_attachment') = $2
		  AND attachment_type <> $3
	`, pointID, sourceContext, attachmentType).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point attachments excluding type: %w", err)
	}
	return count, nil
}

func countLearningPointMaterialFileAttachmentsTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID, sourceContext string) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND COALESCE(source_context, 'work_attachment') = $2
		  AND attachment_type NOT IN ('video', 'subtitle', 'thumbnail')
	`, pointID, sourceContext).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point material file attachments: %w", err)
	}
	return count, nil
}

func ensureLearningPointArtifactTx(ctx context.Context, tx pgx.Tx, pointID, courseID, artifactID uuid.UUID) error {
	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM course_point_artifacts cpa
			JOIN course_points cp ON cp.id = cpa.course_point_id
			WHERE cpa.id = $1
			  AND cpa.course_point_id = $2
			  AND cp.course_id = $3
		)
	`, artifactID, pointID, courseID).Scan(&exists); err != nil {
		return fmt.Errorf("check learning point artifact: %w", err)
	}
	if !exists {
		return errLearningPointNotFound
	}
	return nil
}

func countLearningPointArtifactAttachmentsByTypeTx(ctx context.Context, tx pgx.Tx, pointID, artifactID uuid.UUID, attachmentType string) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND artifact_id = $2
		  AND COALESCE(source_context, 'work_attachment') = 'artifact'
		  AND attachment_type = $3
	`, pointID, artifactID, attachmentType).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point artifact attachments by type: %w", err)
	}
	return count, nil
}

func countLearningPointArtifactFileAttachmentsTx(ctx context.Context, tx pgx.Tx, pointID, artifactID uuid.UUID) (int, error) {
	var count int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_attachments
		WHERE course_point_id = $1
		  AND artifact_id = $2
		  AND COALESCE(source_context, 'work_attachment') = 'artifact'
		  AND attachment_type NOT IN ('video', 'subtitle', 'thumbnail')
	`, pointID, artifactID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count learning point artifact file attachments: %w", err)
	}
	return count, nil
}
