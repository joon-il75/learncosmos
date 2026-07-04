package curriculum

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ListLearningPointArtifacts(ctx context.Context, userID, planetID, pointID uuid.UUID) ([]CoursePointArtifact, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin list point artifacts tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT id, course_point_id, artifact_type, title, url, description,
		       point_category, production_process, learned_points, difficult_points, visibility,
		       order_index, created_at, updated_at
		FROM course_point_artifacts
		WHERE course_point_id = $1
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $1 AND course_id = $2
		  )
		ORDER BY order_index ASC, created_at ASC
	`, pointID, courseID)
	if err != nil {
		return nil, fmt.Errorf("list learning point artifacts: %w", err)
	}
	defer rows.Close()

	artifacts := make([]CoursePointArtifact, 0)
	for rows.Next() {
		var artifact CoursePointArtifact
		if err := rows.Scan(
			&artifact.ID,
			&artifact.CoursePointID,
			&artifact.ArtifactType,
			&artifact.Title,
			&artifact.URL,
			&artifact.Description,
			&artifact.PointCategory,
			&artifact.ProductionProcess,
			&artifact.LearnedPoints,
			&artifact.DifficultPoints,
			&artifact.Visibility,
			&artifact.OrderIndex,
			&artifact.CreatedAt,
			&artifact.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan learning point artifact: %w", err)
		}
		artifacts = append(artifacts, artifact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning point artifacts: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit list point artifacts tx: %w", err)
	}
	return artifacts, nil
}

func (r *Repository) CreateLearningPointArtifact(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointArtifactRequest) (uuid.UUID, CoursePointArtifact, error) {
	if req.ArtifactType == nil || req.Title == nil || req.URL == nil || req.Description == nil {
		return uuid.Nil, CoursePointArtifact{}, errLearningPointArtifactInvalidInput
	}
	artifactType, title, url, description, err := normalizePointArtifactInput(*req.ArtifactType, *req.Title, *req.URL, *req.Description)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}
	pointCategory, productionProcess, learnedPoints, difficultPoints, visibility, err := normalizePointArtifactExtraInput(req.PointCategory, req.ProductionProcess, req.LearnedPoints, req.DifficultPoints, req.Visibility)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}
	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("begin create point artifact tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}

	var artifact CoursePointArtifact
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_artifacts (
			course_point_id, artifact_type, title, url, description,
			point_category, production_process, learned_points, difficult_points, visibility, order_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, course_point_id, artifact_type, title, url, description,
		          point_category, production_process, learned_points, difficult_points, visibility,
		          order_index, created_at, updated_at
	`, pointID, artifactType, title, url, description, pointCategory, productionProcess, learnedPoints, difficultPoints, visibility, orderIndex).Scan(
		&artifact.ID,
		&artifact.CoursePointID,
		&artifact.ArtifactType,
		&artifact.Title,
		&artifact.URL,
		&artifact.Description,
		&artifact.PointCategory,
		&artifact.ProductionProcess,
		&artifact.LearnedPoints,
		&artifact.DifficultPoints,
		&artifact.Visibility,
		&artifact.OrderIndex,
		&artifact.CreatedAt,
		&artifact.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("insert learning point artifact: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "artifact_submitted", map[string]any{"artifact_id": artifact.ID, "title": artifact.Title}); err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("touch learning draft after point artifact create: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("touch learning course after point artifact create: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("commit create point artifact tx: %w", err)
	}
	return courseID, artifact, nil
}

func (r *Repository) UpdateLearningPointArtifactItem(ctx context.Context, userID, planetID, pointID, artifactID uuid.UUID, req UpdateLearningPointArtifactRequest) (uuid.UUID, CoursePointArtifact, error) {
	if req.ArtifactType == nil || req.Title == nil || req.URL == nil || req.Description == nil {
		return uuid.Nil, CoursePointArtifact{}, errLearningPointArtifactInvalidInput
	}
	artifactType, title, url, description, err := normalizePointArtifactInput(*req.ArtifactType, *req.Title, *req.URL, *req.Description)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}
	pointCategory, productionProcess, learnedPoints, difficultPoints, visibility, err := normalizePointArtifactExtraInput(req.PointCategory, req.ProductionProcess, req.LearnedPoints, req.DifficultPoints, req.Visibility)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}
	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("begin update point artifact tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}

	var artifact CoursePointArtifact
	if err := tx.QueryRow(ctx, `
		UPDATE course_point_artifacts
		SET artifact_type = $4,
		    title = $5,
		    url = $6,
		    description = $7,
		    point_category = $8,
		    production_process = $9,
		    learned_points = $10,
		    difficult_points = $11,
		    visibility = $12,
		    order_index = $13,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
		RETURNING id, course_point_id, artifact_type, title, url, description,
		          point_category, production_process, learned_points, difficult_points, visibility,
		          order_index, created_at, updated_at
	`, artifactID, pointID, courseID, artifactType, title, url, description, pointCategory, productionProcess, learnedPoints, difficultPoints, visibility, orderIndex).Scan(
		&artifact.ID,
		&artifact.CoursePointID,
		&artifact.ArtifactType,
		&artifact.Title,
		&artifact.URL,
		&artifact.Description,
		&artifact.PointCategory,
		&artifact.ProductionProcess,
		&artifact.LearnedPoints,
		&artifact.DifficultPoints,
		&artifact.Visibility,
		&artifact.OrderIndex,
		&artifact.CreatedAt,
		&artifact.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointArtifact{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("update learning point artifact: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "artifact_submitted", map[string]any{"artifact_id": artifact.ID, "action": "updated", "title": artifact.Title}); err != nil {
		return uuid.Nil, CoursePointArtifact{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("touch learning draft after point artifact update: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("touch learning course after point artifact update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointArtifact{}, fmt.Errorf("commit update point artifact tx: %w", err)
	}
	return courseID, artifact, nil
}

func (r *Repository) DeleteLearningPointArtifact(ctx context.Context, userID, planetID, pointID, artifactID uuid.UUID) (uuid.UUID, []CoursePointAttachment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("begin delete point artifact tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	rows, err := tx.Query(ctx, `
		SELECT cpa.id, cpa.course_point_id, cpa.artifact_id, cpa.user_id, cpa.provider, cpa.attachment_type,
		       COALESCE(cpa.source_context, 'work_attachment'),
		       COALESCE(cpa.title, ''), COALESCE(cpa.url, ''), COALESCE(cpa.file_path, ''),
		       cpa.file_size, COALESCE(cpa.mime_type, ''), cpa.created_at, cpa.updated_at
		FROM course_point_attachments cpa
		WHERE cpa.course_point_id = $1
		  AND cpa.artifact_id = $2
		  AND COALESCE(cpa.source_context, 'work_attachment') = 'artifact'
		ORDER BY cpa.created_at ASC
	`, pointID, artifactID)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("list artifact attachments before delete: %w", err)
	}
	deletedAttachments := make([]CoursePointAttachment, 0)
	for rows.Next() {
		attachment, err := scanCoursePointAttachment(rows)
		if err != nil {
			rows.Close()
			return uuid.Nil, nil, fmt.Errorf("scan artifact attachment before delete: %w", err)
		}
		deletedAttachments = append(deletedAttachments, attachment)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return uuid.Nil, nil, fmt.Errorf("iterate artifact attachments before delete: %w", err)
	}
	rows.Close()

	var deletedTitle string
	if err := tx.QueryRow(ctx, `
		DELETE FROM course_point_artifacts
		WHERE id = $1
		  AND course_point_id = $2
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $3
		  )
		RETURNING title
	`, artifactID, pointID, courseID).Scan(&deletedTitle); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil, errLearningPointNotFound
		}
		return uuid.Nil, nil, fmt.Errorf("delete learning point artifact: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "artifact_submitted", map[string]any{"artifact_id": artifactID, "action": "deleted", "title": deletedTitle}); err != nil {
		return uuid.Nil, nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, nil, fmt.Errorf("touch learning draft after point artifact delete: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, nil, fmt.Errorf("touch learning course after point artifact delete: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, nil, fmt.Errorf("commit delete point artifact tx: %w", err)
	}
	return courseID, deletedAttachments, nil
}
