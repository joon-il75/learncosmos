package curriculum

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) RecordRecommendationRolloutLearnerExposure(
	ctx context.Context,
	userID uuid.UUID,
	courseDraftIDText string,
	pointIDText string,
	rolloutDebug gin.H,
	candidateCount int,
	requestSnapshot string,
	candidatesSnapshot string,
) error {
	rolloutStateID := parseOptionalUUIDString(rolloutDebugString(rolloutDebug, "rollout_id"))
	courseDraftID := parseOptionalUUIDString(courseDraftIDText)
	pointID := parseOptionalUUIDString(pointIDText)
	var coursePointID *uuid.UUID
	var courseDraftPointID *uuid.UUID
	pointTargetType := ""
	if pointID != nil {
		resolvedCoursePointID, resolvedDraftPointID, resolvedTargetType, resolveErr := r.resolveRecommendationRolloutPointTarget(ctx, userID, *pointID)
		if resolveErr != nil {
			return fmt.Errorf("resolve recommendation rollout point target: %w", resolveErr)
		}
		coursePointID = resolvedCoursePointID
		courseDraftPointID = resolvedDraftPointID
		pointTargetType = resolvedTargetType
	}
	if !json.Valid([]byte(requestSnapshot)) {
		requestSnapshot = "{}"
	}
	if !json.Valid([]byte(candidatesSnapshot)) {
		candidatesSnapshot = "[]"
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO recommendation_rollout_learner_events (
			rollout_state_id,
			user_id,
			course_draft_id,
			course_point_id,
			course_draft_point_id,
			point_target_type,
			event_type,
			rollout_mode,
			traffic_percent,
			bucket,
			would_apply,
			active_applied,
			candidate_count,
			request_snapshot,
			candidates_snapshot,
			event_snapshot
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'recommendation_exposed', $7, $8, $9, $10, $11, $12, $13::jsonb, $14::jsonb, $15::jsonb)
	`, rolloutStateID, userID, courseDraftID, coursePointID, courseDraftPointID, pointTargetType,
		rolloutDebugString(rolloutDebug, "mode"),
		rolloutDebugInt(rolloutDebug, "traffic_percent"),
		rolloutDebugOptionalInt(rolloutDebug, "bucket"),
		rolloutDebugBool(rolloutDebug, "would_apply"),
		rolloutDebugBool(rolloutDebug, "active_applied"),
		candidateCount,
		requestSnapshot,
		candidatesSnapshot,
		rolloutDebugJSON(rolloutDebug),
	)
	if err != nil {
		return fmt.Errorf("record recommendation rollout learner exposure: %w", err)
	}
	return nil
}

func (r *Repository) resolveRecommendationRolloutPointTarget(ctx context.Context, userID uuid.UUID, pointID uuid.UUID) (*uuid.UUID, *uuid.UUID, string, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM course_points cp
			JOIN courses c ON c.id = cp.course_id
			WHERE cp.id = $1
			  AND c.user_id = $2
		)
	`, pointID, userID).Scan(&exists)
	if err != nil {
		return nil, nil, "", fmt.Errorf("check course point target: %w", err)
	}
	if exists {
		return &pointID, nil, "course_point", nil
	}

	err = r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM course_draft_points cdp
			JOIN course_drafts cd ON cd.id = cdp.course_draft_id
			WHERE cdp.id = $1
			  AND cd.user_id = $2
		)
	`, pointID, userID).Scan(&exists)
	if err != nil {
		return nil, nil, "", fmt.Errorf("check draft point target: %w", err)
	}
	if exists {
		return nil, &pointID, "course_draft_point", nil
	}
	return nil, nil, "", nil
}

func (r *Repository) recordRecommendationRolloutLearnerEventTx(
	ctx context.Context,
	tx pgx.Tx,
	userID uuid.UUID,
	pointID uuid.UUID,
	eventType string,
	report CoursePointMaterialReport,
) error {
	var rolloutStateID uuid.UUID
	var rolloutStateIDValue any
	var rolloutMode string
	var trafficPercent int
	var bucket sql.NullInt32
	var wouldApply bool
	var activeApplied bool

	err := tx.QueryRow(ctx, `
		SELECT
			id,
			COALESCE(rollout_mode, ''),
			traffic_percent,
			bucket,
			would_apply,
			active_applied
		FROM recommendation_rollout_learner_events
		WHERE user_id = $1
		  AND course_point_id = $2
		  AND event_type = 'recommendation_exposed'
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, pointID).Scan(&rolloutStateID, &rolloutMode, &trafficPercent, &bucket, &wouldApply, &activeApplied)
	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("load latest recommendation rollout exposure: %w", err)
	}
	if err == nil {
		rolloutStateIDValue = rolloutStateID
	} else {
		if err := tx.QueryRow(ctx, `
			SELECT id, mode, traffic_percent
			FROM recommendation_rollout_states
			WHERE active = true
			ORDER BY updated_at DESC
			LIMIT 1
		`).Scan(&rolloutStateID, &rolloutMode, &trafficPercent); err != nil && err != pgx.ErrNoRows {
			return fmt.Errorf("load active recommendation rollout state for event: %w", err)
		} else if err == nil {
			rolloutStateIDValue = rolloutStateID
		}
	}

	eventSnapshot, _ := json.Marshal(gin.H{
		"report_id":              report.ID,
		"target_type":            report.TargetType,
		"target_content_id":      report.TargetContentID,
		"target_url":             report.TargetURL,
		"replacement_content_id": report.ReplacementContentID,
		"replacement_url":        report.ReplacementURL,
		"replacement_title":      report.ReplacementTitle,
	})

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_rollout_learner_events (
			rollout_state_id,
			user_id,
			course_point_id,
			event_type,
			rollout_mode,
			traffic_percent,
			bucket,
			would_apply,
			active_applied,
			selected_content_id,
			selected_url,
			report_id,
			report_type,
			event_snapshot
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb)
	`, rolloutStateIDValue, userID, pointID, eventType, rolloutMode, trafficPercent, nullInt32ToAny(bucket), wouldApply, activeApplied,
		report.ReplacementContentID, report.ReplacementURL, report.ID, report.ReportType, string(eventSnapshot))
	if err != nil {
		return fmt.Errorf("record recommendation rollout learner event: %w", err)
	}
	return nil
}

func nullInt32ToAny(value sql.NullInt32) any {
	if !value.Valid {
		return nil
	}
	return int(value.Int32)
}

func parseOptionalUUIDString(value string) *uuid.UUID {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil
	}
	return &parsed
}

func rolloutDebugString(debug gin.H, key string) string {
	if value, ok := debug[key].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

func rolloutDebugBool(debug gin.H, key string) bool {
	if value, ok := debug[key].(bool); ok {
		return value
	}
	return false
}

func rolloutDebugInt(debug gin.H, key string) int {
	switch value := debug[key].(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		return 0
	}
}

func rolloutDebugOptionalInt(debug gin.H, key string) *int {
	if _, ok := debug[key]; !ok {
		return nil
	}
	value := rolloutDebugInt(debug, key)
	return &value
}

func rolloutDebugJSON(debug gin.H) string {
	raw, err := json.Marshal(debug)
	if err != nil || !json.Valid(raw) {
		return "{}"
	}
	return string(raw)
}
