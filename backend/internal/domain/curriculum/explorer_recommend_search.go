package curriculum

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func (h *Handler) searchExplorerInternalCandidatesByStage(
	ctx context.Context,
	userID uuid.UUID,
	stages []explorerSearchStage,
	limit int,
	options LessonCandidateSearchOptions,
	embedClient embedder.Embedder,
	provider string,
	billingStatus string,
) ([]ContentSearchCandidate, error) {
	if len(stages) == 0 || limit <= 0 {
		return []ContentSearchCandidate{}, nil
	}
	_ = provider
	_ = billingStatus

	filtered := make([]ContentSearchCandidate, 0, limit)
	seen := make(map[string]struct{}, limit)
	courseFilter := explorerFilterStage{}
	hasCourseFilter := false
	for _, stage := range stages {
		if stage.name == "course" {
			courseFilter = stage.filter
			hasCourseFilter = true
			break
		}
	}

	for _, stage := range stages {
		remaining := limit - len(filtered)
		if remaining <= 0 {
			break
		}

		lexicalQuery := stage.bundle.LexicalQueryExpanded
		if strings.TrimSpace(lexicalQuery) == "" {
			lexicalQuery = stage.bundle.LexicalQueryKO
		}
		if strings.TrimSpace(lexicalQuery) == "" {
			continue
		}

		embeddingText, err := buildExplorerStageEmbedding(ctx, stage.bundle, embedClient)
		if err != nil {
			embeddingText = nil
		}
		stageCandidates, err := h.repo.SearchLessonCandidates(
			ctx,
			userID,
			lexicalQuery,
			embeddingText,
			max(remaining*4, remaining),
			options,
		)
		if err != nil {
			return nil, err
		}

		if hasCourseFilter && stage.name != "course" {
			stageCandidates = filterExplorerCandidatesByAllStages(stageCandidates, []explorerFilterStage{stage.filter, courseFilter}, remaining)
		} else {
			stageCandidates = filterExplorerCandidatesByStage(stageCandidates, []explorerFilterStage{stage.filter}, remaining)
		}
		log.Printf(
			"[explorer/recommend] stage=%s survivors=%d after_course_gate=%t",
			stage.name,
			len(stageCandidates),
			hasCourseFilter && stage.name != "course",
		)
		for _, candidate := range stageCandidates {
			key := explorerCandidateFilterKey(candidate)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			filtered = append(filtered, candidate)
			if len(filtered) >= limit {
				break
			}
		}
	}

	return filtered, nil
}

func buildExplorerStageEmbedding(
	ctx context.Context,
	bundle normalizer.LessonSearchQuery,
	embedClient embedder.Embedder,
) (*string, error) {
	if embedClient == nil || strings.TrimSpace(bundle.DenseQuery) == "" {
		return nil, nil
	}

	embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	vector, err := embedClient.Embed(embedCtx, bundle.DenseQuery)
	if err != nil || len(vector) == 0 {
		return nil, err
	}

	pgvec := float32VectorToPGVector(vector)
	return &pgvec, nil
}

// explorerTokenMatchesText checks whether token appears in searchText.
// Short tokens (≤2 runes, e.g. "코드", "기타") require whole-word matching
// so that "코드" does not match compound words like "qr코드" or "상품코드".
// Longer tokens use substring matching which is safe for distinctive terms.
func explorerTokenMatchesText(searchText, token string) bool {
	if len([]rune(token)) > 2 {
		return strings.Contains(searchText, token)
	}
	for _, word := range strings.Fields(searchText) {
		if normalizeExplorerFilterToken(word) == token {
			return true
		}
	}
	return false
}
