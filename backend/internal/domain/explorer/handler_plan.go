package explorer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/safehttp"
)

func (h *Handler) CheckExplorationNodeLink(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = userID

	var payload struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawURL := strings.TrimSpace(payload.URL)
	if !(strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://")) {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "http:// 또는 https:// 로 시작하는 URL을 입력해주세요."})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 7*time.Second)
	defer cancel()

	client := safehttp.NewClient(safehttp.ClientConfig{Timeout: 7 * time.Second, MaxRedirects: 5})

	statusCode, err := checkURL(ctx, client, http.MethodHead, rawURL)
	if err != nil || statusCode >= 400 {
		statusCode, err = checkURL(ctx, client, http.MethodGet, rawURL)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "링크에 접근할 수 없습니다."})
		return
	}
	if statusCode >= 200 && statusCode < 400 {
		c.JSON(http.StatusOK, gin.H{"valid": true, "url": rawURL, "message": "유효한 링크입니다."})
	} else {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": fmt.Sprintf("링크가 응답하지 않습니다. (상태: %d)", statusCode)})
	}
}

func checkURL(ctx context.Context, client *http.Client, method, rawURL string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 LearnWeaver-LinkChecker/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 512))
	resp.Body.Close()
	return resp.StatusCode, nil
}

func (h *Handler) SavePlan(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SavePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !h.verifyCourseOwner(c, req.CourseDraftID, userID) {
		return
	}

	if err := h.repo.SaveRegionOrder(c.Request.Context(), req.CourseDraftID, req.RegionOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
