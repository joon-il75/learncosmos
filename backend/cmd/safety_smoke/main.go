package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/learnweaver/backend/internal/pkg/db"
)

type smokeCase struct {
	ID             string `json:"id"`
	Locale         string `json:"locale"`
	Text           string `json:"text"`
	ExpectedAction string `json:"expected_action"`
	Note           string `json:"note"`
}

type smokeResult struct {
	ID       string `json:"id"`
	Locale   string `json:"locale"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	RiskType string `json:"risk_type"`
	Passed   bool   `json:"passed"`
}

func main() {
	_ = godotenv.Load()

	casesPath := flag.String("cases", "../docs/data/safety-moderation/qa/false_positive_smoke.jsonl", "JSONL smoke cases path")
	jsonOutput := flag.Bool("json", false, "print JSON lines for each result")
	flag.Parse()

	cases, err := readCases(*casesPath)
	if err != nil {
		log.Fatalf("read smoke cases: %v", err)
	}
	ctx := context.Background()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer pool.Close()

	svc := safety.NewService(safety.NewRepository(pool))
	failed := 0
	for _, item := range cases {
		result, err := svc.Moderate(ctx, safety.ModerateInput{Text: item.Text, Locale: item.Locale, TargetType: "safety_smoke"})
		if err != nil {
			log.Fatalf("moderate %s: %v", item.ID, err)
		}
		actual := string(result.Action)
		passed := actual == strings.TrimSpace(item.ExpectedAction)
		if !passed {
			failed++
		}
		row := smokeResult{ID: item.ID, Locale: item.Locale, Expected: item.ExpectedAction, Actual: actual, RiskType: result.RiskType, Passed: passed}
		if *jsonOutput {
			b, _ := json.Marshal(row)
			fmt.Println(string(b))
			continue
		}
		status := "PASS"
		if !passed {
			status = "FAIL"
		}
		fmt.Printf("%s id=%s locale=%s expected=%s actual=%s risk=%s\n", status, item.ID, item.Locale, item.ExpectedAction, actual, result.RiskType)
	}
	fmt.Printf("summary total=%d failed=%d\n", len(cases), failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func readCases(path string) ([]smokeCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cases := []smokeCase{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var item smokeCase
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Text) == "" || strings.TrimSpace(item.ExpectedAction) == "" {
			return nil, fmt.Errorf("line %d: id/text/expected_action are required", lineNo)
		}
		cases = append(cases, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no smoke cases")
	}
	return cases, nil
}
