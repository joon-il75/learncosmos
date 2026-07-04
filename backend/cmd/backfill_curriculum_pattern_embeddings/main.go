package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type patternEmbeddingTarget struct {
	EmbeddingID       uuid.UUID
	PatternID         uuid.UUID
	PatternKey        string
	PatternVersion    string
	Title             string
	EmbeddingText     string
	EmbeddingTextHash string
}

func main() {
	_ = godotenv.Load()

	var (
		limit             = flag.Int("limit", 20, "maximum pattern embeddings to backfill")
		provider          = flag.String("provider", "embedding_gemma", "pattern embedding provider")
		model             = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
		endpoint          = flag.String("endpoint", defaultString("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma HTTP endpoint")
		expectedDimension = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "expected embedding dimension")
		dryRun            = flag.Bool("dry-run", false, "list target patterns without embedding")
		probe             = flag.Bool("probe", false, "call the embedding endpoint once without touching the database")
		probeText         = flag.String("probe-text", "LearnWeaver curriculum pattern embedding probe", "text used for -probe")
		retryFailed       = flag.Bool("retry-failed", false, "include failed rows in target list")
	)
	flag.Parse()

	ctx := context.Background()
	if *probe {
		if strings.TrimSpace(*endpoint) == "" {
			log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required for -probe")
		}
		client, err := embedder.NewProviderClient(*provider, strings.TrimSpace(*endpoint))
		if err != nil {
			log.Fatalf("failed to create embedder: %v", err)
		}
		probeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		vector, err := client.Embed(probeCtx, strings.TrimSpace(*probeText))
		cancel()
		if err != nil {
			log.Fatalf("probe failed: %v", err)
		}
		if len(vector) == 0 {
			log.Fatal("probe failed: empty embedding vector")
		}
		log.Printf("probe ok provider=%s model=%s dimension=%d endpoint=%s", *provider, *model, len(vector), redactEndpoint(*endpoint))
		return
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if strings.TrimSpace(*endpoint) == "" && !*dryRun {
		log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required unless -dry-run is set")
	}
	if *limit <= 0 {
		log.Fatal("-limit must be greater than 0")
	}
	if *expectedDimension <= 0 {
		log.Fatal("-dimension must be greater than 0")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	targets, err := listPatternEmbeddingTargets(ctx, pool, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, *limit, *retryFailed)
	if err != nil {
		log.Fatalf("failed to list pattern embedding targets: %v", err)
	}
	log.Printf("curriculum pattern embedding targets=%d provider=%s model=%s dimension=%d dry_run=%v retry_failed=%v", len(targets), *provider, *model, *expectedDimension, *dryRun, *retryFailed)
	if *dryRun {
		for _, item := range targets {
			log.Printf("[target] embedding_id=%s pattern_key=%s version=%s hash=%s title=%s", item.EmbeddingID, item.PatternKey, item.PatternVersion, item.EmbeddingTextHash, item.Title)
		}
		return
	}

	client, err := embedder.NewProviderClient(*provider, strings.TrimSpace(*endpoint))
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	successCount := 0
	failCount := 0
	for _, item := range targets {
		embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		vector, err := client.Embed(embedCtx, item.EmbeddingText)
		cancel()
		if err != nil {
			failCount++
			if saveErr := markPatternEmbeddingFailure(ctx, pool, item.EmbeddingID, err); saveErr != nil {
				log.Printf("[failed][save_error] embedding_id=%s pattern_key=%s err=%v save_err=%v", item.EmbeddingID, item.PatternKey, err, saveErr)
				continue
			}
			log.Printf("[failed] embedding_id=%s pattern_key=%s err=%v", item.EmbeddingID, item.PatternKey, err)
			continue
		}
		if len(vector) != *expectedDimension {
			failCount++
			err := fmt.Errorf("unexpected embedding dimension: got=%d want=%d", len(vector), *expectedDimension)
			if saveErr := markPatternEmbeddingFailure(ctx, pool, item.EmbeddingID, err); saveErr != nil {
				log.Printf("[dimension_failed][save_error] embedding_id=%s pattern_key=%s err=%v save_err=%v", item.EmbeddingID, item.PatternKey, err, saveErr)
				continue
			}
			log.Printf("[dimension_failed] embedding_id=%s pattern_key=%s err=%v", item.EmbeddingID, item.PatternKey, err)
			continue
		}

		embeddingText := vectorToPGVector(vector)
		if err := markPatternEmbeddingReady(ctx, pool, item.EmbeddingID, embeddingText); err != nil {
			failCount++
			log.Printf("[db_failed] embedding_id=%s pattern_key=%s err=%v", item.EmbeddingID, item.PatternKey, err)
			continue
		}
		successCount++
		log.Printf("[ready] embedding_id=%s pattern_key=%s dim=%d", item.EmbeddingID, item.PatternKey, len(vector))
	}

	log.Printf("curriculum pattern embedding backfill completed ready=%d failed=%d", successCount, failCount)
}

func listPatternEmbeddingTargets(ctx context.Context, pool *pgxpool.Pool, provider, model string, dimension, limit int, retryFailed bool) ([]patternEmbeddingTarget, error) {
	statuses := []string{"pending"}
	if retryFailed {
		statuses = append(statuses, "failed")
	}

	rows, err := pool.Query(ctx, `
		SELECT cpe.id, cp.id, cp.pattern_key, cp.pattern_version, cp.title, cp.embedding_text, cpe.embedding_text_hash
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cpe.provider = $1
		  AND cpe.model = $2
		  AND cpe.dimension = $3
		  AND cpe.status = ANY($4)
		ORDER BY cp.pattern_key ASC, cpe.updated_at ASC
		LIMIT $5
	`, provider, model, dimension, statuses, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []patternEmbeddingTarget{}
	for rows.Next() {
		var item patternEmbeddingTarget
		if err := rows.Scan(&item.EmbeddingID, &item.PatternID, &item.PatternKey, &item.PatternVersion, &item.Title, &item.EmbeddingText, &item.EmbeddingTextHash); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func markPatternEmbeddingReady(ctx context.Context, pool *pgxpool.Pool, embeddingID uuid.UUID, embedding string) error {
	_, err := pool.Exec(ctx, `
		UPDATE curriculum_pattern_embeddings
		SET embedding = $2::vector,
		    status = 'ready',
		    error_message = '',
		    embedded_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`, embeddingID, embedding)
	return err
}

func markPatternEmbeddingFailure(ctx context.Context, pool *pgxpool.Pool, embeddingID uuid.UUID, embedErr error) error {
	_, err := pool.Exec(ctx, `
		UPDATE curriculum_pattern_embeddings
		SET status = 'failed',
		    error_message = $2,
		    updated_at = NOW()
		WHERE id = $1
	`, embeddingID, truncateString(embedErr.Error(), 1000))
	return err
}

func vectorToPGVector(vector []float32) string {
	parts := make([]string, 0, len(vector))
	for _, value := range vector {
		parts = append(parts, fmt.Sprintf("%f", value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func truncateString(value string, maxRunes int) string {
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	return string(runes[:maxRunes])
}

func redactEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return ""
	}
	if len(endpoint) <= 16 {
		return "***"
	}
	return endpoint[:12] + "..." + endpoint[len(endpoint)-4:]
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
