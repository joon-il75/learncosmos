//go:build manual

package curriculum

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestManualRemainingRebuildE2E(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()

	var userID uuid.UUID
	if err := pool.QueryRow(ctx, `SELECT user_id FROM course_drafts LIMIT 1`).Scan(&userID); err != nil {
		t.Fatalf("load user id: %v", err)
	}

	repo := NewRepository(pool)
	draftAggregate := buildManualRebuildDraftAggregate(userID, "old")
	if err := repo.CreateDraft(ctx, draftAggregate, 0); err != nil {
		t.Fatalf("CreateDraft() error = %v", err)
	}

	courseID := uuid.New()
	confirmed := &CourseAggregate{
		Course: Course{
			ID:                 courseID,
			UserID:             userID,
			SourceDraftID:      &draftAggregate.Draft.ID,
			SourceQuery:        draftAggregate.Draft.SourceQuery,
			LearningGoal:       draftAggregate.Draft.LearningGoal,
			CurrentLevel:       draftAggregate.Draft.CurrentLevel,
			DurationWeeks:      draftAggregate.Draft.DurationWeeks,
			StudyHoursPerWeek:  draftAggregate.Draft.StudyHoursPerWeek,
			PreferredFormat:    draftAggregate.Draft.PreferredFormat,
			Title:              draftAggregate.Draft.Title,
			Description:        draftAggregate.Draft.Description,
			CompletionCriteria: draftAggregate.Draft.CompletionCriteria,
			Status:             CourseStatusActive,
			PlanetTypeID:       draftAggregate.Draft.PlanetTypeID,
		},
		Lessons: buildPlanetLessonTreesFromDraftLessons(draftAggregate.Lessons, courseID),
	}

	if err := repo.ConfirmDraft(ctx, confirmed, draftAggregate.Draft.ID, userID); err != nil {
		t.Fatalf("ConfirmDraft() error = %v", err)
	}
	if err := repo.StartLearningPlanet(ctx, userID, courseID); err != nil {
		t.Fatalf("StartLearningPlanet() error = %v", err)
	}

	touchedPointID := confirmed.Lessons[0].SubLessons[0].Points[0].Point.ID
	if _, err := repo.UpdateLearningPointRuntime(ctx, userID, courseID, touchedPointID, UpdateLearningPointRuntimeRequest{
		Status: LessonRuntimeInProgress,
	}); err != nil {
		t.Fatalf("UpdateLearningPointRuntime() error = %v", err)
	}

	currentDraft, err := repo.GetDraftByID(ctx, draftAggregate.Draft.ID, userID)
	if err != nil {
		t.Fatalf("GetDraftByID() error = %v", err)
	}
	currentCourse, err := repo.getCourseAggregateByID(ctx, courseID)
	if err != nil {
		t.Fatalf("getCourseAggregateByID() error = %v", err)
	}

	mask := buildRemainingRebuildPreserveMaskFromCourse(currentCourse.Lessons)
	if len(mask) != 2 || !mask[0] || mask[1] {
		t.Fatalf("unexpected preserve mask: %+v", mask)
	}

	generated := buildManualRebuildDraftAggregate(userID, "new")
	mergedDraftLessons, err := mergeRemainingDraftLessons(currentDraft.Lessons, generated.Lessons, mask)
	if err != nil {
		t.Fatalf("mergeRemainingDraftLessons() error = %v", err)
	}
	mergedCourseLessons, err := mergeRemainingCourseLessons(currentCourse.Lessons, generated.Lessons, mask, courseID)
	if err != nil {
		t.Fatalf("mergeRemainingCourseLessons() error = %v", err)
	}

	mergedDraft := &DraftAggregate{
		Draft:   currentDraft.Draft,
		Lessons: mergedDraftLessons,
	}
	mergedDraft.Draft.Title = "merged manual rebuild"
	mergedDraft.Draft.Status = DraftStatusLearning
	mergedDraft.Draft.ConfirmedCourseID = &courseID
	for mainIdx := range mergedDraft.Lessons {
		mergedDraft.Lessons[mainIdx].Lesson.CourseDraftID = mergedDraft.Draft.ID
		mergedDraft.Lessons[mainIdx].Lesson.ParentLessonID = nil
		mergedDraft.Lessons[mainIdx].Lesson.OrderIndex = mainIdx
		for subIdx := range mergedDraft.Lessons[mainIdx].SubLessons {
			mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Lesson.CourseDraftID = mergedDraft.Draft.ID
			mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Lesson.ParentLessonID = &mergedDraft.Lessons[mainIdx].Lesson.ID
			mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Lesson.OrderIndex = subIdx
			for pointIdx := range mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Points {
				mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseDraftID = mergedDraft.Draft.ID
				mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseDraftLessonID = mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Lesson.ID
				mergedDraft.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.OrderIndex = pointIdx
			}
		}
	}

	mergedCourse := &CourseAggregate{
		Course:  currentCourse.Course,
		Lessons: mergedCourseLessons,
	}
	mergedCourse.Course.Title = mergedDraft.Draft.Title
	for mainIdx := range mergedCourse.Lessons {
		mergedCourse.Lessons[mainIdx].Lesson.CourseID = courseID
		mergedCourse.Lessons[mainIdx].Lesson.ParentLessonID = nil
		mergedCourse.Lessons[mainIdx].Lesson.OrderIndex = mainIdx
		for subIdx := range mergedCourse.Lessons[mainIdx].SubLessons {
			mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.CourseID = courseID
			mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.ParentLessonID = &mergedCourse.Lessons[mainIdx].Lesson.ID
			mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.OrderIndex = subIdx
			if mergedCourse.Lessons[mainIdx].SubLessons[subIdx].RuntimeEntry != nil {
				mergedCourse.Lessons[mainIdx].SubLessons[subIdx].RuntimeEntry.CourseID = courseID
				mergedCourse.Lessons[mainIdx].SubLessons[subIdx].RuntimeEntry.CourseLessonID = mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.ID
			}
			for pointIdx := range mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Points {
				mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseID = courseID
				mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseLessonID = mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.ID
				mergedCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.OrderIndex = pointIdx
			}
		}
	}

	if err := repo.ReplaceDraftAndCoursePlan(ctx, mergedDraft, mergedCourse, 0, "rebuild_remaining"); err != nil {
		t.Fatalf("ReplaceDraftAndCoursePlan() error = %v", err)
	}

	updatedDraft, err := repo.GetDraftByID(ctx, draftAggregate.Draft.ID, userID)
	if err != nil {
		t.Fatalf("reload draft error = %v", err)
	}
	updatedCourse, err := repo.getCourseAggregateByID(ctx, courseID)
	if err != nil {
		t.Fatalf("reload course error = %v", err)
	}

	if updatedDraft.Lessons[0].Lesson.Title != "old-main-1" {
		t.Fatalf("first draft lesson should be preserved, got %q", updatedDraft.Lessons[0].Lesson.Title)
	}
	if updatedDraft.Lessons[1].Lesson.Title != "new-main-1" {
		t.Fatalf("second draft lesson should be replaced, got %q", updatedDraft.Lessons[1].Lesson.Title)
	}
	if updatedCourse.Lessons[0].Lesson.Title != "old-main-1" {
		t.Fatalf("first course lesson should be preserved, got %q", updatedCourse.Lessons[0].Lesson.Title)
	}
	if updatedCourse.Lessons[1].Lesson.Title != "new-main-1" {
		t.Fatalf("second course lesson should be replaced, got %q", updatedCourse.Lessons[1].Lesson.Title)
	}
	if updatedCourse.Lessons[0].SubLessons[0].Points[0].Point.Status != PointStatusLearning {
		t.Fatalf("touched point status should be preserved, got %q", updatedCourse.Lessons[0].SubLessons[0].Points[0].Point.Status)
	}

	cleanupManualRebuildArtifacts(t, pool, draftAggregate.Draft.ID, courseID)
}

func buildManualRebuildDraftAggregate(userID uuid.UUID, prefix string) *DraftAggregate {
	goal := prefix + "-goal"
	level := "beginner"
	duration := 4
	hours := 3
	format := "video"
	description := prefix + "-description"
	now := time.Now()

	main1ID := uuid.New()
	sub1ID := uuid.New()
	point1ID := uuid.New()
	main2ID := uuid.New()
	sub2ID := uuid.New()
	point2ID := uuid.New()

	return &DraftAggregate{
		Draft: CourseDraft{
			ID:                 uuid.New(),
			UserID:             userID,
			SourceQuery:        prefix + "-source",
			LearningGoal:       &goal,
			CurrentLevel:       &level,
			DurationWeeks:      &duration,
			StudyHoursPerWeek:  &hours,
			PreferredFormat:    &format,
			Title:              prefix + "-course",
			Description:        &description,
			CompletionCriteria: []string{},
			Status:             DraftStatusDraft,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		Lessons: []DraftLessonTree{
			{
				Lesson: CourseDraftLesson{
					ID:         main1ID,
					Title:      prefix + "-main-1",
					LessonRole: LessonRoleCore,
					SourceType: LessonSourceAI,
					OrderIndex: 0,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
				SubLessons: []DraftLessonTree{
					{
						Lesson: CourseDraftLesson{
							ID:             sub1ID,
							ParentLessonID: &main1ID,
							Title:          prefix + "-sub-1",
							LessonRole:     LessonRoleCore,
							SourceType:     LessonSourceAI,
							OrderIndex:     0,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						Points: []DraftPointAggregate{
							{
								Point: CourseDraftPoint{
									ID:                  point1ID,
									PointType:           PointTypeExploration,
									Status:              PointStatusDraft,
									Title:               prefix + "-point-1",
									OrderIndex:          0,
									CreatedAt:           now,
									UpdatedAt:           now,
									SelectionState:      ptrResourceSelection(ResourceSelectionCandidate),
									ContentID:           nil,
									CourseDraftLessonID: sub1ID,
								},
							},
						},
					},
				},
			},
			{
				Lesson: CourseDraftLesson{
					ID:         main2ID,
					Title:      prefix + "-main-2",
					LessonRole: LessonRoleCore,
					SourceType: LessonSourceAI,
					OrderIndex: 1,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
				SubLessons: []DraftLessonTree{
					{
						Lesson: CourseDraftLesson{
							ID:             sub2ID,
							ParentLessonID: &main2ID,
							Title:          prefix + "-sub-2",
							LessonRole:     LessonRoleCore,
							SourceType:     LessonSourceAI,
							OrderIndex:     0,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						Points: []DraftPointAggregate{
							{
								Point: CourseDraftPoint{
									ID:                  point2ID,
									PointType:           PointTypeExploration,
									Status:              PointStatusDraft,
									Title:               prefix + "-point-2",
									OrderIndex:          0,
									CreatedAt:           now,
									UpdatedAt:           now,
									SelectionState:      ptrResourceSelection(ResourceSelectionCandidate),
									ContentID:           nil,
									CourseDraftLessonID: sub2ID,
								},
							},
						},
					},
				},
			},
		},
	}
}

func ptrResourceSelection(v ResourceSelectionState) *ResourceSelectionState {
	return &v
}

func cleanupManualRebuildArtifacts(t *testing.T, pool *pgxpool.Pool, draftID, courseID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	queries := []string{
		`DELETE FROM course_point_blocks WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`,
		`DELETE FROM course_point_ai_summaries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`,
		`DELETE FROM course_point_journal_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`,
		`DELETE FROM course_point_record_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`,
		`DELETE FROM course_point_artifact_entries WHERE course_point_id IN (SELECT id FROM course_points WHERE course_id = $1)`,
		`DELETE FROM course_lesson_runtime_entries WHERE course_id = $1`,
		`DELETE FROM course_points WHERE course_id = $1`,
		`DELETE FROM course_lessons WHERE course_id = $1`,
		`DELETE FROM courses WHERE id = $1`,
	}
	for _, q := range queries {
		if _, err := pool.Exec(ctx, q, courseID); err != nil {
			t.Fatalf("cleanup course query failed: %v", err)
		}
	}

	draftQueries := []string{
		`DELETE FROM course_draft_point_blocks WHERE course_draft_point_id IN (SELECT id FROM course_draft_points WHERE course_draft_id = $1)`,
		`DELETE FROM course_draft_journal_entries WHERE course_draft_id = $1`,
		`DELETE FROM course_draft_record_entries WHERE course_draft_id = $1`,
		`DELETE FROM course_draft_artifact_entries WHERE course_draft_id = $1`,
		`DELETE FROM course_draft_detail_notes WHERE course_draft_id = $1`,
		`DELETE FROM course_draft_points WHERE course_draft_id = $1`,
		`DELETE FROM course_draft_lessons WHERE course_draft_id = $1`,
		`DELETE FROM explorer_nodes WHERE parent_kind = 'subregion' AND parent_id IN (SELECT es.id FROM explorer_subregions es JOIN explorer_regions er ON er.id = es.region_id WHERE er.course_draft_id = $1)`,
		`DELETE FROM explorer_nodes WHERE parent_kind = 'region' AND parent_id IN (SELECT id FROM explorer_regions WHERE course_draft_id = $1)`,
		`DELETE FROM explorer_subregions WHERE region_id IN (SELECT id FROM explorer_regions WHERE course_draft_id = $1)`,
		`DELETE FROM explorer_regions WHERE course_draft_id = $1`,
		`DELETE FROM recommendation_events WHERE course_draft_id = $1`,
		`DELETE FROM course_drafts WHERE id = $1`,
	}
	for _, q := range draftQueries {
		if _, err := pool.Exec(ctx, q, draftID); err != nil {
			t.Fatalf("cleanup draft query failed: %v", err)
		}
	}
}
