package curriculum

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func updateCourseMetadataForRemainingRebuildTx(ctx context.Context, tx pgx.Tx, aggregate *CourseAggregate) error {
	if _, err := tx.Exec(ctx, `
		UPDATE courses
		SET source_query = $1,
		    learning_goal = $2,
		    current_level = $3,
		    duration_weeks = $4,
		    study_hours_per_week = $5,
		    preferred_format = $6,
		    title = $7,
		    description = $8,
		    completion_criteria = $9,
		    status = $10,
		    planet_type_id = $11,
		    updated_at = NOW()
		WHERE id = $12 AND user_id = $13
	`,
		aggregate.Course.SourceQuery,
		aggregate.Course.LearningGoal,
		aggregate.Course.CurrentLevel,
		aggregate.Course.DurationWeeks,
		aggregate.Course.StudyHoursPerWeek,
		aggregate.Course.PreferredFormat,
		aggregate.Course.Title,
		aggregate.Course.Description,
		aggregate.Course.CompletionCriteria,
		aggregate.Course.Status,
		aggregate.Course.PlanetTypeID,
		aggregate.Course.ID,
		aggregate.Course.UserID,
	); err != nil {
		return fmt.Errorf("update course metadata for remaining rebuild: %w", err)
	}
	return nil
}

func insertRebuiltCourseTreesTx(ctx context.Context, tx pgx.Tx, aggregate *CourseAggregate) error {
	for _, mainTree := range aggregate.Lessons {
		if err := insertRebuiltCourseLessonTx(ctx, tx, aggregate.Course.ID, mainTree.Lesson); err != nil {
			return fmt.Errorf("insert rebuilt course main lesson: %w", err)
		}
		for _, point := range mainTree.Points {
			if err := insertRebuiltCoursePointTx(ctx, tx, aggregate.Course.ID, point, "main"); err != nil {
				return err
			}
		}
		for _, subTree := range mainTree.SubLessons {
			if err := insertRebuiltCourseSubTreeTx(ctx, tx, aggregate.Course.ID, subTree); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertRebuiltCourseSubTreeTx(ctx context.Context, tx pgx.Tx, courseID uuid.UUID, subTree CourseLessonTree) error {
	if err := insertRebuiltCourseLessonTx(ctx, tx, courseID, subTree.Lesson); err != nil {
		return fmt.Errorf("insert rebuilt course sub lesson: %w", err)
	}
	if subTree.RuntimeEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_lesson_runtime_entries (
				id, user_id, course_id, course_lesson_id, status, started_at, completed_at, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			subTree.RuntimeEntry.ID,
			subTree.RuntimeEntry.UserID,
			subTree.RuntimeEntry.CourseID,
			subTree.RuntimeEntry.CourseLessonID,
			subTree.RuntimeEntry.Status,
			subTree.RuntimeEntry.StartedAt,
			subTree.RuntimeEntry.CompletedAt,
			subTree.RuntimeEntry.CreatedAt,
			subTree.RuntimeEntry.UpdatedAt,
		); err != nil {
			return fmt.Errorf("insert rebuilt course runtime entry: %w", err)
		}
	}
	for _, point := range subTree.Points {
		if err := insertRebuiltCoursePointTx(ctx, tx, courseID, point, "sub"); err != nil {
			return err
		}
	}
	return nil
}

func insertRebuiltCourseLessonTx(ctx context.Context, tx pgx.Tx, courseID uuid.UUID, lesson CourseLesson) error {
	searchSpecJSON, marshalErr := MarshalLessonRecommendationSearchSpec(lesson.RecommendationSearchSpec)
	if marshalErr != nil {
		return fmt.Errorf("marshal rebuilt course lesson recommendation search spec: %w", marshalErr)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO course_lessons (
			id, course_id, parent_lesson_id, title, objective, summary,
			difficulty_level, lesson_role, source_type, order_index, recommendation_search_spec
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
	`,
		lesson.ID,
		courseID,
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

func insertRebuiltCoursePointTx(ctx context.Context, tx pgx.Tx, courseID uuid.UUID, point CoursePointAggregate, label string) error {
	status := point.Point.Status
	if status == PointStatusDraft {
		status = PointStatusReady
	}
	priceType := PriceTypeFree
	if point.Point.PriceType != nil {
		priceType = *point.Point.PriceType
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO course_points (
			id, course_id, course_lesson_id, point_type,
			content_id, external_url, title, description, thumbnail_url,
			price_type, rank_score, template_type, status, order_index, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`,
		point.Point.ID,
		courseID,
		point.Point.CourseLessonID,
		point.Point.PointType,
		point.Point.ContentID,
		point.Point.ExternalURL,
		point.Point.Title,
		point.Point.Description,
		point.Point.ThumbnailURL,
		priceType,
		point.Point.RankScore,
		point.Point.TemplateType,
		status,
		point.Point.OrderIndex,
		point.Point.CompletedAt,
	); err != nil {
		return fmt.Errorf("insert rebuilt course %s point: %w", label, err)
	}
	return insertRebuiltCoursePointEntriesTx(ctx, tx, point)
}

func insertRebuiltCoursePointEntriesTx(ctx context.Context, tx pgx.Tx, point CoursePointAggregate) error {
	for _, block := range point.Blocks {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_point_blocks (
				id, course_point_id, user_id, block_type, content, order_index
			) VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		`, block.ID, point.Point.ID, block.UserID, block.BlockType, string(block.Content), block.OrderIndex); err != nil {
			return fmt.Errorf("insert rebuilt course point block: %w", err)
		}
	}
	if point.AISummary != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_point_ai_summaries (
				course_point_id, source_title, source_description, summary, learning_language
			) VALUES ($1, $2, $3, $4, $5)
		`, point.Point.ID, point.AISummary.SourceTitle, point.AISummary.SourceDescription, point.AISummary.Summary, normalizeLearningLanguage(point.AISummary.LearningLanguage)); err != nil {
			return fmt.Errorf("insert rebuilt course ai summary: %w", err)
		}
	}
	if point.JournalEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_point_journal_entries (
				course_point_id, observation, reflection, next_step
			) VALUES ($1, $2, $3, $4)
		`, point.Point.ID, point.JournalEntry.Observation, point.JournalEntry.Reflection, point.JournalEntry.NextStep); err != nil {
			return fmt.Errorf("insert rebuilt course journal: %w", err)
		}
	}
	if point.RecordEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_point_record_entries (
				course_point_id, study_minutes, practice_count, confidence_level, application_note
			) VALUES ($1, $2, $3, $4, $5)
		`, point.Point.ID, point.RecordEntry.StudyMinutes, point.RecordEntry.PracticeCount, point.RecordEntry.ConfidenceLevel, point.RecordEntry.ApplicationNote); err != nil {
			return fmt.Errorf("insert rebuilt course record: %w", err)
		}
	}
	if point.ArtifactEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_point_artifact_entries (
				course_point_id, artifact_type, title, url, description
			) VALUES ($1, $2, $3, $4, $5)
		`, point.Point.ID, point.ArtifactEntry.ArtifactType, point.ArtifactEntry.Title, point.ArtifactEntry.URL, point.ArtifactEntry.Description); err != nil {
			return fmt.Errorf("insert rebuilt course artifact: %w", err)
		}
	}
	return nil
}
