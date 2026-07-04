package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	"github.com/learnweaver/backend/internal/pkg/ranker"
)

func buildRecommendationDebugShadowCandidateSummaries(candidates []curriculum.ContentSearchCandidate) []gin.H {
	items := make([]gin.H, 0, len(candidates))
	for idx, candidate := range candidates {
		items = append(items, gin.H{
			"content_id":       candidate.ContentID,
			"title":            candidate.Title,
			"description":      candidate.Description,
			"thumbnail_url":    candidate.ThumbnailURL,
			"external_url":     candidate.ExternalURL,
			"rank_score":       candidate.RankScore,
			"content_type":     candidate.ContentType,
			"resource_type":    candidate.ResourceType,
			"quality_score":    candidate.QualityScore,
			"language":         candidate.Language,
			"selection_reason": "EmbeddingGemma primary vector similarity",
			"feature_snapshot": recommendationDebugRankerFeatureSnapshot(candidate, "shadow_vector", false, true, idx+1, "ko", time.Now()),
		})
	}
	return items
}

func applyRecommendationDebugRankerScores(ctx context.Context, candidates []gin.H, scorer ranker.Scorer) gin.H {
	result := gin.H{
		"enabled":                 true,
		"provider":                scorer.Provider(),
		"model_version":           scorer.ModelVersion(),
		"feature_schema_version":  scorer.FeatureSchemaVersion(),
		"candidate_count":         len(candidates),
		"reason":                  "no-op ranker preserves route order",
		"rerank_applied":          true,
		"candidate_order_mutated": false,
		"scores":                  []gin.H{},
	}
	if len(candidates) == 0 {
		result["reason"] = "no candidates to score"
		return result
	}

	features := make([]ranker.CandidateFeature, 0, len(candidates))
	for idx, candidate := range candidates {
		feature, ok := candidate["feature_snapshot"].(gin.H)
		if !ok || feature == nil {
			feature = gin.H{
				"schema_version": scorer.FeatureSchemaVersion(),
				"route_rank":     idx + 1,
			}
		}
		features = append(features, ranker.CandidateFeature(feature))
	}

	scores, err := scorer.Score(ctx, features)
	if err != nil {
		result["enabled"] = false
		result["reason"] = logsafe.PersistedError(err.Error())
		return result
	}

	validScores := make([]ranker.Score, 0, len(scores))
	for _, score := range scores {
		if score.Index < 0 || score.Index >= len(candidates) {
			continue
		}
		validScores = append(validScores, score)
	}
	sort.SliceStable(validScores, func(i, j int) bool {
		if validScores[i].Score == validScores[j].Score {
			return validScores[i].Index < validScores[j].Index
		}
		return validScores[i].Score > validScores[j].Score
	})

	scoreRows := make([]gin.H, 0, len(validScores))
	for rerankRank, score := range validScores {
		candidate := candidates[score.Index]
		routeRank := recommendationDebugRouteRank(candidate, score.Index+1)
		candidate["ranker_score"] = score.Score
		candidate["ranker_route_rank"] = routeRank
		candidate["ranker_rerank_rank"] = rerankRank + 1
		candidate["ranker_rank_delta"] = routeRank - (rerankRank + 1)
		candidate["ranker_rank"] = rerankRank + 1
		candidate["ranker_reason"] = score.Reason
		candidate["ranker_provider"] = scorer.Provider()
		candidate["ranker_model_version"] = scorer.ModelVersion()
		scoreRows = append(scoreRows, gin.H{
			"index":       score.Index,
			"route_rank":  routeRank,
			"rerank_rank": rerankRank + 1,
			"rank_delta":  routeRank - (rerankRank + 1),
			"score":       score.Score,
			"reason":      score.Reason,
		})
	}
	result["scores"] = scoreRows
	return result
}

func recommendationDebugRouteRank(candidate gin.H, fallback int) int {
	feature, ok := candidate["feature_snapshot"].(gin.H)
	if !ok || feature == nil {
		return fallback
	}
	switch typed := feature["route_rank"].(type) {
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	default:
		return fallback
	}
}

func recommendationDebugRankerFeatureSnapshot(candidate curriculum.ContentSearchCandidate, route string, lexicalMatched, vectorMatched bool, routeRank int, preferredLanguage string, now time.Time) gin.H {
	languageMatch := strings.TrimSpace(preferredLanguage) != "" && strings.EqualFold(strings.TrimSpace(candidate.Language), strings.TrimSpace(preferredLanguage))
	ageDays := 0.0
	if !candidate.CreatedAt.IsZero() {
		ageDays = now.Sub(candidate.CreatedAt).Hours() / 24
		if ageDays < 0 {
			ageDays = 0
		}
	}
	return gin.H{
		"schema_version":  "ranker-feature-v1",
		"route":           route,
		"route_rank":      routeRank,
		"rank_score":      candidate.RankScore,
		"quality_score":   candidate.QualityScore,
		"age_days":        ageDays,
		"lexical_matched": lexicalMatched,
		"vector_matched":  vectorMatched,
		"language":        candidate.Language,
		"language_match":  languageMatch,
		"content_type":    candidate.ContentType,
		"resource_type":   candidate.ResourceType,
		"has_thumbnail":   candidate.ThumbnailURL != nil && strings.TrimSpace(*candidate.ThumbnailURL) != "",
		"has_url":         candidate.ExternalURL != nil && strings.TrimSpace(*candidate.ExternalURL) != "",
	}
}

func recommendationDebugRankerSnapshotReason(modelVersion string) string {
	if strings.TrimSpace(modelVersion) == "" {
		return "LightGBM Ranker model version is not configured"
	}
	return "LightGBM Ranker model version is configured; ranker execution is not connected yet"
}

func (h *AdminHandler) buildRecommendationRankerArtifactStatus(ctx context.Context, state *recommendationRolloutStateRow) recommendationRankerArtifactStatus {
	status := recommendationRankerArtifactStatus{
		Ready:                false,
		Provider:             "lightgbm",
		ConfigSource:         "none",
		FeatureSchemaVersion: "ranker-feature-v1",
		CheckedPaths:         []string{},
	}
	defer sanitizeRecommendationRankerArtifactStatus(&status)
	if state == nil {
		status.Reason = "active rollout state not found"
		return status
	}
	status.RankerModelVersion = strings.TrimSpace(state.RankerModelVersion)
	status.FeatureSchemaVersion = strings.TrimSpace(state.FeatureSchemaVersion)
	if status.FeatureSchemaVersion == "" {
		status.FeatureSchemaVersion = "ranker-feature-v1"
	}
	status.ExpectedModelVersion = status.RankerModelVersion
	status.ExpectedSchemaVersion = status.FeatureSchemaVersion
	if status.RankerModelVersion == "" {
		status.Reason = "ranker model version is not configured"
		return status
	}

	modelPath := strings.TrimSpace(os.Getenv("RANKER_MODEL_PATH"))
	if modelPath != "" {
		status.ConfigSource = "RANKER_MODEL_PATH"
		status.CheckedPaths = append(status.CheckedPaths, modelPath)
	} else {
		modelDir := strings.TrimSpace(os.Getenv("RANKER_MODEL_DIR"))
		if modelDir == "" {
			modelDir = filepath.Join("models", "ranker")
		}
		status.ConfigSource = "RANKER_MODEL_DIR"
		for _, ext := range []string{".txt", ".lgb", ".model"} {
			candidatePath := filepath.Join(modelDir, status.RankerModelVersion+ext)
			status.CheckedPaths = append(status.CheckedPaths, candidatePath)
			if _, err := os.Stat(candidatePath); err == nil {
				modelPath = candidatePath
				break
			}
		}
		if modelPath == "" && len(status.CheckedPaths) > 0 {
			modelPath = status.CheckedPaths[0]
		}
	}
	status.ModelPath = modelPath
	if status.ModelPath == "" {
		status.Reason = "ranker model path is not configured"
		return status
	}

	modelInfo, err := os.Stat(status.ModelPath)
	if err != nil {
		if os.IsNotExist(err) {
			status.Reason = "ranker model file does not exist"
		} else {
			status.Reason = "ranker model file is not readable: " + logsafe.PersistedError(err.Error())
		}
		return status
	}
	if modelInfo.IsDir() {
		status.Reason = "ranker model path is a directory"
		return status
	}
	status.ModelFileExists = true
	status.ModelSizeBytes = modelInfo.Size()
	status.ModelUpdatedAt = modelInfo.ModTime().Format(time.RFC3339)
	if status.ModelSizeBytes == 0 {
		status.Reason = "ranker model file is empty"
		return status
	}

	status.ManifestPath = strings.TrimSpace(os.Getenv("RANKER_MODEL_MANIFEST_PATH"))
	if status.ManifestPath == "" {
		status.ManifestPath = status.ModelPath + ".manifest.json"
	}
	manifestInfo, err := os.Stat(status.ManifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			status.Reason = "ranker model file exists; manifest is optional and not found"
			status.Ready = true
			return status
		}
		status.Reason = "ranker manifest is not readable: " + logsafe.PersistedError(err.Error())
		return status
	}
	if manifestInfo.IsDir() {
		status.Reason = "ranker manifest path is a directory"
		return status
	}
	status.ManifestExists = true
	manifestBytes, err := os.ReadFile(status.ManifestPath)
	if err != nil {
		status.Reason = "failed to read ranker manifest: " + logsafe.PersistedError(err.Error())
		return status
	}
	var manifest struct {
		ModelVersion         string `json:"model_version"`
		RankerModelVersion   string `json:"ranker_model_version"`
		FeatureSchemaVersion string `json:"feature_schema_version"`
		Provider             string `json:"provider"`
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		status.Reason = "ranker manifest is invalid JSON: " + logsafe.PersistedError(err.Error())
		return status
	}
	manifestModelVersion := strings.TrimSpace(firstNonEmpty(manifest.RankerModelVersion, manifest.ModelVersion))
	manifestSchemaVersion := strings.TrimSpace(manifest.FeatureSchemaVersion)
	if strings.TrimSpace(manifest.Provider) != "" {
		status.Provider = strings.TrimSpace(manifest.Provider)
	}
	status.ManifestVersionMatch = manifestModelVersion == "" || manifestModelVersion == status.RankerModelVersion
	status.ManifestSchemaMatch = manifestSchemaVersion == "" || manifestSchemaVersion == status.FeatureSchemaVersion
	if !status.ManifestVersionMatch {
		status.Reason = "ranker manifest model version does not match rollout state"
		return status
	}
	if !status.ManifestSchemaMatch {
		status.Reason = "ranker manifest feature schema version does not match rollout state"
		return status
	}
	status.Ready = true
	status.Reason = "ranker model artifact is readable and matches rollout state"
	_ = ctx
	return status
}

func recommendationDebugTopKOverlap(baseline []gin.H, shadow []gin.H, k int) *float64 {
	if k <= 0 {
		k = 5
	}
	if len(baseline) == 0 || len(shadow) == 0 {
		return nil
	}
	baselineLimit := minInt(k, len(baseline))
	shadowLimit := minInt(k, len(shadow))
	keys := map[string]struct{}{}
	for idx := 0; idx < baselineLimit; idx++ {
		if key := recommendationDebugSummaryKey(baseline[idx]); key != "" {
			keys[key] = struct{}{}
		}
	}
	matches := 0
	for idx := 0; idx < shadowLimit; idx++ {
		if _, ok := keys[recommendationDebugSummaryKey(shadow[idx])]; ok {
			matches++
		}
	}
	denominator := baselineLimit
	if shadowLimit < denominator {
		denominator = shadowLimit
	}
	if denominator <= 0 {
		return nil
	}
	value := float64(matches) / float64(denominator)
	return &value
}

func recommendationDebugSummaryKey(item gin.H) string {
	if value, ok := item["content_id"]; ok && value != nil {
		return "content:" + strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := item["external_url"]; ok && value != nil && strings.TrimSpace(fmt.Sprint(value)) != "" {
		return "url:" + strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := item["title"]; ok && value != nil {
		return "title:" + strings.TrimSpace(strings.ToLower(fmt.Sprint(value)))
	}
	return ""
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sanitizeRecommendationRankerArtifactStatus(status *recommendationRankerArtifactStatus) {
	if status == nil {
		return
	}
	status.ModelPath = displayPathLabel(status.ModelPath)
	status.ManifestPath = displayPathLabel(status.ManifestPath)
	if len(status.CheckedPaths) > 0 {
		labels := make([]string, 0, len(status.CheckedPaths))
		for _, checkedPath := range status.CheckedPaths {
			labels = append(labels, displayPathLabel(checkedPath))
		}
		status.CheckedPaths = labels
	}
}

func displayPathLabel(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	base := filepath.Base(trimmed)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return "configured"
	}
	return base
}
