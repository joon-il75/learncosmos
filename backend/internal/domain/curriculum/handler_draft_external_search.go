package curriculum

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func (h *Handler) searchExternalLessonCandidates(ctx context.Context, _ uuid.UUID, limit int, queryBundle normalizer.LessonSearchQuery, preferredFormat *string) ([]ContentSearchCandidate, error) {
	return h.searchExternalLessonCandidatesWithQuery(ctx, limit, BuildExternalSearchQuery(queryBundle), queryBundle, preferredFormat, "ko")
}

func (h *Handler) searchExternalLessonCandidatesWithQuery(ctx context.Context, limit int, query string, queryBundle normalizer.LessonSearchQuery, preferredFormat *string, preferredLanguage string) ([]ContentSearchCandidate, error) {
	return h.searchExternalLessonCandidatesWithQueryFiltered(ctx, limit, query, queryBundle, preferredFormat, preferredLanguage, nil)
}

func (h *Handler) searchExternalLessonCandidatesWithQueryFiltered(ctx context.Context, limit int, query string, queryBundle normalizer.LessonSearchQuery, preferredFormat *string, preferredLanguage string, reject func(extsearch.SearchResult) bool) ([]ContentSearchCandidate, error) {
	if limit <= 0 {
		return nil, nil
	}

	query = strings.TrimSpace(query)
	if query == "" {
		query = BuildExternalSearchQuery(queryBundle)
	}
	preferredLanguage = normalizeLearningLanguage(preferredLanguage)
	providers := []extsearch.SearchProvider{}
	if strings.TrimSpace(h.youtubeAPIKey) != "" {
		providers = append(providers, extsearch.NewYouTubeProviderWithLanguage(h.youtubeAPIKey, preferredLanguage))
	}
	if preferredLanguage != "en" && strings.TrimSpace(h.naverClientID) != "" && strings.TrimSpace(h.naverSecret) != "" {
		providers = append(providers, extsearch.NewNaverBlogProvider(h.naverClientID, h.naverSecret))
	}
	if len(providers) == 0 {
		return nil, nil
	}

	blockedURLKeys, blockErr := h.repo.ListActiveRecommendationBlockedURLKeys(ctx)
	if blockErr != nil {
		return nil, blockErr
	}

	acceptedKeys := map[string]struct{}{}
	results, providerStats := CollectFilteredExternalSearchResults(ctx, providers, query, func(result extsearch.SearchResult) bool {
		if _, blocked := blockedURLKeys[recommendationBlockURLKey(result.URL)]; blocked {
			return false
		}
		if _, blocked := blockedURLKeys[recommendationBlockURLKey(result.CanonicalURL)]; blocked {
			return false
		}
		if reject != nil && reject(result) {
			return false
		}
		if !isInstructionalExternalSearchResultForQuery(result, queryBundle) {
			return false
		}
		key := externalSearchResultKey(result)
		if _, exists := acceptedKeys[key]; exists {
			return false
		}
		acceptedKeys[key] = struct{}{}
		return true
	})
	for _, stat := range providerStats {
		log.Printf(
			"[lesson/external-search] provider=%s query=%q pages=%d raw=%d accepted=%d filtered=%d error=%q",
			stat.Provider,
			stat.Query,
			stat.PagesAttempted,
			stat.RawCount,
			stat.AcceptedCount,
			stat.FilteredCount,
			stat.ErrorReason,
		)
	}

	candidates := make([]ContentSearchCandidate, 0, limit)
	for idx, result := range results {
		if len(candidates) >= limit {
			break
		}

		rankScore := scoreExternalSearchResult(queryBundle, result, preferredFormat, idx)
		candidates = append(candidates, ContentSearchCandidate{
			ContentID:    nil,
			ResourceType: ResourceTypeExternal,
			ContentType:  result.Source,
			Title:        result.Title,
			Description:  stringPointer(result.Description),
			ThumbnailURL: stringPointer(result.ThumbnailURL),
			ExternalURL:  stringPointer(result.URL),
			PriceType:    PriceTypeFree,
			Language:     result.Language,
			RankScore:    rankScore,
		})
	}

	return candidates, nil
}

func mergeExternalCandidates(internal, external []ContentSearchCandidate, limit int) []ContentSearchCandidate {
	if limit <= 0 {
		limit = 5
	}

	type mergedCandidate struct {
		candidate    ContentSearchCandidate
		internal     bool
		originalRank int
	}

	merged := make([]mergedCandidate, 0, len(internal)+len(external))
	seen := map[string]struct{}{}
	appendUnique := func(items []ContentSearchCandidate, internal bool) {
		for idx, item := range items {
			key := recommendationCandidateDedupeKey(item)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, mergedCandidate{
				candidate:    item,
				internal:     internal,
				originalRank: idx,
			})
		}
	}
	appendUnique(internal, true)
	appendUnique(external, false)

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].candidate.RankScore == merged[j].candidate.RankScore {
			if merged[i].internal == merged[j].internal {
				return merged[i].originalRank < merged[j].originalRank
			}
			return merged[i].internal
		}
		return merged[i].candidate.RankScore > merged[j].candidate.RankScore
	})

	if len(merged) > limit {
		merged = merged[:limit]
	}

	result := make([]ContentSearchCandidate, 0, len(merged))
	for _, item := range merged {
		result = append(result, item.candidate)
	}
	return result
}
