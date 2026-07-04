package explorer

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) InsertNode(ctx context.Context, node Node) (*Node, error) {
	rows, err := r.pool.Query(ctx, `
		INSERT INTO explorer_nodes
			(parent_kind, parent_id, node_type, draft_point_id, title, order_index,
			 source_type, source_url, content_id, summary,
			 research_type, layout_type, block_count)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		          source_type, source_url, content_id, summary,
		          research_type, layout_type, block_count
	`,
		node.ParentKind, node.ParentID, node.NodeType, node.DraftPointID, node.Title, node.OrderIndex,
		node.SourceType, node.SourceURL, node.ContentID, node.Summary,
		node.ResearchType, node.LayoutType, node.BlockCount,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	n := nodes[0]
	return &n, nil
}

func (r *Repository) GetNodeByID(ctx context.Context, nodeID uuid.UUID) (*Node, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		       source_type, source_url, content_id, summary,
		       research_type, layout_type, block_count
		FROM explorer_nodes WHERE id=$1
	`, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	n := nodes[0]
	return &n, nil
}

func (r *Repository) UpdateNode(ctx context.Context, nodeID uuid.UUID, req UpdateNodeRequest) (*Node, error) {
	rows, err := r.pool.Query(ctx, `
		UPDATE explorer_nodes SET
			parent_kind   = COALESCE($2::text, parent_kind),
			parent_id     = COALESCE($3, parent_id),
			title         = COALESCE($4, title),
			source_type   = COALESCE($5::text, source_type),
			source_url    = COALESCE($6, source_url),
			content_id    = COALESCE($7, content_id),
			summary       = COALESCE($8, summary),
			research_type = COALESCE($9::text, research_type),
			layout_type   = COALESCE($10::text, layout_type),
			order_index   = COALESCE($11, order_index),
			updated_at    = now()
		WHERE id=$1
		RETURNING id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		          source_type, source_url, content_id, summary,
		          research_type, layout_type, block_count
	`, nodeID,
		req.ParentKind, req.ParentID, req.Title, req.SourceType, req.SourceURL, req.ContentID, req.Summary,
		req.ResearchType, req.LayoutType, req.OrderIndex,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	n := nodes[0]
	return &n, nil
}

type draftPointLinkTarget struct {
	DraftID           uuid.UUID
	DraftStatus       string
	ConfirmedCourseID *uuid.UUID
}

func (r *Repository) loadDraftPointLinkTargetTx(ctx context.Context, tx pgx.Tx, courseDraftID uuid.UUID) (*draftPointLinkTarget, error) {
	var target draftPointLinkTarget
	if err := tx.QueryRow(ctx, `
		SELECT id, status, confirmed_course_id
		FROM course_drafts
		WHERE id = $1
	`, courseDraftID).Scan(&target.DraftID, &target.DraftStatus, &target.ConfirmedCourseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &target, nil
}

func supportsDraftPointLink(node Node) bool {
	return (node.NodeType == NodeTypeExploration || node.NodeType == NodeTypeResearch) &&
		(node.ParentKind == ParentKindRegion || node.ParentKind == ParentKindSubRegion)
}

func trimmedOrNil(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (r *Repository) resolveDraftLessonIDForNodeParentTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID, parentKind ParentKind, parentID uuid.UUID) (*uuid.UUID, error) {
	var lessonID uuid.UUID
	switch parentKind {
	case ParentKindRegion:
		err := tx.QueryRow(ctx, `
			SELECT course_draft_lesson_id
			FROM explorer_regions
			WHERE id = $1
			  AND course_draft_id = $2
			  AND course_draft_lesson_id IS NOT NULL
		`, parentID, draftID).Scan(&lessonID)
		if err == nil {
			return &lessonID, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	case ParentKindSubRegion:
		err := tx.QueryRow(ctx, `
			SELECT es.course_draft_lesson_id
			FROM explorer_subregions es
			JOIN explorer_regions er ON er.id = es.region_id
			WHERE es.id = $1
			  AND er.course_draft_id = $2
			  AND es.course_draft_lesson_id IS NOT NULL
		`, parentID, draftID).Scan(&lessonID)
		if err == nil {
			return &lessonID, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	exists, err := r.draftLessonExistsTx(ctx, tx, draftID, parentID)
	if err != nil || !exists {
		return nil, err
	}
	legacyID := parentID
	return &legacyID, nil
}

func (r *Repository) draftLessonExistsTx(ctx context.Context, tx pgx.Tx, draftID, lessonID uuid.UUID) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM course_draft_lessons
			WHERE id = $1 AND course_draft_id = $2
		)
	`, lessonID, draftID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) courseLessonExistsTx(ctx context.Context, tx pgx.Tx, courseID, lessonID uuid.UUID) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM course_lessons
			WHERE id = $1 AND course_id = $2
		)
	`, lessonID, courseID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) resolveRuntimeCourseLessonIDTx(ctx context.Context, tx pgx.Tx, courseID, draftLessonID uuid.UUID) (*uuid.UUID, error) {
	var (
		title          string
		orderIndex     int
		parentLessonID *uuid.UUID
	)
	if err := tx.QueryRow(ctx, `
		SELECT title, order_index, parent_lesson_id
		FROM course_draft_lessons
		WHERE id = $1
	`, draftLessonID).Scan(&title, &orderIndex, &parentLessonID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var runtimeParentID *uuid.UUID
	if parentLessonID != nil {
		resolvedParentID, err := r.resolveRuntimeCourseLessonIDTx(ctx, tx, courseID, *parentLessonID)
		if err != nil {
			return nil, err
		}
		if resolvedParentID == nil {
			return nil, nil
		}
		runtimeParentID = resolvedParentID
	}

	var runtimeLessonID uuid.UUID
	err := tx.QueryRow(ctx, `
		SELECT id
		FROM course_lessons
		WHERE course_id = $1
		  AND title = $2
		  AND order_index = $3
		  AND (
			($4::uuid IS NULL AND parent_lesson_id IS NULL)
			OR parent_lesson_id = $4
		  )
		ORDER BY created_at ASC
		LIMIT 1
	`, courseID, title, orderIndex, runtimeParentID).Scan(&runtimeLessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &runtimeLessonID, nil
}

func (r *Repository) nextDraftPointOrderIndexTx(ctx context.Context, tx pgx.Tx, draftID, lessonID uuid.UUID) (int, error) {
	var nextIndex int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(order_index) + 1, 0)
		FROM course_draft_points
		WHERE course_draft_id = $1
		  AND course_draft_lesson_id = $2
	`, draftID, lessonID).Scan(&nextIndex); err != nil {
		return 0, err
	}
	return nextIndex, nil
}

func (r *Repository) nextCoursePointOrderIndexTx(ctx context.Context, tx pgx.Tx, courseID, lessonID uuid.UUID) (int, error) {
	var nextIndex int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(order_index) + 1, 0)
		FROM course_points
		WHERE course_id = $1
		  AND course_lesson_id = $2
	`, courseID, lessonID).Scan(&nextIndex); err != nil {
		return 0, err
	}
	return nextIndex, nil
}

func (r *Repository) createLinkedPointForNodeTx(ctx context.Context, tx pgx.Tx, target *draftPointLinkTarget, node Node) (*uuid.UUID, error) {
	if target == nil || !supportsDraftPointLink(node) {
		return nil, nil
	}

	parentLessonID, err := r.resolveDraftLessonIDForNodeParentTx(ctx, tx, target.DraftID, node.ParentKind, node.ParentID)
	if err != nil || parentLessonID == nil {
		return nil, err
	}

	pointID := uuid.New()
	pointType := string(node.NodeType)
	orderIndex, err := r.nextDraftPointOrderIndexTx(ctx, tx, target.DraftID, *parentLessonID)
	if err != nil {
		return nil, err
	}

	switch node.NodeType {
	case NodeTypeExploration:
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_points (
				id, course_draft_id, course_draft_lesson_id, point_type, status,
				title, content_id, external_url, order_index
			) VALUES ($1, $2, $3, $4, 'draft', $5, $6, $7, $8)
		`, pointID, target.DraftID, *parentLessonID, pointType, node.Title, node.ContentID, trimmedOrNil(node.SourceURL), orderIndex); err != nil {
			return nil, err
		}
	case NodeTypeResearch:
		if _, err := tx.Exec(ctx, `
			INSERT INTO course_draft_points (
				id, course_draft_id, course_draft_lesson_id, point_type, status,
				title, template_type, order_index
			) VALUES ($1, $2, $3, $4, 'draft', $5, 'free_research', $6)
		`, pointID, target.DraftID, *parentLessonID, pointType, node.Title, orderIndex); err != nil {
			return nil, err
		}
	}

	if target.DraftStatus == "learning" && target.ConfirmedCourseID != nil {
		runtimeLessonID, err := r.resolveRuntimeCourseLessonIDTx(ctx, tx, *target.ConfirmedCourseID, *parentLessonID)
		if err != nil {
			return nil, err
		}
		if runtimeLessonID != nil {
			courseOrderIndex, err := r.nextCoursePointOrderIndexTx(ctx, tx, *target.ConfirmedCourseID, *runtimeLessonID)
			if err != nil {
				return nil, err
			}
			switch node.NodeType {
			case NodeTypeExploration:
				if _, err := tx.Exec(ctx, `
					INSERT INTO course_points (
						id, course_id, course_lesson_id, point_type, status,
						title, content_id, external_url, order_index
					) VALUES ($1, $2, $3, $4, 'ready', $5, $6, $7, $8)
				`, pointID, *target.ConfirmedCourseID, *runtimeLessonID, pointType, node.Title, node.ContentID, trimmedOrNil(node.SourceURL), courseOrderIndex); err != nil {
					return nil, err
				}
			case NodeTypeResearch:
				if _, err := tx.Exec(ctx, `
					INSERT INTO course_points (
						id, course_id, course_lesson_id, point_type, status,
						title, template_type, order_index
					) VALUES ($1, $2, $3, $4, 'ready', $5, 'free_research', $6)
				`, pointID, *target.ConfirmedCourseID, *runtimeLessonID, pointType, node.Title, courseOrderIndex); err != nil {
					return nil, err
				}
			}
		}
	}

	return &pointID, nil
}

func (r *Repository) syncLinkedPointForNodeTx(ctx context.Context, tx pgx.Tx, target *draftPointLinkTarget, node Node) error {
	if target == nil || node.DraftPointID == nil || !supportsDraftPointLink(node) {
		return nil
	}

	parentLessonID, err := r.resolveDraftLessonIDForNodeParentTx(ctx, tx, target.DraftID, node.ParentKind, node.ParentID)
	if err != nil || parentLessonID == nil {
		return err
	}

	switch node.NodeType {
	case NodeTypeExploration:
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_points
			SET course_draft_lesson_id = $1,
			    title = $2,
			    content_id = $3,
			    external_url = $4,
			    updated_at = NOW()
			WHERE id = $5
			  AND course_draft_id = $6
		`, *parentLessonID, node.Title, node.ContentID, trimmedOrNil(node.SourceURL), *node.DraftPointID, target.DraftID); err != nil {
			return err
		}
	case NodeTypeResearch:
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_points
			SET course_draft_lesson_id = $1,
			    title = $2,
			    updated_at = NOW()
			WHERE id = $3
			  AND course_draft_id = $4
		`, *parentLessonID, node.Title, *node.DraftPointID, target.DraftID); err != nil {
			return err
		}
	}

	if target.DraftStatus == "learning" && target.ConfirmedCourseID != nil {
		runtimeLessonID, err := r.resolveRuntimeCourseLessonIDTx(ctx, tx, *target.ConfirmedCourseID, *parentLessonID)
		if err != nil {
			return err
		}
		if runtimeLessonID == nil {
			return nil
		}
		switch node.NodeType {
		case NodeTypeExploration:
			if _, err := tx.Exec(ctx, `
				UPDATE course_points
				SET course_lesson_id = $1,
				    title = $2,
				    content_id = $3,
				    external_url = $4,
				    updated_at = NOW()
				WHERE id = $5
				  AND course_id = $6
			`, *runtimeLessonID, node.Title, node.ContentID, trimmedOrNil(node.SourceURL), *node.DraftPointID, *target.ConfirmedCourseID); err != nil {
				return err
			}
		case NodeTypeResearch:
			if _, err := tx.Exec(ctx, `
				UPDATE course_points
				SET course_lesson_id = $1,
				    title = $2,
				    updated_at = NOW()
				WHERE id = $3
				  AND course_id = $4
			`, *runtimeLessonID, node.Title, *node.DraftPointID, *target.ConfirmedCourseID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Repository) syncLinkedPointContentTx(ctx context.Context, tx pgx.Tx, target *draftPointLinkTarget, node Node) error {
	if target == nil || node.DraftPointID == nil || !supportsDraftPointLink(node) {
		return nil
	}

	switch node.NodeType {
	case NodeTypeExploration:
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_points
			SET content_id = $1,
			    external_url = $2,
			    updated_at = NOW()
			WHERE id = $3
			  AND course_draft_id = $4
		`, node.ContentID, trimmedOrNil(node.SourceURL), *node.DraftPointID, target.DraftID); err != nil {
			return err
		}
	case NodeTypeResearch:
		if _, err := tx.Exec(ctx, `
			UPDATE course_draft_points
			SET content_id = $1,
			    updated_at = NOW()
			WHERE id = $2
			  AND course_draft_id = $3
		`, node.ContentID, *node.DraftPointID, target.DraftID); err != nil {
			return err
		}
	}

	if target.DraftStatus == "learning" && target.ConfirmedCourseID != nil {
		switch node.NodeType {
		case NodeTypeExploration:
			if _, err := tx.Exec(ctx, `
				UPDATE course_points
				SET content_id = $1,
				    external_url = $2,
				    updated_at = NOW()
				WHERE id = $3
				  AND course_id = $4
			`, node.ContentID, trimmedOrNil(node.SourceURL), *node.DraftPointID, *target.ConfirmedCourseID); err != nil {
				return err
			}
		case NodeTypeResearch:
			if _, err := tx.Exec(ctx, `
				UPDATE course_points
				SET content_id = $1,
				    updated_at = NOW()
				WHERE id = $2
				  AND course_id = $3
			`, node.ContentID, *node.DraftPointID, *target.ConfirmedCourseID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Repository) InsertNodeWithLinkedPoint(ctx context.Context, courseDraftID uuid.UUID, node Node) (*Node, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	target, err := r.loadDraftPointLinkTargetTx(ctx, tx, courseDraftID)
	if err != nil {
		return nil, err
	}
	if node.DraftPointID == nil {
		pointID, err := r.createLinkedPointForNodeTx(ctx, tx, target, node)
		if err != nil {
			return nil, err
		}
		node.DraftPointID = pointID
	}

	rows, err := tx.Query(ctx, `
		INSERT INTO explorer_nodes
			(parent_kind, parent_id, node_type, draft_point_id, title, order_index,
			 source_type, source_url, content_id, summary,
			 research_type, layout_type, block_count)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		          source_type, source_url, content_id, summary,
		          research_type, layout_type, block_count
	`,
		node.ParentKind, node.ParentID, node.NodeType, node.DraftPointID, node.Title, node.OrderIndex,
		node.SourceType, node.SourceURL, node.ContentID, node.Summary,
		node.ResearchType, node.LayoutType, node.BlockCount,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	n := nodes[0]
	return &n, nil
}

func (r *Repository) UpdateNodeWithLinkedPoint(ctx context.Context, courseDraftID, nodeID uuid.UUID, req UpdateNodeRequest) (*Node, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	target, err := r.loadDraftPointLinkTargetTx(ctx, tx, courseDraftID)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		UPDATE explorer_nodes SET
			parent_kind   = COALESCE($2::text, parent_kind),
			parent_id     = COALESCE($3, parent_id),
			title         = COALESCE($4, title),
			source_type   = COALESCE($5::text, source_type),
			source_url    = COALESCE($6, source_url),
			content_id    = COALESCE($7, content_id),
			summary       = COALESCE($8, summary),
			research_type = COALESCE($9::text, research_type),
			layout_type   = COALESCE($10::text, layout_type),
			order_index   = COALESCE($11, order_index),
			updated_at    = now()
		WHERE id=$1
		RETURNING id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		          source_type, source_url, content_id, summary,
		          research_type, layout_type, block_count
	`, nodeID,
		req.ParentKind, req.ParentID, req.Title, req.SourceType, req.SourceURL, req.ContentID, req.Summary,
		req.ResearchType, req.LayoutType, req.OrderIndex,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	updated := nodes[0]

	if err := r.syncLinkedPointForNodeTx(ctx, tx, target, updated); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *Repository) DeleteNodeAndLinkedDraftPoint(ctx context.Context, courseDraftID, nodeID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		       source_type, source_url, content_id, summary,
		       research_type, layout_type, block_count
		FROM explorer_nodes
		WHERE id = $1
	`, nodeID)
	if err != nil {
		return err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return ErrNotFound
	}
	node := nodes[0]

	if node.DraftPointID != nil {
		if _, err := tx.Exec(ctx, `
			DELETE FROM course_draft_points
			WHERE id = $1
			  AND course_draft_id = $2
		`, *node.DraftPointID, courseDraftID); err != nil {
			return err
		}
	}

	tag, err := tx.Exec(ctx, `DELETE FROM explorer_nodes WHERE id=$1`, nodeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return tx.Commit(ctx)
}

func (r *Repository) SyncNodeContentWithLinkedPoint(ctx context.Context, courseDraftID, nodeID uuid.UUID, contentID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	target, err := r.loadDraftPointLinkTargetTx(ctx, tx, courseDraftID)
	if err != nil {
		return err
	}

	rows, err := tx.Query(ctx, `
		UPDATE explorer_nodes
		SET content_id = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, parent_kind, parent_id, node_type, draft_point_id, title, order_index, status, created_at, updated_at,
		          source_type, source_url, content_id, summary,
		          research_type, layout_type, block_count
	`, nodeID, contentID)
	if err != nil {
		return err
	}
	defer rows.Close()

	nodes, err := scanNodes(rows)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return ErrNotFound
	}

	if err := r.syncLinkedPointContentTx(ctx, tx, target, nodes[0]); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) UpdateNodeStatus(ctx context.Context, nodeID uuid.UUID, status ItemStatus) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE explorer_nodes SET status=$1, updated_at=now() WHERE id=$2`,
		status, nodeID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteNode(ctx context.Context, nodeID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM explorer_nodes WHERE id=$1`,
		nodeID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
