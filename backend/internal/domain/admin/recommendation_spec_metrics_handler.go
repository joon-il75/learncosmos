package admin

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	recommendationSpecMetricsDefaultDays = 7
	recommendationSpecMetricsMaxDays     = 90
	recommendationSpecMetricsShortRunes  = 8
)

type recommendationSpecMetricQueryRow struct {
	Query string `json:"query"`
	Count int64  `json:"count"`
}

type recommendationSpecRolloutMetrics struct {
	RecommendationExposed int64    `json:"recommendation_exposed"`
	MaterialReplaced      int64    `json:"material_replaced"`
	BrokenLinkReported    int64    `json:"broken_link_reported"`
	WrongContentReported  int64    `json:"wrong_content_reported"`
	AverageCandidateCount *float64 `json:"average_candidate_count,omitempty"`
}

type recommendationSpecMetricsResponse struct {
	WindowDays              int                                `json:"window_days"`
	WindowStartedAt         time.Time                          `json:"window_started_at"`
	WindowEndedAt           time.Time                          `json:"window_ended_at"`
	RecommendationSearches  int64                              `json:"recommendation_searches"`
	SearchSpecSourceCounts  map[string]int64                   `json:"search_spec_source_counts"`
	ResolutionSourceCounts  map[string]int64                   `json:"resolution_source_counts"`
	EmptyEffectiveQueryRate *float64                           `json:"empty_effective_query_rate,omitempty"`
	ShortEffectiveQueryRate *float64                           `json:"short_effective_query_rate,omitempty"`
	TopEffectiveQueries     []recommendationSpecMetricQueryRow `json:"top_effective_queries"`
	RolloutEvents           recommendationSpecRolloutMetrics   `json:"rollout_events"`
}

func (h *AdminHandler) GetRecommendationSpecMetrics(c *gin.Context) {
	days, err := parseRecommendationSpecMetricsDays(c.Query("days"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endedAt := time.Now().UTC()
	startedAt := endedAt.AddDate(0, 0, -days)

	summary, err := h.collectRecommendationSpecMetrics(c.Request.Context(), startedAt, endedAt, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load recommendation spec metrics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": summary})
}

func parseRecommendationSpecMetricsDays(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return recommendationSpecMetricsDefaultDays, nil
	}
	days, err := strconv.Atoi(raw)
	if err != nil || days <= 0 {
		return 0, fmt.Errorf("days must be a positive integer")
	}
	if days > recommendationSpecMetricsMaxDays {
		return 0, fmt.Errorf("days must be %d or less", recommendationSpecMetricsMaxDays)
	}
	return days, nil
}

func (h *AdminHandler) collectRecommendationSpecMetrics(ctx context.Context, startedAt, endedAt time.Time, days int) (recommendationSpecMetricsResponse, error) {
	searches, emptyQueries, shortQueries, err := h.countRecommendationSpecSearches(ctx, startedAt, endedAt)
	if err != nil {
		return recommendationSpecMetricsResponse{}, err
	}
	specSources, err := h.countRecommendationSpecMetadataValue(ctx, startedAt, endedAt, "search_spec_source")
	if err != nil {
		return recommendationSpecMetricsResponse{}, err
	}
	resolutionSources, err := h.countRecommendationSpecMetadataValue(ctx, startedAt, endedAt, "resolution_source")
	if err != nil {
		return recommendationSpecMetricsResponse{}, err
	}
	topQueries, err := h.listRecommendationSpecTopEffectiveQueries(ctx, startedAt, endedAt, 10)
	if err != nil {
		return recommendationSpecMetricsResponse{}, err
	}
	rolloutEvents, err := h.collectRecommendationSpecRolloutMetrics(ctx, startedAt, endedAt)
	if err != nil {
		return recommendationSpecMetricsResponse{}, err
	}

	return recommendationSpecMetricsResponse{
		WindowDays:              days,
		WindowStartedAt:         startedAt,
		WindowEndedAt:           endedAt,
		RecommendationSearches:  searches,
		SearchSpecSourceCounts:  specSources,
		ResolutionSourceCounts:  resolutionSources,
		EmptyEffectiveQueryRate: ratioFloat64Ptr(emptyQueries, searches),
		ShortEffectiveQueryRate: ratioFloat64Ptr(shortQueries, searches),
		TopEffectiveQueries:     topQueries,
		RolloutEvents:           rolloutEvents,
	}, nil
}

func (h *AdminHandler) countRecommendationSpecSearches(ctx context.Context, startedAt, endedAt time.Time) (int64, int64, int64, error) {
	var total int64
	var emptyQuery int64
	var shortQuery int64
	err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE btrim(COALESCE(metadata->>'effective_query', '')) = ''),
			COUNT(*) FILTER (
				WHERE btrim(COALESCE(metadata->>'effective_query', '')) <> ''
				  AND char_length(btrim(COALESCE(metadata->>'effective_query', ''))) <= $3
			)
		FROM ai_point_transactions
		WHERE feature = 'explorer_recommendation_search'
		  AND created_at >= $1
		  AND created_at < $2
	`, startedAt, endedAt, recommendationSpecMetricsShortRunes).Scan(&total, &emptyQuery, &shortQuery)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("count recommendation spec searches: %w", err)
	}
	return total, emptyQuery, shortQuery, nil
}

func (h *AdminHandler) countRecommendationSpecMetadataValue(ctx context.Context, startedAt, endedAt time.Time, key string) (map[string]int64, error) {
	rows, err := h.db.Query(ctx, `
		SELECT COALESCE(NULLIF(btrim(metadata->>$3), ''), 'empty') AS value, COUNT(*) AS count
		FROM ai_point_transactions
		WHERE feature = 'explorer_recommendation_search'
		  AND created_at >= $1
		  AND created_at < $2
		GROUP BY value
		ORDER BY count DESC, value ASC
	`, startedAt, endedAt, key)
	if err != nil {
		return nil, fmt.Errorf("count recommendation spec metadata %s: %w", key, err)
	}
	defer rows.Close()

	counts := map[string]int64{}
	for rows.Next() {
		var value string
		var count int64
		if err := rows.Scan(&value, &count); err != nil {
			return nil, fmt.Errorf("scan recommendation spec metadata %s: %w", key, err)
		}
		counts[value] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recommendation spec metadata %s: %w", key, err)
	}
	return counts, nil
}

func (h *AdminHandler) listRecommendationSpecTopEffectiveQueries(ctx context.Context, startedAt, endedAt time.Time, limit int) ([]recommendationSpecMetricQueryRow, error) {
	rows, err := h.db.Query(ctx, `
		SELECT btrim(metadata->>'effective_query') AS query, COUNT(*) AS count
		FROM ai_point_transactions
		WHERE feature = 'explorer_recommendation_search'
		  AND created_at >= $1
		  AND created_at < $2
		  AND btrim(COALESCE(metadata->>'effective_query', '')) <> ''
		GROUP BY query
		ORDER BY count DESC, query ASC
		LIMIT $3
	`, startedAt, endedAt, limit)
	if err != nil {
		return nil, fmt.Errorf("list recommendation spec top effective queries: %w", err)
	}
	defer rows.Close()

	queries := []recommendationSpecMetricQueryRow{}
	for rows.Next() {
		var item recommendationSpecMetricQueryRow
		if err := rows.Scan(&item.Query, &item.Count); err != nil {
			return nil, fmt.Errorf("scan recommendation spec effective query: %w", err)
		}
		queries = append(queries, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recommendation spec effective queries: %w", err)
	}
	return queries, nil
}

func (h *AdminHandler) collectRecommendationSpecRolloutMetrics(ctx context.Context, startedAt, endedAt time.Time) (recommendationSpecRolloutMetrics, error) {
	var metrics recommendationSpecRolloutMetrics
	var averageCandidateCount sql.NullFloat64
	err := h.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE event_type = 'recommendation_exposed'),
			COUNT(*) FILTER (WHERE event_type = 'material_replaced'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'broken_link'),
			COUNT(*) FILTER (WHERE event_type = 'material_reported' AND report_type = 'wrong_content'),
			(AVG(candidate_count) FILTER (WHERE event_type = 'recommendation_exposed'))::float8
		FROM recommendation_rollout_learner_events
		WHERE created_at >= $1
		  AND created_at < $2
	`, startedAt, endedAt).Scan(
		&metrics.RecommendationExposed,
		&metrics.MaterialReplaced,
		&metrics.BrokenLinkReported,
		&metrics.WrongContentReported,
		&averageCandidateCount,
	)
	if err != nil {
		return recommendationSpecRolloutMetrics{}, fmt.Errorf("collect recommendation spec rollout metrics: %w", err)
	}
	if averageCandidateCount.Valid {
		metrics.AverageCandidateCount = &averageCandidateCount.Float64
	}
	return metrics, nil
}

func ratioFloat64Ptr(numerator, denominator int64) *float64 {
	if denominator <= 0 {
		return nil
	}
	value := float64(numerator) / float64(denominator)
	return &value
}
