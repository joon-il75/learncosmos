package curriculum

import "testing"

func TestRerankExplorerCandidatesWithPrimaryQueryDemotesPhotoBookReview(t *testing.T) {
	candidates := []ContentSearchCandidate{
		{
			Title:       "[책리뷰] 인생샷을 만드는 스마트폰 사진구도와 보정법 TMI",
			ContentType: "naver_blog",
			RankScore:   0.097,
		},
		{
			Title:       "스마트폰으로 역광 사진 예쁘게 찍는 법",
			ContentType: "youtube",
			RankScore:   0.083,
		},
	}

	got := rerankExplorerCandidatesWithPrimaryQuery(candidates, nil, "스마트폰 사진 구도 빛 보정 방법", 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(got))
	}
	if got[0].Title == candidates[0].Title {
		t.Fatalf("expected book review candidate to be demoted, got first=%q", got[0].Title)
	}
}

func TestRerankExplorerCandidatesWithPrimaryQueryDemotesDroneCodingForFlightPractice(t *testing.T) {
	candidates := []ContentSearchCandidate{
		{
			Title:       "초보자도 쉬운 파이썬 드론 코딩 #1: 3초 만에 이륙하기",
			ContentType: "naver_blog",
			RankScore:   0.115,
		},
		{
			Title:       "드론 초보자를 위한 이륙부터 착륙까지 연습하는 방법",
			ContentType: "youtube",
			RankScore:   0.117,
		},
	}

	got := rerankExplorerCandidatesWithPrimaryQuery(candidates, nil, "드론 조종 기초 이륙 착륙 연습 방법", 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(got))
	}
	if got[0].Title == candidates[0].Title {
		t.Fatalf("expected drone coding candidate to be demoted, got first=%q", got[0].Title)
	}
}

func TestRerankExplorerCandidatesWithPrimaryQueryDemotesNonDroneAircraftForDronePractice(t *testing.T) {
	candidates := []ContentSearchCandidate{
		{
			Title:       "EPP 입문용 전투기 초보자도 4축 조종으로 쉽게 비행하는 방법",
			ContentType: "naver_blog",
			RankScore:   0.112,
		},
		{
			Title:       "드론 조종법 완벽 정리",
			ContentType: "youtube",
			RankScore:   0.117,
		},
	}

	got := rerankExplorerCandidatesWithPrimaryQuery(candidates, nil, "드론 조종 기초 이륙 착륙 연습 방법", 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(got))
	}
	if got[0].Title == candidates[0].Title {
		t.Fatalf("expected non-drone aircraft candidate to be demoted, got first=%q", got[0].Title)
	}
}
