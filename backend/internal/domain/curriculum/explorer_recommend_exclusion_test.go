package curriculum

import (
	"testing"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
)

func TestFilterExcludedExplorerCandidatesRemovesContentID(t *testing.T) {
	blockedID := uuid.New()
	keptID := uuid.New()
	exclusion := buildExplorerRecommendationExclusion(RecommendExplorerContentRequest{
		ExcludedContentIDs: []string{blockedID.String()},
	})

	candidates := []ContentSearchCandidate{
		{ContentID: &blockedID, Title: "blocked"},
		{ContentID: &keptID, Title: "kept"},
	}

	filtered := filterExcludedExplorerCandidates(candidates, exclusion)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(filtered))
	}
	if filtered[0].Title != "kept" {
		t.Fatalf("expected kept candidate, got %q", filtered[0].Title)
	}
}

func TestExplorerRecommendationExclusionBlocksExternalResultURL(t *testing.T) {
	exclusion := buildExplorerRecommendationExclusion(RecommendExplorerContentRequest{
		ExcludedURLs: []string{"https://example.com/watch/abc#section"},
	})

	if !exclusion.blockedExternalResult(extsearch.SearchResult{URL: "https://example.com/watch/abc"}) {
		t.Fatal("expected external result URL to be blocked")
	}
	if exclusion.blockedExternalResult(extsearch.SearchResult{URL: "https://example.com/watch/other"}) {
		t.Fatal("did not expect different external result URL to be blocked")
	}
}

func TestFilterExcludedExplorerCandidatesRemovesNormalizedURL(t *testing.T) {
	blockedURL := "https://www.youtube.com/watch?v=abc123"
	keptURL := "https://www.youtube.com/watch?v=def456"
	exclusion := buildExplorerRecommendationExclusion(RecommendExplorerContentRequest{
		ExcludedURLs: []string{"HTTPS://www.youtube.com/watch?v=abc123#watch"},
	})

	candidates := []ContentSearchCandidate{
		{ExternalURL: &blockedURL, Title: "blocked"},
		{ExternalURL: &keptURL, Title: "kept"},
	}

	filtered := filterExcludedExplorerCandidates(candidates, exclusion)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(filtered))
	}
	if filtered[0].Title != "kept" {
		t.Fatalf("expected kept candidate, got %q", filtered[0].Title)
	}
}
