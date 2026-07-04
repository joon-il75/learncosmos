package curriculum

import (
	"strings"
	"testing"
)

func TestBuildCurriculumPatternQueryText(t *testing.T) {
	goal := "편집한 영상을 포트폴리오와 유튜브에 공개하고 싶다"
	usage := "유튜브 업로드"
	output := "포트폴리오 영상"
	got := BuildCurriculumPatternQueryText(CreateCourseDraftRequest{
		SourceQuery:      "영상편집 배우기",
		LearningGoal:     &goal,
		GoalUsageContext: &usage,
		GoalOutputType:   &output,
	})

	for _, want := range []string{"영상편집 배우기", goal, usage, output} {
		if !strings.Contains(got, want) {
			t.Fatalf("query text missing %q: %q", want, got)
		}
	}
}

func TestFormatFloat32VectorForPG(t *testing.T) {
	got := FormatFloat32VectorForPG([]float32{0.5, -1.25})
	if got != "[0.500000,-1.250000]" {
		t.Fatalf("vector = %q", got)
	}
}

func TestBuildCurriculumPatternMatchGuidance(t *testing.T) {
	got := BuildCurriculumPatternMatchGuidance([]CurriculumPatternMatch{
		{
			PatternKey:              "instrument_performance:performance_execution:song_completion",
			Title:                   "한 곡 완주형",
			Summary:                 "목표 곡 완주 흐름에 초점을 둔다.",
			RecommendedSequence:     []string{"setup_intro", "first_output", "performance_prep"},
			StageRules:              []string{"기술 소개보다 목표 곡의 구간 연결을 우선합니다."},
			CompletionCriteriaRules: []string{"실제 공연은 completion_criteria로 분리합니다."},
			BadPatterns:             []string{"기본 자세", "스케일 소개"},
			RecommendationSearchSpecTemplate: CurriculumPatternRecommendationSearchSpecTemplate{
				Intent:       "tutorial",
				ContentTypes: []string{"video"},
				Source:       SearchSpecSourcePatternTemplate,
				StageTemplates: []LessonRecommendationSearchSpec{
					{StageRole: "setup_intro", MustInclude: []string{"악기", "도구"}, NiceToHave: []string{"기초"}},
				},
			},
		},
	}, 1)

	for _, want := range []string{
		"한 곡 완주형",
		"instrument_performance:performance_execution:song_completion",
		"setup_intro -> first_output -> performance_prep",
		"기술 소개보다 목표 곡의 구간 연결을 우선합니다.",
		"실제 공연은 completion_criteria로 분리합니다.",
		"기본 자세",
		"추천 검색 명세 template",
		"setup_intro(악기, 도구, 기초)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("guidance missing %q: %q", want, got)
		}
	}
}

func TestBuildCurriculumPatternMatchGuidanceForEnglish(t *testing.T) {
	got := BuildCurriculumPatternMatchGuidanceForLanguage([]CurriculumPatternMatch{
		{
			PatternKey:              "visual_art:artifact_creation:photography_fundamentals_project",
			Title:                   "Photography Fundamentals Project",
			Summary:                 "Use composition and light to complete a small photo project.",
			RecommendedSequence:     []string{"setup_intro", "core_pattern", "artifact_finish"},
			StageRules:              []string{"Design gateway stages around actual photo practice."},
			CompletionCriteriaRules: []string{"Keep public exhibition as completion criteria."},
			BadPatterns:             []string{"ISO definition only"},
		},
	}, 1, "en")

	for _, want := range []string{
		"Summary:",
		"Recommended flow:",
		"Stage rule:",
		"Completion criteria rule:",
		"Lesson patterns to avoid:",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("English guidance missing %q: %q", want, got)
		}
	}
	for _, forbidden := range []string{"요약", "권장 흐름", "단계 규칙", "완료 기준 규칙", "피할 리슨 예"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("English guidance leaked Korean label %q: %q", forbidden, got)
		}
	}
}

func TestHashCurriculumPatternMatchQueryStable(t *testing.T) {
	first := hashCurriculumPatternMatchQuery("목표", "ko", "embedding_gemma", "embedding-gemma", 768)
	second := hashCurriculumPatternMatchQuery("목표", "ko", "embedding_gemma", "embedding-gemma", 768)
	if first == "" || first != second {
		t.Fatalf("hash not stable: %q %q", first, second)
	}
	changed := hashCurriculumPatternMatchQuery("목표", "en", "embedding_gemma", "embedding-gemma", 768)
	if changed == first {
		t.Fatalf("hash should include language")
	}
}

func TestFilterCurriculumPatternGuidanceMatches(t *testing.T) {
	matches := []CurriculumPatternMatch{
		{PatternKey: "digital_creation:presentation_publish:portfolio_publish", Domain: "digital_creation", GoalType: "presentation_publish", SimilarityScore: 0.44},
		{PatternKey: "craft_making:teaching_instruction:craft_teach_youtube", Domain: "craft_making", GoalType: "teaching_instruction", SimilarityScore: 0.35},
		{PatternKey: "visual_art:presentation_publish:calligraphy_publish_showcase", Domain: "visual_art", GoalType: "presentation_publish", SimilarityScore: 0.34},
	}

	got := FilterCurriculumPatternGuidanceMatches(matches, "digital_creation", "presentation_publish", "digital_creation:presentation_publish:portfolio_publish", 2, 0.30)
	if len(got) != 1 {
		t.Fatalf("filtered len = %d, want 1", len(got))
	}
	if got[0].PatternKey != "digital_creation:presentation_publish:portfolio_publish" {
		t.Fatalf("filtered key = %q", got[0].PatternKey)
	}

	fallback := FilterCurriculumPatternGuidanceMatches(matches, "visual_art", "presentation_publish", "", 2, 0.30)
	if len(fallback) != 1 {
		t.Fatalf("fallback len = %d, want 1", len(fallback))
	}
	if fallback[0].PatternKey != "visual_art:presentation_publish:calligraphy_publish_showcase" {
		t.Fatalf("fallback key = %q", fallback[0].PatternKey)
	}
}
