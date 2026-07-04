package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
	"github.com/learnweaver/backend/internal/pkg/keyvalidator"
)

type UpdateUserAISettingsRequest struct {
	Provider    string  `json:"provider"`
	APIKey      string  `json:"api_key"`
	EndpointURL *string `json:"endpoint_url"`
	IsEnabled   *bool   `json:"is_enabled"`
}

type UpdateUserAISettingsEnabledRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

func (h *Handler) getMyAISettingsHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	settings, err := h.svc.repo.GetUserAISettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_settings_load_failed", "AI 설정 조회 실패"))
		return
	}
	if settings == nil {
		c.JSON(http.StatusOK, UserAISettings{Mode: "managed_credit", Provider: "openai", HasAPIKey: false, IsEnabled: false})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *Handler) patchMyAISettingsHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req UpdateUserAISettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}

	provider := strings.TrimSpace(strings.ToLower(req.Provider))
	if provider == "" {
		c.JSON(http.StatusBadRequest, apiError("ai_provider_required", "AI 제공자를 선택해주세요"))
		return
	}

	allowedProviders := map[string]bool{
		"openai":     true,
		"anthropic":  true,
		"google":     true,
		"grok":       true,
		"solar":      true,
		"hyperclova": true,
		"llama":      true,
		"exaone":     true,
	}
	if !allowedProviders[provider] {
		c.JSON(http.StatusBadRequest, apiError("ai_provider_unsupported", "지원하지 않는 AI 제공자입니다"))
		return
	}

	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, apiError("ai_api_key_required", "API 키를 입력해주세요"))
		return
	}

	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	if err := h.svc.repo.UpsertUserAISettings(c.Request.Context(), userID, provider, apiKey, req.EndpointURL, isEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_settings_save_failed", "AI 설정 저장 실패"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI 설정이 저장되었습니다"})
}

func (h *Handler) patchMyAISettingsEnabledHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req UpdateUserAISettingsEnabledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}

	settings, err := h.svc.repo.GetUserAISettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_settings_load_failed", "AI 설정 조회 실패"))
		return
	}
	if settings == nil || !settings.HasAPIKey {
		c.JSON(http.StatusBadRequest, apiError("ai_settings_missing_key", "저장된 BYOK 키가 없습니다"))
		return
	}

	if err := h.svc.repo.UpdateUserAISettingsEnabled(c.Request.Context(), userID, req.IsEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_settings_enabled_update_failed", "BYOK 사용 상태 변경 실패"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "BYOK 사용 상태가 변경되었습니다",
		"is_enabled": req.IsEnabled,
	})
}

func (h *Handler) deleteMyAISettingsHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	if err := h.svc.repo.DeleteUserAISettings(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_settings_delete_failed", "BYOK 키 삭제 실패"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "저장된 BYOK 키가 삭제되었습니다"})
}

func (h *Handler) validateMyAISettingsHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req UpdateUserAISettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}

	result := keyvalidator.Validate(c.Request.Context(), req.Provider, req.APIKey, req.EndpointURL)
	_ = h.svc.repo.RecordUserAIValidation(c.Request.Context(), userID, strings.TrimSpace(strings.ToLower(req.Provider)), result.Valid, result.Message)
	status := http.StatusOK
	if !result.Valid {
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{
		"valid":   result.Valid,
		"message": result.Message,
	})
}
