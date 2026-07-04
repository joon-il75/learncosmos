package admin

import (
	"os"
	"strings"
	"testing"
)

func TestRecommendationDebugEmbeddingMetadataUsesPersistedErrors(t *testing.T) {
	source, err := os.ReadFile("recommendation_debug_embedding_helpers.go")
	if err != nil {
		t.Fatalf("read recommendation_debug_embedding_helpers.go: %v", err)
	}
	text := string(source)
	for _, want := range []string{
		`"error":              logsafe.PersistedError(err.Error())`,
		`"reason":      logsafe.PersistedError(err.Error())`,
		`"endpoint":           logsafe.URL(endpoint)`,
		`"endpoint":    logsafe.URL(endpoint)`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("recommendation debug embedding metadata missing %s", want)
		}
	}
	for _, forbidden := range []string{
		`"error":              err.Error()`,
		`"reason":      err.Error()`,
		`"endpoint":           endpoint`,
		`"endpoint":    endpoint`,
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("recommendation debug embedding metadata must not use raw %s", forbidden)
		}
	}
}
