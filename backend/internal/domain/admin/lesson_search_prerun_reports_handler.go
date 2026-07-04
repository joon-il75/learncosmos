package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
)

const (
	lessonSearchPrerunReportsDefaultLimit = 10
	lessonSearchPrerunReportsMaxLimit     = 50
)

type lessonSearchPrerunReportLister interface {
	ListRecentReports(ctx context.Context, limit int) ([]lessonsearchprerun.StoredReport, error)
}

type lessonSearchPrerunReportResponse struct {
	ID                   string          `json:"id"`
	RunType              string          `json:"run_type"`
	Status               string          `json:"status"`
	StartedAt            string          `json:"started_at"`
	FinishedAt           *string         `json:"finished_at,omitempty"`
	ProviderFilter       string          `json:"provider_filter"`
	TopN                 int             `json:"top_n"`
	RiskTopN             int             `json:"risk_top_n"`
	Summary              json.RawMessage `json:"summary"`
	NoiseTaxonomySummary json.RawMessage `json:"noise_taxonomy_summary"`
	Report               json.RawMessage `json:"report,omitempty"`
	ErrorCode            *string         `json:"error_code,omitempty"`
	ErrorMessage         *string         `json:"error_message,omitempty"`
	CreatedAt            string          `json:"created_at"`
}

// GET /api/v1/super-admin/recommendation-debug/lesson-search-prerun-reports
func (h *AdminHandler) ListLessonSearchPrerunReports(c *gin.Context) {
	limit, err := parseLessonSearchPrerunReportsLimit(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	includeReport := parseLessonSearchPrerunIncludeReport(c.Query("include_report"))

	if h.lessonSearchPrerunReports == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lesson search prerun reports are not configured"})
		return
	}
	reports, err := h.lessonSearchPrerunReports.ListRecentReports(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load lesson search prerun reports"})
		return
	}

	items := make([]lessonSearchPrerunReportResponse, 0, len(reports))
	for _, report := range reports {
		if err := lessonsearchprerun.EnsureNoSecretMarkers(report.Summary, report.NoiseTaxonomySummary, report.ErrorCode, report.ErrorMessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lesson search prerun report failed safety check"})
			return
		}
		item := lessonSearchPrerunReportResponse{
			ID:                   report.ID.String(),
			RunType:              report.RunType,
			Status:               report.Status,
			StartedAt:            report.StartedAt.Format(time.RFC3339Nano),
			ProviderFilter:       report.ProviderFilter,
			TopN:                 report.TopN,
			RiskTopN:             report.RiskTopN,
			Summary:              report.Summary,
			NoiseTaxonomySummary: report.NoiseTaxonomySummary,
			ErrorCode:            report.ErrorCode,
			ErrorMessage:         report.ErrorMessage,
			CreatedAt:            report.CreatedAt.Format(time.RFC3339Nano),
		}
		if report.FinishedAt != nil {
			finishedAt := report.FinishedAt.Format(time.RFC3339Nano)
			item.FinishedAt = &finishedAt
		}
		if includeReport {
			if err := lessonsearchprerun.EnsureNoSecretMarkers(report.Report); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "lesson search prerun report failed safety check"})
				return
			}
			item.Report = report.Report
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":          limit,
		"include_report": includeReport,
		"reports":        items,
	})
}

func parseLessonSearchPrerunReportsLimit(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return lessonSearchPrerunReportsDefaultLimit, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit <= 0 {
		return 0, errLessonSearchPrerunInvalidLimit
	}
	if limit > lessonSearchPrerunReportsMaxLimit {
		return lessonSearchPrerunReportsMaxLimit, nil
	}
	return limit, nil
}

func parseLessonSearchPrerunIncludeReport(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

var errLessonSearchPrerunInvalidLimit = lessonSearchPrerunRequestError("limit must be a positive integer")

type lessonSearchPrerunRequestError string

func (err lessonSearchPrerunRequestError) Error() string {
	return string(err)
}
