package explorer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestHandlerRejectsCrossUserRegionAccess(t *testing.T) {
	ctx, handler, repo, pool := setupExplorerHandlerAccessTest(t)
	ownerID, draftID, regionID := insertExplorerRegionAccessFixture(t, ctx, pool)
	otherID := insertExplorerAccessUser(t, ctx, pool)

	tests := []struct {
		name   string
		method string
		body   string
		params gin.Params
		invoke func(*gin.Context)
	}{
		{
			name:   "get course",
			method: http.MethodGet,
			params: gin.Params{{Key: "courseId", Value: draftID.String()}},
			invoke: handler.GetCourse,
		},
		{
			name:   "update region",
			method: http.MethodPatch,
			body:   `{"name":"Cross User Region"}`,
			params: gin.Params{{Key: "regionId", Value: regionID.String()}},
			invoke: handler.UpdateRegion,
		},
		{
			name:   "delete region",
			method: http.MethodDelete,
			params: gin.Params{{Key: "regionId", Value: regionID.String()}},
			invoke: handler.DeleteRegion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := newExplorerAccessContext(tt.method, tt.body, otherID, tt.params)
			tt.invoke(c)
			if w.Code != http.StatusForbidden {
				t.Fatalf("%s status = %d body = %s, want 403", tt.name, w.Code, w.Body.String())
			}

			region, err := repo.GetRegionByID(ctx, regionID)
			if err != nil {
				t.Fatalf("GetRegionByID(after %s) error = %v", tt.name, err)
			}
			if region.CourseDraftID != draftID || region.Name != "Owner Region" {
				t.Fatalf("region after %s = %+v, want unchanged owner region", tt.name, region)
			}
		})
	}

	if gotOwner, err := repo.GetCourseOwner(ctx, draftID); err != nil || gotOwner != ownerID {
		t.Fatalf("GetCourseOwner() = %s, %v; want %s, nil", gotOwner, err, ownerID)
	}
}

func setupExplorerHandlerAccessTest(t *testing.T) (context.Context, *Handler, *Repository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := explorerHandlerAccessTestDatabaseURL(t)
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

	repo := NewRepository(pool)
	return ctx, NewHandler(repo, NewService(), nil, nil, nil, "", ""), repo, pool
}

func explorerHandlerAccessTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}
	for _, envPath := range explorerHandlerAccessTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}
	t.Skip("DATABASE_URL or backend .env is required for explorer handler access integration tests")
	return ""
}

func explorerHandlerAccessTestEnvCandidates(t *testing.T) []string {
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

func insertExplorerAccessUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "explorer-access-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'Explorer Access Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func insertExplorerRegionAccessFixture(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	userID := insertExplorerAccessUser(t, ctx, pool)
	draftID := uuid.New()
	regionID := uuid.New()

	if _, err := pool.Exec(ctx, `
		INSERT INTO course_drafts (
			id, user_id, source_query, learning_goal, generation_language, title, description, status
		) VALUES ($1, $2, 'explorer access', 'explorer access goal', 'ko', 'Explorer Access Draft', 'explorer access draft', 'draft')
	`, draftID, userID); err != nil {
		t.Fatalf("insert explorer access draft: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO explorer_regions (id, course_draft_id, name, description, order_index)
		VALUES ($1, $2, 'Owner Region', 'owner region', 0)
	`, regionID, draftID); err != nil {
		t.Fatalf("insert explorer access region: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM explorer_regions WHERE id = $1`, regionID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM course_drafts WHERE id = $1`, draftID)
	})
	return userID, draftID, regionID
}

func newExplorerAccessContext(method, body string, userID uuid.UUID, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/", strings.NewReader(body))
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/json")
	}
	c.Set("user_id", userID.String())
	c.Params = params
	return c, w
}
