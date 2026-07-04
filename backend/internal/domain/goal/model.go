package goal

import (
	"time"

	"github.com/google/uuid"
)

type InterviewState string

const (
	StateListening               InterviewState = "listening"
	StateClarifying              InterviewState = "clarifying"
	StateProposingGoal           InterviewState = "proposing_goal"
	StateConfirmed               InterviewState = "confirmed"
	StateRevisingGoal            InterviewState = "revising_goal"
	StateAwaitingRebuildDecision InterviewState = "awaiting_rebuild_decision"
)

type RebuildDecision string

const (
	RebuildKeepStructure RebuildDecision = "keep_structure"
	// RebuildRemaining is kept only for backward compatibility with older saved goals.
	RebuildRemaining RebuildDecision = "rebuild_remaining"
	RebuildAll       RebuildDecision = "rebuild_all"
)

func NormalizeRebuildDecision(decision RebuildDecision) RebuildDecision {
	if decision == RebuildRemaining {
		return RebuildAll
	}
	return decision
}

type ReplyIntent string

const (
	ReplyExplainPurpose        ReplyIntent = "explain_purpose"
	ReplyAskClarifyingQuestion ReplyIntent = "ask_clarifying_question"
	ReplyNarrowDirection       ReplyIntent = "narrow_direction"
	ReplyProposeGoal           ReplyIntent = "propose_goal"
	ReplyConfirmGoal           ReplyIntent = "confirm_goal"
	ReplyReviseGoal            ReplyIntent = "revise_goal"
	ReplyAskRebuildDecision    ReplyIntent = "ask_rebuild_decision"
	ReplyAcknowledgeContinue   ReplyIntent = "acknowledge_and_continue"
)

type InterviewMessage struct {
	Role    string `json:"role"` // "lumi" | "user"
	Content string `json:"content"`
}

type GoalProfile struct {
	ID                uuid.UUID             `json:"id"`
	CourseDraftID     *uuid.UUID            `json:"course_draft_id"`
	UserID            uuid.UUID             `json:"user_id"`
	UserIntent        string                `json:"user_intent"`
	Motivation        *string               `json:"motivation"`
	UsageContext      *string               `json:"usage_context"`
	ConfirmedGoal     *string               `json:"confirmed_goal"`
	GoalType          *string               `json:"goal_type"`
	OutputType        *string               `json:"output_type"`
	DifficultyLevel   *string               `json:"difficulty_level"`
	TimeHorizon       *string               `json:"time_horizon"`
	Language          string                `json:"language"`
	SummarizedContext string                `json:"summarized_context"`
	LearningIntent    LearningIntentProfile `json:"learning_intent_profile"`
	RebuildDecision   *RebuildDecision      `json:"rebuild_decision"`
	RevisionSnapshot  *GoalRevisionState    `json:"revision_snapshot,omitempty"`
	InterviewState    InterviewState        `json:"interview_state"`
	Messages          []InterviewMessage    `json:"messages"`
	Version           int                   `json:"version"`
	IsActive          bool                  `json:"is_active"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

type LearningIntentProfile struct {
	ConfirmedGoal       string   `json:"confirmed_goal,omitempty"`
	LearnerLevel        string   `json:"learner_level,omitempty"`
	Purpose             string   `json:"purpose,omitempty"`
	DesiredOutput       string   `json:"desired_output,omitempty"`
	CurrentBlockers     []string `json:"current_blockers,omitempty"`
	PreferredActivities []string `json:"preferred_activities,omitempty"`
	SuccessCriteria     []string `json:"success_criteria,omitempty"`
	Source              string   `json:"source,omitempty"`
	Confidence          string   `json:"confidence,omitempty"`
}

type GoalRevisionState struct {
	UserIntent        string                `json:"user_intent"`
	Motivation        *string               `json:"motivation"`
	UsageContext      *string               `json:"usage_context"`
	ConfirmedGoal     *string               `json:"confirmed_goal"`
	GoalType          *string               `json:"goal_type"`
	OutputType        *string               `json:"output_type"`
	DifficultyLevel   *string               `json:"difficulty_level"`
	TimeHorizon       *string               `json:"time_horizon"`
	SummarizedContext string                `json:"summarized_context"`
	LearningIntent    LearningIntentProfile `json:"learning_intent_profile,omitempty"`
	InterviewState    InterviewState        `json:"interview_state"`
	Messages          []InterviewMessage    `json:"messages"`
	Version           int                   `json:"version"`
}

type StartInterviewRequest struct {
	UserIntent string `json:"user_intent" binding:"required,max=500"`
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required,max=1000"`
}

type ConfirmGoalRequest struct {
	ConfirmedGoal string `json:"confirmed_goal" binding:"required,max=500"`
}

type AttachDraftRequest struct {
	DraftID uuid.UUID `json:"draft_id" binding:"required"`
}

type ReviseGoalRequest struct {
	Message string `json:"message" binding:"required,max=1000"`
}

type RebuildDecisionRequest struct {
	Decision RebuildDecision `json:"decision" binding:"required"`
}

// Step A: 구조화 판단 결과
type IntentAnalysisResult struct {
	Extracted         ExtractedInfo `json:"extracted"`
	UserStance        string        `json:"user_stance"` // accepting | rejecting | revising | neutral
	GoalCandidate     *string       `json:"goal_candidate"`
	IsGoalClearEnough bool          `json:"is_goal_clear_enough"`
	ReasonSummary     string        `json:"reason_summary"`
}

// Step B: 자연어 응답 생성용 컨텍스트
type ResponseContext struct {
	State               InterviewState
	Language            string
	LatestUserMessage   string
	ConversationSummary string
	Extracted           ExtractedInfo
	MissingFields       []string
	GoalCandidate       *string
	ConfirmedGoal       *string
	ReplyIntent         ReplyIntent
	RecentMessages      []string
}

// Step B: LLM 자연어 응답 결과
type AssistantReplyPayload struct {
	ReplyIntent      ReplyIntent `json:"reply_intent"`
	AssistantMessage string      `json:"assistant_message"`
}

// LLM이 반환하는 인터뷰 응답 구조체
type InterviewAIResponse struct {
	Message        string         `json:"message"`
	NextState      InterviewState `json:"next_state"`
	Extracted      ExtractedInfo  `json:"extracted"`
	ProposedGoal   *string        `json:"proposed_goal"`
	GoalCandidate  *string        `json:"goal_candidate"`
	MissingInfo    []string       `json:"missing_info"`
	Route          string         `json:"-"`
	LLMCallCount   int            `json:"-"`
	FallbackReason string         `json:"-"`
}

type ExtractedInfo struct {
	Motivation      *string `json:"motivation"`
	UsageContext    *string `json:"usage_context"`
	GoalType        *string `json:"goal_type"`
	OutputType      *string `json:"output_type"`
	DifficultyLevel *string `json:"difficulty_level"`
	TimeHorizon     *string `json:"time_horizon"`
}
