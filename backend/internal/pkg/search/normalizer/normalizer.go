package normalizer

import (
	"regexp"
	"slices"
	"strings"
)

type LessonSearchQuery struct {
	RawQuery             string
	DenseQuery           string
	LexicalQueryKO       string
	LexicalQueryExpanded string
	Tokens               []string
	ExpandedTokens       []string
}

var tokenPattern = regexp.MustCompile(`[\p{Hangul}\p{Latin}\p{Nd}]+`)

var canonicalReplacer = strings.NewReplacer(
	"코테", "코딩테스트",
	"자소서", "자기소개서",
	"포폴", "포트폴리오",
	"포토 샵", "포토샵",
	"일러", "일러스트레이터",
	"프리미어 프로", "프리미어프로",
	"애프터 이펙트", "애프터이펙트",
	"애펙", "애프터이펙트",
	"캡 컷", "capcut",
	"캡컷", "capcut",
	"미디", "midi",
	"js", "javascript",
	"ts", "typescript",
	"리액트js", "react",
	"넥스트js", "nextjs",
	"노드js", "nodejs",
	"영단어", "영어 단어",
	"영문법", "영어 문법",
	"회화표현", "회화 표현",
	"시장조사", "시장 조사",
	"강의주제선정", "강의 주제 선정",
	"강의주제", "강의 주제",
	"주제선정", "주제 선정",
	"수익창출", "수익 창출",
)

var stopwords = map[string]struct{}{
	"관련":   {},
	"알아보기": {},
	"소개":   {},
	"정리":   {},
	"모음":   {},
	"가이드":  {},
	"팁":    {},
	"요약":   {},
	"해보기":  {},
	"해보는":  {},
	"배워보기": {},
	"알기":   {},
	"이해하기": {},
}

var synonymExpansion = map[string][]string{
	"독학":    {"혼자공부", "셀프스터디"},
	"입문":    {"초급", "처음", "기초부터"},
	"실전":    {"응용", "활용", "실습"},
	"기타솔로":  {"기타리드", "리드기타"},
	"기타리드":  {"기타솔로", "리드기타"},
	"코드진행":  {"코드전개", "코드패턴"},
	"반주":    {"반주법", "반주패턴"},
	"드로잉":   {"그림그리기", "스케치"},
	"목공예":   {"목공", "woodworking"},
	"자막":    {"캡션", "subtitle"},
	"컷편집":   {"클립편집"},
	"프론트엔드": {"frontend", "웹프론트엔드"},
	"백엔드":   {"backend", "서버개발"},
	"api":   {"restapi", "백엔드api"},
	"회화":    {"말하기", "스피킹"},
	"발음":    {"발음교정"},
	"코딩테스트": {"알고리즘"},
	"영어회화":  {"회화", "스피킹"},
	"비트메이킹": {"비트메이킹", "비트메이킹기초"},
}

func BuildLessonSearchQuery(sourceQuery, levelTitle, levelObjective, lessonTitle, lessonObjective, lessonSummary string) LessonSearchQuery {
	rawParts := compactParts(
		sourceQuery,
		levelTitle,
		lessonTitle,
		levelObjective,
		lessonObjective,
		lessonSummary,
	)
	rawQuery := strings.Join(rawParts, " ")
	denseQuery := normalizeText(rawQuery)
	tokens := tokenize(denseQuery)
	expandedTokens := expandTokens(tokens)

	return LessonSearchQuery{
		RawQuery:             rawQuery,
		DenseQuery:           denseQuery,
		LexicalQueryKO:       strings.Join(tokens, " "),
		LexicalQueryExpanded: strings.Join(expandedTokens, " "),
		Tokens:               tokens,
		ExpandedTokens:       expandedTokens,
	}
}

func BuildSearchTextKO(fields ...string) string {
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		normalized := normalizeText(field)
		if normalized != "" {
			parts = append(parts, normalized)
		}
	}
	return strings.Join(parts, " ")
}

func compactParts(parts ...string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeText(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return ""
	}
	normalized = canonicalReplacer.Replace(normalized)
	normalized = strings.Join(strings.Fields(normalized), " ")
	return normalized
}

func tokenize(value string) []string {
	matches := tokenPattern.FindAllString(value, -1)
	tokens := make([]string, 0, len(matches))
	seen := map[string]struct{}{}
	for _, match := range matches {
		token := strings.ToLower(strings.TrimSpace(match))
		if token == "" {
			continue
		}
		if _, isStopword := stopwords[token]; isStopword {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens
}

func expandTokens(tokens []string) []string {
	expanded := slices.Clone(tokens)
	seen := map[string]struct{}{}
	for _, token := range expanded {
		seen[token] = struct{}{}
	}
	for _, token := range tokens {
		for _, synonym := range synonymExpansion[token] {
			normalized := normalizeText(strings.ReplaceAll(synonym, " ", ""))
			if normalized == "" {
				continue
			}
			if _, exists := seen[normalized]; exists {
				continue
			}
			seen[normalized] = struct{}{}
			expanded = append(expanded, normalized)
		}
	}
	return expanded
}
