package goal

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

func TestRepositoryKeepsGoalProfilesUserScoped(t *testing.T) {
	ctx, repo, pool := setupGoalRepositoryAccessTest(t)
	ownerID := insertGoalRepositoryAccessUser(t, ctx, pool)
	otherID := insertGoalRepositoryAccessUser(t, ctx, pool)
	ownerDraftID := insertGoalRepositoryAccessDraft(t, ctx, pool, ownerID)

	ownerGoal := createGoalRepositoryAccessProfile(t, ctx, repo, ownerID, nil, "owner wants watercolor")
	otherGoal := createGoalRepositoryAccessProfile(t, ctx, repo, otherID, nil, "other wants guitar")

	active, err := repo.GetActiveGoalByUser(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetActiveGoalByUser(owner) error = %v", err)
	}
	if active.ID != ownerGoal.ID || active.UserID != ownerID {
		t.Fatalf("GetActiveGoalByUser(owner) = %+v, want owner predraft goal", active)
	}

	ownerMatches, err := repo.VerifyCourseDraftOwner(ctx, ownerDraftID, ownerID)
	if err != nil {
		t.Fatalf("VerifyCourseDraftOwner(owner) error = %v", err)
	}
	if !ownerMatches {
		t.Fatal("VerifyCourseDraftOwner(owner) = false, want true")
	}

	otherMatches, err := repo.VerifyCourseDraftOwner(ctx, ownerDraftID, otherID)
	if err != nil {
		t.Fatalf("VerifyCourseDraftOwner(other) error = %v", err)
	}
	if otherMatches {
		t.Fatal("VerifyCourseDraftOwner(other) = true, want false")
	}

	if err := repo.AttachDraftToGoal(ctx, ownerGoal.ID, otherID, ownerDraftID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("AttachDraftToGoal(cross-user) error = %v, want ErrNotFound", err)
	}

	active, err = repo.GetActiveGoalByUser(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetActiveGoalByUser(owner after cross-user attach) error = %v", err)
	}
	if active.ID != ownerGoal.ID || active.CourseDraftID != nil {
		t.Fatalf("owner goal after cross-user attach = %+v, want still predraft owner goal", active)
	}

	if err := repo.DeactivatePredraftGoalsByUser(ctx, ownerID); err != nil {
		t.Fatalf("DeactivatePredraftGoalsByUser(owner) error = %v", err)
	}
	if _, err := repo.GetActiveGoalByUser(ctx, ownerID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetActiveGoalByUser(owner after deactivate) error = %v, want ErrNotFound", err)
	}

	active, err = repo.GetActiveGoalByUser(ctx, otherID)
	if err != nil {
		t.Fatalf("GetActiveGoalByUser(other after owner deactivate) error = %v", err)
	}
	if active.ID != otherGoal.ID || active.UserID != otherID {
		t.Fatalf("other goal after owner deactivate = %+v, want untouched other goal", active)
	}

	if err := repo.AttachDraftToGoal(ctx, otherGoal.ID, ownerID, ownerDraftID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("AttachDraftToGoal(owner user with other goal) error = %v, want ErrNotFound", err)
	}

	ownerAttachGoal := createGoalRepositoryAccessProfile(t, ctx, repo, ownerID, nil, "owner ready for draft")
	if err := repo.AttachDraftToGoal(ctx, ownerAttachGoal.ID, ownerID, ownerDraftID); err != nil {
		t.Fatalf("AttachDraftToGoal(owner goal) error = %v", err)
	}
	attached, err := repo.GetActiveGoal(ctx, ownerDraftID)
	if err != nil {
		t.Fatalf("GetActiveGoal(owner draft after attach) error = %v", err)
	}
	if attached.ID != ownerAttachGoal.ID || attached.UserID != ownerID {
		t.Fatalf("attached goal = %+v, want owner attached goal", attached)
	}
}

func setupGoalRepositoryAccessTest(t *testing.T) (context.Context, *Repository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := goalRepositoryAccessTestDatabaseURL(t)
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

func goalRepositoryAccessTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}
	for _, envPath := range goalRepositoryAccessTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}
	t.Skip("DATABASE_URL or backend .env is required for goal repository access integration tests")
	return ""
}

func goalRepositoryAccessTestEnvCandidates(t *testing.T) []string {
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

func insertGoalRepositoryAccessUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "goal-access-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'Goal Access Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func insertGoalRepositoryAccessDraft(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) uuid.UUID {
	t.Helper()
	draftID := uuid.New()
	if _, err := pool.Exec(ctx, `
		INSERT INTO course_drafts (
			id, user_id, source_query, learning_goal, generation_language, title, description, status
		) VALUES ($1, $2, 'goal access', 'goal access learning goal', 'ko', 'Goal Access Draft', 'goal access draft', 'draft')
	`, draftID, userID); err != nil {
		t.Fatalf("insert goal access draft: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM course_drafts WHERE id = $1`, draftID)
	})
	return draftID
}

func createGoalRepositoryAccessProfile(t *testing.T, ctx context.Context, repo *Repository, userID uuid.UUID, draftID *uuid.UUID, intent string) *GoalProfile {
	t.Helper()
	profile := &GoalProfile{
		CourseDraftID:  draftID,
		UserID:         userID,
		UserIntent:     intent,
		Language:       "ko",
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: intent},
		},
		Version: 1,
	}
	if err := repo.CreateGoalProfile(ctx, profile); err != nil {
		t.Fatalf("CreateGoalProfile(%s) error = %v", intent, err)
	}
	t.Cleanup(func() {
		_, _ = repo.pool.Exec(context.Background(), `DELETE FROM course_goal_profiles WHERE id = $1`, profile.ID)
	})
	return profile
}
