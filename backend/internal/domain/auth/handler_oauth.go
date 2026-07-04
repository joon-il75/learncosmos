package auth

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

func (h *Handler) loginHandler(c *gin.Context, provider Provider) {
	locale := normalizeOAuthLocale(c.Query("locale"))
	authURL, err := h.svc.GenerateAuthURL(c.Request.Context(), provider, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth url"})
		return
	}

	if redirectAfter := normalizeRedirectAfter(c.Query("redirect_after")); redirectAfter != "" {
		c.SetCookie("redirect_after", redirectAfter, 600, "/", "", true, true)
	}
	c.SetCookie("visitor_locale", locale, 600, "/", "", true, false)

	c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) callbackHandler(c *gin.Context, provider Provider) {
	frontendURL := h.frontendURL(c)

	code := c.Query("code")
	state := c.Query("state")

	locale, err := h.svc.ValidateState(c.Request.Context(), state)
	if err != nil {
		c.Redirect(http.StatusFound, frontendURL+"/login?error=invalid_state")
		return
	}
	c.SetCookie("visitor_locale", locale, 600, "/", "", true, false)

	resp, refreshToken, err := h.svc.HandleOAuthCallback(c.Request.Context(), provider, code, state)
	if err != nil {
		log.Printf("[callback error] provider=%s err=%v", provider, err)
		c.Redirect(http.StatusFound, frontendURL+"/login?error=oauth_failed")
		return
	}

	maxAge := int(30 * 24 * time.Hour / time.Second)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookie, refreshToken, maxAge, "/", "", true, true)
	c.SetCookie("access_token", resp.AccessToken, 60*15, "/", "", true, true)
	c.SetCookie("is_logged_in", "1", maxAge, "/", "", true, false) // JS 읽기 가능 (비HttpOnly)

	// role에 따라 분기 (admin은 항상 admin/dashboard로)
	if resp.User.Role == RoleAdmin {
		c.SetCookie("redirect_after", "", -1, "/", "", true, true)
		c.Redirect(http.StatusFound, frontendURL+"/admin/dashboard")
		return
	}

	redirectAfter, err := c.Cookie("redirect_after")
	if err != nil || redirectAfter == "" {
		redirectAfter = "/dashboard"
	}
	redirectAfter = normalizeRedirectAfter(redirectAfter)
	if redirectAfter == "" {
		redirectAfter = "/dashboard"
	}

	if resp.User.LanguageSetupRequired {
		c.SetCookie("redirect_after", "", -1, "/", "", true, true)
		setupPath := "/language-setup"
		if locale == "en" {
			setupPath = "/en/language-setup"
		}
		c.Redirect(http.StatusFound, frontendURL+setupPath+"?redirect_after="+url.QueryEscape(redirectAfter))
		return
	}

	if resp.User.RequiredConsentPending {
		agreementsPath := "/agreements"
		if resp.User.UILocale == "en" || locale == "en" {
			agreementsPath = "/en/agreements"
		}
		c.SetCookie("redirect_after", "", -1, "/", "", true, true)
		c.Redirect(http.StatusFound, frontendURL+agreementsPath+"?redirect_after="+url.QueryEscape(redirectAfter))
		return
	}

	if redirectAfter != "" {
		c.SetCookie("redirect_after", "", -1, "/", "", true, true)
		c.Redirect(http.StatusFound, redirectAfter)
		return
	}

	c.Redirect(http.StatusFound, frontendURL+"/dashboard")
}

func (h *Handler) refreshHandler(c *gin.Context) {
	token, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	accessToken, err := h.svc.RefreshAccessToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// 쿠키 보안 정책(SameSite)을 콜백과 일치시킴
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 60*15, "/", "", true, true)
	c.SetCookie("is_logged_in", "1", int(30*24*time.Hour/time.Second), "/", "", true, false)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func (h *Handler) logoutHandler(c *gin.Context) {
	// 1. 서버 사이드 세션 무효화 (필요 시)
	userID, _ := c.Get(ContextKeyUserID)
	if uid, ok := userID.(string); ok && uid != "" {
		_ = h.svc.Logout(c.Request.Context(), uid)
	}

	// 2. 브라우저의 모든 인증 관련 쿠키 삭제
	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) meHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	user, err := h.svc.repo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("user_fetch_failed", "failed to fetch user"))
		return
	}

	displayID := ""
	if user.DisplayID != nil {
		displayID = *user.DisplayID
	}
	consentStatus, err := h.svc.repo.GetConsentStatus(c.Request.Context(), userID, normalizeOAuthLocale(user.UILocale))
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("consent_status_load_failed", "failed to fetch consent status"))
		return
	}
	freePoints, paidPoints, err := h.svc.repo.GetPointBalances(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("points_load_failed", "failed to fetch user points"))
		return
	}
	log.Printf("[auth/me] user_id=%s free_points=%d paid_points=%d total_points=%d", user.ID, freePoints, paidPoints, freePoints+paidPoints)
	c.JSON(http.StatusOK, UserInfo{
		ID:                     user.ID,
		Email:                  user.Email,
		DisplayID:              displayID,
		Nickname:               user.Nickname,
		AvatarURL:              user.AvatarURL,
		Role:                   user.Role,
		PremiumAccess:          user.PremiumAccess,
		Provider:               user.Provider,
		CreatedAt:              user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		FreePoints:             freePoints,
		PaidPoints:             paidPoints,
		TotalPoints:            freePoints + paidPoints,
		TermsAgreed:            consentStatus.TermsAgreed,
		PrivacyAgreed:          consentStatus.PrivacyAgreed,
		RequiredConsentPending: consentStatus.RequiredConsentPending,
		UILocale:               user.UILocale,
		LearningLanguage:       user.LearningLanguage,
		LanguageSetupRequired:  user.LanguageSetupCompletedAt == nil,
	})
}

func (h *Handler) clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie(refreshTokenCookie, "", -1, "/", "", true, true)
	c.SetCookie("is_logged_in", "", -1, "/", "", true, false)
	c.SetCookie("admin_totp_access_token", "", -1, "/", "", false, true)
	c.SetCookie("admin_totp_access_token", "", -1, "/", "", true, true)
	c.SetCookie(superAdminAccessTokenCookieName, "", -1, "/", "", false, true)
	c.SetCookie(superAdminAccessTokenCookieName, "", -1, "/", "", true, true)
}

func normalizeRedirectAfter(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.ContainsAny(trimmed, "\\\r\n\t") {
		return ""
	}
	if !strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, "//") {
		return ""
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed == nil || parsed.Path == "" || !strings.HasPrefix(parsed.Path, "/") {
		return ""
	}
	if strings.HasPrefix(parsed.Path, "//") {
		return ""
	}
	return parsed.RequestURI()
}
