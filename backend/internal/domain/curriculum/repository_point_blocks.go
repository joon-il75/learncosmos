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

func researchBlockEventTitle(block CoursePointBlock) string {
	var content map[string]any
	if err := json.Unmarshal(block.Content, &content); err != nil {
		return ""
	}
	for _, key := range []string{"title", "caption", "url"} {
		if value, ok := content[key].(string); ok {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
}

func (r *Repository) CreateLearningPointBlock(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateResearchNodeBlockRequest) (uuid.UUID, CoursePointBlock, error) {
	blockType, ok := normalizeResearchNodeBlockType(req.BlockType)
	if !ok {
		return uuid.Nil, CoursePointBlock{}, errLearningPointBlockInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("begin learning point block create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	if confirmed {
		return uuid.Nil, CoursePointBlock{}, errLearningPointMaterialLocked
	}

	var existingBlockCount int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_point_blocks
		WHERE course_point_id = $1
	`, pointID).Scan(&existingBlockCount); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("count learning point blocks: %w", err)
	}
	if existingBlockCount >= maxLearningResearchBlockCount {
		return uuid.Nil, CoursePointBlock{}, errLearningPointBlockLimitExceeded
	}

	blockContent := req.Content
	if len(blockContent) == 0 {
		blockContent = json.RawMessage("{}")
	}
	blockContent, err = sanitizeResearchBlockContent(blockContent)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	if blockType == ResearchNodeBlockTypeText {
		if err := validateLearningResearchTextBlockContent(blockContent); err != nil {
			return uuid.Nil, CoursePointBlock{}, err
		}
	}

	var block CoursePointBlock
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_blocks (
			course_point_id, user_id, block_type, content, order_index
		) VALUES ($1, $2, $3, $4::jsonb, $5)
		RETURNING id, course_point_id, user_id, block_type, content, order_index, created_at, updated_at
	`, pointID, userID, blockType, string(blockContent), req.OrderIndex).Scan(
		&block.ID,
		&block.CoursePointID,
		&block.UserID,
		&block.BlockType,
		&block.Content,
		&block.OrderIndex,
		&block.CreatedAt,
		&block.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("create learning point block: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "research_material_saved", map[string]any{"block_id": block.ID, "block_type": block.BlockType, "title": researchBlockEventTitle(block)}); err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("touch learning draft after point block create: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("touch learning course after point block create: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("commit learning point block create tx: %w", err)
	}
	return courseID, block, nil
}

func (r *Repository) UpdateLearningPointBlock(ctx context.Context, userID, planetID, pointID, blockID uuid.UUID, req UpdateResearchNodeBlockRequest) (uuid.UUID, CoursePointBlock, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("begin learning point block update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	if confirmed {
		return uuid.Nil, CoursePointBlock{}, errLearningPointMaterialLocked
	}

	var block CoursePointBlock
	blockContent, err := sanitizeResearchBlockContent(req.Content)
	if err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	if err := validateLearningResearchTextBlockContent(blockContent); err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}
	if err := tx.QueryRow(ctx, `
		UPDATE course_point_blocks cpb
		SET
			content = CASE WHEN $1::text IS NOT NULL THEN $1::jsonb ELSE cpb.content END,
			order_index = COALESCE($2, cpb.order_index),
			updated_at = NOW()
		FROM course_points cp
		WHERE cpb.id = $3
		  AND cpb.course_point_id = $4
		  AND cp.id = $4
		  AND cp.course_id = $5
		  AND cp.point_type = 'research'
		RETURNING cpb.id, cpb.course_point_id, cpb.user_id, cpb.block_type, cpb.content, cpb.order_index, cpb.created_at, cpb.updated_at
	`, nullableJSON(blockContent), req.OrderIndex, blockID, pointID, courseID).Scan(
		&block.ID,
		&block.CoursePointID,
		&block.UserID,
		&block.BlockType,
		&block.Content,
		&block.OrderIndex,
		&block.CreatedAt,
		&block.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointBlock{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("update learning point block: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "research_material_saved", map[string]any{"block_id": block.ID, "block_type": block.BlockType, "action": "updated", "title": researchBlockEventTitle(block)}); err != nil {
		return uuid.Nil, CoursePointBlock{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("touch learning draft after point block update: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("touch learning course after point block update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointBlock{}, fmt.Errorf("commit learning point block update tx: %w", err)
	}
	return courseID, block, nil
}

func (r *Repository) DeleteLearningPointBlock(ctx context.Context, userID, planetID, pointID, blockID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin learning point block delete tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningResearchPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	confirmed, err := learningResearchMaterialConfirmedTx(ctx, tx, pointID)
	if err != nil {
		return uuid.Nil, err
	}
	if confirmed {
		return uuid.Nil, errLearningPointMaterialLocked
	}

	var deletedBlock CoursePointBlock
	if err := tx.QueryRow(ctx, `
		DELETE FROM course_point_blocks cpb
		USING course_points cp
		WHERE cpb.id = $1
		  AND cpb.course_point_id = $2
		  AND cp.id = $2
		  AND cp.course_id = $3
		  AND cp.point_type = 'research'
		RETURNING cpb.id, cpb.course_point_id, cpb.user_id, cpb.block_type, cpb.content, cpb.order_index, cpb.created_at, cpb.updated_at
	`, blockID, pointID, courseID).Scan(
		&deletedBlock.ID,
		&deletedBlock.CoursePointID,
		&deletedBlock.UserID,
		&deletedBlock.BlockType,
		&deletedBlock.Content,
		&deletedBlock.OrderIndex,
		&deletedBlock.CreatedAt,
		&deletedBlock.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, fmt.Errorf("delete learning point block: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "research_material_saved", map[string]any{"block_id": blockID, "block_type": deletedBlock.BlockType, "action": "deleted", "title": researchBlockEventTitle(deletedBlock)}); err != nil {
		return uuid.Nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after point block delete: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after point block delete: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit learning point block delete tx: %w", err)
	}
	return courseID, nil
}

func (r *Repository) UpsertLearningPointAISummary(ctx context.Context, userID, planetID, pointID uuid.UUID, sourceTitle, sourceDescription, summary, learningLanguage string) (uuid.UUID, PointAISummaryEntry, error) {
	sourceTitle = strings.TrimSpace(sourceTitle)
	sourceDescription = strings.TrimSpace(sourceDescription)
	summary = strings.TrimSpace(summary)
	learningLanguage = normalizeLearningLanguage(learningLanguage)
	if summary == "" {
		return uuid.Nil, PointAISummaryEntry{}, errLearningPointAISummaryInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, PointAISummaryEntry{}, fmt.Errorf("begin learning point ai summary tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, _, err := r.getLearningExplorationPointSourceTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, PointAISummaryEntry{}, err
	}

	var entry PointAISummaryEntry
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_ai_summaries (
			course_point_id, source_title, source_description, summary, learning_language
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_point_id)
		DO UPDATE SET
			source_title = EXCLUDED.source_title,
			source_description = EXCLUDED.source_description,
			summary = EXCLUDED.summary,
			learning_language = EXCLUDED.learning_language,
			updated_at = NOW()
		RETURNING source_title, source_description, summary, learning_language, updated_at
	`, pointID, sourceTitle, sourceDescription, summary, learningLanguage).Scan(
		&entry.SourceTitle,
		&entry.SourceDescription,
		&entry.Summary,
		&entry.LearningLanguage,
		&entry.UpdatedAt,
	); err != nil {
		return uuid.Nil, PointAISummaryEntry{}, fmt.Errorf("upsert learning point ai summary: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "ai_hint_used", map[string]any{"kind": "summary", "learning_language": learningLanguage}); err != nil {
		return uuid.Nil, PointAISummaryEntry{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, PointAISummaryEntry{}, fmt.Errorf("touch learning draft after point ai summary: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, PointAISummaryEntry{}, fmt.Errorf("touch learning course after point ai summary: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, PointAISummaryEntry{}, fmt.Errorf("commit learning point ai summary tx: %w", err)
	}
	return courseID, entry, nil
}
