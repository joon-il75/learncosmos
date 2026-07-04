package admin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AdminHandler) GetRecommendationRolloutCheckpoints(c *gin.Context) {
	ctx := c.Request.Context()
	state, err := h.loadActiveRecommendationRolloutState(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rollout state"})
		return
	}

	provider := "embedding_gemma"
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	endpointConfigured := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")) != ""

	var readyRows int64
	var failedRows int64
	if err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'ready'),
			COUNT(*) FILTER (WHERE status = 'failed')
		FROM content_embeddings
		WHERE provider = $1
		  AND model = $2
	`, provider, model).Scan(&readyRows, &failedRows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shadow counts"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shadow target count"})
		return
	}

	var compareRuns int64
	var shadowAttempted int64
	var shadowSucceeded int64
	if err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN metrics_snapshot ? 'shadow_attempted' THEN (metrics_snapshot->>'shadow_attempted')::integer ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN metrics_snapshot ? 'shadow_succeeded' THEN (metrics_snapshot->>'shadow_succeeded')::integer ELSE 0 END), 0)
		FROM recommendation_debug_runs
		WHERE run_type = 'recommendation_compare'
	`).Scan(&compareRuns, &shadowAttempted, &shadowSucceeded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load recommendation run metrics"})
		return
	}

	var totalLabels int64
	var goodFit int64
	var irrelevant int64
	var duplicate int64
	var problem int64
	if err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE label = 'good_fit'),
			COUNT(*) FILTER (WHERE label = 'irrelevant'),
			COUNT(*) FILTER (WHERE label = 'duplicate'),
			COUNT(*) FILTER (WHERE label IN ('broken_link', 'low_quality', 'unsafe'))
		FROM recommendation_debug_labels
	`).Scan(&totalLabels, &goodFit, &irrelevant, &duplicate, &problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load label metrics"})
		return
	}

	goodFitRate := safeRate(goodFit, totalLabels)
	irrelevantRate := safeRate(irrelevant, totalLabels)
	duplicateRate := safeRate(duplicate, totalLabels)
	problemRate := safeRate(problem, totalLabels)
	qualityStatus := "pending"
	qualityReason := "운영자 라벨이 아직 부족합니다."
	if totalLabels > 0 {
		qualityStatus = "warning"
	}
	if totalLabels >= 10 && goodFitRate >= 0.70 && irrelevantRate <= 0.10 && duplicateRate <= 0.10 {
		qualityStatus = "passed"
		qualityReason = "운영자 라벨 1차 품질 기준을 만족합니다."
	}
	if totalLabels >= 10 && (irrelevantRate > 0.15 || problemRate > 0.15) {
		qualityStatus = "failed"
		qualityReason = "무관 또는 문제 후보 비율이 rollback 기준에 가깝습니다."
	}

	rankerArtifact := h.buildRecommendationRankerArtifactStatus(ctx, state)
	rankerConfigured := rankerArtifact.Ready
	rollbackPolicyConfigured := state != nil && strings.TrimSpace(state.RollbackPolicySnapshot) != "" && strings.TrimSpace(state.RollbackPolicySnapshot) != "{}"
	learnerSignal, err := h.evaluateRecommendationLearnerSignalReadiness(ctx, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learner rollout signal"})
		return
	}
	checkpoints := []gin.H{
		buildRolloutCheckpoint("embedding_endpoint", "EmbeddingGemma endpoint", endpointConfigured, "endpoint가 설정되어 있습니다.", "EMBEDDING_GEMMA_ENDPOINT 설정이 필요합니다.", gin.H{
			"provider": provider,
			"model":    model,
		}),
		buildRolloutCheckpoint("shadow_backfill", "Primary embedding backfill", targetCount == 0, "남은 embedding backfill 대상이 없습니다.", "embedding backfill 대상이 남아 있습니다.", gin.H{
			"ready_rows":        readyRows,
			"failed_rows":       failedRows,
			"remaining_targets": targetCount,
		}),
		buildRolloutCheckpoint("shadow_runs", "Primary 비교 실행", shadowSucceeded >= 50, "Embedding 비교 실행 수가 기준을 만족합니다.", "추천 비교 실행 또는 embedding 성공 수가 부족합니다.", gin.H{
			"compare_runs":     compareRuns,
			"shadow_attempted": shadowAttempted,
			"shadow_succeeded": shadowSucceeded,
			"minimum_required": int64(50),
		}),
		buildRolloutCheckpoint("label_quality", "운영자 라벨 품질", qualityStatus == "passed", qualityReason, qualityReason, gin.H{
			"total_labels":     totalLabels,
			"minimum_required": int64(10),
			"good_fit_rate":    goodFitRate,
			"irrelevant_rate":  irrelevantRate,
			"duplicate_rate":   duplicateRate,
			"problem_rate":     problemRate,
			"quality_status":   qualityStatus,
		}),
		buildRolloutCheckpoint("ranker_model", "LightGBM Ranker 모델", rankerConfigured, "ranker model artifact가 준비되어 있습니다.", rankerArtifact.Reason, gin.H{
			"ranker_model_version":   rankerArtifact.RankerModelVersion,
			"feature_schema_version": rankerArtifact.FeatureSchemaVersion,
			"model_path":             rankerArtifact.ModelPath,
			"manifest_path":          rankerArtifact.ManifestPath,
			"model_file_exists":      rankerArtifact.ModelFileExists,
			"manifest_exists":        rankerArtifact.ManifestExists,
			"manifest_version_match": rankerArtifact.ManifestVersionMatch,
			"manifest_schema_match":  rankerArtifact.ManifestSchemaMatch,
		}),
		buildRolloutCheckpoint("rollback_policy", "Rollback 정책", rollbackPolicyConfigured, "rollback 정책이 저장되어 있습니다.", "rollback 정책 snapshot이 필요합니다.", gin.H{}),
		buildRolloutCheckpoint("learner_signal", "Learner signal", learnerSignal.Ready, learnerSignal.Reason, learnerSignal.Reason, gin.H{
			"required":             learnerSignal.Required,
			"exposure_count":       learnerSignal.ExposureCount,
			"minimum_exposure":     learnerSignal.MinimumExposure,
			"replacement_count":    learnerSignal.ReplacementCount,
			"broken_link_count":    learnerSignal.BrokenLinkCount,
			"wrong_content_count":  learnerSignal.WrongContentCount,
			"selection_rate":       learnerSignal.SelectionRate,
			"min_selection_rate":   learnerSignal.MinSelectionRate,
			"broken_link_rate":     learnerSignal.BrokenLinkRate,
			"max_broken_link_rate": learnerSignal.MaxBrokenLinkRate,
			"wrong_content_rate":   learnerSignal.WrongContentRate,
			"max_wrong_rate":       learnerSignal.MaxWrongRate,
		}),
	}

	passedCount := 0
	failedCount := 0
	for _, checkpoint := range checkpoints {
		switch checkpoint["status"] {
		case "passed":
			passedCount++
		case "failed":
			failedCount++
		}
	}
	canPromote := failedCount == 0 && passedCount == len(checkpoints)
	nextMode := nextRecommendationRolloutMode("")
	currentMode := ""
	if state != nil {
		currentMode = state.Mode
		nextMode = nextRecommendationRolloutMode(state.Mode)
	}

	c.JSON(http.StatusOK, gin.H{
		"state": state,
		"summary": gin.H{
			"current_mode":        currentMode,
			"next_mode":           nextMode,
			"can_promote":         canPromote,
			"passed_count":        passedCount,
			"total_count":         len(checkpoints),
			"quality_gate_status": qualityStatus,
			"quality_reason":      qualityReason,
		},
		"metrics": gin.H{
			"shadow": gin.H{
				"endpoint_configured": endpointConfigured,
				"provider":            provider,
				"model":               model,
				"ready_rows":          readyRows,
				"failed_rows":         failedRows,
				"target_count":        targetCount,
				"compare_runs":        compareRuns,
				"shadow_attempted":    shadowAttempted,
				"shadow_succeeded":    shadowSucceeded,
			},
			"labels": gin.H{
				"total_labels":    totalLabels,
				"good_fit_rate":   goodFitRate,
				"irrelevant_rate": irrelevantRate,
				"duplicate_rate":  duplicateRate,
				"problem_rate":    problemRate,
			},
			"learner_signal": gin.H{
				"required":             learnerSignal.Required,
				"ready":                learnerSignal.Ready,
				"reason":               learnerSignal.Reason,
				"exposure_count":       learnerSignal.ExposureCount,
				"minimum_exposure":     learnerSignal.MinimumExposure,
				"replacement_count":    learnerSignal.ReplacementCount,
				"broken_link_count":    learnerSignal.BrokenLinkCount,
				"wrong_content_count":  learnerSignal.WrongContentCount,
				"selection_rate":       learnerSignal.SelectionRate,
				"min_selection_rate":   learnerSignal.MinSelectionRate,
				"broken_link_rate":     learnerSignal.BrokenLinkRate,
				"max_broken_link_rate": learnerSignal.MaxBrokenLinkRate,
				"wrong_content_rate":   learnerSignal.WrongContentRate,
				"max_wrong_rate":       learnerSignal.MaxWrongRate,
			},
		},
		"checkpoints": checkpoints,
	})
}

func (h *AdminHandler) evaluateRecommendationRolloutReadiness(ctx context.Context) (gin.H, error) {
	state, err := h.loadActiveRecommendationRolloutState(ctx)
	if err != nil {
		return nil, err
	}
	provider := "embedding_gemma"
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		model = "embedding-gemma"
	}
	endpointConfigured := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")) != ""

	var readyRows int64
	var targetCount int64
	var shadowSucceeded int64
	var totalLabels int64
	var goodFit int64
	var irrelevant int64
	var duplicate int64
	var problem int64
	if err := h.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM content_embeddings
		WHERE provider = $1
		  AND model = $2
		  AND status = 'ready'
	`, provider, model).Scan(&readyRows); err != nil {
		return nil, err
	}
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
		return nil, err
	}
	if err := h.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN metrics_snapshot ? 'shadow_succeeded' THEN (metrics_snapshot->>'shadow_succeeded')::integer ELSE 0 END), 0)
		FROM recommendation_debug_runs
		WHERE run_type = 'recommendation_compare'
	`).Scan(&shadowSucceeded); err != nil {
		return nil, err
	}
	if err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE label = 'good_fit'),
			COUNT(*) FILTER (WHERE label = 'irrelevant'),
			COUNT(*) FILTER (WHERE label = 'duplicate'),
			COUNT(*) FILTER (WHERE label IN ('broken_link', 'low_quality', 'unsafe'))
		FROM recommendation_debug_labels
	`).Scan(&totalLabels, &goodFit, &irrelevant, &duplicate, &problem); err != nil {
		return nil, err
	}

	goodFitRate := safeRate(goodFit, totalLabels)
	irrelevantRate := safeRate(irrelevant, totalLabels)
	duplicateRate := safeRate(duplicate, totalLabels)
	problemRate := safeRate(problem, totalLabels)
	embeddingReady := endpointConfigured
	backfillReady := targetCount == 0
	shadowReady := shadowSucceeded >= 50
	labelReady := totalLabels >= 10 && goodFitRate >= 0.70 && irrelevantRate <= 0.10 && duplicateRate <= 0.10 && problemRate <= 0.15
	rankerReady := state != nil && strings.TrimSpace(state.RankerModelVersion) != ""
	rollbackReady := state != nil && strings.TrimSpace(state.RollbackPolicySnapshot) != "" && strings.TrimSpace(state.RollbackPolicySnapshot) != "{}"
	learnerSignal, err := h.evaluateRecommendationLearnerSignalReadiness(ctx, state)
	if err != nil {
		return nil, err
	}
	canPromote := embeddingReady && backfillReady && shadowReady && labelReady && rankerReady && rollbackReady && learnerSignal.Ready

	return gin.H{
		"can_promote": canPromote,
		"snapshot": gin.H{
			"embedding_endpoint": embeddingReady,
			"shadow_backfill":    backfillReady,
			"shadow_runs":        shadowReady,
			"label_quality":      labelReady,
			"ranker_model":       rankerReady,
			"rollback_policy":    rollbackReady,
			"learner_signal":     learnerSignal.Ready,
			"metrics": gin.H{
				"provider":          provider,
				"model":             model,
				"ready_rows":        readyRows,
				"remaining_targets": targetCount,
				"shadow_succeeded":  shadowSucceeded,
				"total_labels":      totalLabels,
				"good_fit_rate":     goodFitRate,
				"irrelevant_rate":   irrelevantRate,
				"duplicate_rate":    duplicateRate,
				"problem_rate":      problemRate,
				"learner_signal": gin.H{
					"required":             learnerSignal.Required,
					"exposure_count":       learnerSignal.ExposureCount,
					"minimum_exposure":     learnerSignal.MinimumExposure,
					"selection_rate":       learnerSignal.SelectionRate,
					"min_selection_rate":   learnerSignal.MinSelectionRate,
					"broken_link_rate":     learnerSignal.BrokenLinkRate,
					"max_broken_link_rate": learnerSignal.MaxBrokenLinkRate,
					"wrong_content_rate":   learnerSignal.WrongContentRate,
					"max_wrong_rate":       learnerSignal.MaxWrongRate,
				},
			},
		},
	}, nil
}

func (h *AdminHandler) evaluateRecommendationLearnerSignalReadiness(ctx context.Context, state *recommendationRolloutStateRow) (recommendationLearnerSignalReadiness, error) {
	result := recommendationLearnerSignalReadiness{
		Ready:             true,
		Reason:            "learner traffic 단계 전에는 learner signal을 승격 조건으로 요구하지 않습니다.",
		MinSelectionRate:  0.05,
		MaxBrokenLinkRate: 0.10,
		MaxWrongRate:      0.05,
	}
	if state == nil {
		return result, nil
	}

	switch state.Mode {
	case "learner_10_percent":
		result.Required = true
		result.MinimumExposure = 50
	case "learner_50_percent":
		result.Required = true
		result.MinimumExposure = 200
		result.MaxBrokenLinkRate = 0.08
		result.MaxWrongRate = 0.03
	case "learner_100_percent":
		result.Required = true
		result.MinimumExposure = 500
		result.MaxBrokenLinkRate = 0.08
		result.MaxWrongRate = 0.03
	default:
		return result, nil
	}

	rolloutStateID, err := uuid.Parse(state.ID)
	if err != nil {
		return result, fmt.Errorf("parse rollout state id for learner signal: %w", err)
	}
	if err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE event_type = 'recommendation_exposed'),
			COUNT(*) FILTER (WHERE event_type = 'material_replaced'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'broken_link'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'wrong_content')
		FROM recommendation_rollout_learner_events
		WHERE rollout_state_id = $1
	`, rolloutStateID).Scan(&result.ExposureCount, &result.ReplacementCount, &result.BrokenLinkCount, &result.WrongContentCount); err != nil {
		return result, err
	}

	result.SelectionRate = safeRate(result.ReplacementCount, result.ExposureCount)
	result.BrokenLinkRate = safeRate(result.BrokenLinkCount, result.ExposureCount)
	result.WrongContentRate = safeRate(result.WrongContentCount, result.ExposureCount)
	result.Ready = result.ExposureCount >= result.MinimumExposure &&
		result.SelectionRate >= result.MinSelectionRate &&
		result.BrokenLinkRate <= result.MaxBrokenLinkRate &&
		result.WrongContentRate <= result.MaxWrongRate
	if result.Ready {
		result.Reason = "learner signal 기준을 만족합니다."
		return result, nil
	}
	result.Reason = fmt.Sprintf(
		"learner signal 부족: exposure %d/%d, selection %.1f%% 이상 필요, broken %.1f%% 이하, wrong %.1f%% 이하 기준입니다.",
		result.ExposureCount,
		result.MinimumExposure,
		result.MinSelectionRate*100,
		result.MaxBrokenLinkRate*100,
		result.MaxWrongRate*100,
	)
	return result, nil
}

func buildRolloutCheckpoint(id, label string, passed bool, passedReason, pendingReason string, metrics gin.H) gin.H {
	status := "pending"
	reason := pendingReason
	if passed {
		status = "passed"
		reason = passedReason
	}
	return gin.H{
		"id":      id,
		"label":   label,
		"status":  status,
		"reason":  reason,
		"metrics": metrics,
	}
}
