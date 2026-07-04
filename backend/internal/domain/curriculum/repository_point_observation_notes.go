package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	observationNoteTypeCoreSummary       = "core_summary"
	observationNoteTypeRevisitPart       = "revisit_part"
	observationNoteTypeReferenceMaterial = "reference_material"
)

func normalizeObservationNoteType(value string) (string, bool) {
	switch strings.TrimSpace(value) {
	case observationNoteTypeCoreSummary:
		return observationNoteTypeCoreSummary, true
	case observationNoteTypeRevisitPart:
		return observationNoteTypeRevisitPart, true
	case observationNoteTypeReferenceMaterial:
		return observationNoteTypeReferenceMaterial, true
	default:
		return "", false
	}
}

func scanCoursePointObservationNote(scanner interface {
	Scan(dest ...any) error
}) (CoursePointObservationNote, error) {
	var note CoursePointObservationNote
	err := scanner.Scan(
		&note.ID,
		&note.CoursePointID,
		&note.UserID,
		&note.NoteType,
		&note.Content,
		&note.OrderIndex,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	return note, err
}

func (r *Repository) CreateLearningPointObservationNote(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointObservationNoteRequest) (uuid.UUID, CoursePointObservationNote, error) {
	noteType, ok := normalizeObservationNoteType(req.NoteType)
	content := strings.TrimSpace(req.Content)
	if !ok || content == "" {
		return uuid.Nil, CoursePointObservationNote{}, errLearningPointObservationInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("begin learning point observation tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, err
	}

	var nextOrder int
	if err := tx.QueryRow(ctx, `
        SELECT COALESCE(MAX(order_index), -1) + 1
        FROM course_point_observation_notes
        WHERE course_point_id = $1
    `, pointID).Scan(&nextOrder); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("get next observation note order: %w", err)
	}

	note, err := scanCoursePointObservationNote(tx.QueryRow(ctx, `
        INSERT INTO course_point_observation_notes (course_point_id, user_id, note_type, content, order_index)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, course_point_id, user_id, note_type, content, order_index, created_at, updated_at
    `, pointID, userID, noteType, content, nextOrder))
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("insert learning point observation note: %w", err)
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "observation_note_saved", map[string]any{"note_id": note.ID, "note_type": note.NoteType}); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, err
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("touch learning draft after observation note create: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("touch learning course after observation note create: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("commit learning point observation tx: %w", err)
	}
	return courseID, note, nil
}

func (r *Repository) UpdateLearningPointObservationNote(ctx context.Context, userID, planetID, pointID, noteID uuid.UUID, req UpdateLearningPointObservationNoteRequest) (uuid.UUID, CoursePointObservationNote, error) {
	if req.NoteType == nil && req.Content == nil {
		return uuid.Nil, CoursePointObservationNote{}, errLearningPointObservationInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("begin update learning point observation tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, err
	}

	var currentType string
	var currentContent string
	if err := tx.QueryRow(ctx, `
        SELECT note_type, content
        FROM course_point_observation_notes
        WHERE id = $1 AND course_point_id = $2 AND user_id = $3
        FOR UPDATE
    `, noteID, pointID, userID).Scan(&currentType, &currentContent); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, errLearningPointNotFound
	}

	noteType := currentType
	if req.NoteType != nil {
		normalized, ok := normalizeObservationNoteType(*req.NoteType)
		if !ok {
			return uuid.Nil, CoursePointObservationNote{}, errLearningPointObservationInvalidInput
		}
		noteType = normalized
	}
	content := strings.TrimSpace(currentContent)
	if req.Content != nil {
		content = strings.TrimSpace(*req.Content)
	}
	if content == "" {
		return uuid.Nil, CoursePointObservationNote{}, errLearningPointObservationInvalidInput
	}

	note, err := scanCoursePointObservationNote(tx.QueryRow(ctx, `
        UPDATE course_point_observation_notes
        SET note_type = $1,
            content = $2,
            updated_at = NOW()
        WHERE id = $3 AND course_point_id = $4 AND user_id = $5
        RETURNING id, course_point_id, user_id, note_type, content, order_index, created_at, updated_at
    `, noteType, content, noteID, pointID, userID))
	if err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("update learning point observation note: %w", err)
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "observation_note_saved", map[string]any{"note_id": note.ID, "note_type": note.NoteType, "action": "updated"}); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, err
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("touch learning draft after observation note update: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("touch learning course after observation note update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointObservationNote{}, fmt.Errorf("commit update learning point observation tx: %w", err)
	}
	return courseID, note, nil
}

func (r *Repository) DeleteLearningPointObservationNote(ctx context.Context, userID, planetID, pointID, noteID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin delete learning point observation tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	var noteType string
	if err := tx.QueryRow(ctx, `
        DELETE FROM course_point_observation_notes
        WHERE id = $1 AND course_point_id = $2 AND user_id = $3
        RETURNING note_type
    `, noteID, pointID, userID).Scan(&noteType); err != nil {
		return uuid.Nil, errLearningPointNotFound
	}

	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "observation_note_saved", map[string]any{"note_id": noteID, "note_type": noteType, "action": "deleted"}); err != nil {
		return uuid.Nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after observation note delete: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after observation note delete: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit delete learning point observation tx: %w", err)
	}
	return courseID, nil
}
