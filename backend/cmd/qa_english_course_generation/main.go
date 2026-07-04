package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/llm"
	"github.com/learnweaver/backend/internal/pkg/systemsettings"
)

type qaCase struct {
	ID           string `json:"id"`
	SourceQuery  string `json:"source_query"`
	LearningGoal string `json:"learning_goal"`
	Language     string `json:"language"`
}

type qaReport struct {
	GeneratedAt       string     `json:"generated_at"`
	Provider          string     `json:"provider"`
	Model             string     `json:"model"`
	TotalCases        int        `json:"total_cases"`
	EvaluatedCases    int        `json:"evaluated_cases"`
	KoreanLeakCases   int        `json:"korean_leak_cases"`
	GenericIssueCases int        `json:"generic_issue_cases"`
	ErrorCases        int        `json:"error_cases"`
	Results           []qaResult `json:"results"`
}

type qaResult struct {
	ID                   string                       `json:"id"`
	SourceQuery          string                       `json:"source_query"`
	LearningGoal         string                       `json:"learning_goal"`
	CodeSubpattern       string                       `json:"code_subpattern"`
	GuidanceMatchKeys    []string                     `json:"guidance_match_keys"`
	Title                string                       `json:"title,omitempty"`
	Description          string                       `json:"description,omitempty"`
	Lessons              []qaLesson                   `json:"lessons,omitempty"`
	CompletionCriteria   []string                     `json:"completion_criteria,omitempty"`
	HasKorean            bool                         `json:"has_korean"`
	GenericIssues        []string                     `json:"generic_issues,omitempty"`
	GenerationDurationMS int64                        `json:"generation_duration_ms,omitempty"`
	PatternGuidance      string                       `json:"pattern_guidance,omitempty"`
	EmbeddingTopN        []curriculumPatternMatchView `json:"embedding_top_n,omitempty"`
	Error                string                       `json:"error,omitempty"`
}

type qaLesson struct {
	Title     string `json:"title"`
	Objective string `json:"objective"`
}

type curriculumPatternMatchView struct {
	PatternKey      string  `json:"pattern_key"`
	Title           string  `json:"title"`
	SimilarityScore float64 `json:"similarity_score"`
}

var koreanRe = regexp.MustCompile(`[가-힣]`)

func main() {
	_ = godotenv.Load()

	casesPath := flag.String("cases", "../docs/data/curriculum-patterns/evaluation_curriculum_pattern_prompt_en_v1_cases.json", "QA cases JSON path")
	outPath := flag.String("out", "../docs/data/curriculum-patterns/evaluation_curriculum_actual_generation_en_v1.json", "output report JSON path")
	limit := flag.Int("limit", 6, "max cases to run")
	provider := flag.String("provider", "", "override LLM provider")
	model := flag.String("model", "", "override LLM model")
	feature := flag.String("feature", "default", "system_llm_settings feature")
	embedProvider := flag.String("embed-provider", "embedding_gemma", "pattern embedding provider")
	embedModel := flag.String("embed-model", envDefault("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
	embedEndpoint := flag.String("embed-endpoint", envDefault("EMBEDDING_GEMMA_ENDPOINT", ""), "EmbeddingGemma endpoint")
	dimension := flag.Int("dimension", envDefaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "embedding dimension")
	flag.Parse()

	ctx := context.Background()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if strings.TrimSpace(*embedEndpoint) == "" {
		log.Fatal("EMBEDDING_GEMMA_ENDPOINT is required")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	llmProvider, llmModel := strings.TrimSpace(*provider), strings.TrimSpace(*model)
	if llmProvider == "" || llmModel == "" {
		dbProvider, dbModel, err := goal.NewRepository(pool).GetSystemLLMSetting(ctx, strings.TrimSpace(*feature))
		if err != nil {
			log.Fatalf("system llm setting unavailable: %v", err)
		}
		if llmProvider == "" {
			llmProvider = strings.TrimSpace(dbProvider)
		}
		if llmModel == "" {
			llmModel = strings.TrimSpace(dbModel)
		}
	}

	apiKey, err := systemsettings.NewStore(pool).ResolveAPIKey(ctx, llmProvider, os.Getenv("OPENAI_API_KEY"))
	if err != nil || strings.TrimSpace(apiKey) == "" {
		log.Fatalf("system llm api key unavailable provider=%s err=%v", llmProvider, err)
	}
	llmClient, err := llm.NewClient(strings.ToLower(llmProvider), apiKey, llmModel)
	if err != nil {
		log.Fatalf("llm client: %v", err)
	}

	embedClient, err := embedder.NewProviderClient(strings.TrimSpace(*embedProvider), strings.TrimSpace(*embedEndpoint))
	if err != nil {
		log.Fatalf("embedder: %v", err)
	}

	cases, err := loadCases(*casesPath)
	if err != nil {
		log.Fatalf("load cases: %v", err)
	}
	if *limit > 0 && *limit < len(cases) {
		cases = cases[:*limit]
	}

	repo := curriculum.NewRepository(pool)
	service := curriculum.NewService()
	report := qaReport{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Provider:    llmProvider,
		Model:       llmModel,
		TotalCases:  len(cases),
		Results:     make([]qaResult, 0, len(cases)),
	}

	for _, item := range cases {
		result := runCase(ctx, repo, service, embedClient, llmClient, item, strings.TrimSpace(*embedProvider), strings.TrimSpace(*embedModel), *dimension)
		report.Results = append(report.Results, result)
		if result.Error != "" {
			report.ErrorCases++
			continue
		}
		report.EvaluatedCases++
		if result.HasKorean {
			report.KoreanLeakCases++
		}
		if len(result.GenericIssues) > 0 {
			report.GenericIssueCases++
		}
	}

	payload, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("marshal report: %v", err)
	}
	payload = append(payload, '\n')
	if err := os.WriteFile(strings.TrimSpace(*outPath), payload, 0o644); err != nil {
		log.Fatalf("write report: %v", err)
	}
	log.Printf("english actual course QA completed evaluated=%d korean_leak=%d generic_issue=%d errors=%d out=%s", report.EvaluatedCases, report.KoreanLeakCases, report.GenericIssueCases, report.ErrorCases, *outPath)
}

func runCase(ctx context.Context, repo *curriculum.Repository, service *curriculum.Service, embedClient embedder.Embedder, llmClient llm.Client, item qaCase, embedProvider, embedModel string, dimension int) qaResult {
	goalText := strings.TrimSpace(item.LearningGoal)
	req := curriculum.CreateCourseDraftRequest{
		SourceQuery:      strings.TrimSpace(item.SourceQuery),
		LearningGoal:     &goalText,
		LearningLanguage: "en",
	}
	result := qaResult{
		ID:           strings.TrimSpace(item.ID),
		SourceQuery:  req.SourceQuery,
		LearningGoal: goalText,
	}
	if req.SourceQuery == "" || goalText == "" {
		result.Error = "source_query or learning_goal is empty"
		return result
	}

	query := curriculum.BuildCurriculumPatternQueryText(req)
	vector, err := embedClient.Embed(ctx, query)
	if err != nil {
		result.Error = "embed query: " + err.Error()
		return result
	}
	if len(vector) != dimension {
		result.Error = fmt.Sprintf("unexpected embedding dimension: got=%d want=%d", len(vector), dimension)
		return result
	}
	matches, err := repo.SearchCurriculumPatternMatches(ctx, curriculum.FormatFloat32VectorForPG(vector), embedProvider, embedModel, dimension, "en", 3)
	if err != nil {
		result.Error = "search pattern matches: " + err.Error()
		return result
	}
	routing := curriculum.AnalyzeCurriculumPatternRoutingPreview(req.SourceQuery, goalText, "en")
	result.CodeSubpattern = routing.SubpatternKey
	for _, match := range matches {
		result.EmbeddingTopN = append(result.EmbeddingTopN, curriculumPatternMatchView{
			PatternKey:      match.PatternKey,
			Title:           match.Title,
			SimilarityScore: match.SimilarityScore,
		})
	}
	guidanceMatches := curriculum.FilterCurriculumPatternGuidanceMatches(matches, routing.DomainAxis, routing.GoalType, routing.SubpatternKey, 2, 0.30)
	for _, match := range guidanceMatches {
		result.GuidanceMatchKeys = append(result.GuidanceMatchKeys, match.PatternKey)
	}
	req.PatternMatchGuidance = curriculum.BuildCurriculumPatternMatchGuidanceForLanguage(guidanceMatches, 2, "en")
	result.PatternGuidance = req.PatternMatchGuidance

	startedAt := time.Now()
	llmCtx, cancel := context.WithTimeout(ctx, 75*time.Second)
	defer cancel()
	draft, err := service.BuildDraftWithLLM(llmCtx, llmClient, uuid.New(), req, curriculum.DraftBuildOptions{
		SkipReview:        true,
		GenerationTimeout: 70 * time.Second,
	})
	result.GenerationDurationMS = time.Since(startedAt).Milliseconds()
	if err != nil {
		result.Error = "generate draft: " + err.Error()
		return result
	}

	result.Title = draft.Draft.Title
	if draft.Draft.Description != nil {
		result.Description = strings.TrimSpace(*draft.Draft.Description)
	}
	result.CompletionCriteria = append([]string(nil), draft.Draft.CompletionCriteria...)
	for _, lesson := range draft.Lessons {
		item := qaLesson{Title: strings.TrimSpace(lesson.Lesson.Title)}
		if lesson.Lesson.Objective != nil {
			item.Objective = strings.TrimSpace(*lesson.Lesson.Objective)
		}
		result.Lessons = append(result.Lessons, item)
	}

	joined := strings.Join(append([]string{result.Title, result.Description}, append(result.CompletionCriteria, lessonTexts(result.Lessons)...)...), " ")
	result.HasKorean = koreanRe.MatchString(joined)
	result.GenericIssues = detectGenericIssues(result)
	return result
}

func loadCases(path string) ([]qaCase, error) {
	payload, err := os.ReadFile(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}
	var cases []qaCase
	if err := json.Unmarshal(payload, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}

func lessonTexts(lessons []qaLesson) []string {
	items := make([]string, 0, len(lessons)*2)
	for _, lesson := range lessons {
		items = append(items, lesson.Title, lesson.Objective)
	}
	return items
}

func detectGenericIssues(result qaResult) []string {
	joined := strings.ToLower(strings.Join(lessonTexts(result.Lessons), " "))
	issues := []string{}
	for _, phrase := range []string{
		"understand the basics",
		"learn the basics",
		"introduction to",
		"complete the pre-goal practice flow",
		"setup and introduction",
		"first performance output",
		"first output",
		"core pattern practice",
		"integration practice",
		"performance prep",
		"explore concepts",
		"study theory",
	} {
		if strings.Contains(joined, phrase) {
			issues = append(issues, phrase)
		}
	}
	return issues
}

func envDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envDefaultInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
