package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type patternMatch struct {
	PatternKey      string  `json:"pattern_key"`
	PatternVersion  string  `json:"pattern_version"`
	Domain          string  `json:"domain"`
	GoalType        string  `json:"goal_type"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
	SimilarityScore float64 `json:"similarity_score"`
	Distance        float64 `json:"distance"`
}

type previewOutput struct {
	Query          string                                     `json:"query"`
	Language       string                                     `json:"language"`
	Provider       string                                     `json:"provider"`
	Model          string                                     `json:"model"`
	Dimension      int                                        `json:"dimension"`
	CodeRouting    curriculum.CurriculumPatternRoutingPreview `json:"code_routing"`
	EmbeddingTopN  []patternMatch                             `json:"embedding_top_n"`
	ReadyRowCount  int                                        `json:"ready_row_count"`
	FallbackReason string                                     `json:"fallback_reason,omitempty"`
}

func main() {
	_ = godotenv.Load()

	var (
		goal              = flag.String("goal", "", "goal or source query text to match")
		sourceQuery       = flag.String("source-query", "", "optional source query used for code routing comparison")
		learningGoal      = flag.String("learning-goal", "", "optional confirmed learning goal used for code routing comparison")
		language          = flag.String("language", "ko", "pattern language")
		top               = flag.Int("top", 5, "number of pattern matches to return")
		provider          = flag.String("provider", "embedding_gemma", "pattern embedding provider")
		model             = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
		endpoint          = flag.String("endpoint", defaultString("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma HTTP endpoint")
		expectedDimension = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "expected embedding dimension")
		jsonOutput        = flag.Bool("json", false, "print machine-readable JSON")
	)
	flag.Parse()

	query := strings.TrimSpace(*goal)
	if query == "" {
		query = strings.TrimSpace(strings.Join(compactParts(*sourceQuery, *learningGoal), "\n"))
	}
	if query == "" {
		log.Fatal("missing -goal or -source-query/-learning-goal")
	}
	if *top <= 0 {
		log.Fatal("-top must be greater than 0")
	}
	if *expectedDimension <= 0 {
		log.Fatal("-dimension must be greater than 0")
	}
	if strings.TrimSpace(*endpoint) == "" {
		log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required")
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	client, err := embedder.NewProviderClient(*provider, strings.TrimSpace(*endpoint))
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	vector, err := client.Embed(embedCtx, query)
	cancel()
	if err != nil {
		log.Fatalf("embed query: %v", err)
	}
	if len(vector) != *expectedDimension {
		log.Fatalf("unexpected query embedding dimension: got=%d want=%d", len(vector), *expectedDimension)
	}

	pgVector := curriculum.FormatFloat32VectorForPG(vector)
	readyCount, err := countReadyPatternEmbeddings(ctx, pool, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, normalizeLanguage(*language))
	if err != nil {
		log.Fatalf("count ready pattern embeddings: %v", err)
	}
	matches, err := searchPatternMatches(ctx, pool, pgVector, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, normalizeLanguage(*language), *top)
	if err != nil {
		log.Fatalf("search pattern matches: %v", err)
	}

	routingSource := strings.TrimSpace(*sourceQuery)
	if routingSource == "" {
		routingSource = query
	}
	routingGoal := strings.TrimSpace(*learningGoal)
	if routingGoal == "" {
		routingGoal = query
	}
	output := previewOutput{
		Query:         query,
		Language:      normalizeLanguage(*language),
		Provider:      strings.TrimSpace(*provider),
		Model:         strings.TrimSpace(*model),
		Dimension:     len(vector),
		CodeRouting:   curriculum.AnalyzeCurriculumPatternRoutingPreview(routingSource, routingGoal, normalizeLanguage(*language)),
		EmbeddingTopN: matches,
		ReadyRowCount: readyCount,
	}
	if readyCount == 0 {
		output.FallbackReason = "ready pattern embeddings not found"
	}

	if *jsonOutput {
		payload, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("marshal output: %v", err)
		}
		fmt.Println(string(payload))
		return
	}

	fmt.Printf("query: %s\n", output.Query)
	fmt.Printf("language: %s provider=%s model=%s dimension=%d ready_rows=%d\n", output.Language, output.Provider, output.Model, output.Dimension, output.ReadyRowCount)
	fmt.Printf("code_routing: pattern=%s subpattern=%s mode=%s domain=%s goal_type=%s\n",
		output.CodeRouting.PatternKey,
		emptyAsNone(output.CodeRouting.SubpatternKey),
		emptyAsNone(output.CodeRouting.RefinementMode),
		emptyAsNone(output.CodeRouting.DomainAxis),
		emptyAsNone(output.CodeRouting.GoalType),
	)
	for idx, match := range output.EmbeddingTopN {
		fmt.Printf("%d. score=%.4f distance=%.4f key=%s version=%s title=%s\n", idx+1, match.SimilarityScore, match.Distance, match.PatternKey, match.PatternVersion, match.Title)
	}
	if output.FallbackReason != "" {
		fmt.Printf("fallback_reason: %s\n", output.FallbackReason)
	}
}

func countReadyPatternEmbeddings(ctx context.Context, pool *pgxpool.Pool, provider, model string, dimension int, language string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cp.language = $1
		  AND cpe.provider = $2
		  AND cpe.model = $3
		  AND cpe.dimension = $4
		  AND cpe.status = 'ready'
		  AND cpe.embedding IS NOT NULL
	`, language, provider, model, dimension).Scan(&count)
	return count, err
}

func searchPatternMatches(ctx context.Context, pool *pgxpool.Pool, pgVector, provider, model string, dimension int, language string, limit int) ([]patternMatch, error) {
	rows, err := pool.Query(ctx, `
		SELECT
		  cp.pattern_key,
		  cp.pattern_version,
		  cp.domain,
		  cp.goal_type,
		  cp.language,
		  cp.title,
		  cp.summary,
		  cpe.embedding <=> $1::vector AS distance,
		  1 - (cpe.embedding <=> $1::vector) AS similarity_score
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cp.language = $2
		  AND cpe.provider = $3
		  AND cpe.model = $4
		  AND cpe.dimension = $5
		  AND cpe.status = 'ready'
		  AND cpe.embedding IS NOT NULL
		ORDER BY cpe.embedding <=> $1::vector ASC
		LIMIT $6
	`, pgVector, language, provider, model, dimension, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []patternMatch{}
	for rows.Next() {
		var item patternMatch
		if err := rows.Scan(
			&item.PatternKey,
			&item.PatternVersion,
			&item.Domain,
			&item.GoalType,
			&item.Language,
			&item.Title,
			&item.Summary,
			&item.Distance,
			&item.SimilarityScore,
		); err != nil {
			return nil, err
		}
		matches = append(matches, item)
	}
	return matches, rows.Err()
}

func compactParts(parts ...string) []string {
	result := []string{}
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeLanguage(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func emptyAsNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}

func defaultString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func defaultInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
