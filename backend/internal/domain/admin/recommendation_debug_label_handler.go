package admin

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AdminHandler) SaveRecommendationDebugLabel(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CandidateKey = strings.TrimSpace(req.CandidateKey)
	req.ContentID = strings.TrimSpace(req.ContentID)
	req.URL = strings.TrimSpace(req.URL)
	req.Label = strings.TrimSpace(req.Label)
	req.Note = strings.TrimSpace(req.Note)
	if req.CandidateKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "candidate_key is required"})
		return
	}
	if !isRecommendationDebugLabelAllowed(req.Label) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid label"})
		return
	}
	if len([]rune(req.Note)) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note too long"})
		return
	}

	scenarioID, _ := uuid.Parse(scenario.ID)
	var runID *uuid.UUID
	if strings.TrimSpace(req.RunID) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(req.RunID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
			return
		}
		var exists bool
		if err := h.db.QueryRow(c.Request.Context(), `
			SELECT EXISTS (
				SELECT 1
				FROM recommendation_debug_runs
				WHERE id = $1
				  AND scenario_id = $2
				  AND run_type = 'recommendation_compare'
			)
		`, parsed, scenarioID).Scan(&exists); err != nil || !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "recommendation run not found"})
			return
		}
		runID = &parsed
	}

	var contentID *uuid.UUID
	if req.ContentID != "" {
		parsed, err := uuid.Parse(req.ContentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content_id"})
			return
		}
		contentID = &parsed
	}
	featureSnapshot := strings.TrimSpace(string(req.FeatureSnapshot))
	if featureSnapshot == "" || !json.Valid([]byte(featureSnapshot)) {
		featureSnapshot = "{}"
	}
	var baselineRank any
	if req.BaselineRank != nil && *req.BaselineRank > 0 {
		baselineRank = *req.BaselineRank
	}
	var shadowRank any
	if req.ShadowRank != nil && *req.ShadowRank > 0 {
		shadowRank = *req.ShadowRank
	}

	var item recommendationDebugLabelRow
	if err := h.db.QueryRow(c.Request.Context(), `
		INSERT INTO recommendation_debug_labels (
			scenario_id,
			run_id,
			created_by,
			candidate_key,
			content_id,
			url,
			label,
			note,
			feature_snapshot,
			baseline_rank,
			shadow_rank
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10, $11)
		ON CONFLICT (run_id, candidate_key)
		WHERE run_id IS NOT NULL
		DO UPDATE SET
			created_by = EXCLUDED.created_by,
			content_id = EXCLUDED.content_id,
			url = EXCLUDED.url,
			label = EXCLUDED.label,
			note = EXCLUDED.note,
			feature_snapshot = EXCLUDED.feature_snapshot,
			baseline_rank = EXCLUDED.baseline_rank,
			shadow_rank = EXCLUDED.shadow_rank,
			updated_at = NOW()
		RETURNING
			id::text,
			scenario_id::text,
			COALESCE(run_id::text, ''),
			COALESCE(created_by::text, ''),
			candidate_key,
			COALESCE(content_id::text, ''),
			url,
			label,
			note,
			feature_snapshot::text,
			baseline_rank,
			shadow_rank,
			created_at,
			updated_at
	`, scenarioID, runID, recommendationDebugNullableAdminID(adminID), req.CandidateKey, contentID, req.URL, req.Label, req.Note, featureSnapshot, baselineRank, shadowRank).Scan(
		&item.ID,
		&item.ScenarioID,
		&item.RunID,
		&item.CreatedBy,
		&item.CandidateKey,
		&item.ContentID,
		&item.URL,
		&item.Label,
		&item.Note,
		&item.FeatureSnapshot,
		&item.BaselineRank,
		&item.ShadowRank,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save label"})
		return
	}

	if _, err := h.db.Exec(c.Request.Context(), `
		UPDATE recommendation_debug_scenarios
		SET status = CASE
				WHEN status IN ('recommendation_tested', 'labelled') THEN 'labelled'
				ELSE status
			END,
			updated_at = NOW()
		WHERE id = $1
	`, scenarioID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update scenario label status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": item})
}

func (h *AdminHandler) GetRecommendationDebugLabelSummary(c *gin.Context) {
	scenario, _, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}
	scenarioID, _ := uuid.Parse(scenario.ID)

	rows, err := h.db.Query(c.Request.Context(), `
		SELECT label, COUNT(*)
		FROM recommendation_debug_labels
		WHERE scenario_id = $1
		GROUP BY label
		ORDER BY label
	`, scenarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load label summary"})
		return
	}
	defer rows.Close()

	labelCounts := gin.H{}
	var total int64
	for rows.Next() {
		var label string
		var count int64
		if err := rows.Scan(&label, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan label summary"})
			return
		}
		labelCounts[label] = count
		total += count
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load label summary"})
		return
	}

	goodFit := labelCountValue(labelCounts, "good_fit")
	irrelevant := labelCountValue(labelCounts, "irrelevant")
	duplicate := labelCountValue(labelCounts, "duplicate")
	problem := labelCountValue(labelCounts, "broken_link") + labelCountValue(labelCounts, "low_quality") + labelCountValue(labelCounts, "unsafe")

	goodFitRate := safeRate(goodFit, total)
	irrelevantRate := safeRate(irrelevant, total)
	duplicateRate := safeRate(duplicate, total)
	problemRate := safeRate(problem, total)
	gateStatus := "pending"
	gateReason := "라벨이 아직 없습니다."
	if total > 0 {
		gateStatus = "warning"
		gateReason = "운영자 라벨 샘플이 더 필요합니다."
	}
	if total >= 10 && goodFitRate >= 0.70 && irrelevantRate <= 0.10 && duplicateRate <= 0.10 {
		gateStatus = "passed"
		gateReason = "1차 라벨 품질 기준을 만족합니다."
	}
	if total >= 10 && (irrelevantRate > 0.15 || problemRate > 0.15) {
		gateStatus = "failed"
		gateReason = "무관 또는 문제 후보 비율이 높습니다."
	}

	var latestRunID string
	var latestRunAt sql.NullTime
	_ = h.db.QueryRow(c.Request.Context(), `
		SELECT id::text, created_at
		FROM recommendation_debug_runs
		WHERE scenario_id = $1
		  AND run_type = 'recommendation_compare'
		ORDER BY created_at DESC
		LIMIT 1
	`, scenarioID).Scan(&latestRunID, &latestRunAt)

	summary := gin.H{
		"scenario_id":      scenario.ID,
		"latest_run_id":    latestRunID,
		"total_labels":     total,
		"label_counts":     labelCounts,
		"good_fit_rate":    goodFitRate,
		"irrelevant_rate":  irrelevantRate,
		"duplicate_rate":   duplicateRate,
		"problem_rate":     problemRate,
		"quality_status":   gateStatus,
		"quality_reason":   gateReason,
		"minimum_required": int64(10),
	}
	if latestRunAt.Valid {
		summary["latest_run_at"] = latestRunAt.Time
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (h *AdminHandler) ExportRecommendationDebugLabels(c *gin.Context) {
	scenario, _, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}
	scenarioID, _ := uuid.Parse(scenario.ID)
	format := strings.ToLower(strings.TrimSpace(c.DefaultQuery("format", "jsonl")))
	if format != "jsonl" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be jsonl or csv"})
		return
	}

	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			l.id::text,
			l.scenario_id::text,
			s.course_title,
			s.initial_user_intent,
			s.status,
			COALESCE(l.run_id::text, ''),
			r.created_at,
			l.candidate_key,
			COALESCE(l.content_id::text, ''),
			l.url,
			l.label,
			l.note,
			l.feature_snapshot::text,
			l.baseline_rank,
			l.shadow_rank,
			COALESCE(r.request_snapshot::text, '{}'::text),
			COALESCE(r.baseline_result_snapshot::text, '{}'::text),
			COALESCE(r.shadow_result_snapshot::text, '{}'::text),
			COALESCE(r.metrics_snapshot::text, '{}'::text),
			l.updated_at
		FROM recommendation_debug_labels l
		JOIN recommendation_debug_scenarios s ON s.id = l.scenario_id
		LEFT JOIN recommendation_debug_runs r ON r.id = l.run_id
		WHERE l.scenario_id = $1
		ORDER BY l.updated_at DESC, l.created_at DESC
	`, scenarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export labels"})
		return
	}
	defer rows.Close()

	items := make([]recommendationDebugLabelExportRow, 0)
	for rows.Next() {
		var item recommendationDebugLabelExportRow
		if err := rows.Scan(
			&item.LabelID,
			&item.ScenarioID,
			&item.CourseTitle,
			&item.InitialUserIntent,
			&item.ScenarioStatus,
			&item.RunID,
			&item.RunCreatedAt,
			&item.CandidateKey,
			&item.ContentID,
			&item.URL,
			&item.Label,
			&item.Note,
			&item.FeatureSnapshot,
			&item.BaselineRank,
			&item.ShadowRank,
			&item.RequestSnapshot,
			&item.BaselineResultSnapshot,
			&item.ShadowResultSnapshot,
			&item.MetricsSnapshot,
			&item.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan label export"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export labels"})
		return
	}

	filename := fmt.Sprintf("recommendation-labels-%s.%s", scenario.ID, format)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	if format == "csv" {
		var builder strings.Builder
		writer := csv.NewWriter(&builder)
		_ = writer.Write([]string{
			"label_id",
			"scenario_id",
			"course_title",
			"initial_user_intent",
			"scenario_status",
			"run_id",
			"run_created_at",
			"candidate_key",
			"content_id",
			"url",
			"label",
			"note",
			"baseline_rank",
			"shadow_rank",
			"label_route",
			"label_ranker_score",
			"label_ranker_route_rank",
			"label_ranker_rerank_rank",
			"label_ranker_rank_delta",
			"label_ranker_provider",
			"label_ranker_model_version",
			"baseline_feature_snapshot",
			"shadow_feature_snapshot",
			"baseline_ranker_score",
			"baseline_ranker_rerank_rank",
			"baseline_ranker_rank_delta",
			"shadow_ranker_score",
			"shadow_ranker_rerank_rank",
			"shadow_ranker_rank_delta",
			"feature_snapshot",
			"request_snapshot",
			"baseline_result_snapshot",
			"shadow_result_snapshot",
			"metrics_snapshot",
			"updated_at",
		})
		for _, item := range items {
			rankerExport := recommendationDebugLabelExportRankerSnapshot(item)
			_ = writer.Write([]string{
				item.LabelID,
				item.ScenarioID,
				item.CourseTitle,
				item.InitialUserIntent,
				item.ScenarioStatus,
				item.RunID,
				formatNullTimeRFC3339(item.RunCreatedAt),
				item.CandidateKey,
				item.ContentID,
				item.URL,
				item.Label,
				item.Note,
				formatNullInt32(item.BaselineRank),
				formatNullInt32(item.ShadowRank),
				stringValue(rankerExport["label_route"]),
				stringValue(rankerExport["label_ranker_score"]),
				stringValue(rankerExport["label_ranker_route_rank"]),
				stringValue(rankerExport["label_ranker_rerank_rank"]),
				stringValue(rankerExport["label_ranker_rank_delta"]),
				stringValue(rankerExport["label_ranker_provider"]),
				stringValue(rankerExport["label_ranker_model_version"]),
				jsonStringValue(rankerExport["baseline_feature_snapshot"]),
				jsonStringValue(rankerExport["shadow_feature_snapshot"]),
				stringValue(rankerExport["baseline_ranker_score"]),
				stringValue(rankerExport["baseline_ranker_rerank_rank"]),
				stringValue(rankerExport["baseline_ranker_rank_delta"]),
				stringValue(rankerExport["shadow_ranker_score"]),
				stringValue(rankerExport["shadow_ranker_rerank_rank"]),
				stringValue(rankerExport["shadow_ranker_rank_delta"]),
				item.FeatureSnapshot,
				item.RequestSnapshot,
				item.BaselineResultSnapshot,
				item.ShadowResultSnapshot,
				item.MetricsSnapshot,
				item.UpdatedAt.Format(time.RFC3339),
			})
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode csv"})
			return
		}
		c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(builder.String()))
		return
	}

	var builder strings.Builder
	for _, item := range items {
		rankerExport := recommendationDebugLabelExportRankerSnapshot(item)
		line, err := json.Marshal(gin.H{
			"label_id":                 item.LabelID,
			"scenario_id":              item.ScenarioID,
			"course_title":             item.CourseTitle,
			"initial_user_intent":      item.InitialUserIntent,
			"scenario_status":          item.ScenarioStatus,
			"run_id":                   item.RunID,
			"run_created_at":           formatNullTimeRFC3339(item.RunCreatedAt),
			"candidate_key":            item.CandidateKey,
			"content_id":               item.ContentID,
			"url":                      item.URL,
			"label":                    item.Label,
			"note":                     item.Note,
			"baseline_rank":            nullInt32Value(item.BaselineRank),
			"shadow_rank":              nullInt32Value(item.ShadowRank),
			"feature_snapshot":         rawJSONOrEmpty(item.FeatureSnapshot),
			"ranker_export":            rankerExport,
			"request_snapshot":         rawJSONOrEmpty(item.RequestSnapshot),
			"baseline_result_snapshot": rawJSONOrEmpty(item.BaselineResultSnapshot),
			"shadow_result_snapshot":   rawJSONOrEmpty(item.ShadowResultSnapshot),
			"metrics_snapshot":         rawJSONOrEmpty(item.MetricsSnapshot),
			"updated_at":               item.UpdatedAt.Format(time.RFC3339),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode jsonl"})
			return
		}
		builder.Write(line)
		builder.WriteByte('\n')
	}
	c.Data(http.StatusOK, "application/x-ndjson; charset=utf-8", []byte(builder.String()))
}
