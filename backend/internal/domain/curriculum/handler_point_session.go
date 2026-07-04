package curriculum

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateLearningPointSession(c *gin.Context) {
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
	var req StartLearningPointSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.UserAgent) == "" {
		req.UserAgent = c.Request.UserAgent()
	}
	_, session, err := h.repo.CreateLearningPointSession(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point session"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"session": session})
}

func (h *Handler) HeartbeatLearningPointSession(c *gin.Context) {
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
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	var req UpdateLearningPointSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, session, err := h.repo.HeartbeatLearningPointSession(c.Request.Context(), userID, planetID, pointID, sessionID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to heartbeat learning point session"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *Handler) EndLearningPointSession(c *gin.Context) {
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
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	var req UpdateLearningPointSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, session, err := h.repo.EndLearningPointSession(c.Request.Context(), userID, planetID, pointID, sessionID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to end learning point session"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": session})
}
