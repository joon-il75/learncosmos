package main

import (
	"encoding/json"
	"testing"
)

func TestAuditSummaryCountsSanitizerImpactWithoutRawOutput(t *testing.T) {
	var summary auditSummary
	summary.add(json.RawMessage(`{"text":"<p>safe <strong>text</strong></p>"}`))
	summary.add(json.RawMessage(`{"text":"<p onclick=\"bad()\">hello<script>bad()</script><a href=\"https://example.com\">ok</a></p>"}`))
	summary.add(json.RawMessage(`{"text":"<script>bad()</script>"}`))
	summary.add(json.RawMessage(`{"text":`))

	if summary.Scanned != 4 {
		t.Fatalf("Scanned = %d, want 4", summary.Scanned)
	}
	if summary.ChangedCandidates != 2 {
		t.Fatalf("ChangedCandidates = %d, want 2", summary.ChangedCandidates)
	}
	if summary.EmptyAfterSanitize != 1 {
		t.Fatalf("EmptyAfterSanitize = %d, want 1", summary.EmptyAfterSanitize)
	}
	if summary.RiskyFragmentRemoved != 2 {
		t.Fatalf("RiskyFragmentRemoved = %d, want 2", summary.RiskyFragmentRemoved)
	}
	if summary.ParseErrors != 1 {
		t.Fatalf("ParseErrors = %d, want 1", summary.ParseErrors)
	}
}

func TestRemovedRiskyFragment(t *testing.T) {
	if !removedRiskyFragment(`<img src="x" onerror="bad()">`, `<p></p>`) {
		t.Fatal("removedRiskyFragment() = false, want true")
	}
	if removedRiskyFragment(`<p>safe</p>`, `<p>safe</p>`) {
		t.Fatal("removedRiskyFragment() = true, want false")
	}
}
