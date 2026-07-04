package curriculum

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) SearchLessonCandidates(ctx context.Context, userID uuid.UUID, query string, embedding *string, limit int, options LessonCandidateSearchOptions) ([]ContentSearchCandidate, error) {
	diagnostics, err := r.SearchLessonCandidatesDiagnostics(ctx, userID, query, embedding, limit, options)
	if err != nil {
		return nil, err
	}
	return diagnostics.MergedCandidates, nil
}

func (r *Repository) SearchLessonCandidatesDiagnostics(ctx context.Context, userID uuid.UUID, query string, embedding *string, limit int, options LessonCandidateSearchOptions) (*LessonCandidateSearchDiagnostics, error) {
	if limit <= 0 {
		limit = 5
	}

	retrievalLimit := lessonRetrievalLimit(limit)

	lexicalCandidates, err := r.searchLexicalLessonCandidates(ctx, userID, query, retrievalLimit)
	if err != nil {
		return nil, err
	}

	vectorCandidates, err := r.searchVectorLessonCandidates(ctx, userID, embedding, retrievalLimit)
	if err != nil {
		return nil, err
	}

	return &LessonCandidateSearchDiagnostics{
		LexicalCandidates: lexicalCandidates,
		VectorCandidates:  vectorCandidates,
		MergedCandidates:  mergeLessonCandidates(lexicalCandidates, vectorCandidates, limit, options),
		RetrievalLimit:    retrievalLimit,
	}, nil
}

func lessonRetrievalLimit(limit int) int {
	if limit <= 0 {
		return 30
	}
	if limit < 30 {
		return 30
	}
	return limit
}

func primaryEmbeddingProvider() string {
	return "embedding_gemma"
}

func primaryEmbeddingModel() string {
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		return "embedding-gemma"
	}
	return model
}

func primaryEmbeddingDimension() int {
	raw := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_DIMENSION"))
	if raw == "" {
		return 768
	}
	dimension, err := strconv.Atoi(raw)
	if err != nil || dimension <= 0 {
		return 768
	}
	return dimension
}
