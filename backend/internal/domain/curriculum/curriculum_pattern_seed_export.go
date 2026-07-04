package curriculum

import (
	"fmt"
	"sort"
	"strings"
)

type CurriculumPatternSeed struct {
	PatternKey                       string                                            `json:"pattern_key"`
	PatternVersion                   string                                            `json:"pattern_version"`
	Domain                           string                                            `json:"domain"`
	GoalType                         string                                            `json:"goal_type"`
	LearnerLevel                     string                                            `json:"learner_level"`
	Language                         string                                            `json:"language"`
	Title                            string                                            `json:"title"`
	Summary                          string                                            `json:"summary"`
	TriggerKeywords                  []string                                          `json:"trigger_keywords"`
	Guidance                         map[string]any                                    `json:"guidance"`
	StageRules                       []string                                          `json:"stage_rules"`
	RecommendedSequence              []string                                          `json:"recommended_sequence"`
	BadPatterns                      []string                                          `json:"bad_patterns"`
	RecommendationSearchSpecTemplate CurriculumPatternRecommendationSearchSpecTemplate `json:"recommendation_search_spec_template"`
	EmbeddingText                    string                                            `json:"embedding_text"`
	IsActive                         bool                                              `json:"is_active"`
	Source                           CurriculumSeedSource                              `json:"source"`
}

type CurriculumSeedSource struct {
	CatalogFile        string `json:"catalog_file"`
	RoutingFile        string `json:"routing_file"`
	RecommendationFile string `json:"recommendation_file"`
}

type CurriculumPatternRecommendationSearchSpecTemplate struct {
	Intent         string                           `json:"intent,omitempty"`
	ContentTypes   []string                         `json:"content_types,omitempty"`
	Language       string                           `json:"language,omitempty"`
	Source         string                           `json:"source,omitempty"`
	MustInclude    []string                         `json:"must_include,omitempty"`
	NiceToHave     []string                         `json:"nice_to_have,omitempty"`
	Avoid          []string                         `json:"avoid,omitempty"`
	StageTemplates []LessonRecommendationSearchSpec `json:"stage_templates,omitempty"`
}

func BuildCurriculumPatternSeeds(version string) ([]CurriculumPatternSeed, error) {
	return BuildCurriculumPatternSeedsForLanguage(version, "ko")
}

func BuildCurriculumPatternSeedsForLanguage(version, language string) ([]CurriculumPatternSeed, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		version = "v1"
	}
	languages, err := normalizeSeedExportLanguages(language)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(goalSubpatternSpecs))
	for key := range goalSubpatternSpecs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	seeds := make([]CurriculumPatternSeed, 0, len(keys)*len(languages))
	seen := map[string]struct{}{}
	for _, key := range keys {
		spec := goalSubpatternSpecs[key]
		domain, goalType, err := splitSubpatternKey(key)
		if err != nil {
			return nil, err
		}
		for _, lang := range languages {
			seenKey := key + "\x00" + lang
			if _, ok := seen[seenKey]; ok {
				return nil, fmt.Errorf("duplicate pattern key/language: %s %s", key, lang)
			}
			seen[seenKey] = struct{}{}
			seed := buildCurriculumPatternSeedForLanguage(key, version, domain, goalType, lang, spec)
			if strings.TrimSpace(seed.EmbeddingText) == "" {
				return nil, fmt.Errorf("empty embedding text: %s language=%s", key, lang)
			}
			seeds = append(seeds, seed)
		}
	}
	return seeds, nil
}

func normalizeSeedExportLanguages(language string) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "", "all":
		return []string{"ko", "en"}, nil
	case "ko", "kr", "kor", "korean":
		return []string{"ko"}, nil
	case "en", "eng", "english":
		return []string{"en"}, nil
	default:
		return nil, fmt.Errorf("unsupported seed language: %s", language)
	}
}

func buildCurriculumPatternSeedForLanguage(key, version, domain, goalType, language string, spec goalSubpatternSpec) CurriculumPatternSeed {
	if language == "en" {
		return buildEnglishCurriculumPatternSeed(key, version, domain, goalType, spec)
	}
	return buildKoreanCurriculumPatternSeed(key, version, domain, goalType, spec)
}

func buildKoreanCurriculumPatternSeed(key, version, domain, goalType string, spec goalSubpatternSpec) CurriculumPatternSeed {
	queryTerms := buildRecommendationGoalTerms(goalPatternKey{
		DomainAxis:      domain,
		GoalModePrimary: goalType,
	}, key, "")
	triggerKeywords := dedupeRecommendationTerms(append(append(queryTerms, tokenizedPatternKeyTerms(key)...), extraCurriculumPatternSeedTerms(key)...))
	summary := buildSeedSummary(spec)
	embeddingText := buildCurriculumPatternEmbeddingText(key, domain, goalType, spec, triggerKeywords, summary)
	return CurriculumPatternSeed{
		PatternKey:      key,
		PatternVersion:  version,
		Domain:          domain,
		GoalType:        goalType,
		LearnerLevel:    "any",
		Language:        "ko",
		Title:           spec.Name,
		Summary:         summary,
		TriggerKeywords: triggerKeywords,
		Guidance: map[string]any{
			"last_lesson_rule":          spec.LastLessonRule,
			"completion_criteria_rules": spec.CompletionCriteriaRules,
			"atomic_disallow_examples":  spec.AtomicDisallowExamples,
		},
		StageRules:                       spec.StageRules,
		RecommendedSequence:              spec.StepRoles,
		BadPatterns:                      spec.AtomicDisallowExamples,
		RecommendationSearchSpecTemplate: buildCurriculumPatternRecommendationSearchSpecTemplate(key, domain, goalType, "ko", spec.StepRoles, triggerKeywords),
		EmbeddingText:                    embeddingText,
		IsActive:                         true,
		Source:                           defaultCurriculumSeedSource(),
	}
}

func buildEnglishCurriculumPatternSeed(key, version, domain, goalType string, spec goalSubpatternSpec) CurriculumPatternSeed {
	triggerKeywords := dedupeRecommendationTerms(append(tokenizedPatternKeyTerms(key), extraEnglishCurriculumPatternSeedTerms(key)...))
	title := titleizePatternKey(key)
	summary := buildEnglishSeedSummary(key, domain, goalType, spec)
	stageRules := buildEnglishStageRules(key, domain, goalType, spec)
	completionRules := buildEnglishCompletionCriteriaRules(goalType)
	badPatterns := buildEnglishBadPatterns()
	if override, ok := englishCurriculumPatternOverrides[key]; ok {
		if strings.TrimSpace(override.Title) != "" {
			title = strings.TrimSpace(override.Title)
		}
		if strings.TrimSpace(override.Summary) != "" {
			summary = strings.TrimSpace(override.Summary)
		}
		if len(override.TriggerKeywords) > 0 {
			triggerKeywords = dedupeRecommendationTerms(append(triggerKeywords, override.TriggerKeywords...))
		}
		if len(override.StageRules) > 0 {
			stageRules = append([]string(nil), override.StageRules...)
		}
		if len(override.CompletionCriteriaRules) > 0 {
			completionRules = append([]string(nil), override.CompletionCriteriaRules...)
		}
		if len(override.BadPatterns) > 0 {
			badPatterns = append([]string(nil), override.BadPatterns...)
		}
	}
	embeddingText := buildEnglishCurriculumPatternEmbeddingText(key, domain, goalType, title, summary, triggerKeywords, spec.StepRoles, stageRules, completionRules, badPatterns)
	return CurriculumPatternSeed{
		PatternKey:                       key,
		PatternVersion:                   version,
		Domain:                           domain,
		GoalType:                         goalType,
		LearnerLevel:                     "any",
		Language:                         "en",
		Title:                            title,
		Summary:                          summary,
		TriggerKeywords:                  triggerKeywords,
		Guidance:                         map[string]any{"last_lesson_rule": buildEnglishLastLessonRule(goalType), "completion_criteria_rules": completionRules, "atomic_disallow_examples": badPatterns},
		StageRules:                       stageRules,
		RecommendedSequence:              spec.StepRoles,
		BadPatterns:                      badPatterns,
		RecommendationSearchSpecTemplate: buildCurriculumPatternRecommendationSearchSpecTemplate(key, domain, goalType, "en", spec.StepRoles, triggerKeywords),
		EmbeddingText:                    embeddingText,
		IsActive:                         true,
		Source:                           defaultCurriculumSeedSource(),
	}
}

func buildCurriculumPatternRecommendationSearchSpecTemplate(key, domain, goalType, language string, stageRoles, triggerKeywords []string) CurriculumPatternRecommendationSearchSpecTemplate {
	language = normalizeLearningLanguage(language)
	if template, ok := buildSpecificCurriculumPatternRecommendationSearchSpecTemplate(key, language); ok {
		return template
	}
	baseSourceTerms := append([]string{}, curriculumPatternSearchDomainTerms(domain, goalType, language)...)
	switch language {
	case "ko":
		baseSourceTerms = append([]string{}, prioritizeKoreanSearchTerms(extraCurriculumPatternSeedTerms(key))...)
		baseSourceTerms = append(baseSourceTerms, prioritizeKoreanSearchTerms(triggerKeywords)...)
		baseSourceTerms = append(baseSourceTerms, curriculumPatternSearchDomainTerms(domain, goalType, language)...)
	case "en":
		baseSourceTerms = append([]string{}, englishPatternSearchKeyPhrase(key))
		baseSourceTerms = append(baseSourceTerms, prioritizeEnglishSearchTerms(extraEnglishCurriculumPatternSeedTerms(key))...)
		baseSourceTerms = append(baseSourceTerms, prioritizeEnglishSearchTerms(triggerKeywords)...)
		baseSourceTerms = append(baseSourceTerms, curriculumPatternSearchDomainTerms(domain, goalType, language)...)
	}
	baseTerms := firstNNonEmpty(dedupeRecommendationTerms(baseSourceTerms), 4)
	niceTerms := curriculumPatternSearchNiceTerms(language)
	avoidTerms := curriculumPatternSearchAvoidTermsForKey(key, goalType, language)
	roles := firstNNonEmpty(stageRoles, len(stageRoles))
	if len(roles) == 0 {
		roles = []string{"setup_intro", "core_pattern", "artifact_finish"}
	}
	stageTemplates := make([]LessonRecommendationSearchSpec, 0, len(roles))
	for _, role := range roles {
		stageTerms := curriculumPatternSearchStageTerms(role, domain, goalType, language)
		mustInclude := firstNNonEmpty(dedupeRecommendationTerms(append(append([]string{}, baseTerms...), firstNNonEmpty(stageTerms, 2)...)), 5)
		niceToHave := firstNNonEmpty(filterRecommendationTermsNotIn(dedupeRecommendationTerms(append(stageTerms, niceTerms...)), mustInclude), 8)
		stageTemplates = append(stageTemplates, LessonRecommendationSearchSpec{
			Intent:       "tutorial",
			MustInclude:  mustInclude,
			NiceToHave:   niceToHave,
			Avoid:        avoidTerms,
			ContentTypes: curriculumPatternSearchContentTypesForKey(key, domain, goalType, language),
			Language:     language,
			StageRole:    normalizeSearchSpecStageRole(role),
			Source:       SearchSpecSourcePatternTemplate,
		})
	}
	return CurriculumPatternRecommendationSearchSpecTemplate{
		Intent:         "tutorial",
		ContentTypes:   curriculumPatternSearchContentTypesForKey(key, domain, goalType, language),
		Language:       language,
		Source:         SearchSpecSourcePatternTemplate,
		MustInclude:    baseTerms,
		NiceToHave:     niceTerms,
		Avoid:          avoidTerms,
		StageTemplates: stageTemplates,
	}
}

func filterRecommendationTermsNotIn(terms, excluded []string) []string {
	if len(terms) == 0 || len(excluded) == 0 {
		return terms
	}
	excludedSet := make(map[string]struct{}, len(excluded))
	for _, term := range excluded {
		if trimmed := strings.ToLower(strings.TrimSpace(term)); trimmed != "" {
			excludedSet[trimmed] = struct{}{}
		}
	}
	filtered := make([]string, 0, len(terms))
	for _, term := range terms {
		trimmed := strings.TrimSpace(term)
		if trimmed == "" {
			continue
		}
		if _, found := excludedSet[strings.ToLower(trimmed)]; found {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	return filtered
}

func curriculumPatternSearchDomainTerms(domain, goalType, language string) []string {
	if normalizeLearningLanguage(language) == "en" {
		return dedupeRecommendationTerms([]string{strings.ReplaceAll(domain, "_", " "), strings.ReplaceAll(goalType, "_", " ")})
	}
	switch strings.TrimSpace(domain) {
	case "craft_making":
		return []string{"공예", "만들기"}
	case "language_communication":
		return []string{"회화", "표현"}
	case "instrument_performance":
		return []string{"악기", "연주"}
	case "visual_art":
		return []string{"미술", "작품"}
	case "digital_creation":
		return []string{"디지털 제작", "실습"}
	case "body_movement":
		return []string{"동작", "루틴"}
	case "cooking_baking":
		return []string{"요리", "베이킹"}
	case "knowledge_hobby":
		return []string{"학습", "실천"}
	case "maker_technical_hobby":
		return []string{"메이커", "기술 실습"}
	case "writing_storytelling":
		return []string{"글쓰기", "스토리텔링"}
	case "knowledge_study":
		return []string{"학습", "정리"}
	default:
		return dedupeRecommendationTerms([]string{strings.ReplaceAll(domain, "_", " "), strings.ReplaceAll(goalType, "_", " ")})
	}
}

func englishPatternSearchKeyPhrase(key string) string {
	parts := strings.Split(strings.TrimSpace(key), ":")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(parts[len(parts)-1], "_", " "))
}

func prioritizeEnglishSearchTerms(terms []string) []string {
	if len(terms) == 0 {
		return nil
	}
	prioritized := make([]string, 0, len(terms))
	deferred := make([]string, 0, len(terms))
	for _, term := range terms {
		if isGenericEnglishSearchAxisTerm(term) {
			deferred = append(deferred, term)
			continue
		}
		prioritized = append(prioritized, term)
	}
	return append(prioritized, deferred...)
}

func isGenericEnglishSearchAxisTerm(term string) bool {
	normalized := strings.ToLower(strings.TrimSpace(term))
	if normalized == "" {
		return true
	}
	genericTerms := map[string]struct{}{
		"artifact": {}, "assessment": {}, "baking": {}, "body": {}, "build": {}, "certification": {},
		"communication": {}, "concept": {}, "cooking": {}, "craft": {}, "creation": {}, "digital": {},
		"execution": {}, "foundation": {}, "habit": {}, "hobby": {}, "instruction": {}, "instrument": {},
		"knowledge": {}, "language": {}, "lifestyle": {}, "maker": {}, "making": {}, "mastery": {},
		"movement": {}, "participation": {}, "performance": {}, "presentation": {}, "professional": {},
		"publish": {}, "skill": {}, "storytelling": {}, "teaching": {}, "technical": {},
		"transition": {}, "visual": {}, "writing": {},
	}
	_, ok := genericTerms[normalized]
	return ok
}

func prioritizeKoreanSearchTerms(terms []string) []string {
	if len(terms) == 0 {
		return nil
	}
	prioritized := make([]string, 0, len(terms))
	deferred := make([]string, 0, len(terms))
	for _, term := range terms {
		if searchTermContainsHangul(term) {
			prioritized = append(prioritized, term)
			continue
		}
		deferred = append(deferred, term)
	}
	return append(prioritized, deferred...)
}

func searchTermContainsHangul(term string) bool {
	for _, r := range term {
		if r >= '가' && r <= '힣' {
			return true
		}
	}
	return false
}

func curriculumPatternSearchAvoidTerms(goalType, language string) []string {
	goalType = strings.TrimSpace(goalType)
	if normalizeLearningLanguage(language) == "en" {
		switch goalType {
		case "certification_assessment":
			return []string{"shopping", "sale", "advertisement", "certificate agency", "score guarantee", "pass guarantee"}
		case "professional_transition":
			return []string{"advertisement", "course sale", "recruitment", "consultation", "job guarantee", "agency"}
		case "presentation_publish":
			return []string{"admission portfolio", "job portfolio", "views guarantee", "monetization only", "ghostwriting", "contest ad"}
		}
		return defaultSearchSpecAvoidTerms(language)
	}
	switch goalType {
	case "certification_assessment":
		return []string{"광고", "판매", "쇼핑몰", "자격 대행", "점수 보장", "합격 보장", "학원 광고"}
	case "professional_transition":
		return []string{"광고", "강의 판매", "국비 모집", "상담 신청", "취업 보장", "대행", "환급 광고"}
	case "presentation_publish":
		return []string{"입시 포트폴리오", "취업 포트폴리오", "조회수 보장", "수익화만", "대필", "공모전 광고"}
	}
	return defaultSearchSpecAvoidTerms(language)
}

func curriculumPatternSearchAvoidTermsForKey(key, goalType, language string) []string {
	avoidTerms := append([]string{}, curriculumPatternSearchAvoidTerms(goalType, language)...)
	switch strings.TrimSpace(key) {
	case "knowledge_hobby:foundation_build:audit_course_practice_path":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "degree credit", "certificate ad", "private tutoring", "admission transfer", "course signup", "recruitment")
		} else {
			avoidTerms = append(avoidTerms, "학점은행제", "자격증", "과외", "편입", "취득방법", "수강신청", "모집", "해커스", "토익", "학원")
		}
	case "knowledge_hobby:habit_lifestyle:public_lifelong_online_course":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "enrollment notice", "student recruitment", "benefit coupon", "certificate print")
		} else {
			avoidTerms = append(avoidTerms, "수강신청", "수강생 모집", "모집", "가입 방법", "혜택", "쿠폰", "수료증 출력")
		}
	case "knowledge_hobby:participation_service:blended_lifelong_participation":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "enrollment only", "student recruitment", "tuition payment", "venue notice")
		} else {
			avoidTerms = append(avoidTerms, "수강신청", "수강생 모집", "접수", "장소 확인", "수강료", "모집")
		}
	case "maker_technical_hobby:certification_assessment:equipment_operation_practical_certification":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "exam schedule only", "eligibility only", "textbook sale", "academy ad")
		} else {
			avoidTerms = append(avoidTerms, "일정 안내", "응시 조건", "이론만", "교재 판매", "학원 광고")
		}
	case "knowledge_hobby:certification_assessment:environmental_engineering_certification":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "eligibility only", "degree credit", "career outlook only", "course ad", "textbook sale", "review only")
		} else {
			avoidTerms = append(avoidTerms, "응시자격", "학점은행제", "취업 방향", "인강 찾기", "교재 판매", "후기만")
		}
	case "digital_creation:artifact_creation:creative_coding_visual_project":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "coding academy", "private tutoring", "Python class ad", "tool definition only")
		} else {
			avoidTerms = append(avoidTerms, "코딩학원", "과외", "파이썬 과외", "어원", "학원 광고")
		}
	case "craft_making:artifact_creation:leathercraft_dimensional_bag":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "DIY kit sale", "finished product", "workshop ad", "shopping")
		} else {
			avoidTerms = append(avoidTerms, "DIY키트", "키트", "완제품", "공방 수업", "판매")
		}
	case "digital_creation:artifact_creation:claude_code_agentic_workflow":
		if normalizeLearningLanguage(language) == "en" {
			avoidTerms = append(avoidTerms, "benchmark only", "news only", "tool review only", "stock analysis")
		} else {
			avoidTerms = append(avoidTerms, "성능 분석만", "뉴스", "도구 리뷰만", "주가")
		}
	}
	return dedupeRecommendationTerms(avoidTerms)
}

func curriculumPatternSearchContentTypes(domain, goalType, language string) []string {
	return curriculumPatternSearchContentTypesForKey("", domain, goalType, language)
}

func curriculumPatternSearchContentTypesForKey(key, domain, goalType, language string) []string {
	switch strings.TrimSpace(key) {
	case "knowledge_hobby:habit_lifestyle:houseplant_repotting_care",
		"knowledge_hobby:habit_lifestyle:dog_basic_training_routine",
		"writing_storytelling:habit_lifestyle:daily_writing_habit":
		return []string{"article", "video"}
	}
	switch strings.TrimSpace(goalType) {
	case "foundation_build", "professional_transition", "presentation_publish":
		return []string{"article", "video"}
	}
	if strings.TrimSpace(domain) == "knowledge_hobby" && strings.TrimSpace(goalType) == "concept_mastery" {
		return []string{"article", "video"}
	}
	return []string{"video"}
}

func curriculumPatternSearchNiceTerms(language string) []string {
	if normalizeLearningLanguage(language) == "en" {
		return []string{"beginner", "tutorial", "practice", "step by step"}
	}
	return []string{"기초", "입문", "튜토리얼", "실습", "따라하기"}
}

func curriculumPatternUsesPhysicalMaterials(domain, goalType string) bool {
	switch strings.TrimSpace(domain) {
	case "craft_making", "cooking_baking", "maker_technical_hobby", "visual_art":
		return true
	}
	return false
}

func curriculumPatternSearchStageTerms(role, domain, goalType, language string) []string {
	english := normalizeLearningLanguage(language) == "en"
	switch normalizeSearchSpecStageRole(role) {
	case "setup_intro":
		if english {
			if curriculumPatternUsesPhysicalMaterials(domain, goalType) {
				return []string{"tools", "materials", "basics", "first setup"}
			}
			return []string{"getting started", "beginner", "learning plan", "first steps"}
		}
		if curriculumPatternUsesPhysicalMaterials(domain, goalType) {
			return []string{"도구", "재료", "입문", "기초"}
		}
		return []string{"시작 방법", "학습 계획", "입문", "기초"}
	case "first_output":
		if english {
			return []string{"first result", "simple practice", "beginner output"}
		}
		return []string{"첫 결과물", "간단한 실습", "기초 연습"}
	case "core_pattern":
		if english {
			return []string{"core technique", "pattern practice", "hands on"}
		}
		return []string{"핵심 기법", "반복 연습", "실습"}
	case "integration_practice":
		if english {
			return []string{"combine steps", "guided practice", "workflow"}
		}
		return []string{"단계 연결", "통합 연습", "작업 흐름"}
	case "artifact_finish":
		if english {
			return []string{"finish project", "final touches", "review"}
		}
		return []string{"완성", "마감", "점검"}
	case "performance_prep":
		if english {
			return []string{"rehearsal", "final practice", "performance ready"}
		}
		return []string{"리허설", "최종 연습", "실전 준비"}
	default:
		if english {
			return []string{"tutorial", "practice"}
		}
		return []string{"튜토리얼", "실습"}
	}
}

func buildSpecificCurriculumPatternRecommendationSearchSpecTemplate(key, language string) (CurriculumPatternRecommendationSearchSpecTemplate, bool) {
	switch strings.TrimSpace(key) {
	case "visual_art:artifact_creation:watercolor_beginner_landscape":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"watercolor landscape", "watercolor postcard", "beginner watercolor"}, []string{"wash technique", "gradient", "water control", "sky", "trees", "small landscape"}, []string{"supply sale", "gallery sale", "asmr", "advanced portrait only", "oil pastel", "makeup", "smoky makeup", "kit ad"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"beginner watercolor", "paper", "brush"}, []string{"palette", "water control", "materials"}),
				patternStageTemplate("en", "first_output", []string{"watercolor wash", "gradient", "beginner"}, []string{"wet on wet", "color mixing", "simple practice"}),
				patternStageTemplate("en", "core_pattern", []string{"watercolor landscape", "sky", "trees"}, []string{"layering", "drying", "composition"}),
				patternStageTemplate("en", "artifact_finish", []string{"watercolor postcard", "small landscape", "finish"}, []string{"final touches", "review", "scan"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"수채화", "풍경 엽서", "초보 수채화"}, []string{"번짐", "그라데이션", "물 조절", "하늘", "나무", "작은 풍경"}, []string{"재료 판매", "작품 판매", "asmr", "고급 인물화만", "오일파스텔", "메이크업", "스모키", "키트 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"초보 수채화", "종이", "붓"}, []string{"팔레트", "물 조절", "재료"}),
			{PrimaryQuery: "초보 수채화 풍경 엽서 물조절 번짐 워시 과정", Intent: "tutorial", MustInclude: []string{"수채화", "물조절", "번짐"}, NiceToHave: []string{"풍경 엽서", "워시", "붓", "그라데이션"}, Avoid: []string{"재료 판매", "작품 판매", "asmr", "고급 인물화만", "오일파스텔", "메이크업", "스모키", "키트 광고"}, ContentTypes: []string{"article", "video"}, Language: "ko", StageRole: "first_output", Source: SearchSpecSourcePatternTemplate},
			patternStageTemplate("ko", "core_pattern", []string{"수채화 풍경", "하늘", "나무"}, []string{"레이어", "건조", "구도"}),
			patternStageTemplate("ko", "artifact_finish", []string{"수채화 엽서", "작은 풍경", "완성"}, []string{"마무리", "점검", "스캔"}),
		}), true
	case "body_movement:habit_lifestyle:running_5k_beginner_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"beginner running", "5K", "running plan"}, []string{"walk run", "running form", "breathing", "interval", "pace", "recovery"}, []string{"weight loss ad", "calorie burn", "treadmill only", "shoe review only", "marathon vlog", "extreme challenge"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"beginner running", "walk run", "safe start"}, []string{"warm up", "shoes", "injury prevention"}),
				patternStageTemplate("en", "core_pattern", []string{"running form", "breathing", "pace"}, []string{"cadence", "easy run", "posture"}),
				patternStageTemplate("en", "integration_practice", []string{"interval running", "5K training", "beginner"}, []string{"walk breaks", "distance", "recovery"}),
				patternStageTemplate("en", "routine_plan", []string{"5K running plan", "weekly running", "routine"}, []string{"rest day", "tracking", "progression"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"초보 러닝", "5km", "러닝 루틴"}, []string{"걷기 달리기", "달리기 자세", "호흡", "인터벌", "페이스", "회복"}, []string{"다이어트 광고", "다이어트", "살 빼", "칼로리", "런닝머신", "런닝머신만", "체지방", "러닝화 리뷰만", "마라톤 브이로그", "챌린지", "무리한 챌린지"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"초보 러닝", "걷기 달리기", "안전 시작"}, []string{"워밍업", "러닝화", "부상 예방"}),
			patternStageTemplate("ko", "core_pattern", []string{"달리기 자세", "호흡", "페이스"}, []string{"케이던스", "천천히 달리기", "자세"}),
			{PrimaryQuery: "5km 완주 초보 러닝 훈련 계획 걷기 달리기", Intent: "tutorial", MustInclude: []string{"5km", "초보러너", "완주"}, NiceToHave: []string{"훈련 계획", "걷기 달리기", "페이스", "루틴", "인터벌"}, Avoid: []string{"다이어트 광고", "다이어트", "살 빼", "칼로리", "런닝머신", "런닝머신만", "체지방", "러닝화 리뷰만", "마라톤 브이로그", "챌린지", "무리한 챌린지"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "integration_practice", Source: SearchSpecSourcePatternTemplate},
			patternStageTemplate("ko", "routine_plan", []string{"5km 러닝 계획", "주간 러닝", "루틴"}, []string{"휴식일", "기록", "거리 증가"}),
		}), true
	case "language_communication:performance_execution:daily_english_conversation_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"daily English", "English conversation", "speaking practice"}, []string{"everyday phrases", "shadowing", "role play", "self introduction", "daily routine", "listening"}, []string{"toeic", "exam only", "grammar only", "study abroad ad", "tutoring ad", "school exam", "private lesson"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"daily English", "basic phrases", "beginner"}, []string{"routine", "listening", "pronunciation"}),
				patternStageTemplate("en", "core_pattern", []string{"English conversation", "everyday phrases", "shadowing"}, []string{"question answer", "short sentences", "pronunciation"}),
				patternStageTemplate("en", "integration_practice", []string{"speaking practice", "role play", "daily conversation"}, []string{"morning routine", "work", "shopping"}),
				patternStageTemplate("en", "performance_prep", []string{"daily English routine", "self introduction", "conversation"}, []string{"recording", "feedback", "repeat practice"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"생활 영어", "말하기 연습", "쉐도잉 루틴"}, []string{"일상 영어", "영어 회화", "생활 표현", "역할극", "자기소개", "하루 루틴", "듣기"}, []string{"토익", "오픽", "영어학원", "성인영어학원", "직장인", "과외", "내신", "할인", "쿠폰", "멤버십"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"생활 영어", "기초 표현", "초보"}, []string{"루틴", "듣기", "발음"}),
			patternStageTemplate("ko", "core_pattern", []string{"생활 표현", "쉐도잉", "말하기 연습"}, []string{"일상 영어", "영어 회화", "질문 답변", "짧은 문장", "발음"}),
			patternStageTemplate("ko", "integration_practice", []string{"말하기 연습", "역할극", "일상 대화"}, []string{"아침 루틴", "쇼핑", "짧은 문장"}),
			patternStageTemplate("ko", "performance_prep", []string{"영어 회화 루틴", "자기소개", "대화"}, []string{"녹음", "피드백", "반복 연습"}),
		}), true
	case "maker_technical_hobby:artifact_creation:arduino_sensor_project":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"Arduino", "sensor project", "breadboard"}, []string{"LED", "wiring", "serial monitor", "input output", "servo", "beginner project"}, []string{"kit sale", "shopping", "advanced electronics only", "course sale", "kit review", "special price", "sale"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"Arduino", "breadboard", "beginner"}, []string{"wiring", "USB", "IDE"}),
				patternStageTemplate("en", "first_output", []string{"Arduino LED", "blink", "wiring"}, []string{"resistor", "digital output", "upload"}),
				patternStageTemplate("en", "core_pattern", []string{"Arduino sensor", "serial monitor", "input"}, []string{"ultrasonic", "temperature", "analog read"}),
				patternStageTemplate("en", "artifact_finish", []string{"Arduino sensor project", "servo", "final check"}, []string{"troubleshooting", "case", "demo"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"아두이노", "센서 프로젝트", "브레드보드"}, []string{"LED", "배선", "시리얼 모니터", "입출력", "서보모터", "초보 프로젝트"}, []string{"키트 판매", "키트", "쇼핑몰", "고급 전자이론만", "강의 판매", "세트", "키트 후기", "특가", "알리", "판매"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"아두이노", "브레드보드", "초보"}, []string{"배선", "USB", "IDE"}),
			patternStageTemplate("ko", "first_output", []string{"아두이노 LED", "깜빡이기", "배선"}, []string{"저항", "디지털 출력", "업로드"}),
			patternStageTemplate("ko", "core_pattern", []string{"아두이노 센서", "시리얼 모니터", "입력"}, []string{"초음파", "온습도", "아날로그 읽기"}),
			patternStageTemplate("ko", "artifact_finish", []string{"아두이노 센서 프로젝트", "서보모터", "최종 점검"}, []string{"문제 해결", "케이스", "시연"}),
		}), true
	case "cooking_baking:artifact_creation:korean_home_cooking_basics":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"Korean home cooking", "banchan", "doenjang jjigae"}, []string{"knife prep", "seasoning", "egg roll", "kimchi jjigae", "one meal", "beginner recipe"}, []string{"restaurant review", "mukbang", "diet ad", "certificate exam", "restaurant", "local diner"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"Korean home cooking", "kitchen basics", "food safety"}, []string{"ingredients", "knife", "pan"}),
				patternStageTemplate("en", "structure_plan", []string{"home meal plan", "ingredients prep", "cooking order"}, []string{"rice", "soup", "side dish"}),
				patternStageTemplate("en", "core_pattern", []string{"doenjang jjigae", "banchan", "seasoning"}, []string{"egg roll", "namul", "taste adjustment"}),
				patternStageTemplate("en", "artifact_finish", []string{"Korean home meal", "plating", "taste check"}, []string{"cleanup", "leftovers", "review"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"집밥", "한식 기초", "밑반찬"}, []string{"재료 손질", "간 맞추기", "된장찌개", "계란말이", "한 끼", "초보 레시피"}, []string{"맛집 후기", "먹방", "다이어트 광고", "자격증 시험", "맛집", "식당", "백반 맛집", "행궁동", "의정부"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"집밥", "주방 기초", "식품 위생"}, []string{"재료", "칼질", "팬"}),
			patternStageTemplate("ko", "structure_plan", []string{"한 끼 구성", "재료 손질", "조리 순서"}, []string{"밥", "국", "반찬"}),
			patternStageTemplate("ko", "core_pattern", []string{"된장찌개", "밑반찬", "간 맞추기"}, []string{"계란말이", "나물", "맛 조절"}),
			patternStageTemplate("ko", "artifact_finish", []string{"집밥 한 끼", "플레이팅", "맛 점검"}, []string{"정리", "남은 반찬", "복기"}),
		}), true
	case "body_movement:habit_lifestyle:home_strength_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"home workout", "strength routine", "bodyweight workout"}, []string{"beginner", "safe form", "squat", "push-up", "core", "weekly plan"}, []string{"weight loss ad", "supplement", "extreme challenge", "gym vlog"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"home workout", "beginner", "safe form"}, []string{"warm up", "equipment", "bodyweight"}),
				patternStageTemplate("en", "core_pattern", []string{"bodyweight workout", "squat", "push-up"}, []string{"core", "dumbbell", "form check"}),
				patternStageTemplate("en", "integration_practice", []string{"strength routine", "full body", "beginner"}, []string{"sets", "reps", "rest"}),
				patternStageTemplate("en", "routine_plan", []string{"weekly workout plan", "home strength", "routine"}, []string{"tracking", "recovery", "progression"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"홈트", "근력 루틴", "맨몸운동"}, []string{"초보", "안전 자세", "스쿼트", "푸쉬업", "코어", "주간 계획"}, []string{"다이어트 광고", "보충제", "고강도 챌린지", "헬스장 브이로그"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"홈트", "초보", "안전 자세"}, []string{"워밍업", "준비물", "맨몸"}),
			patternStageTemplate("ko", "core_pattern", []string{"맨몸운동", "스쿼트", "푸쉬업"}, []string{"코어", "덤벨", "자세 점검"}),
			patternStageTemplate("ko", "integration_practice", []string{"근력 루틴", "전신 운동", "초보"}, []string{"세트", "횟수", "휴식"}),
			patternStageTemplate("ko", "routine_plan", []string{"주간 운동 계획", "홈 근력", "루틴"}, []string{"기록", "회복", "강도 조절"}),
		}), true
	case "instrument_performance:performance_execution:vocal_song_practice":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"vocal practice", "singing practice", "cover song"}, []string{"breath support", "pitch", "rhythm", "warm-up", "phrase practice", "full song"}, []string{"audition ad", "reaction video", "performance only", "vocal coach sale"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"vocal warm-up", "breath support", "beginner"}, []string{"posture", "range", "safe voice"}),
				patternStageTemplate("en", "first_output", []string{"singing practice", "first verse", "pitch"}, []string{"slow practice", "lyrics", "melody"}),
				patternStageTemplate("en", "core_pattern", []string{"vocal practice", "pitch", "rhythm"}, []string{"breath", "tone", "phrase"}),
				patternStageTemplate("en", "performance_prep", []string{"cover song", "full song", "rehearsal"}, []string{"recording check", "expression", "mistake fix"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"보컬 연습", "노래 연습", "커버곡"}, []string{"호흡", "음정", "박자", "발성", "구간 연습", "한 곡"}, []string{"오디션 광고", "리액션 영상", "공연 영상만", "보컬 학원 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"보컬 워밍업", "호흡", "초보"}, []string{"자세", "음역", "목 보호"}),
			patternStageTemplate("ko", "first_output", []string{"노래 연습", "첫 소절", "음정"}, []string{"느린 연습", "가사", "멜로디"}),
			patternStageTemplate("ko", "core_pattern", []string{"보컬 연습", "음정", "박자"}, []string{"호흡", "톤", "프레이즈"}),
			patternStageTemplate("ko", "performance_prep", []string{"커버곡", "한 곡", "리허설"}, []string{"녹음 점검", "표현", "실수 보완"}),
		}), true
	case "digital_creation:artifact_creation:video_editing_shortform_project":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"video editing", "shorts", "capcut"}, []string{"cut editing", "subtitles", "music", "vertical video", "youtube shorts", "export", "beginner project"}, []string{"growth hack", "views guarantee", "reaction video", "course sale"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"video editing", "shortform", "source clips"}, []string{"project setup", "vertical format", "import"}),
				patternStageTemplate("en", "first_output", []string{"cut editing", "short video", "beginner"}, []string{"trim", "timeline", "rough cut"}),
				patternStageTemplate("en", "core_pattern", []string{"subtitles", "music", "shorts editing"}, []string{"caption", "transition", "sound level"}),
				patternStageTemplate("en", "artifact_finish", []string{"export", "vertical video", "shorts"}, []string{"final check", "thumbnail", "format"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"영상 편집", "영상편집", "쇼츠"}, []string{"캡컷", "컷편집", "자막", "음악", "세로 영상", "내보내기", "초급 프로젝트"}, []string{"조회수 비법", "조회수 보장", "리액션 영상", "강의 판매"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"영상 편집", "영상편집", "숏폼"}, []string{"소스 정리", "프로젝트 설정", "세로 화면", "불러오기"}),
			patternStageTemplate("ko", "first_output", []string{"컷편집", "짧은 영상", "초보"}, []string{"자르기", "타임라인", "러프컷"}),
			patternStageTemplate("ko", "core_pattern", []string{"자막", "영상편집", "쇼츠 편집"}, []string{"음악", "캡션", "전환", "음량"}),
			patternStageTemplate("ko", "artifact_finish", []string{"내보내기", "썸네일", "쇼츠"}, []string{"영상편집", "캡컷", "최종 점검", "세로 영상", "형식"}),
		}), true
	case "digital_creation:technical_skill:ai_tool_productivity_workflow":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"AI productivity", "ChatGPT workflow", "document workflow"}, []string{"document summary", "spreadsheet automation", "review", "privacy", "repeatable task"}, []string{"prompt list only", "affiliate tool review", "make money fast", "course sale", "Flow AI product demo"}, []LessonRecommendationSearchSpec{
				{PrimaryQuery: "ChatGPT productivity document summary workflow tutorial", Intent: "tutorial", MustInclude: []string{"ChatGPT", "document summary", "workflow"}, NiceToHave: []string{"task selection", "privacy", "input data"}, Avoid: []string{"prompt list only", "affiliate tool review", "make money fast", "course sale", "Flow AI product demo"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "setup_intro", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "ChatGPT prompt document summary review workflow tutorial", Intent: "tutorial", MustInclude: []string{"ChatGPT", "prompt", "document summary"}, NiceToHave: []string{"draft", "review", "rewrite"}, Avoid: []string{"prompt list only", "affiliate tool review", "make money fast", "course sale", "viral hacks"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "core_pattern", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "AI spreadsheet automation workflow checklist tutorial", Intent: "tutorial", MustInclude: []string{"AI tools", "spreadsheet automation", "workflow"}, NiceToHave: []string{"checklist", "quality check", "handoff"}, Avoid: []string{"affiliate tool review", "make money fast", "course sale"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "integration_practice", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "AI productivity repeatable workflow review checklist", Intent: "tutorial", MustInclude: []string{"AI productivity", "repeatable workflow", "review"}, NiceToHave: []string{"template", "guardrails", "next task"}, Avoid: []string{"affiliate tool review", "make money fast", "course sale"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "routine_plan", Source: SearchSpecSourcePatternTemplate},
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"ChatGPT 업무 자동화", "문서 요약", "프롬프트"}, []string{"AI 생산성", "검토", "개인정보", "반복 업무", "구글시트 자동화"}, []string{"프롬프트 모음만", "제휴 도구 리뷰", "돈버는 법", "강의 판매", "플로우 AI"}, []LessonRecommendationSearchSpec{
			{PrimaryQuery: "챗GPT 문서 요약 업무 자동화 실습", Intent: "tutorial", MustInclude: []string{"챗GPT", "문서 요약", "업무 자동화"}, NiceToHave: []string{"사용 사례", "개인정보", "입력 자료"}, Avoid: []string{"프롬프트 모음만", "제휴 도구 리뷰", "돈버는 법", "강의 판매", "플로우 AI", "자동수익화", "무료 특강"}, ContentTypes: []string{"article", "video"}, Language: "ko", StageRole: "setup_intro", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "챗GPT 프롬프트 문서 요약 검토 워크플로우", Intent: "tutorial", MustInclude: []string{"챗GPT", "프롬프트", "문서 요약"}, NiceToHave: []string{"초안", "검토", "재작성"}, Avoid: []string{"프롬프트 모음만", "제휴 도구 리뷰", "돈버는 법", "강의 판매", "미친 사용법", "꿀기능", "상위 0.1%"}, ContentTypes: []string{"article", "video"}, Language: "ko", StageRole: "core_pattern", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "AI 도구 구글시트 자동화 업무 워크플로우", Intent: "tutorial", MustInclude: []string{"AI 도구", "구글시트 자동화", "업무"}, NiceToHave: []string{"체크리스트", "품질 점검", "인수인계"}, Avoid: []string{"제휴 도구 리뷰", "돈버는 법", "강의 판매"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "integration_practice", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "AI 생산성 반복 업무 검토 체크리스트", Intent: "tutorial", MustInclude: []string{"AI 생산성", "반복 업무", "검토"}, NiceToHave: []string{"템플릿", "주의 기준", "다음 업무"}, Avoid: []string{"제휴 도구 리뷰", "돈버는 법", "강의 판매"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "routine_plan", Source: SearchSpecSourcePatternTemplate},
		}), true
	case "craft_making:artifact_creation:candle_soap_resin_project":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"candle making", "soap making", "resin art"}, []string{"mold", "fragrance oil", "wax", "color", "curing", "safety"}, []string{"supply sale", "bulk business", "asmr", "kit only"}, []LessonRecommendationSearchSpec{
				{PrimaryQuery: "beginner candle making materials safety tutorial", Intent: "tutorial", MustInclude: []string{"candle making", "beginner", "materials"}, NiceToHave: []string{"safety", "mold", "ventilation"}, Avoid: []string{"supply sale", "bulk business", "asmr", "kit only"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "setup_intro", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "beginner resin keychain making tutorial", Intent: "tutorial", MustInclude: []string{"resin keychain", "beginner", "resin art"}, NiceToHave: []string{"pouring", "color", "small project"}, Avoid: []string{"supply sale", "bulk business", "asmr", "kit only"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "first_output", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "soy candle wax fragrance oil mold temperature tutorial", Intent: "tutorial", MustInclude: []string{"soy candle", "wax", "mold"}, NiceToHave: []string{"fragrance oil", "temperature", "curing"}, Avoid: []string{"supply sale", "bulk business", "asmr", "kit only"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "core_pattern", Source: SearchSpecSourcePatternTemplate},
				{PrimaryQuery: "handmade soap demold packaging finish tutorial", Intent: "tutorial", MustInclude: []string{"handmade soap", "demold", "packaging"}, NiceToHave: []string{"surface check", "finish", "small batch"}, Avoid: []string{"supply sale", "bulk business", "asmr", "kit only"}, ContentTypes: []string{"video"}, Language: "en", StageRole: "artifact_finish", Source: SearchSpecSourcePatternTemplate},
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"캔들 만들기", "수제비누", "레진아트"}, []string{"몰드", "향료", "왁스", "색 조합", "경화", "안전"}, []string{"재료 판매", "대량 창업", "asmr", "키트만"}, []LessonRecommendationSearchSpec{
			{PrimaryQuery: "캔들 만들기 초보 재료 안전 튜토리얼", Intent: "tutorial", MustInclude: []string{"캔들 만들기", "초보", "재료"}, NiceToHave: []string{"안전", "몰드", "환기"}, Avoid: []string{"재료 판매", "대량 창업", "asmr", "키트만"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "setup_intro", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "초보 레진 키링 만들기 튜토리얼", Intent: "tutorial", MustInclude: []string{"레진 키링", "초보", "레진아트"}, NiceToHave: []string{"붓기", "색 조합", "작은 소품"}, Avoid: []string{"재료 판매", "대량 창업", "asmr", "키트만"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "first_output", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "캔들 만들기 왁스 향료 몰드 온도 튜토리얼", Intent: "tutorial", MustInclude: []string{"캔들 만들기", "왁스", "몰드"}, NiceToHave: []string{"향료", "온도", "경화"}, Avoid: []string{"재료 판매", "대량 창업", "asmr", "키트만"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "core_pattern", Source: SearchSpecSourcePatternTemplate},
			{PrimaryQuery: "수제비누 탈형 포장 마감 튜토리얼", Intent: "tutorial", MustInclude: []string{"수제비누", "탈형", "포장"}, NiceToHave: []string{"표면 점검", "마감", "소량 제작"}, Avoid: []string{"재료 판매", "대량 창업", "asmr", "키트만"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "artifact_finish", Source: SearchSpecSourcePatternTemplate},
		}), true
	case "craft_making:artifact_creation:leathercraft_wallet_project":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"leathercraft", "wallet", "card wallet"},
				[]string{"beginner", "tutorial", "pattern", "cutting", "saddle stitch", "edge finishing"},
				[]string{"shopping", "sale", "kit only", "asmr"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"leathercraft", "wallet", "tools"}, []string{"beginner", "materials", "pattern"}),
					patternStageTemplate("en", "structure_plan", []string{"leathercraft", "wallet", "pattern"}, []string{"card slot", "outer leather", "inner leather"}),
					patternStageTemplate("en", "core_pattern", []string{"leathercraft", "saddle stitch", "cutting"}, []string{"punching", "glue", "card wallet"}),
					patternStageTemplate("en", "artifact_finish", []string{"leathercraft", "wallet", "edge finishing"}, []string{"burnish", "edge paint", "quality check"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"가죽공예", "지갑", "카드지갑"},
			[]string{"기초", "입문", "패턴", "재단", "새들스티치", "엣지마감"},
			[]string{"DIY키트", "키트", "완제품", "재료 판매", "공구 쇼핑몰", "광고", "asmr"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"가죽공예", "지갑", "도구"}, []string{"입문", "재료", "패턴"}),
				patternStageTemplate("ko", "structure_plan", []string{"가죽공예", "지갑", "패턴"}, []string{"카드 슬롯", "외피", "내피"}),
				{PrimaryQuery: "가죽공예 카드지갑 새들스티치 패턴 재단 강좌", Intent: "tutorial", MustInclude: []string{"가죽공예", "카드지갑", "새들스티치"}, NiceToHave: []string{"패턴", "재단", "바느질", "타공"}, Avoid: []string{"DIY키트", "키트", "완제품", "재료 판매", "공구 쇼핑몰", "광고", "asmr"}, ContentTypes: []string{"video"}, Language: "ko", StageRole: "core_pattern", Source: SearchSpecSourcePatternTemplate},
				patternStageTemplate("ko", "artifact_finish", []string{"가죽공예", "지갑", "엣지마감"}, []string{"코바", "엣지코트", "품질 점검"}),
			}), true
	case "language_communication:performance_execution:travel_conversation":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"travel english", "travel conversation", "airport"},
				[]string{"hotel check-in", "restaurant ordering", "directions", "role play", "phrases"},
				[]string{"toeic", "exam", "grammar only", "study abroad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"travel english", "basic phrases", "airport"}, []string{"hotel", "restaurant", "beginner"}),
					patternStageTemplate("en", "core_pattern", []string{"travel english", "airport", "hotel"}, []string{"check-in", "immigration", "ordering"}),
					patternStageTemplate("en", "integration_practice", []string{"travel conversation", "role play", "directions"}, []string{"restaurant", "transport", "shopping"}),
					patternStageTemplate("en", "performance_prep", []string{"travel english", "travel-day simulation", "conversation"}, []string{"listening", "speaking", "polite requests"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"여행영어", "여행 회화", "공항"},
			[]string{"호텔 체크인", "식당 주문", "길 묻기", "역할극", "필수 표현"},
			[]string{"토익", "시험", "문법 강의", "유학"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"여행영어", "기초 표현", "공항"}, []string{"호텔", "식당", "입문"}),
				patternStageTemplate("ko", "core_pattern", []string{"여행영어", "공항", "호텔"}, []string{"체크인", "입국심사", "주문"}),
				patternStageTemplate("ko", "integration_practice", []string{"여행 회화", "역할극", "길 묻기"}, []string{"식당", "교통", "쇼핑"}),
				patternStageTemplate("ko", "performance_prep", []string{"여행영어", "여행 상황", "회화"}, []string{"듣기", "말하기", "정중한 요청"}),
			}), true
	case "instrument_performance:performance_execution:song_completion":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"instrument", "song", "play through"},
				[]string{"section practice", "chord transition", "rhythm", "rehearsal", "full run-through"},
				[]string{"instrument shopping", "gear review", "music theory only", "performance video only"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"instrument", "target song", "beginner arrangement"}, []string{"tuning", "basic posture", "song level"}),
					patternStageTemplate("en", "first_output", []string{"instrument", "song section", "first phrase"}, []string{"intro", "verse", "slow practice"}),
					patternStageTemplate("en", "integration_practice", []string{"song practice", "section transitions", "rhythm"}, []string{"loop practice", "tempo", "play along"}),
					patternStageTemplate("en", "performance_prep", []string{"song run-through", "rehearsal", "full song"}, []string{"recording check", "mistake fixes", "performance ready"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"악기", "한 곡", "완주"},
			[]string{"구간 연습", "코드 전환", "리듬", "리허설", "전체 연주"},
			[]string{"악기 쇼핑", "장비 리뷰", "음악이론만", "연주 영상만"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"악기", "목표곡", "초보 편곡"}, []string{"튜닝", "기본 자세", "난이도"}),
				patternStageTemplate("ko", "first_output", []string{"악기", "곡 구간", "첫 소절"}, []string{"인트로", "벌스", "느린 연습"}),
				patternStageTemplate("ko", "integration_practice", []string{"곡 연습", "구간 연결", "리듬"}, []string{"반복 연습", "템포", "반주 맞추기"}),
				patternStageTemplate("ko", "performance_prep", []string{"한 곡 완주", "리허설", "전체 연주"}, []string{"녹음 점검", "실수 보완", "실전 준비"}),
			}), true
	case "digital_creation:artifact_creation:vibe_coding_mvp_app":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"web app", "MVP", "full-stack"},
				[]string{"AI coding", "CRUD", "database", "authentication", "deploy", "responsive UI"},
				[]string{"prompt list only", "framework comparison", "course sale", "no-code ad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"web app", "MVP scope", "user flow"}, []string{"AI coding tool", "requirements", "screens"}),
					patternStageTemplate("en", "structure_plan", []string{"web app", "data model", "screen plan"}, []string{"API routes", "database", "auth boundary"}),
					patternStageTemplate("en", "first_output", []string{"web app", "first working page", "CRUD"}, []string{"form", "list view", "local state"}),
					patternStageTemplate("en", "core_pattern", []string{"web app", "database", "authentication"}, []string{"validation", "backend", "error fix"}),
					patternStageTemplate("en", "artifact_finish", []string{"web app", "deploy", "smoke test"}, []string{"responsive check", "bug fix", "feedback checklist"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"웹앱", "MVP", "풀스택"},
			[]string{"AI 코딩", "CRUD", "데이터베이스", "인증", "배포", "반응형 UI"},
			[]string{"프롬프트 모음", "프레임워크 비교", "강의 판매", "노코드 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"웹앱", "MVP 범위", "사용자 흐름"}, []string{"AI 코딩 도구", "요구사항", "화면 설계"}),
				patternStageTemplate("ko", "structure_plan", []string{"웹앱", "데이터 모델", "화면 설계"}, []string{"API 라우트", "데이터베이스", "인증 경계"}),
				patternStageTemplate("ko", "first_output", []string{"웹앱", "작동 페이지", "CRUD"}, []string{"폼", "목록 화면", "상태 관리"}),
				patternStageTemplate("ko", "core_pattern", []string{"웹앱", "데이터베이스", "인증"}, []string{"검증", "백엔드", "오류 수정"}),
				patternStageTemplate("ko", "artifact_finish", []string{"웹앱", "배포", "스모크 테스트"}, []string{"반응형 점검", "버그 수정", "피드백 체크리스트"}),
			}), true
	case "digital_creation:artifact_creation:web_dev_project_foundation":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"},
				[]string{"web development", "HTML CSS JavaScript", "beginner project"},
				[]string{"responsive layout", "DOM interaction", "portfolio project", "forms", "accessibility", "deploy"},
				[]string{"AI coding prompt", "no-code tool", "framework comparison", "course sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"web development", "HTML CSS JavaScript", "beginner setup"}, []string{"VS Code", "browser devtools", "file structure"}),
					patternStageTemplate("en", "first_output", []string{"HTML CSS", "first web page", "beginner project"}, []string{"layout", "text image", "style basics"}),
					patternStageTemplate("en", "core_pattern", []string{"JavaScript", "DOM interaction", "web page"}, []string{"button event", "form input", "state"}),
					patternStageTemplate("en", "integration_practice", []string{"responsive website", "project workflow", "HTML CSS JavaScript"}, []string{"mobile layout", "accessibility", "debugging"}),
					patternStageTemplate("en", "artifact_finish", []string{"portfolio project", "deploy", "responsive check"}, []string{"GitHub Pages", "final polish", "project review"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"},
			[]string{"웹개발", "HTML CSS JavaScript", "초급 프로젝트"},
			[]string{"반응형 레이아웃", "DOM 조작", "포트폴리오 프로젝트", "폼", "접근성", "배포"},
			[]string{"AI 코딩 프롬프트", "노코드 도구", "프레임워크 비교", "강의 판매", "학원", "코딩학원", "취업 보장", "국비 모집", "상담 신청"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"웹개발", "HTML CSS JavaScript", "초급 환경설정"}, []string{"VS Code", "브라우저 개발자도구", "파일 구조"}),
				patternStageTemplate("ko", "first_output", []string{"HTML CSS", "첫 웹페이지", "초급 프로젝트"}, []string{"레이아웃", "텍스트 이미지", "스타일 기초"}),
				patternStageTemplate("ko", "core_pattern", []string{"JavaScript", "DOM 조작", "웹페이지"}, []string{"버튼 이벤트", "폼 입력", "상태"}),
				patternStageTemplate("ko", "integration_practice", []string{"반응형 웹사이트", "프로젝트 흐름", "HTML CSS JavaScript"}, []string{"모바일 레이아웃", "접근성", "디버깅"}),
				patternStageTemplate("ko", "artifact_finish", []string{"포트폴리오 프로젝트", "배포", "반응형 점검"}, []string{"GitHub Pages", "최종 다듬기", "프로젝트 리뷰"}),
			}), true
	case "visual_art:artifact_creation:structured_drawing_foundation":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"drawing foundation", "urban sketch", "watercolor sketch"},
				[]string{"contour drawing", "shading", "perspective", "composition", "wash", "postcard"},
				[]string{"art sale", "gallery opening", "entrance exam", "paid class ad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"drawing tools", "sketchbook", "beginner drawing"}, []string{"pencil", "watercolor set", "paper"}),
					patternStageTemplate("en", "first_output", []string{"contour drawing", "simple sketch", "beginner practice"}, []string{"line drawing", "still life", "small scene"}),
					patternStageTemplate("en", "core_pattern", []string{"shading", "perspective", "composition"}, []string{"value study", "depth", "thumbnail sketch"}),
					patternStageTemplate("en", "integration_practice", []string{"urban sketch", "watercolor wash", "scene sketch"}, []string{"sky", "trees", "street scene"}),
					patternStageTemplate("en", "artifact_finish", []string{"finished sketch", "watercolor postcard", "review"}, []string{"final touches", "scan", "portfolio"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"드로잉 기초", "어반스케치", "수채화 스케치"},
			[]string{"선 드로잉", "명암", "원근", "구도", "수채화 워시", "엽서"},
			[]string{"작품 판매", "전시 모집", "입시 미술", "유료 클래스 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"드로잉 도구", "스케치북", "초보 드로잉"}, []string{"연필", "수채화 도구", "종이"}),
				patternStageTemplate("ko", "first_output", []string{"선 드로잉", "간단한 스케치", "기초 연습"}, []string{"윤곽선", "정물", "작은 풍경"}),
				patternStageTemplate("ko", "core_pattern", []string{"명암", "원근", "구도"}, []string{"톤 연습", "깊이감", "썸네일 스케치"}),
				patternStageTemplate("ko", "integration_practice", []string{"어반스케치", "수채화 워시", "풍경 스케치"}, []string{"하늘", "나무", "거리 풍경"}),
				patternStageTemplate("ko", "artifact_finish", []string{"스케치 완성", "수채화 엽서", "작품 점검"}, []string{"마무리", "스캔", "포트폴리오"}),
			}), true
	case "cooking_baking:artifact_creation:cooking_basics_foundation":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"cooking basics", "knife skills", "food safety"}, []string{"ingredient prep", "pan heat", "sauce", "one dish", "beginner recipe"}, []string{"restaurant review", "product sale", "diet ad", "certificate exam"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"cooking basics", "kitchen tools", "food safety"}, []string{"knife safety", "pan", "ingredients"}),
				patternStageTemplate("en", "structure_plan", []string{"simple recipe", "ingredient prep", "cooking order"}, []string{"mise en place", "timing", "seasoning"}),
				patternStageTemplate("en", "core_pattern", []string{"knife skills", "pan cooking", "sauce"}, []string{"saute", "boil", "taste adjustment"}),
				patternStageTemplate("en", "artifact_finish", []string{"beginner dish", "plating", "taste check"}, []string{"one meal", "review", "cleanup"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"요리 기초", "칼질", "식품 위생"}, []string{"재료 손질", "팬 조리", "소스", "한 그릇 요리", "초보 레시피"}, []string{"맛집 후기", "제품 판매", "다이어트 광고", "자격증 시험"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"요리 기초", "주방 도구", "식품 위생"}, []string{"칼 안전", "팬", "재료"}),
			patternStageTemplate("ko", "structure_plan", []string{"간단한 레시피", "재료 손질", "조리 순서"}, []string{"미장플라스", "시간 배분", "간 맞추기"}),
			patternStageTemplate("ko", "core_pattern", []string{"칼질", "팬 조리", "소스"}, []string{"볶기", "삶기", "맛 조절"}),
			patternStageTemplate("ko", "artifact_finish", []string{"초보 요리", "플레이팅", "맛 점검"}, []string{"한 끼", "복기", "정리"}),
		}), true
	case "cooking_baking:teaching_instruction:baking_class_demo":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"baking class", "recipe demo", "teaching flow"}, []string{"mise en place", "dough", "oven", "explanation", "class script"}, []string{"bakery sale", "class recruitment", "franchise", "exam only"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"baking class", "recipe selection", "tools"}, []string{"ingredients", "oven", "prep checklist"}),
				patternStageTemplate("en", "structure_plan", []string{"class script", "demo sequence", "recipe steps"}, []string{"timing", "explanation", "learner mistakes"}),
				patternStageTemplate("en", "core_pattern", []string{"baking technique", "dough", "oven temperature"}, []string{"mixing", "proofing", "texture"}),
				patternStageTemplate("en", "artifact_finish", []string{"finished bake", "quality check", "serving"}, []string{"cooling", "plating", "taste review"}),
				patternStageTemplate("en", "teaching_prep", []string{"baking tutorial", "demo rehearsal", "class delivery"}, []string{"camera", "voice", "Q&A"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"베이킹 클래스", "레시피 시연", "수업 흐름"}, []string{"재료 준비", "반죽", "오븐", "설명", "수업 대본"}, []string{"빵집 판매", "클래스 모집", "프랜차이즈", "시험만"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"베이킹 클래스", "레시피 선정", "도구"}, []string{"재료", "오븐", "준비 체크리스트"}),
			patternStageTemplate("ko", "structure_plan", []string{"수업 대본", "시연 순서", "레시피 단계"}, []string{"시간 배분", "설명", "실수 포인트"}),
			patternStageTemplate("ko", "core_pattern", []string{"베이킹 기법", "반죽", "오븐 온도"}, []string{"믹싱", "발효", "식감"}),
			patternStageTemplate("ko", "artifact_finish", []string{"베이킹 완성", "품질 점검", "담아내기"}, []string{"식힘", "플레이팅", "맛 평가"}),
			patternStageTemplate("ko", "teaching_prep", []string{"베이킹 튜토리얼", "시연 리허설", "수업 진행"}, []string{"촬영", "목소리", "질문 답변"}),
		}), true
	case "craft_making:teaching_instruction:craft_teach_youtube":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"craft tutorial", "making process", "YouTube demo"}, []string{"project steps", "camera setup", "voiceover", "editing", "thumbnail"}, []string{"craft supply sale", "channel growth hack", "paid course", "shorts only"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "first_output", []string{"craft project", "first sample", "making process"}, []string{"simple item", "materials", "steps"}),
				patternStageTemplate("en", "core_pattern", []string{"craft technique", "hands-on demo", "close-up"}, []string{"repeatable steps", "mistakes", "tips"}),
				patternStageTemplate("en", "structure_plan", []string{"tutorial script", "shot list", "YouTube lesson"}, []string{"intro", "chapter", "explanation"}),
				patternStageTemplate("en", "artifact_finish", []string{"finished craft", "final review", "thumbnail"}, []string{"before after", "description", "upload prep"}),
				patternStageTemplate("en", "teaching_prep", []string{"craft tutorial video", "recording", "editing"}, []string{"voiceover", "captions", "publish checklist"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"공예 튜토리얼", "제작 과정", "유튜브 시연"}, []string{"작업 단계", "촬영 세팅", "내레이션", "편집", "썸네일"}, []string{"공예 재료 판매", "채널 성장 비법", "유료 강의", "쇼츠만"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "first_output", []string{"공예 작품", "첫 샘플", "제작 과정"}, []string{"간단한 소품", "재료", "순서"}),
			patternStageTemplate("ko", "core_pattern", []string{"공예 기법", "손작업 시연", "클로즈업"}, []string{"반복 단계", "실수", "팁"}),
			patternStageTemplate("ko", "structure_plan", []string{"튜토리얼 대본", "촬영 컷", "유튜브 강의"}, []string{"도입", "챕터", "설명"}),
			patternStageTemplate("ko", "artifact_finish", []string{"공예 완성", "최종 점검", "썸네일"}, []string{"전후 비교", "설명란", "업로드 준비"}),
			patternStageTemplate("ko", "teaching_prep", []string{"공예 튜토리얼 영상", "촬영", "편집"}, []string{"내레이션", "자막", "게시 체크리스트"}),
		}), true
	case "writing_storytelling:presentation_publish:brunch_serial_publish":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"essay writing", "serial writing", "publishing plan"}, []string{"topic list", "draft", "revision", "platform publishing", "reader hook"}, []string{"marketing agency", "paid writing class", "monetization only", "AI article spam"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "first_output", []string{"essay draft", "first post", "writing prompt"}, []string{"personal essay", "opening", "short draft"}),
				patternStageTemplate("en", "structure_plan", []string{"serial writing", "topic plan", "publishing schedule"}, []string{"series outline", "reader", "category"}),
				patternStageTemplate("en", "integration_practice", []string{"revise essay", "writing flow", "title"}, []string{"paragraph", "voice", "feedback"}),
				patternStageTemplate("en", "publish_prep", []string{"publish essay", "platform setup", "final edit"}, []string{"thumbnail", "introduction", "release checklist"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"브런치 글쓰기", "연재 글", "발행 계획"}, []string{"주제 목록", "초안", "퇴고", "플랫폼 발행", "독자 후킹"}, []string{"마케팅 대행", "유료 글쓰기 강의", "수익화만", "AI 글 대량생성", "대필", "출판 대행", "조회수 보장"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "first_output", []string{"에세이 초안", "첫 글", "글쓰기 주제"}, []string{"개인 에세이", "도입부", "짧은 초안"}),
			patternStageTemplate("ko", "structure_plan", []string{"연재 글쓰기", "주제 기획", "발행 일정"}, []string{"시리즈 목차", "독자", "카테고리"}),
			patternStageTemplate("ko", "integration_practice", []string{"글 퇴고", "문장 흐름", "제목"}, []string{"문단", "문체", "피드백"}),
			patternStageTemplate("ko", "publish_prep", []string{"글 발행", "브런치 설정", "최종 교정"}, []string{"대표 이미지", "소개글", "공개 체크리스트"}),
		}), true
	case "language_communication:certification_assessment:english_score_exam_certification":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"TOEIC", "TEPS", "English exam"}, []string{"listening", "reading", "vocabulary", "grammar", "practice questions", "mock test"}, []string{"study abroad", "conversation only", "course sale", "score guarantee ad"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"TOEIC", "exam structure", "study plan"}, []string{"score target", "part overview", "diagnostic"}),
				patternStageTemplate("en", "core_pattern", []string{"TOEIC listening", "TOEIC reading", "practice questions"}, []string{"part 5", "part 7", "dictation"}),
				patternStageTemplate("en", "integration_practice", []string{"mock test", "time management", "wrong answers"}, []string{"review", "weak points", "strategy"}),
				patternStageTemplate("en", "performance_prep", []string{"exam day", "final review", "mock test"}, []string{"checklist", "pace", "mistake prevention"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"토익", "텝스", "영어 시험"}, []string{"청해", "독해", "어휘", "문법", "문제풀이", "모의고사"}, []string{"유학", "회화만", "강의 판매", "점수 보장 광고", "합격 보장", "학원 모집", "학원 광고", "자격 대행", "수학", "수학학원", "편입", "학원"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"토익", "시험 구조", "학습 계획"}, []string{"목표 점수", "파트 구성", "진단"}),
			patternStageTemplate("ko", "core_pattern", []string{"토익 LC", "토익 RC", "문제풀이"}, []string{"파트5", "파트7", "딕테이션"}),
			patternStageTemplate("ko", "integration_practice", []string{"토익 모의고사", "영어 오답 정리", "시간 관리"}, []string{"복습", "약점", "풀이 전략"}),
			patternStageTemplate("ko", "performance_prep", []string{"시험 직전", "최종 점검", "모의고사"}, []string{"체크리스트", "페이스", "실수 방지"}),
		}), true
	case "cooking_baking:certification_assessment:baking_written_practical_certification":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"baking certification", "written exam", "practical test"}, []string{"bread making", "pastry", "recipe process", "sanitation", "mock test"}, []string{"bakery sale", "class recruitment", "franchise", "certificate agency", "pass guarantee", "school ad"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"baking certification", "exam structure", "tools"}, []string{"written", "practical", "sanitation"}),
				patternStageTemplate("en", "core_pattern", []string{"bread making", "pastry technique", "practical test"}, []string{"mixing", "proofing", "oven"}),
				patternStageTemplate("en", "artifact_finish", []string{"finished product", "shape check", "quality criteria"}, []string{"texture", "color", "scoring"}),
				patternStageTemplate("en", "performance_prep", []string{"mock practical test", "time management", "final checklist"}, []string{"exam station", "mistake fixes", "review"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"제과제빵 자격증", "필기", "실기"}, []string{"제빵", "제과", "공정", "위생", "모의고사"}, []string{"빵집 판매", "클래스 모집", "프랜차이즈", "자격 대행", "합격 보장", "학원 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"제과제빵 자격증", "시험 구조", "도구"}, []string{"필기", "실기", "위생"}),
			patternStageTemplate("ko", "core_pattern", []string{"제빵", "제과 기법", "실기 시험"}, []string{"믹싱", "발효", "오븐"}),
			patternStageTemplate("ko", "artifact_finish", []string{"제품 완성", "모양 점검", "채점 기준"}, []string{"식감", "색", "스코어링"}),
			patternStageTemplate("ko", "performance_prep", []string{"실기 모의고사", "시간 관리", "최종 체크"}, []string{"시험장", "실수 보완", "복습"}),
		}), true
	case "cooking_baking:certification_assessment:korean_cooking_practical_certification":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"Korean cooking practical", "certification", "exam dish"}, []string{"knife skills", "garnish", "time management", "sanitation", "plating"}, []string{"restaurant review", "class recruitment", "recipe sale", "certificate agency", "pass guarantee", "school ad"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"Korean cooking certification", "practical exam", "tools"}, []string{"sanitation", "ingredients", "exam list"}),
				patternStageTemplate("en", "core_pattern", []string{"Korean dish", "knife skills", "practical recipe"}, []string{"julienne", "seasoning", "garnish"}),
				patternStageTemplate("en", "artifact_finish", []string{"exam dish", "plating", "quality criteria"}, []string{"shape", "temperature", "presentation"}),
				patternStageTemplate("en", "performance_prep", []string{"mock practical test", "time management", "final checklist"}, []string{"sequence", "mistake fixes", "exam prep"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"한식조리기능사", "실기", "시험 메뉴"}, []string{"칼질", "고명", "시간 관리", "위생", "담음새"}, []string{"맛집 후기", "클래스 모집", "레시피 판매", "자격 대행", "합격 보장", "학원 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"한식조리기능사", "실기 시험", "도구"}, []string{"위생", "재료", "시험 메뉴"}),
			patternStageTemplate("ko", "core_pattern", []string{"한식 메뉴", "칼질", "실기 레시피"}, []string{"채썰기", "양념", "고명"}),
			patternStageTemplate("ko", "artifact_finish", []string{"시험 메뉴 완성", "담음새", "채점 기준"}, []string{"모양", "온도", "제출 상태"}),
			patternStageTemplate("ko", "performance_prep", []string{"실기 모의고사", "시간 관리", "최종 체크"}, []string{"순서", "실수 보완", "시험 준비"}),
		}), true
	case "body_movement:performance_execution:dance_cover_song_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"k-pop dance", "dance cover", "choreography tutorial"},
				[]string{"mirror mode", "easy choreography", "section practice", "count practice", "rehearsal"},
				[]string{"dance workout", "fitness only", "academy ad", "audition ad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"k-pop dance", "easy choreography", "mirror mode"}, []string{"count practice", "beginner", "warm up"}),
					patternStageTemplate("en", "first_output", []string{"dance cover", "chorus choreography", "mirror mode"}, []string{"point move", "slow practice", "beginner"}),
					patternStageTemplate("en", "core_pattern", []string{"choreography tutorial", "section practice", "dance cover"}, []string{"counts", "transition", "repeat practice"}),
					patternStageTemplate("en", "integration_practice", []string{"dance cover", "section transition", "mirror practice"}, []string{"timing", "angle", "clean move"}),
					patternStageTemplate("en", "performance_prep", []string{"dance cover", "full run-through", "recording rehearsal"}, []string{"expression", "self review", "final practice"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"댄스 커버", "안무 배우기", "거울모드"},
			[]string{"케이팝 댄스", "쉬운 안무", "구간 연습", "박자 연습", "리허설"},
			[]string{"다이어트댄스", "운동 챌린지", "학원 광고", "오디션 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"케이팝 댄스", "쉬운 안무", "거울모드"}, []string{"박자 연습", "입문", "준비운동"}),
				patternStageTemplate("ko", "first_output", []string{"댄스 커버", "후렴 안무", "거울모드"}, []string{"포인트 안무", "느린 연습", "초보"}),
				patternStageTemplate("ko", "core_pattern", []string{"안무 배우기", "구간 연습", "댄스 커버"}, []string{"카운트", "동작 연결", "반복 연습"}),
				patternStageTemplate("ko", "integration_practice", []string{"댄스 커버", "구간 연결", "거울 연습"}, []string{"타이밍", "각도", "동작 정리"}),
				patternStageTemplate("ko", "performance_prep", []string{"댄스 커버", "전체 루틴", "촬영 리허설"}, []string{"표정", "자가 점검", "최종 연습"}),
			}), true
	case "cooking_baking:artifact_creation:home_cafe_latte_foundation":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"home cafe", "latte", "milk steaming"},
				[]string{"espresso", "latte recipe", "milk texture", "beginner", "home barista"},
				[]string{"cafe startup", "machine sale", "academy ad", "product sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"home cafe", "latte", "beginner setup"}, []string{"beans", "espresso", "milk"}),
					patternStageTemplate("en", "first_output", []string{"latte recipe", "make latte", "home cafe"}, []string{"espresso base", "milk", "ratio"}),
					patternStageTemplate("en", "core_pattern", []string{"milk steaming", "latte", "milk texture"}, []string{"velvet milk", "foam", "steam wand"}),
					patternStageTemplate("en", "integration_practice", []string{"home cafe latte", "latte ratio", "taste adjustment"}, []string{"sweetness", "strength", "repeat recipe"}),
					patternStageTemplate("en", "artifact_finish", []string{"home cafe latte", "repeatable recipe", "serving"}, []string{"recipe card", "taste check", "latte art basics"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"홈카페", "라떼 만들기", "카페라떼"},
			[]string{"우유 스티밍", "우유거품", "에스프레소", "라떼 레시피", "우유 질감", "초보", "홈바리스타"},
			[]string{"카페 창업", "머신 판매", "학원 광고", "제품 판매"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"홈카페", "카페라떼", "초보 세팅"}, []string{"라떼", "우유거품", "원두", "에스프레소", "우유"}),
				patternStageTemplate("ko", "first_output", []string{"라떼 만들기", "라떼 레시피", "홈카페"}, []string{"에스프레소 베이스", "우유", "비율"}),
				patternStageTemplate("ko", "core_pattern", []string{"우유 스티밍", "우유거품", "우유 질감"}, []string{"라떼", "카페라떼", "벨벳밀크", "거품", "스팀 노즐"}),
				patternStageTemplate("ko", "integration_practice", []string{"홈카페 라떼", "라떼 비율", "맛 조절"}, []string{"단맛", "농도", "반복 레시피"}),
				patternStageTemplate("ko", "artifact_finish", []string{"홈카페 라떼", "반복 레시피", "서빙"}, []string{"레시피 카드", "맛 점검", "라떼아트 기초"}),
			}), true
	case "maker_technical_hobby:certification_assessment:drone_operator_basic_4class":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"drone license", "basic drone", "safety rules"},
				[]string{"online education", "pre-flight checklist", "flight rules", "exam prep", "beginner"},
				[]string{"answer sheet", "guaranteed pass", "agency", "drone sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"drone license", "basic drone", "online education"}, []string{"beginner", "registration", "course overview"}),
					patternStageTemplate("en", "core_pattern", []string{"drone rules", "safety rules", "basic drone"}, []string{"airspace", "pre-flight", "battery"}),
					patternStageTemplate("en", "integration_practice", []string{"drone safety", "rule scenarios", "online test"}, []string{"quiz", "mistake review", "checklist"}),
					patternStageTemplate("en", "performance_prep", []string{"drone license", "exam prep", "final checklist"}, []string{"online education", "test review", "safety"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"드론 4종", "자격증", "안전 법규"},
			[]string{"온라인 교육", "교육이수", "비행 안전", "시험 준비", "초보"},
			[]string{"답안지", "무조건합격", "자격 대행", "드론 판매", "합격 보장", "학원 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"드론 4종", "자격증", "온라인 교육"}, []string{"초보", "등록", "교육 과정"}),
				patternStageTemplate("ko", "core_pattern", []string{"드론 법규", "안전 법규", "드론 4종"}, []string{"비행 금지", "비행 전 점검", "배터리"}),
				patternStageTemplate("ko", "integration_practice", []string{"드론 안전", "법규 문제", "온라인 시험"}, []string{"퀴즈", "오답 정리", "체크리스트"}),
				patternStageTemplate("ko", "performance_prep", []string{"드론 4종", "시험 준비", "최종 점검"}, []string{"온라인 교육", "시험 복습", "안전"}),
			}), true
	case "visual_art:artifact_creation:sequential_creative_class_path":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"illustration", "portfolio", "creative project"},
				[]string{"beginner", "art project", "series", "project planning", "showcase"},
				[]string{"job portfolio", "admission portfolio", "outsourcing", "course sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"illustration", "beginner project", "tools"}, []string{"style reference", "simple theme", "sketchbook"}),
					patternStageTemplate("en", "first_output", []string{"illustration", "first artwork", "beginner"}, []string{"sketch", "color", "simple character"}),
					patternStageTemplate("en", "core_pattern", []string{"illustration technique", "art project", "composition"}, []string{"line", "color", "shape"}),
					patternStageTemplate("en", "integration_practice", []string{"illustration series", "creative project", "portfolio"}, []string{"theme", "consistency", "review"}),
					patternStageTemplate("en", "artifact_finish", []string{"portfolio", "finished artwork", "project review"}, []string{"showcase", "selection", "presentation"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"},
			[]string{"그림 독학", "일러스트 기초", "작품 만들기"},
			[]string{"작품 주제", "캐릭터 그리기", "스케치", "채색", "작품 정리"},
			[]string{"취업 포트폴리오", "입시 포트폴리오", "외주", "클래스 오픈", "강의 판매"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"그림 독학", "일러스트 기초", "작품 주제"}, []string{"스케치북", "참고자료", "간단한 주제"}),
				patternStageTemplate("ko", "first_output", []string{"그림 그리기", "첫 작품", "일러스트 기초"}, []string{"스케치", "채색", "간단한 캐릭터"}),
				patternStageTemplate("ko", "core_pattern", []string{"일러스트 기법", "캐릭터 그리기", "구도"}, []string{"선", "색", "형태"}),
				patternStageTemplate("ko", "integration_practice", []string{"작품 시리즈", "창작 프로젝트", "작품 점검"}, []string{"주제", "일관성", "피드백"}),
				patternStageTemplate("ko", "artifact_finish", []string{"작품 정리", "포트폴리오", "완성 작품"}, []string{"선별", "소개글", "발표"}),
			}), true
	case "digital_creation:teaching_instruction:creator_tutorial_publish":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"lesson video production", "tutorial video editing", "screen recording lesson"},
				[]string{"demo script", "clear explanation", "screen recording", "voiceover", "basic editing", "captions", "upload checklist"},
				[]string{"gear sale", "course sale", "growth hack", "views guarantee", "monetization", "shorts automation", "subscriber growth"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"lesson video production", "recording setup", "audio check"}, []string{"clear shots", "beginner", "simple script"}),
					patternStageTemplate("en", "structure_plan", []string{"how-to video", "demo script", "lesson flow"}, []string{"problem", "steps", "result"}),
					patternStageTemplate("en", "core_pattern", []string{"screen recording", "voiceover", "basic editing"}, []string{"cutting", "captions", "demo clarity"}),
					patternStageTemplate("en", "artifact_finish", []string{"tutorial video", "captions", "upload checklist"}, []string{"description", "final review", "publish prep"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"교육 영상 제작", "튜토리얼 영상 편집", "화면 녹화 강의"},
			[]string{"시연 흐름", "스크립트", "화면 녹화", "보이스오버", "기본 편집", "자막", "업로드 체크리스트"},
			[]string{"장비 판매", "강의 판매", "수익화 강의", "조회수 보장", "조회수", "구독자", "알고리즘", "쇼츠 자동화", "인스타"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"교육 영상 제작", "촬영 세팅", "오디오 체크"}, []string{"구도", "초보", "간단한 대본"}),
				patternStageTemplate("ko", "structure_plan", []string{"강의 영상 제작", "시연 흐름", "스크립트"}, []string{"문제", "단계", "결과물"}),
				patternStageTemplate("ko", "core_pattern", []string{"화면 녹화", "보이스오버", "기본 편집"}, []string{"컷 편집", "자막", "시연 전달"}),
				patternStageTemplate("ko", "artifact_finish", []string{"튜토리얼 영상", "자막", "업로드 체크리스트"}, []string{"설명문", "최종 점검", "공개 준비"}),
			}), true
	case "digital_creation:presentation_publish:portfolio_publish":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article"},
				[]string{"project upload", "showcase page", "project description"},
				[]string{"export settings", "thumbnail", "description", "public checklist", "portfolio publish", "final package"},
				[]string{"course sale", "job guarantee", "views guarantee", "template sale", "admission portfolio", "portfolio academy", "investment portfolio", "interview portfolio"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"portfolio publish", "choose project", "platform"}, []string{"audience", "success criteria", "checklist"}),
					patternStageTemplate("en", "core_pattern", []string{"project upload", "export settings", "description"}, []string{"thumbnail", "cover", "accessibility"}),
					patternStageTemplate("en", "artifact_finish", []string{"showcase page", "final package", "portfolio review"}, []string{"selection", "presentation", "public checklist"}),
					patternStageTemplate("en", "publish_prep", []string{"portfolio publish", "upload checklist", "public release"}, []string{"final review", "description", "thumbnail"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article"},
			[]string{"작품 업로드", "프로젝트 페이지", "작품 소개"},
			[]string{"내보내기 설정", "썸네일", "설명문", "공개 체크리스트", "포트폴리오 공개", "최종 패키지"},
			[]string{"강의 판매", "취업 보장", "조회수 보장", "템플릿 판매", "취업 포트폴리오", "입시 포트폴리오", "합격 포트폴리오", "포트폴리오 학원", "투자 포트폴리오", "국민연금", "면접", "디자이너 취업"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"포트폴리오 공개", "작품 선택", "플랫폼"}, []string{"대상 독자", "성공 기준", "체크리스트"}),
				patternStageTemplate("ko", "core_pattern", []string{"작품 업로드", "내보내기 설정", "설명문"}, []string{"썸네일", "커버", "접근성"}),
				patternStageTemplate("ko", "artifact_finish", []string{"프로젝트 페이지", "최종 패키지", "포트폴리오 점검"}, []string{"선별", "발표", "공개 체크리스트"}),
				patternStageTemplate("ko", "publish_prep", []string{"포트폴리오 공개", "업로드 체크리스트", "공개 준비"}, []string{"최종 점검", "설명문", "썸네일"}),
			}), true
	case "digital_creation:artifact_creation:freeware_midi_full_track":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"free DAW", "MIDI composition", "finish a track"},
				[]string{"Cakewalk", "LMMS", "drum pattern", "bassline", "arrangement", "mixing", "export"},
				[]string{"plugin sale", "sample pack sale", "beat sale", "course sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"free DAW", "MIDI setup", "beginner"}, []string{"Cakewalk", "LMMS", "instrument track"}),
					patternStageTemplate("en", "first_output", []string{"MIDI beat", "drum pattern", "bassline"}, []string{"loop", "quantize", "simple melody"}),
					patternStageTemplate("en", "core_pattern", []string{"MIDI arrangement", "chord progression", "track structure"}, []string{"intro", "verse", "chorus"}),
					patternStageTemplate("en", "artifact_finish", []string{"finish a track", "mixing", "export"}, []string{"volume balance", "wav", "mp3"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"무료 DAW", "MIDI 작곡", "트랙 완성"},
			[]string{"Cakewalk", "LMMS", "드럼 패턴", "베이스라인", "코드 진행", "편곡", "믹싱", "내보내기"},
			[]string{"플러그인 판매", "샘플팩 판매", "비트 판매", "유료 강의"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"무료 DAW", "MIDI 세팅", "초보"}, []string{"Cakewalk", "LMMS", "악기 트랙"}),
				patternStageTemplate("ko", "first_output", []string{"MIDI 비트", "드럼 패턴", "베이스라인"}, []string{"루프", "퀀타이즈", "간단한 멜로디"}),
				patternStageTemplate("ko", "core_pattern", []string{"MIDI 편곡", "코드 진행", "트랙 구조"}, []string{"인트로", "벌스", "코러스"}),
				patternStageTemplate("ko", "artifact_finish", []string{"트랙 완성", "믹싱", "내보내기"}, []string{"볼륨 밸런스", "wav", "mp3"}),
			}), true
	case "digital_creation:artifact_creation:notion_study_dashboard":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"},
				[]string{"Notion study dashboard", "task database", "weekly review"},
				[]string{"notes page", "study planner", "template build", "relation database", "progress tracker"},
				[]string{"template sale", "course sale", "CRM", "team workspace ad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"Notion study dashboard", "study planner", "beginner"}, []string{"pages", "database", "weekly review"}),
					patternStageTemplate("en", "structure_plan", []string{"Notion task database", "notes page", "study dashboard"}, []string{"subjects", "deadlines", "properties"}),
					patternStageTemplate("en", "core_pattern", []string{"Notion database", "progress tracker", "weekly review"}, []string{"filters", "relations", "views"}),
					patternStageTemplate("en", "artifact_finish", []string{"Notion study template", "dashboard review", "habit tracker"}, []string{"cleanup", "duplicate template", "routine"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"},
			[]string{"노션 학습 대시보드", "노션 공부 관리", "할일 데이터베이스"},
			[]string{"노트 페이지", "주간 리뷰", "학습 플래너", "템플릿 만들기", "진도 관리", "데이터베이스 보기"},
			[]string{"템플릿 판매", "강의 판매", "업무용 CRM", "팀 협업 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"노션 학습 대시보드", "학습 플래너", "초보"}, []string{"페이지", "데이터베이스", "주간 리뷰"}),
				patternStageTemplate("ko", "structure_plan", []string{"노션 할일 데이터베이스", "노트 페이지", "학습 대시보드"}, []string{"과목", "마감일", "속성"}),
				patternStageTemplate("ko", "core_pattern", []string{"노션 데이터베이스", "진도 관리", "주간 리뷰"}, []string{"필터", "관계형", "보기"}),
				patternStageTemplate("ko", "artifact_finish", []string{"노션 학습 템플릿", "대시보드 점검", "습관 트래커"}, []string{"정리", "복제 템플릿", "루틴"}),
			}), true
	case "knowledge_hobby:habit_lifestyle:personal_finance_budget_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"},
				[]string{"personal budget", "expense tracker", "monthly budget"},
				[]string{"spending categories", "Google Sheets budget", "Excel budget", "saving plan", "weekly review"},
				[]string{"stock recommendation", "crypto", "loan ad", "insurance consultation", "get rich quick"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"personal budget", "expense tracker", "beginner"}, []string{"income", "fixed cost", "categories"}),
					patternStageTemplate("en", "structure_plan", []string{"monthly budget", "spending categories", "budget template"}, []string{"Google Sheets", "Excel", "cash flow"}),
					patternStageTemplate("en", "core_pattern", []string{"track spending", "weekly review", "saving plan"}, []string{"overspending", "category review", "habit"}),
					patternStageTemplate("en", "routine_plan", []string{"budget routine", "monthly review", "expense tracker"}, []string{"repeatable", "checklist", "next month"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"},
			[]string{"가계부", "예산 관리", "지출 기록"},
			[]string{"월간 예산", "소비 분류", "엑셀 가계부", "구글시트 가계부", "저축 계획", "주간 점검"},
			[]string{"투자 추천", "코인", "주식 종목", "대출 광고", "보험 상담", "고수익 보장"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"가계부", "지출 기록", "초보"}, []string{"수입", "고정비", "소비 분류"}),
				patternStageTemplate("ko", "structure_plan", []string{"월간 예산", "소비 분류", "가계부 양식"}, []string{"구글시트", "엑셀", "현금 흐름"}),
				patternStageTemplate("ko", "core_pattern", []string{"지출 기록", "주간 점검", "저축 계획"}, []string{"과소비", "분류 점검", "습관"}),
				patternStageTemplate("ko", "routine_plan", []string{"예산 루틴", "월간 점검", "가계부"}, []string{"반복", "체크리스트", "다음 달"}),
			}), true
	case "craft_making:artifact_creation:pottery_handbuilding_cup_project":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"handbuilding pottery", "ceramic cup", "clay mug"},
				[]string{"pinch pot", "slab cup", "coil cup", "attach handle", "score and slip", "drying"},
				[]string{"workshop booking", "pottery sale", "class recruitment", "studio ad"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"handbuilding pottery", "ceramic cup", "beginner"}, []string{"clay prep", "workspace", "moisture"}),
					patternStageTemplate("en", "first_output", []string{"pinch pot", "ceramic cup", "basic shape"}, []string{"wall thickness", "rim", "simple cup"}),
					patternStageTemplate("en", "core_pattern", []string{"slab cup", "attach handle", "score and slip"}, []string{"smooth surface", "join", "drying"}),
					patternStageTemplate("en", "artifact_finish", []string{"ceramic mug", "surface cleanup", "final check"}, []string{"glaze basics", "drying", "inspection"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"도자기 핸드빌딩", "도자기 컵 만들기", "도자기 머그컵"},
			[]string{"핀칭", "판 성형", "코일링", "손잡이 붙이기", "흙 붙이기", "건조"},
			[]string{"공방 예약", "도자기 판매", "체험 클래스", "공방 광고"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"도자기 핸드빌딩", "도자기 컵 만들기", "초보"}, []string{"흙 준비", "작업 공간", "수분 조절"}),
				patternStageTemplate("ko", "first_output", []string{"핀칭", "도자기 컵", "기본 형태"}, []string{"두께", "입구 정리", "간단한 컵"}),
				patternStageTemplate("ko", "core_pattern", []string{"판 성형", "손잡이 붙이기", "흙 붙이기"}, []string{"표면 정리", "접합", "건조"}),
				patternStageTemplate("ko", "artifact_finish", []string{"도자기 머그컵", "표면 정리", "최종 점검"}, []string{"유약 기초", "건조", "검수"}),
			}), true
	case "digital_creation:artifact_creation:freeware_midi_beatmaking":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"free DAW", "MIDI beatmaking", "beginner music production"},
				[]string{"Cakewalk", "LMMS", "drum pattern", "loop", "quantize", "bassline"},
				[]string{"plugin sale", "sample pack sale", "beat sale", "course sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"free DAW", "MIDI beatmaking", "beginner"}, []string{"Cakewalk", "LMMS", "first setup"}),
					patternStageTemplate("en", "first_output", []string{"MIDI beat", "drum pattern", "loop"}, []string{"quantize", "bassline", "simple melody"}),
					patternStageTemplate("en", "artifact_finish", []string{"beatmaking", "arrangement", "export"}, []string{"volume balance", "wav", "mp3"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"무료 DAW", "MIDI 비트메이킹", "초보 작곡"},
			[]string{"Cakewalk", "LMMS", "드럼 패턴", "루프 만들기", "퀀타이즈", "베이스라인"},
			[]string{"플러그인 판매", "샘플팩 판매", "비트 판매", "유료 강의"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"무료 DAW", "MIDI 비트메이킹", "초보"}, []string{"Cakewalk", "LMMS", "시작 방법"}),
				patternStageTemplate("ko", "first_output", []string{"MIDI 비트", "드럼 패턴", "루프"}, []string{"퀀타이즈", "베이스라인", "간단한 멜로디"}),
				patternStageTemplate("ko", "artifact_finish", []string{"비트메이킹", "편곡", "내보내기"}, []string{"볼륨 밸런스", "wav", "mp3"}),
			}), true
	case "digital_creation:artifact_creation:university_studio":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"},
				[]string{"digital artwork project", "design tool", "portfolio"},
				[]string{"project brief", "process", "export", "presentation", "portfolio review"},
				[]string{"admission portfolio", "job guarantee", "template sale", "course sale"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"digital artwork project", "design tool", "portfolio"}, []string{"project brief", "reference", "workflow"}),
					patternStageTemplate("en", "core_pattern", []string{"digital artwork", "design process", "portfolio"}, []string{"composition", "export", "review"}),
					patternStageTemplate("en", "artifact_finish", []string{"digital artwork", "presentation", "portfolio review"}, []string{"description", "final export", "critique"}),
				}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"},
			[]string{"디지털 작품 만들기", "디자인 툴", "포트폴리오"},
			[]string{"작품 기획", "제작 과정", "내보내기", "작품 발표", "포트폴리오 점검"},
			[]string{"입시 포트폴리오", "취업 보장", "템플릿 판매", "강의 판매"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"디지털 작품 만들기", "디자인 툴", "포트폴리오"}, []string{"작품 기획", "참고자료", "작업 흐름"}),
				patternStageTemplate("ko", "core_pattern", []string{"디지털 작품", "제작 과정", "포트폴리오"}, []string{"구도", "내보내기", "점검"}),
				patternStageTemplate("ko", "artifact_finish", []string{"디지털 작품", "작품 발표", "포트폴리오 점검"}, []string{"설명문", "최종 내보내기", "크리틱"}),
			}), true
	case "digital_creation:professional_transition:korean_ncs_web_publisher_training":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"web publisher", "HTML CSS", "responsive web"}, []string{"web standards", "accessibility", "layout", "portfolio", "publishing practice"}, []string{"bootcamp ad", "job guarantee", "course sale", "consultation", "agency"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"web publisher", "HTML CSS", "responsive web"}, []string{"beginner", "web standards", "accessibility"}),
				patternStageTemplate("en", "core_pattern", []string{"CSS layout", "responsive web", "publishing"}, []string{"flex", "grid", "media query"}),
				patternStageTemplate("en", "artifact_finish", []string{"web publisher", "portfolio page", "responsive check"}, []string{"accessibility", "validation", "publish"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"웹퍼블리셔", "HTML CSS", "반응형 웹"}, []string{"웹표준", "웹접근성", "레이아웃", "퍼블리싱 실무", "포트폴리오"}, []string{"학원 광고", "취업 보장", "국비 모집", "강의 판매", "상담 신청", "대행"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"웹퍼블리셔", "HTML CSS", "반응형 웹"}, []string{"기초", "웹표준", "웹접근성"}),
			patternStageTemplate("ko", "core_pattern", []string{"CSS 레이아웃", "반응형 웹", "퍼블리싱"}, []string{"flex", "grid", "미디어쿼리"}),
			patternStageTemplate("ko", "artifact_finish", []string{"웹퍼블리셔", "포트폴리오 페이지", "반응형 점검"}, []string{"접근성", "검증", "배포"}),
		}), true
	case "knowledge_hobby:foundation_build:korean_open_lecture_weekly_survey", "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"K-MOOC", "online course", "study plan"}, []string{"course registration", "weekly study", "lecture notes", "review", "completion"}, []string{"paid course sale", "certificate agency", "exam cram ad"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"K-MOOC", "online course", "registration"}, []string{"study plan", "beginner", "course guide"}),
				patternStageTemplate("en", "core_pattern", []string{"online course", "weekly study", "lecture notes"}, []string{"quiz review", "summary", "routine"}),
				patternStageTemplate("en", "routine_plan", []string{"online course", "completion", "study routine"}, []string{"schedule", "review", "next course"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"K-MOOC", "온라인 강좌", "수강 방법"}, []string{"수강신청", "주차별 학습", "강의 노트", "복습", "수료"}, []string{"유료 강의 판매", "자격 대행", "시험 광고", "학원 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"K-MOOC", "온라인 강좌", "수강 방법"}, []string{"수강신청", "학습 계획", "입문"}),
			patternStageTemplate("ko", "core_pattern", []string{"온라인 강좌", "주차별 학습", "강의 노트"}, []string{"퀴즈 복습", "요약", "루틴"}),
			patternStageTemplate("ko", "routine_plan", []string{"온라인 강좌", "수료", "학습 루틴"}, []string{"일정", "복습", "다음 강좌"}),
		}), true
	case "knowledge_hobby:professional_transition:credit_bank_practicum", "knowledge_hobby:professional_transition:credit_bank_standard_theory":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"academic credit bank", "online class", "credit recognition"}, []string{"practicum", "application", "coursework", "assignment", "completion"}, []string{"agency ad", "tuition sale", "job guarantee", "consultation", "recruitment"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"academic credit bank", "online class", "credit recognition"}, []string{"application", "requirements", "beginner"}),
				patternStageTemplate("en", "core_pattern", []string{"credit recognition", "online coursework", "assignment"}, []string{"discussion", "exam", "completion"}),
				patternStageTemplate("en", "artifact_finish", []string{"credit bank", "completion checklist", "application"}, []string{"documents", "review", "next term"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"학점은행제", "온라인 수업", "학점 인정"}, []string{"실습 과목", "학점인정 신청", "과제 작성", "토론", "수료 기준"}, []string{"대행 광고", "학원 광고", "취업 보장", "상담 신청", "국비 모집", "환급 광고"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"학점은행제", "온라인 수업", "학점 인정"}, []string{"신청 방법", "요건", "입문"}),
			patternStageTemplate("ko", "core_pattern", []string{"학점 인정", "온라인 수업", "과제 작성"}, []string{"토론", "시험", "수료"}),
			patternStageTemplate("ko", "artifact_finish", []string{"학점은행제", "수료 체크리스트", "학점인정 신청"}, []string{"서류", "점검", "다음 학기"}),
		}), true
	case "knowledge_hobby:teaching_instruction:credit_bank_instructional_design":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"instructional design", "learning objective", "lesson plan"}, []string{"assessment design", "teaching material", "activity", "feedback"}, []string{"course sale", "template sale", "job guarantee", "consultation", "agency"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"instructional design", "learning objective", "lesson plan"}, []string{"beginner", "course outline", "activity"}),
				patternStageTemplate("en", "core_pattern", []string{"lesson plan", "assessment design", "teaching material"}, []string{"rubric", "feedback", "activity"}),
				patternStageTemplate("en", "artifact_finish", []string{"lesson plan", "review", "teaching material"}, []string{"checklist", "revision", "presentation"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"교수설계", "학습목표", "수업지도안"}, []string{"평가 설계", "강의안 작성", "학습 활동", "피드백"}, []string{"템플릿 판매", "강의 판매", "취업 보장", "학원 광고", "상담 신청", "대행"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"교수설계", "학습목표", "수업지도안"}, []string{"기초", "수업 설계", "활동"}),
			patternStageTemplate("ko", "core_pattern", []string{"수업지도안", "평가 설계", "강의안 작성"}, []string{"루브릭", "피드백", "학습 활동"}),
			patternStageTemplate("ko", "artifact_finish", []string{"수업지도안", "점검", "강의 자료"}, []string{"체크리스트", "수정", "발표"}),
		}), true
	case "maker_technical_hobby:teaching_instruction:maker_workshop_demo":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en", []string{"maker workshop", "Arduino demo", "DIY prototype"}, []string{"materials", "safety", "circuit", "step demonstration", "project presentation"}, []string{"kit sale", "workshop recruitment", "equipment ad", "course sale"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"Arduino project", "circuit demo", "maker workshop"}, []string{"materials", "safety", "beginner"}),
				patternStageTemplate("en", "core_pattern", []string{"Arduino project", "circuit demo", "prototype"}, []string{"wiring", "sensor", "code"}),
				patternStageTemplate("en", "artifact_finish", []string{"DIY prototype", "project presentation", "final check"}, []string{"troubleshooting", "showcase", "reflection"}),
			}), true
		}
		return buildPatternTemplateFromStages("ko", []string{"메이커 워크숍", "아두이노", "DIY 프로토타입"}, []string{"재료 준비", "안전 안내", "전자회로", "제작 시연", "작품 발표"}, []string{"키트 판매", "체험 모집", "장비 광고", "강의 판매"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"아두이노 프로젝트", "전자회로 시연", "메이커 워크숍"}, []string{"재료 준비", "안전 안내", "입문"}),
			patternStageTemplate("ko", "core_pattern", []string{"아두이노 프로젝트", "전자회로 시연", "프로토타입"}, []string{"배선", "센서", "코드"}),
			patternStageTemplate("ko", "artifact_finish", []string{"DIY 프로토타입", "작품 발표", "최종 점검"}, []string{"문제 해결", "쇼케이스", "회고"}),
		}), true
	case "visual_art:presentation_publish:art_publish_showcase", "visual_art:presentation_publish:calligraphy_publish_showcase":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article", "video"}, []string{"art portfolio", "artwork description", "online exhibition"}, []string{"artwork photo", "artist statement", "portfolio upload", "showcase"}, []string{"admission portfolio", "job portfolio", "gallery sale", "course sale", "views guarantee", "reaction guarantee"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"art portfolio", "artwork description", "online exhibition"}, []string{"artwork photo", "artist statement", "beginner"}),
				patternStageTemplate("en", "core_pattern", []string{"artwork photo", "portfolio upload", "artwork description"}, []string{"lighting", "layout", "caption"}),
				patternStageTemplate("en", "artifact_finish", []string{"online exhibition", "showcase", "artist statement"}, []string{"final review", "publish checklist", "portfolio"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article", "video"}, []string{"작품 소개", "온라인 전시", "아트 포트폴리오"}, []string{"작품 사진", "작가노트", "작품 업로드", "쇼케이스"}, []string{"입시 포트폴리오", "취업 포트폴리오", "갤러리 판매", "강의 판매", "조회수 보장", "반응 보장"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"작품 소개", "온라인 전시", "아트 포트폴리오"}, []string{"작품 사진", "작가노트", "입문"}),
			patternStageTemplate("ko", "core_pattern", []string{"작품 사진", "작품 업로드", "작품 설명"}, []string{"조명", "레이아웃", "캡션"}),
			patternStageTemplate("ko", "artifact_finish", []string{"온라인 전시", "쇼케이스", "작가노트"}, []string{"최종 점검", "공개 체크리스트", "포트폴리오"}),
		}), true
	case "writing_storytelling:presentation_publish:university_studio":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStagesWithContentTypes("en", []string{"article"}, []string{"creative writing", "essay revision", "writing workshop"}, []string{"draft", "peer critique", "revision", "portfolio", "publication prep"}, []string{"admission essay", "ghostwriting", "course sale", "contest ad", "publishing agency", "monetization course"}, []LessonRecommendationSearchSpec{
				patternStageTemplate("en", "setup_intro", []string{"creative writing", "essay draft", "writing workshop"}, []string{"topic", "voice", "beginner"}),
				patternStageTemplate("en", "core_pattern", []string{"essay revision", "peer critique", "draft"}, []string{"structure", "sentence", "feedback"}),
				patternStageTemplate("en", "artifact_finish", []string{"writing portfolio", "revision", "publication prep"}, []string{"selection", "final edit", "author note"}),
			}), true
		}
		return buildPatternTemplateFromStagesWithContentTypes("ko", []string{"article"}, []string{"창작 글쓰기", "에세이 퇴고", "글쓰기 워크숍"}, []string{"초안", "합평", "퇴고", "글 모음", "발표 준비"}, []string{"입시 자소서", "대필", "강의 판매", "공모전 광고", "출판 대행", "수익화 강의", "조회수 보장", "취업 포트폴리오"}, []LessonRecommendationSearchSpec{
			patternStageTemplate("ko", "setup_intro", []string{"창작 글쓰기", "에세이 초안", "글쓰기 워크숍"}, []string{"글감", "문체", "입문"}),
			patternStageTemplate("ko", "core_pattern", []string{"에세이 퇴고", "합평", "초안"}, []string{"구조", "문장", "피드백"}),
			patternStageTemplate("ko", "artifact_finish", []string{"글 모음", "퇴고", "발표 준비"}, []string{"선별", "최종 교정", "작가노트"}),
		}), true
	case "body_movement:habit_lifestyle:yoga_daily_routine":
		if normalizeLearningLanguage(language) == "en" {
			return buildPatternTemplateFromStages("en",
				[]string{"beginner yoga", "daily yoga", "yoga routine"},
				[]string{"breathing", "gentle flow", "morning yoga", "15-minute yoga", "stretching", "safe posture"},
				[]string{"weight loss promise", "advanced pose", "teacher training", "workout challenge"},
				[]LessonRecommendationSearchSpec{
					patternStageTemplate("en", "setup_intro", []string{"beginner yoga", "breathing", "safe posture"}, []string{"morning routine", "gentle stretch", "10 minute"}),
					patternStageTemplate("en", "core_pattern", []string{"yoga routine", "basic poses", "gentle flow"}, []string{"sun salutation", "stretching", "beginner sequence"}),
					patternStageTemplate("en", "integration_practice", []string{"yoga flow", "pose transition", "breath"}, []string{"cool down", "balance", "flexibility"}),
					patternStageTemplate("en", "routine_plan", []string{"daily yoga", "15-minute routine", "repeatable practice"}, []string{"weekly plan", "self check", "habit"}),
				}), true
		}
		return buildPatternTemplateFromStages("ko",
			[]string{"초보 요가", "데일리 요가", "요가 루틴"},
			[]string{"요가 호흡", "부드러운 플로우", "아침 요가", "15분 요가", "스트레칭", "안전한 자세"},
			[]string{"다이어트 광고", "고난도 자세", "강사 자격증", "운동 챌린지"},
			[]LessonRecommendationSearchSpec{
				patternStageTemplate("ko", "setup_intro", []string{"초보 요가", "호흡", "안전한 자세"}, []string{"아침 루틴", "가벼운 스트레칭", "10분"}),
				patternStageTemplate("ko", "core_pattern", []string{"요가 루틴", "기본 자세", "부드러운 플로우"}, []string{"태양경배", "스트레칭", "초보 시퀀스"}),
				patternStageTemplate("ko", "integration_practice", []string{"요가 플로우", "자세 전환", "호흡"}, []string{"마무리 스트레칭", "균형", "유연성"}),
				patternStageTemplate("ko", "routine_plan", []string{"데일리 요가", "15분 루틴", "반복 실천"}, []string{"주간 계획", "자기 점검", "습관"}),
			}), true
	default:
		return CurriculumPatternRecommendationSearchSpecTemplate{}, false
	}
}

func buildPatternTemplateFromStages(language string, mustInclude, niceToHave, avoid []string, stages []LessonRecommendationSearchSpec) CurriculumPatternRecommendationSearchSpecTemplate {
	return buildPatternTemplateFromStagesWithContentTypes(language, []string{"video"}, mustInclude, niceToHave, avoid, stages)
}

func buildPatternTemplateFromStagesWithContentTypes(language string, contentTypes, mustInclude, niceToHave, avoid []string, stages []LessonRecommendationSearchSpec) CurriculumPatternRecommendationSearchSpecTemplate {
	language = normalizeLearningLanguage(language)
	contentTypes = normalizeSearchSpecContentTypes(contentTypes)
	if len(contentTypes) == 0 {
		contentTypes = []string{"video"}
	}
	for idx := range stages {
		stages[idx] = NormalizeLessonRecommendationSearchSpec(stages[idx], LessonRecommendationSearchSpec{
			Intent:       "tutorial",
			Avoid:        avoid,
			ContentTypes: contentTypes,
			Language:     language,
			Source:       SearchSpecSourcePatternTemplate,
		})
		stages[idx].ContentTypes = append([]string{}, contentTypes...)
	}
	return CurriculumPatternRecommendationSearchSpecTemplate{
		Intent:         "tutorial",
		ContentTypes:   contentTypes,
		Language:       language,
		Source:         SearchSpecSourcePatternTemplate,
		MustInclude:    normalizeSearchSpecTokens(mustInclude, 6),
		NiceToHave:     normalizeSearchSpecTokens(niceToHave, 8),
		Avoid:          normalizeSearchSpecTokens(avoid, 10),
		StageTemplates: stages,
	}
}

func patternStageTemplate(language, stageRole string, mustInclude, niceToHave []string) LessonRecommendationSearchSpec {
	return LessonRecommendationSearchSpec{
		Intent:       "tutorial",
		MustInclude:  mustInclude,
		NiceToHave:   niceToHave,
		ContentTypes: []string{"video"},
		Language:     normalizeLearningLanguage(language),
		StageRole:    stageRole,
		Source:       SearchSpecSourcePatternTemplate,
	}
}

type englishCurriculumPatternOverride struct {
	Title                   string
	Summary                 string
	TriggerKeywords         []string
	StageRules              []string
	CompletionCriteriaRules []string
	BadPatterns             []string
}

var englishCurriculumPatternOverrides = map[string]englishCurriculumPatternOverride{
	"visual_art:artifact_creation:photography_fundamentals_project": {
		Title:   "Smartphone Photography Fundamentals Project",
		Summary: "A beginner-friendly photography path that turns composition, natural light, focus, and simple editing into a small finished photo set.",
		TriggerKeywords: []string{
			"take better photos", "phone camera", "smartphone camera", "everyday photography", "photo composition",
			"natural light portraits", "before and after edit", "mini photo series", "photo walk",
		},
		StageRules: []string{
			"Start with camera handling, focus, exposure, and clean framing using a phone or entry-level camera.",
			"Move quickly into composition, natural light, background control, and repeatable shooting practice.",
			"Include one lesson that compares several attempts and selects stronger photos for a small set.",
			"End with simple editing, export, and review of a coherent mini photo project.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can shoot, select, lightly edit, and explain a small photo set, not just define camera terms."},
		BadPatterns:             []string{"camera-spec lectures only", "ISO/shutter/aperture definitions with no shooting task", "editing before any photo practice"},
	},
	"instrument_performance:performance_execution:song_driven_instrument_path": {
		Title:   "Song-Driven Instrument Practice Path",
		Summary: "An instrument-learning path anchored in a favorite song, using sections, riffs, rhythm patterns, and rehearsal loops to reach playable song parts.",
		TriggerKeywords: []string{
			"learn guitar through a song", "favorite song cover", "play along", "riff practice", "bass groove", "drum groove",
			"song sections", "verse chorus bridge", "song-based practice", "cover performance",
		},
		StageRules: []string{
			"Choose a realistic target song or section before expanding into general technique.",
			"Break the target into intro, verse, chorus, riff, groove, or accompaniment sections that can be practiced separately.",
			"Teach the minimum technique needed for each section at the moment it is used.",
			"End with a connected play-through or rehearsal of the selected song section rather than unrelated drills.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can play through the chosen song section or simplified arrangement with steady timing."},
		BadPatterns:             []string{"months of technique before touching the target song", "generic scale lessons unrelated to the song", "instrument maintenance as a core lesson"},
	},
	"instrument_performance:performance_execution:song_completion": {
		Title:   "One-Song Completion Path",
		Summary: "A focused performance path for finishing one complete song or piece through section practice, transitions, rehearsal, and a final run-through.",
		TriggerKeywords: []string{
			"finish one song", "complete one piece", "play a full song", "section practice", "transition practice",
			"rehearse the whole piece", "full run-through", "performance ready",
		},
		StageRules: []string{
			"Confirm the exact song, difficulty level, and simplified arrangement before planning lessons.",
			"Build from hardest sections and transitions toward a stable full run-through.",
			"Include rehearsal habits such as slow tempo, looped trouble spots, and recorded self-checks.",
			"Keep the final lesson focused on a complete practice run and refinement checklist.",
		},
		CompletionCriteriaRules: []string{"Completion means one complete song can be performed or recorded at the learner's target level."},
		BadPatterns:             []string{"broad instrument foundation with no target song", "song history as a main lesson", "public performance as the lesson itself"},
	},
	"language_communication:performance_execution:travel_conversation": {
		Title:   "Travel Conversation Practice Path",
		Summary: "A practical language path for handling airport, hotel, restaurant, transport, shopping, and problem-solving conversations during travel.",
		TriggerKeywords: []string{
			"travel English", "travel conversation", "airport conversation", "hotel check-in", "restaurant ordering",
			"asking directions", "travel phrases", "role play", "situational dialogue",
		},
		StageRules: []string{
			"Organize lessons by real travel situations rather than grammar chapters.",
			"Pair essential phrases with short role-play and response practice in each lesson.",
			"Include listening and speaking turns for common variations, misunderstandings, and polite requests.",
			"End with a connected travel-day simulation that combines several situations.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can handle a short realistic travel conversation without reading a full script."},
		BadPatterns:             []string{"grammar terminology as standalone lessons", "memorizing long phrase lists with no role-play", "country trivia instead of conversation practice"},
	},
	"craft_making:artifact_creation:woodworking_small_project": {
		Title:   "Small Woodworking Project Path",
		Summary: "A project-based woodworking path that moves from measuring and tool safety to cutting, sanding, assembly, finishing, and checking a small wooden item.",
		TriggerKeywords: []string{
			"beginner woodworking", "small woodworking project", "wooden shelf", "cutting board", "wood craft",
			"measure cut sand", "wood assembly", "wood oil finish", "hand tools", "maker project",
		},
		StageRules: []string{
			"Start with safety, material choice, measuring, and a simple plan for one small object.",
			"Sequence cutting, sanding, assembly, and finishing as hands-on project stages.",
			"Keep tool explanations attached to the next project action the learner must perform.",
			"End with finishing, checking stability, and documenting the completed object.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has made a small usable wooden object with acceptable fit, sanding, and finish."},
		BadPatterns:             []string{"wood species encyclopedia lessons", "advanced joinery before a first object", "tool shopping as a standalone course goal"},
	},
	"craft_making:artifact_creation:sewing_project_foundation": {
		Title:   "Beginner Sewing Project Path",
		Summary: "A sewing path for completing a simple usable item such as a pouch, tote, or cushion cover through fabric prep, straight stitching, seams, and finishing.",
		TriggerKeywords: []string{
			"beginner sewing project", "sewing machine basics", "make a pouch", "make a tote bag", "fabric cutting",
			"straight stitch practice", "seam allowance", "hemming", "finish a sewing project",
		},
		StageRules: []string{
			"Begin with machine setup, fabric choice, measurements, and pattern or cutting preparation.",
			"Introduce straight stitching and seam allowance through the actual item being made.",
			"Use assembly steps such as side seams, corners, handles, lining, or closure only when needed for the selected item.",
			"End with seam finishing, pressing, inspection, and a usable completed sewing project.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has completed one simple sewing item that can be used or gifted."},
		BadPatterns:             []string{"all stitch types before one project", "fashion design theory only", "fabric shopping as a main lesson"},
	},
	"cooking_baking:artifact_creation:home_cafe_latte_foundation": {
		Title:   "Home Cafe Latte Foundation Path",
		Summary: "A home cafe path for making repeatable lattes by learning beans, extraction or strong coffee base, milk texture, ratios, taste adjustment, and serving.",
		TriggerKeywords: []string{
			"make latte at home", "home cafe latte", "espresso at home", "milk frothing", "milk texture",
			"latte ratio", "coffee recipe", "repeatable cup", "home barista",
		},
		StageRules: []string{
			"Start with equipment, beans, grind or coffee base, and a simple repeatable setup.",
			"Teach extraction or strong coffee preparation before milk texture and ratio adjustment.",
			"Include side-by-side tasting so the learner can adjust strength, sweetness, and milk balance.",
			"End with a repeatable latte recipe card and one served cup that matches the learner's taste.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can reproduce a latte recipe at home with stable taste and texture."},
		BadPatterns:             []string{"cafe business planning", "latte art before milk texture", "coffee history as the main path"},
	},
	"writing_storytelling:habit_lifestyle:daily_writing_habit": {
		Title:   "Daily Writing Habit Path",
		Summary: "A habit-building writing path that turns prompts, short sessions, reflection, and revision into a repeatable daily or weekly writing rhythm.",
		TriggerKeywords: []string{
			"daily writing", "writing habit", "journal habit", "short writing prompt", "creative writing routine",
			"write every day", "reflection writing", "weekly writing review",
		},
		StageRules: []string{
			"Start with a small writing time, topic boundary, and friction-free setup.",
			"Use short prompts and low-pressure drafts before asking for polished writing.",
			"Include reflection and revision only after several raw writing sessions exist.",
			"End with a sustainable routine and a small collection of selected pieces or entries.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can keep a repeatable writing rhythm and review a small body of writing."},
		BadPatterns:             []string{"publishing advice before any writing habit", "literary theory only", "perfect prose as the first lesson"},
	},
	"digital_creation:presentation_publish:portfolio_publish": {
		Title:   "Digital Portfolio Publishing Path",
		Summary: "A publishing path for selecting, polishing, packaging, and releasing a small digital portfolio item such as a video, short, reel, design, or project page.",
		TriggerKeywords: []string{
			"publish portfolio", "upload to youtube", "youtube shorts", "reels", "project page", "thumbnail",
			"portfolio description", "export settings", "digital showcase", "public release checklist",
		},
		StageRules: []string{
			"Start by choosing one publishable work and defining the audience, platform, and success criteria.",
			"Include polish, export, thumbnail or cover, description, and accessibility checks as separate practical stages.",
			"Keep each lesson tied to preparing the actual portfolio item for release.",
			"End with a pre-publish checklist and final package rather than the public reaction itself.",
		},
		CompletionCriteriaRules: []string{"Completion means the work is ready to publish or has been published according to the learner's chosen platform."},
		BadPatterns:             []string{"platform theory only", "analytics growth hacks before a first publish", "unrelated design tool tutorials"},
	},
	"digital_creation:teaching_instruction:creator_tutorial_publish": {
		Title:   "Creator Tutorial Video Path",
		Summary: "A creator path for turning a hobby skill into a clear tutorial video through topic choice, demonstration flow, recording, editing, and publishing prep.",
		TriggerKeywords: []string{
			"make a tutorial video", "how-to content", "teach a hobby skill", "record a lesson", "demo script",
			"screen recording", "voiceover", "edit an instructional video", "creator tutorial",
		},
		StageRules: []string{
			"Start with one narrow skill the creator can demonstrate from start to finish.",
			"Plan the explanation as problem, steps, demonstration, common mistakes, and result.",
			"Include recording setup, clear shots or screen capture, audio check, and basic editing as production stages.",
			"End with a finished tutorial draft, description, thumbnail idea, and upload checklist.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has a clear tutorial video draft or publish-ready tutorial package."},
		BadPatterns:             []string{"general creator branding before one tutorial", "gear shopping as a core lesson", "long lectures with no demonstration"},
	},
	"maker_technical_hobby:certification_assessment:drone_operator_basic_4class": {
		Title:   "Basic Drone Operator Safety Path",
		Summary: "A beginner drone license path focused on safe pre-flight checks, basic controls, hover and landing practice, rule awareness, and exam readiness.",
		TriggerKeywords: []string{
			"beginner drone operator", "basic drone license", "safe drone operation", "pre-flight checklist",
			"hover practice", "takeoff and landing", "drone rules", "drone safety test",
		},
		StageRules: []string{
			"Use concrete lesson titles around pre-flight checks, takeoff and landing, hover control, rule scenarios, and license readiness; avoid generic basics titles.",
			"Start with safety boundaries, local flight rules, battery checks, propeller checks, and site selection.",
			"Move into takeoff, hovering, landing, orientation, and emergency stop practice before advanced maneuvers.",
			"Include short rule and safety scenario checks tied to beginner license preparation.",
			"End with a license-readiness checklist and safe practice routine rather than a public flight event.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can explain basic drone safety rules and complete a safe beginner flight practice checklist."},
		BadPatterns:             []string{"generic understand the basics lesson", "advanced aerial filming before safety checks", "drone shopping as a core lesson"},
	},
	"body_movement:habit_lifestyle:yoga_daily_routine": {
		Title:   "Daily Yoga Routine Path",
		Summary: "A routine-building yoga path that combines breathing, safe posture setup, short flows, reflection, and a repeatable daily practice.",
		TriggerKeywords: []string{
			"beginner yoga routine", "morning yoga", "15-minute yoga", "flexibility routine", "breathing practice",
			"daily yoga habit", "gentle flow", "yoga sequence", "stretching routine",
		},
		StageRules: []string{
			"Use concrete lesson titles around breathing setup, pose sequence, flexibility flow, and weekly repeat plan; avoid generic introduction titles.",
			"Start with breathing, safety boundaries, and a few beginner poses that can be repeated comfortably.",
			"Build a short sequence by connecting warm-up, core stretch, balance or strength, and cool-down.",
			"Include self-checks for comfort, breath, posture, and consistency instead of pushing intensity.",
			"End with a realistic 10- to 15-minute routine plan the learner can repeat for a week.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can follow a short safe yoga routine independently and adjust it to their body condition."},
		BadPatterns:             []string{"generic introduction to yoga lesson", "advanced poses before a safe routine", "anatomy lectures only", "weight-loss promises as the main goal"},
	},
	"body_movement:performance_execution:dance_cover_song_routine": {
		Title:   "Dance Cover Routine Path",
		Summary: "A performance path for learning a song-based dance cover through rhythm, section practice, transitions, expression, and recording rehearsal.",
		TriggerKeywords: []string{
			"k-pop dance cover", "dance cover routine", "chorus choreography", "point choreography", "mirror practice",
			"section practice", "record a dance cover", "clean transitions", "performance rehearsal",
		},
		StageRules: []string{
			"Use concrete lesson titles such as rhythm and counts, point moves, section transitions, mirror review, and recording rehearsal; avoid generic titles like learning the basics.",
			"Choose the target song section and map the beat, counts, and key movement accents first.",
			"Break choreography into intro, verse, chorus, bridge, or point-move sections and practice them separately.",
			"Use mirror or recording review to clean timing, angles, transitions, and expression.",
			"End with a full section run-through or recording rehearsal, not a public upload requirement.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can perform and review a clean cover of the selected dance section."},
		BadPatterns:             []string{"generic learn the basics lesson", "fitness conditioning only", "artist biography as lessons", "full public upload as the final lesson"},
	},
	"craft_making:artifact_creation:crochet_small_doll_project": {
		Title:   "Crochet Amigurumi Doll Path",
		Summary: "A crochet project path for making a small amigurumi doll through basic stitches, round shaping, increases, decreases, stuffing, assembly, and finishing.",
		TriggerKeywords: []string{
			"crochet doll", "amigurumi", "small crochet toy", "single crochet", "magic ring",
			"increase decrease", "stuffing", "attach parts", "embroider face", "finish amigurumi",
		},
		StageRules: []string{
			"Start with yarn, hook, tension, magic ring, and basic stitches needed for the chosen small doll.",
			"Build the body shape through rounds, increases, decreases, and simple counting habits.",
			"Add stuffing, parts, face details, and assembly as project stages rather than separate theory lessons.",
			"End with finishing, shaping, loose-end cleanup, and a completed small crochet doll.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has finished one small crochet doll with stable shape and attached parts."},
		BadPatterns:             []string{"all crochet stitches before any doll shape", "pattern reading theory only", "large complex doll as the first project"},
	},
	"craft_making:artifact_creation:knitting_basic_wearable": {
		Title:   "Beginner Knitted Scarf Path",
		Summary: "A knitting path for completing a simple wearable item by practicing cast-on, knit or purl rhythm, edge control, length consistency, and finishing.",
		TriggerKeywords: []string{
			"beginner knitting", "knit scarf", "scarf knitting", "cast on", "knit stitch", "purl stitch",
			"edge stitches", "bind off", "weave in ends", "simple wearable",
		},
		StageRules: []string{
			"Start with yarn, needles, cast-on, and a small swatch to stabilize tension.",
			"Move into the stitch pattern used by the chosen scarf or wearable item.",
			"Include checkpoints for width, edges, dropped stitches, and consistent length.",
			"End with binding off, weaving in ends, blocking or shaping, and checking the finished wearable.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has completed a simple scarf or wearable knit with stable stitches and clean finishing."},
		BadPatterns:             []string{"complex cable patterns before a first wearable", "yarn shopping as a lesson", "all knitting theory before starting rows"},
	},
	"craft_making:artifact_creation:pottery_handbuilding_cup_project": {
		Title:   "Handbuilt Pottery Cup Path",
		Summary: "A pottery path for making a handbuilt cup or mug through clay prep, forming, handle attachment, surface cleanup, glazing basics, and final review.",
		TriggerKeywords: []string{
			"handbuilding pottery", "ceramic mug", "pinch pot", "slab mug", "coil cup", "attach handle",
			"score and slip", "surface cleanup", "basic glazing", "handbuilt cup",
		},
		StageRules: []string{
			"Start with clay handling, moisture control, workspace setup, and a simple cup form choice.",
			"Build the cup body with pinch, slab, or coil techniques appropriate for beginners.",
			"Include handle attachment, scoring, smoothing, drying awareness, and surface cleanup.",
			"End with glazing preparation or basic decoration and a final inspection checklist.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has shaped and finished a handbuilt cup or mug ready for firing or display according to the available studio process."},
		BadPatterns:             []string{"kiln chemistry as a first lesson", "advanced wheel throwing before handbuilding", "ceramic history only"},
	},
	"cooking_baking:artifact_creation:cooking_basics_foundation": {
		Title:   "Beginner Home Baking Project Path",
		Summary: "A home baking path for completing a simple cake or dessert through ingredient prep, mixing, baking, cooling, decorating, and taste review.",
		TriggerKeywords: []string{
			"home baking", "beginner cake", "sponge cake", "whipped cream cake", "cake batter", "oven temperature",
			"cooling cake", "cake decorating", "simple dessert", "bake a cake",
		},
		StageRules: []string{
			"Start with ingredients, tools, oven setup, and the target dessert structure.",
			"Teach mixing, batter texture, pan preparation, and baking checks through the chosen recipe.",
			"Include cooling, simple cream or topping work, and troubleshooting common beginner mistakes.",
			"End with decorating, serving, tasting notes, and one repeatable recipe card.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has baked and finished one simple dessert with a repeatable process and basic troubleshooting notes."},
		BadPatterns:             []string{"professional pastry theory before one bake", "recipe collection without baking", "business planning for a first dessert"},
	},
	"cooking_baking:habit_lifestyle:meal_prep_weekday_lunch": {
		Title:   "Weekday Lunch Meal Prep Path",
		Summary: "A routine-building cooking path for planning healthy weekday lunches, shopping and prepping ingredients, batch cooking safely, storing meals, and reviewing the next weekly cycle.",
		TriggerKeywords: []string{
			"meal prep", "healthy lunch", "weekday lunches", "lunch prep", "batch cooking",
			"meal storage", "meal prep routine", "weekly menu", "reheat safely",
		},
		StageRules: []string{
			"Start with the learner's weekday lunch goal, time limits, storage needs, and simple balanced menu choices.",
			"Plan shopping and ingredient prep before cooking so the routine can repeat each week.",
			"Include batch cooking, portioning, storage safety, and reheating decisions as practical stages.",
			"End with a review of the first meal prep cycle and a next-week adjustment plan.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has planned, cooked, stored, and reviewed a simple set of weekday lunches."},
		BadPatterns:             []string{"recipe collection only", "diet promises", "nutrition theory without cooking"},
	},
	"knowledge_hobby:habit_lifestyle:garden_design_maintenance_path": {
		Title:   "Small Garden Design And Care Path",
		Summary: "A practical gardening path for planning a small balcony or container garden, choosing plants, setting up soil and light, and maintaining a simple care routine.",
		TriggerKeywords: []string{
			"balcony herb garden", "container garden", "small garden plan", "grow herbs", "plant care routine",
			"watering schedule", "sunlight check", "potting mix", "garden maintenance", "herb care",
		},
		StageRules: []string{
			"Start by assessing sunlight, space, containers, water access, and the learner's target plants.",
			"Plan plant placement, soil or potting mix, drainage, and a manageable planting schedule.",
			"Include watering, pruning, pest observation, and recovery actions as routine practice.",
			"End with a weekly care checklist and a small garden layout the learner can maintain.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can maintain a small herb or container garden with a clear care routine and observation notes."},
		BadPatterns:             []string{"botany taxonomy only", "large farm planning for a balcony goal", "buying plants without a care plan"},
	},
	"knowledge_hobby:habit_lifestyle:houseplant_repotting_care": {
		Title:   "Houseplant Repotting Care Path",
		Summary: "A practical houseplant care path for checking plant condition, choosing a pot and soil mix, repotting safely, monitoring recovery, and setting a simple follow-up routine.",
		TriggerKeywords: []string{
			"houseplant repotting", "repot a plant", "soil mix", "potting mix", "plant recovery",
			"root check", "watering after repotting", "indoor plant care", "monitor recovery",
		},
		StageRules: []string{
			"Start with plant condition, root crowding signs, container choice, drainage, and basic safety.",
			"Choose a simple soil or potting mix and prepare the workspace before removing the plant.",
			"Practice the repotting flow, root handling, settling soil, and cleanup as one connected task.",
			"End with recovery monitoring, watering timing, light placement, and a short care checklist.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can repot one small houseplant safely and monitor the recovery routine for the next week."},
		BadPatterns:             []string{"plant taxonomy only", "shopping for pots only", "rare plant collecting as the main goal"},
	},
	"knowledge_hobby:habit_lifestyle:personal_finance_budget_routine": {
		Title:   "Personal Budget Routine Path",
		Summary: "A beginner personal finance path for listing income and fixed costs, building spending categories, tracking one week of expenses, and reviewing a monthly budget.",
		TriggerKeywords: []string{
			"personal finance", "monthly budget", "budget beginner", "spending categories", "track spending",
			"weekly expenses", "expense review", "budget routine", "household budget",
		},
		StageRules: []string{
			"Start with current income, fixed expenses, variable spending, and a simple tracking method.",
			"Build practical spending categories that match the learner's actual week instead of generic finance theory.",
			"Include one week of expense tracking and a review of where the first budget needs adjustment.",
			"End with a repeatable weekly budget check and next-week adjustment plan.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has built a simple monthly budget and reviewed one week of real spending categories."},
		BadPatterns:             []string{"investment products as the first lesson", "tax advice", "budget app setup only"},
	},
	"knowledge_hobby:habit_lifestyle:dog_basic_training_routine": {
		Title:   "Dog Basic Training Routine Path",
		Summary: "A positive reinforcement dog training path for short practice sessions, reward timing, sit, stay, recall, and combining commands in a safe daily routine.",
		TriggerKeywords: []string{
			"dog basic training", "sit stay come", "recall training", "positive reinforcement",
			"short training session", "reward timing", "dog commands", "puppy training basics",
		},
		StageRules: []string{
			"Start with reward choice, session length, safety boundaries, and consistent cue words.",
			"Teach sit, stay, and come through short positive reinforcement sessions.",
			"Practice commands in slightly different rooms or distractions without forcing long sessions.",
			"End with a short real-life session plan that combines commands and records what to repeat next.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can run a short positive reinforcement session for sit, stay, and come and record what improved."},
		BadPatterns:             []string{"punishment-based training", "aggression diagnosis", "treat shopping only"},
	},
	"visual_art:artifact_creation:structured_drawing_foundation": {
		Title:   "Urban Sketching Foundation Path",
		Summary: "A drawing foundation path for sketching real scenes through observation, line confidence, perspective, composition, value, and a small finished sketch.",
		TriggerKeywords: []string{
			"urban sketching", "street sketch", "beginner drawing", "line confidence", "basic perspective",
			"composition sketch", "observational drawing", "value shading", "scene sketch",
		},
		StageRules: []string{
			"Start with observation, loose line practice, simple shapes, and confidence-building sketch exercises.",
			"Move into perspective, proportion, and composition using real streets, buildings, or interiors.",
			"Include value, focal point, and selective detail so the sketch communicates a scene without overworking it.",
			"End with one small finished urban or observational sketch and a review of what to improve next.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner can complete a simple scene sketch with readable perspective, composition, and confident lines."},
		BadPatterns:             []string{"tool reviews only", "perspective theory without sketching", "copying isolated objects unrelated to the target scene"},
	},
	"digital_creation:artifact_creation:vibe_coding_mvp_app": {
		Title:   "AI-Assisted Web App MVP Path",
		Summary: "A build-and-ship path for creating a small web app MVP with AI coding tools, focusing on scope, core screens, data model, auth, testing, and deployment.",
		TriggerKeywords: []string{
			"web app MVP", "full-stack app", "AI coding tools", "build and deploy", "authentication",
			"database", "CRUD app", "MVP scope", "deploy small app", "vibe coding",
		},
		StageRules: []string{
			"Start by narrowing the MVP problem, core user flow, and smallest shippable feature set.",
			"Plan screens, data model, auth boundary, and API or backend responsibilities before generating code.",
			"Use AI coding tools in short verifyable steps, with manual review, tests, and error fixes after each feature.",
			"End with deployment, smoke testing, and a simple feedback or improvement checklist.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has a small deployed app that supports the core flow with auth or data persistence working at a basic level."},
		BadPatterns:             []string{"AI prompt lists without building", "architecture diagrams only", "adding many features before one working flow"},
	},
	"digital_creation:artifact_creation:notion_study_dashboard": {
		Title:   "Notion Study Dashboard Path",
		Summary: "A practical Notion setup path for building a study dashboard with tasks, notes, weekly review pages, and a repeatable learning workflow.",
		TriggerKeywords: []string{
			"notion study dashboard", "notion dashboard", "study dashboard", "notion workspace",
			"notion productivity dashboard", "productivity dashboard", "notion template",
			"tasks database", "task database", "task tracker", "notes database", "notes page",
			"weekly review", "study planner", "learning dashboard",
		},
		StageRules: []string{
			"Start with the learner's study workflow, task types, note types, and weekly review needs.",
			"Build a simple tasks database and connect it to current courses or study goals.",
			"Create note pages and a weekly review template that are useful before decorative polish.",
			"End by using the dashboard for one weekly review and recording what to improve next.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has a usable Notion study dashboard with tasks, notes, and a weekly review flow."},
		BadPatterns:             []string{"template copying only", "icon decoration only", "advanced automation before a usable dashboard"},
	},
	"digital_creation:artifact_creation:web_dev_project_foundation": {
		Title:   "Full-Stack Web App Foundation Path",
		Summary: "A web development project path for building a small app through page structure, styling, interaction, backend data, authentication boundary, and deployment checks.",
		TriggerKeywords: []string{
			"web app foundation", "full-stack web app", "authentication", "database", "frontend backend",
			"CRUD app", "responsive UI", "API routes", "deploy web app", "small web project",
		},
		StageRules: []string{
			"Start with the core user flow, page map, and data model for one small app.",
			"Build the UI and interaction in thin vertical slices rather than separate tool lectures.",
			"Add backend persistence, authentication boundary, and validation only where the app flow needs them.",
			"End with deployment, smoke testing, responsive checks, and a small backlog of next improvements.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has a deployed small web app with one working user flow and basic data persistence or authentication."},
		BadPatterns:             []string{"HTML/CSS/JS definitions with no app flow", "framework comparison as lessons", "large architecture before a working page"},
	},
	"digital_creation:artifact_creation:animation_fundamentals_shot_progression": {
		Title:   "Beginner Blender Animation Shot Path",
		Summary: "An animation path for creating a short beginner shot by practicing timing, spacing, keyframes, arcs, review, and final polish.",
		TriggerKeywords: []string{
			"beginner Blender animation", "bouncing ball", "keyframes", "timing and spacing", "animation shot",
			"character movement", "graph editor basics", "playblast", "polish animation",
		},
		StageRules: []string{
			"Start with Blender scene setup, keyframes, playback, and a very small motion exercise.",
			"Use a bouncing ball or simple object to practice timing, spacing, arcs, and squash or stretch if appropriate.",
			"Move into a short character or object movement shot with blocking, spline refinement, and camera framing.",
			"End with playblast review, timing adjustments, and a polished short animation export.",
		},
		CompletionCriteriaRules: []string{"Completion means the learner has exported a short animation shot that shows intentional timing, spacing, and readable motion."},
		BadPatterns:             []string{"Blender interface tour only", "modeling lessons unrelated to animation", "complex character rigging before a first shot"},
	},
}

func defaultCurriculumSeedSource() CurriculumSeedSource {
	return CurriculumSeedSource{
		CatalogFile:        "backend/internal/domain/curriculum/curriculum_subpattern_catalog.go",
		RoutingFile:        "backend/internal/domain/curriculum/curriculum_pattern_routing.go",
		RecommendationFile: "backend/internal/domain/curriculum/recommendation_goal_context.go",
	}
}

func splitSubpatternKey(key string) (string, string, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid subpattern key: %s", key)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func tokenizedPatternKeyTerms(key string) []string {
	parts := strings.FieldsFunc(strings.TrimSpace(key), func(r rune) bool {
		return r == ':' || r == '_'
	})
	return dedupeRecommendationTerms(parts)
}

func extraCurriculumPatternSeedTerms(key string) []string {
	switch strings.TrimSpace(key) {
	case "visual_art:artifact_creation:watercolor_beginner_landscape":
		return []string{"수채화", "초보 수채화", "수채화 풍경", "수채화 엽서", "번짐", "워시", "그라데이션", "물 조절", "하늘 그리기", "나무 그리기"}
	case "body_movement:habit_lifestyle:running_5k_beginner_routine":
		return []string{"러닝", "초보 러닝", "5km", "달리기", "조깅", "러닝 루틴", "걷기 달리기", "인터벌 러닝", "페이스 조절", "부상 예방"}
	case "language_communication:performance_execution:daily_english_conversation_routine":
		return []string{"일상 영어", "영어 회화", "생활 영어", "말하기 연습", "쉐도잉", "하루 10문장", "상황별 영어", "자기소개", "짧은 문장", "회화 루틴"}
	case "maker_technical_hobby:artifact_creation:arduino_sensor_project":
		return []string{"아두이노", "센서 프로젝트", "브레드보드", "LED", "배선", "시리얼 모니터", "초음파 센서", "온습도 센서", "서보모터", "전자회로"}
	case "cooking_baking:artifact_creation:korean_home_cooking_basics":
		return []string{"집밥", "한식 기초", "생활요리", "밑반찬", "된장찌개", "김치찌개", "계란말이", "재료 손질", "간 맞추기", "초보 요리"}
	case "body_movement:habit_lifestyle:home_strength_routine":
		return []string{"홈트", "홈 트레이닝", "근력 루틴", "맨몸운동", "스쿼트", "푸쉬업", "코어 운동", "덤벨 운동", "초보 운동", "주간 운동 계획"}
	case "instrument_performance:performance_execution:vocal_song_practice":
		return []string{"보컬 연습", "노래 연습", "발성", "호흡", "음정", "박자", "커버곡", "구간 연습", "한 곡 부르기", "녹음 점검"}
	case "digital_creation:artifact_creation:video_editing_shortform_project":
		return []string{"영상 편집", "숏폼", "쇼츠", "릴스", "캡컷", "프리미어", "컷편집", "자막", "세로 영상", "내보내기"}
	case "digital_creation:technical_skill:ai_tool_productivity_workflow":
		return []string{"AI 생산성", "ChatGPT", "챗지피티", "프롬프트", "업무 자동화", "문서 요약", "구글시트 자동화", "노션 AI", "검토", "반복 업무"}
	case "craft_making:artifact_creation:candle_soap_resin_project":
		return []string{"캔들 만들기", "향초", "비누 만들기", "수제비누", "레진아트", "몰드", "향료", "왁스", "경화", "탈형"}
	case "craft_making:artifact_creation:leathercraft_basic_accessory":
		return []string{
			"가죽공예", "가죽 소품", "키링", "카드 홀더", "가죽 팔찌", "가죽 재단", "타공", "새들스티치", "엣지 마감", "가죽 공예 초보",
		}
	case "digital_creation:artifact_creation:claude_code_agentic_workflow":
		return []string{
			"클로드 코드", "Claude Code", "AI 코딩", "에이전트 코딩", "코딩 워크플로우", "프롬프트 계획", "코드 수정", "테스트 실행", "자동화 작업", "개발 생산성",
		}
	case "digital_creation:artifact_creation:freeware_midi_beatmaking":
		return []string{
			"무료 DAW", "MIDI 비트", "비트메이킹", "드럼 패턴", "베이스라인", "Cakewalk", "LMMS", "루프 만들기", "퀀타이즈", "초보 작곡",
		}
	case "digital_creation:artifact_creation:midi_foundation_workflow":
		return []string{
			"MIDI 기초", "미디 작곡", "DAW 기초", "악기 트랙", "피아노롤", "드럼 패턴", "코드 진행", "퀀타이즈", "루프", "트랙 구조",
		}
	case "digital_creation:artifact_creation:university_studio":
		return []string{
			"디지털 스튜디오", "디지털 작품", "스튜디오 과제", "작품 제작", "디자인 툴", "제작 과정", "결과물 정리", "작품 발표", "포트폴리오", "창작 프로젝트",
		}
	case "digital_creation:artifact_creation:vibe_coding_fullstack_ship":
		return []string{
			"바이브코딩", "풀스택 앱", "웹앱 만들기", "AI 코딩", "프론트엔드", "백엔드", "데이터베이스", "인증", "배포", "서비스 출시",
		}
	case "digital_creation:professional_transition:korean_ncs_web_publisher_training":
		return []string{
			"웹퍼블리셔", "NCS 웹퍼블리셔", "HTML CSS", "반응형 웹", "웹표준", "웹접근성", "퍼블리싱 실무", "레이아웃 구현", "포트폴리오", "취업 준비",
		}
	case "instrument_performance:participation_service:worship_team_support":
		return []string{
			"찬양팀", "예배 반주", "워십 기타", "워십 피아노", "코드 반주", "콘티 연습", "합주 준비", "박자 맞추기", "인트로", "엔딩",
		}
	case "knowledge_hobby:foundation_build:korean_open_lecture_weekly_survey":
		return []string{
			"K-MOOC", "온라인 공개강좌", "공개강좌", "강의 노트", "주차별 학습", "강의 요약", "학습 계획", "토론 정리", "복습", "수강 루틴",
		}
	case "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey":
		return []string{
			"K-MOOC", "온라인 강좌", "수강 계획", "주차별 학습", "강의 듣기", "학습 노트", "퀴즈 복습", "수료", "학습 루틴", "강의 정리",
		}
	case "knowledge_hobby:habit_lifestyle:university_survey":
		return []string{
			"교양 강의", "대학 강의", "인문학 강의", "강의 노트", "개론", "입문 강좌", "주차별 학습", "읽기 자료", "토론", "요약 정리",
		}
	case "knowledge_hobby:professional_transition:credit_bank_practicum":
		return []string{
			"학점은행제", "실습 과목", "현장실습", "실습일지", "평생교육실습", "사회복지현장실습", "실습 준비", "기관 찾기", "실습 기록", "평가서",
		}
	case "knowledge_hobby:professional_transition:credit_bank_standard_theory":
		return []string{
			"학점은행제", "이론 과목", "온라인 수업", "학습자등록", "학점 인정", "과제 작성", "토론 참여", "중간고사", "기말고사", "수료 기준",
		}
	case "knowledge_hobby:teaching_instruction:credit_bank_instructional_design":
		return []string{
			"교육과정 설계", "교수설계", "학습목표", "수업지도안", "평가 설계", "강의안 작성", "학습 활동", "피드백", "온라인 수업", "수업 자료",
		}
	case "maker_technical_hobby:teaching_instruction:maker_workshop_demo":
		return []string{
			"메이커 워크숍", "DIY 시연", "제작 시연", "아두이노", "전자회로", "프로토타입", "재료 준비", "안전 안내", "단계 설명", "작품 발표",
		}
	case "visual_art:artifact_creation:calligraphy_basic_lettering":
		return []string{
			"캘리그라피", "손글씨", "붓펜", "레터링", "획 연습", "자음 모음", "글자 균형", "문장 쓰기", "엽서 만들기", "작품 만들기",
		}
	case "visual_art:artifact_creation:university_studio":
		return []string{
			"미술 스튜디오", "기초 조형", "드로잉", "색채", "구도", "작품 제작", "크리틱", "스케치", "재료 실험", "작품 발표",
		}
	case "visual_art:presentation_publish:art_publish_showcase":
		return []string{
			"작품 공개", "온라인 전시", "작품 소개", "아트 포트폴리오", "작품 사진", "작품 설명", "전시 준비", "작품 업로드", "작가노트", "쇼케이스",
		}
	case "visual_art:presentation_publish:calligraphy_publish_showcase":
		return []string{
			"캘리그라피 작품 공개", "손글씨 작품", "작품 사진", "문구 설명", "작품 소개", "온라인 전시", "엽서 작품", "작가노트", "게시 준비", "쇼케이스",
		}
	case "writing_storytelling:habit_lifestyle:daily_writing_habit":
		return []string{
			"매일 글쓰기", "글쓰기 습관", "짧은 글쓰기", "글감", "일기 쓰기", "문장 연습", "초안 쓰기", "퇴고", "주간 회고", "글 모음",
		}
	case "writing_storytelling:presentation_publish:university_studio":
		return []string{
			"글쓰기 워크숍", "창작 글쓰기", "에세이", "소설 초안", "합평", "퇴고", "작품 발표", "문장 다듬기", "글 모음", "출간 준비",
		}
	case "cooking_baking:habit_lifestyle:meal_prep_weekday_lunch":
		return []string{"밀프렙", "도시락 준비", "주중 도시락", "건강 도시락", "일주일 식단", "반찬 준비", "식재료 보관", "소분", "점심 도시락", "식단 루틴"}
	case "digital_creation:artifact_creation:adobe_creative_learning_path", "digital_creation:artifact_creation:creative_tool_beginner_output":
		return []string{"포토샵", "일러스트레이터", "어도비", "디자인 툴", "이미지 편집", "벡터 드로잉", "레이어", "마스크", "초보 디자인", "작품 만들기"}
	case "digital_creation:artifact_creation:animation_fundamentals_shot_progression":
		return []string{"블렌더 애니메이션", "애니메이션 기초", "바운싱볼", "키프레임", "타이밍", "스페이싱", "그래프 에디터", "짧은 애니메이션", "플레이블라스트", "동작 다듬기"}
	case "digital_creation:artifact_creation:creative_coding_visual_project":
		return []string{"크리에이티브 코딩", "자바스크립트 드로잉", "p5.js", "인터랙션", "캔버스 애니메이션", "제너레이티브 아트", "도형 그리기", "마우스 반응", "비주얼 프로젝트", "웹 아트"}
	case "digital_creation:artifact_creation:game_3d_guided_pathway":
		return []string{"유니티", "3D 게임", "게임 프로토타입", "실시간 3D", "스크립트", "플레이어 이동", "충돌 처리", "카메라", "씬 구성", "게임 만들기"}
	case "digital_creation:artifact_creation:interactive_music_production_foundation", "digital_creation:artifact_creation:music_production_full_pipeline":
		return []string{"음악 제작", "비트메이킹", "사운드 디자인", "편곡", "믹싱", "마스터링", "루프", "송 구조", "DAW", "트랙 완성"}
	case "digital_creation:artifact_creation:template_design_quick_output":
		return []string{"캔바", "템플릿 디자인", "SNS 디자인", "브랜드 키트", "인스타그램 카드뉴스", "썸네일", "디자인 템플릿", "색상 조합", "폰트", "빠른 결과물"}
	case "digital_creation:foundation_build:mooc_programming_foundation", "digital_creation:foundation_build:problem_set_to_final_project":
		return []string{"프로그래밍 기초", "컴퓨터과학", "파이썬", "문제풀이", "과제 풀이", "최종 프로젝트", "알고리즘 기초", "코딩 연습", "프로젝트 설계", "학습 노트"}
	case "digital_creation:foundation_build:web_platform_core_curriculum":
		return []string{"웹 플랫폼", "HTML", "CSS", "JavaScript", "DOM", "브라우저", "웹표준", "반응형", "접근성", "기초 웹개발"}
	case "digital_creation:professional_transition:cloud_lab_skill_path":
		return []string{"클라우드 실습", "Google Cloud", "핸즈온랩", "스킬배지", "클라우드 기초", "가상머신", "스토리지", "IAM", "실습 환경", "클라우드 학습"}
	case "digital_creation:professional_transition:frontend_competency_map":
		return []string{"프론트엔드 역량", "웹 접근성", "반응형 디자인", "HTML CSS", "JavaScript", "컴포넌트", "성능 최적화", "포트폴리오", "역량 체크", "실무 준비"}
	case "digital_creation:professional_transition:learning_path_role_skill":
		return []string{"직무 학습 로드맵", "역할 기반 학습", "스킬맵", "실무 과제", "학습 모듈", "역량 개발", "업무 시나리오", "포트폴리오", "전환 준비", "학습 계획"}
	case "digital_creation:technical_skill:computational_music_analysis_project":
		return []string{"음악 분석", "music21", "악보 분석", "코퍼스 분석", "컴퓨터 음악이론", "파이썬 음악", "상징 악보", "화성 분석", "데이터 분석", "분석 프로젝트"}
	case "knowledge_hobby:concept_mastery:climate_data_graphing_explanation":
		return []string{"기후 데이터", "기온 편차", "그래프 그리기", "지구 온도", "데이터 해석", "선그래프", "시각화", "기후 변화", "활동지", "설명하기"}
	case "knowledge_hobby:concept_mastery:guided_simulation_activity_sheet", "knowledge_hobby:concept_mastery:science_simulation_inquiry_path":
		return []string{"PhET", "과학 시뮬레이션", "탐구 활동", "활동지", "변인 실험", "가설 세우기", "관찰 기록", "실험 결과", "개념 학습", "수업 활동"}
	case "knowledge_hobby:concept_mastery:primary_source_inquiry_note":
		return []string{"사료 분석", "역사 탐구", "1차 자료", "관찰 질문", "출처 분석", "근거 찾기", "탐구 노트", "자료 읽기", "역사 수업", "질문 만들기"}
	case "knowledge_hobby:foundation_build:audit_course_practice_path", "knowledge_hobby:foundation_build:mooc_humanities_survey", "knowledge_hobby:foundation_build:short_course_weekly_discussion":
		return []string{"온라인 강좌", "강의 듣기", "읽기 자료", "토론", "주차별 학습", "강의 노트", "복습", "요약", "학습 계획", "수강 루틴"}
	case "knowledge_hobby:foundation_build:interactive_tutor_concept_practice":
		return []string{"개념 연습", "인터랙티브 튜터", "생물학 기초", "피드백", "문제 연습", "개념 확인", "학습 활동", "오답 점검", "짧은 퀴즈", "복습"}
	case "knowledge_hobby:foundation_build:oer_learning_object_cluster", "knowledge_hobby:foundation_build:open_textbook_chapter_practice":
		return []string{"공개 교육자료", "오픈 교재", "학습 객체", "챕터 연습", "스터디 가이드", "활동지", "개념 정리", "시뮬레이션", "연습문제", "학습 노트"}
	case "knowledge_hobby:habit_lifestyle:citizen_science_observation_record":
		return []string{"시민과학", "iNaturalist", "생물 관찰", "종 동정", "관찰 기록", "BioBlitz", "야외 관찰", "사진 기록", "자연 탐사", "관찰 노트"}
	case "knowledge_hobby:habit_lifestyle:digital_literacy_badged_course":
		return []string{"디지털 리터러시", "온라인 안전", "디지털 역량", "정보 과부하", "스마트폰 활용", "인터넷 사용", "개인정보 보호", "디지털 배지", "생활 디지털", "기초 교육"}
	case "knowledge_hobby:habit_lifestyle:dog_basic_training_routine":
		return []string{"강아지 훈련", "앉아 기다려 이리와", "긍정 강화", "리콜 훈련", "보상 타이밍", "짧은 훈련", "반려견 교육", "강아지 명령어", "훈련 루틴", "퍼피 트레이닝"}
	case "knowledge_hobby:habit_lifestyle:houseplant_repotting_care":
		return []string{"반려식물 분갈이", "흙 배합", "식물 회복", "분갈이 후 물주기", "뿌리 정리", "화분 배수", "식물 관리", "분갈이 방법", "식물 상태", "관리 루틴"}
	case "knowledge_hobby:habit_lifestyle:organic_growing_cycle_plan":
		return []string{"유기농 텃밭", "작물 재배", "파종", "재배 주기", "계절 작물", "텃밭 계획", "모종", "물주기", "수확", "재배 기록"}
	case "knowledge_hobby:professional_transition:completion_refund_online_course", "knowledge_hobby:professional_transition:short_diploma_assessment_path":
		return []string{"온라인 과정", "진도 관리", "수료", "환급 과정", "과제 제출", "평가 준비", "수료 기준", "학습 일정", "직무 역량", "증명서"}
	case "knowledge_hobby:professional_transition:midlife_transition_learning_support":
		return []string{"중장년 직업역량", "직업 전환", "취창업", "평생학습", "역량 개발", "재취업 준비", "진로 탐색", "학습 상담", "이력서 준비", "전환 계획"}
	case "knowledge_hobby:professional_transition:rural_return_agriculture_foundation":
		return []string{"귀농귀촌", "농업 기초", "품목기술", "정착 준비", "농업 교육", "작물 선택", "농촌 생활", "영농 계획", "귀농 준비", "재배 기초"}
	case "language_communication:certification_assessment:english_integrated_four_skills_certification":
		return []string{"아이엘츠", "토플", "영어 4技能", "리딩", "리스닝", "스피킹", "라이팅", "문제풀이", "시험 준비", "모의고사"}
	case "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path":
		return []string{"종이 회로", "카드보드 회로", "전자회로", "로봇 만들기", "프로토타입", "LED 회로", "구리테이프", "메이커 프로젝트", "발명", "회로 실습"}
	case "maker_technical_hobby:artifact_creation:design_lab_prototype_project", "maker_technical_hobby:artifact_creation:maker_project_class":
		return []string{"메이커 프로젝트", "프로토타입", "제작 실습", "디자인 프로세스", "전자회로", "목공", "DIY", "문제 해결", "작품 만들기", "메이커 수업"}
	case "maker_technical_hobby:artifact_creation:design_thinking_problem_framing_path", "maker_technical_hobby:artifact_creation:no_build_prototype_test":
		return []string{"디자인씽킹", "문제 정의", "공감", "프로토타입 테스트", "사용자 피드백", "아이디어 검증", "노빌드 프로토타입", "가설 검증", "인터뷰", "실험 계획"}
	case "visual_art:artifact_creation:photography_studio_project":
		return []string{"사진 스튜디오", "스튜디오 조명", "디지털 이미징", "촬영 세팅", "인물 사진", "조명 배치", "배경지", "사진 크리틱", "보정", "촬영 프로젝트"}
	case "writing_storytelling:presentation_publish:critical_reading_to_writing_draft":
		return []string{"비평적 읽기", "창작 글쓰기", "글 초안", "소설 쓰기", "읽고 쓰기", "인물 분석", "구조 분석", "초안 작성", "퇴고", "글쓰기 발표"}
	case "digital_creation:presentation_publish:portfolio_publish":
		return []string{
			"영상편집",
			"영상 편집",
			"유튜브",
			"youtube",
			"업로드",
			"게시",
			"공개",
			"릴스",
			"shorts",
			"포트폴리오",
			"디지털 결과물",
			"콘텐츠 공개",
			"게시 설명",
			"썸네일",
		}
	case "digital_creation:teaching_instruction:creator_tutorial_publish":
		return []string{
			"강의 영상",
			"짧은 강의 영상",
			"튜토리얼 영상",
			"설명 영상",
			"시연 영상",
			"교육 영상",
			"유튜브 강의",
			"내가 아는 취미 기술",
			"취미 기술 설명",
			"기술을 설명하는 영상",
			"강의 영상을 촬영",
			"강의 영상을 촬영하고 편집",
			"영상으로 설명",
			"교육 콘텐츠",
			"학습 영상",
			"촬영",
			"촬영 구성",
			"편집",
			"편집해 게시",
			"게시",
			"설명 흐름",
			"설명 스크립트",
			"제작 예시",
			"creator tutorial",
			"tutorial video",
			"how-to video",
		}
	case "cooking_baking:artifact_creation:cooking_basics_foundation":
		return []string{
			"홈베이킹",
			"케이크",
			"생크림 케이크",
			"제누와즈",
			"아이싱",
			"크림",
			"디저트",
			"반죽",
			"오븐",
			"굽기",
			"완성",
			"초급 dish",
		}
	case "craft_making:artifact_creation:leathercraft_wallet_project":
		return []string{
			"카드지갑",
			"카드 지갑",
			"반지갑",
			"지갑",
			"wallet",
			"card wallet",
			"wallet project",
			"수납칸",
			"카드 슬롯",
			"외피",
			"내피",
			"패턴",
			"재단",
			"타공",
			"새들스티치",
			"saddle stitch",
			"엣지 마감",
			"edge finishing",
		}
	case "craft_making:artifact_creation:leathercraft_dimensional_bag":
		return []string{
			"미니가방",
			"미니 가방",
			"크로스백",
			"파우치",
			"가방",
			"bag",
			"crossbody bag",
			"pouch",
			"gusset",
			"거싯",
			"스트랩",
			"strap",
			"안감",
			"포켓",
			"여밈",
			"하드웨어",
			"입체 구조",
			"다층 조립",
			"dimensional assembly",
		}
	case "visual_art:artifact_creation:calligraphy_quote_art_project":
		return []string{
			"문장 작품",
			"문구 작품",
			"짧은 문장",
			"문장 캘리그라피",
			"캘리그라피 문장",
			"명언",
			"quote art",
			"quote project",
			"캘리그라피 작품",
			"붓펜 획",
			"획과 글자 균형",
			"작품 구성",
			"글자 균형",
			"여백",
			"배치",
			"문구 선택",
			"감정 전달",
			"구도 안정",
			"composition",
			"카드",
			"엽서",
			"액자",
			"frame",
			"postcard",
			"선물",
			"마감",
		}
	case "visual_art:artifact_creation:photography_fundamentals_project":
		return []string{
			"사진",
			"사진 기초",
			"스마트폰 사진",
			"휴대폰 사진",
			"폰 사진",
			"일상 사진",
			"인물 사진",
			"사진 잘 찍기",
			"촬영 구도",
			"사진 구도",
			"빛 활용",
			"자연광",
			"초점",
			"노출",
			"보정",
			"사진 보정",
			"배경 정리",
			"포즈 유도",
			"작은 사진 프로젝트",
			"mobile photography",
			"portrait photography",
			"composition",
			"lighting",
			"focus",
		}
	case "cooking_baking:habit_lifestyle:food_preservation_safety_path":
		return []string{
			"김치 담그기",
			"김치",
			"발효 관리",
			"발효 보관",
			"배추 절이기",
			"양념 비율",
			"식품 보존",
			"저장 발효",
			"소량 보존",
			"위생 기준",
			"보관 온도",
			"home food preservation",
			"fermentation safety",
		}
	case "knowledge_hobby:concept_mastery:horticulture_diagnostic_foundation":
		return []string{
			"반려식물",
			"분갈이",
			"흙 배합",
			"물주기 기준",
			"식물 진단",
			"식물 관리",
			"화분 관리",
			"뿌리 상태",
			"배수",
			"병해충 관찰",
			"care plan",
			"plant diagnostics",
		}
	case "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path":
		return []string{
			"베란다 허브",
			"허브 키우기",
			"허브 화분",
			"베란다 화분",
			"햇빛과 물주기",
			"물주기 루틴",
			"작은 정원",
			"컨테이너 정원",
			"화분 가꾸기",
			"분갈이",
			"허브 관리",
			"balcony herb garden",
			"plant care routine",
		}
	case "instrument_performance:performance_execution:song_driven_instrument_path":
		return []string{
			"곡 카피",
			"원곡 카피",
			"원곡에 맞춰",
			"원곡 느낌",
			"커버",
			"좋아하는 곡 커버",
			"좋아하는 팝송",
			"팝송 반주",
			"피아노 반주",
			"코드 반주",
			"왼손 반주",
			"밴드곡",
			"리프",
			"솔로",
			"베이스 라인",
			"베이스라인",
			"그루브",
			"필인",
			"반주 패턴",
			"코드 진행",
			"반주 코드 진행",
			"구간별 카피",
			"song cover",
			"cover song",
			"song-driven instrument practice",
		}
	case "visual_art:artifact_creation:sequential_creative_class_path":
		return []string{
			"창작 클래스",
			"일러스트",
			"포트폴리오",
			"작품 프로젝트",
			"연속 프로젝트",
			"갤러리",
			"작품 완성",
			"creative class",
			"illustration",
			"portfolio",
			"project",
		}
	case "body_movement:performance_execution:dance_cover_song_routine":
		return []string{
			"케이팝 댄스",
			"kpop dance",
			"k-pop dance",
			"댄스 커버",
			"dance cover",
			"안무 커버",
			"안무 한 곡",
			"포인트 안무",
			"후렴 안무",
			"동작 구간",
			"구간 연습",
			"리듬과 동작",
			"음악 박자",
			"전체 루틴",
			"리허설",
		}
	case "body_movement:habit_lifestyle:yoga_daily_routine":
		return []string{
			"아침 요가",
			"요가 루틴",
			"데일리 요가",
			"매일 요가",
			"15분 요가",
			"10분 요가",
			"요가 자세",
			"요가 호흡",
			"기본 호흡",
			"스트레칭 자세",
			"자세 전환",
			"짧은 루틴",
			"반복 가능한 루틴",
			"routine yoga",
			"daily yoga",
		}
	case "craft_making:artifact_creation:woodworking_small_project":
		return []string{
			"목공예",
			"목공",
			"woodworking",
			"woodcraft",
			"작은 선반",
			"벽 선반",
			"선반 만들기",
			"도마 만들기",
			"주방용 도마",
			"목재 선택",
			"목재 재단",
			"나무 재단",
			"치수 계획",
			"샌딩",
			"사포질",
			"모서리 다듬기",
			"오일 마감",
			"목재 마감",
			"조립",
			"wooden shelf",
			"cutting board",
			"wood finishing",
		}
	case "craft_making:artifact_creation:knitting_basic_wearable":
		return []string{
			"대바늘",
			"뜨개",
			"뜨개질",
			"knitting",
			"knit",
			"목도리",
			"스카프",
			"니트 소품",
			"착용 소품",
			"코 잡기",
			"겉뜨기",
			"안뜨기",
			"단수",
			"폭 유지",
			"반복 무늬",
			"실 정리",
			"마감",
			"scarf knitting",
			"wearable knit project",
		}
	case "craft_making:artifact_creation:crochet_small_doll_project":
		return []string{
			"코바늘",
			"crochet",
			"amigurumi",
			"아미구루미",
			"코바늘 인형",
			"인형 만들기",
			"작은 인형",
			"원형뜨기",
			"짧은뜨기",
			"증감",
			"솜 넣기",
			"부품 연결",
			"몸통",
			"얼굴 마감",
			"형태 잡기",
			"crochet doll",
			"small crochet toy",
		}
	case "craft_making:artifact_creation:pottery_handbuilding_cup_project":
		return []string{
			"도자기",
			"pottery",
			"ceramic",
			"ceramics",
			"머그컵",
			"컵 만들기",
			"도자기 컵",
			"핸드빌딩",
			"흙 성형",
			"기본 성형",
			"손잡이 붙이기",
			"손잡이 접합",
			"표면 정리",
			"유약",
			"소성",
			"mug project",
			"handbuilding cup",
		}
	case "craft_making:artifact_creation:sewing_project_foundation":
		return []string{
			"재봉틀",
			"재봉",
			"봉제",
			"sewing machine",
			"에코백",
			"eco bag",
			"파우치",
			"쿠션커버",
			"원단 재단",
			"fabric cutting",
			"직선 박기",
			"박음질",
			"seam",
			"hem",
			"stitch",
			"초급 봉제",
			"봉제 결과물",
		}
	case "cooking_baking:artifact_creation:home_cafe_latte_foundation":
		return []string{
			"홈카페",
			"home cafe",
			"라떼",
			"latte",
			"에스프레소 추출",
			"espresso",
			"우유 스티밍",
			"milk steaming",
			"우유 질감",
			"라떼 비율",
			"맛 균형",
			"반복 레시피",
			"커피",
			"coffee",
			"집에서 라떼",
		}
	default:
		return nil
	}
}

func extraEnglishCurriculumPatternSeedTerms(key string) []string {
	base := tokenizedPatternKeyTerms(key)
	switch strings.TrimSpace(key) {
	case "visual_art:artifact_creation:watercolor_beginner_landscape":
		return append(base, "watercolor landscape", "beginner watercolor", "watercolor postcard", "wash technique", "gradient wash", "water control", "sky painting", "tree painting")
	case "body_movement:habit_lifestyle:running_5k_beginner_routine":
		return append(base, "beginner running", "5K", "couch to 5k", "walk run", "running form", "interval running", "pace", "recovery")
	case "language_communication:performance_execution:daily_english_conversation_routine":
		return append(base, "daily English", "English conversation", "everyday phrases", "shadowing", "speaking practice", "role play", "daily routine")
	case "maker_technical_hobby:artifact_creation:arduino_sensor_project":
		return append(base, "Arduino", "sensor project", "breadboard", "LED", "wiring", "serial monitor", "servo", "beginner project")
	case "cooking_baking:artifact_creation:korean_home_cooking_basics":
		return append(base, "Korean home cooking", "banchan", "doenjang jjigae", "kimchi jjigae", "egg roll", "seasoning", "beginner recipe")
	case "body_movement:habit_lifestyle:home_strength_routine":
		return append(base, "home workout", "home strength", "bodyweight workout", "strength routine", "squat", "push-up", "core workout", "dumbbell workout", "weekly workout plan")
	case "instrument_performance:performance_execution:vocal_song_practice":
		return append(base, "vocal practice", "singing practice", "breath support", "pitch practice", "rhythm", "cover song", "full song", "recording check")
	case "digital_creation:artifact_creation:video_editing_shortform_project":
		return append(base, "video editing", "shortform video", "youtube shorts", "reels", "capcut", "premiere pro", "cut editing", "subtitles", "vertical video", "export")
	case "digital_creation:technical_skill:ai_tool_productivity_workflow":
		return append(base, "AI productivity", "ChatGPT workflow", "prompt workflow", "document summary", "spreadsheet automation", "notion ai", "work automation", "review workflow")
	case "craft_making:artifact_creation:candle_soap_resin_project":
		return append(base, "candle making", "soy candle", "soap making", "handmade soap", "resin art", "epoxy resin", "mold", "fragrance oil", "curing")
	case "visual_art:artifact_creation:photography_fundamentals_project":
		return append(base,
			"photography fundamentals", "mobile photography", "smartphone photography", "portrait photography",
			"composition", "lighting", "natural light", "focus", "exposure", "photo editing", "small photo project",
		)
	case "craft_making:artifact_creation:sewing_project_foundation":
		return append(base,
			"sewing machine", "beginner sewing", "eco bag", "pouch", "fabric cutting", "straight stitch", "seam", "hem", "stitch",
		)
	case "craft_making:artifact_creation:woodworking_small_project":
		return append(base,
			"woodworking", "small shelf", "cutting board", "wood selection", "measuring", "sanding", "wood finishing", "oil finish",
		)
	case "cooking_baking:artifact_creation:home_cafe_latte_foundation":
		return append(base,
			"home cafe", "latte", "espresso", "milk steaming", "milk texture", "coffee", "repeatable recipe",
		)
	case "instrument_performance:performance_execution:song_driven_instrument_path":
		return append(base,
			"song cover", "favorite song", "riff", "solo", "bass line", "groove", "fill", "accompaniment", "accompaniment pattern",
			"piano accompaniment", "chord accompaniment", "left hand accompaniment", "chord progression",
		)
	case "instrument_performance:performance_execution:song_completion":
		return append(base,
			"complete a song", "perform one song", "play through", "practice sections", "rehearsal",
		)
	case "digital_creation:presentation_publish:portfolio_publish":
		return append(base,
			"video editing", "youtube", "upload", "publish", "portfolio", "thumbnail", "shorts", "reels", "digital work",
		)
	case "digital_creation:teaching_instruction:creator_tutorial_publish":
		return append(base,
			"tutorial video", "how-to video", "teaching video", "explain a hobby skill", "record and edit a lesson", "instructional content",
		)
	default:
		return base
	}
}

func titleizePatternKey(key string) string {
	parts := strings.Split(strings.TrimSpace(key), ":")
	if len(parts) == 0 {
		return "Curriculum Pattern"
	}
	last := strings.ReplaceAll(parts[len(parts)-1], "_", " ")
	words := strings.Fields(last)
	for idx, word := range words {
		if len(word) == 0 {
			continue
		}
		words[idx] = strings.ToUpper(word[:1]) + word[1:]
	}
	title := strings.Join(words, " ")
	if title == "" {
		return "Curriculum Pattern"
	}
	return title
}

func buildEnglishSeedSummary(key, domain, goalType string, spec goalSubpatternSpec) string {
	flow := strings.Join(describeStepRolesEnglish(spec.StepRoles), " -> ")
	if flow == "" {
		flow = "setup -> core practice -> integrated output"
	}
	return fmt.Sprintf("%s pattern for %s/%s goals. Use a gateway progression: %s.", titleizePatternKey(key), strings.ReplaceAll(domain, "_", " "), strings.ReplaceAll(goalType, "_", " "), flow)
}

func buildEnglishStageRules(_ string, _ string, goalType string, spec goalSubpatternSpec) []string {
	flow := strings.Join(describeStepRolesEnglish(spec.StepRoles), " -> ")
	if flow == "" {
		flow = "setup -> core practice -> integration -> final preparation"
	}
	rules := []string{
		"Design lessons as learner-facing gateway stages, not as a copied list of topics, resources, or definitions.",
		"Recommended progression: " + flow + ".",
		"Each lesson should end with a practical ability, draft, output, rehearsal, or decision that moves the learner closer to the confirmed goal.",
	}
	if goalType == "artifact_creation" {
		rules = append(rules, "Keep the final lesson focused on completing and checking a small concrete artifact or project.")
	}
	if goalType == "performance_execution" {
		rules = append(rules, "Keep the final lesson focused on integrated practice or rehearsal before the target performance.")
	}
	return rules
}

func buildEnglishLastLessonRule(goalType string) string {
	switch goalType {
	case "certification_assessment":
		return "The last lesson should prepare for the exam or assessment; the actual exam result belongs in completion criteria."
	case "presentation_publish":
		return "The last lesson should prepare the work for publishing or presentation; the actual publishing event belongs in completion criteria."
	case "performance_execution":
		return "The last lesson should be integrated practice or rehearsal before the target performance."
	default:
		return "The last lesson should integrate the key skills into a small final output or ready-to-use practice result."
	}
}

func buildEnglishCompletionCriteriaRules(goalType string) []string {
	switch goalType {
	case "certification_assessment":
		return []string{"Keep taking the official exam, passing, or receiving a certificate in completion criteria, not as a lesson title."}
	case "presentation_publish":
		return []string{"Keep the real publishing, exhibition, upload, or public release in completion criteria, not as a lesson title."}
	case "professional_transition":
		return []string{"Keep paid work, job change, client acquisition, or formal role transition in completion criteria."}
	default:
		return []string{"Keep external final events outside the lesson list and express them as completion criteria."}
	}
}

func buildEnglishBadPatterns() []string {
	return []string{
		"one tiny setup detail as a standalone lesson",
		"a single glossary term as a lesson",
		"the final event itself as the final lesson",
		"a copied list of resource chapter titles",
	}
}

func buildEnglishCurriculumPatternEmbeddingText(key, domain, goalType, title, summary string, triggerKeywords, stepRoles, stageRules, completionRules, badPatterns []string) string {
	segments := []string{
		"pattern_key: " + key,
		"language: en",
		"domain: " + domain,
		"goal_type: " + goalType,
		"title: " + title,
	}
	if len(triggerKeywords) > 0 {
		segments = append(segments, "trigger_keywords: "+strings.Join(triggerKeywords, ", "))
	}
	if len(stepRoles) > 0 {
		segments = append(segments, "recommended_sequence: "+strings.Join(describeStepRolesEnglish(stepRoles), " -> "))
	}
	if strings.TrimSpace(summary) != "" {
		segments = append(segments, "summary: "+strings.TrimSpace(summary))
	}
	if len(stageRules) > 0 {
		segments = append(segments, "stage_rules: "+strings.Join(stageRules, " "))
	}
	if len(completionRules) > 0 {
		segments = append(segments, "completion_criteria_rules: "+strings.Join(completionRules, " "))
	}
	if len(badPatterns) > 0 {
		segments = append(segments, "bad_patterns: "+strings.Join(badPatterns, ", "))
	}
	return strings.Join(segments, "\n")
}

func buildSeedSummary(spec goalSubpatternSpec) string {
	parts := []string{}
	if strings.TrimSpace(spec.Name) != "" {
		parts = append(parts, strings.TrimSpace(spec.Name))
	}
	if len(spec.StageRules) > 0 {
		parts = append(parts, strings.TrimSpace(spec.StageRules[0]))
	}
	if strings.TrimSpace(spec.LastLessonRule) != "" {
		parts = append(parts, strings.TrimSpace(spec.LastLessonRule))
	}
	return strings.Join(parts, " ")
}

func buildCurriculumPatternEmbeddingText(key, domain, goalType string, spec goalSubpatternSpec, triggerKeywords []string, summary string) string {
	segments := []string{
		"pattern_key: " + key,
		"domain: " + domain,
		"goal_type: " + goalType,
		"title: " + strings.TrimSpace(spec.Name),
	}
	if len(triggerKeywords) > 0 {
		segments = append(segments, "trigger_keywords: "+strings.Join(triggerKeywords, ", "))
	}
	if len(spec.StepRoles) > 0 {
		segments = append(segments, "recommended_sequence: "+strings.Join(spec.StepRoles, " -> "))
	}
	if strings.TrimSpace(summary) != "" {
		segments = append(segments, "summary: "+strings.TrimSpace(summary))
	}
	if len(spec.StageRules) > 0 {
		segments = append(segments, "stage_rules: "+strings.Join(spec.StageRules, " "))
	}
	if strings.TrimSpace(spec.LastLessonRule) != "" {
		segments = append(segments, "last_lesson_rule: "+strings.TrimSpace(spec.LastLessonRule))
	}
	if len(spec.CompletionCriteriaRules) > 0 {
		segments = append(segments, "completion_criteria_rules: "+strings.Join(spec.CompletionCriteriaRules, " "))
	}
	if len(spec.AtomicDisallowExamples) > 0 {
		segments = append(segments, "bad_patterns: "+strings.Join(spec.AtomicDisallowExamples, ", "))
	}
	return strings.Join(segments, "\n")
}
