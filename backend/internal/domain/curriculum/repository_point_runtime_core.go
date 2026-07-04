package curriculum

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UpdateLearningLessonRuntime(ctx context.Context, userID, planetID, lessonID uuid.UUID, req UpdateLearningLessonRuntimeRequest) (uuid.UUID, *CourseLessonRuntimeEntry, error) {
	if req.Status != LessonRuntimeInProgress && req.Status != LessonRuntimeCompleted {
		return uuid.Nil, nil, errLessonRuntimeInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("begin lesson runtime tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var draftID uuid.UUID
	var courseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id, confirmed_course_id
		FROM course_drafts
		WHERE user_id = $1
		  AND status = 'learning'
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, planetID).Scan(&draftID, &courseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil, errLearningLessonNotFound
		}
		return uuid.Nil, nil, fmt.Errorf("get learning draft for lesson runtime: %w", err)
	}
	if courseID == nil {
		return uuid.Nil, nil, errLearningLessonNotFound
	}

	var courseLessonID uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT cl.id
		FROM course_lessons cl
		WHERE cl.course_id = $1
		  AND cl.id = $2
		LIMIT 1
	`, *courseID, lessonID).Scan(&courseLessonID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil, errLearningLessonNotFound
		}
		return uuid.Nil, nil, fmt.Errorf("get learning lesson: %w", err)
	}

	var entry CourseLessonRuntimeEntry
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_lesson_runtime_entries (
			user_id, course_id, course_lesson_id, status, started_at, completed_at
		) VALUES (
			$1, $2, $3, $4, NOW(), CASE WHEN $4 = 'completed' THEN NOW() ELSE NULL END
		)
		ON CONFLICT (user_id, course_lesson_id) DO UPDATE
		SET status = EXCLUDED.status,
		    started_at = COALESCE(course_lesson_runtime_entries.started_at, NOW()),
		    completed_at = CASE
		      WHEN EXCLUDED.status = 'completed' THEN COALESCE(course_lesson_runtime_entries.completed_at, NOW())
		      ELSE NULL
		    END,
		    updated_at = NOW()
		RETURNING id, user_id, course_id, course_lesson_id, status, started_at, completed_at, created_at, updated_at
	`, userID, *courseID, courseLessonID, req.Status).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.CourseID,
		&entry.CourseLessonID,
		&entry.Status,
		&entry.StartedAt,
		&entry.CompletedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	); err != nil {
		return uuid.Nil, nil, fmt.Errorf("upsert lesson runtime entry: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, nil, fmt.Errorf("touch learning draft updated_at: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, *courseID, userID); err != nil {
		return uuid.Nil, nil, fmt.Errorf("touch learning course updated_at: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, nil, fmt.Errorf("commit lesson runtime tx: %w", err)
	}

	return *courseID, &entry, nil
}

func (r *Repository) UpdateLearningPointRuntime(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpdateLearningPointRuntimeRequest) (uuid.UUID, error) {
	if req.Status != LessonRuntimeInProgress && req.Status != LessonRuntimeCompleted {
		return uuid.Nil, errPointRuntimeInvalidInput
	}

	nextPointStatus := PointStatusLearning
	if req.Status == LessonRuntimeCompleted {
		nextPointStatus = PointStatusCompleted
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin point runtime tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var draftID uuid.UUID
	var courseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id, confirmed_course_id
		FROM course_drafts
		WHERE user_id = $1
		  AND status = 'learning'
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, planetID).Scan(&draftID, &courseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, fmt.Errorf("get learning draft for point runtime: %w", err)
	}
	if courseID == nil {
		return uuid.Nil, errLearningPointNotFound
	}

	var pointType PointType
	if err := tx.QueryRow(ctx, `
		SELECT point_type
		FROM course_points
		WHERE id = $1
		  AND course_id = $2
		LIMIT 1
	`, pointID, *courseID).Scan(&pointType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, fmt.Errorf("get learning point type for runtime: %w", err)
	}
	if pointType == PointTypeResearch {
		confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
		if err != nil {
			return uuid.Nil, err
		}
		if !confirmed {
			return uuid.Nil, errLearningPointMaterialNotReady
		}
	}

	if req.Status == LessonRuntimeCompleted {
		questions, err := r.getLearningPointCompletionQuestionsTx(ctx, tx, *courseID, pointID)
		if err != nil {
			return uuid.Nil, err
		}
		selfEvaluation, err := r.getLearningPointSelfEvaluationTx(ctx, tx, *courseID, pointID)
		if err != nil {
			return uuid.Nil, err
		}
		if !learningPointCompletionReady(questions, selfEvaluation) {
			return uuid.Nil, errLearningPointCompletionNotReady
		}
	}

	tag, err := tx.Exec(ctx, `
		UPDATE course_points
		SET status = $3,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_id = $2
		  AND status IN ('ready', 'learning', 'completed')
	`, pointID, *courseID, nextPointStatus)
	if err != nil {
		return uuid.Nil, fmt.Errorf("update learning point runtime: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return uuid.Nil, errLearningPointNotFound
	}
	eventType := "point_started"
	if nextPointStatus == PointStatusCompleted {
		eventType = "point_completed"
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, eventType, map[string]any{"status": nextPointStatus}); err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft updated_at after point runtime: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, *courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course updated_at after point runtime: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit point runtime tx: %w", err)
	}

	return *courseID, nil
}
