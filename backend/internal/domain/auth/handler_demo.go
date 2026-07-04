package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/demoaccount"
)

type demoLoginRequest struct {
	Code string `json:"code"`
}

func (h *Handler) demoLoginHandler(c *gin.Context) {
	if h.demoRepo == nil || !h.demoLogin.Enabled || strings.TrimSpace(h.demoLogin.Secret) == "" {
		c.JSON(http.StatusForbidden, apiError("demo_login_disabled", "demo login is disabled"))
		return
	}
	if !h.isDemoLoginHostAllowed(c) {
		c.JSON(http.StatusForbidden, apiError("demo_login_disabled", "demo login is disabled"))
		return
	}

	var req demoLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Code) == "" {
		c.JSON(http.StatusBadRequest, apiError("invalid_demo_code", "invalid demo code"))
		return
	}

	demoUser, err := h.demoRepo.AuthenticateByCode(c.Request.Context(), req.Code, h.demoLogin.Secret)
	if err != nil {
		switch {
		case errors.Is(err, demoaccount.ErrInvalidCode):
			c.JSON(http.StatusUnauthorized, apiError("invalid_demo_code", "invalid demo code"))
		case errors.Is(err, demoaccount.ErrExpiredAccount):
			c.JSON(http.StatusForbidden, apiError("expired_demo_code", "expired demo code"))
		case errors.Is(err, demoaccount.ErrDisabledAccount):
			c.JSON(http.StatusForbidden, apiError("disabled_demo_account", "disabled demo account"))
		default:
			c.JSON(http.StatusForbidden, apiError("demo_account_unavailable", "demo account unavailable"))
		}
		return
	}

	role := Role(demoUser.Role)
	if role == "" {
		role = RoleLearner
	}
	user := &User{
		ID:                       demoUser.ID,
		Email:                    demoUser.Email,
		Nickname:                 demoUser.Nickname,
		Role:                     role,
		PremiumAccess:            demoUser.PremiumAccess,
		AvatarURL:                demoUser.AvatarURL,
		DisplayID:                demoUser.DisplayID,
		Status:                   demoUser.Status,
		Provider:                 "demo",
		UILocale:                 demoUser.UILocale,
		LearningLanguage:         demoUser.LearningLanguage,
		LanguageSetupCompletedAt: demoUser.LanguageSetupCompletedAt,
		LastLoginAt:              demoUser.LastLoginAt,
		WithdrawnAt:              demoUser.WithdrawnAt,
		ReactivatedAt:            demoUser.ReactivatedAt,
		CreatedAt:                demoUser.CreatedAt,
	}

	resp, refreshToken, err := h.svc.IssueSessionForUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("demo_account_unavailable", "demo account unavailable"))
		return
	}

	maxAge := int(30 * 24 * time.Hour / time.Second)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookie, refreshToken, maxAge, "/", "", true, true)
	c.SetCookie("access_token", resp.AccessToken, 60*15, "/", "", true, true)
	c.SetCookie("is_logged_in", "1", maxAge, "/", "", true, false)

	redirectURL := "/dashboard"
	if resp.User.RequiredConsentPending {
		redirectURL = "/agreements"
		if resp.User.UILocale == "en" {
			redirectURL = "/en/agreements"
		}
	}
	if resp.User.LanguageSetupRequired {
		redirectURL = "/language-setup"
		if resp.User.UILocale == "en" {
			redirectURL = "/en/language-setup"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectURL,
		"user":         resp.User,
	})
}

func (h *Handler) isDemoLoginHostAllowed(c *gin.Context) bool {
	if len(h.demoLogin.AllowedHosts) == 0 {
		return true
	}

	host := normalizedHost(c.Request.Host)
	if host == "" {
		return false
	}
	for _, allowed := range h.demoLogin.AllowedHosts {
		if normalizedHost(allowed) == host {
			return true
		}
	}
	return false
}
