package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/learnweaver/backend/internal/domain/curriculum"
)

type inputReport struct {
	Results []inputCourseResult `json:"results"`
}

type inputCourseResult struct {
	ID           string        `json:"id"`
	SourceQuery  string        `json:"source_query"`
	LearningGoal string        `json:"learning_goal"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Lessons      []inputLesson `json:"lessons"`
	Error        string        `json:"error,omitempty"`
}

type inputLesson struct {
	Title     string `json:"title"`
	Objective string `json:"objective"`
}

type qaReport struct {
	GeneratedAt        string     `json:"generated_at"`
	InputPath          string     `json:"input_path"`
	TotalCourses       int        `json:"total_courses"`
	EvaluatedCourses   int        `json:"evaluated_courses"`
	EvaluatedLessons   int        `json:"evaluated_lessons"`
	KoreanLeakCases    int        `json:"korean_leak_cases"`
	WeakQueryCases     int        `json:"weak_query_cases"`
	MissingIntentCases int        `json:"missing_intent_cases"`
	ErrorCases         int        `json:"error_cases"`
	Results            []qaResult `json:"results"`
}

type qaResult struct {
	ID            string           `json:"id"`
	SourceQuery   string           `json:"source_query"`
	LearningGoal  string           `json:"learning_goal"`
	SubpatternKey string           `json:"subpattern_key"`
	QueryHint     string           `json:"query_hint"`
	Lessons       []qaLessonResult `json:"lessons,omitempty"`
	Error         string           `json:"error,omitempty"`
}

type qaLessonResult struct {
	Title             string   `json:"title"`
	Objective         string   `json:"objective"`
	ContextQuery      string   `json:"context_query"`
	PrimaryQuery      string   `json:"primary_query"`
	ExternalQuery     string   `json:"external_query"`
	HasKorean         bool     `json:"has_korean"`
	WeakQuery         bool     `json:"weak_query"`
	MissingIntent     bool     `json:"missing_intent"`
	DiagnosticReasons []string `json:"diagnostic_reasons,omitempty"`
}

var koreanRe = regexp.MustCompile(`[가-힣]`)

func main() {
	inPath := flag.String("in", "../docs/data/curriculum-patterns/evaluation_curriculum_actual_generation_en_v3.json", "actual generation QA report JSON")
	outPath := flag.String("out", "../docs/data/curriculum-patterns/evaluation_curriculum_recommendation_context_en_v1.json", "output recommendation context QA report JSON")
	maxLessons := flag.Int("max-lessons", 3, "max lessons per course to evaluate; <=0 means all")
	flag.Parse()

	report, err := loadInputReport(*inPath)
	if err != nil {
		log.Fatalf("load input report: %v", err)
	}

	output := qaReport{
		GeneratedAt:  time.Now().Format(time.RFC3339),
		InputPath:    strings.TrimSpace(*inPath),
		TotalCourses: len(report.Results),
		Results:      make([]qaResult, 0, len(report.Results)),
	}

	for _, course := range report.Results {
		result := evaluateCourse(course, *maxLessons)
		output.Results = append(output.Results, result)
		if result.Error != "" {
			output.ErrorCases++
			continue
		}
		output.EvaluatedCourses++
		for _, lesson := range result.Lessons {
			output.EvaluatedLessons++
			if lesson.HasKorean {
				output.KoreanLeakCases++
			}
			if lesson.WeakQuery {
				output.WeakQueryCases++
			}
			if lesson.MissingIntent {
				output.MissingIntentCases++
			}
		}
	}

	payload, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("marshal report: %v", err)
	}
	payload = append(payload, '\n')
	if err := os.WriteFile(strings.TrimSpace(*outPath), payload, 0o644); err != nil {
		log.Fatalf("write report: %v", err)
	}
	log.Printf(
		"english recommendation context QA completed courses=%d lessons=%d korean_leak=%d weak_query=%d missing_intent=%d errors=%d out=%s",
		output.EvaluatedCourses,
		output.EvaluatedLessons,
		output.KoreanLeakCases,
		output.WeakQueryCases,
		output.MissingIntentCases,
		output.ErrorCases,
		strings.TrimSpace(*outPath),
	)
}

func loadInputReport(path string) (inputReport, error) {
	payload, err := os.ReadFile(strings.TrimSpace(path))
	if err != nil {
		return inputReport{}, err
	}
	var report inputReport
	if err := json.Unmarshal(payload, &report); err != nil {
		return inputReport{}, err
	}
	return report, nil
}

func evaluateCourse(course inputCourseResult, maxLessons int) qaResult {
	goal := strings.TrimSpace(course.LearningGoal)
	req := curriculum.CreateCourseDraftRequest{
		SourceQuery:      strings.TrimSpace(course.SourceQuery),
		LearningGoal:     &goal,
		LearningLanguage: "en",
	}
	result := qaResult{
		ID:           strings.TrimSpace(course.ID),
		SourceQuery:  req.SourceQuery,
		LearningGoal: goal,
	}
	if strings.TrimSpace(course.Error) != "" {
		result.Error = course.Error
		return result
	}
	if req.SourceQuery == "" || goal == "" {
		result.Error = "source_query or learning_goal is empty"
		return result
	}
	meta := curriculum.BuildRecommendationGoalMetadata(req)
	result.SubpatternKey = meta.SubpatternKey
	result.QueryHint = meta.QueryHint

	lessons := course.Lessons
	if maxLessons > 0 && len(lessons) > maxLessons {
		lessons = lessons[:maxLessons]
	}
	for _, lesson := range lessons {
		result.Lessons = append(result.Lessons, evaluateLesson(req, meta, lesson))
	}
	if len(result.Lessons) == 0 {
		result.Error = "no lessons to evaluate"
	}
	return result
}

func evaluateLesson(req curriculum.CreateCourseDraftRequest, meta curriculum.RecommendationGoalMetadata, lesson inputLesson) qaLessonResult {
	lessonTitle := strings.TrimSpace(lesson.Title)
	lessonObjective := strings.TrimSpace(lesson.Objective)
	context := curriculum.ExplorerRecommendationContext{
		SourceQuery:     strings.TrimSpace(req.SourceQuery),
		LearningGoal:    strings.TrimSpace(derefString(req.LearningGoal)),
		GoalQueryHint:   strings.TrimSpace(meta.QueryHint),
		RegionTitle:     lessonTitle,
		RegionObjective: lessonObjective,
		NodeTitle:       lessonTitle,
		NodeSummary:     lessonObjective,
		FallbackQuery:   lessonTitle,
	}
	contextQuery := curriculum.BuildExplorerRecommendationExternalSearchQuery(context)
	primaryQuery := buildEnglishPrimaryRecommendationQuery(req, lessonTitle, lessonObjective)
	externalQuery := curriculum.BuildExplorerRecommendationExternalSearchQueryWithPrimaryQuery(context, primaryQuery)

	result := qaLessonResult{
		Title:         lessonTitle,
		Objective:     lessonObjective,
		ContextQuery:  contextQuery,
		PrimaryQuery:  primaryQuery,
		ExternalQuery: externalQuery,
	}
	joined := strings.Join([]string{contextQuery, primaryQuery, externalQuery}, " ")
	result.HasKorean = koreanRe.MatchString(joined)
	if result.HasKorean {
		result.DiagnosticReasons = append(result.DiagnosticReasons, "query contains Korean text")
	}
	result.WeakQuery = isWeakRecommendationQuery(externalQuery)
	if result.WeakQuery {
		result.DiagnosticReasons = append(result.DiagnosticReasons, "query has too few useful tokens")
	}
	result.MissingIntent = !hasEnglishSearchIntent(externalQuery)
	if result.MissingIntent {
		result.DiagnosticReasons = append(result.DiagnosticReasons, "query is missing English search intent")
	}
	return result
}

func buildEnglishPrimaryRecommendationQuery(req curriculum.CreateCourseDraftRequest, lessonTitle, lessonObjective string) string {
	words := tokenizeRecommendationQuery(strings.Join([]string{
		strings.TrimSpace(req.SourceQuery),
		strings.TrimSpace(derefString(req.LearningGoal)),
		lessonTitle,
		lessonObjective,
	}, " "))
	if len(words) > 7 {
		words = words[:7]
	}
	if len(words) == 0 {
		return "tutorial"
	}
	if !hasEnglishSearchIntent(strings.Join(words, " ")) {
		words = append(words, "tutorial")
	}
	return strings.Join(words, " ")
}

func tokenizeRecommendationQuery(text string) []string {
	replacer := strings.NewReplacer(
		"'", " ",
		"\"", " ",
		",", " ",
		".", " ",
		":", " ",
		";", " ",
		"(", " ",
		")", " ",
		"/", " ",
		"-", " ",
	)
	cleaned := strings.ToLower(replacer.Replace(text))
	stopWords := map[string]struct{}{
		"a": {}, "an": {}, "and": {}, "are": {}, "as": {}, "be": {}, "by": {}, "can": {}, "for": {}, "from": {}, "i": {},
		"in": {}, "into": {}, "it": {}, "learn": {}, "learning": {}, "lesson": {}, "of": {}, "on": {}, "or": {}, "step": {},
		"the": {}, "them": {}, "this": {}, "to": {}, "want": {}, "with": {},
	}
	result := []string{}
	seen := map[string]struct{}{}
	for _, token := range strings.Fields(cleaned) {
		token = strings.TrimSpace(token)
		if len([]rune(token)) < 3 {
			continue
		}
		if _, skip := stopWords[token]; skip {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		result = append(result, token)
		if len(result) >= 8 {
			break
		}
	}
	return result
}

func isWeakRecommendationQuery(query string) bool {
	useful := 0
	for _, token := range strings.Fields(strings.ToLower(query)) {
		if len([]rune(strings.TrimSpace(token))) >= 3 {
			useful++
		}
	}
	return useful < 3
}

func hasEnglishSearchIntent(query string) bool {
	lower := strings.ToLower(query)
	return strings.Contains(lower, "tutorial") ||
		strings.Contains(lower, "guide") ||
		strings.Contains(lower, "practice") ||
		strings.Contains(lower, "project") ||
		strings.Contains(lower, "how")
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
