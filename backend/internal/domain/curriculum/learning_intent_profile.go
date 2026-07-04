package curriculum

import (
	"strings"

	goaldomain "github.com/learnweaver/backend/internal/domain/goal"
)

func learningIntentProfileFromGoal(profile goaldomain.LearningIntentProfile) LearningIntentProfile {
	return NormalizeLearningIntentProfile(LearningIntentProfile{
		ConfirmedGoal:       profile.ConfirmedGoal,
		LearnerLevel:        profile.LearnerLevel,
		Purpose:             profile.Purpose,
		DesiredOutput:       profile.DesiredOutput,
		CurrentBlockers:     append([]string(nil), profile.CurrentBlockers...),
		PreferredActivities: append([]string(nil), profile.PreferredActivities...),
		SuccessCriteria:     append([]string(nil), profile.SuccessCriteria...),
		Source:              profile.Source,
		Confidence:          profile.Confidence,
	})
}

func NormalizeLearningIntentProfile(profile LearningIntentProfile) LearningIntentProfile {
	profile.ConfirmedGoal = limitPromptRunes(cleanPromptText(profile.ConfirmedGoal), 120)
	profile.LearnerLevel = normalizeIntentLearnerLevel(profile.LearnerLevel)
	profile.Purpose = limitPromptRunes(cleanPromptText(profile.Purpose), 120)
	profile.DesiredOutput = limitPromptRunes(cleanPromptText(profile.DesiredOutput), 60)
	profile.CurrentBlockers = normalizePromptList(profile.CurrentBlockers, 2, 40)
	profile.PreferredActivities = normalizePromptList(profile.PreferredActivities, 2, 40)
	profile.SuccessCriteria = normalizePromptList(profile.SuccessCriteria, 2, 60)
	profile.Source = strings.TrimSpace(profile.Source)
	profile.Confidence = strings.TrimSpace(profile.Confidence)
	if isEmptyLearningIntentProfile(profile) {
		return LearningIntentProfile{}
	}
	return profile
}

func buildLearningIntentPromptSummary(profile LearningIntentProfile) string {
	profile = NormalizeLearningIntentProfile(profile)
	if isEmptyLearningIntentProfile(profile) {
		return ""
	}
	parts := make([]string, 0, 6)
	if profile.LearnerLevel != "" {
		parts = append(parts, "수준="+profile.LearnerLevel)
	}
	if profile.Purpose != "" {
		parts = append(parts, "목적="+profile.Purpose)
	}
	if profile.DesiredOutput != "" {
		parts = append(parts, "산출물="+profile.DesiredOutput)
	}
	if len(profile.CurrentBlockers) > 0 {
		parts = append(parts, "막힘="+strings.Join(profile.CurrentBlockers, ", "))
	}
	if len(profile.PreferredActivities) > 0 {
		parts = append(parts, "선호활동="+strings.Join(profile.PreferredActivities, ", "))
	}
	if len(profile.SuccessCriteria) > 0 {
		parts = append(parts, "성공기준="+strings.Join(profile.SuccessCriteria, ", "))
	}
	return limitPromptRunes(strings.Join(parts, " | "), 260)
}

func isEmptyLearningIntentProfile(profile LearningIntentProfile) bool {
	return profile.ConfirmedGoal == "" && profile.LearnerLevel == "" && profile.Purpose == "" && profile.DesiredOutput == "" &&
		len(profile.CurrentBlockers) == 0 && len(profile.PreferredActivities) == 0 && len(profile.SuccessCriteria) == 0
}

func normalizeIntentLearnerLevel(value string) string {
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

func normalizePromptList(values []string, maxItems, maxRunes int) []string {
	if maxItems <= 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = limitPromptRunes(cleanPromptText(value), maxRunes)
		if value == "" || containsPromptText(out, value) {
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

func cleanPromptText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func limitPromptRunes(value string, limit int) string {
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

func containsPromptText(values []string, target string) bool {
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
