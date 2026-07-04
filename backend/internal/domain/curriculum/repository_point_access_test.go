package curriculum

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type learningPointAccessFixture struct {
	UserID   uuid.UUID
	DraftID  uuid.UUID
	CourseID uuid.UUID
	LevelID  uuid.UUID
	LessonID uuid.UUID
	PointID  uuid.UUID
}

func TestLearningPointContextRejectsCrossUserAndMismatchedPoint(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	owner := insertLearningPointAccessFixture(t, ctx, pool)
	other := insertLearningPointAccessFixture(t, ctx, pool)

	courseID, err := repo.ResolveLearningPointCourseID(ctx, owner.UserID, owner.CourseID, owner.PointID)
	if err != nil {
		t.Fatalf("ResolveLearningPointCourseID(owner) error = %v", err)
	}
	if courseID != owner.CourseID {
		t.Fatalf("ResolveLearningPointCourseID(owner) = %s, want %s", courseID, owner.CourseID)
	}

	_, err = repo.ResolveLearningPointCourseID(ctx, other.UserID, owner.CourseID, owner.PointID)
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("ResolveLearningPointCourseID(cross-user) error = %v, want errLearningPointNotFound", err)
	}

	_, err = repo.ResolveLearningPointCourseID(ctx, owner.UserID, owner.CourseID, other.PointID)
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("ResolveLearningPointCourseID(mismatched point) error = %v, want errLearningPointNotFound", err)
	}

	_, err = repo.CountLearningPointAttachments(ctx, other.UserID, owner.CourseID, owner.PointID, "work_attachment")
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("CountLearningPointAttachments(cross-user) error = %v, want errLearningPointNotFound", err)
	}
}

func TestReadablePointContextAllowsArchivedOwnerButRejectsCrossUser(t *testing.T) {
	ctx, repo, pool := setupDraftNodeRepositoryTest(t)
	owner := insertLearningPointAccessFixture(t, ctx, pool)
	other := insertLearningPointAccessFixture(t, ctx, pool)

	if _, err := pool.Exec(ctx, `UPDATE course_drafts SET status = 'archived' WHERE id = $1`, owner.DraftID); err != nil {
		t.Fatalf("archive owner draft: %v", err)
	}

	_, _, err := repo.GetReadablePointAttachment(ctx, owner.UserID, owner.CourseID, owner.PointID, uuid.New())
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("GetReadablePointAttachment(owner archived missing attachment) error = %v, want errLearningPointNotFound", err)
	}

	_, _, err = repo.GetReadablePointAttachment(ctx, other.UserID, owner.CourseID, owner.PointID, uuid.New())
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("GetReadablePointAttachment(cross-user archived) error = %v, want errLearningPointNotFound", err)
	}

	_, err = repo.ResolveLearningPointCourseID(ctx, owner.UserID, owner.CourseID, owner.PointID)
	if !errors.Is(err, errLearningPointNotFound) {
		t.Fatalf("ResolveLearningPointCourseID(owner archived write context) error = %v, want errLearningPointNotFound", err)
	}
}

func insertLearningPointAccessFixture(t *testing.T, ctx context.Context, pool *pgxpool.Pool) learningPointAccessFixture {
	t.Helper()
	userID := insertDraftNodeTestUser(t, ctx, pool)
	draftID := uuid.New()
	courseID := uuid.New()
	levelID := uuid.New()
	lessonID := uuid.New()
	pointID := uuid.New()

	if _, err := pool.Exec(ctx, `
		INSERT INTO course_drafts (
			id, user_id, source_query, learning_goal, generation_language, title, description, status
		) VALUES ($1, $2, 'security point access', 'security point access goal', 'ko', 'Security Access Draft', 'security access draft', 'draft')
	`, draftID, userID); err != nil {
		t.Fatalf("insert access draft: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO courses (
			id, user_id, source_draft_id, source_query, learning_goal, title, description, status
		) VALUES ($1, $2, $3, 'security point access', 'security point access goal', 'Security Access Course', 'security access course', 'active')
	`, courseID, userID, draftID); err != nil {
		t.Fatalf("insert access course: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		UPDATE course_drafts
		SET status = 'learning', confirmed_course_id = $2
		WHERE id = $1
	`, draftID, courseID); err != nil {
		t.Fatalf("attach access course to draft: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO course_levels (id, course_id, title, level_code, order_index)
		VALUES ($1, $2, 'Security Access Level', 'custom', 0)
	`, levelID, courseID); err != nil {
		t.Fatalf("insert access level: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO course_lessons (
			id, course_level_id, course_id, title, lesson_role, source_type, order_index
		) VALUES ($1, $2, $3, 'Security Access Lesson', 'core', 'manual', 0)
	`, lessonID, levelID, courseID); err != nil {
		t.Fatalf("insert access lesson: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO course_points (
			id, course_id, course_lesson_id, point_type, status, title, order_index
		) VALUES ($1, $2, $3, 'research', 'learning', 'Security Access Point', 0)
	`, pointID, courseID, lessonID); err != nil {
		t.Fatalf("insert access point: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM course_points WHERE id = $1`, pointID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM course_lessons WHERE id = $1`, lessonID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM course_levels WHERE id = $1`, levelID)
		_, _ = pool.Exec(cleanupCtx, `UPDATE course_drafts SET confirmed_course_id = NULL WHERE id = $1`, draftID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM courses WHERE id = $1`, courseID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM course_drafts WHERE id = $1`, draftID)
	})

	return learningPointAccessFixture{
		UserID:   userID,
		DraftID:  draftID,
		CourseID: courseID,
		LevelID:  levelID,
		LessonID: lessonID,
		PointID:  pointID,
	}
}
