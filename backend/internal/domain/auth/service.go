package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	accessTokenTTL     = 15 * time.Minute
	refreshTokenTTL    = 30 * 24 * time.Hour
	refreshTokenCookie = "refresh_token"
	redisRefreshPrefix = "refresh:"
	redisStatePrefix   = "oauth_state:"
	initialFreePoints  = 30
)

type User struct {
	ID                       string     `db:"id"`
	Email                    string     `db:"email"`
	Nickname                 string     `db:"nickname"`
	Role                     Role       `db:"role"`
	PremiumAccess            bool       `db:"premium_access"`
	AvatarURL                string     `db:"avatar_url"`
	DisplayID                *string    `db:"display_id"`
	Status                   string     `db:"status"`
	Provider                 string     `db:"provider"`
	UILocale                 string     `db:"ui_locale"`
	LearningLanguage         string     `db:"learning_language"`
	LanguageSetupCompletedAt *time.Time `db:"language_setup_completed_at"`
	LastLoginAt              *time.Time `db:"last_login_at"`
	WithdrawnAt              *time.Time `db:"withdrawn_at"`
	ReactivatedAt            *time.Time `db:"reactivated_at"`
	CreatedAt                time.Time  `db:"created_at"`
}

type SocialAccount struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	Provider   Provider  `db:"provider"`
	ProviderID string    `db:"provider_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type UserRepository interface {
	FindBySocialAccount(ctx context.Context, provider Provider, providerID string) (*User, error)
	Create(ctx context.Context, user *User) error
	CreateSocialAccount(ctx context.Context, sa *SocialAccount) error
	GrantFreePoints(ctx context.Context, userID string, amount int) error
	FindByID(ctx context.Context, id string) (*User, error)
	GetPointBalances(ctx context.Context, userID string) (freeBalance int, paidBalance int, err error)
	GetPointSetting(ctx context.Context, key string) (int, error)
	GetConsentStatus(ctx context.Context, userID string, locale string) (*ConsentStatus, error)
	RecordRequiredConsents(ctx context.Context, userID, locale, ipAddress, userAgent string) error
	GetUserAISettings(ctx context.Context, userID string) (*UserAISettings, error)
	ListUserAIUsageEvents(ctx context.Context, userID string, startAt, endAt time.Time, limit, offset int) ([]UserAIUsageEvent, int, UserAIUsageSummary, error)
	ListUserPointTransactions(ctx context.Context, userID string, startAt, endAt time.Time, limit, offset int) ([]UserPointTransaction, int, UserPointUsageSummary, error)
	UpsertUserAISettings(ctx context.Context, userID, provider, apiKey string, endpointURL *string, isEnabled bool) error
	UpdateUserAISettingsEnabled(ctx context.Context, userID string, isEnabled bool) error
	DeleteUserAISettings(ctx context.Context, userID string) error
	RecordUserAIValidation(ctx context.Context, userID, provider string, valid bool, message string) error
	DisplayIDExists(ctx context.Context, displayID string) (bool, error)
	UpdateNickname(ctx context.Context, userID, nickname string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	EmailExists(ctx context.Context, email, excludeUserID string) (bool, error)
	UpdateLanguagePreferences(ctx context.Context, userID, uiLocale, learningLanguage string) error
	UpdateAvatarURL(ctx context.Context, userID, avatarURL string) error
	// UpdateOnLogin: 재로그인 시 닉네임·아바타 갱신, 이메일은 기존 값이 있으면 보호
	UpdateOnLogin(ctx context.Context, userID, email, nickname, avatarURL string) error
	ReactivateUser(ctx context.Context, userID, email, nickname, avatarURL string) error
	WithdrawUser(ctx context.Context, userID string) error
}

type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

type OAuthConfig struct {
	Google ProviderConfig
	Kakao  ProviderConfig
	Naver  ProviderConfig
}

type Service struct {
	repo      UserRepository
	redis     *redis.Client
	jwtSecret []byte
	oauth     OAuthConfig
}

func NewService(repo UserRepository, redisClient *redis.Client, jwtSecret string, oauth OAuthConfig) *Service {
	return &Service{
		repo:      repo,
		redis:     redisClient,
		jwtSecret: []byte(jwtSecret),
		oauth:     oauth,
	}
}

func (s *Service) WithdrawUser(ctx context.Context, userID string) error {
	if err := s.repo.WithdrawUser(ctx, userID); err != nil {
		return err
	}
	return s.Logout(ctx, userID)
}
