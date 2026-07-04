package content

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestRepositoryKeepsCRUDUserScoped(t *testing.T) {
	ctx, repo, pool := setupContentRepositoryAccessTest(t)
	ownerID := insertContentRepositoryAccessUser(t, ctx, pool)
	otherID := insertContentRepositoryAccessUser(t, ctx, pool)

	ownerContent := createContentRepositoryAccessContent(t, ctx, repo, ownerID, "Owner Content")
	otherContent := createContentRepositoryAccessContent(t, ctx, repo, otherID, "Other Content")

	got, err := repo.GetByID(ctx, ownerContent.ID, ownerID)
	if err != nil {
		t.Fatalf("GetByID(owner) error = %v", err)
	}
	if got.ID != ownerContent.ID || got.UserID != ownerID {
		t.Fatalf("GetByID(owner) = %+v, want owner content", got)
	}

	if _, err := repo.GetByID(ctx, ownerContent.ID, otherID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("GetByID(cross-user) error = %v, want pgx.ErrNoRows", err)
	}

	list, err := repo.List(ctx, ownerID, ListContentsQuery{Limit: 10})
	if err != nil {
		t.Fatalf("List(owner) error = %v", err)
	}
	if list.Total != 1 || len(list.Contents) != 1 || list.Contents[0].ID != ownerContent.ID {
		t.Fatalf("List(owner) = total %d contents %+v, want only owner content", list.Total, list.Contents)
	}

	otherTitle := "Cross User Update"
	if _, err := repo.Update(ctx, ownerContent.ID, otherID, UpdateContentRequest{Title: &otherTitle}); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("Update(cross-user) error = %v, want pgx.ErrNoRows", err)
	}

	got, err = repo.GetByID(ctx, ownerContent.ID, ownerID)
	if err != nil {
		t.Fatalf("GetByID(owner after cross-user update) error = %v", err)
	}
	if got.Title != "Owner Content" {
		t.Fatalf("owner content title after cross-user update = %q, want unchanged", got.Title)
	}

	ownerTitle := "Owner Updated Content"
	updated, err := repo.Update(ctx, ownerContent.ID, ownerID, UpdateContentRequest{Title: &ownerTitle})
	if err != nil {
		t.Fatalf("Update(owner) error = %v", err)
	}
	if updated.Title != ownerTitle {
		t.Fatalf("Update(owner) title = %q, want %q", updated.Title, ownerTitle)
	}

	if err := repo.Delete(ctx, ownerContent.ID, otherID); err == nil {
		t.Fatal("Delete(cross-user) error = nil, want error")
	}

	got, err = repo.GetByID(ctx, ownerContent.ID, ownerID)
	if err != nil {
		t.Fatalf("GetByID(owner after cross-user delete) error = %v", err)
	}
	if got.ID != ownerContent.ID {
		t.Fatalf("owner content after cross-user delete = %s, want %s", got.ID, ownerContent.ID)
	}

	if err := repo.Delete(ctx, ownerContent.ID, ownerID); err != nil {
		t.Fatalf("Delete(owner) error = %v", err)
	}
	if _, err := repo.GetByID(ctx, ownerContent.ID, ownerID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("GetByID(owner after owner delete) error = %v, want pgx.ErrNoRows", err)
	}

	got, err = repo.GetByID(ctx, otherContent.ID, otherID)
	if err != nil {
		t.Fatalf("GetByID(other after owner delete) error = %v", err)
	}
	if got.ID != otherContent.ID || got.UserID != otherID {
		t.Fatalf("other content after owner delete = %+v, want untouched", got)
	}
}

func setupContentRepositoryAccessTest(t *testing.T) (context.Context, *Repository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := contentRepositoryAccessTestDatabaseURL(t)
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

func contentRepositoryAccessTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}
	for _, envPath := range contentRepositoryAccessTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}
	t.Skip("DATABASE_URL or backend .env is required for content repository access integration tests")
	return ""
}

func contentRepositoryAccessTestEnvCandidates(t *testing.T) []string {
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

func insertContentRepositoryAccessUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "content-access-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'Content Access Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func createContentRepositoryAccessContent(t *testing.T, ctx context.Context, repo *Repository, userID uuid.UUID, title string) *Content {
	t.Helper()
	rawURL := "https://example.com/" + strings.ReplaceAll(strings.ToLower(title), " ", "-")
	content, err := repo.Create(ctx, userID, CreateContentRequest{
		ContentType: ContentTypeArticle,
		URL:         &rawURL,
		Title:       title,
		Language:    "ko",
	})
	if err != nil {
		t.Fatalf("Create(%s) error = %v", title, err)
	}
	return content
}
