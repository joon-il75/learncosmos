package main

import (
	"strings"
	"testing"

	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func TestVisualArtEvaluationContentSeedsCoverWatercolorAndPencilDrawing(t *testing.T) {
	if len(visualArtEvaluationContentSeeds) != 6 {
		t.Fatalf("seed count = %d, want 6", len(visualArtEvaluationContentSeeds))
	}

	joined := strings.Builder{}
	seen := map[string]bool{}
	for _, seed := range visualArtEvaluationContentSeeds {
		if seed.ExternalID == "" || seed.Title == "" || seed.ContentType == "" {
			t.Fatalf("seed has required field missing: %#v", seed)
		}
		if seen[seed.ExternalID] {
			t.Fatalf("duplicate seed external id: %s", seed.ExternalID)
		}
		seen[seed.ExternalID] = true
		joined.WriteString(normalizer.BuildSearchTextKO(seed.Title, seed.Description, seed.Author, seed.ContentType, seed.Language))
		joined.WriteByte(' ')
	}

	searchText := joined.String()
	for _, want := range []string{"수채화", "물조절", "붓", "번짐", "드로잉", "명암", "연필", "스케치"} {
		if !strings.Contains(searchText, want) {
			t.Fatalf("visual-art seed search text missing %q: %s", want, searchText)
		}
	}
}
