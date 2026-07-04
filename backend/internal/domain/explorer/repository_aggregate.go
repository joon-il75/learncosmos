package explorer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetCourseOwner: course_drafts.user_id 반환 (소유권 검증용)
func (r *Repository) GetCourseOwner(ctx context.Context, courseDraftID uuid.UUID) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT user_id FROM course_drafts WHERE id = $1`, courseDraftID,
	).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotFound
	}
	return userID, err
}

func (r *Repository) GetCourseStatus(ctx context.Context, courseDraftID uuid.UUID) (string, error) {
	var status string
	err := r.pool.QueryRow(ctx,
		`SELECT status FROM course_drafts WHERE id = $1`, courseDraftID,
	).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return status, err
}

// GetCourseAggregate: 전체 집계 조회 (지역 → 서브지역 → 지점 중첩 구조)
// includeInactive=true 시 비활성 항목 포함
func (r *Repository) GetCourseAggregate(ctx context.Context, courseDraftID uuid.UUID, includeInactive bool) (*CourseAggregate, error) {
	var title string
	var status string
	var isInactive bool
	var updatedAt time.Time
	var confirmedCourseID *uuid.UUID
	var textureMapID *uuid.UUID
	var textureMapName *string
	var textureMapAsset *string
	var textureMapRotationDurationSeconds *int
	var textureMapRotationDirection *string
	err := r.pool.QueryRow(ctx,
		`SELECT cd.title,
		        cd.status,
		        COALESCE(
		            c.status = 'archived',
		            cd.status = 'archived'
		            AND cd.confirmed_course_id IS NULL
		            AND NOT EXISTS (
		                SELECT 1
		                FROM recommendation_events re
		                WHERE re.course_draft_id = cd.id
		                  AND re.event_type = 'curriculum_completed'
		            )
		        ),
		        cd.updated_at,
		        cd.confirmed_course_id,
		        COALESCE(c.planet_texture_map_id, cd.planet_texture_map_id),
		        ptmp.name,
		        ptmp.asset_path,
		        ptmp.rotation_duration_seconds,
		        ptmp.rotation_direction
		 FROM course_drafts cd
		 LEFT JOIN courses c ON c.id = cd.confirmed_course_id AND c.user_id = cd.user_id
		 LEFT JOIN planet_texture_maps ptmp ON ptmp.id = COALESCE(c.planet_texture_map_id, cd.planet_texture_map_id)
		 WHERE cd.id = $1`, courseDraftID,
	).Scan(
		&title,
		&status,
		&isInactive,
		&updatedAt,
		&confirmedCourseID,
		&textureMapID,
		&textureMapName,
		&textureMapAsset,
		&textureMapRotationDurationSeconds,
		&textureMapRotationDirection,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	regions, err := r.queryRegions(ctx, courseDraftID, includeInactive)
	if err != nil {
		return nil, err
	}
	if len(regions) == 0 {
		return &CourseAggregate{
			CourseDraftID:                           courseDraftID,
			Title:                                   title,
			Status:                                  status,
			IsInactive:                              isInactive,
			Progress:                                float64Pointer(0),
			PlanetTextureMapID:                      textureMapID,
			PlanetTextureMapName:                    textureMapName,
			PlanetTextureMapAsset:                   textureMapAsset,
			PlanetTextureMapRotationDurationSeconds: textureMapRotationDurationSeconds,
			PlanetTextureMapRotationDirection:       textureMapRotationDirection,
			UpdatedAt:                               updatedAt,
			Regions:                                 []RegionAggregate{},
		}, nil
	}

	regionIDs := make([]uuid.UUID, len(regions))
	for i, reg := range regions {
		regionIDs[i] = reg.ID
	}

	subregions, err := r.querySubRegionsByRegions(ctx, regionIDs, includeInactive)
	if err != nil {
		return nil, err
	}

	subregionIDs := make([]uuid.UUID, len(subregions))
	for i, sr := range subregions {
		subregionIDs[i] = sr.ID
	}

	nodes, err := r.queryNodesByParents(ctx, regionIDs, subregionIDs, includeInactive, confirmedCourseID)
	if err != nil {
		return nil, err
	}

	agg := assembleCourseAggregate(
		courseDraftID,
		title,
		status,
		isInactive,
		updatedAt,
		computeExplorerProgress(regions, subregions, nodes),
		textureMapID,
		textureMapName,
		textureMapAsset,
		textureMapRotationDurationSeconds,
		textureMapRotationDirection,
		regions,
		subregions,
		nodes,
	)
	return &agg, nil
}

func (r *Repository) queryRegions(ctx context.Context, courseDraftID uuid.UUID, includeInactive bool) ([]Region, error) {
	q := `SELECT id, course_draft_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	      FROM explorer_regions
	      WHERE course_draft_id = $1`
	if !includeInactive {
		q += ` AND status = 'active'`
	}
	q += ` ORDER BY order_index ASC, created_at ASC`

	rows, err := r.pool.Query(ctx, q, courseDraftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Region, 0)
	for rows.Next() {
		var reg Region
		if err := rows.Scan(
			&reg.ID, &reg.CourseDraftID, &reg.CourseDraftLessonID, &reg.Name, &reg.Description,
			&reg.OrderIndex, &reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, reg)
	}
	return result, rows.Err()
}

func (r *Repository) querySubRegionsByRegions(ctx context.Context, regionIDs []uuid.UUID, includeInactive bool) ([]SubRegion, error) {
	if len(regionIDs) == 0 {
		return []SubRegion{}, nil
	}

	q := `SELECT id, region_id, course_draft_lesson_id, name, description, order_index, status, created_at, updated_at
	      FROM explorer_subregions
	      WHERE region_id = ANY($1)`
	if !includeInactive {
		q += ` AND status = 'active'`
	}
	q += ` ORDER BY order_index ASC, created_at ASC`

	rows, err := r.pool.Query(ctx, q, regionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]SubRegion, 0)
	for rows.Next() {
		var sr SubRegion
		if err := rows.Scan(
			&sr.ID, &sr.RegionID, &sr.CourseDraftLessonID, &sr.Name, &sr.Description,
			&sr.OrderIndex, &sr.Status, &sr.CreatedAt, &sr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, sr)
	}
	return result, rows.Err()
}

func (r *Repository) queryNodesByParents(ctx context.Context, regionIDs, subregionIDs []uuid.UUID, includeInactive bool, courseID *uuid.UUID) ([]Node, error) {
	if len(regionIDs) == 0 && len(subregionIDs) == 0 {
		return []Node{}, nil
	}

	statusCond := ""
	if !includeInactive {
		statusCond = " AND n.status = 'active'"
	}
	q := fmt.Sprintf(`
		SELECT id, parent_kind, parent_id, node_type, draft_point_id, learning_status,
		       title, order_index, status, created_at, updated_at,
		       source_type, source_url, content_id, summary,
		       research_type, layout_type, block_count
		FROM (
			SELECT n.id, n.parent_kind, n.parent_id, n.node_type, n.draft_point_id, cp.status AS learning_status,
			       n.title, n.order_index, n.status, n.created_at, n.updated_at,
			       n.source_type, n.source_url, n.content_id, n.summary,
			       n.research_type, n.layout_type, n.block_count
			FROM explorer_nodes n
			LEFT JOIN course_points cp
			  ON cp.id = n.draft_point_id
			 AND cp.course_id = $3
			WHERE (
				(n.parent_kind = 'region' AND n.parent_id = ANY($1))
				OR
				(n.parent_kind = 'subregion' AND n.parent_id = ANY($2))
			)%s
		) explorer_nodes
		ORDER BY order_index ASC, created_at ASC
	`, statusCond)

	rows, err := r.pool.Query(ctx, q, regionIDs, subregionIDs, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNodesWithLearningStatus(rows)
}

func assembleCourseAggregate(
	courseDraftID uuid.UUID, title string, status string, isInactive bool, updatedAt time.Time,
	progress *float64,
	textureMapID *uuid.UUID,
	textureMapName *string,
	textureMapAsset *string,
	textureMapRotationDurationSeconds *int,
	textureMapRotationDirection *string,
	regions []Region, subregions []SubRegion, nodes []Node,
) CourseAggregate {
	subByRegion := make(map[uuid.UUID][]SubRegion, len(subregions))
	for _, sr := range subregions {
		subByRegion[sr.RegionID] = append(subByRegion[sr.RegionID], sr)
	}

	nodesByParent := make(map[uuid.UUID][]Node, len(nodes))
	for _, n := range nodes {
		nodesByParent[n.ParentID] = append(nodesByParent[n.ParentID], n)
	}

	sortNodes := func(ns []Node) {
		sort.Slice(ns, func(i, j int) bool {
			if ns[i].OrderIndex == ns[j].OrderIndex {
				return ns[i].CreatedAt.Before(ns[j].CreatedAt)
			}
			return ns[i].OrderIndex < ns[j].OrderIndex
		})
	}

	regionAggregates := make([]RegionAggregate, 0, len(regions))
	for _, reg := range regions {
		subs := subByRegion[reg.ID]
		sort.Slice(subs, func(i, j int) bool {
			if subs[i].OrderIndex == subs[j].OrderIndex {
				return subs[i].CreatedAt.Before(subs[j].CreatedAt)
			}
			return subs[i].OrderIndex < subs[j].OrderIndex
		})

		subAggregates := make([]SubRegionAggregate, 0, len(subs))
		for _, sr := range subs {
			ns := nodesByParent[sr.ID]
			if ns == nil {
				ns = []Node{}
			}
			sortNodes(ns)
			subAggregates = append(subAggregates, SubRegionAggregate{
				SubRegion: sr,
				Nodes:     ns,
			})
		}

		regionNodes := nodesByParent[reg.ID]
		if regionNodes == nil {
			regionNodes = []Node{}
		}
		sortNodes(regionNodes)

		regionAggregates = append(regionAggregates, RegionAggregate{
			Region:     reg,
			SubRegions: subAggregates,
			Nodes:      regionNodes,
		})
	}

	return CourseAggregate{
		CourseDraftID:                           courseDraftID,
		Title:                                   title,
		Status:                                  status,
		IsInactive:                              isInactive,
		Progress:                                progress,
		PlanetTextureMapID:                      textureMapID,
		PlanetTextureMapName:                    textureMapName,
		PlanetTextureMapAsset:                   textureMapAsset,
		PlanetTextureMapRotationDurationSeconds: textureMapRotationDurationSeconds,
		PlanetTextureMapRotationDirection:       textureMapRotationDirection,
		UpdatedAt:                               updatedAt,
		Regions:                                 regionAggregates,
	}
}

func computeExplorerProgress(regions []Region, subregions []SubRegion, nodes []Node) *float64 {
	activeRegionIDs := make(map[uuid.UUID]struct{}, len(regions))
	for _, region := range regions {
		if region.Status != ItemStatusActive {
			continue
		}
		activeRegionIDs[region.ID] = struct{}{}
	}
	if len(activeRegionIDs) == 0 {
		return float64Pointer(0)
	}

	activeSubRegionParent := make(map[uuid.UUID]uuid.UUID, len(subregions))
	for _, subregion := range subregions {
		if subregion.Status != ItemStatusActive {
			continue
		}
		if _, ok := activeRegionIDs[subregion.RegionID]; !ok {
			continue
		}
		activeSubRegionParent[subregion.ID] = subregion.RegionID
	}

	nodesByRegion := make(map[uuid.UUID][]Node, len(activeRegionIDs))
	for _, node := range nodes {
		if node.Status != ItemStatusActive {
			continue
		}
		switch node.ParentKind {
		case ParentKindRegion:
			if _, ok := activeRegionIDs[node.ParentID]; ok {
				nodesByRegion[node.ParentID] = append(nodesByRegion[node.ParentID], node)
			}
		case ParentKindSubRegion:
			if regionID, ok := activeSubRegionParent[node.ParentID]; ok {
				nodesByRegion[regionID] = append(nodesByRegion[regionID], node)
			}
		}
	}

	totalRegions := 0
	completedRegions := 0
	for _, region := range regions {
		if region.Status != ItemStatusActive {
			continue
		}
		totalRegions++
		regionNodes := nodesByRegion[region.ID]
		if len(regionNodes) == 0 {
			continue
		}
		allCompleted := true
		for _, node := range regionNodes {
			if node.LearningStatus == nil || *node.LearningStatus != "completed" {
				allCompleted = false
				break
			}
		}
		if allCompleted {
			completedRegions++
		}
	}
	if totalRegions == 0 {
		return float64Pointer(0)
	}
	progress := float64(completedRegions) / float64(totalRegions)
	return &progress
}

func float64Pointer(value float64) *float64 {
	return &value
}
