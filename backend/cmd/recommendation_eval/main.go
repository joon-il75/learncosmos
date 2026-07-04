package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/evaluation"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	loadEnv()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	endpoint := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT"))
	if endpoint == "" {
		log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required")
	}
	embedClient, err := embedder.NewProviderClient("embedding_gemma", endpoint)
	if err != nil {
		log.Fatalf("failed to build embedding_gemma client: %v", err)
	}

	repo := curriculum.NewRepository(pool)
	evalUserID, err := resolveEvaluationUserID(ctx)
	if err != nil {
		log.Fatalf("failed to resolve evaluation user: %v", err)
	}
	if evalUserID != uuid.Nil {
		log.Printf("using recommendation evaluation user_id=%s", evalUserID.String())
	} else {
		log.Printf("using recommendation evaluation scope=all_indexed_content")
	}
	summary, err := evaluation.RunRecommendationEvaluationForUser(ctx, repo, embedClient, evalUserID, nil)
	if err != nil {
		log.Fatalf("failed to run recommendation evaluation: %v", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		log.Fatalf("failed to print evaluation result: %v", err)
	}
}

func loadEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}

func resolveEvaluationUserID(ctx context.Context) (uuid.UUID, error) {
	_ = ctx
	if raw := strings.TrimSpace(os.Getenv("RECOMMENDATION_EVAL_USER_ID")); raw != "" {
		return uuid.Parse(raw)
	}
	return uuid.Nil, nil
}
