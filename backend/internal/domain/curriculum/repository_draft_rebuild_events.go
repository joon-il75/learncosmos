package curriculum

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func insertRebuildEventTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID, draftID uuid.UUID, sourceQuery string, payload map[string]any, label string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s event payload: %w", label, err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`, userID, draftID, EventCurriculumEdited, sourceQuery, string(body)); err != nil {
		return fmt.Errorf("insert %s event: %w", label, err)
	}
	return nil
}
