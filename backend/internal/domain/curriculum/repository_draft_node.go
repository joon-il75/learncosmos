package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func trimmedStringOrNil(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (r *Repository) CreateMainLesson(ctx context.Context, userID, draftID uuid.UUID, req CreateMainLessonRequest) (*CourseDraftLesson, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create main lesson tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		sourceQuery        string
		learningGoal       *string
		generationLanguage string
	)
	if err := tx.QueryRow(ctx, `
		SELECT source_query, learning_goal, COALESCE(generation_language, 'ko')
		FROM course_drafts
		WHERE id = $1
		  AND user_id = $2
		  AND status = ANY($3)
	`, draftID, userID, getEditableDraftStatuses(true)).Scan(&sourceQuery, &learningGoal, &generationLanguage); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errDraftResourceNotFound
		}
		return nil, fmt.Errorf("get draft context for main lesson search spec: %w", err)
	}

	title := strings.TrimSpace(req.Title)
	objective := trimmedStringOrNil(req.Objective)
	orderIndex := req.OrderIndex
	searchSpec := BuildFallbackLessonRecommendationSearchSpecForLesson(sourceQuery, learningGoal, generationLanguage, title, objective, orderIndex)
	searchSpecJSON, err := MarshalLessonRecommendationSearchSpec(searchSpec)
	if err != nil {
		return nil, fmt.Errorf("marshal main lesson search spec: %w", err)
	}

	var lesson CourseDraftLesson
	var recommendationSearchSpecJSON []byte
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_draft_lessons
			(course_draft_id, parent_lesson_id, title, objective, summary, lesson_role, source_type, order_index, recommendation_search_spec)
		VALUES ($1, NULL, $2, $3, $3, 'core', 'manual', $4, $5::jsonb)
		RETURNING id, course_draft_id, parent_lesson_id, title, objective, summary,
		          difficulty_level, lesson_role, source_type, order_index,
		          COALESCE(recommendation_search_spec, '{}'::jsonb), created_at, updated_at
	`, draftID, title, objective, orderIndex, searchSpecJSON).Scan(
		&lesson.ID, &lesson.CourseDraftID, &lesson.ParentLessonID,
		&lesson.Title, &lesson.Objective, &lesson.Summary,
		&lesson.DifficultyLevel, &lesson.LessonRole, &lesson.SourceType,
		&lesson.OrderIndex, &recommendationSearchSpecJSON, &lesson.CreatedAt, &lesson.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create main lesson: %w", err)
	}
	lesson.RecommendationSearchSpec = ParseLessonRecommendationSearchSpec(recommendationSearchSpecJSON)
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create main lesson tx: %w", err)
	}
	return &lesson, nil
}

func (r *Repository) CreateSubLesson(ctx context.Context, userID, draftID, mainLessonID uuid.UUID, req CreateSubLessonRequest) (*CourseDraftLesson, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create sub lesson tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := fetchDraftMainLessonContext(ctx, tx, draftID, userID, mainLessonID, true); err != nil {
		if errors.Is(err, errDraftLevelNotFound) {
			return nil, errDraftResourceNotFound
		}
		return nil, err
	}

	// 현재 서브리슨 개수로 order_index 결정
	var existingCount int
	_ = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM course_draft_lessons
		WHERE course_draft_id = $1 AND parent_lesson_id = $2
	`, draftID, mainLessonID).Scan(&existingCount)

	orderIndex := req.OrderIndex
	if orderIndex == 0 {
		orderIndex = existingCount
	}

	var (
		sourceQuery        string
		learningGoal       *string
		generationLanguage string
	)
	if err := tx.QueryRow(ctx, `
		SELECT source_query, learning_goal, COALESCE(generation_language, 'ko')
		FROM course_drafts
		WHERE id = $1 AND user_id = $2
	`, draftID, userID).Scan(&sourceQuery, &learningGoal, &generationLanguage); err != nil {
		return nil, fmt.Errorf("get draft context for sub lesson search spec: %w", err)
	}

	title := strings.TrimSpace(req.Title)
	objective := trimmedStringOrNil(req.Objective)
	searchSpec := BuildFallbackLessonRecommendationSearchSpecForLesson(sourceQuery, learningGoal, generationLanguage, title, objective, orderIndex)
	searchSpecJSON, err := MarshalLessonRecommendationSearchSpec(searchSpec)
	if err != nil {
		return nil, fmt.Errorf("marshal sub lesson search spec: %w", err)
	}

	var lesson CourseDraftLesson
	var recommendationSearchSpecJSON []byte
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_draft_lessons
			(course_draft_id, parent_lesson_id, title, objective, summary, lesson_role, source_type, order_index, recommendation_search_spec)
		VALUES ($1, $2, $3, $4, $4, 'support', 'manual', $5, $6::jsonb)
		RETURNING id, course_draft_id, parent_lesson_id, title, objective, summary,
		          difficulty_level, lesson_role, source_type, order_index,
		          COALESCE(recommendation_search_spec, '{}'::jsonb), created_at, updated_at
	`, draftID, mainLessonID, title, objective, orderIndex, searchSpecJSON).Scan(
		&lesson.ID, &lesson.CourseDraftID, &lesson.ParentLessonID,
		&lesson.Title, &lesson.Objective, &lesson.Summary,
		&lesson.DifficultyLevel, &lesson.LessonRole, &lesson.SourceType,
		&lesson.OrderIndex, &recommendationSearchSpecJSON, &lesson.CreatedAt, &lesson.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create sub lesson: %w", err)
	}
	lesson.RecommendationSearchSpec = ParseLessonRecommendationSearchSpec(recommendationSearchSpecJSON)
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create sub lesson tx: %w", err)
	}
	return &lesson, nil
}

func (r *Repository) CreateResearchNode(ctx context.Context, userID, draftID, mainLessonID uuid.UUID, req CreateResearchNodeRequest) (*CourseDraftPoint, error) {
	templateType := req.TemplateType
	if templateType == "" {
		templateType = ResearchNodeTemplateFreeResearch
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create research node tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, _, err := fetchDraftMainLessonContext(ctx, tx, draftID, userID, mainLessonID, true); err != nil {
		if errors.Is(err, errDraftLevelNotFound) {
			return nil, errDraftResourceNotFound
		}
		return nil, err
	}
	var point CourseDraftPoint
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_draft_points
			(course_draft_id, course_draft_lesson_id, point_type, status, title, template_type, order_index)
		VALUES ($1, $2, 'research', 'draft', $3, $4, $5)
		RETURNING id, course_draft_id, course_draft_lesson_id, point_type, status,
		          title, description, template_type, order_index, created_at, updated_at
	`, draftID, mainLessonID, req.Title, templateType, req.OrderIndex).Scan(
		&point.ID, &point.CourseDraftID, &point.CourseDraftLessonID,
		&point.PointType, &point.Status, &point.Title, &point.Description, &point.TemplateType,
		&point.OrderIndex, &point.CreatedAt, &point.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create research node: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create research node tx: %w", err)
	}
	return &point, nil
}

func (r *Repository) UpdateResearchNode(ctx context.Context, userID, draftID, pointID uuid.UUID, req UpdateResearchNodeRequest) (*CourseDraftPoint, error) {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, pointID); err != nil {
		return nil, err
	}
	var point CourseDraftPoint
	if err := r.pool.QueryRow(ctx, `
		UPDATE course_draft_points
		SET
			title         = COALESCE($1, title),
			template_type = COALESCE($2, template_type),
			order_index   = COALESCE($3, order_index),
			updated_at    = NOW()
		WHERE id = $4 AND course_draft_id = $5 AND point_type = 'research'
		RETURNING id, course_draft_id, course_draft_lesson_id, point_type, status,
		          title, description, template_type, order_index, created_at, updated_at
	`, req.Title, req.TemplateType, req.OrderIndex, pointID, draftID).Scan(
		&point.ID, &point.CourseDraftID, &point.CourseDraftLessonID,
		&point.PointType, &point.Status, &point.Title, &point.Description, &point.TemplateType,
		&point.OrderIndex, &point.CreatedAt, &point.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errDraftResourceNotFound
		}
		return nil, fmt.Errorf("update research node: %w", err)
	}
	return &point, nil
}

func (r *Repository) DeleteResearchNode(ctx context.Context, userID, draftID, pointID uuid.UUID) error {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, pointID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM course_draft_points
		WHERE id = $1 AND course_draft_id = $2 AND point_type = 'research'
	`, pointID, draftID)
	if err != nil {
		return fmt.Errorf("delete research node: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errDraftResourceNotFound
	}
	return nil
}

func (r *Repository) GetResearchNodeBlocks(ctx context.Context, userID, draftID, pointID uuid.UUID) ([]CourseDraftPointBlock, error) {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, pointID); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.course_draft_point_id, b.block_type, b.content, b.order_index, b.created_at, b.updated_at
		FROM course_draft_point_blocks b
		JOIN course_draft_points p ON p.id = b.course_draft_point_id
		WHERE b.course_draft_point_id = $1 AND p.course_draft_id = $2
		ORDER BY b.order_index ASC, b.created_at ASC
	`, pointID, draftID)
	if err != nil {
		return nil, fmt.Errorf("get research node blocks: %w", err)
	}
	defer rows.Close()

	var blocks []CourseDraftPointBlock
	for rows.Next() {
		var b CourseDraftPointBlock
		if err := rows.Scan(&b.ID, &b.CourseDraftPointID, &b.BlockType, &b.Content, &b.OrderIndex, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan research node block: %w", err)
		}
		blocks = append(blocks, b)
	}
	if blocks == nil {
		blocks = []CourseDraftPointBlock{}
	}
	return blocks, nil
}

func (r *Repository) CreateResearchNodeBlock(ctx context.Context, userID, draftID, nodeID uuid.UUID, req CreateResearchNodeBlockRequest) (*CourseDraftPointBlock, error) {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, nodeID); err != nil {
		return nil, err
	}
	blockContent := req.Content
	if len(blockContent) == 0 {
		blockContent = json.RawMessage("{}")
	}
	var b CourseDraftPointBlock
	if err := r.pool.QueryRow(ctx, `
		INSERT INTO course_draft_point_blocks
			(course_draft_point_id, block_type, content, order_index)
		VALUES ($1, $2, $3::jsonb, $4)
		RETURNING id, course_draft_point_id, block_type, content, order_index, created_at, updated_at
	`, nodeID, req.BlockType, string(blockContent), req.OrderIndex).Scan(
		&b.ID, &b.CourseDraftPointID, &b.BlockType, &b.Content, &b.OrderIndex, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errDraftResourceNotFound
		}
		return nil, fmt.Errorf("create research node block: %w", err)
	}
	return &b, nil
}

func (r *Repository) UpdateResearchNodeBlock(ctx context.Context, userID, draftID, nodeID, blockID uuid.UUID, req UpdateResearchNodeBlockRequest) (*CourseDraftPointBlock, error) {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, nodeID); err != nil {
		return nil, err
	}
	var b CourseDraftPointBlock
	if err := r.pool.QueryRow(ctx, `
		UPDATE course_draft_point_blocks b
		SET
			content     = CASE WHEN $1::text IS NOT NULL THEN $1::jsonb ELSE b.content END,
			order_index = COALESCE($2, b.order_index),
			updated_at  = NOW()
		FROM course_draft_points p
		WHERE b.id = $3 AND b.course_draft_point_id = $4 AND p.id = $4 AND p.course_draft_id = $5
		RETURNING b.id, b.course_draft_point_id, b.block_type, b.content, b.order_index, b.created_at, b.updated_at
	`, nullableJSON(req.Content), req.OrderIndex, blockID, nodeID, draftID).Scan(
		&b.ID, &b.CourseDraftPointID, &b.BlockType, &b.Content, &b.OrderIndex, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errDraftResourceNotFound
		}
		return nil, fmt.Errorf("update research node block: %w", err)
	}
	return &b, nil
}

func (r *Repository) DeleteResearchNodeBlock(ctx context.Context, userID, draftID, nodeID, blockID uuid.UUID) error {
	if _, err := r.fetchDraftResearchPointContext(ctx, userID, draftID, nodeID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM course_draft_point_blocks b
		USING course_draft_points p
		WHERE b.id = $1 AND b.course_draft_point_id = $2 AND p.id = $2 AND p.course_draft_id = $3
	`, blockID, nodeID, draftID)
	if err != nil {
		return fmt.Errorf("delete research node block: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errDraftResourceNotFound
	}
	return nil
}

// CompleteResearchPoint sets status='completed' on a course_draft_points row.
// Phase 4 bridge: operates on course_draft_points (populated after migration 042 runs).
func (r *Repository) CompleteResearchPoint(ctx context.Context, userID, draftID, pointID uuid.UUID) error {
	if err := r.verifyDraftOwner(ctx, userID, draftID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_points
		SET status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND course_draft_id = $2
	`, pointID, draftID)
	if err != nil {
		return fmt.Errorf("complete research point: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errDraftResourceNotFound
	}
	return nil
}

// UncompleteResearchPoint resets status='draft' and clears completed_at.
func (r *Repository) UncompleteResearchPoint(ctx context.Context, userID, draftID, pointID uuid.UUID) error {
	if err := r.verifyDraftOwner(ctx, userID, draftID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_points
		SET status = 'draft', completed_at = NULL, updated_at = NOW()
		WHERE id = $1 AND course_draft_id = $2
	`, pointID, draftID)
	if err != nil {
		return fmt.Errorf("uncomplete research point: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errDraftResourceNotFound
	}
	return nil
}

// nullableJSON returns nil if the raw message is empty, otherwise returns the string pointer.
