package curriculum

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListSharedPlanets(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planets, err := h.repo.ListPlanetsByDraftStatuses(c.Request.Context(), userID, []DraftStatus{
		DraftStatusArchived,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shared planets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planets": planets})
}

func (h *Handler) GetSharedPlanet(c *gin.Context) {
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
		DraftStatusArchived,
	})
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "shared planet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": planet})
}

func (h *Handler) GetSharedPlanetPoint(c *gin.Context) {
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

	point, err := h.repo.GetPlanetPointDetailByIDAndDraftStatuses(c.Request.Context(), userID, courseID, pointID, sharedPointStatuses())
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "shared point not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planet": point.Planet, "point": point})
}

func (h *Handler) GetSharedPlanetRecords(c *gin.Context) {
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

	records, err := h.repo.GetPlanetRecordsByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusArchived,
	})
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "shared planet records not found"})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *Handler) GetSharedPlanetRecordFeed(c *gin.Context) {
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

	feed, err := h.repo.GetPlanetRecordFeedByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusArchived,
	}, "shared")
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "shared planet record feed not found"})
		return
	}

	c.JSON(http.StatusOK, feed)
}

func (h *Handler) GetSharedPlanetResults(c *gin.Context) {
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

	results, err := h.repo.GetPlanetResultsByIDAndDraftStatuses(c.Request.Context(), userID, courseID, []DraftStatus{
		DraftStatusArchived,
	})
	if err != nil {
		if errors.Is(err, errInactivePlanetAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "inactive planet only allows planning page"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "shared planet results not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}
