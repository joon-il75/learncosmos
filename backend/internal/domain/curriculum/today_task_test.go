package curriculum

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildTodayTaskFromPlanetAggregatesContinuesFirstIncompletePoint(t *testing.T) {
	planetID := uuid.New()
	completedPointID := uuid.New()
	nextPointID := uuid.New()

	task := buildTodayTaskFromPlanetAggregates([]*PlanetAggregate{{
		Planet: PlanetListItem{
			ID:     planetID,
			Title:  "드로잉",
			Status: "learning",
		},
		Lessons: []CourseLessonTree{{
			Lesson: CourseLesson{ID: uuid.New(), Title: "기초"},
			Points: []CoursePointAggregate{
				{Point: CoursePoint{ID: completedPointID, Title: "선 긋기", Status: "completed"}},
				{Point: CoursePoint{ID: nextPointID, Title: "명암 연습"}},
			},
		}},
	}}, nil)

	if task == nil {
		t.Fatal("expected today task")
	}
	if task.Kind != TodayTaskKindContinuePoint {
		t.Fatalf("expected %q task, got %q", TodayTaskKindContinuePoint, task.Kind)
	}
	if task.PointID == nil || *task.PointID != nextPointID {
		t.Fatalf("expected next incomplete point id %s, got %#v", nextPointID, task.PointID)
	}
}

func TestBuildTodayTaskFromPlanetAggregatesStartsReadyPlanet(t *testing.T) {
	planetID := uuid.New()

	task := buildTodayTaskFromPlanetAggregates([]*PlanetAggregate{{
		Planet: PlanetListItem{
			ID:     planetID,
			Title:  "피아노",
			Status: "ready",
		},
	}}, nil)

	if task == nil {
		t.Fatal("expected today task")
	}
	if task.Kind != TodayTaskKindStartPlanet {
		t.Fatalf("expected %q task, got %q", TodayTaskKindStartPlanet, task.Kind)
	}
	if task.PlanetID == nil || *task.PlanetID != planetID {
		t.Fatalf("expected planet id %s, got %#v", planetID, task.PlanetID)
	}
}

func TestBuildTodayTaskFromPlanetAggregatesFallsBackToCreateCourse(t *testing.T) {
	task := buildTodayTaskFromPlanetAggregates(nil, nil)

	if task == nil {
		t.Fatal("expected today task")
	}
	if task.Kind != TodayTaskKindCreateCourse {
		t.Fatalf("expected %q task, got %q", TodayTaskKindCreateCourse, task.Kind)
	}
}
