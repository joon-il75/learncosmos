package auth

import "github.com/golang-jwt/jwt/v5"

type Role string

const (
	RoleLearner    Role = "learner"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "super_admin"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderKakao  Provider = "kakao"
	ProviderNaver  Provider = "naver"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}

type OAuthUserInfo struct {
	Provider   Provider
	ProviderID string
	Email      string
	Nickname   string
	AvatarURL  string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	ID                     string `json:"id"`
	Email                  string `json:"email,omitempty"`
	DisplayID              string `json:"display_id,omitempty"`
	Nickname               string `json:"nickname,omitempty"`
	AvatarURL              string `json:"avatar_url,omitempty"`
	Role                   Role   `json:"role"`
	PremiumAccess          bool   `json:"premium_access"`
	Provider               string `json:"provider,omitempty"`
	CreatedAt              string `json:"created_at,omitempty"`
	FreePoints             int    `json:"free_points"`
	PaidPoints             int    `json:"paid_points"`
	TotalPoints            int    `json:"total_points"`
	TermsAgreed            bool   `json:"terms_agreed"`
	PrivacyAgreed          bool   `json:"privacy_agreed"`
	RequiredConsentPending bool   `json:"required_consent_pending"`
	UILocale               string `json:"ui_locale,omitempty"`
	LearningLanguage       string `json:"learning_language,omitempty"`
	LanguageSetupRequired  bool   `json:"language_setup_required"`
	IsNewSocialSignup      bool   `json:"is_new_social_signup,omitempty"`
}

type PolicyDocumentSummary struct {
	ID                string `json:"id"`
	SetID             string `json:"set_id,omitempty"`
	Type              string `json:"type"`
	Title             string `json:"title"`
	Version           int    `json:"version"`
	Locale            string `json:"locale,omitempty"`
	TranslationStatus string `json:"translation_status,omitempty"`
	RequestedLocale   string `json:"requested_locale,omitempty"`
	FallbackUsed      bool   `json:"fallback_used,omitempty"`
	EffectiveAt       string `json:"effective_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

type PolicyDocument struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Version     int    `json:"version"`
	Content     string `json:"content"`
	Required    bool   `json:"required"`
	Active      bool   `json:"active"`
	EffectiveAt string `json:"effective_at,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ConsentStatus struct {
	TermsAgreed            bool                    `json:"terms_agreed"`
	PrivacyAgreed          bool                    `json:"privacy_agreed"`
	TermsAgreedAt          string                  `json:"terms_agreed_at,omitempty"`
	PrivacyAgreedAt        string                  `json:"privacy_agreed_at,omitempty"`
	RequiredConsentPending bool                    `json:"required_consent_pending"`
	RequiredDocuments      []PolicyDocumentSummary `json:"required_documents,omitempty"`
}

type UserAISettings struct {
	Mode                 string  `json:"mode"`
	Provider             string  `json:"provider"`
	HasAPIKey            bool    `json:"has_api_key"`
	IsEnabled            bool    `json:"is_enabled"`
	EndpointURL          *string `json:"endpoint_url,omitempty"`
	UpdatedAt            string  `json:"updated_at,omitempty"`
	LastValidationStatus string  `json:"last_validation_status,omitempty"`
	LastValidatedAt      string  `json:"last_validated_at,omitempty"`
	LastValidationError  string  `json:"last_validation_error,omitempty"`
	LastFailedAt         string  `json:"last_failed_at,omitempty"`
	NextRetryAt          string  `json:"next_retry_at,omitempty"`
}

type UserAIUsageEvent struct {
	ID               string   `json:"id"`
	Source           string   `json:"source"`
	Provider         string   `json:"provider"`
	Model            string   `json:"model,omitempty"`
	Feature          string   `json:"feature"`
	BillingStatus    string   `json:"billing_status,omitempty"`
	InputTokens      *int     `json:"input_tokens,omitempty"`
	OutputTokens     *int     `json:"output_tokens,omitempty"`
	EstimatedCostUSD *float64 `json:"estimated_cost_usd,omitempty"`
	Success          bool     `json:"success"`
	ErrorCode        string   `json:"error_code,omitempty"`
	CreatedAt        string   `json:"created_at"`
}

type UserAIUsageSummary struct {
	TotalTokens                  int     `json:"total_tokens"`
	AverageDailyTokens           float64 `json:"average_daily_tokens"`
	TotalEstimatedCostUSD        float64 `json:"total_estimated_cost_usd"`
	AverageDailyEstimatedCostUSD float64 `json:"average_daily_estimated_cost_usd"`
}

type UserPointTransaction struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Amount        int            `json:"amount"`
	Feature       string         `json:"feature,omitempty"`
	ReferenceType string         `json:"reference_type,omitempty"`
	ReferenceID   *string        `json:"reference_id,omitempty"`
	Description   string         `json:"description,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	CreatedAt     string         `json:"created_at"`
}

type UserPointUsageSummary struct {
	GrantedPoints   int `json:"granted_points"`
	PurchasedPoints int `json:"purchased_points"`
	UsedPoints      int `json:"used_points"`
	RefundedPoints  int `json:"refunded_points"`
	NetChange       int `json:"net_change"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserInfo `json:"user"`
}
