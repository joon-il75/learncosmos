package llmjobs

import (
	"testing"
	"time"
)

func TestEvaluateStuckJobRecovery(t *testing.T) {
	now := time.Date(2026, 5, 25, 12, 0, 0, 0, time.UTC)
	policy := RecoveryPolicy{
		QueuedTimeout:  10 * time.Minute,
		RunningTimeout: 15 * time.Minute,
	}

	tests := []struct {
		name string
		job  *Job
		want RecoveryDecision
	}{
		{
			name: "nil job",
			job:  nil,
			want: RecoveryDecision{Action: RecoveryActionNone},
		},
		{
			name: "terminal status ignored",
			job: &Job{
				Status:    StatusSucceeded,
				CreatedAt: now.Add(-time.Hour),
				UpdatedAt: now.Add(-time.Hour),
			},
			want: RecoveryDecision{Action: RecoveryActionNone},
		},
		{
			name: "fresh queued ignored",
			job: &Job{
				Status:    StatusQueued,
				CreatedAt: now.Add(-2 * time.Minute),
				QueuedAt:  timePtr(now.Add(-2 * time.Minute)),
				UpdatedAt: now.Add(-2 * time.Minute),
			},
			want: RecoveryDecision{Action: RecoveryActionNone},
		},
		{
			name: "stale queued becomes retryable failed candidate",
			job: &Job{
				Status:    StatusQueued,
				CreatedAt: now.Add(-11 * time.Minute),
				QueuedAt:  timePtr(now.Add(-11 * time.Minute)),
				UpdatedAt: now.Add(-11 * time.Minute),
			},
			want: RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckQueued,
				ErrorMessage: "llm job remained queued beyond recovery timeout",
				Retryable:    true,
			},
		},
		{
			name: "stale queued falls back to created at",
			job: &Job{
				Status:    StatusQueued,
				CreatedAt: now.Add(-11 * time.Minute),
				UpdatedAt: now.Add(-11 * time.Minute),
			},
			want: RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckQueued,
				ErrorMessage: "llm job remained queued beyond recovery timeout",
				Retryable:    true,
			},
		},
		{
			name: "fresh running ignored",
			job: &Job{
				Status:    StatusRunning,
				CreatedAt: now.Add(-20 * time.Minute),
				StartedAt: timePtr(now.Add(-5 * time.Minute)),
				UpdatedAt: now.Add(-5 * time.Minute),
			},
			want: RecoveryDecision{Action: RecoveryActionNone},
		},
		{
			name: "stale running becomes retryable failed candidate",
			job: &Job{
				Status:    StatusRunning,
				CreatedAt: now.Add(-30 * time.Minute),
				StartedAt: timePtr(now.Add(-16 * time.Minute)),
				UpdatedAt: now.Add(-16 * time.Minute),
			},
			want: RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckRunning,
				ErrorMessage: "llm job remained running beyond recovery timeout",
				Retryable:    true,
			},
		},
		{
			name: "stale running falls back to updated at",
			job: &Job{
				Status:    StatusRunning,
				CreatedAt: now.Add(-30 * time.Minute),
				UpdatedAt: now.Add(-16 * time.Minute),
			},
			want: RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckRunning,
				ErrorMessage: "llm job remained running beyond recovery timeout",
				Retryable:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateStuckJobRecovery(tt.job, now, policy)
			if got != tt.want {
				t.Fatalf("EvaluateStuckJobRecovery() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestEvaluateStuckJobRecoveryDisabledTimeouts(t *testing.T) {
	now := time.Date(2026, 5, 25, 12, 0, 0, 0, time.UTC)
	job := &Job{
		Status:    StatusRunning,
		CreatedAt: now.Add(-time.Hour),
		StartedAt: timePtr(now.Add(-time.Hour)),
		UpdatedAt: now.Add(-time.Hour),
	}
	got := EvaluateStuckJobRecovery(job, now, RecoveryPolicy{})
	if got.Action != RecoveryActionNone {
		t.Fatalf("action = %s, want %s", got.Action, RecoveryActionNone)
	}
}

func timePtr(value time.Time) *time.Time {
	return &value
}
