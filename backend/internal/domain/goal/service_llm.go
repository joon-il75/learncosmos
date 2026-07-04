package goal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/llm"
)

func (s *Service) ParseAIResponse(raw string) (*InterviewAIResponse, error) {
	raw = strings.TrimSpace(raw)
	if tagged, ok := s.parseTaggedInterviewResponse(raw); ok {
		return tagged, nil
	}
	// JSON 블록 추출
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("invalid AI response format")
	}
	raw = raw[start : end+1]

	var resp InterviewAIResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}
	if resp.Message == "" {
		return nil, fmt.Errorf("AI response missing message")
	}
	if resp.NextState == "" {
		resp.NextState = StateClarifying
	}
	return &resp, nil
}

func (s *Service) CallLLM(ctx context.Context, provider, apiKey, model, prompt string) (string, error) {
	client, err := llm.NewClient(provider, apiKey, model)
	if err != nil {
		return "", fmt.Errorf("create llm client: %w", err)
	}
	return client.Complete(ctx, prompt)
}

func (s *Service) CallLLMWithClient(ctx context.Context, client llm.Client, prompt string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("llm client is nil")
	}
	return client.Complete(ctx, prompt)
}

func (s *Service) ParseAnalysisResult(raw string) (*IntentAnalysisResult, error) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("no JSON in analysis response")
	}
	var result IntentAnalysisResult
	if err := json.Unmarshal([]byte(raw[start:end+1]), &result); err != nil {
		return nil, fmt.Errorf("parse analysis result: %w", err)
	}
	return &result, nil
}

// BuildResponsePrompt는 Step B: 자연어 응답 생성 프롬프트를 생성한다.

func (s *Service) ParseReplyPayload(raw string) (string, ReplyIntent, error) {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end < 0 || end <= start {
		return "", "", fmt.Errorf("no JSON in reply response")
	}
	var payload AssistantReplyPayload
	if err := json.Unmarshal([]byte(raw[start:end+1]), &payload); err != nil {
		return "", "", fmt.Errorf("parse reply payload: %w", err)
	}
	if strings.TrimSpace(payload.AssistantMessage) == "" {
		return "", "", fmt.Errorf("empty assistant_message")
	}
	return payload.AssistantMessage, payload.ReplyIntent, nil
}

func (s *Service) RunTwoStepLLM(ctx context.Context, profile *GoalProfile, userMessage, provider, apiKey, model string) (*InterviewAIResponse, error) {
	client, err := llm.NewClient(provider, apiKey, model)
	if err != nil {
		return nil, fmt.Errorf("create llm client: %w", err)
	}
	return s.RunTwoStepLLMWithClient(ctx, profile, userMessage, client)
}

func (s *Service) RunTwoStepLLMWithClient(ctx context.Context, profile *GoalProfile, userMessage string, client llm.Client) (*InterviewAIResponse, error) {
	// Step A: 구조화 판단
	analysisPrompt := s.BuildAnalysisPrompt(profile, userMessage)
	rawAnalysis, err := s.CallLLMWithClient(ctx, client, analysisPrompt)
	if err != nil {
		return nil, fmt.Errorf("step A LLM: %w", err)
	}
	analysis, err := s.ParseAnalysisResult(rawAnalysis)
	if err != nil {
		return nil, fmt.Errorf("step A parse: %w", err)
	}

	// 서버 주도 상태 결정
	s.ApplyExtracted(profile, analysis.Extracted)
	s.ApplyExtractedFromUserHistory(profile)
	missing := s.ComputeMissingFields(profile)
	missing = s.PreventRepeatedQuestion(profile, missing)
	replyIntent := s.DetermineReplyIntent(profile, analysis, missing)
	if replyIntent == ReplyProposeGoal && (analysis.GoalCandidate == nil || s.isSimpleHobbyGoalReady(profile)) {
		goalCandidate := s.BuildGoalProposal(profile)
		analysis.GoalCandidate = &goalCandidate
	}

	// 요약 컨텍스트 갱신
	summary := s.BuildConversationSummary(profile)
	if replyIntent == ReplyProposeGoal && s.readyToPropose(profile, userMessage) {
		profile.SummarizedContext = summary
		return s.BuildTemplateInterviewReply(ResponseContext{
			State:               profile.InterviewState,
			Language:            profile.Language,
			LatestUserMessage:   userMessage,
			ConversationSummary: summary,
			Extracted:           analysis.Extracted,
			MissingFields:       missing,
			GoalCandidate:       analysis.GoalCandidate,
			ConfirmedGoal:       profile.ConfirmedGoal,
			ReplyIntent:         replyIntent,
			RecentMessages:      buildRecentMessageStrings(profile.Messages, 6),
		}), nil
	}

	// Step B: 자연어 응답 생성
	rctx := ResponseContext{
		State:               profile.InterviewState,
		Language:            profile.Language,
		LatestUserMessage:   userMessage,
		ConversationSummary: summary,
		Extracted:           analysis.Extracted,
		MissingFields:       missing,
		GoalCandidate:       analysis.GoalCandidate,
		ConfirmedGoal:       profile.ConfirmedGoal,
		ReplyIntent:         replyIntent,
		RecentMessages:      buildRecentMessageStrings(profile.Messages, 6),
	}
	responsePrompt := s.BuildResponsePrompt(rctx)
	rawReply, err := s.CallLLMWithClient(ctx, client, responsePrompt)
	if err != nil {
		return nil, fmt.Errorf("step B LLM: %w", err)
	}
	message, _, err := s.ParseReplyPayload(rawReply)
	if err != nil {
		return nil, fmt.Errorf("step B parse: %w", err)
	}

	// 상태 전이 결정
	nextState := profile.InterviewState
	switch replyIntent {
	case ReplyProposeGoal:
		nextState = StateProposingGoal
	case ReplyConfirmGoal:
		nextState = StateConfirmed
	case ReplyReviseGoal:
		nextState = StateProposingGoal
	case ReplyAskRebuildDecision:
		nextState = StateAwaitingRebuildDecision
	}

	// summarized_context 갱신
	profile.SummarizedContext = summary

	return &InterviewAIResponse{
		Message:       message,
		NextState:     nextState,
		Extracted:     analysis.Extracted,
		GoalCandidate: analysis.GoalCandidate,
		ProposedGoal:  analysis.GoalCandidate,
	}, nil
}

func (s *Service) RunAnalysisTemplateReplyWithClient(ctx context.Context, profile *GoalProfile, userMessage string, client llm.Client) (*InterviewAIResponse, error) {
	analysisPrompt := s.BuildAnalysisPrompt(profile, userMessage)
	rawAnalysis, err := s.CallLLMWithClient(ctx, client, analysisPrompt)
	if err != nil {
		return nil, fmt.Errorf("step A LLM: %w", err)
	}
	analysis, err := s.ParseAnalysisResult(rawAnalysis)
	if err != nil {
		return nil, fmt.Errorf("step A parse: %w", err)
	}

	s.ApplyExtracted(profile, analysis.Extracted)
	s.ApplyExtractedFromUserHistory(profile)
	missing := s.ComputeMissingFields(profile)
	missing = s.PreventRepeatedQuestion(profile, missing)
	replyIntent := s.DetermineReplyIntent(profile, analysis, missing)
	if replyIntent == ReplyProposeGoal && (analysis.GoalCandidate == nil || s.isSimpleHobbyGoalReady(profile)) {
		goalCandidate := s.BuildGoalProposal(profile)
		analysis.GoalCandidate = &goalCandidate
	}
	summary := s.BuildConversationSummary(profile)
	profile.SummarizedContext = summary

	return s.BuildTemplateInterviewReply(ResponseContext{
		State:               profile.InterviewState,
		Language:            profile.Language,
		LatestUserMessage:   userMessage,
		ConversationSummary: summary,
		Extracted:           analysis.Extracted,
		MissingFields:       missing,
		GoalCandidate:       analysis.GoalCandidate,
		ConfirmedGoal:       profile.ConfirmedGoal,
		ReplyIntent:         replyIntent,
		RecentMessages:      buildRecentMessageStrings(profile.Messages, 6),
	}), nil
}
