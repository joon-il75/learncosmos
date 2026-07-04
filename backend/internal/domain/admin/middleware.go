package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/config"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/pkg/apperr"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

const superAdminAccessTokenCookieName = "super_admin_access_token"
const adminTOTPSessionCookieName = "admin_totp_access_token"

// TOTPMiddleware: super_admin 전용 라우트 보호 (기존 super-admin 패널용)
func TOTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(auth.ContextKeyRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperr.ErrForbidden.Message})
			return
		}
		role, ok := roleVal.(auth.Role)
		if !ok || role != auth.RoleSuperAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperr.ErrForbidden.Message})
			return
		}
		c.Next()
	}
}

// RequireAdminTOTPSetup: admin 역할 확인 + TOTP 설정 여부 확인.
// TOTP 미설정 시 401 + redirect 반환.
func RequireAdminTOTPSetup(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := auth.GetCurrentRole(c)
		if !ok || role != auth.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperr.ErrForbidden.Message})
			return
		}

		userID, ok := auth.GetCurrentUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperr.ErrUnauthorized.Message})
			return
		}

		if !IsTOTPEnabled(c.Request.Context(), userID, db) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "TOTP setup required",
				"redirect": "/admin/totp-setup",
			})
			return
		}

		c.Next()
	}
}

// RequireTOTPValidated: X-TOTP-Code 헤더의 6자리 코드 검증.
func RequireTOTPValidated(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.GetHeader("X-TOTP-Code")
		if code == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "TOTP code required"})
			return
		}

		userID, ok := auth.GetCurrentUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperr.ErrUnauthorized.Message})
			return
		}

		var encryptedSecret string
		err := db.QueryRow(c.Request.Context(),
			"SELECT totp_secret FROM users WHERE id = $1", userID,
		).Scan(&encryptedSecret)
		if err != nil || encryptedSecret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
			return
		}

		if !ValidateTOTP(encryptedSecret, code) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
			return
		}

		c.Next()
	}
}

func (h *AdminHandler) requireAdminTOTPSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := auth.GetCurrentUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperr.ErrUnauthorized.Message})
			return
		}
		role, ok := auth.GetCurrentRole(c)
		if !ok || role != auth.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperr.ErrForbidden.Message})
			return
		}

		token, err := c.Cookie(adminTOTPSessionCookieName)
		if err != nil || strings.TrimSpace(token) == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin_totp_required"})
			return
		}

		claims, err := h.authSvc.ParseToken(token)
		if err != nil || claims.UserID != userID || claims.Role != auth.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin_totp_required"})
			return
		}

		c.Next()
	}
}

type loginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	TOTP     string `json:"totp"`
}

func SuperAdminAuth(cfg *config.Config, authSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if req.ID != cfg.SuperAdminID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(cfg.SuperAdminPWHash), []byte(req.Password)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		valid := totp.Validate(req.TOTP, cfg.SuperAdminTOTPSecret)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid totp"})
			return
		}

		token, err := authSvc.IssueSuperAdminToken(auth.User{
			ID:        "super_admin",
			Role:      auth.RoleSuperAdmin,
			CreatedAt: time.Now(),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
			return
		}

		sessionCookieSecure := strings.EqualFold(cfg.Env, "production")
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			superAdminAccessTokenCookieName,
			token,
			int((4 * time.Hour).Seconds()),
			"/",
			"",
			sessionCookieSecure,
			true,
		)

		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
