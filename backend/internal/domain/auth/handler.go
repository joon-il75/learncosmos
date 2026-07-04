package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/demoaccount"
)

type Handler struct {
	svc       *Service
	demoRepo  *demoaccount.Repository
	demoLogin DemoLoginConfig
}

type DemoLoginConfig struct {
	Enabled      bool
	Secret       string
	AllowedHosts []string
}

const (
	profileAvatarMaxBytes = int64(1 << 20)
	profileAvatarMaxSide  = 1024
	profileAvatarDir      = "/home/weaver/learnweaver/backend/storage/avatars"
	profileAvatarURLBase  = "/api/v1/users/avatars"
)

func NewHandler(svc *Service, demoRepo *demoaccount.Repository, demoLogin DemoLoginConfig) *Handler {
	return &Handler{svc: svc, demoRepo: demoRepo, demoLogin: demoLogin}
}

func apiError(errorCode, message string) gin.H {
	return gin.H{"error": message, "error_code": errorCode}
}

var providers = []Provider{ProviderGoogle, ProviderKakao, ProviderNaver}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	requiredConsents := h.svc.RequireRequiredConsents()

	for _, p := range providers {
		provider := p
		r.GET("/auth/"+string(provider)+"/login", func(c *gin.Context) {
			h.loginHandler(c, provider)
		})
		r.GET("/auth/"+string(provider)+"/callback", func(c *gin.Context) {
			h.callbackHandler(c, provider)
		})
	}
	r.POST("/auth/demo/login", demoLoginRateLimitMiddleware(), h.demoLoginHandler)
	r.POST("/auth/refresh", h.refreshHandler)
	r.POST("/auth/logout", h.logoutHandler)
	r.GET("/auth/me", h.svc.JWTMiddleware(), h.meHandler)
	r.GET("/users/avatars/:filename", h.getUserAvatarHandler)
	r.PATCH("/users/me", h.svc.JWTMiddleware(), requiredConsents, h.patchMeHandler)
	r.POST("/users/me/avatar", h.svc.JWTMiddleware(), requiredConsents, h.uploadMyAvatarHandler)
	r.DELETE("/users/me", h.svc.JWTMiddleware(), h.withdrawMeHandler)
	r.GET("/users/me/consents", h.svc.JWTMiddleware(), h.getMyConsentStatusHandler)
	r.POST("/users/me/consents", h.svc.JWTMiddleware(), h.recordMyConsentsHandler)
	r.GET("/users/me/preferences", h.svc.JWTMiddleware(), h.getMyPreferencesHandler)
	r.PATCH("/users/me/preferences", h.svc.JWTMiddleware(), h.patchMyPreferencesHandler)
	r.GET("/users/me/ai-settings", h.svc.JWTMiddleware(), requiredConsents, h.getMyAISettingsHandler)
	r.GET("/users/me/ai-usage", h.svc.JWTMiddleware(), requiredConsents, h.getMyAIUsageHandler)
	r.GET("/users/me/point-usage", h.svc.JWTMiddleware(), requiredConsents, h.getMyPointUsageHandler)
	r.PATCH("/users/me/ai-settings", h.svc.JWTMiddleware(), requiredConsents, h.patchMyAISettingsHandler)
	r.DELETE("/users/me/ai-settings", h.svc.JWTMiddleware(), requiredConsents, h.deleteMyAISettingsHandler)
	r.PATCH("/users/me/ai-settings/enabled", h.svc.JWTMiddleware(), requiredConsents, h.patchMyAISettingsEnabledHandler)
	r.POST("/users/me/ai-settings/validate", h.svc.JWTMiddleware(), requiredConsents, aiValidationRateLimitMiddleware(), h.validateMyAISettingsHandler)
}
