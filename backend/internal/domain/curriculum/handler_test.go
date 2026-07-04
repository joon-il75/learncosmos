package curriculum

import (
	"testing"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/pkg/extsearch"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func TestDecideLessonSearchExecution(t *testing.T) {
	got := decideLessonSearchExecution(1)
	if got.Provider != primaryEmbeddingProvider() {
		t.Fatalf("provider mismatch: got %q want %q", got.Provider, primaryEmbeddingProvider())
	}
	if got.APIKey != "" {
		t.Fatalf("api key mismatch: got %q want empty", got.APIKey)
	}
	if got.BillingStatus != "charged" {
		t.Fatalf("billing status mismatch: got %q", got.BillingStatus)
	}
	if got.EffectiveCost != 1 {
		t.Fatalf("effective cost mismatch: got %d want 1", got.EffectiveCost)
	}
	if !got.ShouldCharge {
		t.Fatalf("should charge mismatch: got false want true")
	}
}

func TestGetEditableDraftStatuses(t *testing.T) {
	tests := []struct {
		name            string
		includeLearning bool
		want            []DraftStatus
	}{
		{
			name:            "pre learning editing allows draft and confirmed",
			includeLearning: false,
			want:            []DraftStatus{DraftStatusDraft, DraftStatusConfirmed},
		},
		{
			name:            "learning edit contexts also keep learning",
			includeLearning: true,
			want:            []DraftStatus{DraftStatusDraft, DraftStatusConfirmed, DraftStatusLearning},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEditableDraftStatuses(tt.includeLearning)
			if len(got) != len(tt.want) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildRemainingRebuildPreserveMaskFromCourse(t *testing.T) {
	lessons := []CourseLessonTree{
		{
			SubLessons: []CourseLessonTree{
				{
					RuntimeEntry: &CourseLessonRuntimeEntry{Status: LessonRuntimeInProgress},
				},
			},
		},
		{
			SubLessons: []CourseLessonTree{
				{
					Points: []CoursePointAggregate{
						{Point: CoursePoint{Status: PointStatusReady}},
					},
				},
			},
		},
		{
			SubLessons: []CourseLessonTree{
				{
					Points: []CoursePointAggregate{
						{Point: CoursePoint{Status: PointStatusCompleted}},
					},
				},
			},
		},
	}

	got := buildRemainingRebuildPreserveMaskFromCourse(lessons)
	want := []bool{true, false, true}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %t, want %t", i, got[i], want[i])
		}
	}
}

func TestIsDraftPointTouched(t *testing.T) {
	candidate := ResourceSelectionCandidate
	selected := ResourceSelectionSelected
	tests := []struct {
		name  string
		point DraftPointAggregate
		want  bool
	}{
		{
			name:  "candidate exploration point is untouched",
			point: DraftPointAggregate{Point: CourseDraftPoint{Status: PointStatusDraft, SelectionState: &candidate}},
			want:  false,
		},
		{
			name:  "selected exploration point is touched",
			point: DraftPointAggregate{Point: CourseDraftPoint{Status: PointStatusDraft, SelectionState: &selected}},
			want:  true,
		},
		{
			name:  "research block makes point touched",
			point: DraftPointAggregate{Point: CourseDraftPoint{Status: PointStatusDraft}, Blocks: []CourseDraftPointBlock{{ID: uuid.New()}}},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDraftPointTouched(tt.point); got != tt.want {
				t.Fatalf("isDraftPointTouched() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestMergeRemainingDraftLessons(t *testing.T) {
	existing := []DraftLessonTree{
		{Lesson: CourseDraftLesson{ID: uuid.New(), Title: "preserve"}},
		{Lesson: CourseDraftLesson{ID: uuid.New(), Title: "replace-a"}},
		{Lesson: CourseDraftLesson{ID: uuid.New(), Title: "replace-b"}},
	}
	generated := []DraftLessonTree{
		{Lesson: CourseDraftLesson{ID: uuid.New(), Title: "new-a"}},
		{Lesson: CourseDraftLesson{ID: uuid.New(), Title: "new-b"}},
	}
	mask := []bool{true, false, false}

	got, err := mergeRemainingDraftLessons(existing, generated, mask)
	if err != nil {
		t.Fatalf("mergeRemainingDraftLessons returned error: %v", err)
	}
	if got[0].Lesson.Title != "preserve" || got[1].Lesson.Title != "new-a" || got[2].Lesson.Title != "new-b" {
		t.Fatalf("unexpected merge result: %+v", got)
	}
}

func TestBuildModeForSystemFeature(t *testing.T) {
	tests := []struct {
		feature string
		want    string
	}{
		{feature: "default", want: "system_llm_default"},
		{feature: "pro_curriculum", want: "system_llm_pro_curriculum"},
		{feature: "PRO_CURRICULUM", want: "system_llm_pro_curriculum"},
		{feature: "", want: "system_llm_default"},
	}

	for _, tt := range tests {
		if got := buildModeForSystemFeature(tt.feature); got != tt.want {
			t.Fatalf("build mode mismatch for %q: got %q want %q", tt.feature, got, tt.want)
		}
	}
}

func TestBuildExternalSearchQuery(t *testing.T) {
	query := BuildExternalSearchQuery(normalizer.LessonSearchQuery{
		RawQuery: "코바늘 입문자를 위한 기초 도구 소개와 종류 설명",
		Tokens:   []string{"코바늘", "입문", "기초", "도구", "종류", "뜨개"},
	})
	if query == "" {
		t.Fatal("query should not be empty")
	}
	if got, want := query, "코바늘 입문 도구 소개 종류 용도"; got != want {
		t.Fatalf("unexpected query: got %q want %q", got, want)
	}
}

func TestBuildExternalSearchQuery_ConceptIntro(t *testing.T) {
	query := BuildExternalSearchQuery(normalizer.LessonSearchQuery{
		RawQuery: "코바늘이 무엇인지 이해하고 기본 원리를 설명합니다.",
		Tokens:   []string{"코바늘", "이해", "기본", "원리", "설명"},
	})
	if got, want := query, "코바늘 이해 개념 설명"; got != want {
		t.Fatalf("unexpected concept_intro query: got %q want %q", got, want)
	}
}

func TestBuildExternalSearchQuery_HowTo(t *testing.T) {
	query := BuildExternalSearchQuery(normalizer.LessonSearchQuery{
		RawQuery: "사슬뜨기를 직접 따라 하며 시작할 수 있습니다.",
		Tokens:   []string{"사슬뜨기", "따라하기", "시작", "방법"},
	})
	if got, want := query, "사슬뜨기 따라하기 시작 하는 법 방법"; got != want {
		t.Fatalf("unexpected how_to query: got %q want %q", got, want)
	}
}

func TestBuildExternalSearchQuery_Comparison(t *testing.T) {
	query := BuildExternalSearchQuery(normalizer.LessonSearchQuery{
		RawQuery: "우드 코바늘과 알루미늄 코바늘의 차이를 비교합니다.",
		Tokens:   []string{"우드", "코바늘", "알루미늄", "차이", "비교"},
	})
	if got, want := query, "우드 코바늘 알루미늄 비교 차이"; got != want {
		t.Fatalf("unexpected comparison query: got %q want %q", got, want)
	}
}

func TestBuildExternalSearchQuery_ComparisonPreservesSpecificDomainTokens(t *testing.T) {
	query := BuildExternalSearchQuery(normalizer.LessonSearchQuery{
		RawQuery: "홈카페 커피 입문에서 에스프레소 머신과 원두, 그라인더 차이를 비교합니다.",
		Tokens:   []string{"홈카페", "커피", "입문", "집에서", "에스프레소", "머신", "원두", "그라인더", "차이", "비교"},
	})
	if got, want := query, "홈카페 커피 에스프레소 머신 원두 비교 차이"; got != want {
		t.Fatalf("unexpected comparison query preserving domain tokens: got %q want %q", got, want)
	}
}

func TestBuildYouTubeExternalSearchQueryLanguageIntent(t *testing.T) {
	if got, want := buildYouTubeExternalSearchQuery("AI coding SaaS MVP"), "AI coding SaaS MVP tutorial"; got != want {
		t.Fatalf("unexpected english youtube query: got %q want %q", got, want)
	}
	if got, want := buildYouTubeExternalSearchQuery("사슬뜨기 기본 연습"), "사슬뜨기 기본 연습 영상"; got != want {
		t.Fatalf("unexpected korean youtube query: got %q want %q", got, want)
	}
}

func TestMergeExternalCandidates(t *testing.T) {
	internalID := uuid.New()
	externalID := uuid.New()

	merged := mergeExternalCandidates(
		[]ContentSearchCandidate{
			{ContentID: &internalID, Title: "internal 1"},
			{ContentID: &externalID, Title: "duplicate"},
		},
		[]ContentSearchCandidate{
			{ContentID: &externalID, Title: "duplicate"},
			{ContentID: nil, ExternalURL: stringPointer("https://www.youtube.com/watch?v=abc"), Title: "external 2"},
		},
		3,
	)

	if len(merged) != 3 {
		t.Fatalf("merged length mismatch: got %d want %d", len(merged), 3)
	}
	if merged[0].ContentID == nil || *merged[0].ContentID != internalID {
		t.Fatalf("first candidate should preserve internal order")
	}
	if merged[1].ContentID == nil || *merged[1].ContentID != externalID {
		t.Fatalf("second candidate should keep unique duplicate once")
	}
	if merged[2].Title != "external 2" {
		t.Fatalf("third candidate mismatch: got %q", merged[2].Title)
	}
}

func TestScoreExternalSearchResult(t *testing.T) {
	preferredFormat := "video"
	score := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			Tokens:         []string{"코바늘", "사슬뜨기"},
			ExpandedTokens: []string{"뜨개질"},
		},
		extsearch.SearchResult{
			Source:      "youtube",
			Title:       "코바늘 사슬뜨기 기초 배우기",
			Description: "입문자를 위한 뜨개질 기본 영상",
		},
		&preferredFormat,
		0,
	)
	if score <= 0.06 {
		t.Fatalf("expected strong external score, got %f", score)
	}
}

func TestInstructionalExternalSearchResultForQueryRequiresYoutubeRelevance(t *testing.T) {
	queryBundle := normalizer.BuildLessonSearchQuery("기초 수채화 연습", "", "", "", "", "")

	if !isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "youtube",
		Title:       "기초 수채화 연습 방법",
		Description: "초보자를 위한 수채화 단계별 영상",
	}, queryBundle) {
		t.Fatal("expected relevant instructional youtube result to pass")
	}

	if isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "youtube",
		Title:       "카페 브이로그와 일상 기록",
		Description: "잔잔한 하루를 담은 영상",
	}, queryBundle) {
		t.Fatal("expected unrelated youtube result to be filtered")
	}

	if isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "youtube",
		Title:       "수채화 감성 브이로그",
		Description: "수채화 그림을 배경으로 한 일상 영상",
	}, queryBundle) {
		t.Fatal("expected topical but non-instructional youtube result to be filtered")
	}
}

func TestInstructionalExternalSearchResultForQueryFiltersGenericBlogMismatch(t *testing.T) {
	queryBundle := normalizer.BuildLessonSearchQuery("목공예 간단 기법 작은 결과물 완성 재료 방법", "", "", "", "", "")

	if !isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "naver_blog",
		Title:       "목공예 초보를 위한 작은 선반 만들기 방법",
		Description: "목공 재료와 기초 기법을 단계별로 정리합니다.",
	}, queryBundle) {
		t.Fatal("expected relevant woodworking blog result to pass")
	}

	if isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "naver_blog",
		Title:       "기본김밥 만들기 재료 간단 꼬마김밥",
		Description: "간단한 재료로 결과물을 완성하는 방법입니다.",
	}, queryBundle) {
		t.Fatal("expected generic cooking blog result to be filtered")
	}
}

func TestInstructionalExternalSearchResultForQueryRequiresPrimaryAnchor(t *testing.T) {
	queryBundle := normalizer.BuildLessonSearchQuery("피아노 팝송 반주 카피", "", "", "", "", "")

	if !isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "naver_blog",
		Title:       "피아노 배우기 코드기초 보고 따라치기",
		Description: "초보자를 위한 피아노 반주 방법을 단계별로 정리합니다.",
	}, queryBundle) {
		t.Fatal("expected piano instructional blog result to pass")
	}

	if isInstructionalExternalSearchResultForQuery(extsearch.SearchResult{
		Source:      "naver_blog",
		Title:       "수원 취미 보컬 - K-발라드부터 팝송까지",
		Description: "취미 보컬 레슨에서 팝송을 연습하는 방법을 소개합니다.",
	}, queryBundle) {
		t.Fatal("expected blog matching only weak context terms to be filtered")
	}
}

func TestScoreExternalSearchResult_ToolIntroPenalizesHowTo(t *testing.T) {
	scoreTool := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "기초 도구 소개, 코바늘 소개, 코바늘의 다양한 종류와 각각의 용도에 대해 설명합니다.",
			Tokens:   []string{"코바늘", "도구", "소개", "종류", "용도"},
		},
		extsearch.SearchResult{
			Source:      "naver_blog",
			Title:       "코바늘 종류와 용도 정리",
			Description: "코바늘 호수와 재질, 실 종류를 소개합니다.",
		},
		nil,
		0,
	)
	scoreHowTo := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "기초 도구 소개, 코바늘 소개, 코바늘의 다양한 종류와 각각의 용도에 대해 설명합니다.",
			Tokens:   []string{"코바늘", "도구", "소개", "종류", "용도"},
		},
		extsearch.SearchResult{
			Source:      "youtube",
			Title:       "코바늘 사슬뜨기 하는 방법",
			Description: "입문자를 위한 사슬뜨기 실습 영상",
		},
		nil,
		0,
	)
	if scoreTool <= scoreHowTo {
		t.Fatalf("expected tool-intro content to outrank how-to content: tool=%f howto=%f", scoreTool, scoreHowTo)
	}
}

func TestScoreExternalSearchResult_ConceptIntroPrefersExplanation(t *testing.T) {
	scoreConcept := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "코바늘이 무엇인지 이해하고 기본 원리를 설명합니다.",
			Tokens:   []string{"코바늘", "이해", "기본", "원리", "설명"},
		},
		extsearch.SearchResult{
			Source:      "naver_blog",
			Title:       "코바늘 기본 개념과 원리 설명",
			Description: "코바늘의 구조와 기초 이론을 설명합니다.",
		},
		nil,
		0,
	)
	scoreHowTo := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "코바늘이 무엇인지 이해하고 기본 원리를 설명합니다.",
			Tokens:   []string{"코바늘", "이해", "기본", "원리", "설명"},
		},
		extsearch.SearchResult{
			Source:      "youtube",
			Title:       "코바늘 사슬뜨기 하는 법",
			Description: "실습 중심의 입문 영상",
		},
		nil,
		0,
	)
	if scoreConcept <= scoreHowTo {
		t.Fatalf("expected concept_intro content to outrank how-to content: concept=%f howto=%f", scoreConcept, scoreHowTo)
	}
}

func TestScoreExternalSearchResult_ComparisonPrefersComparison(t *testing.T) {
	scoreCompare := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "우드 코바늘과 알루미늄 코바늘의 차이를 비교합니다.",
			Tokens:   []string{"우드", "코바늘", "알루미늄", "차이", "비교"},
		},
		extsearch.SearchResult{
			Source:      "naver_blog",
			Title:       "우드 코바늘과 알루미늄 코바늘 비교",
			Description: "재질별 차이와 선택 기준 정리",
		},
		nil,
		0,
	)
	scorePractice := scoreExternalSearchResult(
		normalizer.LessonSearchQuery{
			RawQuery: "우드 코바늘과 알루미늄 코바늘의 차이를 비교합니다.",
			Tokens:   []string{"우드", "코바늘", "알루미늄", "차이", "비교"},
		},
		extsearch.SearchResult{
			Source:      "youtube",
			Title:       "코바늘 사슬뜨기 실습",
			Description: "초보자를 위한 따라하기",
		},
		nil,
		0,
	)
	if scoreCompare <= scorePractice {
		t.Fatalf("expected comparison content to outrank practice content: compare=%f practice=%f", scoreCompare, scorePractice)
	}
}

func TestMergeExternalCandidates_SortsByRankScore(t *testing.T) {
	internalID := uuid.New()
	merged := mergeExternalCandidates(
		[]ContentSearchCandidate{
			{ContentID: &internalID, Title: "internal", RankScore: 0.04},
		},
		[]ContentSearchCandidate{
			{ExternalURL: stringPointer("https://www.youtube.com/watch?v=high"), Title: "external-high", RankScore: 0.07},
			{ExternalURL: stringPointer("https://www.youtube.com/watch?v=low"), Title: "external-low", RankScore: 0.03},
		},
		3,
	)
	if len(merged) != 3 {
		t.Fatalf("merged length mismatch: got %d want 3", len(merged))
	}
	if merged[0].Title != "external-high" {
		t.Fatalf("expected highest-score external candidate first, got %q", merged[0].Title)
	}
	if merged[1].Title != "internal" {
		t.Fatalf("expected internal candidate second, got %q", merged[1].Title)
	}
}

func TestExternalSearchResultKey(t *testing.T) {
	got := externalSearchResultKey(extsearch.SearchResult{
		Source:            "naver_blog",
		ExternalContentID: "blog.naver.com/test/123",
		CanonicalURL:      "https://blog.naver.com/test/123",
		Title:             "테스트 글",
	})
	want := "url:https://blog.naver.com/test/123"
	if got != want {
		t.Fatalf("externalSearchResultKey mismatch: got %q want %q", got, want)
	}
}

func TestCourseGenerationAsyncMaxActiveJobsDefaultAndEnv(t *testing.T) {
	t.Setenv("COURSE_GENERATION_ASYNC_MAX_ACTIVE_JOBS", "")
	if got := courseGenerationAsyncMaxActiveJobs(); got != defaultCourseGenerationAsyncMaxActiveJobs {
		t.Fatalf("default max active jobs = %d, want %d", got, defaultCourseGenerationAsyncMaxActiveJobs)
	}

	t.Setenv("COURSE_GENERATION_ASYNC_MAX_ACTIVE_JOBS", "24")
	if got := courseGenerationAsyncMaxActiveJobs(); got != 24 {
		t.Fatalf("env max active jobs = %d, want 24", got)
	}

	t.Setenv("COURSE_GENERATION_ASYNC_MAX_ACTIVE_JOBS", "invalid")
	if got := courseGenerationAsyncMaxActiveJobs(); got != defaultCourseGenerationAsyncMaxActiveJobs {
		t.Fatalf("invalid env max active jobs = %d, want %d", got, defaultCourseGenerationAsyncMaxActiveJobs)
	}
}

func TestBuildCourseGenerationWaitEstimate(t *testing.T) {
	position := 9
	estimate := buildCourseGenerationWaitEstimate(llmjobs.ActiveJobEstimate{
		ActiveCount:   12,
		QueuePosition: &position,
	}, 40)
	if estimate == nil {
		t.Fatal("estimate is nil")
	}
	if estimate.ActiveJobs != 12 || estimate.QueuePosition == nil || *estimate.QueuePosition != 9 || estimate.MaxActiveJobs != 40 {
		t.Fatalf("unexpected estimate: %+v", estimate)
	}
	if estimate.MinSeconds != 15 || estimate.MaxSeconds != 30 || estimate.Label != "약 15~30초" {
		t.Fatalf("unexpected wait range: %+v", estimate)
	}

	running := 0
	estimate = buildCourseGenerationWaitEstimate(llmjobs.ActiveJobEstimate{
		ActiveCount:   8,
		QueuePosition: &running,
	}, 40)
	if estimate.MinSeconds != 5 || estimate.MaxSeconds != 20 {
		t.Fatalf("running wait range = %d~%d, want 5~20", estimate.MinSeconds, estimate.MaxSeconds)
	}
}
