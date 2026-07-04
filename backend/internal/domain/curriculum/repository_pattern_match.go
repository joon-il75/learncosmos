package curriculum

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type CurriculumPatternMatch struct {
	PatternKey                       string
	PatternVersion                   string
	Domain                           string
	GoalType                         string
	Language                         string
	Title                            string
	Summary                          string
	StageRules                       []string
	RecommendedSequence              []string
	CompletionCriteriaRules          []string
	BadPatterns                      []string
	RecommendationSearchSpecTemplate CurriculumPatternRecommendationSearchSpecTemplate
	SimilarityScore                  float64
	Distance                         float64
}

func (r *Repository) SearchCurriculumPatternMatches(ctx context.Context, pgVector, provider, model string, dimension int, language string, limit int) ([]CurriculumPatternMatch, error) {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "embedding_gemma"
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = "embedding-gemma"
	}
	language = normalizeLearningLanguage(language)
	if dimension <= 0 {
		dimension = 768
	}
	if limit <= 0 {
		limit = 3
	}
	if strings.TrimSpace(pgVector) == "" {
		return nil, fmt.Errorf("query vector is required")
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
		  cp.pattern_key,
		  cp.pattern_version,
		  cp.domain,
		  cp.goal_type,
		  cp.language,
		  cp.title,
		  cp.summary,
		  cp.stage_rules,
		  cp.recommended_sequence,
		  COALESCE(cp.guidance->'completion_criteria_rules', '[]'::jsonb),
		  cp.bad_patterns,
		  cp.recommendation_search_spec_template,
		  cpe.embedding <=> $1::vector AS distance,
		  1 - (cpe.embedding <=> $1::vector) AS similarity_score
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cp.language = $2
		  AND cpe.provider = $3
		  AND cpe.model = $4
		  AND cpe.dimension = $5
		  AND cpe.status = 'ready'
		  AND cpe.embedding IS NOT NULL
		ORDER BY cpe.embedding <=> $1::vector ASC
		LIMIT $6
	`, pgVector, language, provider, model, dimension, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []CurriculumPatternMatch{}
	for rows.Next() {
		var item CurriculumPatternMatch
		var stageRulesJSON, recommendedSequenceJSON, completionCriteriaRulesJSON, badPatternsJSON, searchSpecTemplateJSON []byte
		if err := rows.Scan(
			&item.PatternKey,
			&item.PatternVersion,
			&item.Domain,
			&item.GoalType,
			&item.Language,
			&item.Title,
			&item.Summary,
			&stageRulesJSON,
			&recommendedSequenceJSON,
			&completionCriteriaRulesJSON,
			&badPatternsJSON,
			&searchSpecTemplateJSON,
			&item.Distance,
			&item.SimilarityScore,
		); err != nil {
			return nil, err
		}
		item.StageRules = decodePatternStringArray(stageRulesJSON)
		item.RecommendedSequence = decodePatternStringArray(recommendedSequenceJSON)
		item.CompletionCriteriaRules = decodePatternStringArray(completionCriteriaRulesJSON)
		item.BadPatterns = decodePatternStringArray(badPatternsJSON)
		item.RecommendationSearchSpecTemplate = decodeCurriculumPatternSearchSpecTemplate(searchSpecTemplateJSON)
		matches = append(matches, item)
	}
	return matches, rows.Err()
}

func BuildCurriculumPatternQueryText(req CreateCourseDraftRequest) string {
	return strings.Join(compactNonEmptyStrings(
		req.SourceQuery,
		valueOrEmpty(req.LearningGoal),
		valueOrEmpty(req.GoalUserIntent),
		valueOrEmpty(req.GoalUsageContext),
		valueOrEmpty(req.GoalOutputType),
		valueOrEmpty(req.GoalDifficultyLevel),
	), "\n")
}

func FormatFloat32VectorForPG(vector []float32) string {
	parts := make([]string, 0, len(vector))
	for _, value := range vector {
		parts = append(parts, fmt.Sprintf("%f", value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func BuildCurriculumPatternMatchGuidance(matches []CurriculumPatternMatch, limit int) string {
	return BuildCurriculumPatternMatchGuidanceForLanguage(matches, limit, "ko")
}

func BuildCurriculumPatternMatchGuidanceForLanguage(matches []CurriculumPatternMatch, limit int, language string) string {
	if limit <= 0 || limit > len(matches) {
		limit = len(matches)
	}
	if limit == 0 {
		return ""
	}

	var b strings.Builder
	for idx := 0; idx < limit; idx++ {
		match := matches[idx]
		title := strings.TrimSpace(match.Title)
		if title == "" {
			title = strings.TrimSpace(match.PatternKey)
		}
		if title == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("%d. %s", idx+1, title))
		if strings.TrimSpace(match.PatternKey) != "" {
			b.WriteString(" (" + strings.TrimSpace(match.PatternKey) + ")")
		}
		b.WriteString("\n")
		if strings.TrimSpace(match.Summary) != "" {
			b.WriteString("- " + patternGuidanceLabel(language, "summary") + ": " + strings.TrimSpace(match.Summary) + "\n")
		}
		if len(match.RecommendedSequence) > 0 {
			b.WriteString("- " + patternGuidanceLabel(language, "recommended_flow") + ": " + strings.Join(match.RecommendedSequence, " -> ") + "\n")
		}
		for _, rule := range firstNNonEmpty(match.StageRules, 2) {
			b.WriteString("- " + patternGuidanceLabel(language, "stage_rule") + ": " + rule + "\n")
		}
		for _, rule := range firstNNonEmpty(match.CompletionCriteriaRules, 2) {
			b.WriteString("- " + patternGuidanceLabel(language, "completion_rule") + ": " + rule + "\n")
		}
		if len(match.BadPatterns) > 0 {
			b.WriteString("- " + patternGuidanceLabel(language, "bad_patterns") + ": " + strings.Join(firstNNonEmpty(match.BadPatterns, 4), ", ") + "\n")
		}
		if summary := summarizeCurriculumPatternSearchSpecTemplate(match.RecommendationSearchSpecTemplate, language); summary != "" {
			b.WriteString("- " + patternGuidanceLabel(language, "search_spec_template") + ": " + summary + "\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func patternGuidanceLabel(language, key string) string {
	if normalizeLearningLanguage(language) == "en" {
		switch key {
		case "summary":
			return "Summary"
		case "recommended_flow":
			return "Recommended flow"
		case "stage_rule":
			return "Stage rule"
		case "completion_rule":
			return "Completion criteria rule"
		case "bad_patterns":
			return "Lesson patterns to avoid"
		case "search_spec_template":
			return "Recommendation search spec template"
		}
	}
	switch key {
	case "summary":
		return "요약"
	case "recommended_flow":
		return "권장 흐름"
	case "stage_rule":
		return "단계 규칙"
	case "completion_rule":
		return "완료 기준 규칙"
	case "bad_patterns":
		return "피할 리슨 예"
	case "search_spec_template":
		return "추천 검색 명세 template"
	default:
		return key
	}
}

func FilterCurriculumPatternGuidanceMatches(matches []CurriculumPatternMatch, domain, goalType, preferredPatternKey string, limit int, minScore float64) []CurriculumPatternMatch {
	if limit <= 0 || len(matches) == 0 {
		return nil
	}
	domain = strings.TrimSpace(domain)
	goalType = strings.TrimSpace(goalType)
	preferredPatternKey = strings.TrimSpace(preferredPatternKey)
	if preferredPatternKey != "" {
		for _, match := range matches {
			if minScore > 0 && match.SimilarityScore < minScore {
				continue
			}
			if strings.TrimSpace(match.PatternKey) == preferredPatternKey {
				return []CurriculumPatternMatch{match}
			}
		}
	}
	if domain == "" || goalType == "" {
		for _, match := range matches {
			if strings.TrimSpace(match.Domain) != "" && strings.TrimSpace(match.GoalType) != "" {
				domain = strings.TrimSpace(match.Domain)
				goalType = strings.TrimSpace(match.GoalType)
				break
			}
		}
	}

	result := make([]CurriculumPatternMatch, 0, limit)
	for _, match := range matches {
		if minScore > 0 && match.SimilarityScore < minScore {
			continue
		}
		if domain != "" && strings.TrimSpace(match.Domain) != domain {
			continue
		}
		if goalType != "" && strings.TrimSpace(match.GoalType) != goalType {
			continue
		}
		result = append(result, match)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func decodeCurriculumPatternSearchSpecTemplate(raw []byte) CurriculumPatternRecommendationSearchSpecTemplate {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "" || strings.TrimSpace(string(raw)) == "{}" {
		return CurriculumPatternRecommendationSearchSpecTemplate{}
	}
	var value CurriculumPatternRecommendationSearchSpecTemplate
	if err := json.Unmarshal(raw, &value); err != nil {
		return CurriculumPatternRecommendationSearchSpecTemplate{}
	}
	return value
}

func summarizeCurriculumPatternSearchSpecTemplate(template CurriculumPatternRecommendationSearchSpecTemplate, language string) string {
	if strings.TrimSpace(template.Intent) == "" && len(template.StageTemplates) == 0 {
		return ""
	}
	parts := []string{}
	if strings.TrimSpace(template.Intent) != "" {
		parts = append(parts, "intent="+strings.TrimSpace(template.Intent))
	}
	if len(template.ContentTypes) > 0 {
		parts = append(parts, "content_types="+strings.Join(firstNNonEmpty(template.ContentTypes, 3), "/"))
	}
	stageParts := []string{}
	for _, stage := range template.StageTemplates {
		role := strings.TrimSpace(stage.StageRole)
		if role == "" {
			continue
		}
		tokens := append(firstNNonEmpty(stage.MustInclude, 2), firstNNonEmpty(stage.NiceToHave, 2)...)
		if len(tokens) > 0 {
			stageParts = append(stageParts, role+"("+strings.Join(tokens, ", ")+")")
		} else {
			stageParts = append(stageParts, role)
		}
		if len(stageParts) >= 3 {
			break
		}
	}
	if len(stageParts) > 0 {
		parts = append(parts, "stages="+strings.Join(stageParts, " -> "))
	}
	return strings.Join(parts, "; ")
}

func decodePatternStringArray(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil
	}
	return firstNNonEmpty(values, len(values))
}

func firstNNonEmpty(values []string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	result := make([]string, 0, limit)
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func compactNonEmptyStrings(values ...string) []string {
	result := []string{}
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
