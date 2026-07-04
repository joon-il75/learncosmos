package curriculum

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func isLearningPointStartEvent(eventType string, payload map[string]any) bool {
	if strings.TrimSpace(eventType) == "" {
		return false
	}
	if action, ok := payload["action"].(string); ok && strings.EqualFold(strings.TrimSpace(action), "deleted") {
		return false
	}
	switch strings.TrimSpace(eventType) {
	case "journal_saved",
		"record_saved",
		"artifact_submitted",
		"practice_log_added",
		"question_added",
		"question_answered",
		"self_evaluation_saved",
		"attachment_added",
		"research_material_confirmed":
		return true
	default:
		return false
	}
}

func (r *Repository) buildPointEventContextSnapshotTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID) (map[string]any, error) {
	var (
		courseID          uuid.UUID
		courseTitle       string
		lessonID          uuid.UUID
		lessonTitle       string
		parentLessonID    *uuid.UUID
		parentLessonTitle *string
		pointTitle        string
		pointType         string
	)
	if err := tx.QueryRow(ctx, `
		SELECT c.id, c.title,
		       lesson.id, lesson.title,
		       parent_lesson.id, parent_lesson.title,
		       cp.title, cp.point_type
		FROM course_points cp
		JOIN courses c ON c.id = cp.course_id
		JOIN course_lessons lesson ON lesson.id = cp.course_lesson_id
		LEFT JOIN course_lessons parent_lesson ON parent_lesson.id = lesson.parent_lesson_id
		WHERE cp.id = $1
		LIMIT 1
	`, pointID).Scan(
		&courseID,
		&courseTitle,
		&lessonID,
		&lessonTitle,
		&parentLessonID,
		&parentLessonTitle,
		&pointTitle,
		&pointType,
	); err != nil {
		return nil, fmt.Errorf("build point event context snapshot: %w", err)
	}

	snapshot := map[string]any{
		"course_id":    courseID.String(),
		"course_title": courseTitle,
		"lesson_id":    lessonID.String(),
		"lesson_title": lessonTitle,
		"point_id":     pointID.String(),
		"point_title":  pointTitle,
		"point_type":   pointType,
	}
	if parentLessonID != nil {
		snapshot["parent_lesson_id"] = parentLessonID.String()
	}
	if parentLessonTitle != nil {
		snapshot["parent_lesson_title"] = *parentLessonTitle
	}
	return snapshot, nil
}

func (r *Repository) attachPointEventContextSnapshotTx(ctx context.Context, tx pgx.Tx, pointID uuid.UUID, payload map[string]any) error {
	if payload == nil {
		return nil
	}
	if _, exists := payload["context_snapshot"]; exists {
		return nil
	}
	snapshot, err := r.buildPointEventContextSnapshotTx(ctx, tx, pointID)
	if err != nil {
		return err
	}
	payload["context_snapshot"] = snapshot
	return nil
}

func (r *Repository) markLearningPointStartedByEventTx(ctx context.Context, tx pgx.Tx, userID, pointID uuid.UUID, triggerEvent string) error {
	tag, err := tx.Exec(ctx, `
		UPDATE course_points
		SET status = 'learning',
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'ready'
	`, pointID)
	if err != nil {
		return fmt.Errorf("mark learning point started: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	payload := map[string]any{
		"status":        PointStatusLearning,
		"trigger_event": triggerEvent,
	}
	if err := r.attachPointEventContextSnapshotTx(ctx, tx, pointID, payload); err != nil {
		return err
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal point started payload: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO course_point_events (course_point_id, user_id, event_type, event_payload)
		VALUES ($1, $2, 'point_started', $3)
	`, pointID, userID, payloadBytes); err != nil {
		return fmt.Errorf("insert point started event: %w", err)
	}
	return nil
}

func (r *Repository) insertPointEventTx(ctx context.Context, tx pgx.Tx, userID, pointID uuid.UUID, eventType string, payload map[string]any) error {
	if strings.TrimSpace(eventType) == "" {
		return nil
	}
	if payload == nil {
		payload = map[string]any{}
	}
	if isLearningPointStartEvent(eventType, payload) {
		if err := r.markLearningPointStartedByEventTx(ctx, tx, userID, pointID, strings.TrimSpace(eventType)); err != nil {
			return err
		}
	}
	if err := r.attachPointEventContextSnapshotTx(ctx, tx, pointID, payload); err != nil {
		return err
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal point event payload: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO course_point_events (course_point_id, user_id, event_type, event_payload)
		VALUES ($1, $2, $3, $4)
	`, pointID, userID, strings.TrimSpace(eventType), payloadBytes); err != nil {
		return fmt.Errorf("insert point event: %w", err)
	}
	return nil
}

func (r *Repository) ListLearningPointEvents(ctx context.Context, userID, planetID, pointID uuid.UUID) ([]CoursePointEvent, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin list point events tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT cpe.id, cpe.course_point_id, cpe.user_id, cpe.event_type, cpe.event_payload, cpe.created_at
		FROM course_point_events cpe
		JOIN course_points cp ON cp.id = cpe.course_point_id
		WHERE cpe.course_point_id = $1
		  AND cp.course_id = $2
		ORDER BY cpe.created_at DESC
		LIMIT 50
	`, pointID, courseID)
	if err != nil {
		return nil, fmt.Errorf("list learning point events: %w", err)
	}
	defer rows.Close()

	reversed := make([]CoursePointEvent, 0)
	for rows.Next() {
		var event CoursePointEvent
		if err := rows.Scan(
			&event.ID,
			&event.CoursePointID,
			&event.UserID,
			&event.EventType,
			&event.EventPayload,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan learning point event: %w", err)
		}
		reversed = append(reversed, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning point events: %w", err)
	}
	events := make([]CoursePointEvent, 0, len(reversed))
	for i := len(reversed) - 1; i >= 0; i-- {
		events = append(events, reversed[i])
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit list point events tx: %w", err)
	}
	return events, nil
}
