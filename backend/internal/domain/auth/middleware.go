package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"

	superAdminAccessTokenCookieName = "super_admin_access_token"
)

func (s *Service) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		claims, err := s.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		user, err := s.repo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil || user.Status == "withdrawn" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

// JWTMiddlewareHeaderOnly: 슈퍼관리자 라우트에서 사용되는 전용 미들웨어.
// super-admin access_token HttpOnly 쿠키를 우선 사용하고, Authorization fallback은 명시 opt-in일 때만 허용한다.
func (s *Service) JWTMiddlewareHeaderOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractSuperAdminToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		claims, err := s.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

// extractToken: 쿠키를 우선 확인하고, Authorization Bearer fallback은 production에서 기본 비활성화한다.
func extractToken(c *gin.Context) (string, error) {
	// 1. 쿠키에서 access_token 확인 (가장 우선)
	if token, err := c.Cookie("access_token"); err == nil && token != "" {
		return token, nil
	}

	if authBearerFallbackEnabled() {
		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			return strings.TrimPrefix(header, "Bearer "), nil
		}
	}

	return "", apperr.ErrUnauthorized
}

// extractSuperAdminToken은 super-admin 인증에 특화된 토큰 추출기.
// super_admin_access_token 쿠키를 우선 사용하고, fallback은 명시 opt-in일 때만 허용한다.
func extractSuperAdminToken(c *gin.Context) (string, error) {
	if token, err := c.Cookie(superAdminAccessTokenCookieName); err == nil && token != "" {
		return token, nil
	}
	if authBearerFallbackEnabled() {
		return extractToken(c)
	}
	return "", apperr.ErrUnauthorized
}

func authBearerFallbackEnabled() bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv("AUTH_BEARER_FALLBACK_ENABLED")))
	switch raw {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	}
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production")
}

func RequireRole(required Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextKeyRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		role, ok := roleVal.(Role)
		if !ok || roleLevel(role) < roleLevel(required) {
			c.AbortWithStatusJSON(http.StatusForbidden, apiError("forbidden", apperr.ErrForbidden.Message))
			return
		}
		c.Next()
	}
}

func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	ipMap := make(map[string]struct{}, len(allowedIPs))
	for _, ip := range allowedIPs {
		ipMap[strings.TrimSpace(ip)] = struct{}{}
	}
	return func(c *gin.Context) {
		if _, ok := ipMap[c.ClientIP()]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, apiError("forbidden", apperr.ErrForbidden.Message))
			return
		}
		c.Next()
	}
}

func GetCurrentUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(ContextKeyUserID)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

func GetCurrentRole(c *gin.Context) (Role, bool) {
	v, ok := c.Get(ContextKeyRole)
	if !ok {
		return "", false
	}
	role, ok := v.(Role)
	return role, ok
}

func (s *Service) RequireRequiredConsents() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetCurrentUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		role, _ := GetCurrentRole(c)
		if role != RoleLearner {
			c.Next()
			return
		}

		status, err := s.repo.GetConsentStatus(c.Request.Context(), userID, "ko")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("consent_status_load_failed", apperr.ErrUnauthorized.Message))
			return
		}

		if status.RequiredConsentPending {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":                    "required_consents_pending",
				"error_code":               "required_consents_pending",
				"required_consent_pending": true,
			})
			return
		}

		c.Next()
	}
}

func roleLevel(r Role) int {
	switch r {
	case RoleSuperAdmin:
		return 3
	case RoleAdmin:
		return 2
	case RoleLearner:
		return 1
	}
	return 0
}
