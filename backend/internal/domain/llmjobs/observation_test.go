package llmjobs

import (
	"testing"
	"time"
)

func TestBuildObservationCutoffs(t *testing.T) {
	now := time.Date(2026, 5, 25, 12, 0, 0, 0, time.UTC)
	cutoffs := buildObservationCutoffs(now, RecoveryPolicy{
		QueuedTimeout:  10 * time.Minute,
		RunningTimeout: 15 * time.Minute,
	})
	if cutoffs.QueuedCutoff == nil || !cutoffs.QueuedCutoff.Equal(now.Add(-10*time.Minute)) {
		t.Fatalf("queued cutoff = %v, want %v", cutoffs.QueuedCutoff, now.Add(-10*time.Minute))
	}
	if cutoffs.RunningCutoff == nil || !cutoffs.RunningCutoff.Equal(now.Add(-15*time.Minute)) {
		t.Fatalf("running cutoff = %v, want %v", cutoffs.RunningCutoff, now.Add(-15*time.Minute))
	}
}

func TestBuildObservationCutoffsDisabled(t *testing.T) {
	now := time.Date(2026, 5, 25, 12, 0, 0, 0, time.UTC)
	cutoffs := buildObservationCutoffs(now, RecoveryPolicy{})
	if cutoffs.QueuedCutoff != nil {
		t.Fatalf("queued cutoff = %v, want nil", cutoffs.QueuedCutoff)
	}
	if cutoffs.RunningCutoff != nil {
		t.Fatalf("running cutoff = %v, want nil", cutoffs.RunningCutoff)
	}
}

func TestParseObservationDuration(t *testing.T) {
	fallback := 10 * time.Minute
	tests := []struct {
		name string
		raw  string
		want time.Duration
	}{
		{name: "empty", raw: "", want: fallback},
		{name: "duration", raw: "30m", want: 30 * time.Minute},
		{name: "milliseconds", raw: "1500", want: 1500 * time.Millisecond},
		{name: "zero disables", raw: "0", want: 0},
		{name: "negative duration rejected", raw: "-1s", want: fallback},
		{name: "invalid rejected", raw: "soon", want: fallback},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseObservationDuration(tt.raw, fallback); got != tt.want {
				t.Fatalf("parseObservationDuration(%q) = %s, want %s", tt.raw, got, tt.want)
			}
		})
	}
}
