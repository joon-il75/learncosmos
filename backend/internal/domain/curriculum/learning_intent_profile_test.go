package curriculum

import (
	"strings"
	"testing"
)

func TestBuildLearningIntentPromptSummaryBoundsProfile(t *testing.T) {
	profile := LearningIntentProfile{
		LearnerLevel:        "beginner",
		Purpose:             strings.Repeat("목적", 80),
		DesiredOutput:       "손 지갑",
		CurrentBlockers:     []string{"도구 선택", "바느질 기초", "초과 항목"},
		PreferredActivities: []string{"프로젝트 실습", "기초 따라하기", "초과 항목"},
		SuccessCriteria:     []string{"손 지갑 완성", "기초 도구 사용 이해", "초과 항목"},
	}

	summary := buildLearningIntentPromptSummary(profile)

	if summary == "" {
		t.Fatalf("summary is empty")
	}
	if len([]rune(summary)) > 260 {
		t.Fatalf("summary length = %d, want <= 260", len([]rune(summary)))
	}
	if !strings.Contains(summary, "수준=beginner") || !strings.Contains(summary, "산출물=손 지갑") {
		t.Fatalf("summary missing key fields: %q", summary)
	}
	if strings.Contains(summary, "초과 항목") {
		t.Fatalf("summary should cap list items: %q", summary)
	}
}

func TestBuildGoalPromptDetailsIncludesLearningIntentSummary(t *testing.T) {
	goal := "가죽공예로 손 지갑을 완성한다"
	req := CreateCourseDraftRequest{
		SourceQuery:  goal,
		LearningGoal: &goal,
		LearningIntent: LearningIntentProfile{
			LearnerLevel:        "beginner",
			Purpose:             "기초부터 따라 하며 완성 경험을 얻기",
			DesiredOutput:       "손 지갑",
			PreferredActivities: []string{"프로젝트 실습"},
			SuccessCriteria:     []string{"손 지갑 완성"},
		},
	}

	details := buildGoalPromptDetails(req)
	if len(details) == 0 || !strings.HasPrefix(details[0], "학습 의도 요약: ") {
		t.Fatalf("first detail = %#v, want learning intent summary first", details)
	}
	if !strings.Contains(details[0], "산출물=손 지갑") || !strings.Contains(details[0], "선호활동=프로젝트 실습") {
		t.Fatalf("learning intent detail missing signals: %q", details[0])
	}
}

func TestCurriculumGenerationPromptDoesNotIncludeRawProfileMetadata(t *testing.T) {
	goal := "가죽공예로 손 지갑을 완성한다"
	req := CreateCourseDraftRequest{
		SourceQuery:  goal,
		LearningGoal: &goal,
		LearningIntent: LearningIntentProfile{
			LearnerLevel:        "beginner",
			Purpose:             "기초부터 따라 하며 완성 경험을 얻기",
			DesiredOutput:       "손 지갑",
			CurrentBlockers:     []string{"도구 선택"},
			PreferredActivities: []string{"프로젝트 실습"},
			SuccessCriteria:     []string{"손 지갑 완성"},
			Source:              "goal_chat",
			Confidence:          "high",
		},
	}

	prompt := buildCurriculumGenerationPrompt(req)
	if !strings.Contains(prompt, "학습 의도 요약:") {
		t.Fatalf("prompt missing learning intent summary: %q", prompt)
	}
	if strings.Contains(prompt, "goal_chat") || strings.Contains(prompt, "confidence") || strings.Contains(prompt, "confirmed_goal") {
		t.Fatalf("prompt includes raw profile metadata: %q", prompt)
	}
}
