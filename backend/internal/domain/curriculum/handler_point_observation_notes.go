package curriculum

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateLearningPointObservationNote(c *gin.Context) {
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

	var req CreateLearningPointObservationNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	if h.enforceSafety(c, userID, "point_observation_note", &pointID, req.Content, language, map[string]any{"field": "content", "note_type": req.NoteType, "planet_id": planetID.String(), "point_id": pointID.String()}) {
		return
	}

	courseID, note, err := h.repo.CreateLearningPointObservationNote(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointObservationInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "observation note type and content are required"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create observation note"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"observation_note": note}, "observation note created but failed to reload")
}

func (h *Handler) UpdateLearningPointObservationNote(c *gin.Context) {
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
	noteID, err := uuid.Parse(c.Param("note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req UpdateLearningPointObservationNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Content != nil {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		noteType := ""
		if req.NoteType != nil {
			noteType = *req.NoteType
		}
		if h.enforceSafety(c, userID, "point_observation_note", &noteID, *req.Content, language, map[string]any{"field": "content", "note_type": noteType, "planet_id": planetID.String(), "point_id": pointID.String(), "note_id": noteID.String()}) {
			return
		}
	}

	courseID, note, err := h.repo.UpdateLearningPointObservationNote(c.Request.Context(), userID, planetID, pointID, noteID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointObservationInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "observation note type and content are required"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "observation note not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update observation note"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"observation_note": note}, "observation note updated but failed to reload")
}

func (h *Handler) DeleteLearningPointObservationNote(c *gin.Context) {
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
	noteID, err := uuid.Parse(c.Param("note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	courseID, err := h.repo.DeleteLearningPointObservationNote(c.Request.Context(), userID, planetID, pointID, noteID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "observation note not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete observation note"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"deleted": true}, "observation note deleted but failed to reload")
}
