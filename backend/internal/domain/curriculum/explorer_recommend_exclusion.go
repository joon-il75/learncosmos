package curriculum

import (
	"strings"

	"github.com/learnweaver/backend/internal/pkg/extsearch"
)

type explorerRecommendationExclusion struct {
	contentIDs map[string]struct{}
	urlKeys    map[string]struct{}
}

func buildExplorerRecommendationExclusion(req RecommendExplorerContentRequest) explorerRecommendationExclusion {
	exclusion := explorerRecommendationExclusion{
		contentIDs: make(map[string]struct{}),
		urlKeys:    make(map[string]struct{}),
	}
	for _, rawID := range req.ExcludedContentIDs {
		id := strings.TrimSpace(strings.ToLower(rawID))
		if id != "" {
			exclusion.contentIDs[id] = struct{}{}
		}
	}
	for _, rawURL := range req.ExcludedURLs {
		key := recommendationBlockURLKey(rawURL)
		if key != "" {
			exclusion.urlKeys[key] = struct{}{}
		}
	}
	return exclusion
}

func (e explorerRecommendationExclusion) empty() bool {
	return len(e.contentIDs) == 0 && len(e.urlKeys) == 0
}

func (e explorerRecommendationExclusion) blocked(candidate ContentSearchCandidate) bool {
	if candidate.ContentID != nil {
		if _, ok := e.contentIDs[strings.ToLower(candidate.ContentID.String())]; ok {
			return true
		}
	}
	if candidate.ExternalURL != nil {
		if _, ok := e.urlKeys[recommendationBlockURLKey(*candidate.ExternalURL)]; ok {
			return true
		}
	}
	return false
}

func (e explorerRecommendationExclusion) blockedExternalResult(result extsearch.SearchResult) bool {
	if e.empty() {
		return false
	}
	if _, ok := e.urlKeys[recommendationBlockURLKey(result.URL)]; ok {
		return true
	}
	if _, ok := e.urlKeys[recommendationBlockURLKey(result.CanonicalURL)]; ok {
		return true
	}
	return false
}

func filterExcludedExplorerCandidates(candidates []ContentSearchCandidate, exclusion explorerRecommendationExclusion) []ContentSearchCandidate {
	if exclusion.empty() || len(candidates) == 0 {
		return candidates
	}
	filtered := make([]ContentSearchCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if exclusion.blocked(candidate) {
			continue
		}
		filtered = append(filtered, candidate)
	}
	return filtered
}
