package curriculum

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
)

func normalizeRebuiltDraftAggregate(existing CourseDraft, rebuilt *DraftAggregate) {
	if rebuilt == nil {
		return
	}
	rebuilt.Draft.ID = existing.ID
	rebuilt.Draft.UserID = existing.UserID
	rebuilt.Draft.Status = existing.Status
	rebuilt.Draft.ConfirmedCourseID = existing.ConfirmedCourseID
	rebuilt.Draft.PlanetTypeID = existing.PlanetTypeID
	rebuilt.Draft.PlanetTypeName = existing.PlanetTypeName
	rebuilt.Draft.PlanetTypeAsset = existing.PlanetTypeAsset
	for mainIdx := range rebuilt.Lessons {
		rebuilt.Lessons[mainIdx].Lesson.CourseDraftID = existing.ID
		rebuilt.Lessons[mainIdx].Lesson.ParentLessonID = nil
		rebuilt.Lessons[mainIdx].Lesson.OrderIndex = mainIdx
		for subIdx := range rebuilt.Lessons[mainIdx].SubLessons {
			rebuilt.Lessons[mainIdx].SubLessons[subIdx].Lesson.CourseDraftID = existing.ID
			rebuilt.Lessons[mainIdx].SubLessons[subIdx].Lesson.ParentLessonID = &rebuilt.Lessons[mainIdx].Lesson.ID
			rebuilt.Lessons[mainIdx].SubLessons[subIdx].Lesson.OrderIndex = subIdx
			for pointIdx := range rebuilt.Lessons[mainIdx].SubLessons[subIdx].Points {
				rebuilt.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseDraftID = existing.ID
				rebuilt.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseDraftLessonID = rebuilt.Lessons[mainIdx].SubLessons[subIdx].Lesson.ID
				rebuilt.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.OrderIndex = pointIdx
			}
		}
		for pointIdx := range rebuilt.Lessons[mainIdx].Points {
			rebuilt.Lessons[mainIdx].Points[pointIdx].Point.CourseDraftID = existing.ID
			rebuilt.Lessons[mainIdx].Points[pointIdx].Point.CourseDraftLessonID = rebuilt.Lessons[mainIdx].Lesson.ID
			rebuilt.Lessons[mainIdx].Points[pointIdx].Point.OrderIndex = pointIdx
		}
	}
}

func buildFullRebuildCourseAggregate(rebuiltDraft *DraftAggregate, existingCourse *CourseAggregate) *CourseAggregate {
	if rebuiltDraft == nil || existingCourse == nil {
		return nil
	}

	rebuiltCourse := &CourseAggregate{
		Course:  existingCourse.Course,
		Lessons: buildPlanetLessonTreesFromDraftLessons(rebuiltDraft.Lessons, existingCourse.Course.ID),
	}
	rebuiltCourse.Course.SourceQuery = rebuiltDraft.Draft.SourceQuery
	rebuiltCourse.Course.LearningGoal = rebuiltDraft.Draft.LearningGoal
	rebuiltCourse.Course.CurrentLevel = rebuiltDraft.Draft.CurrentLevel
	rebuiltCourse.Course.DurationWeeks = rebuiltDraft.Draft.DurationWeeks
	rebuiltCourse.Course.StudyHoursPerWeek = rebuiltDraft.Draft.StudyHoursPerWeek
	rebuiltCourse.Course.PreferredFormat = rebuiltDraft.Draft.PreferredFormat
	rebuiltCourse.Course.Title = rebuiltDraft.Draft.Title
	rebuiltCourse.Course.Description = rebuiltDraft.Draft.Description
	rebuiltCourse.Course.CompletionCriteria = rebuiltDraft.Draft.CompletionCriteria
	rebuiltCourse.Course.PlanetTypeID = existingCourse.Course.PlanetTypeID
	rebuiltCourse.Course.PlanetTypeName = existingCourse.Course.PlanetTypeName
	rebuiltCourse.Course.PlanetTypeAsset = existingCourse.Course.PlanetTypeAsset
	for mainIdx := range rebuiltCourse.Lessons {
		rebuiltCourse.Lessons[mainIdx].Lesson.CourseID = existingCourse.Course.ID
		rebuiltCourse.Lessons[mainIdx].Lesson.ParentLessonID = nil
		rebuiltCourse.Lessons[mainIdx].Lesson.OrderIndex = mainIdx
		for subIdx := range rebuiltCourse.Lessons[mainIdx].SubLessons {
			rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.CourseID = existingCourse.Course.ID
			rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.ParentLessonID = &rebuiltCourse.Lessons[mainIdx].Lesson.ID
			rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.OrderIndex = subIdx
			for pointIdx := range rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Points {
				rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseID = existingCourse.Course.ID
				rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.CourseLessonID = rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Lesson.ID
				rebuiltCourse.Lessons[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.OrderIndex = pointIdx
			}
		}
		for pointIdx := range rebuiltCourse.Lessons[mainIdx].Points {
			rebuiltCourse.Lessons[mainIdx].Points[pointIdx].Point.CourseID = existingCourse.Course.ID
			rebuiltCourse.Lessons[mainIdx].Points[pointIdx].Point.CourseLessonID = rebuiltCourse.Lessons[mainIdx].Lesson.ID
			rebuiltCourse.Lessons[mainIdx].Points[pointIdx].Point.OrderIndex = pointIdx
		}
	}

	return rebuiltCourse
}

func buildRebuildDraftRequest(draft CourseDraft, activeGoal *goal.GoalProfile, confirmedGoal string) CreateCourseDraftRequest {
	req := CreateCourseDraftRequest{
		SourceQuery:         draft.SourceQuery,
		LearningGoal:        &confirmedGoal,
		CurrentLevel:        draft.CurrentLevel,
		DurationWeeks:       draft.DurationWeeks,
		StudyHoursPerWeek:   draft.StudyHoursPerWeek,
		PreferredFormat:     draft.PreferredFormat,
		GoalProfileID:       &activeGoal.ID,
		GoalProfileVersion:  &activeGoal.Version,
		GoalUserIntent:      &activeGoal.UserIntent,
		GoalMotivation:      activeGoal.Motivation,
		GoalUsageContext:    activeGoal.UsageContext,
		GoalDifficultyLevel: activeGoal.DifficultyLevel,
		GoalTimeHorizon:     activeGoal.TimeHorizon,
		GoalOutputType:      activeGoal.OutputType,
		GoalType:            activeGoal.GoalType,
		LearningIntent:      learningIntentProfileFromGoal(activeGoal.LearningIntent),
		LearningLanguage:    normalizeLearningLanguage(activeGoal.Language),
	}
	return req
}

func isDraftPointTouched(point DraftPointAggregate) bool {
	if point.Point.Status == PointStatusLearning || point.Point.Status == PointStatusCompleted || point.Point.CompletedAt != nil {
		return true
	}
	if point.Point.SelectionState != nil && *point.Point.SelectionState != ResourceSelectionCandidate {
		return true
	}
	if len(point.Blocks) > 0 || point.JournalEntry != nil || point.RecordEntry != nil || point.ArtifactEntry != nil {
		return true
	}
	return false
}

func isCoursePointTouched(point CoursePointAggregate) bool {
	if point.Point.Status == PointStatusLearning || point.Point.Status == PointStatusCompleted || point.Point.CompletedAt != nil {
		return true
	}
	if len(point.Blocks) > 0 || point.AISummary != nil || point.JournalEntry != nil || point.RecordEntry != nil || point.ArtifactEntry != nil {
		return true
	}
	return false
}

func mainLessonHasTouchedDraftPoint(main DraftLessonTree) bool {
	for _, point := range main.Points {
		if isDraftPointTouched(point) {
			return true
		}
	}
	for _, sub := range main.SubLessons {
		for _, point := range sub.Points {
			if isDraftPointTouched(point) {
				return true
			}
		}
	}
	return false
}

func mainLessonHasTouchedCoursePoint(main CourseLessonTree) bool {
	for _, point := range main.Points {
		if isCoursePointTouched(point) {
			return true
		}
	}
	for _, sub := range main.SubLessons {
		if sub.RuntimeEntry != nil {
			return true
		}
		for _, point := range sub.Points {
			if isCoursePointTouched(point) {
				return true
			}
		}
	}
	return false
}

func buildRemainingRebuildPreserveMaskFromDraft(lessons []DraftLessonTree) []bool {
	mask := make([]bool, len(lessons))
	for i, lesson := range lessons {
		mask[i] = mainLessonHasTouchedDraftPoint(lesson)
	}
	return mask
}

func buildRemainingRebuildPreserveMaskFromCourse(lessons []CourseLessonTree) []bool {
	mask := make([]bool, len(lessons))
	for i, lesson := range lessons {
		mask[i] = mainLessonHasTouchedCoursePoint(lesson)
	}
	return mask
}

func mergeRemainingDraftLessons(existing []DraftLessonTree, generated []DraftLessonTree, preserveMask []bool) ([]DraftLessonTree, error) {
	if len(existing) != len(preserveMask) {
		return nil, fmt.Errorf("preserve mask length mismatch")
	}

	replacementCount := 0
	for _, preserve := range preserveMask {
		if !preserve {
			replacementCount++
		}
	}
	if replacementCount == 0 {
		return existing, nil
	}
	if len(generated) < replacementCount {
		return nil, fmt.Errorf("not enough generated lessons for remaining rebuild")
	}

	merged := make([]DraftLessonTree, 0, len(existing))
	generatedIdx := 0
	for i, existingLesson := range existing {
		if preserveMask[i] {
			merged = append(merged, existingLesson)
			continue
		}
		merged = append(merged, generated[generatedIdx])
		generatedIdx++
	}
	return merged, nil
}

func mergeRemainingCourseLessons(existing []CourseLessonTree, generated []DraftLessonTree, preserveMask []bool, courseID uuid.UUID) ([]CourseLessonTree, error) {
	if len(existing) != len(preserveMask) {
		return nil, fmt.Errorf("preserve mask length mismatch")
	}

	replacementCount := 0
	for _, preserve := range preserveMask {
		if !preserve {
			replacementCount++
		}
	}
	if replacementCount == 0 {
		return existing, nil
	}
	if len(generated) < replacementCount {
		return nil, fmt.Errorf("not enough generated lessons for remaining rebuild")
	}

	generatedCourse := buildPlanetLessonTreesFromDraftLessons(generated, courseID)
	merged := make([]CourseLessonTree, 0, len(existing))
	generatedIdx := 0
	for i, existingLesson := range existing {
		if preserveMask[i] {
			merged = append(merged, existingLesson)
			continue
		}
		merged = append(merged, generatedCourse[generatedIdx])
		generatedIdx++
	}
	return merged, nil
}

func buildDraftLessonTreesFromCourseLessons(courseLessons []CourseLessonTree) []DraftLessonTree {
	trees := make([]DraftLessonTree, 0, len(courseLessons))
	for _, courseLesson := range courseLessons {
		lesson := CourseDraftLesson{
			ID:              courseLesson.Lesson.ID,
			CourseDraftID:   uuid.Nil,
			ParentLessonID:  courseLesson.Lesson.ParentLessonID,
			Title:           courseLesson.Lesson.Title,
			Objective:       courseLesson.Lesson.Objective,
			Summary:         courseLesson.Lesson.Summary,
			DifficultyLevel: courseLesson.Lesson.DifficultyLevel,
			LessonRole:      courseLesson.Lesson.LessonRole,
			SourceType:      courseLesson.Lesson.SourceType,
			OrderIndex:      courseLesson.Lesson.OrderIndex,
			CreatedAt:       courseLesson.Lesson.CreatedAt,
			UpdatedAt:       courseLesson.Lesson.UpdatedAt,
		}
		tree := DraftLessonTree{
			Lesson:     lesson,
			Points:     make([]DraftPointAggregate, 0, len(courseLesson.Points)),
			SubLessons: buildDraftLessonTreesFromCourseLessons(courseLesson.SubLessons),
		}
		for _, coursePoint := range courseLesson.Points {
			point := DraftPointAggregate{
				Point: CourseDraftPoint{
					ID:                  coursePoint.Point.ID,
					CourseDraftID:       uuid.Nil,
					CourseDraftLessonID: coursePoint.Point.CourseLessonID,
					PointType:           coursePoint.Point.PointType,
					Status:              coursePoint.Point.Status,
					Title:               coursePoint.Point.Title,
					Description:         coursePoint.Point.Description,
					TemplateType:        coursePoint.Point.TemplateType,
					ContentID:           coursePoint.Point.ContentID,
					ExternalURL:         coursePoint.Point.ExternalURL,
					ThumbnailURL:        coursePoint.Point.ThumbnailURL,
					PriceType:           coursePoint.Point.PriceType,
					RankScore:           coursePoint.Point.RankScore,
					OrderIndex:          coursePoint.Point.OrderIndex,
					CompletedAt:         coursePoint.Point.CompletedAt,
					CreatedAt:           coursePoint.Point.CreatedAt,
					UpdatedAt:           coursePoint.Point.UpdatedAt,
				},
				JournalEntry:  coursePoint.JournalEntry,
				RecordEntry:   coursePoint.RecordEntry,
				ArtifactEntry: coursePoint.ArtifactEntry,
			}
			for _, block := range coursePoint.Blocks {
				point.Blocks = append(point.Blocks, CourseDraftPointBlock{
					ID:                 block.ID,
					CourseDraftPointID: block.CoursePointID,
					BlockType:          block.BlockType,
					Content:            block.Content,
					OrderIndex:         block.OrderIndex,
					CreatedAt:          block.CreatedAt,
					UpdatedAt:          block.UpdatedAt,
				})
			}
			tree.Points = append(tree.Points, point)
		}
		trees = append(trees, tree)
	}
	return trees
}
