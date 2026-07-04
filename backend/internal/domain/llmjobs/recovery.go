package llmjobs

import "time"

const (
	RecoveryActionNone          = "none"
	RecoveryActionFailRetryable = "fail_retryable"

	RecoveryErrorCodeStuckQueued  = "llm_job_stuck_queued"
	RecoveryErrorCodeStuckRunning = "llm_job_stuck_running"
)

type RecoveryPolicy struct {
	QueuedTimeout  time.Duration
	RunningTimeout time.Duration
}

type RecoveryDecision struct {
	Action       string
	ErrorCode    string
	ErrorMessage string
	Retryable    bool
}

func EvaluateStuckJobRecovery(job *Job, now time.Time, policy RecoveryPolicy) RecoveryDecision {
	if job == nil || terminalStatus(job.Status) {
		return RecoveryDecision{Action: RecoveryActionNone}
	}
	if now.IsZero() {
		now = time.Now()
	}
	switch job.Status {
	case StatusQueued:
		if policy.QueuedTimeout <= 0 {
			return RecoveryDecision{Action: RecoveryActionNone}
		}
		startedAt := job.QueuedAt
		if startedAt == nil {
			startedAt = &job.CreatedAt
		}
		if isOlderThan(now, *startedAt, policy.QueuedTimeout) {
			return RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckQueued,
				ErrorMessage: "llm job remained queued beyond recovery timeout",
				Retryable:    true,
			}
		}
	case StatusRunning:
		if policy.RunningTimeout <= 0 {
			return RecoveryDecision{Action: RecoveryActionNone}
		}
		startedAt := job.StartedAt
		if startedAt == nil {
			startedAt = &job.UpdatedAt
		}
		if isOlderThan(now, *startedAt, policy.RunningTimeout) {
			return RecoveryDecision{
				Action:       RecoveryActionFailRetryable,
				ErrorCode:    RecoveryErrorCodeStuckRunning,
				ErrorMessage: "llm job remained running beyond recovery timeout",
				Retryable:    true,
			}
		}
	}
	return RecoveryDecision{Action: RecoveryActionNone}
}

func isOlderThan(now, startedAt time.Time, timeout time.Duration) bool {
	if startedAt.IsZero() || timeout <= 0 {
		return false
	}
	return !startedAt.After(now) && now.Sub(startedAt) >= timeout
}
