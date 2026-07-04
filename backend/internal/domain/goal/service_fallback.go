package goal

import (
	"fmt"
	"strings"
)

func (s *Service) ApplyExtracted(profile *GoalProfile, info ExtractedInfo) {
	if info.Motivation != nil && *info.Motivation != "" {
		profile.Motivation = info.Motivation
	}
	if info.UsageContext != nil && *info.UsageContext != "" {
		profile.UsageContext = info.UsageContext
	}
	if info.GoalType != nil && *info.GoalType != "" {
		profile.GoalType = info.GoalType
	}
	if info.OutputType != nil && *info.OutputType != "" {
		profile.OutputType = info.OutputType
	}
	if info.DifficultyLevel != nil && *info.DifficultyLevel != "" {
		profile.DifficultyLevel = info.DifficultyLevel
	}
	if info.TimeHorizon != nil && *info.TimeHorizon != "" {
		profile.TimeHorizon = info.TimeHorizon
	}
}

func (s *Service) ApplyExtractedFromUserHistory(profile *GoalProfile) {
	if profile == nil {
		return
	}
	for _, msg := range profile.Messages {
		if msg.Role != "user" {
			continue
		}
		s.ApplyExtracted(profile, s.InferExtractedInfo(msg.Content))
	}
}

func (s *Service) ShouldOverrideAIResponse(profile *GoalProfile, userMessage string, resp *InterviewAIResponse) bool {
	if resp == nil {
		return true
	}
	if strings.TrimSpace(resp.Message) == "" {
		return true
	}
	if last := lastLumiMessage(profile); last != "" && normalizeInterviewText(last) == normalizeInterviewText(resp.Message) {
		return true
	}
	return strings.Contains(resp.Message, genericInterviewFallbackMessage)
}

func (s *Service) BuildFallbackInterviewResponse(profile *GoalProfile, userMessage string) *InterviewAIResponse {
	temp := cloneGoalProfile(profile)
	extracted := s.InferExtractedInfo(userMessage)
	s.ApplyExtracted(temp, extracted)
	// 개별 메시지에서 추출 (전체 history 연결 금지 — raw 문자열이 목표에 삽입되는 문제 방지)
	for _, msg := range profile.Messages {
		if msg.Role == "user" {
			s.ApplyExtracted(temp, s.InferExtractedInfo(msg.Content))
		}
	}

	if s.readyToPropose(temp, userMessage) {
		proposedGoal := s.BuildGoalProposal(temp)
		message := fmt.Sprintf("지금까지 말씀해주신 내용을 바탕으로 이런 목표를 제안드릴게요. \"%s\" 이 목표로 탐험계획을 만들어볼까요?", proposedGoal)
		if s.isSummaryRequest(userMessage) {
			message = fmt.Sprintf("지금까지 내용을 정리하면, %s\n이 목표로 탐험계획을 만들어볼까요?", proposedGoal)
		}
		return &InterviewAIResponse{
			Message:       message,
			NextState:     StateProposingGoal,
			Extracted:     extracted,
			GoalCandidate: &proposedGoal,
			ProposedGoal:  &proposedGoal,
		}
	}

	return &InterviewAIResponse{
		Message:   s.BuildClarifyingReply(temp, userMessage),
		NextState: StateClarifying,
		Extracted: extracted,
	}
}

func (s *Service) InferExtractedInfo(userMessage string) ExtractedInfo {
	trimmed := strings.TrimSpace(userMessage)
	lower := strings.ToLower(trimmed)
	info := ExtractedInfo{}

	if trimmed == "" {
		return info
	}
	if s.isGoalProceedRequest(trimmed) {
		return info
	}

	if (s.isGoalSignal(lower) && !s.isTopicOnlyIntent(lower)) || containsAny(lower, "작가", "강사", "실력", "잘 그리고", "잘 그리", "표현", "완성도", "자신감", "자신 있게", "익숙", "익히", "기록을 남기", "남기고 싶", "자유롭게") {
		info.Motivation = stringPtr(trimmed)
	}
	if containsAny(lower, "힐링", "취미", "취미생활", "취미 생활", "즐기", "휴식", "마음이 편", "편안") {
		info.Motivation = stringPtr(trimmed)
	}
	if s.isUsageSignal(lower) {
		info.UsageContext = stringPtr(trimmed)
	} else if loc := extractLocationContext(lower); loc != "" {
		info.UsageContext = stringPtr(loc)
	}
	if containsAny(lower, "강사", "수업", "클래스", "가르치", "티칭") {
		info.GoalType = stringPtr("teaching")
	} else if containsAny(lower, "작가", "아티스트", "프로", "전시", "포트폴리오") {
		info.GoalType = stringPtr("creative_output")
	}
	if containsAny(lower, "작품", "그림", "드로잉", "스케치") {
		info.OutputType = stringPtr("artwork")
	}
	if containsAny(lower, "초보", "처음", "입문", "기초", "기본") {
		info.DifficultyLevel = stringPtr("beginner")
	}
	if containsAny(lower, "3개월", "세 달") {
		info.TimeHorizon = stringPtr("3개월")
	} else if containsAny(lower, "6개월", "여섯 달", "반년") {
		info.TimeHorizon = stringPtr("6개월")
	} else if containsAny(lower, "1년", "일 년") {
		info.TimeHorizon = stringPtr("1년")
	}

	return info
}

func (s *Service) isGoalSignal(lower string) bool {
	return containsAny(lower,
		"되고 싶", "하고 싶", "고 싶", "할 수 있게", "할수 있게", "할 수 있",
		"되는게 목표", "되는 게 목표", "목표야", "목표는", "목표로",
		"가르치고 싶", "수업하고 싶", "클래스 열고 싶", "강의하고 싶",
	)
}

func (s *Service) isTopicOnlyIntent(lower string) bool {
	return containsAny(lower,
		"배우고 싶어", "배우고 싶다", "배우고 싶어요",
		"해보고 싶어", "해보고 싶다", "해보고 싶어요",
		"굽고 싶어", "굽고 싶다", "굽고 싶어요",
	) && !containsAny(lower,
		"만들고 싶", "되고 싶", "할 수 있게", "기록을 남기", "남기고 싶",
		"가르치고 싶", "수업하고 싶", "강의하고 싶", "클래스 열고 싶",
	)
}

func (s *Service) isUsageSignal(lower string) bool {
	return containsAny(lower,
		"여행", "일상", "기록", "카페", "거리", "현장", "sns", "인스타",
		"전시", "작품", "포트폴리오", "브런치", "수업", "클래스", "강의",
		"사람들한테", "남들에게", "가르치", "업로드", "집에서", "취미생활", "취미 생활",
	)
}

func (s *Service) BuildClarifyingReply(profile *GoalProfile, userMessage string) string {
	if s.isNeedStatusQuestion(userMessage) {
		missing := s.describeMissingFields(profile)
		if missing == "" {
			return "이제 목표를 제안드릴 수 있어요. 지금까지 말씀해주신 방향으로 목표를 정리해볼까요?"
		}
		return fmt.Sprintf("거의 다 잡혔어요. 목표를 더 분명하게 만들려면 %s 정도만 더 알면 좋아요. 가장 먼저 어디에 써보고 싶으신가요?", missing)
	}
	if s.isSummaryRequest(userMessage) {
		summary := s.BuildInterviewSummary(profile)
		if summary != "" {
			return summary + " 이 방향에서 가장 이루고 싶은 장면이나 결과 한 가지를 말해주실래요?"
		}
	}
	if profile.Motivation == nil || s.needsOutcomeClarification(profile) {
		return s.buildFirstOutcomeQuestion(profile)
	}
	if profile.Motivation != nil && profile.UsageContext != nil && profile.TimeHorizon == nil {
		return "좋아요. 그 목표는 언제쯤 이루고 싶으세요?"
	}
	if s.isWellbeingMotivation(derefString(profile.Motivation)) {
		return "좋아요. 힐링이 되는 취미로 만들고 싶으시군요. 가장 먼저 어떤 물건이나 작품을 직접 만들어보고 싶으세요?"
	}
	if profile.GoalType != nil && strings.TrimSpace(*profile.GoalType) == "teaching" && profile.UsageContext == nil {
		return "좋아요. 그럼 가장 먼저 어떤 형태의 수업을 해보고 싶으세요?"
	}
	if profile.UsageContext == nil {
		return "좋아요. 그 실력은 가장 먼저 어디에서 써보고 싶으세요?"
	}
	if profile.TimeHorizon == nil {
		return "좋아요. 그 목표는 어느 정도 기간 안에 이루고 싶으세요?"
	}
	return "좋아요. 그 목표가 이루어졌다고 느끼려면 어떤 결과가 나오면 좋을까요?"
}

func (s *Service) BuildInterviewSummary(profile *GoalProfile) string {
	parts := []string{}
	if profile.UserIntent != "" {
		parts = append(parts, fmt.Sprintf("어반스케치처럼 '%s'를 배우고 싶으시고", profile.UserIntent))
	}
	if profile.Motivation != nil && strings.TrimSpace(*profile.Motivation) != "" {
		parts = append(parts, fmt.Sprintf("궁극적으로는 '%s'는 마음이 있으세요", strings.TrimSpace(*profile.Motivation)))
	}
	if profile.UsageContext != nil && strings.TrimSpace(*profile.UsageContext) != "" {
		parts = append(parts, fmt.Sprintf("이 실력은 '%s' 같은 맥락에 쓰고 싶으신 것 같아요", strings.TrimSpace(*profile.UsageContext)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "지금까지 내용을 정리하면, " + strings.Join(parts, ", ") + "."
}

func (s *Service) BuildGoalProposal(profile *GoalProfile) string {
	topic := normalizeTopicPhrase(strings.TrimSpace(profile.UserIntent))
	if topic == "" {
		topic = "이 학습 주제"
	}

	motivation := strings.TrimSpace(derefString(profile.Motivation))
	usageContext := normalizeUsageContext(strings.TrimSpace(derefString(profile.UsageContext)))
	timeHorizon := strings.TrimSpace(derefString(profile.TimeHorizon))
	prefix := ""
	if timeHorizon != "" {
		prefix = timeHorizon + " 안에 "
	}
	if s.isWellbeingMotivation(motivation) || s.isWellbeingMotivation(usageContext) {
		topicObject := topic + koreanObjectParticle(topic)
		if isInstrumentTopic(topic) {
			return fmt.Sprintf("%s%s 취미로 꾸준히 연주할 수 있게 된다.", prefix, topicObject)
		}
		return fmt.Sprintf("%s%s 취미로 꾸준히 즐길 수 있게 된다.", prefix, topicObject)
	}

	outcome := "자신만의 결과물을 만들 수 있게 된다"
	if motivation != "" {
		outcome = normalizeOutcomePhrase(motivation)
	}

	// outcome에 topic이 포함된 경우 중복 제거
	if strings.Contains(outcome, topic) {
		outcome = strings.TrimSpace(strings.ReplaceAll(outcome, topic+"를 ", ""))
		outcome = strings.TrimSpace(strings.ReplaceAll(outcome, topic+"을 ", ""))
		outcome = strings.TrimSpace(strings.ReplaceAll(outcome, topic+" ", ""))
	}

	// 기간은 사용자가 명시한 경우에만 붙인다
	if usageContext != "" && !strings.Contains(outcome, usageContext) {
		return fmt.Sprintf("%s%s 통해 %s에서 %s.", prefix, topic+koreanObjectParticle(topic), usageContext, outcome)
	}
	return fmt.Sprintf("%s%s 통해 %s.", prefix, topic+koreanObjectParticle(topic), outcome)
}

func (s *Service) readyToPropose(profile *GoalProfile, userMessage string) bool {
	if profile == nil || strings.TrimSpace(profile.UserIntent) == "" {
		return false
	}
	if profile.Motivation != nil && profile.UsageContext != nil {
		return true
	}
	if s.isSummaryRequest(userMessage) {
		return profile.Motivation != nil || profile.UsageContext != nil
	}
	if s.isGoalProceedRequest(userMessage) {
		return profile.Motivation != nil || profile.UsageContext != nil || countUserMessages(profile.Messages) >= 2
	}
	if s.isSimpleHobbyGoalReady(profile) {
		return true
	}
	return countUserMessages(profile.Messages) >= 3 && (profile.Motivation != nil || profile.UsageContext != nil)
}

func (s *Service) isSimpleHobbyGoalReady(profile *GoalProfile) bool {
	if profile == nil {
		return false
	}
	topic := normalizeTopicPhrase(strings.TrimSpace(profile.UserIntent))
	if topic == "" || (!isInstrumentTopic(topic) && !isVisualArtHobbyTopic(topic)) {
		return false
	}
	return s.isWellbeingMotivation(derefString(profile.Motivation)) || s.isWellbeingMotivation(derefString(profile.UsageContext))
}

func (s *Service) collectUserHistory(profile *GoalProfile) string {
	if profile == nil {
		return ""
	}
	parts := make([]string, 0, len(profile.Messages))
	for _, message := range profile.Messages {
		if message.Role != "user" {
			continue
		}
		trimmed := strings.TrimSpace(message.Content)
		if trimmed == "" {
			continue
		}
		parts = append(parts, trimmed)
	}
	return strings.Join(parts, " ")
}

func (s *Service) isSummaryRequest(userMessage string) bool {
	lower := strings.ToLower(strings.TrimSpace(userMessage))
	return containsAny(lower, "정리", "요약", "이제까지", "지금까지", "목표로 생성", "목표로 만들", "이 방향")
}

func (s *Service) isNeedStatusQuestion(userMessage string) bool {
	lower := strings.ToLower(strings.TrimSpace(userMessage))
	return containsAny(lower, "부족", "왜 아직", "아직 목표", "더 필요한", "왜 계속", "뭐가 필요")
}

func (s *Service) describeMissingFields(profile *GoalProfile) string {
	missing := []string{}
	if profile.Motivation == nil {
		missing = append(missing, "왜 이루고 싶은지")
	}
	if profile.UsageContext == nil {
		missing = append(missing, "어디에 써보고 싶은지")
	}
	if profile.TimeHorizon == nil {
		missing = append(missing, "언제까지 이루고 싶은지")
	}
	return strings.Join(missing, ", ")
}

func (s *Service) needsOutcomeClarification(profile *GoalProfile) bool {
	if profile == nil || profile.Motivation == nil {
		return false
	}

	motivation := normalizeInterviewText(derefString(profile.Motivation))
	topic := normalizeInterviewText(normalizeTopicPhrase(profile.UserIntent))
	if motivation == "" || topic == "" {
		return false
	}
	if !strings.Contains(motivation, topic) {
		return false
	}
	if containsAny(motivation, "강사", "수업", "클래스", "작가", "여행", "기록", "포트폴리오", "sns", "힐링", "취미", "집에서", "선물") {
		return false
	}
	return containsAny(motivation, "배우고싶", "해보고싶", "굽고싶", "그리고싶", "만들고싶")
}

func (s *Service) isWellbeingMotivation(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return containsAny(lower, "힐링", "취미", "취미생활", "취미 생활", "즐기", "휴식", "마음이 편", "편안")
}

func (s *Service) buildFirstOutcomeQuestion(profile *GoalProfile) string {
	topic := normalizeTopicPhrase(strings.TrimSpace(profile.UserIntent))
	if topic == "" {
		return "좋아요. 그걸 배우고 나서 가장 해내고 싶은 건 무엇인가요?"
	}

	lower := strings.ToLower(topic)
	switch {
	case containsAny(lower, "스케치", "드로잉", "그림", "수채화", "펜화"):
		return fmt.Sprintf("좋아요. %s로 가장 먼저 그리고 싶은 장면이나 순간이 있나요?", topic)
	case containsAny(lower, "도자기", "도예", "목공", "가죽", "뜨개", "연 만들기"):
		return fmt.Sprintf("좋아요. %s를 배우고 나면 가장 먼저 직접 만들어보고 싶은 게 있나요?", topic)
	case containsAny(lower, "영어", "일본어", "중국어", "회화"):
		return fmt.Sprintf("좋아요. %s를 배워서 가장 먼저 할 수 있게 되고 싶은 건 무엇인가요?", topic)
	default:
		return fmt.Sprintf("좋아요. %s를 통해 가장 해내고 싶은 건 무엇인가요?", topic)
	}
}

// extractLocationContext detects "X에서 [verb]하고 싶어" patterns and returns X.
// e.g., "찬양단에서 연주하고 싶어" → "찬양단"
// ─── 2-step LLM ──────────────────────────────────────────────────────────────

// BuildConversationSummary는 최근 메시지와 extracted 정보를 한 문단 요약으로 만든다.

func (s *Service) BuildConversationSummary(profile *GoalProfile) string {
	if profile == nil {
		return ""
	}
	parts := []string{}
	topic := normalizeTopicPhrase(strings.TrimSpace(profile.UserIntent))
	if topic != "" {
		parts = append(parts, fmt.Sprintf("학습자는 '%s'를 배우고 싶어한다", topic))
	}
	if profile.Motivation != nil && strings.TrimSpace(*profile.Motivation) != "" {
		parts = append(parts, fmt.Sprintf("동기는 '%s'이다", strings.TrimSpace(*profile.Motivation)))
	}
	if profile.UsageContext != nil && strings.TrimSpace(*profile.UsageContext) != "" {
		parts = append(parts, fmt.Sprintf("활용 맥락은 '%s'이다", strings.TrimSpace(*profile.UsageContext)))
	}
	if profile.TimeHorizon != nil && strings.TrimSpace(*profile.TimeHorizon) != "" {
		parts = append(parts, fmt.Sprintf("목표 기간은 %s이다", strings.TrimSpace(*profile.TimeHorizon)))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ". ") + "."
}

// PreventRepeatedQuestion은 최근 대화에서 이미 확보된 필드를 missingFields에서 제거한다.

func (s *Service) PreventRepeatedQuestion(profile *GoalProfile, missingFields []string) []string {
	if profile == nil {
		return missingFields
	}
	recentAssistantMessages := make([]string, 0, 3)
	count := 0
	for i := len(profile.Messages) - 1; i >= 0 && count < 6; i-- {
		m := profile.Messages[i]
		if m.Role == "lumi" {
			recentAssistantMessages = append(recentAssistantMessages, strings.ToLower(m.Content))
			count++
		}
	}
	alreadyAsked := map[string]bool{}
	for _, msg := range recentAssistantMessages {
		if containsAny(msg, "동기", "왜", "이루고 싶") {
			alreadyAsked["motivation"] = true
		}
		if containsAny(msg, "어디에", "맥락", "써보고", "활용") {
			alreadyAsked["usage_context"] = true
		}
		if containsAny(msg, "언제", "기간", "얼마나") {
			alreadyAsked["time_horizon"] = true
		}
	}
	filtered := missingFields[:0]
	for _, f := range missingFields {
		if !alreadyAsked[f] {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// ComputeMissingFields는 아직 수집되지 않은 필드 목록을 반환한다.

func (s *Service) ComputeMissingFields(profile *GoalProfile) []string {
	missing := []string{}
	if profile.Motivation == nil {
		missing = append(missing, "motivation")
	}
	if profile.UsageContext == nil {
		missing = append(missing, "usage_context")
	}
	return missing
}

// DetermineReplyIntent는 서버 주도로 reply_intent를 결정한다.

func (s *Service) DetermineReplyIntent(profile *GoalProfile, analysis *IntentAnalysisResult, missingFields []string) ReplyIntent {
	if profile == nil || analysis == nil {
		return ReplyAskClarifyingQuestion
	}
	if profile.InterviewState == StateRevisingGoal {
		if analysis.GoalCandidate != nil {
			return ReplyReviseGoal
		}
		return ReplyAskClarifyingQuestion
	}
	if analysis.UserStance == "accepting" && profile.InterviewState == StateProposingGoal {
		return ReplyConfirmGoal
	}
	if analysis.IsGoalClearEnough && analysis.GoalCandidate != nil {
		return ReplyProposeGoal
	}
	if s.readyToPropose(profile, "") {
		return ReplyProposeGoal
	}
	if len(profile.Messages) <= 1 {
		return ReplyExplainPurpose
	}
	if len(missingFields) > 0 {
		return ReplyAskClarifyingQuestion
	}
	return ReplyNarrowDirection
}

// BuildAnalysisPrompt는 Step A: 구조화 판단 프롬프트를 생성한다.
