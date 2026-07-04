package ranker

import (
	"context"
	"sort"
	"strings"
)

type CandidateFeature map[string]any

type Score struct {
	Index  int     `json:"index"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

type Scorer interface {
	Score(ctx context.Context, features []CandidateFeature) ([]Score, error)
	Provider() string
	ModelVersion() string
	FeatureSchemaVersion() string
}

type NoOpScorer struct {
	modelVersion         string
	featureSchemaVersion string
}

func NewNoOpScorer(modelVersion, featureSchemaVersion string) NoOpScorer {
	return NoOpScorer{
		modelVersion:         strings.TrimSpace(modelVersion),
		featureSchemaVersion: strings.TrimSpace(featureSchemaVersion),
	}
}

func (s NoOpScorer) Provider() string {
	return "noop"
}

func (s NoOpScorer) ModelVersion() string {
	if s.modelVersion == "" {
		return "noop-ranker"
	}
	return s.modelVersion
}

func (s NoOpScorer) FeatureSchemaVersion() string {
	if s.featureSchemaVersion == "" {
		return "ranker-feature-v1"
	}
	return s.featureSchemaVersion
}

func (s NoOpScorer) Score(_ context.Context, features []CandidateFeature) ([]Score, error) {
	scores := make([]Score, 0, len(features))
	for idx, feature := range features {
		routeRank := intFeature(feature, "route_rank", idx+1)
		if routeRank <= 0 {
			routeRank = idx + 1
		}
		scores = append(scores, Score{
			Index:  idx,
			Score:  1 / float64(routeRank),
			Reason: "noop ranker preserves route order",
		})
	}
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Index < scores[j].Index
	})
	return scores, nil
}

func intFeature(feature CandidateFeature, key string, fallback int) int {
	value, ok := feature[key]
	if !ok {
		return fallback
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	default:
		return fallback
	}
}
