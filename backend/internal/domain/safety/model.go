package safety

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	ActionAllow    Action = "allow"
	ActionSoftWarn Action = "soft_warn"
	ActionBlock    Action = "block"
)

const (
	ErrorCodeBlocked  = "safety_input_blocked"
	ErrorCodeSoftWarn = "safety_input_soft_warn"
)

const (
	TargetAIGoalInterviewOutput = "ai_goal_interview_output"
	TargetAICourseDraftOutput   = "ai_course_draft_output"
	TargetAIPointSummaryOutput  = "ai_point_summary_output"
	TargetAIPointQuestionOutput = "ai_point_question_output"
	TargetAIPointFeedbackOutput = "ai_point_feedback_output"
	TargetAIPointSelfEvalOutput = "ai_point_self_evaluation_output"
)

var AIOutputTargetTypes = []string{
	TargetAIGoalInterviewOutput,
	TargetAICourseDraftOutput,
	TargetAIPointSummaryOutput,
	TargetAIPointQuestionOutput,
	TargetAIPointFeedbackOutput,
	TargetAIPointSelfEvalOutput,
}

type Rule struct {
	ID          *uuid.UUID
	RuleType    string
	Pattern     string
	RiskType    string
	Action      Action
	Locale      *string
	Description string
}

type ModerateInput struct {
	UserID     *uuid.UUID
	TargetType string
	TargetID   *uuid.UUID
	Text       string
	Locale     string
	Route      string
	Metadata   map[string]any
}

type ModerationResult struct {
	Action        Action     `json:"action"`
	RiskType      string     `json:"risk_type"`
	ErrorCode     *string    `json:"error_code,omitempty"`
	Message       *string    `json:"message,omitempty"`
	MatchedRuleID *uuid.UUID `json:"matched_rule_id,omitempty"`
}

type LogEntry struct {
	ID             uuid.UUID      `json:"id"`
	UserID         *uuid.UUID     `json:"user_id,omitempty"`
	UserEmail      *string        `json:"user_email,omitempty"`
	TargetType     string         `json:"target_type"`
	TargetID       *uuid.UUID     `json:"target_id,omitempty"`
	InputTextHash  string         `json:"input_text_hash"`
	RiskType       string         `json:"risk_type"`
	Action         Action         `json:"action"`
	MatchedRuleID  *uuid.UUID     `json:"matched_rule_id,omitempty"`
	Route          string         `json:"route"`
	Locale         *string        `json:"locale,omitempty"`
	Metadata       map[string]any `json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	MatchedPattern *string        `json:"matched_pattern,omitempty"`
	RuleType       *string        `json:"rule_type,omitempty"`
}

type ListLogsInput struct {
	Action     string
	RiskType   string
	TargetType string
	Limit      int
	Offset     int
}

type AIOutputSummary struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Windows     []AIOutputSummaryWindow `json:"windows"`
}

type AIOutputSummaryWindow struct {
	Label          string                   `json:"label"`
	Hours          int                      `json:"hours"`
	Since          time.Time                `json:"since"`
	Total          int                      `json:"total"`
	RiskyTotal     int                      `json:"risky_total"`
	RiskyRate      float64                  `json:"risky_rate"`
	Severity       string                   `json:"severity"`
	Message        string                   `json:"message"`
	ActionCounts   AIOutputActionCounts     `json:"action_counts"`
	Targets        []AIOutputTargetSummary  `json:"targets"`
	RiskTypes      []AIOutputDimensionCount `json:"risk_types"`
	SourceFeatures []AIOutputDimensionCount `json:"source_features"`
	MatchedRules   []AIOutputRuleCount      `json:"matched_rules"`
}

type AIOutputActionCounts struct {
	Allow    int `json:"allow"`
	SoftWarn int `json:"soft_warn"`
	Block    int `json:"block"`
}

type AIOutputTargetSummary struct {
	TargetType   string               `json:"target_type"`
	Total        int                  `json:"total"`
	RiskyTotal   int                  `json:"risky_total"`
	RiskyRate    float64              `json:"risky_rate"`
	Severity     string               `json:"severity"`
	ActionCounts AIOutputActionCounts `json:"action_counts"`
}

type AIOutputDimensionCount struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type AIOutputRuleCount struct {
	RuleID  string `json:"rule_id,omitempty"`
	Pattern string `json:"pattern,omitempty"`
	Count   int    `json:"count"`
}

type AIOutputReviewCandidates struct {
	GeneratedAt     time.Time                 `json:"generated_at"`
	WindowHours     int                       `json:"window_hours"`
	Since           time.Time                 `json:"since"`
	TotalCandidates int                       `json:"total_candidates"`
	Candidates      []AIOutputReviewCandidate `json:"candidates"`
}

type AIOutputReviewCandidate struct {
	Key            string                 `json:"key"`
	TargetType     string                 `json:"target_type"`
	RiskType       string                 `json:"risk_type"`
	SourceFeature  string                 `json:"source_feature"`
	MatchedRuleID  string                 `json:"matched_rule_id,omitempty"`
	MatchedPattern string                 `json:"matched_pattern,omitempty"`
	ActionCounts   AIOutputActionCounts   `json:"action_counts"`
	Total          int                    `json:"total"`
	LatestSeenAt   time.Time              `json:"latest_seen_at"`
	SampleTargetID *uuid.UUID             `json:"sample_target_id,omitempty"`
	Review         AIOutputReviewDecision `json:"review"`
}

type AIOutputReviewStatus string

const (
	AIOutputReviewStatusUnreviewed       AIOutputReviewStatus = "unreviewed"
	AIOutputReviewStatusFalsePositive    AIOutputReviewStatus = "false_positive"
	AIOutputReviewStatusNeedsPromptGuard AIOutputReviewStatus = "needs_prompt_guard"
	AIOutputReviewStatusNeedsRuleTuning  AIOutputReviewStatus = "needs_rule_tuning"
	AIOutputReviewStatusNeedsMasking     AIOutputReviewStatus = "needs_masking"
	AIOutputReviewStatusNeedsRegenerate  AIOutputReviewStatus = "needs_regeneration"
	AIOutputReviewStatusResolved         AIOutputReviewStatus = "resolved"
)

type AIOutputReviewDecision struct {
	CandidateKey string               `json:"candidate_key"`
	Status       AIOutputReviewStatus `json:"status"`
	Note         string               `json:"note"`
	ReviewedBy   string               `json:"reviewed_by"`
	ReviewedAt   *time.Time           `json:"reviewed_at,omitempty"`
	Metadata     map[string]any       `json:"metadata"`
	CreatedAt    *time.Time           `json:"created_at,omitempty"`
	UpdatedAt    *time.Time           `json:"updated_at,omitempty"`
}

type AIOutputReviewDecisionInput struct {
	CandidateKey string
	Status       string
	Note         string
	ReviewedBy   string
	Metadata     map[string]any
}

type AIOutputLogObservation struct {
	TargetType     string
	TargetID       *uuid.UUID
	Action         Action
	RiskType       string
	SourceFeature  string
	MatchedRuleID  string
	MatchedPattern string
	CreatedAt      time.Time
}

// ModerationRule is the super-admin facing representation of a DB-backed safety rule.
type ModerationRule struct {
	ID          uuid.UUID `json:"id"`
	RuleType    string    `json:"rule_type"`
	Pattern     string    `json:"pattern"`
	RiskType    string    `json:"risk_type"`
	Action      Action    `json:"action"`
	Locale      *string   `json:"locale,omitempty"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListRulesInput struct {
	IsActive *bool
	Action   string
	RiskType string
	RuleType string
	Locale   string
	Query    string
	Limit    int
	Offset   int
}

type RuleMutationInput struct {
	RuleType    *string
	Pattern     *string
	RiskType    *string
	Action      *string
	Locale      *string
	Description *string
}

type RuleEvent struct {
	ID          uuid.UUID      `json:"id"`
	RuleID      uuid.UUID      `json:"rule_id"`
	AdminUserID string         `json:"admin_user_id"`
	Action      string         `json:"action"`
	Before      map[string]any `json:"before"`
	After       map[string]any `json:"after"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
}
