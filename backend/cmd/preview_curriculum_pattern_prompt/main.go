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

	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
)

type promptPreviewOutput struct {
	SourceQuery            string                                     `json:"source_query"`
	LearningGoal           string                                     `json:"learning_goal"`
	Language               string                                     `json:"language"`
	Provider               string                                     `json:"provider"`
	Model                  string                                     `json:"model"`
	Dimension              int                                        `json:"dimension"`
	CodeRouting            curriculum.CurriculumPatternRoutingPreview `json:"code_routing"`
	EmbeddingTopN          []curriculum.CurriculumPatternMatch        `json:"embedding_top_n"`
	CodeRank               int                                        `json:"code_rank"`
	TopNMatched            bool                                       `json:"top_n_matched"`
	PatternMatchGuidance   string                                     `json:"pattern_match_guidance"`
	CodeOnlyPrompt         string                                     `json:"code_only_prompt"`
	EmbeddingGuidedPrompt  string                                     `json:"embedding_guided_prompt"`
	PromptChanged          bool                                       `json:"prompt_changed"`
	FallbackReason         string                                     `json:"fallback_reason,omitempty"`
	GuidanceMatchCount     int                                        `json:"guidance_match_count"`
	HasCrossDomainGuidance bool                                       `json:"has_cross_domain_guidance"`
	UsedFallbackGuidance   bool                                       `json:"used_fallback_guidance"`
}

type promptPreviewCase struct {
	ID           string `json:"id"`
	SourceQuery  string `json:"source_query"`
	LearningGoal string `json:"learning_goal"`
	CurrentLevel string `json:"current_level,omitempty"`
	Language     string `json:"language"`
}

type promptPreviewCaseResult struct {
	ID                     string                                     `json:"id"`
	SourceQuery            string                                     `json:"source_query"`
	LearningGoal           string                                     `json:"learning_goal"`
	Language               string                                     `json:"language"`
	CodeRouting            curriculum.CurriculumPatternRoutingPreview `json:"code_routing"`
	EmbeddingTopN          []curriculum.CurriculumPatternMatch        `json:"embedding_top_n"`
	CodeRank               int                                        `json:"code_rank"`
	TopNMatched            bool                                       `json:"top_n_matched"`
	GuidanceMatchKeys      []string                                   `json:"guidance_match_keys"`
	PatternMatchGuidance   string                                     `json:"pattern_match_guidance"`
	PromptChanged          bool                                       `json:"prompt_changed"`
	HasCrossDomainGuidance bool                                       `json:"has_cross_domain_guidance"`
	UsedFallbackGuidance   bool                                       `json:"used_fallback_guidance"`
	SkippedReason          string                                     `json:"skipped_reason,omitempty"`
	Error                  string                                     `json:"error,omitempty"`
}

type promptPreviewReport struct {
	GeneratedAt         string                    `json:"generated_at"`
	Provider            string                    `json:"provider"`
	Model               string                    `json:"model"`
	Dimension           int                       `json:"dimension"`
	TopN                int                       `json:"top_n"`
	GuidanceLimit       int                       `json:"guidance_limit"`
	TotalCases          int                       `json:"total_cases"`
	EvaluatedCases      int                       `json:"evaluated_cases"`
	SkippedCases        int                       `json:"skipped_cases"`
	CrossDomainGuidance int                       `json:"cross_domain_guidance"`
	PromptChangedCases  int                       `json:"prompt_changed_cases"`
	TopNMatchedCases    int                       `json:"top_n_matched_cases"`
	CodeRank2PlusCases  int                       `json:"code_rank_2_plus_cases"`
	CodeUnmatchedCases  int                       `json:"code_unmatched_cases"`
	FallbackGuidance    int                       `json:"fallback_guidance_cases"`
	Results             []promptPreviewCaseResult `json:"results"`
}

func main() {
	_ = godotenv.Load()

	var (
		sourceQuery       = flag.String("source-query", "", "source query")
		learningGoal      = flag.String("learning-goal", "", "confirmed learning goal")
		currentLevel      = flag.String("current-level", "", "optional current learner level")
		casesPath         = flag.String("cases", "", "optional prompt preview cases JSON path")
		outPath           = flag.String("out", defaultOutputPath(), "output JSON report path for -cases mode")
		language          = flag.String("language", "ko", "learning language")
		top               = flag.Int("top", 3, "number of embedding matches")
		guidanceLimit     = flag.Int("guidance-limit", 2, "number of matches to inject into guidance")
		provider          = flag.String("provider", "embedding_gemma", "pattern embedding provider")
		model             = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
		endpoint          = flag.String("endpoint", defaultString("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma HTTP endpoint")
		expectedDimension = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "expected embedding dimension")
		jsonOutput        = flag.Bool("json", false, "print machine-readable JSON")
	)
	flag.Parse()

	if strings.TrimSpace(*casesPath) == "" && strings.TrimSpace(*sourceQuery) == "" {
		log.Fatal("-source-query is required")
	}
	if strings.TrimSpace(*casesPath) == "" && strings.TrimSpace(*learningGoal) == "" {
		log.Fatal("-learning-goal is required")
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

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	client, err := embedder.NewProviderClient(strings.TrimSpace(*provider), strings.TrimSpace(*endpoint))
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	if strings.TrimSpace(*casesPath) != "" {
		cases, err := loadCases(*casesPath)
		if err != nil {
			log.Fatalf("load cases: %v", err)
		}
		report := promptPreviewReport{
			GeneratedAt:   time.Now().Format(time.RFC3339),
			Provider:      strings.TrimSpace(*provider),
			Model:         strings.TrimSpace(*model),
			Dimension:     *expectedDimension,
			TopN:          *top,
			GuidanceLimit: *guidanceLimit,
			TotalCases:    len(cases),
			Results:       make([]promptPreviewCaseResult, 0, len(cases)),
		}
		repo := curriculum.NewRepository(pool)
		for _, item := range cases {
			result := previewCase(ctx, repo, client, item, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, *top, *guidanceLimit)
			report.Results = append(report.Results, result)
			if result.SkippedReason != "" {
				report.SkippedCases++
				continue
			}
			report.EvaluatedCases++
			if result.PromptChanged {
				report.PromptChangedCases++
			}
			if result.HasCrossDomainGuidance {
				report.CrossDomainGuidance++
			}
			if result.TopNMatched {
				report.TopNMatchedCases++
			}
			if result.CodeRank >= 2 {
				report.CodeRank2PlusCases++
			}
			if result.CodeRank == 0 {
				report.CodeUnmatchedCases++
			}
			if result.UsedFallbackGuidance {
				report.FallbackGuidance++
			}
		}
		if err := writeReport(*outPath, report); err != nil {
			log.Fatalf("write report: %v", err)
		}
		log.Printf("curriculum pattern prompt preview completed cases=%d evaluated=%d cross_domain_guidance=%d out=%s", report.TotalCases, report.EvaluatedCases, report.CrossDomainGuidance, *outPath)
		return
	}

	goal := strings.TrimSpace(*learningGoal)
	req := curriculum.CreateCourseDraftRequest{
		SourceQuery:      strings.TrimSpace(*sourceQuery),
		LearningGoal:     &goal,
		LearningLanguage: normalizeLanguage(*language),
	}
	if strings.TrimSpace(*currentLevel) != "" {
		level := strings.TrimSpace(*currentLevel)
		req.CurrentLevel = &level
	}

	repo := curriculum.NewRepository(pool)
	output, err := buildPreview(ctx, repo, client, req, strings.TrimSpace(*provider), strings.TrimSpace(*model), *expectedDimension, *top, *guidanceLimit)
	if err != nil {
		log.Fatal(err)
	}

	if *jsonOutput {
		payload, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("marshal output: %v", err)
		}
		fmt.Println(string(payload))
		return
	}

	printTextOutput(output)
}

func buildPreview(ctx context.Context, repo *curriculum.Repository, client embedder.Embedder, req curriculum.CreateCourseDraftRequest, provider, model string, expectedDimension, top, guidanceLimit int) (promptPreviewOutput, error) {
	query := curriculum.BuildCurriculumPatternQueryText(req)
	embedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	vector, err := client.Embed(embedCtx, query)
	cancel()
	if err != nil {
		return promptPreviewOutput{}, fmt.Errorf("embed query: %w", err)
	}
	if len(vector) != expectedDimension {
		return promptPreviewOutput{}, fmt.Errorf("unexpected query embedding dimension: got=%d want=%d", len(vector), expectedDimension)
	}

	language := normalizeLanguage(req.LearningLanguage)
	matches, err := repo.SearchCurriculumPatternMatches(ctx, curriculum.FormatFloat32VectorForPG(vector), provider, model, expectedDimension, language, top)
	if err != nil {
		return promptPreviewOutput{}, fmt.Errorf("search pattern matches: %w", err)
	}
	goal := strings.TrimSpace(derefString(req.LearningGoal))
	codeRouting := curriculum.AnalyzeCurriculumPatternRoutingPreview(req.SourceQuery, goal, req.LearningLanguage)
	codeRank := findPatternRank(matches, codeRouting.SubpatternKey)
	guidanceMatches := curriculum.FilterCurriculumPatternGuidanceMatches(matches, codeRouting.DomainAxis, codeRouting.GoalType, codeRouting.SubpatternKey, guidanceLimit, 0.30)
	guidance := curriculum.BuildCurriculumPatternMatchGuidanceForLanguage(guidanceMatches, guidanceLimit, language)
	codeOnlyPrompt := curriculum.BuildCurriculumGenerationPromptPreview(req)
	req.PatternMatchGuidance = guidance
	embeddingGuidedPrompt := curriculum.BuildCurriculumGenerationPromptPreview(req)

	output := promptPreviewOutput{
		SourceQuery:            req.SourceQuery,
		LearningGoal:           goal,
		Language:               language,
		Provider:               provider,
		Model:                  model,
		Dimension:              len(vector),
		CodeRouting:            codeRouting,
		EmbeddingTopN:          matches,
		CodeRank:               codeRank,
		TopNMatched:            codeRank > 0,
		PatternMatchGuidance:   guidance,
		CodeOnlyPrompt:         codeOnlyPrompt,
		EmbeddingGuidedPrompt:  embeddingGuidedPrompt,
		PromptChanged:          codeOnlyPrompt != embeddingGuidedPrompt,
		GuidanceMatchCount:     len(guidanceMatches),
		HasCrossDomainGuidance: hasCrossDomainGuidance(guidanceMatches, codeRouting.DomainAxis, codeRouting.GoalType),
		UsedFallbackGuidance:   usedFallbackGuidance(guidanceMatches, codeRouting.DomainAxis, codeRouting.GoalType),
	}
	if len(matches) == 0 {
		output.FallbackReason = "ready pattern matches not found"
	}
	return output, nil
}

func previewCase(ctx context.Context, repo *curriculum.Repository, client embedder.Embedder, item promptPreviewCase, provider, model string, dimension, top, guidanceLimit int) promptPreviewCaseResult {
	goal := strings.TrimSpace(item.LearningGoal)
	req := curriculum.CreateCourseDraftRequest{
		SourceQuery:      strings.TrimSpace(item.SourceQuery),
		LearningGoal:     &goal,
		LearningLanguage: normalizeLanguage(item.Language),
	}
	if strings.TrimSpace(item.CurrentLevel) != "" {
		level := strings.TrimSpace(item.CurrentLevel)
		req.CurrentLevel = &level
	}
	result := promptPreviewCaseResult{
		ID:           strings.TrimSpace(item.ID),
		SourceQuery:  req.SourceQuery,
		LearningGoal: goal,
		Language:     req.LearningLanguage,
	}
	if req.SourceQuery == "" || goal == "" {
		result.SkippedReason = "source_query or learning_goal is empty"
		return result
	}
	output, err := buildPreview(ctx, repo, client, req, provider, model, dimension, top, guidanceLimit)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.CodeRouting = output.CodeRouting
	result.EmbeddingTopN = output.EmbeddingTopN
	result.CodeRank = output.CodeRank
	result.TopNMatched = output.TopNMatched
	result.PatternMatchGuidance = output.PatternMatchGuidance
	result.PromptChanged = output.PromptChanged
	result.HasCrossDomainGuidance = output.HasCrossDomainGuidance
	result.UsedFallbackGuidance = output.UsedFallbackGuidance
	result.GuidanceMatchKeys = extractGuidanceMatchKeys(output.PatternMatchGuidance, output.EmbeddingTopN)
	return result
}

func printTextOutput(output promptPreviewOutput) {
	fmt.Printf("source_query: %s\n", output.SourceQuery)
	fmt.Printf("learning_goal: %s\n", output.LearningGoal)
	fmt.Printf("language: %s provider=%s model=%s dimension=%d\n", output.Language, output.Provider, output.Model, output.Dimension)
	fmt.Printf("code_routing: pattern=%s subpattern=%s mode=%s\n", output.CodeRouting.PatternKey, emptyAsNone(output.CodeRouting.SubpatternKey), emptyAsNone(output.CodeRouting.RefinementMode))
	fmt.Printf("code_rank: %d top_n_matched=%t\n", output.CodeRank, output.TopNMatched)
	fmt.Println("embedding_top_n:")
	for idx, match := range output.EmbeddingTopN {
		fmt.Printf("%d. score=%.4f key=%s title=%s\n", idx+1, match.SimilarityScore, match.PatternKey, match.Title)
	}
	fmt.Printf("prompt_changed: %t\n", output.PromptChanged)
	fmt.Println("\npattern_match_guidance:")
	fmt.Println(output.PatternMatchGuidance)
}

func loadCases(path string) ([]promptPreviewCase, error) {
	payload, err := os.ReadFile(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}
	var cases []promptPreviewCase
	if err := json.Unmarshal(payload, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}

func writeReport(path string, report promptPreviewReport) error {
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

func hasCrossDomainGuidance(matches []curriculum.CurriculumPatternMatch, domain, goalType string) bool {
	if strings.TrimSpace(domain) == "" || strings.TrimSpace(goalType) == "" {
		return false
	}
	for _, match := range matches {
		if strings.TrimSpace(match.Domain) != strings.TrimSpace(domain) || strings.TrimSpace(match.GoalType) != strings.TrimSpace(goalType) {
			return true
		}
	}
	return false
}

func usedFallbackGuidance(matches []curriculum.CurriculumPatternMatch, domain, goalType string) bool {
	return len(matches) > 0 && (strings.TrimSpace(domain) == "" || strings.TrimSpace(goalType) == "")
}

func extractGuidanceMatchKeys(guidance string, matches []curriculum.CurriculumPatternMatch) []string {
	if strings.TrimSpace(guidance) == "" {
		return nil
	}
	keys := []string{}
	for _, match := range matches {
		if strings.Contains(guidance, match.PatternKey) {
			keys = append(keys, match.PatternKey)
		}
	}
	return keys
}

func findPatternRank(matches []curriculum.CurriculumPatternMatch, patternKey string) int {
	patternKey = strings.TrimSpace(patternKey)
	if patternKey == "" {
		return 0
	}
	for idx, match := range matches {
		if strings.TrimSpace(match.PatternKey) == patternKey {
			return idx + 1
		}
	}
	return 0
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizeLanguage(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func emptyAsNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}

func defaultString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func defaultInt(key string, fallback int) int {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func defaultOutputPath() string {
	return filepath.Clean("../docs/data/curriculum-patterns/evaluation_curriculum_pattern_prompt_v1.json")
}
