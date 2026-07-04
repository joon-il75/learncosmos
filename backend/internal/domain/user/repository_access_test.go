package user

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestPostgresRepositoryKeepsAISettingsUserScoped(t *testing.T) {
	ctx, repo, pool := setupUserRepositoryAccessTest(t)
	ownerID := insertUserRepositoryAccessUser(t, ctx, pool)
	otherID := insertUserRepositoryAccessUser(t, ctx, pool)

	insertUserRepositoryAPIKey(t, ctx, pool, ownerID, "openai", true)
	insertUserRepositoryAPIKey(t, ctx, pool, otherID, "anthropic", true)

	if err := repo.UpdateUserAISettingsEnabled(ctx, ownerID.String(), false); err != nil {
		t.Fatalf("UpdateUserAISettingsEnabled(owner) error = %v", err)
	}

	ownerSettings, err := repo.GetUserAISettings(ctx, ownerID.String())
	if err != nil {
		t.Fatalf("GetUserAISettings(owner) error = %v", err)
	}
	if ownerSettings == nil || ownerSettings.Provider != "openai" || ownerSettings.IsEnabled {
		t.Fatalf("owner settings = %+v, want disabled openai settings", ownerSettings)
	}

	otherSettings, err := repo.GetUserAISettings(ctx, otherID.String())
	if err != nil {
		t.Fatalf("GetUserAISettings(other) error = %v", err)
	}
	if otherSettings == nil || otherSettings.Provider != "anthropic" || !otherSettings.IsEnabled {
		t.Fatalf("other settings = %+v, want still-enabled anthropic settings", otherSettings)
	}

	if err := repo.DeleteUserAISettings(ctx, ownerID.String()); err != nil {
		t.Fatalf("DeleteUserAISettings(owner) error = %v", err)
	}
	ownerSettings, err = repo.GetUserAISettings(ctx, ownerID.String())
	if err != nil {
		t.Fatalf("GetUserAISettings(owner after delete) error = %v", err)
	}
	if ownerSettings != nil {
		t.Fatalf("owner settings after delete = %+v, want nil", ownerSettings)
	}
	otherSettings, err = repo.GetUserAISettings(ctx, otherID.String())
	if err != nil {
		t.Fatalf("GetUserAISettings(other after owner delete) error = %v", err)
	}
	if otherSettings == nil || otherSettings.Provider != "anthropic" {
		t.Fatalf("other settings after owner delete = %+v, want untouched", otherSettings)
	}
}

func TestPostgresRepositoryKeepsUsageAndPointHistoryUserScoped(t *testing.T) {
	ctx, repo, pool := setupUserRepositoryAccessTest(t)
	ownerID := insertUserRepositoryAccessUser(t, ctx, pool)
	otherID := insertUserRepositoryAccessUser(t, ctx, pool)
	now := time.Now().UTC().Truncate(time.Second)
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)

	insertUserRepositoryAIUsage(t, ctx, pool, ownerID, "byok", 10, 5, now)
	insertUserRepositoryAIUsage(t, ctx, pool, ownerID, "system", 100, 100, now)
	insertUserRepositoryAIUsage(t, ctx, pool, otherID, "byok", 30, 7, now)

	events, total, summary, err := repo.ListUserAIUsageEvents(ctx, ownerID.String(), start, end, 10, 0)
	if err != nil {
		t.Fatalf("ListUserAIUsageEvents(owner) error = %v", err)
	}
	if total != 1 || len(events) != 1 {
		t.Fatalf("owner ai usage total=%d len=%d, want 1", total, len(events))
	}
	if summary.TotalTokens != 15 {
		t.Fatalf("owner ai usage total tokens = %d, want 15", summary.TotalTokens)
	}

	insertUserRepositoryWallet(t, ctx, pool, ownerID, 3, 2)
	insertUserRepositoryWallet(t, ctx, pool, otherID, 70, 20)
	free, paid, err := repo.GetPointBalances(ctx, ownerID.String())
	if err != nil {
		t.Fatalf("GetPointBalances(owner) error = %v", err)
	}
	if free != 3 || paid != 2 {
		t.Fatalf("owner balances free=%d paid=%d, want 3/2", free, paid)
	}

	insertUserRepositoryPointTransaction(t, ctx, pool, ownerID, "grant_free", 3, now)
	insertUserRepositoryPointTransaction(t, ctx, pool, otherID, "purchase", 70, now)
	transactions, total, pointSummary, err := repo.ListUserPointTransactions(ctx, ownerID.String(), start, end, 10, 0)
	if err != nil {
		t.Fatalf("ListUserPointTransactions(owner) error = %v", err)
	}
	if total != 1 || len(transactions) != 1 {
		t.Fatalf("owner point transactions total=%d len=%d, want 1", total, len(transactions))
	}
	if pointSummary.GrantedPoints != 3 || pointSummary.NetChange != 3 {
		t.Fatalf("owner point summary = %+v, want grant/net 3", pointSummary)
	}
}

func setupUserRepositoryAccessTest(t *testing.T) (context.Context, *PostgresRepository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := userRepositoryAccessTestDatabaseURL(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping test database: %v", err)
	}
	t.Cleanup(pool.Close)
	return ctx, NewPostgresRepository(pool), pool
}

func userRepositoryAccessTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}
	for _, envPath := range userRepositoryAccessTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}
	t.Skip("DATABASE_URL or backend .env is required for user repository access integration tests")
	return ""
}

func userRepositoryAccessTestEnvCandidates(t *testing.T) []string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	var candidates []string
	for dir := cwd; ; dir = filepath.Dir(dir) {
		candidates = append(candidates, filepath.Join(dir, ".env"))
		candidates = append(candidates, filepath.Join(dir, "backend", ".env"))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return candidates
}

func insertUserRepositoryAccessUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "user-access-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'User Access Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func insertUserRepositoryAPIKey(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, provider string, enabled bool) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO user_api_keys (id, user_id, provider, api_key_encrypted, is_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, 'encrypted-test-key', $4, NOW(), NOW())
	`, uuid.New(), userID, provider, enabled); err != nil {
		t.Fatalf("insert user api key: %v", err)
	}
}

func insertUserRepositoryAIUsage(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, source string, inputTokens, outputTokens int, createdAt time.Time) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO ai_usage_events (
			id, user_id, source, provider, model, feature, billing_status, input_tokens, output_tokens, estimated_cost_usd, success, created_at
		) VALUES ($1, $2, $3, 'openai', 'gpt-test', 'security_access_test', 'byok', $4, $5, 0.010000, true, $6)
	`, uuid.New(), userID, source, inputTokens, outputTokens, createdAt); err != nil {
		t.Fatalf("insert ai usage event: %v", err)
	}
}

func insertUserRepositoryWallet(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, freeBalance, paidBalance int) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO ai_point_wallets (id, user_id, free_balance, paid_balance, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, uuid.New(), userID, freeBalance, paidBalance); err != nil {
		t.Fatalf("insert point wallet: %v", err)
	}
}

func insertUserRepositoryPointTransaction(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, transactionType string, amount int, createdAt time.Time) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO ai_point_transactions (id, user_id, type, amount, feature, description, metadata, created_at)
		VALUES ($1, $2, $3, $4, 'security_access_test', 'security access test', '{}'::jsonb, $5)
	`, uuid.New(), userID, transactionType, amount, createdAt); err != nil {
		t.Fatalf("insert point transaction: %v", err)
	}
}

func TestRecordUserAIValidationStoresProviderErrorCode(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatalf("read repository.go: %v", err)
	}
	text := string(source)
	if !strings.Contains(text, "trimmed := logsafe.ProviderErrorCode(message)") {
		t.Fatal("RecordUserAIValidation must store a provider error code instead of raw validation message")
	}
	if strings.Contains(text, "trimmed := strings.TrimSpace(message)") || strings.Contains(text, "trimmed := message") {
		t.Fatal("RecordUserAIValidation must not persist raw validation messages")
	}
}
