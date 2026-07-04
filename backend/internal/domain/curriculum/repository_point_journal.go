package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) UpdateLearningPointGoal(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpdateLearningPointGoalRequest) (uuid.UUID, error) {
	if req.PointGoal == nil && req.PointCategory == nil {
		return uuid.Nil, errPointRuntimeInvalidInput
	}

	pointGoal := trimOptionalString(req.PointGoal)
	pointCategory := trimOptionalString(req.PointCategory)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin learning point goal tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_points
		SET point_goal = $1,
		    point_category = $2,
		    updated_at = NOW()
		WHERE id = $3
		  AND course_id = $4
	`, nullableTrimmedString(pointGoal), nullableTrimmedString(pointCategory), pointID, courseID); err != nil {
		return uuid.Nil, fmt.Errorf("update learning point goal: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "point_goal_saved", nil); err != nil {
		return uuid.Nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after point goal: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after point goal: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit learning point goal tx: %w", err)
	}
	return courseID, nil
}

func (r *Repository) UpdateLearningPointJournal(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpdateDraftLessonJournalRequest) (uuid.UUID, DraftLessonJournal, error) {
	if req.Observation == nil && req.CoreConcept == nil {
		return uuid.Nil, DraftLessonJournal{}, errLearningPointJournalInvalidInput
	}

	observation := trimOptionalString(req.Observation)
	reflection := trimOptionalString(req.Reflection)
	nextStep := trimOptionalString(req.NextStep)
	coreConcept := trimOptionalString(req.CoreConcept)
	myExplanation := trimOptionalString(req.MyExplanation)
	examples := trimOptionalString(req.Examples)
	confusedParts := trimOptionalString(req.ConfusedParts)
	referenceLinks := trimOptionalString(req.ReferenceLinks)
	if coreConcept == "" {
		coreConcept = observation
	}
	if myExplanation == "" {
		myExplanation = reflection
	}
	if confusedParts == "" {
		confusedParts = nextStep
	}
	if observation == "" {
		observation = coreConcept
	}
	if reflection == "" {
		reflection = myExplanation
	}
	if nextStep == "" {
		nextStep = confusedParts
	}
	if observation == "" && reflection == "" && nextStep == "" && coreConcept == "" && myExplanation == "" && examples == "" && confusedParts == "" && referenceLinks == "" {
		return uuid.Nil, DraftLessonJournal{}, errLearningPointJournalInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, DraftLessonJournal{}, fmt.Errorf("begin learning point journal tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, DraftLessonJournal{}, err
	}

	var journal DraftLessonJournal
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_journal_entries (
			course_point_id, observation, reflection, next_step,
			core_concept, my_explanation, examples, confused_parts, reference_links
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (course_point_id)
		DO UPDATE SET
			observation = EXCLUDED.observation,
			reflection = EXCLUDED.reflection,
			next_step = EXCLUDED.next_step,
			core_concept = EXCLUDED.core_concept,
			my_explanation = EXCLUDED.my_explanation,
			examples = EXCLUDED.examples,
			confused_parts = EXCLUDED.confused_parts,
			reference_links = EXCLUDED.reference_links,
			updated_at = NOW()
		RETURNING course_point_id, observation, reflection, next_step,
		          core_concept, my_explanation, examples, confused_parts, reference_links, updated_at
	`, pointID, observation, reflection, nextStep, coreConcept, myExplanation, examples, confusedParts, referenceLinks).Scan(
		&journal.TargetID,
		&journal.Observation,
		&journal.Reflection,
		&journal.NextStep,
		&journal.CoreConcept,
		&journal.MyExplanation,
		&journal.Examples,
		&journal.ConfusedParts,
		&journal.ReferenceLinks,
		&journal.UpdatedAt,
	); err != nil {
		return uuid.Nil, DraftLessonJournal{}, fmt.Errorf("upsert learning point journal: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "journal_saved", nil); err != nil {
		return uuid.Nil, DraftLessonJournal{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, DraftLessonJournal{}, fmt.Errorf("touch learning draft after point journal: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, DraftLessonJournal{}, fmt.Errorf("touch learning course after point journal: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, DraftLessonJournal{}, fmt.Errorf("commit learning point journal tx: %w", err)
	}
	return courseID, journal, nil
}

func (r *Repository) UpdateLearningPointRecord(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpdateDraftLessonRecordRequest) (uuid.UUID, DraftLessonRecord, error) {
	if req.StudyMinutes == nil || req.PracticeCount == nil || req.ConfidenceLevel == nil || req.ApplicationNote == nil {
		return uuid.Nil, DraftLessonRecord{}, errLearningPointRecordInvalidInput
	}

	studyMinutes := *req.StudyMinutes
	practiceCount := *req.PracticeCount
	confidenceLevel := *req.ConfidenceLevel
	applicationNote := strings.TrimSpace(*req.ApplicationNote)
	if studyMinutes < 0 || practiceCount < 0 || confidenceLevel < 1 || confidenceLevel > 5 {
		return uuid.Nil, DraftLessonRecord{}, errLearningPointRecordInvalidInput
	}
	if studyMinutes == 0 && practiceCount == 0 && confidenceLevel == 3 && applicationNote == "" {
		return uuid.Nil, DraftLessonRecord{}, errLearningPointRecordInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, DraftLessonRecord{}, fmt.Errorf("begin learning point record tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, DraftLessonRecord{}, err
	}

	var record DraftLessonRecord
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_record_entries (
			course_point_id, study_minutes, practice_count, confidence_level, application_note
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_point_id)
		DO UPDATE SET
			study_minutes = EXCLUDED.study_minutes,
			practice_count = EXCLUDED.practice_count,
			confidence_level = EXCLUDED.confidence_level,
			application_note = EXCLUDED.application_note,
			updated_at = NOW()
		RETURNING course_point_id, study_minutes, practice_count, confidence_level, application_note, updated_at
	`, pointID, studyMinutes, practiceCount, confidenceLevel, applicationNote).Scan(
		&record.TargetID,
		&record.StudyMinutes,
		&record.PracticeCount,
		&record.ConfidenceLevel,
		&record.ApplicationNote,
		&record.UpdatedAt,
	); err != nil {
		return uuid.Nil, DraftLessonRecord{}, fmt.Errorf("upsert learning point record: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "record_saved", map[string]any{
		"study_minutes":  studyMinutes,
		"practice_count": practiceCount,
	}); err != nil {
		return uuid.Nil, DraftLessonRecord{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, DraftLessonRecord{}, fmt.Errorf("touch learning draft after point record: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, DraftLessonRecord{}, fmt.Errorf("touch learning course after point record: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, DraftLessonRecord{}, fmt.Errorf("commit learning point record tx: %w", err)
	}
	return courseID, record, nil
}

func (r *Repository) UpdateLearningPointArtifact(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpdateDraftLessonArtifactRequest) (uuid.UUID, DraftLessonArtifact, error) {
	if req.ArtifactType == nil || req.Title == nil || req.URL == nil || req.Description == nil {
		return uuid.Nil, DraftLessonArtifact{}, errLearningPointArtifactInvalidInput
	}

	artifactType, title, url, description, err := normalizePointArtifactInput(*req.ArtifactType, *req.Title, *req.URL, *req.Description)
	if err != nil {
		return uuid.Nil, DraftLessonArtifact{}, errLearningPointArtifactInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, DraftLessonArtifact{}, fmt.Errorf("begin learning point artifact tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, DraftLessonArtifact{}, err
	}

	var artifact DraftLessonArtifact
	if err := tx.QueryRow(ctx, `
		WITH first_artifact AS (
			SELECT id
			FROM course_point_artifacts
			WHERE course_point_id = $1
			ORDER BY order_index ASC, created_at ASC
			LIMIT 1
		),
		updated AS (
			UPDATE course_point_artifacts
			SET artifact_type = $2,
			    title = $3,
			    url = $4,
			    description = $5,
			    point_category = COALESCE(point_category, ''),
			    production_process = COALESCE(production_process, ''),
			    learned_points = COALESCE(learned_points, ''),
			    difficult_points = COALESCE(difficult_points, ''),
			    visibility = COALESCE(NULLIF(visibility, ''), 'private'),
			    updated_at = NOW()
			WHERE id IN (SELECT id FROM first_artifact)
			RETURNING course_point_id, artifact_type, title, url, description, updated_at
		),
		inserted AS (
			INSERT INTO course_point_artifacts (
				course_point_id, artifact_type, title, url, description,
				point_category, production_process, learned_points, difficult_points, visibility, order_index
			)
			SELECT $1, $2, $3, $4, $5, '', '', '', '', 'private', 0
			WHERE NOT EXISTS (SELECT 1 FROM updated)
			RETURNING course_point_id, artifact_type, title, url, description, updated_at
		)
		SELECT course_point_id, artifact_type, title, url, description, updated_at FROM updated
		UNION ALL
		SELECT course_point_id, artifact_type, title, url, description, updated_at FROM inserted
	`, pointID, artifactType, title, url, description).Scan(
		&artifact.TargetID,
		&artifact.ArtifactType,
		&artifact.Title,
		&artifact.URL,
		&artifact.Description,
		&artifact.UpdatedAt,
	); err != nil {
		return uuid.Nil, DraftLessonArtifact{}, fmt.Errorf("upsert learning point artifact: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, DraftLessonArtifact{}, fmt.Errorf("touch learning draft after point artifact: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, DraftLessonArtifact{}, fmt.Errorf("touch learning course after point artifact: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, DraftLessonArtifact{}, fmt.Errorf("commit learning point artifact tx: %w", err)
	}
	return courseID, artifact, nil
}
