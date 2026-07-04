package curriculum

import (
	"time"

	"github.com/google/uuid"
)

type PlanetListItem struct {
	ID                                      uuid.UUID    `json:"id"`
	DraftID                                 uuid.UUID    `json:"draft_id"`
	Title                                   string       `json:"title"`
	Status                                  PlanetStatus `json:"status"`
	IsInactive                              bool         `json:"is_inactive"`
	PlanetTypeID                            *uuid.UUID   `json:"planet_type_id"`
	PlanetTypeName                          *string      `json:"planet_type_name,omitempty"`
	PlanetTypeAsset                         *string      `json:"planet_type_asset,omitempty"`
	PlanetTextureMapID                      *uuid.UUID   `json:"planet_texture_map_id,omitempty"`
	PlanetTextureMapName                    *string      `json:"planet_texture_map_name,omitempty"`
	PlanetTextureMapAsset                   *string      `json:"planet_texture_map_asset,omitempty"`
	PlanetTextureMapRotationDurationSeconds *int         `json:"planet_texture_map_rotation_duration_seconds,omitempty"`
	PlanetTextureMapRotationDirection       *string      `json:"planet_texture_map_rotation_direction,omitempty"`
	LessonCount                             int          `json:"lesson_count"`
	CompletedLessons                        int          `json:"completed_lesson_count"`
	Progress                                *float64     `json:"progress,omitempty"`
	CanComplete                             bool         `json:"can_complete"`
	LastAccessedAt                          *time.Time   `json:"last_accessed_at,omitempty"`
	UpdatedAt                               time.Time    `json:"updated_at"`
	Destination                             string       `json:"destination"`
	ShareCount                              int          `json:"share_count,omitempty"`
}

type PlanetAggregate struct {
	Planet      PlanetListItem     `json:"planet"`
	GoalContext *PlanetGoalContext `json:"goal_context,omitempty"`
	Lessons     []CourseLessonTree `json:"lessons"`
}

type PlanetPointDetail struct {
	Planet          PlanetListItem       `json:"planet"`
	GoalContext     *PlanetGoalContext   `json:"goal_context,omitempty"`
	LevelTitle      string               `json:"level_title"`
	LessonTitle     string               `json:"lesson_title"`
	LessonID        uuid.UUID            `json:"lesson_id"`
	Point           CoursePointAggregate `json:"point"`
	PreviousPointID *uuid.UUID           `json:"previous_point_id,omitempty"`
	NextPointID     *uuid.UUID           `json:"next_point_id,omitempty"`
}

type PlanetGoalContext struct {
	LearningGoal       *string    `json:"learning_goal,omitempty"`
	GoalProfileID      *uuid.UUID `json:"goal_profile_id,omitempty"`
	GoalProfileVersion *int       `json:"goal_profile_version,omitempty"`
	ConfirmedGoal      *string    `json:"confirmed_goal,omitempty"`
	UsageContext       *string    `json:"usage_context,omitempty"`
	Motivation         *string    `json:"motivation,omitempty"`
	GoalReadiness      *float64   `json:"goal_readiness,omitempty"`
}

type CoursePointAggregate struct {
	Point                       CoursePoint                  `json:"point"`
	Blocks                      []CoursePointBlock           `json:"blocks,omitempty"`
	ResearchMaterialConfirmed   bool                         `json:"research_material_confirmed"`
	ResearchMaterialConfirmedAt *time.Time                   `json:"research_material_confirmed_at,omitempty"`
	Questions                   []CoursePointQuestion        `json:"questions,omitempty"`
	SelfEvaluation              *CoursePointSelfEvaluation   `json:"self_evaluation,omitempty"`
	AISummary                   *PointAISummaryEntry         `json:"ai_summary_entry,omitempty"`
	JournalEntry                *DraftJournalEntry           `json:"journal_entry,omitempty"`
	RecordEntry                 *DraftRecordEntry            `json:"record_entry,omitempty"`
	ArtifactEntry               *DraftArtifactEntry          `json:"artifact_entry,omitempty"`
	Artifacts                   []CoursePointArtifact        `json:"artifacts,omitempty"`
	Attachments                 []CoursePointAttachment      `json:"attachments,omitempty"`
	MaterialReports             []CoursePointMaterialReport  `json:"material_reports,omitempty"`
	PracticeLogs                []CoursePointPracticeLog     `json:"practice_logs,omitempty"`
	Events                      []CoursePointEvent           `json:"events,omitempty"`
	ObservationNotes            []CoursePointObservationNote `json:"observation_notes,omitempty"`
}

type PlanetRecordAggregate struct {
	Course  PlanetRecordCourse   `json:"course"`
	Lessons []PlanetRecordLesson `json:"lessons"`
}

type PlanetRecordFeedResponse struct {
	RouteKind string                   `json:"route_kind"`
	Course    PlanetRecordFeedCourse   `json:"course"`
	Filters   []PlanetRecordFeedFilter `json:"filters"`
	Cards     []PlanetRecordCard       `json:"cards"`
}

type PlanetRecordFeedCourse struct {
	ID                  uuid.UUID `json:"id"`
	Title               string    `json:"title"`
	Status              string    `json:"status"`
	TotalPoints         int       `json:"total_points"`
	CompletedPoints     int       `json:"completed_points"`
	TotalLessons        int       `json:"total_lessons"`
	CompletedLessons    int       `json:"completed_lessons"`
	RecordCardCount     int       `json:"record_card_count"`
	ShareCandidateCount int       `json:"share_candidate_count"`
}

type PlanetRecordFeedFilter struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type PlanetRecordCard struct {
	ID                  string                 `json:"id"`
	CardType            string                 `json:"card_type"`
	Title               string                 `json:"title"`
	Summary             string                 `json:"summary"`
	Badge               string                 `json:"badge"`
	Category            string                 `json:"category"`
	Visibility          string                 `json:"visibility"`
	CourseID            uuid.UUID              `json:"course_id"`
	LessonID            *uuid.UUID             `json:"lesson_id,omitempty"`
	PointID             *uuid.UUID             `json:"point_id,omitempty"`
	LessonTitle         string                 `json:"lesson_title,omitempty"`
	PointTitle          string                 `json:"point_title,omitempty"`
	PointType           PointType              `json:"point_type,omitempty"`
	OccurredAt          time.Time              `json:"occurred_at"`
	ShareCandidateScore int                    `json:"share_candidate_score"`
	ShareText           string                 `json:"share_text"`
	Detail              PlanetRecordCardDetail `json:"detail"`
	Actions             PlanetRecordCardAction `json:"actions"`
}

type PlanetRecordCardDetail struct {
	PrimaryText           string                      `json:"primary_text,omitempty"`
	JournalExcerpt        string                      `json:"journal_excerpt,omitempty"`
	QuestionExcerpt       string                      `json:"question_excerpt,omitempty"`
	ArtifactTitle         string                      `json:"artifact_title,omitempty"`
	SelfEvaluationSummary string                      `json:"self_evaluation_summary,omitempty"`
	PracticeSummary       string                      `json:"practice_summary,omitempty"`
	LessonProgress        *PlanetRecordLessonProgress `json:"lesson_progress,omitempty"`
}

type PlanetRecordLessonProgress struct {
	CompletedPoints int `json:"completed_points"`
	TotalPoints     int `json:"total_points"`
	Percent         int `json:"percent"`
}

type PlanetRecordCardAction struct {
	CanOpenLearningPage bool `json:"can_open_learning_page"`
	CanCopyShareText    bool `json:"can_copy_share_text"`
	CanCreateShareCard  bool `json:"can_create_share_card"`
}

type PlanetRecordCourse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	TotalPoints    int       `json:"total_points"`
	RecordedPoints int       `json:"recorded_points"`
}

type PlanetRecordLesson struct {
	LessonID    uuid.UUID           `json:"lesson_id"`
	LessonTitle string              `json:"lesson_title"`
	Points      []PlanetRecordPoint `json:"points"`
}

type PlanetRecordPoint struct {
	PointID                       uuid.UUID                  `json:"point_id"`
	PointTitle                    string                     `json:"point_title"`
	PointType                     PointType                  `json:"point_type"`
	Blocks                        []CoursePointBlock         `json:"blocks,omitempty"`
	JournalEntry                  *DraftJournalEntry         `json:"journal_entry,omitempty"`
	ApplicationNote               string                     `json:"application_note"`
	SelfEvaluationApplicationNote string                     `json:"self_evaluation_application_note"`
	GoalAlignmentNote             string                     `json:"goal_alignment_note"`
	Questions                     []CoursePointQuestion      `json:"questions,omitempty"`
	RecordEntry                   *DraftRecordEntry          `json:"record_entry,omitempty"`
	SelfEvaluation                *CoursePointSelfEvaluation `json:"self_evaluation,omitempty"`
}

type PlanetResultAggregate struct {
	Course  PlanetResultCourse   `json:"course"`
	Lessons []PlanetResultLesson `json:"lessons"`
}

type PlanetResultCourse struct {
	ID                  uuid.UUID `json:"id"`
	Title               string    `json:"title"`
	TotalPoints         int       `json:"total_points"`
	PointsWithArtifacts int       `json:"points_with_artifacts"`
}

type PlanetResultLesson struct {
	LessonID    uuid.UUID           `json:"lesson_id"`
	LessonTitle string              `json:"lesson_title"`
	Points      []PlanetResultPoint `json:"points"`
}

type PlanetResultPoint struct {
	PointID    uuid.UUID             `json:"point_id"`
	PointTitle string                `json:"point_title"`
	PointType  PointType             `json:"point_type"`
	Artifacts  []CoursePointArtifact `json:"artifacts"`
}
