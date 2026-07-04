package curriculum

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) GetTodayTask(ctx context.Context, userID uuid.UUID) (*TodayTask, error) {
	activePlanets, err := r.ListPlanetsByDraftStatuses(ctx, userID, []DraftStatus{DraftStatusConfirmed, DraftStatusLearning})
	if err != nil {
		return nil, err
	}
	sortPlanetListItemsByRecentActivity(activePlanets)

	planetAggregates := make([]*PlanetAggregate, 0, len(activePlanets))
	for _, planet := range activePlanets {
		aggregate, err := r.GetPlanetByIDAndDraftStatuses(ctx, userID, planet.ID, []DraftStatus{DraftStatusConfirmed, DraftStatusLearning})
		if err != nil {
			return nil, err
		}
		planetAggregates = append(planetAggregates, aggregate)
	}

	completedPlanets, err := r.ListPlanetsByDraftStatuses(ctx, userID, []DraftStatus{DraftStatusArchived})
	if err != nil {
		return nil, err
	}
	sortPlanetListItemsByRecentActivity(completedPlanets)

	return buildTodayTaskFromPlanetAggregates(planetAggregates, completedPlanets), nil
}

func buildTodayTaskFromPlanetAggregates(activePlanets []*PlanetAggregate, completedPlanets []PlanetListItem) *TodayTask {
	for _, planet := range activePlanets {
		if planet == nil {
			continue
		}

		switch normalizeTodayTaskStatus(planet.Planet.Status) {
		case "ready":
			return buildStartPlanetTodayTask(planet)
		case "learning":
			if task := buildContinuePointTodayTask(planet); task != nil {
				return task
			}
			if planet.Planet.CanComplete {
				return buildCompletePlanetTodayTask(planet)
			}
		}
	}

	if len(completedPlanets) > 0 {
		return buildReviewRecordsTodayTask(completedPlanets[0])
	}

	return buildCreateCourseTodayTask()
}

func buildContinuePointTodayTask(planet *PlanetAggregate) *TodayTask {
	pointDetails := make([]PlanetPointDetail, 0)
	for _, lesson := range planet.Lessons {
		pointDetails = append(pointDetails, flattenPlanetPointDetails([]CourseLessonTree{lesson}, lesson.Lesson.Title)...)
	}

	for _, detail := range pointDetails {
		if normalizeTodayTaskStatus(detail.Point.Point.Status) == "completed" {
			continue
		}
		pointID := detail.Point.Point.ID
		planetID := planet.Planet.ID
		pointTitle := strings.TrimSpace(detail.Point.Point.Title)
		if pointTitle == "" {
			pointTitle = "다음 포인트"
		}

		return &TodayTask{
			Kind:        TodayTaskKindContinuePoint,
			Title:       fmt.Sprintf("%s 이어가기", pointTitle),
			Description: fmt.Sprintf("%s에서 다음 학습 포인트를 이어가세요.", planet.Planet.Title),
			CTALabel:    "바로 이어가기",
			Href:        fmt.Sprintf("/dashboard/planets/learning/%s/points/%s", planetID.String(), pointID.String()),
			PlanetID:    &planetID,
			PointID:     &pointID,
		}
	}

	return nil
}

func buildStartPlanetTodayTask(planet *PlanetAggregate) *TodayTask {
	planetID := planet.Planet.ID
	return &TodayTask{
		Kind:        TodayTaskKindStartPlanet,
		Title:       fmt.Sprintf("%s 시작하기", planet.Planet.Title),
		Description: "계획이 준비되었습니다. 첫 학습 포인트부터 가볍게 출발해 보세요.",
		CTALabel:    "학습 시작하기",
		Href:        fmt.Sprintf("/dashboard/planets/learning/%s", planetID.String()),
		PlanetID:    &planetID,
	}
}

func buildCompletePlanetTodayTask(planet *PlanetAggregate) *TodayTask {
	planetID := planet.Planet.ID
	return &TodayTask{
		Kind:        TodayTaskKindCompletePlanet,
		Title:       fmt.Sprintf("%s 마무리하기", planet.Planet.Title),
		Description: "필수 포인트가 완료되었습니다. 기록을 정리하고 코스를 마무리할 수 있어요.",
		CTALabel:    "코스 완료하기",
		Href:        fmt.Sprintf("/dashboard/planets/learning/%s", planetID.String()),
		PlanetID:    &planetID,
	}
}

func buildReviewRecordsTodayTask(planet PlanetListItem) *TodayTask {
	planetID := planet.ID
	return &TodayTask{
		Kind:        TodayTaskKindReviewRecords,
		Title:       fmt.Sprintf("%s 기록 돌아보기", planet.Title),
		Description: "완료한 학습 기록을 다시 보며 다음 여정을 고를 수 있어요.",
		CTALabel:    "기록 보기",
		Href:        fmt.Sprintf("/dashboard/planets/shared/%s", planetID.String()),
		PlanetID:    &planetID,
	}
}

func buildCreateCourseTodayTask() *TodayTask {
	return &TodayTask{
		Kind:        TodayTaskKindCreateCourse,
		Title:       "새로운 학습 우주 만들기",
		Description: "관심 있는 취미를 입력하면 오늘 시작할 수 있는 학습 흐름을 만들어 드릴게요.",
		CTALabel:    "코스 만들기",
		Href:        "/dashboard/goal",
	}
}

func sortPlanetListItemsByRecentActivity(planets []PlanetListItem) {
	sort.SliceStable(planets, func(i, j int) bool {
		return todayTaskRecentAt(planets[i]).After(todayTaskRecentAt(planets[j]))
	})
}

func todayTaskRecentAt(planet PlanetListItem) time.Time {
	if planet.LastAccessedAt != nil {
		return *planet.LastAccessedAt
	}
	return planet.UpdatedAt
}

func normalizeTodayTaskStatus[T ~string](status T) string {
	return strings.TrimSpace(strings.ToLower(string(status)))
}
