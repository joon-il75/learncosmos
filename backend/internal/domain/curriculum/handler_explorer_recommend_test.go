package curriculum

import (
	"strings"
	"testing"

	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func TestFilterExplorerCandidatesByContext_StrongTopicTokenRequired(t *testing.T) {
	queryBundle := normalizer.BuildLessonSearchQuery(
		"어반스케치",
		"",
		"",
		"도시 스케치 연습",
		"어반스케치 재료 이해",
		"",
	)

	baristaDesc := "에스프레소 머신 구조 이해와 사용법을 설명합니다."
	sketchDesc := "어반스케치 재료와 도시 스케치 구도를 연습합니다."
	candidates := []ContentSearchCandidate{
		{
			Title:       "[언택트복지관] 1. 에스프레소 머신의 명칭과 사용법",
			Description: &baristaDesc,
		},
		{
			Title:       "어반스케치 재료와 도시 풍경 스케치 입문",
			Description: &sketchDesc,
		},
	}

	filtered := filterExplorerCandidatesByContext(candidates, queryBundle, 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "어반스케치 재료와 도시 풍경 스케치 입문" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByContext_ShortTokenWholeWordOnly(t *testing.T) {
	// "코드" (2 runes) must NOT match "qr코드" in unrelated content.
	queryBundle := normalizer.BuildLessonSearchQuery(
		"통기타 배우기 코드 연습",
		"",
		"",
		"코드 연습",
		"",
		"",
	)

	espressoDesc := "에스프레소 머신 구조 이해와 qr코드 스캔으로 사용법을 알아봅니다."
	guitarDesc := "통기타 오픈 코드 연습과 스트로크 패턴을 익힙니다."
	candidates := []ContentSearchCandidate{
		{
			Title:       "[언택트복지관] 1. 에스프레소 머신의 명칭과 사용법",
			Description: &espressoDesc,
		},
		{
			Title:       "통기타 코드 연습 입문",
			Description: &guitarDesc,
		},
	}

	filtered := filterExplorerCandidatesByContext(candidates, queryBundle, 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1 (only guitar content)", len(filtered))
	}
	if filtered[0].Title != "통기타 코드 연습 입문" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByStage_PrefersNodeThenLevelThenCourse(t *testing.T) {
	nodeBundle := normalizer.BuildLessonSearchQuery("", "", "", "머그컵 손잡이 만들기", "", "머그컵 손잡이 연결")
	levelBundle := normalizer.BuildLessonSearchQuery("", "도자기 머그컵", "", "", "", "")
	courseBundle := normalizer.BuildLessonSearchQuery("도자기 배우기", "", "", "", "", "")

	mugNodeDesc := "도자기 머그컵 손잡이 연결과 형태 보정을 설명합니다."
	mugLevelDesc := "도자기 머그컵 물레 성형과 컵 형태 잡기를 다룹니다."
	craftCourseDesc := "도자기 기초 재료와 소성 흐름을 설명합니다."
	coffeeDesc := "카페 브이로그와 커피 머신 사용법을 소개합니다."

	candidates := []ContentSearchCandidate{
		{Title: "도자기 머그컵 손잡이 연결", Description: &mugNodeDesc},
		{Title: "도자기 머그컵 성형 입문", Description: &mugLevelDesc},
		{Title: "도자기 기초 재료와 소성", Description: &craftCourseDesc},
		{Title: "카페 브이로그 머신 사용법", Description: &coffeeDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(nodeBundle, levelBundle, courseBundle), 5)
	if len(filtered) != 3 {
		t.Fatalf("filtered length = %d, want 3", len(filtered))
	}
	if filtered[0].Title != "도자기 머그컵 손잡이 연결" {
		t.Fatalf("expected node match first, got %q", filtered[0].Title)
	}
	if filtered[1].Title != "도자기 머그컵 성형 입문" {
		t.Fatalf("expected level match second, got %q", filtered[1].Title)
	}
	if filtered[2].Title != "도자기 기초 재료와 소성" {
		t.Fatalf("expected course match third, got %q", filtered[2].Title)
	}
}

func TestFilterExplorerCandidatesByStage_RejectsUnrelatedSavedExternalEvenOnCourseFallback(t *testing.T) {
	nodeBundle := normalizer.BuildLessonSearchQuery("", "", "", "통기타 코드 전환", "", "")
	levelBundle := normalizer.BuildLessonSearchQuery("", "통기타 반주", "", "", "", "")
	courseBundle := normalizer.BuildLessonSearchQuery("통기타 배우기", "", "", "", "", "")

	guitarDesc := "통기타 반주에서 오픈 코드 전환과 스트로크를 연습합니다."
	baristaDesc := "에스프레소 머신 구조 이해와 카페 장비 사용법을 설명합니다."

	candidates := []ContentSearchCandidate{
		{Title: "통기타 코드 전환 연습", Description: &guitarDesc},
		{Title: "카페 장비 사용법 브이로그", Description: &baristaDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(nodeBundle, levelBundle, courseBundle), 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "통기타 코드 전환 연습" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByStage_WeakTokensOnlyDoNotReviveUnrelatedContent(t *testing.T) {
	nodeBundle := normalizer.BuildLessonSearchQuery("", "", "", "입문 이해 연습", "", "")
	levelBundle := normalizer.BuildLessonSearchQuery("", "기초 연습", "", "", "", "")

	genericDesc := "입문자를 위한 학습 이해와 연습 방법을 설명합니다."
	unrelatedDesc := "기초 연습 루틴과 입문 가이드를 설명합니다."

	candidates := []ContentSearchCandidate{
		{Title: "학습 입문 이해 가이드", Description: &genericDesc},
		{Title: "기초 연습 루틴 소개", Description: &unrelatedDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(nodeBundle, levelBundle), 5)
	if len(filtered) != 0 {
		t.Fatalf("weak tokens only should not revive candidates, got %d", len(filtered))
	}
}

func TestFilterExplorerCandidatesByStage_JapaneseTravelConversationRejectsSavedUnrelatedContent(t *testing.T) {
	levelBundle := normalizer.BuildLessonSearchQuery(
		"",
		"핵심 표현으로 짧은 여행 대화를 시작한다",
		"공항 식당 길 묻기처럼 자주 마주치는 여행 상황에서 핵심 표현으로 짧은 대화를 시작할 수 있게 한다",
		"",
		"",
		"",
	)
	courseBundle := normalizer.BuildLessonSearchQuery(
		"일본 여행 회화 준비 일본어 회화 일본 여행을 원활하게 준비한다",
		"",
		"",
		"",
		"",
		"",
	)

	japaneseDesc := "일본 식당 주문과 공항 길 묻기 상황에서 바로 쓸 수 있는 일본어 회화를 다룹니다."
	guitarDesc := "통기타 초기본기 연습과 오른손 피킹 감을 익히는 입문 영상입니다."
	baristaDesc := "에스프레소 머신 구조와 카페 장비 사용법을 설명합니다."

	candidates := []ContentSearchCandidate{
		{Title: "일본 식당에서 주문하는 일본어 회화", Description: &japaneseDesc},
		{Title: "기타 1일차 초기본기 연습", Description: &guitarDesc},
		{Title: "에스프레소 머신의 명칭과 사용법", Description: &baristaDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(normalizer.LessonSearchQuery{}, levelBundle, courseBundle), 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "일본 식당에서 주문하는 일본어 회화" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByStage_OnlineLectureMaterialRejectsSingleTokenSavedContent(t *testing.T) {
	nodeBundle := normalizer.BuildLessonSearchQuery(
		"",
		"",
		"",
		"온라인 강의 자료 제작",
		"",
		"",
	)
	levelBundle := normalizer.BuildLessonSearchQuery(
		"",
		"시장 조사 및 강의 주제 선정",
		"목표 시장을 이해하고 적합한 강의 주제를 결정할 수 있다",
		"",
		"",
		"",
	)
	courseBundle := normalizer.BuildLessonSearchQuery(
		"온라인 강의를 만들고 싶어 101클래스에 강의를 제공하여 수익을 창출한다",
		"",
		"",
		"",
		"",
		"",
	)

	lectureDesc := "온라인 강의 준비를 위해 파워포인트로 동영상 강의 자료를 만들고 제작 흐름을 설명합니다."
	guitarDesc := "기본 코드 마지막 강의로 기타 초보가 코드 잡는 요령을 익힙니다."

	candidates := []ContentSearchCandidate{
		{Title: "[온라인강의] 파워포인트로 동영상 강의자료 만들기", Description: &lectureDesc},
		{Title: "기타 기초 코드 강의", Description: &guitarDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(nodeBundle, levelBundle, courseBundle), 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "[온라인강의] 파워포인트로 동영상 강의자료 만들기" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByStage_UrbanSketchNodeStillRequiresCourseTheme(t *testing.T) {
	nodeBundle := normalizer.BuildLessonSearchQuery(
		"",
		"",
		"",
		"기본 표현으로 작은 그림을 완성한다",
		"",
		"선 형태 명암 색 같은 기본 표현을 사용해 작은 그림 한 점을 직접 완성할 수 있게 한다",
	)
	courseBundle := normalizer.BuildLessonSearchQuery(
		"어반스케치 그림 실력 향상 커리큘럼 어반스케치 그림 실력을 향상시키다",
		"",
		"",
		"",
		"",
		"",
	)

	sketchDesc := "어반스케치에서 선, 형태, 명암을 사용해 작은 그림 한 장면을 완성하는 연습입니다."
	japaneseDesc := "일본 식당에서 기본 표현으로 짧은 회화를 시작하고 주문 대화를 이어가는 방법입니다."

	candidates := []ContentSearchCandidate{
		{Title: "어반스케치 기본 표현으로 작은 그림 완성", Description: &sketchDesc},
		{Title: "일본 식당에서 기본 표현으로 주문하기", Description: &japaneseDesc},
	}

	filtered := filterExplorerCandidatesByStage(candidates, []explorerFilterStage{
		mustBuildExplorerFilterStage(t, "node", nodeBundle),
		mustBuildExplorerFilterStage(t, "course", courseBundle),
	}, 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "어반스케치 기본 표현으로 작은 그림 완성" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestFilterExplorerCandidatesByAllStages_RequiresIntersection(t *testing.T) {
	levelStage := mustBuildExplorerFilterStage(t, "level", normalizer.BuildLessonSearchQuery(
		"",
		"기본 표현으로 작은 그림을 완성한다",
		"선 형태 명암 색 같은 표현을 사용해 그림 한 점을 직접 완성할 수 있게 한다",
		"",
		"",
		"",
	))
	courseStage := mustBuildExplorerFilterStage(t, "course", normalizer.BuildLessonSearchQuery(
		"어반스케치 그림 실력을 향상시키다",
		"",
		"",
		"",
		"",
		"",
	))

	sketchDesc := "어반스케치 선과 형태, 명암으로 작은 그림 한 장면을 완성하는 드로잉 연습입니다."
	japaneseDesc := "기본 표현으로 일본 식당 주문 대화를 시작하는 회화 연습입니다."

	candidates := []ContentSearchCandidate{
		{Title: "어반스케치 작은 그림 완성", Description: &sketchDesc},
		{Title: "일본 식당 기본 표현", Description: &japaneseDesc},
	}

	filtered := filterExplorerCandidatesByAllStages(candidates, []explorerFilterStage{levelStage, courseStage}, 5)
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].Title != "어반스케치 작은 그림 완성" {
		t.Fatalf("unexpected surviving candidate: %q", filtered[0].Title)
	}
}

func TestBuildExplorerExternalSearchQuery_PrefersConcreteTravelTokens(t *testing.T) {
	stages := buildExplorerSearchStages(
		normalizer.LessonSearchQuery{},
		normalizer.BuildLessonSearchQuery(
			"",
			"핵심 표현으로 짧은 여행 대화를 시작한다",
			"공항 식당 길 묻기처럼 자주 마주치는 여행 상황에서 핵심 표현으로 짧은 대화를 시작할 수 있게 한다",
			"",
			"",
			"",
		),
		normalizer.BuildLessonSearchQuery(
			"일본 여행 회화 준비 일본어 회화 일본 여행을 원활하게 준비한다",
			"",
			"",
			"",
			"",
			"",
		),
	)

	query := buildExplorerExternalSearchQuery(stages, "핵심 표현으로 짧은 여행 대화를 시작한다 일본 여행 회화 준비")
	if !strings.HasPrefix(query, "일본") && !strings.HasPrefix(query, "일본어") {
		t.Fatalf("expected travel query to start with course-specific token, got %q", query)
	}
	if !strings.Contains(query, "공항") {
		t.Fatalf("expected travel query to keep concrete situation token, got %q", query)
	}
	if strings.Contains(query, "준비한다") || strings.Contains(query, "말하기") {
		t.Fatalf("unexpected generic course token in travel query: %q", query)
	}
}

func TestBuildExplorerExternalSearchQuery_PrefersCourseSpecificTokenForUrbanSketch(t *testing.T) {
	stages := buildExplorerSearchStages(
		normalizer.LessonSearchQuery{},
		normalizer.BuildLessonSearchQuery(
			"",
			"기본 표현으로 작은 그림을 완성한다",
			"선 형태 명암 색 같은 표현을 사용해 그림 한 점을 직접 완성할 수 있게 한다",
			"",
			"",
			"",
		),
		normalizer.BuildLessonSearchQuery(
			"어반스케치 그림 실력을 향상시키다",
			"",
			"",
			"",
			"",
			"",
		),
	)

	query := buildExplorerExternalSearchQuery(stages, "기본 표현으로 작은 그림을 완성한다 어반스케치 그림 실력 향상 커리큘럼")
	if !strings.HasPrefix(query, "어반스케치") {
		t.Fatalf("expected urban sketch query to start with course-specific token, got %q", query)
	}
	if strings.Contains(query, "실력") || strings.Contains(query, "향상") {
		t.Fatalf("unexpected generic course token in urban sketch query: %q", query)
	}
	if !strings.Contains(query, "작은") {
		t.Fatalf("expected urban sketch query to retain level token, got %q", query)
	}
}

func TestBuildExplorerExternalSearchQueryWithPrimaryQuery_UsesDirectQueryOnly(t *testing.T) {
	stages := buildExplorerSearchStages(
		normalizer.BuildLessonSearchQuery("", "", "", "시장 조사 및 강의 주제 선정", "", "수익화에 적합한 주제를 선정한다"),
		normalizer.BuildLessonSearchQuery("", "온라인 강의를 통한 수익 창출", "수익을 창출하고 싶다", "", "", ""),
		normalizer.BuildLessonSearchQuery("온라인 강의를 통한 수익 창출", "", "", "", "", ""),
	)

	query := buildExplorerExternalSearchQueryWithPrimaryQuery(stages, "온라인 강의 주제 선정 방법", "온라인 강의를 통한 수익 창출")
	if query != "온라인 강의 주제 선정 방법" {
		t.Fatalf("direct primary query should be preserved, got %q", query)
	}
	if strings.Contains(query, "수익") || strings.Contains(query, "시장") {
		t.Fatalf("direct query should not be mixed with context tokens, got %q", query)
	}
}

func TestRerankExplorerCandidatesWithPrimaryQuery_PrioritizesDirectQueryOverMonetizationContext(t *testing.T) {
	ctx := ExplorerRecommendationContext{
		SourceQuery:     "온라인 강의를 통한 수익 창출",
		LearningGoal:    "온라인 강의로 수익을 창출하고 싶다",
		RegionTitle:     "시장 조사 및 강의 주제 선정",
		RegionObjective: "수익화에 적합한 주제를 선정하고 시장을 분석한다",
		NodeTitle:       "시장 조사 및 강의 주제 선정",
		NodeSummary:     "수익성이 높은 강의 주제를 선택한다",
	}
	directDesc := "온라인 강의 주제 선정 방법과 교육 콘텐츠 기획 기준을 설명합니다."
	monetizationDesc := "온라인 강의 수익화와 마케팅 판매 전략으로 매출을 높이는 방법입니다."
	candidates := []ContentSearchCandidate{
		{Title: "온라인 강의 수익화 전략", Description: &monetizationDesc},
		{Title: "온라인 강의 주제 선정 방법", Description: &directDesc},
	}

	reranked := RerankExplorerRecommendationCandidatesWithPrimaryQuery(candidates, ctx, "온라인 강의 주제 선정 방법", 5)
	if len(reranked) != 2 {
		t.Fatalf("reranked length = %d, want 2", len(reranked))
	}
	if reranked[0].Title != "온라인 강의 주제 선정 방법" {
		t.Fatalf("direct-query candidate should rank first, got %q", reranked[0].Title)
	}
}

func TestNormalizeExplorerFilterToken_StripsParticles(t *testing.T) {
	tests := map[string]string{
		"어반스케치를": "어반스케치",
		"도시에서":   "도시",
		"연습":     "연습",
	}

	for input, want := range tests {
		if got := normalizeExplorerFilterToken(input); got != want {
			t.Fatalf("normalizeExplorerFilterToken(%q) = %q, want %q", input, got, want)
		}
	}
}

func mustBuildExplorerFilterStage(t *testing.T, name string, bundle normalizer.LessonSearchQuery) explorerFilterStage {
	t.Helper()
	stage, ok := buildExplorerFilterStage(name, bundle)
	if !ok {
		t.Fatalf("failed to build explorer filter stage: %s", name)
	}
	return stage
}
