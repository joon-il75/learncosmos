package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateLearningPointQuestion(ctx context.Context, userID, planetID, pointID uuid.UUID, req CreateLearningPointQuestionRequest) (uuid.UUID, CoursePointQuestion, error) {
	return r.createLearningPointQuestionWithAuthor(ctx, userID, planetID, pointID, req.Title, req.Question, req.QuestionType, req.AnswerMethod, "learner")
}

func (r *Repository) createLearningPointQuestionWithAuthor(ctx context.Context, userID, planetID, pointID uuid.UUID, titleValue *string, questionText string, questionType PointQuestionType, answerMethodValue *string, createdBy string) (uuid.UUID, CoursePointQuestion, error) {
	title := trimOptionalString(titleValue)
	question := strings.TrimSpace(questionText)
	answerMethod := normalizePointQuestionAnswerMethod(answerMethodValue)
	normalizedCreatedBy := strings.TrimSpace(createdBy)
	if normalizedCreatedBy == "" {
		normalizedCreatedBy = "learner"
	}
	questionType, ok := normalizePointQuestionType(questionType)
	if !ok || question == "" {
		return uuid.Nil, CoursePointQuestion{}, errLearningPointQuestionInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("begin learning point question tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointQuestion{}, err
	}
	goalProfileVersion, err := r.getDraftGoalProfileVersionTx(ctx, tx, draftID)
	if err != nil {
		return uuid.Nil, CoursePointQuestion{}, err
	}

	var entry CoursePointQuestion
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_questions (
			course_point_id, goal_profile_version, title, question, question_type, answer_method, answer, status, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, NULL, 'pending', $7)
		RETURNING id, course_point_id, goal_profile_version, title, question, question_type, answer_method, answer, status, created_by, created_at, updated_at
	`, pointID, goalProfileVersion, nullableTrimmedString(title), question, questionType, answerMethod, normalizedCreatedBy).Scan(
		&entry.ID,
		&entry.CoursePointID,
		&entry.GoalProfileVersion,
		&entry.Title,
		&entry.Question,
		&entry.QuestionType,
		&entry.AnswerMethod,
		&entry.Answer,
		&entry.Status,
		&entry.CreatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("insert learning point question: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "question_added", map[string]any{"question_id": entry.ID, "created_by": entry.CreatedBy, "title": entry.Title, "question": entry.Question}); err != nil {
		return uuid.Nil, CoursePointQuestion{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("touch learning draft after point question: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("touch learning course after point question: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("commit learning point question tx: %w", err)
	}

	return courseID, entry, nil
}

func (r *Repository) UpdateLearningPointQuestion(ctx context.Context, userID, planetID, pointID, questionID uuid.UUID, req UpdateLearningPointQuestionRequest) (uuid.UUID, CoursePointQuestion, error) {
	if req.Title == nil && req.Question == nil && req.QuestionType == nil && req.AnswerMethod == nil && req.Answer == nil && req.AIFeedback == nil && req.Status == nil {
		return uuid.Nil, CoursePointQuestion{}, errLearningPointQuestionInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("begin learning point question update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointQuestion{}, err
	}

	var answer *string
	var title *string
	var question *string
	var questionType *PointQuestionType
	var answerMethod *string
	answerProvided := req.Answer != nil
	statusProvided := req.Status != nil
	status := PointQuestionStatusPending
	if req.Title != nil {
		title = nullableTrimmedString(*req.Title)
	}
	if req.Question != nil {
		trimmed := strings.TrimSpace(*req.Question)
		if trimmed == "" {
			return uuid.Nil, CoursePointQuestion{}, errLearningPointQuestionInvalidInput
		}
		question = &trimmed
	}
	if req.QuestionType != nil {
		normalizedType, ok := normalizePointQuestionType(*req.QuestionType)
		if !ok {
			return uuid.Nil, CoursePointQuestion{}, errLearningPointQuestionInvalidInput
		}
		questionType = &normalizedType
	}
	if req.AnswerMethod != nil {
		answerMethod = normalizePointQuestionAnswerMethod(req.AnswerMethod)
	}
	if req.Answer != nil {
		trimmed := strings.TrimSpace(*req.Answer)
		if trimmed != "" {
			answer = &trimmed
			status = PointQuestionStatusAnswered
		}
	}
	if req.Status != nil {
		normalizedStatus, ok := normalizePointQuestionStatus(*req.Status)
		if !ok {
			return uuid.Nil, CoursePointQuestion{}, errLearningPointQuestionInvalidInput
		}
		status = normalizedStatus
		if status == PointQuestionStatusPending && req.Answer == nil {
			answer = nil
		}
	}

	var entry CoursePointQuestion
	if err := tx.QueryRow(ctx, `
		UPDATE course_point_questions q
		SET title = COALESCE($1, q.title),
		    question = COALESCE($2, q.question),
		    question_type = COALESCE($3, q.question_type),
		    answer_method = COALESCE($4, q.answer_method),
		    answer = CASE WHEN $10 THEN $5 ELSE q.answer END,
		    status = CASE WHEN $11 THEN $6 ELSE q.status END,
		    updated_at = NOW()
		FROM course_points cp
		WHERE q.id = $7
		  AND q.course_point_id = cp.id
		  AND cp.id = $8
		  AND cp.course_id = $9
		RETURNING q.id, q.course_point_id, q.goal_profile_version, q.title, q.question, q.question_type, q.answer_method, q.answer, q.status, q.created_by, q.created_at, q.updated_at
	`, title, question, questionType, answerMethod, answer, status, questionID, pointID, courseID, answerProvided, statusProvided || answerProvided).Scan(
		&entry.ID,
		&entry.CoursePointID,
		&entry.GoalProfileVersion,
		&entry.Title,
		&entry.Question,
		&entry.QuestionType,
		&entry.AnswerMethod,
		&entry.Answer,
		&entry.Status,
		&entry.CreatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, CoursePointQuestion{}, errLearningPointNotFound
		}
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("update learning point question: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "question_answered", map[string]any{"question_id": entry.ID, "status": entry.Status, "title": entry.Title, "question": entry.Question}); err != nil {
		return uuid.Nil, CoursePointQuestion{}, err
	}
	if req.AIFeedback != nil {
		trimmedFeedback := strings.TrimSpace(*req.AIFeedback)
		if trimmedFeedback == "" {
			if _, err := tx.Exec(ctx, `
				DELETE FROM course_point_question_ai_feedbacks
				WHERE course_point_question_id = $1
			`, entry.ID); err != nil {
				return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("delete learning point question ai feedback: %w", err)
			}
			entry.AIFeedback = nil
		} else {
			if err := tx.QueryRow(ctx, `
				INSERT INTO course_point_question_ai_feedbacks (
					course_point_question_id, feedback
				) VALUES ($1, $2)
				ON CONFLICT (course_point_question_id)
				DO UPDATE SET
					feedback = EXCLUDED.feedback,
					updated_at = NOW()
				RETURNING feedback
			`, entry.ID, trimmedFeedback).Scan(&trimmedFeedback); err != nil {
				return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("upsert learning point question ai feedback during update: %w", err)
			}
			entry.AIFeedback = &trimmedFeedback
		}
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("touch learning draft after point question update: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("touch learning course after point question update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointQuestion{}, fmt.Errorf("commit learning point question update tx: %w", err)
	}

	return courseID, entry, nil
}

func (r *Repository) DeleteLearningPointQuestion(ctx context.Context, userID, planetID, pointID, questionID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin learning point question delete tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, err
	}

	var deletedTitle *string
	var deletedQuestion string
	if err := tx.QueryRow(ctx, `
		SELECT q.title, q.question
		FROM course_point_questions q
		JOIN course_points cp ON cp.id = q.course_point_id
		WHERE q.id = $1
		  AND cp.id = $2
		  AND cp.course_id = $3
	`, questionID, pointID, courseID).Scan(&deletedTitle, &deletedQuestion); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errLearningPointNotFound
		}
		return uuid.Nil, fmt.Errorf("get learning point question before delete: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM course_point_question_ai_feedbacks
		WHERE course_point_question_id = $1
	`, questionID); err != nil {
		return uuid.Nil, fmt.Errorf("delete learning point question ai feedback before question delete: %w", err)
	}

	tag, err := tx.Exec(ctx, `
		DELETE FROM course_point_questions q
		USING course_points cp
		WHERE q.id = $1
		  AND q.course_point_id = cp.id
		  AND cp.id = $2
		  AND cp.course_id = $3
	`, questionID, pointID, courseID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("delete learning point question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return uuid.Nil, errLearningPointNotFound
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "question_added", map[string]any{"question_id": questionID, "action": "deleted", "title": deletedTitle, "question": deletedQuestion}); err != nil {
		return uuid.Nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning draft after point question delete: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, fmt.Errorf("touch learning course after point question delete: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit learning point question delete tx: %w", err)
	}
	return courseID, nil
}

func (r *Repository) UpsertLearningPointQuestionAIFeedback(ctx context.Context, userID, planetID, pointID, questionID uuid.UUID, feedback string) (uuid.UUID, string, error) {
	trimmedFeedback := strings.TrimSpace(feedback)
	if trimmedFeedback == "" {
		return uuid.Nil, "", errLearningPointQuestionInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("begin learning point question ai feedback tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, "", err
	}

	var savedFeedback string
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_question_ai_feedbacks (
			course_point_question_id, feedback
		)
		SELECT q.id, $1
		FROM course_point_questions q
		JOIN course_points cp ON cp.id = q.course_point_id
		WHERE q.id = $2
		  AND cp.id = $3
		  AND cp.course_id = $4
		ON CONFLICT (course_point_question_id)
		DO UPDATE SET
			feedback = EXCLUDED.feedback,
			updated_at = NOW()
		RETURNING feedback
	`, trimmedFeedback, questionID, pointID, courseID).Scan(&savedFeedback); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, "", errLearningPointNotFound
		}
		return uuid.Nil, "", fmt.Errorf("upsert learning point question ai feedback: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "ai_hint_used", map[string]any{"question_id": questionID}); err != nil {
		return uuid.Nil, "", err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, "", fmt.Errorf("touch learning draft after point question ai feedback: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, "", fmt.Errorf("touch learning course after point question ai feedback: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, "", fmt.Errorf("commit learning point question ai feedback tx: %w", err)
	}

	return courseID, savedFeedback, nil
}

func (r *Repository) UpsertLearningPointSelfEvaluation(ctx context.Context, userID, planetID, pointID uuid.UUID, req UpsertLearningPointSelfEvaluationRequest) (uuid.UUID, CoursePointSelfEvaluation, error) {
	understandingScore := req.UnderstandingScore
	if understandingScore == nil {
		understandingScore = req.Understanding
	}
	proficiencyScore := req.ProficiencyScore
	if proficiencyScore == nil {
		proficiencyScore = req.Proficiency
	}
	if understandingScore == nil || proficiencyScore == nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, errLearningPointSelfEvalInvalidInput
	}

	understanding := *understandingScore
	proficiency := *proficiencyScore
	if !validScorePtr(understandingScore) ||
		!validScorePtr(req.ApplicationScore) ||
		!validScorePtr(proficiencyScore) ||
		!validScorePtr(req.ProblemSolvingScore) ||
		!validScorePtr(req.ExpressionScore) ||
		!validScorePtr(req.FinalScore) {
		return uuid.Nil, CoursePointSelfEvaluation{}, errLearningPointSelfEvalInvalidInput
	}

	applicationNote := trimOptionalString(req.ApplicationNote)
	understandingReason := trimOptionalString(req.UnderstandingReason)
	applicationReason := trimOptionalString(req.ApplicationReason)
	proficiencyReason := trimOptionalString(req.ProficiencyReason)
	problemSolvingReason := trimOptionalString(req.ProblemSolvingReason)
	expressionReason := trimOptionalString(req.ExpressionReason)
	if applicationNote == "" {
		applicationNote = applicationReason
	}
	goalAlignmentNote := trimOptionalString(req.GoalAlignmentNote)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, fmt.Errorf("begin learning point self evaluation tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, courseID, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, err
	}
	goalProfileVersion, err := r.getDraftGoalProfileVersionTx(ctx, tx, draftID)
	if err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, err
	}

	var entry CoursePointSelfEvaluation
	if err := tx.QueryRow(ctx, `
		INSERT INTO course_point_self_evaluations (
			course_point_id, goal_profile_version, understanding, application_note, proficiency,
			understanding_score, understanding_reason, application_score, application_reason,
			proficiency_score, proficiency_reason, problem_solving_score, problem_solving_reason,
			expression_score, expression_reason, goal_alignment_note, final_score, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW())
		ON CONFLICT (course_point_id)
		DO UPDATE SET
			goal_profile_version = EXCLUDED.goal_profile_version,
			understanding = EXCLUDED.understanding,
			application_note = EXCLUDED.application_note,
			proficiency = EXCLUDED.proficiency,
			understanding_score = EXCLUDED.understanding_score,
			understanding_reason = EXCLUDED.understanding_reason,
			application_score = EXCLUDED.application_score,
			application_reason = EXCLUDED.application_reason,
			proficiency_score = EXCLUDED.proficiency_score,
			proficiency_reason = EXCLUDED.proficiency_reason,
			problem_solving_score = EXCLUDED.problem_solving_score,
			problem_solving_reason = EXCLUDED.problem_solving_reason,
			expression_score = EXCLUDED.expression_score,
			expression_reason = EXCLUDED.expression_reason,
			goal_alignment_note = EXCLUDED.goal_alignment_note,
			final_score = EXCLUDED.final_score,
			updated_at = NOW()
		RETURNING course_point_id, goal_profile_version, understanding, application_note, proficiency,
		          understanding_score, COALESCE(understanding_reason, ''),
		          application_score, COALESCE(application_reason, ''),
		          proficiency_score, COALESCE(proficiency_reason, ''),
		          problem_solving_score, COALESCE(problem_solving_reason, ''),
		          expression_score, COALESCE(expression_reason, ''),
		          goal_alignment_note, final_score, updated_at
	`, pointID, goalProfileVersion, understanding, applicationNote, proficiency,
		understandingScore, understandingReason, req.ApplicationScore, applicationReason,
		proficiencyScore, proficiencyReason, req.ProblemSolvingScore, problemSolvingReason,
		req.ExpressionScore, expressionReason, goalAlignmentNote, req.FinalScore).Scan(
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
		return uuid.Nil, CoursePointSelfEvaluation{}, fmt.Errorf("upsert learning point self evaluation: %w", err)
	}
	if err := r.insertPointEventTx(ctx, tx, userID, pointID, "self_evaluation_saved", map[string]any{
		"understanding": entry.Understanding,
		"proficiency":   entry.Proficiency,
		"final_score":   entry.FinalScore,
	}); err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE course_drafts SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, draftID, userID); err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, fmt.Errorf("touch learning draft after self evaluation: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE courses SET updated_at = NOW() WHERE id = $1 AND user_id = $2`, courseID, userID); err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, fmt.Errorf("touch learning course after self evaluation: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, CoursePointSelfEvaluation{}, fmt.Errorf("commit learning point self evaluation tx: %w", err)
	}

	return courseID, entry, nil
}

func (r *Repository) UpsertLearningPointSelfEvaluationApplicationAnswers(ctx context.Context, userID, planetID, pointID uuid.UUID, answers []LearningPointSelfEvaluationApplicationAnswer) error {
	normalized := normalizeSelfEvaluationApplicationAnswers(answers)
	if len(normalized) == 0 {
		return errLearningPointSelfEvalInvalidInput
	}
	payload, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("marshal self evaluation application answers: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin self evaluation application answer tx: %w", err)
	}
	defer tx.Rollback(ctx)

	draftID, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID)
	if err != nil {
		return err
	}
	goalProfileVersion, err := r.getDraftGoalProfileVersionTx(ctx, tx, draftID)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO course_point_self_evaluation_application_answers (
			id, course_point_id, user_id, goal_profile_version, answers, updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (course_point_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			goal_profile_version = EXCLUDED.goal_profile_version,
			answers = EXCLUDED.answers,
			updated_at = NOW()
	`, uuid.New(), pointID, userID, goalProfileVersion, payload); err != nil {
		return fmt.Errorf("upsert self evaluation application answers: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit self evaluation application answer tx: %w", err)
	}
	return nil
}

func (r *Repository) GetLearningPointSelfEvaluationApplicationAnswers(ctx context.Context, userID, planetID, pointID uuid.UUID) ([]LearningPointSelfEvaluationApplicationAnswer, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin self evaluation application answer read tx: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, _, err := r.getLearningPointContextTx(ctx, tx, userID, planetID, pointID); err != nil {
		return nil, err
	}
	var raw []byte
	if err := tx.QueryRow(ctx, `
		SELECT answers
		FROM course_point_self_evaluation_application_answers
		WHERE course_point_id = $1 AND user_id = $2
	`, pointID, userID).Scan(&raw); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("load self evaluation application answers: %w", err)
	}
	var answers []LearningPointSelfEvaluationApplicationAnswer
	if err := json.Unmarshal(raw, &answers); err != nil {
		return nil, fmt.Errorf("decode self evaluation application answers: %w", err)
	}
	return normalizeSelfEvaluationApplicationAnswers(answers), nil
}

func normalizeSelfEvaluationApplicationAnswers(answers []LearningPointSelfEvaluationApplicationAnswer) []LearningPointSelfEvaluationApplicationAnswer {
	if len(answers) == 0 {
		return nil
	}
	out := make([]LearningPointSelfEvaluationApplicationAnswer, 0, len(answers))
	for _, item := range answers {
		question := strings.TrimSpace(item.Question)
		answer := strings.TrimSpace(item.Answer)
		if question == "" || answer == "" {
			continue
		}
		if len([]rune(question)) > 1000 {
			question = string([]rune(question)[:1000])
		}
		if len([]rune(answer)) > 5000 {
			answer = string([]rune(answer)[:5000])
		}
		out = append(out, LearningPointSelfEvaluationApplicationAnswer{Question: question, Answer: answer})
		if len(out) >= 3 {
			break
		}
	}
	return out
}
