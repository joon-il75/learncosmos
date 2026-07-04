package curriculum

type DraftStatus string

const (
	DraftStatusDraft     DraftStatus = "draft"
	DraftStatusLearning  DraftStatus = "learning"
	DraftStatusConfirmed DraftStatus = "confirmed"
	DraftStatusArchived  DraftStatus = "archived"
)

type CourseStatus string

const (
	CourseStatusActive   CourseStatus = "active"
	CourseStatusArchived CourseStatus = "archived"
)

type PlanetStatus string

const (
	PlanetStatusDraft     PlanetStatus = "draft"
	PlanetStatusReady     PlanetStatus = "ready"
	PlanetStatusLearning  PlanetStatus = "learning"
	PlanetStatusCompleted PlanetStatus = "completed"
)

type LessonRuntimeStatus string

const (
	LessonRuntimeInProgress LessonRuntimeStatus = "in_progress"
	LessonRuntimeCompleted  LessonRuntimeStatus = "completed"
)

type PointStatus string

const (
	PointStatusDraft     PointStatus = "draft"
	PointStatusReady     PointStatus = "ready"
	PointStatusLearning  PointStatus = "learning"
	PointStatusCompleted PointStatus = "completed"
)

type PointType string

const (
	PointTypeExploration PointType = "exploration"
	PointTypeResearch    PointType = "research"
)

type PointQuestionStatus string

const (
	PointQuestionStatusPending  PointQuestionStatus = "pending"
	PointQuestionStatusAnswered PointQuestionStatus = "answered"
)

type PointQuestionType string

const (
	PointQuestionTypeReflection    PointQuestionType = "reflection"
	PointQuestionTypeApplication   PointQuestionType = "application"
	PointQuestionTypeGoalAlignment PointQuestionType = "goal_alignment"
)

type LessonRole string

const (
	LessonRoleCore        LessonRole = "core"
	LessonRoleSupport     LessonRole = "support"
	LessonRoleDemo        LessonRole = "demo"
	LessonRoleReference   LessonRole = "reference"
	LessonRoleInspiration LessonRole = "inspiration"
)

type LessonSourceType string

const (
	LessonSourceAI          LessonSourceType = "ai_generated"
	LessonSourceManual      LessonSourceType = "manual"
	LessonSourceRecommended LessonSourceType = "recommended"
)

type ResourceType string

const (
	ResourceTypeContent     ResourceType = "content"
	ResourceTypeExternal    ResourceType = "external"
	ResourceTypeCreator     ResourceType = "creator"
	ResourceTypeAIGenerated ResourceType = "ai_generated"
)

type ResourceSelectionState string

const (
	ResourceSelectionCandidate ResourceSelectionState = "candidate"
	ResourceSelectionSelected  ResourceSelectionState = "selected"
	ResourceSelectionRejected  ResourceSelectionState = "rejected"
)

type PriceType string

const (
	PriceTypeFree  PriceType = "free"
	PriceTypePaid  PriceType = "paid"
	PriceTypeOwned PriceType = "owned"
)

type RecommendationEventType string

const (
	EventCurriculumGenerated  RecommendationEventType = "curriculum_generated"
	EventCurriculumDeleted    RecommendationEventType = "curriculum_deleted"
	EventLessonRecommended    RecommendationEventType = "lesson_recommended"
	EventLessonSelected       RecommendationEventType = "lesson_selected"
	EventLessonRejected       RecommendationEventType = "lesson_rejected"
	EventLessonAddedManually  RecommendationEventType = "lesson_added_manually"
	EventLessonGeneratedByAI  RecommendationEventType = "lesson_generated_by_ai"
	EventCurriculumConfirmed  RecommendationEventType = "curriculum_confirmed"
	EventCurriculumStarted    RecommendationEventType = "curriculum_started"
	EventCurriculumEdited     RecommendationEventType = "curriculum_edited"
	EventCurriculumEditedLive RecommendationEventType = "curriculum_edited_after_start"
	EventCurriculumCompleted  RecommendationEventType = "curriculum_completed"
)

type ResearchNodeTemplateType string

const (
	ResearchNodeTemplateConceptSummary   ResearchNodeTemplateType = "concept_summary"
	ResearchNodeTemplatePracticeStrategy ResearchNodeTemplateType = "practice_strategy"
	ResearchNodeTemplateProblemSolving   ResearchNodeTemplateType = "problem_solving"
	ResearchNodeTemplateFreeResearch     ResearchNodeTemplateType = "free_research"
)

type ResearchNodeBlockType string

const (
	ResearchNodeBlockTypeText  ResearchNodeBlockType = "text"
	ResearchNodeBlockTypeImage ResearchNodeBlockType = "image"
	ResearchNodeBlockTypeLink  ResearchNodeBlockType = "link"
)
