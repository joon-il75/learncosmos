package goal

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) applyInterviewResponse(profile *GoalProfile, aiResp *InterviewAIResponse) {
	if profile == nil || aiResp == nil {
		return
	}
	profile.InterviewState = aiResp.NextState
	h.service.ApplyExtracted(profile, aiResp.Extracted)
	proposedGoal := strings.TrimSpace(derefString(aiResp.ProposedGoal))
	if extracted := h.service.ExtractProposedGoalFromReply(aiResp.Message); aiResp.NextState == StateProposingGoal && extracted != "" {
		proposedGoal = extracted
	}
	if proposedGoal != "" {
		profile.ConfirmedGoal = &proposedGoal
	}
	h.service.RefreshLearningIntentProfile(profile)
	profile.Messages = append(profile.Messages, InterviewMessage{
		Role: "lumi", Content: aiResp.Message,
	})
}

func goalConfirmedMessage(language string) string {
	if normalizeLearningLanguage(language) == "en" {
		return "Your goal is set. Shall we create your exploration plan now?"
	}
	return "목표가 확정됐어요! 이제 탐험계획을 만들어볼까요?"
}

// GET /api/v1/goals/active
func (h *Handler) GetActiveGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusOK, gin.H{"goal": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// DELETE /api/v1/goals/active
func (h *Handler) ResetActiveGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.repo.DeactivatePredraftGoalsByUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reset": true})
}

// GET /api/v1/course-drafts/:id/goal
func (h *Handler) GetGoal(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusOK, gin.H{"goal": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/start
func (h *Handler) StartInterviewForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req StartInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"goal": existing})
		return
	}

	language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	if h.enforceGoalSafety(c, userID, "goal_start", nil, req.UserIntent, language, map[string]any{"field": "user_intent"}) {
		return
	}

	profile := &GoalProfile{
		CourseDraftID:  nil,
		UserID:         userID,
		UserIntent:     req.UserIntent,
		Language:       language,
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: req.UserIntent},
		},
		Version: 1,
	}

	if h.shouldUseGoalInterviewWorker() {
		if err := h.repo.CreateGoalProfile(c.Request.Context(), profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create goal profile"})
			return
		}
		if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.UserIntent), "failed to save interview") {
			return
		}
		if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save interview"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"goal": profile})
		return
	}

	llmCtx, llmCancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer llmCancel()
	resp := h.resolveInterviewResponse(llmCtx, userID, profile, req.UserIntent)
	h.applyInterviewResponse(profile, resp)

	if err := h.repo.CreateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create goal profile"})
		return
	}
	h.observeGoalInterviewAIOutput(c.Request.Context(), userID, profile, resp, "goal_interview_start_sync")
	c.JSON(http.StatusCreated, gin.H{"goal": profile})
}

// POST /api/v1/course-drafts/:id/goal/start
func (h *Handler) StartInterview(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	var req StartInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 기존 인터뷰가 이미 있으면 반환
	existing, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"goal": existing})
		return
	}

	language := h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	if h.enforceGoalSafety(c, userID, "course_draft_goal_start", &draftID, req.UserIntent, language, map[string]any{"field": "user_intent", "draft_id": draftID.String()}) {
		return
	}

	// 새 goal profile 생성
	profile := &GoalProfile{
		CourseDraftID:  &draftID,
		UserID:         userID,
		UserIntent:     req.UserIntent,
		Language:       language,
		InterviewState: StateClarifying,
		Messages: []InterviewMessage{
			{Role: "user", Content: req.UserIntent},
		},
		Version: 1,
	}

	if h.shouldUseGoalInterviewWorker() {
		if err := h.repo.CreateGoalProfile(c.Request.Context(), profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create goal profile"})
			return
		}
		if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.UserIntent), "failed to save interview") {
			return
		}
		if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save interview"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"goal": profile})
		return
	}

	llmCtx, llmCancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer llmCancel()
	resp := h.resolveInterviewResponse(llmCtx, userID, profile, req.UserIntent)
	h.applyInterviewResponse(profile, resp)

	if err := h.repo.CreateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create goal profile"})
		return
	}
	h.observeGoalInterviewAIOutput(c.Request.Context(), userID, profile, resp, "course_draft_goal_interview_start_sync")
	c.JSON(http.StatusCreated, gin.H{"goal": profile})
}

// POST /api/v1/goals/interview
func (h *Handler) SendMessageForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if profile.InterviewState == StateConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal already confirmed"})
		return
	}
	if strings.TrimSpace(profile.Language) == "" {
		profile.Language = h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	}
	if h.enforceGoalSafety(c, userID, "goal_interview", &profile.ID, req.Message, profile.Language, map[string]any{"field": "message", "goal_profile_id": profile.ID.String()}) {
		return
	}

	profile.Messages = append(profile.Messages, InterviewMessage{Role: "user", Content: req.Message})

	if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.Message), "failed to save interview") {
		return
	}

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save interview"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/course-drafts/:id/goal/interview
func (h *Handler) SendMessage(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if profile.InterviewState == StateConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal already confirmed"})
		return
	}
	if strings.TrimSpace(profile.Language) == "" {
		profile.Language = h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	}
	if h.enforceGoalSafety(c, userID, "course_draft_goal_interview", &draftID, req.Message, profile.Language, map[string]any{"field": "message", "draft_id": draftID.String(), "goal_profile_id": profile.ID.String()}) {
		return
	}

	profile.Messages = append(profile.Messages, InterviewMessage{Role: "user", Content: req.Message})

	if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.Message), "failed to save interview") {
		return
	}

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save interview"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/confirm
func (h *Handler) ConfirmGoalForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ConfirmGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}

	profile.ConfirmedGoal = &req.ConfirmedGoal
	profile.InterviewState = StateConfirmed
	h.service.RefreshLearningIntentProfile(profile)
	profile.Messages = append(profile.Messages, InterviewMessage{
		Role:    "lumi",
		Content: goalConfirmedMessage(profile.Language),
	})

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/course-drafts/:id/goal/confirm
func (h *Handler) ConfirmGoal(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	var req ConfirmGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}

	profile.ConfirmedGoal = &req.ConfirmedGoal
	profile.InterviewState = StateConfirmed
	h.service.RefreshLearningIntentProfile(profile)
	if profile.Version > 1 {
		profile.InterviewState = StateAwaitingRebuildDecision
	}
	profile.Messages = append(profile.Messages, InterviewMessage{
		Role:    "lumi",
		Content: goalConfirmedMessage(profile.Language),
	})

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm goal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/attach-draft
func (h *Handler) AttachDraft(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AttachDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), req.DraftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		attached, attachedErr := h.repo.GetActiveGoal(c.Request.Context(), req.DraftID)
		if attachedErr == nil {
			c.JSON(http.StatusOK, gin.H{"goal": attached})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if profile.InterviewState != StateConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal not confirmed"})
		return
	}

	if err := h.repo.AttachDraftToGoal(c.Request.Context(), profile.ID, userID, req.DraftID); err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach draft"})
		return
	}

	profile.CourseDraftID = &req.DraftID
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/revise
func (h *Handler) ReviseGoalForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ReviseGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if profile.InterviewState != StateConfirmed && profile.InterviewState != StateAwaitingRebuildDecision {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal not yet confirmed"})
		return
	}

	BeginGoalRevision(profile, req.Message)

	if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.Message), "failed to save revision") {
		return
	}

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save revision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/course-drafts/:id/goal/revise
func (h *Handler) ReviseGoal(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	var req ReviseGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal interview not started"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if profile.InterviewState != StateConfirmed && profile.InterviewState != StateAwaitingRebuildDecision {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal not yet confirmed"})
		return
	}

	BeginGoalRevision(profile, req.Message)

	if handleGoalInterviewTurnError(c, profile, h.applyGoalInterviewTurn(c.Request.Context(), userID, profile, req.Message), "failed to save revision") {
		return
	}

	if err := h.repo.UpdateGoalProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save revision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/cancel-revision
func (h *Handler) CancelGoalRevisionForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	profile, err = h.repo.CancelGoalRevision(c.Request.Context(), profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel goal revision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/course-drafts/:id/goal/cancel-revision
func (h *Handler) CancelGoalRevision(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	profile, err = h.repo.CancelGoalRevision(c.Request.Context(), profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel goal revision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"goal": profile})
}

// POST /api/v1/goals/rebuild-decision
func (h *Handler) RebuildDecisionForUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req RebuildDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Decision = NormalizeRebuildDecision(req.Decision)
	if req.Decision != RebuildKeepStructure && req.Decision != RebuildAll {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid decision value"})
		return
	}

	profile, err := h.repo.GetActiveGoalByUser(c.Request.Context(), userID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}

	if err := h.repo.UpdateRebuildDecision(c.Request.Context(), profile.ID, req.Decision); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rebuild decision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "decision": req.Decision})
}

// POST /api/v1/course-drafts/:id/goal/rebuild-decision
func (h *Handler) RebuildDecision(c *gin.Context) {
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
	if ok, _ := h.repo.VerifyCourseDraftOwner(c.Request.Context(), draftID, userID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	var req RebuildDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Decision = NormalizeRebuildDecision(req.Decision)
	if req.Decision != RebuildKeepStructure && req.Decision != RebuildAll {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid decision value"})
		return
	}

	profile, err := h.repo.GetActiveGoal(c.Request.Context(), draftID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}

	if err := h.repo.UpdateRebuildDecision(c.Request.Context(), profile.ID, req.Decision); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rebuild decision"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "decision": req.Decision})
}
