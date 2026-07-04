package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
)

func TestBuildRefreshDecisionAppliesProfileOverlayToFallbackSpec(t *testing.T) {
	row := refreshLessonRow{
		LessonID:    uuid.New(),
		DraftID:     uuid.New(),
		SourceQuery: "가죽공예",
		Language:    "ko",
		Title:       "도구와 재료 준비하기",
		OrderIndex:  0,
		RawSpec: []byte(`{
			"primary_query":"가죽공예 입문 도구 사용법 기초 튜토리얼",
			"intent":"tutorial",
			"must_include":["가죽공예","도구","기초"],
			"nice_to_have":["입문"],
			"content_types":["video"],
			"language":"ko",
			"stage_role":"setup_intro",
			"source":"fallback"
		}`),
		LearningIntent: curriculum.LearningIntentProfile{
			DesiredOutput:       "카드지갑",
			PreferredActivities: []string{"따라 만들기"},
		},
	}

	decision := buildRefreshDecision(row)
	if decision.Skip != "" {
		t.Fatalf("skip = %q, want empty", decision.Skip)
	}
	if !decision.Changed {
		t.Fatal("expected changed decision")
	}
	if decision.Spec.Source != curriculum.SearchSpecSourceFallbackOverlay {
		t.Fatalf("source = %q, want fallback overlay", decision.Spec.Source)
	}
	if !containsString(decision.Spec.MustInclude, "카드지갑") {
		t.Fatalf("must_include = %#v, want 카드지갑", decision.Spec.MustInclude)
	}
}

func TestBuildRefreshDecisionProtectsManualGeneratedAndOverlaySpecs(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lessonType string
		wantSkip   string
	}{
		{name: "manual spec", source: curriculum.SearchSpecSourceManual, wantSkip: "manual_spec"},
		{name: "manual lesson", source: curriculum.SearchSpecSourceFallback, lessonType: string(curriculum.LessonSourceManual), wantSkip: "manual_lesson"},
		{name: "generated", source: curriculum.SearchSpecSourceGenerated, wantSkip: "generated"},
		{name: "fallback overlay", source: curriculum.SearchSpecSourceFallbackOverlay, wantSkip: "already_overlay"},
		{name: "template overlay", source: curriculum.SearchSpecSourcePatternTemplateOverlay, wantSkip: "already_overlay"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := refreshLessonRow{
				LessonID:         uuid.New(),
				DraftID:          uuid.New(),
				SourceQuery:      "가죽공예",
				Language:         "ko",
				Title:            "도구와 재료 준비하기",
				LessonSourceType: tt.lessonType,
				RawSpec:          []byte(`{"primary_query":"가죽공예 기초 튜토리얼","intent":"tutorial","source":"` + tt.source + `"}`),
				LearningIntent:   curriculum.LearningIntentProfile{DesiredOutput: "카드지갑"},
			}
			decision := buildRefreshDecision(row)
			if decision.Skip != tt.wantSkip {
				t.Fatalf("skip = %q, want %q", decision.Skip, tt.wantSkip)
			}
		})
	}
}

func TestBuildRefreshDecisionSkipsEmptyProfile(t *testing.T) {
	decision := buildRefreshDecision(refreshLessonRow{
		LessonID:    uuid.New(),
		DraftID:     uuid.New(),
		SourceQuery: "가죽공예",
		Language:    "ko",
		Title:       "도구와 재료 준비하기",
		RawSpec:     []byte(`{"primary_query":"가죽공예 기초 튜토리얼","intent":"tutorial","source":"fallback"}`),
	})
	if decision.Skip != "no_profile" {
		t.Fatalf("skip = %q, want no_profile", decision.Skip)
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestBuildRefreshDecisionTreatsSentenceLikeSuccessCriteriaAsUnchanged(t *testing.T) {
	row := refreshLessonRow{
		LessonID:    uuid.New(),
		DraftID:     uuid.New(),
		SourceQuery: "악기",
		Language:    "ko",
		Title:       "목표 연주를 끝까지 이어갈 준비를 한다",
		OrderIndex:  0,
		RawSpec: []byte(`{
			"primary_query":"곡 연습 구간 연결 리듬 반복 연습 템포 튜토리얼",
			"intent":"tutorial",
			"must_include":["곡 연습","구간 연결","리듬"],
			"nice_to_have":["반복 연습","템포","반주 맞추기"],
			"content_types":["video"],
			"language":"ko",
			"stage_role":"core_pattern",
			"source":"pattern_template"
		}`),
		LearningIntent: curriculum.LearningIntentProfile{
			SuccessCriteria: []string{"3개월 안에 좋아하는 곡을 연주한다"},
		},
	}

	decision := buildRefreshDecision(row)
	if decision.Skip != "" {
		t.Fatalf("skip = %q, want empty", decision.Skip)
	}
	if decision.Changed {
		t.Fatalf("expected unchanged decision, got %#v", decision.Spec)
	}
}
