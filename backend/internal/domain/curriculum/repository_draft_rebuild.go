package curriculum

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ReplaceDraftPlan(ctx context.Context, aggregate *DraftAggregate, chargeAmount int) error {
	if aggregate == nil {
		return fmt.Errorf("draft aggregate is nil")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin replace draft tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentStatus DraftStatus
	var confirmedCourseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT status, confirmed_course_id
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, aggregate.Draft.ID, aggregate.Draft.UserID).Scan(&currentStatus, &confirmedCourseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("draft not found")
		}
		return fmt.Errorf("load draft for replace: %w", err)
	}

	if currentStatus != DraftStatusDraft && currentStatus != DraftStatusConfirmed {
		return fmt.Errorf("draft not rebuildable")
	}

	if err := chargePointsTxWithDetail(ctx, tx, aggregate.Draft.UserID, chargeAmount, "use_course_gen", &PointTransactionDetail{
		Feature:       "course_draft_rebuild_all",
		ReferenceType: "course_draft",
		ReferenceID:   &aggregate.Draft.ID,
		Description:   strings.TrimSpace(aggregate.Draft.Title),
		Metadata: map[string]any{
			"mode": "rebuild_all",
		},
	}); err != nil {
		return err
	}

	if confirmedCourseID != nil {
		if err := deleteConfirmedCourseForRebuildTx(ctx, tx, *confirmedCourseID, aggregate.Draft.UserID); err != nil {
			return err
		}
	}
	if err := deleteDraftPlanDataForRebuildTx(ctx, tx, aggregate.Draft.ID, "rebuild"); err != nil {
		return err
	}
	if err := markExplorerPlanInactiveTx(ctx, tx, aggregate.Draft.ID); err != nil {
		return fmt.Errorf("deactivate explorer plan for rebuild: %w", err)
	}
	if err := updateDraftMetadataForRebuildAllTx(ctx, tx, aggregate); err != nil {
		return err
	}
	if err := insertRebuiltDraftLessonSkeletonsTx(ctx, tx, aggregate); err != nil {
		return err
	}
	if err := seedExplorerPlanFromDraftLessons(ctx, tx, aggregate); err != nil {
		return fmt.Errorf("seed rebuilt explorer plan: %w", err)
	}
	if err := insertRebuildEventTx(ctx, tx, aggregate.Draft.UserID, aggregate.Draft.ID, aggregate.Draft.SourceQuery, map[string]any{
		"title":                aggregate.Draft.Title,
		"main_lesson_count":    len(aggregate.Lessons),
		"goal_profile_id":      aggregate.Draft.GoalProfileID,
		"goal_profile_version": aggregate.Draft.GoalProfileVersion,
		"mode":                 "rebuild_all",
	}, "rebuild"); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit replace draft tx: %w", err)
	}
	return nil
}

func (r *Repository) ReplaceDraftAndCoursePlan(ctx context.Context, draftAggregate *DraftAggregate, courseAggregate *CourseAggregate, chargeAmount int, rebuildMode string) error {
	if draftAggregate == nil {
		return fmt.Errorf("draft aggregate is nil")
	}
	if courseAggregate == nil {
		return fmt.Errorf("course aggregate is nil")
	}
	if strings.TrimSpace(rebuildMode) == "" {
		rebuildMode = "rebuild_all"
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin replace draft+course tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentStatus DraftStatus
	var confirmedCourseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT status, confirmed_course_id
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, draftAggregate.Draft.ID, draftAggregate.Draft.UserID).Scan(&currentStatus, &confirmedCourseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("draft not found")
		}
		return fmt.Errorf("load draft for replace remaining: %w", err)
	}

	if currentStatus != DraftStatusConfirmed && currentStatus != DraftStatusLearning {
		return fmt.Errorf("draft not replaceable with course")
	}
	if confirmedCourseID == nil || *confirmedCourseID != courseAggregate.Course.ID {
		return fmt.Errorf("confirmed course mismatch")
	}

	if err := chargePointsTxWithDetail(ctx, tx, draftAggregate.Draft.UserID, chargeAmount, "use_course_gen", &PointTransactionDetail{
		Feature:       "course_draft_rebuild",
		ReferenceType: "course_draft",
		ReferenceID:   &draftAggregate.Draft.ID,
		Description:   strings.TrimSpace(draftAggregate.Draft.Title),
		Metadata: map[string]any{
			"mode":      rebuildMode,
			"course_id": courseAggregate.Course.ID.String(),
		},
	}); err != nil {
		return err
	}

	if err := deleteCoursePlanDataForRebuildTx(ctx, tx, courseAggregate.Course.ID, "remaining rebuild"); err != nil {
		return err
	}
	if err := deleteDraftPlanDataForRebuildTx(ctx, tx, draftAggregate.Draft.ID, "remaining rebuild"); err != nil {
		return err
	}
	if err := markExplorerPlanInactiveTx(ctx, tx, draftAggregate.Draft.ID); err != nil {
		return fmt.Errorf("deactivate explorer plan for remaining rebuild: %w", err)
	}
	if err := updateDraftMetadataForRemainingRebuildTx(ctx, tx, draftAggregate, currentStatus, confirmedCourseID); err != nil {
		return err
	}
	if err := insertRebuiltDraftTreesTx(ctx, tx, draftAggregate); err != nil {
		return err
	}
	if err := seedExplorerPlanFromDraftLessons(ctx, tx, draftAggregate); err != nil {
		return fmt.Errorf("seed explorer plan for remaining rebuild: %w", err)
	}
	if err := updateCourseMetadataForRemainingRebuildTx(ctx, tx, courseAggregate); err != nil {
		return err
	}
	if err := insertRebuiltCourseTreesTx(ctx, tx, courseAggregate); err != nil {
		return err
	}
	if err := insertRebuildEventTx(ctx, tx, draftAggregate.Draft.UserID, draftAggregate.Draft.ID, draftAggregate.Draft.SourceQuery, map[string]any{
		"title":                draftAggregate.Draft.Title,
		"main_lesson_count":    len(draftAggregate.Lessons),
		"goal_profile_id":      draftAggregate.Draft.GoalProfileID,
		"goal_profile_version": draftAggregate.Draft.GoalProfileVersion,
		"generation_language":  normalizeLearningLanguage(draftAggregate.Draft.GenerationLanguage),
		"mode":                 rebuildMode,
	}, "remaining rebuild"); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit replace draft+course tx: %w", err)
	}
	return nil
}
