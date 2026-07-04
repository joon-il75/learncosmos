package curriculum

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestBuildFallbackLessonRecommendationSearchSpecLeathercraftKorean(t *testing.T) {
	req := CreateCourseDraftRequest{
		SourceQuery:      "가죽공예를 기초부터 배워서 나만의 손 지갑을 만들어 보자",
		LearningLanguage: "ko",
	}
	lesson := generatedMainLesson{Title: "도구와 재료 익히기", Objective: "가죽공예 입문 도구를 익힌다"}

	spec := BuildFallbackLessonRecommendationSearchSpec(req, lesson, 0)

	if spec.PrimaryQuery != "가죽공예 입문 도구 사용법 기초 튜토리얼" {
		t.Fatalf("unexpected primary query: %q", spec.PrimaryQuery)
	}
	if spec.StageRole != "setup_intro" {
		t.Fatalf("unexpected stage role: %q", spec.StageRole)
	}
	if len(spec.MustInclude) < 2 {
		t.Fatalf("expected multiple must_include tokens, got %#v", spec.MustInclude)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecTravelEnglishKorean(t *testing.T) {
	req := CreateCourseDraftRequest{
		SourceQuery:      "필수 여행영어 120문장",
		LearningLanguage: "ko",
	}
	lesson := generatedMainLesson{Title: "공항과 호텔에서 바로 쓰는 표현", Objective: "체크인과 입국심사 표현을 익힌다"}

	spec := BuildFallbackLessonRecommendationSearchSpec(req, lesson, 0)

	if spec.PrimaryQuery != "여행 영어 공항 호텔 회화 튜토리얼" {
		t.Fatalf("unexpected primary query: %q", spec.PrimaryQuery)
	}
	if spec.Intent != "tutorial" {
		t.Fatalf("unexpected intent: %q", spec.Intent)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecUsesPatternTemplate(t *testing.T) {
	template, ok := buildSpecificCurriculumPatternRecommendationSearchSpecTemplate("body_movement:habit_lifestyle:yoga_daily_routine", "ko")
	if !ok {
		t.Fatal("expected yoga pattern template")
	}
	req := CreateCourseDraftRequest{
		SourceQuery:               "초보 요가 데일리 루틴",
		LearningLanguage:          "ko",
		PatternSearchSpecTemplate: template,
	}
	lesson := generatedMainLesson{Title: "몸을 안전하게 풀고 기본 흐름을 시작한다", Objective: "호흡과 가벼운 스트레칭으로 루틴을 시작한다"}
	spec := BuildFallbackLessonRecommendationSearchSpec(req, lesson, 0)
	if spec.Source != SearchSpecSourcePatternTemplate {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourcePatternTemplate)
	}
	if spec.StageRole != "setup_intro" {
		t.Fatalf("stage_role = %q", spec.StageRole)
	}
	if got := strings.Join(spec.MustInclude, ","); got != "초보 요가,호흡,안전한 자세" {
		t.Fatalf("must_include = %q", got)
	}
	if !strings.Contains(spec.PrimaryQuery, "초보 요가") || !strings.Contains(spec.PrimaryQuery, "호흡") {
		t.Fatalf("primary_query does not include template tokens: %q", spec.PrimaryQuery)
	}
}

func TestNormalizeGeneratedLessonRecommendationSearchSpecKeepsArtifactFinishTemplate(t *testing.T) {
	template, ok := buildSpecificCurriculumPatternRecommendationSearchSpecTemplate("craft_making:artifact_creation:leathercraft_wallet_project", "ko")
	if !ok {
		t.Fatal("expected leathercraft wallet template")
	}
	req := CreateCourseDraftRequest{
		SourceQuery:               "가죽공예 카드지갑 만들기",
		LearningLanguage:          "ko",
		PatternSearchSpecTemplate: template,
	}
	lesson := generatedMainLesson{
		Title:                    "엣지 마감과 수납 품질을 점검한다",
		Objective:                "엣지코트와 코바 마감 상태를 확인한다",
		RecommendationSearchSpec: LessonRecommendationSearchSpec{StageRole: "core_pattern"},
	}
	spec := normalizeGeneratedLessonRecommendationSearchSpec(req, lesson, 3)
	if spec.StageRole != "artifact_finish" {
		t.Fatalf("stage_role = %q, want artifact_finish", spec.StageRole)
	}
	if got := strings.Join(spec.MustInclude, ","); got != "가죽공예,지갑,엣지마감" {
		t.Fatalf("must_include = %q", got)
	}
	if !strings.Contains(spec.PrimaryQuery, "엣지마감") || !strings.Contains(spec.PrimaryQuery, "엣지코트") {
		t.Fatalf("primary_query does not include finish template tokens: %q", spec.PrimaryQuery)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecAppliesLearningIntentOverlay(t *testing.T) {
	template, ok := buildSpecificCurriculumPatternRecommendationSearchSpecTemplate("craft_making:artifact_creation:leathercraft_wallet_project", "ko")
	if !ok {
		t.Fatal("expected leathercraft wallet template")
	}
	req := CreateCourseDraftRequest{
		SourceQuery:               "가죽공예 카드지갑 만들기",
		LearningLanguage:          "ko",
		PatternSearchSpecTemplate: template,
		LearningIntent: LearningIntentProfile{
			DesiredOutput:       "손 지갑",
			CurrentBlockers:     []string{"도구 선택"},
			PreferredActivities: []string{"프로젝트 실습"},
			SuccessCriteria:     []string{"손 지갑 완성"},
		},
	}
	lesson := generatedMainLesson{Title: "패턴을 재단하고 첫 지갑 구조를 만든다", Objective: "카드지갑 패턴과 재단 흐름을 익힌다"}

	spec := BuildFallbackLessonRecommendationSearchSpec(req, lesson, 1)

	if spec.Source != SearchSpecSourcePatternTemplateOverlay {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourcePatternTemplateOverlay)
	}
	if !containsString(spec.MustInclude, "손 지갑") {
		t.Fatalf("must_include missing profile output: %v", spec.MustInclude)
	}
	if !strings.Contains(spec.PrimaryQuery, "손 지갑") {
		t.Fatalf("primary_query missing profile output: %q", spec.PrimaryQuery)
	}
	if containsString(spec.NiceToHave, "프로젝트 실습") {
		t.Fatalf("generic preferred activity should not pollute nice_to_have: %v", spec.NiceToHave)
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecSkipsGenericActivityHints(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "가죽공예 재단 카드지갑 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"가죽공예", "재단", "카드지갑"},
		NiceToHave:   []string{"패턴"},
		Language:     "ko",
		StageRole:    "core_pattern",
		Source:       SearchSpecSourcePatternTemplate,
	}

	got := ApplyLearningIntentProfileToLessonSearchSpec(spec, LearningIntentProfile{
		PreferredActivities: []string{"따라 만들기", "단계별", "쉬운"},
	})

	if !reflect.DeepEqual(got, spec) {
		t.Fatalf("generic activity hints should not change spec\ngot=%#v\nwant=%#v", got, spec)
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecAddsSpecificBlockerButSkipsGenericHints(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "가죽공예 입문 도구 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"가죽공예", "도구", "기초"},
		NiceToHave:   []string{"입문"},
		Language:     "ko",
		StageRole:    "setup_intro",
		Source:       SearchSpecSourcePatternTemplate,
	}

	got := ApplyLearningIntentProfileToLessonSearchSpec(spec, LearningIntentProfile{
		CurrentBlockers:     []string{"도구 선택"},
		PreferredActivities: []string{"따라 만들기"},
	})

	if got.Source != SearchSpecSourcePatternTemplateOverlay {
		t.Fatalf("source = %q, want %q", got.Source, SearchSpecSourcePatternTemplateOverlay)
	}
	if !containsString(got.NiceToHave, "도구 선택") {
		t.Fatalf("nice_to_have = %v, want 도구 선택", got.NiceToHave)
	}
	if containsString(got.NiceToHave, "따라 만들기") {
		t.Fatalf("generic activity hint polluted nice_to_have: %v", got.NiceToHave)
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecSkipsGeneratedAndManual(t *testing.T) {
	profile := LearningIntentProfile{DesiredOutput: "손 지갑"}
	for _, source := range []string{SearchSpecSourceGenerated, SearchSpecSourceManual} {
		spec := ApplyLearningIntentProfileToLessonSearchSpec(LessonRecommendationSearchSpec{Source: source, MustInclude: []string{"가죽공예"}}, profile)
		if spec.Source != source {
			t.Fatalf("source = %q, want %q", spec.Source, source)
		}
		if containsString(spec.MustInclude, "손 지갑") {
			t.Fatalf("%s spec should not be overlaid: %v", source, spec.MustInclude)
		}
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecUsesFallbackOverlaySource(t *testing.T) {
	spec := ApplyLearningIntentProfileToLessonSearchSpec(LessonRecommendationSearchSpec{
		PrimaryQuery: "여행 영어 공항 회화 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"여행영어", "공항"},
		NiceToHave:   []string{"회화"},
		Language:     "ko",
		StageRole:    "core_pattern",
		Source:       SearchSpecSourceFallback,
	}, LearningIntentProfile{DesiredOutput: "호텔 체크인 회화", SuccessCriteria: []string{"호텔 체크인 회화 완성"}})

	if spec.Source != SearchSpecSourceFallbackOverlay {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourceFallbackOverlay)
	}
	if !containsString(spec.MustInclude, "호텔 체크인 회화") {
		t.Fatalf("must_include = %v, want profile output", spec.MustInclude)
	}
}

func TestLessonRecommendationSearchSpecUnmarshalCoercesStringLists(t *testing.T) {
	raw := []byte(`{"primary_query":"초보 요가 루틴 튜토리얼","intent":"tutorial","must_include":"초보 요가, 호흡","nice_to_have":["15분 요가", "스트레칭/플로우"],"avoid":"광고; 고난도 자세","content_types":"video","language":"ko","stage_role":"setup_intro","source":"generated"}`)
	var spec LessonRecommendationSearchSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("unmarshal search spec: %v", err)
	}
	spec = NormalizeLessonRecommendationSearchSpec(spec, LessonRecommendationSearchSpec{})
	if got := strings.Join(spec.MustInclude, ","); got != "초보 요가,호흡" {
		t.Fatalf("must_include = %q", got)
	}
	if got := strings.Join(spec.NiceToHave, ","); got != "15분 요가,스트레칭,플로우" {
		t.Fatalf("nice_to_have = %q", got)
	}
	if got := strings.Join(spec.Avoid, ","); got != "광고,고난도 자세" {
		t.Fatalf("avoid = %q", got)
	}
	if got := strings.Join(spec.ContentTypes, ","); got != "video" {
		t.Fatalf("content_types = %q", got)
	}
}

func TestScoreCandidateWithLessonSearchSpecAvoidNegation(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "가죽공예 입문 도구 사용법 기초 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"가죽공예", "도구", "기초"},
		NiceToHave:   []string{"입문"},
		Avoid:        []string{"광고"},
		StageRole:    "setup_intro",
	}
	candidate := ContentSearchCandidate{
		Title: "가죽공예 기초 도구 입문 강좌 광고 하나도 없음",
	}

	score := ScoreCandidateWithLessonSearchSpec(candidate, spec)
	if score <= 0 {
		t.Fatalf("expected positive score without avoid penalty, got %d", score)
	}
}

func TestScoreCandidateWithLessonSearchSpecPrefersVideoContentType(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "가죽공예 입문 도구 사용법 기초 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"가죽공예", "도구", "기초"},
		NiceToHave:   []string{"입문"},
		ContentTypes: []string{"video"},
		StageRole:    "setup_intro",
	}
	blog := ContentSearchCandidate{
		Title:       "가죽공예 기초 강좌 도구 사용법",
		ContentType: "naver_blog",
	}
	youtube := ContentSearchCandidate{
		Title:       "[가죽공예 기초 강좌] 입문자를 위한 기초 도구 안내",
		ContentType: "youtube",
	}

	blogScore := ScoreCandidateWithLessonSearchSpec(blog, spec)
	youtubeScore := ScoreCandidateWithLessonSearchSpec(youtube, spec)
	if youtubeScore <= blogScore {
		t.Fatalf("expected youtube score > blog score, got youtube=%d blog=%d", youtubeScore, blogScore)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecForLessonUsesEditedTitle(t *testing.T) {
	goal := "손 지갑을 완성한다"
	objective := "카드지갑 바느질을 연습한다"
	spec := BuildFallbackLessonRecommendationSearchSpecForLesson(
		"가죽공예",
		&goal,
		"ko",
		"새들스티치로 카드지갑 바느질하기",
		&objective,
		2,
	)

	if spec.PrimaryQuery == "" {
		t.Fatal("expected fallback primary query")
	}
	if spec.Source != SearchSpecSourceFallback {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourceFallback)
	}
	if spec.StageRole != "core_pattern" {
		t.Fatalf("stage role = %q, want core_pattern", spec.StageRole)
	}
	if spec.PrimaryQuery != "가죽공예 바느질 재단 방법 기초 튜토리얼" {
		t.Fatalf("primary query = %q", spec.PrimaryQuery)
	}
	if got := strings.Join(spec.MustInclude, ","); got != "가죽공예,바느질,재단" {
		t.Fatalf("must_include = %q", got)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecForLessonWithProfileAppliesOverlay(t *testing.T) {
	goal := "가죽공예를 기초부터 배워서 카드지갑을 완성한다"
	objective := "도구와 재료를 준비한다"
	spec := BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(
		"가죽공예",
		&goal,
		"ko",
		"입문 도구와 재료 준비하기",
		&objective,
		0,
		LearningIntentProfile{
			DesiredOutput:       "카드지갑",
			PreferredActivities: []string{"따라 만들기"},
		},
	)

	if spec.Source != SearchSpecSourceFallbackOverlay {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourceFallbackOverlay)
	}
	if !containsString(spec.MustInclude, "카드지갑") {
		t.Fatalf("must_include = %#v, want 카드지갑", spec.MustInclude)
	}
	if !strings.Contains(spec.PrimaryQuery, "카드지갑") {
		t.Fatalf("primary query = %q, want 카드지갑", spec.PrimaryQuery)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecForLessonPreservesNoProfileBehavior(t *testing.T) {
	goal := "손 지갑을 완성한다"
	objective := "카드지갑 바느질을 연습한다"
	withoutProfile := BuildFallbackLessonRecommendationSearchSpecForLesson(
		"가죽공예",
		&goal,
		"ko",
		"새들스티치로 카드지갑 바느질하기",
		&objective,
		2,
	)
	withEmptyProfile := BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(
		"가죽공예",
		&goal,
		"ko",
		"새들스티치로 카드지갑 바느질하기",
		&objective,
		2,
		LearningIntentProfile{},
	)

	if !reflect.DeepEqual(withoutProfile, withEmptyProfile) {
		t.Fatalf("empty profile changed fallback spec\nwithout=%#v\nwith=%#v", withoutProfile, withEmptyProfile)
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecDoesNotChangeSourceWhenNoTokenAdded(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "악기 목표곡 초보 편곡 튜닝 기본 자세 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"악기", "목표곡", "초보", "편곡", "튜닝", "기본 자세"},
		NiceToHave:   []string{"반복 연습", "천천히", "박자", "리듬", "첫 소절", "인트로", "벌스", "템포"},
		Language:     "ko",
		StageRole:    "setup_intro",
		Source:       SearchSpecSourcePatternTemplate,
	}

	got := ApplyLearningIntentProfileToLessonSearchSpec(spec, LearningIntentProfile{
		DesiredOutput:       "연주",
		PreferredActivities: []string{"반복 연습"},
	})

	if !reflect.DeepEqual(got, spec) {
		t.Fatalf("expected unchanged spec when profile adds no usable token\ngot=%#v\nwant=%#v", got, spec)
	}
}

func TestApplyLearningIntentProfileToLessonSearchSpecSkipsSentenceLikeSuccessCriteria(t *testing.T) {
	spec := LessonRecommendationSearchSpec{
		PrimaryQuery: "곡 연습 구간 연결 리듬 튜토리얼",
		Intent:       "tutorial",
		MustInclude:  []string{"곡 연습", "구간 연결", "리듬"},
		NiceToHave:   []string{"반복 연습", "템포"},
		Language:     "ko",
		StageRole:    "core_pattern",
		Source:       SearchSpecSourcePatternTemplate,
	}

	got := ApplyLearningIntentProfileToLessonSearchSpec(spec, LearningIntentProfile{
		SuccessCriteria: []string{"3개월 안에 좋아하는 곡을 연주한다"},
	})

	if !reflect.DeepEqual(got, spec) {
		t.Fatalf("expected sentence-like success criteria to be skipped\ngot=%#v\nwant=%#v", got, spec)
	}
}

func TestBuildFallbackLessonRecommendationSearchSpecUsesCompactSourceQueryToken(t *testing.T) {
	spec := BuildFallbackLessonRecommendationSearchSpec(CreateCourseDraftRequest{
		SourceQuery:      "미디 음악을 만들고 싶어",
		LearningLanguage: "ko",
		LearningIntent: LearningIntentProfile{
			DesiredOutput: "미디 음악 만들기",
		},
	}, generatedMainLesson{Title: "짧은 MIDI beat를 완성하고 저장한다"}, 3)

	if containsString(spec.MustInclude, "미디 음악을 만들고 싶어") {
		t.Fatalf("must_include should not include sentence-like source query: %#v", spec.MustInclude)
	}
	if !containsString(spec.MustInclude, "미디") {
		t.Fatalf("must_include = %#v, want compact source token 미디", spec.MustInclude)
	}
	if !strings.Contains(spec.PrimaryQuery, "미디 음악 만들기") {
		t.Fatalf("primary query = %q, want desired output", spec.PrimaryQuery)
	}
	if containsString(spec.NiceToHave, "미디 음악 만들기") {
		t.Fatalf("nice_to_have should not duplicate must_include token: %#v", spec.NiceToHave)
	}
}

func TestLessonRecommendationSearchContextFallbackSpecAppliesLearningIntent(t *testing.T) {
	goal := "미디 음악을 만들고 싶어"
	ctx := LessonRecommendationSearchContext{
		SourceQuery:      "미디 음악을 만들고 싶어",
		LearningGoal:     &goal,
		LearningLanguage: "ko",
		LessonTitle:      "드럼 패턴과 첫 루프를 만든다",
		OrderIndex:       2,
		LearningIntent: LearningIntentProfile{
			DesiredOutput: "미디 음악 만들기",
		},
	}

	spec := ctx.FallbackSpec()
	if spec.Source != SearchSpecSourceFallbackOverlay {
		t.Fatalf("source = %q, want %q", spec.Source, SearchSpecSourceFallbackOverlay)
	}
	if !containsString(spec.MustInclude, "미디 음악 만들기") {
		t.Fatalf("must_include = %#v, want profile output", spec.MustInclude)
	}
	if containsString(spec.MustInclude, "미디 음악을 만들고 싶어") {
		t.Fatalf("must_include should not include sentence-like source query: %#v", spec.MustInclude)
	}
}

func TestParseLearningIntentProfileNormalizesAndIgnoresInvalidJSON(t *testing.T) {
	profile := ParseLearningIntentProfile([]byte(`{"desired_output":" 카드지갑 ","learner_level":"초보","preferred_activities":["따라 만들기","따라 만들기"]}`))
	if profile.DesiredOutput != "카드지갑" {
		t.Fatalf("desired output = %q, want 카드지갑", profile.DesiredOutput)
	}
	if profile.LearnerLevel != "beginner" {
		t.Fatalf("learner level = %q, want beginner", profile.LearnerLevel)
	}
	if len(profile.PreferredActivities) != 1 {
		t.Fatalf("preferred activities = %#v, want deduped single item", profile.PreferredActivities)
	}

	invalid := ParseLearningIntentProfile([]byte(`{"desired_output"`))
	if !reflect.DeepEqual(invalid, LearningIntentProfile{}) {
		t.Fatalf("invalid profile = %#v, want empty", invalid)
	}
}
