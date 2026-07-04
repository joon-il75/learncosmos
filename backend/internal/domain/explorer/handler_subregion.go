package explorer

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
)

func (h *Handler) CreateSubRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateSubRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), req.RegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	if !h.verifyCourseOwner(c, region.CourseDraftID, userID) {
		return
	}

	count, err := h.repo.CountSubRegionsByRegion(c.Request.Context(), req.RegionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subregion count"})
		return
	}

	if err := h.service.ValidateSubRegionLimit(count); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mainLessonID, err := h.repo.GetMainLessonIDByRegion(
		c.Request.Context(),
		region.ID,
		region.CourseDraftID,
		region.OrderIndex,
		region.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "main lesson not found for region"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve parent main lesson"})
		return
	}

	newLesson, err := h.curriculumRepo.CreateSubLesson(c.Request.Context(), userID, region.CourseDraftID, mainLessonID, curriculum.CreateSubLessonRequest{
		Title:      req.Name,
		Objective:  req.Description,
		OrderIndex: req.OrderIndex,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subregion lesson"})
		return
	}

	subRegionLessonID := newLesson.ID
	lessonSpec := h.buildSubRegionLessonSpec(c.Request.Context(), userID, subRegionLessonID, req.Name, req.Description, newLesson.OrderIndex)

	lessonTitle := req.Name
	if lessonTitle == "" && newLesson.Title != "" {
		lessonTitle = newLesson.Title
	}
	lessonObjective := req.Description
	if lessonObjective != nil {
		trimmed := strings.TrimSpace(*lessonObjective)
		lessonObjective = &trimmed
	}
	lessonSummary := lessonObjective
	if err := h.repo.UpdateSubRegionLesson(c.Request.Context(), region.CourseDraftID, userID, subRegionLessonID, &lessonTitle, lessonObjective, lessonSummary, &newLesson.OrderIndex, &mainLessonID, lessonSpec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync subregion lesson"})
		if err := h.repo.DeleteSubRegionDraftLesson(c.Request.Context(), region.CourseDraftID, userID, subRegionLessonID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rollback created subregion lesson"})
			return
		}
		return
	}

	subregion, err := h.repo.InsertSubRegion(c.Request.Context(), req, &subRegionLessonID)
	if err != nil {
		err = h.repo.DeleteSubRegionDraftLesson(c.Request.Context(), region.CourseDraftID, userID, subRegionLessonID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subregion and rollback lesson"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subregion"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"subregion": subregion})
}

func (h *Handler) UpdateSubRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subRegionID, err := uuid.Parse(c.Param("subRegionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subregion id"})
		return
	}

	subregion, err := h.repo.GetSubRegionByID(c.Request.Context(), subRegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subregion not found"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), subregion.RegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "parent region not found"})
		return
	}

	if !h.verifyCourseOwner(c, region.CourseDraftID, userID) {
		return
	}

	var req UpdateSubRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRegionID := subregion.RegionID
	newParentLessonID := (*uuid.UUID)(nil)
	if req.RegionID != nil && *req.RegionID != subregion.RegionID {
		newRegion, err := h.repo.GetRegionByID(c.Request.Context(), *req.RegionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "target region not found"})
			return
		}
		if newRegion.CourseDraftID != region.CourseDraftID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target region belongs to a different course"})
			return
		}
		count, err := h.repo.CountSubRegionsByRegion(c.Request.Context(), *req.RegionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subregion count"})
			return
		}
		if err := h.service.ValidateSubRegionLimit(count); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mainLessonID, err := h.repo.GetMainLessonIDByRegion(c.Request.Context(), *req.RegionID, region.CourseDraftID, newRegion.OrderIndex, newRegion.Name)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "target main lesson not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve target main lesson"})
			return
		}
		newRegionID = *req.RegionID
		req.RegionID = &newRegionID
		newParentLessonID = &mainLessonID
	}

	updated, err := h.repo.UpdateSubRegion(c.Request.Context(), subRegionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subregion"})
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

		syncLesson := req.RegionID != nil && *req.RegionID != subregion.RegionID
		if req.Name != nil {
			syncLesson = true
		}
		if req.Description != nil {
			syncLesson = true
		}
		if req.OrderIndex != nil {
			syncLesson = true
		}

		var lessonSpec *string
		if syncLesson {
			lessonSpec = h.buildSubRegionLessonSpec(c.Request.Context(), userID, *updated.CourseDraftLessonID, nextTitleText, nextObjective, newOrderIndex)
		}

		if err := h.repo.UpdateSubRegionLesson(c.Request.Context(), region.CourseDraftID, userID, *updated.CourseDraftLessonID, nextTitle, nextObjective, nextSummary, nextOrder, newParentLessonID, lessonSpec); err != nil {
			if errors.Is(err, ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "linked lesson not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync lesson detail"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"subregion": updated})
}

func (h *Handler) DeleteSubRegion(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subRegionID, err := uuid.Parse(c.Param("subRegionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subregion id"})
		return
	}

	subregion, err := h.repo.GetSubRegionByID(c.Request.Context(), subRegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subregion not found"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), subregion.RegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "parent region not found"})
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
		err = h.repo.DeleteSubRegionCascade(c.Request.Context(), subRegionID)
		if err == nil && subregion.CourseDraftLessonID != nil {
			err = h.repo.DeleteSubRegionDraftLesson(c.Request.Context(), region.CourseDraftID, userID, *subregion.CourseDraftLessonID)
			if errors.Is(err, errSubRegionLessonHasPoints) {
				c.JSON(http.StatusConflict, gin.H{"error": "subregion lesson has linked points"})
				return
			}
			if errors.Is(err, ErrNotFound) {
				err = nil
			}
		}
	} else {
		err = h.repo.UpdateSubRegionStatus(c.Request.Context(), subRegionID, ItemStatusInactive)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subregion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) UpdateSubRegionStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subRegionID, err := uuid.Parse(c.Param("subRegionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subregion id"})
		return
	}

	subregion, err := h.repo.GetSubRegionByID(c.Request.Context(), subRegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subregion not found"})
		return
	}

	region, err := h.repo.GetRegionByID(c.Request.Context(), subregion.RegionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "parent region not found"})
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

	if req.Status == ItemStatusActive && region.Status == ItemStatusInactive {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrParentInactive.Error()})
		return
	}

	if err := h.repo.UpdateSubRegionStatus(c.Request.Context(), subRegionID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subregion status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) buildSubRegionLessonSpec(ctx context.Context, userID, lessonID uuid.UUID, title string, objective *string, orderIndex int) *string {
	if lessonID == uuid.Nil {
		return nil
	}
	searchContext, err := h.curriculumRepo.GetLessonRecommendationSearchContextByLessonID(ctx, userID, lessonID)
	if err != nil {
		return nil
	}
	searchContext.LessonTitle = strings.TrimSpace(title)
	searchContext.LessonObjective = objective
	if orderIndex > 0 {
		searchContext.OrderIndex = orderIndex
	}
	fallback := searchContext.FallbackSpec()
	specJSON, err := curriculum.MarshalLessonRecommendationSearchSpec(fallback)
	if err != nil || specJSON == "{}" {
		return nil
	}
	return &specJSON
}
