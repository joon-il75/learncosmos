package curriculum

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CourseDraft struct {
	ID                                      uuid.UUID   `json:"id"`
	UserID                                  uuid.UUID   `json:"user_id"`
	SourceQuery                             string      `json:"source_query"`
	LearningGoal                            *string     `json:"learning_goal"`
	GoalProfileID                           *uuid.UUID  `json:"goal_profile_id,omitempty"`
	GoalProfileVersion                      *int        `json:"goal_profile_version,omitempty"`
	CurrentLevel                            *string     `json:"current_level"`
	DurationWeeks                           *int        `json:"duration_weeks"`
	StudyHoursPerWeek                       *int        `json:"study_hours_per_week"`
	PreferredFormat                         *string     `json:"preferred_format"`
	GenerationLanguage                      string      `json:"generation_language,omitempty"`
	Title                                   string      `json:"title"`
	Description                             *string     `json:"description"`
	CompletionCriteria                      []string    `json:"completion_criteria,omitempty"`
	Status                                  DraftStatus `json:"status"`
	ConfirmedCourseID                       *uuid.UUID  `json:"confirmed_course_id"`
	IsInactive                              bool        `json:"is_inactive"`
	PlanetTypeID                            *uuid.UUID  `json:"planet_type_id"`
	PlanetTypeName                          *string     `json:"planet_type_name,omitempty"`
	PlanetTypeAsset                         *string     `json:"planet_type_asset,omitempty"`
	PlanetTextureMapID                      *uuid.UUID  `json:"planet_texture_map_id,omitempty"`
	PlanetTextureMapName                    *string     `json:"planet_texture_map_name,omitempty"`
	PlanetTextureMapAsset                   *string     `json:"planet_texture_map_asset,omitempty"`
	PlanetTextureMapRotationDurationSeconds *int        `json:"planet_texture_map_rotation_duration_seconds,omitempty"`
	PlanetTextureMapRotationDirection       *string     `json:"planet_texture_map_rotation_direction,omitempty"`
	LastAccessedAt                          *time.Time  `json:"last_accessed_at,omitempty"`
	CreatedAt                               time.Time   `json:"created_at"`
	UpdatedAt                               time.Time   `json:"updated_at"`
}

type CourseDraftLesson struct {
	ID                       uuid.UUID                      `json:"id"`
	CourseDraftID            uuid.UUID                      `json:"course_draft_id"`
	ParentLessonID           *uuid.UUID                     `json:"parent_lesson_id,omitempty"`
	Title                    string                         `json:"title"`
	Objective                *string                        `json:"objective"`
	Summary                  *string                        `json:"summary"`
	DifficultyLevel          *string                        `json:"difficulty_level"`
	RecommendationSearchSpec LessonRecommendationSearchSpec `json:"recommendation_search_spec,omitempty"`
	OperationNote            *string                        `json:"operation_note,omitempty"`
	JournalEntry             *DraftJournalEntry             `json:"journal_entry,omitempty"`
	RecordEntry              *DraftRecordEntry              `json:"record_entry,omitempty"`
	ArtifactEntry            *DraftArtifactEntry            `json:"artifact_entry,omitempty"`
	LessonRole               LessonRole                     `json:"lesson_role"`
	SourceType               LessonSourceType               `json:"source_type"`
	OrderIndex               int                            `json:"order_index"`
	CreatedAt                time.Time                      `json:"created_at"`
	UpdatedAt                time.Time                      `json:"updated_at"`
}

type DraftJournalEntry struct {
	Observation    string    `json:"observation"`
	Reflection     string    `json:"reflection"`
	NextStep       string    `json:"next_step"`
	CoreConcept    string    `json:"core_concept"`
	MyExplanation  string    `json:"my_explanation"`
	Examples       string    `json:"examples"`
	ConfusedParts  string    `json:"confused_parts"`
	ReferenceLinks string    `json:"reference_links"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DraftRecordEntry struct {
	StudyMinutes    int       `json:"study_minutes"`
	PracticeCount   int       `json:"practice_count"`
	ConfidenceLevel int       `json:"confidence_level"`
	ApplicationNote string    `json:"application_note"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DraftArtifactEntry struct {
	ArtifactType string    `json:"artifact_type"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Description  string    `json:"description"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PointAISummaryEntry struct {
	SourceTitle       string    `json:"source_title"`
	SourceDescription string    `json:"source_description"`
	Summary           string    `json:"summary"`
	LearningLanguage  string    `json:"learning_language,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CourseDraftPoint struct {
	ID                  uuid.UUID                 `json:"id"`
	CourseDraftID       uuid.UUID                 `json:"course_draft_id"`
	CourseDraftLessonID uuid.UUID                 `json:"course_draft_lesson_id"`
	PointType           PointType                 `json:"point_type"`
	Status              PointStatus               `json:"status"`
	Title               string                    `json:"title"`
	Description         *string                   `json:"description"`
	TemplateType        *ResearchNodeTemplateType `json:"template_type,omitempty"`
	SelectionState      *ResourceSelectionState   `json:"selection_state,omitempty"`
	ContentID           *uuid.UUID                `json:"content_id,omitempty"`
	ExternalURL         *string                   `json:"external_url,omitempty"`
	ThumbnailURL        *string                   `json:"thumbnail_url,omitempty"`
	PriceType           *PriceType                `json:"price_type,omitempty"`
	RankScore           *float64                  `json:"rank_score,omitempty"`
	OrderIndex          int                       `json:"order_index"`
	CompletedAt         *time.Time                `json:"completed_at,omitempty"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

type CourseDraftPointBlock struct {
	ID                 uuid.UUID             `json:"id"`
	CourseDraftPointID uuid.UUID             `json:"course_draft_point_id"`
	BlockType          ResearchNodeBlockType `json:"block_type"`
	Content            json.RawMessage       `json:"content"`
	OrderIndex         int                   `json:"order_index"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

type CreateSubLessonRequest struct {
	Title      string  `json:"title"`
	Objective  *string `json:"objective"`
	OrderIndex int     `json:"order_index"`
}

type CreateMainLessonRequest struct {
	Title      string  `json:"title"`
	Objective  *string `json:"objective"`
	OrderIndex int     `json:"order_index"`
}

type CreateResearchNodeRequest struct {
	Title        string                   `json:"title"`
	TemplateType ResearchNodeTemplateType `json:"template_type"`
	OrderIndex   int                      `json:"order_index"`
}

type UpdateResearchNodeRequest struct {
	Title        *string                   `json:"title"`
	TemplateType *ResearchNodeTemplateType `json:"template_type"`
	OrderIndex   *int                      `json:"order_index"`
}

type CreateResearchNodeBlockRequest struct {
	BlockType  ResearchNodeBlockType `json:"block_type"`
	Content    json.RawMessage       `json:"content"`
	OrderIndex int                   `json:"order_index"`
}

type UpdateResearchNodeBlockRequest struct {
	Content    json.RawMessage `json:"content"`
	OrderIndex *int            `json:"order_index"`
}

type ReorderResearchNodeBlocksRequest struct {
	BlockIDs []uuid.UUID `json:"block_ids"`
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

type CreateCourseDraftRequest struct {
	SourceQuery               string                                            `json:"source_query" binding:"required,max=300"`
	LearningGoal              *string                                           `json:"learning_goal"`
	CurrentLevel              *string                                           `json:"current_level"`
	DurationWeeks             *int                                              `json:"duration_weeks"`
	StudyHoursPerWeek         *int                                              `json:"study_hours_per_week"`
	PreferredFormat           *string                                           `json:"preferred_format"`
	SkipGeneration            bool                                              `json:"skip_generation"`
	GoalProfileID             *uuid.UUID                                        `json:"-"`
	GoalProfileVersion        *int                                              `json:"-"`
	GoalUserIntent            *string                                           `json:"-"`
	GoalMotivation            *string                                           `json:"-"`
	GoalUsageContext          *string                                           `json:"-"`
	GoalDifficultyLevel       *string                                           `json:"-"`
	GoalTimeHorizon           *string                                           `json:"-"`
	GoalOutputType            *string                                           `json:"-"`
	GoalType                  *string                                           `json:"-"`
	LearningIntent            LearningIntentProfile                             `json:"-"`
	LearningLanguage          string                                            `json:"-"`
	PatternMatchGuidance      string                                            `json:"-"`
	PatternSearchSpecTemplate CurriculumPatternRecommendationSearchSpecTemplate `json:"-"`
}

type DeleteCourseDraftRequest struct {
	ConfirmTitle string `json:"confirm_title" binding:"required"`
}

type CourseDraftPointPreview struct {
	Cost         int `json:"cost"`
	FreeBalance  int `json:"free_balance"`
	PaidBalance  int `json:"paid_balance"`
	TotalBalance int `json:"total_balance"`
}

type UpdateCourseDraftRequest struct {
	Title              *string    `json:"title"`
	Description        *string    `json:"description"`
	LearningGoal       *string    `json:"learning_goal"`
	CurrentLevel       *string    `json:"current_level"`
	DurationWeeks      *int       `json:"duration_weeks"`
	StudyHoursPerWeek  *int       `json:"study_hours_per_week"`
	PreferredFormat    *string    `json:"preferred_format"`
	PlanetTypeID       *uuid.UUID `json:"planet_type_id"`
	PlanetTextureMapID *uuid.UUID `json:"planet_texture_map_id"`
}

type SearchCourseDraftLessonsRequest struct {
	LessonID     *uuid.UUID `json:"lesson_id"`
	MaxPerLesson int        `json:"max_per_lesson"`
}

type AttachDraftLessonResourceRequest struct {
	ContentID      uuid.UUID              `json:"content_id" binding:"required"`
	SelectionState ResourceSelectionState `json:"selection_state"`
}

type UpdateDraftMainLessonRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Objective   *string `json:"objective"`
}

type UpdateDraftLessonRequest struct {
	Title           *string `json:"title"`
	Objective       *string `json:"objective"`
	Summary         *string `json:"summary"`
	DifficultyLevel *string `json:"difficulty_level"`
}

type UpdateDraftDetailMemoRequest struct {
	Note *string `json:"note"`
}

type UpdateDraftLessonJournalRequest struct {
	Observation    *string `json:"observation"`
	Reflection     *string `json:"reflection"`
	NextStep       *string `json:"next_step"`
	CoreConcept    *string `json:"core_concept"`
	MyExplanation  *string `json:"my_explanation"`
	Examples       *string `json:"examples"`
	ConfusedParts  *string `json:"confused_parts"`
	ReferenceLinks *string `json:"reference_links"`
}

type UpdateDraftLessonRecordRequest struct {
	StudyMinutes    *int    `json:"study_minutes"`
	PracticeCount   *int    `json:"practice_count"`
	ConfidenceLevel *int    `json:"confidence_level"`
	ApplicationNote *string `json:"application_note"`
}

type UpdateDraftLessonArtifactRequest struct {
	ArtifactType *string `json:"artifact_type"`
	Title        *string `json:"title"`
	URL          *string `json:"url"`
	Description  *string `json:"description"`
}

type UpdateDraftStructureRequest struct {
	MainLessons []UpdateDraftStructureMainLesson `json:"main_lessons"`
}

type UpdateDraftStructureMainLesson struct {
	ID         uuid.UUID                          `json:"id" binding:"required"`
	OrderIndex int                                `json:"order_index"`
	Lessons    []UpdateDraftStructureLesson       `json:"lessons"`
	Points     []UpdateDraftStructureResearchNode `json:"points"`
}

func (l UpdateDraftStructureMainLesson) StructurePoints() []UpdateDraftStructureResearchNode {
	return l.Points
}

type UpdateDraftStructureLesson struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	OrderIndex int       `json:"order_index"`
}

type UpdateDraftStructureResearchNode struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	OrderIndex int       `json:"order_index"`
}

type DraftDetailMemo struct {
	LessonID  uuid.UUID `json:"lesson_id"`
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DraftLessonJournal struct {
	TargetID       uuid.UUID `json:"target_id"`
	Observation    string    `json:"observation"`
	Reflection     string    `json:"reflection"`
	NextStep       string    `json:"next_step"`
	CoreConcept    string    `json:"core_concept"`
	MyExplanation  string    `json:"my_explanation"`
	Examples       string    `json:"examples"`
	ConfusedParts  string    `json:"confused_parts"`
	ReferenceLinks string    `json:"reference_links"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DraftLessonRecord struct {
	TargetID        uuid.UUID `json:"target_id"`
	StudyMinutes    int       `json:"study_minutes"`
	PracticeCount   int       `json:"practice_count"`
	ConfidenceLevel int       `json:"confidence_level"`
	ApplicationNote string    `json:"application_note"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DraftLessonArtifact struct {
	TargetID     uuid.UUID `json:"target_id"`
	ArtifactType string    `json:"artifact_type"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Description  string    `json:"description"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DraftAggregate struct {
	Draft   CourseDraft       `json:"draft"`
	Lessons []DraftLessonTree `json:"lessons"`
}

type DraftPointAggregate struct {
	Point         CourseDraftPoint        `json:"point"`
	Blocks        []CourseDraftPointBlock `json:"blocks,omitempty"`
	JournalEntry  *DraftJournalEntry      `json:"journal_entry,omitempty"`
	RecordEntry   *DraftRecordEntry       `json:"record_entry,omitempty"`
	ArtifactEntry *DraftArtifactEntry     `json:"artifact_entry,omitempty"`
}

type DraftLessonTree struct {
	Lesson     CourseDraftLesson     `json:"lesson"`
	Points     []DraftPointAggregate `json:"points,omitempty"`
	SubLessons []DraftLessonTree     `json:"sub_lessons,omitempty"`
}
