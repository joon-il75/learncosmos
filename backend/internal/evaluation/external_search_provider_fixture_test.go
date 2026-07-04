package evaluation

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/learnweaver/backend/internal/domain/curriculum"
)

type externalSearchProviderFixture struct {
	ID         string                                    `json:"id"`
	Provider   string                                    `json:"provider"`
	Query      string                                    `json:"query"`
	SearchSpec curriculum.LessonRecommendationSearchSpec `json:"search_spec"`
	Candidates []externalSearchProviderFixtureCandidate  `json:"candidates"`
}

type externalSearchProviderFixtureCandidate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Relevant    bool   `json:"relevant"`
}

func TestExternalSearchProviderFixturesPreferRelevantCandidates(t *testing.T) {
	fixtures := loadExternalSearchProviderFixtures(t)
	if len(fixtures) < 5 {
		t.Fatalf("fixture count = %d, want at least 5", len(fixtures))
	}

	providers := map[string]bool{}
	stageRoles := map[string]bool{}
	for _, fixture := range fixtures {
		providers[fixture.Provider] = true
		stageRoles[fixture.SearchSpec.StageRole] = true
		if fixture.Query == "" {
			t.Fatalf("%s missing query", fixture.ID)
		}
		if fixture.SearchSpec.PrimaryQuery == "" {
			t.Fatalf("%s missing search spec primary query", fixture.ID)
		}

		bestRelevant := -1 << 30
		bestNoise := -1 << 30
		for _, candidate := range fixture.Candidates {
			description := candidate.Description
			score := curriculum.ScoreCandidateWithLessonSearchSpec(curriculum.ContentSearchCandidate{
				Title:       candidate.Title,
				Description: &description,
				ContentType: candidate.ContentType,
				Language:    fixture.SearchSpec.Language,
			}, fixture.SearchSpec)
			if candidate.Relevant {
				if score > bestRelevant {
					bestRelevant = score
				}
			} else if score > bestNoise {
				bestNoise = score
			}
		}
		if bestRelevant <= 0 {
			t.Fatalf("%s best relevant score = %d, want positive", fixture.ID, bestRelevant)
		}
		if bestRelevant <= bestNoise {
			t.Fatalf("%s relevant candidate should outrank noise, relevant=%d noise=%d", fixture.ID, bestRelevant, bestNoise)
		}
	}

	for _, provider := range []string{"youtube", "naver_blog"} {
		if !providers[provider] {
			t.Fatalf("missing provider fixture %q", provider)
		}
	}
	for _, stageRole := range []string{"setup_intro", "first_output", "core_pattern", "practice_loop", "artifact_finish", "troubleshooting", "portfolio_publish"} {
		if !stageRoles[stageRole] {
			t.Fatalf("missing stage role fixture %q", stageRole)
		}
	}
}

func loadExternalSearchProviderFixtures(t *testing.T) []externalSearchProviderFixture {
	t.Helper()
	raw, err := os.ReadFile("testdata/external_search_provider_fixtures.json")
	if err != nil {
		t.Fatalf("read provider fixtures: %v", err)
	}
	var fixtures []externalSearchProviderFixture
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		t.Fatalf("parse provider fixtures: %v", err)
	}
	return fixtures
}
