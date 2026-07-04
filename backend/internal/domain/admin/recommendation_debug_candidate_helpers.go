package admin

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func (h *AdminHandler) buildRecommendationDebugExternalStage(ctx context.Context, ownerUserID uuid.UUID, maxPerLesson int, queryBundle normalizer.LessonSearchQuery, internalCandidates []curriculum.ContentSearchCandidate) gin.H {
	return h.buildRecommendationDebugExternalStageWithQuery(ctx, ownerUserID, maxPerLesson, curriculum.BuildExternalSearchQuery(queryBundle), internalCandidates)
}

func (h *AdminHandler) buildRecommendationDebugExternalStageWithQuery(ctx context.Context, ownerUserID uuid.UUID, maxPerLesson int, externalQuery string, internalCandidates []curriculum.ContentSearchCandidate) gin.H {
	const externalSearchLimit = 10
	externalQuery = strings.TrimSpace(externalQuery)

	youtubeAPIKey := strings.TrimSpace(os.Getenv("YOUTUBE_API_KEY"))
	naverClientID := strings.TrimSpace(os.Getenv("NAVER_CLIENT_ID"))
	naverClientSecret := strings.TrimSpace(os.Getenv("NAVER_CLIENT_SECRET"))
	providers := []extsearch.SearchProvider{}
	if youtubeAPIKey != "" {
		providers = append(providers, extsearch.NewYouTubeProvider(youtubeAPIKey))
	}
	if naverClientID != "" && naverClientSecret != "" {
		providers = append(providers, extsearch.NewNaverBlogProvider(naverClientID, naverClientSecret))
	}
	if len(providers) == 0 {
		return gin.H{
			"provider":                         "youtube,naver_blog",
			"attempted":                        false,
			"reason":                           "external_provider_unavailable",
			"external_query":                   externalQuery,
			"raw_hit_count":                    0,
			"existing_content_match_count":     0,
			"dedupe_dropped_count":             0,
			"supplement_hit_count":             0,
			"final_merged_with_external_count": len(internalCandidates),
			"preview_results":                  []gin.H{},
		}
	}

	results, providerStats := curriculum.CollectFilteredExternalSearchResults(ctx, providers, externalQuery, curriculum.IsInstructionalExternalSearchResult)
	if len(results) == 0 {
		return gin.H{
			"provider":                         "youtube,naver_blog",
			"attempted":                        true,
			"reason":                           "external_results_empty",
			"external_query":                   externalQuery,
			"raw_hit_count":                    0,
			"existing_content_match_count":     0,
			"dedupe_dropped_count":             0,
			"supplement_hit_count":             0,
			"final_merged_with_external_count": len(internalCandidates),
			"provider_stats":                   recommendationDebugExternalProviderStats(providerStats),
			"preview_results":                  []gin.H{},
		}
	}

	contentRepo := content.NewRepository(h.db)
	internalIDs := make(map[uuid.UUID]struct{}, len(internalCandidates))
	for _, candidate := range internalCandidates {
		if candidate.ContentID != nil {
			internalIDs[*candidate.ContentID] = struct{}{}
		}
	}

	seenKeys := map[string]struct{}{}
	existingMatchCount := 0
	dedupeDroppedCount := 0
	supplementCount := 0
	previewResults := make([]gin.H, 0, len(results))

	for _, result := range results {
		existing, key := h.lookupRecommendationDebugExternalContent(ctx, contentRepo, result)
		alreadySaved := existing != nil
		if alreadySaved {
			existingMatchCount++
		}

		duplicateWithInternal := false
		if existing != nil {
			_, duplicateWithInternal = internalIDs[existing.ID]
		}

		dedupeKey := key
		if dedupeKey == "" {
			dedupeKey = recommendationDebugExternalResultKey(result)
		}

		duplicateWithinExternal := false
		if _, exists := seenKeys[dedupeKey]; exists {
			duplicateWithinExternal = true
		} else if dedupeKey != "" {
			seenKeys[dedupeKey] = struct{}{}
		}

		eligible := !duplicateWithInternal && !duplicateWithinExternal && supplementCount < externalSearchLimit
		if duplicateWithInternal || duplicateWithinExternal {
			dedupeDroppedCount++
		}
		if eligible {
			supplementCount++
		}

		previewResults = append(previewResults, gin.H{
			"source":                  result.Source,
			"title":                   result.Title,
			"description":             result.Description,
			"url":                     result.URL,
			"external_url":            result.URL,
			"canonical_url":           result.CanonicalURL,
			"thumbnail_url":           result.ThumbnailURL,
			"author":                  result.Author,
			"language":                result.Language,
			"external_content_id":     result.ExternalContentID,
			"content_id":              recommendationDebugExistingContentID(existing),
			"already_saved":           alreadySaved,
			"duplicate_with_internal": duplicateWithInternal,
			"duplicate_with_external": duplicateWithinExternal,
			"eligible_for_supplement": eligible,
			"reason":                  recommendationDebugExternalReason(alreadySaved, duplicateWithInternal, duplicateWithinExternal, eligible),
		})
	}

	finalCount := len(internalCandidates) + supplementCount
	if finalCount > maxPerLesson {
		finalCount = maxPerLesson
	}

	return gin.H{
		"provider":                         "youtube,naver_blog",
		"attempted":                        true,
		"reason":                           "",
		"external_query":                   externalQuery,
		"raw_hit_count":                    len(results),
		"existing_content_match_count":     existingMatchCount,
		"dedupe_dropped_count":             dedupeDroppedCount,
		"supplement_hit_count":             supplementCount,
		"final_merged_with_external_count": finalCount,
		"provider_stats":                   recommendationDebugExternalProviderStats(providerStats),
		"preview_results":                  previewResults,
	}
}

func recommendationDebugExternalProviderStats(stats []curriculum.ExternalSearchProviderStats) []gin.H {
	items := make([]gin.H, 0, len(stats))
	for _, stat := range stats {
		items = append(items, gin.H{
			"provider":        stat.Provider,
			"query":           stat.Query,
			"pages_attempted": stat.PagesAttempted,
			"raw_count":       stat.RawCount,
			"accepted_count":  stat.AcceptedCount,
			"filtered_count":  stat.FilteredCount,
			"error_reason":    stat.ErrorReason,
		})
	}
	return items
}

func recommendationDebugExternalReason(alreadySaved, duplicateWithInternal, duplicateWithinExternal, eligible bool) string {
	reasons := make([]string, 0, 4)
	if alreadySaved {
		reasons = append(reasons, "이미 저장된 콘텐츠와 매칭")
	}
	if duplicateWithInternal {
		reasons = append(reasons, "내부 후보와 중복되어 제외")
	}
	if duplicateWithinExternal {
		reasons = append(reasons, "외부 결과끼리 중복되어 제외")
	}
	if eligible {
		reasons = append(reasons, "부족한 내부 후보를 보강할 수 있음")
	}
	if len(reasons) == 0 {
		return "외부 후보로 조회됐지만 최종 보강 조건은 충족하지 않음"
	}
	return strings.Join(reasons, ", ")
}

func recommendationDebugExistingContentID(existing *content.Content) string {
	if existing == nil {
		return ""
	}
	return existing.ID.String()
}

func (h *AdminHandler) lookupRecommendationDebugExternalContent(ctx context.Context, repo *content.Repository, result extsearch.SearchResult) (*content.Content, string) {
	externalSource := strings.TrimSpace(result.Source)
	externalContentID := strings.TrimSpace(result.ExternalContentID)
	if externalSource != "" && externalContentID != "" {
		existing, err := repo.FindByExternalSource(ctx, externalSource, externalContentID)
		if err == nil && existing != nil {
			return existing, "content:" + existing.ID.String()
		}
	}
	canonicalURL := strings.TrimSpace(result.CanonicalURL)
	if canonicalURL != "" {
		existing, err := repo.FindByCanonicalURL(ctx, canonicalURL)
		if err == nil && existing != nil {
			return existing, "content:" + existing.ID.String()
		}
		return nil, "url:" + canonicalURL
	}
	if externalSource != "" && externalContentID != "" {
		return nil, "ext:" + externalSource + ":" + externalContentID
	}
	return nil, ""
}

func recommendationDebugExternalResultKey(result extsearch.SearchResult) string {
	if trimmed := strings.TrimSpace(result.CanonicalURL); trimmed != "" {
		return "url:" + trimmed
	}
	if trimmed := strings.TrimSpace(result.ExternalContentID); trimmed != "" {
		return "ext:" + strings.TrimSpace(result.Source) + ":" + trimmed
	}
	return "title:" + strings.TrimSpace(strings.ToLower(result.Title))
}

func buildRecommendationDebugCandidateSummaries(diagnostics *curriculum.LessonCandidateSearchDiagnostics, preferredFormat *string, preferredLanguage string, now time.Time) []gin.H {
	if diagnostics == nil {
		return []gin.H{}
	}

	lexicalKeys := make(map[string]struct{}, len(diagnostics.LexicalCandidates))
	vectorKeys := make(map[string]struct{}, len(diagnostics.VectorCandidates))
	for _, candidate := range diagnostics.LexicalCandidates {
		lexicalKeys[recommendationDebugCandidateKey(candidate)] = struct{}{}
	}
	for _, candidate := range diagnostics.VectorCandidates {
		vectorKeys[recommendationDebugCandidateKey(candidate)] = struct{}{}
	}

	items := make([]gin.H, 0, len(diagnostics.MergedCandidates))
	for idx, candidate := range diagnostics.MergedCandidates {
		key := recommendationDebugCandidateKey(candidate)
		_, lexicalMatched := lexicalKeys[key]
		_, vectorMatched := vectorKeys[key]
		signals := recommendationDebugCandidateSignals(candidate, lexicalMatched, vectorMatched, preferredFormat, preferredLanguage, now)
		items = append(items, gin.H{
			"content_id":        candidate.ContentID,
			"title":             candidate.Title,
			"description":       candidate.Description,
			"thumbnail_url":     candidate.ThumbnailURL,
			"external_url":      candidate.ExternalURL,
			"rank_score":        candidate.RankScore,
			"content_type":      candidate.ContentType,
			"resource_type":     candidate.ResourceType,
			"quality_score":     candidate.QualityScore,
			"language":          candidate.Language,
			"created_at":        candidate.CreatedAt,
			"lexical_matched":   lexicalMatched,
			"vector_matched":    vectorMatched,
			"selection_reason":  recommendationDebugCandidateReason(candidate, lexicalMatched, vectorMatched, preferredFormat, preferredLanguage, now),
			"selection_signals": signals,
			"feature_snapshot":  recommendationDebugRankerFeatureSnapshot(candidate, "baseline", lexicalMatched, vectorMatched, idx+1, preferredLanguage, now),
		})
	}
	return items
}

func buildRecommendationDebugExternalCandidateSummaries(externalStage gin.H, maxPerLesson int) []gin.H {
	return buildRecommendationDebugExternalCandidateSummariesWithContext(externalStage, maxPerLesson, nil)
}

func buildRecommendationDebugExternalCandidateSummariesWithContext(externalStage gin.H, maxPerLesson int, explorerContext *curriculum.ExplorerRecommendationContext, primaryQuery ...string) []gin.H {
	rawResults, _ := externalStage["preview_results"].([]gin.H)
	if len(rawResults) == 0 {
		if anyResults, ok := externalStage["preview_results"].([]any); ok {
			rawResults = make([]gin.H, 0, len(anyResults))
			for _, value := range anyResults {
				rawResults = append(rawResults, recommendationDebugAnyMap(value))
			}
		}
	}
	candidates := make([]curriculum.ContentSearchCandidate, 0, len(rawResults))
	candidateMeta := make(map[string]gin.H, len(rawResults))
	for idx, result := range rawResults {
		eligible, _ := result["eligible_for_supplement"].(bool)
		if !eligible {
			continue
		}
		source := stringValue(result["source"])
		url := firstNonEmpty(stringValue(result["external_url"]), stringValue(result["url"]), stringValue(result["canonical_url"]))
		title := stringValue(result["title"])
		if title == "" || url == "" {
			continue
		}
		contentType := "external"
		switch source {
		case "youtube":
			contentType = "youtube"
		case "naver_blog":
			contentType = "blog"
		}
		rankScore := 0.01
		if idx < 20 {
			rankScore = float64(20-idx) / 21.0
		}
		candidate := curriculum.ContentSearchCandidate{
			ResourceType: curriculum.ResourceTypeExternal,
			ContentType:  contentType,
			Title:        title,
			Description:  recommendationDebugStringPointer(stringValue(result["description"])),
			ThumbnailURL: recommendationDebugStringPointer(stringValue(result["thumbnail_url"])),
			ExternalURL:  recommendationDebugStringPointer(url),
			PriceType:    curriculum.PriceTypeFree,
			Language:     firstNonEmpty(stringValue(result["language"]), "ko"),
			RankScore:    rankScore,
		}
		candidates = append(candidates, candidate)
		candidateMeta[recommendationDebugCandidateKey(candidate)] = result
	}
	if explorerContext != nil {
		query := ""
		if len(primaryQuery) > 0 {
			query = primaryQuery[0]
		}
		candidates = curriculum.RerankExplorerRecommendationCandidatesWithPrimaryQuery(candidates, *explorerContext, query, maxPerLesson)
	}
	if len(candidates) > maxPerLesson {
		candidates = candidates[:maxPerLesson]
	}

	items := make([]gin.H, 0, len(candidates))
	for _, candidate := range candidates {
		result := candidateMeta[recommendationDebugCandidateKey(candidate)]
		source := stringValue(result["source"])
		contentID := stringValue(result["content_id"])
		alreadySaved, _ := result["already_saved"].(bool)
		items = append(items, gin.H{
			"content_id":          contentID,
			"title":               candidate.Title,
			"description":         derefString(candidate.Description),
			"thumbnail_url":       derefString(candidate.ThumbnailURL),
			"external_url":        derefString(candidate.ExternalURL),
			"rank_score":          candidate.RankScore,
			"content_type":        candidate.ContentType,
			"resource_type":       "external",
			"quality_score":       0,
			"language":            candidate.Language,
			"lexical_matched":     false,
			"vector_matched":      false,
			"source_type":         "external",
			"external_source":     source,
			"external_content_id": stringValue(result["external_content_id"]),
			"author":              stringValue(result["author"]),
			"already_saved":       alreadySaved,
			"eligible_for_save":   !alreadySaved,
			"selection_reason":    "내부 추천 후보가 없어 외부 검색으로 보강된 후보",
			"selection_signals":   []string{"external_search"},
			"feature_snapshot": gin.H{
				"route":               "external_fallback",
				"rank_position":       len(items) + 1,
				"external_source":     source,
				"external_content_id": stringValue(result["external_content_id"]),
				"already_saved":       alreadySaved,
			},
		})
	}
	return items
}

func mergeRecommendationDebugCandidateSummaries(internalCandidates, externalCandidates []gin.H, limit int) []gin.H {
	if limit <= 0 {
		limit = 5
	}

	type rankedCandidate struct {
		item     gin.H
		internal bool
		index    int
		score    float64
	}

	seen := map[string]struct{}{}
	merged := make([]rankedCandidate, 0, len(internalCandidates)+len(externalCandidates))
	appendItems := func(items []gin.H, internal bool) {
		for idx, item := range items {
			key := recommendationDebugSummaryCandidateKey(item)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			score, _ := item["rank_score"].(float64)
			merged = append(merged, rankedCandidate{
				item:     item,
				internal: internal,
				index:    idx,
				score:    score,
			})
		}
	}

	appendItems(internalCandidates, true)
	appendItems(externalCandidates, false)

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].score == merged[j].score {
			if merged[i].internal == merged[j].internal {
				return merged[i].index < merged[j].index
			}
			return merged[i].internal
		}
		return merged[i].score > merged[j].score
	})

	if len(merged) > limit {
		merged = merged[:limit]
	}
	result := make([]gin.H, 0, len(merged))
	for idx, candidate := range merged {
		candidate.item["rank_position"] = idx + 1
		if feature, ok := candidate.item["feature_snapshot"].(gin.H); ok {
			feature["rank_position"] = idx + 1
		}
		result = append(result, candidate.item)
	}
	return result
}

func recommendationDebugSummaryCandidateKey(item gin.H) string {
	if contentID := stringValue(item["content_id"]); contentID != "" {
		return "content:" + contentID
	}
	if url := stringValue(item["external_url"]); url != "" {
		return "url:" + strings.TrimSpace(strings.ToLower(url))
	}
	if url := stringValue(item["url"]); url != "" {
		return "url:" + strings.TrimSpace(strings.ToLower(url))
	}
	return "title:" + strings.TrimSpace(strings.ToLower(stringValue(item["title"])))
}

func recommendationDebugStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func recommendationDebugCandidateKey(candidate curriculum.ContentSearchCandidate) string {
	if candidate.ContentID != nil {
		return "content:" + candidate.ContentID.String()
	}
	if candidate.ExternalURL != nil && strings.TrimSpace(*candidate.ExternalURL) != "" {
		return "url:" + strings.TrimSpace(*candidate.ExternalURL)
	}
	return "title:" + strings.TrimSpace(strings.ToLower(candidate.Title))
}

func recommendationDebugCandidateReason(candidate curriculum.ContentSearchCandidate, lexicalMatched, vectorMatched bool, preferredFormat *string, preferredLanguage string, now time.Time) string {
	reasons := make([]string, 0, 5)
	if lexicalMatched {
		reasons = append(reasons, "lexical 매칭으로 후보에 포함")
	}
	if vectorMatched {
		reasons = append(reasons, "임베딩 유사도로 보강")
	}
	if candidate.QualityScore >= 0.8 {
		reasons = append(reasons, "quality_score가 높음")
	}
	if strings.EqualFold(strings.TrimSpace(candidate.Language), strings.TrimSpace(preferredLanguage)) {
		reasons = append(reasons, "선호 언어와 일치")
	}
	if recommendationDebugPreferredFormatMatch(candidate.ContentType, preferredFormat) {
		reasons = append(reasons, "선호 형식과 일치")
	}
	if recommendationDebugRecentContent(candidate.CreatedAt, now) {
		reasons = append(reasons, "최근 콘텐츠 가산점")
	}
	if len(reasons) == 0 {
		return "기본 랭킹 점수로 최종 후보에 포함"
	}
	return strings.Join(reasons, ", ")
}

func recommendationDebugCandidateSignals(candidate curriculum.ContentSearchCandidate, lexicalMatched, vectorMatched bool, preferredFormat *string, preferredLanguage string, now time.Time) []string {
	signals := make([]string, 0, 6)
	if lexicalMatched {
		signals = append(signals, "lexical")
	}
	if vectorMatched {
		signals = append(signals, "vector")
	}
	if candidate.QualityScore >= 0.8 {
		signals = append(signals, "quality_high")
	}
	if strings.EqualFold(strings.TrimSpace(candidate.Language), strings.TrimSpace(preferredLanguage)) {
		signals = append(signals, "language_match")
	}
	if recommendationDebugPreferredFormatMatch(candidate.ContentType, preferredFormat) {
		signals = append(signals, "preferred_format_match")
	}
	if recommendationDebugRecentContent(candidate.CreatedAt, now) {
		signals = append(signals, "recent")
	}
	if len(signals) == 0 {
		signals = append(signals, "base_rank")
	}
	return signals
}

func recommendationDebugPreferredFormatMatch(contentType string, preferredFormat *string) bool {
	if preferredFormat == nil {
		return false
	}
	format := strings.ToLower(strings.TrimSpace(*preferredFormat))
	if format == "" {
		return false
	}
	switch {
	case strings.Contains(format, "video"), strings.Contains(format, "영상"), strings.Contains(format, "유튜브"):
		return contentType == "youtube"
	case strings.Contains(format, "text"), strings.Contains(format, "글"), strings.Contains(format, "읽기"), strings.Contains(format, "article"), strings.Contains(format, "blog"):
		return contentType == "blog" || contentType == "article" || contentType == "internal"
	default:
		return false
	}
}

func recommendationDebugRecentContent(createdAt, now time.Time) bool {
	if createdAt.IsZero() {
		return false
	}
	if now.IsZero() {
		now = time.Now()
	}
	return now.Sub(createdAt) <= 90*24*time.Hour
}
