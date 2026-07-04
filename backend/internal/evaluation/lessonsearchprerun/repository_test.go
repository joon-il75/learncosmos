package lessonsearchprerun

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestEnsureNoSecretMarkersRejectsKnownMarkers(t *testing.T) {
	for _, value := range []any{
		map[string]any{"error": "YOUTUBE_API_KEY missing"},
		`{"header":"Authorization Bearer token"}`,
		[]byte(`{"key":"sk-test123"}`),
	} {
		if err := EnsureNoSecretMarkers(value); err == nil {
			t.Fatalf("EnsureNoSecretMarkers(%#v) = nil, want error", value)
		}
	}
}

func TestEnsureNoSecretMarkersAllowsNormalReport(t *testing.T) {
	report := map[string]any{
		"summary": map[string]any{
			"provider_errors": 0,
		},
		"noise_taxonomy_summary": map[string]any{
			"registered_avoid_hits": []map[string]any{{"token": "키트", "count": 1}},
		},
	}
	if err := EnsureNoSecretMarkers(report); err != nil {
		t.Fatalf("EnsureNoSecretMarkers() error = %v", err)
	}
}

func TestNormalizeSaveReportInputAppliesDefaults(t *testing.T) {
	got, err := normalizeSaveReportInput(SaveReportInput{})
	if err != nil {
		t.Fatalf("normalizeSaveReportInput() error = %v", err)
	}
	if got.RunType != "daily" {
		t.Fatalf("run type = %q, want daily", got.RunType)
	}
	if got.Status != ReportStatusSucceeded {
		t.Fatalf("status = %q, want %q", got.Status, ReportStatusSucceeded)
	}
	if got.ProviderFilter != "all" {
		t.Fatalf("provider filter = %q, want all", got.ProviderFilter)
	}
	if got.TopN != 5 {
		t.Fatalf("top_n = %d, want 5", got.TopN)
	}
	if got.Summary == nil || got.NoiseTaxonomySummary == nil || got.Report == nil {
		t.Fatalf("json object defaults were not set: %#v", got)
	}
}

func TestNormalizeSaveReportInputRejectsInvalidStatusAndRiskTopN(t *testing.T) {
	if _, err := normalizeSaveReportInput(SaveReportInput{Status: "done"}); err == nil {
		t.Fatal("expected invalid status error")
	}
	if _, err := normalizeSaveReportInput(SaveReportInput{RiskTopN: -1}); err == nil {
		t.Fatal("expected invalid risk_top_n error")
	}
}

func TestMarshalReportJSONRequiresObject(t *testing.T) {
	if _, err := marshalReportJSON([]string{"not-object"}, "report"); err == nil {
		t.Fatal("expected non-object JSON error")
	}
	if raw, err := marshalReportJSON(map[string]any{"ok": true}, "report"); err != nil {
		t.Fatalf("marshalReportJSON() error = %v", err)
	} else if !json.Valid(raw) {
		t.Fatalf("marshalReportJSON() produced invalid JSON: %s", raw)
	}
}

func TestReportRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is required for lesson search prerun report repository integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	defer pool.Close()

	repo := NewReportRepository(pool)
	startedAt := time.Now().UTC().Add(-time.Minute)
	finishedAt := startedAt.Add(10 * time.Second)
	stored, err := repo.SaveReport(ctx, SaveReportInput{
		RunType:        "test",
		Status:         ReportStatusSucceeded,
		StartedAt:      startedAt,
		FinishedAt:     &finishedAt,
		ProviderFilter: "all",
		TopN:           5,
		RiskTopN:       10,
		Summary: map[string]any{
			"provider_errors": 0,
			"gap_count":       0,
		},
		NoiseTaxonomySummary: map[string]any{
			"registered_avoid_hits": []map[string]any{{"token": "키트", "count": 1}},
		},
		Report: map[string]any{
			"mode": "refresh-existing",
		},
	})
	if err != nil {
		t.Fatalf("SaveReport() error = %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM lesson_search_prerun_reports WHERE id = $1`, stored.ID)
	})
	if stored.ID.String() == "" {
		t.Fatal("stored report id is empty")
	}
	if stored.Status != ReportStatusSucceeded {
		t.Fatalf("stored status = %q", stored.Status)
	}

	recent, err := repo.ListRecentReports(ctx, 5)
	if err != nil {
		t.Fatalf("ListRecentReports() error = %v", err)
	}
	found := false
	for _, item := range recent {
		if item.ID == stored.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("stored report %s not found in recent reports", stored.ID)
	}

	deleted, err := repo.DeleteReportsCreatedBefore(ctx, time.Now().UTC().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("DeleteReportsCreatedBefore() error = %v", err)
	}
	if deleted < 0 {
		t.Fatalf("deleted rows = %d, want non-negative", deleted)
	}
}
