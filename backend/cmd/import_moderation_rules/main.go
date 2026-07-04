package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/pkg/db"
)

type reviewedRule struct {
	Locale         string `json:"locale"`
	RuleType       string `json:"rule_type"`
	Pattern        string `json:"pattern"`
	RiskType       string `json:"risk_type"`
	Action         string `json:"action"`
	Description    string `json:"description"`
	SourceID       string `json:"source_id"`
	SourceLicense  string `json:"source_license"`
	SourceRepo     string `json:"source_repo"`
	SourceURL      string `json:"source_url"`
	SourceLine     string `json:"source_line"`
	PatternSHA256  string `json:"pattern_sha256"`
	ReviewDecision string `json:"review_decision"`
	ReviewReason   string `json:"review_reason"`
}

type importSummary struct {
	Total       int
	Inserted    int64
	Duplicates  int
	ActionCount map[string]int
	RiskCount   map[string]int
	LocaleCount map[string]int
}

func main() {
	_ = godotenv.Load()

	var (
		reviewCSV = flag.String("review-csv", "", "review CSV path; when set with -out, generate reviewed JSONL")
		outPath   = flag.String("out", "", "output JSONL path for approved review rows")
		rulesPath = flag.String("rules-jsonl", "", "reviewed rules JSONL path to dry-run or apply")
		dryRun    = flag.Bool("dry-run", true, "summarize without inserting")
		apply     = flag.Bool("apply", false, "insert reviewed rules into moderation_rules")
	)
	flag.Parse()

	switch {
	case strings.TrimSpace(*reviewCSV) != "":
		if strings.TrimSpace(*outPath) == "" {
			log.Fatal("-out is required with -review-csv")
		}
		rules, err := approvedRulesFromCSV(*reviewCSV)
		if err != nil {
			log.Fatalf("read review csv: %v", err)
		}
		if err := writeRulesJSONL(*outPath, rules); err != nil {
			log.Fatalf("write reviewed rules jsonl: %v", err)
		}
		log.Printf("generated reviewed rules jsonl path=%s count=%d", *outPath, len(rules))
	case strings.TrimSpace(*rulesPath) != "":
		rules, err := readRulesJSONL(*rulesPath)
		if err != nil {
			log.Fatalf("read rules jsonl: %v", err)
		}
		if err := validateRules(rules); err != nil {
			log.Fatalf("validate rules: %v", err)
		}
		ctx := context.Background()
		pool, err := openPool(ctx)
		if err != nil {
			log.Fatalf("open db: %v", err)
		}
		defer pool.Close()
		summary, err := summarizeAndMaybeImport(ctx, pool, rules, *apply && !*dryRun)
		if err != nil {
			log.Fatalf("import moderation rules: %v", err)
		}
		mode := "dry-run"
		if *apply && !*dryRun {
			mode = "apply"
		}
		log.Printf("mode=%s total=%d inserted=%d duplicates=%d locale=%v action=%v risk=%v", mode, summary.Total, summary.Inserted, summary.Duplicates, summary.LocaleCount, summary.ActionCount, summary.RiskCount)
	default:
		log.Fatal("set either -review-csv ... -out ... or -rules-jsonl ...")
	}
}

func approvedRulesFromCSV(path string) ([]reviewedRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty csv")
	}
	header := map[string]int{}
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}
	required := []string{"locale", "pattern", "source_id", "source_line", "source_license", "suggested_risk_type", "suggested_action", "review_decision", "final_risk_type", "final_action", "import_ready", "review_reason", "pattern_sha256"}
	for _, col := range required {
		if _, ok := header[col]; !ok {
			return nil, fmt.Errorf("missing column %q", col)
		}
	}

	rules := make([]reviewedRule, 0)
	for rowIndex, row := range rows[1:] {
		get := func(col string) string {
			i := header[col]
			if i >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[i])
		}
		if !strings.EqualFold(get("review_decision"), "approve") || !strings.EqualFold(get("import_ready"), "true") {
			continue
		}
		riskType := defaultIfBlank(get("final_risk_type"), get("suggested_risk_type"))
		action := defaultIfBlank(get("final_action"), get("suggested_action"))
		rule := reviewedRule{
			Locale:         get("locale"),
			RuleType:       "keyword",
			Pattern:        get("pattern"),
			RiskType:       riskType,
			Action:         action,
			Description:    buildDescription(get("source_id"), get("source_license"), get("source_line"), get("review_reason")),
			SourceID:       get("source_id"),
			SourceLicense:  get("source_license"),
			SourceRepo:     "https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words",
			SourceURL:      sourceURL(get("locale")),
			SourceLine:     get("source_line"),
			PatternSHA256:  get("pattern_sha256"),
			ReviewDecision: get("review_decision"),
			ReviewReason:   get("review_reason"),
		}
		if err := validateRule(rule); err != nil {
			return nil, fmt.Errorf("row %d: %w", rowIndex+2, err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func writeRulesJSONL(path string, rules []reviewedRule) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	for _, rule := range rules {
		if err := encoder.Encode(rule); err != nil {
			return err
		}
	}
	return nil
}

func readRulesJSONL(path string) ([]reviewedRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	rules := make([]reviewedRule, 0)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rule reviewedRule
		if err := json.Unmarshal([]byte(line), &rule); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		rules = append(rules, rule)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func validateRules(rules []reviewedRule) error {
	seen := map[string]struct{}{}
	for i, rule := range rules {
		if err := validateRule(rule); err != nil {
			return fmt.Errorf("rule %d: %w", i+1, err)
		}
		key := strings.TrimSpace(rule.RuleType) + "\x00" + strings.TrimSpace(rule.Pattern) + "\x00" + strings.TrimSpace(rule.Locale)
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate jsonl rule rule_type/pattern/locale=%q/%q/%q", rule.RuleType, rule.Pattern, rule.Locale)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func validateRule(rule reviewedRule) error {
	if strings.TrimSpace(rule.Pattern) == "" {
		return fmt.Errorf("pattern is required")
	}
	if rule.RuleType != "keyword" && rule.RuleType != "regex" {
		return fmt.Errorf("invalid rule_type %q", rule.RuleType)
	}
	if strings.TrimSpace(rule.RiskType) == "" {
		return fmt.Errorf("risk_type is required")
	}
	if rule.Action != "soft_warn" && rule.Action != "block" {
		return fmt.Errorf("invalid action %q", rule.Action)
	}
	return nil
}

func openPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return db.NewPool(ctx, databaseURL)
}

func summarizeAndMaybeImport(ctx context.Context, pool *pgxpool.Pool, rules []reviewedRule, apply bool) (importSummary, error) {
	summary := importSummary{
		Total:       len(rules),
		ActionCount: map[string]int{},
		RiskCount:   map[string]int{},
		LocaleCount: map[string]int{},
	}
	if len(rules) == 0 {
		return summary, nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return summary, err
	}
	defer tx.Rollback(ctx)

	for _, rule := range rules {
		summary.ActionCount[rule.Action]++
		summary.RiskCount[rule.RiskType]++
		summary.LocaleCount[defaultIfBlank(rule.Locale, "global")]++

		var exists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM moderation_rules
				WHERE rule_type = $1
				  AND pattern = $2
				  AND COALESCE(locale, '') = COALESCE($3, '')
			)
		`, rule.RuleType, rule.Pattern, nullableLocale(rule.Locale)).Scan(&exists); err != nil {
			return summary, err
		}
		if exists {
			summary.Duplicates++
			continue
		}
		if !apply {
			continue
		}
		tag, err := tx.Exec(ctx, `
			INSERT INTO moderation_rules (rule_type, pattern, risk_type, action, locale, description, is_active)
			VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, TRUE)
			ON CONFLICT DO NOTHING
		`, rule.RuleType, rule.Pattern, rule.RiskType, rule.Action, strings.TrimSpace(rule.Locale), rule.Description)
		if err != nil {
			return summary, err
		}
		summary.Inserted += tag.RowsAffected()
	}

	if apply {
		if err := tx.Commit(ctx); err != nil {
			return summary, err
		}
	}
	return summary, nil
}

func buildDescription(sourceID, license, sourceLine, reason string) string {
	parts := []string{sourceID, license}
	if strings.TrimSpace(sourceLine) != "" {
		parts = append(parts, "source_line="+strings.TrimSpace(sourceLine))
	}
	parts = append(parts, "reviewed=2026-05-29")
	if strings.TrimSpace(reason) != "" {
		parts = append(parts, "reason="+strings.TrimSpace(reason))
	}
	return strings.Join(parts, " ")
}

func sourceURL(locale string) string {
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return "https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words"
	}
	return "https://raw.githubusercontent.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/master/" + locale
}

func nullableLocale(locale string) any {
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return nil
	}
	return locale
}

func defaultIfBlank(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return strings.TrimSpace(fallback)
	}
	return value
}
