package curriculum

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) getLearningPointContextTx(ctx context.Context, tx pgx.Tx, userID, planetID, pointID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	var draftID uuid.UUID
	var courseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id, confirmed_course_id
		FROM course_drafts
		WHERE user_id = $1
		  AND status = 'learning'
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, planetID).Scan(&draftID, &courseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("get learning draft for point action: %w", err)
	}
	if courseID == nil {
		return uuid.Nil, uuid.Nil, errLearningPointNotFound
	}

	var verifiedPointID uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM course_points
		WHERE id = $1
		  AND course_id = $2
		LIMIT 1
	`, pointID, *courseID).Scan(&verifiedPointID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("verify learning point: %w", err)
	}

	return draftID, *courseID, nil
}

func (r *Repository) getReadablePointContextTx(ctx context.Context, tx pgx.Tx, userID, planetID, pointID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	var draftID uuid.UUID
	var courseID *uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id, confirmed_course_id
		FROM course_drafts
		WHERE user_id = $1
		  AND status IN ('learning', 'archived')
		  AND (id = $2 OR confirmed_course_id = $2)
		ORDER BY updated_at DESC
		LIMIT 1
	`, userID, planetID).Scan(&draftID, &courseID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("get readable draft for point action: %w", err)
	}
	if courseID == nil {
		return uuid.Nil, uuid.Nil, errLearningPointNotFound
	}

	var verifiedPointID uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM course_points
		WHERE id = $1
		  AND course_id = $2
		LIMIT 1
	`, pointID, *courseID).Scan(&verifiedPointID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("verify readable point: %w", err)
	}
	return draftID, *courseID, nil
}

func (r *Repository) getDraftGoalProfileVersionTx(ctx context.Context, tx pgx.Tx, draftID uuid.UUID) (*int, error) {
	var version *int
	if err := tx.QueryRow(ctx, `
		SELECT goal_profile_version
		FROM course_drafts
		WHERE id = $1
		LIMIT 1
	`, draftID).Scan(&version); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get draft goal profile version: %w", err)
	}
	return version, nil
}

func (r *Repository) getLearningResearchPointContextTx(ctx context.Context, tx pgx.Tx, userID, planetID, pointID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	var verifiedPointID uuid.UUID
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM course_points
		WHERE id = $1
		  AND course_id = $2
		  AND point_type = 'research'
		LIMIT 1
	`, pointID, courseID).Scan(&verifiedPointID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("verify learning research point: %w", err)
	}

	return draftID, courseID, nil
}

func (r *Repository) getLearningExplorationPointSourceTx(ctx context.Context, tx pgx.Tx, userID, planetID, pointID uuid.UUID) (uuid.UUID, uuid.UUID, CoursePoint, error) {
	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, uuid.Nil, CoursePoint{}, err
	}

	var point CoursePoint
	point.CourseID = courseID
	if err := tx.QueryRow(ctx, `
		SELECT id, course_lesson_id, point_type, title, description, point_goal, point_category,
		       external_url, thumbnail_url, order_index, created_at, updated_at
		FROM course_points
		WHERE id = $1
		  AND course_id = $2
		  AND point_type = 'exploration'
		LIMIT 1
	`, pointID, courseID).Scan(
		&point.ID,
		&point.CourseLessonID,
		&point.PointType,
		&point.Title,
		&point.Description,
		&point.PointGoal,
		&point.PointCategory,
		&point.ExternalURL,
		&point.ThumbnailURL,
		&point.OrderIndex,
		&point.CreatedAt,
		&point.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, CoursePoint{}, errLearningPointNotFound
		}
		return uuid.Nil, uuid.Nil, CoursePoint{}, fmt.Errorf("get learning exploration point source: %w", err)
	}

	return draftID, courseID, point, nil
}

func (r *Repository) getLearningPointQuestionsTx(ctx context.Context, tx pgx.Tx, courseID, pointID uuid.UUID) ([]CoursePointQuestion, error) {
	rows, err := tx.Query(ctx, `
		SELECT cpq.id, cpq.course_point_id, cpq.goal_profile_version, cpq.title, cpq.question, cpq.question_type,
		       cpq.answer_method, cpq.answer, cpqf.feedback, cpq.status, cpq.created_by, cpq.created_at, cpq.updated_at
		FROM course_point_questions cpq
		LEFT JOIN course_point_question_ai_feedbacks cpqf
		  ON cpqf.course_point_question_id = cpq.id
		WHERE cpq.course_point_id = $1
		  AND cpq.course_point_id IN (
		    SELECT id
		    FROM course_points
		    WHERE id = $1
		      AND course_id = $2
		  )
		ORDER BY created_at ASC, updated_at ASC
	`, pointID, courseID)
	if err != nil {
		return nil, fmt.Errorf("get learning point questions: %w", err)
	}
	defer rows.Close()

	questions := make([]CoursePointQuestion, 0)
	for rows.Next() {
		var question CoursePointQuestion
		if err := rows.Scan(
			&question.ID,
			&question.CoursePointID,
			&question.GoalProfileVersion,
			&question.Title,
			&question.Question,
			&question.QuestionType,
			&question.AnswerMethod,
			&question.Answer,
			&question.AIFeedback,
			&question.Status,
			&question.CreatedBy,
			&question.CreatedAt,
			&question.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan learning point question: %w", err)
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning point questions: %w", err)
	}
	return questions, nil
}

func (r *Repository) getLearningPointCompletionQuestionsTx(ctx context.Context, tx pgx.Tx, courseID, pointID uuid.UUID) ([]CoursePointQuestion, error) {
	rows, err := tx.Query(ctx, `
		SELECT cpq.question, cpq.answer, cpq.status
		FROM course_point_questions cpq
		WHERE cpq.course_point_id = $1
		  AND cpq.course_point_id IN (
		    SELECT id
		    FROM course_points
		    WHERE id = $1
		      AND course_id = $2
		  )
	`, pointID, courseID)
	if err != nil {
		return nil, fmt.Errorf("get learning point completion questions: %w", err)
	}
	defer rows.Close()

	questions := make([]CoursePointQuestion, 0)
	for rows.Next() {
		var question CoursePointQuestion
		if err := rows.Scan(
			&question.Question,
			&question.Answer,
			&question.Status,
		); err != nil {
			return nil, fmt.Errorf("scan learning point completion question: %w", err)
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate learning point completion questions: %w", err)
	}
	return questions, nil
}

func (r *Repository) getLearningPointSelfEvaluationTx(ctx context.Context, tx pgx.Tx, courseID, pointID uuid.UUID) (*CoursePointSelfEvaluation, error) {
	var entry CoursePointSelfEvaluation
	if err := tx.QueryRow(ctx, `
		SELECT course_point_id, goal_profile_version, understanding, application_note,
		       proficiency, understanding_score, COALESCE(understanding_reason, ''),
		       application_score, COALESCE(application_reason, ''),
		       proficiency_score, COALESCE(proficiency_reason, ''),
		       problem_solving_score, COALESCE(problem_solving_reason, ''),
		       expression_score, COALESCE(expression_reason, ''),
		       goal_alignment_note, final_score, updated_at
		FROM course_point_self_evaluations
		WHERE course_point_id = $1
		  AND course_point_id IN (
		    SELECT id
		    FROM course_points
		    WHERE id = $1
		      AND course_id = $2
		  )
		LIMIT 1
	`, pointID, courseID).Scan(
		&entry.CoursePointID,
		&entry.GoalProfileVersion,
		&entry.Understanding,
		&entry.ApplicationNote,
		&entry.Proficiency,
		&entry.UnderstandingScore,
		&entry.UnderstandingReason,
		&entry.ApplicationScore,
		&entry.ApplicationReason,
		&entry.ProficiencyScore,
		&entry.ProficiencyReason,
		&entry.ProblemSolvingScore,
		&entry.ProblemSolvingReason,
		&entry.ExpressionScore,
		&entry.ExpressionReason,
		&entry.GoalAlignmentNote,
		&entry.FinalScore,
		&entry.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get learning point self evaluation: %w", err)
	}
	return &entry, nil
}
