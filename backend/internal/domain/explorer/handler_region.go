package explorer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
)

func (h *Handler) GetCourse(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	includeInactive := c.Query("include_inactive") == "true"
	aggregate, err := h.repo.GetCourseAggregate(c.Request.Context(), courseID, includeInactive)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": aggregate})
}

func (h *Handler) CreateRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !h.verifyCourseOwner(c, req.CourseDraftID, userID) {
		return
	}

	newLesson, err := h.curriculumRepo.CreateMainLesson(c.Request.Context(), userID, req.CourseDraftID, curriculum.CreateMainLessonRequest{
		Title:      req.Name,
		Objective:  req.Description,
		OrderIndex: req.OrderIndex,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create region lesson"})
		return
	}

	regionLessonID := newLesson.ID
	region, err := h.repo.InsertRegion(c.Request.Context(), req, &regionLessonID)
	if err != nil {
		_ = h.repo.DeleteRegionDraftLesson(c.Request.Context(), req.CourseDraftID, userID, regionLessonID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create region"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"region": region})
}

func (h *Handler) UpdateRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	regionID, err := uuid.Parse(c.Param("regionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), regionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	if !h.verifyCourseOwner(c, region.CourseDraftID, userID) {
		return
	}

	var req UpdateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.UpdateRegion(c.Request.Context(), regionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update region"})
		return
	}

	if updated.CourseDraftLessonID != nil {
		nextTitle := (*string)(nil)
		nextTitleText := strings.TrimSpace(updated.Name)
		if req.Name != nil {
			title := strings.TrimSpace(*req.Name)
			if title != "" {
				nextTitleText = title
				nextTitle = &title
			}
		}

		nextObjective := (*string)(nil)
		nextSummary := (*string)(nil)
		if req.Description != nil {
			desc := strings.TrimSpace(*req.Description)
			nextObjective = &desc
			nextSummary = &desc
		}

		var nextOrder *int
		newOrderIndex := updated.OrderIndex
		if req.OrderIndex != nil {
			nextOrder = req.OrderIndex
			newOrderIndex = *req.OrderIndex
		}

		syncLesson := req.Name != nil || req.Description != nil || req.OrderIndex != nil
		var lessonSpec *string
		if syncLesson {
			lessonSpec = h.buildSubRegionLessonSpec(c.Request.Context(), userID, *updated.CourseDraftLessonID, nextTitleText, nextObjective, newOrderIndex)
		}
		if syncLesson {
			if err := h.repo.UpdateRegionLesson(c.Request.Context(), updated.CourseDraftID, userID, *updated.CourseDraftLessonID, nextTitle, nextObjective, nextSummary, nextOrder, lessonSpec); err != nil {
				if errors.Is(err, ErrNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "linked lesson not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync lesson detail"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"region": updated})
}

func (h *Handler) DeleteRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	regionID, err := uuid.Parse(c.Param("regionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), regionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	if !h.verifyCourseOwner(c, region.CourseDraftID, userID) {
		return
	}

	preLearning, ok := h.isPreLearningCourse(c, region.CourseDraftID)
	if !ok {
		return
	}

	if preLearning {
		err = h.repo.DeleteRegionCascade(c.Request.Context(), regionID)
		if err == nil && region.CourseDraftLessonID != nil {
			err = h.repo.DeleteRegionDraftLesson(c.Request.Context(), region.CourseDraftID, userID, *region.CourseDraftLessonID)
			if errors.Is(err, errRegionLessonHasPoints) {
				c.JSON(http.StatusConflict, gin.H{"error": "region lesson has linked points"})
				return
			}
			if errors.Is(err, ErrNotFound) {
				err = nil
			}
		}
	} else {
		err = h.repo.UpdateRegionStatusWithChildren(c.Request.Context(), regionID, ItemStatusInactive)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete region"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) UpdateRegionStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	regionID, err := uuid.Parse(c.Param("regionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), regionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	if !h.verifyCourseOwner(c, region.CourseDraftID, userID) {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.WithChildren {
		err = h.repo.UpdateRegionStatusWithChildren(c.Request.Context(), regionID, req.Status)
	} else {
		err = h.repo.UpdateRegionStatus(c.Request.Context(), regionID, req.Status)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update region status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
