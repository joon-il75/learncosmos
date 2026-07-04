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

type artifactAuditSummary struct {
	ScannedPointArtifacts int `json:"scanned_point_artifacts"`
	ScannedDraftArtifacts int `json:"scanned_draft_artifacts"`
	ChangedCandidates     int `json:"changed_candidates"`
	EmptyAfterSanitize    int `json:"empty_after_sanitize"`
	RiskyFragmentRemoved  int `json:"risky_fragment_removed"`
}

func main() {
	limit := flag.Int("limit", 10000, "maximum rows per artifact table to scan")
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

	var summary artifactAuditSummary
	if err := scanArtifactDescriptions(ctx, pool, *limit, `
		SELECT description
		FROM course_point_artifacts
		ORDER BY updated_at DESC
		LIMIT $1
	`, func() { summary.ScannedPointArtifacts++ }, &summary); err != nil {
		log.Fatalf("scan course_point_artifacts: %v", err)
	}
	if err := scanArtifactDescriptions(ctx, pool, *limit, `
		SELECT description
		FROM course_draft_artifact_entries
		ORDER BY updated_at DESC
		LIMIT $1
	`, func() { summary.ScannedDraftArtifacts++ }, &summary); err != nil {
		log.Fatalf("scan course_draft_artifact_entries: %v", err)
	}

	encoded, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		log.Fatalf("encode summary: %v", err)
	}
	fmt.Println(string(encoded))
}

func scanArtifactDescriptions(ctx context.Context, pool *pgxpool.Pool, limit int, query string, incrementScanned func(), summary *artifactAuditSummary) error {
	rows, err := pool.Query(ctx, query, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var description string
		if err := rows.Scan(&description); err != nil {
			return err
		}
		incrementScanned()
		summary.add(description)
	}
	return rows.Err()
}

func (s *artifactAuditSummary) add(description string) {
	cleaned := curriculum.SanitizeTiptapHTMLForStorage(description)
	if cleaned != description {
		s.ChangedCandidates++
	}
	if cleaned == "<p></p>" && strings.TrimSpace(description) != "" {
		s.EmptyAfterSanitize++
	}
	if artifactRemovedRiskyFragment(description, cleaned) {
		s.RiskyFragmentRemoved++
	}
}

func artifactRemovedRiskyFragment(before, after string) bool {
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
