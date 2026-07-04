package goal

import (
	"fmt"
	"strings"
)

func (s *Service) BuildTemplateInterviewReply(rctx ResponseContext) *InterviewAIResponse {
	goalCandidate := strings.TrimSpace(derefString(rctx.GoalCandidate))
	if goalCandidate == "" && (rctx.ReplyIntent == ReplyProposeGoal || rctx.ReplyIntent == ReplyReviseGoal) {
		temp := &GoalProfile{
			UserIntent:      normalizeTopicPhrase(rctx.ConversationSummary),
			Motivation:      rctx.Extracted.Motivation,
			UsageContext:    rctx.Extracted.UsageContext,
			GoalType:        rctx.Extracted.GoalType,
			OutputType:      rctx.Extracted.OutputType,
			DifficultyLevel: rctx.Extracted.DifficultyLevel,
			TimeHorizon:     rctx.Extracted.TimeHorizon,
			Language:        rctx.Language,
			InterviewState:  rctx.State,
		}
		if temp.UserIntent == "" {
			temp.UserIntent = strings.TrimSpace(rctx.LatestUserMessage)
		}
		goalCandidate = strings.TrimSpace(s.BuildGoalProposal(temp))
	}

	nextState := rctx.State
	message := ""
	language := normalizeLearningLanguage(rctx.Language)
	switch rctx.ReplyIntent {
	case ReplyExplainPurpose:
		nextState = StateClarifying
		if language == "en" {
			message = "Good. I will help you turn that learning direction into one clear goal. What result do you most want from learning it?"
		} else {
			message = "좋아요. 그 학습 방향을 하나의 분명한 목표로 정리해볼게요. 배우고 나서 가장 얻고 싶은 결과는 무엇인가요?"
		}
	case ReplyProposeGoal:
		if goalCandidate == "" {
			return s.BuildFallbackInterviewResponse(&GoalProfile{Language: rctx.Language}, rctx.LatestUserMessage)
		}
		nextState = StateProposingGoal
		if language == "en" {
			message = fmt.Sprintf("Based on what you shared, we can set your goal like this:\n\n\"%s\"\n\nShall we create your exploration plan with this goal?", goalCandidate)
		} else {
			message = fmt.Sprintf("말씀해주신 방향을 바탕으로 목표를 이렇게 잡아볼 수 있어요.\n\n\"%s\"\n\n이 목표로 탐험계획을 만들어볼까요?", goalCandidate)
		}
	case ReplyConfirmGoal:
		nextState = StateConfirmed
		if language == "en" {
			message = "Your goal is set. Shall we create your exploration plan now?"
		} else {
			message = "목표가 확정됐어요. 이제 이 목표에 맞춰 탐험계획을 만들어볼까요?"
		}
	case ReplyReviseGoal:
		if goalCandidate == "" {
			return s.BuildFallbackInterviewResponse(&GoalProfile{Language: rctx.Language}, rctx.LatestUserMessage)
		}
		nextState = StateProposingGoal
		if language == "en" {
			message = fmt.Sprintf("Good. I revised the direction into this goal:\n\n\"%s\"\n\nShall we use this goal for the exploration plan?", goalCandidate)
		} else {
			message = fmt.Sprintf("좋아요. 바뀐 방향을 반영해서 목표를 이렇게 다시 잡아볼 수 있어요.\n\n\"%s\"\n\n이 목표로 탐험계획을 만들어볼까요?", goalCandidate)
		}
	case ReplyAskRebuildDecision:
		nextState = StateAwaitingRebuildDecision
		if language == "en" {
			message = "The goal has changed. Would you like to keep your edited structure and adjust from here, or rebuild the whole exploration plan?"
		} else {
			message = "목표가 바뀌었어요. 지금까지 편집한 구조를 유지하며 이어갈까요, 아니면 탐험계획을 전체 재구성할까요?"
		}
	case ReplyAcknowledgeContinue:
		if language == "en" {
			message = "Good. I will keep that in mind. What part should we make clearer next?"
		} else {
			message = "좋아요. 그 부분을 반영해둘게요. 다음으로 어떤 점을 더 분명히 해볼까요?"
		}
	case ReplyNarrowDirection, ReplyAskClarifyingQuestion:
		fallthrough
	default:
		nextState = StateClarifying
		message = s.buildTemplateClarifyingQuestion(rctx)
	}

	if strings.TrimSpace(message) == "" {
		return s.BuildFallbackInterviewResponse(&GoalProfile{Language: rctx.Language}, rctx.LatestUserMessage)
	}
	var proposedGoal *string
	var goalPtr *string
	if goalCandidate != "" && (nextState == StateProposingGoal || rctx.ReplyIntent == ReplyConfirmGoal) {
		goalPtr = stringPtr(goalCandidate)
		proposedGoal = goalPtr
	}
	return &InterviewAIResponse{
		Message:       message,
		NextState:     nextState,
		Extracted:     rctx.Extracted,
		GoalCandidate: goalPtr,
		ProposedGoal:  proposedGoal,
		MissingInfo:   rctx.MissingFields,
	}
}

func (s *Service) TryBuildGoalInterviewFastPath(profile *GoalProfile, userMessage string) (*InterviewAIResponse, bool) {
	if profile == nil {
		return nil, false
	}
	trimmed := strings.TrimSpace(userMessage)
	if trimmed == "" {
		return nil, false
	}
	if profile.InterviewState == StateProposingGoal && s.isClearGoalAcceptance(trimmed) {
		confirmedGoal := strings.TrimSpace(derefString(profile.ConfirmedGoal))
		if confirmedGoal == "" {
			return nil, false
		}
		nextState := StateConfirmed
		message := goalConfirmedMessage(profile.Language)
		if profile.Version > 1 {
			nextState = StateAwaitingRebuildDecision
			if normalizeLearningLanguage(profile.Language) == "en" {
				message = "Your revised goal is set. Would you like to keep your edited structure and adjust from here, or rebuild the whole exploration plan?"
			} else {
				message = "수정한 목표가 확정됐어요. 지금까지 편집한 구조를 유지하며 이어갈까요, 아니면 탐험계획을 전체 재구성할까요?"
			}
		}
		return &InterviewAIResponse{
			Message:       message,
			NextState:     nextState,
			ProposedGoal:  &confirmedGoal,
			GoalCandidate: &confirmedGoal,
		}, true
	}
	if profile.InterviewState != StateClarifying {
		return nil, false
	}
	if countUserMessages(profile.Messages) > 1 {
		return nil, false
	}
	if !s.isShortTopicOnlyInput(trimmed) {
		return nil, false
	}
	extracted := s.InferExtractedInfo(trimmed)
	temp := cloneGoalProfile(profile)
	s.ApplyExtracted(temp, extracted)
	message := s.BuildClarifyingReply(temp, trimmed)
	if strings.TrimSpace(message) == "" {
		return nil, false
	}
	return &InterviewAIResponse{
		Message:   message,
		NextState: StateClarifying,
		Extracted: extracted,
	}, true
}

func (s *Service) isShortTopicOnlyInput(message string) bool {
	normalized := strings.ToLower(strings.TrimSpace(message))
	if normalized == "" {
		return false
	}
	if len([]rune(normalized)) > 60 {
		return false
	}
	if containsAny(normalized, "?", "？", "\n") {
		return false
	}
	if containsAny(normalized,
		"가르치", "강의", "수업", "클래스", "작가", "전시", "포트폴리오", "취업", "이직",
		"여행에서", "회사에서", "학교에서", "아이에게", "사람들한테",
		"because", "so that", "for my job", "portfolio", "teach", "class",
	) {
		return false
	}
	return containsAny(normalized,
		"배우고 싶", "해보고 싶", "입문", "기초", "처음", "연습하고 싶",
		"learn", "want to learn", "try", "beginner", "start",
	)
}

func (s *Service) isClearGoalAcceptance(message string) bool {
	normalized := normalizeInterviewText(message)
	if normalized == "" {
		return false
	}
	if containsAny(normalized,
		"아니", "싫", "다시", "바꿀", "수정", "모르겠", "잘모르", "다른",
		"no", "not", "change", "revise", "again", "different", "notsure",
	) {
		return false
	}
	switch normalized {
	case "응", "네", "예", "좋아", "좋아요", "좋습니다", "그래", "그렇게해", "그걸로해", "그걸로", "진행해", "진행", "확정", "맞아", "맞아요",
		"yes", "yeah", "yep", "ok", "okay", "sure", "soundsgood", "letsdoit", "confirm":
		return true
	}
	return containsAny(normalized, "그목표로", "그걸로진행", "그렇게진행", "이대로진행", "이목표로", "좋습니다진행")
}

func (s *Service) isGoalProceedRequest(message string) bool {
	normalized := normalizeInterviewText(message)
	if normalized == "" {
		return false
	}
	if len([]rune(normalized)) > 24 {
		return false
	}
	switch normalized {
	case "그목표로", "이목표로", "그걸로", "그걸로해", "그렇게해", "그렇게진행", "이대로진행", "그목표로해줘", "이목표로해줘":
		return true
	}
	return containsAny(normalized, "그목표로진행", "이목표로진행", "그걸로진행")
}

func (s *Service) buildTemplateClarifyingQuestion(rctx ResponseContext) string {
	language := normalizeLearningLanguage(rctx.Language)
	missing := ""
	if len(rctx.MissingFields) > 0 {
		missing = rctx.MissingFields[0]
	}
	topic := normalizeTopicPhrase(strings.TrimSpace(rctx.LatestUserMessage))
	if topic == "" && strings.TrimSpace(rctx.ConversationSummary) != "" {
		topic = "that topic"
		if language != "en" {
			topic = "그 주제"
		}
	}
	if language == "en" {
		switch missing {
		case "usage_context":
			return "Good. Where would you most like to use this skill first?"
		case "time_horizon":
			return "Good. When would you like to reach this goal?"
		default:
			if topic != "" {
				return fmt.Sprintf("Good. What result do you most want from learning %s?", topic)
			}
			return "Good. What result do you most want from learning this?"
		}
	}
	switch missing {
	case "usage_context":
		return "좋아요. 이 실력을 가장 먼저 어디에 써보고 싶으세요?"
	case "time_horizon":
		return "좋아요. 그 목표는 언제쯤 이루고 싶으세요?"
	default:
		if topic != "" {
			return fmt.Sprintf("좋아요. %s를 통해 가장 얻고 싶은 결과는 무엇인가요?", topic)
		}
		return "좋아요. 배우고 나서 가장 얻고 싶은 결과는 무엇인가요?"
	}
}

// buildRecentMessageStrings는 최근 N개 메시지를 문자열 슬라이스로 반환한다.
// RunTwoStepLLM은 Step A(분석) → Step B(응답 생성)를 순서대로 실행한다.
// 실패 시 기존 fallback 로직으로 위임한다.
