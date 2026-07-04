package llmjobs

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsDuplicateJobError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "llm job idempotency unique violation",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "idx_llm_jobs_idempotency_key",
			},
			want: true,
		},
		{
			name: "other unique violation",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "other_constraint",
			},
			want: false,
		},
		{
			name: "generic error",
			err:  errors.New("boom"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDuplicateJobError(tt.err); got != tt.want {
				t.Fatalf("isDuplicateJobError() = %v, want %v", got, tt.want)
			}
		})
	}
}
