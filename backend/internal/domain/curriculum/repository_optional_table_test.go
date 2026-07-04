package curriculum

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUndefinedTableError(t *testing.T) {
	err := &pgconn.PgError{
		Code:    "42P01",
		Message: `relation "course_point_question_ai_feedbacks" does not exist`,
	}

	if !isUndefinedTableError(err, "course_point_question_ai_feedbacks") {
		t.Fatal("expected undefined table error to be detected")
	}

	if isUndefinedTableError(err, "course_points") {
		t.Fatal("expected table name mismatch to return false")
	}
}
