package curriculum

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CoursePointArtifact struct {
	ID                uuid.UUID `json:"id"`
	CoursePointID     uuid.UUID `json:"course_point_id"`
	ArtifactType      string    `json:"artifact_type"`
	Title             string    `json:"title"`
	URL               string    `json:"url"`
	Description       string    `json:"description"`
	PointCategory     string    `json:"point_category"`
	ProductionProcess string    `json:"production_process"`
	LearnedPoints     string    `json:"learned_points"`
	DifficultPoints   string    `json:"difficult_points"`
	Visibility        string    `json:"visibility"`
	OrderIndex        int       `json:"order_index"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CoursePointAttachment struct {
	ID             uuid.UUID  `json:"id"`
	CoursePointID  uuid.UUID  `json:"course_point_id"`
	ArtifactID     *uuid.UUID `json:"artifact_id,omitempty"`
	UserID         uuid.UUID  `json:"user_id"`
	Provider       string     `json:"provider"`
	AttachmentType string     `json:"attachment_type"`
	SourceContext  string     `json:"source_context"`
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	FilePath       string     `json:"file_path"`
	FileSize       *int64     `json:"file_size,omitempty"`
	MimeType       string     `json:"mime_type"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CoursePointMaterialReport struct {
	ID                   uuid.UUID  `json:"id"`
	CoursePointID        uuid.UUID  `json:"course_point_id"`
	UserID               uuid.UUID  `json:"user_id"`
	AttachmentID         *uuid.UUID `json:"attachment_id,omitempty"`
	TargetType           string     `json:"target_type"`
	ReportType           string     `json:"report_type"`
	Message              string     `json:"message"`
	Status               string     `json:"status"`
	TargetContentID      *uuid.UUID `json:"target_content_id,omitempty"`
	TargetURL            string     `json:"target_url,omitempty"`
	TargetTitle          string     `json:"target_title,omitempty"`
	ReplacementContentID *uuid.UUID `json:"replacement_content_id,omitempty"`
	ReplacementURL       string     `json:"replacement_url,omitempty"`
	ReplacementTitle     string     `json:"replacement_title,omitempty"`
	ReplacedAt           *time.Time `json:"replaced_at,omitempty"`
	AdminNote            string     `json:"admin_note,omitempty"`
	ReviewedBy           *string    `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type CoursePointPracticeLog struct {
	ID              uuid.UUID `json:"id"`
	CoursePointID   uuid.UUID `json:"course_point_id"`
	UserID          uuid.UUID `json:"user_id"`
	Title           string    `json:"title"`
	ActivityName    string    `json:"activity_name"`
	AttemptCount    *int      `json:"attempt_count,omitempty"`
	SuccessCount    *int      `json:"success_count,omitempty"`
	FailureCount    *int      `json:"failure_count,omitempty"`
	DurationMinutes *int      `json:"duration_minutes,omitempty"`
	BlockedPart     string    `json:"blocked_part"`
	ChangedMethod   string    `json:"changed_method"`
	AchievementNote string    `json:"achievement_note"`
	Achievement     string    `json:"achievement"`
	NextPlan        string    `json:"next_plan"`
	NextPractice    string    `json:"next_practice"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CoursePointEvent struct {
	ID            uuid.UUID       `json:"id"`
	CoursePointID uuid.UUID       `json:"course_point_id"`
	UserID        uuid.UUID       `json:"user_id"`
	EventType     string          `json:"event_type"`
	EventPayload  json.RawMessage `json:"event_payload"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CoursePointObservationNote struct {
	ID            uuid.UUID `json:"id"`
	CoursePointID uuid.UUID `json:"course_point_id"`
	UserID        uuid.UUID `json:"user_id"`
	NoteType      string    `json:"note_type"`
	Content       string    `json:"content"`
	OrderIndex    int       `json:"order_index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateLearningPointObservationNoteRequest struct {
	NoteType string `json:"note_type"`
	Content  string `json:"content"`
}

type UpdateLearningPointObservationNoteRequest struct {
	NoteType *string `json:"note_type"`
	Content  *string `json:"content"`
}

type CoursePointLearningSession struct {
	ID             uuid.UUID  `json:"id"`
	CoursePointID  uuid.UUID  `json:"course_point_id"`
	UserID         uuid.UUID  `json:"user_id"`
	StartedAt      time.Time  `json:"started_at"`
	LastSeenAt     time.Time  `json:"last_seen_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	ActiveSeconds  int        `json:"active_seconds"`
	HeartbeatCount int        `json:"heartbeat_count"`
	EndReason      *string    `json:"end_reason,omitempty"`
	LastVisibility string     `json:"last_visibility"`
	UserAgent      *string    `json:"user_agent,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type StartLearningPointSessionRequest struct {
	Visibility string `json:"visibility"`
	UserAgent  string `json:"user_agent"`
}

type UpdateLearningPointSessionRequest struct {
	ActiveSecondsDelta int    `json:"active_seconds_delta"`
	Visibility         string `json:"visibility"`
	Reason             string `json:"reason"`
}

type UpdateLearningPointGoalRequest struct {
	PointGoal     *string `json:"point_goal"`
	PointCategory *string `json:"point_category"`
}

type CreateLearningPointQuestionRequest struct {
	Title        *string           `json:"title"`
	Question     string            `json:"question"`
	QuestionType PointQuestionType `json:"question_type"`
	AnswerMethod *string           `json:"answer_method"`
}

type UpdateLearningPointQuestionRequest struct {
	Title        *string              `json:"title"`
	Question     *string              `json:"question"`
	QuestionType *PointQuestionType   `json:"question_type"`
	AnswerMethod *string              `json:"answer_method"`
	Answer       *string              `json:"answer"`
	AIFeedback   *string              `json:"ai_feedback"`
	Status       *PointQuestionStatus `json:"status"`
}

type UpsertLearningPointSelfEvaluationRequest struct {
	Understanding        *int    `json:"understanding"`
	ApplicationNote      *string `json:"application_note"`
	Proficiency          *int    `json:"proficiency"`
	UnderstandingScore   *int    `json:"understanding_score"`
	UnderstandingReason  *string `json:"understanding_reason"`
	ApplicationScore     *int    `json:"application_score"`
	ApplicationReason    *string `json:"application_reason"`
	ProficiencyScore     *int    `json:"proficiency_score"`
	ProficiencyReason    *string `json:"proficiency_reason"`
	ProblemSolvingScore  *int    `json:"problem_solving_score"`
	ProblemSolvingReason *string `json:"problem_solving_reason"`
	ExpressionScore      *int    `json:"expression_score"`
	ExpressionReason     *string `json:"expression_reason"`
	GoalAlignmentNote    *string `json:"goal_alignment_note"`
	FinalScore           *int    `json:"final_score"`
}

type LearningPointSelfEvaluationAIDraftRequest struct {
	ApplicationAnswers []LearningPointSelfEvaluationApplicationAnswer `json:"application_answers"`
}

type LearningPointSelfEvaluationApplicationAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LearningPointSelfEvaluationApplicationQuestion struct {
	Question string `json:"question"`
	Intent   string `json:"intent"`
}

type LearningPointSelfEvaluationAIDraft struct {
	ApplicationQuestions []LearningPointSelfEvaluationApplicationQuestion `json:"application_questions"`
	UnderstandingScore   int                                              `json:"understanding_score"`
	UnderstandingReason  string                                           `json:"understanding_reason"`
	ApplicationScore     int                                              `json:"application_score"`
	ApplicationReason    string                                           `json:"application_reason"`
	ProficiencyScore     int                                              `json:"proficiency_score"`
	ProficiencyReason    string                                           `json:"proficiency_reason"`
	ProblemSolvingScore  int                                              `json:"problem_solving_score"`
	ProblemSolvingReason string                                           `json:"problem_solving_reason"`
	ExpressionScore      int                                              `json:"expression_score"`
	ExpressionReason     string                                           `json:"expression_reason"`
	GoalAlignmentNote    string                                           `json:"goal_alignment_note"`
}

type CreateLearningPointArtifactRequest struct {
	ArtifactType      *string `json:"artifact_type"`
	Title             *string `json:"title"`
	URL               *string `json:"url"`
	Description       *string `json:"description"`
	PointCategory     *string `json:"point_category"`
	ProductionProcess *string `json:"production_process"`
	LearnedPoints     *string `json:"learned_points"`
	DifficultPoints   *string `json:"difficult_points"`
	Visibility        *string `json:"visibility"`
	OrderIndex        *int    `json:"order_index"`
}

type UpdateLearningPointArtifactRequest struct {
	ArtifactType      *string `json:"artifact_type"`
	Title             *string `json:"title"`
	URL               *string `json:"url"`
	Description       *string `json:"description"`
	PointCategory     *string `json:"point_category"`
	ProductionProcess *string `json:"production_process"`
	LearnedPoints     *string `json:"learned_points"`
	DifficultPoints   *string `json:"difficult_points"`
	Visibility        *string `json:"visibility"`
	OrderIndex        *int    `json:"order_index"`
}

type CreateLearningPointAttachmentRequest struct {
	Provider       *string    `json:"provider"`
	AttachmentType *string    `json:"attachment_type"`
	SourceContext  *string    `json:"source_context"`
	ArtifactID     *uuid.UUID `json:"artifact_id"`
	Title          *string    `json:"title"`
	URL            *string    `json:"url"`
	FilePath       *string    `json:"file_path"`
	FileSize       *int64     `json:"file_size"`
	MimeType       *string    `json:"mime_type"`
}

type UpdateLearningPointAttachmentRequest struct {
	Provider       *string    `json:"provider"`
	AttachmentType *string    `json:"attachment_type"`
	SourceContext  *string    `json:"source_context"`
	ArtifactID     *uuid.UUID `json:"artifact_id"`
	Title          *string    `json:"title"`
	URL            *string    `json:"url"`
	FilePath       *string    `json:"file_path"`
	FileSize       *int64     `json:"file_size"`
	MimeType       *string    `json:"mime_type"`
}

type CreateLearningPointMaterialReportRequest struct {
	TargetType   *string    `json:"target_type"`
	ReportType   *string    `json:"report_type"`
	Message      *string    `json:"message"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
}

type ReplaceLearningPointMaterialRequest struct {
	ReportID     *uuid.UUID `json:"report_id"`
	ContentID    *uuid.UUID `json:"content_id"`
	Title        *string    `json:"title"`
	URL          *string    `json:"url"`
	ThumbnailURL *string    `json:"thumbnail_url"`
}

type CreateLearningPointPracticeLogRequest struct {
	Title           *string `json:"title"`
	ActivityName    *string `json:"activity_name"`
	AttemptCount    *int    `json:"attempt_count"`
	SuccessCount    *int    `json:"success_count"`
	FailureCount    *int    `json:"failure_count"`
	DurationMinutes *int    `json:"duration_minutes"`
	BlockedPart     *string `json:"blocked_part"`
	ChangedMethod   *string `json:"changed_method"`
	AchievementNote *string `json:"achievement_note"`
	Achievement     *string `json:"achievement"`
	NextPlan        *string `json:"next_plan"`
	NextPractice    *string `json:"next_practice"`
}

type UpdateLearningPointPracticeLogRequest struct {
	Title           *string `json:"title"`
	ActivityName    *string `json:"activity_name"`
	AttemptCount    *int    `json:"attempt_count"`
	SuccessCount    *int    `json:"success_count"`
	FailureCount    *int    `json:"failure_count"`
	DurationMinutes *int    `json:"duration_minutes"`
	BlockedPart     *string `json:"blocked_part"`
	ChangedMethod   *string `json:"changed_method"`
	AchievementNote *string `json:"achievement_note"`
	Achievement     *string `json:"achievement"`
	NextPlan        *string `json:"next_plan"`
	NextPractice    *string `json:"next_practice"`
}
