package curriculum

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListDrafts(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	drafts, err := h.repo.ListDrafts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list drafts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"drafts": drafts})
}

func (h *Handler) GetDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": draft})
}

func (h *Handler) UpdateDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCourseDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	updated, err := h.service.ApplyDraftUpdate(*existing, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.resolveSelectedExternalResources(c.Request.Context(), userID, updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save selected external resources"})
		return
	}

	if err := h.repo.UpdateDraft(c.Request.Context(), updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": updated})
}

func (h *Handler) DeleteDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req DeleteCourseDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	if draft.Draft.Status == DraftStatusArchived {
		if draft.Draft.IsInactive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already deactivated"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "completed_course_cannot_be_deactivated"})
		return
	}
	if strings.TrimSpace(req.ConfirmTitle) != strings.TrimSpace(draft.Draft.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirm_title_mismatch"})
		return
	}

	if err := h.repo.DeleteDraft(c.Request.Context(), draftID, userID, draft.Draft.SourceQuery, draft.Draft.Title); err != nil {
		if err.Error() == "draft_not_deletable" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "draft_not_deletable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true, "deactivated": true})
}

func (h *Handler) ActivateDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req DeleteCourseDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	if !draft.Draft.IsInactive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "draft_not_inactive"})
		return
	}
	if strings.TrimSpace(req.ConfirmTitle) != strings.TrimSpace(draft.Draft.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirm_title_mismatch"})
		return
	}

	destination, err := h.repo.ActivateDraft(c.Request.Context(), draftID, userID, draft.Draft.SourceQuery, draft.Draft.Title)
	if err != nil {
		switch err.Error() {
		case "draft_not_found":
			c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		case "draft_not_inactive":
			c.JSON(http.StatusBadRequest, gin.H{"error": "draft_not_inactive"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to activate draft"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"activated": true, "destination": destination})
}

func (h *Handler) UpdateDraftLesson(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req UpdateDraftLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateDraftLessonDetail(c.Request.Context(), draftID, userID, lessonID, req); err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftLessonInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "lesson title is required"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lesson"})
		}
		return
	}

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lesson updated but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": fetchedDraft})
}

func (h *Handler) UpdateDraftLevel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	levelID, err := uuid.Parse(c.Param("main_lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid main lesson id"})
		return
	}

	var req UpdateDraftMainLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateDraftMainLessonDetail(c.Request.Context(), draftID, userID, levelID, req); err != nil {
		switch {
		case errors.Is(err, errDraftLevelNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "main lesson not found"})
		case errors.Is(err, errDraftLevelInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "main lesson title is required"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update main lesson"})
		}
		return
	}

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "main lesson updated but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": fetchedDraft})
}

func (h *Handler) UpdateDraftStructure(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateDraftStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateDraftStructure(c.Request.Context(), draftID, userID, req); err != nil {
		switch {
		case errors.Is(err, errDraftStructureNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "draft structure not found"})
		case errors.Is(err, errDraftStructureInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft structure"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update draft structure"})
		}
		return
	}

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "draft structure updated but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": fetchedDraft})
}

func (h *Handler) UpdateDraftLevelMemo(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	levelID, err := uuid.Parse(c.Param("main_lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid main lesson id"})
		return
	}

	var req UpdateDraftDetailMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memo, err := h.repo.UpdateDraftMainLessonMemo(c.Request.Context(), draftID, userID, levelID, req)
	if err != nil {
		switch {
		case errors.Is(err, errDraftLevelNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "main lesson not found"})
		case errors.Is(err, errDraftMemoInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "memo note is required"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update main lesson memo"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"memo": memo})
}

func (h *Handler) UpdateDraftLessonMemo(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req UpdateDraftDetailMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memo, err := h.repo.UpdateDraftLessonMemo(c.Request.Context(), draftID, userID, lessonID, req)
	if err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftMemoInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "memo note is required"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lesson memo"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"memo": memo})
}

func (h *Handler) UpdateDraftLessonJournal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req UpdateDraftLessonJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	journal, err := h.repo.UpdateDraftLessonJournal(c.Request.Context(), draftID, userID, lessonID, req)
	if err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftJournalInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "journal content is required"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lesson journal"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"journal": journal})
}

func (h *Handler) UpdateDraftLessonRecord(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req UpdateDraftLessonRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.repo.UpdateDraftLessonRecord(c.Request.Context(), draftID, userID, lessonID, req)
	if err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftRecordInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "record content is invalid"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lesson record"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func (h *Handler) UpdateDraftLessonArtifact(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req UpdateDraftLessonArtifactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artifact, err := h.repo.UpdateDraftLessonArtifact(c.Request.Context(), draftID, userID, lessonID, req)
	if err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftArtifactInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "artifact content is invalid"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lesson artifact"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"artifact": artifact})
}
