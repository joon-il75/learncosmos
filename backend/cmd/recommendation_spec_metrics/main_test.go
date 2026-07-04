package main

import "testing"

func TestRatioPtr(t *testing.T) {
	if got := ratioPtr(1, 0); got != nil {
		t.Fatalf("ratioPtr denominator zero = %v, want nil", *got)
	}
	got := ratioPtr(1, 4)
	if got == nil {
		t.Fatal("ratioPtr returned nil")
	}
	if *got != 0.25 {
		t.Fatalf("ratioPtr = %v, want 0.25", *got)
	}
}

func TestSortedMetricKeys(t *testing.T) {
	got := sortedMetricKeys(map[string]int64{"fallback": 1, "empty": 2, "pattern_template": 3})
	want := []string{"empty", "fallback", "pattern_template"}
	if len(got) != len(want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %#v, want %#v", got, want)
		}
	}
}
