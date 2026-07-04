package llmjobs

import (
	"context"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultObservationQueuedTimeout  = 10 * time.Minute
	DefaultObservationRunningTimeout = 30 * time.Minute
)

type JobObservationCount struct {
	Feature        string     `json:"feature"`
	Status         Status     `json:"status"`
	TotalCount     int        `json:"total_count"`
	StaleCount     int        `json:"stale_count"`
	OldestCreated  *time.Time `json:"oldest_created_at,omitempty"`
	OldestActivity *time.Time `json:"oldest_activity_at,omitempty"`
}

type observationCutoffs struct {
	QueuedCutoff  *time.Time
	RunningCutoff *time.Time
}

func buildObservationCutoffs(now time.Time, policy RecoveryPolicy) observationCutoffs {
	if now.IsZero() {
		now = time.Now()
	}
	var cutoffs observationCutoffs
	if policy.QueuedTimeout > 0 {
		queued := now.Add(-policy.QueuedTimeout)
		cutoffs.QueuedCutoff = &queued
	}
	if policy.RunningTimeout > 0 {
		running := now.Add(-policy.RunningTimeout)
		cutoffs.RunningCutoff = &running
	}
	return cutoffs
}

func (r *Repository) CountJobObservations(ctx context.Context, now time.Time, policy RecoveryPolicy) ([]JobObservationCount, error) {
	cutoffs := buildObservationCutoffs(now, policy)
	rows, err := r.pool.Query(ctx, `
		SELECT
			feature,
			status,
			count(*)::int AS total_count,
			count(*) FILTER (
				WHERE (status = $1 AND $3::timestamptz IS NOT NULL AND COALESCE(queued_at, created_at) <= $3)
				   OR (status = $2 AND $4::timestamptz IS NOT NULL AND COALESCE(started_at, updated_at) <= $4)
			)::int AS stale_count,
			min(created_at) AS oldest_created_at,
			min(COALESCE(started_at, queued_at, updated_at, created_at)) AS oldest_activity_at
		FROM llm_jobs
		WHERE status IN ($1, $2, $5, $6, $7, $8)
		GROUP BY feature, status
		ORDER BY feature, status
	`, StatusQueued, StatusRunning, cutoffs.QueuedCutoff, cutoffs.RunningCutoff, StatusSucceeded, StatusFailed, StatusCanceled, StatusExpired)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := []JobObservationCount{}
	for rows.Next() {
		var item JobObservationCount
		if err := rows.Scan(&item.Feature, &item.Status, &item.TotalCount, &item.StaleCount, &item.OldestCreated, &item.OldestActivity); err != nil {
			return nil, err
		}
		counts = append(counts, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return counts, nil
}

func DefaultObservationPolicy() RecoveryPolicy {
	return RecoveryPolicy{
		QueuedTimeout:  DefaultObservationQueuedTimeout,
		RunningTimeout: DefaultObservationRunningTimeout,
	}
}

func parseObservationDuration(raw string, fallback time.Duration) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	if duration, err := time.ParseDuration(raw); err == nil && duration >= 0 {
		return duration
	}
	if ms, err := strconv.Atoi(raw); err == nil && ms >= 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return fallback
}
