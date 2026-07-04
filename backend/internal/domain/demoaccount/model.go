package demoaccount

import "time"

const (
	DefaultEmailPrefix = "demo"
	DefaultEmailDomain = "learnweavr.local"
	DefaultLabelPrefix = "demo"
	DefaultUILocale    = "ko"
	DefaultPoints      = 300
	DefaultAvatarURL   = "/api/v1/users/avatars/demo-default-profile.png"
)

type CreateBatchRequest struct {
	Count            int
	StartNumber      int
	EmailPrefix      string
	EmailDomain      string
	LabelPrefix      string
	UILocale         string
	LearningLanguage string
	InitialPoints    int
	ExpiresAt        *time.Time
	AssignedTo       string
	AssignmentNote   string
	CreatedBy        *string
	CreatedByActor   string
}

type CreatedAccount struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Label     string `json:"label"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	LoginCode string `json:"login_code"`
}

type ListOptions struct {
	Query  string
	Status string
	Page   int
	Limit  int
}

type ListResult struct {
	Items []AccountSummary `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type AccountSummary struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Label          string  `json:"label"`
	Email          string  `json:"email"`
	Nickname       string  `json:"nickname"`
	AssignedTo     string  `json:"assigned_to"`
	AssignmentNote string  `json:"assignment_note"`
	LoginCode      *string `json:"login_code,omitempty"`
	Status         string  `json:"status"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
	ActivatedAt    *string `json:"activated_at,omitempty"`
	LastUsedAt     *string `json:"last_used_at,omitempty"`
	DisabledAt     *string `json:"disabled_at,omitempty"`
	FreePoints     int     `json:"free_points"`
	PaidPoints     int     `json:"paid_points"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type UpdateRequest struct {
	AssignedTo     *string
	AssignmentNote *string
	ExpiresAt      **time.Time
	Disabled       *bool
	DisabledReason *string
}

type PurgePreview struct {
	DemoAccountID               string   `json:"demo_account_id"`
	UserID                      string   `json:"user_id"`
	Label                       string   `json:"label"`
	CourseDrafts                int      `json:"course_drafts"`
	Courses                     int      `json:"courses"`
	GoalProfiles                int      `json:"goal_profiles"`
	GoalRevisionLogs            int      `json:"goal_revision_logs"`
	CoursePoints                int      `json:"course_points"`
	CourseDraftPoints           int      `json:"course_draft_points"`
	Attachments                 int      `json:"attachments"`
	ObjectStorageFiles          int      `json:"object_storage_files"`
	AIUsageEvents               int      `json:"ai_usage_events"`
	RecommendationEvents        int      `json:"recommendation_events"`
	RecommendationLearnerEvents int      `json:"recommendation_learner_events"`
	Contents                    int      `json:"contents"`
	SafeDeletableContents       int      `json:"safe_deletable_contents"`
	ObjectKeys                  []string `json:"object_keys,omitempty"`
}

type PurgeResult struct {
	PurgePreview
	DeletedObjectKeys       []string `json:"deleted_object_keys"`
	ObjectDeleteFailedKeys  []string `json:"object_delete_failed_keys"`
	ObjectDeleteFailureText []string `json:"object_delete_failure_text"`
	AccountDeleted          bool     `json:"account_deleted"`
}

type AuthenticatedUser struct {
	DemoAccountID            string
	ID                       string
	Email                    string
	Nickname                 string
	Role                     string
	PremiumAccess            bool
	AvatarURL                string
	DisplayID                *string
	Status                   string
	UILocale                 string
	LearningLanguage         string
	LanguageSetupCompletedAt *time.Time
	LastLoginAt              *time.Time
	WithdrawnAt              *time.Time
	ReactivatedAt            *time.Time
	CreatedAt                time.Time
}
