package curriculum

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

func TestCreateMainLessonStoresManualCoreFallbackSpec(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	userID, draftID := insertDraftNodeTestDraft(t, ctx, pool, DraftStatusDraft)

	objective := "  공통 main helper smoke  "
	lesson, err := repo.CreateMainLesson(ctx, userID, draftID, CreateMainLessonRequest{
		Title:      "  Smoke Main Helper Region  ",
		Objective:  &objective,
		OrderIndex: 7,
	})
	if err != nil {
		t.Fatalf("CreateMainLesson() error = %v", err)
	}

	if lesson.CourseDraftID != draftID {
		t.Fatalf("CourseDraftID = %v, want %v", lesson.CourseDraftID, draftID)
	}
	if lesson.ParentLessonID != nil {
		t.Fatalf("ParentLessonID = %v, want nil", *lesson.ParentLessonID)
	}
	if lesson.Title != "Smoke Main Helper Region" {
		t.Fatalf("Title = %q", lesson.Title)
	}
	if lesson.Objective == nil || *lesson.Objective != "공통 main helper smoke" {
		t.Fatalf("Objective = %v", lesson.Objective)
	}
	if lesson.Summary == nil || *lesson.Summary != "공통 main helper smoke" {
		t.Fatalf("Summary = %v", lesson.Summary)
	}
	if lesson.LessonRole != LessonRoleCore {
		t.Fatalf("LessonRole = %q, want %q", lesson.LessonRole, LessonRoleCore)
	}
	if lesson.SourceType != LessonSourceManual {
		t.Fatalf("SourceType = %q, want %q", lesson.SourceType, LessonSourceManual)
	}
	if lesson.OrderIndex != 7 {
		t.Fatalf("OrderIndex = %d, want 7", lesson.OrderIndex)
	}
	assertFallbackSearchSpec(t, lesson.RecommendationSearchSpec)

	var dbRole LessonRole
	var dbSource LessonSourceType
	var dbSpecSource string
	if err := pool.QueryRow(ctx, `
		SELECT lesson_role, source_type, recommendation_search_spec->>'source'
		FROM course_draft_lessons
		WHERE id = $1
	`, lesson.ID).Scan(&dbRole, &dbSource, &dbSpecSource); err != nil {
		t.Fatalf("load created main lesson: %v", err)
	}
	if dbRole != LessonRoleCore || dbSource != LessonSourceManual || dbSpecSource != "fallback" {
		t.Fatalf("db contract = role:%q source:%q spec:%q", dbRole, dbSource, dbSpecSource)
	}
}

func TestCreateSubLessonStoresManualSupportFallbackSpecAndAutoOrder(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	userID, draftID := insertDraftNodeTestDraft(t, ctx, pool, DraftStatusDraft)

	mainLesson, err := repo.CreateMainLesson(ctx, userID, draftID, CreateMainLessonRequest{
		Title:      "Main Lesson",
		OrderIndex: 0,
	})
	if err != nil {
		t.Fatalf("CreateMainLesson() error = %v", err)
	}

	firstObjective := "첫 번째 보조 리슨"
	first, err := repo.CreateSubLesson(ctx, userID, draftID, mainLesson.ID, CreateSubLessonRequest{
		Title:      "First Support Lesson",
		Objective:  &firstObjective,
		OrderIndex: 0,
	})
	if err != nil {
		t.Fatalf("CreateSubLesson(first) error = %v", err)
	}
	secondObjective := "두 번째 보조 리슨"
	second, err := repo.CreateSubLesson(ctx, userID, draftID, mainLesson.ID, CreateSubLessonRequest{
		Title:      "  Second Support Lesson  ",
		Objective:  &secondObjective,
		OrderIndex: 0,
	})
	if err != nil {
		t.Fatalf("CreateSubLesson(second) error = %v", err)
	}

	if first.ParentLessonID == nil || *first.ParentLessonID != mainLesson.ID {
		t.Fatalf("first parent = %v, want %v", first.ParentLessonID, mainLesson.ID)
	}
	if first.LessonRole != LessonRoleSupport {
		t.Fatalf("first LessonRole = %q, want %q", first.LessonRole, LessonRoleSupport)
	}
	if first.SourceType != LessonSourceManual {
		t.Fatalf("first SourceType = %q, want %q", first.SourceType, LessonSourceManual)
	}
	if first.OrderIndex != 0 {
		t.Fatalf("first OrderIndex = %d, want 0", first.OrderIndex)
	}
	assertFallbackSearchSpec(t, first.RecommendationSearchSpec)

	if second.Title != "Second Support Lesson" {
		t.Fatalf("second Title = %q", second.Title)
	}
	if second.OrderIndex != 1 {
		t.Fatalf("second OrderIndex = %d, want 1", second.OrderIndex)
	}
	if second.Objective == nil || *second.Objective != secondObjective {
		t.Fatalf("second Objective = %v", second.Objective)
	}

	var supportCount int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM course_draft_lessons
		WHERE course_draft_id = $1
		  AND parent_lesson_id = $2
		  AND lesson_role = 'support'
		  AND source_type = 'manual'
		  AND recommendation_search_spec->>'source' = 'fallback'
	`, draftID, mainLesson.ID).Scan(&supportCount); err != nil {
		t.Fatalf("count support lessons: %v", err)
	}
	if supportCount != 2 {
		t.Fatalf("supportCount = %d, want 2", supportCount)
	}
}

func TestCreateLessonHelpersRejectWrongOwner(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	userID, draftID := insertDraftNodeTestDraft(t, ctx, pool, DraftStatusDraft)
	wrongUserID := insertDraftNodeTestUser(t, ctx, pool)

	_, err := repo.CreateMainLesson(ctx, wrongUserID, draftID, CreateMainLessonRequest{Title: "Wrong Owner"})
	if !errors.Is(err, errDraftResourceNotFound) {
		t.Fatalf("CreateMainLesson wrong owner error = %v, want errDraftResourceNotFound", err)
	}

	mainLesson, err := repo.CreateMainLesson(ctx, userID, draftID, CreateMainLessonRequest{Title: "Main Lesson"})
	if err != nil {
		t.Fatalf("CreateMainLesson() error = %v", err)
	}
	_, err = repo.CreateSubLesson(ctx, wrongUserID, draftID, mainLesson.ID, CreateSubLessonRequest{Title: "Wrong Owner Sub"})
	if !errors.Is(err, errDraftResourceNotFound) {
		t.Fatalf("CreateSubLesson wrong owner error = %v, want errDraftResourceNotFound", err)
	}
}

func TestCreateMainLessonRejectsArchivedDraft(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	userID, draftID := insertDraftNodeTestDraft(t, ctx, pool, DraftStatusArchived)

	_, err := repo.CreateMainLesson(ctx, userID, draftID, CreateMainLessonRequest{Title: "Archived Draft Lesson"})
	if !errors.Is(err, errDraftResourceNotFound) {
		t.Fatalf("CreateMainLesson archived draft error = %v, want errDraftResourceNotFound", err)
	}
}

func setupDraftNodeRepositoryTest(t *testing.T) (context.Context, *Repository, *pgxpool.Pool) {
	t.Helper()
	databaseURL := draftNodeTestDatabaseURL(t)

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

func draftNodeTestDatabaseURL(t *testing.T) string {
	t.Helper()
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}

	for _, envPath := range draftNodeTestEnvCandidates(t) {
		_ = godotenv.Load(envPath)
		if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
			return databaseURL
		}
	}

	t.Skip("DATABASE_URL or LEARNWEAVER_TEST_DATABASE_URL is required for draft node repository integration tests")
	return ""
}

func draftNodeTestEnvCandidates(t *testing.T) []string {
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

func insertDraftNodeTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	email := "draft-node-helper-" + userID.String() + "@example.test"
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, nickname, role)
		VALUES ($1, $2, 'Draft Node Helper Test', 'learner')
	`, userID, email); err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func insertDraftNodeTestDraft(t *testing.T, ctx context.Context, pool *pgxpool.Pool, status DraftStatus) (uuid.UUID, uuid.UUID) {
	t.Helper()
	userID := insertDraftNodeTestUser(t, ctx, pool)
	draftID := uuid.New()
	learningGoal := "가죽공예를 기초부터 배워서 나만의 손 지갑을 만들어 보자"
	if _, err := pool.Exec(ctx, `
		INSERT INTO course_drafts (
			id, user_id, source_query, learning_goal, generation_language, title, description, status
		) VALUES ($1, $2, $3, $4, 'ko', $5, $6, $7)
	`, draftID, userID, "가죽공예 손지갑 만들기", learningGoal, "Helper Test Draft", "helper test draft", status); err != nil {
		t.Fatalf("insert test draft: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM course_draft_lessons WHERE course_draft_id = $1`, draftID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM course_drafts WHERE id = $1`, draftID)
	})
	return userID, draftID
}

func assertFallbackSearchSpec(t *testing.T, spec LessonRecommendationSearchSpec) {
	t.Helper()
	if spec.Source != "fallback" {
		t.Fatalf("spec.Source = %q, want fallback", spec.Source)
	}
	if strings.TrimSpace(spec.PrimaryQuery) == "" {
		t.Fatalf("spec.PrimaryQuery is empty")
	}
	if spec.Language != "ko" {
		t.Fatalf("spec.Language = %q, want ko", spec.Language)
	}
	if len(spec.MustInclude) == 0 {
		t.Fatalf("spec.MustInclude is empty")
	}
}
