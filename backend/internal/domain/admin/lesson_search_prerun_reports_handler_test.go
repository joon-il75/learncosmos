package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
)

type fakeLessonSearchPrerunReportLister struct {
	limit   int
	reports []lessonsearchprerun.StoredReport
}

func (lister *fakeLessonSearchPrerunReportLister) ListRecentReports(ctx context.Context, limit int) ([]lessonsearchprerun.StoredReport, error) {
	lister.limit = limit
	return lister.reports, nil
}

func TestListLessonSearchPrerunReportsHidesReportByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeLessonSearchPrerunReportLister{reports: []lessonsearchprerun.StoredReport{testLessonSearchPrerunStoredReport(t)}}
	handler := &AdminHandler{lessonSearchPrerunReports: store}
	router := gin.New()
	router.GET("/reports", handler.ListLessonSearchPrerunReports)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reports?limit=2", nil)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if store.limit != 2 {
		t.Fatalf("limit = %d, want 2", store.limit)
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body["include_report"] != false {
		t.Fatalf("include_report = %v, want false", body["include_report"])
	}
	reports := body["reports"].([]any)
	item := reports[0].(map[string]any)
	if _, ok := item["report"]; ok {
		t.Fatal("report should be hidden by default")
	}
	if item["status"] != lessonsearchprerun.ReportStatusSucceeded {
		t.Fatalf("status = %v", item["status"])
	}
}

func TestListLessonSearchPrerunReportsIncludesReportWhenRequested(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeLessonSearchPrerunReportLister{reports: []lessonsearchprerun.StoredReport{testLessonSearchPrerunStoredReport(t)}}
	handler := &AdminHandler{lessonSearchPrerunReports: store}
	router := gin.New()
	router.GET("/reports", handler.ListLessonSearchPrerunReports)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reports?include_report=true&limit=999", nil)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if store.limit != lessonSearchPrerunReportsMaxLimit {
		t.Fatalf("limit = %d, want max limit", store.limit)
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	reports := body["reports"].([]any)
	item := reports[0].(map[string]any)
	if _, ok := item["report"]; !ok {
		t.Fatal("report should be included when requested")
	}
}

func TestParseLessonSearchPrerunReportsLimitRejectsInvalid(t *testing.T) {
	if _, err := parseLessonSearchPrerunReportsLimit("abc"); err == nil {
		t.Fatal("expected invalid limit error")
	}
	if got, err := parseLessonSearchPrerunReportsLimit(""); err != nil || got != lessonSearchPrerunReportsDefaultLimit {
		t.Fatalf("default limit = %d, err = %v", got, err)
	}
}

func testLessonSearchPrerunStoredReport(t *testing.T) lessonsearchprerun.StoredReport {
	t.Helper()
	startedAt := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Minute)
	return lessonsearchprerun.StoredReport{
		ID:                   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		RunType:              "daily",
		Status:               lessonsearchprerun.ReportStatusSucceeded,
		StartedAt:            startedAt,
		FinishedAt:           &finishedAt,
		ProviderFilter:       "all",
		TopN:                 5,
		RiskTopN:             10,
		Summary:              json.RawMessage(`{"provider_errors":0,"estimated_provider_requests":2}`),
		NoiseTaxonomySummary: json.RawMessage(`{"registered_avoid_hits":[],"unregistered_noise_hints":[]}`),
		Report:               json.RawMessage(`{"mode":"refresh-existing"}`),
		CreatedAt:            finishedAt,
	}
}
