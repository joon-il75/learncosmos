package admin

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

func (h *AdminHandler) buildRecommendationDebugEmbedding(ctx context.Context, queryText string) (gin.H, *string) {
	endpoint := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT"))
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	if endpoint == "" {
		return gin.H{
			"provider":           "embedding_gemma",
			"model":              model,
			"used":               false,
			"fallback_text_only": true,
			"dimension":          0,
			"error":              "EMBEDDING_GEMMA_ENDPOINT is not configured",
		}, nil
	}

	embedClient, err := embedder.NewProviderClient("embedding_gemma", endpoint)
	if err != nil {
		return gin.H{
			"provider":           "embedding_gemma",
			"model":              model,
			"endpoint":           logsafe.URL(endpoint),
			"used":               false,
			"fallback_text_only": true,
			"dimension":          0,
			"error":              logsafe.PersistedError(err.Error()),
		}, nil
	}

	embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	vector, err := embedClient.Embed(embedCtx, queryText)
	if err != nil {
		return gin.H{
			"provider":           "embedding_gemma",
			"model":              model,
			"endpoint":           logsafe.URL(endpoint),
			"used":               false,
			"fallback_text_only": true,
			"dimension":          0,
			"error":              logsafe.PersistedError(err.Error()),
		}, nil
	}

	embeddingText := recommendationDebugVectorToPGVector(vector)
	return gin.H{
		"provider":           "embedding_gemma",
		"model":              model,
		"endpoint":           logsafe.URL(endpoint),
		"used":               true,
		"fallback_text_only": false,
		"dimension":          len(vector),
		"preview":            recommendationDebugEmbeddingPreview(vector, 6),
	}, &embeddingText
}

func (h *AdminHandler) buildRecommendationDebugShadowEmbedding(ctx context.Context, queryText string) (gin.H, *string) {
	endpoint := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT"))
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	if endpoint == "" {
		return gin.H{
			"provider":    "embedding_gemma",
			"model":       model,
			"attempted":   false,
			"used":        false,
			"dimension":   0,
			"latency_ms":  0,
			"reason":      "EMBEDDING_GEMMA_ENDPOINT is not configured",
			"admin_only":  true,
			"search_used": false,
		}, nil
	}

	client, err := embedder.NewProviderClient("embedding_gemma", endpoint)
	if err != nil {
		return gin.H{
			"provider":    "embedding_gemma",
			"model":       model,
			"attempted":   false,
			"used":        false,
			"dimension":   0,
			"latency_ms":  0,
			"reason":      logsafe.PersistedError(err.Error()),
			"admin_only":  true,
			"search_used": false,
		}, nil
	}

	startedAt := time.Now()
	embedCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	vector, err := client.Embed(embedCtx, queryText)
	latencyMS := int(time.Since(startedAt).Milliseconds())
	if err != nil {
		return gin.H{
			"provider":    "embedding_gemma",
			"model":       model,
			"endpoint":    logsafe.URL(endpoint),
			"attempted":   true,
			"used":        false,
			"dimension":   0,
			"latency_ms":  latencyMS,
			"reason":      logsafe.PersistedError(err.Error()),
			"admin_only":  true,
			"search_used": false,
		}, nil
	}

	embeddingText := recommendationDebugVectorToPGVector(vector)
	return gin.H{
		"provider":    "embedding_gemma",
		"model":       model,
		"endpoint":    logsafe.URL(endpoint),
		"attempted":   true,
		"used":        true,
		"dimension":   len(vector),
		"latency_ms":  latencyMS,
		"preview":     recommendationDebugEmbeddingPreview(vector, 6),
		"reason":      "query_embedding_ready; local content vector index is not connected yet",
		"admin_only":  true,
		"search_used": false,
	}, &embeddingText
}
