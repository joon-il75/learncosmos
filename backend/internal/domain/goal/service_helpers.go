package goal

import (
	"encoding/json"
	"fmt"
	"strings"
)

func cloneGoalProfile(profile *GoalProfile) *GoalProfile {
	if profile == nil {
		return &GoalProfile{}
	}
	cloned := *profile
	cloned.Messages = append([]InterviewMessage(nil), profile.Messages...)
	return &cloned
}

func countUserMessages(messages []InterviewMessage) int {
	count := 0
	for _, message := range messages {
		if message.Role == "user" && strings.TrimSpace(message.Content) != "" {
			count++
		}
	}
	return count
}

func containsAny(target string, patterns ...string) bool {
	for _, pattern := range patterns {
		if strings.Contains(target, pattern) {
			return true
		}
	}
	return false
}

func lastLumiMessage(profile *GoalProfile) string {
	if profile == nil {
		return ""
	}
	for idx := len(profile.Messages) - 1; idx >= 0; idx-- {
		if profile.Messages[idx].Role == "lumi" {
			return strings.TrimSpace(profile.Messages[idx].Content)
		}
	}
	return ""
}

func normalizeInterviewText(value string) string {
	replacer := strings.NewReplacer(" ", "", "\n", "", "\t", "", ".", "", "!", "", "?", "", ",", "", "\"", "", "'", "")
	return replacer.Replace(strings.TrimSpace(strings.ToLower(value)))
}

func normalizeTopicPhrase(topic string) string {
	normalized := strings.TrimSpace(topic)
	for _, suffix := range []string{
		" 배우고 싶어요", " 배우고 싶어", " 배우고 싶다",
		" 배우기", " 해보고 싶어요", " 해보고 싶어", " 해보고 싶다",
		" 배우자", " 해보자", " 하자",
	} {
		if strings.HasSuffix(normalized, suffix) {
			normalized = strings.TrimSpace(strings.TrimSuffix(normalized, suffix))
			break
		}
	}
	for _, particle := range []string{"를", "을"} {
		if strings.HasSuffix(normalized, particle) {
			normalized = strings.TrimSpace(strings.TrimSuffix(normalized, particle))
			break
		}
	}
	return normalized
}

func normalizeOutcomePhrase(outcome string) string {
	normalized := strings.TrimSpace(outcome)
	normalized = strings.TrimSuffix(normalized, ".")
	normalized = strings.TrimSuffix(normalized, "!")
	normalized = strings.TrimSuffix(normalized, "?")

	switch {
	case strings.HasSuffix(normalized, "하고 싶어요"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "하고 싶어요")) + "할 수 있게 된다"
	case strings.HasSuffix(normalized, "하고 싶어"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "하고 싶어")) + "할 수 있게 된다"
	case strings.HasSuffix(normalized, "하고 싶다"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "하고 싶다")) + "할 수 있게 된다"
	case strings.HasSuffix(normalized, "되고 싶어요"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "되고 싶어요")) + "될 수 있게 된다"
	case strings.HasSuffix(normalized, "되고 싶어"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "되고 싶어")) + "될 수 있게 된다"
	case strings.HasSuffix(normalized, "되고 싶다"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "되고 싶다")) + "될 수 있게 된다"
	case strings.HasSuffix(normalized, "고 싶어요"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "고 싶어요")) + "고 싶다는 목표를 이룰 수 있게 된다"
	case strings.HasSuffix(normalized, "고 싶어"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "고 싶어")) + "고 싶다는 목표를 이룰 수 있게 된다"
	case strings.HasSuffix(normalized, "고 싶다"):
		return strings.TrimSpace(strings.TrimSuffix(normalized, "고 싶다")) + "고 싶다는 목표를 이룰 수 있게 된다"
	default:
		return normalized
	}
}

func normalizeUsageContext(context string) string {
	normalized := strings.TrimSpace(context)
	normalized = strings.TrimSuffix(normalized, ".")
	normalized = strings.TrimSuffix(normalized, "!")
	normalized = strings.TrimSuffix(normalized, "?")
	normalized = strings.TrimPrefix(normalized, "연 만들기 배우고 싶어 ")
	normalized = strings.TrimPrefix(normalized, "목공예 배우고 싶어 ")
	return strings.TrimSpace(normalized)
}

func koreanObjectParticle(word string) string {
	trimmed := strings.TrimSpace(word)
	if trimmed == "" {
		return "를"
	}
	runes := []rune(trimmed)
	last := runes[len(runes)-1]
	if last < 0xAC00 || last > 0xD7A3 {
		return "를"
	}
	if (last-0xAC00)%28 == 0 {
		return "를"
	}
	return "을"
}

func isInstrumentTopic(topic string) bool {
	lower := strings.ToLower(strings.TrimSpace(topic))
	return containsAny(lower,
		"바이올린", "비올라", "첼로", "콘트라베이스", "더블베이스",
		"피아노", "기타", "일렉기타", "베이스", "드럼",
		"플룻", "플루트", "색소폰", "리코더", "하모니카",
		"거문고", "가야금", "해금", "대금", "소금", "장구",
		"violin", "cello", "piano", "guitar", "drum", "saxophone",
	)
}

func isVisualArtHobbyTopic(topic string) bool {
	lower := strings.ToLower(strings.TrimSpace(topic))
	return containsAny(lower,
		"수채화", "스케치", "어반스케치", "드로잉", "그림", "펜화", "캘리그라피",
		"watercolor", "watercolour", "sketch", "drawing", "calligraphy",
	)
}

func (s *Service) parseTaggedInterviewResponse(raw string) (*InterviewAIResponse, bool) {
	analysisJSON, ok := extractTaggedBlock(raw, "analysis")
	if !ok {
		return nil, false
	}

	replyText, hasReplyTag := extractTaggedBlock(raw, "reply")
	if !hasReplyTag {
		replyText = strings.TrimSpace(strings.Replace(raw, "<analysis>"+analysisJSON+"</analysis>", "", 1))
	}

	var resp InterviewAIResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(analysisJSON)), &resp); err != nil {
		return nil, false
	}
	if strings.TrimSpace(replyText) != "" {
		resp.Message = strings.TrimSpace(replyText)
	}
	if strings.TrimSpace(resp.Message) == "" {
		return nil, false
	}
	if resp.NextState == "" {
		resp.NextState = StateClarifying
	}
	return &resp, true
}

func extractTaggedBlock(raw string, tag string) (string, bool) {
	openTag := "<" + tag + ">"
	closeTag := "</" + tag + ">"
	start := strings.Index(raw, openTag)
	end := strings.Index(raw, closeTag)
	if start < 0 || end < 0 || end <= start {
		return "", false
	}
	start += len(openTag)
	return strings.TrimSpace(raw[start:end]), true
}

func extractLocationContext(lower string) string {
	idx := strings.Index(lower, "에서 ")
	if idx <= 0 {
		return ""
	}
	suffix := lower[idx+len("에서 "):]
	if !containsAny(suffix, "하고 싶", "연주하고", "쓰고", "활동하고", "사용하고", "써보고", "보여주고") {
		return ""
	}
	// extract the word/phrase before 에서
	prefix := lower[:idx]
	words := strings.Fields(prefix)
	if len(words) == 0 {
		return ""
	}
	return words[len(words)-1]
}

func stringPtr(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	value := strings.TrimSpace(v)
	return &value
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (s *Service) ExtractProposedGoalFromReply(message string) string {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return ""
	}

	for _, pair := range [][2]string{
		{"'", "'"},
		{"\"", "\""},
		{"“", "”"},
		{"‘", "’"},
	} {
		if extracted := extractGoalBetweenDelimiters(trimmed, pair[0], pair[1]); extracted != "" {
			return extracted
		}
	}

	for _, marker := range []string{
		"목표를 다음과 같이 제안합니다:",
		"목표를 이렇게 제안합니다:",
		"목표를 제안드릴게요.",
		"목표를 제안드릴게요:",
		"목표를 잡아볼 수 있겠어요.",
	} {
		if extracted := extractGoalAfterMarker(trimmed, marker); extracted != "" {
			return extracted
		}
	}
	return ""
}

func extractGoalBetweenDelimiters(message, open, close string) string {
	start := strings.Index(message, open)
	if start < 0 {
		return ""
	}
	rest := message[start+len(open):]
	end := strings.Index(rest, close)
	if end < 0 {
		return ""
	}
	return cleanExtractedGoal(rest[:end])
}

func extractGoalAfterMarker(message, marker string) string {
	idx := strings.Index(message, marker)
	if idx < 0 {
		return ""
	}
	candidate := strings.TrimSpace(message[idx+len(marker):])
	if cut := strings.Index(candidate, "이 목표로"); cut >= 0 {
		candidate = candidate[:cut]
	}
	return cleanExtractedGoal(candidate)
}

func cleanExtractedGoal(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, " \t\n\r'\"“”‘’")
	value = strings.TrimSuffix(value, ".")
	return strings.TrimSpace(value)
}

func buildRecentMessageStrings(messages []InterviewMessage, n int) []string {
	start := len(messages) - n
	if start < 0 {
		start = 0
	}
	result := make([]string, 0, n)
	for _, m := range messages[start:] {
		role := "루미"
		if m.Role == "user" {
			role = "학습자"
		}
		result = append(result, fmt.Sprintf("%s: %s", role, m.Content))
	}
	return result
}
