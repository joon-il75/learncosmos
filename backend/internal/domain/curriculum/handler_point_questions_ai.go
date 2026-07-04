package curriculum

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

func (h *Handler) GenerateLearningPointQuestion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}

	if h.pointQuestionWorkerEnabled && h.llmGateway != nil {
		learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		job, err := h.llmGateway.SubmitAndWait(c.Request.Context(), CreatePointQuestionGenerateJobInput(userID, planetID, pointID, learningLanguage), h.pointQuestionWorkerSyncWait)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
			return
		}
		if job != nil && job.Status == "succeeded" {
			h.respondWithPointQuestionGenerateJobResult(c, userID, pointID, job)
			return
		}
		if job != nil && !terminalLLMJobStatus(job.Status) {
			c.JSON(http.StatusAccepted, gin.H{
				"error":      "point_question_generate_pending",
				"error_code": "point_question_generate_pending",
				"job_id":     job.ID.String(),
				"feature":    job.Feature,
				"status":     job.Status,
				"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
			})
			return
		}
		code := "point_question_generate_failed"
		if job != nil && job.ErrorCode != nil && strings.TrimSpace(*job.ErrorCode) != "" {
			code = strings.TrimSpace(*job.ErrorCode)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": code, "error_code": code})
		return
	}

	tx, err := h.repo.pool.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare ai question"})
		return
	}

	_, courseID, err := h.repo.getLearningPointContextTx(c.Request.Context(), tx, userID, planetID, pointID)
	if err != nil {
		tx.Rollback(c.Request.Context())
		if errors.Is(err, errLearningPointNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point context"})
		return
	}
	if err := tx.Rollback(c.Request.Context()); err != nil {
		log.Printf("[point-question] rollback after context load failed point=%s err=%v", pointID, err)
	}

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusConfirmed,
		DraftStatusLearning,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learning planet"})
		return
	}

	resolvedPoint, lessonTitle, found := findPointInPlanetLessons(planet.Lessons, pointID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		return
	}

	exec, err := h.resolveEvaluationExecution(c.Request.Context(), userID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve llm"})
		return
	}
	if strings.TrimSpace(exec.APIKey) == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no_llm_available"})
		return
	}
	learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)

	client, clientErr := llm.NewClient(exec.Provider, exec.APIKey, "")
	if clientErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create llm client"})
		return
	}
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_question_generate", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
	})

	llmCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	raw, llmErr := client.Complete(llmCtx, buildPointAIQuestionPrompt(planet.GoalContext, resolvedPoint, lessonTitle, learningLanguage))
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_question_failed"})
		return
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	suggestion := parsePointAIQuestionSuggestion(raw, resolvedPoint.Point.PointType, learningLanguage)
	courseID, question, err := h.repo.createLearningPointQuestionWithAuthor(
		c.Request.Context(),
		userID,
		planetID,
		pointID,
		nil,
		suggestion.Question,
		PointQuestionType(suggestion.QuestionType),
		nil,
		"ai",
	)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointQuestionInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point question is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ai question"})
		}
		return
	}
	h.observePointQuestionAIOutput(c.Request.Context(), userID, planetID, pointID, question.ID, question.Question, question.QuestionType, learningLanguage, "point_question_generate_sync")

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"question": question}, "ai question saved but failed to reload")
}

func (h *Handler) GenerateLearningPointQuestionFeedback(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	questionID, err := uuid.Parse(c.Param("question_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	if h.pointFeedbackWorkerEnabled && h.llmGateway != nil && c.Request.ContentLength == 0 {
		learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		job, err := h.llmGateway.SubmitAndWait(c.Request.Context(), CreatePointFeedbackGenerateJobInput(userID, planetID, pointID, questionID, learningLanguage), h.pointFeedbackWorkerSyncWait)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
			return
		}
		if job != nil && job.Status == "succeeded" {
			h.respondWithPointFeedbackGenerateJobResult(c, job)
			return
		}
		if job != nil && !terminalLLMJobStatus(job.Status) {
			c.JSON(http.StatusAccepted, gin.H{
				"error":      "point_feedback_generate_pending",
				"error_code": "point_feedback_generate_pending",
				"job_id":     job.ID.String(),
				"feature":    job.Feature,
				"status":     job.Status,
				"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
			})
			return
		}
		code := "point_feedback_generate_failed"
		if job != nil && job.ErrorCode != nil && strings.TrimSpace(*job.ErrorCode) != "" {
			code = strings.TrimSpace(*job.ErrorCode)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": code, "error_code": code})
		return
	}

	tx, err := h.repo.pool.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare ai feedback"})
		return
	}

	_, courseID, err := h.repo.getLearningPointContextTx(c.Request.Context(), tx, userID, planetID, pointID)
	if err != nil {
		tx.Rollback(c.Request.Context())
		if errors.Is(err, errLearningPointNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point context"})
		return
	}
	if err := tx.Rollback(c.Request.Context()); err != nil {
		log.Printf("[point-feedback] rollback after context load failed point=%s err=%v", pointID, err)
	}

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusConfirmed,
		DraftStatusLearning,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learning planet"})
		return
	}

	resolvedPoint, lessonTitle, found := findPointInPlanetLessons(planet.Lessons, pointID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		return
	}

	var targetQuestion *CoursePointQuestion
	for idx := range resolvedPoint.Questions {
		if resolvedPoint.Questions[idx].ID == questionID {
			targetQuestion = &resolvedPoint.Questions[idx]
			break
		}
	}
	if targetQuestion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning point question not found"})
		return
	}
	var req UpdateLearningPointQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	promptQuestion := *targetQuestion
	if req.Title != nil {
		promptQuestion.Title = nullableTrimmedString(*req.Title)
	}
	if req.Question != nil {
		promptQuestion.Question = strings.TrimSpace(*req.Question)
	}
	if req.QuestionType != nil {
		promptQuestion.QuestionType = *req.QuestionType
	}
	if req.Answer != nil {
		promptQuestion.Answer = nullableTrimmedString(*req.Answer)
	}

	exec, err := h.resolveEvaluationExecution(c.Request.Context(), userID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve llm"})
		return
	}
	if strings.TrimSpace(exec.APIKey) == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no_llm_available"})
		return
	}
	learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)

	client, clientErr := llm.NewClient(exec.Provider, exec.APIKey, "")
	if clientErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create llm client"})
		return
	}
	client = h.trackBYOKLLMUsageWithMetadata(client, userID, exec.Provider, "point_question_feedback", exec.BillingStatus, map[string]any{
		"learning_language": learningLanguage,
		"planet_id":         planetID.String(),
		"point_id":          pointID.String(),
		"question_id":       questionID.String(),
	})

	llmCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	raw, llmErr := client.Complete(llmCtx, buildPointAIFeedbackPrompt(planet.GoalContext, resolvedPoint, lessonTitle, promptQuestion, learningLanguage))
	if llmErr != nil {
		if exec.BillingStatus == "byok_no_charge" {
			go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, false, llmErr.Error())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_feedback_failed"})
		return
	}
	if exec.BillingStatus == "byok_no_charge" {
		go h.repo.RecordBYOKValidation(context.Background(), userID, exec.Provider, true, "")
	}

	feedback := normalizePointAIFeedback(raw, learningLanguage)
	h.observePointFeedbackAIOutput(c.Request.Context(), userID, planetID, pointID, questionID, feedback, learningLanguage, "point_feedback_generate_sync")
	c.JSON(http.StatusOK, gin.H{"feedback": feedback, "learning_language": learningLanguage})
}

func (h *Handler) GenerateLearningPointSelfEvaluationDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req LearningPointSelfEvaluationAIDraftRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid self evaluation draft request"})
			return
		}
	}

	planetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}

	if text, fields := learningPointSelfEvaluationApplicationAnswersModerationText(req.ApplicationAnswers); text != "" {
		learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_self_evaluation", &pointID, text, learningLanguage, map[string]any{
			"field":     "self_evaluation_application_answers",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	if h.pointSelfEvalDraftWorkerEnabled && h.llmGateway != nil {
		useStoredAnswers := len(req.ApplicationAnswers) > 0
		if useStoredAnswers {
			if err := h.repo.UpsertLearningPointSelfEvaluationApplicationAnswers(c.Request.Context(), userID, planetID, pointID, req.ApplicationAnswers); err != nil {
				if errors.Is(err, errLearningPointSelfEvalInvalidInput) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "self evaluation application answers invalid", "error_code": "self_evaluation_application_answers_invalid"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save self evaluation application answers", "error_code": "self_evaluation_application_answers_save_failed"})
				return
			}
		}
		learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		job, err := h.llmGateway.SubmitAndWait(c.Request.Context(), CreatePointSelfEvaluationDraftJobInput(userID, planetID, pointID, learningLanguage, useStoredAnswers), h.pointSelfEvalDraftWorkerSyncWait)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
			return
		}
		if job != nil && job.Status == "succeeded" {
			h.respondWithPointSelfEvaluationDraftJobResult(c, job)
			return
		}
		if job != nil && !terminalLLMJobStatus(job.Status) {
			c.JSON(http.StatusAccepted, gin.H{
				"error":      "point_self_evaluation_draft_pending",
				"error_code": "point_self_evaluation_draft_pending",
				"job_id":     job.ID.String(),
				"feature":    job.Feature,
				"status":     job.Status,
				"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
			})
			return
		}
		code := "point_self_evaluation_draft_failed"
		if job != nil && job.ErrorCode != nil && strings.TrimSpace(*job.ErrorCode) != "" {
			code = strings.TrimSpace(*job.ErrorCode)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": code, "error_code": code})
		return
	}

	result, err := h.executePointSelfEvaluationDraft(c.Request.Context(), userID, planetID, pointID, req.ApplicationAnswers)
	if err != nil {
		switch {
		case errors.Is(err, errNoPointAILLMAvailable):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no_llm_available"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		case errors.Is(err, errPointSelfEvaluationDraftFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_self_evaluation_draft_failed"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate self evaluation draft"})
		}
		return
	}

	h.observePointSelfEvaluationAIOutput(c.Request.Context(), userID, planetID, pointID, result.Draft, result.LearningLanguage, "point_self_evaluation_draft_sync")
	c.JSON(http.StatusOK, gin.H{
		"draft":             result.Draft,
		"learning_language": result.LearningLanguage,
	})
}
