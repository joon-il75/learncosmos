package curriculum

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListLearningPointEvents(c *gin.Context) {
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

	events, err := h.repo.ListLearningPointEvents(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list learning point events"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *Handler) CreateLearningPointPracticeLog(c *gin.Context) {
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

	var req CreateLearningPointPracticeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := createLearningPointPracticeLogModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_practice_log", &pointID, text, language, map[string]any{
			"field":     "practice_log_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, practiceLog, err := h.repo.CreateLearningPointPracticeLog(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointPracticeLogInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point practice log content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point practice log"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"practice_log": practiceLog}, "point practice log created but failed to reload")
}

func (h *Handler) UpdateLearningPointPracticeLog(c *gin.Context) {
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
	logID, err := uuid.Parse(c.Param("log_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid practice log id"})
		return
	}

	var req UpdateLearningPointPracticeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := updateLearningPointPracticeLogModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_practice_log", &logID, text, language, map[string]any{
			"field":     "practice_log_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
			"log_id":    logID.String(),
		}) {
			return
		}
	}

	courseID, practiceLog, err := h.repo.UpdateLearningPointPracticeLog(c.Request.Context(), userID, planetID, pointID, logID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointPracticeLogInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point practice log content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point practice log not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point practice log"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"practice_log": practiceLog}, "point practice log updated but failed to reload")
}

func createLearningPointPracticeLogModerationText(req CreateLearningPointPracticeLogRequest) (string, []string) {
	parts := make([]string, 0, 8)
	fields := make([]string, 0, 8)
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
	add("activity_name", req.ActivityName)
	add("blocked_part", req.BlockedPart)
	add("changed_method", req.ChangedMethod)
	add("achievement_note", req.AchievementNote)
	add("achievement", req.Achievement)
	add("next_plan", req.NextPlan)
	add("next_practice", req.NextPractice)
	return strings.Join(parts, "\n"), fields
}

func updateLearningPointPracticeLogModerationText(req UpdateLearningPointPracticeLogRequest) (string, []string) {
	parts := make([]string, 0, 8)
	fields := make([]string, 0, 8)
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
	add("activity_name", req.ActivityName)
	add("blocked_part", req.BlockedPart)
	add("changed_method", req.ChangedMethod)
	add("achievement_note", req.AchievementNote)
	add("achievement", req.Achievement)
	add("next_plan", req.NextPlan)
	add("next_practice", req.NextPractice)
	return strings.Join(parts, "\n"), fields
}

func (h *Handler) DeleteLearningPointPracticeLog(c *gin.Context) {
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
	logID, err := uuid.Parse(c.Param("log_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid practice log id"})
		return
	}

	courseID, err := h.repo.DeleteLearningPointPracticeLog(c.Request.Context(), userID, planetID, pointID, logID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point practice log not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete learning point practice log"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, nil, "point practice log deleted but failed to reload")
}
