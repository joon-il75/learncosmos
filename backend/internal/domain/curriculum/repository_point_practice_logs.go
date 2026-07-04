package curriculum

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateLearningPointPracticeLog(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointPracticeLogRequest) (uuid.UUID, CoursePointPracticeLog, error) {
	title, activityName, attemptCount, successCount, failureCount, durationMinutes, blockedPart, changedMethod, achievementNote, achievement, nextPlan, nextPractice, err := normalizePracticeLogInput(
		req.Title,
		req.ActivityName,
		req.AttemptCount,
		req.SuccessCount,
		req.FailureCount,
		req.DurationMinutes,
		req.BlockedPart,
		req.ChangedMethod,
		req.AchievementNote,
		req.Achievement,
		req.NextPlan,
		req.NextPractice,
	)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("begin create point practice log tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	log, err := scanCoursePointPracticeLog(tx.QueryRow(ctx, `
		INSERT INTO course_point_practice_logs (
			course_point_id, user_id, title, attempt_count, success_count, failure_count,
			duration_minutes, blocked_part, changed_method, achievement_note, next_plan,
			activity_name, achievement, next_practice
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, course_point_id, user_id, title, activity_name, attempt_count, success_count, failure_count,
		          duration_minutes, blocked_part, changed_method, achievement_note, achievement, next_plan, next_practice, created_at, updated_at
	`, pointID, userID, title, attemptCount, successCount, failureCount, durationMinutes, blockedPart, changedMethod, achievementNote, nextPlan, activityName, achievement, nextPractice))
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("insert learning point practice log: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "practice_log_added", map[string]any{"practice_log_id": log.ID, "title": log.Title}); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("touch learning draft after point practice log create: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("touch learning course after point practice log create: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("commit create point practice log tx: %w", err)
	}
	return courseID, log, nil
}

func (r *Repository) UpdateLearningPointPracticeLog(ctx context.Context, userID, planetID, pointID, logID uuid.UUID, req UpdateLearningPointPracticeLogRequest) (uuid.UUID, CoursePointPracticeLog, error) {
	title, activityName, attemptCount, successCount, failureCount, durationMinutes, blockedPart, changedMethod, achievementNote, achievement, nextPlan, nextPractice, err := normalizePracticeLogInput(
		req.Title,
		req.ActivityName,
		req.AttemptCount,
		req.SuccessCount,
		req.FailureCount,
		req.DurationMinutes,
		req.BlockedPart,
		req.ChangedMethod,
		req.AchievementNote,
		req.Achievement,
		req.NextPlan,
		req.NextPractice,
	)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("begin update point practice log tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	log, err := scanCoursePointPracticeLog(tx.QueryRow(ctx, `
		UPDATE course_point_practice_logs
		SET title = $5,
		    attempt_count = $6,
		    success_count = $7,
		    failure_count = $8,
		    duration_minutes = $9,
		    blocked_part = $10,
		    changed_method = $11,
		    achievement_note = $12,
		    next_plan = $13,
		    activity_name = $14,
		    achievement = $15,
		    next_practice = $16,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $4
		  )
		RETURNING id, course_point_id, user_id, title, activity_name, attempt_count, success_count, failure_count,
		          duration_minutes, blocked_part, changed_method, achievement_note, achievement, next_plan, next_practice, created_at, updated_at
	`, logID, pointID, userID, courseID, title, attemptCount, successCount, failureCount, durationMinutes, blockedPart, changedMethod, achievementNote, nextPlan, activityName, achievement, nextPractice))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointPracticeLog{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("update learning point practice log: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "practice_log_added", map[string]any{"practice_log_id": log.ID, "action": "updated", "title": log.Title}); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("touch learning draft after point practice log update: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("touch learning course after point practice log update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointPracticeLog{}, fmt.Errorf("commit update point practice log tx: %w", err)
	}
	return courseID, log, nil
}

func (r *Repository) DeleteLearningPointPracticeLog(ctx context.Context, userID, planetID, pointID, logID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin delete point practice log tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	var deletedTitle string
	if err := tx.QueryRow(ctx, `
		DELETE FROM course_point_practice_logs
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $4
		  )
		RETURNING title
	`, logID, pointID, userID, courseID).Scan(&deletedTitle); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, fmt.Errorf("delete learning point practice log: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "practice_log_added", map[string]any{"practice_log_id": logID, "action": "deleted", "title": deletedTitle}); err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after point practice log delete: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after point practice log delete: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit delete point practice log tx: %w", err)
	}
	return courseID, nil
}
