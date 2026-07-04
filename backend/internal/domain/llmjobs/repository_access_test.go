package llmjobs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestRepositoryKeepsUserJobPollingScoped(t *testing.T) {
	ctx, repo, pool := setupLLMJobRepositoryAccessTest(t)
	ownerID := insertLLMJobRepositoryAccessUser(t, ctx, pool)
	otherID := insertLLMJobRepositoryAccessUser(t, ctx, pool)

	ownerJob := createLLMJobRepositoryAccessJob(t, ctx, repo, &ownerID, "owner")
	otherJob := createLLMJobRepositoryAccessJob(t, ctx, repo, &otherID, "other")
	systemJob := createLLMJobRepositoryAccessJob(t, ctx, repo, nil, "system")

	got, err := repo.GetJobForUser(ctx, ownerJob.ID, ownerID)
	if err != nil {
		t.Fatalf("GetJobForUser(owner job, owner user) error = %v", err)
	}
	if got.ID != ownerJob.ID || got.UserID == nil || *got.UserID != ownerID {
		t.Fatalf("GetJobForUser(owner job, owner user) = %+v, want owner job", got)
	}

	if _, err := repo.GetJobForUser(ctx, ownerJob.ID, otherID); !errors.Is(err, ErrJobAccessDenied) {
		t.Fatalf("GetJobForUser(owner job, other user) error = %v, want ErrJobAccessDenied", err)
	}

	if _, err := repo.GetJobForUser(ctx, otherJob.ID, ownerID); !errors.Is(err, ErrJobAccessDenied) {
		t.Fatalf("GetJobForUser(other job, owner user) error = %v, want ErrJobAccessDenied", err)
	}

	if _, err := repo.GetJobForUser(ctx, systemJob.ID, ownerID); !errors.Is(err, ErrJobAccessDenied) {
		t.Fatalf("GetJobForUser(system job, owner user) error = %v, want ErrJobAccessDenied", err)
	}

	got, err = repo.GetJob(ctx, otherJob.ID)
	if err != nil {
		t.Fatalf("GetJob(other job after denied owner poll) error = %v", err)
	}
	if got.ID != otherJob.ID || got.UserID == nil || *got.UserID != otherID {
		t.Fatalf("other job after denied owner poll = %+v, want untouched other job", got)
	}
}

func setupLLMJobRepositoryAccessTest(t *testing.T) (context.Context, *Repository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := llmJobRepositoryAccessTestDatabaseURL(t)
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
	return ctx, NewRepository(pool), pool
}

func llmJobRepositoryAccessTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}
	for _, envPath := range llmJobRepositoryAccessTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}
	t.Skip("DATABASE_URL or backend .env is required for llm job repository access integration tests")
	return ""
}

func llmJobRepositoryAccessTestEnvCandidates(t *testing.T) []string {
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

func insertLLMJobRepositoryAccessUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "llm-job-access-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'LLM Job Access Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func createLLMJobRepositoryAccessJob(t *testing.T, ctx context.Context, repo *Repository, userID *uuid.UUID, label string) *Job {
	t.Helper()
	idempotencyKey := "llm-job-access-" + label + "-" + uuid.NewString()
	job, err := repo.CreateJob(ctx, CreateJobInput{
		UserID:         userID,
		Feature:        FeatureDummy,
		IdempotencyKey: &idempotencyKey,
		RequestRef: map[string]any{
			"test":  label,
			"scope": "repository_access",
		},
		PromptInputRef: map[string]any{
			"label": label,
		},
	})
	if err != nil {
		t.Fatalf("CreateJob(%s) error = %v", label, err)
	}
	t.Cleanup(func() {
		_, _ = repo.pool.Exec(context.Background(), `DELETE FROM llm_jobs WHERE id = $1`, job.ID)
	})
	return job
}
