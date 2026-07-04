package curriculum

import (
	"strings"
	"testing"
)

func TestBuildCurriculumPatternSeedsForLanguageIncludesEnglishOverrides(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	counts := map[string]int{}
	seen := map[string]struct{}{}
	var photography CurriculumPatternSeed
	for _, seed := range seeds {
		counts[seed.Language]++
		identity := seed.PatternKey + "\x00" + seed.PatternVersion + "\x00" + seed.Language
		if _, ok := seen[identity]; ok {
			t.Fatalf("duplicate seed identity: %s", identity)
		}
		seen[identity] = struct{}{}
		if seed.PatternKey == "visual_art:artifact_creation:photography_fundamentals_project" && seed.Language == "en" {
			photography = seed
		}
	}

	if counts["ko"] == 0 || counts["en"] == 0 || counts["ko"] != counts["en"] {
		t.Fatalf("language counts = %#v, want equal non-zero ko/en counts", counts)
	}
	if photography.PatternKey == "" {
		t.Fatal("English photography override seed not found")
	}

	if len(photography.RecommendationSearchSpecTemplate.StageTemplates) == 0 {
		t.Fatal("English photography seed search spec template is empty")
	}
	if photography.RecommendationSearchSpecTemplate.Source != SearchSpecSourcePatternTemplate {
		t.Fatalf("search spec template source = %q", photography.RecommendationSearchSpecTemplate.Source)
	}

	for _, want := range []string{
		"Smartphone Photography Fundamentals Project",
		"smartphone camera",
		"natural light",
		"mini photo project",
		"Completion means the learner can shoot",
	} {
		if !strings.Contains(photography.EmbeddingText, want) && !containsString(photography.TriggerKeywords, want) && !strings.Contains(photography.Summary, want) {
			t.Fatalf("English photography seed missing %q: %#v", want, photography)
		}
	}
	for _, forbidden := range []string{"요약", "권장", "단계", "완료 기준", "피할 리슨"} {
		if strings.Contains(photography.EmbeddingText, forbidden) || strings.Contains(photography.Summary, forbidden) {
			t.Fatalf("English photography seed leaked Korean phrase %q: %q", forbidden, photography.EmbeddingText)
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP0SearchSpecOverrides(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+"\x00"+seed.Language] = seed
	}

	leather := byKeyLang["craft_making:artifact_creation:leathercraft_wallet_project\x00ko"]
	if leather.PatternKey == "" {
		t.Fatal("Korean leathercraft wallet seed not found")
	}
	leatherSummary := summarizeCurriculumPatternSearchSpecTemplate(leather.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"가죽공예", "지갑", "새들스티치", "엣지마감"} {
		if !strings.Contains(leatherSummary, want) && !containsString(leather.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(leather.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("leather template missing %q: %v summary=%q", want, leather.RecommendationSearchSpecTemplate, leatherSummary)
		}
	}

	travel := byKeyLang["language_communication:performance_execution:travel_conversation\x00ko"]
	if travel.PatternKey == "" {
		t.Fatal("Korean travel conversation seed not found")
	}
	travelSummary := summarizeCurriculumPatternSearchSpecTemplate(travel.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"여행영어", "공항", "호텔", "역할극"} {
		if !strings.Contains(travelSummary, want) && !containsString(travel.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(travel.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("travel template missing %q: %v summary=%q", want, travel.RecommendationSearchSpecTemplate, travelSummary)
		}
	}
	for _, forbidden := range []string{"language", "communication", "도구", "재료"} {
		if strings.Contains(travelSummary, forbidden) || containsString(travel.RecommendationSearchSpecTemplate.MustInclude, forbidden) {
			t.Fatalf("travel template kept generic token %q: %v summary=%q", forbidden, travel.RecommendationSearchSpecTemplate, travelSummary)
		}
	}

	song := byKeyLang["instrument_performance:performance_execution:song_completion\x00ko"]
	if song.PatternKey == "" {
		t.Fatal("Korean song completion seed not found")
	}
	songSummary := summarizeCurriculumPatternSearchSpecTemplate(song.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"한 곡", "구간 연습", "리허설", "전체 연주"} {
		if !strings.Contains(songSummary, want) && !containsString(song.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(song.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("song template missing %q: %v summary=%q", want, song.RecommendationSearchSpecTemplate, songSummary)
		}
	}

	webapp := byKeyLang["digital_creation:artifact_creation:vibe_coding_mvp_app\x00ko"]
	if webapp.PatternKey == "" {
		t.Fatal("Korean web app MVP seed not found")
	}
	webappSummary := summarizeCurriculumPatternSearchSpecTemplate(webapp.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"웹앱", "MVP", "데이터베이스", "배포"} {
		if !strings.Contains(webappSummary, want) && !containsString(webapp.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(webapp.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("webapp template missing %q: %v summary=%q", want, webapp.RecommendationSearchSpecTemplate, webappSummary)
		}
	}

	webFoundation := byKeyLang["digital_creation:artifact_creation:web_dev_project_foundation\x00ko"]
	if webFoundation.PatternKey == "" {
		t.Fatal("Korean web development foundation seed not found")
	}
	webFoundationSummary := summarizeCurriculumPatternSearchSpecTemplate(webFoundation.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"웹개발", "HTML CSS JavaScript", "DOM 조작", "반응형 레이아웃"} {
		if !strings.Contains(webFoundationSummary, want) && !containsString(webFoundation.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(webFoundation.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("web foundation template missing %q: %v summary=%q", want, webFoundation.RecommendationSearchSpecTemplate, webFoundationSummary)
		}
	}
	for _, forbidden := range []string{"MVP", "풀스택", "AI 코딩 도구", "프롬프트 모음"} {
		if strings.Contains(webFoundationSummary, forbidden) || containsString(webFoundation.RecommendationSearchSpecTemplate.MustInclude, forbidden) || containsString(webFoundation.RecommendationSearchSpecTemplate.NiceToHave, forbidden) {
			t.Fatalf("web foundation template kept vibe-coding token %q: %v summary=%q", forbidden, webFoundation.RecommendationSearchSpecTemplate, webFoundationSummary)
		}
	}

	yoga := byKeyLang["body_movement:habit_lifestyle:yoga_daily_routine\x00ko"]
	if yoga.PatternKey == "" {
		t.Fatal("Korean yoga routine seed not found")
	}
	yogaSummary := summarizeCurriculumPatternSearchSpecTemplate(yoga.RecommendationSearchSpecTemplate, "ko")
	for _, want := range []string{"초보 요가", "요가 루틴", "요가 호흡", "15분 요가"} {
		if !strings.Contains(yogaSummary, want) && !containsString(yoga.RecommendationSearchSpecTemplate.MustInclude, want) && !containsString(yoga.RecommendationSearchSpecTemplate.NiceToHave, want) {
			t.Fatalf("yoga template missing %q: %v summary=%q", want, yoga.RecommendationSearchSpecTemplate, yogaSummary)
		}
	}
	for _, forbidden := range []string{"body movement", "habit lifestyle"} {
		if strings.Contains(yogaSummary, forbidden) || containsString(yoga.RecommendationSearchSpecTemplate.MustInclude, forbidden) {
			t.Fatalf("yoga template kept generic token %q: %v summary=%q", forbidden, yoga.RecommendationSearchSpecTemplate, yogaSummary)
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP1SearchSpecOverrides(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+"\x00"+seed.Language] = seed
	}

	cases := []struct {
		key       string
		wants     []string
		forbidden []string
		wantAvoid []string
	}{
		{
			key:       "visual_art:artifact_creation:structured_drawing_foundation",
			wants:     []string{"드로잉 기초", "어반스케치", "수채화 스케치", "명암"},
			forbidden: []string{"drawing foundation", "contour drawing", "미술", "작품"},
		},
		{
			key:       "cooking_baking:artifact_creation:cooking_basics_foundation",
			wants:     []string{"요리 기초", "칼질", "식품 위생", "재료 손질"},
			forbidden: []string{"cooking baking", "artifact creation"},
		},
		{
			key:       "cooking_baking:teaching_instruction:baking_class_demo",
			wants:     []string{"베이킹 클래스", "레시피 시연", "수업 흐름", "시연 리허설"},
			forbidden: []string{"cooking baking", "teaching instruction"},
		},
		{
			key:       "craft_making:teaching_instruction:craft_teach_youtube",
			wants:     []string{"공예 튜토리얼", "제작 과정", "유튜브 시연", "촬영"},
			forbidden: []string{"craft", "making"},
		},
		{
			key:       "writing_storytelling:presentation_publish:brunch_serial_publish",
			wants:     []string{"브런치 글쓰기", "연재 글", "발행 계획", "글 발행"},
			forbidden: []string{"writing storytelling", "presentation publish"},
		},
		{
			key:       "language_communication:certification_assessment:english_score_exam_certification",
			wants:     []string{"토익", "텝스", "영어 시험", "모의고사"},
			forbidden: []string{"회화", "표현"},
			wantAvoid: []string{"수학", "수학학원", "학원 모집"},
		},
		{
			key:       "cooking_baking:certification_assessment:baking_written_practical_certification",
			wants:     []string{"제과제빵 자격증", "필기", "실기", "채점 기준"},
			forbidden: []string{"cooking baking", "certification assessment"},
		},
		{
			key:       "cooking_baking:certification_assessment:korean_cooking_practical_certification",
			wants:     []string{"한식조리기능사", "실기", "시험 메뉴", "시간 관리"},
			forbidden: []string{"cooking baking", "certification assessment"},
		},
	}

	for _, tc := range cases {
		seed := byKeyLang[tc.key+"\x00ko"]
		if seed.PatternKey == "" {
			t.Fatalf("Korean seed not found: %s", tc.key)
		}
		template := seed.RecommendationSearchSpecTemplate
		summary := summarizeCurriculumPatternSearchSpecTemplate(template, "ko")
		for _, want := range tc.wants {
			if !strings.Contains(summary, want) && !curriculumPatternSearchSpecTemplateContains(template, want) {
				t.Fatalf("%s template missing %q: %v summary=%q", tc.key, want, template, summary)
			}
		}
		for _, forbidden := range tc.forbidden {
			if strings.Contains(summary, forbidden) || curriculumPatternSearchSpecTemplateContains(template, forbidden) {
				t.Fatalf("%s template kept bad token %q: %v summary=%q", tc.key, forbidden, template, summary)
			}
		}
		avoid := collectCurriculumPatternSearchSpecTemplateAvoidTerms(template)
		for _, want := range tc.wantAvoid {
			if !containsString(avoid, want) {
				t.Fatalf("%s template avoid missing %q: %v", tc.key, want, avoid)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP2SearchSpecOverrides(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+"\x00"+seed.Language] = seed
	}

	cases := []struct {
		key       string
		wants     []string
		forbidden []string
	}{
		{
			key:       "digital_creation:teaching_instruction:creator_tutorial_publish",
			wants:     []string{"교육 영상 제작", "튜토리얼 영상 편집", "화면 녹화 강의", "오디오 체크", "자막"},
			forbidden: []string{"강의 영상 짧은 강의 영상", "도구 재료", "조회수", "구독자"},
		},
		{
			key:       "digital_creation:presentation_publish:portfolio_publish",
			wants:     []string{"포트폴리오 공개", "작품 업로드", "프로젝트 페이지", "공개 체크리스트"},
			forbidden: []string{"영상편집 영상 편집", "쇼핑몰", "취업 포트폴리오", "입시 포트폴리오"},
		},
		{
			key:       "digital_creation:artifact_creation:freeware_midi_full_track",
			wants:     []string{"무료 DAW", "MIDI 작곡", "트랙 완성", "드럼 패턴"},
			forbidden: []string{"digital", "creation", "artifact", "freeware"},
		},
		{
			key:       "digital_creation:artifact_creation:notion_study_dashboard",
			wants:     []string{"노션 학습 대시보드", "노션 공부 관리", "할일 데이터베이스", "주간 리뷰"},
			forbidden: []string{"Notion study dashboard", "tasks database", "notes page"},
		},
		{
			key:       "knowledge_hobby:habit_lifestyle:personal_finance_budget_routine",
			wants:     []string{"가계부", "예산 관리", "지출 기록", "월간 예산"},
			forbidden: []string{"personal finance", "monthly budget", "track spending"},
		},
		{
			key:       "craft_making:artifact_creation:pottery_handbuilding_cup_project",
			wants:     []string{"도자기 핸드빌딩", "도자기 컵 만들기", "도자기 머그컵", "판 성형"},
			forbidden: []string{"공방 예약", "도자기 판매"},
		},
	}

	for _, tc := range cases {
		seed := byKeyLang[tc.key+"\x00ko"]
		if seed.PatternKey == "" {
			t.Fatalf("Korean seed not found: %s", tc.key)
		}
		template := seed.RecommendationSearchSpecTemplate
		summary := summarizeCurriculumPatternSearchSpecTemplate(template, "ko")
		for _, want := range tc.wants {
			if !strings.Contains(summary, want) && !curriculumPatternSearchSpecTemplateContains(template, want) {
				t.Fatalf("%s template missing %q: %v summary=%q", tc.key, want, template, summary)
			}
		}
		for _, forbidden := range tc.forbidden {
			if strings.Contains(summary, forbidden) || curriculumPatternSearchSpecTemplateContainsNonAvoid(template, forbidden) {
				t.Fatalf("%s template kept bad token %q: %v summary=%q", tc.key, forbidden, template, summary)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP3PrelaunchAuditOverrides(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+"\x00"+seed.Language] = seed
	}

	cases := []struct {
		key          string
		wants        []string
		contentTypes []string
		forbidden    []string
	}{
		{
			key:       "digital_creation:artifact_creation:freeware_midi_beatmaking",
			wants:     []string{"무료 DAW", "MIDI 비트메이킹", "드럼 패턴"},
			forbidden: []string{"digital creation", "artifact creation", "도구 재료"},
		},
		{
			key:          "digital_creation:professional_transition:korean_ncs_web_publisher_training",
			wants:        []string{"웹퍼블리셔", "HTML CSS", "반응형 웹"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"digital creation", "professional transition", "도구 재료"},
		},
		{
			key:          "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey",
			wants:        []string{"K-MOOC", "온라인 강좌", "수강 방법"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"knowledge hobby", "habit lifestyle", "도구 재료"},
		},
		{
			key:          "knowledge_hobby:professional_transition:credit_bank_standard_theory",
			wants:        []string{"학점은행제", "온라인 수업", "학점 인정"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"knowledge hobby", "professional transition", "도구 재료"},
		},
		{
			key:          "knowledge_hobby:teaching_instruction:credit_bank_instructional_design",
			wants:        []string{"교수설계", "학습목표", "수업지도안"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"knowledge hobby", "teaching instruction", "도구 재료"},
		},
		{
			key:       "maker_technical_hobby:teaching_instruction:maker_workshop_demo",
			wants:     []string{"메이커 워크숍", "아두이노", "DIY 프로토타입"},
			forbidden: []string{"maker technical hobby", "teaching instruction"},
		},
		{
			key:          "visual_art:presentation_publish:art_publish_showcase",
			wants:        []string{"작품 소개", "온라인 전시", "아트 포트폴리오"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"visual art", "presentation publish"},
		},
		{
			key:          "writing_storytelling:presentation_publish:university_studio",
			wants:        []string{"창작 글쓰기", "에세이 퇴고", "글쓰기 워크숍"},
			contentTypes: []string{"article"},
			forbidden:    []string{"writing storytelling", "presentation publish", "도구 재료"},
		},
	}

	for _, tc := range cases {
		seed := byKeyLang[tc.key+"\x00ko"]
		if seed.PatternKey == "" {
			t.Fatalf("Korean seed not found: %s", tc.key)
		}
		template := seed.RecommendationSearchSpecTemplate
		summary := summarizeCurriculumPatternSearchSpecTemplate(template, "ko")
		for _, want := range tc.wants {
			if !strings.Contains(summary, want) && !curriculumPatternSearchSpecTemplateContains(template, want) {
				t.Fatalf("%s template missing %q: %v summary=%q", tc.key, want, template, summary)
			}
		}
		for _, contentType := range tc.contentTypes {
			if !containsString(template.ContentTypes, contentType) {
				t.Fatalf("%s content_types missing %q: %v", tc.key, contentType, template.ContentTypes)
			}
		}
		for _, forbidden := range tc.forbidden {
			if strings.Contains(summary, forbidden) || curriculumPatternSearchSpecTemplateContainsNonAvoid(template, forbidden) {
				t.Fatalf("%s template kept bad token %q: %v summary=%q", tc.key, forbidden, template, summary)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP4BroadPrelaunchAuditTerms(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+"\x00"+seed.Language] = seed
	}

	cases := []struct {
		key          string
		wants        []string
		contentTypes []string
		forbidden    []string
	}{
		{
			key:       "cooking_baking:habit_lifestyle:meal_prep_weekday_lunch",
			wants:     []string{"밀프렙", "도시락 준비", "주중 도시락"},
			forbidden: []string{"meal prep healthy lunch", "batch cooking"},
		},
		{
			key:       "digital_creation:artifact_creation:animation_fundamentals_shot_progression",
			wants:     []string{"블렌더 애니메이션", "바운싱볼", "키프레임"},
			forbidden: []string{"animation fundamentals", "bouncing ball"},
		},
		{
			key:          "digital_creation:foundation_build:mooc_programming_foundation",
			wants:        []string{"프로그래밍 기초", "컴퓨터과학", "파이썬"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"programming computer science", "problem set"},
		},
		{
			key:          "digital_creation:professional_transition:frontend_competency_map",
			wants:        []string{"프론트엔드 역량", "웹 접근성", "반응형 디자인"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"front-end developer", "competency"},
		},
		{
			key:          "knowledge_hobby:concept_mastery:climate_data_graphing_explanation",
			wants:        []string{"기후 데이터", "기온 편차", "그래프 그리기"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"climate data", "temperature anomaly"},
		},
		{
			key:          "knowledge_hobby:foundation_build:open_textbook_chapter_practice",
			wants:        []string{"공개 교육자료", "오픈 교재", "챕터 연습"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"open textbook", "chapter practice"},
		},
		{
			key:       "knowledge_hobby:habit_lifestyle:dog_basic_training_routine",
			wants:     []string{"강아지 훈련", "긍정 강화", "리콜 훈련"},
			forbidden: []string{"dog basic training", "positive reinforcement"},
		},
		{
			key:       "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path",
			wants:     []string{"종이 회로", "카드보드 회로", "전자회로"},
			forbidden: []string{"cardboard circuits", "electronics robotics"},
		},
		{
			key:       "visual_art:artifact_creation:photography_studio_project",
			wants:     []string{"사진 스튜디오", "스튜디오 조명", "촬영 세팅"},
			forbidden: []string{"photography studio", "digital imaging"},
		},
		{
			key:          "writing_storytelling:presentation_publish:critical_reading_to_writing_draft",
			wants:        []string{"비평적 읽기", "창작 글쓰기", "글 초안"},
			contentTypes: []string{"article", "video"},
			forbidden:    []string{"creative writing", "critical reading"},
		},
	}

	for _, tc := range cases {
		seed := byKeyLang[tc.key+"\x00ko"]
		if seed.PatternKey == "" {
			t.Fatalf("Korean seed not found: %s", tc.key)
		}
		template := seed.RecommendationSearchSpecTemplate
		summary := summarizeCurriculumPatternSearchSpecTemplate(template, "ko")
		for _, want := range tc.wants {
			if !strings.Contains(summary, want) && !curriculumPatternSearchSpecTemplateContains(template, want) {
				t.Fatalf("%s template missing %q: %v summary=%q", tc.key, want, template, summary)
			}
		}
		for _, contentType := range tc.contentTypes {
			if !containsString(template.ContentTypes, contentType) {
				t.Fatalf("%s content_types missing %q: %v", tc.key, contentType, template.ContentTypes)
			}
		}
		for _, forbidden := range tc.forbidden {
			if strings.Contains(summary, forbidden) || curriculumPatternSearchSpecTemplateContainsNonAvoid(template, forbidden) {
				t.Fatalf("%s template kept bad token %q: %v summary=%q", tc.key, forbidden, template, summary)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP5AvoidQualityGuards(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	for _, seed := range seeds {
		template := seed.RecommendationSearchSpecTemplate
		avoid := strings.Join(collectCurriculumPatternSearchSpecTemplateAvoidTerms(template), " ")
		switch seed.GoalType {
		case "professional_transition":
			for _, want := range []string{"취업 보장", "상담 신청"} {
				if !strings.Contains(avoid, want) {
					t.Fatalf("%s professional_transition avoid missing %q: %q", seed.PatternKey, want, avoid)
				}
			}
		case "presentation_publish":
			if !strings.Contains(avoid, "조회수 보장") {
				t.Fatalf("%s presentation_publish avoid missing 조회수 보장: %q", seed.PatternKey, avoid)
			}
			if !strings.Contains(avoid, "취업 포트폴리오") && !strings.Contains(avoid, "대필") {
				t.Fatalf("%s presentation_publish avoid missing portfolio/ghostwriting guard: %q", seed.PatternKey, avoid)
			}
		case "certification_assessment":
			for _, want := range []string{"자격 대행", "합격 보장"} {
				if !strings.Contains(avoid, want) {
					t.Fatalf("%s certification_assessment avoid missing %q: %q", seed.PatternKey, want, avoid)
				}
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsGenericKoreanFallbackUsesSpecificTerms(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	byKey := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKey[seed.PatternKey] = seed
	}

	cases := []struct {
		key       string
		wants     []string
		forbidden []string
	}{
		{
			key:       "body_movement:performance_execution:dance_cover_song_routine",
			wants:     []string{"댄스 커버", "안무 배우기", "거울모드"},
			forbidden: []string{"body movement", "performance execution", "도구"},
		},
		{
			key:       "cooking_baking:artifact_creation:home_cafe_latte_foundation",
			wants:     []string{"홈카페", "라떼 만들기", "카페라떼"},
			forbidden: []string{"cooking baking", "artifact creation"},
		},
		{
			key:       "maker_technical_hobby:certification_assessment:drone_operator_basic_4class",
			wants:     []string{"드론 4종", "자격증", "안전 법규"},
			forbidden: []string{"maker technical hobby", "certification assessment"},
		},
		{
			key:       "visual_art:artifact_creation:sequential_creative_class_path",
			wants:     []string{"그림 독학", "일러스트 기초", "작품 만들기"},
			forbidden: []string{"creative class"},
		},
	}

	for _, tc := range cases {
		seed := byKey[tc.key]
		if seed.PatternKey == "" {
			t.Fatalf("Korean seed not found: %s", tc.key)
		}
		template := seed.RecommendationSearchSpecTemplate
		summary := summarizeCurriculumPatternSearchSpecTemplate(template, "ko")
		for _, want := range tc.wants {
			if !strings.Contains(summary, want) && !curriculumPatternSearchSpecTemplateContains(template, want) {
				t.Fatalf("%s template missing %q: %v summary=%q", tc.key, want, template, summary)
			}
		}
		for _, forbidden := range tc.forbidden {
			if strings.Contains(summary, forbidden) || curriculumPatternSearchSpecTemplateContains(template, forbidden) {
				t.Fatalf("%s template kept bad token %q: %v summary=%q", tc.key, forbidden, template, summary)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsSearchSpecQualityAudit(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	genericAxisTokens := []string{
		"artifact creation",
		"body movement",
		"certification assessment",
		"cooking baking",
		"craft making",
		"digital creation",
		"foundation build",
		"habit lifestyle",
		"instrument performance",
		"knowledge hobby",
		"language communication",
		"maker technical hobby",
		"performance execution",
		"presentation publish",
		"teaching instruction",
		"visual art",
		"writing storytelling",
	}
	splitGenericAxisTokens := map[string]bool{
		"artifact":      true,
		"assessment":    true,
		"body":          true,
		"build":         true,
		"certification": true,
		"communication": true,
		"craft":         true,
		"creation":      true,
		"digital":       true,
		"execution":     true,
		"foundation":    true,
		"habit":         true,
		"hobby":         true,
		"instruction":   true,
		"instrument":    true,
		"knowledge":     true,
		"language":      true,
		"lifestyle":     true,
		"maker":         true,
		"making":        true,
		"movement":      true,
		"participation": true,
		"performance":   true,
		"presentation":  true,
		"professional":  true,
		"publish":       true,
		"service":       true,
		"storytelling":  true,
		"teaching":      true,
		"technical":     true,
		"transition":    true,
		"visual":        true,
		"writing":       true,
	}
	certificationCoreTerms := []string{"시험", "자격증", "필기", "실기"}

	for _, seed := range seeds {
		template := seed.RecommendationSearchSpecTemplate
		if template.Source != SearchSpecSourcePatternTemplate {
			t.Fatalf("%s source = %q, want %q", seed.PatternKey, template.Source, SearchSpecSourcePatternTemplate)
		}
		if template.Language != "ko" {
			t.Fatalf("%s language = %q, want ko", seed.PatternKey, template.Language)
		}
		if len(template.ContentTypes) == 0 {
			t.Fatalf("%s content_types is empty", seed.PatternKey)
		}
		if len(template.StageTemplates) == 0 {
			t.Fatalf("%s stage_templates is empty", seed.PatternKey)
		}

		terms := collectCurriculumPatternSearchSpecTemplateTerms(template)
		joined := strings.ToLower(strings.Join(terms, "\n"))
		for _, token := range genericAxisTokens {
			if strings.Contains(joined, token) {
				t.Fatalf("%s search spec kept generic axis token %q: %q", seed.PatternKey, token, joined)
			}
		}

		splitAxisRootCount := 0
		for _, term := range template.MustInclude {
			if splitGenericAxisTokens[strings.ToLower(strings.TrimSpace(term))] {
				splitAxisRootCount++
			}
		}
		if splitAxisRootCount >= 2 {
			t.Fatalf("%s root must_include kept split generic axis tokens: %v", seed.PatternKey, template.MustInclude)
		}

		for idx, stage := range template.StageTemplates {
			if stage.StageRole == "" {
				t.Fatalf("%s stage[%d] role is empty", seed.PatternKey, idx)
			}
			if stage.Language != "ko" {
				t.Fatalf("%s stage[%d] language = %q, want ko", seed.PatternKey, idx, stage.Language)
			}
			if len(stage.MustInclude) == 0 {
				t.Fatalf("%s stage[%d:%s] must_include is empty", seed.PatternKey, idx, stage.StageRole)
			}
			if len(stage.ContentTypes) == 0 {
				t.Fatalf("%s stage[%d:%s] content_types is empty", seed.PatternKey, idx, stage.StageRole)
			}
			if query := buildPatternTemplatePrimaryQuery(stage, stage); strings.TrimSpace(query) == "" {
				t.Fatalf("%s stage[%d:%s] effective query would be empty", seed.PatternKey, idx, stage.StageRole)
			}
		}

		if strings.Contains(seed.PatternKey, ":certification_assessment:") {
			avoidTerms := collectCurriculumPatternSearchSpecTemplateAvoidTerms(template)
			for _, avoid := range avoidTerms {
				for _, core := range certificationCoreTerms {
					if strings.Contains(avoid, core) {
						t.Fatalf("%s avoid term %q contains certification core term %q", seed.PatternKey, avoid, core)
					}
				}
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP6EnglishGenericRootAudit(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "en")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}

	genericRootTerms := map[string]bool{
		"artifact": true, "assessment": true, "baking": true, "body": true, "build": true, "certification": true,
		"communication": true, "concept": true, "cooking": true, "craft": true, "creation": true, "digital": true,
		"execution": true, "foundation": true, "habit": true, "hobby": true, "instruction": true, "instrument": true,
		"knowledge": true, "language": true, "lifestyle": true, "maker": true, "making": true, "mastery": true,
		"movement": true, "participation": true, "performance": true, "presentation": true, "professional": true,
		"publish": true, "skill": true, "storytelling": true, "teaching": true, "technical": true,
		"transition": true, "visual": true, "writing": true,
	}

	for _, seed := range seeds {
		template := seed.RecommendationSearchSpecTemplate
		if template.Language != "en" {
			t.Fatalf("%s template language = %q, want en", seed.PatternKey, template.Language)
		}
		if len(template.MustInclude) >= 2 {
			first := strings.ToLower(strings.TrimSpace(template.MustInclude[0]))
			second := strings.ToLower(strings.TrimSpace(template.MustInclude[1]))
			if genericRootTerms[first] && genericRootTerms[second] {
				t.Fatalf("%s English root must_include starts with generic axis terms: %v", seed.PatternKey, template.MustInclude)
			}
		}
		if query := strings.Join(template.MustInclude, " "); strings.TrimSpace(query) == "" {
			t.Fatalf("%s English root effective query would be empty", seed.PatternKey)
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP7StageTermDeduplication(t *testing.T) {
	for _, language := range []string{"ko", "en"} {
		seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", language)
		if err != nil {
			t.Fatalf("BuildCurriculumPatternSeedsForLanguage(%q) error = %v", language, err)
		}
		for _, seed := range seeds {
			for _, stage := range seed.RecommendationSearchSpecTemplate.StageTemplates {
				must := make(map[string]struct{}, len(stage.MustInclude))
				for _, term := range stage.MustInclude {
					if trimmed := strings.ToLower(strings.TrimSpace(term)); trimmed != "" {
						must[trimmed] = struct{}{}
					}
				}
				for _, term := range stage.NiceToHave {
					trimmed := strings.ToLower(strings.TrimSpace(term))
					if trimmed == "" {
						continue
					}
					if _, found := must[trimmed]; found {
						t.Fatalf("%s %s stage %s duplicates term across must/nice: %q", language, seed.PatternKey, stage.StageRole, term)
					}
				}
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesP8NewCoveragePatterns(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}
	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+":"+seed.Language] = seed
	}

	tests := []struct {
		key      string
		language string
		want     []string
	}{
		{"body_movement:habit_lifestyle:home_strength_routine", "ko", []string{"홈트", "근력 루틴", "맨몸운동"}},
		{"body_movement:habit_lifestyle:home_strength_routine", "en", []string{"home workout", "strength routine", "bodyweight workout"}},
		{"instrument_performance:performance_execution:vocal_song_practice", "ko", []string{"보컬 연습", "노래 연습", "커버곡"}},
		{"instrument_performance:performance_execution:vocal_song_practice", "en", []string{"vocal practice", "singing practice", "cover song"}},
		{"digital_creation:artifact_creation:video_editing_shortform_project", "ko", []string{"영상 편집", "영상편집", "캡컷"}},
		{"digital_creation:artifact_creation:video_editing_shortform_project", "en", []string{"video editing", "youtube shorts", "cut editing"}},
		{"digital_creation:technical_skill:ai_tool_productivity_workflow", "ko", []string{"챗GPT", "문서 요약", "프롬프트"}},
		{"digital_creation:technical_skill:ai_tool_productivity_workflow", "en", []string{"AI productivity", "ChatGPT workflow", "document summary"}},
		{"craft_making:artifact_creation:candle_soap_resin_project", "ko", []string{"캔들 만들기", "수제비누", "레진아트"}},
		{"craft_making:artifact_creation:candle_soap_resin_project", "en", []string{"candle making", "soap making", "resin art"}},
	}

	for _, tt := range tests {
		seed, ok := byKeyLang[tt.key+":"+tt.language]
		if !ok {
			t.Fatalf("missing seed %s:%s", tt.key, tt.language)
		}
		for _, token := range tt.want {
			if !curriculumPatternSearchSpecTemplateContains(seed.RecommendationSearchSpecTemplate, token) {
				t.Fatalf("%s:%s search spec missing %q", tt.key, tt.language, token)
			}
		}
	}

	aiSeed := byKeyLang["digital_creation:technical_skill:ai_tool_productivity_workflow:ko"]
	aiSetup := curriculumPatternStageTemplateByRole(aiSeed.RecommendationSearchSpecTemplate, "setup_intro")
	if aiSetup.PrimaryQuery != "챗GPT 문서 요약 업무 자동화 실습" {
		t.Fatalf("AI productivity setup primary_query = %q", aiSetup.PrimaryQuery)
	}
	if containsString(aiSetup.MustInclude, "워크플로우") {
		t.Fatalf("AI productivity setup must_include kept broad workflow token: %v", aiSetup.MustInclude)
	}
	if !containsString(aiSetup.Avoid, "플로우 AI") {
		t.Fatalf("AI productivity setup avoid missing flow AI product guard: %v", aiSetup.Avoid)
	}

	if !containsString(aiSetup.ContentTypes, "article") || !containsString(aiSetup.ContentTypes, "video") {
		t.Fatalf("AI productivity setup content types = %v", aiSetup.ContentTypes)
	}

	craftSeed := byKeyLang["craft_making:artifact_creation:candle_soap_resin_project:ko"]
	craftSetup := curriculumPatternStageTemplateByRole(craftSeed.RecommendationSearchSpecTemplate, "setup_intro")
	if craftSetup.PrimaryQuery != "캔들 만들기 초보 재료 안전 튜토리얼" {
		t.Fatalf("candle setup primary_query = %q", craftSetup.PrimaryQuery)
	}
	if containsString(craftSetup.MustInclude, "비누 만들기") || containsString(craftSetup.MustInclude, "레진 안전") {
		t.Fatalf("candle setup mixed multiple material axes in must_include: %v", craftSetup.MustInclude)
	}
	craftFirst := curriculumPatternStageTemplateByRole(craftSeed.RecommendationSearchSpecTemplate, "first_output")
	if craftFirst.PrimaryQuery != "초보 레진 키링 만들기 튜토리얼" {
		t.Fatalf("candle first_output primary_query = %q", craftFirst.PrimaryQuery)
	}
}

func TestBuildCurriculumPatternSeedsIncludesP9AdditionalCoveragePatterns(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "all")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}
	byKeyLang := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKeyLang[seed.PatternKey+":"+seed.Language] = seed
	}

	tests := []struct {
		key      string
		language string
		want     []string
	}{
		{"visual_art:artifact_creation:watercolor_beginner_landscape", "ko", []string{"수채화", "풍경 엽서", "그라데이션"}},
		{"visual_art:artifact_creation:watercolor_beginner_landscape", "en", []string{"watercolor landscape", "watercolor postcard", "gradient"}},
		{"body_movement:habit_lifestyle:running_5k_beginner_routine", "ko", []string{"초보 러닝", "5km", "인터벌"}},
		{"body_movement:habit_lifestyle:running_5k_beginner_routine", "en", []string{"beginner running", "5K", "interval"}},
		{"language_communication:performance_execution:daily_english_conversation_routine", "ko", []string{"일상 영어", "영어 회화", "쉐도잉"}},
		{"language_communication:performance_execution:daily_english_conversation_routine", "en", []string{"daily English", "English conversation", "shadowing"}},
		{"maker_technical_hobby:artifact_creation:arduino_sensor_project", "ko", []string{"아두이노", "센서 프로젝트", "브레드보드"}},
		{"maker_technical_hobby:artifact_creation:arduino_sensor_project", "en", []string{"Arduino", "sensor project", "breadboard"}},
		{"cooking_baking:artifact_creation:korean_home_cooking_basics", "ko", []string{"집밥", "한식 기초", "된장찌개"}},
		{"cooking_baking:artifact_creation:korean_home_cooking_basics", "en", []string{"Korean home cooking", "banchan", "doenjang jjigae"}},
	}

	for _, tt := range tests {
		seed, ok := byKeyLang[tt.key+":"+tt.language]
		if !ok {
			t.Fatalf("missing seed %s:%s", tt.key, tt.language)
		}
		for _, token := range tt.want {
			if !curriculumPatternSearchSpecTemplateContains(seed.RecommendationSearchSpecTemplate, token) {
				t.Fatalf("%s:%s search spec missing %q", tt.key, tt.language, token)
			}
		}
	}

	runningSeed := byKeyLang["body_movement:habit_lifestyle:running_5k_beginner_routine:ko"]
	runningSetup := curriculumPatternStageTemplateByRole(runningSeed.RecommendationSearchSpecTemplate, "setup_intro")
	if containsString(runningSetup.Avoid, "5km") {
		t.Fatalf("running setup avoid contains core target: %v", runningSetup.Avoid)
	}
	for _, wantAvoid := range []string{"살 빼", "칼로리", "챌린지"} {
		if !containsString(runningSetup.Avoid, wantAvoid) {
			t.Fatalf("running setup avoid missing %q: %v", wantAvoid, runningSetup.Avoid)
		}
	}
	if !containsString(runningSetup.MustInclude, "초보 러닝") {
		t.Fatalf("running setup must_include = %v", runningSetup.MustInclude)
	}

	englishSeed := byKeyLang["language_communication:performance_execution:daily_english_conversation_routine:ko"]
	englishCore := curriculumPatternStageTemplateByRole(englishSeed.RecommendationSearchSpecTemplate, "core_pattern")
	for _, forbiddenTarget := range []string{"영어", "회화", "생활"} {
		if containsString(englishCore.Avoid, forbiddenTarget) {
			t.Fatalf("daily english avoid contains core target %q: %v", forbiddenTarget, englishCore.Avoid)
		}
	}
	for _, want := range []string{"생활 표현", "쉐도잉", "말하기 연습"} {
		if !containsString(englishCore.MustInclude, want) {
			t.Fatalf("daily english core must_include missing %q: %v", want, englishCore.MustInclude)
		}
	}
	for _, wantAvoid := range []string{"토익", "오픽", "영어학원", "성인영어학원", "직장인", "할인", "쿠폰", "멤버십"} {
		if !containsString(englishCore.Avoid, wantAvoid) {
			t.Fatalf("daily english core avoid missing %q: %v", wantAvoid, englishCore.Avoid)
		}
	}

	watercolorSeed := byKeyLang["visual_art:artifact_creation:watercolor_beginner_landscape:ko"]
	if !containsString(watercolorSeed.RecommendationSearchSpecTemplate.Avoid, "오일파스텔") || !containsString(watercolorSeed.RecommendationSearchSpecTemplate.Avoid, "메이크업") {
		t.Fatalf("watercolor avoid missing cross-medium guards: %v", watercolorSeed.RecommendationSearchSpecTemplate.Avoid)
	}
	arduinoSeed := byKeyLang["maker_technical_hobby:artifact_creation:arduino_sensor_project:ko"]
	if !containsString(arduinoSeed.RecommendationSearchSpecTemplate.Avoid, "키트 후기") || !containsString(arduinoSeed.RecommendationSearchSpecTemplate.Avoid, "특가") {
		t.Fatalf("arduino avoid missing shopping guards: %v", arduinoSeed.RecommendationSearchSpecTemplate.Avoid)
	}
	homeCookingSeed := byKeyLang["cooking_baking:artifact_creation:korean_home_cooking_basics:ko"]
	if !containsString(homeCookingSeed.RecommendationSearchSpecTemplate.Avoid, "맛집") || !containsString(homeCookingSeed.RecommendationSearchSpecTemplate.Avoid, "식당") {
		t.Fatalf("home cooking avoid missing restaurant guards: %v", homeCookingSeed.RecommendationSearchSpecTemplate.Avoid)
	}
	if !containsString(englishCore.Avoid, "과외") || !containsString(englishCore.Avoid, "내신") {
		t.Fatalf("daily english avoid missing tutoring guards: %v", englishCore.Avoid)
	}
}

func TestBuildCurriculumPatternSeedsIncludesP11LongTailSearchNoiseGuards(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}
	byKey := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKey[seed.PatternKey] = seed
	}

	tests := []struct {
		key   string
		avoid []string
	}{
		{
			key:   "knowledge_hobby:foundation_build:audit_course_practice_path",
			avoid: []string{"학점은행제", "과외", "편입", "수강신청", "토익"},
		},
		{
			key:   "knowledge_hobby:habit_lifestyle:public_lifelong_online_course",
			avoid: []string{"수강신청", "수강생 모집", "쿠폰", "수료증 출력"},
		},
		{
			key:   "knowledge_hobby:participation_service:blended_lifelong_participation",
			avoid: []string{"수강신청", "접수", "장소 확인", "수강료"},
		},
		{
			key:   "maker_technical_hobby:certification_assessment:equipment_operation_practical_certification",
			avoid: []string{"일정 안내", "응시 조건", "교재 판매"},
		},
		{
			key:   "knowledge_hobby:certification_assessment:environmental_engineering_certification",
			avoid: []string{"응시자격", "학점은행제", "취업 방향", "후기만"},
		},
		{
			key:   "digital_creation:artifact_creation:creative_coding_visual_project",
			avoid: []string{"코딩학원", "과외", "어원"},
		},
		{
			key:   "craft_making:artifact_creation:leathercraft_dimensional_bag",
			avoid: []string{"DIY키트", "완제품", "공방 수업"},
		},
		{
			key:   "digital_creation:artifact_creation:claude_code_agentic_workflow",
			avoid: []string{"성능 분석만", "뉴스", "도구 리뷰만"},
		},
	}

	for _, tt := range tests {
		seed, ok := byKey[tt.key]
		if !ok {
			t.Fatalf("missing seed %s", tt.key)
		}
		avoid := collectCurriculumPatternSearchSpecTemplateAvoidTerms(seed.RecommendationSearchSpecTemplate)
		for _, want := range tt.avoid {
			if !containsString(avoid, want) {
				t.Fatalf("%s avoid missing %q: %v", tt.key, want, avoid)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsIncludesArticleVideoSearchQualityExpansions(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}
	byKey := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKey[seed.PatternKey] = seed
	}

	for _, key := range []string{
		"language_communication:certification_assessment:english_score_exam_certification",
		"knowledge_hobby:habit_lifestyle:houseplant_repotting_care",
		"knowledge_hobby:habit_lifestyle:dog_basic_training_routine",
		"writing_storytelling:habit_lifestyle:daily_writing_habit",
		"digital_creation:artifact_creation:web_dev_project_foundation",
	} {
		seed, ok := byKey[key]
		if !ok {
			t.Fatalf("missing seed %s", key)
		}
		if !containsString(seed.RecommendationSearchSpecTemplate.ContentTypes, "article") || !containsString(seed.RecommendationSearchSpecTemplate.ContentTypes, "video") {
			t.Fatalf("%s root content_types = %v", key, seed.RecommendationSearchSpecTemplate.ContentTypes)
		}
		for _, stage := range seed.RecommendationSearchSpecTemplate.StageTemplates {
			if !containsString(stage.ContentTypes, "article") || !containsString(stage.ContentTypes, "video") {
				t.Fatalf("%s stage %s content_types = %v", key, stage.StageRole, stage.ContentTypes)
			}
		}
	}
}

func TestBuildCurriculumPatternSeedsPreservesSpecificTemplateStageContentTypes(t *testing.T) {
	seeds, err := BuildCurriculumPatternSeedsForLanguage("v1", "ko")
	if err != nil {
		t.Fatalf("BuildCurriculumPatternSeedsForLanguage() error = %v", err)
	}
	byKey := map[string]CurriculumPatternSeed{}
	for _, seed := range seeds {
		byKey[seed.PatternKey] = seed
	}

	tests := []struct {
		key  string
		want []string
	}{
		{"digital_creation:presentation_publish:portfolio_publish", []string{"article"}},
		{"writing_storytelling:presentation_publish:brunch_serial_publish", []string{"article", "video"}},
		{"knowledge_hobby:professional_transition:credit_bank_standard_theory", []string{"article", "video"}},
		{"knowledge_hobby:foundation_build:korean_open_lecture_weekly_survey", []string{"article", "video"}},
	}

	for _, tt := range tests {
		seed, ok := byKey[tt.key]
		if !ok {
			t.Fatalf("missing seed %s", tt.key)
		}
		for _, stage := range seed.RecommendationSearchSpecTemplate.StageTemplates {
			for _, want := range tt.want {
				if !containsString(stage.ContentTypes, want) {
					t.Fatalf("%s stage %s content types = %v, want %v", tt.key, stage.StageRole, stage.ContentTypes, tt.want)
				}
			}
			if len(stage.ContentTypes) != len(tt.want) {
				t.Fatalf("%s stage %s content types = %v, want exactly %v", tt.key, stage.StageRole, stage.ContentTypes, tt.want)
			}
		}
	}
}

func curriculumPatternStageTemplateByRole(template CurriculumPatternRecommendationSearchSpecTemplate, role string) LessonRecommendationSearchSpec {
	for _, stage := range template.StageTemplates {
		if stage.StageRole == role {
			return stage
		}
	}
	return LessonRecommendationSearchSpec{}
}

func curriculumPatternSearchSpecTemplateContains(template CurriculumPatternRecommendationSearchSpecTemplate, token string) bool {
	if containsString(template.MustInclude, token) || containsString(template.NiceToHave, token) || containsString(template.Avoid, token) || containsString(template.ContentTypes, token) {
		return true
	}
	for _, stage := range template.StageTemplates {
		if containsString(stage.MustInclude, token) || containsString(stage.NiceToHave, token) || containsString(stage.Avoid, token) || containsString(stage.ContentTypes, token) {
			return true
		}
	}
	return false
}

func collectCurriculumPatternSearchSpecTemplateTerms(template CurriculumPatternRecommendationSearchSpecTemplate) []string {
	terms := make([]string, 0, len(template.MustInclude)+len(template.NiceToHave)+len(template.Avoid)+len(template.ContentTypes))
	terms = append(terms, template.Intent, template.Language, template.Source)
	terms = append(terms, template.MustInclude...)
	terms = append(terms, template.NiceToHave...)
	terms = append(terms, template.Avoid...)
	terms = append(terms, template.ContentTypes...)
	for _, stage := range template.StageTemplates {
		terms = append(terms, stage.PrimaryQuery, stage.Intent, stage.Language, stage.StageRole, stage.Source)
		terms = append(terms, stage.MustInclude...)
		terms = append(terms, stage.NiceToHave...)
		terms = append(terms, stage.Avoid...)
		terms = append(terms, stage.ContentTypes...)
	}
	return terms
}

func collectCurriculumPatternSearchSpecTemplateAvoidTerms(template CurriculumPatternRecommendationSearchSpecTemplate) []string {
	terms := append([]string{}, template.Avoid...)
	for _, stage := range template.StageTemplates {
		terms = append(terms, stage.Avoid...)
	}
	return terms
}

func curriculumPatternSearchSpecTemplateContainsNonAvoid(template CurriculumPatternRecommendationSearchSpecTemplate, token string) bool {
	if containsString(template.MustInclude, token) || containsString(template.NiceToHave, token) || containsString(template.ContentTypes, token) {
		return true
	}
	for _, stage := range template.StageTemplates {
		if containsString(stage.MustInclude, token) || containsString(stage.NiceToHave, token) || containsString(stage.ContentTypes, token) {
			return true
		}
	}
	return false
}
