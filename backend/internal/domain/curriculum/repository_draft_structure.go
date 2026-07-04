package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UpdateDraftStructure(ctx context.Context, draftID, userID uuid.UUID, req UpdateDraftStructureRequest) error {
	if len(req.MainLessons) == 0 {
		return errDraftStructureInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update draft structure tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	if err := tx.QueryRow(ctx, `
		SELECT source_query
		FROM course_drafts
		WHERE id = $1
		  AND user_id = $2
		  AND status = ANY($3)
		FOR UPDATE
	`, draftID, userID, getEditableDraftStatuses(false)).Scan(&sourceQuery); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftStructureNotFound
		}
		return fmt.Errorf("get draft for structure update: %w", err)
	}

	mainLessonRows, err := tx.Query(ctx, `
		SELECT id
		FROM course_draft_lessons
		WHERE course_draft_id = $1 AND parent_lesson_id IS NULL
		ORDER BY order_index, created_at
		FOR UPDATE
	`, draftID)
	if err != nil {
		return fmt.Errorf("list draft main lessons for structure update: %w", err)
	}
	defer mainLessonRows.Close()

	existingMainLessons := make(map[uuid.UUID]struct{})
	for mainLessonRows.Next() {
		var lessonID uuid.UUID
		if err := mainLessonRows.Scan(&lessonID); err != nil {
			return fmt.Errorf("scan draft main lesson for structure update: %w", err)
		}
		existingMainLessons[lessonID] = struct{}{}
	}
	if err := mainLessonRows.Err(); err != nil {
		return fmt.Errorf("iterate draft main lessons for structure update: %w", err)
	}
	mainLessonRows.Close()
	if len(req.MainLessons) != len(existingMainLessons) {
		return errDraftStructureInvalidInput
	}

	subLessonRows, err := tx.Query(ctx, `
		SELECT id
		FROM course_draft_lessons
		WHERE course_draft_id = $1 AND parent_lesson_id IS NOT NULL
		ORDER BY order_index, created_at
		FOR UPDATE
	`, draftID)
	if err != nil {
		return fmt.Errorf("list draft sub lessons for structure update: %w", err)
	}
	defer subLessonRows.Close()

	existingSubLessons := make(map[uuid.UUID]struct{})
	for subLessonRows.Next() {
		var lessonID uuid.UUID
		if err := subLessonRows.Scan(&lessonID); err != nil {
			return fmt.Errorf("scan draft sub lesson for structure update: %w", err)
		}
		existingSubLessons[lessonID] = struct{}{}
	}
	if err := subLessonRows.Err(); err != nil {
		return fmt.Errorf("iterate draft sub lessons for structure update: %w", err)
	}
	subLessonRows.Close()

	pointRows, err := tx.Query(ctx, `
		SELECT id
		FROM course_draft_points
		WHERE course_draft_id = $1
		  AND point_type = 'research'
		ORDER BY order_index, created_at
		FOR UPDATE
	`, draftID)
	if err != nil {
		return fmt.Errorf("list draft points for structure update: %w", err)
	}
	defer pointRows.Close()

	existingPoints := make(map[uuid.UUID]struct{})
	for pointRows.Next() {
		var pointID uuid.UUID
		if err := pointRows.Scan(&pointID); err != nil {
			return fmt.Errorf("scan draft point for structure update: %w", err)
		}
		existingPoints[pointID] = struct{}{}
	}
	if err := pointRows.Err(); err != nil {
		return fmt.Errorf("iterate draft points for structure update: %w", err)
	}
	pointRows.Close()

	seenMainLessons := make(map[uuid.UUID]struct{}, len(req.MainLessons))
	seenMainOrders := make(map[int]struct{}, len(req.MainLessons))
	seenSubLessons := make(map[uuid.UUID]struct{}, len(existingSubLessons))
	seenPoints := make(map[uuid.UUID]struct{}, len(existingPoints))
	for _, mainLesson := range req.MainLessons {
		if mainLesson.ID == uuid.Nil || mainLesson.OrderIndex < 0 {
			return errDraftStructureInvalidInput
		}
		if _, ok := existingMainLessons[mainLesson.ID]; !ok {
			return errDraftStructureInvalidInput
		}
		if _, exists := seenMainLessons[mainLesson.ID]; exists {
			return errDraftStructureInvalidInput
		}
		seenMainLessons[mainLesson.ID] = struct{}{}
		if _, exists := seenMainOrders[mainLesson.OrderIndex]; exists {
			return errDraftStructureInvalidInput
		}
		seenMainOrders[mainLesson.OrderIndex] = struct{}{}

		seenSubOrders := make(map[int]struct{}, len(mainLesson.Lessons))
		for _, sub := range mainLesson.Lessons {
			if sub.ID == uuid.Nil || sub.OrderIndex < 0 {
				return errDraftStructureInvalidInput
			}
			if _, ok := existingSubLessons[sub.ID]; !ok {
				return errDraftStructureInvalidInput
			}
			if _, exists := seenSubLessons[sub.ID]; exists {
				return errDraftStructureInvalidInput
			}
			seenSubLessons[sub.ID] = struct{}{}
			if _, exists := seenSubOrders[sub.OrderIndex]; exists {
				return errDraftStructureInvalidInput
			}
			seenSubOrders[sub.OrderIndex] = struct{}{}
		}

		seenPointOrders := make(map[int]struct{}, len(mainLesson.Points))
		for _, point := range mainLesson.StructurePoints() {
			if point.ID == uuid.Nil || point.OrderIndex < 0 {
				return errDraftStructureInvalidInput
			}
			if _, ok := existingPoints[point.ID]; !ok {
				return errDraftStructureInvalidInput
			}
			if _, exists := seenPoints[point.ID]; exists {
				return errDraftStructureInvalidInput
			}
			seenPoints[point.ID] = struct{}{}
			if _, exists := seenPointOrders[point.OrderIndex]; exists {
				return errDraftStructureInvalidInput
			}
			seenPointOrders[point.OrderIndex] = struct{}{}
		}
	}
	if len(seenMainLessons) != len(existingMainLessons) ||
		len(seenSubLessons) != len(existingSubLessons) ||
		len(seenPoints) != len(existingPoints) {
		return errDraftStructureInvalidInput
	}

	// Temporarily set negative order_index to avoid unique constraint conflicts during reorder
	tempMainIndex := -1
	for lessonID := range existingMainLessons {
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_lessons SET order_index = $1, updated_at = NOW() WHERE id = $2
		`, tempMainIndex, lessonID); err != nil {
			return fmt.Errorf("set temporary main lesson order: %w", err)
		}
		tempMainIndex--
	}

	tempSubIndex := -1
	for lessonID := range existingSubLessons {
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_lessons SET order_index = $1, updated_at = NOW() WHERE id = $2
		`, tempSubIndex, lessonID); err != nil {
			return fmt.Errorf("set temporary sub lesson order: %w", err)
		}
		tempSubIndex--
	}

	tempPointIndex := -1
	for pointID := range existingPoints {
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_points SET order_index = $1, updated_at = NOW() WHERE id = $2
		`, tempPointIndex, pointID); err != nil {
			return fmt.Errorf("set temporary draft point order: %w", err)
		}
		tempPointIndex--
	}

	movedLessons := 0
	movedPoints := 0
	for _, mainLesson := range req.MainLessons {
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_lessons
			SET order_index = $1, updated_at = NOW()
			WHERE id = $2
		`, mainLesson.OrderIndex, mainLesson.ID); err != nil {
			return fmt.Errorf("update main lesson structure: %w", err)
		}

		for _, sub := range mainLesson.Lessons {
			if _, err := tx.Exec(ctx, `
				UPDATE course_draft_lessons
				SET parent_lesson_id = $1, order_index = $2, updated_at = NOW()
				WHERE id = $3
			`, mainLesson.ID, sub.OrderIndex, sub.ID); err != nil {
				return fmt.Errorf("update sub lesson structure: %w", err)
			}
			movedLessons++
		}

		for _, point := range mainLesson.StructurePoints() {
			if _, err := tx.Exec(ctx, `
				UPDATE course_draft_points
				SET course_draft_lesson_id = $1, order_index = $2, updated_at = NOW()
				WHERE id = $3
			`, mainLesson.ID, point.OrderIndex, point.ID); err != nil {
				return fmt.Errorf("update draft point structure: %w", err)
			}
			movedPoints++
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("touch draft after structure update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"main_lesson_count": len(req.MainLessons),
		"lesson_count":      movedLessons,
		"point_count":       movedPoints,
	})
	if err != nil {
		return fmt.Errorf("marshal draft structure update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert draft structure update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update draft structure tx: %w", err)
	}
	return nil
}
