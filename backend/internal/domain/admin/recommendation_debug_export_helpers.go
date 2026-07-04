package admin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func recommendationDebugLabelExportRankerSnapshot(item recommendationDebugLabelExportRow) gin.H {
	labelFeature := recommendationDebugDecodeJSONObject(item.FeatureSnapshot)
	baselineCandidate, shadowCandidate := recommendationDebugFindExportCandidates(item)
	baselineFeature := recommendationDebugObjectField(baselineCandidate, "feature_snapshot")
	shadowFeature := recommendationDebugObjectField(shadowCandidate, "feature_snapshot")
	return gin.H{
		"label_route":                 stringValue(labelFeature["route"]),
		"label_ranker_score":          labelFeature["ranker_score"],
		"label_ranker_route_rank":     labelFeature["ranker_route_rank"],
		"label_ranker_rerank_rank":    labelFeature["ranker_rerank_rank"],
		"label_ranker_rank_delta":     labelFeature["ranker_rank_delta"],
		"label_ranker_provider":       stringValue(labelFeature["ranker_provider"]),
		"label_ranker_model_version":  stringValue(labelFeature["ranker_model_version"]),
		"baseline_candidate_snapshot": baselineCandidate,
		"shadow_candidate_snapshot":   shadowCandidate,
		"baseline_feature_snapshot":   baselineFeature,
		"shadow_feature_snapshot":     shadowFeature,
		"baseline_ranker_score":       baselineCandidate["ranker_score"],
		"baseline_ranker_route_rank":  baselineCandidate["ranker_route_rank"],
		"baseline_ranker_rerank_rank": baselineCandidate["ranker_rerank_rank"],
		"baseline_ranker_rank_delta":  baselineCandidate["ranker_rank_delta"],
		"shadow_ranker_score":         shadowCandidate["ranker_score"],
		"shadow_ranker_route_rank":    shadowCandidate["ranker_route_rank"],
		"shadow_ranker_rerank_rank":   shadowCandidate["ranker_rerank_rank"],
		"shadow_ranker_rank_delta":    shadowCandidate["ranker_rank_delta"],
	}
}

func recommendationDebugFindExportCandidates(item recommendationDebugLabelExportRow) (gin.H, gin.H) {
	runSnapshot := recommendationDebugDecodeJSONObject(item.BaselineResultSnapshot)
	lessons, _ := runSnapshot["lessons"].([]any)
	for _, lessonValue := range lessons {
		lessonObj := recommendationDebugAnyMap(lessonValue)
		lessonMeta := recommendationDebugObjectField(lessonObj, "lesson")
		lessonID := stringValue(lessonMeta["lesson_id"])
		baselineCandidates := recommendationDebugSectionCandidates(lessonObj, "baseline")
		for idx, candidateValue := range baselineCandidates {
			candidate := recommendationDebugAnyMap(candidateValue)
			if recommendationDebugExportCandidateKey(lessonID, candidate, idx) != item.CandidateKey {
				continue
			}
			return candidate, recommendationDebugFindMatchingShadowCandidate(lessonObj, candidate)
		}
	}
	return gin.H{}, gin.H{}
}

func recommendationDebugFindMatchingShadowCandidate(lessonObj gin.H, baselineCandidate gin.H) gin.H {
	identity := recommendationDebugExportCandidateIdentity(baselineCandidate)
	if identity == "" {
		return gin.H{}
	}
	for _, candidateValue := range recommendationDebugSectionCandidates(lessonObj, "shadow") {
		candidate := recommendationDebugAnyMap(candidateValue)
		if recommendationDebugExportCandidateIdentity(candidate) == identity {
			return candidate
		}
	}
	return gin.H{}
}

func recommendationDebugSectionCandidates(lessonObj gin.H, section string) []any {
	sectionObj := recommendationDebugObjectField(lessonObj, section)
	candidates, _ := sectionObj["candidates"].([]any)
	return candidates
}

func recommendationDebugExportCandidateKey(lessonID string, candidate gin.H, idx int) string {
	return fmt.Sprintf("%s:%s:%d", lessonID, recommendationDebugExportCandidateIdentity(candidate), idx+1)
}

func recommendationDebugExportCandidateIdentity(candidate gin.H) string {
	if value := stringValue(candidate["content_id"]); value != "" {
		return value
	}
	if value := stringValue(candidate["external_url"]); value != "" {
		return value
	}
	return stringValue(candidate["title"])
}

func recommendationDebugDecodeJSONObject(value string) gin.H {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return gin.H{}
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil {
		return gin.H{}
	}
	return gin.H(decoded)
}

func recommendationDebugObjectField(obj gin.H, key string) gin.H {
	return recommendationDebugAnyMap(obj[key])
}

func recommendationDebugAnyMap(value any) gin.H {
	switch typed := value.(type) {
	case gin.H:
		return typed
	case map[string]any:
		return gin.H(typed)
	default:
		return gin.H{}
	}
}
