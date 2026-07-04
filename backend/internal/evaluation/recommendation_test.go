package evaluation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func TestTop5Relevant(t *testing.T) {
	contentID := uuid.New()
	candidates := []curriculum.ContentSearchCandidate{
		{
			ContentID:   &contentID,
			Title:       "기타 코드 전환 기본 연습",
			ContentType: "youtube",
			Language:    "ko",
		},
	}

	ok, matched := top5Relevant(candidates, []string{"기타", "코드", "수채화"})
	if !ok {
		t.Fatal("expected relevance match")
	}
	if len(matched) != 2 {
		t.Fatalf("expected 2 matched keywords, got %d", len(matched))
	}
}

func TestEvaluationFailureReason(t *testing.T) {
	emptyQuery := normalizer.LessonSearchQuery{}
	if got := evaluationFailureReason(emptyQuery, false, 0, 0, 0); got != "NO_QUERY_TOKENS" {
		t.Fatalf("expected NO_QUERY_TOKENS, got %s", got)
	}

	query := normalizer.LessonSearchQuery{Tokens: []string{"기타"}}
	if got := evaluationFailureReason(query, false, 0, 0, 0); got != "EMBEDDING_FAILED" {
		t.Fatalf("expected EMBEDDING_FAILED, got %s", got)
	}
	if got := evaluationFailureReason(query, true, 0, 0, 0); got != "LEXICAL_ZERO_VECTOR_ZERO" {
		t.Fatalf("expected LEXICAL_ZERO_VECTOR_ZERO, got %s", got)
	}
	if got := evaluationFailureReason(query, true, 0, 2, 0); got != "LEXICAL_ZERO_VECTOR_LOW" {
		t.Fatalf("expected LEXICAL_ZERO_VECTOR_LOW, got %s", got)
	}
	if got := evaluationFailureReason(query, true, 1, 0, 0); got != "SEED_CONTENT_SHORTAGE" {
		t.Fatalf("expected SEED_CONTENT_SHORTAGE, got %s", got)
	}
	if got := evaluationFailureReason(query, true, 1, 1, 2); got != "" {
		t.Fatalf("expected empty failure reason, got %s", got)
	}
}

func TestRankRecommendationEvalCandidatesWithSpec(t *testing.T) {
	spec := curriculum.LessonRecommendationSearchSpec{
		PrimaryQuery: "드로잉 기초 명암 연필 스케치 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"드로잉", "명암", "연필"},
		NiceToHave:   []string{"스케치"},
		Avoid:        []string{"작품 판매"},
		ContentTypes: []string{"video"},
		Language:     "ko",
		StageRole:    "core_pattern",
		Source:       curriculum.SearchSpecSourcePatternTemplate,
	}
	candidates := []curriculum.ContentSearchCandidate{
		{Title: "연필 작품 판매 안내", ContentType: "blog", Language: "ko", RankScore: 10},
		{Title: "연필 드로잉 명암 기초 튜토리얼", ContentType: "youtube", Language: "ko", RankScore: 1},
	}

	ranked := rankRecommendationEvalCandidatesWithSpec(candidates, spec)
	if len(ranked) != 2 {
		t.Fatalf("ranked length = %d", len(ranked))
	}
	if ranked[0].Title != "연필 드로잉 명암 기초 튜토리얼" {
		t.Fatalf("expected spec-relevant candidate first, got %q", ranked[0].Title)
	}
}

func TestBuildSpecQueryFromTokens(t *testing.T) {
	spec := curriculum.LessonRecommendationSearchSpec{
		MustInclude: []string{"토익", "영어 시험", "토익"},
		NiceToHave:  []string{"문제풀이", "모의고사", "오답"},
	}
	if got := buildSpecQueryFromTokens(spec); got != "토익 영어 시험 문제풀이 모의고사 오답" {
		t.Fatalf("spec query = %q", got)
	}
}

func TestExternalSearchQualityFixturesPreferInstructionalCandidates(t *testing.T) {
	type fixture struct {
		id         string
		spec       curriculum.LessonRecommendationSearchSpec
		goodTitle  string
		noiseTitle string
	}

	fixtures := []fixture{
		{
			id: "leathercraft-wallet",
			spec: curriculum.LessonRecommendationSearchSpec{
				PrimaryQuery: "가죽공예 카드지갑 새들스티치 패턴 재단 강좌",
				Intent:       "tutorial",
				MustInclude:  []string{"가죽공예", "카드지갑", "새들스티치"},
				NiceToHave:   []string{"패턴", "재단", "바느질", "타공"},
				Avoid:        []string{"DIY키트", "키트", "완제품", "공방 모집"},
				ContentTypes: []string{"video"},
				Language:     "ko",
				StageRole:    "core_pattern",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			},
			goodTitle:  "가죽공예 카드지갑 새들스티치 패턴 재단 강좌",
			noiseTitle: "가죽공예 카드지갑 DIY키트 완제품 판매 후기",
		},
		{
			id: "watercolor-landscape-postcard",
			spec: curriculum.LessonRecommendationSearchSpec{
				PrimaryQuery: "초보 수채화 풍경 엽서 물조절 번짐 워시 과정",
				Intent:       "tutorial",
				MustInclude:  []string{"수채화", "물조절", "번짐"},
				NiceToHave:   []string{"풍경 엽서", "워시", "붓", "그라데이션"},
				Avoid:        []string{"오일파스텔", "메이크업", "키트"},
				ContentTypes: []string{"video"},
				Language:     "ko",
				StageRole:    "first_output",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			},
			goodTitle:  "초보 수채화 풍경 엽서 물조절 번짐 워시 과정",
			noiseTitle: "오일파스텔 수채화 느낌 메이크업 키트 후기",
		},
		{
			id: "daily-english-conversation",
			spec: curriculum.LessonRecommendationSearchSpec{
				PrimaryQuery: "생활 영어 표현 쉐도잉 말하기 연습 루틴",
				Intent:       "tutorial",
				MustInclude:  []string{"생활 영어", "쉐도잉", "말하기 연습"},
				NiceToHave:   []string{"일상 영어", "영어 회화", "필수 표현", "반복"},
				Avoid:        []string{"과외", "내신", "토익", "오픽", "영어학원", "성인영어학원", "비즈니스영어", "직장인"},
				ContentTypes: []string{"video", "article"},
				Language:     "ko",
				StageRole:    "core_pattern",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			},
			goodTitle:  "생활 영어 표현 쉐도잉 말하기 연습 루틴",
			noiseTitle: "직장인 비즈니스영어 성인영어학원 오픽 토익 과외",
		},
		{
			id: "arduino-sensor-project",
			spec: curriculum.LessonRecommendationSearchSpec{
				PrimaryQuery: "아두이노 센서 프로젝트 브레드보드 초보 튜토리얼",
				Intent:       "tutorial",
				MustInclude:  []string{"아두이노", "센서 프로젝트", "브레드보드"},
				NiceToHave:   []string{"회로", "코드"},
				Avoid:        []string{"키트 후기", "특가", "판매"},
				ContentTypes: []string{"video", "article"},
				Language:     "ko",
				StageRole:    "first_output",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			},
			goodTitle:  "아두이노 센서 프로젝트 브레드보드 회로 연결 초보 튜토리얼",
			noiseTitle: "아두이노 센서 키트 후기 특가 판매 모음",
		},
		{
			id: "beginner-5k-running",
			spec: curriculum.LessonRecommendationSearchSpec{
				PrimaryQuery: "5km 완주 초보 러닝 훈련 계획 걷기 달리기",
				Intent:       "tutorial",
				MustInclude:  []string{"5km", "초보러너", "완주"},
				NiceToHave:   []string{"훈련 계획", "걷기 달리기", "페이스", "루틴", "인터벌"},
				Avoid:        []string{"살 빼", "칼로리", "챌린지", "런닝머신", "체지방"},
				ContentTypes: []string{"video"},
				Language:     "ko",
				StageRole:    "integration_practice",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			},
			goodTitle:  "러닝 초보 프로그램 5KM 완주 달리기 루틴",
			noiseTitle: "살 빼는 러닝머신 칼로리 챌린지 후기",
		},
	}

	for _, item := range fixtures {
		good := curriculum.ContentSearchCandidate{Title: item.goodTitle, ContentType: "youtube", Language: "ko"}
		noise := curriculum.ContentSearchCandidate{Title: item.noiseTitle, ContentType: "naver_blog", Language: "ko"}
		goodScore := curriculum.ScoreCandidateWithLessonSearchSpec(good, item.spec)
		noiseScore := curriculum.ScoreCandidateWithLessonSearchSpec(noise, item.spec)
		if goodScore <= 0 {
			t.Fatalf("%s good score = %d, want positive", item.id, goodScore)
		}
		if goodScore <= noiseScore {
			t.Fatalf("%s expected instructional candidate above noise, good=%d noise=%d", item.id, goodScore, noiseScore)
		}
	}
}

func TestDefaultRecommendationEvalCasesUseSeedAlignedTopics(t *testing.T) {
	wantIDs := []string{
		"leathercraft-wallet-stitching",
		"vibe-coding-ai-studio-intro",
		"urban-sketch-beginner-house",
		"watercolor-control",
		"pencil-shading",
		"guitar-chord-transition",
	}
	if len(DefaultRecommendationEvalCases) != len(wantIDs) {
		t.Fatalf("default eval case count = %d, want %d", len(DefaultRecommendationEvalCases), len(wantIDs))
	}
	for i, wantID := range wantIDs {
		item := DefaultRecommendationEvalCases[i]
		if item.ID != wantID {
			t.Fatalf("case[%d] id = %q, want %q", i, item.ID, wantID)
		}
		if item.SearchSpec.PrimaryQuery == "" {
			t.Fatalf("case[%d] %s missing search spec primary query", i, item.ID)
		}
		if len(item.ExpectedKeywords) == 0 {
			t.Fatalf("case[%d] %s missing expected keywords", i, item.ID)
		}
	}
}
