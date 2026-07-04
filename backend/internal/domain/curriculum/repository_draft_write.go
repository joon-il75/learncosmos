package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UpdateDraft(ctx context.Context, aggregate *DraftAggregate) error {
	if aggregate == nil {
		return fmt.Errorf("draft aggregate is nil")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE course_drafts
		SET learning_goal = $1,
		    goal_profile_id = $2,
		    goal_profile_version = $3,
		    current_level = $4,
		    duration_weeks = $5,
		    study_hours_per_week = $6,
		    preferred_format = $7,
		    title = $8,
		    description = $9,
		    completion_criteria = $10,
		    planet_type_id = COALESCE($11, planet_type_id),
		    planet_texture_map_id = COALESCE($12, planet_texture_map_id),
		    updated_at = NOW()
		WHERE id = $13 AND user_id = $14
	`,
		aggregate.Draft.LearningGoal,
		aggregate.Draft.GoalProfileID,
		aggregate.Draft.GoalProfileVersion,
		aggregate.Draft.CurrentLevel,
		aggregate.Draft.DurationWeeks,
		aggregate.Draft.StudyHoursPerWeek,
		aggregate.Draft.PreferredFormat,
		aggregate.Draft.Title,
		aggregate.Draft.Description,
		aggregate.Draft.CompletionCriteria,
		aggregate.Draft.PlanetTypeID,
		aggregate.Draft.PlanetTextureMapID,
		aggregate.Draft.ID,
		aggregate.Draft.UserID,
	)
	if err != nil {
		return fmt.Errorf("update draft: %w", err)
	}

	var confirmedCourseID *uuid.UUID
	var draftStatus DraftStatus
	if err := tx.QueryRow(ctx, `
		SELECT confirmed_course_id, status
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
	`, aggregate.Draft.ID, aggregate.Draft.UserID).Scan(&confirmedCourseID, &draftStatus); err != nil {
		return fmt.Errorf("load draft sync context: %w", err)
	}

	if aggregate.Draft.PlanetTypeID != nil && confirmedCourseID != nil &&
		(draftStatus == DraftStatusConfirmed || draftStatus == DraftStatusLearning) {
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET planet_type_id = $1,
			    updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, aggregate.Draft.PlanetTypeID, *confirmedCourseID, aggregate.Draft.UserID); err != nil {
			return fmt.Errorf("sync confirmed course planet type: %w", err)
		}
	}
	if aggregate.Draft.PlanetTextureMapID != nil && confirmedCourseID != nil &&
		(draftStatus == DraftStatusConfirmed || draftStatus == DraftStatusLearning) {
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET planet_texture_map_id = $1,
			    updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, aggregate.Draft.PlanetTextureMapID, *confirmedCourseID, aggregate.Draft.UserID); err != nil {
			return fmt.Errorf("sync confirmed course planet texture map: %w", err)
		}
	}

	payload, err := json.Marshal(map[string]any{
		"title":             aggregate.Draft.Title,
		"main_lesson_count": len(aggregate.Lessons),
	})
	if err != nil {
		return fmt.Errorf("marshal update event payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, aggregate.Draft.UserID, aggregate.Draft.ID, EventCurriculumEdited, aggregate.Draft.SourceQuery, string(payload))
	if err != nil {
		return fmt.Errorf("insert draft update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update draft: %w", err)
	}
	return nil
}

func (r *Repository) SyncDraftGoalContext(
	ctx context.Context,
	draftID uuid.UUID,
	userID uuid.UUID,
	learningGoal string,
	goalProfileID uuid.UUID,
	goalProfileVersion int,
) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_drafts
		SET learning_goal = $1,
		    goal_profile_id = $2,
		    goal_profile_version = $3,
		    updated_at = NOW()
		WHERE id = $4 AND user_id = $5
	`,
		learningGoal,
		goalProfileID,
		goalProfileVersion,
		draftID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("sync draft goal context: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("draft not found")
	}
	return nil
}

func (r *Repository) DeleteDraft(ctx context.Context, draftID, userID uuid.UUID, sourceQuery, title string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var confirmedCourseID *uuid.UUID
	var status DraftStatus
	if err := tx.QueryRow(ctx, `
		SELECT confirmed_course_id, status
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, draftID, userID).Scan(&confirmedCourseID, &status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("draft_not_deletable")
		}
		return fmt.Errorf("get draft for delete: %w", err)
	}
	if status == DraftStatusArchived {
		return fmt.Errorf("draft_not_deletable")
	}

	if confirmedCourseID != nil {
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET status = 'archived',
			    updated_at = NOW()
			WHERE id = $1 AND user_id = $2
		`, *confirmedCourseID, userID); err != nil {
			return fmt.Errorf("archive confirmed course: %w", err)
		}
	}

	tag, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET status = 'archived',
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status <> 'archived'
	`, draftID, userID)
	if err != nil {
		return fmt.Errorf("archive draft: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("draft_not_deletable")
	}

	payload, err := json.Marshal(map[string]any{
		"deleted_draft_id": draftID.String(),
		"title":            title,
	})
	if err != nil {
		return fmt.Errorf("marshal delete event payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4::jsonb)
	`, userID, EventCurriculumDeleted, sourceQuery, string(payload))
	if err != nil {
		return fmt.Errorf("insert delete event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delete draft: %w", err)
	}
	return nil
}

func (r *Repository) ActivateDraft(ctx context.Context, draftID, userID uuid.UUID, sourceQuery, title string) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin activate tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var confirmedCourseID *uuid.UUID
	var status DraftStatus
	if err := tx.QueryRow(ctx, `
		SELECT confirmed_course_id, status
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, draftID, userID).Scan(&confirmedCourseID, &status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("draft_not_found")
		}
		return "", fmt.Errorf("get draft for activate: %w", err)
	}
	if status != DraftStatusArchived {
		return "", fmt.Errorf("draft_not_inactive")
	}

	destination := fmt.Sprintf("/dashboard/course-drafts/%s?section=planning", draftID.String())
	if confirmedCourseID != nil {
		var courseStatus CourseStatus
		if err := tx.QueryRow(ctx, `
			SELECT status
			FROM courses
			WHERE id = $1 AND user_id = $2
			FOR UPDATE
		`, *confirmedCourseID, userID).Scan(&courseStatus); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", fmt.Errorf("draft_not_inactive")
			}
			return "", fmt.Errorf("get course for activate: %w", err)
		}
		if courseStatus != CourseStatusArchived {
			return "", fmt.Errorf("draft_not_inactive")
		}
		if _, err := tx.Exec(ctx, `
			UPDATE courses
			SET status = 'active',
			    updated_at = NOW()
			WHERE id = $1 AND user_id = $2
		`, *confirmedCourseID, userID); err != nil {
			return "", fmt.Errorf("activate confirmed course: %w", err)
		}
		destination = fmt.Sprintf("/dashboard/planets/learning/%s", confirmedCourseID.String())
	}

	nextStatus := DraftStatusLearning
	if confirmedCourseID == nil {
		nextStatus = DraftStatusConfirmed
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET status = $3,
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID, nextStatus); err != nil {
		return "", fmt.Errorf("activate draft: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"activated_draft_id": draftID.String(),
		"title":              title,
	})
	if err != nil {
		return "", fmt.Errorf("marshal activate event payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumEdited, sourceQuery, string(payload))
	if err != nil {
		return "", fmt.Errorf("insert activate event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit activate draft: %w", err)
	}
	return destination, nil
}

func insertConfirmedCourseTx(ctx context.Context, tx pgx.Tx, confirmed *CourseAggregate, draftID, userID uuid.UUID) error {
	if confirmed == nil {
		return fmt.Errorf("confirmed course aggregate is nil")
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO courses (
			id, user_id, source_draft_id, source_query, learning_goal, current_level,
			duration_weeks, study_hours_per_week, preferred_format, title, description, completion_criteria, status, planet_type_id, planet_texture_map_id
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`,
		confirmed.Course.ID,
		confirmed.Course.UserID,
		confirmed.Course.SourceDraftID,
		confirmed.Course.SourceQuery,
		confirmed.Course.LearningGoal,
		confirmed.Course.CurrentLevel,
		confirmed.Course.DurationWeeks,
		confirmed.Course.StudyHoursPerWeek,
		confirmed.Course.PreferredFormat,
		confirmed.Course.Title,
		confirmed.Course.Description,
		confirmed.Course.CompletionCriteria,
		confirmed.Course.Status,
		confirmed.Course.PlanetTypeID,
		confirmed.Course.PlanetTextureMapID,
	)
	if err != nil {
		return fmt.Errorf("insert confirmed course: %w", err)
	}

	insertCourseLesson := func(lesson CourseLesson) error {
		searchSpecJSON, marshalErr := MarshalLessonRecommendationSearchSpec(lesson.RecommendationSearchSpec)
		if marshalErr != nil {
			return fmt.Errorf("marshal course lesson recommendation search spec: %w", marshalErr)
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO course_lessons (
				id, course_id, parent_lesson_id, title, objective, summary,
				difficulty_level, lesson_role, source_type, order_index, recommendation_search_spec
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		`,
			lesson.ID,
			confirmed.Course.ID,
			lesson.ParentLessonID,
			lesson.Title,
			lesson.Objective,
			lesson.Summary,
			lesson.DifficultyLevel,
			lesson.LessonRole,
			lesson.SourceType,
			lesson.OrderIndex,
			searchSpecJSON,
		)
		return err
	}

	insertCoursePoint := func(lessonID uuid.UUID, point CoursePointAggregate) error {
		priceType := PriceTypeFree
		if point.Point.PriceType != nil {
			priceType = *point.Point.PriceType
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO course_points (
				id, course_id, course_lesson_id, point_type,
				content_id, external_url, title, description, thumbnail_url,
				price_type, rank_score, template_type, status, order_index
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 'ready', $13)
		`,
			point.Point.ID,
			confirmed.Course.ID,
			lessonID,
			point.Point.PointType,
			point.Point.ContentID,
			point.Point.ExternalURL,
			point.Point.Title,
			point.Point.Description,
			point.Point.ThumbnailURL,
			priceType,
			point.Point.RankScore,
			point.Point.TemplateType,
			point.Point.OrderIndex,
		)
		return err
	}

	insertCoursePointEntries := func(point CoursePointAggregate) error {
		for _, block := range point.Blocks {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_blocks (
					id, course_point_id, user_id, block_type, content, order_index
				) VALUES ($1, $2, $3, $4, $5::jsonb, $6)
			`, block.ID, point.Point.ID, block.UserID, block.BlockType, string(block.Content), block.OrderIndex); err != nil {
				return fmt.Errorf("insert confirmed point block: %w", err)
			}
		}
		if point.AISummary != nil {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_ai_summaries (
					course_point_id, source_title, source_description, summary, learning_language
				) VALUES ($1, $2, $3, $4, $5)
			`, point.Point.ID, point.AISummary.SourceTitle, point.AISummary.SourceDescription, point.AISummary.Summary, normalizeLearningLanguage(point.AISummary.LearningLanguage)); err != nil {
				return fmt.Errorf("insert confirmed point ai summary: %w", err)
			}
		}
		if point.JournalEntry != nil {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_journal_entries (
					course_point_id, observation, reflection, next_step
				) VALUES ($1, $2, $3, $4)
			`, point.Point.ID, point.JournalEntry.Observation, point.JournalEntry.Reflection, point.JournalEntry.NextStep); err != nil {
				return fmt.Errorf("insert confirmed point journal entry: %w", err)
			}
		}
		if point.RecordEntry != nil {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_record_entries (
					course_point_id, study_minutes, practice_count, confidence_level, application_note
				) VALUES ($1, $2, $3, $4, $5)
			`, point.Point.ID, point.RecordEntry.StudyMinutes, point.RecordEntry.PracticeCount, point.RecordEntry.ConfidenceLevel, point.RecordEntry.ApplicationNote); err != nil {
				return fmt.Errorf("insert confirmed point record entry: %w", err)
			}
		}
		if point.ArtifactEntry != nil {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_artifact_entries (
					course_point_id, artifact_type, title, url, description
				) VALUES ($1, $2, $3, $4, $5)
			`, point.Point.ID, point.ArtifactEntry.ArtifactType, point.ArtifactEntry.Title, point.ArtifactEntry.URL, point.ArtifactEntry.Description); err != nil {
				return fmt.Errorf("insert confirmed point artifact entry: %w", err)
			}
		}
		for _, artifact := range point.Artifacts {
			if _, err := tx.Exec(ctx, `
				INSERT INTO course_point_artifacts (
					id, course_point_id, artifact_type, title, url, description, order_index
				) VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (id) DO NOTHING
			`, artifact.ID, point.Point.ID, artifact.ArtifactType, artifact.Title, artifact.URL, artifact.Description, artifact.OrderIndex); err != nil {
				return fmt.Errorf("insert confirmed point artifact: %w", err)
			}
		}
		return nil
	}

	for _, mainTree := range confirmed.Lessons {
		if err := insertCourseLesson(mainTree.Lesson); err != nil {
			return fmt.Errorf("insert confirmed main lesson: %w", err)
		}
		for _, point := range mainTree.Points {
			if err := insertCoursePoint(mainTree.Lesson.ID, point); err != nil {
				return fmt.Errorf("insert confirmed main lesson point: %w", err)
			}
			if err := insertCoursePointEntries(point); err != nil {
				return err
			}
		}
		for _, subTree := range mainTree.SubLessons {
			if err := insertCourseLesson(subTree.Lesson); err != nil {
				return fmt.Errorf("insert confirmed sub lesson: %w", err)
			}
			for _, point := range subTree.Points {
				if err := insertCoursePoint(subTree.Lesson.ID, point); err != nil {
					return fmt.Errorf("insert confirmed sub lesson point: %w", err)
				}
				if err := insertCoursePointEntries(point); err != nil {
					return err
				}
			}
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE course_drafts
		SET status = 'confirmed', confirmed_course_id = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`, confirmed.Course.ID, draftID, userID)
	if err != nil {
		return fmt.Errorf("update draft confirmed state: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"course_id": confirmed.Course.ID.String(),
		"title":     confirmed.Course.Title,
	})
	if err != nil {
		return fmt.Errorf("marshal confirm event payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumConfirmed, confirmed.Course.SourceQuery, string(payload))
	if err != nil {
		return fmt.Errorf("insert confirm event: %w", err)
	}

	return nil
}

func setDraftLearningStateTx(ctx context.Context, tx pgx.Tx, userID, draftID, courseID uuid.UUID, sourceQuery string) error {
	_, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET status = 'learning', updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID)
	if err != nil {
		return fmt.Errorf("update draft learning state: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"planet_id":           draftID.String(),
		"confirmed_course_id": courseID.String(),
	})
	if err != nil {
		return fmt.Errorf("marshal start event payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumStarted, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert start learning event: %w", err)
	}

	return nil
}
