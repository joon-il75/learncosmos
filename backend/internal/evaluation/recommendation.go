package evaluation

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

type RecommendationEvalCase struct {
	ID               string
	SourceQuery      string
	LevelTitle       string
	LevelObjective   string
	LessonTitle      string
	LessonObjective  string
	LessonSummary    string
	PreferredFormat  string
	ExpectedKeywords []string
	SearchSpec       curriculum.LessonRecommendationSearchSpec
}

type RecommendationEvalCaseResult struct {
	ID                  string   `json:"id"`
	RawQuery            string   `json:"raw_query"`
	DenseQuery          string   `json:"dense_query"`
	LexicalQuery        string   `json:"lexical_query"`
	LexicalHitCount     int      `json:"lexical_hit_count"`
	VectorHitCount      int      `json:"vector_hit_count"`
	MergedHitCount      int      `json:"merged_hit_count"`
	Top1Relevant        bool     `json:"top1_relevant"`
	Top5Relevant        bool     `json:"top5_relevant"`
	MatchedKeywords     []string `json:"matched_keywords"`
	TopTitles           []string `json:"top_titles,omitempty"`
	FailureReason       string   `json:"failure_reason,omitempty"`
	SpecQuery           string   `json:"spec_query,omitempty"`
	SpecSource          string   `json:"spec_source,omitempty"`
	SpecMergedHitCount  int      `json:"spec_merged_hit_count,omitempty"`
	SpecTop1Relevant    bool     `json:"spec_top1_relevant,omitempty"`
	SpecTop5Relevant    bool     `json:"spec_top5_relevant,omitempty"`
	SpecMatchedKeywords []string `json:"spec_matched_keywords,omitempty"`
	SpecTopTitles       []string `json:"spec_top_titles,omitempty"`
	SpecScoreApplied    bool     `json:"spec_score_applied,omitempty"`
	SpecImproved        bool     `json:"spec_improved,omitempty"`
}

type RecommendationEvalSummary struct {
	UserID                  string                         `json:"user_id,omitempty"`
	TotalCases              int                            `json:"total_cases"`
	CandidateHitRate        float64                        `json:"candidate_hit_rate"`
	EmptyResultRate         float64                        `json:"empty_result_rate"`
	Top1RelevanceRate       float64                        `json:"top1_relevance_rate"`
	Top5RelevanceRate       float64                        `json:"top5_relevance_rate"`
	LexicalRecall           float64                        `json:"lexical_recall"`
	VectorRecall            float64                        `json:"vector_recall"`
	SpecCaseCount           int                            `json:"spec_case_count"`
	SpecCandidateHitRate    float64                        `json:"spec_candidate_hit_rate"`
	SpecTop1RelevanceRate   float64                        `json:"spec_top1_relevance_rate"`
	SpecTop5RelevanceRate   float64                        `json:"spec_top5_relevance_rate"`
	SpecTop1ImprovementRate float64                        `json:"spec_top1_improvement_rate"`
	SpecImprovementRate     float64                        `json:"spec_improvement_rate"`
	Cases                   []RecommendationEvalCaseResult `json:"cases"`
}

var DefaultRecommendationEvalCases = []RecommendationEvalCase{
	{
		ID:               "leathercraft-wallet-stitching",
		SourceQuery:      "가죽공예 카드지갑 만들기",
		LevelTitle:       "카드지갑 제작 흐름 익히기",
		LevelObjective:   "카드지갑을 만들기 위한 재단, 바느질, 마감 흐름을 익힌다.",
		LessonTitle:      "가죽공예 카드지갑 재단과 바느질",
		LessonObjective:  "카드지갑 부품을 재단하고 새들스티치로 연결하는 기본 작업을 연습한다.",
		LessonSummary:    "가죽공예 카드지갑 제작 과정에서 패턴, 재단, 타공, 바느질을 순서대로 확인한다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"가죽공예", "카드지갑", "지갑"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "가죽공예 카드지갑 새들스티치 패턴 재단 강좌",
			Intent:       "tutorial",
			MustInclude:  []string{"가죽공예", "카드지갑", "새들스티치"},
			NiceToHave:   []string{"패턴", "재단", "바느질", "타공"},
			Avoid:        []string{"DIY키트", "키트", "완제품", "공방 모집", "원데이 클래스 광고"},
			ContentTypes: []string{"video"},
			Language:     "ko",
			StageRole:    "core_pattern",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
	{
		ID:               "vibe-coding-ai-studio-intro",
		SourceQuery:      "바이브코딩 앱 만들기",
		LevelTitle:       "AI 코딩 도구로 MVP 시작하기",
		LevelObjective:   "AI 코딩 도구를 활용해 작은 앱 MVP의 시작점을 만든다.",
		LessonTitle:      "바이브코딩 입문과 Google AI Studio 활용",
		LessonObjective:  "바이브코딩 흐름과 무료 AI 코딩 도구의 기본 사용법을 익힌다.",
		LessonSummary:    "AI Studio 같은 도구로 프롬프트를 입력하고 작은 앱 아이디어를 코드로 구체화한다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"바이브코딩", "AI Studio", "코딩"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "바이브코딩 AI Studio 입문 앱 만들기 튜토리얼",
			Intent:       "tutorial",
			MustInclude:  []string{"바이브코딩", "AI Studio", "코딩"},
			NiceToHave:   []string{"입문", "무료", "앱 만들기"},
			Avoid:        []string{"채용", "뉴스", "광고"},
			ContentTypes: []string{"video", "article"},
			Language:     "ko",
			StageRole:    "setup_intro",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
	{
		ID:               "urban-sketch-beginner-house",
		SourceQuery:      "어반스케치 기초",
		LevelTitle:       "왕초보 어반스케치 시작",
		LevelObjective:   "간단한 집과 풍경을 보며 선과 형태를 따라 그린다.",
		LessonTitle:      "초보자를 위한 집 그림 어반스케치",
		LessonObjective:  "정면 집 모양을 따라 그리며 어반스케치의 기본 구도를 익힌다.",
		LessonSummary:    "초보자가 쉽게 시작할 수 있는 집 그림과 간단한 풍경 스케치 과정을 따라 한다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"어반스케치", "초보", "그림"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "어반스케치 초보 집그림 기초 튜토리얼",
			Intent:       "tutorial",
			MustInclude:  []string{"어반스케치", "초보", "그림"},
			NiceToHave:   []string{"집그림", "기초", "풍경"},
			Avoid:        []string{"작품 판매", "전시 모집", "입시 미술", "유료 클래스 광고"},
			ContentTypes: []string{"video"},
			Language:     "ko",
			StageRole:    "first_output",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
	{
		ID:               "watercolor-control",
		SourceQuery:      "수채화 독학",
		LevelTitle:       "기초 붓 사용과 물 조절",
		LevelObjective:   "수채화에서 물 양과 붓 압력을 조절하는 감각을 익힌다.",
		LessonTitle:      "수채화 물 조절 연습",
		LessonObjective:  "번짐을 통제하면서 명암을 나누는 연습을 한다.",
		LessonSummary:    "물의 양과 붓의 수분량에 따라 색 번짐이 달라지는 차이를 비교한다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"수채화", "물조절", "붓", "번짐"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "수채화 물조절 붓 번짐 기초 튜토리얼",
			Intent:       "tutorial",
			MustInclude:  []string{"수채화", "물조절", "붓"},
			NiceToHave:   []string{"번짐", "기초", "연습"},
			Avoid:        []string{"작품 판매", "전시 모집", "입시 미술", "유료 클래스 광고"},
			ContentTypes: []string{"video"},
			Language:     "ko",
			StageRole:    "core_pattern",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
	{
		ID:               "pencil-shading",
		SourceQuery:      "연필 드로잉 입문",
		LevelTitle:       "형태와 명암의 기초",
		LevelObjective:   "기본 도형을 명암으로 입체감 있게 그린다.",
		LessonTitle:      "연필 드로잉 명암 기초",
		LessonObjective:  "밝기 단계를 나눠 명암 스케일을 만든다.",
		LessonSummary:    "구와 원기둥에 빛 방향을 설정하고 단계별 톤을 쌓는다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"드로잉", "명암", "연필", "스케치"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "드로잉 기초 명암 연필 스케치 튜토리얼",
			Intent:       "tutorial",
			MustInclude:  []string{"드로잉 기초", "명암", "연필"},
			NiceToHave:   []string{"스케치", "구도", "기초 연습"},
			Avoid:        []string{"작품 판매", "전시 모집", "입시 미술", "유료 클래스 광고"},
			ContentTypes: []string{"video", "article"},
			Language:     "ko",
			StageRole:    "core_pattern",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
	{
		ID:               "guitar-chord-transition",
		SourceQuery:      "통기타 독학",
		LevelTitle:       "기본 코드와 리듬",
		LevelObjective:   "자주 쓰는 기본 코드를 자연스럽게 연결한다.",
		LessonTitle:      "기타 코드 전환 연습",
		LessonObjective:  "느린 템포에서 코드 이동을 안정적으로 반복한다.",
		LessonSummary:    "기타 코드를 안정적으로 잡고 천천히 전환하는 연습 방법을 확인한다.",
		PreferredFormat:  "video",
		ExpectedKeywords: []string{"기타", "코드", "연습"},
		SearchSpec: curriculum.LessonRecommendationSearchSpec{
			PrimaryQuery: "기타 코드 전환 연습 기초 튜토리얼",
			Intent:       "tutorial",
			MustInclude:  []string{"기타", "코드", "연습"},
			NiceToHave:   []string{"전환", "기초", "통기타"},
			Avoid:        []string{"악기 판매", "공연 영상", "커버곡"},
			ContentTypes: []string{"video"},
			Language:     "ko",
			StageRole:    "core_pattern",
			Source:       curriculum.SearchSpecSourcePatternTemplate,
		},
	},
}

func RunRecommendationEvaluation(
	ctx context.Context,
	repo *curriculum.Repository,
	embedClient embedder.Embedder,
	cases []RecommendationEvalCase,
) (*RecommendationEvalSummary, error) {
	return RunRecommendationEvaluationForUser(ctx, repo, embedClient, uuid.Nil, cases)
}

func RunRecommendationEvaluationForUser(
	ctx context.Context,
	repo *curriculum.Repository,
	embedClient embedder.Embedder,
	userID uuid.UUID,
	cases []RecommendationEvalCase,
) (*RecommendationEvalSummary, error) {
	if len(cases) == 0 {
		cases = DefaultRecommendationEvalCases
	}

	results := make([]RecommendationEvalCaseResult, 0, len(cases))
	var candidateHitCases int
	var top1RelevantCases int
	var top5RelevantCases int
	var lexicalRecallCases int
	var vectorRecallCases int
	var specCaseCount int
	var specCandidateHitCases int
	var specTop1RelevantCases int
	var specTop5RelevantCases int
	var specTop1ImprovedCases int
	var specImprovedCases int

	for _, item := range cases {
		queryBundle := normalizer.BuildLessonSearchQuery(
			item.SourceQuery,
			item.LevelTitle,
			item.LevelObjective,
			item.LessonTitle,
			item.LessonObjective,
			item.LessonSummary,
		)

		lexicalQuery := queryBundle.LexicalQueryExpanded
		if strings.TrimSpace(lexicalQuery) == "" {
			lexicalQuery = queryBundle.LexicalQueryKO
		}

		preferredFormat := strings.TrimSpace(item.PreferredFormat)
		diagnostics, embeddingUsed, err := runRecommendationEvalSearch(ctx, repo, embedClient, userID, lexicalQuery, queryBundle.DenseQuery, preferredFormat)
		if err != nil {
			return nil, err
		}

		result := RecommendationEvalCaseResult{
			ID:              item.ID,
			RawQuery:        queryBundle.RawQuery,
			DenseQuery:      queryBundle.DenseQuery,
			LexicalQuery:    lexicalQuery,
			LexicalHitCount: len(diagnostics.LexicalCandidates),
			VectorHitCount:  len(diagnostics.VectorCandidates),
			MergedHitCount:  len(diagnostics.MergedCandidates),
		}

		if result.MergedHitCount > 0 {
			candidateHitCases++
		}
		if result.LexicalHitCount > 0 {
			lexicalRecallCases++
		}
		if result.VectorHitCount > 0 {
			vectorRecallCases++
		}

		result.Top1Relevant, _ = topNRelevant(diagnostics.MergedCandidates, item.ExpectedKeywords, 1)
		result.Top5Relevant, result.MatchedKeywords = top5Relevant(diagnostics.MergedCandidates, item.ExpectedKeywords)
		result.TopTitles = topCandidateTitles(diagnostics.MergedCandidates, 3)
		if result.Top1Relevant {
			top1RelevantCases++
		}
		if result.Top5Relevant {
			top5RelevantCases++
		}
		result.FailureReason = evaluationFailureReason(queryBundle, embeddingUsed, result.LexicalHitCount, result.VectorHitCount, result.MergedHitCount)

		if hasRecommendationSearchSpec(item.SearchSpec) {
			spec := curriculum.NormalizeLessonRecommendationSearchSpec(item.SearchSpec, curriculum.LessonRecommendationSearchSpec{
				Intent:       "tutorial",
				ContentTypes: []string{"video"},
				Language:     "ko",
				Source:       curriculum.SearchSpecSourcePatternTemplate,
			})
			specQuery := strings.TrimSpace(spec.PrimaryQuery)
			if specQuery == "" {
				specQuery = buildSpecQueryFromTokens(spec)
			}
			if specQuery != "" {
				specCaseCount++
				specDiagnostics, _, err := runRecommendationEvalSearch(ctx, repo, embedClient, userID, specQuery, specQuery, preferredFormat)
				if err != nil {
					return nil, err
				}
				specCandidates := rankRecommendationEvalCandidatesWithSpec(specDiagnostics.MergedCandidates, spec)
				result.SpecQuery = specQuery
				result.SpecSource = spec.Source
				result.SpecMergedHitCount = len(specCandidates)
				result.SpecTop1Relevant, _ = topNRelevant(specCandidates, item.ExpectedKeywords, 1)
				result.SpecTop5Relevant, result.SpecMatchedKeywords = top5Relevant(specCandidates, item.ExpectedKeywords)
				result.SpecTopTitles = topCandidateTitles(specCandidates, 3)
				result.SpecScoreApplied = true
				if result.SpecMergedHitCount > 0 {
					specCandidateHitCases++
				}
				if result.SpecTop1Relevant {
					specTop1RelevantCases++
				}
				if result.SpecTop5Relevant {
					specTop5RelevantCases++
				}
				if !result.Top1Relevant && result.SpecTop1Relevant {
					specTop1ImprovedCases++
				}
				result.SpecImproved = !result.Top5Relevant && result.SpecTop5Relevant
				if result.SpecImproved {
					specImprovedCases++
				}
			}
		}

		results = append(results, result)
	}

	total := float64(len(results))
	if total == 0 {
		total = 1
	}

	specTotal := float64(specCaseCount)
	if specTotal == 0 {
		specTotal = 1
	}

	return &RecommendationEvalSummary{
		UserID:                  recommendationEvalUserIDString(userID),
		TotalCases:              len(results),
		CandidateHitRate:        float64(candidateHitCases) / total,
		EmptyResultRate:         float64(len(results)-candidateHitCases) / total,
		Top1RelevanceRate:       float64(top1RelevantCases) / total,
		Top5RelevanceRate:       float64(top5RelevantCases) / total,
		LexicalRecall:           float64(lexicalRecallCases) / total,
		VectorRecall:            float64(vectorRecallCases) / total,
		SpecCaseCount:           specCaseCount,
		SpecCandidateHitRate:    float64(specCandidateHitCases) / specTotal,
		SpecTop1RelevanceRate:   float64(specTop1RelevantCases) / specTotal,
		SpecTop5RelevanceRate:   float64(specTop5RelevantCases) / specTotal,
		SpecTop1ImprovementRate: float64(specTop1ImprovedCases) / specTotal,
		SpecImprovementRate:     float64(specImprovedCases) / specTotal,
		Cases:                   results,
	}, nil
}

func runRecommendationEvalSearch(
	ctx context.Context,
	repo *curriculum.Repository,
	embedClient embedder.Embedder,
	userID uuid.UUID,
	lexicalQuery string,
	denseQuery string,
	preferredFormat string,
) (*curriculum.LessonCandidateSearchDiagnostics, bool, error) {
	var embeddingText *string
	if embedClient != nil && strings.TrimSpace(denseQuery) != "" {
		vector, err := embedClient.Embed(ctx, denseQuery)
		if err == nil && len(vector) > 0 {
			pgvec := float32VectorToPGVector(vector)
			embeddingText = &pgvec
		}
	}
	diagnostics, err := repo.SearchLessonCandidatesDiagnostics(ctx, userID, lexicalQuery, embeddingText, 5, curriculum.LessonCandidateSearchOptions{
		PreferredLanguage: "ko",
		PreferredFormat:   stringPointerOrNil(preferredFormat),
		Now:               time.Now(),
	})
	if err != nil {
		return nil, embeddingText != nil, err
	}
	return diagnostics, embeddingText != nil, nil
}

func hasRecommendationSearchSpec(spec curriculum.LessonRecommendationSearchSpec) bool {
	return strings.TrimSpace(spec.PrimaryQuery) != "" || len(spec.MustInclude) > 0 || len(spec.NiceToHave) > 0
}

func buildSpecQueryFromTokens(spec curriculum.LessonRecommendationSearchSpec) string {
	tokens := append([]string{}, spec.MustInclude...)
	tokens = append(tokens, spec.NiceToHave...)
	return strings.Join(firstNNonEmptyStrings(tokens, 5), " ")
}

func firstNNonEmptyStrings(values []string, limit int) []string {
	result := make([]string, 0, limit)
	seen := map[string]struct{}{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

func rankRecommendationEvalCandidatesWithSpec(candidates []curriculum.ContentSearchCandidate, spec curriculum.LessonRecommendationSearchSpec) []curriculum.ContentSearchCandidate {
	if len(candidates) == 0 {
		return nil
	}
	ranked := append([]curriculum.ContentSearchCandidate(nil), candidates...)
	sort.SliceStable(ranked, func(i, j int) bool {
		left := curriculum.ScoreCandidateWithLessonSearchSpec(ranked[i], spec)
		right := curriculum.ScoreCandidateWithLessonSearchSpec(ranked[j], spec)
		if left == right {
			return ranked[i].RankScore > ranked[j].RankScore
		}
		return left > right
	})
	return ranked
}

func topCandidateTitles(candidates []curriculum.ContentSearchCandidate, limit int) []string {
	if limit <= 0 || len(candidates) == 0 {
		return nil
	}
	if limit > len(candidates) {
		limit = len(candidates)
	}
	titles := make([]string, 0, limit)
	for idx := 0; idx < limit; idx++ {
		title := strings.TrimSpace(candidates[idx].Title)
		if title == "" {
			continue
		}
		titles = append(titles, title)
	}
	return titles
}

func top5Relevant(candidates []curriculum.ContentSearchCandidate, expectedKeywords []string) (bool, []string) {
	return topNRelevant(candidates, expectedKeywords, 5)
}

func topNRelevant(candidates []curriculum.ContentSearchCandidate, expectedKeywords []string, limit int) (bool, []string) {
	if len(candidates) == 0 || len(expectedKeywords) == 0 || limit <= 0 {
		return false, nil
	}
	if limit > len(candidates) {
		limit = len(candidates)
	}

	threshold := relevantKeywordThreshold(expectedKeywords)
	for idx := 0; idx < limit; idx++ {
		candidate := candidates[idx]
		candidateText := normalizer.BuildSearchTextKO(candidate.Title, valueOrEmpty(candidate.Description), candidate.ContentType, candidate.Language)
		matched := make([]string, 0, len(expectedKeywords))
		for _, keyword := range expectedKeywords {
			normalizedKeyword := normalizer.BuildSearchTextKO(keyword)
			if normalizedKeyword == "" {
				continue
			}
			if strings.Contains(candidateText, normalizedKeyword) {
				matched = append(matched, keyword)
			}
		}
		if len(matched) >= threshold {
			return true, matched
		}
	}

	return false, nil
}

func relevantKeywordThreshold(expectedKeywords []string) int {
	nonEmpty := 0
	for _, keyword := range expectedKeywords {
		if strings.TrimSpace(keyword) != "" {
			nonEmpty++
		}
	}
	if nonEmpty >= 2 {
		return 2
	}
	return 1
}

func evaluationFailureReason(query normalizer.LessonSearchQuery, embeddingUsed bool, lexicalHits, vectorHits, mergedHits int) string {
	if len(query.Tokens) == 0 {
		return "NO_QUERY_TOKENS"
	}
	if mergedHits > 0 {
		return ""
	}
	if !embeddingUsed {
		return "EMBEDDING_FAILED"
	}
	if lexicalHits == 0 && vectorHits == 0 {
		return "LEXICAL_ZERO_VECTOR_ZERO"
	}
	if lexicalHits == 0 && vectorHits > 0 {
		return "LEXICAL_ZERO_VECTOR_LOW"
	}
	return "SEED_CONTENT_SHORTAGE"
}

func stringPointerOrNil(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func float32VectorToPGVector(vector []float32) string {
	parts := make([]string, len(vector))
	for i, v := range vector {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func recommendationEvalUserIDString(userID uuid.UUID) string {
	if userID == uuid.Nil {
		return ""
	}
	return userID.String()
}
