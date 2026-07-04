package admin

import (
	"database/sql"
	"encoding/json"
	"time"
)

type recommendationDebugPreviewRequest struct {
	DraftID      string `json:"draft_id"`
	LessonID     string `json:"lesson_id"`
	MaxPerLesson int    `json:"max_per_lesson"`
}

type recommendationDebugDirectPreviewRequest struct {
	OwnerUserID     string `json:"owner_user_id"`
	SourceQuery     string `json:"source_query"`
	GroupTitle      string `json:"group_title"`
	GroupObjective  string `json:"group_objective"`
	LessonTitle     string `json:"lesson_title"`
	LessonObjective string `json:"lesson_objective"`
	LessonSummary   string `json:"lesson_summary"`
	MaxPerLesson    int    `json:"max_per_lesson"`
}

type recommendationDebugScenarioRequest struct {
	CourseTitle       string `json:"course_title"`
	InitialUserIntent string `json:"initial_user_intent"`
	Notes             string `json:"notes"`
}

type recommendationDebugGoalMessageRequest struct {
	Message string `json:"message"`
}

type recommendationDebugGoalConfirmRequest struct {
	ConfirmedGoal string `json:"confirmed_goal"`
}

type recommendationDebugScenarioCompareRequest struct {
	MaxPerLesson        int    `json:"max_per_lesson"`
	LessonID            string `json:"lesson_id"`
	OrderIndex          *int   `json:"order_index"`
	RecommendationQuery string `json:"recommendation_query"`
}

type recommendationDebugExternalCandidateSaveRequest struct {
	LessonID          string `json:"lesson_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	URL               string `json:"url"`
	Source            string `json:"source"`
	ExternalContentID string `json:"external_content_id"`
	ThumbnailURL      string `json:"thumbnail_url"`
	Author            string `json:"author"`
	Language          string `json:"language"`
}

type recommendationDebugLabelRequest struct {
	RunID           string          `json:"run_id"`
	CandidateKey    string          `json:"candidate_key"`
	ContentID       string          `json:"content_id"`
	URL             string          `json:"url"`
	Label           string          `json:"label"`
	Note            string          `json:"note"`
	FeatureSnapshot json.RawMessage `json:"feature_snapshot"`
	BaselineRank    *int            `json:"baseline_rank"`
	ShadowRank      *int            `json:"shadow_rank"`
}

type recommendationDebugShadowProbeRequest struct {
	Endpoint  string `json:"endpoint"`
	Model     string `json:"model"`
	ProbeText string `json:"probe_text"`
}

type recommendationRolloutTransitionRequest struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type recommendationRolloutMetricSnapshotRequest struct {
	WindowHours int `json:"window_hours"`
}

type recommendationRolloutRankerRequest struct {
	RankerModelVersion   string `json:"ranker_model_version"`
	FeatureSchemaVersion string `json:"feature_schema_version"`
}

type recommendationRankerArtifactStatus struct {
	Ready                 bool     `json:"ready"`
	Reason                string   `json:"reason"`
	Provider              string   `json:"provider"`
	RankerModelVersion    string   `json:"ranker_model_version"`
	FeatureSchemaVersion  string   `json:"feature_schema_version"`
	ExpectedModelVersion  string   `json:"expected_model_version"`
	ExpectedSchemaVersion string   `json:"expected_schema_version"`
	ConfigSource          string   `json:"config_source"`
	ModelPath             string   `json:"model_path"`
	ManifestPath          string   `json:"manifest_path"`
	ModelFileExists       bool     `json:"model_file_exists"`
	ManifestExists        bool     `json:"manifest_exists"`
	ModelSizeBytes        int64    `json:"model_size_bytes"`
	ModelUpdatedAt        string   `json:"model_updated_at"`
	ManifestVersionMatch  bool     `json:"manifest_version_match"`
	ManifestSchemaMatch   bool     `json:"manifest_schema_match"`
	CheckedPaths          []string `json:"checked_paths"`
}

type recommendationLearnerSignalReadiness struct {
	Required          bool
	Ready             bool
	Reason            string
	ExposureCount     int64
	ReplacementCount  int64
	BrokenLinkCount   int64
	WrongContentCount int64
	MinimumExposure   int64
	SelectionRate     float64
	BrokenLinkRate    float64
	WrongContentRate  float64
	MinSelectionRate  float64
	MaxBrokenLinkRate float64
	MaxWrongRate      float64
}

type recommendationDebugGeneratedLessonsSnapshot struct {
	CourseTitle       string                                       `json:"course_title"`
	InitialUserIntent string                                       `json:"initial_user_intent"`
	ConfirmedGoal     string                                       `json:"confirmed_goal"`
	DraftTitle        string                                       `json:"draft_title"`
	DraftDescription  string                                       `json:"draft_description"`
	Lessons           []recommendationDebugGeneratedLessonSnapshot `json:"lessons"`
}

type recommendationDebugGeneratedLessonSnapshot struct {
	LessonID   string `json:"lesson_id"`
	Title      string `json:"title"`
	Objective  string `json:"objective"`
	OrderIndex int    `json:"order_index"`
	SourceType string `json:"source_type"`
	LessonRole string `json:"lesson_role"`
}

type recommendationRolloutStateRow struct {
	ID                     string       `json:"id"`
	Mode                   string       `json:"mode"`
	TrafficPercent         int          `json:"traffic_percent"`
	EmbeddingProvider      string       `json:"embedding_provider"`
	EmbeddingModel         string       `json:"embedding_model"`
	EmbeddingDimension     int          `json:"embedding_dimension"`
	RankerModelVersion     string       `json:"ranker_model_version"`
	FeatureSchemaVersion   string       `json:"feature_schema_version"`
	QualityGateStatus      string       `json:"quality_gate_status"`
	QualityGateSnapshot    string       `json:"quality_gate_snapshot"`
	RollbackPolicySnapshot string       `json:"rollback_policy_snapshot"`
	ApprovedBy             string       `json:"approved_by"`
	ApprovedAt             sql.NullTime `json:"approved_at"`
	ActivatedAt            sql.NullTime `json:"activated_at"`
	RolledBackAt           sql.NullTime `json:"rolled_back_at"`
	RollbackReason         string       `json:"rollback_reason"`
	Active                 bool         `json:"active"`
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
}

type recommendationDebugScenarioListRow struct {
	ID                string    `json:"id"`
	CreatedBy         string    `json:"created_by"`
	CourseTitle       string    `json:"course_title"`
	InitialUserIntent string    `json:"initial_user_intent"`
	Status            string    `json:"status"`
	Notes             string    `json:"notes"`
	RunCount          int       `json:"run_count"`
	LabelCount        int       `json:"label_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type recommendationDebugScenarioDetailRow struct {
	ID                       string                        `json:"id"`
	CreatedBy                string                        `json:"created_by"`
	CourseTitle              string                        `json:"course_title"`
	InitialUserIntent        string                        `json:"initial_user_intent"`
	Status                   string                        `json:"status"`
	GoalProfileSnapshot      string                        `json:"goal_profile_snapshot"`
	GeneratedLessonsSnapshot string                        `json:"generated_lessons_snapshot"`
	Notes                    string                        `json:"notes"`
	CreatedAt                time.Time                     `json:"created_at"`
	UpdatedAt                time.Time                     `json:"updated_at"`
	Runs                     []recommendationDebugRunRow   `json:"runs"`
	Labels                   []recommendationDebugLabelRow `json:"labels"`
}

type recommendationDebugRunRow struct {
	ID                     string        `json:"id"`
	ScenarioID             string        `json:"scenario_id"`
	CreatedBy              string        `json:"created_by"`
	RunType                string        `json:"run_type"`
	RequestSnapshot        string        `json:"request_snapshot"`
	BaselineResultSnapshot string        `json:"baseline_result_snapshot"`
	ShadowResultSnapshot   string        `json:"shadow_result_snapshot"`
	FeatureSnapshot        string        `json:"feature_snapshot"`
	MetricsSnapshot        string        `json:"metrics_snapshot"`
	ProviderSnapshot       string        `json:"provider_snapshot"`
	LatencyMS              sql.NullInt32 `json:"latency_ms"`
	CreatedAt              time.Time     `json:"created_at"`
}

type recommendationDebugLabelRow struct {
	ID              string        `json:"id"`
	ScenarioID      string        `json:"scenario_id"`
	RunID           string        `json:"run_id"`
	CreatedBy       string        `json:"created_by"`
	CandidateKey    string        `json:"candidate_key"`
	ContentID       string        `json:"content_id"`
	URL             string        `json:"url"`
	Label           string        `json:"label"`
	Note            string        `json:"note"`
	FeatureSnapshot string        `json:"feature_snapshot"`
	BaselineRank    sql.NullInt32 `json:"baseline_rank"`
	ShadowRank      sql.NullInt32 `json:"shadow_rank"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type recommendationDebugLabelExportRow struct {
	LabelID                string        `json:"label_id"`
	ScenarioID             string        `json:"scenario_id"`
	CourseTitle            string        `json:"course_title"`
	InitialUserIntent      string        `json:"initial_user_intent"`
	ScenarioStatus         string        `json:"scenario_status"`
	RunID                  string        `json:"run_id"`
	RunCreatedAt           sql.NullTime  `json:"run_created_at"`
	CandidateKey           string        `json:"candidate_key"`
	ContentID              string        `json:"content_id"`
	URL                    string        `json:"url"`
	Label                  string        `json:"label"`
	Note                   string        `json:"note"`
	FeatureSnapshot        string        `json:"feature_snapshot"`
	BaselineRank           sql.NullInt32 `json:"baseline_rank"`
	ShadowRank             sql.NullInt32 `json:"shadow_rank"`
	RequestSnapshot        string        `json:"request_snapshot"`
	BaselineResultSnapshot string        `json:"baseline_result_snapshot"`
	ShadowResultSnapshot   string        `json:"shadow_result_snapshot"`
	MetricsSnapshot        string        `json:"metrics_snapshot"`
	UpdatedAt              time.Time     `json:"updated_at"`
}

type recommendationRolloutEventRow struct {
	ID             string    `json:"id"`
	RolloutStateID string    `json:"rollout_state_id"`
	CreatedBy      string    `json:"created_by"`
	EventType      string    `json:"event_type"`
	FromMode       string    `json:"from_mode"`
	ToMode         string    `json:"to_mode"`
	Payload        string    `json:"payload"`
	CreatedAt      time.Time `json:"created_at"`
}

type recommendationRolloutMetricRow struct {
	ID                  string          `json:"id"`
	RolloutStateID      string          `json:"rollout_state_id"`
	WindowStartedAt     time.Time       `json:"window_started_at"`
	WindowEndedAt       time.Time       `json:"window_ended_at"`
	RecommendationCount int             `json:"recommendation_count"`
	SelectionRate       sql.NullFloat64 `json:"selection_rate"`
	BrokenLinkRate      sql.NullFloat64 `json:"broken_link_rate"`
	WrongContentRate    sql.NullFloat64 `json:"wrong_content_rate"`
	FallbackRate        sql.NullFloat64 `json:"fallback_rate"`
	P95LatencyMS        sql.NullInt32   `json:"p95_latency_ms"`
	Top5Overlap         sql.NullFloat64 `json:"top5_overlap"`
	GoodFitRate         sql.NullFloat64 `json:"good_fit_rate"`
	IrrelevantRate      sql.NullFloat64 `json:"irrelevant_rate"`
	DuplicateRate       sql.NullFloat64 `json:"duplicate_rate"`
	MetricsSnapshot     string          `json:"metrics_snapshot"`
	CreatedAt           time.Time       `json:"created_at"`
}

const recommendationDictionaryVersion = "builtin-v1"
