package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

type Config struct {
	Port                                  string
	Env                                   string
	DatabaseURL                           string
	RedisURL                              string
	JWTSecret                             string
	ContentHealthCheckEnabled             bool
	LessonSearchPrerunSchedulerEnabled    bool
	LessonSearchPrerunIntervalHours       int
	LessonSearchPrerunStartDelayMinutes   int
	LessonSearchPrerunTopN                int
	LessonSearchPrerunRiskTopN            int
	LessonSearchPrerunTimeoutSeconds      int
	LessonSearchPrerunProvider            string
	LessonSearchPrerunReportRetentionDays int
	LessonSearchPrerunRiskFixtureIDs      []string

	OAuth OAuthConfig

	SuperAdminID         string
	SuperAdminPWHash     string
	SuperAdminTOTPSecret string
	AdminAllowedIPs      []string

	DemoLoginEnabled      bool
	DemoLoginSecret       string
	DemoLoginAllowedHosts []string

	GoalChatSinglePromptFallbackEnabled bool
	GoalChatMode                        string
	GoalChatFastPathEnabled             bool

	LLMWorkerEnabled                bool
	LLMWorkerCount                  int
	LLMProviderOpenAIMaxConcurrency int
	LLMProviderGoogleMaxConcurrency int
	LLMProviderBYOKMaxConcurrency   int
	CourseGenerationAsyncEnabled    bool

	OpenAIAPIKey  string
	YouTubeAPIKey string

	NCPObjectStorageAccessKey string
	NCPObjectStorageSecretKey string
	NCPObjectStorageBucket    string
	NCPObjectStorageEndpoint  string
	NCPObjectStorageRegion    string
	EmbeddingGemmaEndpoint    string
	EncryptionKey             string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:                                  getEnv("PORT", "8080"),
		Env:                                   getEnv("APP_ENV", "development"),
		DatabaseURL:                           os.Getenv("DATABASE_URL"),
		RedisURL:                              getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:                             os.Getenv("JWT_SECRET"),
		ContentHealthCheckEnabled:             strings.EqualFold(getEnv("CONTENT_HEALTH_CHECK_ENABLED", "false"), "true"),
		LessonSearchPrerunSchedulerEnabled:    strings.EqualFold(getEnv("LESSON_SEARCH_PRERUN_SCHEDULER_ENABLED", "false"), "true"),
		LessonSearchPrerunIntervalHours:       getEnvInt("LESSON_SEARCH_PRERUN_INTERVAL_HOURS", 24),
		LessonSearchPrerunStartDelayMinutes:   getEnvInt("LESSON_SEARCH_PRERUN_START_DELAY_MINUTES", 10),
		LessonSearchPrerunTopN:                getEnvInt("LESSON_SEARCH_PRERUN_TOP_N", 5),
		LessonSearchPrerunRiskTopN:            getEnvInt("LESSON_SEARCH_PRERUN_RISK_TOP_N", 10),
		LessonSearchPrerunTimeoutSeconds:      getEnvInt("LESSON_SEARCH_PRERUN_TIMEOUT_SECONDS", 60),
		LessonSearchPrerunProvider:            getEnv("LESSON_SEARCH_PRERUN_PROVIDER", "all"),
		LessonSearchPrerunReportRetentionDays: getEnvInt("LESSON_SEARCH_PRERUN_REPORT_RETENTION_DAYS", 30),
		LessonSearchPrerunRiskFixtureIDs:      splitAndTrim(getEnv("LESSON_SEARCH_PRERUN_RISK_FIXTURE_IDS", "naver-english-score-exam-article youtube-leathercraft-wallet naver-watercolor-landscape-postcard youtube-beginner-5k-running youtube-shortform-video-publish")),

		OAuth: OAuthConfig{
			Google: ProviderConfig{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
				AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:     "https://oauth2.googleapis.com/token",
				UserInfoURL:  "https://openidconnect.googleapis.com/v1/userinfo",
			},
			Kakao: ProviderConfig{
				ClientID:     os.Getenv("KAKAO_CLIENT_ID"),
				ClientSecret: os.Getenv("KAKAO_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("KAKAO_REDIRECT_URL"),
				AuthURL:      "https://kauth.kakao.com/oauth/authorize",
				TokenURL:     "https://kauth.kakao.com/oauth/token",
				UserInfoURL:  "https://kapi.kakao.com/v2/user/me",
			},
			Naver: ProviderConfig{
				ClientID:     os.Getenv("NAVER_CLIENT_ID"),
				ClientSecret: os.Getenv("NAVER_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("NAVER_REDIRECT_URL"),
				AuthURL:      "https://nid.naver.com/oauth2.0/authorize",
				TokenURL:     "https://nid.naver.com/oauth2.0/token",
				UserInfoURL:  "https://openapi.naver.com/v1/nid/me",
			},
		},

		SuperAdminID:         os.Getenv("SUPER_ADMIN_ID"),
		SuperAdminPWHash:     os.Getenv("SUPER_ADMIN_PW_HASH"),
		SuperAdminTOTPSecret: os.Getenv("SUPER_ADMIN_TOTP_SECRET"),

		DemoLoginEnabled: strings.EqualFold(getEnv("DEMO_LOGIN_ENABLED", "false"), "true"),
		DemoLoginSecret:  os.Getenv("DEMO_LOGIN_SECRET"),

		GoalChatSinglePromptFallbackEnabled: strings.EqualFold(getEnv("GOAL_CHAT_SINGLE_PROMPT_FALLBACK_ENABLED", "false"), "true"),
		GoalChatMode:                        getEnv("GOAL_CHAT_MODE", "two_step"),
		GoalChatFastPathEnabled:             strings.EqualFold(getEnv("GOAL_CHAT_FAST_PATH_ENABLED", "false"), "true"),

		LLMWorkerEnabled:                strings.EqualFold(getEnv("LLM_WORKER_ENABLED", "false"), "true"),
		LLMWorkerCount:                  getEnvInt("LLM_WORKER_COUNT", 6),
		LLMProviderOpenAIMaxConcurrency: getEnvInt("LLM_PROVIDER_OPENAI_MAX_CONCURRENCY", 4),
		LLMProviderGoogleMaxConcurrency: getEnvInt("LLM_PROVIDER_GOOGLE_MAX_CONCURRENCY", 2),
		LLMProviderBYOKMaxConcurrency:   getEnvInt("LLM_PROVIDER_BYOK_MAX_CONCURRENCY", 4),
		CourseGenerationAsyncEnabled:    strings.EqualFold(getEnv("COURSE_GENERATION_ASYNC_ENABLED", "false"), "true"),

		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		YouTubeAPIKey: os.Getenv("YOUTUBE_API_KEY"),

		NCPObjectStorageAccessKey: os.Getenv("NCP_OBJECT_STORAGE_ACCESS_KEY"),
		NCPObjectStorageSecretKey: os.Getenv("NCP_OBJECT_STORAGE_SECRET_KEY"),
		NCPObjectStorageBucket:    os.Getenv("NCP_OBJECT_STORAGE_BUCKET"),
		NCPObjectStorageEndpoint:  getEnv("NCP_OBJECT_STORAGE_ENDPOINT", "https://kr.object.ncloudstorage.com"),
		NCPObjectStorageRegion:    getEnv("NCP_OBJECT_STORAGE_REGION", "kr-standard"),
		EmbeddingGemmaEndpoint:    os.Getenv("EMBEDDING_GEMMA_ENDPOINT"),
		EncryptionKey:             os.Getenv("ENCRYPTION_KEY"),
	}

	if ips := os.Getenv("ADMIN_ALLOWED_IPS"); ips != "" {
		cfg.AdminAllowedIPs = strings.Split(ips, ",")
	}
	if hosts := os.Getenv("DEMO_LOGIN_ALLOWED_HOSTS"); hosts != "" {
		cfg.DemoLoginAllowedHosts = splitAndTrim(hosts)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	if c.Env != "production" {
		return nil
	}

	var missing []string
	required := map[string]string{
		"DATABASE_URL":                  c.DatabaseURL,
		"REDIS_URL":                     os.Getenv("REDIS_URL"),
		"JWT_SECRET":                    c.JWTSecret,
		"FRONTEND_URL":                  os.Getenv("FRONTEND_URL"),
		"GOOGLE_CLIENT_ID":              c.OAuth.Google.ClientID,
		"GOOGLE_CLIENT_SECRET":          c.OAuth.Google.ClientSecret,
		"GOOGLE_REDIRECT_URL":           c.OAuth.Google.RedirectURL,
		"KAKAO_CLIENT_ID":               c.OAuth.Kakao.ClientID,
		"KAKAO_CLIENT_SECRET":           c.OAuth.Kakao.ClientSecret,
		"KAKAO_REDIRECT_URL":            c.OAuth.Kakao.RedirectURL,
		"NAVER_CLIENT_ID":               c.OAuth.Naver.ClientID,
		"NAVER_CLIENT_SECRET":           c.OAuth.Naver.ClientSecret,
		"NAVER_REDIRECT_URL":            c.OAuth.Naver.RedirectURL,
		"SUPER_ADMIN_ID":                c.SuperAdminID,
		"SUPER_ADMIN_PW_HASH":           c.SuperAdminPWHash,
		"SUPER_ADMIN_TOTP_SECRET":       c.SuperAdminTOTPSecret,
		"OPENAI_API_KEY":                c.OpenAIAPIKey,
		"EMBEDDING_GEMMA_ENDPOINT":      c.EmbeddingGemmaEndpoint,
		"ENCRYPTION_KEY":                c.EncryptionKey,
		"NCP_OBJECT_STORAGE_BUCKET":     c.NCPObjectStorageBucket,
		"NCP_OBJECT_STORAGE_ACCESS_KEY": c.NCPObjectStorageAccessKey,
		"NCP_OBJECT_STORAGE_SECRET_KEY": c.NCPObjectStorageSecretKey,
	}
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required production env: %s", strings.Join(missing, ", "))
	}

	keyBytes, err := hex.DecodeString(strings.TrimSpace(c.EncryptionKey))
	if err != nil || len(keyBytes) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be a 64-character hex encoded 32-byte key")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}
	if c.DemoLoginEnabled && strings.TrimSpace(c.DemoLoginSecret) == "" {
		return fmt.Errorf("DEMO_LOGIN_SECRET is required when DEMO_LOGIN_ENABLED=true")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func splitAndTrim(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		switch r {
		case ',', ' ', '\t', '\n':
			return true
		default:
			return false
		}
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
