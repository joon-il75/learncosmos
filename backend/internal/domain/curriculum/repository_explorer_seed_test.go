package curriculum

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildExplorerRegionSeed(t *testing.T) {
	mainID := uuid.New()
	subID := uuid.New()
	mainObjective := "도시 장면을 관찰하며 스케치 방향을 잡는다."
	subSummary := "재료 특성을 비교하고 기본 도구를 정리한다."

	seed := buildExplorerRegionSeed(DraftLessonTree{
		Lesson: CourseDraftLesson{
			ID:        mainID,
			Title:     "어반스케치로 도시 장면 관찰하기",
			Objective: &mainObjective,
		},
		SubLessons: []DraftLessonTree{
			{
				Lesson: CourseDraftLesson{
					ID:      subID,
					Title:   "어반스케치 재료 이해",
					Summary: &subSummary,
				},
			},
		},
	}, 0)

	if seed.ID != mainID {
		t.Fatalf("seed.ID = %s, want %s", seed.ID, mainID)
	}
	if seed.Description == nil || *seed.Description != mainObjective {
		t.Fatalf("seed.Description = %#v, want %q", seed.Description, mainObjective)
	}
	if len(seed.SubRegions) != 1 {
		t.Fatalf("sub region len = %d, want 1", len(seed.SubRegions))
	}
	if seed.SubRegions[0].ID != subID {
		t.Fatalf("sub seed ID = %s, want %s", seed.SubRegions[0].ID, subID)
	}
	if seed.SubRegions[0].Description == nil || *seed.SubRegions[0].Description != subSummary {
		t.Fatalf("sub seed Description = %#v, want %q", seed.SubRegions[0].Description, subSummary)
	}
}

func TestFirstNonEmptyString(t *testing.T) {
	blank := "   "
	value := "keep me"

	got := firstNonEmptyString(nil, &blank, &value)
	if got == nil || *got != value {
		t.Fatalf("firstNonEmptyString returned %#v, want %q", got, value)
	}
}
