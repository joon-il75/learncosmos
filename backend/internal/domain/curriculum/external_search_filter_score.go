package curriculum

import (
	"net/url"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/extsearch"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func recommendationBlockURLKey(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Host != "" {
		parsed.Fragment = ""
		parsed.Scheme = strings.ToLower(parsed.Scheme)
		parsed.Host = strings.ToLower(parsed.Host)
		parsed.Path = strings.TrimRight(parsed.EscapedPath(), "/")
		if parsed.Path == "" {
			parsed.Path = "/"
		}
		return strings.TrimRight(parsed.String(), "/")
	}
	return strings.TrimRight(strings.ToLower(trimmed), "/")
}

func externalSearchResultKey(result extsearch.SearchResult) string {
	if trimmed := strings.TrimSpace(result.CanonicalURL); trimmed != "" {
		return "url:" + trimmed
	}
	if trimmed := strings.TrimSpace(result.ExternalContentID); trimmed != "" {
		return "ext:" + strings.TrimSpace(result.Source) + ":" + trimmed
	}
	return "title:" + strings.TrimSpace(strings.ToLower(result.Title))
}

var instructionalPositiveTerms = []string{"강의", "방법", "하는 법", "알려", "배우", "튜토리얼", "가이드", "노하우", "실전", "단계", "팁", "사례", "초보", "입문", "기초", "연습", "정리", "설명", "핵심", "전략", "만들기", "제작", "선정", "분석", "실습", "how to", "tutorial", "guide"}

var externalPromoTerms = []string{"상담", "문의", "신청", "모집", "수강료", "할인", "이벤트", "국비지원", "대행", "업체", "오픈", "공지", "가격", "비용", "렌탈", "예약", "자격증", "취득", "합격", "학원", "후기", "스튜디오", "대관"}

var externalHardPromoTerms = []string{"상담", "문의", "신청", "모집", "수강료", "할인", "이벤트", "국비지원", "대행", "업체", "오픈", "공지", "자격증", "취득", "합격", "학원", "후기", "스튜디오", "대관"}

func isInstructionalExternalSearchResult(result extsearch.SearchResult) bool {
	return isInstructionalExternalSearchResultForQuery(result, normalizer.LessonSearchQuery{})
}

func isInstructionalExternalSearchResultForQuery(result extsearch.SearchResult, queryBundle normalizer.LessonSearchQuery) bool {
	instructionalScore, promoScore := scoreExternalInstructionalSignals(result)
	source := strings.ToLower(strings.TrimSpace(result.Source))

	if source == "naver_blog" || source == "blog" {
		if containsExternalHardPromoSignal(result) {
			return false
		}
		if promoScore >= 2 {
			return false
		}
		if promoScore >= 1 && instructionalScore <= 1 {
			return false
		}
		if hasExternalSpecificQueryTokens(queryBundle) && countExternalQueryTokenMatches(queryBundle, result) == 0 {
			return false
		}
		if anchors := externalPrimaryQueryAnchors(queryBundle); len(anchors) > 0 && countExternalQueryTokenMatchesForTokens(anchors, result) == 0 {
			return false
		}
		return instructionalScore >= 1
	}

	if promoScore >= 2 && instructionalScore <= 1 {
		return false
	}
	if hasExternalSpecificQueryTokens(queryBundle) && countExternalQueryTokenMatches(queryBundle, result) == 0 {
		return false
	}
	if anchors := externalPrimaryQueryAnchors(queryBundle); len(anchors) > 0 && countExternalQueryTokenMatchesForTokens(anchors, result) == 0 {
		return false
	}
	if source == "youtube" {
		return instructionalScore >= 1
	}
	return instructionalScore >= 1
}

func hasExternalSpecificQueryTokens(queryBundle normalizer.LessonSearchQuery) bool {
	for _, token := range append(queryBundle.Tokens, queryBundle.ExpandedTokens...) {
		token = strings.ToLower(strings.TrimSpace(token))
		if len([]rune(token)) < 2 {
			continue
		}
		if _, generic := externalGenericQueryTerms[token]; generic {
			continue
		}
		return true
	}
	return false
}

func countExternalQueryTokenMatches(queryBundle normalizer.LessonSearchQuery, result extsearch.SearchResult) int {
	return countExternalQueryTokenMatchesForTokens(append(queryBundle.Tokens, queryBundle.ExpandedTokens...), result)
}

func countExternalQueryTokenMatchesForTokens(tokens []string, result extsearch.SearchResult) int {
	text := strings.ToLower(strings.TrimSpace(strings.Join([]string{
		result.Title,
		result.Description,
		result.Author,
	}, " ")))
	if text == "" {
		return 0
	}

	count := 0
	seen := map[string]struct{}{}
	for _, token := range tokens {
		token = strings.ToLower(strings.TrimSpace(token))
		if len([]rune(token)) < 2 {
			continue
		}
		if _, generic := externalGenericQueryTerms[token]; generic {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		if strings.Contains(text, token) {
			count++
		}
	}
	return count
}

func externalPrimaryQueryAnchors(queryBundle normalizer.LessonSearchQuery) []string {
	anchors := make([]string, 0, 4)
	seen := map[string]struct{}{}
	for _, token := range append(queryBundle.Tokens, queryBundle.ExpandedTokens...) {
		token = strings.ToLower(strings.TrimSpace(token))
		if len([]rune(token)) < 2 {
			continue
		}
		if _, ok := externalPrimaryQueryAnchorTerms[token]; !ok {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		anchors = append(anchors, token)
	}
	return anchors
}

func containsExternalHardPromoSignal(result extsearch.SearchResult) bool {
	text := strings.ToLower(strings.TrimSpace(strings.Join([]string{
		result.Title,
		result.Description,
		result.Author,
	}, " ")))
	for _, term := range externalHardPromoTerms {
		if strings.Contains(text, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func scoreExternalInstructionalSignals(result extsearch.SearchResult) (int, int) {
	text := strings.ToLower(strings.TrimSpace(strings.Join([]string{
		result.Title,
		result.Description,
		result.Author,
	}, " ")))
	instructionalScore := 0
	for _, term := range instructionalPositiveTerms {
		if strings.Contains(text, strings.ToLower(term)) {
			instructionalScore++
		}
	}

	promoScore := 0
	for _, term := range externalPromoTerms {
		if strings.Contains(text, strings.ToLower(term)) {
			promoScore++
		}
	}
	return instructionalScore, promoScore
}

func stringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func scoreExternalSearchResult(queryBundle normalizer.LessonSearchQuery, result extsearch.SearchResult, preferredFormat *string, rankIndex int) float64 {
	text := strings.ToLower(strings.TrimSpace(result.Title + " " + result.Description))
	if text == "" {
		text = strings.ToLower(strings.TrimSpace(result.Title))
	}

	matched := 0
	seen := map[string]struct{}{}
	for _, token := range append(queryBundle.Tokens, queryBundle.ExpandedTokens...) {
		token = strings.ToLower(strings.TrimSpace(token))
		if len([]rune(token)) < 2 {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		if strings.Contains(text, token) {
			matched++
		}
	}

	score := 0.025
	if matched > 0 {
		score += float64(minInt(matched, 4)) * 0.012
	}
	if preferredFormatMatchesExternal(preferredFormat, result.Source) {
		score += 0.02
	}
	instructionalScore, promoScore := scoreExternalInstructionalSignals(result)
	score += float64(minInt(instructionalScore, 4)) * 0.008
	score -= float64(minInt(promoScore, 3)) * 0.015
	switch detectExternalSearchIntent(queryBundle) {
	case externalSearchIntentConceptIntro:
		score += scoreTermMatches(text, conceptIntroPositiveTerms, 0.012)
		score -= scoreTermMatches(text, conceptIntroNegativeTerms, 0.018)
	case externalSearchIntentToolIntro:
		score += scoreTermMatches(text, toolIntroPositiveTerms, 0.015)
		score -= scoreTermMatches(text, toolIntroNegativeTerms, 0.02)
	case externalSearchIntentHowTo:
		score += scoreTermMatches(text, howToPositiveTerms, 0.015)
		score -= scoreTermMatches(text, howToNegativeTerms, 0.02)
	case externalSearchIntentComparison:
		score += scoreTermMatches(text, comparisonPositiveTerms, 0.014)
		score -= scoreTermMatches(text, comparisonNegativeTerms, 0.02)
	}
	score -= float64(rankIndex) * 0.001
	if score < 0.005 {
		score = 0.005
	}
	return score
}

func preferredFormatMatchesExternal(preferredFormat *string, source string) bool {
	if preferredFormat == nil {
		return false
	}
	format := strings.ToLower(strings.TrimSpace(*preferredFormat))
	if format == "" {
		return false
	}
	switch {
	case strings.Contains(format, "video"), strings.Contains(format, "영상"), strings.Contains(format, "유튜브"):
		return source == "youtube"
	case strings.Contains(format, "text"), strings.Contains(format, "글"), strings.Contains(format, "읽기"), strings.Contains(format, "article"), strings.Contains(format, "blog"):
		return source == "blog" || source == "article"
	default:
		return false
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func scoreTermMatches(text string, terms []string, weight float64) float64 {
	score := 0.0
	seen := map[string]struct{}{}
	for _, term := range terms {
		key := strings.ToLower(strings.TrimSpace(term))
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		if strings.Contains(text, key) {
			score += weight
		}
	}
	return score
}
