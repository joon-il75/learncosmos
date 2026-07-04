package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const maxLearningSessionDeltaSeconds = 300

func normalizeLearningSessionVisibility(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "hidden":
		return "hidden"
	case "visible":
		return "visible"
	default:
		return "visible"
	}
}

func normalizeLearningSessionEndReason(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		trimmed = "heartbeat"
	}
	if len(trimmed) > 40 {
		trimmed = trimmed[:40]
	}
	return &trimmed
}

func clampLearningSessionDelta(value int) int {
	if value < 0 {
		return 0
	}
	if value > maxLearningSessionDeltaSeconds {
		return maxLearningSessionDeltaSeconds
	}
	return value
}

func scanCoursePointLearningSession(row pgx.Row) (CoursePointLearningSession, error) {
	var session CoursePointLearningSession
	if err := row.Scan(
		&session.ID,
		&session.CoursePointID,
		&session.UserID,
		&session.StartedAt,
		&session.LastSeenAt,
		&session.EndedAt,
		&session.ActiveSeconds,
		&session.HeartbeatCount,
		&session.EndReason,
		&session.LastVisibility,
		&session.UserAgent,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return CoursePointLearningSession{}, err
	}
	return session, nil
}

func (r *Repository) CreateLearningPointSession(ctx context.Context, userID, planetID, pointID uuid.UUID, req StartLearningPointSessionRequest) (uuid.UUID, CoursePointLearningSession, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("begin create point learning session tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, err
	}

	visibility := normalizeLearningSessionVisibility(req.Visibility)
	userAgent := strings.TrimSpace(req.UserAgent)
	if userAgent == "" {
		userAgent = "unknown"
	}
	if len(userAgent) > 512 {
		userAgent = userAgent[:512]
	}

	session, err := scanCoursePointLearningSession(tx.QueryRow(ctx, `
		INSERT INTO course_point_learning_sessions (
			course_point_id, user_id, last_visibility, user_agent
		) VALUES ($1, $2, $3, $4)
		RETURNING id, course_point_id, user_id, started_at, last_seen_at, ended_at,
		          active_seconds, heartbeat_count, end_reason, last_visibility, user_agent, created_at, updated_at
	`, pointID, userID, visibility, userAgent))
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("insert point learning session: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("commit create point learning session tx: %w", err)
	}
	return courseID, session, nil
}

func (r *Repository) HeartbeatLearningPointSession(ctx context.Context, userID, planetID, pointID, sessionID uuid.UUID, req UpdateLearningPointSessionRequest) (uuid.UUID, CoursePointLearningSession, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("begin heartbeat point learning session tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, err
	}

	session, err := scanCoursePointLearningSession(tx.QueryRow(ctx, `
		UPDATE course_point_learning_sessions
		SET last_seen_at = NOW(),
		    active_seconds = active_seconds + $5,
		    heartbeat_count = heartbeat_count + 1,
		    last_visibility = $6,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $4
		  )
		RETURNING id, course_point_id, user_id, started_at, last_seen_at, ended_at,
		          active_seconds, heartbeat_count, end_reason, last_visibility, user_agent, created_at, updated_at
	`, sessionID, pointID, userID, courseID, clampLearningSessionDelta(req.ActiveSecondsDelta), normalizeLearningSessionVisibility(req.Visibility)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, CoursePointLearningSession{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("heartbeat point learning session: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("commit heartbeat point learning session tx: %w", err)
	}
	return courseID, session, nil
}

func (r *Repository) EndLearningPointSession(ctx context.Context, userID, planetID, pointID, sessionID uuid.UUID, req UpdateLearningPointSessionRequest) (uuid.UUID, CoursePointLearningSession, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("begin end point learning session tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointLearningSession{}, err
	}

	session, err := scanCoursePointLearningSession(tx.QueryRow(ctx, `
		UPDATE course_point_learning_sessions
		SET last_seen_at = NOW(),
		    ended_at = COALESCE(ended_at, NOW()),
		    active_seconds = active_seconds + $5,
		    heartbeat_count = heartbeat_count + 1,
		    end_reason = COALESCE(end_reason, $6),
		    last_visibility = $7,
		    updated_at = NOW()
		WHERE id = $1
		  AND course_point_id = $2
		  AND user_id = $3
		  AND course_point_id IN (
			SELECT id FROM course_points WHERE id = $2 AND course_id = $4
		  )
		RETURNING id, course_point_id, user_id, started_at, last_seen_at, ended_at,
		          active_seconds, heartbeat_count, end_reason, last_visibility, user_agent, created_at, updated_at
	`, sessionID, pointID, userID, courseID, clampLearningSessionDelta(req.ActiveSecondsDelta), normalizeLearningSessionEndReason(req.Reason), normalizeLearningSessionVisibility(req.Visibility)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, CoursePointLearningSession{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("end point learning session: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointLearningSession{}, fmt.Errorf("commit end point learning session tx: %w", err)
	}
	return courseID, session, nil
}
