package curriculum

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func mapDraftStatusToPlanetStatus(status DraftStatus) PlanetStatus {
	switch status {
	case DraftStatusConfirmed:
		return PlanetStatusReady
	case DraftStatusLearning:
		return PlanetStatusLearning
	case DraftStatusArchived:
		return PlanetStatusCompleted
	default:
		return PlanetStatusDraft
	}
}

func containsDraftStatus(statuses []DraftStatus, target DraftStatus) bool {
	for _, status := range statuses {
		if status == target {
			return true
		}
	}
	return false
}

func buildPlanetListItemProgress(status DraftStatus, regionCount, completedRegionCount int) *float64 {
	if status != DraftStatusLearning || regionCount <= 0 {
		return nil
	}
	progress := float64(completedRegionCount) / float64(regionCount)
	return &progress
}

type planetProgressSummary struct {
	RegionCount          int
	CompletedRegionCount int
	PointCount           int
	CompletedPointCount  int
	Progress             *float64
	CanComplete          bool
}

type planetProgressQueryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func buildPlanetProgressSummary(status DraftStatus, regionCount, completedRegionCount, pointCount, completedPointCount int) planetProgressSummary {
	return planetProgressSummary{
		RegionCount:          regionCount,
		CompletedRegionCount: completedRegionCount,
		PointCount:           pointCount,
		CompletedPointCount:  completedPointCount,
		Progress:             buildPlanetListItemProgress(status, regionCount, completedRegionCount),
		CanComplete:          regionCount > 0 && completedRegionCount == regionCount,
	}
}

func (r *Repository) loadExplorerProgressSummary(ctx context.Context, q planetProgressQueryer, draftID uuid.UUID, status DraftStatus, confirmedCourseID *uuid.UUID) (*planetProgressSummary, error) {
	if confirmedCourseID != nil {
		return r.loadExplorerCourseProgressSummary(ctx, q, draftID, *confirmedCourseID, status)
	}
	return r.loadExplorerDraftProgressSummary(ctx, q, draftID, status)
}

func (r *Repository) loadExplorerCourseProgressSummary(ctx context.Context, q planetProgressQueryer, draftID, courseID uuid.UUID, status DraftStatus) (*planetProgressSummary, error) {
	var regionCount, completedRegionCount, pointCount, completedPointCount int
	err := q.QueryRow(ctx, `
		WITH active_regions AS (
			SELECT r.id
			FROM explorer_regions r
			WHERE r.course_draft_id = $1
			  AND r.status = 'active'
		),
		region_nodes AS (
			SELECT
				r.id AS region_id,
				n.id AS node_id,
				cp.status AS point_status
			FROM active_regions r
			LEFT JOIN explorer_nodes n
			  ON n.parent_kind = 'region'
			 AND n.parent_id = r.id
			 AND n.status = 'active'
			LEFT JOIN course_points cp
			  ON cp.id = n.draft_point_id
			 AND cp.course_id = $2
			UNION ALL
			SELECT
				r.id AS region_id,
				n.id AS node_id,
				cp.status AS point_status
			FROM explorer_subregions sr
			JOIN active_regions r ON r.id = sr.region_id
			LEFT JOIN explorer_nodes n
			  ON n.parent_kind = 'subregion'
			 AND n.parent_id = sr.id
			 AND n.status = 'active'
			LEFT JOIN course_points cp
			  ON cp.id = n.draft_point_id
			 AND cp.course_id = $2
			WHERE sr.status = 'active'
		),
		region_points AS (
			SELECT
				region_id,
				COUNT(node_id) AS point_count,
				COUNT(node_id) FILTER (WHERE point_status = 'completed') AS completed_count
			FROM region_nodes
			GROUP BY region_id
		)
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE point_count > 0 AND point_count = completed_count),
			COALESCE(SUM(point_count), 0),
			COALESCE(SUM(completed_count), 0)
		FROM region_points
	`, draftID, courseID).Scan(&regionCount, &completedRegionCount, &pointCount, &completedPointCount)
	if err != nil {
		return nil, err
	}
	if regionCount == 0 {
		return nil, nil
	}
	summary := buildPlanetProgressSummary(status, regionCount, completedRegionCount, pointCount, completedPointCount)
	return &summary, nil
}

func (r *Repository) loadExplorerDraftProgressSummary(ctx context.Context, q planetProgressQueryer, draftID uuid.UUID, status DraftStatus) (*planetProgressSummary, error) {
	var regionCount, completedRegionCount, pointCount, completedPointCount int
	err := q.QueryRow(ctx, `
		WITH active_regions AS (
			SELECT r.id
			FROM explorer_regions r
			WHERE r.course_draft_id = $1
			  AND r.status = 'active'
		),
		region_nodes AS (
			SELECT
				r.id AS region_id,
				n.id AS node_id,
				cdp.status AS point_status
			FROM active_regions r
			LEFT JOIN explorer_nodes n
			  ON n.parent_kind = 'region'
			 AND n.parent_id = r.id
			 AND n.status = 'active'
			LEFT JOIN course_draft_points cdp
			  ON cdp.id = n.draft_point_id
			 AND cdp.course_draft_id = $1
			UNION ALL
			SELECT
				r.id AS region_id,
				n.id AS node_id,
				cdp.status AS point_status
			FROM explorer_subregions sr
			JOIN active_regions r ON r.id = sr.region_id
			LEFT JOIN explorer_nodes n
			  ON n.parent_kind = 'subregion'
			 AND n.parent_id = sr.id
			 AND n.status = 'active'
			LEFT JOIN course_draft_points cdp
			  ON cdp.id = n.draft_point_id
			 AND cdp.course_draft_id = $1
			WHERE sr.status = 'active'
		),
		region_points AS (
			SELECT
				region_id,
				COUNT(node_id) AS point_count,
				COUNT(node_id) FILTER (WHERE point_status = 'completed') AS completed_count
			FROM region_nodes
			GROUP BY region_id
		)
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE point_count > 0 AND point_count = completed_count),
			COALESCE(SUM(point_count), 0),
			COALESCE(SUM(completed_count), 0)
		FROM region_points
	`, draftID).Scan(&regionCount, &completedRegionCount, &pointCount, &completedPointCount)
	if err != nil {
		return nil, err
	}
	if regionCount == 0 {
		return nil, nil
	}
	summary := buildPlanetProgressSummary(status, regionCount, completedRegionCount, pointCount, completedPointCount)
	return &summary, nil
}

func applyPlanetProgressSummary(item *PlanetListItem, summary planetProgressSummary) {
	item.LessonCount = summary.RegionCount
	item.CompletedLessons = summary.CompletedRegionCount
	item.Progress = summary.Progress
	item.CanComplete = summary.CanComplete
}

func (r *Repository) ListPlanetsByDraftStatuses(ctx context.Context, userID uuid.UUID, statuses []DraftStatus) ([]PlanetListItem, error) {
	drafts, err := r.listDraftsByStatuses(ctx, userID, statuses)
	if err != nil {
		return nil, err
	}

	planets := make([]PlanetListItem, 0, len(drafts))
	for _, draft := range drafts {
		planet, err := r.buildPlanetAggregateFromDraftState(ctx, userID, draft)
		if err != nil {
			return nil, err
		}
		planets = append(planets, planet.Planet)
	}
	return planets, nil
}
