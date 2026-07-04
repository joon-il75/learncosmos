package curriculum

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/objectstorage"
)

func (h *Handler) ListLearningPointAttachments(c *gin.Context) {
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

	attachments, err := h.repo.ListLearningPointAttachments(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list learning point attachments"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"attachments": attachments})
}

func (h *Handler) CreateLearningPointAttachment(c *gin.Context) {
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

	var req CreateLearningPointAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := createLearningPointAttachmentModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_attachment", &pointID, text, language, map[string]any{
			"field":     "attachment_metadata",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, attachment, err := h.repo.CreateLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point attachment content is invalid"})
		case errors.Is(err, errLearningPointAttachmentLimitExceeded):
			c.JSON(http.StatusConflict, gin.H{"error": "research material attachment limit exceeded"})
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point attachment"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"attachment": attachment}, "point attachment created but failed to reload")
}

func (h *Handler) UploadLearningPointAttachment(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if h.objectStorage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object storage is not configured"})
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

	maxBytes := h.objectStorage.Config().MaxUploadBytes
	requestMaxBytes := maxBytes
	if researchMaterialVideoMaxBytes > requestMaxBytes {
		requestMaxBytes = researchMaterialVideoMaxBytes
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestMaxBytes)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	rawAttachmentType := strings.TrimSpace(c.PostForm("attachment_type"))
	sourceContext := strings.TrimSpace(c.PostForm("source_context"))
	if sourceContext == "" {
		sourceContext = "work_attachment"
	}
	var replaceAttachmentID *uuid.UUID
	if rawReplaceAttachmentID := strings.TrimSpace(c.PostForm("replace_attachment_id")); rawReplaceAttachmentID != "" {
		parsedReplaceAttachmentID, err := uuid.Parse(rawReplaceAttachmentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid replace attachment id"})
			return
		}
		if sourceContext != "work_attachment" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "attachment replacement is only supported for work attachments"})
			return
		}
		replaceAttachmentID = &parsedReplaceAttachmentID
	}
	var artifactID *uuid.UUID
	if sourceContext == "research_material" {
		confirmed, err := h.repo.LearningPointResearchMaterialConfirmed(c.Request.Context(), userID, planetID, pointID)
		if err != nil {
			switch {
			case errors.Is(err, errLearningPointNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check research material status"})
			}
			return
		}
		if confirmed {
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
			return
		}
	}
	if sourceContext == "artifact" {
		parsedArtifactID, err := uuid.Parse(strings.TrimSpace(c.PostForm("artifact_id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artifact id"})
			return
		}
		artifactID = &parsedArtifactID
	}
	isTiptapInlineImageUpload := sourceContext == "research_material" && rawAttachmentType == "image"
	uploadLimit := maxBytes
	if sourceContext == "research_material" || sourceContext == "artifact" {
		uploadLimit = researchMaterialUploadLimitBytes(rawAttachmentType, fileHeader.Filename)
	}
	if fileHeader.Size > uploadLimit {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file is too large"})
		return
	}
	if sourceContext == "research_material" || sourceContext == "artifact" {
		if err := validateResearchMaterialUploadMetadata(fileHeader, rawAttachmentType); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		plannedAttachmentType := inferResearchMaterialAttachmentType(rawAttachmentType, fileHeader.Filename)
		count := 0
		var err error
		if sourceContext == "artifact" {
			if artifactID == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artifact id"})
				return
			}
			if plannedAttachmentType == "video" || plannedAttachmentType == "subtitle" || plannedAttachmentType == "thumbnail" {
				count, err = h.repo.CountLearningPointArtifactAttachmentsByType(c.Request.Context(), userID, planetID, pointID, *artifactID, plannedAttachmentType)
			} else {
				count, err = h.repo.CountLearningPointArtifactFileAttachments(c.Request.Context(), userID, planetID, pointID, *artifactID)
			}
		} else if plannedAttachmentType == "video" || plannedAttachmentType == "subtitle" || plannedAttachmentType == "thumbnail" {
			count, err = h.repo.CountLearningPointAttachmentsByType(c.Request.Context(), userID, planetID, pointID, sourceContext, plannedAttachmentType)
		} else {
			count, err = h.repo.CountLearningPointMaterialFileAttachments(c.Request.Context(), userID, planetID, pointID, sourceContext)
		}
		if err != nil {
			switch {
			case errors.Is(err, errLearningPointNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count learning point attachments"})
			}
			return
		}
		if plannedAttachmentType == "video" && count >= 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "research material video limit exceeded"})
			return
		}
		if plannedAttachmentType == "subtitle" && count >= 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "research material subtitle limit exceeded"})
			return
		}
		if plannedAttachmentType == "thumbnail" && count >= 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "research material thumbnail limit exceeded"})
			return
		}
		if plannedAttachmentType != "video" && plannedAttachmentType != "subtitle" && plannedAttachmentType != "thumbnail" && count >= maxLearningResearchMaterialAttachmentCount {
			c.JSON(http.StatusConflict, gin.H{"error": "research material attachment limit exceeded"})
			return
		}
	}
	if isTiptapInlineImageUpload {
		if err := validateTiptapInlineImageUploadMetadata(fileHeader); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	contentType, err := detectUploadedAttachmentContentType(fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	if isTiptapInlineImageUpload {
		if err := validateTiptapInlineImageUploadContent(contentType); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if sourceContext == "research_material" || sourceContext == "artifact" {
		if err := validateResearchMaterialUploadContent(fileHeader.Filename, rawAttachmentType, contentType); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(fileHeader.Filename), filepath.Ext(fileHeader.Filename))
	}
	if text, fields := uploadLearningPointAttachmentModerationText(title, fileHeader.Filename, rawAttachmentType, sourceContext); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		targetID := &pointID
		if replaceAttachmentID != nil {
			targetID = replaceAttachmentID
		}
		metadata := map[string]any{
			"field":          "attachment_metadata",
			"fields":         fields,
			"planet_id":      planetID.String(),
			"point_id":       pointID.String(),
			"source_context": sourceContext,
		}
		if artifactID != nil {
			metadata["artifact_id"] = artifactID.String()
		}
		if replaceAttachmentID != nil {
			metadata["attachment_id"] = replaceAttachmentID.String()
		}
		if h.enforceSafety(c, userID, "point_attachment", targetID, text, language, metadata) {
			return
		}
	}

	courseID, err := h.repo.ResolveLearningPointCourseID(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve learning point"})
		}
		return
	}

	objectKey := objectstorage.BuildLearningPointObjectKey(userID.String(), courseID.String(), pointID.String(), fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	defer file.Close()
	if err := h.objectStorage.PutObjectReader(c.Request.Context(), objectKey, file, fileHeader.Size, contentType); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to upload attachment"})
		return
	}

	attachmentType := normalizeUploadedAttachmentType(rawAttachmentType, contentType)
	fileSize := fileHeader.Size
	req := CreateLearningPointAttachmentRequest{
		Provider:       stringPtr("learner"),
		AttachmentType: &attachmentType,
		SourceContext:  &sourceContext,
		ArtifactID:     artifactID,
		Title:          &title,
		FilePath:       &objectKey,
		FileSize:       &fileSize,
		MimeType:       &contentType,
	}

	if replaceAttachmentID != nil {
		_, previousAttachment, err := h.repo.GetLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, *replaceAttachmentID)
		if err != nil {
			switch {
			case errors.Is(err, errLearningPointMaterialLocked):
				c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
			case errors.Is(err, errLearningPointNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learning point attachment"})
			}
			return
		}
		if previousAttachment.SourceContext != "work_attachment" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "attachment replacement is only supported for work attachments"})
			return
		}
		updateReq := UpdateLearningPointAttachmentRequest{
			Provider:       stringPtr("learner"),
			AttachmentType: &attachmentType,
			SourceContext:  &sourceContext,
			Title:          &title,
			FilePath:       &objectKey,
			FileSize:       &fileSize,
			MimeType:       &contentType,
		}
		courseID, attachment, err := h.repo.UpdateLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, *replaceAttachmentID, updateReq)
		if err != nil {
			switch {
			case errors.Is(err, errLearningPointMaterialLocked):
				c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
			case errors.Is(err, errLearningPointNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point attachment"})
			}
			return
		}
		extra := gin.H{"attachment": attachment, "replaced_attachment": previousAttachment}
		if previousObjectKey := strings.TrimSpace(previousAttachment.FilePath); previousObjectKey != "" && previousObjectKey != objectKey {
			if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), previousObjectKey) {
				extra["previous_attachment_object_delete_skipped"] = true
				extra["previous_attachment_object_delete_reason"] = "object_key_not_managed"
				log.Printf("[point-attachment] previous object delete skipped attachment=%s point=%s key=%q reason=object_key_not_managed", *replaceAttachmentID, pointID, previousObjectKey)
			} else if h.objectStorage == nil {
				extra["previous_attachment_object_delete_failed"] = true
				extra["previous_attachment_object_delete_reason"] = "object_storage_not_configured"
				log.Printf("[point-attachment] previous object delete skipped attachment=%s point=%s key=%q reason=not_configured", *replaceAttachmentID, pointID, previousObjectKey)
			} else if err := h.objectStorage.DeleteObject(c.Request.Context(), previousObjectKey); err != nil {
				extra["previous_attachment_object_delete_failed"] = true
				extra["previous_attachment_object_delete_reason"] = "object_storage_delete_failed"
				log.Printf("[point-attachment] previous object delete failed attachment=%s point=%s key=%q err=%v", *replaceAttachmentID, pointID, previousObjectKey, err)
			} else {
				extra["previous_attachment_object_deleted"] = true
			}
		}
		h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, extra, "point attachment replaced but failed to reload")
		return
	}

	courseID, attachment, err := h.repo.CreateLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point attachment content is invalid"})
		case errors.Is(err, errLearningPointAttachmentLimitExceeded):
			c.JSON(http.StatusConflict, gin.H{"error": "research material attachment limit exceeded"})
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point attachment"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"attachment": attachment}, "point attachment uploaded but failed to reload")
}

func (h *Handler) UpdateLearningPointAttachment(c *gin.Context) {
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
	attachmentID, err := uuid.Parse(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	var req UpdateLearningPointAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := updateLearningPointAttachmentModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_attachment", &attachmentID, text, language, map[string]any{
			"field":         "attachment_metadata",
			"fields":        fields,
			"planet_id":     planetID.String(),
			"point_id":      pointID.String(),
			"attachment_id": attachmentID.String(),
		}) {
			return
		}
	}

	courseID, attachment, err := h.repo.UpdateLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, attachmentID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point attachment content is invalid"})
		case errors.Is(err, errLearningPointAttachmentLimitExceeded):
			c.JSON(http.StatusConflict, gin.H{"error": "research material attachment limit exceeded"})
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point attachment"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"attachment": attachment}, "point attachment updated but failed to reload")
}

func createLearningPointAttachmentModerationText(req CreateLearningPointAttachmentRequest) (string, []string) {
	parts := make([]string, 0, 7)
	fields := make([]string, 0, 7)
	appendOptionalModerationPart(&parts, &fields, req.Provider, "provider")
	appendOptionalModerationPart(&parts, &fields, req.AttachmentType, "attachment_type")
	appendOptionalModerationPart(&parts, &fields, req.SourceContext, "source_context")
	appendOptionalModerationPart(&parts, &fields, req.Title, "title")
	appendOptionalModerationPart(&parts, &fields, req.URL, "url")
	appendOptionalModerationPart(&parts, &fields, req.FilePath, "file_path")
	appendOptionalModerationPart(&parts, &fields, req.MimeType, "mime_type")
	return strings.Join(parts, "\n"), fields
}

func updateLearningPointAttachmentModerationText(req UpdateLearningPointAttachmentRequest) (string, []string) {
	parts := make([]string, 0, 7)
	fields := make([]string, 0, 7)
	appendOptionalModerationPart(&parts, &fields, req.Provider, "provider")
	appendOptionalModerationPart(&parts, &fields, req.AttachmentType, "attachment_type")
	appendOptionalModerationPart(&parts, &fields, req.SourceContext, "source_context")
	appendOptionalModerationPart(&parts, &fields, req.Title, "title")
	appendOptionalModerationPart(&parts, &fields, req.URL, "url")
	appendOptionalModerationPart(&parts, &fields, req.FilePath, "file_path")
	appendOptionalModerationPart(&parts, &fields, req.MimeType, "mime_type")
	return strings.Join(parts, "\n"), fields
}

func uploadLearningPointAttachmentModerationText(title, filename, attachmentType, sourceContext string) (string, []string) {
	parts := make([]string, 0, 4)
	fields := make([]string, 0, 4)
	appendModerationPart(&parts, &fields, title, "title")
	appendModerationPart(&parts, &fields, filename, "filename")
	appendModerationPart(&parts, &fields, attachmentType, "attachment_type")
	appendModerationPart(&parts, &fields, sourceContext, "source_context")
	return strings.Join(parts, "\n"), fields
}

func appendModerationPart(parts *[]string, fields *[]string, value string, field string) {
	text := strings.TrimSpace(value)
	if text == "" {
		return
	}
	*parts = append(*parts, text)
	*fields = append(*fields, field)
}

func (h *Handler) OpenLearningPointAttachment(c *gin.Context) {
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
	attachmentID, err := uuid.Parse(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	courseID, attachment, err := h.repo.GetReadablePointAttachment(c.Request.Context(), userID, planetID, pointID, attachmentID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learning point attachment"})
		}
		return
	}
	if strings.TrimSpace(attachment.FilePath) == "" {
		if strings.TrimSpace(attachment.URL) == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "attachment file not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"url": attachment.URL, "external": true})
		return
	}
	if h.objectStorage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object storage is not configured"})
		return
	}
	if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), attachment.FilePath) {
		log.Printf("[point-attachment] object open denied attachment=%s point=%s key=%q reason=object_key_not_managed", attachmentID, pointID, attachment.FilePath)
		c.JSON(http.StatusForbidden, gin.H{"error": "attachment file is not accessible"})
		return
	}

	presignedURL, expiresAt, err := h.objectStorage.PresignGetObject(attachment.FilePath, 0)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create attachment url"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": presignedURL, "expires_at": expiresAt})
}

func (h *Handler) InlineLearningPointAttachment(c *gin.Context) {
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
	attachmentID, err := uuid.Parse(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	courseID, attachment, err := h.repo.GetReadablePointAttachment(c.Request.Context(), userID, planetID, pointID, attachmentID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load learning point attachment"})
		}
		return
	}
	attachmentType := strings.ToLower(strings.TrimSpace(attachment.AttachmentType))
	mimeType := strings.ToLower(strings.TrimSpace(attachment.MimeType))
	if attachmentType == "subtitle" {
		if strings.TrimSpace(attachment.FilePath) == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "attachment subtitle file not found"})
			return
		}
		if h.objectStorage == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object storage is not configured"})
			return
		}
		if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), attachment.FilePath) {
			log.Printf("[point-attachment] subtitle object denied attachment=%s point=%s key=%q reason=object_key_not_managed", attachmentID, pointID, attachment.FilePath)
			c.JSON(http.StatusForbidden, gin.H{"error": "attachment file is not accessible"})
			return
		}
		body, contentType, err := h.objectStorage.GetObjectBytes(c.Request.Context(), attachment.FilePath, 2*1024*1024)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to load attachment subtitle"})
			return
		}
		if strings.TrimSpace(contentType) == "" {
			contentType = "text/vtt; charset=utf-8"
		}
		c.Header("Cache-Control", "private, no-store")
		c.Data(http.StatusOK, contentType, body)
		return
	}
	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "attachment is not an image"})
		return
	}
	if strings.TrimSpace(attachment.FilePath) == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "attachment image file not found"})
		return
	}
	if h.objectStorage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object storage is not configured"})
		return
	}
	if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), attachment.FilePath) {
		log.Printf("[point-attachment] image object denied attachment=%s point=%s key=%q reason=object_key_not_managed", attachmentID, pointID, attachment.FilePath)
		c.JSON(http.StatusForbidden, gin.H{"error": "attachment file is not accessible"})
		return
	}

	presignedURL, _, err := h.objectStorage.PresignGetObject(attachment.FilePath, 0)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create attachment image url"})
		return
	}
	c.Header("Cache-Control", "private, no-store")
	c.Redirect(http.StatusTemporaryRedirect, presignedURL)
}

func (h *Handler) DeleteLearningPointAttachment(c *gin.Context) {
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
	attachmentID, err := uuid.Parse(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	courseID, attachment, err := h.repo.DeleteLearningPointAttachment(c.Request.Context(), userID, planetID, pointID, attachmentID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "confirmed research material cannot be deleted"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point attachment not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete learning point attachment"})
		}
		return
	}

	extra := gin.H{"deleted_attachment": attachment}
	if objectKey := strings.TrimSpace(attachment.FilePath); objectKey != "" {
		if !objectstorage.IsLearningPointObjectKeyFor(userID.String(), courseID.String(), pointID.String(), objectKey) {
			extra["attachment_object_delete_skipped"] = true
			extra["attachment_object_delete_reason"] = "object_key_not_managed"
			log.Printf("[point-attachment] object delete skipped attachment=%s point=%s key=%q reason=object_key_not_managed", attachmentID, pointID, objectKey)
		} else if h.objectStorage == nil {
			extra["attachment_object_delete_failed"] = true
			extra["attachment_object_delete_reason"] = "object_storage_not_configured"
			log.Printf("[point-attachment] object delete skipped attachment=%s point=%s key=%q reason=not_configured", attachmentID, pointID, objectKey)
		} else if err := h.objectStorage.DeleteObject(c.Request.Context(), objectKey); err != nil {
			extra["attachment_object_delete_failed"] = true
			extra["attachment_object_delete_reason"] = "object_storage_delete_failed"
			log.Printf("[point-attachment] object delete failed attachment=%s point=%s key=%q err=%v", attachmentID, pointID, objectKey, err)
		} else {
			extra["attachment_object_deleted"] = true
		}
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, extra, "point attachment deleted but failed to reload")
}
