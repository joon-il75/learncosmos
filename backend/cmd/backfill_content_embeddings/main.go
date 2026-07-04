package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type contentRow struct {
	ID          uuid.UUID
	Title       string
	Description string
}

func main() {
	var (
		limit             = flag.Int("limit", 20, "maximum contents to backfill")
		provider          = flag.String("provider", "embedding_gemma", "content embedding provider")
		model             = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "content embedding model")
		endpoint          = flag.String("endpoint", defaultString("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma HTTP endpoint")
		expectedDimension = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "expected dimension used for failed rows")
		dryRun            = flag.Bool("dry-run", false, "list target contents without embedding")
		probe             = flag.Bool("probe", false, "call the embedding endpoint once without touching the database")
		probeText         = flag.String("probe-text", "LearnWeaver EmbeddingGemma probe", "text used for -probe")
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

	items, err := listContentEmbeddingTargets(ctx, pool, strings.TrimSpace(*provider), strings.TrimSpace(*model), *limit)
	if err != nil {
		log.Fatalf("failed to list targets: %v", err)
	}
	log.Printf("content embedding targets: %d provider=%s model=%s dry_run=%v", len(items), *provider, *model, *dryRun)
	if *dryRun {
		for _, item := range items {
			log.Printf("[target] id=%s title=%s", item.ID, item.Title)
		}
		return
	}

	client, err := embedder.NewProviderClient(*provider, strings.TrimSpace(*endpoint))
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	successCount := 0
	failCount := 0
	for _, item := range items {
		text := embedder.BuildText(item.Title, item.Description)
		sourceHash := sourceTextHash(text)
		preview := sourceTextPreview(text, 240)

		embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		vector, err := client.Embed(embedCtx, text)
		cancel()
		if err != nil {
			failCount++
			if saveErr := upsertContentEmbeddingFailure(ctx, pool, item, *provider, *model, *expectedDimension, sourceHash, preview, err); saveErr != nil {
				log.Printf("[failed][save_error] id=%s err=%v save_err=%v", item.ID, err, saveErr)
				continue
			}
			log.Printf("[failed] id=%s err=%v", item.ID, err)
			continue
		}

		embeddingText := vectorToPGVector(vector)
		if err := upsertContentEmbeddingReady(ctx, pool, item, *provider, *model, len(vector), sourceHash, preview, embeddingText); err != nil {
			failCount++
			log.Printf("[db_failed] id=%s err=%v", item.ID, err)
			continue
		}
		successCount++
		log.Printf("[ready] id=%s title=%s dim=%d", item.ID, item.Title, len(vector))
	}

	log.Printf("content embedding backfill completed ready=%d failed=%d", successCount, failCount)
}

func listContentEmbeddingTargets(ctx context.Context, pool *pgxpool.Pool, provider, model string, limit int) ([]contentRow, error) {
	rows, err := pool.Query(ctx, `
		SELECT c.id, c.title, COALESCE(c.description, '')
		FROM contents c
		WHERE c.content_status IN ('active', 'redirect')
		  AND c.health_score > 0.6
		  AND NOT EXISTS (
		    SELECT 1
		    FROM content_embeddings ces
		    WHERE ces.content_id = c.id
		      AND ces.provider = $1
		      AND ces.model = $2
		      AND ces.status = 'ready'
		  )
		ORDER BY c.updated_at DESC, c.created_at DESC
		LIMIT $3
	`, provider, model, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []contentRow{}
	for rows.Next() {
		var item contentRow
		if err := rows.Scan(&item.ID, &item.Title, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func upsertContentEmbeddingReady(ctx context.Context, pool *pgxpool.Pool, item contentRow, provider, model string, dimension int, sourceHash, preview, embedding string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO content_embeddings (
			content_id, provider, model, dimension, embedding,
			source_text_hash, source_text_preview, status, error_message, embedded_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5::vector, $6, $7, 'ready', '', NOW(), NOW(), NOW())
		ON CONFLICT (content_id, provider, model, dimension)
		DO UPDATE SET
			embedding = EXCLUDED.embedding,
			source_text_hash = EXCLUDED.source_text_hash,
			source_text_preview = EXCLUDED.source_text_preview,
			status = 'ready',
			error_message = '',
			embedded_at = NOW(),
			updated_at = NOW()
	`, item.ID, provider, model, dimension, embedding, sourceHash, preview)
	return err
}

func upsertContentEmbeddingFailure(ctx context.Context, pool *pgxpool.Pool, item contentRow, provider, model string, dimension int, sourceHash, preview string, embedErr error) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO content_embeddings (
			content_id, provider, model, dimension,
			source_text_hash, source_text_preview, status, error_message,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'failed', $7, NOW(), NOW())
		ON CONFLICT (content_id, provider, model, dimension)
		DO UPDATE SET
			source_text_hash = EXCLUDED.source_text_hash,
			source_text_preview = EXCLUDED.source_text_preview,
			status = 'failed',
			error_message = EXCLUDED.error_message,
			updated_at = NOW()
	`, item.ID, provider, model, dimension, sourceHash, preview, truncateString(embedErr.Error(), 1000))
	return err
}

func vectorToPGVector(vector []float32) string {
	parts := make([]string, 0, len(vector))
	for _, value := range vector {
		parts = append(parts, fmt.Sprintf("%f", value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func sourceTextHash(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}

func sourceTextPreview(text string, maxRunes int) string {
	return truncateString(strings.TrimSpace(text), maxRunes)
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
