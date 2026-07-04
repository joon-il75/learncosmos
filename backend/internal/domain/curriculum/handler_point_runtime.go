package curriculum

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

func (h *Handler) UpdateLearningPointRuntime(c *gin.Context) {
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

	var req UpdateLearningPointRuntimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID, err := h.repo.UpdateLearningPointRuntime(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errPointRuntimeInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point runtime status is invalid"})
		case errors.Is(err, errLearningPointMaterialNotReady):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is not confirmed"})
		case errors.Is(err, errLearningPointCompletionNotReady):
			c.JSON(http.StatusConflict, gin.H{"error": "point completion requires at least one answered question and six self evaluation inputs"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point runtime"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, nil, "point runtime saved but failed to reload")
}

func (h *Handler) ConfirmLearningPointResearchMaterial(c *gin.Context) {
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

	courseID, err := h.repo.ConfirmLearningPointResearchMaterial(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointMaterialNotReady):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is not ready"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning research point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm research material"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"research_material_confirmed": true}, "research material confirmed but failed to reload")
}

func (h *Handler) UnconfirmLearningPointResearchMaterial(c *gin.Context) {
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

	courseID, err := h.repo.UnconfirmLearningPointResearchMaterial(c.Request.Context(), userID, planetID, pointID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning research point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unconfirm research material"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"research_material_confirmed": false}, "research material unconfirmed but failed to reload")
}

func (h *Handler) UpdateLearningPointJournal(c *gin.Context) {
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

	var req UpdateDraftLessonJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := learningPointJournalModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_journal", &pointID, text, language, map[string]any{
			"field":     "journal_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, journal, err := h.repo.UpdateLearningPointJournal(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointJournalInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point journal content is required"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point journal"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"journal": journal}, "point journal saved but failed to reload")
}

func learningPointJournalModerationText(req UpdateDraftLessonJournalRequest) (string, []string) {
	parts := make([]string, 0, 8)
	fields := make([]string, 0, 8)
	add := func(field string, value *string) {
		if value == nil {
			return
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			return
		}
		parts = append(parts, trimmed)
		fields = append(fields, field)
	}
	add("observation", req.Observation)
	add("reflection", req.Reflection)
	add("next_step", req.NextStep)
	add("core_concept", req.CoreConcept)
	add("my_explanation", req.MyExplanation)
	add("examples", req.Examples)
	add("confused_parts", req.ConfusedParts)
	add("reference_links", req.ReferenceLinks)
	return strings.Join(parts, "\n"), fields
}

func (h *Handler) UpdateLearningPointGoal(c *gin.Context) {
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

	var req UpdateLearningPointGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if text, fields := learningPointGoalModerationText(req); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "point_goal", &pointID, text, language, map[string]any{
			"field":     "goal_content",
			"fields":    fields,
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, err := h.repo.UpdateLearningPointGoal(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errPointRuntimeInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point goal content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point goal"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, nil, "point goal saved but failed to reload")
}

func learningPointGoalModerationText(req UpdateLearningPointGoalRequest) (string, []string) {
	parts := make([]string, 0, 2)
	fields := make([]string, 0, 2)
	add := func(field string, value *string) {
		if value == nil {
			return
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			return
		}
		parts = append(parts, trimmed)
		fields = append(fields, field)
	}
	add("point_goal", req.PointGoal)
	add("point_category", req.PointCategory)
	return strings.Join(parts, "\n"), fields
}

func (h *Handler) UpdateLearningPointRecord(c *gin.Context) {
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

	var req UpdateDraftLessonRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID, record, err := h.repo.UpdateLearningPointRecord(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointRecordInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point record content is invalid"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point record"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"record": record}, "point record saved but failed to reload")
}

func (h *Handler) GenerateLearningPointAISummary(c *gin.Context) {
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

	if h.shouldUsePointAISummaryWorker() {
		h.handlePointAISummaryWithWorker(c, userID, planetID, pointID)
		return
	}

	result, err := h.executePointAISummary(c.Request.Context(), userID, planetID, pointID, "point_ai_summary_sync")
	if err != nil {
		h.writePointAISummaryError(c, err)
		return
	}
	h.respondWithLearningPointDetail(c, userID, result.CourseID, pointID, http.StatusOK, gin.H{}, "ai summary saved but failed to reload")
}

func (h *Handler) shouldUsePointAISummaryWorker() bool {
	return h != nil && h.pointAISummaryWorkerEnabled && h.llmGateway != nil
}

func (h *Handler) handlePointAISummaryWithWorker(c *gin.Context, userID, planetID, pointID uuid.UUID) {
	learningLanguage := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	job, err := h.llmGateway.SubmitAndWait(c.Request.Context(), CreatePointAISummaryJobInput(userID, planetID, pointID, learningLanguage), h.pointAISummaryWorkerSyncWait)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
		return
	}
	if job == nil || job.Status != llmjobs.StatusSucceeded {
		if job != nil && !terminalLLMJobStatus(job.Status) {
			c.JSON(http.StatusAccepted, gin.H{
				"error":      "point_ai_summary_pending",
				"error_code": "point_ai_summary_pending",
				"job_id":     job.ID,
				"feature":    job.Feature,
				"status":     job.Status,
				"poll_url":   "/api/v1/llm-jobs/" + job.ID.String(),
			})
			return
		}
		code := "point_ai_summary_failed"
		if job != nil && job.ErrorCode != nil && *job.ErrorCode != "" {
			code = *job.ErrorCode
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": code, "error_code": code})
		return
	}
	var result pointAISummaryJobResult
	if err := json.Unmarshal(job.ResultRef, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_ai_summary_result_invalid", "error_code": "point_ai_summary_result_invalid"})
		return
	}
	courseID, err := uuid.Parse(result.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "point_ai_summary_result_invalid", "error_code": "point_ai_summary_result_invalid"})
		return
	}
	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{}, "ai summary saved but failed to reload")
}

func terminalLLMJobStatus(status llmjobs.Status) bool {
	return status == llmjobs.StatusSucceeded || status == llmjobs.StatusFailed || status == llmjobs.StatusCanceled || status == llmjobs.StatusExpired
}

func (h *Handler) writePointAISummaryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errNoPointAILLMAvailable):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no_llm_available"})
	case errors.Is(err, errPointAISummaryFailed):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ai_summary_failed"})
	case errors.Is(err, errLearningPointAISummaryInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": "point ai summary is empty"})
	case errors.Is(err, errLearningPointNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save learning point ai summary"})
	}
}

func (h *Handler) CreateLearningPointBlock(c *gin.Context) {
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

	var req CreateResearchNodeBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if text, fields := learningPointResearchBlockModerationText(req.Content); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "research_material", &pointID, text, language, map[string]any{
			"field":      "research_block_content",
			"fields":     fields,
			"block_type": req.BlockType,
			"planet_id":  planetID.String(),
			"point_id":   pointID.String(),
		}) {
			return
		}
	}

	courseID, block, err := h.repo.CreateLearningPointBlock(c.Request.Context(), userID, planetID, pointID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointBlockInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point block content is invalid"})
		case errors.Is(err, errLearningPointBlockLimitExceeded):
			c.JSON(http.StatusConflict, gin.H{"error": "research material post limit exceeded"})
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create learning point block"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusCreated, gin.H{"block": block}, "point block created but failed to reload")
}

func (h *Handler) UpdateLearningPointBlock(c *gin.Context) {
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

	if text, fields := learningPointResearchBlockModerationText(req.Content); text != "" {
		language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
		if h.enforceSafety(c, userID, "research_material", &pointID, text, language, map[string]any{
			"field":     "research_block_content",
			"fields":    fields,
			"block_id":  blockID.String(),
			"planet_id": planetID.String(),
			"point_id":  pointID.String(),
		}) {
			return
		}
	}

	courseID, block, err := h.repo.UpdateLearningPointBlock(c.Request.Context(), userID, planetID, pointID, blockID, req)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointBlockInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "point block content is invalid"})
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "research material is locked"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update learning point block"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"block": block}, "point block updated but failed to reload")
}

func learningPointResearchBlockModerationText(raw json.RawMessage) (string, []string) {
	if len(raw) == 0 {
		return "", nil
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", nil
	}
	fields := []string{}
	parts := []string{}
	for _, key := range []string{"title", "text", "caption", "url"} {
		value, ok := payload[key].(string)
		if !ok {
			continue
		}
		text := strings.TrimSpace(value)
		if key == "text" {
			text = researchBlockPlainText(text)
		}
		if text == "" {
			continue
		}
		fields = append(fields, key)
		parts = append(parts, text)
	}
	return strings.Join(parts, "\n"), fields
}

func (h *Handler) DeleteLearningPointBlock(c *gin.Context) {
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
	blockID, err := uuid.Parse(c.Param("block_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid block id"})
		return
	}

	courseID, err := h.repo.DeleteLearningPointBlock(c.Request.Context(), userID, planetID, pointID, blockID)
	if err != nil {
		switch {
		case errors.Is(err, errLearningPointMaterialLocked):
			c.JSON(http.StatusConflict, gin.H{"error": "confirmed research material cannot be deleted"})
		case errors.Is(err, errLearningPointNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "learning point not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete learning point block"})
		}
		return
	}

	h.respondWithLearningPointDetail(c, userID, courseID, pointID, http.StatusOK, gin.H{"deleted": true}, "point block deleted but failed to reload")
}
