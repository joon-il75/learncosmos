package explorer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	errSubRegionLessonHasPoints = errors.New("subregion draft lesson has linked points")
)

func (r *Repository) CountSubRegionsByRegion(ctx context.Context, regionID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM explorer_subregions WHERE region_id=$1 AND status='active'`,
		regionID,
	).Scan(&count)
	return count, err
}

func (r *Repository) InsertSubRegion(ctx context.Context, req CreateSubRegionRequest, lessonID *uuid.UUID) (*SubRegion, error) {
	var sr SubRegion
	err := r.pool.QueryRow(ctx, `
		INSERT INTO explorer_subregions (region_id, name, description, order_index, course_draft_lesson_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, region_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	`, req.RegionID, req.Name, req.Description, req.OrderIndex, lessonID).Scan(
		&sr.ID, &sr.RegionID, &sr.CourseDraftLessonID, &sr.Name, &sr.Description,
		&sr.OrderIndex, &sr.Status, &sr.CreatedAt, &sr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (r *Repository) SetSubRegionLessonID(ctx context.Context, subRegionID, lessonID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE explorer_subregions SET course_draft_lesson_id = $1, updated_at = now() WHERE id = $2`,
		lessonID, subRegionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetSubRegionByID(ctx context.Context, subRegionID uuid.UUID) (*SubRegion, error) {
	var sr SubRegion
	err := r.pool.QueryRow(ctx, `
		SELECT id, region_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
		FROM explorer_subregions WHERE id=$1
	`, subRegionID).Scan(
		&sr.ID, &sr.RegionID, &sr.CourseDraftLessonID, &sr.Name, &sr.Description,
		&sr.OrderIndex, &sr.Status, &sr.CreatedAt, &sr.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &sr, err
}

func (r *Repository) UpdateSubRegion(ctx context.Context, subRegionID uuid.UUID, req UpdateSubRegionRequest) (*SubRegion, error) {
	var sr SubRegion
	err := r.pool.QueryRow(ctx, `
		UPDATE explorer_subregions SET
			region_id   = COALESCE($2, region_id),
			name        = COALESCE($3, name),
			description = COALESCE($4, description),
			order_index = COALESCE($5, order_index),
			updated_at  = now()
		WHERE id=$1
		RETURNING id, region_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	`, subRegionID, req.RegionID, req.Name, req.Description, req.OrderIndex,
	).Scan(
		&sr.ID, &sr.RegionID, &sr.CourseDraftLessonID, &sr.Name, &sr.Description,
		&sr.OrderIndex, &sr.Status, &sr.CreatedAt, &sr.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &sr, err
}

func (r *Repository) UpdateSubRegionLesson(
	ctx context.Context,
	draftID,
	userID,
	lessonID uuid.UUID,
	title *string,
	objective *string,
	summary *string,
	orderIndex *int,
	parentLessonID *uuid.UUID,
	recommendationSearchSpec *string,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_lessons cdl
		SET title = COALESCE($1, cdl.title),
			objective = COALESCE($2, cdl.objective),
			summary = COALESCE($3, cdl.summary),
			order_index = COALESCE($4, cdl.order_index),
			parent_lesson_id = COALESCE($5, cdl.parent_lesson_id),
			recommendation_search_spec = COALESCE($6::jsonb, cdl.recommendation_search_spec),
			updated_at = NOW()
		FROM course_drafts cd
		WHERE cdl.id = $7
		  AND cdl.course_draft_id = $8
		  AND cdl.parent_lesson_id IS NOT NULL
		  AND cdl.course_draft_id = cd.id
		  AND cd.user_id = $9
	`,
		title, objective, summary, orderIndex, parentLessonID, recommendationSearchSpec,
		lessonID, draftID, userID,
	)
	if err != nil {
		return fmt.Errorf("update subregion lesson: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateSubRegionStatus(ctx context.Context, subRegionID uuid.UUID, status ItemStatus) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE explorer_subregions SET status=$1, updated_at=now() WHERE id=$2`,
		status, subRegionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteSubRegionDraftLesson(ctx context.Context, draftID, userID, lessonID uuid.UUID) error {
	var linkedPoints int
	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_draft_points
		WHERE course_draft_id = $1
		  AND course_draft_lesson_id = $2
	`, draftID, lessonID).Scan(&linkedPoints); err != nil {
		return fmt.Errorf("count subregion draft lesson points: %w", err)
	}
	if linkedPoints > 0 {
		return errSubRegionLessonHasPoints
	}

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM course_draft_lessons cdl
		USING course_drafts cd
		WHERE cdl.id = $1
		  AND cdl.course_draft_id = $2
		  AND cdl.course_draft_id = cd.id
		  AND cd.user_id = $3
	`, lessonID, draftID, userID)
	if err != nil {
		return fmt.Errorf("delete subregion draft lesson: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteSubRegionCascade(ctx context.Context, subRegionID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		DELETE FROM course_draft_points
		WHERE id IN (
			SELECT draft_point_id
			FROM explorer_nodes
			WHERE parent_kind='subregion'
			  AND parent_id = $1
			  AND draft_point_id IS NOT NULL
		)
		  AND course_draft_id = (
			SELECT r.course_draft_id
			FROM explorer_subregions sr
			JOIN explorer_regions r ON r.id = sr.region_id
			WHERE sr.id = $1
		  )
	`, subRegionID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`DELETE FROM explorer_nodes WHERE parent_kind='subregion' AND parent_id=$1`,
		subRegionID,
	); err != nil {
		return err
	}

	tag, err := tx.Exec(ctx,
		`DELETE FROM explorer_subregions WHERE id=$1`,
		subRegionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetMainLessonIDByRegion(ctx context.Context, regionID, courseDraftID uuid.UUID, regionOrder int, regionName string) (uuid.UUID, error) {
	var lessonID uuid.UUID
	if err := r.pool.QueryRow(ctx, `
		SELECT cdl.id
		FROM explorer_regions er
		JOIN course_draft_lessons cdl ON cdl.id = er.course_draft_lesson_id
		WHERE er.id = $1
		  AND er.course_draft_id = $2
		  AND cdl.course_draft_id = er.course_draft_id
		  AND cdl.parent_lesson_id IS NULL
	`, regionID, courseDraftID).Scan(&lessonID); err == nil {
		return lessonID, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, err
	}

	if err := r.pool.QueryRow(ctx, `
		SELECT id
		FROM course_draft_lessons
		WHERE id = $1
		  AND course_draft_id = $2
		  AND parent_lesson_id IS NULL
	`, regionID, courseDraftID).Scan(&lessonID); err == nil {
		return lessonID, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, err
	}

	if err := r.pool.QueryRow(ctx, `
		SELECT id
		FROM course_draft_lessons
		WHERE course_draft_id = $1
		  AND parent_lesson_id IS NULL
		  AND order_index = $2
		ORDER BY id
		LIMIT 1
	`, courseDraftID, regionOrder).Scan(&lessonID); err == nil {
		return lessonID, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, err
	}

	trimmedName := strings.TrimSpace(regionName)
	if trimmedName != "" {
		if err := r.pool.QueryRow(ctx, `
			SELECT id
			FROM course_draft_lessons
			WHERE course_draft_id = $1
			  AND parent_lesson_id IS NULL
			  AND title = $2
			ORDER BY id
			LIMIT 1
		`, courseDraftID, trimmedName).Scan(&lessonID); err == nil {
			return lessonID, nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, err
		}
	}

	return uuid.Nil, ErrNotFound
}
