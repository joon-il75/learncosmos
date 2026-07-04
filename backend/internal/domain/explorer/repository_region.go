package explorer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	errRegionLessonHasPoints = errors.New("region draft lesson has linked points")
)

func (r *Repository) InsertRegion(ctx context.Context, req CreateRegionRequest, lessonID *uuid.UUID) (*Region, error) {
	var reg Region
	err := r.pool.QueryRow(ctx, `
		INSERT INTO explorer_regions (course_draft_id, course_draft_lesson_id, name, description, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, course_draft_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	`, req.CourseDraftID, lessonID, req.Name, req.Description, req.OrderIndex,
	).Scan(
		&reg.ID, &reg.CourseDraftID, &reg.CourseDraftLessonID, &reg.Name, &reg.Description,
		&reg.OrderIndex, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *Repository) UpdateRegionLesson(
	ctx context.Context,
	draftID,
	userID,
	lessonID uuid.UUID,
	title *string,
	objective *string,
	summary *string,
	orderIndex *int,
	recommendationSearchSpec *string,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_lessons cdl
		SET title = COALESCE($1, cdl.title),
			objective = COALESCE($2, cdl.objective),
			summary = COALESCE($3, cdl.summary),
			order_index = COALESCE($4, cdl.order_index),
			recommendation_search_spec = COALESCE($5::jsonb, cdl.recommendation_search_spec),
			updated_at = NOW()
		FROM course_drafts cd
		WHERE cdl.id = $6
		  AND cdl.course_draft_id = $7
		  AND cdl.parent_lesson_id IS NULL
		  AND cdl.course_draft_id = cd.id
		  AND cd.user_id = $8
	`,
		title, objective, summary, orderIndex, recommendationSearchSpec,
		lessonID, draftID, userID,
	)
	if err != nil {
		return fmt.Errorf("update region lesson: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteRegionDraftLesson(ctx context.Context, draftID, userID, lessonID uuid.UUID) error {
	var hasPoints bool
	if err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM course_draft_points cdp
			JOIN course_drafts cd ON cd.id = cdp.course_draft_id
			WHERE cdp.course_draft_lesson_id = $1
			  AND cdp.course_draft_id = $2
			  AND cd.user_id = $3
		)
	`, lessonID, draftID, userID).Scan(&hasPoints); err != nil {
		return fmt.Errorf("check region lesson points: %w", err)
	}
	if hasPoints {
		return errRegionLessonHasPoints
	}

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM course_draft_lessons cdl
		USING course_drafts cd
		WHERE cdl.id = $1
		  AND cdl.course_draft_id = $2
		  AND cdl.parent_lesson_id IS NULL
		  AND cdl.course_draft_id = cd.id
		  AND cd.user_id = $3
	`, lessonID, draftID, userID)
	if err != nil {
		return fmt.Errorf("delete region draft lesson: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetRegionByID(ctx context.Context, regionID uuid.UUID) (*Region, error) {
	var reg Region
	err := r.pool.QueryRow(ctx, `
		SELECT id, course_draft_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
		FROM explorer_regions WHERE id = $1
	`, regionID).Scan(
		&reg.ID, &reg.CourseDraftID, &reg.CourseDraftLessonID, &reg.Name, &reg.Description,
		&reg.OrderIndex, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &reg, err
}

func (r *Repository) UpdateRegion(ctx context.Context, regionID uuid.UUID, req UpdateRegionRequest) (*Region, error) {
	var reg Region
	err := r.pool.QueryRow(ctx, `
		UPDATE explorer_regions SET
			name        = COALESCE($2, name),
			description = COALESCE($3, description),
			order_index = COALESCE($4, order_index),
			updated_at  = now()
		WHERE id = $1
		RETURNING id, course_draft_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	`, regionID, req.Name, req.Description, req.OrderIndex,
	).Scan(
		&reg.ID, &reg.CourseDraftID, &reg.CourseDraftLessonID, &reg.Name, &reg.Description,
		&reg.OrderIndex, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &reg, err
}

func (r *Repository) UpdateRegionStatus(ctx context.Context, regionID uuid.UUID, status ItemStatus) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE explorer_regions SET status=$1, updated_at=now() WHERE id=$2`,
		status, regionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateRegionStatusWithChildren(ctx context.Context, regionID uuid.UUID, status ItemStatus) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx,
		`UPDATE explorer_regions SET status=$1, updated_at=now() WHERE id=$2`,
		status, regionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	rows, err := tx.Query(ctx,
		`SELECT id FROM explorer_subregions WHERE region_id=$1`, regionID,
	)
	if err != nil {
		return err
	}
	var subIDs []uuid.UUID
	for rows.Next() {
		var sid uuid.UUID
		if err := rows.Scan(&sid); err != nil {
			rows.Close()
			return err
		}
		subIDs = append(subIDs, sid)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for _, sid := range subIDs {
		if _, err := tx.Exec(ctx,
			`UPDATE explorer_subregions SET status=$1, updated_at=now() WHERE id=$2`,
			status, sid,
		); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`UPDATE explorer_nodes SET status=$1, updated_at=now()
			 WHERE parent_kind='subregion' AND parent_id=$2`,
			status, sid,
		); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(ctx,
		`UPDATE explorer_nodes SET status=$1, updated_at=now()
		 WHERE parent_kind='region' AND parent_id=$2`,
		status, regionID,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) DeleteRegionCascade(ctx context.Context, regionID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx,
		`SELECT id, course_draft_lesson_id FROM explorer_subregions WHERE region_id=$1`, regionID,
	)
	if err != nil {
		return err
	}
	var subIDs []uuid.UUID
	var subLessonIDs []uuid.UUID
	for rows.Next() {
		var sid uuid.UUID
		var lessonID *uuid.UUID
		if err := rows.Scan(&sid, &lessonID); err != nil {
			rows.Close()
			return err
		}
		subIDs = append(subIDs, sid)
		if lessonID != nil {
			subLessonIDs = append(subLessonIDs, *lessonID)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	if len(subIDs) > 0 {
		if _, err := tx.Exec(ctx, `
			DELETE FROM course_draft_points
			WHERE id IN (
				SELECT draft_point_id
				FROM explorer_nodes
				WHERE parent_kind='subregion'
				  AND parent_id = ANY($1)
				  AND draft_point_id IS NOT NULL
			)
			  AND course_draft_id = (
				SELECT course_draft_id
				FROM explorer_regions
				WHERE id = $2
			  )
		`, subIDs, regionID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM explorer_nodes WHERE parent_kind='subregion' AND parent_id=ANY($1)`,
			subIDs,
		); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM course_draft_points
		WHERE id IN (
			SELECT draft_point_id
			FROM explorer_nodes
			WHERE parent_kind='region'
			  AND parent_id = $1
			  AND draft_point_id IS NOT NULL
		)
		  AND course_draft_id = (
			SELECT course_draft_id
			FROM explorer_regions
			WHERE id = $1
		  )
	`, regionID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`DELETE FROM explorer_nodes WHERE parent_kind='region' AND parent_id=$1`,
		regionID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`DELETE FROM explorer_subregions WHERE region_id=$1`,
		regionID,
	); err != nil {
		return err
	}

	if len(subLessonIDs) > 0 {
		if _, err := tx.Exec(ctx, `
			DELETE FROM course_draft_lessons
			WHERE id = ANY($1)
			  AND parent_lesson_id IS NOT NULL
			  AND NOT EXISTS (
				SELECT 1 FROM course_draft_points WHERE course_draft_lesson_id = course_draft_lessons.id
			  )
		`, subLessonIDs); err != nil {
			return err
		}
	}

	tag, err := tx.Exec(ctx,
		`DELETE FROM explorer_regions WHERE id=$1`,
		regionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *Repository) SaveRegionOrder(ctx context.Context, courseDraftID uuid.UUID, regionOrder []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, rid := range regionOrder {
		tag, err := tx.Exec(ctx,
			`UPDATE explorer_regions SET order_index=$1, updated_at=now()
			 WHERE id=$2 AND course_draft_id=$3`,
			i, rid, courseDraftID,
		)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
	}

	return tx.Commit(ctx)
}
