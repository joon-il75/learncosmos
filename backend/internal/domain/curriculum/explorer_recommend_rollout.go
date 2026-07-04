package curriculum

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) buildRecommendationRolloutReadOnlyDebug(ctx context.Context, userID uuid.UUID) gin.H {
	var rolloutID string
	var mode string
	var trafficPercent int
	var embeddingProvider string
	var embeddingModel string
	var featureSchemaVersion string
	var rankerModelVersion string
	if err := h.repo.pool.QueryRow(ctx, `
		SELECT
			id::text,
			mode,
			traffic_percent,
			embedding_provider,
			embedding_model,
			feature_schema_version,
			ranker_model_version
		FROM recommendation_rollout_states
		WHERE active = true
		ORDER BY updated_at DESC
		LIMIT 1
	`).Scan(&rolloutID, &mode, &trafficPercent, &embeddingProvider, &embeddingModel, &featureSchemaVersion, &rankerModelVersion); err != nil {
		return gin.H{
			"instrumentation_only": true,
			"available":            false,
			"active_applied":       false,
			"reason":               "active rollout state not found",
		}
	}

	bucket := recommendationRolloutReadOnlyBucket(userID.String(), rolloutID)
	wouldApply := strings.HasPrefix(mode, "learner_") && bucket < trafficPercent
	reason := "outside rollout traffic"
	fallbackReason := "not in active learner cohort"
	routeGuard := h.evaluateRecommendationLearnerRankerRouteGuard(ctx, embeddingProvider, embeddingModel, rankerModelVersion, featureSchemaVersion)
	activeRoute := "baseline"
	if wouldApply {
		if ready, _ := routeGuard["ready"].(bool); ready {
			activeRoute = "ranker_shadow_guarded"
			reason = "inside rollout traffic; ranker shadow route is guard-ready"
			fallbackReason = "candidate order mutation is disabled; baseline recommendation used"
		} else {
			reason = "inside rollout traffic; active ranking shell is baseline fallback"
			fallbackReason = strings.TrimSpace(rolloutDebugString(routeGuard, "reason"))
			if fallbackReason == "" {
				fallbackReason = "active ranking guard failed; baseline recommendation used"
			}
		}
	}
	if mode == "shadow_only" || mode == "admin_preview" {
		reason = "mode is admin/shadow only"
		fallbackReason = "rollout mode is not learner active"
		activeRoute = "baseline"
	}
	if mode == "rolled_back" {
		reason = "rollout is rolled back"
		fallbackReason = "rollout is rolled back"
		activeRoute = "baseline"
	}
	return gin.H{
		"instrumentation_only":    true,
		"available":               true,
		"routing_shell":           true,
		"guarded_route":           true,
		"active_applied":          false,
		"candidate_order_mutated": false,
		"active_route":            activeRoute,
		"fallback_reason":         fallbackReason,
		"would_apply":             wouldApply,
		"reason":                  reason,
		"rollout_id":              rolloutID,
		"mode":                    mode,
		"traffic_percent":         trafficPercent,
		"bucket":                  bucket,
		"embedding_provider":      embeddingProvider,
		"embedding_model":         embeddingModel,
		"feature_schema_version":  featureSchemaVersion,
		"ranker_model_version":    rankerModelVersion,
		"ranker_shadow_guard":     routeGuard,
	}
}

func (h *Handler) evaluateRecommendationLearnerRankerRouteGuard(ctx context.Context, embeddingProvider, embeddingModel, rankerModelVersion, featureSchemaVersion string) gin.H {
	embeddingProvider = strings.TrimSpace(embeddingProvider)
	if embeddingProvider == "" {
		embeddingProvider = "embedding_gemma"
	}
	embeddingModel = strings.TrimSpace(embeddingModel)
	if embeddingModel == "" {
		embeddingModel = strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	}
	if embeddingModel == "" {
		embeddingModel = "embedding-gemma"
	}
	rankerModelVersion = strings.TrimSpace(rankerModelVersion)
	featureSchemaVersion = strings.TrimSpace(featureSchemaVersion)
	if featureSchemaVersion == "" {
		featureSchemaVersion = "ranker-feature-v1"
	}

	checks := gin.H{
		"embedding_provider":       embeddingProvider,
		"embedding_model":          embeddingModel,
		"feature_schema_version":   featureSchemaVersion,
		"ranker_model_version":     rankerModelVersion,
		"embedding_endpoint_ready": strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")) != "",
	}
	if !checks["embedding_endpoint_ready"].(bool) {
		return gin.H{"ready": false, "reason": "EmbeddingGemma endpoint is not configured", "checks": checks}
	}

	var readyShadowRows int64
	if err := h.repo.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM content_embeddings
		WHERE provider = $1
		  AND model = $2
		  AND status = 'ready'
	`, embeddingProvider, embeddingModel).Scan(&readyShadowRows); err != nil {
		checks["shadow_vector_ready_rows"] = int64(0)
		return gin.H{"ready": false, "reason": "failed to count shadow vector rows", "checks": checks}
	}
	checks["shadow_vector_ready_rows"] = readyShadowRows
	if readyShadowRows == 0 {
		return gin.H{"ready": false, "reason": "ready shadow vector rows not found", "checks": checks}
	}

	checks["ranker_model_configured"] = rankerModelVersion != ""
	if rankerModelVersion == "" {
		return gin.H{"ready": false, "reason": "ranker model version is not configured", "checks": checks}
	}
	modelPath := recommendationLearnerRankerModelPath(rankerModelVersion)
	checks["ranker_model_path"] = modelPath
	modelInfo, err := os.Stat(modelPath)
	if err != nil || modelInfo.IsDir() || modelInfo.Size() == 0 {
		return gin.H{"ready": false, "reason": "ranker model artifact is not ready", "checks": checks}
	}
	checks["ranker_model_size_bytes"] = modelInfo.Size()

	return gin.H{
		"ready":  true,
		"reason": "ranker shadow route guard passed; candidate order mutation is still disabled",
		"checks": checks,
	}
}

func recommendationLearnerRankerModelPath(rankerModelVersion string) string {
	modelPath := strings.TrimSpace(os.Getenv("RANKER_MODEL_PATH"))
	if modelPath != "" {
		return modelPath
	}
	modelDir := strings.TrimSpace(os.Getenv("RANKER_MODEL_DIR"))
	if modelDir == "" {
		modelDir = filepath.Join("models", "ranker")
	}
	for _, ext := range []string{".txt", ".lgb", ".model"} {
		candidatePath := filepath.Join(modelDir, strings.TrimSpace(rankerModelVersion)+ext)
		if _, err := os.Stat(candidatePath); err == nil {
			return candidatePath
		}
	}
	return filepath.Join(modelDir, strings.TrimSpace(rankerModelVersion)+".txt")
}

func recommendationRolloutReadOnlyBucket(userID, rolloutID string) int {
	sum := sha256.Sum256([]byte(strings.TrimSpace(userID) + ":" + strings.TrimSpace(rolloutID)))
	return int(binary.BigEndian.Uint32(sum[:4]) % 100)
}

func (h *Handler) recordRecommendationRolloutExposure(ctx context.Context, userID uuid.UUID, req RecommendExplorerContentRequest, items []RecommendExplorerContentItem, rolloutDebug gin.H) {
	requestSnapshot, _ := json.Marshal(gin.H{
		"query":                  strings.TrimSpace(req.Query),
		"query_language":         normalizeLearningLanguage(req.QueryLanguage),
		"max_results":            req.MaxResults,
		"course_draft_id":        strings.TrimSpace(req.CourseDraftID),
		"point_id":               strings.TrimSpace(req.PointID),
		"region_title":           strings.TrimSpace(req.RegionTitle),
		"subregion_title":        strings.TrimSpace(req.SubRegionTitle),
		"node_title":             strings.TrimSpace(req.NodeTitle),
		"node_summary_present":   strings.TrimSpace(req.NodeSummary) != "",
		"region_context_present": strings.TrimSpace(req.RegionDescription) != "" || strings.TrimSpace(req.SubRegionDescription) != "",
	})
	candidatesSnapshot, _ := json.Marshal(items)
	if err := h.repo.RecordRecommendationRolloutLearnerExposure(ctx, userID, req.CourseDraftID, req.PointID, rolloutDebug, len(items), string(requestSnapshot), string(candidatesSnapshot)); err != nil {
		log.Printf("[recommendation/rollout] exposure telemetry skipped: %v", err)
	}
}
