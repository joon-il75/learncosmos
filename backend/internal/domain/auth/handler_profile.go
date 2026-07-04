package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

type UpdateProfileRequest struct {
	Nickname *string `json:"nickname"`
	Email    *string `json:"email"`
}

func (h *Handler) getUserAvatarHandler(c *gin.Context) {
	filename := filepath.Base(strings.TrimSpace(c.Param("filename")))
	if filename == "." || filename == "" || strings.Contains(filename, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid avatar filename"})
		return
	}
	c.File(filepath.Join(profileAvatarDir, filename))
}

func (h *Handler) patchMeHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiError("invalid_request", "잘못된 요청입니다"))
		return
	}

	if req.Nickname != nil {
		nickname := strings.TrimSpace(*req.Nickname)
		if nickname == "" {
			c.JSON(http.StatusBadRequest, apiError("profile_nickname_required", "닉네임은 비워둘 수 없습니다"))
			return
		}
		if utf8.RuneCountInString(nickname) > 20 {
			c.JSON(http.StatusBadRequest, apiError("profile_nickname_too_long", "닉네임은 20자 이하여야 합니다"))
			return
		}
		if err := h.svc.repo.UpdateNickname(c.Request.Context(), userID, nickname); err != nil {
			c.JSON(http.StatusInternalServerError, apiError("profile_update_failed", "닉네임 업데이트 실패"))
			return
		}
	}

	if req.Email != nil {
		user, err := h.svc.repo.FindByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, apiError("user_not_found", "사용자를 찾을 수 없습니다"))
			return
		}
		if !(user.Provider == "kakao" && user.Email == "") {
			c.JSON(http.StatusForbidden, apiError("social_email_locked", "소셜 로그인 연동 이메일은 변경할 수 없습니다"))
			return
		}
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			c.JSON(http.StatusBadRequest, apiError("email_required", "이메일을 입력해주세요"))
			return
		}
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			c.JSON(http.StatusBadRequest, apiError("email_invalid", "올바른 이메일 형식이 아닙니다"))
			return
		}
		exists, _ := h.svc.repo.EmailExists(c.Request.Context(), email, userID)
		if exists {
			c.JSON(http.StatusConflict, apiError("email_in_use", "이미 사용 중인 이메일입니다"))
			return
		}
		if err := h.svc.repo.UpdateEmail(c.Request.Context(), userID, email); err != nil {
			c.JSON(http.StatusInternalServerError, apiError("email_update_failed", "이메일 업데이트 실패"))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "프로필이 업데이트되었습니다"})
}

func (h *Handler) uploadMyAvatarHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, profileAvatarMaxBytes+1024)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError("avatar_file_required", "프로필 이미지 파일이 필요합니다"))
		return
	}
	if fileHeader.Size > profileAvatarMaxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, apiError("avatar_file_too_large", "프로필 이미지는 1MB 이하여야 합니다"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, apiError("avatar_open_failed", "프로필 이미지를 열 수 없습니다"))
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(io.LimitReader(file, profileAvatarMaxBytes+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_read_failed", "프로필 이미지를 읽지 못했습니다"))
		return
	}
	if int64(len(payload)) > profileAvatarMaxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, apiError("avatar_file_too_large", "프로필 이미지는 1MB 이하여야 합니다"))
		return
	}

	contentType := http.DetectContentType(payload)
	ext, requiresDimensionCheck, ok := normalizeProfileAvatarContentType(contentType)
	if !ok {
		c.JSON(http.StatusBadRequest, apiError("avatar_type_unsupported", "프로필 이미지는 JPG, PNG, WebP만 사용할 수 있습니다"))
		return
	}
	if requiresDimensionCheck {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(payload))
		if err != nil {
			c.JSON(http.StatusBadRequest, apiError("avatar_dimension_unreadable", "프로필 이미지 크기를 확인할 수 없습니다"))
			return
		}
		if cfg.Width > profileAvatarMaxSide || cfg.Height > profileAvatarMaxSide {
			c.JSON(http.StatusBadRequest, apiError("avatar_dimension_too_large", "프로필 이미지는 가로/세로 1024px 이하여야 합니다"))
			return
		}
	}

	if err := os.MkdirAll(profileAvatarDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_storage_failed", "프로필 이미지 저장소를 준비하지 못했습니다"))
		return
	}
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_filename_failed", "프로필 이미지 이름을 만들지 못했습니다"))
		return
	}
	filename := userID + "-" + time.Now().UTC().Format("20060102150405") + "-" + hex.EncodeToString(random) + ext
	target := filepath.Join(profileAvatarDir, filename)
	tmpPath := target + ".uploading"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_write_failed", "프로필 이미지를 저장하지 못했습니다"))
		return
	}
	if err := os.Rename(tmpPath, target); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_apply_failed", "프로필 이미지를 반영하지 못했습니다"))
		return
	}

	avatarURL := profileAvatarURLBase + "/" + filename
	if err := h.svc.repo.UpdateAvatarURL(c.Request.Context(), userID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("avatar_metadata_failed", "프로필 이미지 정보를 저장하지 못했습니다"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "프로필 이미지가 저장되었습니다",
		"avatar_url": avatarURL,
	})
}

func normalizeProfileAvatarContentType(contentType string) (string, bool, bool) {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg":
		return ".jpg", true, true
	case "image/png":
		return ".png", true, true
	case "image/webp":
		return ".webp", false, true
	default:
		return "", false, false
	}
}

func (h *Handler) withdrawMeHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	if err := h.svc.WithdrawUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, apiError("withdrawal_failed", "계정 탈퇴 처리에 실패했습니다"))
		return
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "계정이 탈퇴 처리되었습니다. 같은 소셜 계정으로 다시 가입하면 보관 기간 내 학습 기록을 다시 확인할 수 있습니다.",
	})
}
