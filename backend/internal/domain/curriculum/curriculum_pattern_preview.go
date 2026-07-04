package curriculum

import "strings"

type CurriculumPatternRoutingPreview struct {
	PatternKey     string
	SubpatternKey  string
	RefinementMode string
	DomainAxis     string
	GoalType       string
}

func AnalyzeCurriculumPatternRoutingPreview(sourceQuery, learningGoal, learningLanguage string) CurriculumPatternRoutingPreview {
	goal := strings.TrimSpace(learningGoal)
	req := CreateCourseDraftRequest{
		SourceQuery:      strings.TrimSpace(sourceQuery),
		LearningLanguage: normalizeLearningLanguage(learningLanguage),
	}
	if goal != "" {
		req.LearningGoal = &goal
	}

	patternKey := inferGoalPatternKey(req)
	summary := analyzeGoalRouting(req)
	return CurriculumPatternRoutingPreview{
		PatternKey:     formatGoalPatternKeyForLog(patternKey),
		SubpatternKey:  summary.SubpatternKey,
		RefinementMode: summary.RefinementMode,
		DomainAxis:     patternKey.DomainAxis,
		GoalType:       patternKey.GoalModePrimary,
	}
}

func BuildCurriculumGenerationPromptPreview(req CreateCourseDraftRequest) string {
	req.SourceQuery = strings.TrimSpace(req.SourceQuery)
	req.LearningLanguage = normalizeLearningLanguage(req.LearningLanguage)
	return buildCurriculumGenerationPrompt(req)
}
