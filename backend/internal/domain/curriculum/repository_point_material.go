package curriculum

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) LearningPointResearchMaterialConfirmed(ctx context.Context, userID, planetID, pointID uuid.UUID) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin research material confirmation state tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return false, err
	}
	confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
	if err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit research material confirmation state tx: %w", err)
	}
	return confirmed, nil
}

func learningResearchMaterialConfirmedTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID) (bool, error) {
	var confirmed bool
	if err := tx.QueryRow(ctx, `
		WITH block_state AS (
			SELECT
				COUNT(*) AS block_count
			FROM course_point_blocks
			WHERE course_point_id = $1
		),
		attachment_state AS (
			SELECT
				COUNT(*) AS attachment_count
			FROM course_point_attachments
			WHERE course_point_id = $1
			  AND source_context = 'research_material'
		),
		latest_confirmation_event AS (
			SELECT event_type, created_at
			FROM course_point_events
			WHERE course_point_id = $1
			  AND event_type IN ('research_material_confirmed', 'research_material_unconfirmed')
			ORDER BY created_at DESC
			LIMIT 1
		)
		SELECT COALESCE(
			(block_state.block_count + attachment_state.attachment_count) > 0
			AND latest_confirmation_event.event_type = 'research_material_confirmed',
			false
		)
		FROM block_state, attachment_state
		LEFT JOIN latest_confirmation_event ON true
	`, pointID).Scan(&confirmed); err != nil {
		return false, fmt.Errorf("get research material confirmation state: %w", err)
	}
	return confirmed, nil
}

func (r *Repository) ConfirmLearningPointResearchMaterial(ctx context.Context, userID, planetID, pointID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin research material confirm tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	var materialCount int
	if err := tx.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM course_point_blocks WHERE course_point_id = $1)
			+
			(SELECT COUNT(*) FROM course_point_attachments WHERE course_point_id = $1 AND source_context = 'research_material')
	`, pointID).Scan(&materialCount); err != nil {
		return uuid.Nil, fmt.Errorf("count research materials: %w", err)
	}
	if materialCount == 0 {
		return uuid.Nil, errLearningPointMaterialNotReady
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "research_material_confirmed", map[string]any{"material_count": materialCount}); err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after research material confirm: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after research material confirm: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit research material confirm tx: %w", err)
	}
	return courseID, nil
}

func (r *Repository) UnconfirmLearningPointResearchMaterial(ctx context.Context, userID, planetID, pointID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin research material unconfirm tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "research_material_unconfirmed", map[string]any{"action": "unconfirmed"}); err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after research material unconfirm: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after research material unconfirm: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit research material unconfirm tx: %w", err)
	}
	return courseID, nil
}
