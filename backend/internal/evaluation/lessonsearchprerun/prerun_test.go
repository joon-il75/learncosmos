package lessonsearchprerun

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
)

type fakeProvider struct {
	candidates []Candidate
	err        error
}

func (provider fakeProvider) Search(context.Context, Fixture, int) ([]Candidate, error) {
	return provider.candidates, provider.err
}

type forbiddenProvider struct {
	called *bool
}

func (provider forbiddenProvider) Search(context.Context, Fixture, int) ([]Candidate, error) {
	*provider.called = true
	return nil, errors.New("provider should not be called")
}

type fakeExtSearchProvider struct {
	results []extsearch.SearchResult
}

func (provider fakeExtSearchProvider) Search(ctx context.Context, query string, limit int) ([]extsearch.SearchResult, error) {
	if len(provider.results) > limit {
		return provider.results[:limit], nil
	}
	return provider.results, nil
}

func TestBuildReportScoresProviderCandidates(t *testing.T) {
	fixtures := []Fixture{{
		ID:       "youtube-leather-wallet-setup",
		Provider: "youtube",
		Query:    "가죽공예 카드지갑 도구 기초",
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "가죽공예 카드지갑 도구 기초",
			Intent:       "tutorial",
			MustInclude:  []string{"가죽공예", "카드지갑"},
			NiceToHave:   []string{"도구", "기초"},
			Avoid:        []string{"판매", "키트"},
			ContentTypes: []string{"video"},
			Language:     "ko",
			StageRole:    "setup_intro",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
		Candidates: []FixtureCandidate{
			{Title: "가죽공예 카드지갑 도구 기초", Relevant: true},
			{Title: "카드지갑 키트 판매", Relevant: false},
		},
	}}

	report := BuildReport(context.Background(), Options{
		Provider: "all",
		Mode:     DefaultMode,
		TopN:     5,
	}, fixtures, map[string]Provider{
		"youtube": fakeProvider{candidates: []Candidate{
			{
				Title:       "가죽공예 카드지갑 입문 도구와 기초 튜토리얼",
				Description: "초보자를 위한 준비물과 제작 순서",
				ContentType: "youtube",
			},
			{
				Title:       "카드지갑 DIY 키트 판매 후기",
				Description: "구매 링크와 구성품 소개",
				ContentType: "youtube",
			},
		}},
	}, time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC))

	if report.Summary.FixturesProcessed != 1 {
		t.Fatalf("processed = %d, want 1", report.Summary.FixturesProcessed)
	}
	if report.Summary.NoiseSuspects != 1 {
		t.Fatalf("noise suspects = %d, want 1", report.Summary.NoiseSuspects)
	}
	if report.Summary.RelevantLikeCandidates != 1 {
		t.Fatalf("relevant-like = %d, want 1", report.Summary.RelevantLikeCandidates)
	}
	if len(report.Fixtures) != 1 || len(report.Fixtures[0].Candidates) != 2 {
		t.Fatalf("candidates = %#v, want 2 candidates", report.Fixtures)
	}
	good := report.Fixtures[0].Candidates[0]
	if good.NeedsReview {
		t.Fatalf("good candidate needs review: %#v", good)
	}
	if good.Score <= 0 {
		t.Fatalf("good score = %d, want positive", good.Score)
	}
	if len(good.MustIncludeHits) != 2 {
		t.Fatalf("must hits = %#v, want 2 hits", good.MustIncludeHits)
	}
	noise := report.Fixtures[0].Candidates[1]
	if !noise.NeedsReview {
		t.Fatalf("noise candidate should need review: %#v", noise)
	}
	if noise.Warning != "avoid_terms_present" {
		t.Fatalf("noise warning = %q, want avoid_terms_present", noise.Warning)
	}
	if countNoiseToken(report.NoiseTaxonomySummary.RegisteredAvoidHits, "판매") != 1 {
		t.Fatalf("registered avoid hits = %#v, want 판매 count 1", report.NoiseTaxonomySummary.RegisteredAvoidHits)
	}
}

func TestBuildReportSummarizesNoiseTaxonomy(t *testing.T) {
	fixtures := []Fixture{{
		ID:       "naver-english-exam-noise",
		Provider: "naver_blog",
		Query:    "토익 오답노트 학습법",
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "토익 오답노트 학습법",
			MustInclude:  []string{"토익", "오답노트"},
			Avoid:        []string{"학원"},
			ContentTypes: []string{"article"},
			Language:     "ko",
			StageRole:    "practice_loop",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	}}

	report := BuildReport(context.Background(), Options{
		Provider: "all",
		Mode:     DefaultMode,
		TopN:     5,
	}, fixtures, map[string]Provider{
		"naver_blog": fakeProvider{candidates: []Candidate{
			{
				Title:       "토익 오답노트 학습법",
				Description: "틀린 문제를 복습하는 방법",
				ContentType: "naver_blog",
			},
			{
				Title:       "토익 합격보장 수강신청",
				Description: "응시자격과 편입 상담",
				ContentType: "naver_blog",
			},
			{
				Title:       "토익 학원 광고",
				Description: "수업 안내",
				ContentType: "naver_blog",
			},
		}},
	}, time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC))

	if countNoiseToken(report.NoiseTaxonomySummary.RegisteredAvoidHits, "학원") != 1 {
		t.Fatalf("registered avoid hits = %#v, want 학원 count 1", report.NoiseTaxonomySummary.RegisteredAvoidHits)
	}
	for _, token := range []string{"합격보장", "수강신청", "응시자격", "편입"} {
		if countNoiseToken(report.NoiseTaxonomySummary.UnregisteredNoiseHints, token) != 1 {
			t.Fatalf("unregistered hints = %#v, want %s count 1", report.NoiseTaxonomySummary.UnregisteredNoiseHints, token)
		}
	}
	if countNoiseToken(report.NoiseTaxonomySummary.UnregisteredNoiseHints, "학원") != 0 {
		t.Fatalf("unregistered hints include registered avoid token: %#v", report.NoiseTaxonomySummary.UnregisteredNoiseHints)
	}
}

func TestBuildReportSkipsProviderNotWired(t *testing.T) {
	report := BuildReport(context.Background(), Options{
		Provider: "all",
		Mode:     DefaultMode,
		TopN:     5,
	}, []Fixture{{
		ID:       "naver-writing",
		Provider: "naver_blog",
		Query:    "브런치 글쓰기",
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "브런치 글쓰기",
			MustInclude:  []string{"브런치", "글쓰기"},
		},
	}}, map[string]Provider{}, time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC))

	if len(report.Fixtures) != 1 {
		t.Fatalf("fixtures = %d, want 1", len(report.Fixtures))
	}
	if report.Fixtures[0].Skipped != "skipped_provider_not_wired" {
		t.Fatalf("skipped = %q, want skipped_provider_not_wired", report.Fixtures[0].Skipped)
	}
	if report.Summary.FixturesSkipped != 1 {
		t.Fatalf("fixtures skipped = %d, want 1", report.Summary.FixturesSkipped)
	}
	if report.Summary.ProviderErrors != 0 {
		t.Fatalf("provider errors = %d, want 0", report.Summary.ProviderErrors)
	}
}

func TestBuildReportSkipsMissingCredentials(t *testing.T) {
	report := BuildReport(context.Background(), Options{
		Provider: "all",
		Mode:     DefaultMode,
		TopN:     5,
	}, []Fixture{{
		ID:       "youtube-test",
		Provider: "youtube",
		Query:    "가죽공예 카드지갑",
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "가죽공예 카드지갑",
			MustInclude:  []string{"가죽공예", "카드지갑"},
		},
	}}, map[string]Provider{
		"youtube": missingCredentialsProvider{provider: "youtube"},
	}, time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC))

	if report.Fixtures[0].Skipped != "skipped_missing_credentials" {
		t.Fatalf("skipped = %q, want skipped_missing_credentials", report.Fixtures[0].Skipped)
	}
	if report.Summary.ProviderErrors != 0 {
		t.Fatalf("provider errors = %d, want 0", report.Summary.ProviderErrors)
	}
	if report.Summary.FixturesSkipped != 1 {
		t.Fatalf("fixtures skipped = %d, want 1", report.Summary.FixturesSkipped)
	}
}

func TestBuildGapReportSummarizesCoverageAndGaps(t *testing.T) {
	fixtures := []Fixture{
		{
			ID:       "youtube-setup",
			Provider: "youtube",
			SearchSpec: curriculum.LessonRecommendationSearchSpec{
				StageRole:    "setup_intro",
				ContentTypes: []string{"video"},
			},
		},
		{
			ID:       "naver-core",
			Provider: "naver_blog",
			SearchSpec: curriculum.LessonRecommendationSearchSpec{
				StageRole:    "core_pattern",
				ContentTypes: []string{"article", "video"},
			},
		},
	}

	report := BuildGapReport(Options{Provider: "all", Mode: DiscoverGapMode}, fixtures, time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC))

	if report.Mode != DiscoverGapMode {
		t.Fatalf("mode = %q, want %q", report.Mode, DiscoverGapMode)
	}
	if report.Summary.FixturesTotal != 2 {
		t.Fatalf("fixtures total = %d, want 2", report.Summary.FixturesTotal)
	}
	if report.Summary.Providers["youtube"] != 1 || report.Summary.Providers["naver_blog"] != 1 {
		t.Fatalf("providers = %#v", report.Summary.Providers)
	}
	if report.Summary.StageRoles["setup_intro"] != 1 || report.Summary.StageRoles["core_pattern"] != 1 {
		t.Fatalf("stage roles = %#v", report.Summary.StageRoles)
	}
	if report.Summary.ContentTypes["video"] != 1 || report.Summary.ContentTypes["article+video"] != 1 {
		t.Fatalf("content types = %#v", report.Summary.ContentTypes)
	}
	if !hasProviderStageCoverage(report, "youtube", "setup_intro", 1) {
		t.Fatalf("missing youtube setup coverage: %#v", report.Coverage.ByProviderStage)
	}
	if !hasGap(report, "high", "stage_role", "", "first_output", "") {
		t.Fatalf("missing high first_output gap: %#v", report.GapCandidates)
	}
	if !hasGap(report, "medium", "provider_stage", "naver_blog", "setup_intro", "") {
		t.Fatalf("missing naver setup provider-stage gap: %#v", report.GapCandidates)
	}
	if !hasGap(report, "low", "content_type", "", "", "article") {
		t.Fatalf("missing article content type gap: %#v", report.GapCandidates)
	}
	for _, candidate := range report.GapCandidates {
		if candidate.FixtureAutoApply {
			t.Fatalf("gap candidate should not auto-apply: %#v", candidate)
		}
	}
	if report.Summary.GapCount != len(report.GapCandidates) {
		t.Fatalf("gap count = %d, candidates = %d", report.Summary.GapCount, len(report.GapCandidates))
	}
}

func TestBuildGapReportRespectsProviderFilter(t *testing.T) {
	fixtures := []Fixture{
		{
			ID:       "youtube-setup",
			Provider: "youtube",
			SearchSpec: curriculum.LessonRecommendationSearchSpec{
				StageRole:    "setup_intro",
				ContentTypes: []string{"video"},
			},
		},
		{
			ID:       "naver-setup",
			Provider: "naver_blog",
			SearchSpec: curriculum.LessonRecommendationSearchSpec{
				StageRole:    "setup_intro",
				ContentTypes: []string{"article"},
			},
		},
	}

	report := BuildGapReport(Options{Provider: "youtube", Mode: DiscoverGapMode}, fixtures, time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC))

	if report.Summary.Providers["youtube"] != 1 {
		t.Fatalf("providers = %#v, want youtube only", report.Summary.Providers)
	}
	if _, ok := report.Summary.Providers["naver_blog"]; ok {
		t.Fatalf("providers = %#v, did not expect naver_blog under youtube filter", report.Summary.Providers)
	}
	for _, coverage := range report.Coverage.ByProviderStage {
		if coverage.Provider != "youtube" {
			t.Fatalf("coverage includes provider %q under youtube filter", coverage.Provider)
		}
	}
	for _, candidate := range report.GapCandidates {
		if candidate.GapType == "provider_stage" && candidate.Provider != "youtube" {
			t.Fatalf("gap includes provider %q under youtube filter", candidate.Provider)
		}
	}
}

func TestRunDiscoverGapWritesReportWithoutCallingProviders(t *testing.T) {
	dir := t.TempDir()
	fixturePath := filepath.Join(dir, "fixtures.json")
	outPath := filepath.Join(dir, "gap-report.json")
	payload := []byte(`[
		{
			"id":"youtube-test",
			"provider":"youtube",
			"query":"가죽공예 카드지갑",
			"search_spec":{
				"primary_query":"가죽공예 카드지갑",
				"content_types":["video"],
				"stage_role":"setup_intro"
			},
			"candidates":[
				{"title":"가죽공예 카드지갑","content_type":"youtube","relevant":true}
			]
		}
	]`)
	if err := os.WriteFile(fixturePath, payload, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	called := false
	err := Run(context.Background(), Options{
		FixturesPath: fixturePath,
		OutputPath:   outPath,
		Provider:     "all",
		Mode:         DiscoverGapMode,
		Timeout:      time.Second,
	}, map[string]Provider{
		"youtube": forbiddenProvider{called: &called},
	})
	if err != nil {
		t.Fatalf("Run discover-gap: %v", err)
	}
	if called {
		t.Fatal("discover-gap called provider")
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	var result GapReport
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("parse gap report: %v", err)
	}
	if result.Mode != DiscoverGapMode {
		t.Fatalf("mode = %q, want %q", result.Mode, DiscoverGapMode)
	}
	if len(result.GapCandidates) == 0 {
		t.Fatal("expected at least one gap candidate")
	}
}

func TestExtsearchProviderAdapterMapsResults(t *testing.T) {
	adapter := extsearchProviderAdapter{provider: fakeExtSearchProvider{results: []extsearch.SearchResult{{
		Source:      "youtube",
		Title:       "가죽공예 카드지갑 기초 튜토리얼",
		Description: "초보 제작 순서",
	}}}}

	candidates, err := adapter.Search(context.Background(), Fixture{Query: "가죽공예 카드지갑"}, 5)
	if err != nil {
		t.Fatalf("adapter search: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates = %d, want 1", len(candidates))
	}
	if candidates[0].ContentType != "youtube" {
		t.Fatalf("content type = %q, want youtube", candidates[0].ContentType)
	}
}

func TestBuildDefaultProvidersUsesMissingCredentialsAdapters(t *testing.T) {
	t.Setenv("YOUTUBE_API_KEY", "")
	t.Setenv("NAVER_CLIENT_ID", "")
	t.Setenv("NAVER_CLIENT_SECRET", "")

	providers := BuildDefaultProviders()
	for _, providerName := range []string{"youtube", "naver_blog"} {
		provider, ok := providers[providerName]
		if !ok {
			t.Fatalf("missing provider %q", providerName)
		}
		_, err := provider.Search(context.Background(), Fixture{Provider: providerName, Query: "test"}, 5)
		var missingCredentials providerMissingCredentialsError
		if !errors.As(err, &missingCredentials) {
			t.Fatalf("provider %q error = %v, want missing credentials", providerName, err)
		}
	}
}

func TestBuildDefaultProvidersUsesConfiguredAdapters(t *testing.T) {
	t.Setenv("YOUTUBE_API_KEY", "youtube-test-key")
	t.Setenv("NAVER_CLIENT_ID", "naver-test-id")
	t.Setenv("NAVER_CLIENT_SECRET", "naver-test-secret")

	providers := BuildDefaultProviders()
	if _, ok := providers["youtube"].(extsearchProviderAdapter); !ok {
		t.Fatalf("youtube provider = %T, want extsearchProviderAdapter", providers["youtube"])
	}
	if _, ok := providers["naver_blog"].(extsearchProviderAdapter); !ok {
		t.Fatalf("naver provider = %T, want extsearchProviderAdapter", providers["naver_blog"])
	}
}

func TestRunWritesValidReport(t *testing.T) {
	dir := t.TempDir()
	fixturePath := filepath.Join(dir, "fixtures.json")
	outPath := filepath.Join(dir, "report.json")
	payload := []byte(`[
		{
			"id":"youtube-test",
			"provider":"youtube",
			"query":"가죽공예 카드지갑",
			"search_spec":{
				"primary_query":"가죽공예 카드지갑",
				"intent":"tutorial",
				"must_include":["가죽공예","카드지갑"],
				"content_types":["video"],
				"language":"ko",
				"stage_role":"core_pattern",
				"source":"pattern_template"
			},
			"candidates":[
				{"title":"가죽공예 카드지갑","content_type":"youtube","relevant":true}
			]
		}
	]`)
	if err := os.WriteFile(fixturePath, payload, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	err := Run(context.Background(), Options{
		FixturesPath: fixturePath,
		OutputPath:   outPath,
		Provider:     "youtube",
		Mode:         DefaultMode,
		TopN:         5,
		Timeout:      time.Second,
	}, map[string]Provider{
		"youtube": fakeProvider{candidates: []Candidate{{
			Title:       "가죽공예 카드지갑 기초 튜토리얼",
			Description: "재단과 바느질 준비",
			ContentType: "youtube",
		}}},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	var result Report
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("parse report: %v", err)
	}
	if result.Summary.FixturesProcessed != 1 {
		t.Fatalf("processed = %d, want 1", result.Summary.FixturesProcessed)
	}
	if len(result.Fixtures) != 1 || len(result.Fixtures[0].Candidates) != 1 {
		t.Fatalf("report fixtures = %#v", result.Fixtures)
	}
}

func TestValidateOptionsAllowsDiscoverGapAndRejectsUnsupportedValues(t *testing.T) {
	if err := ValidateOptions(Options{Mode: DiscoverGapMode, Provider: "all"}); err != nil {
		t.Fatalf("discover-gap validation failed: %v", err)
	}
	if err := ValidateOptions(Options{Mode: "unsupported", Provider: "all"}); err == nil {
		t.Fatal("expected unsupported mode error")
	}
	if err := ValidateOptions(Options{Mode: DefaultMode, Provider: "web"}); err == nil {
		t.Fatal("expected unsupported provider error")
	}
	if err := ValidateOptions(Options{Mode: DefaultMode, Provider: "youtube", Limit: -1}); err == nil {
		t.Fatal("expected negative limit error")
	}
}

func TestSanitizeProviderErrorRedactsKnownSecretMarkers(t *testing.T) {
	message := sanitizeProviderError(errors.New("Authorization Bearer abc key=secret NAVER_CLIENT_SECRET"))
	for _, value := range []string{"Authorization", "Bearer", "key=", "NAVER_CLIENT_SECRET"} {
		if contains(value, message) {
			t.Fatalf("message %q still contains %q", message, value)
		}
	}
}

func contains(needle, haystack string) bool {
	return len(needle) == 0 || strings.Contains(haystack, needle)
}

func countNoiseToken(values []NoiseTokenCount, token string) int {
	for _, value := range values {
		if value.Token == token {
			return value.Count
		}
	}
	return 0
}

func hasProviderStageCoverage(report GapReport, provider string, stageRole string, count int) bool {
	for _, coverage := range report.Coverage.ByProviderStage {
		if coverage.Provider == provider && coverage.StageRole == stageRole && coverage.Count == count {
			return true
		}
	}
	return false
}

func hasGap(report GapReport, priority string, gapType string, provider string, stageRole string, contentType string) bool {
	for _, candidate := range report.GapCandidates {
		if candidate.Priority == priority && candidate.GapType == gapType && candidate.Provider == provider && candidate.StageRole == stageRole && candidate.ContentType == contentType {
			return true
		}
	}
	return false
}
