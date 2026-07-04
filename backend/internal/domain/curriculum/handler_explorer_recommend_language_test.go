package curriculum

import "testing"

func TestResolveExplorerRecommendationLanguage_PrefersDraftGenerationLanguage(t *testing.T) {
	draft := &CourseDraft{GenerationLanguage: "en"}

	got := resolveExplorerRecommendationLanguage("ko", draft)
	if got != "en" {
		t.Fatalf("language = %q, want en", got)
	}
}

func TestResolveExplorerRecommendationLanguage_FallsBackToUserLanguage(t *testing.T) {
	got := resolveExplorerRecommendationLanguage("en", nil)
	if got != "en" {
		t.Fatalf("language = %q, want en", got)
	}
}

func TestSharedPointStatuses_UsesArchivedDrafts(t *testing.T) {
	statuses := sharedPointStatuses()
	if len(statuses) != 1 || statuses[0] != DraftStatusArchived {
		t.Fatalf("sharedPointStatuses = %#v, want archived only", statuses)
	}
}
