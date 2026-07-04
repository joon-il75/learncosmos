package normalizer

import "testing"

func TestBuildLessonSearchQuery(t *testing.T) {
	got := BuildLessonSearchQuery(
		"기타 독학 블루스 솔로 입문",
		"기본 리듬과 펜타토닉",
		"기초부터 리듬을 익힌다",
		"마이너 펜타토닉 1포지션 익히기",
		"기타 솔로 연결을 위한 핵심 포지션 학습",
		"독학 기준으로 손가락 위치와 프레이즈 감각을 익힌다",
	)

	if got.RawQuery == "" || got.DenseQuery == "" {
		t.Fatalf("expected raw and dense query to be populated: %+v", got)
	}
	if got.LexicalQueryKO == "" || got.LexicalQueryExpanded == "" {
		t.Fatalf("expected lexical queries to be populated: %+v", got)
	}
	if len(got.Tokens) == 0 || len(got.ExpandedTokens) == 0 {
		t.Fatalf("expected tokens and expanded tokens: %+v", got)
	}
}

func TestBuildSearchTextKO(t *testing.T) {
	got := BuildSearchTextKO("포토 샵 기초", "일러", "ko")
	if got == "" {
		t.Fatal("expected normalized search text")
	}
	if got != "포토샵 기초 일러스트레이터 ko" {
		t.Fatalf("unexpected normalized search text: %q", got)
	}
}

func TestBuildLessonSearchQuerySplitsCommonCompoundLearningTerms(t *testing.T) {
	got := BuildLessonSearchQuery("시장조사 강의주제선정 수익창출", "", "", "", "", "")
	want := []string{"시장", "조사", "강의", "주제", "선정", "수익", "창출"}
	for _, token := range want {
		if !containsToken(got.Tokens, token) {
			t.Fatalf("expected token %q in %v", token, got.Tokens)
		}
	}
}

func containsToken(tokens []string, token string) bool {
	for _, item := range tokens {
		if item == token {
			return true
		}
	}
	return false
}
