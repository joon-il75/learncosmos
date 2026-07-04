package curriculum

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func deleteConfirmedCourseForRebuildTx(ctx context.Context, tx pgx.Tx, courseID uuid.UUID, userID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		DELETE FROM courses
		WHERE id = $1 AND user_id = $2
	`, courseID, userID); err != nil {
		return fmt.Errorf("delete confirmed course for rebuild: %w", err)
	}
	return nil
}

func deleteCoursePlanDataForRebuildTx(ctx context.Context, tx pgx.Tx, courseID uuid.UUID, label string) error {
	queries := []struct {
		sql string
		err string
	}{
		{`DELETE FROM course_point_blocks WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`, "delete course point blocks for %s: %w"},
		{`DELETE FROM course_point_ai_summaries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`, "delete course point ai summaries for %s: %w"},
		{`DELETE FROM course_point_journal_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`, "delete course point journal entries for %s: %w"},
		{`DELETE FROM course_point_record_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`, "delete course point record entries for %s: %w"},
		{`DELETE FROM course_point_artifact_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`, "delete course point artifact entries for %s: %w"},
		{`DELETE FROM course_lesson_runtime_entries WHERE course_id = $1`, "delete course runtime entries for %s: %w"},
		{`DELETE FROM course_points WHERE course_id = $1`, "delete course points for %s: %w"},
		{`DELETE FROM course_lessons WHERE course_id = $1`, "delete course lessons for %s: %w"},
	}
	for _, query := range queries {
		if _, err := tx.Exec(ctx, query.sql, courseID); err != nil {
			return fmt.Errorf(query.err, label, err)
		}
	}
	return nil
}

func deleteDraftPlanDataForRebuildTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, label string) error {
	queries := []struct {
		sql string
		err string
	}{
		{`DELETE FROM course_draft_point_blocks WHERE course_draft_point_id IN (SELECT id FROM course_draft_points WHERE course_draft_id = $1)`, "delete draft point blocks for %s: %w"},
		{`DELETE FROM course_draft_points WHERE course_draft_id = $1`, "delete draft points for %s: %w"},
		{`DELETE FROM course_draft_journal_entries WHERE course_draft_id = $1`, "delete draft journal entries for %s: %w"},
		{`DELETE FROM course_draft_record_entries WHERE course_draft_id = $1`, "delete draft record entries for %s: %w"},
		{`DELETE FROM course_draft_artifact_entries WHERE course_draft_id = $1`, "delete draft artifact entries for %s: %w"},
		{`DELETE FROM course_draft_detail_notes WHERE course_draft_id = $1`, "delete draft detail notes for %s: %w"},
		{`DELETE FROM course_draft_lessons WHERE course_draft_id = $1`, "delete draft lessons for %s: %w"},
	}
	for _, query := range queries {
		if _, err := tx.Exec(ctx, query.sql, draftID); err != nil {
			return fmt.Errorf(query.err, label, err)
		}
	}
	return nil
}

func markExplorerPlanInactiveTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		UPDATE explorer_nodes
		SET status = 'inactive',
		    updated_at = NOW()
		WHERE parent_kind = 'subregion'
		  AND parent_id IN (
			SELECT es.id
			FROM explorer_subregions es
			JOIN explorer_regions er ON er.id = es.region_id
			WHERE er.course_draft_id = $1
		  )
	`, draftID); err != nil {
		return fmt.Errorf("deactivate explorer subregion nodes: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE explorer_nodes
		SET status = 'inactive',
		    updated_at = NOW()
		WHERE parent_kind = 'region'
		  AND parent_id IN (
			SELECT id FROM explorer_regions WHERE course_draft_id = $1
		  )
	`, draftID); err != nil {
		return fmt.Errorf("deactivate explorer region nodes: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE explorer_subregions
		SET status = 'inactive',
		    updated_at = NOW()
		WHERE region_id IN (
			SELECT id FROM explorer_regions WHERE course_draft_id = $1
		)
	`, draftID); err != nil {
		return fmt.Errorf("deactivate explorer subregions: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE explorer_regions
		SET status = 'inactive',
		    updated_at = NOW()
		WHERE course_draft_id = $1
	`, draftID); err != nil {
		return fmt.Errorf("deactivate explorer regions: %w", err)
	}

	return nil
}
