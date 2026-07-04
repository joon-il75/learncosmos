package main

import "testing"

func TestArtifactAuditSummaryCountsSanitizerImpactWithoutRawOutput(t *testing.T) {
	var summary artifactAuditSummary
	summary.add(`<p>safe <strong>text</strong></p>`)
	summary.add(`<p onclick="bad()">hello<script>bad()</script><a href="https://example.com">ok</a></p>`)
	summary.add(`<script>bad()</script>`)

	if summary.ChangedCandidates != 2 {
		t.Fatalf("ChangedCandidates = %d, want 2", summary.ChangedCandidates)
	}
	if summary.EmptyAfterSanitize != 1 {
		t.Fatalf("EmptyAfterSanitize = %d, want 1", summary.EmptyAfterSanitize)
	}
	if summary.RiskyFragmentRemoved != 2 {
		t.Fatalf("RiskyFragmentRemoved = %d, want 2", summary.RiskyFragmentRemoved)
	}
}

func TestArtifactRemovedRiskyFragment(t *testing.T) {
	if !artifactRemovedRiskyFragment(`<a href="javascript:bad()">bad</a>`, `bad`) {
		t.Fatal("artifactRemovedRiskyFragment() = false, want true")
	}
	if artifactRemovedRiskyFragment(`<p>safe</p>`, `<p>safe</p>`) {
		t.Fatal("artifactRemovedRiskyFragment() = true, want false")
	}
}
