package jobs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
)

type fakePrerunReportStore struct {
	inputs       []lessonsearchprerun.SaveReportInput
	deleteCutoff *time.Time
}

func (store *fakePrerunReportStore) SaveReport(ctx context.Context, input lessonsearchprerun.SaveReportInput) (*lessonsearchprerun.StoredReport, error) {
	store.inputs = append(store.inputs, input)
	return &lessonsearchprerun.StoredReport{Status: input.Status}, nil
}

func (store *fakePrerunReportStore) DeleteReportsCreatedBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	store.deleteCutoff = &cutoff
	return 0, nil
}

type fakePrerunLocker struct {
	acquired     bool
	unlockCalled bool
}

func (locker *fakePrerunLocker) TryLock(ctx context.Context) (bool, func(context.Context) error, error) {
	if !locker.acquired {
		return false, nil, nil
	}
	return true, func(context.Context) error {
		locker.unlockCalled = true
		return nil
	}, nil
}

type fakePrerunProvider struct {
	called bool
	calls  []int
}

func (provider *fakePrerunProvider) Search(ctx context.Context, fixture lessonsearchprerun.Fixture, topN int) ([]lessonsearchprerun.Candidate, error) {
	provider.called = true
	provider.calls = append(provider.calls, topN)
	return []lessonsearchprerun.Candidate{{Title: "기초 강좌", Description: "입문 튜토리얼", ContentType: "video"}}, nil
}

func TestLessonSearchPrerunSchedulerRunOnceSkipsWhenLockNotAcquired(t *testing.T) {
	store := &fakePrerunReportStore{}
	scheduler := NewLessonSearchPrerunScheduler(LessonSearchPrerunSchedulerConfig{Provider: "youtube", YouTubeAPIKey: "key"}, store, &fakePrerunLocker{}, nil)

	if err := scheduler.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if len(store.inputs) != 1 {
		t.Fatalf("saved reports = %d, want 1", len(store.inputs))
	}
	if got := store.inputs[0].Status; got != lessonsearchprerun.ReportStatusSkippedLockNotAcquired {
		t.Fatalf("status = %q, want lock skip", got)
	}
}

func TestLessonSearchPrerunSchedulerRunOnceSkipsMissingCredentials(t *testing.T) {
	store := &fakePrerunReportStore{}
	provider := &fakePrerunProvider{}
	scheduler := NewLessonSearchPrerunScheduler(LessonSearchPrerunSchedulerConfig{Provider: "youtube"}, store, &fakePrerunLocker{acquired: true}, nil)
	scheduler.providers["youtube"] = provider

	if err := scheduler.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if provider.called {
		t.Fatal("provider was called despite missing credentials")
	}
	if got := store.inputs[0].Status; got != lessonsearchprerun.ReportStatusSkippedMissingCredentials {
		t.Fatalf("status = %q, want missing credentials skip", got)
	}
}

func TestLessonSearchPrerunSchedulerRunOnceSavesSucceededReport(t *testing.T) {
	fixturePath := writePrerunFixture(t)
	store := &fakePrerunReportStore{}
	locker := &fakePrerunLocker{acquired: true}
	provider := &fakePrerunProvider{}
	scheduler := NewLessonSearchPrerunScheduler(LessonSearchPrerunSchedulerConfig{
		Provider:       "youtube",
		YouTubeAPIKey:  "key",
		FixturesPath:   fixturePath,
		TopN:           5,
		RiskTopN:       10,
		RiskFixtureIDs: []string{"youtube-setup"},
		Retention:      24 * time.Hour,
	}, store, locker, nil)
	scheduler.providers["youtube"] = provider

	if err := scheduler.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if !provider.called {
		t.Fatal("provider was not called")
	}
	if len(provider.calls) != 2 {
		t.Fatalf("provider calls = %d, want 2", len(provider.calls))
	}
	if provider.calls[0] != 5 || provider.calls[1] != 10 {
		t.Fatalf("provider top_n calls = %v, want [5 10]", provider.calls)
	}
	if !locker.unlockCalled {
		t.Fatal("advisory lock was not released")
	}
	if len(store.inputs) != 1 {
		t.Fatalf("saved reports = %d, want 1", len(store.inputs))
	}
	input := store.inputs[0]
	if input.Status != lessonsearchprerun.ReportStatusSucceeded {
		t.Fatalf("status = %q, want succeeded", input.Status)
	}
	if input.TopN != 5 || input.RiskTopN != 10 {
		t.Fatalf("top_n/risk_top_n = %d/%d, want 5/10", input.TopN, input.RiskTopN)
	}
	if store.deleteCutoff == nil {
		t.Fatal("retention delete was not called")
	}
	summary, ok := input.Summary.(map[string]any)
	if !ok {
		t.Fatalf("summary type = %T, want map", input.Summary)
	}
	if summary["estimated_provider_requests"] != 2 {
		t.Fatalf("estimated_provider_requests = %v, want 2", summary["estimated_provider_requests"])
	}
	if summary["estimated_candidates_max"] != 15 {
		t.Fatalf("estimated_candidates_max = %v, want 15", summary["estimated_candidates_max"])
	}
}

func TestNormalizeLessonSearchPrerunSchedulerConfigDefaults(t *testing.T) {
	got := normalizeLessonSearchPrerunSchedulerConfig(LessonSearchPrerunSchedulerConfig{Provider: "unknown", TopN: -1, RiskTopN: -1, RiskFixtureIDs: []string{"a", "a,b"}})
	if got.Provider != "all" {
		t.Fatalf("provider = %q, want all", got.Provider)
	}
	if got.TopN != 5 {
		t.Fatalf("top_n = %d, want 5", got.TopN)
	}
	if got.RiskTopN != 0 {
		t.Fatalf("risk_top_n = %d, want 0", got.RiskTopN)
	}
	if len(got.RiskFixtureIDs) != 2 || got.RiskFixtureIDs[0] != "a" || got.RiskFixtureIDs[1] != "b" {
		t.Fatalf("risk fixture ids = %#v, want [a b]", got.RiskFixtureIDs)
	}
	if got.Interval != 24*time.Hour {
		t.Fatalf("interval = %s, want 24h", got.Interval)
	}
	if got.Timeout != 60*time.Second {
		t.Fatalf("timeout = %s, want 60s", got.Timeout)
	}
}

func writePrerunFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fixtures.json")
	data := `[
  {
    "id": "youtube-setup",
    "provider": "youtube",
    "query": "기초 강좌",
    "search_spec": {
      "primary_query": "기초 강좌",
      "must_include": ["기초"],
      "nice_to_have": ["입문"],
      "avoid": ["학원"],
      "content_types": ["video"],
      "language": "ko",
      "stage_role": "setup_intro"
    },
    "candidates": []
  }
]`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}
