package curriculum

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type curriculumPatternMatchResult struct {
	Matches   []CurriculumPatternMatch
	CacheHit  bool
	QueryHash string
}

func (h *Handler) resolveCurriculumPatternMatches(ctx context.Context, req CreateCourseDraftRequest, limit int) (*curriculumPatternMatchResult, error) {
	queryText := BuildCurriculumPatternQueryText(req)
	if strings.TrimSpace(queryText) == "" {
		return nil, fmt.Errorf("empty_query")
	}
	if limit <= 0 {
		limit = 3
	}

	provider := "embedding_gemma"
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	dimension := getEnvInt("EMBEDDING_GEMMA_DIMENSION", 768)
	language := normalizeLearningLanguage(req.LearningLanguage)
	queryHash := hashCurriculumPatternMatchQuery(queryText, language, provider, model, dimension)
	cacheKey := "curriculum_pattern_match:v1:" + queryHash

	if h.redisClient != nil {
		if cached, err := h.redisClient.Get(ctx, cacheKey).Result(); err == nil && strings.TrimSpace(cached) != "" {
			var matches []CurriculumPatternMatch
			if jsonErr := json.Unmarshal([]byte(cached), &matches); jsonErr == nil {
				return &curriculumPatternMatchResult{Matches: matches, CacheHit: true, QueryHash: queryHash}, nil
			}
		}
	}

	endpoint := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT"))
	if endpoint == "" {
		return nil, fmt.Errorf("embedding_endpoint_missing")
	}
	client, err := embedder.NewProviderClient(provider, endpoint)
	if err != nil {
		return nil, fmt.Errorf("embedder_init_failed: %w", err)
	}
	vector, err := client.Embed(ctx, queryText)
	if err != nil {
		return nil, fmt.Errorf("query_embedding_failed: %w", err)
	}
	if len(vector) != dimension {
		return nil, fmt.Errorf("query_embedding_dimension_mismatch: got=%d want=%d", len(vector), dimension)
	}

	matches, err := h.repo.SearchCurriculumPatternMatches(ctx, FormatFloat32VectorForPG(vector), provider, model, dimension, language, limit)
	if err != nil {
		return nil, fmt.Errorf("pattern_search_failed: %w", err)
	}
	if h.redisClient != nil && len(matches) > 0 {
		if payload, err := json.Marshal(matches); err == nil {
			ttl := time.Duration(getEnvInt("CURRICULUM_PATTERN_CACHE_TTL_SECONDS", 86400)) * time.Second
			if ttl > 0 {
				if err := h.redisClient.Set(ctx, cacheKey, string(payload), ttl).Err(); err != nil {
					log.Printf("[curriculum/pattern_match] mode=cache_set_failed err=%q", truncateShadowError(err))
				}
			}
		}
	}
	return &curriculumPatternMatchResult{Matches: matches, CacheHit: false, QueryHash: queryHash}, nil
}

func (h *Handler) logCurriculumPatternMatchResult(userID uuid.UUID, mode, codePatternKey, codeSubpatternKey, refinementMode string, result *curriculumPatternMatchResult, startedAt time.Time) {
	topKeys := make([]string, 0, len(result.Matches))
	codeRank := 0
	for idx, match := range result.Matches {
		topKeys = append(topKeys, match.PatternKey)
		if strings.TrimSpace(codeSubpatternKey) != "" && match.PatternKey == codeSubpatternKey {
			codeRank = idx + 1
		}
	}
	top1 := ""
	if len(topKeys) > 0 {
		top1 = topKeys[0]
	}
	log.Printf(
		"[curriculum/pattern_shadow] user=%s mode=%s code_pattern=%s code_subpattern=%s refinement_mode=%s embedding_top1=%s embedding_top3=%s code_rank=%d top3_matched=%t cache_hit=%t duration_ms=%d",
		userID,
		mode,
		codePatternKey,
		codeSubpatternKey,
		refinementMode,
		top1,
		strings.Join(topKeys, ","),
		codeRank,
		codeRank > 0 && codeRank <= 3,
		result.CacheHit,
		time.Since(startedAt).Milliseconds(),
	)
}

func hashCurriculumPatternMatchQuery(queryText, language, provider, model string, dimension int) string {
	segments := []string{
		strings.TrimSpace(queryText),
		"language=" + normalizeLearningLanguage(language),
		"provider=" + strings.TrimSpace(provider),
		"model=" + strings.TrimSpace(model),
		fmt.Sprintf("dimension=%d", dimension),
	}
	sum := sha256.Sum256([]byte(strings.Join(segments, "\n")))
	return hex.EncodeToString(sum[:])
}
