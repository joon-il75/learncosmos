package curriculum

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func (r *Repository) searchLexicalLessonCandidates(ctx context.Context, userID uuid.UUID, query string, limit int) ([]ContentSearchCandidate, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []ContentSearchCandidate{}, nil
	}

	fallbackTokens := lexicalFallbackTokens(query)
	fallbackThreshold := lexicalFallbackThreshold(fallbackTokens)

	rows, err := r.pool.Query(ctx, `
		WITH lexical_query AS (
		  SELECT $2::text AS query_text,
		         $4::text[] AS fallback_tokens,
		         $5::int AS fallback_threshold
		)
		SELECT id, content_type, title, description, thumbnail_url, url, language,
		       COALESCE(quality_score, 0.5) AS quality_score, created_at,
		       (
		         similarity(search_text_ko, lq.query_text) * 0.75
		         + CASE
		             WHEN lq.fallback_threshold > 0 THEN
		               LEAST((SELECT COUNT(*) FROM unnest(lq.fallback_tokens) AS token WHERE search_text_ko LIKE '%' || token || '%'), lq.fallback_threshold)::float
		               / lq.fallback_threshold::float * 0.20
		             ELSE 0
		           END
		         + COALESCE(quality_score, 0.5) * 0.05
		       ) AS rank_score
		FROM contents
		CROSS JOIN lexical_query lq
		WHERE ($1 = '00000000-0000-0000-0000-000000000000'::uuid OR user_id = $1 OR is_public = true)
		  AND content_status IN ('active', 'redirect')
		  AND health_score > 0.6
		  AND (
		    search_text_ko % lq.query_text
		    OR (
		      lq.fallback_threshold > 0
		      AND (SELECT COUNT(*) FROM unnest(lq.fallback_tokens) AS token WHERE search_text_ko LIKE '%' || token || '%') >= lq.fallback_threshold
		    )
		  )
		  AND NOT EXISTS (
		    SELECT 1
		    FROM content_recommendation_blocks crb
		    WHERE crb.active = true
		      AND (
		        crb.content_id = contents.id
		        OR (COALESCE(crb.url_key, '') <> '' AND crb.url_key IN (
		          lower(trim(COALESCE(contents.url, ''))),
		          lower(trim(COALESCE(contents.canonical_url, '')))
		        ))
		      )
		  )
		ORDER BY rank_score DESC, created_at DESC
		LIMIT $3
	`, userID, query, limit, fallbackTokens, fallbackThreshold)
	if err != nil {
		return nil, fmt.Errorf("search lexical lesson candidates: %w", err)
	}
	defer rows.Close()

	candidates := []ContentSearchCandidate{}
	for rows.Next() {
		candidate, err := scanContentSearchCandidate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan lexical search candidate: %w", err)
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}

func (r *Repository) searchVectorLessonCandidates(ctx context.Context, userID uuid.UUID, embedding *string, limit int) ([]ContentSearchCandidate, error) {
	if embedding == nil || strings.TrimSpace(*embedding) == "" {
		return []ContentSearchCandidate{}, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.content_type, c.title, c.description, c.thumbnail_url, c.url, c.language,
		       COALESCE(c.quality_score, 0.5) AS quality_score, c.created_at,
		       (
		         (1 - (ce.embedding <=> $5::vector)) * 0.9
		         + COALESCE(c.quality_score, 0.5) * 0.1
		       ) AS rank_score
		FROM content_embeddings ce
		JOIN contents c ON c.id = ce.content_id
		WHERE ($1 = '00000000-0000-0000-0000-000000000000'::uuid OR c.user_id = $1 OR c.is_public = true)
		  AND c.content_status IN ('active', 'redirect')
		  AND c.health_score > 0.6
		  AND ce.provider = $2
		  AND ce.model = $3
		  AND ce.dimension = $4
		  AND ce.status = 'ready'
		  AND ce.embedding IS NOT NULL
		  AND NOT EXISTS (
		    SELECT 1
		    FROM content_recommendation_blocks crb
		    WHERE crb.active = true
		      AND (
		        crb.content_id = c.id
		        OR (COALESCE(crb.url_key, '') <> '' AND crb.url_key IN (
		          lower(trim(COALESCE(c.url, ''))),
		          lower(trim(COALESCE(c.canonical_url, '')))
		        ))
		      )
		  )
		ORDER BY rank_score DESC, c.created_at DESC
		LIMIT $6
	`, userID, primaryEmbeddingProvider(), primaryEmbeddingModel(), primaryEmbeddingDimension(), *embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("search vector lesson candidates: %w", err)
	}
	defer rows.Close()

	candidates := []ContentSearchCandidate{}
	for rows.Next() {
		candidate, err := scanContentSearchCandidate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan vector search candidate: %w", err)
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}

func (r *Repository) SearchShadowVectorLessonCandidates(ctx context.Context, userID uuid.UUID, provider, model string, dimension int, embedding *string, limit int) ([]ContentSearchCandidate, error) {
	if embedding == nil || strings.TrimSpace(*embedding) == "" || strings.TrimSpace(provider) == "" || strings.TrimSpace(model) == "" || dimension <= 0 {
		return []ContentSearchCandidate{}, nil
	}
	if limit <= 0 {
		limit = 5
	}
	retrievalLimit := lessonRetrievalLimit(limit)

	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.content_type, c.title, c.description, c.thumbnail_url, c.url, c.language,
		       COALESCE(c.quality_score, 0.5) AS quality_score, c.created_at,
		       (
		         (1 - (ces.embedding <=> $5::vector)) * 0.9
		         + COALESCE(c.quality_score, 0.5) * 0.1
		       ) AS rank_score
		FROM content_embeddings ces
		JOIN contents c ON c.id = ces.content_id
		WHERE (c.user_id = $1 OR c.is_public = true)
		  AND c.content_status IN ('active', 'redirect')
		  AND c.health_score > 0.6
		  AND ces.provider = $2
		  AND ces.model = $3
		  AND ces.dimension = $4
		  AND ces.status = 'ready'
		  AND ces.embedding IS NOT NULL
		  AND NOT EXISTS (
		    SELECT 1
		    FROM content_recommendation_blocks crb
		    WHERE crb.active = true
		      AND (
		        crb.content_id = c.id
		        OR (COALESCE(crb.url_key, '') <> '' AND crb.url_key IN (
		          lower(trim(COALESCE(c.url, ''))),
		          lower(trim(COALESCE(c.canonical_url, '')))
		        ))
		      )
		  )
		ORDER BY rank_score DESC, c.created_at DESC
		LIMIT $6
	`, userID, strings.TrimSpace(provider), strings.TrimSpace(model), dimension, *embedding, retrievalLimit)
	if err != nil {
		return nil, fmt.Errorf("search shadow vector lesson candidates: %w", err)
	}
	defer rows.Close()

	candidates := []ContentSearchCandidate{}
	for rows.Next() {
		candidate, err := scanContentSearchCandidate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan shadow vector search candidate: %w", err)
		}
		candidates = append(candidates, candidate)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shadow vector search candidates: %w", err)
	}
	return candidates, nil
}

func (r *Repository) ListActiveRecommendationBlockedURLKeys(ctx context.Context) (map[string]struct{}, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT url_key
		FROM content_recommendation_blocks
		WHERE active = true
		  AND COALESCE(url_key, '') <> ''
	`)
	if err != nil {
		return nil, fmt.Errorf("list active recommendation blocked urls: %w", err)
	}
	defer rows.Close()

	keys := map[string]struct{}{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("scan active recommendation blocked url: %w", err)
		}
		keys[key] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active recommendation blocked urls: %w", err)
	}
	return keys, nil
}

func lexicalFallbackTokens(query string) []string {
	normalized := normalizer.BuildSearchTextKO(query)
	if normalized == "" {
		return nil
	}
	seen := map[string]struct{}{}
	tokens := []string{}
	for _, token := range strings.Fields(normalized) {
		token = strings.TrimSpace(strings.ToLower(token))
		if token == "" || utf8.RuneCountInString(token) < 2 {
			continue
		}
		if isGenericLexicalFallbackToken(token) {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens
}

func lexicalFallbackThreshold(tokens []string) int {
	if len(tokens) >= 2 {
		return 2
	}
	return len(tokens)
}

func isGenericLexicalFallbackToken(token string) bool {
	_, ok := genericLexicalFallbackTokens[token]
	return ok
}

var genericLexicalFallbackTokens = map[string]struct{}{
	"기초": {}, "입문": {}, "초보": {}, "연습": {}, "방법": {}, "튜토리얼": {}, "가이드": {}, "독학": {},
	"활용": {}, "사용": {}, "사용법": {}, "기본": {}, "핵심": {}, "단계": {}, "비교": {}, "확인": {}, "정리": {},
	"만들기": {}, "배우기": {}, "시작": {}, "흐름": {}, "학습": {}, "강의": {}, "영상": {}, "콘텐츠": {},
	"tutorial": {}, "beginner": {}, "basic": {}, "guide": {}, "intro": {}, "practice": {},
}

func scanContentSearchCandidate(row interface{ Scan(dest ...any) error }) (ContentSearchCandidate, error) {
	var candidate ContentSearchCandidate
	if err := row.Scan(
		&candidate.ContentID,
		&candidate.ContentType,
		&candidate.Title,
		&candidate.Description,
		&candidate.ThumbnailURL,
		&candidate.ExternalURL,
		&candidate.Language,
		&candidate.QualityScore,
		&candidate.CreatedAt,
		&candidate.RankScore,
	); err != nil {
		return ContentSearchCandidate{}, err
	}
	candidate.ResourceType = ResourceTypeContent
	candidate.PriceType = PriceTypeFree
	return candidate, nil
}

func mergeLessonCandidates(lexicalCandidates, vectorCandidates []ContentSearchCandidate, limit int, options LessonCandidateSearchOptions) []ContentSearchCandidate {
	if limit <= 0 {
		limit = 5
	}

	type fusedCandidate struct {
		candidate     ContentSearchCandidate
		rrfScore      float64
		finalScore    float64
		bestRawScore  float64
		bestRankIndex int
	}

	const rrfK = 60.0

	fused := map[string]*fusedCandidate{}
	applyRanking := func(candidates []ContentSearchCandidate) {
		for idx, candidate := range candidates {
			key := recommendationCandidateDedupeKey(candidate)
			entry, exists := fused[key]
			if !exists {
				candidateCopy := candidate
				entry = &fusedCandidate{
					candidate:     candidateCopy,
					bestRawScore:  candidate.RankScore,
					bestRankIndex: idx,
				}
				fused[key] = entry
			} else {
				if candidate.RankScore > entry.bestRawScore {
					entry.candidate = candidate
					entry.bestRawScore = candidate.RankScore
				}
				if idx < entry.bestRankIndex {
					entry.bestRankIndex = idx
				}
			}
			entry.rrfScore += 1.0 / (rrfK + float64(idx+1))
		}
	}

	applyRanking(lexicalCandidates)
	applyRanking(vectorCandidates)

	merged := make([]fusedCandidate, 0, len(fused))
	for _, entry := range fused {
		entry.finalScore = entry.rrfScore + metadataBoost(entry.candidate, options)
		entry.candidate.RankScore = entry.finalScore
		merged = append(merged, *entry)
	}

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].finalScore == merged[j].finalScore {
			if merged[i].bestRawScore == merged[j].bestRawScore {
				return merged[i].bestRankIndex < merged[j].bestRankIndex
			}
			return merged[i].bestRawScore > merged[j].bestRawScore
		}
		return merged[i].finalScore > merged[j].finalScore
	})

	if len(merged) > limit {
		merged = merged[:limit]
	}

	result := make([]ContentSearchCandidate, 0, len(merged))
	for _, entry := range merged {
		result = append(result, entry.candidate)
	}
	return result
}

func recommendationCandidateDedupeKey(candidate ContentSearchCandidate) string {
	if candidate.ContentID != nil {
		return "content:" + candidate.ContentID.String()
	}
	if candidate.ExternalURL != nil {
		if urlKey := normalizeRecommendationURLKey(*candidate.ExternalURL); urlKey != "" {
			return "url:" + urlKey
		}
	}
	return "title:" + strings.TrimSpace(strings.ToLower(candidate.Title))
}

func normalizeRecommendationURLKey(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return strings.TrimRight(strings.ToLower(rawURL), "/")
	}
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Fragment = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return strings.TrimRight(parsed.String(), "/")
}

func metadataBoost(candidate ContentSearchCandidate, options LessonCandidateSearchOptions) float64 {
	boost := candidate.QualityScore * 0.08

	if preferredLanguage := strings.TrimSpace(strings.ToLower(options.PreferredLanguage)); preferredLanguage != "" {
		if strings.TrimSpace(strings.ToLower(candidate.Language)) == preferredLanguage {
			boost += 0.03
		}
	}

	if matchesPreferredFormat(candidate.ContentType, options.PreferredFormat) {
		boost += 0.05
	}

	now := options.Now
	if now.IsZero() {
		now = time.Now()
	}
	age := now.Sub(candidate.CreatedAt)
	switch {
	case age <= 30*24*time.Hour:
		boost += 0.03
	case age <= 90*24*time.Hour:
		boost += 0.02
	case age <= 365*24*time.Hour:
		boost += 0.01
	}

	return boost
}

func matchesPreferredFormat(contentType string, preferredFormat *string) bool {
	if preferredFormat == nil {
		return false
	}

	format := strings.TrimSpace(strings.ToLower(*preferredFormat))
	if format == "" {
		return false
	}

	switch {
	case strings.Contains(format, "video"), strings.Contains(format, "영상"), strings.Contains(format, "유튜브"):
		return contentType == "youtube"
	case strings.Contains(format, "text"), strings.Contains(format, "글"), strings.Contains(format, "읽기"), strings.Contains(format, "article"), strings.Contains(format, "blog"):
		return contentType == "blog" || contentType == "article" || contentType == "internal"
	case strings.Contains(format, "internal"), strings.Contains(format, "문서"):
		return contentType == "internal"
	default:
		return false
	}
}
