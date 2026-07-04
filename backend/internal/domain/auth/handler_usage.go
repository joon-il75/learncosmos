package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

func (h *Handler) getMyAIUsageHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate
	if rawStart := strings.TrimSpace(c.Query("start_date")); rawStart != "" {
		parsed, err := time.ParseInLocation("2006-01-02", rawStart, now.Location())
		if err != nil {
			c.JSON(http.StatusBadRequest, apiError("ai_usage_invalid_start_date", "invalid_start_date"))
			return
		}
		startDate = parsed
	}
	if rawEnd := strings.TrimSpace(c.Query("end_date")); rawEnd != "" {
		parsed, err := time.ParseInLocation("2006-01-02", rawEnd, now.Location())
		if err != nil {
			c.JSON(http.StatusBadRequest, apiError("ai_usage_invalid_end_date", "invalid_end_date"))
			return
		}
		endDate = parsed
	}
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, apiError("ai_usage_invalid_date_range", "invalid_date_range"))
		return
	}
	maxEndExclusive := startDate.AddDate(0, 3, 0)
	endExclusive := endDate.AddDate(0, 0, 1)
	if endExclusive.After(maxEndExclusive) {
		c.JSON(http.StatusBadRequest, apiError("ai_usage_date_range_too_large", "date_range_too_large"))
		return
	}

	page := 1
	if rawPage := strings.TrimSpace(c.Query("page")); rawPage != "" {
		if parsed, err := strconv.Atoi(rawPage); err == nil && parsed > 0 {
			page = parsed
		}
	}
	limit := 10
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		if parsed, err := strconv.Atoi(rawLimit); err == nil && parsed > 0 && parsed <= 10 {
			limit = parsed
		}
	}
	offset := (page - 1) * limit

	events, total, summary, err := h.svc.repo.ListUserAIUsageEvents(c.Request.Context(), userID, startDate, endExclusive, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("ai_usage_load_failed", "AI 사용량 조회 실패"))
		return
	}
	rangeDays := int(endDate.Sub(startDate).Hours()/24) + 1
	if rangeDays < 1 {
		rangeDays = 1
	}
	summary.AverageDailyTokens = float64(summary.TotalTokens) / float64(rangeDays)
	summary.AverageDailyEstimatedCostUSD = summary.TotalEstimatedCostUSD / float64(rangeDays)
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}
	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"start_date":  startDate.Format("2006-01-02"),
		"end_date":    endDate.Format("2006-01-02"),
		"summary":     summary,
	})
}

func (h *Handler) getMyPointUsageHandler(c *gin.Context) {
	userID, ok := GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, apiError("unauthorized", apperr.ErrUnauthorized.Message))
		return
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate
	if rawStart := strings.TrimSpace(c.Query("start_date")); rawStart != "" {
		parsed, err := time.ParseInLocation("2006-01-02", rawStart, now.Location())
		if err != nil {
			c.JSON(http.StatusBadRequest, apiError("point_usage_invalid_start_date", "invalid_start_date"))
			return
		}
		startDate = parsed
	}
	if rawEnd := strings.TrimSpace(c.Query("end_date")); rawEnd != "" {
		parsed, err := time.ParseInLocation("2006-01-02", rawEnd, now.Location())
		if err != nil {
			c.JSON(http.StatusBadRequest, apiError("point_usage_invalid_end_date", "invalid_end_date"))
			return
		}
		endDate = parsed
	}
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, apiError("point_usage_invalid_date_range", "invalid_date_range"))
		return
	}
	maxEndExclusive := startDate.AddDate(0, 3, 0)
	endExclusive := endDate.AddDate(0, 0, 1)
	if endExclusive.After(maxEndExclusive) {
		c.JSON(http.StatusBadRequest, apiError("point_usage_date_range_too_large", "date_range_too_large"))
		return
	}

	page := 1
	if rawPage := strings.TrimSpace(c.Query("page")); rawPage != "" {
		if parsed, err := strconv.Atoi(rawPage); err == nil && parsed > 0 {
			page = parsed
		}
	}
	limit := 10
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		if parsed, err := strconv.Atoi(rawLimit); err == nil && parsed > 0 && parsed <= 10 {
			limit = parsed
		}
	}
	offset := (page - 1) * limit

	transactions, total, summary, err := h.svc.repo.ListUserPointTransactions(c.Request.Context(), userID, startDate, endExclusive, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("point_usage_load_failed", "포인트 사용량 조회 실패"))
		return
	}
	freePoints, paidPoints, err := h.svc.repo.GetPointBalances(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiError("points_load_failed", "failed to fetch user points"))
		return
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"total_pages":  totalPages,
		"start_date":   startDate.Format("2006-01-02"),
		"end_date":     endDate.Format("2006-01-02"),
		"summary":      summary,
		"balance": gin.H{
			"free_points":  freePoints,
			"paid_points":  paidPoints,
			"total_points": freePoints + paidPoints,
		},
	})
}
