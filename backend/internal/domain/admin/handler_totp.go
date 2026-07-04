package admin

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/auth"
)

// totpSetupHandler godoc
// POST /api/v1/admin/totp/setup
// JWT(admin) 필수. TOTP QR URL과 1회용 평문 시크릿 반환. DB에 암호화된 시크릿 저장.
func (h *AdminHandler) totpSetupHandler(c *gin.Context) {
	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// DB에서 이메일 조회
	var email string
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT email FROM users WHERE id = $1", userID,
	).Scan(&email)
	if err != nil {
		log.Printf("[admin] totpSetup: failed to fetch email for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	encryptedSecret, plaintextSecret, qrURL, err := GenerateTOTPSecret(email)
	if err != nil {
		log.Printf("[admin] totpSetup: GenerateTOTPSecret error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// totp_secret 저장 (totp_enabled는 아직 false)
	_, err = h.db.Exec(c.Request.Context(),
		"UPDATE users SET totp_secret = $1 WHERE id = $2",
		encryptedSecret, userID,
	)
	if err != nil {
		log.Printf("[admin] totpSetup: failed to save totp_secret for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_url": qrURL,
		"secret": plaintextSecret, // 1회만 노출
	})
}

// totpVerifyHandler godoc
// POST /api/v1/admin/totp/verify
// QR 스캔 후 첫 6자리 코드 입력으로 TOTP 설정 완료.
func (h *AdminHandler) totpVerifyHandler(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var encryptedSecret string
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT totp_secret FROM users WHERE id = $1", userID,
	).Scan(&encryptedSecret)
	if err != nil || encryptedSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "totp not initialized"})
		return
	}

	if !ValidateTOTP(encryptedSecret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	_, err = h.db.Exec(c.Request.Context(),
		"UPDATE users SET totp_enabled = true WHERE id = $1", userID,
	)
	if err != nil {
		log.Printf("[admin] totpVerify: failed to enable totp for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if err := h.setAdminTOTPSessionCookie(c, userID); err != nil {
		log.Printf("[admin] totpVerify: failed to issue totp session for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// totpValidateHandler godoc
// POST /api/v1/admin/totp/validate
// 매 로그인 시 TOTP 코드 검증.
func (h *AdminHandler) totpValidateHandler(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var encryptedSecret string
	var enabled bool
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT COALESCE(totp_secret,''), totp_enabled FROM users WHERE id = $1", userID,
	).Scan(&encryptedSecret, &enabled)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}
	if !enabled || encryptedSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "totp_not_initialized"})
		return
	}

	if !ValidateTOTP(encryptedSecret, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	if err := h.setAdminTOTPSessionCookie(c, userID); err != nil {
		log.Printf("[admin] totpValidate: failed to issue totp session for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func (h *AdminHandler) setAdminTOTPSessionCookie(c *gin.Context, userID string) error {
	token, err := h.authSvc.IssueAdminTOTPToken(userID)
	if err != nil {
		return err
	}

	secure := strings.EqualFold(h.cfg.Env, "production")
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		adminTOTPSessionCookieName,
		token,
		int((4 * time.Hour).Seconds()),
		"/",
		"",
		secure,
		true,
	)
	return nil
}

func (h *AdminHandler) patchRoleHandler(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if body.Role != string(auth.RoleLearner) && body.Role != string(auth.RoleAdmin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	_, err := h.db.Exec(c.Request.Context(),
		"UPDATE users SET role=$1 WHERE id=$2",
		body.Role, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"user_id":  id,
		"new_role": body.Role,
	})
}

// totpStatusHandler godoc
// GET /api/v1/admin/totp/status
// 관리자의 TOTP 설정 상태를 반환합니다.
func (h *AdminHandler) totpStatusHandler(c *gin.Context) {
	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var enabled bool
	var resetRequested bool
	err := h.db.QueryRow(c.Request.Context(),
		"SELECT totp_enabled, totp_reset_requested FROM users WHERE id = $1", userID,
	).Scan(&enabled, &resetRequested)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totp_enabled":    enabled,
		"reset_requested": resetRequested,
	})
}

// totpResetRequestHandler godoc
// POST /api/v1/admin/totp/reset-request
// 관리자가 기기 분실 등의 이유로 슈퍼관리자에게 TOTP 초기화를 신청합니다.
func (h *AdminHandler) totpResetRequestHandler(c *gin.Context) {
	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	_, err := h.db.Exec(c.Request.Context(),
		"UPDATE users SET totp_reset_requested = true WHERE id = $1", userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "신청 처리 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "초기화 신청이 완료되었습니다. 슈퍼관리자 확인 후 처리됩니다."})
}

// totpResetHandler godoc
// POST /api/v1/super-admin/users/:id/totp-reset
// 슈퍼관리자가 관리자의 TOTP를 초기화합니다.
func (h *AdminHandler) totpResetHandler(c *gin.Context) {
	targetID := c.Param("id")

	_, err := h.db.Exec(c.Request.Context(),
		`UPDATE users
		 SET totp_secret = NULL, totp_enabled = false, totp_reset_requested = false
		 WHERE id = $1 AND role = 'admin'`,
		targetID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "초기화 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TOTP가 초기화되었습니다."})
}
