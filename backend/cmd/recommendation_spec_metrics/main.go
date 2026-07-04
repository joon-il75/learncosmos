package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/pkg/db"
)

const (
	defaultWindowDays = 7
	shortQueryRunes   = 8
)

type metricsSummary struct {
	WindowDays              int                 `json:"window_days"`
	WindowStartedAt         time.Time           `json:"window_started_at"`
	WindowEndedAt           time.Time           `json:"window_ended_at"`
	RecommendationSearches  int64               `json:"recommendation_searches"`
	SearchSpecSourceCounts  map[string]int64    `json:"search_spec_source_counts"`
	ResolutionSourceCounts  map[string]int64    `json:"resolution_source_counts"`
	EmptyEffectiveQueryRate *float64            `json:"empty_effective_query_rate,omitempty"`
	ShortEffectiveQueryRate *float64            `json:"short_effective_query_rate,omitempty"`
	TopEffectiveQueries     []effectiveQueryRow `json:"top_effective_queries"`
	RolloutEvents           rolloutEventMetrics `json:"rollout_events"`
}

type effectiveQueryRow struct {
	Query string `json:"query"`
	Count int64  `json:"count"`
}

type rolloutEventMetrics struct {
	RecommendationExposed int64    `json:"recommendation_exposed"`
	MaterialReplaced      int64    `json:"material_replaced"`
	BrokenLinkReported    int64    `json:"broken_link_reported"`
	WrongContentReported  int64    `json:"wrong_content_reported"`
	AverageCandidateCount *float64 `json:"average_candidate_count,omitempty"`
}

func main() {
	var windowDays int
	flag.IntVar(&windowDays, "days", defaultWindowDays, "metrics window in days")
	flag.Parse()

	if windowDays <= 0 {
		log.Fatal("-days must be greater than 0")
	}

	loadEnv()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	windowEndedAt := time.Now().UTC()
	windowStartedAt := windowEndedAt.AddDate(0, 0, -windowDays)

	summary, err := collectRecommendationSpecMetrics(ctx, pool, windowStartedAt, windowEndedAt, windowDays)
	if err != nil {
		log.Fatalf("failed to collect recommendation spec metrics: %v", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		log.Fatalf("failed to print metrics: %v", err)
	}
}

func loadEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}

func collectRecommendationSpecMetrics(ctx context.Context, pool *pgxpool.Pool, startedAt, endedAt time.Time, windowDays int) (metricsSummary, error) {
	searches, emptyQueries, shortQueries, err := countRecommendationSearches(ctx, pool, startedAt, endedAt)
	if err != nil {
		return metricsSummary{}, err
	}

	specSources, err := countMetadataValue(ctx, pool, startedAt, endedAt, "search_spec_source")
	if err != nil {
		return metricsSummary{}, err
	}
	resolutionSources, err := countMetadataValue(ctx, pool, startedAt, endedAt, "resolution_source")
	if err != nil {
		return metricsSummary{}, err
	}
	topQueries, err := listTopEffectiveQueries(ctx, pool, startedAt, endedAt, 10)
	if err != nil {
		return metricsSummary{}, err
	}
	rolloutEvents, err := collectRolloutEventMetrics(ctx, pool, startedAt, endedAt)
	if err != nil {
		return metricsSummary{}, err
	}

	return metricsSummary{
		WindowDays:              windowDays,
		WindowStartedAt:         startedAt,
		WindowEndedAt:           endedAt,
		RecommendationSearches:  searches,
		SearchSpecSourceCounts:  specSources,
		ResolutionSourceCounts:  resolutionSources,
		EmptyEffectiveQueryRate: ratioPtr(emptyQueries, searches),
		ShortEffectiveQueryRate: ratioPtr(shortQueries, searches),
		TopEffectiveQueries:     topQueries,
		RolloutEvents:           rolloutEvents,
	}, nil
}

func countRecommendationSearches(ctx context.Context, pool *pgxpool.Pool, startedAt, endedAt time.Time) (int64, int64, int64, error) {
	var total int64
	var emptyQuery int64
	var shortQuery int64
	err := pool.QueryRow(ctx, `
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
	`, startedAt, endedAt, shortQueryRunes).Scan(&total, &emptyQuery, &shortQuery)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("count recommendation searches: %w", err)
	}
	return total, emptyQuery, shortQuery, nil
}

func countMetadataValue(ctx context.Context, pool *pgxpool.Pool, startedAt, endedAt time.Time, key string) (map[string]int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT COALESCE(NULLIF(btrim(metadata->>$3), ''), 'empty') AS value, COUNT(*) AS count
		FROM ai_point_transactions
		WHERE feature = 'explorer_recommendation_search'
		  AND created_at >= $1
		  AND created_at < $2
		GROUP BY value
		ORDER BY count DESC, value ASC
	`, startedAt, endedAt, key)
	if err != nil {
		return nil, fmt.Errorf("count metadata value %s: %w", key, err)
	}
	defer rows.Close()

	counts := map[string]int64{}
	for rows.Next() {
		var value string
		var count int64
		if err := rows.Scan(&value, &count); err != nil {
			return nil, fmt.Errorf("scan metadata value %s: %w", key, err)
		}
		counts[value] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metadata value %s: %w", key, err)
	}
	return counts, nil
}

func listTopEffectiveQueries(ctx context.Context, pool *pgxpool.Pool, startedAt, endedAt time.Time, limit int) ([]effectiveQueryRow, error) {
	rows, err := pool.Query(ctx, `
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
		return nil, fmt.Errorf("list top effective queries: %w", err)
	}
	defer rows.Close()

	queries := []effectiveQueryRow{}
	for rows.Next() {
		var item effectiveQueryRow
		if err := rows.Scan(&item.Query, &item.Count); err != nil {
			return nil, fmt.Errorf("scan effective query: %w", err)
		}
		queries = append(queries, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate effective queries: %w", err)
	}
	return queries, nil
}

func collectRolloutEventMetrics(ctx context.Context, pool *pgxpool.Pool, startedAt, endedAt time.Time) (rolloutEventMetrics, error) {
	var metrics rolloutEventMetrics
	var averageCandidateCount *float64
	err := pool.QueryRow(ctx, `
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
		return rolloutEventMetrics{}, fmt.Errorf("collect rollout event metrics: %w", err)
	}
	metrics.AverageCandidateCount = averageCandidateCount
	return metrics, nil
}

func ratioPtr(numerator, denominator int64) *float64 {
	if denominator <= 0 {
		return nil
	}
	value := float64(numerator) / float64(denominator)
	return &value
}

func sortedMetricKeys(values map[string]int64) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
