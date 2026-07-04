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

func (r *Repository) UpdateDraftLessonJournal(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftLessonJournalRequest) (DraftLessonJournal, error) {
	if req.Observation == nil || req.Reflection == nil || req.NextStep == nil {
		return DraftLessonJournal{}, errDraftJournalInvalidInput
	}

	observation := strings.TrimSpace(*req.Observation)
	reflection := strings.TrimSpace(*req.Reflection)
	nextStep := strings.TrimSpace(*req.NextStep)
	if observation == "" && reflection == "" && nextStep == "" {
		return DraftLessonJournal{}, errDraftJournalInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DraftLessonJournal{}, fmt.Errorf("begin update draft lesson journal tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	var currentTitle string
	err = tx.QueryRow(ctx, `
		SELECT cd.source_query, cdl.title
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(true)).Scan(&sourceQuery, &currentTitle)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DraftLessonJournal{}, errDraftLessonNotFound
		}
		return DraftLessonJournal{}, fmt.Errorf("get draft lesson for journal update: %w", err)
	}

	var journal DraftLessonJournal
	err = tx.QueryRow(ctx, `
		INSERT INTO course_draft_journal_entries (
			course_draft_id, course_draft_lesson_id, observation, reflection, next_step
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_draft_lesson_id)
		DO UPDATE SET
			observation = EXCLUDED.observation,
			reflection = EXCLUDED.reflection,
			next_step = EXCLUDED.next_step,
			updated_at = NOW()
		RETURNING course_draft_lesson_id, observation, reflection, next_step, updated_at
	`, draftID, lessonID, observation, reflection, nextStep).Scan(
		&journal.TargetID,
		&journal.Observation,
		&journal.Reflection,
		&journal.NextStep,
		&journal.UpdatedAt,
	)
	if err != nil {
		return DraftLessonJournal{}, fmt.Errorf("upsert draft lesson journal: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return DraftLessonJournal{}, fmt.Errorf("touch draft after lesson journal update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"lesson_id":          lessonID.String(),
		"title":              currentTitle,
		"observation_length": len(observation),
		"reflection_length":  len(reflection),
		"next_step_length":   len(nextStep),
	})
	if err != nil {
		return DraftLessonJournal{}, fmt.Errorf("marshal lesson journal update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return DraftLessonJournal{}, fmt.Errorf("insert lesson journal update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DraftLessonJournal{}, fmt.Errorf("commit update draft lesson journal tx: %w", err)
	}
	return journal, nil
}

func (r *Repository) UpdateDraftLessonRecord(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftLessonRecordRequest) (DraftLessonRecord, error) {
	if req.StudyMinutes == nil || req.PracticeCount == nil || req.ConfidenceLevel == nil || req.ApplicationNote == nil {
		return DraftLessonRecord{}, errDraftRecordInvalidInput
	}

	studyMinutes := *req.StudyMinutes
	practiceCount := *req.PracticeCount
	confidenceLevel := *req.ConfidenceLevel
	applicationNote := strings.TrimSpace(*req.ApplicationNote)
	if studyMinutes < 0 || practiceCount < 0 || confidenceLevel < 1 || confidenceLevel > 5 {
		return DraftLessonRecord{}, errDraftRecordInvalidInput
	}
	if studyMinutes == 0 && practiceCount == 0 && confidenceLevel == 3 && applicationNote == "" {
		return DraftLessonRecord{}, errDraftRecordInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DraftLessonRecord{}, fmt.Errorf("begin update draft lesson record tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	var currentTitle string
	err = tx.QueryRow(ctx, `
		SELECT cd.source_query, cdl.title
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(true)).Scan(&sourceQuery, &currentTitle)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DraftLessonRecord{}, errDraftLessonNotFound
		}
		return DraftLessonRecord{}, fmt.Errorf("get draft lesson for record update: %w", err)
	}

	var record DraftLessonRecord
	err = tx.QueryRow(ctx, `
		INSERT INTO course_draft_record_entries (
			course_draft_id, course_draft_lesson_id, study_minutes, practice_count, confidence_level, application_note
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (course_draft_lesson_id)
		DO UPDATE SET
			study_minutes = EXCLUDED.study_minutes,
			practice_count = EXCLUDED.practice_count,
			confidence_level = EXCLUDED.confidence_level,
			application_note = EXCLUDED.application_note,
			updated_at = NOW()
		RETURNING course_draft_lesson_id, study_minutes, practice_count, confidence_level, application_note, updated_at
	`, draftID, lessonID, studyMinutes, practiceCount, confidenceLevel, applicationNote).Scan(
		&record.TargetID,
		&record.StudyMinutes,
		&record.PracticeCount,
		&record.ConfidenceLevel,
		&record.ApplicationNote,
		&record.UpdatedAt,
	)
	if err != nil {
		return DraftLessonRecord{}, fmt.Errorf("upsert draft lesson record: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return DraftLessonRecord{}, fmt.Errorf("touch draft after lesson record update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"lesson_id":        lessonID.String(),
		"title":            currentTitle,
		"study_minutes":    studyMinutes,
		"practice_count":   practiceCount,
		"confidence_level": confidenceLevel,
		"application_note": len(applicationNote),
	})
	if err != nil {
		return DraftLessonRecord{}, fmt.Errorf("marshal lesson record update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return DraftLessonRecord{}, fmt.Errorf("insert lesson record update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DraftLessonRecord{}, fmt.Errorf("commit update draft lesson record tx: %w", err)
	}
	return record, nil
}

func (r *Repository) UpdateDraftLessonArtifact(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftLessonArtifactRequest) (DraftLessonArtifact, error) {
	if req.ArtifactType == nil || req.Title == nil || req.URL == nil || req.Description == nil {
		return DraftLessonArtifact{}, errDraftArtifactInvalidInput
	}

	artifactType := strings.TrimSpace(*req.ArtifactType)
	title := strings.TrimSpace(*req.Title)
	url := strings.TrimSpace(*req.URL)
	description := sanitizePointArtifactDescription(*req.Description)
	if title == "" {
		return DraftLessonArtifact{}, errDraftArtifactInvalidInput
	}
	if artifactType == "" {
		artifactType = "note"
	}
	if url != "" && !strings.HasPrefix(strings.ToLower(url), "http://") && !strings.HasPrefix(strings.ToLower(url), "https://") {
		return DraftLessonArtifact{}, errDraftArtifactInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("begin update draft lesson artifact tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	var currentTitle string
	err = tx.QueryRow(ctx, `
		SELECT cd.source_query, cdl.title
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(true)).Scan(&sourceQuery, &currentTitle)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DraftLessonArtifact{}, errDraftLessonNotFound
		}
		return DraftLessonArtifact{}, fmt.Errorf("get draft lesson for artifact update: %w", err)
	}

	var artifact DraftLessonArtifact
	err = tx.QueryRow(ctx, `
		INSERT INTO course_draft_artifact_entries (
			course_draft_id, course_draft_lesson_id, artifact_type, title, url, description
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (course_draft_lesson_id)
		DO UPDATE SET
			artifact_type = EXCLUDED.artifact_type,
			title = EXCLUDED.title,
			url = EXCLUDED.url,
			description = EXCLUDED.description,
			updated_at = NOW()
		RETURNING course_draft_lesson_id, artifact_type, title, url, description, updated_at
	`, draftID, lessonID, artifactType, title, url, description).Scan(
		&artifact.TargetID,
		&artifact.ArtifactType,
		&artifact.Title,
		&artifact.URL,
		&artifact.Description,
		&artifact.UpdatedAt,
	)
	if err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("upsert draft lesson artifact: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("touch draft after lesson artifact update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"lesson_id":          lessonID.String(),
		"title":              currentTitle,
		"artifact_type":      artifactType,
		"artifact_title":     title,
		"url_present":        url != "",
		"description_length": len(description),
	})
	if err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("marshal lesson artifact update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("insert lesson artifact update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DraftLessonArtifact{}, fmt.Errorf("commit update draft lesson artifact tx: %w", err)
	}
	return artifact, nil
}
