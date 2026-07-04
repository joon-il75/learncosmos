package lessonsearchprerun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
)

const (
	DefaultFixturePath = "internal/evaluation/testdata/external_search_provider_fixtures.json"
	DefaultOutputPath  = "/tmp/lesson-search-fixture-refresh-report.json"
	DefaultMode        = "refresh-existing"
	DiscoverGapMode    = "discover-gap"
)

type Options struct {
	FixturesPath string
	OutputPath   string
	Provider     string
	FixtureID    string
	Limit        int
	TopN         int
	Mode         string
	Timeout      time.Duration
}

type Fixture struct {
	ID         string                                    `json:"id"`
	Provider   string                                    `json:"provider"`
	Query      string                                    `json:"query"`
	SearchSpec curriculum.LessonRecommendationSearchSpec `json:"search_spec"`
	Candidates []FixtureCandidate                        `json:"candidates"`
}

type FixtureCandidate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Relevant    bool   `json:"relevant"`
}

type Candidate struct {
	Title       string
	Description string
	ContentType string
}

type Provider interface {
	Search(ctx context.Context, fixture Fixture, topN int) ([]Candidate, error)
}

type Report struct {
	GeneratedAt          time.Time            `json:"generated_at"`
	Mode                 string               `json:"mode"`
	SourceFixturePath    string               `json:"source_fixture_path"`
	ProviderFilter       string               `json:"provider_filter"`
	TopN                 int                  `json:"top_n"`
	Summary              ReportSummary        `json:"summary"`
	NoiseTaxonomySummary NoiseTaxonomySummary `json:"noise_taxonomy_summary,omitempty"`
	Fixtures             []ReportFixture      `json:"fixtures"`
	Errors               []ReportError        `json:"errors,omitempty"`
}

type ReportSummary struct {
	FixturesTotal          int `json:"fixtures_total"`
	FixturesProcessed      int `json:"fixtures_processed"`
	FixturesSkipped        int `json:"fixtures_skipped"`
	ProviderErrors         int `json:"provider_errors"`
	NoiseSuspects          int `json:"noise_suspects"`
	RelevantLikeCandidates int `json:"relevant_like_candidates"`
}

type ReportFixture struct {
	ID                string            `json:"id"`
	Provider          string            `json:"provider"`
	Query             string            `json:"query"`
	StageRole         string            `json:"stage_role"`
	ContentTypes      []string          `json:"content_types,omitempty"`
	SearchSpecSummary SearchSpecSummary `json:"search_spec_summary"`
	ExistingFixture   ExistingFixture   `json:"existing_fixture"`
	Candidates        []ReportCandidate `json:"candidates,omitempty"`
	Warnings          []string          `json:"warnings,omitempty"`
	Skipped           string            `json:"skipped,omitempty"`
	Error             string            `json:"error,omitempty"`
}

type SearchSpecSummary struct {
	MustInclude []string `json:"must_include,omitempty"`
	NiceToHave  []string `json:"nice_to_have,omitempty"`
	Avoid       []string `json:"avoid,omitempty"`
}

type ExistingFixture struct {
	RelevantCount int `json:"relevant_count"`
	NoiseCount    int `json:"noise_count"`
}

type ReportCandidate struct {
	Rank            int      `json:"rank"`
	Provider        string   `json:"provider"`
	ContentType     string   `json:"content_type"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	Score           int      `json:"score"`
	MustIncludeHits []string `json:"must_include_hits,omitempty"`
	NiceToHaveHits  []string `json:"nice_to_have_hits,omitempty"`
	AvoidHits       []string `json:"avoid_hits,omitempty"`
	NoiseHintHits   []string `json:"noise_hint_hits,omitempty"`
	NeedsReview     bool     `json:"needs_review"`
	Warning         string   `json:"warning,omitempty"`
}

type ReportError struct {
	FixtureID string `json:"fixture_id,omitempty"`
	Provider  string `json:"provider,omitempty"`
	Message   string `json:"message"`
}

type NoiseTaxonomySummary struct {
	RegisteredAvoidHits    []NoiseTokenCount `json:"registered_avoid_hits,omitempty"`
	UnregisteredNoiseHints []NoiseTokenCount `json:"unregistered_noise_hints,omitempty"`
}

type NoiseTokenCount struct {
	Token string `json:"token"`
	Count int    `json:"count"`
}

type GapReport struct {
	GeneratedAt       time.Time      `json:"generated_at"`
	Mode              string         `json:"mode"`
	SourceFixturePath string         `json:"source_fixture_path"`
	ProviderFilter    string         `json:"provider_filter"`
	Summary           GapSummary     `json:"summary"`
	Coverage          GapCoverage    `json:"coverage"`
	GapCandidates     []GapCandidate `json:"gap_candidates"`
}

type GapSummary struct {
	FixturesTotal int            `json:"fixtures_total"`
	Providers     map[string]int `json:"providers"`
	StageRoles    map[string]int `json:"stage_roles"`
	ContentTypes  map[string]int `json:"content_types"`
	GapCount      int            `json:"gap_count"`
}

type GapCoverage struct {
	ByProviderStage []ProviderStageCoverage `json:"by_provider_stage"`
}

type ProviderStageCoverage struct {
	Provider  string `json:"provider"`
	StageRole string `json:"stage_role"`
	Count     int    `json:"count"`
}

type GapCandidate struct {
	Priority         string `json:"priority"`
	GapType          string `json:"gap_type"`
	Provider         string `json:"provider,omitempty"`
	StageRole        string `json:"stage_role,omitempty"`
	ContentType      string `json:"content_type,omitempty"`
	Reason           string `json:"reason"`
	SuggestedAction  string `json:"suggested_action"`
	FixtureAutoApply bool   `json:"fixture_auto_apply"`
}

type providerNotWiredError struct {
	provider string
}

func (err providerNotWiredError) Error() string {
	return "provider_not_wired"
}

type providerMissingCredentialsError struct {
	provider string
}

func (err providerMissingCredentialsError) Error() string {
	return "missing_credentials"
}

type missingCredentialsProvider struct {
	provider string
}

func (provider missingCredentialsProvider) Search(context.Context, Fixture, int) ([]Candidate, error) {
	return nil, providerMissingCredentialsError{provider: provider.provider}
}

type extsearchProviderAdapter struct {
	provider extsearch.SearchProvider
}

func (adapter extsearchProviderAdapter) Search(ctx context.Context, fixture Fixture, topN int) ([]Candidate, error) {
	results, err := adapter.provider.Search(ctx, strings.TrimSpace(fixture.Query), topN)
	if err != nil {
		return nil, err
	}
	candidates := make([]Candidate, 0, len(results))
	for _, result := range results {
		candidates = append(candidates, Candidate{
			Title:       result.Title,
			Description: result.Description,
			ContentType: result.Source,
		})
	}
	return candidates, nil
}

func Run(ctx context.Context, opts Options, providers map[string]Provider) error {
	opts = NormalizeOptions(opts)
	if err := ValidateOptions(opts); err != nil {
		return err
	}
	fixtures, err := LoadFixtures(opts.FixturesPath)
	if err != nil {
		return err
	}

	if opts.Mode == DiscoverGapMode {
		result := BuildGapReport(opts, fixtures, time.Now().UTC())
		if err := WriteReport(opts.OutputPath, result); err != nil {
			return err
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()
	result := BuildReport(ctx, opts, fixtures, providers, time.Now().UTC())
	if err := WriteReport(opts.OutputPath, result); err != nil {
		return err
	}
	return nil
}

func NormalizeOptions(opts Options) Options {
	opts.Provider = strings.ToLower(strings.TrimSpace(opts.Provider))
	if opts.Provider == "" {
		opts.Provider = "all"
	}
	opts.Mode = strings.ToLower(strings.TrimSpace(opts.Mode))
	if opts.Mode == "" {
		opts.Mode = DefaultMode
	}
	if opts.TopN <= 0 {
		opts.TopN = 5
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 10 * time.Second
	}
	if strings.TrimSpace(opts.FixturesPath) == "" {
		opts.FixturesPath = DefaultFixturePath
	}
	if strings.TrimSpace(opts.OutputPath) == "" {
		opts.OutputPath = DefaultOutputPath
	}
	return opts
}

func ValidateOptions(opts Options) error {
	switch opts.Mode {
	case DefaultMode, DiscoverGapMode:
	default:
		return fmt.Errorf("unsupported -mode %q", opts.Mode)
	}
	switch opts.Provider {
	case "all", "youtube", "naver_blog":
	default:
		return fmt.Errorf("unsupported -provider %q", opts.Provider)
	}
	if opts.Limit < 0 {
		return errors.New("-limit must be greater than or equal to 0")
	}
	return nil
}

func LoadFixtures(path string) ([]Fixture, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fixtures: %w", err)
	}
	var fixtures []Fixture
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		return nil, fmt.Errorf("parse fixtures: %w", err)
	}
	return fixtures, nil
}

func BuildReport(ctx context.Context, opts Options, fixtures []Fixture, providers map[string]Provider, now time.Time) Report {
	result := Report{
		GeneratedAt:       now,
		Mode:              opts.Mode,
		SourceFixturePath: opts.FixturesPath,
		ProviderFilter:    opts.Provider,
		TopN:              opts.TopN,
		Summary: ReportSummary{
			FixturesTotal: len(fixtures),
		},
	}
	registeredAvoidHits := map[string]int{}
	unregisteredNoiseHints := map[string]int{}

	for _, fixture := range fixtures {
		if !matchesFilter(fixture, opts) {
			result.Summary.FixturesSkipped++
			continue
		}
		if opts.Limit > 0 && result.Summary.FixturesProcessed >= opts.Limit {
			result.Summary.FixturesSkipped++
			continue
		}
		result.Summary.FixturesProcessed++
		ReportFixture := buildFixtureReport(ctx, fixture, opts, providers)
		if ReportFixture.Skipped != "" {
			result.Summary.FixturesSkipped++
		}
		if ReportFixture.Error != "" {
			result.Summary.ProviderErrors++
			result.Errors = append(result.Errors, ReportError{
				FixtureID: fixture.ID,
				Provider:  fixture.Provider,
				Message:   ReportFixture.Error,
			})
		}
		for _, candidate := range ReportFixture.Candidates {
			if len(candidate.AvoidHits) > 0 {
				result.Summary.NoiseSuspects++
			}
			for _, token := range candidate.AvoidHits {
				incrementTokenCount(registeredAvoidHits, token)
			}
			for _, token := range candidate.NoiseHintHits {
				incrementTokenCount(unregisteredNoiseHints, token)
			}
			if !candidate.NeedsReview && candidate.Score > 0 {
				result.Summary.RelevantLikeCandidates++
			}
		}
		result.Fixtures = append(result.Fixtures, ReportFixture)
	}
	result.NoiseTaxonomySummary = buildNoiseTaxonomySummary(registeredAvoidHits, unregisteredNoiseHints)
	return result
}

func BuildGapReport(opts Options, fixtures []Fixture, now time.Time) GapReport {
	result := GapReport{
		GeneratedAt:       now,
		Mode:              DiscoverGapMode,
		SourceFixturePath: opts.FixturesPath,
		ProviderFilter:    opts.Provider,
		Summary: GapSummary{
			FixturesTotal: len(fixtures),
			Providers:     map[string]int{},
			StageRoles:    map[string]int{},
			ContentTypes:  map[string]int{},
		},
		Coverage: GapCoverage{
			ByProviderStage: []ProviderStageCoverage{},
		},
		GapCandidates: []GapCandidate{},
	}

	providerStageCounts := map[string]int{}
	processed := 0
	for _, fixture := range fixtures {
		if !matchesFilter(fixture, opts) {
			continue
		}
		if opts.Limit > 0 && processed >= opts.Limit {
			continue
		}
		processed++

		spec := curriculum.NormalizeLessonRecommendationSearchSpec(fixture.SearchSpec, curriculum.LessonRecommendationSearchSpec{})
		provider := normalizeGapValue(fixture.Provider, "unknown")
		stageRole := normalizeGapValue(spec.StageRole, "unknown")
		result.Summary.Providers[provider]++
		result.Summary.StageRoles[stageRole]++
		for _, contentType := range normalizedContentTypes(spec.ContentTypes) {
			result.Summary.ContentTypes[contentType]++
		}
		providerStageCounts[providerStageKey(provider, stageRole)]++
	}

	providers := gapProviders(opts, result.Summary.Providers)
	stages := coreStageRoles()
	for _, provider := range providers {
		for _, stage := range stages {
			result.Coverage.ByProviderStage = append(result.Coverage.ByProviderStage, ProviderStageCoverage{
				Provider:  provider,
				StageRole: stage,
				Count:     providerStageCounts[providerStageKey(provider, stage)],
			})
		}
	}

	for _, stage := range stages {
		if result.Summary.StageRoles[stage] == 0 {
			result.GapCandidates = append(result.GapCandidates, GapCandidate{
				Priority:         "high",
				GapType:          "stage_role",
				StageRole:        stage,
				Reason:           "core stage role has no provider fixture",
				SuggestedAction:  "manual sample collection",
				FixtureAutoApply: false,
			})
		}
	}
	for _, provider := range providers {
		for _, stage := range stages {
			if result.Summary.StageRoles[stage] > 0 && providerStageCounts[providerStageKey(provider, stage)] == 0 {
				result.GapCandidates = append(result.GapCandidates, GapCandidate{
					Priority:         "medium",
					GapType:          "provider_stage",
					Provider:         provider,
					StageRole:        stage,
					Reason:           "stage is covered overall but not by this provider",
					SuggestedAction:  "manual sample collection",
					FixtureAutoApply: false,
				})
			}
		}
	}
	for _, contentType := range []string{"video", "article", "article+video"} {
		if result.Summary.ContentTypes[contentType] == 0 {
			result.GapCandidates = append(result.GapCandidates, GapCandidate{
				Priority:         "low",
				GapType:          "content_type",
				ContentType:      contentType,
				Reason:           "content type has no provider fixture",
				SuggestedAction:  "manual sample collection",
				FixtureAutoApply: false,
			})
		}
	}
	result.Summary.GapCount = len(result.GapCandidates)
	return result
}

func coreStageRoles() []string {
	return []string{
		"setup_intro",
		"first_output",
		"core_pattern",
		"practice_loop",
		"artifact_finish",
		"troubleshooting",
		"portfolio_publish",
	}
}

func gapProviders(opts Options, counts map[string]int) []string {
	if opts.Provider != "all" {
		return []string{opts.Provider}
	}
	values := []string{"youtube", "naver_blog"}
	seen := map[string]struct{}{"youtube": {}, "naver_blog": {}}
	for provider := range counts {
		if _, ok := seen[provider]; ok {
			continue
		}
		values = append(values, provider)
		seen[provider] = struct{}{}
	}
	sort.Strings(values[2:])
	return values
}

func normalizedContentTypes(contentTypes []string) []string {
	if len(contentTypes) == 0 {
		return []string{"unknown"}
	}
	seen := map[string]struct{}{}
	for _, contentType := range contentTypes {
		normalized := normalizeGapValue(contentType, "unknown")
		seen[normalized] = struct{}{}
	}
	if _, hasArticle := seen["article"]; hasArticle {
		if _, hasVideo := seen["video"]; hasVideo {
			return []string{"article+video"}
		}
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func normalizeGapValue(value string, fallback string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return fallback
	}
	return normalized
}

func providerStageKey(provider string, stageRole string) string {
	return provider + "\x00" + stageRole
}

func matchesFilter(fixture Fixture, opts Options) bool {
	if opts.Provider != "all" && !strings.EqualFold(fixture.Provider, opts.Provider) {
		return false
	}
	if strings.TrimSpace(opts.FixtureID) != "" && fixture.ID != strings.TrimSpace(opts.FixtureID) {
		return false
	}
	return true
}

func buildFixtureReport(ctx context.Context, fixture Fixture, opts Options, providers map[string]Provider) ReportFixture {
	spec := curriculum.NormalizeLessonRecommendationSearchSpec(fixture.SearchSpec, curriculum.LessonRecommendationSearchSpec{})
	item := ReportFixture{
		ID:           fixture.ID,
		Provider:     fixture.Provider,
		Query:        fixture.Query,
		StageRole:    spec.StageRole,
		ContentTypes: append([]string(nil), spec.ContentTypes...),
		SearchSpecSummary: SearchSpecSummary{
			MustInclude: append([]string(nil), spec.MustInclude...),
			NiceToHave:  append([]string(nil), spec.NiceToHave...),
			Avoid:       append([]string(nil), spec.Avoid...),
		},
		ExistingFixture: summarizeExistingFixture(fixture.Candidates),
	}

	provider := providers[strings.ToLower(strings.TrimSpace(fixture.Provider))]
	if provider == nil {
		item.Skipped = "skipped_provider_not_wired"
		item.Warnings = append(item.Warnings, "provider adapter is not wired in CLI skeleton")
		return item
	}

	candidates, err := provider.Search(ctx, fixture, opts.TopN)
	if err != nil {
		var missingCredentials providerMissingCredentialsError
		if errors.As(err, &missingCredentials) {
			item.Skipped = "skipped_missing_credentials"
			item.Warnings = append(item.Warnings, "provider credentials are not configured")
			return item
		}
		var notWired providerNotWiredError
		if errors.As(err, &notWired) {
			item.Skipped = "skipped_provider_not_wired"
			item.Warnings = append(item.Warnings, "provider adapter is not wired in CLI skeleton")
			return item
		}
		item.Error = sanitizeProviderError(err)
		return item
	}
	if len(candidates) == 0 {
		item.Warnings = append(item.Warnings, "empty_provider_result")
		return item
	}
	for idx, candidate := range trimCandidates(candidates, opts.TopN) {
		item.Candidates = append(item.Candidates, scoreProviderCandidate(fixture.Provider, idx+1, candidate, spec))
	}
	item.Warnings = append(item.Warnings, fixtureWarnings(item.Candidates, spec)...)
	return item
}

func summarizeExistingFixture(candidates []FixtureCandidate) ExistingFixture {
	var result ExistingFixture
	for _, candidate := range candidates {
		if candidate.Relevant {
			result.RelevantCount++
		} else {
			result.NoiseCount++
		}
	}
	return result
}

func trimCandidates(candidates []Candidate, topN int) []Candidate {
	if topN <= 0 || len(candidates) <= topN {
		return candidates
	}
	return candidates[:topN]
}

func scoreProviderCandidate(provider string, rank int, candidate Candidate, spec curriculum.LessonRecommendationSearchSpec) ReportCandidate {
	description := strings.TrimSpace(candidate.Description)
	contentType := strings.TrimSpace(candidate.ContentType)
	if contentType == "" {
		contentType = provider
	}
	curriculumCandidate := curriculum.ContentSearchCandidate{
		Title:       strings.TrimSpace(candidate.Title),
		Description: &description,
		ContentType: contentType,
		Language:    spec.Language,
	}
	score := curriculum.ScoreCandidateWithLessonSearchSpec(curriculumCandidate, spec)
	searchText := strings.ToLower(strings.Join([]string{candidate.Title, description}, " "))
	mustHits := matchingTokens(searchText, spec.MustInclude)
	niceHits := matchingTokens(searchText, spec.NiceToHave)
	avoidHits := matchingTokens(searchText, spec.Avoid)
	noiseHintHits := unregisteredNoiseHintHits(searchText, spec.Avoid)
	needsReview := len(mustHits) == 0 || len(avoidHits) > 0
	warning := ""
	if len(avoidHits) > 0 {
		warning = "avoid_terms_present"
	} else if len(mustHits) == 0 {
		warning = "missing_must_include"
	}
	return ReportCandidate{
		Rank:            rank,
		Provider:        provider,
		ContentType:     contentType,
		Title:           strings.TrimSpace(candidate.Title),
		Description:     description,
		Score:           score,
		MustIncludeHits: mustHits,
		NiceToHaveHits:  niceHits,
		AvoidHits:       avoidHits,
		NoiseHintHits:   noiseHintHits,
		NeedsReview:     needsReview,
		Warning:         warning,
	}
}

func knownNoiseHintTokens() []string {
	return []string{
		"과외", "수강신청", "합격보장", "합격 보장", "응시자격", "편입", "수학학원",
		"키트", "완제품", "판매", "공방 모집", "클래스 모집",
		"맛집", "배달", "식당 후기", "밀키트",
		"다이어트", "칼로리", "챌린지", "런닝머신",
		"부업", "조회수", "수익 인증", "전자책 판매", "대행",
		"제품 리뷰", "할인", "특가",
	}
}

func buildNoiseTaxonomySummary(registeredAvoidHits map[string]int, unregisteredNoiseHints map[string]int) NoiseTaxonomySummary {
	return NoiseTaxonomySummary{
		RegisteredAvoidHits:    tokenCountList(registeredAvoidHits),
		UnregisteredNoiseHints: tokenCountList(unregisteredNoiseHints),
	}
}

func tokenCountList(counts map[string]int) []NoiseTokenCount {
	values := make([]NoiseTokenCount, 0, len(counts))
	for token, count := range counts {
		if strings.TrimSpace(token) == "" || count <= 0 {
			continue
		}
		values = append(values, NoiseTokenCount{Token: token, Count: count})
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].Count == values[j].Count {
			return values[i].Token < values[j].Token
		}
		return values[i].Count > values[j].Count
	})
	return values
}

func incrementTokenCount(counts map[string]int, token string) {
	normalized := strings.TrimSpace(token)
	if normalized == "" {
		return
	}
	counts[normalized]++
}

func unregisteredNoiseHintHits(searchText string, avoidTokens []string) []string {
	registered := map[string]struct{}{}
	for _, token := range avoidTokens {
		normalized := strings.ToLower(strings.TrimSpace(token))
		if normalized != "" {
			registered[normalized] = struct{}{}
		}
	}
	hits := matchingTokens(searchText, knownNoiseHintTokens())
	values := make([]string, 0, len(hits))
	for _, hit := range hits {
		normalized := strings.ToLower(strings.TrimSpace(hit))
		if _, ok := registered[normalized]; ok {
			continue
		}
		values = append(values, hit)
	}
	return values
}

func matchingTokens(searchText string, tokens []string) []string {
	seen := map[string]struct{}{}
	var hits []string
	for _, token := range tokens {
		normalized := strings.ToLower(strings.TrimSpace(token))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		if strings.Contains(searchText, normalized) {
			seen[normalized] = struct{}{}
			hits = append(hits, token)
		}
	}
	sort.Strings(hits)
	return hits
}

func fixtureWarnings(candidates []ReportCandidate, spec curriculum.LessonRecommendationSearchSpec) []string {
	warnings := map[string]struct{}{}
	for _, candidate := range candidates {
		if len(candidate.AvoidHits) > 0 {
			warnings["noise_suspect_present"] = struct{}{}
		}
		if spec.StageRole == "setup_intro" && containsAny(strings.ToLower(candidate.Title+" "+candidate.Description), "새들스티치", "saddle stitch", "고급", "advanced") {
			warnings["setup_intro_deep_practice_candidate"] = struct{}{}
		}
	}
	if len(warnings) == 0 {
		return nil
	}
	values := make([]string, 0, len(warnings))
	for warning := range warnings {
		values = append(values, warning)
	}
	sort.Strings(values)
	return values
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

func sanitizeProviderError(err error) string {
	message := strings.TrimSpace(err.Error())
	if message == "" {
		return "provider_error"
	}
	for _, sensitive := range []string{"YOUTUBE_API_KEY", "NAVER_CLIENT_ID", "NAVER_CLIENT_SECRET", "Authorization", "Bearer", "key="} {
		message = strings.ReplaceAll(message, sensitive, "[redacted]")
	}
	if len(message) > 200 {
		message = message[:200]
	}
	return message
}

func BuildDefaultProviders() map[string]Provider {
	return BuildProvidersWithCredentials(
		os.Getenv("YOUTUBE_API_KEY"),
		os.Getenv("NAVER_CLIENT_ID"),
		os.Getenv("NAVER_CLIENT_SECRET"),
	)
}

func BuildProvidersWithCredentials(youtubeAPIKey, naverClientID, naverClientSecret string) map[string]Provider {
	youtubeAPIKey = strings.TrimSpace(youtubeAPIKey)
	naverClientID = strings.TrimSpace(naverClientID)
	naverClientSecret = strings.TrimSpace(naverClientSecret)

	providers := map[string]Provider{}
	if youtubeAPIKey == "" {
		providers["youtube"] = missingCredentialsProvider{provider: "youtube"}
	} else {
		providers["youtube"] = extsearchProviderAdapter{
			provider: extsearch.NewYouTubeProviderWithLanguage(youtubeAPIKey, "ko"),
		}
	}
	if naverClientID == "" || naverClientSecret == "" {
		providers["naver_blog"] = missingCredentialsProvider{provider: "naver_blog"}
	} else {
		providers["naver_blog"] = extsearchProviderAdapter{
			provider: extsearch.NewNaverBlogProvider(naverClientID, naverClientSecret),
		}
	}
	return providers
}

func LoadEnv() {
	if path := strings.TrimSpace(os.Getenv("LEARNWEAVER_BACKEND_ENV_FILE")); path != "" {
		_ = godotenv.Load(path)
	}
	_ = godotenv.Load("/Run/learnweaver/backend.env")
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}

func WriteReport(path string, result any) error {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create report directory: %w", err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create report: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}
