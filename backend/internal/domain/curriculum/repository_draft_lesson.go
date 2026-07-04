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

func (r *Repository) UpdateDraftLessonDetail(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftLessonRequest) error {
	nextTitle := req.Title
	if nextTitle != nil {
		trimmed := strings.TrimSpace(*nextTitle)
		if trimmed == "" {
			return errDraftLessonInvalidInput
		}
		nextTitle = &trimmed
	}

	nextObjective := req.Objective
	if nextObjective != nil {
		trimmed := strings.TrimSpace(*nextObjective)
		nextObjective = &trimmed
	}

	nextSummary := req.Summary
	if nextSummary != nil {
		trimmed := strings.TrimSpace(*nextSummary)
		nextSummary = &trimmed
	}

	nextDifficultyLevel := req.DifficultyLevel
	if nextDifficultyLevel != nil {
		trimmed := strings.TrimSpace(*nextDifficultyLevel)
		nextDifficultyLevel = &trimmed
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update draft lesson tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	var generationLanguage string
	var currentTitle string
	var learningGoal *string
	var currentObjective *string
	var orderIndex int
	err = tx.QueryRow(ctx, `
		SELECT cd.source_query, cd.learning_goal, COALESCE(cd.generation_language, 'ko'),
		       cdl.title, cdl.objective, cdl.order_index
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(false)).Scan(
		&sourceQuery,
		&learningGoal,
		&generationLanguage,
		&currentTitle,
		&currentObjective,
		&orderIndex,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftLessonNotFound
		}
		return fmt.Errorf("get draft lesson for detail update: %w", err)
	}

	var nextSearchSpecJSON *string
	if nextTitle != nil || nextObjective != nil {
		effectiveTitle := currentTitle
		if nextTitle != nil {
			effectiveTitle = *nextTitle
		}
		effectiveObjective := currentObjective
		if nextObjective != nil {
			effectiveObjective = nextObjective
		}
		searchSpec := BuildFallbackLessonRecommendationSearchSpecForLesson(sourceQuery, learningGoal, generationLanguage, effectiveTitle, effectiveObjective, orderIndex)
		searchSpecJSON, marshalErr := MarshalLessonRecommendationSearchSpec(searchSpec)
		if marshalErr != nil {
			return fmt.Errorf("marshal updated lesson recommendation search spec: %w", marshalErr)
		}
		nextSearchSpecJSON = &searchSpecJSON
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_draft_lessons
		SET title = COALESCE($1, title),
		    objective = COALESCE($2, objective),
		    summary = COALESCE($3, summary),
		    difficulty_level = COALESCE($4, difficulty_level),
		    recommendation_search_spec = COALESCE($5::jsonb, recommendation_search_spec),
		    updated_at = NOW()
		WHERE id = $6
	`, nextTitle, nextObjective, nextSummary, nextDifficultyLevel, nextSearchSpecJSON, lessonID); err != nil {
		return fmt.Errorf("update draft lesson detail: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("touch draft after lesson detail update: %w", err)
	}

	eventTitle := currentTitle
	if nextTitle != nil {
		eventTitle = *nextTitle
	}
	payload, err := json.Marshal(map[string]any{
		"lesson_id": lessonID.String(),
		"title":     eventTitle,
	})
	if err != nil {
		return fmt.Errorf("marshal lesson detail update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert lesson detail update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update draft lesson tx: %w", err)
	}
	return nil
}

func (r *Repository) UpdateDraftMainLessonDetail(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftMainLessonRequest) error {
	nextTitle := req.Title
	if nextTitle != nil {
		trimmed := strings.TrimSpace(*nextTitle)
		if trimmed == "" {
			return errDraftLevelInvalidInput
		}
		nextTitle = &trimmed
	}

	nextDescription := req.Description
	if nextDescription != nil {
		trimmed := strings.TrimSpace(*nextDescription)
		nextDescription = &trimmed
	}

	nextObjective := req.Objective
	if nextObjective != nil {
		trimmed := strings.TrimSpace(*nextObjective)
		nextObjective = &trimmed
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update main lesson detail tx: %w", err)
	}
	defer tx.Rollback(ctx)

	sourceQuery, currentTitle, err := fetchDraftMainLessonContext(ctx, tx, draftID, userID, lessonID, true)
	if err != nil {
		return err
	}

	var generationLanguage string
	var learningGoal *string
	var currentObjective *string
	var orderIndex int
	if err := tx.QueryRow(ctx, `
		SELECT cd.learning_goal, COALESCE(cd.generation_language, 'ko'), cdl.objective, cdl.order_index
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cdl.course_draft_id = $2
		  AND cd.user_id = $3
	`, lessonID, draftID, userID).Scan(&learningGoal, &generationLanguage, &currentObjective, &orderIndex); err != nil {
		return fmt.Errorf("get draft main lesson search spec context: %w", err)
	}

	var nextSearchSpecJSON *string
	if nextTitle != nil || nextObjective != nil {
		effectiveTitle := currentTitle
		if nextTitle != nil {
			effectiveTitle = *nextTitle
		}
		effectiveObjective := currentObjective
		if nextObjective != nil {
			effectiveObjective = nextObjective
		}
		searchSpec := BuildFallbackLessonRecommendationSearchSpecForLesson(sourceQuery, learningGoal, generationLanguage, effectiveTitle, effectiveObjective, orderIndex)
		searchSpecJSON, marshalErr := MarshalLessonRecommendationSearchSpec(searchSpec)
		if marshalErr != nil {
			return fmt.Errorf("marshal updated main lesson recommendation search spec: %w", marshalErr)
		}
		nextSearchSpecJSON = &searchSpecJSON
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_draft_lessons
		SET title = COALESCE($1, title),
		    summary = COALESCE($2, summary),
		    objective = COALESCE($3, objective),
		    recommendation_search_spec = COALESCE($4::jsonb, recommendation_search_spec),
		    updated_at = NOW()
		WHERE id = $5
	`, nextTitle, nextDescription, nextObjective, nextSearchSpecJSON, lessonID); err != nil {
		return fmt.Errorf("update main lesson detail: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("touch draft after main lesson detail update: %w", err)
	}

	eventTitle := currentTitle
	if nextTitle != nil {
		eventTitle = *nextTitle
	}
	payload, err := json.Marshal(map[string]any{
		"main_lesson_id": lessonID.String(),
		"title":          eventTitle,
	})
	if err != nil {
		return fmt.Errorf("marshal main lesson detail update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert main lesson detail update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update main lesson detail tx: %w", err)
	}
	return nil
}

func (r *Repository) UpdateDraftMainLessonMemo(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftDetailMemoRequest) (DraftDetailMemo, error) {
	if req.Note == nil {
		return DraftDetailMemo{}, errDraftMemoInvalidInput
	}
	note := strings.TrimSpace(*req.Note)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("begin update main lesson memo tx: %w", err)
	}
	defer tx.Rollback(ctx)

	sourceQuery, currentTitle, err := fetchDraftMainLessonContext(ctx, tx, draftID, userID, lessonID, true)
	if err != nil {
		return DraftDetailMemo{}, err
	}

	var memo DraftDetailMemo
	err = tx.QueryRow(ctx, `
		INSERT INTO course_draft_detail_notes (
			course_draft_id, course_draft_lesson_id, target_type, note
		) VALUES ($1, $2, 'lesson', $3)
		ON CONFLICT (course_draft_lesson_id) WHERE target_type = 'lesson'
		DO UPDATE SET note = EXCLUDED.note, updated_at = NOW()
		RETURNING course_draft_lesson_id, note, updated_at
	`, draftID, lessonID, note).Scan(&memo.LessonID, &memo.Note, &memo.UpdatedAt)
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("upsert main lesson memo: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("touch draft after main lesson memo update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"main_lesson_id": lessonID.String(),
		"title":          currentTitle,
		"note_length":    len(note),
	})
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("marshal main lesson memo update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("insert main lesson memo update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("commit update main lesson memo tx: %w", err)
	}
	return memo, nil
}

func (r *Repository) UpdateDraftLessonMemo(ctx context.Context, draftID, userID, lessonID uuid.UUID, req UpdateDraftDetailMemoRequest) (DraftDetailMemo, error) {
	if req.Note == nil {
		return DraftDetailMemo{}, errDraftMemoInvalidInput
	}
	note := strings.TrimSpace(*req.Note)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("begin update draft lesson memo tx: %w", err)
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
			return DraftDetailMemo{}, errDraftLessonNotFound
		}
		return DraftDetailMemo{}, fmt.Errorf("get draft lesson for memo update: %w", err)
	}

	var memo DraftDetailMemo
	err = tx.QueryRow(ctx, `
		INSERT INTO course_draft_detail_notes (
			course_draft_id, course_draft_lesson_id, target_type, note
		) VALUES ($1, $2, 'lesson', $3)
		ON CONFLICT (course_draft_lesson_id) WHERE target_type = 'lesson'
		DO UPDATE SET note = EXCLUDED.note, updated_at = NOW()
		RETURNING course_draft_lesson_id, note, updated_at
	`, draftID, lessonID, note).Scan(&memo.LessonID, &memo.Note, &memo.UpdatedAt)
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("upsert draft lesson memo: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("touch draft after lesson memo update: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"lesson_id":   lessonID.String(),
		"title":       currentTitle,
		"note_length": len(note),
	})
	if err != nil {
		return DraftDetailMemo{}, fmt.Errorf("marshal lesson memo update payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, userID, draftID, lessonID, EventCurriculumEdited, sourceQuery, string(payload)); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("insert lesson memo update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DraftDetailMemo{}, fmt.Errorf("commit update draft lesson memo tx: %w", err)
	}
	return memo, nil
}
