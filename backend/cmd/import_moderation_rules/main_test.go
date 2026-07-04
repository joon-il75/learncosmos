package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApprovedRulesFromCSVFiltersOnlyApprovedImportReadyRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "review.csv")
	content := "locale,pattern,source_id,source_line,source_license,suggested_risk_type,suggested_action,review_priority,review_decision,final_risk_type,final_action,import_ready,review_reason,pattern_sha256\n" +
		"ko,검수어,ldnoobw,1,CC-BY-4.0,abuse_or_harassment,soft_warn,normal_review,approve,hate_or_discrimination,block,true,clear severe phrase,abc\n" +
		"ko,보류어,ldnoobw,2,CC-BY-4.0,abuse_or_harassment,soft_warn,normal_review,approve,,soft_warn,false,,def\n" +
		"en,rejected,ldnoobw,3,CC-BY-4.0,abuse_or_harassment,soft_warn,normal_review,reject,,soft_warn,true,,ghi\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	rules, err := approvedRulesFromCSV(path)
	if err != nil {
		t.Fatalf("approvedRulesFromCSV returned error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 approved rule, got %d", len(rules))
	}
	got := rules[0]
	if got.Pattern != "검수어" || got.Action != "block" || got.RiskType != "hate_or_discrimination" {
		t.Fatalf("unexpected rule: %+v", got)
	}
	if got.RuleType != "keyword" {
		t.Fatalf("expected keyword rule, got %q", got.RuleType)
	}
	if got.Description == "" || got.SourceURL == "" {
		t.Fatalf("expected attribution fields, got %+v", got)
	}
}

func TestValidateRulesRejectsDuplicatePatternLocale(t *testing.T) {
	rules := []reviewedRule{
		{Locale: "ko", RuleType: "keyword", Pattern: "중복", RiskType: "abuse_or_harassment", Action: "soft_warn"},
		{Locale: "ko", RuleType: "keyword", Pattern: "중복", RiskType: "abuse_or_harassment", Action: "soft_warn"},
	}
	if err := validateRules(rules); err == nil {
		t.Fatal("expected duplicate validation error")
	}
}
