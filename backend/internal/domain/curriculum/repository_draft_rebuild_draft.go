package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func updateDraftMetadataForRebuildAllTx(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate) error {
	if _, err := tx.Exec(ctx, `
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
		    status = 'draft',
		    confirmed_course_id = NULL,
		    updated_at = NOW()
		WHERE id = $11 AND user_id = $12
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
		aggregate.Draft.ID,
		aggregate.Draft.UserID,
	); err != nil {
		return fmt.Errorf("update draft metadata for rebuild: %w", err)
	}
	return nil
}

func updateDraftMetadataForRemainingRebuildTx(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate, status DraftStatus, confirmedCourseID *uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET source_query = $1,
		    learning_goal = $2,
		    goal_profile_id = $3,
		    goal_profile_version = $4,
		    current_level = $5,
		    duration_weeks = $6,
		    study_hours_per_week = $7,
		    preferred_format = $8,
		    generation_language = $9,
		    title = $10,
		    description = $11,
		    completion_criteria = $12,
		    status = $13,
		    confirmed_course_id = $14,
		    updated_at = NOW()
		WHERE id = $15 AND user_id = $16
	`,
		aggregate.Draft.SourceQuery,
		aggregate.Draft.LearningGoal,
		aggregate.Draft.GoalProfileID,
		aggregate.Draft.GoalProfileVersion,
		aggregate.Draft.CurrentLevel,
		aggregate.Draft.DurationWeeks,
		aggregate.Draft.StudyHoursPerWeek,
		aggregate.Draft.PreferredFormat,
		normalizeLearningLanguage(aggregate.Draft.GenerationLanguage),
		aggregate.Draft.Title,
		aggregate.Draft.Description,
		aggregate.Draft.CompletionCriteria,
		status,
		confirmedCourseID,
		aggregate.Draft.ID,
		aggregate.Draft.UserID,
	); err != nil {
		return fmt.Errorf("update draft metadata for remaining rebuild: %w", err)
	}
	return nil
}

func insertRebuiltDraftLessonSkeletonsTx(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate) error {
	for _, mainTree := range aggregate.Lessons {
		if err := insertRebuiltDraftLessonTx(ctx, tx, aggregate.Draft.ID, mainTree.Lesson); err != nil {
			return fmt.Errorf("insert rebuilt main lesson: %w", err)
		}
		for _, subTree := range mainTree.SubLessons {
			if err := insertRebuiltDraftLessonTx(ctx, tx, aggregate.Draft.ID, subTree.Lesson); err != nil {
				return fmt.Errorf("insert rebuilt sub lesson: %w", err)
			}
		}
	}
	return nil
}

func insertRebuiltDraftTreesTx(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate) error {
	for _, mainTree := range aggregate.Lessons {
		if err := insertRebuiltDraftTreeTx(ctx, tx, aggregate.Draft.ID, mainTree); err != nil {
			return err
		}
	}
	return nil
}

func insertRebuiltDraftTreeTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, mainTree DraftLessonTree) error {
	if err := insertRebuiltDraftLessonTx(ctx, tx, draftID, mainTree.Lesson); err != nil {
		return fmt.Errorf("insert rebuilt draft main lesson: %w", err)
	}
	if err := insertRebuiltDraftLessonNoteTx(ctx, tx, draftID, mainTree.Lesson, "main"); err != nil {
		return err
	}
	for _, point := range mainTree.Points {
		if err := insertRebuiltDraftPointTx(ctx, tx, draftID, point, "main"); err != nil {
			return err
		}
	}
	for _, subTree := range mainTree.SubLessons {
		if err := insertRebuiltDraftSubTreeTx(ctx, tx, draftID, subTree); err != nil {
			return err
		}
	}
	return nil
}

func insertRebuiltDraftSubTreeTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, subTree DraftLessonTree) error {
	if err := insertRebuiltDraftLessonTx(ctx, tx, draftID, subTree.Lesson); err != nil {
		return fmt.Errorf("insert rebuilt draft sub lesson: %w", err)
	}
	if err := insertRebuiltDraftLessonNoteTx(ctx, tx, draftID, subTree.Lesson, "sub"); err != nil {
		return err
	}
	if err := insertRebuiltDraftLessonEntriesTx(ctx, tx, draftID, subTree.Lesson); err != nil {
		return err
	}
	for _, point := range subTree.Points {
		if err := insertRebuiltDraftPointTx(ctx, tx, draftID, point, "sub"); err != nil {
			return err
		}
	}
	return nil
}

func insertRebuiltDraftLessonTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, lesson CourseDraftLesson) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO course_draft_lessons (
			id, course_draft_id, parent_lesson_id, title, objective, summary,
			difficulty_level, lesson_role, source_type, order_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		lesson.ID,
		draftID,
		lesson.ParentLessonID,
		lesson.Title,
		lesson.Objective,
		lesson.Summary,
		lesson.DifficultyLevel,
		lesson.LessonRole,
		lesson.SourceType,
		lesson.OrderIndex,
	)
	return err
}

func insertRebuiltDraftLessonNoteTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, lesson CourseDraftLesson, label string) error {
	if lesson.OperationNote == nil || strings.TrimSpace(*lesson.OperationNote) == "" {
		return nil
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO course_draft_detail_notes (
			course_draft_id, course_draft_lesson_id, target_type, note
		) VALUES ($1, $2, 'lesson', $3)
	`, draftID, lesson.ID, strings.TrimSpace(*lesson.OperationNote)); err != nil {
		return fmt.Errorf("insert rebuilt draft %s lesson note: %w", label, err)
	}
	return nil
}

func insertRebuiltDraftLessonEntriesTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, lesson CourseDraftLesson) error {
	if lesson.JournalEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_journal_entries (
				course_draft_id, course_draft_lesson_id, observation, reflection, next_step
			) VALUES ($1, $2, $3, $4, $5)
		`, draftID, lesson.ID, lesson.JournalEntry.Observation, lesson.JournalEntry.Reflection, lesson.JournalEntry.NextStep); err != nil {
			return fmt.Errorf("insert rebuilt draft sub lesson journal: %w", err)
		}
	}
	if lesson.RecordEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_record_entries (
				course_draft_id, course_draft_lesson_id, study_minutes, practice_count, confidence_level, application_note
			) VALUES ($1, $2, $3, $4, $5, $6)
		`, draftID, lesson.ID, lesson.RecordEntry.StudyMinutes, lesson.RecordEntry.PracticeCount, lesson.RecordEntry.ConfidenceLevel, lesson.RecordEntry.ApplicationNote); err != nil {
			return fmt.Errorf("insert rebuilt draft sub lesson record: %w", err)
		}
	}
	if lesson.ArtifactEntry != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_artifact_entries (
				course_draft_id, course_draft_lesson_id, artifact_type, title, url, description
			) VALUES ($1, $2, $3, $4, $5, $6)
		`, draftID, lesson.ID, lesson.ArtifactEntry.ArtifactType, lesson.ArtifactEntry.Title, lesson.ArtifactEntry.URL, lesson.ArtifactEntry.Description); err != nil {
			return fmt.Errorf("insert rebuilt draft sub lesson artifact: %w", err)
		}
	}
	return nil
}

func insertRebuiltDraftPointTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, point DraftPointAggregate, label string) error {
	priceType := PriceTypeFree
	if point.Point.PriceType != nil {
		priceType = *point.Point.PriceType
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO course_draft_points (
			id, course_draft_id, course_draft_lesson_id, point_type, status, title, description,
			template_type, selection_state, content_id, external_url, thumbnail_url, price_type, rank_score,
			order_index, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`,
		point.Point.ID,
		draftID,
		point.Point.CourseDraftLessonID,
		point.Point.PointType,
		point.Point.Status,
		point.Point.Title,
		point.Point.Description,
		point.Point.TemplateType,
		point.Point.SelectionState,
		point.Point.ContentID,
		point.Point.ExternalURL,
		point.Point.ThumbnailURL,
		priceType,
		point.Point.RankScore,
		point.Point.OrderIndex,
		point.Point.CompletedAt,
	); err != nil {
		return fmt.Errorf("insert rebuilt draft %s point: %w", label, err)
	}
	return insertRebuiltDraftPointEntriesTx(ctx, tx, point)
}

func insertRebuiltDraftPointEntriesTx(ctx context.Context, tx pgx.Tx, point DraftPointAggregate) error {
	for _, block := range point.Blocks {
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_point_blocks (
				id, course_draft_point_id, block_type, content, order_index
			) VALUES ($1, $2, $3, $4::jsonb, $5)
		`, block.ID, point.Point.ID, block.BlockType, string(block.Content), block.OrderIndex); err != nil {
			return fmt.Errorf("insert rebuilt draft point block: %w", err)
		}
	}
	return nil
}
