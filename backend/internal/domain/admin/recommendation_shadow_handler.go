package admin

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

func (h *AdminHandler) GetRecommendationShadowStatus(c *gin.Context) {
	ctx := c.Request.Context()
	endpointConfigured := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")) != ""
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	provider := "embedding_gemma"

	rows, err := h.db.Query(ctx, `
		SELECT status, provider, model, dimension, COUNT(*), MAX(updated_at)
		FROM content_embeddings
		GROUP BY status, provider, model, dimension
		ORDER BY provider, model, dimension, status
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load embedding status"})
		return
	}
	defer rows.Close()

	counts := []gin.H{}
	var totalRows int64
	var readyRows int64
	var failedRows int64
	for rows.Next() {
		var status string
		var rowProvider string
		var rowModel string
		var dimension int
		var count int64
		var updatedAt sql.NullTime
		if err := rows.Scan(&status, &rowProvider, &rowModel, &dimension, &count, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan embedding status"})
			return
		}
		totalRows += count
		if status == "ready" {
			readyRows += count
		}
		if status == "failed" {
			failedRows += count
		}
		item := gin.H{
			"status":    status,
			"provider":  rowProvider,
			"model":     rowModel,
			"dimension": dimension,
			"count":     count,
		}
		if updatedAt.Valid {
			item["updated_at"] = updatedAt.Time
		}
		counts = append(counts, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load embedding status"})
		return
	}

	var targetCount int64
	if err := h.db.QueryRow(ctx, `
		SELECT COUNT(*)
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
	`, provider, model).Scan(&targetCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load embedding target count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoint_configured": endpointConfigured,
		"provider":            provider,
		"model":               model,
		"total_rows":          totalRows,
		"ready_rows":          readyRows,
		"failed_rows":         failedRows,
		"target_count":        targetCount,
		"counts":              counts,
		"cli_command_hint":    `DATABASE_URL="..." go run ./cmd/backfill_content_embeddings -endpoint="http://..." -limit=3`,
	})
}

func (h *AdminHandler) ProbeRecommendationShadowEmbedding(c *gin.Context) {
	var req recommendationDebugShadowProbeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	endpoint := strings.TrimSpace(req.Endpoint)
	endpointSource := "request"
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT"))
		endpointSource = "env"
	}
	if endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EMBEDDING_GEMMA_ENDPOINT is not configured"})
		return
	}
	parsedURL, err := url.Parse(endpoint)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endpoint"})
		return
	}
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	}
	if model == "" {
		model = "embedding-gemma"
	}
	probeText := strings.TrimSpace(req.ProbeText)
	if probeText == "" {
		probeText = "LearnWeaver EmbeddingGemma probe"
	}
	if len([]rune(probeText)) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "probe_text too long"})
		return
	}

	startedAt := time.Now()
	probeCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	vector, err := embedder.NewLocalHTTPClientWithModel(endpoint, model).Embed(probeCtx, probeText)
	latencyMS := int(time.Since(startedAt).Milliseconds())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":           err.Error(),
			"ok":              false,
			"provider":        "embedding_gemma",
			"model":           model,
			"endpoint_source": endpointSource,
			"endpoint_redact": redactEndpointForAdmin(endpoint),
			"latency_ms":      latencyMS,
			"probe_text_size": len([]rune(probeText)),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":              true,
		"provider":        "embedding_gemma",
		"model":           model,
		"dimension":       len(vector),
		"endpoint_source": endpointSource,
		"endpoint_redact": redactEndpointForAdmin(endpoint),
		"latency_ms":      latencyMS,
		"probe_text_size": len([]rune(probeText)),
	})
}

func (h *AdminHandler) ListRecommendationShadowBackfillTargets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	provider := strings.TrimSpace(c.DefaultQuery("provider", "embedding_gemma"))
	if provider == "" {
		provider = "embedding_gemma"
	}
	model := strings.TrimSpace(c.Query("model"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	}
	if model == "" {
		model = "embedding-gemma"
	}

	var totalCount int64
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT COUNT(*)
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
	`, provider, model).Scan(&totalCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count shadow targets"})
		return
	}

	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			c.id::text,
			c.title,
			COALESCE(c.description, ''),
			canonical_url,
			content_type,
			resource_type,
			language,
			quality_score,
			health_score,
			c.updated_at,
			COALESCE(latest.status, ''),
			COALESCE(latest.error_message, ''),
			latest.updated_at
		FROM contents c
		LEFT JOIN LATERAL (
			SELECT status, error_message, updated_at
			FROM content_embeddings ces
			WHERE ces.content_id = c.id
			  AND ces.provider = $1
			  AND ces.model = $2
			ORDER BY ces.updated_at DESC
			LIMIT 1
		) latest ON true
		WHERE c.content_status IN ('active', 'redirect')
		  AND c.health_score > 0.6
		  AND NOT EXISTS (
		    SELECT 1
		    FROM content_embeddings ces_ready
		    WHERE ces_ready.content_id = c.id
		      AND ces_ready.provider = $1
		      AND ces_ready.model = $2
		      AND ces_ready.status = 'ready'
		  )
		ORDER BY c.updated_at DESC, c.created_at DESC
		LIMIT $3
	`, provider, model, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shadow targets"})
		return
	}
	defer rows.Close()

	targets := []gin.H{}
	for rows.Next() {
		var id string
		var title string
		var description string
		var canonicalURL string
		var contentType string
		var resourceType string
		var language string
		var qualityScore float64
		var healthScore float64
		var updatedAt time.Time
		var latestStatus string
		var latestError string
		var latestUpdatedAt sql.NullTime
		if err := rows.Scan(
			&id,
			&title,
			&description,
			&canonicalURL,
			&contentType,
			&resourceType,
			&language,
			&qualityScore,
			&healthScore,
			&updatedAt,
			&latestStatus,
			&latestError,
			&latestUpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan shadow target"})
			return
		}
		item := gin.H{
			"id":                   id,
			"title":                title,
			"description":          description,
			"canonical_url":        canonicalURL,
			"content_type":         contentType,
			"resource_type":        resourceType,
			"language":             language,
			"quality_score":        qualityScore,
			"health_score":         healthScore,
			"updated_at":           updatedAt,
			"latest_shadow_status": latestStatus,
			"latest_error":         latestError,
		}
		if latestUpdatedAt.Valid {
			item["latest_shadow_updated_at"] = latestUpdatedAt.Time
		}
		targets = append(targets, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shadow targets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"provider":    provider,
		"model":       model,
		"limit":       limit,
		"total_count": totalCount,
		"targets":     targets,
	})
}

func redactEndpointForAdmin(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return ""
	}
	if len(endpoint) <= 16 {
		return "***"
	}
	return endpoint[:12] + "..." + endpoint[len(endpoint)-4:]
}
