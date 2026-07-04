package curriculum

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMergeLessonCandidates(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	lexical := []ContentSearchCandidate{
		{ContentID: &id1, Title: "lexical-1"},
		{ContentID: &id2, Title: "lexical-2"},
	}
	vector := []ContentSearchCandidate{
		{ContentID: &id2, Title: "vector-duplicate"},
		{ContentID: &id3, Title: "vector-3"},
	}

	got := mergeLessonCandidates(lexical, vector, 5, LessonCandidateSearchOptions{})
	if len(got) != 3 {
		t.Fatalf("expected 3 merged candidates, got %d", len(got))
	}
	if got[0].ContentID == nil || *got[0].ContentID != id2 || got[1].ContentID == nil || *got[1].ContentID != id1 || got[2].ContentID == nil || *got[2].ContentID != id3 {
		t.Fatalf("unexpected merge order: %+v", got)
	}

	expected := 1.0/62.0 + 1.0/61.0
	if math.Abs(got[0].RankScore-expected) > 1e-9 {
		t.Fatalf("expected fused score %f, got %f", expected, got[0].RankScore)
	}
}

func TestMergeLessonCandidates_MetadataBoost(t *testing.T) {
	now := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	videoID := uuid.New()
	articleID := uuid.New()
	preferredFormat := "video"

	lexical := []ContentSearchCandidate{
		{
			ContentID:    &articleID,
			ContentType:  "article",
			Language:     "ko",
			QualityScore: 0.2,
			CreatedAt:    now.AddDate(0, -6, 0),
			RankScore:    0.9,
		},
		{
			ContentID:    &videoID,
			ContentType:  "youtube",
			Language:     "ko",
			QualityScore: 0.9,
			CreatedAt:    now.AddDate(0, 0, -7),
			RankScore:    0.7,
		},
	}

	got := mergeLessonCandidates(lexical, nil, 5, LessonCandidateSearchOptions{
		PreferredFormat:   &preferredFormat,
		PreferredLanguage: "ko",
		Now:               now,
	})

	if len(got) != 2 {
		t.Fatalf("expected 2 merged candidates, got %d", len(got))
	}
	if got[0].ContentID == nil || *got[0].ContentID != videoID {
		t.Fatalf("expected metadata-boosted video candidate first, got %+v", got)
	}
}

func TestLessonRetrievalLimit(t *testing.T) {
	if got := lessonRetrievalLimit(0); got != 30 {
		t.Fatalf("expected retrieval limit 30 for zero limit, got %d", got)
	}

	if got := lessonRetrievalLimit(5); got != 30 {
		t.Fatalf("expected retrieval limit 30 for small limit, got %d", got)
	}

	if got := lessonRetrievalLimit(30); got != 30 {
		t.Fatalf("expected retrieval limit 30 for exact limit, got %d", got)
	}

	if got := lessonRetrievalLimit(50); got != 50 {
		t.Fatalf("expected retrieval limit 50 for larger limit, got %d", got)
	}
}

func TestLexicalFallbackTokens(t *testing.T) {
	got := lexicalFallbackTokens("수채화 물조절 붓 번짐 기초 튜토리얼 연습")
	want := []string{"수채화", "물조절", "번짐"}
	if len(got) != len(want) {
		t.Fatalf("tokens = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tokens = %#v, want %#v", got, want)
		}
	}
	if threshold := lexicalFallbackThreshold(got); threshold != 2 {
		t.Fatalf("threshold = %d, want 2", threshold)
	}
}

func TestLexicalFallbackTokensSingleSpecificToken(t *testing.T) {
	got := lexicalFallbackTokens("기타 기초 연습")
	if len(got) != 1 || got[0] != "기타" {
		t.Fatalf("tokens = %#v, want [기타]", got)
	}
	if threshold := lexicalFallbackThreshold(got); threshold != 1 {
		t.Fatalf("threshold = %d, want 1", threshold)
	}
}
