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

func (r *Repository) ReplaceLessonCandidateResources(ctx context.Context, draftID, lessonID uuid.UUID, candidates []ContentSearchCandidate, sourceQuery string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin candidate tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		DELETE FROM course_draft_points
		WHERE course_draft_lesson_id = $1
		  AND point_type = 'exploration'
		  AND selection_state = 'candidate'
	`, lessonID)
	if err != nil {
		return fmt.Errorf("delete candidate resources: %w", err)
	}

	for idx, candidate := range candidates {
		resourceID := uuid.New()
		_, err := tx.Exec(ctx, `
			INSERT INTO course_draft_points (
				id, course_draft_id, course_draft_lesson_id, point_type, selection_state,
				content_id, external_url, title, description, thumbnail_url,
				price_type, rank_score, status, order_index
			)
			SELECT $1, cdl.course_draft_id, $2, 'exploration', 'candidate',
			       $3, $4, $5, $6, $7, $8, $9, 'draft', $10
			FROM course_draft_lessons cdl
			WHERE cdl.id = $2
		`,
			resourceID,
			lessonID,
			candidate.ContentID,
			candidate.ExternalURL,
			candidate.Title,
			candidate.Description,
			candidate.ThumbnailURL,
			candidate.PriceType,
			candidate.RankScore,
			idx,
		)
		if err != nil {
			return fmt.Errorf("insert candidate resource: %w", err)
		}

		payload, marshalErr := json.Marshal(map[string]any{
			"title":      candidate.Title,
			"rank_score": candidate.RankScore,
		})
		if marshalErr != nil {
			return fmt.Errorf("marshal lesson recommended payload: %w", marshalErr)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO recommendation_events (
				user_id, course_draft_id, course_draft_lesson_id, content_id, event_type, source_query, payload
			)
			SELECT cd.user_id, $1, $2, $3, $4, $5, $6::jsonb
			FROM course_drafts cd
			WHERE cd.id = $1
		`,
			draftID,
			lessonID,
			candidate.ContentID,
			EventLessonRecommended,
			sourceQuery,
			string(payload),
		)
		if err != nil {
			return fmt.Errorf("insert lesson recommended event: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit candidate tx: %w", err)
	}
	return nil
}

func (r *Repository) SelectDraftLessonResource(ctx context.Context, draftID, userID, resourceID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin select resource tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var lessonID uuid.UUID
	var contentID *uuid.UUID
	var selectionState ResourceSelectionState
	var title string
	var sourceQuery string
	err = tx.QueryRow(ctx, `
		SELECT p.course_draft_lesson_id, p.content_id, p.selection_state,
		       p.title, cd.source_query
		FROM course_draft_points p
		JOIN course_draft_lessons cdl ON cdl.id = p.course_draft_lesson_id
		JOIN course_drafts cd ON cd.id = p.course_draft_id
		WHERE p.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		  AND p.point_type = 'exploration'
		FOR UPDATE OF p
	`, resourceID, draftID, userID, getEditableDraftStatuses(false)).Scan(&lessonID, &contentID, &selectionState, &title, &sourceQuery)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftResourceNotFound
		}
		return fmt.Errorf("get draft resource for selection: %w", err)
	}

	if selectionState == ResourceSelectionSelected {
		return tx.Commit(ctx)
	}
	if selectionState != ResourceSelectionCandidate {
		return errDraftResourceInvalidState
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_draft_points
		SET selection_state = 'selected', updated_at = NOW()
		WHERE id = $1
	`, resourceID); err != nil {
		return fmt.Errorf("update draft resource selection: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("touch draft after resource selection: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"resource_id": resourceID.String(),
		"title":       title,
	})
	if err != nil {
		return fmt.Errorf("marshal lesson selected payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, content_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
	`, userID, draftID, lessonID, contentID, EventLessonSelected, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert lesson selected event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit select resource tx: %w", err)
	}
	return nil
}

func (r *Repository) AttachDraftLessonResource(ctx context.Context, draftID, userID, lessonID uuid.UUID, req AttachDraftLessonResourceRequest) error {
	selectionState := req.SelectionState
	if selectionState == "" {
		selectionState = ResourceSelectionSelected
	}
	if selectionState != ResourceSelectionCandidate && selectionState != ResourceSelectionSelected {
		return errDraftResourceInvalidState
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin attach resource tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sourceQuery string
	err = tx.QueryRow(ctx, `
		SELECT cd.source_query
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(false)).Scan(&sourceQuery)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftLessonNotFound
		}
		return fmt.Errorf("get draft lesson for resource attach: %w", err)
	}

	var title string
	var description *string
	var thumbnailURL *string
	var externalURL *string
	err = tx.QueryRow(ctx, `
		SELECT title, description, thumbnail_url, url
		FROM contents
		WHERE id = $1
		  AND (user_id = $2 OR is_public = true)
	`, req.ContentID, userID).Scan(&title, &description, &thumbnailURL, &externalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errDraftContentNotFound
		}
		return fmt.Errorf("get content for resource attach: %w", err)
	}

	var duplicateExists bool
	if externalURL != nil && strings.TrimSpace(*externalURL) != "" {
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM course_draft_points
				WHERE course_draft_lesson_id = $1
				  AND point_type = 'exploration'
				  AND (
				    content_id = $2
				    OR external_url = $3
				  )
			)
		`, lessonID, req.ContentID, strings.TrimSpace(*externalURL)).Scan(&duplicateExists)
	} else {
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM course_draft_points
				WHERE course_draft_lesson_id = $1
				  AND point_type = 'exploration'
				  AND content_id = $2
			)
		`, lessonID, req.ContentID).Scan(&duplicateExists)
	}
	if err != nil {
		return fmt.Errorf("check duplicate draft resource: %w", err)
	}
	if duplicateExists {
		return errDraftResourceDuplicate
	}

	var orderIndex int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(order_index), -1) + 1
		FROM course_draft_points
		WHERE course_draft_lesson_id = $1
		  AND point_type = 'exploration'
	`, lessonID).Scan(&orderIndex); err != nil {
		return fmt.Errorf("next draft resource order: %w", err)
	}

	resourceID := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO course_draft_points (
			id, course_draft_id, course_draft_lesson_id, point_type, selection_state,
			content_id, external_url, title, description, thumbnail_url,
			price_type, status, order_index
		)
		SELECT $1, cdl.course_draft_id, $2, 'exploration', $3,
		       $4, $5, $6, $7, $8, 'owned', 'draft', $9
		FROM course_draft_lessons cdl
		WHERE cdl.id = $2
	`, resourceID, lessonID, selectionState, req.ContentID, externalURL, title, description, thumbnailURL, orderIndex)
	if err != nil {
		return fmt.Errorf("insert draft lesson resource: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE course_drafts
		SET updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, draftID, userID); err != nil {
		return fmt.Errorf("touch draft after resource attach: %w", err)
	}

	payload, err := json.Marshal(map[string]any{
		"resource_id":     resourceID.String(),
		"content_id":      req.ContentID.String(),
		"title":           title,
		"selection_state": string(selectionState),
	})
	if err != nil {
		return fmt.Errorf("marshal lesson added manually payload: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, course_draft_lesson_id, content_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
	`, userID, draftID, lessonID, req.ContentID, EventLessonAddedManually, sourceQuery, string(payload)); err != nil {
		return fmt.Errorf("insert lesson added manually event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit attach resource tx: %w", err)
	}
	return nil
}
