package safety

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestObserveAIOutputDoesNotTreatBlockAsError(t *testing.T) {
	service := NewService(nil)

	result := service.ObserveAIOutput(context.Background(), ModerateInput{
		TargetType: TargetAIPointQuestionOutput,
		Text:       "how to kill",
		Locale:     "en",
		Route:      "test",
	})

	if result == nil {
		t.Fatal("expected moderation result")
	}
	if result.Action != ActionBlock {
		t.Fatalf("expected block result for builtin unsafe phrase, got %q", result.Action)
	}
	if result.RiskType != "unsafe_instruction" {
		t.Fatalf("expected unsafe_instruction risk, got %q", result.RiskType)
	}
}

func TestObserveAIOutputCopiesMetadataAndMarksDirection(t *testing.T) {
	service := NewService(nil)
	metadata := map[string]any{"source_feature": "point_question", "direction": "user_input"}
	input := ModerateInput{
		TargetType: TargetAIPointQuestionOutput,
		Text:       "normal study question",
		Locale:     "en",
		Route:      "test",
		Metadata:   metadata,
	}

	result := service.ObserveAIOutput(context.Background(), input)

	if result == nil || result.Action != ActionAllow {
		t.Fatalf("expected allow result, got %#v", result)
	}
	if metadata["direction"] != "user_input" {
		t.Fatalf("expected original metadata not to be mutated, got %v", metadata["direction"])
	}
	marked := aiOutputMetadata(metadata)
	if marked["direction"] != "ai_output" {
		t.Fatalf("expected ai_output direction, got %v", marked["direction"])
	}
	if marked["source_feature"] != "point_question" {
		t.Fatalf("expected source_feature to be preserved, got %v", marked["source_feature"])
	}
}

func TestObserveAIOutputNilServiceAllowsEmptyInput(t *testing.T) {
	var service *Service

	result := service.ObserveAIOutput(context.Background(), ModerateInput{})

	if result == nil {
		t.Fatal("expected moderation result")
	}
	if result.Action != ActionAllow {
		t.Fatalf("expected allow result, got %q", result.Action)
	}
}

func TestAIOutputTargetTypesAreStable(t *testing.T) {
	expected := []string{
		"ai_goal_interview_output",
		"ai_course_draft_output",
		"ai_point_summary_output",
		"ai_point_question_output",
		"ai_point_feedback_output",
		"ai_point_self_evaluation_output",
	}
	if len(AIOutputTargetTypes) != len(expected) {
		t.Fatalf("expected %d AI output target types, got %d", len(expected), len(AIOutputTargetTypes))
	}
	seen := map[string]bool{}
	for index, targetType := range AIOutputTargetTypes {
		if targetType != expected[index] {
			t.Fatalf("target type index %d: expected %q, got %q", index, expected[index], targetType)
		}
		if targetType == "" {
			t.Fatalf("target type index %d is empty", index)
		}
		if seen[targetType] {
			t.Fatalf("duplicate target type %q", targetType)
		}
		seen[targetType] = true
	}
}

func TestBuildAIOutputSummaryNoData(t *testing.T) {
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	summary := BuildAIOutputSummary(now, nil)
	if len(summary.Windows) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(summary.Windows))
	}
	for _, window := range summary.Windows {
		if window.Total != 0 || window.Severity != "none" {
			t.Fatalf("expected empty %s window to be none, got total=%d severity=%s", window.Label, window.Total, window.Severity)
		}
	}
}

func TestBuildAIOutputSummaryCautionFromRepeatedTargetRisk(t *testing.T) {
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	observations := []AIOutputLogObservation{
		{TargetType: TargetAIPointQuestionOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_question", CreatedAt: now.Add(-time.Hour)},
		{TargetType: TargetAIPointQuestionOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_question", CreatedAt: now.Add(-2 * time.Hour)},
		{TargetType: TargetAIPointQuestionOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_question", CreatedAt: now.Add(-3 * time.Hour)},
	}
	summary := BuildAIOutputSummary(now, observations)
	window := summary.Windows[0]
	if window.Severity != "warning" {
		t.Fatalf("expected repeated source feature risk to be warning, got %s", window.Severity)
	}
	if window.RiskyTotal != 3 || window.RiskyRate != 100 {
		t.Fatalf("unexpected risk totals: total=%d rate=%v", window.RiskyTotal, window.RiskyRate)
	}
}

func TestBuildAIOutputSummaryCriticalFromRepeatedRule(t *testing.T) {
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	observations := make([]AIOutputLogObservation, 0, 10)
	for index := 0; index < 10; index++ {
		observations = append(observations, AIOutputLogObservation{TargetType: TargetAICourseDraftOutput, Action: ActionBlock, RiskType: "unsafe_instruction", SourceFeature: "course_draft", MatchedRuleID: "rule-1", MatchedPattern: "danger", CreatedAt: now.Add(-time.Duration(index) * time.Hour)})
	}
	summary := BuildAIOutputSummary(now, observations)
	if summary.Windows[0].Severity != "critical" {
		t.Fatalf("expected critical severity, got %s", summary.Windows[0].Severity)
	}
	if len(summary.Windows[0].MatchedRules) == 0 || summary.Windows[0].MatchedRules[0].Count != 10 {
		t.Fatalf("expected matched rule count 10, got %#v", summary.Windows[0].MatchedRules)
	}
}

func TestBuildAIOutputReviewCandidatesNoRiskyData(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	observations := []AIOutputLogObservation{
		{TargetType: TargetAIPointQuestionOutput, Action: ActionAllow, RiskType: "none", SourceFeature: "point_question", CreatedAt: now.Add(-time.Hour)},
	}

	result := BuildAIOutputReviewCandidates(now, 24*30, observations)

	if result.TotalCandidates != 0 {
		t.Fatalf("expected no review candidates, got %d", result.TotalCandidates)
	}
	if len(result.Candidates) != 0 {
		t.Fatalf("expected empty candidate list, got %#v", result.Candidates)
	}
}

func TestBuildAIOutputReviewCandidatesGroupsRiskyObservations(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	targetID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	observations := []AIOutputLogObservation{
		{TargetType: TargetAIPointFeedbackOutput, TargetID: &targetID, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_feedback", MatchedRuleID: "rule-1", MatchedPattern: "danger", CreatedAt: now.Add(-2 * time.Hour)},
		{TargetType: TargetAIPointFeedbackOutput, TargetID: &targetID, Action: ActionBlock, RiskType: "unsafe_instruction", SourceFeature: "point_feedback", MatchedRuleID: "rule-1", MatchedPattern: "danger", CreatedAt: now.Add(-time.Hour)},
		{TargetType: TargetAIPointFeedbackOutput, Action: ActionAllow, RiskType: "none", SourceFeature: "point_feedback", CreatedAt: now.Add(-30 * time.Minute)},
	}

	result := BuildAIOutputReviewCandidates(now, 24*30, observations)

	if result.TotalCandidates != 1 || len(result.Candidates) != 1 {
		t.Fatalf("expected one grouped candidate, got %#v", result.Candidates)
	}
	candidate := result.Candidates[0]
	if candidate.Total != 2 || candidate.ActionCounts.SoftWarn != 1 || candidate.ActionCounts.Block != 1 || candidate.ActionCounts.Allow != 0 {
		t.Fatalf("unexpected candidate action counts: %#v", candidate)
	}
	if candidate.TargetType != TargetAIPointFeedbackOutput || candidate.RiskType != "unsafe_instruction" || candidate.SourceFeature != "point_feedback" {
		t.Fatalf("unexpected candidate dimensions: %#v", candidate)
	}
	if candidate.MatchedRuleID != "rule-1" || candidate.MatchedPattern != "danger" {
		t.Fatalf("unexpected matched rule fields: %#v", candidate)
	}
	if candidate.SampleTargetID == nil || *candidate.SampleTargetID != targetID {
		t.Fatalf("expected sample target id %s, got %#v", targetID, candidate.SampleTargetID)
	}
	if !candidate.LatestSeenAt.Equal(now.Add(-time.Hour)) {
		t.Fatalf("expected latest seen at latest risky observation, got %s", candidate.LatestSeenAt)
	}
}

func TestBuildAIOutputReviewCandidatesSortsByTotalThenLatest(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	observations := []AIOutputLogObservation{
		{TargetType: TargetAIPointQuestionOutput, Action: ActionBlock, RiskType: "unsafe_instruction", SourceFeature: "point_question", MatchedRuleID: "rule-b", CreatedAt: now.Add(-time.Hour)},
		{TargetType: TargetAIPointSummaryOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_summary", MatchedRuleID: "rule-a", CreatedAt: now.Add(-3 * time.Hour)},
		{TargetType: TargetAIPointSummaryOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_summary", MatchedRuleID: "rule-a", CreatedAt: now.Add(-2 * time.Hour)},
	}

	result := BuildAIOutputReviewCandidates(now, 24*30, observations)

	if len(result.Candidates) != 2 {
		t.Fatalf("expected two candidates, got %#v", result.Candidates)
	}
	if result.Candidates[0].TargetType != TargetAIPointSummaryOutput || result.Candidates[0].Total != 2 {
		t.Fatalf("expected highest total candidate first, got %#v", result.Candidates)
	}
}

func TestAIOutputReviewDecisionDefaultsToUnreviewed(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	observations := []AIOutputLogObservation{
		{TargetType: TargetAIPointQuestionOutput, Action: ActionSoftWarn, RiskType: "unsafe_instruction", SourceFeature: "point_question", MatchedRuleID: "rule-1", CreatedAt: now.Add(-time.Hour)},
	}

	result := BuildAIOutputReviewCandidates(now, 24*30, observations)

	if len(result.Candidates) != 1 {
		t.Fatalf("expected one candidate, got %#v", result.Candidates)
	}
	review := result.Candidates[0].Review
	if review.Status != AIOutputReviewStatusUnreviewed {
		t.Fatalf("expected unreviewed status, got %q", review.Status)
	}
	if review.CandidateKey != result.Candidates[0].Key {
		t.Fatalf("expected review key to match candidate key")
	}
}

func TestApplyAIOutputReviewDecisionsMergesStoredDecision(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	result := BuildAIOutputReviewCandidates(now, 24*30, []AIOutputLogObservation{
		{TargetType: TargetAIPointSummaryOutput, Action: ActionBlock, RiskType: "unsafe_instruction", SourceFeature: "point_summary", MatchedRuleID: "rule-2", CreatedAt: now.Add(-time.Hour)},
	})
	if len(result.Candidates) != 1 {
		t.Fatalf("expected one candidate, got %#v", result.Candidates)
	}
	key := result.Candidates[0].Key
	reviewedAt := now.Add(time.Minute)
	ApplyAIOutputReviewDecisions(&result, map[string]AIOutputReviewDecision{
		key: {
			CandidateKey: key,
			Status:       AIOutputReviewStatusFalsePositive,
			Note:         "정상 문맥",
			ReviewedBy:   "admin-1",
			ReviewedAt:   &reviewedAt,
			Metadata:     map[string]any{"source": "test"},
		},
	})

	review := result.Candidates[0].Review
	if review.Status != AIOutputReviewStatusFalsePositive || review.Note != "정상 문맥" || review.ReviewedBy != "admin-1" {
		t.Fatalf("unexpected merged review: %#v", review)
	}
	if review.ReviewedAt == nil || !review.ReviewedAt.Equal(reviewedAt) {
		t.Fatalf("expected reviewed_at to merge, got %#v", review.ReviewedAt)
	}
}

func TestValidateAIOutputReviewDecisionInputRejectsInvalidStatus(t *testing.T) {
	_, err := validateAIOutputReviewDecisionInput(AIOutputReviewDecisionInput{
		CandidateKey: "candidate-1",
		Status:       "maybe",
	})
	if err == nil {
		t.Fatal("expected invalid decision error")
	}
	if err != ErrInvalidAIOutputReviewDecision {
		t.Fatalf("expected ErrInvalidAIOutputReviewDecision, got %v", err)
	}
}

func TestValidateAIOutputReviewDecisionInputCleansFields(t *testing.T) {
	input, err := validateAIOutputReviewDecisionInput(AIOutputReviewDecisionInput{
		CandidateKey: " candidate-1 ",
		Status:       " resolved ",
		Note:         " done ",
		ReviewedBy:   " admin-1 ",
	})
	if err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	if input.CandidateKey != "candidate-1" || input.Status != "resolved" || input.Note != "done" || input.ReviewedBy != "admin-1" {
		t.Fatalf("input was not cleaned: %#v", input)
	}
	if input.Metadata == nil {
		t.Fatal("expected metadata map")
	}
}

func TestSanitizeSafetyMetadataAllowlist(t *testing.T) {
	metadata := sanitizeSafetyMetadata(map[string]any{
		"source_feature":   " point_question ",
		"direction":        "ai_output",
		"content_length":   42,
		"prompt":           "private learner prompt",
		"api_key":          "sk-secret",
		"provider_payload": map[string]any{"token": "secret"},
	})
	if metadata["source_feature"] != "point_question" || metadata["direction"] != "ai_output" || metadata["content_length"] != 42 {
		t.Fatalf("expected safe metadata fields to remain, got %#v", metadata)
	}
	for _, forbidden := range []string{"prompt", "api_key", "provider_payload"} {
		if _, ok := metadata[forbidden]; ok {
			t.Fatalf("metadata leaked forbidden key %q: %#v", forbidden, metadata)
		}
	}
}

func TestValidateAIOutputReviewDecisionInputSanitizesMetadata(t *testing.T) {
	input, err := validateAIOutputReviewDecisionInput(AIOutputReviewDecisionInput{
		CandidateKey: "candidate-1",
		Status:       "resolved",
		Metadata: map[string]any{
			"source":     "super_admin_safety",
			"raw_prompt": "private prompt",
		},
	})
	if err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
	if input.Metadata["source"] != "super_admin_safety" {
		t.Fatalf("expected source metadata to remain, got %#v", input.Metadata)
	}
	if _, ok := input.Metadata["raw_prompt"]; ok {
		t.Fatalf("expected raw_prompt metadata to be removed: %#v", input.Metadata)
	}
}

func TestApplyAIOutputReviewDecisionsSanitizesStoredMetadata(t *testing.T) {
	now := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	result := BuildAIOutputReviewCandidates(now, 24*30, []AIOutputLogObservation{
		{TargetType: TargetAIPointQuestionOutput, Action: ActionBlock, RiskType: "unsafe_instruction", SourceFeature: "point_question", CreatedAt: now.Add(-time.Hour)},
	})
	key := result.Candidates[0].Key
	ApplyAIOutputReviewDecisions(&result, map[string]AIOutputReviewDecision{
		key: {
			CandidateKey: key,
			Status:       AIOutputReviewStatusResolved,
			Metadata:     map[string]any{"source": "test", "raw_payload": "secret"},
		},
	})
	metadata := result.Candidates[0].Review.Metadata
	if metadata["source"] != "test" {
		t.Fatalf("expected source metadata to remain, got %#v", metadata)
	}
	if _, ok := metadata["raw_payload"]; ok {
		t.Fatalf("expected raw_payload metadata to be removed: %#v", metadata)
	}
}
