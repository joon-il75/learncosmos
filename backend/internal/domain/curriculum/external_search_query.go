package curriculum

import (
	"strings"

	"github.com/learnweaver/backend/internal/pkg/extsearch"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func externalProviderSearchQuery(provider extsearch.SearchProvider, query string) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
	if trimmed == "" {
		return trimmed
	}
	switch provider.(type) {
	case *extsearch.YouTubeProvider:
		return buildYouTubeExternalSearchQuery(trimmed)
	default:
		return trimmed
	}
}

func buildYouTubeExternalSearchQuery(query string) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
	if trimmed == "" {
		return trimmed
	}
	lowered := strings.ToLower(trimmed)
	if strings.Contains(lowered, "유튜브") || strings.Contains(lowered, "youtube") {
		return trimmed
	}

	bundle := normalizer.BuildLessonSearchQuery(trimmed, "", "", "", "", "")
	tokenSet := make(map[string]struct{}, len(bundle.Tokens)+len(bundle.ExpandedTokens))
	for _, token := range append(bundle.Tokens, bundle.ExpandedTokens...) {
		normalized := strings.ToLower(strings.TrimSpace(token))
		if normalized != "" {
			tokenSet[normalized] = struct{}{}
		}
	}
	hasToken := func(tokens ...string) bool {
		for _, token := range tokens {
			if _, ok := tokenSet[token]; ok {
				return true
			}
			if strings.Contains(lowered, strings.ToLower(token)) {
				return true
			}
		}
		return false
	}

	if hasToken("온라인") && hasToken("강의") {
		switch {
		case hasToken("주제") && hasToken("선정"):
			return "온라인 강의 만들기 주제 선정"
		case hasToken("콘텐츠") && hasToken("구성"):
			return "온라인 강의 콘텐츠 구성 방법"
		case hasToken("촬영") || hasToken("편집"):
			return "온라인 강의 촬영 편집"
		case hasToken("홍보") || hasToken("광고") || hasToken("마케팅"):
			return "온라인 강의 마케팅 홍보"
		case hasToken("플랫폼") || hasToken("게시"):
			return "온라인 강의 플랫폼 게시"
		default:
			return "온라인 강의 만들기"
		}
	}

	if !containsHangul(trimmed) {
		if strings.Contains(lowered, "video") || strings.Contains(lowered, "tutorial") || strings.Contains(lowered, "guide") {
			return trimmed
		}
		return trimmed + " tutorial"
	}
	if strings.Contains(lowered, "영상") {
		return trimmed
	}
	return trimmed + " 영상"
}

func containsHangul(value string) bool {
	for _, ch := range value {
		if ch >= '가' && ch <= '힣' {
			return true
		}
	}
	return false
}

func BuildExternalSearchQuery(queryBundle normalizer.LessonSearchQuery) string {
	intent := detectExternalSearchIntent(queryBundle)
	switch intent {
	case externalSearchIntentConceptIntro:
		return buildConceptIntroExternalSearchQuery(queryBundle)
	case externalSearchIntentToolIntro:
		return buildToolIntroExternalSearchQuery(queryBundle)
	case externalSearchIntentHowTo:
		return buildHowToExternalSearchQuery(queryBundle)
	case externalSearchIntentComparison:
		return buildComparisonExternalSearchQuery(queryBundle)
	}

	parts := make([]string, 0, 6)
	if trimmed := strings.TrimSpace(queryBundle.RawQuery); trimmed != "" {
		parts = append(parts, trimmed)
	}
	for idx, token := range queryBundle.Tokens {
		if idx >= 4 {
			break
		}
		parts = append(parts, token)
	}
	query := strings.Join(parts, " ")
	query = strings.Join(strings.Fields(query), " ")
	runes := []rune(query)
	if len(runes) > 120 {
		return string(runes[:120])
	}
	return query
}

type externalSearchIntent string

const (
	externalSearchIntentGeneral      externalSearchIntent = "general"
	externalSearchIntentConceptIntro externalSearchIntent = "concept_intro"
	externalSearchIntentToolIntro    externalSearchIntent = "tool_intro"
	externalSearchIntentHowTo        externalSearchIntent = "how_to"
	externalSearchIntentComparison   externalSearchIntent = "comparison"
)

var conceptIntroPositiveTerms = []string{"개념", "설명", "원리", "이해", "기초 이론", "기본", "구조", "특징", "입문"}

var conceptIntroNegativeTerms = []string{"하는 법", "따라하기", "실습", "사슬뜨기", "짧은뜨기", "작품 만들기", "비교", "차이"}

var toolIntroPositiveTerms = []string{"도구", "소개", "종류", "용도", "재질", "특징", "호수", "선택", "준비물", "실종류", "실 종류", "코바늘 종류"}

var toolIntroNegativeTerms = []string{"사슬뜨기", "짧은뜨기", "한길긴뜨기", "뜨는법", "뜨는 법", "방법", "실습", "기초뜨기", "기초 뜨기"}

var howToPositiveTerms = []string{"하는 법", "방법", "순서", "따라하기", "시작", "실습", "튜토리얼", "뜨는 법", "배우기"}

var howToNegativeTerms = []string{"개념", "원리", "종류", "용도", "비교", "차이"}

var comparisonPositiveTerms = []string{"비교", "차이", "선택", "장단점", "추천", "vs", "어떤", "호수별", "무엇이 다른가"}

var comparisonNegativeTerms = []string{"따라하기", "실습", "사슬뜨기", "작품 만들기"}

var comparisonWeakContextTerms = []string{"입문", "기초", "기본", "처음", "초급", "집에서", "함께", "중심으로", "이해", "이해하고"}

var externalGenericQueryTerms = map[string]struct{}{
	"강의": {}, "방법": {}, "하는": {}, "법": {}, "영상": {}, "튜토리얼": {}, "가이드": {}, "배우기": {}, "기초": {}, "입문": {}, "연습": {},
	"기본": {}, "간단": {}, "작은": {}, "결과물": {}, "완성": {}, "재료": {}, "도구": {}, "기법": {}, "표현": {}, "형태": {}, "명암": {}, "색": {},
	"tutorial": {}, "guide": {}, "how": {}, "to": {}, "learn": {}, "learning": {}, "basic": {}, "basics": {}, "beginner": {}, "practice": {},
}

var externalPrimaryQueryAnchorTerms = map[string]struct{}{
	"피아노": {}, "piano": {}, "바이올린": {}, "violin": {}, "첼로": {}, "cello": {}, "일렉기타": {}, "베이스기타": {}, "베이스": {}, "기타": {},
	"드럼": {}, "drum": {}, "drums": {}, "플룻": {}, "플루트": {}, "flute": {}, "색소폰": {}, "saxophone": {}, "리코더": {}, "recorder": {},
	"하모니카": {}, "harmonica": {}, "거문고": {}, "가야금": {}, "수채화": {}, "watercolor": {}, "목공예": {}, "목공": {}, "woodworking": {},
	"코바늘": {}, "뜨개": {}, "crochet": {}, "영어회화": {}, "여행영어": {}, "영어": {}, "english": {},
}

func detectExternalSearchIntent(queryBundle normalizer.LessonSearchQuery) externalSearchIntent {
	text := strings.ToLower(strings.TrimSpace(strings.Join([]string{
		queryBundle.RawQuery,
		queryBundle.DenseQuery,
		queryBundle.LexicalQueryKO,
	}, " ")))
	if containsAnyTerm(text, comparisonPositiveTerms) {
		return externalSearchIntentComparison
	}
	if containsAnyTerm(text, howToPositiveTerms) {
		return externalSearchIntentHowTo
	}
	if containsAnyTerm(text, toolIntroPositiveTerms) {
		return externalSearchIntentToolIntro
	}
	if containsAnyTerm(text, conceptIntroPositiveTerms) {
		return externalSearchIntentConceptIntro
	}
	return externalSearchIntentGeneral
}

func DetectExternalSearchIntent(queryBundle normalizer.LessonSearchQuery) string {
	return string(detectExternalSearchIntent(queryBundle))
}

func buildToolIntroExternalSearchQuery(queryBundle normalizer.LessonSearchQuery) string {
	parts := make([]string, 0, 8)
	for _, token := range queryBundle.Tokens {
		trimmed := strings.TrimSpace(token)
		if trimmed == "" {
			continue
		}
		lowered := strings.ToLower(trimmed)
		negative := false
		for _, term := range toolIntroNegativeTerms {
			if strings.Contains(lowered, strings.ToLower(term)) {
				negative = true
				break
			}
		}
		if negative {
			continue
		}
		parts = append(parts, trimmed)
		if len(parts) >= 2 {
			break
		}
	}
	for _, term := range []string{"도구", "소개", "종류", "용도"} {
		if !containsString(parts, term) {
			parts = append(parts, term)
		}
	}
	query := strings.Join(parts, " ")
	query = strings.Join(strings.Fields(query), " ")
	runes := []rune(query)
	if len(runes) > 80 {
		return string(runes[:80])
	}
	return query
}

func buildConceptIntroExternalSearchQuery(queryBundle normalizer.LessonSearchQuery) string {
	base := firstPositiveTokens(queryBundle.Tokens, conceptIntroNegativeTerms, 2)
	base = append(base, "개념", "설명")
	return compactIntentQuery(base, 80)
}

func buildHowToExternalSearchQuery(queryBundle normalizer.LessonSearchQuery) string {
	base := firstPositiveTokens(queryBundle.Tokens, howToNegativeTerms, 3)
	base = append(base, "하는 법", "방법")
	return compactIntentQuery(base, 80)
}

func buildComparisonExternalSearchQuery(queryBundle normalizer.LessonSearchQuery) string {
	base := firstPositiveTokensSkippingWeak(queryBundle.Tokens, comparisonNegativeTerms, comparisonWeakContextTerms, 5)
	if len(base) == 0 {
		base = firstPositiveTokens(queryBundle.Tokens, comparisonNegativeTerms, 3)
	}
	base = filterOutComparisonTerms(base)
	base = append(base, "비교", "차이")
	return compactIntentQuery(base, 80)
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func containsAnyTerm(text string, terms []string) bool {
	for _, term := range terms {
		if strings.Contains(text, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func firstPositiveTokens(tokens []string, negativeTerms []string, maxCount int) []string {
	result := make([]string, 0, maxCount)
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token)
		if trimmed == "" {
			continue
		}
		lowered := strings.ToLower(trimmed)
		negative := false
		for _, term := range negativeTerms {
			if strings.Contains(lowered, strings.ToLower(term)) {
				negative = true
				break
			}
		}
		if negative {
			continue
		}
		result = append(result, trimmed)
		if len(result) >= maxCount {
			break
		}
	}
	return result
}

func firstPositiveTokensSkippingWeak(tokens []string, negativeTerms, weakTerms []string, maxCount int) []string {
	result := make([]string, 0, maxCount)
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token)
		if trimmed == "" {
			continue
		}
		lowered := strings.ToLower(trimmed)
		negative := false
		for _, term := range negativeTerms {
			if strings.Contains(lowered, strings.ToLower(term)) {
				negative = true
				break
			}
		}
		if negative {
			continue
		}
		weak := false
		for _, term := range weakTerms {
			if lowered == strings.ToLower(strings.TrimSpace(term)) {
				weak = true
				break
			}
		}
		if weak {
			continue
		}
		result = append(result, trimmed)
		if len(result) >= maxCount {
			break
		}
	}
	return result
}

func compactIntentQuery(parts []string, maxRunes int) string {
	unique := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		unique = append(unique, trimmed)
	}
	query := strings.Join(unique, " ")
	query = strings.Join(strings.Fields(query), " ")
	runes := []rune(query)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes])
	}
	return query
}

func filterOutComparisonTerms(parts []string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		lowered := strings.ToLower(trimmed)
		if lowered == "비교" || lowered == "차이" || lowered == "선택" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}
