package curriculum

import (
	"sort"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func rerankExplorerCandidates(
	candidates []ContentSearchCandidate,
	stages []explorerSearchStage,
	limit int,
) []ContentSearchCandidate {
	if len(candidates) <= 1 {
		return candidates
	}

	courseTokens := make([]string, 0, 4)
	levelTokens := make([]string, 0, 4)
	for _, stage := range stages {
		switch stage.name {
		case "course":
			for _, token := range stage.filter.strongTokens {
				if isPriorityExplorerCourseToken(token) {
					courseTokens = append(courseTokens, token)
				}
			}
		case "level":
			levelTokens = append(levelTokens, stage.filter.strongTokens...)
		}
	}

	type scoredCandidate struct {
		candidate ContentSearchCandidate
		index     int
		score     int
	}
	scored := make([]scoredCandidate, 0, len(candidates))
	for idx, candidate := range candidates {
		searchText := normalizer.BuildSearchTextKO(
			candidate.Title,
			derefString(candidate.Description),
			derefString(candidate.ExternalURL),
		)
		score := 0
		for _, token := range courseTokens {
			if explorerTokenMatchesText(searchText, token) {
				score += 10
			}
		}
		for _, token := range levelTokens {
			if explorerTokenMatchesText(searchText, token) {
				score += 2
			}
		}
		score += explorerRecommendationSourcePreference(candidate)
		scored = append(scored, scoredCandidate{
			candidate: candidate,
			index:     idx,
			score:     score,
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].index < scored[j].index
		}
		return scored[i].score > scored[j].score
	})

	reranked := make([]ContentSearchCandidate, 0, len(scored))
	for _, item := range scored {
		reranked = append(reranked, item.candidate)
	}
	reranked = ensureExplorerVideoQuota(reranked, limit)
	if limit > 0 && len(reranked) > limit {
		return reranked[:limit]
	}
	return reranked
}

func rerankExplorerCandidatesWithPrimaryQuery(
	candidates []ContentSearchCandidate,
	stages []explorerSearchStage,
	primaryQuery string,
	limit int,
) []ContentSearchCandidate {
	if len(candidates) <= 1 {
		return candidates
	}

	primaryBundle := normalizer.BuildLessonSearchQuery(primaryQuery, "", "", "", "", "")
	primaryTokens := make([]string, 0, len(primaryBundle.Tokens)+len(primaryBundle.ExpandedTokens))
	primaryTokenSet := map[string]struct{}{}
	for _, token := range append(primaryBundle.Tokens, primaryBundle.ExpandedTokens...) {
		trimmed := normalizeExplorerFilterToken(token)
		if trimmed == "" {
			continue
		}
		if _, exists := primaryTokenSet[trimmed]; exists {
			continue
		}
		primaryTokenSet[trimmed] = struct{}{}
		primaryTokens = append(primaryTokens, trimmed)
	}
	if len(primaryTokens) == 0 {
		return rerankExplorerCandidates(candidates, stages, limit)
	}

	contextTokens := make([]string, 0, 8)
	for _, stage := range stages {
		for _, token := range stage.filter.strongTokens {
			trimmed := normalizeExplorerFilterToken(token)
			if trimmed == "" {
				continue
			}
			if _, primary := primaryTokenSet[trimmed]; primary {
				continue
			}
			contextTokens = append(contextTokens, trimmed)
		}
	}
	driftPenaltyTerms := directQueryDriftPenaltyTerms(primaryTokenSet)
	primaryPhrase := strings.ToLower(strings.Join(strings.Fields(primaryQuery), " "))

	type scoredCandidate struct {
		candidate ContentSearchCandidate
		index     int
		score     int
	}
	scored := make([]scoredCandidate, 0, len(candidates))
	for idx, candidate := range candidates {
		searchText := normalizer.BuildSearchTextKO(
			candidate.Title,
			derefString(candidate.Description),
			derefString(candidate.ExternalURL),
		)
		score := 0
		if primaryPhrase != "" && strings.Contains(searchText, primaryPhrase) {
			score += 30
		}
		primaryMatches := 0
		for _, token := range primaryTokens {
			if explorerTokenMatchesText(searchText, token) {
				primaryMatches++
				score += 14
			}
		}
		for _, token := range contextTokens {
			if explorerTokenMatchesText(searchText, token) {
				score += 1
			}
		}
		if primaryMatches == 0 {
			score -= 30
		}
		for _, term := range driftPenaltyTerms {
			if explorerTokenMatchesText(searchText, term) {
				score -= 8
			}
		}
		score += explorerWeakCandidatePenalty(primaryTokenSet, searchText)
		score += explorerRecommendationSourcePreference(candidate)
		scored = append(scored, scoredCandidate{
			candidate: candidate,
			index:     idx,
			score:     score,
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].index < scored[j].index
		}
		return scored[i].score > scored[j].score
	})

	reranked := make([]ContentSearchCandidate, 0, len(scored))
	for _, item := range scored {
		reranked = append(reranked, item.candidate)
	}
	reranked = ensureExplorerVideoQuota(reranked, limit)
	if limit > 0 && len(reranked) > limit {
		return reranked[:limit]
	}
	return reranked
}

func rerankExplorerCandidatesWithSearchSpec(
	candidates []ContentSearchCandidate,
	spec LessonRecommendationSearchSpec,
	limit int,
) []ContentSearchCandidate {
	spec = NormalizeLessonRecommendationSearchSpec(spec, LessonRecommendationSearchSpec{})
	if len(candidates) <= 1 || strings.TrimSpace(spec.PrimaryQuery) == "" {
		return candidates
	}

	type scoredCandidate struct {
		candidate ContentSearchCandidate
		index     int
		score     int
	}
	scored := make([]scoredCandidate, 0, len(candidates))
	for idx, candidate := range candidates {
		scored = append(scored, scoredCandidate{
			candidate: candidate,
			index:     idx,
			score:     ScoreCandidateWithLessonSearchSpec(candidate, spec),
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].index < scored[j].index
		}
		return scored[i].score > scored[j].score
	})

	reranked := make([]ContentSearchCandidate, 0, len(scored))
	for _, item := range scored {
		reranked = append(reranked, item.candidate)
	}
	reranked = ensureExplorerVideoQuota(reranked, limit)
	if limit > 0 && len(reranked) > limit {
		return reranked[:limit]
	}
	return reranked
}

func directQueryDriftPenaltyTerms(primaryTokenSet map[string]struct{}) []string {
	terms := []string{"수익", "수익화", "창출", "매출", "돈", "돈버는", "마케팅", "판매"}
	filtered := make([]string, 0, len(terms))
	for _, term := range terms {
		if _, exists := primaryTokenSet[term]; exists {
			continue
		}
		filtered = append(filtered, term)
	}
	return filtered
}

func explorerWeakCandidatePenalty(primaryTokenSet map[string]struct{}, searchText string) int {
	penalty := 0
	if explorerPrimaryHasAny(primaryTokenSet, "사진", "스마트폰", "구도", "보정") &&
		containsAnyExplorerSearchTerm(searchText, "책리뷰", "책 리뷰", "도서리뷰", "도서 리뷰", "서평") {
		penalty -= 35
	}
	if explorerPrimaryHas(primaryTokenSet, "드론") &&
		explorerPrimaryHasAny(primaryTokenSet, "조종", "비행", "이륙", "착륙", "호버링") &&
		containsAnyExplorerSearchTerm(searchText, "코딩", "파이썬", "python", "프로그래밍") {
		penalty -= 28
	}
	if explorerPrimaryHas(primaryTokenSet, "드론") &&
		containsAnyExplorerSearchTerm(searchText, "epp", "전투기", "비행기", "항공기") {
		penalty -= 32
	}
	return penalty
}

func explorerPrimaryHas(primaryTokenSet map[string]struct{}, token string) bool {
	_, ok := primaryTokenSet[normalizeExplorerFilterToken(token)]
	return ok
}

func explorerPrimaryHasAny(primaryTokenSet map[string]struct{}, tokens ...string) bool {
	for _, token := range tokens {
		if explorerPrimaryHas(primaryTokenSet, token) {
			return true
		}
	}
	return false
}

func containsAnyExplorerSearchTerm(searchText string, terms ...string) bool {
	for _, term := range terms {
		if explorerTokenMatchesText(searchText, strings.ToLower(strings.TrimSpace(term))) {
			return true
		}
	}
	return false
}

func ensureExplorerVideoQuota(candidates []ContentSearchCandidate, limit int) []ContentSearchCandidate {
	if limit < 3 || len(candidates) <= limit {
		return candidates
	}
	quota := minInt(2, limit)
	window := append([]ContentSearchCandidate{}, candidates[:limit]...)
	rest := append([]ContentSearchCandidate{}, candidates[limit:]...)

	videoCount := 0
	for _, candidate := range window {
		if isExplorerVideoCandidate(candidate) {
			videoCount++
		}
	}
	if videoCount >= quota {
		return candidates
	}

	for restIdx := 0; restIdx < len(rest) && videoCount < quota; restIdx++ {
		promote := rest[restIdx]
		if !isExplorerVideoCandidate(promote) {
			continue
		}
		removeIdx := -1
		for idx := len(window) - 1; idx >= 0; idx-- {
			if !isExplorerVideoCandidate(window[idx]) {
				removeIdx = idx
				break
			}
		}
		if removeIdx < 0 {
			break
		}
		demoted := window[removeIdx]
		window = append(window[:removeIdx], window[removeIdx+1:]...)
		insertIdx := minInt(1+videoCount*2, len(window))
		window = append(window[:insertIdx], append([]ContentSearchCandidate{promote}, window[insertIdx:]...)...)
		rest[restIdx] = demoted
		videoCount++
	}

	return append(window, rest...)
}

func isExplorerVideoCandidate(candidate ContentSearchCandidate) bool {
	contentType := strings.ToLower(strings.TrimSpace(candidate.ContentType))
	externalURL := strings.ToLower(strings.TrimSpace(derefString(candidate.ExternalURL)))
	return contentType == "youtube" || strings.Contains(externalURL, "youtube.com/") || strings.Contains(externalURL, "youtu.be/")
}

func explorerRecommendationSourcePreference(candidate ContentSearchCandidate) int {
	switch {
	case isExplorerVideoCandidate(candidate):
		return 18
	case strings.ToLower(strings.TrimSpace(candidate.ContentType)) == "naver_blog", strings.ToLower(strings.TrimSpace(candidate.ContentType)) == "blog":
		return -6
	default:
		return 0
	}
}

func buildExplorerFilterStages(bundles ...normalizer.LessonSearchQuery) []explorerFilterStage {
	stageNames := []string{"node", "level", "course"}
	stages := make([]explorerFilterStage, 0, len(bundles))
	for idx, bundle := range bundles {
		stageName := "course"
		if idx < len(stageNames) {
			stageName = stageNames[idx]
		}
		stage, ok := buildExplorerFilterStage(stageName, bundle)
		if !ok {
			continue
		}
		stages = append(stages, stage)
	}
	return stages
}

func buildExplorerFilterStage(name string, bundle normalizer.LessonSearchQuery) (explorerFilterStage, bool) {
	requiredTokens := collectExplorerFilterTokens(bundle)
	if len(requiredTokens) == 0 {
		return explorerFilterStage{}, false
	}

	strongTokens := make([]string, 0, len(requiredTokens))
	for _, token := range requiredTokens {
		if isWeakExplorerFilterToken(token) {
			continue
		}
		strongTokens = append(strongTokens, token)
	}
	if len(strongTokens) == 0 {
		return explorerFilterStage{}, false
	}

	minStrongMatch := 1
	switch name {
	case "node":
		if len(strongTokens) >= 3 {
			minStrongMatch = 2
		}
	case "course":
		minStrongMatch = minInt(2, len(strongTokens))
	case "level":
		if len(strongTokens) >= 3 {
			minStrongMatch = 2
		}
	}

	return explorerFilterStage{
		requiredTokens: requiredTokens,
		strongTokens:   strongTokens,
		minStrongMatch: minStrongMatch,
	}, true
}

func collectExplorerFilterTokens(queryBundle normalizer.LessonSearchQuery) []string {
	requiredTokens := make([]string, 0, len(queryBundle.Tokens)+len(queryBundle.ExpandedTokens))
	seen := make(map[string]struct{}, len(queryBundle.Tokens)+len(queryBundle.ExpandedTokens))
	for _, token := range append(queryBundle.Tokens, queryBundle.ExpandedTokens...) {
		trimmed := normalizeExplorerFilterToken(token)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		requiredTokens = append(requiredTokens, trimmed)
	}
	return requiredTokens
}

func explorerCandidateFilterKey(candidate ContentSearchCandidate) string {
	if candidate.ContentID != nil {
		return "content:" + candidate.ContentID.String()
	}
	if candidate.ExternalURL != nil && strings.TrimSpace(*candidate.ExternalURL) != "" {
		return "url:" + strings.TrimSpace(*candidate.ExternalURL)
	}
	return "title:" + strings.TrimSpace(strings.ToLower(candidate.Title))
}

var explorerWeakFilterTokens = map[string]struct{}{
	"입문": {}, "기초": {}, "초급": {}, "처음": {}, "처음부터": {},
	"연습": {}, "훈련": {}, "실습": {}, "이해": {}, "설명": {},
	"정리": {}, "소개": {}, "가이드": {}, "배우기": {}, "익히기": {},
	"배우는": {}, "익히는": {}, "사용법": {}, "방법": {}, "재료": {},
	"도구": {}, "준비": {}, "활용": {}, "응용": {}, "기본": {},
	"핵심": {}, "표현": {}, "문장": {}, "패턴": {}, "질문": {}, "응답": {},
	"대화": {}, "상황": {}, "흐름": {}, "연결": {}, "시작": {}, "시작한다": {},
	"이어가기": {}, "이어간다": {}, "연결한다": {}, "자연스럽게": {}, "실제": {},
	"필요": {}, "필요한": {}, "원활하게": {}, "하게": {}, "있게": {}, "수": {},
	"한다": {}, "할": {}, "맞는": {}, "맞춰": {}, "짧은": {}, "자주": {},
}

var explorerParticleSuffixes = []string{
	"으로", "에서", "에게", "까지", "부터", "보다", "처럼",
	"은", "는", "이", "가", "을", "를", "의", "에", "와", "과", "도", "로",
}

func normalizeExplorerFilterToken(token string) string {
	normalized := strings.TrimSpace(strings.ToLower(token))
	if normalized == "" {
		return ""
	}
	for _, suffix := range explorerParticleSuffixes {
		if strings.HasSuffix(normalized, suffix) && len([]rune(normalized)) > len([]rune(suffix))+1 {
			normalized = strings.TrimSuffix(normalized, suffix)
			break
		}
	}
	return normalized
}

func isWeakExplorerFilterToken(token string) bool {
	_, exists := explorerWeakFilterTokens[token]
	return exists
}

func isPriorityExplorerCourseToken(token string) bool {
	if isWeakExplorerFilterToken(token) {
		return false
	}
	switch token {
	case "그림", "실력", "향상", "향상시키다", "준비", "준비한다", "과정", "커리큘럼", "말하기":
		return false
	}
	return len([]rune(token)) >= 3
}

func matchesExplorerCandidateContext(searchText string, strongTokens, requiredTokens []string, minStrongMatch int) bool {
	if len(strongTokens) > 0 {
		if minStrongMatch <= 0 {
			minStrongMatch = 1
		}
		matchedStrong := 0
		for _, token := range strongTokens {
			if explorerTokenMatchesText(searchText, token) {
				matchedStrong++
				if matchedStrong >= minStrongMatch {
					return true
				}
			}
		}
		return false
	}

	matchedWeak := 0
	for _, token := range requiredTokens {
		if explorerTokenMatchesText(searchText, token) {
			matchedWeak++
			if matchedWeak >= 2 {
				return true
			}
		}
	}
	return false
}
