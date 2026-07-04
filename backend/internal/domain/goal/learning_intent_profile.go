package goal

import "strings"

const (
	LearningIntentSourceGoalChat = "goal_chat"
	LearningIntentSourceBackfill = "backfill"
	LearningIntentSourceFallback = "fallback"
)

func BuildLearningIntentProfile(profile *GoalProfile) LearningIntentProfile {
	return BuildLearningIntentProfileWithSource(profile, LearningIntentSourceGoalChat)
}

func BuildLearningIntentProfileWithSource(profile *GoalProfile, source string) LearningIntentProfile {
	if profile == nil {
		return NormalizeLearningIntentProfile(LearningIntentProfile{Source: source})
	}
	confirmedGoal := limitRunes(cleanProfileText(derefString(profile.ConfirmedGoal)), 160)
	level := normalizeLearnerLevel(derefString(profile.DifficultyLevel))
	desiredOutput := deriveDesiredOutput(profile)
	purpose := derivePurpose(profile)
	preferred := derivePreferredActivities(profile)
	success := deriveSuccessCriteria(confirmedGoal, desiredOutput)

	return NormalizeLearningIntentProfile(LearningIntentProfile{
		ConfirmedGoal:       confirmedGoal,
		LearnerLevel:        level,
		Purpose:             purpose,
		DesiredOutput:       desiredOutput,
		PreferredActivities: preferred,
		SuccessCriteria:     success,
		Source:              source,
	})
}

func IsActionableLearningIntentProfile(profile LearningIntentProfile) bool {
	profile = NormalizeLearningIntentProfile(profile)
	return profile.ConfirmedGoal != "" || profile.Purpose != "" || profile.DesiredOutput != "" || len(profile.SuccessCriteria) > 0
}

func NormalizeLearningIntentProfile(profile LearningIntentProfile) LearningIntentProfile {
	profile.ConfirmedGoal = limitRunes(cleanProfileText(profile.ConfirmedGoal), 160)
	profile.LearnerLevel = normalizeLearnerLevel(profile.LearnerLevel)
	profile.Purpose = limitRunes(cleanProfileText(profile.Purpose), 160)
	profile.DesiredOutput = limitRunes(cleanProfileText(profile.DesiredOutput), 80)
	profile.CurrentBlockers = normalizeProfileList(profile.CurrentBlockers, 3, 60)
	profile.PreferredActivities = normalizeProfileList(profile.PreferredActivities, 3, 60)
	profile.SuccessCriteria = normalizeProfileList(profile.SuccessCriteria, 3, 80)
	profile.Source = normalizeLearningIntentSource(profile.Source)
	profile.Confidence = normalizeLearningIntentConfidence(profile.Confidence)
	if isEmptyLearningIntentProfile(profile) {
		return LearningIntentProfile{}
	}
	if profile.Confidence == "" {
		profile.Confidence = inferLearningIntentConfidence(profile)
	}
	if profile.Source == "" {
		profile.Source = LearningIntentSourceGoalChat
	}
	return profile
}

func IsEmptyLearningIntentProfile(profile LearningIntentProfile) bool {
	return isEmptyLearningIntentProfile(profile)
}

func isEmptyLearningIntentProfile(profile LearningIntentProfile) bool {
	return profile.ConfirmedGoal == "" && profile.LearnerLevel == "" && profile.Purpose == "" &&
		profile.DesiredOutput == "" && len(profile.CurrentBlockers) == 0 &&
		len(profile.PreferredActivities) == 0 && len(profile.SuccessCriteria) == 0
}

func derivePurpose(profile *GoalProfile) string {
	parts := []string{}
	for _, value := range []string{
		derefString(profile.Motivation),
		derefString(profile.UsageContext),
		humanizeGoalType(derefString(profile.GoalType)),
	} {
		value = cleanProfileText(value)
		if value != "" && !isNoisyProfileText(value) && !containsProfileText(parts, value) {
			parts = append(parts, value)
		}
	}
	return limitRunes(strings.Join(parts, " / "), 160)
}

func deriveDesiredOutput(profile *GoalProfile) string {
	output := cleanProfileText(derefString(profile.OutputType))
	if output != "" && !isGenericOutputType(output) {
		return limitRunes(output, 80)
	}
	goal := cleanProfileText(derefString(profile.ConfirmedGoal))
	if containsAny(goal, "연주") {
		return "연주"
	}
	if containsAny(goal, "그리", "드로잉", "스케치") {
		return "그림"
	}
	if containsAny(goal, "작성", "글쓰기", "쓰기") {
		return "글"
	}
	for _, marker := range []string{"만들", "완성", "제작", "출시", "공개"} {
		if strings.Contains(goal, marker) {
			return limitRunes(goal, 80)
		}
	}
	return ""
}

func derivePreferredActivities(profile *GoalProfile) []string {
	activities := []string{}
	goalType := strings.ToLower(strings.TrimSpace(derefString(profile.GoalType)))
	outputType := strings.ToLower(strings.TrimSpace(derefString(profile.OutputType)))
	combined := strings.ToLower(strings.Join([]string{
		derefString(profile.OutputType),
		derefString(profile.UsageContext),
		derefString(profile.DifficultyLevel),
		derefString(profile.GoalType),
		derefString(profile.ConfirmedGoal),
	}, " "))
	if strings.Contains(goalType, "artifact") || strings.Contains(goalType, "creation") ||
		(outputType != "" && !isGenericOutputType(outputType) && containsAny(combined, "만들", "제작", "완성", "출시", "공개")) ||
		containsAny(combined, "만들", "제작", "완성품") {
		activities = append(activities, "프로젝트 실습")
	}
	if strings.Contains(goalType, "performance") || strings.Contains(goalType, "habit") ||
		containsAny(combined, "practice", "routine", "연습", "루틴", "습관", "연주", "회화", "말하기") {
		activities = append(activities, "반복 연습")
	}
	if containsAny(combined, "beginner", "basic", "intro", "초보", "입문", "기초", "초급") {
		activities = append(activities, "기초 따라하기")
	}
	return activities
}

func deriveSuccessCriteria(confirmedGoal, desiredOutput string) []string {
	criteria := []string{}
	if desiredOutput != "" && len([]rune(desiredOutput)) <= 20 {
		criteria = append(criteria, desiredOutput+" 완성")
	}
	if confirmedGoal != "" && !containsProfileText(criteria, confirmedGoal) {
		criteria = append(criteria, confirmedGoal)
	}
	return criteria
}

func humanizeGoalType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "artifact_creation":
		return "결과물 만들기"
	case "performance_execution":
		return "실전 수행"
	case "habit_lifestyle":
		return "루틴 형성"
	case "knowledge_foundation", "concept_foundation":
		return "개념 이해"
	case "certification_assessment":
		return "시험 준비"
	case "presentation_publish":
		return "공개/발표"
	case "professional_transition":
		return "전문 전환"
	default:
		return ""
	}
}

func isNoisyProfileText(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return strings.Contains(lower, "bench isolated") || strings.Contains(lower, "smoke test")
}

func normalizeLearnerLevel(value string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	switch normalized {
	case "beginner", "basic", "intro", "초급", "입문", "기초", "초보":
		return "beginner"
	case "intermediate", "중급":
		return "intermediate"
	case "advanced", "고급", "상급":
		return "advanced"
	case "unknown", "미상":
		return "unknown"
	default:
		return ""
	}
}

func normalizeLearningIntentSource(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case LearningIntentSourceGoalChat:
		return LearningIntentSourceGoalChat
	case LearningIntentSourceBackfill:
		return LearningIntentSourceBackfill
	case LearningIntentSourceFallback:
		return LearningIntentSourceFallback
	default:
		return ""
	}
}

func normalizeLearningIntentConfidence(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "low", "medium", "high":
		return strings.TrimSpace(strings.ToLower(value))
	default:
		return ""
	}
}

func inferLearningIntentConfidence(profile LearningIntentProfile) string {
	score := 0
	if profile.ConfirmedGoal != "" {
		score += 2
	}
	if profile.LearnerLevel != "" {
		score++
	}
	if profile.Purpose != "" {
		score++
	}
	if profile.DesiredOutput != "" {
		score++
	}
	if len(profile.PreferredActivities) > 0 || len(profile.SuccessCriteria) > 0 {
		score++
	}
	switch {
	case score >= 5:
		return "high"
	case score >= 3:
		return "medium"
	case score > 0:
		return "low"
	default:
		return ""
	}
}

func normalizeProfileList(values []string, maxItems, maxRunes int) []string {
	if maxItems <= 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = limitRunes(cleanProfileText(value), maxRunes)
		if value == "" || containsProfileText(out, value) {
			continue
		}
		out = append(out, value)
		if len(out) >= maxItems {
			break
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func cleanProfileText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func limitRunes(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return strings.TrimSpace(string(runes[:limit]))
}

func containsProfileText(values []string, target string) bool {
	target = strings.TrimSpace(strings.ToLower(target))
	if target == "" {
		return false
	}
	for _, value := range values {
		if strings.TrimSpace(strings.ToLower(value)) == target {
			return true
		}
	}
	return false
}

func isGenericOutputType(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "artifact", "project", "performance", "habit", "routine", "knowledge", "output", "result", "작품", "프로젝트", "결과물":
		return true
	default:
		return false
	}
}
