package llmjobs

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusCanceled  Status = "canceled"
	StatusExpired   Status = "expired"
)

const (
	FeatureDummy                                   = "dummy_echo"
	FeatureCourseDraftCreate                       = "course_draft_create"
	FeatureGoalInterviewTurn                       = "goal_interview_turn"
	FeaturePointAISummary                          = "point_ai_summary"
	FeaturePointQuestionGenerate                   = "point_question_generate"
	FeaturePointFeedbackGenerate                   = "point_feedback_generate"
	FeaturePointSelfEvaluationDraft                = "point_self_evaluation_draft"
	FeatureAdminRecommendationDebugGoalTurn        = "admin_recommendation_debug_goal_turn"
	FeatureAdminRecommendationDebugLessonsGenerate = "admin_recommendation_debug_lessons_generate"
)

type ProcessorError struct {
	Code      string
	Message   string
	Retryable bool
}

func (e *ProcessorError) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.Message) != "" {
		return e.Message
	}
	return e.Code
}

func NewProcessorError(code, message string, retryable bool) *ProcessorError {
	return &ProcessorError{Code: strings.TrimSpace(code), Message: strings.TrimSpace(message), Retryable: retryable}
}

var (
	ErrInvalidJobInput        = errors.New("invalid llm job input")
	ErrInvalidTransition      = errors.New("invalid llm job status transition")
	ErrJobNotFound            = errors.New("llm job not found")
	ErrJobAccessDenied        = errors.New("llm job access denied")
	ErrUnsupportedFeature     = errors.New("unsupported llm job feature")
	ErrWorkerQueueDisabled    = errors.New("llm worker queue disabled")
	ErrDuplicateJob           = errors.New("duplicate llm job")
	ErrActiveJobLimitExceeded = errors.New("llm active job limit exceeded")
)

type Job struct {
	ID                     uuid.UUID       `json:"id"`
	UserID                 *uuid.UUID      `json:"user_id,omitempty"`
	Feature                string          `json:"feature"`
	IdempotencyKey         *string         `json:"idempotency_key,omitempty"`
	Status                 Status          `json:"status"`
	Priority               int             `json:"priority"`
	Phase                  *string         `json:"phase,omitempty"`
	RequestRef             json.RawMessage `json:"request_ref"`
	PromptTemplateKey      *string         `json:"prompt_template_key,omitempty"`
	PromptInputRef         json.RawMessage `json:"prompt_input_ref"`
	Provider               *string         `json:"provider,omitempty"`
	Model                  *string         `json:"model,omitempty"`
	ProviderMode           *string         `json:"provider_mode,omitempty"`
	APIKeyRef              *string         `json:"api_key_ref,omitempty"`
	ResultRef              json.RawMessage `json:"result_ref"`
	ErrorCode              *string         `json:"error_code,omitempty"`
	ErrorMessage           *string         `json:"error_message,omitempty"`
	Retryable              bool            `json:"retryable"`
	PointCost              *int            `json:"point_cost,omitempty"`
	BillingStatus          *string         `json:"billing_status,omitempty"`
	QueueWaitMS            *int            `json:"queue_wait_ms,omitempty"`
	ProviderWaitMS         *int            `json:"provider_wait_ms,omitempty"`
	LLMGenerationElapsedMS *int            `json:"llm_generation_elapsed_ms,omitempty"`
	ParseValidateElapsedMS *int            `json:"parse_validate_elapsed_ms,omitempty"`
	PersistElapsedMS       *int            `json:"persist_elapsed_ms,omitempty"`
	TotalJobElapsedMS      *int            `json:"total_job_elapsed_ms,omitempty"`
	ProviderInFlightCount  *int            `json:"provider_inflight_count,omitempty"`
	ProviderQueueDepth     *int            `json:"provider_queue_depth,omitempty"`
	WorkerID               *string         `json:"worker_id,omitempty"`
	Attempts               int             `json:"attempts"`
	CreatedAt              time.Time       `json:"created_at"`
	QueuedAt               *time.Time      `json:"queued_at,omitempty"`
	StartedAt              *time.Time      `json:"started_at,omitempty"`
	FinishedAt             *time.Time      `json:"finished_at,omitempty"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

type ActiveJobEstimate struct {
	ActiveCount   int
	QueuePosition *int
}

type CreateJobInput struct {
	UserID            *uuid.UUID
	Feature           string
	IdempotencyKey    *string
	Priority          int
	RequestRef        map[string]any
	PromptTemplateKey *string
	PromptInputRef    map[string]any
	Provider          *string
	Model             *string
	ProviderMode      *string
	APIKeyRef         *string
	PointCost         *int
	BillingStatus     *string
	MaxActiveJobs     int
}

type CompleteJobInput struct {
	ResultRef              map[string]any
	ProviderWaitMS         *int
	LLMGenerationElapsedMS *int
	ParseValidateElapsedMS *int
	PersistElapsedMS       *int
	TotalJobElapsedMS      *int
}

type FailJobInput struct {
	ErrorCode              string
	ErrorMessage           string
	Retryable              bool
	ProviderWaitMS         *int
	LLMGenerationElapsedMS *int
	ParseValidateElapsedMS *int
	PersistElapsedMS       *int
	TotalJobElapsedMS      *int
}

func normalizeFeature(value string) string {
	return strings.TrimSpace(value)
}

func ValidateStatusTransition(from, to Status) bool {
	switch from {
	case StatusQueued:
		return to == StatusRunning || to == StatusCanceled || to == StatusExpired || to == StatusFailed
	case StatusRunning:
		return to == StatusSucceeded || to == StatusFailed || to == StatusCanceled || to == StatusExpired
	case StatusSucceeded, StatusFailed, StatusCanceled, StatusExpired:
		return false
	default:
		return false
	}
}

func terminalStatus(status Status) bool {
	return status == StatusSucceeded || status == StatusFailed || status == StatusCanceled || status == StatusExpired
}
