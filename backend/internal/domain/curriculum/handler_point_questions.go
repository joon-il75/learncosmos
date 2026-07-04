package curriculum

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateLearningPointQuestion(c *gin.Context) {
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

	var req CreateLearningPointQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if text, fields := createLearningPointQuestionModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_question", &pointID, text, language, map[string]any{
			"field":     "question_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, question, err := h.repo.CreateLearningPointQuestion(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointQuestionInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point question is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point question"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"question": question}, "point question created but failed to reload")
}

func (h *Handler) UpdateLearningPointQuestion(c *gin.Context) {
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

	var req UpdateLearningPointQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if text, fields := updateLearningPointQuestionModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_question", &questionID, text, language, map[string]any{
			"field":       "question_content",
			"fields":      fields,
			"planet_id":   planetID.String(),
			"point_id":    pointID.String(),
			"question_id": questionID.String(),
		}) {
			return
		}
	}

	courseID, question, err := h.repo.UpdateLearningPointQuestion(c.Request.Context(), userID, planetID, pointID, questionID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointQuestionInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point question update is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point question"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"question": question}, "point question saved but failed to reload")
}

func createLearningPointQuestionModerationText(req CreateLearningPointQuestionRequest) (string, []string) {
	parts := make([]string, 0, 3)
	fields := make([]string, 0, 3)
	add := func(field string, value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return
		}
		parts = append(parts, trimmed)
		fields = append(fields, field)
	}
	if req.Title != nil {
		add("title", *req.Title)
	}
	add("question", req.Question)
	if req.AnswerMethod != nil {
		add("answer_method", *req.AnswerMethod)
	}
	return strings.Join(parts, "\n"), fields
}

func updateLearningPointQuestionModerationText(req UpdateLearningPointQuestionRequest) (string, []string) {
	parts := make([]string, 0, 4)
	fields := make([]string, 0, 4)
	add := func(field string, value *string) {
		if value == nil {
			return
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			return
		}
		parts = append(parts, trimmed)
		fields = append(fields, field)
	}
	add("title", req.Title)
	add("question", req.Question)
	add("answer_method", req.AnswerMethod)
	add("answer", req.Answer)
	return strings.Join(parts, "\n"), fields
}

func (h *Handler) DeleteLearningPointQuestion(c *gin.Context) {
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

	courseID, err := h.repo.DeleteLearningPointQuestion(c.Request.Context(), userID, planetID, pointID, questionID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point question not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete learning point question"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"deleted": true}, "point question deleted but failed to reload")
}

func (h *Handler) UpsertLearningPointSelfEvaluation(c *gin.Context) {
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

	var req UpsertLearningPointSelfEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if text, fields := upsertLearningPointSelfEvaluationModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_self_evaluation", &pointID, text, language, map[string]any{
			"field":     "self_evaluation_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, selfEvaluation, err := h.repo.UpsertLearningPointSelfEvaluation(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointSelfEvalInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point self evaluation is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save learning point self evaluation"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"self_evaluation": selfEvaluation}, "point self evaluation saved but failed to reload")
}

func upsertLearningPointSelfEvaluationModerationText(req UpsertLearningPointSelfEvaluationRequest) (string, []string) {
	parts := make([]string, 0, 7)
	fields := make([]string, 0, 7)
	add := func(field string, value *string) {
		if value == nil {
			return
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			return
		}
		parts = append(parts, trimmed)
		fields = append(fields, field)
	}
	add("application_note", req.ApplicationNote)
	add("understanding_reason", req.UnderstandingReason)
	add("application_reason", req.ApplicationReason)
	add("proficiency_reason", req.ProficiencyReason)
	add("problem_solving_reason", req.ProblemSolvingReason)
	add("expression_reason", req.ExpressionReason)
	add("goal_alignment_note", req.GoalAlignmentNote)
	return strings.Join(parts, "\n"), fields
}

func learningPointSelfEvaluationApplicationAnswersModerationText(answers []LearningPointSelfEvaluationApplicationAnswer) (string, []string) {
	parts := make([]string, 0, len(answers)*2)
	fields := make([]string, 0, len(answers)*2)
	for _, answer := range answers {
		question := strings.TrimSpace(answer.Question)
		if question != "" {
			parts = append(parts, question)
			fields = append(fields, "application_question")
		}
		text := strings.TrimSpace(answer.Answer)
		if text != "" {
			parts = append(parts, text)
			fields = append(fields, "application_answer")
		}
	}
	return strings.Join(parts, "\n"), fields
}
