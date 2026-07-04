package curriculum

import "strings"

func trimmedGoalField(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func buildGoalInferenceContext(req CreateCourseDraftRequest) string {
	parts := []string{
		trimmedGoalField(req.LearningGoal),
		trimmedGoalField(req.GoalUsageContext),
		trimmedGoalField(req.GoalMotivation),
		trimmedGoalField(req.GoalUserIntent),
		strings.TrimSpace(req.SourceQuery),
	}
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.ToLower(strings.Join(filtered, " "))
}

func buildGoalPromptDetails(req CreateCourseDraftRequest) []string {
	details := make([]string, 0, 7)
	if summary := buildLearningIntentPromptSummary(req.LearningIntent); summary != "" {
		details = append(details, "학습 의도 요약: "+summary)
	}
	if usage := trimmedGoalField(req.GoalUsageContext); usage != "" {
		details = append(details, "활용 맥락: "+usage)
	}
	if motivation := trimmedGoalField(req.GoalMotivation); motivation != "" {
		details = append(details, "동기: "+motivation)
	}
	if level := trimmedGoalField(req.GoalDifficultyLevel); level != "" {
		details = append(details, "목표 난이도: "+level)
	}
	if horizon := trimmedGoalField(req.GoalTimeHorizon); horizon != "" {
		details = append(details, "목표 기간: "+horizon)
	}
	if output := trimmedGoalField(req.GoalOutputType); output != "" {
		details = append(details, "원하는 결과물 형식: "+output)
	}
	if goalType := trimmedGoalField(req.GoalType); goalType != "" {
		details = append(details, "목표 유형: "+goalType)
	}
	return details
}

func buildGoalNarrativeContext(req CreateCourseDraftRequest) string {
	parts := []string{
		trimmedGoalField(req.LearningGoal),
		trimmedGoalField(req.GoalUsageContext),
		trimmedGoalField(req.GoalMotivation),
		trimmedGoalField(req.GoalUserIntent),
	}
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	if len(filtered) == 0 {
		return strings.TrimSpace(req.SourceQuery)
	}
	return strings.Join(filtered, " ")
}
