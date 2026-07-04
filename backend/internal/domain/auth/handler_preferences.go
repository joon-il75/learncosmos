package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

type RecordConsentsRequest struct {
	AgreeTerms   bool   `json:"agree_terms"`
	AgreePrivacy bool   `json:"agree_privacy"`
	Locale       string `json:"locale"`
}

type UpdateLanguagePreferencesRequest struct {
	UILocale         string `json:"ui_locale"`
	LearningLanguage string `json:"learning_language"`
}

func (h *Handler) getMyConsentStatusHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	locale := normalizeOAuthLocale(c.Query("locale"))
	status, err := h.svc.repo.GetConsentStatus(c.Request.Context(), userID, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("consent_status_load_failed", "동의 상태 조회 실패"))
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *Handler) recordMyConsentsHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req RecordConsentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}
	if !req.AgreeTerms || !req.AgreePrivacy {
		c.JSON(http.StatusBadRequest, apiError("consent_required_terms_privacy", "서비스 이용약관과 개인정보처리방침에 모두 동의해야 합니다"))
		return
	}

	locale := normalizeOAuthLocale(req.Locale)
	if err := h.svc.repo.RecordRequiredConsents(c.Request.Context(), userID, locale, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("consent_save_failed", "동의 저장 실패"))
		return
	}

	status, err := h.svc.repo.GetConsentStatus(c.Request.Context(), userID, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("consent_status_load_failed", "동의 상태 조회 실패"))
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *Handler) getMyPreferencesHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	user, err := h.svc.repo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("preferences_load_failed", "사용자 설정 조회 실패"))
		return
	}

	defaultLocale := normalizeOAuthLocale(c.Query("locale"))
	uiLocale := defaultLocale
	if strings.TrimSpace(user.UILocale) != "" {
		uiLocale = normalizeOAuthLocale(user.UILocale)
	}
	learningLanguage := defaultLocale
	if strings.TrimSpace(user.LearningLanguage) != "" {
		learningLanguage = normalizeOAuthLocale(user.LearningLanguage)
	}
	c.JSON(http.StatusOK, gin.H{
		"ui_locale":               uiLocale,
		"learning_language":       learningLanguage,
		"language_setup_required": user.LanguageSetupCompletedAt == nil,
	})
}

func (h *Handler) patchMyPreferencesHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req UpdateLanguagePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}

	uiLocale := normalizeOAuthLocale(req.UILocale)
	learningLanguage := normalizeOAuthLocale(req.LearningLanguage)
	if err := h.svc.repo.UpdateLanguagePreferences(c.Request.Context(), userID, uiLocale, learningLanguage); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("preferences_save_failed", "언어 설정 저장 실패"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ui_locale":               uiLocale,
		"learning_language":       learningLanguage,
		"language_setup_required": false,
	})
}
