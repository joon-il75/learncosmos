package curriculum

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID                                      uuid.UUID    `json:"id"`
	UserID                                  uuid.UUID    `json:"user_id"`
	SourceDraftID                           *uuid.UUID   `json:"source_draft_id"`
	SourceQuery                             string       `json:"source_query"`
	LearningGoal                            *string      `json:"learning_goal"`
	CurrentLevel                            *string      `json:"current_level"`
	DurationWeeks                           *int         `json:"duration_weeks"`
	StudyHoursPerWeek                       *int         `json:"study_hours_per_week"`
	PreferredFormat                         *string      `json:"preferred_format"`
	Title                                   string       `json:"title"`
	Description                             *string      `json:"description"`
	CompletionCriteria                      []string     `json:"completion_criteria,omitempty"`
	Status                                  CourseStatus `json:"status"`
	PlanetTypeID                            *uuid.UUID   `json:"planet_type_id"`
	PlanetTypeName                          *string      `json:"planet_type_name,omitempty"`
	PlanetTypeAsset                         *string      `json:"planet_type_asset,omitempty"`
	PlanetTextureMapID                      *uuid.UUID   `json:"planet_texture_map_id,omitempty"`
	PlanetTextureMapName                    *string      `json:"planet_texture_map_name,omitempty"`
	PlanetTextureMapAsset                   *string      `json:"planet_texture_map_asset,omitempty"`
	PlanetTextureMapRotationDurationSeconds *int         `json:"planet_texture_map_rotation_duration_seconds,omitempty"`
	PlanetTextureMapRotationDirection       *string      `json:"planet_texture_map_rotation_direction,omitempty"`
	CreatedAt                               time.Time    `json:"created_at"`
	UpdatedAt                               time.Time    `json:"updated_at"`
}

type CourseLesson struct {
	ID                       uuid.UUID                      `json:"id"`
	CourseID                 uuid.UUID                      `json:"course_id"`
	ParentLessonID           *uuid.UUID                     `json:"parent_lesson_id,omitempty"`
	Title                    string                         `json:"title"`
	Objective                *string                        `json:"objective"`
	Summary                  *string                        `json:"summary"`
	DifficultyLevel          *string                        `json:"difficulty_level"`
	RecommendationSearchSpec LessonRecommendationSearchSpec `json:"recommendation_search_spec,omitempty"`
	LessonRole               LessonRole                     `json:"lesson_role"`
	SourceType               LessonSourceType               `json:"source_type"`
	OrderIndex               int                            `json:"order_index"`
	CreatedAt                time.Time                      `json:"created_at"`
	UpdatedAt                time.Time                      `json:"updated_at"`
}

type CoursePoint struct {
	ID             uuid.UUID                 `json:"id"`
	CourseID       uuid.UUID                 `json:"course_id"`
	CourseLessonID uuid.UUID                 `json:"course_lesson_id"`
	PointType      PointType                 `json:"point_type"`
	Status         PointStatus               `json:"status"`
	Title          string                    `json:"title"`
	Description    *string                   `json:"description"`
	PointGoal      *string                   `json:"point_goal,omitempty"`
	PointCategory  *string                   `json:"point_category,omitempty"`
	TemplateType   *ResearchNodeTemplateType `json:"template_type,omitempty"`
	ContentID      *uuid.UUID                `json:"content_id,omitempty"`
	ExternalURL    *string                   `json:"external_url,omitempty"`
	ThumbnailURL   *string                   `json:"thumbnail_url,omitempty"`
	PriceType      *PriceType                `json:"price_type,omitempty"`
	RankScore      *float64                  `json:"rank_score,omitempty"`
	OrderIndex     int                       `json:"order_index"`
	CompletedAt    *time.Time                `json:"completed_at,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

type CoursePointBlock struct {
	ID            uuid.UUID             `json:"id"`
	CoursePointID uuid.UUID             `json:"course_point_id"`
	UserID        *uuid.UUID            `json:"user_id,omitempty"`
	BlockType     ResearchNodeBlockType `json:"block_type"`
	Content       json.RawMessage       `json:"content"`
	OrderIndex    int                   `json:"order_index"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

type CoursePointQuestion struct {
	ID                 uuid.UUID           `json:"id"`
	CoursePointID      uuid.UUID           `json:"course_point_id"`
	GoalProfileVersion *int                `json:"goal_profile_version,omitempty"`
	Title              *string             `json:"title,omitempty"`
	Question           string              `json:"question"`
	QuestionType       PointQuestionType   `json:"question_type"`
	AnswerMethod       *string             `json:"answer_method,omitempty"`
	Answer             *string             `json:"answer,omitempty"`
	AIFeedback         *string             `json:"ai_feedback,omitempty"`
	Status             PointQuestionStatus `json:"status"`
	CreatedBy          string              `json:"created_by"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

type CoursePointSelfEvaluation struct {
	CoursePointID        uuid.UUID `json:"course_point_id"`
	GoalProfileVersion   *int      `json:"goal_profile_version,omitempty"`
	Understanding        int       `json:"understanding"`
	ApplicationNote      string    `json:"application_note"`
	Proficiency          int       `json:"proficiency"`
	UnderstandingScore   *int      `json:"understanding_score,omitempty"`
	UnderstandingReason  string    `json:"understanding_reason"`
	ApplicationScore     *int      `json:"application_score,omitempty"`
	ApplicationReason    string    `json:"application_reason"`
	ProficiencyScore     *int      `json:"proficiency_score,omitempty"`
	ProficiencyReason    string    `json:"proficiency_reason"`
	ProblemSolvingScore  *int      `json:"problem_solving_score,omitempty"`
	ProblemSolvingReason string    `json:"problem_solving_reason"`
	ExpressionScore      *int      `json:"expression_score,omitempty"`
	ExpressionReason     string    `json:"expression_reason"`
	GoalAlignmentNote    string    `json:"goal_alignment_note"`
	FinalScore           *int      `json:"final_score,omitempty"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CommunityPost struct {
	ID         uuid.UUID  `json:"id"`
	CourseID   uuid.UUID  `json:"course_id"`
	LessonID   uuid.UUID  `json:"lesson_id"`
	PointID    *uuid.UUID `json:"point_id,omitempty"`
	UserID     uuid.UUID  `json:"user_id"`
	PostType   string     `json:"post_type"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Visibility string     `json:"visibility"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type RecommendationEvent struct {
	ID                  uuid.UUID               `json:"id"`
	UserID              uuid.UUID               `json:"user_id"`
	CourseDraftID       *uuid.UUID              `json:"course_draft_id"`
	CourseDraftLessonID *uuid.UUID              `json:"course_draft_lesson_id"`
	ContentID           *uuid.UUID              `json:"content_id"`
	EventType           RecommendationEventType `json:"event_type"`
	SourceQuery         *string                 `json:"source_query"`
	Payload             map[string]any          `json:"payload"`
	CreatedAt           time.Time               `json:"created_at"`
}

type CourseAggregate struct {
	Course  Course             `json:"course"`
	Lessons []CourseLessonTree `json:"lessons"`
}

type CourseLessonTree struct {
	Lesson       CourseLesson              `json:"lesson"`
	Points       []CoursePointAggregate    `json:"points,omitempty"`
	SubLessons   []CourseLessonTree        `json:"sub_lessons,omitempty"`
	RuntimeEntry *CourseLessonRuntimeEntry `json:"runtime_entry,omitempty"`
}

type CourseLessonRuntimeEntry struct {
	ID             uuid.UUID           `json:"id"`
	UserID         uuid.UUID           `json:"user_id"`
	CourseID       uuid.UUID           `json:"course_id"`
	CourseLessonID uuid.UUID           `json:"course_lesson_id"`
	Status         LessonRuntimeStatus `json:"status"`
	StartedAt      time.Time           `json:"started_at"`
	CompletedAt    *time.Time          `json:"completed_at,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type UpdateLearningLessonRuntimeRequest struct {
	Status LessonRuntimeStatus `json:"status"`
}

type UpdateLearningPointRuntimeRequest struct {
	Status LessonRuntimeStatus `json:"status"`
}
