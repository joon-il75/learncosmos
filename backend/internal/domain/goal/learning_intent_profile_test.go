package goal

import (
	"strings"
	"testing"
)

func TestBuildLearningIntentProfileWithSourceUsesStructuredGoalFields(t *testing.T) {
	confirmed := "가죽공예를 기초부터 배워서 나만의 손 지갑을 만들어 보자"
	motivation := "처음부터 따라 하며 완성 경험을 얻고 싶다"
	usage := "주말에 실습 중심으로 배우고 싶다"
	goalType := "artifact_creation"
	outputType := "손 지갑"
	level := "초보"
	profile := &GoalProfile{
		ConfirmedGoal:   &confirmed,
		Motivation:      &motivation,
		UsageContext:    &usage,
		GoalType:        &goalType,
		OutputType:      &outputType,
		DifficultyLevel: &level,
	}

	intent := BuildLearningIntentProfileWithSource(profile, LearningIntentSourceBackfill)

	if intent.Source != LearningIntentSourceBackfill {
		t.Fatalf("source = %q, want %q", intent.Source, LearningIntentSourceBackfill)
	}
	if intent.ConfirmedGoal != confirmed {
		t.Fatalf("confirmed goal = %q, want %q", intent.ConfirmedGoal, confirmed)
	}
	if intent.LearnerLevel != "beginner" {
		t.Fatalf("learner level = %q, want beginner", intent.LearnerLevel)
	}
	if intent.DesiredOutput != "손 지갑" {
		t.Fatalf("desired output = %q, want 손 지갑", intent.DesiredOutput)
	}
	if !strings.Contains(intent.Purpose, motivation) || !strings.Contains(intent.Purpose, usage) {
		t.Fatalf("purpose %q does not include structured context", intent.Purpose)
	}
	if len(intent.PreferredActivities) == 0 {
		t.Fatalf("preferred activities should not be empty")
	}
	if len(intent.SuccessCriteria) == 0 {
		t.Fatalf("success criteria should not be empty")
	}
	if intent.Confidence != "high" {
		t.Fatalf("confidence = %q, want high", intent.Confidence)
	}
}

func TestBuildLearningIntentProfileDerivesShortPerformanceOutput(t *testing.T) {
	confirmed := "통기타를 취미로 꾸준히 연주할 수 있게 된다"
	profile := &GoalProfile{ConfirmedGoal: &confirmed}

	intent := BuildLearningIntentProfileWithSource(profile, LearningIntentSourceBackfill)

	if intent.DesiredOutput != "연주" {
		t.Fatalf("desired output = %q, want 연주", intent.DesiredOutput)
	}
	if len(intent.SuccessCriteria) == 0 || intent.SuccessCriteria[0] != "연주 완성" {
		t.Fatalf("success criteria = %#v, want first item 연주 완성", intent.SuccessCriteria)
	}
}

func TestActionableLearningIntentProfileRequiresConcreteSignal(t *testing.T) {
	thin := LearningIntentProfile{LearnerLevel: "beginner", PreferredActivities: []string{"기초 따라하기"}, Source: LearningIntentSourceBackfill}
	if IsActionableLearningIntentProfile(thin) {
		t.Fatalf("thin level-only profile should not be actionable")
	}
	concrete := LearningIntentProfile{Purpose: "취미 연주", Source: LearningIntentSourceBackfill}
	if !IsActionableLearningIntentProfile(concrete) {
		t.Fatalf("profile with purpose should be actionable")
	}
}

func TestBuildLearningIntentProfileLeavesEmptyProfileEmpty(t *testing.T) {
	intent := BuildLearningIntentProfileWithSource(&GoalProfile{UserIntent: "기타"}, LearningIntentSourceBackfill)
	if !IsEmptyLearningIntentProfile(intent) {
		t.Fatalf("intent = %#v, want empty profile", intent)
	}
	if intent.Source != "" || intent.Confidence != "" {
		t.Fatalf("empty profile should not keep source/confidence: %#v", intent)
	}
}

func TestNormalizeLearningIntentProfileBoundsAndDeduplicates(t *testing.T) {
	longGoal := strings.Repeat("가", 200)
	intent := NormalizeLearningIntentProfile(LearningIntentProfile{
		ConfirmedGoal:       longGoal,
		LearnerLevel:        "초급",
		CurrentBlockers:     []string{" 도구 선택 ", "도구 선택", strings.Repeat("나", 80)},
		PreferredActivities: []string{"실습", "실습", "따라하기", "프로젝트", "초과"},
		SuccessCriteria:     []string{"완성", "완성"},
		Source:              "bad-source",
		Confidence:          "bad-confidence",
	})

	if len([]rune(intent.ConfirmedGoal)) != 160 {
		t.Fatalf("confirmed goal rune length = %d, want 160", len([]rune(intent.ConfirmedGoal)))
	}
	if intent.LearnerLevel != "beginner" {
		t.Fatalf("learner level = %q, want beginner", intent.LearnerLevel)
	}
	if len(intent.CurrentBlockers) != 2 {
		t.Fatalf("current blockers = %#v, want 2 deduped items", intent.CurrentBlockers)
	}
	if len(intent.PreferredActivities) != 3 {
		t.Fatalf("preferred activities = %#v, want capped 3 items", intent.PreferredActivities)
	}
	if intent.Source != LearningIntentSourceGoalChat {
		t.Fatalf("source = %q, want %q for non-empty profile", intent.Source, LearningIntentSourceGoalChat)
	}
	if intent.Confidence == "" {
		t.Fatalf("confidence should be inferred")
	}
}

func TestBeginGoalRevisionSnapshotsLearningIntent(t *testing.T) {
	goal := "기타 한 곡을 끝까지 연주하기"
	profile := &GoalProfile{
		UserIntent:     "기타 배우기",
		ConfirmedGoal:  &goal,
		InterviewState: StateConfirmed,
		Version:        1,
		LearningIntent: LearningIntentProfile{
			ConfirmedGoal: goal,
			LearnerLevel:  "beginner",
			Source:        LearningIntentSourceGoalChat,
			Confidence:    "medium",
		},
	}

	BeginGoalRevision(profile, "목표를 조금 바꿀게요")

	if profile.RevisionSnapshot == nil {
		t.Fatalf("revision snapshot is nil")
	}
	if profile.RevisionSnapshot.LearningIntent.ConfirmedGoal != goal {
		t.Fatalf("snapshot learning intent goal = %q, want %q", profile.RevisionSnapshot.LearningIntent.ConfirmedGoal, goal)
	}
	if profile.RevisionSnapshot.LearningIntent.Source != LearningIntentSourceGoalChat {
		t.Fatalf("snapshot source = %q, want %q", profile.RevisionSnapshot.LearningIntent.Source, LearningIntentSourceGoalChat)
	}
}
