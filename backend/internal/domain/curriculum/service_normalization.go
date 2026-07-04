package curriculum

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func (s *Service) NormalizeDraftAggregate(draft *DraftAggregate) error {
	if draft == nil {
		return fmt.Errorf("draft aggregate is nil")
	}
	if draft.Draft.UserID == uuid.Nil {
		return fmt.Errorf("draft user_id is required")
	}
	draft.Draft.SourceQuery = strings.TrimSpace(draft.Draft.SourceQuery)
	draft.Draft.Title = strings.TrimSpace(draft.Draft.Title)
	if draft.Draft.SourceQuery == "" {
		return fmt.Errorf("source_query is required")
	}
	if draft.Draft.Title == "" {
		return fmt.Errorf("draft title is required")
	}
	if draft.Draft.Status == "" {
		draft.Draft.Status = DraftStatusDraft
	}

	slices.SortStableFunc(draft.Lessons, func(a, b DraftLessonTree) int {
		return a.Lesson.OrderIndex - b.Lesson.OrderIndex
	})
	for mainIdx := range draft.Lessons {
		tree := &draft.Lessons[mainIdx]
		tree.Lesson.CourseDraftID = draft.Draft.ID
		tree.Lesson.ParentLessonID = nil
		tree.Lesson.Title = strings.TrimSpace(tree.Lesson.Title)
		if tree.Lesson.Title == "" {
			return fmt.Errorf("main lesson title is required")
		}
		if tree.Lesson.LessonRole == "" {
			tree.Lesson.LessonRole = LessonRoleCore
		}
		if tree.Lesson.SourceType == "" {
			tree.Lesson.SourceType = LessonSourceManual
		}
		tree.Lesson.OrderIndex = mainIdx

		slices.SortStableFunc(tree.SubLessons, func(a, b DraftLessonTree) int {
			return a.Lesson.OrderIndex - b.Lesson.OrderIndex
		})
		for subIdx := range tree.SubLessons {
			sub := &tree.SubLessons[subIdx]
			mainID := tree.Lesson.ID
			sub.Lesson.CourseDraftID = draft.Draft.ID
			sub.Lesson.ParentLessonID = &mainID
			sub.Lesson.Title = strings.TrimSpace(sub.Lesson.Title)
			if sub.Lesson.Title == "" {
				return fmt.Errorf("lesson title is required")
			}
			if sub.Lesson.LessonRole == "" {
				sub.Lesson.LessonRole = LessonRoleCore
			}
			if sub.Lesson.SourceType == "" {
				sub.Lesson.SourceType = LessonSourceManual
			}
			sub.Lesson.OrderIndex = subIdx

			slices.SortStableFunc(sub.Points, func(a, b DraftPointAggregate) int {
				return a.Point.OrderIndex - b.Point.OrderIndex
			})
			for resourceIdx := range sub.Points {
				point := &sub.Points[resourceIdx]
				point.Point.OrderIndex = resourceIdx
			}
		}
	}
	return nil
}

func (s *Service) BuildConfirmedCourseFromDraft(draft DraftAggregate) (*CourseAggregate, error) {
	if err := s.NormalizeDraftAggregate(&draft); err != nil {
		return nil, err
	}
	if draft.Draft.Status == DraftStatusArchived {
		return nil, fmt.Errorf("archived draft cannot be confirmed")
	}

	sourceDraftID := draft.Draft.ID
	course := Course{
		ID:                 uuid.New(),
		UserID:             draft.Draft.UserID,
		SourceDraftID:      &sourceDraftID,
		SourceQuery:        draft.Draft.SourceQuery,
		LearningGoal:       draft.Draft.LearningGoal,
		CurrentLevel:       draft.Draft.CurrentLevel,
		DurationWeeks:      draft.Draft.DurationWeeks,
		StudyHoursPerWeek:  draft.Draft.StudyHoursPerWeek,
		PreferredFormat:    draft.Draft.PreferredFormat,
		Title:              draft.Draft.Title,
		Description:        draft.Draft.Description,
		CompletionCriteria: append([]string(nil), draft.Draft.CompletionCriteria...),
		Status:             CourseStatusActive,
		PlanetTypeID:       draft.Draft.PlanetTypeID,
		PlanetTypeName:     draft.Draft.PlanetTypeName,
		PlanetTypeAsset:    draft.Draft.PlanetTypeAsset,
	}

	result := &CourseAggregate{Course: course, Lessons: make([]CourseLessonTree, 0, len(draft.Lessons))}
	for _, mainTree := range draft.Lessons {
		mainID := uuid.New()
		mainLesson := CourseLesson{
			ID:                       mainID,
			CourseID:                 course.ID,
			Title:                    mainTree.Lesson.Title,
			Objective:                mainTree.Lesson.Objective,
			Summary:                  mainTree.Lesson.Summary,
			RecommendationSearchSpec: mainTree.Lesson.RecommendationSearchSpec,
			LessonRole:               mainTree.Lesson.LessonRole,
			SourceType:               mainTree.Lesson.SourceType,
			OrderIndex:               mainTree.Lesson.OrderIndex,
		}
		mainCourseLessonTree := CourseLessonTree{
			Lesson:     mainLesson,
			SubLessons: make([]CourseLessonTree, 0, len(mainTree.SubLessons)),
		}

		for _, subTree := range mainTree.SubLessons {
			subID := uuid.New()
			subLesson := CourseLesson{
				ID:                       subID,
				CourseID:                 course.ID,
				ParentLessonID:           &mainID,
				Title:                    subTree.Lesson.Title,
				Objective:                subTree.Lesson.Objective,
				Summary:                  subTree.Lesson.Summary,
				DifficultyLevel:          subTree.Lesson.DifficultyLevel,
				RecommendationSearchSpec: subTree.Lesson.RecommendationSearchSpec,
				LessonRole:               subTree.Lesson.LessonRole,
				SourceType:               subTree.Lesson.SourceType,
				OrderIndex:               subTree.Lesson.OrderIndex,
			}
			subCourseLessonTree := CourseLessonTree{
				Lesson: subLesson,
				Points: make([]CoursePointAggregate, 0, len(subTree.Points)),
			}
			// resources from draft sub-lesson → points in course lesson
			for _, pt := range subTree.Points {
				subCourseLessonTree.Points = append(subCourseLessonTree.Points, CoursePointAggregate{
					Point: CoursePoint{
						ID:             uuid.New(),
						CourseID:       course.ID,
						CourseLessonID: subID,
						PointType:      pt.Point.PointType,
						Title:          pt.Point.Title,
						Description:    pt.Point.Description,
						ContentID:      pt.Point.ContentID,
						ExternalURL:    pt.Point.ExternalURL,
						ThumbnailURL:   pt.Point.ThumbnailURL,
						PriceType:      pt.Point.PriceType,
						RankScore:      pt.Point.RankScore,
						OrderIndex:     pt.Point.OrderIndex,
					},
				})
			}
			mainCourseLessonTree.SubLessons = append(mainCourseLessonTree.SubLessons, subCourseLessonTree)
		}
		result.Lessons = append(result.Lessons, mainCourseLessonTree)
	}

	return result, nil
}

func normalizeCompletionCriteria(criteria []string) []string {
	result := make([]string, 0, len(criteria))
	seen := make(map[string]struct{}, len(criteria))
	for _, item := range criteria {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
