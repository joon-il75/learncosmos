package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) StartLearningPlanet(ctx context.Context, userID, planetID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin start learning tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var draftID uuid.UUID
	var status DraftStatus
	var confirmedCourseID *uuid.UUID
	var sourceQuery string
	if err := tx.QueryRow(ctx, `
		SELECT id, status, confirmed_course_id, source_query
		FROM course_drafts
		WHERE user_id = $1
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, planetID).Scan(&draftID, &status, &confirmedCourseID, &sourceQuery); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftResourceNotFound
		}
		return fmt.Errorf("get draft for start learning: %w", err)
	}

	switch status {
	case DraftStatusLearning:
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit start learning noop: %w", err)
		}
		return nil
	case DraftStatusConfirmed:
	default:
		return errDraftStartNotAllowed
	}

	if confirmedCourseID != nil {
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET updated_at = NOW()
			WHERE id = $1 AND user_id = $2
		`, *confirmedCourseID, userID); err != nil {
			return fmt.Errorf("touch confirmed course updated_at: %w", err)
		}
	}

	if confirmedCourseID == nil {
		return fmt.Errorf("confirmed course missing for start learning")
	}
	if err := setDraftLearningStateTx(ctx, tx, userID, draftID, *confirmedCourseID, sourceQuery); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit start learning tx: %w", err)
	}
	return nil
}

func (r *Repository) CompleteLearningPlanet(ctx context.Context, userID, planetID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin complete learning tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var draftID uuid.UUID
	var status DraftStatus
	var confirmedCourseID *uuid.UUID
	var sourceQuery string
	if err := tx.QueryRow(ctx, `
		SELECT id, status, confirmed_course_id, source_query
		FROM course_drafts
		WHERE user_id = $1
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, planetID).Scan(&draftID, &status, &confirmedCourseID, &sourceQuery); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftResourceNotFound
		}
		return fmt.Errorf("get draft for complete learning: %w", err)
	}

	if status != DraftStatusLearning {
		return errDraftCompleteNotAllowed
	}

	progressSummary, err := r.loadExplorerProgressSummary(ctx, tx, draftID, status, confirmedCourseID)
	if err != nil {
		return fmt.Errorf("load planet progress for complete: %w", err)
	}
	if progressSummary == nil || !progressSummary.CanComplete {
		return errDraftCompleteNotAllowed
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET status = 'archived', updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("update draft archived state: %w", err)
	}

	if confirmedCourseID != nil {
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET updated_at = NOW()
			WHERE id = $1 AND user_id = $2
		`, *confirmedCourseID, userID); err != nil {
			return fmt.Errorf("touch confirmed course on complete: %w", err)
		}
	}

	payload, err := json.Marshal(map[string]any{
		"planet_id": draftID.String(),
		"confirmed_course_id": func() any {
			if confirmedCourseID == nil {
				return nil
			}
			return confirmedCourseID.String()
		}(),
	})
	if err != nil {
		return fmt.Errorf("marshal complete event payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumCompleted, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert complete learning event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit complete learning tx: %w", err)
	}
	return nil
}

func isCourseLessonCompleted(lesson CourseLessonTree) bool {
	if lesson.RuntimeEntry != nil {
		return lesson.RuntimeEntry.Status == LessonRuntimeCompleted
	}
	summary := ""
	if lesson.Lesson.Summary != nil {
		summary = strings.TrimSpace(*lesson.Lesson.Summary)
	}
	return strings.Contains(summary, "[학습완료]")
}

// isDraftMainLessonCompleted returns true when a main lesson is fully complete.
// Completion requires at least one point or sub-lesson to exist, and:
// - All direct points on the main lesson (if any) must be completed.
// - All sub-lessons (if any) must have all their points completed.
func isDraftMainLessonCompleted(main DraftLessonTree) bool {
	hasAny := false
	if len(main.Points) > 0 {
		hasAny = true
		for _, p := range main.Points {
			if p.Point.Status != PointStatusCompleted {
				return false
			}
		}
	}
	if len(main.SubLessons) > 0 {
		hasAny = true
		for _, sub := range main.SubLessons {
			if !isDraftLessonPointsCompleted(sub) {
				return false
			}
		}
	}
	return hasAny
}

func isDraftLessonPointsCompleted(lesson DraftLessonTree) bool {
	if len(lesson.Points) == 0 {
		return false
	}
	for _, p := range lesson.Points {
		if p.Point.Status != PointStatusCompleted {
			return false
		}
	}
	return true
}

func countDraftLessonPoints(lesson DraftLessonTree) (int, int) {
	total := len(lesson.Points)
	completed := 0
	for _, p := range lesson.Points {
		if p.Point.Status == PointStatusCompleted {
			completed++
		}
	}
	return total, completed
}

func countDraftPointsInMainLesson(main DraftLessonTree) (int, int) {
	total, completed := countDraftLessonPoints(main)
	for _, sub := range main.SubLessons {
		subTotal, subCompleted := countDraftLessonPoints(sub)
		total += subTotal
		completed += subCompleted
	}
	return total, completed
}

func countDraftCompletableRegionsInMainLesson(main DraftLessonTree) (int, int) {
	total := 1
	completed := 0
	if isDraftMainLessonCompleted(main) {
		completed++
	}
	return total, completed
}

func (r *Repository) TouchPlanetLastAccessed(ctx context.Context, userID, planetID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_drafts
		SET last_accessed_at = NOW()
		WHERE user_id = $1
		  AND (id = $2 OR confirmed_course_id = $2)
	`, userID, planetID)
	if err != nil {
		return fmt.Errorf("touch planet last accessed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errDraftResourceNotFound
	}
	return nil
}

// isCourseMainLessonCompleted returns true when a main lesson is fully complete.
// Completion requires at least one point or sub-lesson to exist, and:
// - All direct points on the main lesson (if any) must be completed.
// - All sub-lessons (if any) must have all their points completed.
func isCourseMainLessonCompleted(main CourseLessonTree) bool {
	hasAny := false
	if len(main.Points) > 0 {
		hasAny = true
		for _, p := range main.Points {
			if p.Point.Status != PointStatusCompleted {
				return false
			}
		}
	}
	if len(main.SubLessons) > 0 {
		hasAny = true
		for _, sub := range main.SubLessons {
			if !isCourseLessonPointsCompleted(sub) {
				return false
			}
		}
	}
	return hasAny
}

func isCourseLessonPointsCompleted(lesson CourseLessonTree) bool {
	if len(lesson.Points) == 0 {
		return false
	}
	for _, p := range lesson.Points {
		if p.Point.Status != PointStatusCompleted {
			return false
		}
	}
	return true
}

func countCourseLessonPoints(lesson CourseLessonTree) (int, int) {
	total := len(lesson.Points)
	completed := 0
	for _, p := range lesson.Points {
		if p.Point.Status == PointStatusCompleted {
			completed++
		}
	}
	return total, completed
}

func countCoursePointsInMainLesson(main CourseLessonTree) (int, int) {
	total, completed := countCourseLessonPoints(main)
	for _, sub := range main.SubLessons {
		subTotal, subCompleted := countCourseLessonPoints(sub)
		total += subTotal
		completed += subCompleted
	}
	return total, completed
}

func countCourseCompletableRegionsInMainLesson(main CourseLessonTree) (int, int) {
	total := 1
	completed := 0
	if isCourseMainLessonCompleted(main) {
		completed++
	}
	return total, completed
}
