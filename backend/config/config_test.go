package config

import (
	"strings"
	"testing"
)

func TestValidateSkipsNonProduction(t *testing.T) {
	cfg := &Config{Env: "development"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected non-production config to skip validation, got %v", err)
	}
}

func TestValidateProductionRequiresEnv(t *testing.T) {
	cfg := &Config{Env: "production"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected missing production env error")
	}
	if !strings.Contains(err.Error(), "missing required production env") {
		t.Fatalf("expected missing env error, got %v", err)
	}
}

func TestValidateProductionAcceptsRequiredEnv(t *testing.T) {
	setRequiredProductionEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateProductionRejectsInvalidEncryptionKey(t *testing.T) {
	setRequiredProductionEnv(t)
	t.Setenv("ENCRYPTION_KEY", "not-hex")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	err = cfg.Validate()
	if err == nil {
		t.Fatal("expected invalid ENCRYPTION_KEY error")
	}
	if !strings.Contains(err.Error(), "ENCRYPTION_KEY") {
		t.Fatalf("expected ENCRYPTION_KEY error, got %v", err)
	}
}

func TestLoadGoalChatSinglePromptFallbackDefaultsFalse(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GoalChatSinglePromptFallbackEnabled {
		t.Fatal("expected goal chat single prompt fallback to default false")
	}
}

func TestLoadGoalChatSinglePromptFallbackEnabled(t *testing.T) {
	t.Setenv("GOAL_CHAT_SINGLE_PROMPT_FALLBACK_ENABLED", "true")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.GoalChatSinglePromptFallbackEnabled {
		t.Fatal("expected goal chat single prompt fallback to be enabled")
	}
}

func TestLoadGoalChatModeDefaultsTwoStep(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GoalChatMode != "two_step" {
		t.Fatalf("goal chat mode = %q, want two_step", cfg.GoalChatMode)
	}
}

func TestLoadGoalChatModeFromEnv(t *testing.T) {
	t.Setenv("GOAL_CHAT_MODE", "analysis_template")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GoalChatMode != "analysis_template" {
		t.Fatalf("goal chat mode = %q, want analysis_template", cfg.GoalChatMode)
	}
}

func TestLoadGoalChatFastPathDefaultsFalse(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.GoalChatFastPathEnabled {
		t.Fatal("expected goal chat fast path to default false")
	}
}

func TestLoadGoalChatFastPathEnabled(t *testing.T) {
	t.Setenv("GOAL_CHAT_FAST_PATH_ENABLED", "true")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.GoalChatFastPathEnabled {
		t.Fatal("expected goal chat fast path to be enabled")
	}
}

func setRequiredProductionEnv(t *testing.T) {
	t.Helper()

	values := map[string]string{
		"APP_ENV":                       "production",
		"DATABASE_URL":                  "postgres://localhost:5432/learnweaver",
		"REDIS_URL":                     "redis://localhost:6379/0",
		"JWT_SECRET":                    "12345678901234567890123456789012",
		"FRONTEND_URL":                  "https://www.learnweavr.com",
		"GOOGLE_CLIENT_ID":              "google-client-id",
		"GOOGLE_CLIENT_SECRET":          "google-client-secret",
		"GOOGLE_REDIRECT_URL":           "https://www.learnweavr.com/api/v1/auth/google/callback",
		"KAKAO_CLIENT_ID":               "kakao-client-id",
		"KAKAO_CLIENT_SECRET":           "kakao-client-secret",
		"KAKAO_REDIRECT_URL":            "https://www.learnweavr.com/api/v1/auth/kakao/callback",
		"NAVER_CLIENT_ID":               "naver-client-id",
		"NAVER_CLIENT_SECRET":           "naver-client-secret",
		"NAVER_REDIRECT_URL":            "https://www.learnweavr.com/api/v1/auth/naver/callback",
		"SUPER_ADMIN_ID":                "super-admin",
		"SUPER_ADMIN_PW_HASH":           "hash",
		"SUPER_ADMIN_TOTP_SECRET":       "totp-secret",
		"OPENAI_API_KEY":                "openai-key",
		"EMBEDDING_GEMMA_ENDPOINT":      "http://127.0.0.1:9000/v1/embeddings",
		"ENCRYPTION_KEY":                "0000000000000000000000000000000000000000000000000000000000000000",
		"NCP_OBJECT_STORAGE_BUCKET":     "learnweaver-prod",
		"NCP_OBJECT_STORAGE_ACCESS_KEY": "access-key",
		"NCP_OBJECT_STORAGE_SECRET_KEY": "secret-key",
	}
	for key, value := range values {
		t.Setenv(key, value)
	}
}

func TestLoadLessonSearchPrerunSchedulerDefaultsDisabled(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LessonSearchPrerunSchedulerEnabled {
		t.Fatal("expected lesson search prerun scheduler to default disabled")
	}
	if cfg.LessonSearchPrerunIntervalHours != 24 {
		t.Fatalf("interval hours = %d, want 24", cfg.LessonSearchPrerunIntervalHours)
	}
	if cfg.LessonSearchPrerunStartDelayMinutes != 10 {
		t.Fatalf("start delay minutes = %d, want 10", cfg.LessonSearchPrerunStartDelayMinutes)
	}
	if cfg.LessonSearchPrerunTopN != 5 || cfg.LessonSearchPrerunRiskTopN != 10 {
		t.Fatalf("top_n/risk_top_n = %d/%d, want 5/10", cfg.LessonSearchPrerunTopN, cfg.LessonSearchPrerunRiskTopN)
	}
	if cfg.LessonSearchPrerunTimeoutSeconds != 60 {
		t.Fatalf("timeout seconds = %d, want 60", cfg.LessonSearchPrerunTimeoutSeconds)
	}
	if cfg.LessonSearchPrerunProvider != "all" {
		t.Fatalf("provider = %q, want all", cfg.LessonSearchPrerunProvider)
	}
	if cfg.LessonSearchPrerunReportRetentionDays != 30 {
		t.Fatalf("retention days = %d, want 30", cfg.LessonSearchPrerunReportRetentionDays)
	}
	if len(cfg.LessonSearchPrerunRiskFixtureIDs) == 0 {
		t.Fatal("expected default risk fixture ids")
	}
}

func TestLoadLessonSearchPrerunSchedulerFromEnv(t *testing.T) {
	t.Setenv("LESSON_SEARCH_PRERUN_SCHEDULER_ENABLED", "true")
	t.Setenv("LESSON_SEARCH_PRERUN_INTERVAL_HOURS", "12")
	t.Setenv("LESSON_SEARCH_PRERUN_START_DELAY_MINUTES", "1")
	t.Setenv("LESSON_SEARCH_PRERUN_TOP_N", "3")
	t.Setenv("LESSON_SEARCH_PRERUN_RISK_TOP_N", "7")
	t.Setenv("LESSON_SEARCH_PRERUN_TIMEOUT_SECONDS", "30")
	t.Setenv("LESSON_SEARCH_PRERUN_PROVIDER", "youtube")
	t.Setenv("LESSON_SEARCH_PRERUN_REPORT_RETENTION_DAYS", "14")
	t.Setenv("LESSON_SEARCH_PRERUN_RISK_FIXTURE_IDS", "fixture-a, fixture-b fixture-c")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.LessonSearchPrerunSchedulerEnabled {
		t.Fatal("expected lesson search prerun scheduler enabled")
	}
	if cfg.LessonSearchPrerunIntervalHours != 12 || cfg.LessonSearchPrerunStartDelayMinutes != 1 {
		t.Fatalf("interval/start delay = %d/%d, want 12/1", cfg.LessonSearchPrerunIntervalHours, cfg.LessonSearchPrerunStartDelayMinutes)
	}
	if cfg.LessonSearchPrerunTopN != 3 || cfg.LessonSearchPrerunRiskTopN != 7 {
		t.Fatalf("top_n/risk_top_n = %d/%d, want 3/7", cfg.LessonSearchPrerunTopN, cfg.LessonSearchPrerunRiskTopN)
	}
	if cfg.LessonSearchPrerunTimeoutSeconds != 30 || cfg.LessonSearchPrerunReportRetentionDays != 14 {
		t.Fatalf("timeout/retention = %d/%d, want 30/14", cfg.LessonSearchPrerunTimeoutSeconds, cfg.LessonSearchPrerunReportRetentionDays)
	}
	if cfg.LessonSearchPrerunProvider != "youtube" {
		t.Fatalf("provider = %q, want youtube", cfg.LessonSearchPrerunProvider)
	}
	if got := strings.Join(cfg.LessonSearchPrerunRiskFixtureIDs, ","); got != "fixture-a,fixture-b,fixture-c" {
		t.Fatalf("risk fixture ids = %q, want fixture-a,fixture-b,fixture-c", got)
	}
}
