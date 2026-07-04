package curriculum

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

const (
	SearchSpecSourceGenerated              = "generated"
	SearchSpecSourcePatternTemplate        = "pattern_template"
	SearchSpecSourcePatternTemplateOverlay = "pattern_template_profile_overlay"
	SearchSpecSourceFallback               = "fallback"
	SearchSpecSourceFallbackOverlay        = "fallback_profile_overlay"
	SearchSpecSourceManual                 = "manual"
)

type LessonRecommendationSearchSpec struct {
	PrimaryQuery string   `json:"primary_query,omitempty"`
	Intent       string   `json:"intent,omitempty"`
	MustInclude  []string `json:"must_include,omitempty"`
	NiceToHave   []string `json:"nice_to_have,omitempty"`
	Avoid        []string `json:"avoid,omitempty"`
	ContentTypes []string `json:"content_types,omitempty"`
	Language     string   `json:"language,omitempty"`
	StageRole    string   `json:"stage_role,omitempty"`
	Source       string   `json:"source,omitempty"`
}

func (spec *LessonRecommendationSearchSpec) UnmarshalJSON(raw []byte) error {
	type lessonRecommendationSearchSpecAlias LessonRecommendationSearchSpec
	var payload struct {
		lessonRecommendationSearchSpecAlias
		MustInclude  any `json:"must_include"`
		NiceToHave   any `json:"nice_to_have"`
		Avoid        any `json:"avoid"`
		ContentTypes any `json:"content_types"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return err
	}
	*spec = LessonRecommendationSearchSpec(payload.lessonRecommendationSearchSpecAlias)
	spec.MustInclude = coerceSearchSpecStringList(payload.MustInclude)
	spec.NiceToHave = coerceSearchSpecStringList(payload.NiceToHave)
	spec.Avoid = coerceSearchSpecStringList(payload.Avoid)
	spec.ContentTypes = coerceSearchSpecStringList(payload.ContentTypes)
	return nil
}

func NormalizeLessonRecommendationSearchSpec(spec LessonRecommendationSearchSpec, fallback LessonRecommendationSearchSpec) LessonRecommendationSearchSpec {
	spec.PrimaryQuery = truncateRunes(strings.Join(strings.Fields(spec.PrimaryQuery), " "), 120)
	if spec.PrimaryQuery == "" {
		spec.PrimaryQuery = fallback.PrimaryQuery
	}
	spec.Intent = normalizeSearchSpecIntent(firstNonEmpty(spec.Intent, fallback.Intent))
	spec.MustInclude = normalizeSearchSpecTokens(spec.MustInclude, 6)
	if len(spec.MustInclude) == 0 {
		spec.MustInclude = normalizeSearchSpecTokens(fallback.MustInclude, 6)
	}
	spec.NiceToHave = normalizeSearchSpecTokens(spec.NiceToHave, 8)
	if len(spec.NiceToHave) == 0 {
		spec.NiceToHave = normalizeSearchSpecTokens(fallback.NiceToHave, 8)
	}
	spec.NiceToHave = removeSearchSpecTokens(spec.NiceToHave, spec.MustInclude)
	spec.Avoid = normalizeSearchSpecTokens(spec.Avoid, 10)
	if len(spec.Avoid) == 0 {
		spec.Avoid = normalizeSearchSpecTokens(fallback.Avoid, 10)
	}
	spec.ContentTypes = normalizeSearchSpecContentTypes(spec.ContentTypes)
	if len(spec.ContentTypes) == 0 {
		spec.ContentTypes = normalizeSearchSpecContentTypes(fallback.ContentTypes)
	}
	spec.Language = normalizeLearningLanguage(firstNonEmpty(spec.Language, fallback.Language))
	spec.StageRole = normalizeSearchSpecStageRole(firstNonEmpty(spec.StageRole, fallback.StageRole))
	spec.Source = normalizeSearchSpecSource(firstNonEmpty(spec.Source, fallback.Source))
	if spec.Source == "" {
		spec.Source = SearchSpecSourceFallback
	}
	return spec
}

func BuildFallbackLessonRecommendationSearchSpec(req CreateCourseDraftRequest, lesson generatedMainLesson, index int) LessonRecommendationSearchSpec {
	language := normalizeLearningLanguage(req.LearningLanguage)
	topic := strings.TrimSpace(req.SourceQuery)
	goal := strings.TrimSpace(derefString(req.LearningGoal))
	title := strings.TrimSpace(lesson.Title)
	objective := strings.TrimSpace(lesson.Objective)
	context := strings.ToLower(strings.Join(compactNonEmptyStrings(topic, goal, title, objective), " "))
	lessonContext := strings.ToLower(strings.Join(compactNonEmptyStrings(title, objective), " "))
	stageRole := inferFallbackSearchStageRoleFromLessonText(index, lessonContext)
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: strings.Join(compactNonEmptyStrings(topic, title, searchSpecIntentWord(language, "tutorial")), " "),
		Intent:       "tutorial",
		MustInclude:  firstNNonEmpty([]string{primarySearchTermFromText(topic), primarySearchTermFromText(title), primarySearchTermFromText(goal)}, 3),
		NiceToHave:   []string{searchSpecIntentWord(language, "basic"), searchSpecIntentWord(language, "practice")},
		Avoid:        defaultSearchSpecAvoidTerms(language),
		ContentTypes: []string{"video"},
		Language:     language,
		StageRole:    stageRole,
		Source:       SearchSpecSourceFallback,
	}

	switch {
	case containsAny(context, "가죽공예", "leathercraft", "leather craft", "leather wallet", "wallet"):
		if language == "en" {
			spec.PrimaryQuery = strings.Join(compactNonEmptyStrings("leathercraft", leathercraftSearchStageTermEnglish(stageRole), "beginner tutorial"), " ")
			spec.MustInclude, spec.NiceToHave = leathercraftSearchTokensEnglish(stageRole)
			spec.Avoid = []string{"shopping", "store", "sale", "asmr"}
		} else {
			spec.PrimaryQuery = strings.Join(compactNonEmptyStrings("가죽공예", leathercraftSearchStageTermKorean(stageRole), "기초 튜토리얼"), " ")
			spec.MustInclude, spec.NiceToHave = leathercraftSearchTokensKorean(stageRole)
			spec.Avoid = []string{"재료 판매", "공구 쇼핑몰", "광고", "asmr"}
		}
	case containsAny(context, "여행 영어", "여행영어", "travel english", "airport", "hotel", "공항", "호텔"):
		if language == "en" {
			spec.PrimaryQuery = "travel English airport hotel conversation tutorial"
			spec.MustInclude = []string{"travel english", "airport", "hotel"}
			spec.NiceToHave = []string{"conversation", "phrases", "check in", "immigration"}
			spec.Avoid = []string{"toeic", "exam", "grammar lecture", "study abroad"}
		} else {
			spec.PrimaryQuery = "여행 영어 공항 호텔 회화 튜토리얼"
			spec.MustInclude = []string{"여행영어", "공항", "호텔"}
			spec.NiceToHave = []string{"회화", "필수", "표현", "체크인", "입국심사"}
			spec.Avoid = []string{"광고", "유학", "토익", "시험", "문법 강의"}
		}
	}

	spec = applyPatternSearchSpecTemplateFallback(req, spec, stageRole, index)
	spec = ApplyLearningIntentProfileToLessonSearchSpec(spec, req.LearningIntent)

	return NormalizeLessonRecommendationSearchSpec(spec, LessonRecommendationSearchSpec{
		Intent:       "tutorial",
		ContentTypes: []string{"video"},
		Language:     language,
		StageRole:    stageRole,
		Source:       SearchSpecSourceFallback,
	})
}

func applyPatternSearchSpecTemplateFallback(req CreateCourseDraftRequest, spec LessonRecommendationSearchSpec, stageRole string, index int) LessonRecommendationSearchSpec {
	template := req.PatternSearchSpecTemplate
	if strings.TrimSpace(template.Intent) == "" && len(template.StageTemplates) == 0 {
		return spec
	}
	stage := findPatternSearchSpecStageTemplate(template.StageTemplates, stageRole, index)
	if strings.TrimSpace(stage.Intent) != "" {
		spec.Intent = stage.Intent
	} else if strings.TrimSpace(template.Intent) != "" {
		spec.Intent = template.Intent
	}
	if len(stage.MustInclude) > 0 {
		spec.MustInclude = stage.MustInclude
	} else if len(template.MustInclude) > 0 {
		spec.MustInclude = template.MustInclude
	}
	if len(stage.NiceToHave) > 0 {
		spec.NiceToHave = stage.NiceToHave
	} else if len(template.NiceToHave) > 0 {
		spec.NiceToHave = template.NiceToHave
	}
	if len(stage.Avoid) > 0 {
		spec.Avoid = stage.Avoid
	} else if len(template.Avoid) > 0 {
		spec.Avoid = template.Avoid
	}
	if len(stage.ContentTypes) > 0 {
		spec.ContentTypes = stage.ContentTypes
	} else if len(template.ContentTypes) > 0 {
		spec.ContentTypes = template.ContentTypes
	}
	spec.Language = normalizeLearningLanguage(firstNonEmpty(stage.Language, template.Language, req.LearningLanguage))
	spec.StageRole = firstNonEmpty(stage.StageRole, stageRole)
	spec.Source = SearchSpecSourcePatternTemplate
	spec.PrimaryQuery = buildPatternTemplatePrimaryQuery(spec, stage)
	return spec
}

func ApplyLearningIntentProfileToLessonSearchSpec(spec LessonRecommendationSearchSpec, profile LearningIntentProfile) LessonRecommendationSearchSpec {
	profile = NormalizeLearningIntentProfile(profile)
	if isEmptyLearningIntentProfile(profile) {
		return spec
	}
	if spec.Source == SearchSpecSourceManual || spec.Source == SearchSpecSourceGenerated {
		return spec
	}
	changed := false

	if token := profileSearchSpecOutputToken(profile.DesiredOutput); token != "" {
		if shouldUseProfileOutputAsMust(token) {
			if updated := appendSearchSpecToken(spec.MustInclude, token, 6); !searchSpecTokensEqual(spec.MustInclude, updated) {
				spec.MustInclude = updated
				changed = true
			}
		} else if updated := appendSearchSpecToken(spec.NiceToHave, token, 8); !searchSpecTokensEqual(spec.NiceToHave, updated) {
			spec.NiceToHave = updated
			changed = true
		}
	}
	for _, token := range profile.CurrentBlockers {
		if token = profileSearchSpecHintToken(token); token != "" {
			if updated := appendSearchSpecToken(spec.NiceToHave, token, 8); !searchSpecTokensEqual(spec.NiceToHave, updated) {
				spec.NiceToHave = updated
				changed = true
			}
		}
	}
	for _, token := range profile.PreferredActivities {
		if token = profileSearchSpecHintToken(token); token != "" {
			if updated := appendSearchSpecToken(spec.NiceToHave, token, 8); !searchSpecTokensEqual(spec.NiceToHave, updated) {
				spec.NiceToHave = updated
				changed = true
			}
		}
	}
	for _, token := range profile.SuccessCriteria {
		if token = profileSearchSpecSuccessToken(token); token != "" {
			if updated := appendSearchSpecToken(spec.NiceToHave, token, 8); !searchSpecTokensEqual(spec.NiceToHave, updated) {
				spec.NiceToHave = updated
				changed = true
			}
		}
	}
	if !changed {
		return spec
	}
	spec.MustInclude = normalizeSearchSpecTokens(spec.MustInclude, 6)
	spec.NiceToHave = normalizeSearchSpecTokens(spec.NiceToHave, 8)
	spec.PrimaryQuery = buildProfileOverlayPrimaryQuery(spec)
	switch spec.Source {
	case SearchSpecSourcePatternTemplate:
		spec.Source = SearchSpecSourcePatternTemplateOverlay
	case SearchSpecSourceFallback, "":
		spec.Source = SearchSpecSourceFallbackOverlay
	}
	return spec
}

func searchSpecTokensEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func appendSearchSpecToken(values []string, token string, limit int) []string {
	token = strings.TrimSpace(token)
	if token == "" {
		return values
	}
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), token) {
			return values
		}
	}
	if len(values) >= limit {
		return values
	}
	return append(values, token)
}

func buildProfileOverlayPrimaryQuery(spec LessonRecommendationSearchSpec) string {
	tokens := append([]string{}, spec.MustInclude...)
	tokens = append(tokens, firstNNonEmpty(spec.NiceToHave, 3)...)
	tokens = normalizeSearchSpecTokens(tokens, 5)
	intentWord := searchSpecIntentWord(spec.Language, spec.Intent)
	return strings.Join(compactNonEmptyStrings(append(tokens, intentWord)...), " ")
}

func profileSearchSpecOutputToken(value string) string {
	value = truncateRunes(strings.Join(strings.Fields(value), " "), 32)
	if value == "" || isGenericProfileSearchToken(value) {
		return ""
	}
	return value
}

func profileSearchSpecHintToken(value string) string {
	value = truncateRunes(strings.Join(strings.Fields(value), " "), 24)
	if value == "" || isGenericProfileSearchToken(value) {
		return ""
	}
	return value
}

func profileSearchSpecSuccessToken(value string) string {
	value = strings.TrimSpace(value)
	for _, suffix := range []string{" 완성", " 완료", " 이해", " 사용 이해"} {
		value = strings.TrimSpace(strings.TrimSuffix(value, suffix))
	}
	value = strings.Join(strings.Fields(value), " ")
	if utf8.RuneCountInString(value) > 14 {
		return ""
	}
	value = truncateRunes(value, 24)
	if value == "" || isGenericProfileSearchToken(value) {
		return ""
	}
	return value
}

func shouldUseProfileOutputAsMust(token string) bool {
	return utf8.RuneCountInString(strings.TrimSpace(token)) >= 3 && !isGenericProfileSearchToken(token)
}

func isGenericProfileSearchToken(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	switch lower {
	case "프로젝트", "실습", "프로젝트 실습", "반복 연습", "기초 따라하기", "따라 만들기", "단계별", "쉬운", "쉽게", "연습", "기초", "입문", "초보", "결과물", "작품", "연주", "그림", "글", "project", "practice", "beginner", "basic", "step by step", "easy", "output", "result":
		return true
	default:
		return false
	}
}

func findPatternSearchSpecStageTemplate(stages []LessonRecommendationSearchSpec, stageRole string, index int) LessonRecommendationSearchSpec {
	stageRole = strings.TrimSpace(stageRole)
	for _, stage := range stages {
		if strings.TrimSpace(stage.StageRole) == stageRole {
			return stage
		}
	}
	if index >= 0 && index < len(stages) {
		return stages[index]
	}
	return LessonRecommendationSearchSpec{}
}

func buildPatternTemplatePrimaryQuery(spec LessonRecommendationSearchSpec, stage LessonRecommendationSearchSpec) string {
	tokens := append([]string{}, spec.MustInclude...)
	tokens = append(tokens, firstNNonEmpty(spec.NiceToHave, 3)...)
	tokens = normalizeSearchSpecTokens(tokens, 5)
	intentWord := searchSpecIntentWord(spec.Language, spec.Intent)
	query := strings.Join(compactNonEmptyStrings(append(tokens, intentWord)...), " ")
	if strings.TrimSpace(stage.PrimaryQuery) != "" {
		query = strings.TrimSpace(stage.PrimaryQuery)
	}
	return query
}

func BuildFallbackLessonRecommendationSearchSpecForLesson(sourceQuery string, learningGoal *string, language string, title string, objective *string, index int) LessonRecommendationSearchSpec {
	return BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(
		sourceQuery,
		learningGoal,
		language,
		title,
		objective,
		index,
		LearningIntentProfile{},
	)
}

func BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(sourceQuery string, learningGoal *string, language string, title string, objective *string, index int, profile LearningIntentProfile) LessonRecommendationSearchSpec {
	return BuildFallbackLessonRecommendationSearchSpec(CreateCourseDraftRequest{
		SourceQuery:      sourceQuery,
		LearningGoal:     learningGoal,
		LearningLanguage: language,
		LearningIntent:   profile,
	}, generatedMainLesson{
		Title:     title,
		Objective: derefString(objective),
	}, index)
}

func MarshalLessonRecommendationSearchSpec(spec LessonRecommendationSearchSpec) (string, error) {
	if strings.TrimSpace(spec.PrimaryQuery) == "" {
		return "{}", nil
	}
	payload, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func ParseLessonRecommendationSearchSpec(raw []byte) LessonRecommendationSearchSpec {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "" || strings.TrimSpace(string(raw)) == "{}" {
		return LessonRecommendationSearchSpec{}
	}
	var spec LessonRecommendationSearchSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return LessonRecommendationSearchSpec{}
	}
	return NormalizeLessonRecommendationSearchSpec(spec, LessonRecommendationSearchSpec{})
}

func ScoreCandidateWithLessonSearchSpec(candidate ContentSearchCandidate, spec LessonRecommendationSearchSpec) int {
	spec = NormalizeLessonRecommendationSearchSpec(spec, LessonRecommendationSearchSpec{})
	if strings.TrimSpace(spec.PrimaryQuery) == "" {
		return 0
	}
	searchText := strings.ToLower(strings.Join(compactNonEmptyStrings(candidate.Title, derefString(candidate.Description), derefString(candidate.ExternalURL)), " "))
	score := 0
	mustHits := countSearchSpecTokenHits(searchText, spec.MustInclude)
	if mustHits >= 2 {
		score += 3
	}
	if countSearchSpecTokenHits(searchText, spec.NiceToHave) > 0 {
		score += 1
	}
	if countSearchSpecAvoidHits(searchText, spec.Avoid) > 0 {
		score -= 5
	}
	if strings.EqualFold(spec.Intent, "tutorial") && containsAny(searchText, "튜토리얼", "강좌", "방법", "how to", "tutorial", "lesson") {
		score += 2
	}
	if spec.StageRole == "setup_intro" {
		if containsAny(searchText, "입문", "기초", "초보", "도구", "beginner", "basic", "tools") {
			score += 2
		}
		if containsAny(searchText, "새들스티치", "saddle stitch", "고급", "advanced") {
			score -= 1
		}
	}
	if searchSpecPrefersContentType(spec.ContentTypes, "video") && searchSpecCandidateLooksVideo(candidate) {
		score += 2
	}
	return score
}

func searchSpecPrefersContentType(contentTypes []string, contentType string) bool {
	for _, value := range contentTypes {
		if strings.EqualFold(strings.TrimSpace(value), contentType) {
			return true
		}
	}
	return false
}

func searchSpecCandidateLooksVideo(candidate ContentSearchCandidate) bool {
	contentType := strings.ToLower(strings.TrimSpace(candidate.ContentType))
	if contentType == "youtube" || contentType == "video" {
		return true
	}
	url := strings.ToLower(strings.TrimSpace(derefString(candidate.ExternalURL)))
	return strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")
}

func normalizeSearchSpecIntent(intent string) string {
	switch strings.ToLower(strings.TrimSpace(intent)) {
	case "tutorial", "lecture", "practice", "example", "checklist", "reference", "course":
		return strings.ToLower(strings.TrimSpace(intent))
	default:
		return "tutorial"
	}
}

func normalizeSearchSpecSource(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case SearchSpecSourceGenerated, SearchSpecSourcePatternTemplate, SearchSpecSourcePatternTemplateOverlay, SearchSpecSourceFallback, SearchSpecSourceFallbackOverlay, SearchSpecSourceManual:
		return strings.ToLower(strings.TrimSpace(source))
	default:
		return ""
	}
}

func normalizeSearchSpecStageRole(stageRole string) string {
	stageRole = strings.ToLower(strings.TrimSpace(stageRole))
	if stageRole == "" {
		return "core_pattern"
	}
	return stageRole
}

func normalizeSearchSpecContentTypes(values []string) []string {
	allowed := map[string]struct{}{"video": {}, "article": {}, "course": {}, "document": {}}
	result := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		trimmed := strings.ToLower(strings.TrimSpace(value))
		if _, ok := allowed[trimmed]; !ok {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func coerceSearchSpecStringList(value any) []string {
	switch typed := value.(type) {
	case nil:
		return nil
	case []string:
		return typed
	case []any:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				result = append(result, splitSearchSpecTokenText(text)...)
			}
		}
		return result
	case string:
		return splitSearchSpecTokenText(typed)
	default:
		return nil
	}
}

func splitSearchSpecTokenText(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return strings.ContainsRune(",/|\n\t;", r)
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func removeSearchSpecTokens(values []string, disallowed []string) []string {
	if len(values) == 0 || len(disallowed) == 0 {
		return values
	}
	blocked := make(map[string]struct{}, len(disallowed))
	for _, value := range disallowed {
		key := strings.ToLower(strings.TrimSpace(value))
		if key != "" {
			blocked[key] = struct{}{}
		}
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := blocked[strings.ToLower(strings.TrimSpace(value))]; exists {
			continue
		}
		result = append(result, value)
	}
	return result
}

func normalizeSearchSpecTokens(values []string, limit int) []string {
	result := make([]string, 0, limit)
	seen := map[string]struct{}{}
	for _, value := range values {
		trimmed := truncateRunes(strings.Join(strings.Fields(value), " "), 32)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func countSearchSpecTokenHits(searchText string, tokens []string) int {
	hits := 0
	for _, token := range tokens {
		if strings.TrimSpace(token) == "" {
			continue
		}
		if strings.Contains(searchText, strings.ToLower(strings.TrimSpace(token))) {
			hits++
		}
	}
	return hits
}

func countSearchSpecAvoidHits(searchText string, tokens []string) int {
	hits := 0
	for _, token := range tokens {
		trimmed := strings.ToLower(strings.TrimSpace(token))
		if trimmed == "" {
			continue
		}
		if strings.Contains(searchText, trimmed) && !searchSpecNegatesAvoid(searchText, trimmed) {
			hits++
		}
	}
	return hits
}

func searchSpecNegatesAvoid(searchText, avoid string) bool {
	negativePhrases := []string{
		avoid + " 없음",
		avoid + " 아님",
		avoid + " 아닙니다",
		"no " + avoid,
		"not " + avoid,
		"not sponsored",
		"no sponsor",
		"협찬 아님",
		"협찬없음",
		"광고 없음",
		"광고 하나도 없음",
	}
	for _, phrase := range negativePhrases {
		if strings.Contains(searchText, strings.ToLower(phrase)) {
			return true
		}
	}
	return false
}

func inferFallbackSearchStageRoleFromLessonText(index int, lessonContext string) string {
	lessonContext = strings.ToLower(strings.TrimSpace(lessonContext))
	if containsAny(lessonContext, "새들스티치", "saddle stitch", "바느질", "스티치", "재단", "본드", "stitching", "cutting", "glue") {
		return "core_pattern"
	}
	if containsAny(lessonContext, "마감", "완성", "결과물", "finish", "finishing", "complete", "final") {
		return "artifact_finish"
	}
	if containsAny(lessonContext, "도구", "재료", "입문", "기초", "초보", "tools", "materials", "beginner", "basic") {
		return "setup_intro"
	}
	return inferFallbackSearchStageRole(index)
}

func inferFallbackSearchStageRole(index int) string {
	switch index {
	case 0:
		return "setup_intro"
	case 1:
		return "first_output"
	case 2:
		return "core_pattern"
	default:
		return "integration_practice"
	}
}

func leathercraftSearchStageTermKorean(stageRole string) string {
	switch stageRole {
	case "setup_intro":
		return "입문 도구 사용법"
	case "first_output":
		return "첫 소품 제작"
	case "artifact_finish":
		return "지갑 완성 마감"
	default:
		return "바느질 재단 방법"
	}
}

func leathercraftSearchStageTermEnglish(stageRole string) string {
	switch stageRole {
	case "setup_intro":
		return "basic tools"
	case "first_output":
		return "first small project"
	case "artifact_finish":
		return "wallet finishing"
	default:
		return "cutting stitching"
	}
}

func leathercraftSearchTokensKorean(stageRole string) ([]string, []string) {
	switch stageRole {
	case "setup_intro":
		return []string{"가죽공예", "도구", "기초"}, []string{"입문", "재료", "초보"}
	case "artifact_finish":
		return []string{"가죽공예", "지갑", "마감"}, []string{"완성", "엣지", "코바"}
	default:
		return []string{"가죽공예", "바느질", "재단"}, []string{"새들스티치", "본드", "지갑"}
	}
}

func leathercraftSearchTokensEnglish(stageRole string) ([]string, []string) {
	switch stageRole {
	case "setup_intro":
		return []string{"leathercraft", "tools", "beginner"}, []string{"materials", "basic", "starter"}
	case "artifact_finish":
		return []string{"leathercraft", "wallet", "finishing"}, []string{"edge", "burnish", "complete"}
	default:
		return []string{"leathercraft", "stitching", "cutting"}, []string{"saddle stitch", "glue", "wallet"}
	}
}

func searchSpecIntentWord(language, key string) string {
	if normalizeLearningLanguage(language) == "en" {
		switch key {
		case "tutorial":
			return "tutorial"
		case "basic":
			return "basic"
		case "practice":
			return "practice"
		}
	}
	switch key {
	case "tutorial":
		return "튜토리얼"
	case "basic":
		return "기초"
	case "practice":
		return "연습"
	default:
		return key
	}
}

func defaultSearchSpecAvoidTerms(language string) []string {
	if normalizeLearningLanguage(language) == "en" {
		return []string{"shopping", "sale", "advertisement", "exam"}
	}
	return []string{"광고", "판매", "쇼핑몰", "시험"}
}

func primarySearchTermFromText(value string) string {
	fields := strings.Fields(strings.TrimSpace(value))
	if len(fields) == 0 {
		return ""
	}
	return truncateRunes(fields[0], 32)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 || utf8.RuneCountInString(value) <= limit {
		return value
	}
	runes := []rune(value)
	return string(runes[:limit])
}
