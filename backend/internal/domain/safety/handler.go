package safety

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type aiOutputReviewDecisionRequest struct {
	CandidateKey string `json:"candidate_key"`
	Status       string `json:"status"`
	Note         string `json:"note"`
}

type ruleRequest struct {
	RuleType    *string `json:"rule_type"`
	Pattern     *string `json:"pattern"`
	RiskType    *string `json:"risk_type"`
	Action      *string `json:"action"`
	Locale      *string `json:"locale"`
	Description *string `json:"description"`
}

type moderateRequest struct {
	TargetType string     `json:"target_type"`
	TargetID   *uuid.UUID `json:"target_id"`
	Text       string     `json:"text"`
	Locale     string     `json:"locale"`
}

func (h *Handler) Moderate(c *gin.Context) {
	userID, _ := getUserID(c)
	var req moderateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "error_code": "invalid_request"})
		return
	}
	targetType := strings.TrimSpace(req.TargetType)
	if targetType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_type is required", "error_code": "invalid_request"})
		return
	}
	result, err := h.service.Moderate(c.Request.Context(), ModerateInput{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   req.TargetID,
		Text:       req.Text,
		Locale:     req.Locale,
		Route:      c.FullPath(),
		Metadata:   map[string]any{"content_length": len([]rune(req.Text))},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to moderate input", "error_code": "safety_moderation_failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ListLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.ListLogs(c.Request.Context(), ListLogsInput{
		Action:     c.Query("action"),
		RiskType:   c.Query("risk_type"),
		TargetType: c.Query("target_type"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list moderation logs", "error_code": "safety_logs_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": items})
}

func (h *Handler) AIOutputSummary(c *gin.Context) {
	summary, err := h.service.GetAIOutputSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build AI output summary", "error_code": "safety_ai_output_summary_failed"})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *Handler) AIOutputReviewCandidates(c *gin.Context) {
	candidates, err := h.service.GetAIOutputReviewCandidates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build AI output review candidates", "error_code": "safety_ai_output_review_candidates_failed"})
		return
	}
	c.JSON(http.StatusOK, candidates)
}

func (h *Handler) SaveAIOutputReviewDecision(c *gin.Context) {
	var req aiOutputReviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "error_code": "invalid_request"})
		return
	}
	decision, err := h.service.SaveAIOutputReviewDecision(c.Request.Context(), AIOutputReviewDecisionInput{
		CandidateKey: req.CandidateKey,
		Status:       req.Status,
		Note:         req.Note,
		ReviewedBy:   getActorID(c),
		Metadata:     map[string]any{"source": "super_admin_safety"},
	})
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"review": decision})
		return
	}
	if errors.Is(err, ErrInvalidAIOutputReviewDecision) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid AI output review decision", "error_code": "invalid_ai_output_review_decision"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save AI output review decision", "error_code": "safety_ai_output_review_save_failed"})
}

func getUserID(c *gin.Context) (*uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return nil, false
	}
	idStr, ok := raw.(string)
	if !ok {
		return nil, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, false
	}
	return &id, true
}

func (h *Handler) ListRules(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var active *bool
	if raw := strings.TrimSpace(c.Query("is_active")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid is_active", "error_code": "invalid_request"})
			return
		}
		active = &parsed
	}
	items, err := h.service.ListRules(c.Request.Context(), ListRulesInput{
		IsActive: active,
		Action:   c.Query("action"),
		RiskType: c.Query("risk_type"),
		RuleType: c.Query("rule_type"),
		Locale:   c.Query("locale"),
		Query:    c.Query("q"),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list moderation rules", "error_code": "safety_rules_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rules": items})
}

func (h *Handler) CreateRule(c *gin.Context) {
	input, ok := bindRuleMutation(c, true)
	if !ok {
		return
	}
	rule, err := h.service.CreateRule(c.Request.Context(), input, getActorID(c))
	writeRuleMutationResponse(c, rule, err)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	ruleID, ok := parseRuleID(c)
	if !ok {
		return
	}
	input, ok := bindRuleMutation(c, false)
	if !ok {
		return
	}
	rule, err := h.service.UpdateRule(c.Request.Context(), ruleID, input, getActorID(c))
	writeRuleMutationResponse(c, rule, err)
}

func (h *Handler) DeactivateRule(c *gin.Context) {
	ruleID, ok := parseRuleID(c)
	if !ok {
		return
	}
	rule, err := h.service.SetRuleActive(c.Request.Context(), ruleID, false, getActorID(c))
	writeRuleMutationResponse(c, rule, err)
}

func (h *Handler) ReactivateRule(c *gin.Context) {
	ruleID, ok := parseRuleID(c)
	if !ok {
		return
	}
	rule, err := h.service.SetRuleActive(c.Request.Context(), ruleID, true, getActorID(c))
	writeRuleMutationResponse(c, rule, err)
}

func bindRuleMutation(c *gin.Context, requireAll bool) (RuleMutationInput, bool) {
	var req ruleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "error_code": "invalid_request"})
		return RuleMutationInput{}, false
	}
	input := RuleMutationInput{
		RuleType:    cleanOptional(req.RuleType),
		Pattern:     cleanOptional(req.Pattern),
		RiskType:    cleanOptional(req.RiskType),
		Action:      cleanOptional(req.Action),
		Locale:      cleanOptional(req.Locale),
		Description: cleanOptional(req.Description),
	}
	if requireAll && (input.RuleType == nil || input.Pattern == nil || input.RiskType == nil || input.Action == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rule_type, pattern, risk_type, action are required", "error_code": "invalid_request"})
		return RuleMutationInput{}, false
	}
	if input.RuleType != nil && *input.RuleType != "keyword" && *input.RuleType != "regex" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rule_type must be keyword or regex", "error_code": "invalid_rule_type"})
		return RuleMutationInput{}, false
	}
	if input.Action != nil && *input.Action != string(ActionSoftWarn) && *input.Action != string(ActionBlock) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action must be soft_warn or block", "error_code": "invalid_action"})
		return RuleMutationInput{}, false
	}
	if input.Pattern != nil && *input.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pattern is required", "error_code": "invalid_pattern"})
		return RuleMutationInput{}, false
	}
	if input.RiskType != nil && *input.RiskType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "risk_type is required", "error_code": "invalid_risk_type"})
		return RuleMutationInput{}, false
	}
	if input.RuleType != nil && *input.RuleType == "regex" && input.Pattern != nil {
		if _, err := regexp.Compile(normalizeText(*input.Pattern)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid regex pattern", "error_code": "invalid_regex"})
			return RuleMutationInput{}, false
		}
	}
	return input, true
}

func cleanOptional(value *string) *string {
	if value == nil {
		return nil
	}
	cleaned := strings.TrimSpace(*value)
	return &cleaned
}

func parseRuleID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(strings.TrimSpace(c.Param("rule_id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule_id", "error_code": "invalid_rule_id"})
		return uuid.Nil, false
	}
	return id, true
}

func writeRuleMutationResponse(c *gin.Context, rule ModerationRule, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"rule": rule})
		return
	}
	if errors.Is(err, ErrRuleDuplicate) {
		c.JSON(http.StatusConflict, gin.H{"error": "moderation rule already exists", "error_code": "moderation_rule_duplicate"})
		return
	}
	if errors.Is(err, ErrRuleNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "moderation rule not found", "error_code": "moderation_rule_not_found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save moderation rule", "error_code": "safety_rule_save_failed"})
}

func getActorID(c *gin.Context) string {
	raw, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	if id, ok := raw.(string); ok {
		return strings.TrimSpace(id)
	}
	return ""
}
