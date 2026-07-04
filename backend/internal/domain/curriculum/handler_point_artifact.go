package curriculum

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/objectstorage"
)

func (h *Handler) UpdateLearningPointArtifact(c *gin.Context) {
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

	var req UpdateDraftLessonArtifactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := updateDraftLessonArtifactModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_artifact", &pointID, text, language, map[string]any{
			"field":     "artifact_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, artifact, err := h.repo.UpdateLearningPointArtifact(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point artifact content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point artifact"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"artifact": artifact}, "point artifact saved but failed to reload")
}

func (h *Handler) ListLearningPointArtifacts(c *gin.Context) {
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

	artifacts, err := h.repo.ListLearningPointArtifacts(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointBlockInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point block content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list learning point artifacts"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"artifacts": artifacts})
}

func (h *Handler) CreateLearningPointArtifact(c *gin.Context) {
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

	var req CreateLearningPointArtifactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := createLearningPointArtifactModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_artifact", &pointID, text, language, map[string]any{
			"field":     "artifact_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, artifact, err := h.repo.CreateLearningPointArtifact(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point artifact content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point artifact"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"artifact": artifact}, "point artifact created but failed to reload")
}

func (h *Handler) UpdateLearningPointArtifactItem(c *gin.Context) {
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
	artifactID, err := uuid.Parse(c.Param("artifact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artifact id"})
		return
	}

	var req UpdateLearningPointArtifactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := updateLearningPointArtifactModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_artifact", &artifactID, text, language, map[string]any{
			"field":       "artifact_content",
			"fields":      fields,
			"planet_id":   planetID.String(),
			"point_id":    pointID.String(),
			"artifact_id": artifactID.String(),
		}) {
			return
		}
	}

	courseID, artifact, err := h.repo.UpdateLearningPointArtifactItem(c.Request.Context(), userID, planetID, pointID, artifactID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point artifact content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point artifact not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point artifact"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"artifact": artifact}, "point artifact updated but failed to reload")
}

func (h *Handler) DeleteLearningPointArtifact(c *gin.Context) {
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
	artifactID, err := uuid.Parse(c.Param("artifact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artifact id"})
		return
	}

	courseID, deletedAttachments, err := h.repo.DeleteLearningPointArtifact(c.Request.Context(), userID, planetID, pointID, artifactID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point artifact not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete learning point artifact"})
		}
		return
	}

	extra := gin.H{"deleted_artifact_attachments": deletedAttachments}
	if h.objectStorage != nil {
		deletedCount := 0
		skippedCount := 0
		for _, attachment := range deletedAttachments {
			if objectKey := strings.TrimSpace(attachment.FilePath); objectKey != "" {
				if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), objectKey) {
					skippedCount++
					log.Printf("[point-artifact] attachment object delete skipped artifact=%s attachment=%s point=%s key=%q reason=object_key_not_managed", artifactID, attachment.ID, pointID, objectKey)
				} else if err := h.objectStorage.DeleteObject(c.Request.Context(), objectKey); err != nil {
					log.Printf("[point-artifact] attachment object delete failed artifact=%s attachment=%s point=%s key=%q err=%v", artifactID, attachment.ID, pointID, objectKey, err)
					extra["artifact_attachment_object_delete_failed"] = true
				} else {
					deletedCount++
				}
			}
		}
		extra["artifact_attachment_objects_deleted"] = deletedCount
		if skippedCount > 0 {
			extra["artifact_attachment_objects_delete_skipped"] = skippedCount
			extra["artifact_attachment_object_delete_reason"] = "object_key_not_managed"
		}
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, extra, "point artifact deleted but failed to reload")
}

func updateDraftLessonArtifactModerationText(req UpdateDraftLessonArtifactRequest) (string, []string) {
	parts := make([]string, 0, 4)
	fields := make([]string, 0, 4)
	appendOptionalModerationPart(&parts, &fields, req.ArtifactType, "artifact_type")
	appendOptionalModerationPart(&parts, &fields, req.Title, "title")
	appendOptionalModerationPart(&parts, &fields, req.URL, "url")
	appendOptionalModerationPart(&parts, &fields, req.Description, "description")
	return strings.Join(parts, "\n"), fields
}

func createLearningPointArtifactModerationText(req CreateLearningPointArtifactRequest) (string, []string) {
	parts := make([]string, 0, 9)
	fields := make([]string, 0, 9)
	appendOptionalModerationPart(&parts, &fields, req.ArtifactType, "artifact_type")
	appendOptionalModerationPart(&parts, &fields, req.Title, "title")
	appendOptionalModerationPart(&parts, &fields, req.URL, "url")
	appendOptionalModerationPart(&parts, &fields, req.Description, "description")
	appendOptionalModerationPart(&parts, &fields, req.PointCategory, "point_category")
	appendOptionalModerationPart(&parts, &fields, req.ProductionProcess, "production_process")
	appendOptionalModerationPart(&parts, &fields, req.LearnedPoints, "learned_points")
	appendOptionalModerationPart(&parts, &fields, req.DifficultPoints, "difficult_points")
	appendOptionalModerationPart(&parts, &fields, req.Visibility, "visibility")
	return strings.Join(parts, "\n"), fields
}

func updateLearningPointArtifactModerationText(req UpdateLearningPointArtifactRequest) (string, []string) {
	parts := make([]string, 0, 9)
	fields := make([]string, 0, 9)
	appendOptionalModerationPart(&parts, &fields, req.ArtifactType, "artifact_type")
	appendOptionalModerationPart(&parts, &fields, req.Title, "title")
	appendOptionalModerationPart(&parts, &fields, req.URL, "url")
	appendOptionalModerationPart(&parts, &fields, req.Description, "description")
	appendOptionalModerationPart(&parts, &fields, req.PointCategory, "point_category")
	appendOptionalModerationPart(&parts, &fields, req.ProductionProcess, "production_process")
	appendOptionalModerationPart(&parts, &fields, req.LearnedPoints, "learned_points")
	appendOptionalModerationPart(&parts, &fields, req.DifficultPoints, "difficult_points")
	appendOptionalModerationPart(&parts, &fields, req.Visibility, "visibility")
	return strings.Join(parts, "\n"), fields
}

func appendOptionalModerationPart(parts *[]string, fields *[]string, value *string, field string) {
	if value == nil {
		return
	}
	text := strings.TrimSpace(*value)
	if text == "" {
		return
	}
	*parts = append(*parts, text)
	*fields = append(*fields, field)
}
