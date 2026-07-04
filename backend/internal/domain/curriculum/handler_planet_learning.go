package curriculum

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListLearningPlanets(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planets, err := h.repo.ListPlanetsByDraftStatuses(c.Request.Context(), userID, []DraftStatus{
		DraftStatusConfirmed,
		DraftStatusLearning,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list learning planets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planets": planets})
}

func (h *Handler) GetLearningPlanet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.TouchPlanetLastAccessed(c.Request.Context(), userID, courseID); err != nil && !errors.Is(err, errDraftResourceNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet access"})
		return
	}

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusConfirmed,
		DraftStatusLearning,
	})
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "learning planet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}

func (h *Handler) GetLearningPlanetPoint(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}

	if err := h.repo.TouchPlanetLastAccessed(c.Request.Context(), userID, courseID); err != nil && !errors.Is(err, errDraftResourceNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet access"})
		return
	}

	point, err := h.repo.GetPlanetPointDetailByIDAndDraftStatuses(c.Request.Context(), userID, courseID, pointID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": point.Planet, "point": point})
}

func (h *Handler) GetLearningPlanetRecords(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.TouchPlanetLastAccessed(c.Request.Context(), userID, courseID); err != nil && !errors.Is(err, errDraftResourceNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet access"})
		return
	}

	records, err := h.repo.GetPlanetRecordsByIDAndDraftStatuses(c.Request.Context(), userID, courseID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning planet records not found"})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *Handler) GetLearningPlanetRecordFeed(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.TouchPlanetLastAccessed(c.Request.Context(), userID, courseID); err != nil && !errors.Is(err, errDraftResourceNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet access"})
		return
	}

	feed, err := h.repo.GetPlanetRecordFeedByIDAndDraftStatuses(c.Request.Context(), userID, courseID, learningPointStatuses(), "learning")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning planet record feed not found"})
		return
	}

	c.JSON(http.StatusOK, feed)
}

func (h *Handler) GetLearningPlanetResults(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.TouchPlanetLastAccessed(c.Request.Context(), userID, courseID); err != nil && !errors.Is(err, errDraftResourceNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update planet access"})
		return
	}

	results, err := h.repo.GetPlanetResultsByIDAndDraftStatuses(c.Request.Context(), userID, courseID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "learning planet results not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) StartLearningPlanet(c *gin.Context) {
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

	if err := h.repo.StartLearningPlanet(c.Request.Context(), userID, planetID); err != nil {
		switch {
		case errors.Is(err, errDraftResourceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning planet not found"})
		case errors.Is(err, errDraftStartNotAllowed):
			c.JSON(http.StatusBadRequest, gin.H{"error": "planet is not ready to start"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start learning planet"})
		}
		return
	}

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(c.Request.Context(), userID, planetID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "planet started but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}

func (h *Handler) CompleteLearningPlanet(c *gin.Context) {
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

	if err := h.repo.CompleteLearningPlanet(c.Request.Context(), userID, planetID); err != nil {
		switch {
		case errors.Is(err, errDraftResourceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning planet not found"})
		case errors.Is(err, errDraftCompleteNotAllowed):
			c.JSON(http.StatusBadRequest, gin.H{"error": "planet is not in learning state"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete planet"})
		}
		return
	}

	planet, err := h.repo.GetPlanetByIDAndDraftStatuses(c.Request.Context(), userID, planetID, []DraftStatus{
		DraftStatusArchived,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "planet completed but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}
