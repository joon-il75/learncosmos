package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type evalCase struct {
	ID           string `json:"id"`
	SourceQuery  string `json:"source_query"`
	LearningGoal string `json:"learning_goal"`
	Language     string `json:"language"`
}

type patternMatch struct {
	Rank            int     `json:"rank"`
	PatternKey      string  `json:"pattern_key"`
	PatternVersion  string  `json:"pattern_version"`
	Domain          string  `json:"domain"`
	GoalType        string  `json:"goal_type"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	SimilarityScore float64 `json:"similarity_score"`
	Distance        float64 `json:"distance"`
}

type evalCaseResult struct {
	ID                 string                                     `json:"id"`
	Query              string                                     `json:"query"`
	SourceQuery        string                                     `json:"source_query"`
	LearningGoal       string                                     `json:"learning_goal"`
	Language           string                                     `json:"language"`
	CodeRouting        curriculum.CurriculumPatternRoutingPreview `json:"code_routing"`
	EmbeddingTopN      []patternMatch                             `json:"embedding_top_n"`
	CodeSubpatternRank int                                        `json:"code_subpattern_rank"`
	Top1Matched        bool                                       `json:"top1_matched"`
	Top3Matched        bool                                       `json:"top3_matched"`
	SkippedReason      string                                     `json:"skipped_reason,omitempty"`
	Error              string                                     `json:"error,omitempty"`
}

type evalSummary struct {
	TotalCases        int     `json:"total_cases"`
	EvaluatedCases    int     `json:"evaluated_cases"`
	SkippedCases      int     `json:"skipped_cases"`
	Top1Matches       int     `json:"top1_matches"`
	Top3Matches       int     `json:"top3_matches"`
	Top1MatchRate     float64 `json:"top1_match_rate"`
	Top3MatchRate     float64 `json:"top3_match_rate"`
	ReadyPatternRows  int     `json:"ready_pattern_rows"`
	Provider          string  `json:"provider"`
	Model             string  `json:"model"`
	Dimension         int     `json:"dimension"`
	TopN              int     `json:"top_n"`
	EvaluationVersion string  `json:"evaluation_version"`
}

type evalReport struct {
	GeneratedAt string           `json:"generated_at"`
	Summary     evalSummary      `json:"summary"`
	Results     []evalCaseResult `json:"results"`
}

func main() {
	_ = godotenv.Load()

	var (
		casesPath         = flag.String("cases", "", "optional evaluation cases JSON path")
		outPath           = flag.String("out", defaultOutputPath(), "output JSON report path")
		top               = flag.Int("top", 3, "number of embedding matches per case")
		provider          = flag.String("provider", "embedding_gemma", "pattern embedding provider")
		model             = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
		endpoint          = flag.String("endpoint", defaultString("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma HTTP endpoint")
		expectedDimension = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "expected embedding dimension")
		printCases        = flag.Bool("print-cases", false, "print default cases as JSON and exit")
	)
	flag.Parse()

	if *printCases {
		printDefaultCases()
		return
	}
	if *top <= 0 {
		log.Fatal("-top must be greater than 0")
	}
	if *expectedDimension <= 0 {
		log.Fatal("-dimension must be greater than 0")
	}
	if strings.TrimSpace(*endpoint) == "" {
		log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required")
	}
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	cases, err := loadCases(*casesPath)
	if err != nil {
		log.Fatalf("load cases: %v", err)
	}
	if len(cases) == 0 {
		log.Fatal("no evaluation cases")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	client, err := embedder.NewProviderClient(*provider, strings.TrimSpace(*endpoint))
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	readyCount, err := countReadyPatternEmbeddings(ctx, pool, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, "ko")
	if err != nil {
		log.Fatalf("count ready pattern embeddings: %v", err)
	}

	results := make([]evalCaseResult, 0, len(cases))
	summary := evalSummary{
		TotalCases:        len(cases),
		ReadyPatternRows:  readyCount,
		Provider:          strings.TrimSpace(*provider),
		Model:             strings.TrimSpace(*model),
		Dimension:         *expectedDimension,
		TopN:              *top,
		EvaluationVersion: "curriculum_pattern_match_eval_v1",
	}

	for _, item := range cases {
		result := evaluateCase(ctx, pool, client, item, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, *top)
		results = append(results, result)
		if result.SkippedReason != "" {
			summary.SkippedCases++
			continue
		}
		summary.EvaluatedCases++
		if result.Top1Matched {
			summary.Top1Matches++
		}
		if result.Top3Matched {
			summary.Top3Matches++
		}
	}
	if summary.EvaluatedCases > 0 {
		summary.Top1MatchRate = float64(summary.Top1Matches) / float64(summary.EvaluatedCases)
		summary.Top3MatchRate = float64(summary.Top3Matches) / float64(summary.EvaluatedCases)
	}

	report := evalReport{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Summary:     summary,
		Results:     results,
	}
	if err := writeReport(*outPath, report); err != nil {
		log.Fatalf("write report: %v", err)
	}

	log.Printf(
		"curriculum pattern match evaluation completed cases=%d evaluated=%d top1=%d(%.1f%%) top3=%d(%.1f%%) out=%s",
		summary.TotalCases,
		summary.EvaluatedCases,
		summary.Top1Matches,
		summary.Top1MatchRate*100,
		summary.Top3Matches,
		summary.Top3MatchRate*100,
		*outPath,
	)
}

func evaluateCase(ctx context.Context, pool *pgxpool.Pool, client embedder.Embedder, item evalCase, provider, model string, dimension, top int) evalCaseResult {
	language := normalizeLanguage(item.Language)
	query := strings.Join(compactParts(item.SourceQuery, item.LearningGoal), "\n")
	result := evalCaseResult{
		ID:           item.ID,
		Query:        query,
		SourceQuery:  strings.TrimSpace(item.SourceQuery),
		LearningGoal: strings.TrimSpace(item.LearningGoal),
		Language:     language,
		CodeRouting:  curriculum.AnalyzeCurriculumPatternRoutingPreview(item.SourceQuery, item.LearningGoal, language),
	}
	if strings.TrimSpace(query) == "" {
		result.SkippedReason = "empty query"
		return result
	}
	if strings.TrimSpace(result.CodeRouting.SubpatternKey) == "" {
		result.SkippedReason = "code routing subpattern is empty"
		return result
	}

	embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	vector, err := client.Embed(embedCtx, query)
	cancel()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if len(vector) != dimension {
		result.Error = fmt.Sprintf("unexpected query embedding dimension: got=%d want=%d", len(vector), dimension)
		return result
	}

	matches, err := searchPatternMatches(ctx, pool, vectorToPGVector(vector), provider, model, dimension, language, top)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.EmbeddingTopN = matches
	for _, match := range matches {
		if match.PatternKey == result.CodeRouting.SubpatternKey {
			result.CodeSubpatternRank = match.Rank
			break
		}
	}
	result.Top1Matched = result.CodeSubpatternRank == 1
	result.Top3Matched = result.CodeSubpatternRank > 0 && result.CodeSubpatternRank <= 3
	return result
}

func loadCases(path string) ([]evalCase, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultCases(), nil
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []evalCase
	if err := json.Unmarshal(payload, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}

func writeReport(path string, report evalReport) error {
	payload, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}

func printDefaultCases() {
	payload, err := json.MarshalIndent(defaultCases(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(payload))
}

func countReadyPatternEmbeddings(ctx context.Context, pool *pgxpool.Pool, provider, model string, dimension int, language string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cp.language = $1
		  AND cpe.provider = $2
		  AND cpe.model = $3
		  AND cpe.dimension = $4
		  AND cpe.status = 'ready'
		  AND cpe.embedding IS NOT NULL
	`, language, provider, model, dimension).Scan(&count)
	return count, err
}

func searchPatternMatches(ctx context.Context, pool *pgxpool.Pool, pgVector, provider, model string, dimension int, language string, limit int) ([]patternMatch, error) {
	rows, err := pool.Query(ctx, `
		SELECT
		  cp.pattern_key,
		  cp.pattern_version,
		  cp.domain,
		  cp.goal_type,
		  cp.language,
		  cp.title,
		  cpe.embedding <=> $1::vector AS distance,
		  1 - (cpe.embedding <=> $1::vector) AS similarity_score
		FROM curriculum_pattern_embeddings cpe
		JOIN curriculum_patterns cp ON cp.id = cpe.pattern_id
		WHERE cp.is_active = TRUE
		  AND cp.language = $2
		  AND cpe.provider = $3
		  AND cpe.model = $4
		  AND cpe.dimension = $5
		  AND cpe.status = 'ready'
		  AND cpe.embedding IS NOT NULL
		ORDER BY cpe.embedding <=> $1::vector ASC
		LIMIT $6
	`, pgVector, language, provider, model, dimension, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []patternMatch{}
	rank := 1
	for rows.Next() {
		var item patternMatch
		item.Rank = rank
		if err := rows.Scan(
			&item.PatternKey,
			&item.PatternVersion,
			&item.Domain,
			&item.GoalType,
			&item.Language,
			&item.Title,
			&item.Distance,
			&item.SimilarityScore,
		); err != nil {
			return nil, err
		}
		matches = append(matches, item)
		rank++
	}
	return matches, rows.Err()
}

func vectorToPGVector(vector []float32) string {
	parts := make([]string, 0, len(vector))
	for _, value := range vector {
		parts = append(parts, fmt.Sprintf("%f", value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func compactParts(parts ...string) []string {
	result := []string{}
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeLanguage(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func defaultOutputPath() string {
	return filepath.Clean("../docs/data/curriculum-patterns/evaluation_curriculum_pattern_match_v1.json")
}

func defaultString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func defaultInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func defaultCases() []evalCase {
	return []evalCase{
		{ID: "leather-wallet", SourceQuery: "가죽공예 카드지갑 만들기", LearningGoal: "가죽공예 카드지갑을 직접 만들고 마감 품질을 점검하고 싶다", Language: "ko"},
		{ID: "leather-bag", SourceQuery: "가죽공예 가방 만들기", LearningGoal: "입체 구조가 있는 가죽 가방을 직접 완성하고 싶다", Language: "ko"},
		{ID: "calligraphy-basic", SourceQuery: "붓펜 캘리그라피 엽서", LearningGoal: "붓펜으로 짧은 문구 엽서를 직접 쓰고 싶다", Language: "ko"},
		{ID: "calligraphy-showcase", SourceQuery: "캘리그라피 작품 만들기", LearningGoal: "캘리그라피 문구 작품을 완성해 온라인에 게시하고 싶다", Language: "ko"},
		{ID: "travel-conversation", SourceQuery: "영어 여행 회화", LearningGoal: "여행 중 공항과 식당에서 영어로 자연스럽게 대화하고 싶다", Language: "ko"},
		{ID: "toeic-score", SourceQuery: "토익 800점 준비", LearningGoal: "TOEIC 800점을 목표로 파트별 문제풀이와 시간 대응을 연습하고 싶다", Language: "ko"},
		{ID: "opic-speaking", SourceQuery: "오픽 영어 말하기", LearningGoal: "OPIc 인터뷰에서 경험 답변과 돌발 질문에 자연스럽게 대응하고 싶다", Language: "ko"},
		{ID: "ielts-four-skills", SourceQuery: "IELTS 6.5 시험 준비", LearningGoal: "IELTS 6.5를 목표로 영어 reading listening writing speaking 네 영역을 준비하고 싶다", Language: "ko"},
		{ID: "midi-beat", SourceQuery: "무료 DAW MIDI 비트 만들기", LearningGoal: "무료 DAW로 MIDI 첫 비트를 만들어 짧은 루프를 완성하고 싶다", Language: "ko"},
		{ID: "midi-track", SourceQuery: "BandLab MIDI track 만들기", LearningGoal: "BandLab으로 전체 MIDI 트랙을 구성하고 완성하고 싶다", Language: "ko"},
		{ID: "vibe-mvp", SourceQuery: "vibe coding app 만들기", LearningGoal: "바이브 코딩으로 작은 MVP 웹앱을 만들어 실행해 보고 싶다", Language: "ko"},
		{ID: "vibe-fullstack", SourceQuery: "vibe coding full-stack 배포", LearningGoal: "바이브 코딩으로 full-stack 웹앱을 만들고 Vercel에 배포하고 싶다", Language: "ko"},
		{ID: "claude-code", SourceQuery: "Claude Code agentic workflow", LearningGoal: "Claude Code로 작업 계획을 세우고 agentic workflow로 기능을 구현하고 싶다", Language: "ko"},
		{ID: "korean-open-lecture", SourceQuery: "KOCW 공개강의 수강", LearningGoal: "한국 공개강의를 주차별로 수강하고 핵심 개념을 정리하고 싶다", Language: "ko"},
		{ID: "kmooc-online", SourceQuery: "K-MOOC 온라인 강좌", LearningGoal: "K-MOOC 온라인 강좌를 끝까지 수료하고 생활에 적용하고 싶다", Language: "ko"},
		{ID: "gtq-design-cert", SourceQuery: "GTQ 포토샵 실기 자격", LearningGoal: "GTQ 포토샵 실기 과제를 시간 안에 완성하고 싶다", Language: "ko"},
		{ID: "ncs-web-publisher", SourceQuery: "NCS 웹퍼블리셔 직업훈련", LearningGoal: "웹퍼블리셔 직무에 필요한 HTML CSS JavaScript 포트폴리오를 만들고 싶다", Language: "ko"},
		{ID: "baking-cert", SourceQuery: "제과기능사 제빵기능사 자격증", LearningGoal: "제과기능사와 제빵기능사 실기 과제를 반복 연습해 시험에 대비하고 싶다", Language: "ko"},
		{ID: "korean-cooking-cert", SourceQuery: "한식조리기능사 실기", LearningGoal: "한식 조리 실기 공개과제를 시간 안에 완성하고 싶다", Language: "ko"},
		{ID: "baking-class-demo", SourceQuery: "베이킹 클래스 시연", LearningGoal: "베이킹 수업에서 레시피를 설명하고 시연할 수 있게 준비하고 싶다", Language: "ko"},
		{ID: "drone-basic", SourceQuery: "드론 4종 이론", LearningGoal: "드론 4종 이론을 학습하고 안전 운항 기준을 이해하고 싶다", Language: "ko"},
		{ID: "drone-practical", SourceQuery: "드론 실기 자격", LearningGoal: "드론 실기 조종 절차와 코스 비행을 연습하고 싶다", Language: "ko"},
		{ID: "cloud-foundation", SourceQuery: "AWS cloud practitioner 자격", LearningGoal: "클라우드 입문 자격 시험 범위와 서비스 개념을 정리하고 싶다", Language: "ko"},
		{ID: "cloud-operator", SourceQuery: "Azure Administrator 실무 자격", LearningGoal: "클라우드 운영과 네트워크 보안 시나리오를 연습해 실무 자격을 준비하고 싶다", Language: "ko"},
		{ID: "daily-writing", SourceQuery: "매일 글쓰기 루틴", LearningGoal: "매일 짧은 글을 쓰고 퇴고하는 습관을 만들고 싶다", Language: "ko"},
		{ID: "brunch-serial", SourceQuery: "브런치 연재 글쓰기", LearningGoal: "브런치에 올릴 에세이 연재를 기획하고 첫 글을 발행하고 싶다", Language: "ko"},
		{ID: "urban-agriculture", SourceQuery: "도시농업 실습", LearningGoal: "베란다 텃밭을 계획하고 작물 관리 기록을 남기고 싶다", Language: "ko"},
		{ID: "ai-literacy", SourceQuery: "AI 디지털배움터 스마트폰 키오스크", LearningGoal: "스마트폰과 키오스크 생활앱을 안전하게 쓰는 디지털 문해 역량을 기르고 싶다", Language: "ko"},
		{ID: "drawing-foundation", SourceQuery: "Drawspace drawing lessons", LearningGoal: "Learn contour drawing, shading, and perspective to finish a small sketch", Language: "ko"},
		{ID: "portfolio-publish", SourceQuery: "영상편집 배우기", LearningGoal: "편집한 영상을 포트폴리오와 유튜브에 공개하고 싶다", Language: "ko"},
	}
}
