package alphaaccess

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func apiError(errorCode, message string) gin.H {
	return gin.H{"error": message, "error_code": errorCode}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc) {
	protected := r.Group("/alpha-access")
	protected.Use(jwtMiddleware)
	{
		protected.GET("/status", h.statusHandler)
		protected.POST("/redeem", consentMiddleware, h.redeemHandler)
	}
}

func (h *Handler) RegisterSuperAdminRoutes(r *gin.RouterGroup) {
	r.GET("/alpha-invite-codes", h.listInviteCodesHandler)
	r.POST("/alpha-invite-codes", h.createInviteCodeHandler)
	r.PATCH("/alpha-invite-codes/:id/revoke", h.revokeInviteCodeHandler)
}

func (h *Handler) RequireAlphaAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := auth.GetCurrentUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
			return
		}

		role, _ := auth.GetCurrentRole(c)
		if role != auth.RoleLearner {
			c.Next()
			return
		}

		granted, _, err := h.repo.HasAlphaAccess(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, apiError("alpha_status_load_failed", "failed to fetch alpha access status"))
			return
		}
		if !granted {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":               "alpha_access_required",
				"error_code":          "alpha_access_required",
				"alpha_access_needed": true,
			})
			return
		}
		c.Next()
	}
}

func (h *Handler) statusHandler(c *gin.Context) {
	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	role, _ := auth.GetCurrentRole(c)
	if role != auth.RoleLearner {
		c.JSON(http.StatusOK, StatusResponse{
			Granted:      true,
			RequiresCode: false,
			Role:         string(role),
		})
		return
	}

	granted, grantedAt, err := h.repo.HasAlphaAccess(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("alpha_status_load_failed", "failed to fetch alpha access status"))
		return
	}
	var grantedAtText *string
	if grantedAt != nil {
		text := grantedAt.UTC().Format(time.RFC3339)
		grantedAtText = &text
	}
	c.JSON(http.StatusOK, StatusResponse{
		Granted:      granted,
		GrantedAt:    grantedAtText,
		RequiresCode: !granted,
		Role:         string(role),
	})
}

func (h *Handler) redeemHandler(c *gin.Context) {
	userID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}
	grantedAt, err := h.repo.Redeem(c.Request.Context(), userID, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, ErrCodeNotFound):
			c.JSON(http.StatusNotFound, apiError("alpha_code_not_found", "초대 코드를 찾을 수 없습니다"))
		case errors.Is(err, ErrCodeRevoked):
			c.JSON(http.StatusBadRequest, apiError("alpha_code_revoked", "회수된 초대 코드입니다"))
		case errors.Is(err, ErrCodeExpired):
			c.JSON(http.StatusBadRequest, apiError("alpha_code_expired", "입력 기간이 지난 초대 코드입니다"))
		case errors.Is(err, ErrCodeUsedUp):
			c.JSON(http.StatusBadRequest, apiError("alpha_code_used_up", "이미 사용이 완료된 초대 코드입니다"))
		default:
			c.JSON(http.StatusInternalServerError, apiError("alpha_redeem_failed", "초대 코드 처리에 실패했습니다"))
		}
		return
	}
	c.JSON(http.StatusOK, RedeemResponse{
		Granted:   true,
		GrantedAt: grantedAt.UTC().Format(time.RFC3339),
		Message:   "알파 테스트 접근권이 활성화되었습니다",
	})
}

func (h *Handler) listInviteCodesHandler(c *gin.Context) {
	items, err := h.repo.ListInviteCodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invite codes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) createInviteCodeHandler(c *gin.Context) {
	var req CreateInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}
	expiresAt, err := parseExpiresAt(req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "만료일을 확인해 주세요"})
		return
	}
	if !expiresAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "만료일은 현재 시각 이후여야 합니다"})
		return
	}
	if req.MaxUses < 1 || req.MaxUses > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "사용 가능 횟수는 1~100 사이로 입력해 주세요"})
		return
	}

	var createdBy *string
	if userID, ok := auth.GetCurrentUserID(c); ok {
		if _, err := uuid.Parse(userID); err == nil {
			createdBy = &userID
		}
	}

	item, err := h.repo.CreateInviteCode(c.Request.Context(), createdBy, expiresAt, req.MaxUses, req.SentToNote, req.AdminNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "초대 코드 발급에 실패했습니다"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) revokeInviteCodeHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 코드 ID입니다"})
		return
	}
	if err := h.repo.RevokeInviteCode(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "초대 코드 회수에 실패했습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func parseExpiresAt(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, errors.New("empty expires_at")
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	date, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return date.Add(24*time.Hour - time.Nanosecond), nil
}
