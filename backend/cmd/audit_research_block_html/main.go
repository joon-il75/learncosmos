package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/curriculum"
)

type auditSummary struct {
	Scanned              int `json:"scanned"`
	ChangedCandidates    int `json:"changed_candidates"`
	EmptyAfterSanitize   int `json:"empty_after_sanitize"`
	RiskyFragmentRemoved int `json:"risky_fragment_removed"`
	ParseErrors          int `json:"parse_errors"`
}

func main() {
	limit := flag.Int("limit", 10000, "maximum research text blocks to scan")
	timeout := flag.Duration("timeout", 15*time.Second, "database audit timeout")
	flag.Parse()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if *limit <= 0 {
		log.Fatal("limit must be positive")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	rows, err := pool.Query(ctx, `
		SELECT content
		FROM course_point_blocks
		WHERE block_type = 'text'
		  AND content ? 'text'
		ORDER BY updated_at DESC
		LIMIT $1
	`, *limit)
	if err != nil {
		log.Fatalf("query research blocks: %v", err)
	}
	defer rows.Close()

	var summary auditSummary
	for rows.Next() {
		var raw json.RawMessage
		if err := rows.Scan(&raw); err != nil {
			log.Fatalf("scan research block: %v", err)
		}
		summary.add(raw)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("iterate research blocks: %v", err)
	}

	encoded, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		log.Fatalf("encode summary: %v", err)
	}
	fmt.Println(string(encoded))
}

func (s *auditSummary) add(raw json.RawMessage) {
	s.Scanned++

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		s.ParseErrors++
		return
	}
	text, _ := payload["text"].(string)
	cleaned := curriculum.SanitizeTiptapHTMLForStorage(text)
	if cleaned != text {
		s.ChangedCandidates++
	}
	if cleaned == "<p></p>" && strings.TrimSpace(text) != "" {
		s.EmptyAfterSanitize++
	}
	if removedRiskyFragment(text, cleaned) {
		s.RiskyFragmentRemoved++
	}
}

func removedRiskyFragment(before, after string) bool {
	before = strings.ToLower(before)
	after = strings.ToLower(after)
	for _, fragment := range []string{
		"onerror", "onclick", "onload", "javascript:", "data:image", "<script", "<style", "<iframe", "<svg",
	} {
		if strings.Contains(before, fragment) && !strings.Contains(after, fragment) {
			return true
		}
	}
	return false
}
