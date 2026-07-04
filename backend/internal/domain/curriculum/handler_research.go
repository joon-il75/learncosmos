package curriculum

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ── Sub Lesson ─────────────────────────────────────────────────────────────────

func (h *Handler) CreateSubLesson(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	mainLessonID, err := uuid.Parse(c.Param("main_lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid main lesson id"})
		return
	}
	var req CreateSubLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	lesson, err := h.repo.CreateSubLesson(c.Request.Context(), userID, draftID, mainLessonID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "draft or main lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sub lesson"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"lesson": lesson})
}

// ── Research Node ──────────────────────────────────────────────────────────────

func (h *Handler) CreateResearchNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	levelID, err := uuid.Parse(c.Param("main_lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid main lesson id"})
		return
	}
	var req CreateResearchNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	node, err := h.repo.CreateResearchNode(c.Request.Context(), userID, draftID, levelID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "draft or main lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create point"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"research_node": node, "point": node})
}

func (h *Handler) UpdateResearchNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	var req UpdateResearchNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	node, err := h.repo.UpdateResearchNode(c.Request.Context(), userID, draftID, nodeID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update point"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"research_node": node, "point": node})
}

func (h *Handler) DeleteResearchNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	if err := h.repo.DeleteResearchNode(c.Request.Context(), userID, draftID, nodeID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete point"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetResearchNodeBlocks(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	blocks, err := h.repo.GetResearchNodeBlocks(c.Request.Context(), userID, draftID, nodeID)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blocks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"blocks": blocks, "point_blocks": blocks})
}

func (h *Handler) CreateResearchNodeBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	var req CreateResearchNodeBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	block, err := h.repo.CreateResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create block"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"block": block, "point_block": block})
}

func (h *Handler) UpdateResearchNodeBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	blockID, err := uuid.Parse(c.Param("block_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block id"})
		return
	}
	var req UpdateResearchNodeBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	block, err := h.repo.UpdateResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, blockID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update block"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block": block})
}

func (h *Handler) DeleteResearchNodeBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}
	blockID, err := uuid.Parse(c.Param("block_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block id"})
		return
	}
	if err := h.repo.DeleteResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, blockID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete block"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Research Point ─────────────────────────────────────────────────────────────

func (h *Handler) UpdateResearchPoint(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	var req UpdateResearchNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	node, err := h.repo.UpdateResearchNode(c.Request.Context(), userID, draftID, nodeID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update point"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"point": node})
}

func (h *Handler) DeleteResearchPoint(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	if err := h.repo.DeleteResearchNode(c.Request.Context(), userID, draftID, nodeID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete point"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetResearchPointBlocks(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	blocks, err := h.repo.GetResearchNodeBlocks(c.Request.Context(), userID, draftID, nodeID)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blocks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"blocks": blocks, "point_blocks": blocks})
}

func (h *Handler) CreateResearchPointBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	var req CreateResearchNodeBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	block, err := h.repo.CreateResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create block"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"block": block, "point_block": block})
}

func (h *Handler) UpdateResearchPointBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	blockID, err := uuid.Parse(c.Param("block_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block id"})
		return
	}
	var req UpdateResearchNodeBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	block, err := h.repo.UpdateResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, blockID, req)
	if err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update block"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block": block, "point_block": block})
}

func (h *Handler) DeleteResearchPointBlock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	nodeID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	blockID, err := uuid.Parse(c.Param("block_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block id"})
		return
	}
	if err := h.repo.DeleteResearchNodeBlock(c.Request.Context(), userID, draftID, nodeID, blockID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete block"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) CompleteResearchPoint(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	if err := h.repo.CompleteResearchPoint(c.Request.Context(), userID, draftID, pointID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete point"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) UncompleteResearchPoint(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}
	pointID, err := uuid.Parse(c.Param("point_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point id"})
		return
	}
	if err := h.repo.UncompleteResearchPoint(c.Request.Context(), userID, draftID, pointID); err != nil {
		if errors.Is(err, errDraftResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to uncomplete point"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
