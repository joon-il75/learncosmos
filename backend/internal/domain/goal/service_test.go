package goal

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestBuildFallbackInterviewResponseProposesGoalWhenContextIsEnough(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "어반스케치",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "user", Content: "멋진 작가가 되고 싶어"},
			{Role: "user", Content: "그림 실력을 키우고 싶어"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "이제까지 내용을 정리하면")
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateProposingGoal)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal == "" {
		t.Fatal("expected proposed goal")
	}
	if resp.Message == genericInterviewFallbackMessage {
		t.Fatal("expected meaningful response instead of generic fallback")
	}
}

func TestBuildFallbackInterviewResponseExplainsMissingContext(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "어반스케치",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "아직 목표가 부족한거야?")
	if resp.NextState != StateClarifying {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateClarifying)
	}
	if resp.Message == genericInterviewFallbackMessage {
		t.Fatal("expected missing-context guidance instead of generic fallback")
	}
}

func TestBuildFallbackInterviewResponseAdvancesAfterTeacherGoal(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "어반스케치",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "user", Content: "강사가 되고 싶어"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "강사가 되는게 목표야")
	if resp.Message == "좋아요. 어반스케치를 배워서 가장 이루고 싶은 변화는 무엇인가요? 예를 들면 더 잘 그리고 싶다, 나만의 작품을 만들고 싶다 같은 방향이 있어요." {
		t.Fatal("expected conversation to advance instead of repeating the first motivation question")
	}
	if resp.NextState != StateClarifying && resp.NextState != StateProposingGoal {
		t.Fatalf("unexpected next state: %q", resp.NextState)
	}
}

func TestInferExtractedInfoDetectsTeachingGoal(t *testing.T) {
	svc := NewService()
	info := svc.InferExtractedInfo("강사가 되고 싶어")
	if info.Motivation == nil || *info.Motivation == "" {
		t.Fatal("expected motivation to be extracted")
	}
	if info.GoalType == nil || *info.GoalType != "teaching" {
		t.Fatalf("expected teaching goal type, got %#v", info.GoalType)
	}
}

func TestShouldOverrideAIResponseWhenRepeatingLastLumiQuestion(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "어반스케치",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "lumi", Content: "좋아요. 어반스케치를 배워서 가장 이루고 싶은 변화는 무엇인가요?"},
		},
	}
	resp := &InterviewAIResponse{
		Message:   "좋아요. 어반스케치를 배워서 가장 이루고 싶은 변화는 무엇인가요?",
		NextState: StateClarifying,
	}
	if !svc.ShouldOverrideAIResponse(profile, "강사가 되고 싶어", resp) {
		t.Fatal("expected repeated lumi question to be overridden")
	}
}

func TestBuildFallbackInterviewResponseDoesNotRepeatFirstQuestionAfterUsageGoal(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "어반스케치",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "user", Content: "여행하면서 자유롭게 그림으로 기록을 남기고 싶어"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "여행하면서 자유롭게 그림으로 기록을 남기고 싶어")
	if resp.Message == "좋아요. 어반스케치를 배워서 가장 이루고 싶은 변화는 무엇인가요? 예를 들면 더 잘 그리고 싶다, 나만의 작품을 만들고 싶다 같은 방향이 있어요." {
		t.Fatal("expected conversation to advance instead of repeating the first motivation question")
	}
}

func TestBeginGoalRevisionIncrementsVersionAndClearsRebuildDecision(t *testing.T) {
	decision := RebuildAll
	profile := &GoalProfile{
		InterviewState:  StateConfirmed,
		RebuildDecision: &decision,
		Version:         1,
		Messages: []InterviewMessage{
			{Role: "lumi", Content: "기존 목표를 확인했어요."},
		},
	}

	BeginGoalRevision(profile, "목표를 조금 바꾸고 싶어요.")

	if profile.Version != 2 {
		t.Fatalf("version = %d, want 2", profile.Version)
	}
	if profile.InterviewState != StateRevisingGoal {
		t.Fatalf("state = %q, want %q", profile.InterviewState, StateRevisingGoal)
	}
	if profile.RebuildDecision != nil {
		t.Fatalf("rebuild decision = %v, want nil", *profile.RebuildDecision)
	}
	if profile.RevisionSnapshot == nil {
		t.Fatal("revision snapshot = nil, want previous goal state")
	}
	if profile.RevisionSnapshot.Version != 1 {
		t.Fatalf("snapshot version = %d, want 1", profile.RevisionSnapshot.Version)
	}
	if profile.RevisionSnapshot.InterviewState != StateConfirmed {
		t.Fatalf("snapshot state = %q, want %q", profile.RevisionSnapshot.InterviewState, StateConfirmed)
	}
	if len(profile.Messages) != 2 {
		t.Fatalf("messages len = %d, want 2", len(profile.Messages))
	}
	if profile.Messages[1].Role != "user" || profile.Messages[1].Content != "목표를 조금 바꾸고 싶어요." {
		t.Fatalf("unexpected appended message: %+v", profile.Messages[1])
	}
}

func TestReviseGoalIntentTransitionsToProposingGoal(t *testing.T) {
	svc := NewService()
	candidate := "플랫폼에 강의를 무료로 제공할 수 있는 강의안을 완성한다."
	profile := &GoalProfile{
		UserIntent:     "강의 만들기",
		InterviewState: StateRevisingGoal,
		Messages: []InterviewMessage{
			{Role: "user", Content: "목표를 수정하고 싶어요."},
			{Role: "user", Content: "플랫폼에 강의를 무료로 제공하고 싶어요."},
		},
	}
	analysis := &IntentAnalysisResult{
		GoalCandidate:     &candidate,
		IsGoalClearEnough: true,
		UserStance:        "revising",
	}

	intent := svc.DetermineReplyIntent(profile, analysis, nil)
	if intent != ReplyReviseGoal {
		t.Fatalf("intent = %q, want %q", intent, ReplyReviseGoal)
	}

	nextState := profile.InterviewState
	switch intent {
	case ReplyProposeGoal, ReplyReviseGoal:
		nextState = StateProposingGoal
	case ReplyConfirmGoal:
		nextState = StateConfirmed
	case ReplyAskRebuildDecision:
		nextState = StateAwaitingRebuildDecision
	}
	if nextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", nextState, StateProposingGoal)
	}
}

func TestNormalizeRebuildDecisionMapsLegacyRemainingToAll(t *testing.T) {
	if got := NormalizeRebuildDecision(RebuildRemaining); got != RebuildAll {
		t.Fatalf("NormalizeRebuildDecision(rebuild_remaining) = %q, want %q", got, RebuildAll)
	}
	if got := NormalizeRebuildDecision(RebuildKeepStructure); got != RebuildKeepStructure {
		t.Fatalf("NormalizeRebuildDecision(keep_structure) = %q, want %q", got, RebuildKeepStructure)
	}
}

func TestBuildFallbackInterviewResponseForCraftTopicStartsWithOutcomeQuestion(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "목공예",
		Messages: []InterviewMessage{
			{Role: "user", Content: "목공예 배우고 싶어"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "목공예 배우고 싶어")
	if !strings.Contains(resp.Message, "직접 만들어보고 싶은") {
		t.Fatalf("expected more natural craft first question, got %q", resp.Message)
	}
}

func TestBuildFallbackInterviewResponseForCraftHomeGoalAdvances(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "목공예",
		Messages: []InterviewMessage{
			{Role: "user", Content: "목공예 배우고 싶어"},
			{Role: "user", Content: "집에서 간단한 물건은 내손으로 만들고 싶어. 취미생활로"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "집에서 간단한 물건은 내손으로 만들고 싶어. 취미생활로")
	if strings.Contains(resp.Message, "그 실력을 가장 먼저 어디에 써보고 싶으세요?") {
		t.Fatalf("expected conversation to move past usage question, got %q", resp.Message)
	}
}

func TestBuildFallbackInterviewResponseForPotteryTopicStartsWithOutcomeQuestion(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "도자기",
		Messages: []InterviewMessage{
			{Role: "user", Content: "도자기 굽고 싶다"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "도자기 굽고 싶다")
	if !strings.Contains(resp.Message, "직접 만들어보고 싶은") {
		t.Fatalf("expected more natural pottery first question, got %q", resp.Message)
	}
}

func TestBuildFallbackInterviewResponseForPotteryHealingGoalAsksAboutResult(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "도자기",
		Messages: []InterviewMessage{
			{Role: "user", Content: "도자기 굽고 싶다"},
			{Role: "user", Content: "도자기 만드는걸로 힐링 할래"},
		},
	}

	resp := svc.BuildFallbackInterviewResponse(profile, "도자기 만드는걸로 힐링 할래")
	if strings.Contains(resp.Message, "그 실력을 가장 먼저 어디에 써보고 싶으세요?") {
		t.Fatalf("expected a result-oriented follow-up, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "어떤 물건이나 작품") {
		t.Fatalf("expected healing follow-up to ask about desired creation, got %q", resp.Message)
	}
}

func TestBuildGoalProposalProducesNaturalSentence(t *testing.T) {
	svc := NewService()
	motivation := "초등학교에서 연만들기 강의를 하고 싶어"
	usageContext := "초등학교에서 연만들기 강의를 하고 싶어"
	profile := &GoalProfile{
		UserIntent:     "연 만들기 배우고 싶어",
		Motivation:     &motivation,
		UsageContext:   &usageContext,
		TimeHorizon:    nil,
		InterviewState: StateProposingGoal,
	}

	got := svc.BuildGoalProposal(profile)
	if strings.Contains(got, "배우고 싶어를 활용해") || strings.Contains(got, "같은 상황에서") {
		t.Fatalf("proposal still contains awkward duplicated phrasing: %q", got)
	}
}

func TestParseAIResponseParsesTaggedNaturalReply(t *testing.T) {
	svc := NewService()
	raw := `좋아요. 그럼 초등학교 수업에서 아이들과 함께 연 만들기를 진행할 수 있는 방향으로 목표를 잡아볼게요. 이 목표로 탐험계획을 만들어볼까요?
<analysis>{"next_state":"proposing_goal","goal_candidate":"초등학교에서 연 만들기 강의를 할 수 있게 되기","proposed_goal":"3개월 안에 연 만들기를 바탕으로 초등학교에서 연 만들기 강의를 할 수 있게 된다.","missing_info":[],"extracted":{"motivation":"초등학교에서 연만들기 강의를 하고 싶어","usage_context":"초등학교 수업","goal_type":"teaching","output_type":null,"difficulty_level":null,"time_horizon":"3개월"}}</analysis>`

	resp, err := svc.ParseAIResponse(raw)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if !strings.Contains(resp.Message, "이 목표로 탐험계획을 만들어볼까요?") {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateProposingGoal)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal == "" {
		t.Fatal("expected proposed goal from tagged analysis")
	}
}

func TestExtractProposedGoalFromRevisionReply(t *testing.T) {
	svc := NewService()
	message := "무료 제공을 통해 강의를 제공하려고 하시는군요! 이를 바탕으로 새로운 목표를 다음과 같이 제안합니다: '무료 강의를 제공하여 경험을 쌓고, 나중에 유료 강의를 계획합니다.' 이 목표로 탐험계획을 만들어볼까요?"

	got := svc.ExtractProposedGoalFromReply(message)
	want := "무료 강의를 제공하여 경험을 쌓고, 나중에 유료 강의를 계획합니다"
	if got != want {
		t.Fatalf("extracted goal = %q, want %q", got, want)
	}
}

func TestApplyInterviewResponseRefreshesLearningIntentProfile(t *testing.T) {
	svc := NewService()
	handler := NewHandler(nil, svc)
	oldGoal := "이전 목표"
	newGoal := "가죽공예로 손 지갑을 완성한다"
	level := "초보"
	profile := &GoalProfile{
		ConfirmedGoal:   &oldGoal,
		DifficultyLevel: &level,
		InterviewState:  StateClarifying,
		LearningIntent: LearningIntentProfile{
			ConfirmedGoal: oldGoal,
			Source:        LearningIntentSourceBackfill,
			Confidence:    "low",
		},
	}

	handler.applyInterviewResponse(profile, &InterviewAIResponse{
		Message:      "좋아요. 이 목표로 탐험계획을 만들어볼까요?",
		NextState:    StateProposingGoal,
		ProposedGoal: &newGoal,
	})

	if profile.LearningIntent.ConfirmedGoal != newGoal {
		t.Fatalf("learning intent confirmed goal = %q, want %q", profile.LearningIntent.ConfirmedGoal, newGoal)
	}
	if profile.LearningIntent.Source != LearningIntentSourceGoalChat {
		t.Fatalf("learning intent source = %q, want %q", profile.LearningIntent.Source, LearningIntentSourceGoalChat)
	}
	if profile.LearningIntent.LearnerLevel != "beginner" {
		t.Fatalf("learning intent level = %q, want beginner", profile.LearningIntent.LearnerLevel)
	}
}

func TestRefreshLearningIntentProfileRebuildsFromCurrentGoalFields(t *testing.T) {
	svc := NewService()
	confirmed := "3개월 안에 좋아하는 곡을 연주한다"
	motivation := "개인적으로 좋아하는 곡을 연주하기 위해"
	output := "연주"
	level := "입문"
	profile := &GoalProfile{
		ConfirmedGoal:   &confirmed,
		Motivation:      &motivation,
		OutputType:      &output,
		DifficultyLevel: &level,
		LearningIntent: LearningIntentProfile{
			ConfirmedGoal: "오래된 목표",
			Source:        LearningIntentSourceBackfill,
		},
	}

	svc.RefreshLearningIntentProfile(profile)

	if profile.LearningIntent.ConfirmedGoal != confirmed {
		t.Fatalf("confirmed goal = %q, want %q", profile.LearningIntent.ConfirmedGoal, confirmed)
	}
	if profile.LearningIntent.Source != LearningIntentSourceGoalChat {
		t.Fatalf("source = %q, want %q", profile.LearningIntent.Source, LearningIntentSourceGoalChat)
	}
	if profile.LearningIntent.DesiredOutput != output {
		t.Fatalf("desired output = %q, want %q", profile.LearningIntent.DesiredOutput, output)
	}
	if profile.LearningIntent.Confidence != "high" {
		t.Fatalf("confidence = %q, want high", profile.LearningIntent.Confidence)
	}
}

func TestApplyInterviewResponseUsesGoalInReplyForProposalCard(t *testing.T) {
	svc := NewService()
	handler := NewHandler(nil, svc)
	oldGoal := "이전 목표"
	profile := &GoalProfile{
		ConfirmedGoal:  &oldGoal,
		InterviewState: StateRevisingGoal,
		Version:        2,
	}
	message := "이를 바탕으로 새로운 목표를 다음과 같이 제안합니다: '무료 강의를 제공하여 경험을 쌓고, 나중에 유료 강의를 계획합니다.' 이 목표로 탐험계획을 만들어볼까요?"

	handler.applyInterviewResponse(profile, &InterviewAIResponse{
		Message:   message,
		NextState: StateProposingGoal,
	})

	if profile.ConfirmedGoal == nil {
		t.Fatal("confirmed goal is nil")
	}
	want := "무료 강의를 제공하여 경험을 쌓고, 나중에 유료 강의를 계획합니다"
	if *profile.ConfirmedGoal != want {
		t.Fatalf("confirmed goal = %q, want %q", *profile.ConfirmedGoal, want)
	}
}

func TestShouldOverrideAIResponseOnlyForRepeatOrEmpty(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent: "도자기",
		Messages: []InterviewMessage{
			{Role: "user", Content: "도자기 굽고 싶다"},
			{Role: "lumi", Content: "좋아요. 도자기를 통해 어떤 변화를 가장 이루고 싶으세요?"},
		},
	}

	resp := &InterviewAIResponse{
		Message:   "좋아요. 힐링이 되는 취미로 만들고 싶으시군요. 가장 먼저 어떤 작품을 직접 만들어보고 싶으세요?",
		NextState: StateClarifying,
	}
	if svc.ShouldOverrideAIResponse(profile, "도자기 만드는걸로 힐링 할래", resp) {
		t.Fatal("expected natural non-repeated response to be accepted")
	}
}

func TestNormalizeTopicPhraseStripsConversationStyleSuffix(t *testing.T) {
	got := normalizeTopicPhrase("기초 영어 회화를 배우자")
	if got != "기초 영어 회화" {
		t.Fatalf("normalized topic = %q, want %q", got, "기초 영어 회화")
	}
}

func TestBuildResponsePromptUsesLearningLanguage(t *testing.T) {
	svc := NewService()
	prompt := svc.BuildResponsePrompt(ResponseContext{
		State:             StateClarifying,
		Language:          "en",
		LatestUserMessage: "I want to learn urban sketching",
		ReplyIntent:       ReplyNarrowDirection,
	})

	if !strings.Contains(prompt, "Use concise, natural English") {
		t.Fatalf("prompt missing English language instruction: %q", prompt)
	}
	if !strings.Contains(prompt, "Lumi's natural English reply") {
		t.Fatalf("prompt missing English reply schema hint: %q", prompt)
	}
}

type countingErrorLLMClient struct {
	calls int
}

func (c *countingErrorLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	c.calls++
	return "", errors.New("llm unavailable")
}

type staticLLMClient struct {
	calls    int
	response string
}

func (c *staticLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	c.calls++
	return c.response, nil
}

type sequenceLLMClient struct {
	calls     int
	responses []string
}

func (c *sequenceLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	c.calls++
	if len(c.responses) == 0 {
		return "", errors.New("no response configured")
	}
	idx := c.calls - 1
	if idx >= len(c.responses) {
		idx = len(c.responses) - 1
	}
	return c.responses[idx], nil
}

func TestBuildTemplateInterviewReplyProposesGoalInKorean(t *testing.T) {
	svc := NewService()
	goal := "3개월 안에 어반스케치로 여행 장면을 자신 있게 기록할 수 있게 된다."
	resp := svc.BuildTemplateInterviewReply(ResponseContext{
		State:             StateClarifying,
		Language:          "ko",
		LatestUserMessage: "여행에서 그림으로 기록하고 싶어",
		ReplyIntent:       ReplyProposeGoal,
		GoalCandidate:     &goal,
	})

	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateProposingGoal)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != goal {
		t.Fatalf("proposed goal = %#v, want %q", resp.ProposedGoal, goal)
	}
	if !strings.Contains(resp.Message, "이 목표로 탐험계획을 만들어볼까요?") {
		t.Fatalf("message missing Korean proposal confirmation: %q", resp.Message)
	}
}

func TestBuildTemplateInterviewReplyUsesEnglish(t *testing.T) {
	svc := NewService()
	goal := "Use urban sketching to confidently document travel scenes."
	resp := svc.BuildTemplateInterviewReply(ResponseContext{
		State:             StateClarifying,
		Language:          "en",
		LatestUserMessage: "I want to sketch travel scenes",
		ReplyIntent:       ReplyProposeGoal,
		GoalCandidate:     &goal,
	})

	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateProposingGoal)
	}
	if !strings.Contains(resp.Message, "Shall we create your exploration plan with this goal?") {
		t.Fatalf("message missing English proposal confirmation: %q", resp.Message)
	}
}

func TestRunAnalysisTemplateReplyWithClientCallsLLMOnce(t *testing.T) {
	svc := NewService()
	goal := "어반스케치로 여행 장면을 자신 있게 기록할 수 있게 된다."
	client := &staticLLMClient{response: `{
		"extracted": {
			"motivation": "여행 장면을 그림으로 기록하고 싶다",
			"usage_context": "여행",
			"goal_type": "creative_output",
			"output_type": "artwork",
			"difficulty_level": "beginner",
			"time_horizon": "3개월"
		},
		"user_stance": "neutral",
		"goal_candidate": "` + goal + `",
		"is_goal_clear_enough": true,
		"reason_summary": "목표와 활용 맥락이 분명하다"
	}`}
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "lumi", Content: "좋아요. 어디에 써보고 싶으세요?"},
			{Role: "user", Content: "여행에서 그림으로 기록하고 싶어"},
		},
	}

	resp, err := svc.RunAnalysisTemplateReplyWithClient(context.Background(), profile, "여행에서 그림으로 기록하고 싶어", client)
	if err != nil {
		t.Fatalf("RunAnalysisTemplateReplyWithClient() error = %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want 1", client.calls)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateProposingGoal)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != goal {
		t.Fatalf("proposed goal = %#v, want %q", resp.ProposedGoal, goal)
	}
}

func TestCreateGoalInterviewTurnJobInputStoresReferencesWithoutMessageRaw(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	profile := &GoalProfile{
		ID:                profileID,
		UserID:            userID,
		Language:          "ko",
		InterviewState:    StateClarifying,
		Version:           3,
		SummarizedContext: "사용자 원문이 아닌 요약",
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "lumi", Content: "어디에 쓰고 싶으세요?"},
			{Role: "user", Content: "여행에서 그림으로 기록하고 싶어"},
		},
	}

	input := CreateGoalInterviewTurnJobInput(userID, profile, "analysis_template")

	if input.Feature != "goal_interview_turn" {
		t.Fatalf("feature = %q, want goal_interview_turn", input.Feature)
	}
	if got := input.RequestRef["goal_profile_id"]; got != profileID.String() {
		t.Fatalf("goal_profile_id = %v, want %s", got, profileID)
	}
	if got := input.RequestRef["goal_profile_version"]; got != 3 {
		t.Fatalf("goal_profile_version = %v, want 3", got)
	}
	if got := input.RequestRef["message_count"]; got != 3 {
		t.Fatalf("message_count = %v, want 3", got)
	}
	if got := input.RequestRef["mode"]; got != "analysis_template" {
		t.Fatalf("mode = %v, want analysis_template", got)
	}
	if len(input.PromptInputRef) != 0 {
		t.Fatalf("prompt input ref should be empty, got %#v", input.PromptInputRef)
	}
	for key, value := range input.RequestRef {
		if strings.Contains(key, "message") && key != "message_count" {
			t.Fatalf("unexpected raw message key in request ref: %s", key)
		}
		if value == "여행에서 그림으로 기록하고 싶어" || value == "어반스케치 배우고 싶어" {
			t.Fatalf("request ref stores raw user message in %s", key)
		}
	}
}

func TestResolveInterviewResponseUsesAnalysisTemplateMode(t *testing.T) {
	svc := NewService()
	handler := NewHandlerWithOptions(nil, svc, HandlerOptions{
		GoalChatMode: "analysis_template",
	})
	client := &staticLLMClient{response: `{
		"extracted": {
			"motivation": "여행 장면을 그림으로 기록하고 싶다",
			"usage_context": "여행",
			"goal_type": "creative_output",
			"output_type": "artwork",
			"difficulty_level": "beginner",
			"time_horizon": null
		},
		"user_stance": "neutral",
		"goal_candidate": "어반스케치로 여행 장면을 자신 있게 기록할 수 있게 된다.",
		"is_goal_clear_enough": true,
		"reason_summary": "목표와 활용 맥락이 분명하다"
	}`}
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "여행에서 그림으로 기록하고 싶어"},
		},
	}

	resp := handler.resolveInterviewResponseWithExecution(context.Background(), uuid.Nil, profile, "여행에서 그림으로 기록하고 싶어", &goalLLMExecution{
		Client:        client,
		Provider:      "openai",
		Model:         "gpt-4o-mini",
		BillingStatus: "charged",
	})

	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want 1", client.calls)
	}
	if resp == nil || resp.NextState != StateProposingGoal {
		t.Fatalf("expected proposing template response, got %#v", resp)
	}
	if resp.Route != "analysis_template" || resp.LLMCallCount != 1 {
		t.Fatalf("route=%q calls=%d, want analysis_template/1", resp.Route, resp.LLMCallCount)
	}
}

func TestTryBuildGoalInterviewFastPathHandlesInitialShortTopic(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent:     "어반스케치 배우고 싶어",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
		},
	}

	resp, ok := svc.TryBuildGoalInterviewFastPath(profile, "어반스케치 배우고 싶어")
	if !ok {
		t.Fatal("expected fast path")
	}
	if resp.NextState != StateClarifying {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateClarifying)
	}
	if strings.TrimSpace(resp.Message) == "" {
		t.Fatal("expected clarifying message")
	}
}

func TestTryBuildGoalInterviewFastPathRejectsComplexInitialInput(t *testing.T) {
	svc := NewService()
	profile := &GoalProfile{
		UserIntent:     "어반스케치를 배워서 여행에서 그림으로 기록하고 싶어",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치를 배워서 여행에서 그림으로 기록하고 싶어"},
		},
	}

	if _, ok := svc.TryBuildGoalInterviewFastPath(profile, "어반스케치를 배워서 여행에서 그림으로 기록하고 싶어"); ok {
		t.Fatal("expected complex goal-like input to use LLM path")
	}
}

func TestTryBuildGoalInterviewFastPathConfirmsAcceptedGoal(t *testing.T) {
	svc := NewService()
	goal := "어반스케치로 여행 장면을 자신 있게 기록할 수 있게 된다."
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		ConfirmedGoal:  &goal,
		Language:       "ko",
		InterviewState: StateProposingGoal,
		Version:        1,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
			{Role: "lumi", Content: "이 목표로 탐험계획을 만들어볼까요?"},
			{Role: "user", Content: "좋아"},
		},
	}

	resp, ok := svc.TryBuildGoalInterviewFastPath(profile, "좋아")
	if !ok {
		t.Fatal("expected acceptance fast path")
	}
	if resp.NextState != StateConfirmed {
		t.Fatalf("next state = %q, want %q", resp.NextState, StateConfirmed)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != goal {
		t.Fatalf("proposed goal = %#v, want %q", resp.ProposedGoal, goal)
	}
}

func TestTryBuildGoalInterviewFastPathRejectsRevisionLikeAcceptance(t *testing.T) {
	svc := NewService()
	goal := "어반스케치로 여행 장면을 자신 있게 기록할 수 있게 된다."
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		ConfirmedGoal:  &goal,
		Language:       "ko",
		InterviewState: StateProposingGoal,
	}

	if _, ok := svc.TryBuildGoalInterviewFastPath(profile, "아니 다시 바꿀래"); ok {
		t.Fatal("expected revision-like message to use LLM path")
	}
}

func TestResolveInterviewResponseFastPathDoesNotCallLLMWhenEnabled(t *testing.T) {
	svc := NewService()
	handler := NewHandlerWithOptions(nil, svc, HandlerOptions{
		GoalChatFastPathEnabled: true,
	})
	client := &countingErrorLLMClient{}
	profile := &GoalProfile{
		UserIntent:     "어반스케치 배우고 싶어",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
		},
	}

	resp := handler.resolveInterviewResponseWithExecution(context.Background(), uuid.Nil, profile, "어반스케치 배우고 싶어", &goalLLMExecution{
		Client:        client,
		Provider:      "openai",
		Model:         "gpt-4o-mini",
		BillingStatus: "charged",
	})

	if client.calls != 0 {
		t.Fatalf("llm calls = %d, want 0", client.calls)
	}
	if resp == nil || resp.NextState != StateClarifying {
		t.Fatalf("expected clarifying fast path response, got %#v", resp)
	}
	if resp.Route != "fast_path" || resp.LLMCallCount != 0 {
		t.Fatalf("route=%q calls=%d, want fast_path/0", resp.Route, resp.LLMCallCount)
	}
}

func TestResolveInterviewResponseDoesNotRetrySinglePromptFallbackByDefault(t *testing.T) {
	svc := NewService()
	handler := NewHandlerWithOptions(nil, svc, HandlerOptions{
		SinglePromptFallbackLLMEnabled: false,
	})
	client := &countingErrorLLMClient{}
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
		},
	}

	resp := handler.resolveInterviewResponseWithExecution(context.Background(), uuid.Nil, profile, "어반스케치 배우고 싶어", &goalLLMExecution{
		Client:        client,
		Provider:      "openai",
		Model:         "gpt-4o-mini",
		BillingStatus: "charged",
	})

	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want 1", client.calls)
	}
	if resp == nil || strings.TrimSpace(resp.Message) == "" {
		t.Fatalf("expected local fallback response, got %#v", resp)
	}
	if resp.Route != "local_fallback" || resp.LLMCallCount != 1 || resp.FallbackReason == "" {
		t.Fatalf("route=%q calls=%d reason=%q, want local_fallback/1/reason", resp.Route, resp.LLMCallCount, resp.FallbackReason)
	}
}

func TestResolveInterviewResponseCanUseSinglePromptFallbackWhenEnabled(t *testing.T) {
	svc := NewService()
	handler := NewHandlerWithOptions(nil, svc, HandlerOptions{
		SinglePromptFallbackLLMEnabled: true,
	})
	client := &countingErrorLLMClient{}
	profile := &GoalProfile{
		UserIntent:     "어반스케치",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "어반스케치 배우고 싶어"},
		},
	}

	resp := handler.resolveInterviewResponseWithExecution(context.Background(), uuid.Nil, profile, "어반스케치 배우고 싶어", &goalLLMExecution{
		Client:        client,
		Provider:      "openai",
		Model:         "gpt-4o-mini",
		BillingStatus: "charged",
	})

	if client.calls != 2 {
		t.Fatalf("llm calls = %d, want 2", client.calls)
	}
	if resp == nil || strings.TrimSpace(resp.Message) == "" {
		t.Fatalf("expected local fallback response, got %#v", resp)
	}
	if resp.Route != "local_fallback" || resp.LLMCallCount != 2 {
		t.Fatalf("route=%q calls=%d, want local_fallback/2", resp.Route, resp.LLMCallCount)
	}
}

func TestInferExtractedInfoIgnoresShortGoalProceedRequest(t *testing.T) {
	svc := NewService()
	info := svc.InferExtractedInfo("그 목표로")
	if info.Motivation != nil || info.UsageContext != nil {
		t.Fatalf("expected proceed request not to overwrite extracted fields, got %#v", info)
	}
}

func TestBuildGoalProposalForHobbyInstrument(t *testing.T) {
	svc := NewService()
	motivation := "취미로"
	profile := &GoalProfile{
		UserIntent: "바이올린 배우고 싶어",
		Motivation: &motivation,
	}

	got := svc.BuildGoalProposal(profile)
	want := "바이올린을 취미로 꾸준히 연주할 수 있게 된다."
	if got != want {
		t.Fatalf("goal proposal = %q, want %q", got, want)
	}
}

func TestBuildGoalProposalForHobbyWatercolorFromUsageContext(t *testing.T) {
	svc := NewService()
	usageContext := "취미로"
	profile := &GoalProfile{
		UserIntent:   "기초 수채화",
		UsageContext: &usageContext,
	}

	got := svc.BuildGoalProposal(profile)
	want := "기초 수채화를 취미로 꾸준히 즐길 수 있게 된다."
	if got != want {
		t.Fatalf("goal proposal = %q, want %q", got, want)
	}
}

func TestRunAnalysisTemplateReplyForHobbyInstrumentProposesAfterProceedRequest(t *testing.T) {
	svc := NewService()
	client := &staticLLMClient{response: `{
		"extracted": {
			"motivation": null,
			"usage_context": null,
			"goal_type": null,
			"output_type": null,
			"difficulty_level": null,
			"time_horizon": null
		},
		"user_stance": "neutral",
		"goal_candidate": null,
		"is_goal_clear_enough": false,
		"reason_summary": "추가 정보가 필요하다"
	}`}
	profile := &GoalProfile{
		UserIntent:     "바이올린 배우고 싶어",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "바이올린 배우고 싶어"},
			{Role: "lumi", Content: "왜 바이올린을 배우고 싶은지 말씀해주실 수 있을까요?"},
			{Role: "user", Content: "취미로"},
			{Role: "lumi", Content: "구체적으로 어떤 곡이나 스타일을 연주하고 싶으신가요?"},
			{Role: "user", Content: "그 목표로"},
		},
	}

	resp, err := svc.RunAnalysisTemplateReplyWithClient(context.Background(), profile, "그 목표로", client)
	if err != nil {
		t.Fatalf("RunAnalysisTemplateReplyWithClient() error = %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want 1", client.calls)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q; message=%q", resp.NextState, StateProposingGoal, resp.Message)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != "바이올린을 취미로 꾸준히 연주할 수 있게 된다." {
		t.Fatalf("proposed goal = %#v", resp.ProposedGoal)
	}
	if !strings.Contains(resp.Message, "이 목표로 탐험계획을 만들어볼까요?") {
		t.Fatalf("expected proposal confirmation, got %q", resp.Message)
	}
}

func TestRunAnalysisTemplateReplyForHobbyInstrumentProposesAfterHobbyMotivation(t *testing.T) {
	svc := NewService()
	client := &staticLLMClient{response: `{
		"extracted": {
			"motivation": "취미로",
			"usage_context": null,
			"goal_type": null,
			"output_type": null,
			"difficulty_level": null,
			"time_horizon": null
		},
		"user_stance": "neutral",
		"goal_candidate": null,
		"is_goal_clear_enough": false,
		"reason_summary": "활용 맥락이 더 필요하다"
	}`}
	profile := &GoalProfile{
		UserIntent:     "바이올린",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "바이올린"},
			{Role: "lumi", Content: "배우고 나서 가장 얻고 싶은 결과는 무엇인가요?"},
			{Role: "user", Content: "취미로"},
		},
	}

	resp, err := svc.RunAnalysisTemplateReplyWithClient(context.Background(), profile, "취미로", client)
	if err != nil {
		t.Fatalf("RunAnalysisTemplateReplyWithClient() error = %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want 1", client.calls)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q; message=%q", resp.NextState, StateProposingGoal, resp.Message)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != "바이올린을 취미로 꾸준히 연주할 수 있게 된다." {
		t.Fatalf("proposed goal = %#v", resp.ProposedGoal)
	}
}

func TestRunAnalysisTemplateReplyForHobbyWatercolorOverridesWeakGoalCandidate(t *testing.T) {
	svc := NewService()
	client := &staticLLMClient{response: `{
		"extracted": {
			"motivation": null,
			"usage_context": "취미로",
			"goal_type": null,
			"output_type": null,
			"difficulty_level": null,
			"time_horizon": null
		},
		"user_stance": "neutral",
		"goal_candidate": "기초 수채화를 배우기",
		"is_goal_clear_enough": true,
		"reason_summary": "기초 수채화를 취미로 배우려 한다"
	}`}
	profile := &GoalProfile{
		UserIntent:     "기초 수채화",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "기초 수채화"},
			{Role: "lumi", Content: "가장 먼저 그리고 싶은 장면이나 순간이 있나요?"},
			{Role: "user", Content: "취미로"},
			{Role: "lumi", Content: "이 실력을 가장 먼저 어디에 써보고 싶으세요?"},
			{Role: "user", Content: "취미로"},
		},
	}

	resp, err := svc.RunAnalysisTemplateReplyWithClient(context.Background(), profile, "취미로", client)
	if err != nil {
		t.Fatalf("RunAnalysisTemplateReplyWithClient() error = %v", err)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q; message=%q", resp.NextState, StateProposingGoal, resp.Message)
	}
	if resp.ProposedGoal == nil || *resp.ProposedGoal != "기초 수채화를 취미로 꾸준히 즐길 수 있게 된다." {
		t.Fatalf("proposed goal = %#v", resp.ProposedGoal)
	}
}

func TestRunTwoStepLLMForHobbyInstrumentUsesLocalProposalBeforeRepeatingQuestion(t *testing.T) {
	svc := NewService()
	client := &sequenceLLMClient{responses: []string{
		`{
			"extracted": {
				"motivation": null,
				"usage_context": null,
				"goal_type": null,
				"output_type": null,
				"difficulty_level": null,
				"time_horizon": null
			},
			"user_stance": "neutral",
			"goal_candidate": null,
			"is_goal_clear_enough": false,
			"reason_summary": "추가 정보가 필요하다"
		}`,
		`{"reply_intent":"ask_clarifying_question","assistant_message":"어떤 특정한 곡이나 스타일이 있으신가요?"}`,
	}}
	profile := &GoalProfile{
		UserIntent:     "바이올린 배우고 싶어",
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: "바이올린 배우고 싶어"},
			{Role: "lumi", Content: "왜 바이올린을 배우고 싶은지 말씀해주실 수 있을까요?"},
			{Role: "user", Content: "취미로"},
			{Role: "lumi", Content: "구체적으로 어떤 곡이나 스타일을 연주하고 싶으신가요?"},
			{Role: "user", Content: "그 목표로"},
		},
	}

	resp, err := svc.RunTwoStepLLMWithClient(context.Background(), profile, "그 목표로", client)
	if err != nil {
		t.Fatalf("RunTwoStepLLMWithClient() error = %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("llm calls = %d, want only Step A before local proposal", client.calls)
	}
	if resp.NextState != StateProposingGoal {
		t.Fatalf("next state = %q, want %q; message=%q", resp.NextState, StateProposingGoal, resp.Message)
	}
	if strings.Contains(resp.Message, "특정한 곡") || strings.Contains(resp.Message, "스타일") {
		t.Fatalf("expected local proposal instead of repeated LLM question, got %q", resp.Message)
	}
}
