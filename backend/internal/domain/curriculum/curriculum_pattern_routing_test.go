package curriculum

import (
	"strings"
	"testing"
)

func TestInferGoalPatternKey(t *testing.T) {
	goal := "브런치에 글을 연재하고 싶다"
	key := inferGoalPatternKey(CreateCourseDraftRequest{
		SourceQuery:  "에세이 글쓰기",
		LearningGoal: &goal,
	})

	if key.DomainAxis != "writing_storytelling" {
		t.Fatalf("domain axis = %q, want writing_storytelling", key.DomainAxis)
	}
	if key.GoalModePrimary != "presentation_publish" {
		t.Fatalf("goal mode primary = %q, want presentation_publish", key.GoalModePrimary)
	}
}

func TestInferGoalPatternKey_AmbiguousGoalTweaks(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantDomain  string
		wantPrimary string
	}{
		{
			name:        "visual art weak publish intent",
			sourceQuery: "그림 배우기",
			goal:        "작품을 올려보고 싶다",
			wantDomain:  "visual_art",
			wantPrimary: "presentation_publish",
		},
		{
			name:        "writing self reflection",
			sourceQuery: "글쓰기 배우기",
			goal:        "내 생각을 정리하고 싶다",
			wantDomain:  "writing_storytelling",
			wantPrimary: "habit_lifestyle",
		},
		{
			name:        "digital side hustle",
			sourceQuery: "영상편집 배우기",
			goal:        "부업으로 이어가고 싶다",
			wantDomain:  "digital_creation",
			wantPrimary: "professional_transition",
		},
		{
			name:        "cooking healthy family meal",
			sourceQuery: "요리 배우기",
			goal:        "가족에게 건강한 식사를 해주고 싶다",
			wantDomain:  "cooking_baking",
			wantPrimary: "habit_lifestyle",
		},
		{
			name:        "credit bank practicum",
			sourceQuery: "사회복지현장실습",
			goal:        "학점은행제로 사회복지현장실습을 이수하고 싶다",
			wantDomain:  "knowledge_hobby",
			wantPrimary: "professional_transition",
		},
		{
			name:        "credit bank instructional design",
			sourceQuery: "영유아교수방법론",
			goal:        "학점은행제로 영유아교수방법론을 통해 수업 설계를 배우고 싶다",
			wantDomain:  "knowledge_hobby",
			wantPrimary: "teaching_instruction",
		},
		{
			name:        "credit bank standard theory",
			sourceQuery: "인간행동과사회환경",
			goal:        "학점은행제로 인간행동과사회환경을 배우고 사회복지 기초를 다지고 싶다",
			wantDomain:  "knowledge_hobby",
			wantPrimary: "professional_transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := inferGoalPatternKey(CreateCourseDraftRequest{
				SourceQuery:  tt.sourceQuery,
				LearningGoal: &tt.goal,
			})
			if key.DomainAxis != tt.wantDomain {
				t.Fatalf("domain axis = %q, want %q", key.DomainAxis, tt.wantDomain)
			}
			if key.GoalModePrimary != tt.wantPrimary {
				t.Fatalf("goal mode primary = %q, want %q", key.GoalModePrimary, tt.wantPrimary)
			}
		})
	}
}

func TestInferGoalPatternKey_CertificationAssessmentDomains(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantDomain  string
	}{
		{
			name:        "jlpt n4",
			sourceQuery: "일본어 배우기",
			goal:        "JLPT N4 취득을 목표로 일본어를 공부한다",
			wantDomain:  "language_communication",
		},
		{
			name:        "aws cloud practitioner",
			sourceQuery: "클라우드 배우기",
			goal:        "AWS Certified Cloud Practitioner 자격증을 취득하고 싶다",
			wantDomain:  "digital_creation",
		},
		{
			name:        "forklift practical",
			sourceQuery: "지게차 배우기",
			goal:        "지게차운전기능사 취득을 목표로 공부하고 싶다",
			wantDomain:  "maker_technical_hobby",
		},
		{
			name:        "korean cooking practical",
			sourceQuery: "요리 배우기",
			goal:        "한식조리기능사 자격증을 취득하고 싶다",
			wantDomain:  "cooking_baking",
		},
		{
			name:        "makeup certification",
			sourceQuery: "메이크업 배우기",
			goal:        "미용사(메이크업) 자격증을 취득하고 싶다",
			wantDomain:  "visual_art",
		},
		{
			name:        "toeic certification",
			sourceQuery: "영어 배우기",
			goal:        "TOEIC 850점을 목표로 영어를 공부하고 싶다",
			wantDomain:  "language_communication",
		},
		{
			name:        "opic certification",
			sourceQuery: "영어 말하기 배우기",
			goal:        "OPIc IM2를 목표로 영어 말하기를 준비하고 싶다",
			wantDomain:  "language_communication",
		},
		{
			name:        "ielts certification",
			sourceQuery: "영어 배우기",
			goal:        "IELTS 6.5를 목표로 영어를 준비하고 싶다",
			wantDomain:  "language_communication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := inferGoalPatternKey(CreateCourseDraftRequest{
				SourceQuery:  tt.sourceQuery,
				LearningGoal: &tt.goal,
			})
			if key.DomainAxis != tt.wantDomain {
				t.Fatalf("domain axis = %q, want %q", key.DomainAxis, tt.wantDomain)
			}
			if key.GoalModePrimary != "certification_assessment" {
				t.Fatalf("goal mode primary = %q, want certification_assessment", key.GoalModePrimary)
			}
		})
	}
}

func TestInferGoalSubpatternKey(t *testing.T) {
	tests := []struct {
		name string
		req  CreateCourseDraftRequest
		want string
	}{
		{
			name: "worship team support",
			req: func() CreateCourseDraftRequest {
				goal := "찬양단 봉사에서 자연스럽게 반주하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "통기타 배우기", LearningGoal: &goal}
			}(),
			want: "instrument_performance:participation_service:worship_team_support",
		},
		{
			name: "song completion",
			req: func() CreateCourseDraftRequest {
				goal := "좋아하는 곡 한 곡을 끝까지 완주하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "피아노 배우기", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "cello song completion",
			req: func() CreateCourseDraftRequest {
				goal := "첼로 기본 자세와 활 쓰기를 익혀 쉬운 클래식 곡 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "첼로 배우기", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "bass guitar band song completion",
			req: func() CreateCourseDraftRequest {
				goal := "베이스 기타로 기본 리듬과 루트음을 익혀 밴드 합주 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "베이스 기타 합주 준비", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "electric guitar riff cover song driven",
			req: func() CreateCourseDraftRequest {
				goal := "좋아하는 록 곡의 리프와 솔로를 구간별로 카피해 한 곡을 커버한다"
				return CreateCourseDraftRequest{SourceQuery: "일렉기타 록 곡 카피", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "english guitar favorite song cover song driven",
			req: func() CreateCourseDraftRequest {
				goal := "I want to learn beginner guitar and play through a favorite song"
				return CreateCourseDraftRequest{SourceQuery: "beginner guitar song cover", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "bass band groove cover song driven",
			req: func() CreateCourseDraftRequest {
				goal := "좋아하는 밴드곡의 베이스 라인과 그루브를 익혀 원곡에 맞춰 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "베이스 기타 밴드곡 그루브", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "piano pop accompaniment cover song driven",
			req: func() CreateCourseDraftRequest {
				goal := "좋아하는 팝송의 코드 진행과 반주 패턴을 익혀 원곡 느낌으로 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "피아노 팝송 반주 카피", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "piano left hand chord accompaniment song driven",
			req: func() CreateCourseDraftRequest {
				goal := "피아노 왼손 반주와 코드 진행을 익혀 좋아하는 팝송을 안정적으로 반주한다"
				return CreateCourseDraftRequest{SourceQuery: "피아노 코드 반주 연습", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "drum pop groove cover song driven",
			req: func() CreateCourseDraftRequest {
				goal := "좋아하는 팝송의 드럼 그루브와 필인을 익혀 원곡에 맞춰 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "드럼 팝송 그루브 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name: "geomungo traditional melody completion",
			req: func() CreateCourseDraftRequest {
				goal := "거문고 기본 주법을 익혀 전통 가락 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "거문고 전통 가락 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "gayageum folk song completion",
			req: func() CreateCourseDraftRequest {
				goal := "가야금 기초 주법을 익혀 민요 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "가야금 민요 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "flute classical song completion",
			req: func() CreateCourseDraftRequest {
				goal := "플룻 호흡과 운지를 익혀 쉬운 클래식 곡 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "플룻 클래식 곡 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "saxophone jazz song completion",
			req: func() CreateCourseDraftRequest {
				goal := "색소폰 호흡과 운지를 익혀 쉬운 재즈 스탠더드 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "색소폰 재즈 곡 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "recorder classical song completion",
			req: func() CreateCourseDraftRequest {
				goal := "리코더 운지와 호흡을 익혀 쉬운 클래식 곡 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "리코더 클래식 곡 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "harmonica folk song completion",
			req: func() CreateCourseDraftRequest {
				goal := "하모니카 호흡과 벤딩 기초를 익혀 좋아하는 민요 한 곡을 연주한다"
				return CreateCourseDraftRequest{SourceQuery: "하모니카 민요 연주", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:song_completion",
		},
		{
			name: "brunch serial publish",
			req: func() CreateCourseDraftRequest {
				goal := "브런치에 글을 연재하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "에세이 글쓰기", LearningGoal: &goal}
			}(),
			want: "writing_storytelling:presentation_publish:brunch_serial_publish",
		},
		{
			name: "daily writing habit",
			req: func() CreateCourseDraftRequest {
				goal := "매일 짧게 기록하는 글쓰기 습관을 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "글쓰기 배우기", LearningGoal: &goal}
			}(),
			want: "writing_storytelling:habit_lifestyle:daily_writing_habit",
		},
		{
			name: "english daily writing habit",
			req: func() CreateCourseDraftRequest {
				goal := "I want to write short reflections every day and review them weekly"
				return CreateCourseDraftRequest{SourceQuery: "daily writing habit", LearningGoal: &goal}
			}(),
			want: "writing_storytelling:habit_lifestyle:daily_writing_habit",
		},
		{
			name: "kpop dance cover song routine",
			req: func() CreateCourseDraftRequest {
				goal := "기본 리듬과 동작 구간 연습을 통해 좋아하는 케이팝 안무 한 곡을 커버한다"
				return CreateCourseDraftRequest{SourceQuery: "케이팝 댄스 커버", LearningGoal: &goal}
			}(),
			want: "body_movement:performance_execution:dance_cover_song_routine",
		},
		{
			name: "morning yoga daily routine",
			req: func() CreateCourseDraftRequest {
				goal := "기본 호흡과 스트레칭 자세를 익혀 15분 아침 요가 루틴을 꾸준히 실천한다"
				return CreateCourseDraftRequest{SourceQuery: "아침 요가 루틴 만들기", LearningGoal: &goal}
			}(),
			want: "body_movement:habit_lifestyle:yoga_daily_routine",
		},
		{
			name: "english yoga daily routine",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a 15-minute morning yoga routine for flexibility and breathing"
				return CreateCourseDraftRequest{SourceQuery: "daily beginner yoga routine", LearningGoal: &goal}
			}(),
			want: "body_movement:habit_lifestyle:yoga_daily_routine",
		},
		{
			name: "english home workout routine",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a simple home workout routine and keep a safe weekly exercise habit"
				return CreateCourseDraftRequest{SourceQuery: "home workout beginner routine", LearningGoal: &goal}
			}(),
			want: "body_movement:habit_lifestyle:home_strength_routine",
		},
		{
			name: "korean home strength routine",
			req: func() CreateCourseDraftRequest {
				goal := "맨몸운동과 덤벨 기초로 20분 홈트 근력 루틴을 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "초보 홈트 근력 루틴", LearningGoal: &goal}
			}(),
			want: "body_movement:habit_lifestyle:home_strength_routine",
		},
		{
			name: "vocal song practice",
			req: func() CreateCourseDraftRequest {
				goal := "발성과 호흡, 음정을 잡아 좋아하는 노래 한 곡을 안정적으로 부르고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "보컬 노래 연습", LearningGoal: &goal}
			}(),
			want: "instrument_performance:performance_execution:vocal_song_practice",
		},
		{
			name: "shortform video editing project",
			req: func() CreateCourseDraftRequest {
				goal := "캡컷으로 컷편집과 자막을 익혀 쇼츠 영상 하나를 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "숏폼 영상 편집 초보", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:video_editing_shortform_project",
		},
		{
			name: "ai productivity workflow",
			req: func() CreateCourseDraftRequest {
				goal := "ChatGPT로 문서 요약과 반복 업무 자동화 워크플로우를 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "AI 생산성 업무 자동화", LearningGoal: &goal}
			}(),
			want: "digital_creation:technical_skill:ai_tool_productivity_workflow",
		},
		{
			name: "candle soap resin project",
			req: func() CreateCourseDraftRequest {
				goal := "향초와 수제비누 재료를 안전하게 다뤄 작은 캔들 소품을 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "캔들 만들기 비누 레진아트", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:candle_soap_resin_project",
		},
		{
			name: "watercolor beginner landscape",
			req: func() CreateCourseDraftRequest {
				goal := "수채화 번짐과 그라데이션을 익혀 작은 풍경 엽서를 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "초보 수채화 풍경 엽서", LearningGoal: &goal}
			}(),
			want: "visual_art:artifact_creation:watercolor_beginner_landscape",
		},
		{
			name: "beginner 5k running routine",
			req: func() CreateCourseDraftRequest {
				goal := "걷기 달리기와 인터벌로 초보 5km 러닝 루틴을 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "초보 러닝 5km 계획", LearningGoal: &goal}
			}(),
			want: "body_movement:habit_lifestyle:running_5k_beginner_routine",
		},
		{
			name: "daily english conversation routine",
			req: func() CreateCourseDraftRequest {
				goal := "일상 영어 표현을 쉐도잉하고 매일 짧게 말하는 회화 루틴을 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "일상 영어회화 루틴", LearningGoal: &goal}
			}(),
			want: "language_communication:performance_execution:daily_english_conversation_routine",
		},
		{
			name: "arduino sensor project",
			req: func() CreateCourseDraftRequest {
				goal := "아두이노와 브레드보드로 센서값을 읽고 LED나 서보모터를 움직이는 프로젝트를 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "아두이노 센서 프로젝트 초보", LearningGoal: &goal}
			}(),
			want: "maker_technical_hobby:artifact_creation:arduino_sensor_project",
		},
		{
			name: "korean home cooking basics",
			req: func() CreateCourseDraftRequest {
				goal := "된장찌개와 계란말이, 밑반찬을 익혀 집밥 한 끼를 직접 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "집밥 한식 기초 초보 요리", LearningGoal: &goal}
			}(),
			want: "cooking_baking:artifact_creation:korean_home_cooking_basics",
		},
		{
			name: "craft teach youtube",
			req: func() CreateCourseDraftRequest {
				goal := "나만의 작품을 만들고 유튜브를 통해 강의하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "코바늘 배우기", LearningGoal: &goal}
			}(),
			want: "craft_making:teaching_instruction:craft_teach_youtube",
		},
		{
			name: "woodworking small shelf project",
			req: func() CreateCourseDraftRequest {
				goal := "목공예 기본 도구 사용과 사포질, 마감을 익혀 작은 벽 선반 하나를 완성한다"
				return CreateCourseDraftRequest{SourceQuery: "목공예 작은 선반 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:woodworking_small_project",
		},
		{
			name: "woodworking cutting board project",
			req: func() CreateCourseDraftRequest {
				goal := "나무 재단과 샌딩, 오일 마감을 익혀 주방용 도마 하나를 직접 만든다"
				return CreateCourseDraftRequest{SourceQuery: "목공예 도마 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:woodworking_small_project",
		},
		{
			name: "english woodworking birdhouse project",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a simple wooden birdhouse with safe measuring, cutting, sanding, and assembly"
				return CreateCourseDraftRequest{SourceQuery: "beginner woodworking birdhouse", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:woodworking_small_project",
		},
		{
			name: "knitting scarf wearable project",
			req: func() CreateCourseDraftRequest {
				goal := "대바늘 겉뜨기와 안뜨기를 익혀 겨울 목도리 하나를 완성한다"
				return CreateCourseDraftRequest{SourceQuery: "대바늘 목도리 뜨기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:knitting_basic_wearable",
		},
		{
			name: "crochet small doll project",
			req: func() CreateCourseDraftRequest {
				goal := "코바늘 짧은뜨기와 원형뜨기를 익혀 작은 인형 하나를 완성한다"
				return CreateCourseDraftRequest{SourceQuery: "코바늘 인형 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:crochet_small_doll_project",
		},
		{
			name: "english crochet granny square project",
			req: func() CreateCourseDraftRequest {
				goal := "I want to crochet a neat granny square and join several squares into a small coaster or pouch"
				return CreateCourseDraftRequest{SourceQuery: "crochet granny square beginner", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:crochet_small_doll_project",
		},
		{
			name: "sewing eco bag project",
			req: func() CreateCourseDraftRequest {
				goal := "재봉틀 직선 박기와 원단 재단을 익혀 에코백 하나를 완성한다"
				return CreateCourseDraftRequest{SourceQuery: "재봉틀 에코백 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:sewing_project_foundation",
		},
		{
			name: "pottery mug handbuilding project",
			req: func() CreateCourseDraftRequest {
				goal := "흙 성형과 손잡이 붙이기, 유약 기초를 익혀 머그컵 하나를 완성한다"
				return CreateCourseDraftRequest{SourceQuery: "도자기 머그컵 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:pottery_handbuilding_cup_project",
		},
		{
			name: "english pottery glazing bowl project",
			req: func() CreateCourseDraftRequest {
				goal := "I want to handbuild a small bowl and test simple glazing choices before finishing it"
				return CreateCourseDraftRequest{SourceQuery: "pottery glazing small bowl", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:pottery_handbuilding_cup_project",
		},
		{
			name: "leathercraft basic accessory",
			req: func() CreateCourseDraftRequest {
				goal := "가죽공예를 처음 시작해서 작은 소품을 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "가죽공예 배우기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:leathercraft_basic_accessory",
		},
		{
			name: "leathercraft wallet project",
			req: func() CreateCourseDraftRequest {
				goal := "가죽공예로 카드지갑을 직접 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "가죽공예 카드지갑 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:leathercraft_wallet_project",
		},
		{
			name: "leathercraft dimensional bag",
			req: func() CreateCourseDraftRequest {
				goal := "가죽공예로 미니백과 클러치 같은 입체 소품을 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "가죽공예 가방 만들기", LearningGoal: &goal}
			}(),
			want: "craft_making:artifact_creation:leathercraft_dimensional_bag",
		},
		{
			name: "calligraphy basic lettering",
			req: func() CreateCourseDraftRequest {
				goal := "붓펜 캘리그라피를 처음 시작해서 예쁜 손글씨를 쓰고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "캘리그라피 배우기", LearningGoal: &goal}
			}(),
			want: "visual_art:artifact_creation:calligraphy_basic_lettering",
		},
		{
			name: "calligraphy quote art project",
			req: func() CreateCourseDraftRequest {
				goal := "캘리그라피로 명언 엽서 작품을 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "붓펜 캘리그라피 엽서", LearningGoal: &goal}
			}(),
			want: "visual_art:artifact_creation:calligraphy_quote_art_project",
		},
		{
			name: "calligraphy publish showcase",
			req: func() CreateCourseDraftRequest {
				goal := "캘리그라피 작품을 인스타에 공개하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "캘리그라피 작품 만들기", LearningGoal: &goal}
			}(),
			want: "visual_art:presentation_publish:calligraphy_publish_showcase",
		},
		{
			name: "art publish showcase",
			req: func() CreateCourseDraftRequest {
				goal := "그림을 전시와 인스타 업로드용으로 공개하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "드로잉 배우기", LearningGoal: &goal}
			}(),
			want: "visual_art:presentation_publish:art_publish_showcase",
		},
		{
			name: "travel conversation",
			req: func() CreateCourseDraftRequest {
				goal := "여행 가서 영어로 자연스럽게 말하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영어 회화 배우기", LearningGoal: &goal}
			}(),
			want: "language_communication:performance_execution:travel_conversation",
		},
		{
			name: "english travel conversation",
			req: func() CreateCourseDraftRequest {
				goal := "I want to speak naturally at airports, hotels, and restaurants while traveling"
				return CreateCourseDraftRequest{SourceQuery: "travel English conversation", LearningGoal: &goal}
			}(),
			want: "language_communication:performance_execution:travel_conversation",
		},
		{
			name: "english home baking cake",
			req: func() CreateCourseDraftRequest {
				goal := "I want to bake a simple sponge cake and decorate it with whipped cream"
				return CreateCourseDraftRequest{SourceQuery: "beginner home baking sponge cake", LearningGoal: &goal}
			}(),
			want: "cooking_baking:artifact_creation:cooking_basics_foundation",
		},
		{
			name: "english weekday meal prep lunch",
			req: func() CreateCourseDraftRequest {
				goal := "I want to plan, cook, and store simple healthy lunches for a weekday meal prep routine"
				return CreateCourseDraftRequest{SourceQuery: "healthy meal prep beginner lunch", LearningGoal: &goal}
			}(),
			want: "cooking_baking:habit_lifestyle:meal_prep_weekday_lunch",
		},
		{
			name: "english baking class demo",
			req: func() CreateCourseDraftRequest {
				goal := "I want to demonstrate a baking recipe in a small class"
				return CreateCourseDraftRequest{SourceQuery: "baking class recipe demo", LearningGoal: &goal}
			}(),
			want: "cooking_baking:teaching_instruction:baking_class_demo",
		},
		{
			name: "english balcony herb garden",
			req: func() CreateCourseDraftRequest {
				goal := "I want to grow herbs on my balcony and keep them healthy through a simple care routine"
				return CreateCourseDraftRequest{SourceQuery: "small balcony herb garden", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path",
		},
		{
			name: "english houseplant repotting care",
			req: func() CreateCourseDraftRequest {
				goal := "I want to repot a small houseplant safely, choose a basic soil mix, and monitor recovery after repotting"
				return CreateCourseDraftRequest{SourceQuery: "houseplant repotting soil mix", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:houseplant_repotting_care",
		},
		{
			name: "english personal finance budget routine",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a simple monthly budget, track spending categories, and review one week of expenses"
				return CreateCourseDraftRequest{SourceQuery: "personal finance monthly budget beginner", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:personal_finance_budget_routine",
		},
		{
			name: "english dog basic training routine",
			req: func() CreateCourseDraftRequest {
				goal := "I want to train my dog to sit, stay, and come using short positive reinforcement sessions"
				return CreateCourseDraftRequest{SourceQuery: "dog basic training sit stay recall", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:dog_basic_training_routine",
		},
		{
			name: "english notion study dashboard",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a simple Notion study dashboard with tasks, notes, and weekly review pages"
				return CreateCourseDraftRequest{SourceQuery: "Notion study dashboard beginner", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:notion_study_dashboard",
		},
		{
			name: "english notion productivity dashboard task tracker",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build a Notion productivity dashboard with a task tracker, study template, and weekly review database"
				return CreateCourseDraftRequest{SourceQuery: "Notion dashboard template beginner", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:notion_study_dashboard",
		},
		{
			name: "korean balcony herb garden",
			req: func() CreateCourseDraftRequest {
				goal := "햇빛과 물주기, 분갈이를 익혀 베란다에서 허브 화분을 꾸준히 관리한다"
				return CreateCourseDraftRequest{SourceQuery: "베란다 허브 키우기", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path",
		},
		{
			name: "korean houseplant repotting care",
			req: func() CreateCourseDraftRequest {
				goal := "분갈이 시기와 흙 배합, 물주기 기준을 익혀 반려식물을 건강하게 관리한다"
				return CreateCourseDraftRequest{SourceQuery: "반려식물 분갈이 관리", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:concept_mastery:horticulture_diagnostic_foundation",
		},
		{
			name: "korean kimchi fermentation preservation",
			req: func() CreateCourseDraftRequest {
				goal := "배추 절이기와 양념 비율, 발효 보관을 익혀 배추김치 한 통을 담근다"
				return CreateCourseDraftRequest{SourceQuery: "김치 담그기 발효 관리", LearningGoal: &goal}
			}(),
			want: "cooking_baking:habit_lifestyle:food_preservation_safety_path",
		},
		{
			name: "english office spreadsheet certification",
			req: func() CreateCourseDraftRequest {
				goal := "I want to prepare for an office spreadsheet certification exam"
				return CreateCourseDraftRequest{SourceQuery: "office spreadsheet certification", LearningGoal: &goal}
			}(),
			want: "digital_creation:certification_assessment:office_tool_practical_certification",
		},
		{
			name: "english counseling certification",
			req: func() CreateCourseDraftRequest {
				goal := "I want to prepare for a counseling certification exam"
				return CreateCourseDraftRequest{SourceQuery: "counseling certification exam", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:certification_assessment:counseling_certification",
		},
		{
			name: "english environmental engineering certification",
			req: func() CreateCourseDraftRequest {
				goal := "I want to prepare for an environmental engineering certification exam"
				return CreateCourseDraftRequest{SourceQuery: "environmental engineering certification exam", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:certification_assessment:environmental_engineering_certification",
		},
		{
			name: "english citizen science observation record",
			req: func() CreateCourseDraftRequest {
				goal := "I want to observe local plants, record them with iNaturalist, and review my biodiversity notes"
				return CreateCourseDraftRequest{SourceQuery: "citizen science biodiversity observation", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:citizen_science_observation_record",
		},
		{
			name: "english primary source inquiry note",
			req: func() CreateCourseDraftRequest {
				goal := "I want to analyze a historical primary source and write an observe-reflect-question inquiry note"
				return CreateCourseDraftRequest{SourceQuery: "primary source analysis note", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:concept_mastery:primary_source_inquiry_note",
		},
		{
			name: "english cardboard circuit invention",
			req: func() CreateCourseDraftRequest {
				goal := "I want to make a cardboard circuit invention with LEDs, switches, and a working prototype"
				return CreateCourseDraftRequest{SourceQuery: "cardboard circuit invention", LearningGoal: &goal}
			}(),
			want: "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path",
		},
		{
			name: "english design lab prototype project",
			req: func() CreateCourseDraftRequest {
				goal := "I want to frame a design problem and build a simple prototype for user feedback"
				return CreateCourseDraftRequest{SourceQuery: "design lab prototype project", LearningGoal: &goal}
			}(),
			want: "maker_technical_hobby:artifact_creation:design_lab_prototype_project",
		},
		{
			name: "english organic growing cycle",
			req: func() CreateCourseDraftRequest {
				goal := "I want to grow a small organic vegetable bed and follow the care cycle from seed to harvest"
				return CreateCourseDraftRequest{SourceQuery: "organic vegetable growing cycle", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:habit_lifestyle:organic_growing_cycle_plan",
		},
		{
			name: "english web app mvp",
			req: func() CreateCourseDraftRequest {
				goal := "I want to build and deploy a small full-stack web app with authentication and a database"
				return CreateCourseDraftRequest{SourceQuery: "build a simple web app MVP", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:web_dev_project_foundation",
		},
		{
			name: "english portfolio publish",
			req: func() CreateCourseDraftRequest {
				goal := "I want to edit one short video and publish it as a portfolio piece"
				return CreateCourseDraftRequest{SourceQuery: "publish a video editing portfolio", LearningGoal: &goal}
			}(),
			want: "digital_creation:presentation_publish:portfolio_publish",
		},
		{
			name: "portfolio publish",
			req: func() CreateCourseDraftRequest {
				goal := "편집한 영상을 포트폴리오와 유튜브에 공개하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영상편집 배우기", LearningGoal: &goal}
			}(),
			want: "digital_creation:presentation_publish:portfolio_publish",
		},
		{
			name: "video publish thumbnail does not route to nail beauty certification",
			req: func() CreateCourseDraftRequest {
				goal := "짧은 영상을 편집해 썸네일과 설명을 준비하고 유튜브에 공개한다"
				return CreateCourseDraftRequest{SourceQuery: "영상편집 유튜브 공개", LearningGoal: &goal}
			}(),
			want: "digital_creation:presentation_publish:portfolio_publish",
		},
		{
			name: "english podcast first episode publish",
			req: func() CreateCourseDraftRequest {
				goal := "I want to plan, record, edit, and publish a short first podcast episode with clear audio"
				return CreateCourseDraftRequest{SourceQuery: "beginner podcast first episode", LearningGoal: &goal}
			}(),
			want: "digital_creation:presentation_publish:portfolio_publish",
		},
		{
			name: "freeware midi beatmaking",
			req: func() CreateCourseDraftRequest {
				goal := "무료 BandLab으로 미디 비트를 처음 만들어보고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "무료 미디 음악 제작", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:freeware_midi_beatmaking",
		},
		{
			name: "freeware midi full track",
			req: func() CreateCourseDraftRequest {
				goal := "무료 DAW로 미디 드럼과 멜로디를 편곡해 짧은 트랙을 완성하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "BandLab MIDI track 만들기", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:freeware_midi_full_track",
		},
		{
			name: "midi foundation workflow",
			req: func() CreateCourseDraftRequest {
				goal := "MIDI 기초 개념과 컨트롤러 송수신 흐름을 이해하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "MIDI 101", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:midi_foundation_workflow",
		},
		{
			name: "vibe coding mvp app",
			req: func() CreateCourseDraftRequest {
				goal := "바이브 코딩으로 작은 앱 MVP를 처음 만들어보고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "vibe coding app 만들기", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:vibe_coding_mvp_app",
		},
		{
			name: "vibe coding fullstack ship",
			req: func() CreateCourseDraftRequest {
				goal := "Cursor와 Supabase, Vercel로 AI 코딩 full-stack 앱을 배포 직전까지 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "vibe coding full-stack", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:vibe_coding_fullstack_ship",
		},
		{
			name: "claude code agentic workflow",
			req: func() CreateCourseDraftRequest {
				goal := "Claude Code와 MCP를 사용해 agentic coding workflow를 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "Claude Code vibe coding", LearningGoal: &goal}
			}(),
			want: "digital_creation:artifact_creation:claude_code_agentic_workflow",
		},
		{
			name: "creator tutorial publish",
			req: func() CreateCourseDraftRequest {
				goal := "영상편집 과정을 튜토리얼로 설명하고 강의하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영상편집 배우기", LearningGoal: &goal}
			}(),
			want: "digital_creation:teaching_instruction:creator_tutorial_publish",
		},
		{
			name: "maker workshop demo",
			req: func() CreateCourseDraftRequest {
				goal := "아두이노 프로젝트 만드는 법을 워크숍에서 설명하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "아두이노 배우기", LearningGoal: &goal}
			}(),
			want: "maker_technical_hobby:teaching_instruction:maker_workshop_demo",
		},
		{
			name: "baking class demo",
			req: func() CreateCourseDraftRequest {
				goal := "베이킹 레시피를 설명하며 클래스를 해보고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "베이킹 배우기", LearningGoal: &goal}
			}(),
			want: "cooking_baking:teaching_instruction:baking_class_demo",
		},
		{
			name: "home cafe latte foundation",
			req: func() CreateCourseDraftRequest {
				goal := "에스프레소 추출과 우유 스티밍을 익혀 집에서 라떼를 안정적으로 만든다"
				return CreateCourseDraftRequest{SourceQuery: "홈카페 라떼 만들기", LearningGoal: &goal}
			}(),
			want: "cooking_baking:artifact_creation:home_cafe_latte_foundation",
		},
		{
			name: "jlpt exam targeted language",
			req: func() CreateCourseDraftRequest {
				goal := "JLPT N4 취득을 목표로 일본어를 공부한다"
				return CreateCourseDraftRequest{SourceQuery: "일본어 배우기", LearningGoal: &goal}
			}(),
			want: "language_communication:certification_assessment:exam_targeted_language",
		},
		{
			name: "cloud foundation certification",
			req: func() CreateCourseDraftRequest {
				goal := "AWS Certified Cloud Practitioner 자격증을 취득하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "클라우드 배우기", LearningGoal: &goal}
			}(),
			want: "digital_creation:certification_assessment:cloud_foundation_certification",
		},
		{
			name: "korean design tool practical certification",
			req: func() CreateCourseDraftRequest {
				goal := "GTQ 포토샵 2급 자격증을 실기 모의고사까지 준비하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "GTQ포토샵 기초부터 자격증까지", LearningGoal: &goal}
			}(),
			want: "digital_creation:certification_assessment:korean_design_tool_practical_certification",
		},
		{
			name: "korean ncs web publisher training",
			req: func() CreateCourseDraftRequest {
				goal := "NCS 디지털디자인 웹퍼블리셔 직업훈련을 통해 취업용 포트폴리오를 만들고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "HRD-Net 웹퍼블리셔 프론트엔드 직업훈련", LearningGoal: &goal}
			}(),
			want: "digital_creation:professional_transition:korean_ncs_web_publisher_training",
		},
		{
			name: "equipment operation practical certification",
			req: func() CreateCourseDraftRequest {
				goal := "지게차운전기능사 취득을 목표로 공부하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "지게차 배우기", LearningGoal: &goal}
			}(),
			want: "maker_technical_hobby:certification_assessment:equipment_operation_practical_certification",
		},
		{
			name: "korean cooking practical certification",
			req: func() CreateCourseDraftRequest {
				goal := "한식조리기능사 자격증을 취득하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "요리 배우기", LearningGoal: &goal}
			}(),
			want: "cooking_baking:certification_assessment:korean_cooking_practical_certification",
		},
		{
			name: "beauty practical certification",
			req: func() CreateCourseDraftRequest {
				goal := "미용사(메이크업) 자격증을 취득하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "메이크업 배우기", LearningGoal: &goal}
			}(),
			want: "visual_art:certification_assessment:beauty_service_practical_certification",
		},
		{
			name: "credit bank practicum",
			req: func() CreateCourseDraftRequest {
				goal := "학점은행제로 사회복지현장실습을 이수하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "사회복지현장실습", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:professional_transition:credit_bank_practicum",
		},
		{
			name: "korean education practicum",
			req: func() CreateCourseDraftRequest {
				goal := "외국어로서의한국어교육실습을 통해 실제 수업에 설 준비를 하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "외국어로서의한국어교육실습", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:professional_transition:credit_bank_practicum",
		},
		{
			name: "credit bank instructional design",
			req: func() CreateCourseDraftRequest {
				goal := "학점은행제로 영유아교수방법론을 통해 수업 설계를 배우고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영유아교수방법론", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:teaching_instruction:credit_bank_instructional_design",
		},
		{
			name: "credit bank assessment design",
			req: func() CreateCourseDraftRequest {
				goal := "외국어로서의한국어능력평가론으로 평가 문항을 설계하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "외국어로서의한국어능력평가론", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:teaching_instruction:credit_bank_instructional_design",
		},
		{
			name: "credit bank standard theory",
			req: func() CreateCourseDraftRequest {
				goal := "학점은행제로 인간행동과사회환경을 배우고 사회복지 기초를 다지고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "인간행동과사회환경", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:professional_transition:credit_bank_standard_theory",
		},
		{
			name: "mental health theory",
			req: func() CreateCourseDraftRequest {
				goal := "정신건강론을 배우고 실무 관점으로 핵심 이론을 정리하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "정신건강론", LearningGoal: &goal}
			}(),
			want: "knowledge_hobby:professional_transition:credit_bank_standard_theory",
		},
		{
			name: "english score exam certification",
			req: func() CreateCourseDraftRequest {
				goal := "TOEIC 850점을 목표로 영어를 공부하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영어 배우기", LearningGoal: &goal}
			}(),
			want: "language_communication:certification_assessment:english_score_exam_certification",
		},
		{
			name: "english speaking interview certification",
			req: func() CreateCourseDraftRequest {
				goal := "OPIc IM2를 목표로 영어 말하기를 준비하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영어 말하기 배우기", LearningGoal: &goal}
			}(),
			want: "language_communication:certification_assessment:english_speaking_interview_certification",
		},
		{
			name: "english integrated four skills certification",
			req: func() CreateCourseDraftRequest {
				goal := "IELTS 6.5를 목표로 영어를 준비하고 싶다"
				return CreateCourseDraftRequest{SourceQuery: "영어 배우기", LearningGoal: &goal}
			}(),
			want: "language_communication:certification_assessment:english_integrated_four_skills_certification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferGoalSubpatternKey(tt.req); got != tt.want {
				t.Fatalf("inferGoalSubpatternKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInferGoalSubpatternKey_UsesGoalProfileContext(t *testing.T) {
	goal := "영어로 자연스럽게 말하고 싶다"
	usage := "해외여행에서 식당 주문과 길 묻기를 하고 싶다"
	if got := inferGoalSubpatternKey(CreateCourseDraftRequest{
		SourceQuery:      "영어 회화 배우기",
		LearningGoal:     &goal,
		GoalUsageContext: &usage,
	}); got != "language_communication:performance_execution:travel_conversation" {
		t.Fatalf("inferGoalSubpatternKey() with usage context = %q", got)
	}
}

func TestKoreanLifelongPracticalLocalRouting(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantPattern string
		wantSub     string
	}{
		{
			name:        "public lifelong online course",
			sourceQuery: "GSEEK 온라인학습 취미 건강 생활상식 디지털역량 평생학습",
			goal:        "온라인 강좌를 골라 생활에 적용하는 학습 루틴을 만들고 싶다",
			wantPattern: "knowledge_hobby:habit_lifestyle:foundation_build",
			wantSub:     "knowledge_hobby:habit_lifestyle:public_lifelong_online_course",
		},
		{
			name:        "blended lifelong participation",
			sourceQuery: "경기도 평생학습포털 GSEEK 오프라인 화상학습 지역 강좌",
			goal:        "온라인 선행학습 후 지역 오프라인 강좌에 참여할 준비를 하고 싶다",
			wantPattern: "knowledge_hobby:participation_service:habit_lifestyle",
			wantSub:     "knowledge_hobby:participation_service:blended_lifelong_participation",
		},
		{
			name:        "midlife transition learning support",
			sourceQuery: "서울런4050 직업교육경비 중장년 직업역량 취창업 학습경비",
			goal:        "중장년 직업 전환을 위해 직무역량 학습계획을 세우고 싶다",
			wantPattern: "knowledge_hobby:professional_transition:foundation_build",
			wantSub:     "knowledge_hobby:professional_transition:midlife_transition_learning_support",
		},
		{
			name:        "completion refund online course",
			sourceQuery: "서울런4050 중장년 특화 온라인과정 진도율 환급 수료",
			goal:        "온라인 과정을 끝까지 듣고 수료 직전까지 핵심 내용을 정리하고 싶다",
			wantPattern: "knowledge_hobby:professional_transition:habit_lifestyle",
			wantSub:     "knowledge_hobby:professional_transition:completion_refund_online_course",
		},
		{
			name:        "federated lifelong open resource",
			sourceQuery: "국가평생학습포털 늘배움 공개 평생교육 강좌 묶음 학습경로",
			goal:        "공개 평생학습 자료를 모아 자기주도 학습계획을 만들고 싶다",
			wantPattern: "knowledge_hobby:foundation_build:habit_lifestyle",
			wantSub:     "knowledge_hobby:foundation_build:federated_lifelong_open_resource",
		},
		{
			name:        "korean content creation pipeline",
			sourceQuery: "한국콘텐츠아카데미 방송영상 게임 만화 애니메이션 콘텐츠 제작 포트폴리오",
			goal:        "콘텐츠 장르를 골라 첫 결과물을 만들고 포트폴리오 직전까지 정리하고 싶다",
			wantPattern: "digital_creation:artifact_creation:foundation_build",
			wantSub:     "digital_creation:artifact_creation:korean_content_creation_pipeline",
		},
		{
			name:        "rural return agriculture foundation",
			sourceQuery: "농업교육포털 귀농귀촌 온라인 교육 농업 품목기술 정착 수료",
			goal:        "귀농귀촌을 준비하며 기초 품목기술과 정착 계획을 정리하고 싶다",
			wantPattern: "knowledge_hobby:professional_transition:habit_lifestyle",
			wantSub:     "knowledge_hobby:professional_transition:rural_return_agriculture_foundation",
		},
		{
			name:        "urban agriculture practice series",
			sourceQuery: "인천광역시농업기술센터 도시농업교육 텃밭 가정원예 재배 실습",
			goal:        "도시농업으로 작은 텃밭과 화분 관리 루틴을 만들고 싶다",
			wantPattern: "knowledge_hobby:habit_lifestyle:artifact_creation",
			wantSub:     "knowledge_hobby:habit_lifestyle:urban_agriculture_practice_series",
		},
		{
			name:        "ai digital literacy life practice",
			sourceQuery: "AI·디지털배움터 생활밀착형 스마트폰 키오스크 보이스피싱 디지털 문해",
			goal:        "스마트폰과 키오스크, AI 활용을 안전하게 생활에서 쓰고 싶다",
			wantPattern: "knowledge_hobby:habit_lifestyle:digital_literacy",
			wantSub:     "knowledge_hobby:habit_lifestyle:ai_digital_literacy_life_practice",
		},
		{
			name:        "local hobby workshop series",
			sourceQuery: "계룡시 평생학습포털 생활취미 강좌 가죽공예 입문 지역 평생학습센터",
			goal:        "지역 생활취미 워크숍에서 작은 작품을 완성하고 싶다",
			wantPattern: "craft_making:artifact_creation:participation_service",
			wantSub:     "craft_making:artifact_creation:local_hobby_workshop_series",
		},
		{
			name:        "korean open lecture weekly survey",
			sourceQuery: "K-MOOC 주차별 강의계획서 공개강좌 인문학 학습목표 중간고사 기말고사",
			goal:        "K-MOOC 공개강좌를 따라가며 핵심 개념을 통합 정리하고 싶다",
			wantPattern: "knowledge_hobby:foundation_build:habit_lifestyle",
			wantSub:     "knowledge_hobby:foundation_build:korean_open_lecture_weekly_survey",
		},
		{
			name:        "korean design tool practical certification",
			sourceQuery: "GSEEK GTQ포토샵2급 기초부터 자격증까지 강의계획서",
			goal:        "GTQ 포토샵 2급 실기시험을 제한시간 안에 풀 수 있게 준비하고 싶다",
			wantPattern: "digital_creation:certification_assessment:artifact_creation",
			wantSub:     "digital_creation:certification_assessment:korean_design_tool_practical_certification",
		},
		{
			name:        "korean ncs web publisher training",
			sourceQuery: "Work24 HRD-Net NCS 디지털디자인 웹퍼블리셔 프론트엔드 반응형 웹 직업훈련",
			goal:        "웹퍼블리셔 취업을 목표로 포트폴리오를 만들고 싶다",
			wantPattern: "digital_creation:professional_transition:artifact_creation",
			wantSub:     "digital_creation:professional_transition:korean_ncs_web_publisher_training",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{SourceQuery: tt.sourceQuery, LearningGoal: &tt.goal}
			if got := formatGoalPatternKeyForLog(inferGoalPatternKey(req)); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if got := inferGoalSubpatternKey(req); got != tt.wantSub {
				t.Fatalf("subpattern = %q, want %q", got, tt.wantSub)
			}
			if summary := analyzeGoalRouting(req); summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesKoreanLifelongPracticalGuidance(t *testing.T) {
	goal := "스마트폰과 키오스크, AI 활용을 안전하게 생활에서 쓰고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "AI·디지털배움터 생활밀착형 스마트폰 키오스크 보이스피싱 디지털 문해",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: AI·디지털 생활문해 실천형") {
		t.Fatalf("prompt missing Korean digital literacy subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "디지털 생활 과제를 안전하게 반복") {
		t.Fatalf("prompt missing Korean digital literacy last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesKoreanLectureCertificateNCSGuidance(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantName    string
		wantRule    string
	}{
		{
			name:        "open lecture weekly survey",
			sourceQuery: "KOCW 공개 강의계획서 주차별 강의계획서",
			goal:        "대학 공개강의를 따라 핵심 개념과 적용 흐름을 정리하고 싶다",
			wantName:    "세부 목표형: 한국 공개강의 주차형",
			wantRule:    "15주·16주 주차표를 그대로 복제하지 말고",
		},
		{
			name:        "design tool practical certification",
			sourceQuery: "GTQ포토샵2급 기초부터 자격증까지",
			goal:        "GTQ 포토샵 2급 실기시험을 모의고사까지 준비하고 싶다",
			wantName:    "세부 목표형: 한국 디자인 도구 실기 자격형",
			wantRule:    "제한시간 안에 수행하는 모의고사",
		},
		{
			name:        "ncs web publisher training",
			sourceQuery: "Work24 HRD-Net NCS 디지털디자인 웹퍼블리셔 프론트엔드 직업훈련",
			goal:        "웹퍼블리셔 취업용 포트폴리오를 완성하고 싶다",
			wantName:    "세부 목표형: 한국 NCS 웹퍼블리셔 직업훈련형",
			wantRule:    "포트폴리오와 지원 자료를 점검",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
				SourceQuery:  tt.sourceQuery,
				LearningGoal: &tt.goal,
			})

			if !strings.Contains(prompt, tt.wantName) {
				t.Fatalf("prompt missing subpattern name %q: %q", tt.wantName, prompt)
			}
			if !strings.Contains(prompt, tt.wantRule) {
				t.Fatalf("prompt missing subpattern rule %q: %q", tt.wantRule, prompt)
			}
		})
	}
}

func TestAnalyzeGoalRouting(t *testing.T) {
	t.Run("subpattern strong", func(t *testing.T) {
		goal := "찬양단 봉사에서 자연스럽게 반주하고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "통기타 배우기",
			LearningGoal: &goal,
		})

		if summary.SubpatternKey != "instrument_performance:participation_service:worship_team_support" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q, want subpattern_strong", summary.RefinementMode)
		}
		if got := formatGoalPatternKeyForLog(summary.PatternKey); got != "instrument_performance:participation_service:performance_execution" {
			t.Fatalf("formatted pattern key = %q", got)
		}
	})

	t.Run("pattern soft", func(t *testing.T) {
		goal := "기타 연주를 자연스럽게 이어가고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "기타 배우기",
			LearningGoal: &goal,
		})

		if summary.SubpatternKey != "" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "pattern_soft" {
			t.Fatalf("refinement mode = %q, want pattern_soft", summary.RefinementMode)
		}
	})

	t.Run("credit bank practicum strong", func(t *testing.T) {
		goal := "학점은행제로 평생교육실습을 이수하고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "평생교육실습",
			LearningGoal: &goal,
		})

		if summary.SubpatternKey != "knowledge_hobby:professional_transition:credit_bank_practicum" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q, want subpattern_strong", summary.RefinementMode)
		}
		if got := formatGoalPatternKeyForLog(summary.PatternKey); got != "knowledge_hobby:professional_transition:participation_service" {
			t.Fatalf("formatted pattern key = %q", got)
		}
	})

	t.Run("credit bank instructional design strong", func(t *testing.T) {
		goal := "학점은행제로 외국어로서의한국어능력평가론을 배우고 평가 문항을 설계하고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "외국어로서의한국어능력평가론",
			LearningGoal: &goal,
		})

		if summary.SubpatternKey != "knowledge_hobby:teaching_instruction:credit_bank_instructional_design" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q, want subpattern_strong", summary.RefinementMode)
		}
		if got := formatGoalPatternKeyForLog(summary.PatternKey); got != "knowledge_hobby:teaching_instruction:professional_transition" {
			t.Fatalf("formatted pattern key = %q", got)
		}
	})

	t.Run("credit bank standard theory strong", func(t *testing.T) {
		goal := "학점은행제로 정신건강론을 배우고 실무 관점으로 핵심 이론을 정리하고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "정신건강론",
			LearningGoal: &goal,
		})

		if summary.SubpatternKey != "knowledge_hobby:professional_transition:credit_bank_standard_theory" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q, want subpattern_strong", summary.RefinementMode)
		}
		if got := formatGoalPatternKeyForLog(summary.PatternKey); got != "knowledge_hobby:professional_transition:foundation_build" {
			t.Fatalf("formatted pattern key = %q", got)
		}
	})

	t.Run("unclassified", func(t *testing.T) {
		goal := "새로운 취미를 찾고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "취미 배우기",
			LearningGoal: &goal,
		})

		if summary.RefinementMode != "unclassified" {
			t.Fatalf("refinement mode = %q, want unclassified", summary.RefinementMode)
		}
		if got := formatGoalPatternKeyForLog(summary.PatternKey); got != "unclassified" {
			t.Fatalf("formatted pattern key = %q, want unclassified", got)
		}
	})
}

func TestBuildCurriculumGenerationPromptIncludesPatternMatchGuidance(t *testing.T) {
	goal := "색소폰으로 재즈 스탠더드 한 곡을 연주한다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:          "색소폰 재즈 곡 연주",
		LearningGoal:         &goal,
		PatternMatchGuidance: "1. 한 곡 완주형\n- 권장 흐름: setup_intro -> first_output -> performance_prep",
	})

	for _, want := range []string{
		"데이터 기반 패턴 검색 가이드",
		"한 곡 완주형",
		"setup_intro -> first_output -> performance_prep",
		"사용자 확정 목표와 충돌하면 확정 목표를 우선",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q", want)
		}
	}
}

func TestApplyDomainSpecificLessonRefinement_CookingHabitLifestyle(t *testing.T) {
	goal := "가족에게 건강한 식사를 해주고 싶다"
	doc := applyDomainSpecificLessonRefinement(CreateCourseDraftRequest{
		SourceQuery:  "요리 배우기",
		LearningGoal: &goal,
	}, generatedDraftDocument{
		MainLessons: []generatedMainLesson{
			{Title: "재료 알아보기", Objective: "재료를 이해한다."},
			{Title: "식사 만들기", Objective: "식사를 만든다."},
		},
	})

	if len(doc.MainLessons) != 3 {
		t.Fatalf("lessons len = %d, want 3", len(doc.MainLessons))
	}
	if got := doc.MainLessons[0].Title; got != "기본 재료와 조리 흐름으로 한 끼를 완성한다" {
		t.Fatalf("first lesson title = %q", got)
	}
	if got := doc.MainLessons[2].Title; got != "가족과 일상에 맞는 식사 루틴을 정리한다" {
		t.Fatalf("last lesson title = %q", got)
	}
	found := false
	for _, item := range doc.CompletionCriteria {
		if item == "건강한 식사 준비를 반복 가능한 생활 루틴으로 이어갈 수 있다." {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("completion criteria missing cooking habit criterion: %#v", doc.CompletionCriteria)
	}
}

func TestBuildCurriculumGenerationPromptIncludesPatternGuidance(t *testing.T) {
	goal := "브런치에 글을 연재하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "에세이 글쓰기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "권장 단계 흐름") {
		t.Fatalf("prompt missing pattern guidance: %q", prompt)
	}
	if !strings.Contains(prompt, "공개와 게시 준비") {
		t.Fatalf("prompt missing publish-prep guidance: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesSubpatternGuidance(t *testing.T) {
	goal := "좋아하는 곡 한 곡을 끝까지 완주하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "피아노 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "자주 등장하는 실제 목표형 세부 가이드") {
		t.Fatalf("prompt missing subpattern section: %q", prompt)
	}
	if !strings.Contains(prompt, "세부 목표형: 한 곡 완주형") {
		t.Fatalf("prompt missing song completion subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "목표 곡을 끝까지 이어갈 직전 준비 단계") {
		t.Fatalf("prompt missing song completion last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumReviewPromptIncludesSubpatternGuidance(t *testing.T) {
	goal := "매일 짧게 기록하는 글쓰기 습관을 만들고 싶다"
	prompt, err := buildCurriculumReviewPrompt(CreateCourseDraftRequest{
		SourceQuery:  "글쓰기 배우기",
		LearningGoal: &goal,
	}, generatedDraftDocument{
		Title: "글쓰기 코스",
		MainLessons: []generatedMainLesson{
			{Title: "짧은 글 쓰기", Objective: "짧은 글을 쓴다."},
		},
		CompletionCriteria: []string{"짧은 기록을 남길 수 있다."},
	})
	if err != nil {
		t.Fatalf("buildCurriculumReviewPrompt() error = %v", err)
	}

	if !strings.Contains(prompt, "실사용 목표형 세부 검토 기준") {
		t.Fatalf("review prompt missing subpattern section: %q", prompt)
	}
	if !strings.Contains(prompt, "세부 목표형: 매일 글쓰기 루틴형") {
		t.Fatalf("review prompt missing daily writing subpattern name: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCraftTeachYoutubeSubpatternGuidance(t *testing.T) {
	goal := "나만의 작품을 만들고 유튜브를 통해 강의하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "코바늘 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 공예 유튜브 강의형") {
		t.Fatalf("prompt missing craft-teach-youtube subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "제작 과정을 설명하고 시연할 준비") {
		t.Fatalf("prompt missing craft-teach-youtube last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesTravelConversationSubpatternGuidance(t *testing.T) {
	goal := "여행 가서 영어로 자연스럽게 말하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "영어 회화 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 여행 회화형") {
		t.Fatalf("prompt missing travel-conversation subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "여행 상황에서 실제 대화를 시작하기 직전") {
		t.Fatalf("prompt missing travel-conversation last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCertificationSubpatternGuidance(t *testing.T) {
	goal := "JLPT N4 취득을 목표로 일본어를 공부한다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "일본어 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 언어 시험 준비형") {
		t.Fatalf("prompt missing certification subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "시험 직전까지 영역별 문제 흐름을 정리") {
		t.Fatalf("prompt missing certification last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishCertificationSubpatternGuidance(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		goal      string
		containsA string
		containsB string
	}{
		{
			name:      "toeic",
			source:    "영어 배우기",
			goal:      "TOEIC 850점을 목표로 영어를 공부하고 싶다",
			containsA: "세부 목표형: 영어 점수형 시험 준비형",
			containsB: "파트별 문제풀이와 시간 운영 흐름",
		},
		{
			name:      "opic",
			source:    "영어 말하기 배우기",
			goal:      "OPIc IM2를 목표로 영어 말하기를 준비하고 싶다",
			containsA: "세부 목표형: 영어 말하기 인터뷰 시험형",
			containsB: "인터뷰 흐름과 답변 길이 조절",
		},
		{
			name:      "ielts",
			source:    "영어 배우기",
			goal:      "IELTS 6.5를 목표로 영어를 준비하고 싶다",
			containsA: "세부 목표형: 영어 4영역 통합 시험형",
			containsB: "4영역 task 전환과 시간 배분 흐름",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
				SourceQuery:  tt.source,
				LearningGoal: &tt.goal,
			})
			if !strings.Contains(prompt, tt.containsA) {
				t.Fatalf("prompt missing %q: %q", tt.containsA, prompt)
			}
			if !strings.Contains(prompt, tt.containsB) {
				t.Fatalf("prompt missing %q: %q", tt.containsB, prompt)
			}
		})
	}
}

func TestBuildRecommendationGoalMetadata_EnglishCertifications(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		goal           string
		wantSubpattern string
		wantTerm       string
	}{
		{
			name:           "toeic",
			source:         "영어 배우기",
			goal:           "TOEIC 850점을 목표로 영어를 공부하고 싶다",
			wantSubpattern: "language_communication:certification_assessment:english_score_exam_certification",
			wantTerm:       "toeic",
		},
		{
			name:           "opic",
			source:         "영어 말하기 배우기",
			goal:           "OPIc IM2를 목표로 영어 말하기를 준비하고 싶다",
			wantSubpattern: "language_communication:certification_assessment:english_speaking_interview_certification",
			wantTerm:       "인터뷰",
		},
		{
			name:           "ielts",
			source:         "영어 배우기",
			goal:           "IELTS 6.5를 목표로 영어를 준비하고 싶다",
			wantSubpattern: "language_communication:certification_assessment:english_integrated_four_skills_certification",
			wantTerm:       "writing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
				SourceQuery:  tt.source,
				LearningGoal: &tt.goal,
			})
			if meta.SubpatternKey != tt.wantSubpattern {
				t.Fatalf("subpattern = %q, want %q", meta.SubpatternKey, tt.wantSubpattern)
			}
			if !strings.Contains(strings.ToLower(meta.QueryHint), strings.ToLower(tt.wantTerm)) {
				t.Fatalf("query hint = %q, want term %q", meta.QueryHint, tt.wantTerm)
			}
		})
	}
}

func TestEnglishCreativeLifestyleRouting(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		goal           string
		language       string
		wantPattern    string
		wantSubpattern string
	}{
		{
			name:           "adobe creative learning path",
			source:         "Adobe Learn Photoshop tutorial learning path",
			goal:           "Create a portfolio-ready design output with Adobe creative tools",
			wantPattern:    "digital_creation:artifact_creation:foundation_build",
			wantSubpattern: "digital_creation:artifact_creation:adobe_creative_learning_path",
		},
		{
			name:           "template design quick output",
			source:         "Canva Design School",
			goal:           "Create a quick social media design from templates",
			wantPattern:    "digital_creation:artifact_creation:foundation_build",
			wantSubpattern: "digital_creation:artifact_creation:template_design_quick_output",
		},
		{
			name:           "sequential creative class path",
			source:         "Skillshare Learning Paths illustration",
			goal:           "Build a sequence of creative class projects for a small portfolio",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:sequential_creative_class_path",
		},
		{
			name:           "animation fundamentals shot progression",
			source:         "Blender Studio Animation Fundamentals",
			goal:           "Make a short character animation shot with bouncing ball and body mechanics",
			wantPattern:    "digital_creation:artifact_creation:foundation_build",
			wantSubpattern: "digital_creation:artifact_creation:animation_fundamentals_shot_progression",
		},
		{
			name:           "interactive music production foundation",
			source:         "Ableton Learning Music",
			goal:           "Create a short beat loop and song section with browser music production",
			wantPattern:    "digital_creation:artifact_creation:foundation_build",
			wantSubpattern: "digital_creation:artifact_creation:interactive_music_production_foundation",
		},
		{
			name:           "song driven instrument path",
			source:         "Fender Play path",
			goal:           "Practice guitar riffs and chords to play a simple song",
			wantPattern:    "instrument_performance:performance_execution",
			wantSubpattern: "instrument_performance:performance_execution:song_driven_instrument_path",
		},
		{
			name:           "maker project class",
			source:         "Instructables Classes electronics project",
			goal:           "Build a small DIY maker project and explain the steps",
			wantPattern:    "maker_technical_hobby:artifact_creation:foundation_build",
			wantSubpattern: "maker_technical_hobby:artifact_creation:maker_project_class",
		},
		{
			name:           "short diploma assessment path",
			source:         "Alison Diploma Courses",
			goal:           "Prepare career skills through self-paced modules and assessment",
			wantPattern:    "knowledge_hobby:professional_transition:artifact_creation",
			wantSubpattern: "knowledge_hobby:professional_transition:short_diploma_assessment_path",
		},
		{
			name:           "short course weekly discussion",
			source:         "FutureLearn Short Courses",
			goal:           "Study a short course with weekly activities and discussion",
			wantPattern:    "knowledge_hobby:foundation_build:artifact_creation",
			wantSubpattern: "knowledge_hobby:foundation_build:short_course_weekly_discussion",
		},
		{
			name:           "audit course practice path",
			source:         "edX Free Audit Courses",
			goal:           "Use lectures readings and ungraded practice in an audit course",
			wantPattern:    "knowledge_hobby:foundation_build:artifact_creation",
			wantSubpattern: "knowledge_hobby:foundation_build:audit_course_practice_path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{SourceQuery: tt.source, LearningGoal: &tt.goal, LearningLanguage: tt.language}
			summary := analyzeGoalRouting(req)
			if got := formatGoalPatternKeyForLog(summary.PatternKey); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if summary.SubpatternKey != tt.wantSubpattern {
				t.Fatalf("subpattern = %q, want %q", summary.SubpatternKey, tt.wantSubpattern)
			}
			if summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishCreativeLifestyleGuidance(t *testing.T) {
	goal := "Create a quick social media design from templates"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "Canva Design School",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recurring goal subtype: template design quick output") {
		t.Fatalf("prompt missing English subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "Recommended flow for this subtype") {
		t.Fatalf("prompt missing English subtype flow: %q", prompt)
	}
	if strings.Contains(prompt, "템플릿 기반") {
		t.Fatalf("prompt leaked Korean subtype guidance: %q", prompt)
	}
}

func TestBuildRecommendationGoalMetadata_EnglishCreativeLifestyle(t *testing.T) {
	goal := "Create a short beat loop and song section with browser music production"
	meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
		SourceQuery:  "Ableton Learning Music",
		LearningGoal: &goal,
	})

	if meta.SubpatternKey != "digital_creation:artifact_creation:interactive_music_production_foundation" {
		t.Fatalf("subpattern = %q", meta.SubpatternKey)
	}
	if !strings.Contains(strings.ToLower(meta.QueryHint), "ableton") {
		t.Fatalf("query hint = %q, want ableton", meta.QueryHint)
	}
}

func TestEnglishCraftMakerLifeSkillRouting(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		goal           string
		language       string
		wantPattern    string
		wantSubpattern string
	}{
		{
			name:           "sewing project foundation",
			source:         "Alison Basics of Sewing",
			goal:           "Learn sewing machine basics and finish a small sewing project",
			wantPattern:    "craft_making:artifact_creation:foundation_build",
			wantSubpattern: "craft_making:artifact_creation:sewing_project_foundation",
		},
		{
			name:           "cardboard circuit invention path",
			source:         "Cardboard Circuits open online course",
			goal:           "Build a cardboard invention with electronics robotics and programming",
			wantPattern:    "maker_technical_hobby:artifact_creation:foundation_build",
			wantSubpattern: "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path",
		},
		{
			name:           "music production full pipeline",
			source:         "MIXXIN Academy Music Production Certificate",
			goal:           "Create a full music production track from sound design to mixing and mastering",
			wantPattern:    "digital_creation:artifact_creation:foundation_build",
			wantSubpattern: "digital_creation:artifact_creation:music_production_full_pipeline",
		},
		{
			name:           "photography fundamentals project",
			source:         "Icon Photography School Photography Fundamentals",
			goal:           "Learn exposure composition and lighting to complete a small photo project",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:photography_fundamentals_project",
		},
		{
			name:           "smartphone photography daily project",
			source:         "스마트폰 사진 기초",
			goal:           "스마트폰으로 일상 사진을 잘 찍고 구도와 빛을 활용해 작은 사진 프로젝트를 완성한다",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:photography_fundamentals_project",
		},
		{
			name:           "portrait photography natural light",
			source:         "인물 사진 촬영",
			goal:           "자연광과 배경 정리를 활용해 친구 인물 사진을 안정적으로 촬영한다",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:photography_fundamentals_project",
		},
		{
			name:           "english smartphone photography",
			source:         "smartphone photography",
			goal:           "I want to take better smartphone photos with composition and natural light",
			language:       "en",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:photography_fundamentals_project",
		},
		{
			name:           "garden design maintenance path",
			source:         "Garden Tutor Free Online Gardening Course",
			goal:           "Design a small garden and maintain it with a weekly routine",
			wantPattern:    "knowledge_hobby:habit_lifestyle:foundation_build",
			wantSubpattern: "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path",
		},
		{
			name:           "cooking basics foundation",
			source:         "Alison Cooking Basics Food Hygiene Kitchen Utensils and Seasonings",
			goal:           "Use knife skills sauces herbs and spices to make a simple dish safely",
			wantPattern:    "cooking_baking:artifact_creation:foundation_build",
			wantSubpattern: "cooking_baking:artifact_creation:cooking_basics_foundation",
		},
		{
			name:           "critical reading to writing draft",
			source:         "OpenLearn Creative Writing and Critical Reading",
			goal:           "Use critical reading to revise a fiction poetry or scriptwriting draft",
			wantPattern:    "writing_storytelling:presentation_publish:artifact_creation",
			wantSubpattern: "writing_storytelling:presentation_publish:critical_reading_to_writing_draft",
		},
		{
			name:           "photography studio project",
			source:         "MIT OCW Introduction to Photography and Related Media",
			goal:           "Create a photography studio term project with digital imaging and studio lighting",
			wantPattern:    "visual_art:artifact_creation:artifact_creation",
			wantSubpattern: "visual_art:artifact_creation:photography_studio_project",
		},
		{
			name:           "design lab prototype project",
			source:         "MIT OCW D-Lab II: Design",
			goal:           "Develop a real-world prototype through design process fabrication and review",
			wantPattern:    "maker_technical_hobby:artifact_creation:foundation_build",
			wantSubpattern: "maker_technical_hobby:artifact_creation:design_lab_prototype_project",
		},
		{
			name:           "organic growing cycle plan",
			source:         "Alison Growing Organic Food Sustainably",
			goal:           "Plan crop seasons sowing and care routine for a first organic growing cycle",
			wantPattern:    "knowledge_hobby:habit_lifestyle:foundation_build",
			wantSubpattern: "knowledge_hobby:habit_lifestyle:organic_growing_cycle_plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{SourceQuery: tt.source, LearningGoal: &tt.goal, LearningLanguage: tt.language}
			summary := analyzeGoalRouting(req)
			if got := formatGoalPatternKeyForLog(summary.PatternKey); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if summary.SubpatternKey != tt.wantSubpattern {
				t.Fatalf("subpattern = %q, want %q", summary.SubpatternKey, tt.wantSubpattern)
			}
			if summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishCraftMakerLifeSkillGuidance(t *testing.T) {
	goal := "Learn sewing machine basics and finish a small sewing project"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "Alison Basics of Sewing",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recurring goal subtype: sewing project foundation") {
		t.Fatalf("prompt missing English subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "Recommended flow for this subtype") {
		t.Fatalf("prompt missing English subtype flow: %q", prompt)
	}
	if strings.Contains(prompt, "재봉틀") {
		t.Fatalf("prompt leaked Korean subtype guidance: %q", prompt)
	}
}

func TestBuildRecommendationGoalMetadata_EnglishCraftMakerLifeSkill(t *testing.T) {
	goal := "Develop a real-world prototype through design process fabrication and review"
	meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
		SourceQuery:  "MIT OCW D-Lab II: Design",
		LearningGoal: &goal,
	})

	if meta.SubpatternKey != "maker_technical_hobby:artifact_creation:design_lab_prototype_project" {
		t.Fatalf("subpattern = %q", meta.SubpatternKey)
	}
	if !strings.Contains(strings.ToLower(meta.QueryHint), "prototype") {
		t.Fatalf("query hint = %q, want prototype", meta.QueryHint)
	}
}

func TestEnglishScienceDesignCitizenHumanitiesRouting(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		goal           string
		wantPattern    string
		wantSubpattern string
	}{
		{
			name:           "science simulation inquiry path",
			source:         "PhET Interactive Simulations",
			goal:           "Use an interactive science simulation to predict variables and explain cause-effect relationships",
			wantPattern:    "knowledge_hobby:concept_mastery:foundation_build",
			wantSubpattern: "knowledge_hobby:concept_mastery:science_simulation_inquiry_path",
		},
		{
			name:           "guided simulation activity sheet",
			source:         "PhET Teaching with PhET Small Group Activities",
			goal:           "Use a guided simulation activity sheet for observation discussion and reflection",
			wantPattern:    "knowledge_hobby:concept_mastery:foundation_build",
			wantSubpattern: "knowledge_hobby:concept_mastery:guided_simulation_activity_sheet",
		},
		{
			name:           "food preservation safety path",
			source:         "Penn State Extension Home Food Preservation Water Bath Canning",
			goal:           "Learn water bath canning safety for high acid foods using research-tested recipes",
			wantPattern:    "cooking_baking:habit_lifestyle:foundation_build",
			wantSubpattern: "cooking_baking:habit_lifestyle:food_preservation_safety_path",
		},
		{
			name:           "horticulture diagnostic foundation",
			source:         "University of Minnesota Extension ProHort",
			goal:           "Build horticulture diagnostics from botany soil integrated pest management and plant pathology",
			wantPattern:    "knowledge_hobby:concept_mastery:foundation_build",
			wantSubpattern: "knowledge_hobby:concept_mastery:horticulture_diagnostic_foundation",
		},
		{
			name:           "computational music analysis project",
			source:         "MIT OCW Computational Music Theory and Analysis",
			goal:           "Use Python music21 and symbolic score representation for a small corpus analysis project",
			wantPattern:    "digital_creation:technical_skill:foundation_build",
			wantSubpattern: "digital_creation:technical_skill:computational_music_analysis_project",
		},
		{
			name:           "design thinking problem framing path",
			source:         "Stanford d.school Design Thinking Bootleg",
			goal:           "Use empathize define ideate prototype test to frame a human-centered design problem",
			wantPattern:    "maker_technical_hobby:artifact_creation:foundation_build",
			wantSubpattern: "maker_technical_hobby:artifact_creation:design_thinking_problem_framing_path",
		},
		{
			name:           "no-build prototype test",
			source:         "Stanford d.school The No-Build Hack",
			goal:           "Create a no-build prototype test with scrappy materials and user feedback",
			wantPattern:    "maker_technical_hobby:artifact_creation:foundation_build",
			wantSubpattern: "maker_technical_hobby:artifact_creation:no_build_prototype_test",
		},
		{
			name:           "climate data graphing explanation",
			source:         "NASA Science Graphing Global Temperature Trends Activity",
			goal:           "Graph temperature anomaly data and compare short-term and long-term global temperature trends",
			wantPattern:    "knowledge_hobby:concept_mastery:foundation_build",
			wantSubpattern: "knowledge_hobby:concept_mastery:climate_data_graphing_explanation",
		},
		{
			name:           "citizen science observation record",
			source:         "National Park Service iNaturalist Training Lesson Plan",
			goal:           "Document organisms in a BioBlitz and build an iNaturalist biodiversity record",
			wantPattern:    "knowledge_hobby:habit_lifestyle:foundation_build",
			wantSubpattern: "knowledge_hobby:habit_lifestyle:citizen_science_observation_record",
		},
		{
			name:           "primary source inquiry note",
			source:         "Library of Congress Primary Source Analysis Tool",
			goal:           "Use observe reflect question to write a primary source inquiry note",
			wantPattern:    "knowledge_hobby:concept_mastery:foundation_build",
			wantSubpattern: "knowledge_hobby:concept_mastery:primary_source_inquiry_note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{SourceQuery: tt.source, LearningGoal: &tt.goal}
			summary := analyzeGoalRouting(req)
			if got := formatGoalPatternKeyForLog(summary.PatternKey); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if summary.SubpatternKey != tt.wantSubpattern {
				t.Fatalf("subpattern = %q, want %q", summary.SubpatternKey, tt.wantSubpattern)
			}
			if summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishScienceDesignCitizenHumanitiesGuidance(t *testing.T) {
	goal := "Use an interactive science simulation to predict variables and explain cause-effect relationships"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "PhET Interactive Simulations",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recurring goal subtype: science simulation inquiry path") {
		t.Fatalf("prompt missing English subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "prediction prompt") {
		t.Fatalf("prompt missing English subtype role: %q", prompt)
	}
	if strings.Contains(prompt, "시뮬레이션") {
		t.Fatalf("prompt leaked Korean subtype guidance: %q", prompt)
	}
}

func TestBuildRecommendationGoalMetadata_EnglishScienceDesignCitizenHumanities(t *testing.T) {
	goal := "Graph temperature anomaly data and compare short-term and long-term global temperature trends"
	meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
		SourceQuery:  "NASA Science Graphing Global Temperature Trends Activity",
		LearningGoal: &goal,
	})

	if meta.SubpatternKey != "knowledge_hobby:concept_mastery:climate_data_graphing_explanation" {
		t.Fatalf("subpattern = %q", meta.SubpatternKey)
	}
	if !strings.Contains(strings.ToLower(meta.QueryHint), "temperature anomaly") {
		t.Fatalf("query hint = %q, want temperature anomaly", meta.QueryHint)
	}
}

func TestBuildFallbackDraftPlan_CertificationVariants(t *testing.T) {
	t.Run("language certification", func(t *testing.T) {
		goal := "JLPT N4 취득을 목표로 일본어를 공부한다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "일본어 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "시험 범위와 N4 핵심 표현 범주를 정리한다" {
			t.Fatalf("first lesson title = %q", got)
		}
		if got := plan.MainLessons[len(plan.MainLessons)-1].Title; got != "시험 직전 약점을 보완하고 실전 감각을 맞춘다" {
			t.Fatalf("last lesson title = %q", got)
		}
	})

	t.Run("cloud certification", func(t *testing.T) {
		goal := "AWS Certified Cloud Practitioner 자격증을 취득하고 싶다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "클라우드 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "시험 구조와 핵심 서비스 범주를 정리한다" {
			t.Fatalf("first lesson title = %q", got)
		}
	})

	t.Run("practical certification", func(t *testing.T) {
		goal := "지게차운전기능사 취득을 목표로 공부하고 싶다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "지게차 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "시험 구조와 기본 안전·위생 기준을 정리한다" {
			t.Fatalf("first lesson title = %q", got)
		}
		if got := plan.MainLessons[len(plan.MainLessons)-1].Title; got != "시험 직전 체크포인트를 점검한다" {
			t.Fatalf("last lesson title = %q", got)
		}
	})

	t.Run("english score certification", func(t *testing.T) {
		goal := "TOEIC 850점을 목표로 영어를 공부하고 싶다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "영어 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "시험 구조와 파트 구성을 정리한다" {
			t.Fatalf("first lesson title = %q", got)
		}
		if got := plan.MainLessons[len(plan.MainLessons)-1].Title; got != "시험 직전 점수형 운영 흐름을 정리한다" {
			t.Fatalf("last lesson title = %q", got)
		}
	})

	t.Run("english speaking interview certification", func(t *testing.T) {
		goal := "OPIc IM2를 목표로 영어 말하기를 준비하고 싶다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "영어 말하기 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "인터뷰 구조와 자기소개 흐름을 익힌다" {
			t.Fatalf("first lesson title = %q", got)
		}
	})

	t.Run("english integrated four skills certification", func(t *testing.T) {
		goal := "IELTS 6.5를 목표로 영어를 준비하고 싶다"
		plan := buildFallbackDraftPlan(CreateCourseDraftRequest{
			SourceQuery:  "영어 배우기",
			LearningGoal: &goal,
		})
		if got := plan.MainLessons[0].Title; got != "시험 유형과 4영역 구조를 정리한다" {
			t.Fatalf("first lesson title = %q", got)
		}
	})
}

func TestBuildCurriculumGenerationPromptIncludesGoalProfileDetails(t *testing.T) {
	goal := "외국인 동료와 업무 대화를 자연스럽게 이어가고 싶다"
	usage := "주간 회의와 슬랙 텍스트 대화에서 바로 써야 한다"
	motivation := "프로젝트 협업에서 소극적으로 보이지 않고 싶다"
	timeHorizon := "3개월"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "영어 회화 배우기",
		LearningGoal:     &goal,
		GoalUsageContext: &usage,
		GoalMotivation:   &motivation,
		GoalTimeHorizon:  &timeHorizon,
	})

	if !strings.Contains(prompt, "활용 맥락: "+usage) {
		t.Fatalf("prompt missing goal usage context: %q", prompt)
	}
	if !strings.Contains(prompt, "동기: "+motivation) {
		t.Fatalf("prompt missing goal motivation: %q", prompt)
	}
	if !strings.Contains(prompt, "목표 기간: "+timeHorizon) {
		t.Fatalf("prompt missing goal time horizon: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCreatorTutorialSubpatternGuidance(t *testing.T) {
	goal := "영상편집 과정을 튜토리얼로 설명하고 강의하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 디지털 제작 튜토리얼형") {
		t.Fatalf("prompt missing creator-tutorial subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "짧은 강의 영상이나 튜토리얼 콘텐츠를 촬영·편집하고 게시 직전") {
		t.Fatalf("prompt missing creator-tutorial last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesBakingClassDemoSubpatternGuidance(t *testing.T) {
	goal := "베이킹 레시피를 설명하며 클래스를 해보고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "베이킹 배우기",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 베이킹 클래스 시연형") {
		t.Fatalf("prompt missing baking-class subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "레시피와 시연 흐름을 설명할 준비") {
		t.Fatalf("prompt missing baking-class last-lesson rule: %q", prompt)
	}
}

func TestUniversityPatternRouting(t *testing.T) {
	t.Run("university survey", func(t *testing.T) {
		goal := "대학 강의처럼 현대미술 흐름을 넓게 이해하고 비교해서 정리하고 싶다"
		key := inferGoalPatternKey(CreateCourseDraftRequest{
			SourceQuery:  "현대미술사",
			LearningGoal: &goal,
		})
		if key.DomainAxis != "knowledge_hobby" {
			t.Fatalf("domain axis = %q, want knowledge_hobby", key.DomainAxis)
		}
		if key.GoalModePrimary != "habit_lifestyle" {
			t.Fatalf("goal mode primary = %q, want habit_lifestyle", key.GoalModePrimary)
		}
		if key.GoalModeSecondary != "foundation_build" {
			t.Fatalf("goal mode secondary = %q, want foundation_build", key.GoalModeSecondary)
		}
		if got := inferGoalSubpatternKey(CreateCourseDraftRequest{
			SourceQuery:  "현대미술사",
			LearningGoal: &goal,
		}); got != "knowledge_hobby:habit_lifestyle:university_survey" {
			t.Fatalf("subpattern = %q", got)
		}
	})

	t.Run("university studio writing", func(t *testing.T) {
		goal := "희곡 워크숍처럼 초안부터 써 보고 수정해 제출 가능한 형태까지 만들고 싶다"
		key := inferGoalPatternKey(CreateCourseDraftRequest{
			SourceQuery:  "희곡 쓰기",
			LearningGoal: &goal,
		})
		if key.DomainAxis != "writing_storytelling" {
			t.Fatalf("domain axis = %q, want writing_storytelling", key.DomainAxis)
		}
		if key.GoalModePrimary != "presentation_publish" {
			t.Fatalf("goal mode primary = %q, want presentation_publish", key.GoalModePrimary)
		}
		if got := inferGoalSubpatternKey(CreateCourseDraftRequest{
			SourceQuery:  "희곡 쓰기",
			LearningGoal: &goal,
		}); got != "writing_storytelling:presentation_publish:university_studio" {
			t.Fatalf("subpattern = %q", got)
		}
	})

	t.Run("university studio creative tech", func(t *testing.T) {
		goal := "작곡 수업처럼 작은 스케치에서 출발해 제출 가능한 작업물까지 만들고 싶다"
		key := inferGoalPatternKey(CreateCourseDraftRequest{
			SourceQuery:  "입문 작곡",
			LearningGoal: &goal,
		})
		if key.DomainAxis != "digital_creation" {
			t.Fatalf("domain axis = %q, want digital_creation", key.DomainAxis)
		}
		if key.GoalModePrimary != "artifact_creation" {
			t.Fatalf("goal mode primary = %q, want artifact_creation", key.GoalModePrimary)
		}
		if got := inferGoalSubpatternKey(CreateCourseDraftRequest{
			SourceQuery:  "입문 작곡",
			LearningGoal: &goal,
		}); got != "digital_creation:artifact_creation:university_studio" {
			t.Fatalf("subpattern = %q", got)
		}
	})
}

func TestAnalyzeGoalRouting_UniversityPatterns(t *testing.T) {
	t.Run("survey strong", func(t *testing.T) {
		goal := "대학 강의처럼 현대미술 흐름을 넓게 이해하고 비교해서 정리하고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "현대미술사",
			LearningGoal: &goal,
		})
		if summary.SubpatternKey != "knowledge_hobby:habit_lifestyle:university_survey" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q", summary.RefinementMode)
		}
	})

	t.Run("studio strong", func(t *testing.T) {
		goal := "희곡 워크숍처럼 초안부터 써 보고 수정해 제출 가능한 형태까지 만들고 싶다"
		summary := analyzeGoalRouting(CreateCourseDraftRequest{
			SourceQuery:  "희곡 쓰기",
			LearningGoal: &goal,
		})
		if summary.SubpatternKey != "writing_storytelling:presentation_publish:university_studio" {
			t.Fatalf("subpattern key = %q", summary.SubpatternKey)
		}
		if summary.RefinementMode != "subpattern_strong" {
			t.Fatalf("refinement mode = %q", summary.RefinementMode)
		}
	})
}

func TestBuildCurriculumGenerationPromptIncludesUniversitySubpatternGuidance(t *testing.T) {
	t.Run("survey", func(t *testing.T) {
		goal := "대학 강의처럼 현대미술 흐름을 넓게 이해하고 비교해서 정리하고 싶다"
		prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
			SourceQuery:  "현대미술사",
			LearningGoal: &goal,
		})
		if !strings.Contains(prompt, "세부 목표형: 대학 survey형") {
			t.Fatalf("prompt missing university survey subpattern name: %q", prompt)
		}
		if !strings.Contains(prompt, "비교·비평·종합 관점") {
			t.Fatalf("prompt missing university survey rule: %q", prompt)
		}
	})

	t.Run("studio", func(t *testing.T) {
		goal := "희곡 워크숍처럼 초안부터 써 보고 수정해 제출 가능한 형태까지 만들고 싶다"
		prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
			SourceQuery:  "희곡 쓰기",
			LearningGoal: &goal,
		})
		if !strings.Contains(prompt, "세부 목표형: 대학 studio/workshop형") {
			t.Fatalf("prompt missing university studio subpattern name: %q", prompt)
		}
		if !strings.Contains(prompt, "작은 초안이나 첫 결과물이 빠르게 나오도록") {
			t.Fatalf("prompt missing university studio rule: %q", prompt)
		}
	})
}

func TestKMOOCPatternRouting(t *testing.T) {
	goal := "K-MOOC 온라인 공개강좌처럼 AI와 윤리 핵심 쟁점을 배우고 평가 직전까지 정리하고 싶다"
	req := CreateCourseDraftRequest{
		SourceQuery:  "K-MOOC AI와 윤리",
		LearningGoal: &goal,
	}

	key := inferGoalPatternKey(req)
	if key.DomainAxis != "knowledge_hobby" {
		t.Fatalf("domain axis = %q, want knowledge_hobby", key.DomainAxis)
	}
	if key.GoalModePrimary != "habit_lifestyle" {
		t.Fatalf("goal mode primary = %q, want habit_lifestyle", key.GoalModePrimary)
	}
	if key.GoalModeSecondary != "foundation_build" {
		t.Fatalf("goal mode secondary = %q, want foundation_build", key.GoalModeSecondary)
	}

	if got := inferGoalSubpatternKey(req); got != "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey" {
		t.Fatalf("subpattern = %q", got)
	}
}

func TestAnalyzeGoalRouting_KMOOCPattern(t *testing.T) {
	goal := "케이무크 강좌처럼 클래식 음악 감상을 주차별로 따라가고 평가 직전까지 정리하고 싶다"
	summary := analyzeGoalRouting(CreateCourseDraftRequest{
		SourceQuery:  "K-MOOC 클래식 친구 만들기",
		LearningGoal: &goal,
	})
	if summary.SubpatternKey != "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey" {
		t.Fatalf("subpattern key = %q", summary.SubpatternKey)
	}
	if summary.RefinementMode != "subpattern_strong" {
		t.Fatalf("refinement mode = %q, want subpattern_strong", summary.RefinementMode)
	}
}

func TestBuildCurriculumGenerationPromptIncludesKMOOCSubpatternGuidance(t *testing.T) {
	goal := "K-MOOC 공개강좌 방식으로 소셜 네트워크 분석을 배우고 퀴즈와 총괄평가 직전까지 정리하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "K-MOOC 비정형 데이터 분석 소셜 네트워크",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: K-MOOC 온라인 공개강좌형") {
		t.Fatalf("prompt missing K-MOOC subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "K-MOOC 주차를 그대로 복제하지 말고") {
		t.Fatalf("prompt missing K-MOOC compression rule: %q", prompt)
	}
	if !strings.Contains(prompt, "수강 완료, 시험 제출, 이수증 발급은 completion_criteria") {
		t.Fatalf("prompt missing K-MOOC completion rule: %q", prompt)
	}
}

func TestApplyDomainSpecificLessonRefinement_KMOOC(t *testing.T) {
	goal := "온라인 공개강좌 방식으로 3D 모션캡쳐 기술을 배우고 평가 직전까지 정리하고 싶다"
	doc := applyDomainSpecificLessonRefinement(CreateCourseDraftRequest{
		SourceQuery:  "K-MOOC 3D 모션캡쳐",
		LearningGoal: &goal,
	}, generatedDraftDocument{
		MainLessons: []generatedMainLesson{{Title: "임시 리슨", Objective: "임시 목표"}},
	})

	if len(doc.MainLessons) != 4 {
		t.Fatalf("lesson count = %d, want 4", len(doc.MainLessons))
	}
	if doc.MainLessons[0].Title != "공개강좌의 범위와 주차 흐름을 잡는다" {
		t.Fatalf("first lesson = %q", doc.MainLessons[0].Title)
	}
	if doc.MainLessons[3].Title != "평가 직전까지 전체 흐름을 통합 정리한다" {
		t.Fatalf("last lesson = %q", doc.MainLessons[3].Title)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCreditBankPracticumSubpatternGuidance(t *testing.T) {
	goal := "학점은행제로 사회복지현장실습을 이수하고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "사회복지현장실습",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 학점은행제 실습형") {
		t.Fatalf("prompt missing credit-bank practicum subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "기관실습과 사후정리를 마쳐 실습 인정 직전까지 정리") {
		t.Fatalf("prompt missing credit-bank practicum last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCreditBankInstructionalDesignSubpatternGuidance(t *testing.T) {
	goal := "학점은행제로 영유아교수방법론을 통해 수업 설계를 배우고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "영유아교수방법론",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 학점은행제 교수설계형") {
		t.Fatalf("prompt missing credit-bank instructional-design subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "수업안이나 평가 설계를 정리해 실제 운영 직전까지 준비") {
		t.Fatalf("prompt missing credit-bank instructional-design last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptIncludesCreditBankStandardTheorySubpatternGuidance(t *testing.T) {
	goal := "학점은행제로 인간행동과사회환경을 배우고 사회복지 기초를 다지고 싶다"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:  "인간행동과사회환경",
		LearningGoal: &goal,
	})

	if !strings.Contains(prompt, "세부 목표형: 학점은행제 표준 이론형") {
		t.Fatalf("prompt missing credit-bank standard-theory subpattern name: %q", prompt)
	}
	if !strings.Contains(prompt, "주요 개념을 사례와 실무 관점으로 연결해 정리") {
		t.Fatalf("prompt missing credit-bank standard-theory last-lesson rule: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptUsesLearningLanguage(t *testing.T) {
	goal := "Build a simple SaaS MVP with AI coding tools"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "AI coding SaaS MVP",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Learner-facing output language must be English") {
		t.Fatalf("prompt missing English output instruction: %q", prompt)
	}
	if !strings.Contains(prompt, "completion_criteria must be natural English") {
		t.Fatalf("prompt missing English criteria instruction: %q", prompt)
	}
}

func TestBuildCurriculumGenerationPromptUsesEnglishPatternGuidance(t *testing.T) {
	goal := "Build a simple SaaS MVP with AI coding tools"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "AI coding SaaS MVP",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recommended stage flow") {
		t.Fatalf("prompt missing English pattern guidance: %q", prompt)
	}
	if !strings.Contains(prompt, "Recurring goal subtype") {
		t.Fatalf("prompt missing English subpattern guidance: %q", prompt)
	}
	if strings.Contains(prompt, "권장 단계 흐름") || strings.Contains(prompt, "세부 목표형") || strings.Contains(prompt, "독립 lesson 금지") {
		t.Fatalf("English prompt should not include Korean pattern guidance: %q", prompt)
	}
}

func TestEnglishMOOCOERRouting(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantPattern string
		wantSub     string
	}{
		{
			name:        "programming foundation",
			sourceQuery: "MIT OCW 6.0001 Introduction to Computer Science and Programming in Python",
			goal:        "Learn programming through problem sets and build a small Python program",
			wantPattern: "digital_creation:foundation_build:artifact_creation",
			wantSub:     "digital_creation:foundation_build:mooc_programming_foundation",
		},
		{
			name:        "humanities survey",
			sourceQuery: "MIT OCW Introduction to Western Music",
			goal:        "Follow a humanities survey course and prepare a short discussion essay",
			wantPattern: "knowledge_hobby:foundation_build:artifact_creation",
			wantSub:     "knowledge_hobby:foundation_build:mooc_humanities_survey",
		},
		{
			name:        "digital literacy badged course",
			sourceQuery: "OpenLearn Digital skills: succeeding in a digital world",
			goal:        "Improve digital literacy, online safety, and make a personal practice plan",
			wantPattern: "knowledge_hobby:habit_lifestyle:foundation_build",
			wantSub:     "knowledge_hobby:habit_lifestyle:digital_literacy_badged_course",
		},
		{
			name:        "role skill learning path",
			sourceQuery: "Microsoft Learn Azure learning path",
			goal:        "Use role based training modules to prepare for cloud task scenarios",
			wantPattern: "digital_creation:professional_transition:artifact_creation",
			wantSub:     "digital_creation:professional_transition:learning_path_role_skill",
		},
		{
			name:        "structured drawing foundation",
			sourceQuery: "Drawspace drawing lessons",
			goal:        "Learn contour drawing, shading, and perspective to finish a small sketch",
			wantPattern: "visual_art:artifact_creation:artifact_creation",
			wantSub:     "visual_art:artifact_creation:structured_drawing_foundation",
		},
		{
			name:        "watercolor landscape postcard beginner pattern",
			sourceQuery: "beginner watercolor landscape postcard",
			goal:        "Paint a small watercolor landscape postcard with simple washes, trees, and sky depth",
			wantPattern: "visual_art:artifact_creation:artifact_creation",
			wantSub:     "visual_art:artifact_creation:watercolor_beginner_landscape",
		},
		{
			name:        "urban sketch structured drawing foundation",
			sourceQuery: "어반스케치 동네 풍경",
			goal:        "어반스케치 기초 구도와 선 표현을 익혀 동네 풍경 한 장을 완성한다",
			wantPattern: "visual_art:artifact_creation:artifact_creation",
			wantSub:     "visual_art:artifact_creation:structured_drawing_foundation",
		},
		{
			name:        "creative tool beginner output",
			sourceQuery: "Envato Tuts+ Photoshop tutorial for beginners",
			goal:        "Complete a beginner creative tool output and prepare it for a portfolio",
			wantPattern: "digital_creation:artifact_creation:foundation_build",
			wantSub:     "digital_creation:artifact_creation:creative_tool_beginner_output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{
				SourceQuery:      tt.sourceQuery,
				LearningGoal:     &tt.goal,
				LearningLanguage: "en",
			}
			key := inferGoalPatternKey(req)
			if got := formatGoalPatternKeyForLog(key); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if got := inferGoalSubpatternKey(req); got != tt.wantSub {
				t.Fatalf("subpattern = %q, want %q", got, tt.wantSub)
			}
			summary := analyzeGoalRouting(req)
			if summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestEnglishActualGenerationQARoutingCases(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantSub     string
	}{
		{
			name:        "speaking interview certification",
			sourceQuery: "English speaking interview test",
			goal:        "I want to improve my spoken answers for an English speaking interview test",
			wantSub:     "language_communication:certification_assessment:english_speaking_interview_certification",
		},
		{
			name:        "leather wallet",
			sourceQuery: "beginner leather wallet project",
			goal:        "I want to make a simple leather card wallet with clean stitching and edge finishing",
			wantSub:     "craft_making:artifact_creation:leathercraft_wallet_project",
		},
		{
			name:        "maker arduino demo",
			sourceQuery: "Arduino maker workshop demo",
			goal:        "I want to build a small Arduino sensor demo and explain it in a workshop",
			wantSub:     "maker_technical_hobby:teaching_instruction:maker_workshop_demo",
		},
		{
			name:        "drone basic license",
			sourceQuery: "beginner drone operator basics",
			goal:        "I want to learn safe basic drone operation and prepare for a beginner drone license",
			wantSub:     "maker_technical_hobby:certification_assessment:drone_operator_basic_4class",
		},
		{
			name:        "climate data graph",
			sourceQuery: "climate data graph explanation",
			goal:        "I want to read climate data, make a simple graph, and explain the trend clearly",
			wantSub:     "knowledge_hobby:concept_mastery:climate_data_graphing_explanation",
		},
		{
			name:        "essay series publish",
			sourceQuery: "publish a short essay series",
			goal:        "I want to plan and publish a short weekly essay series online",
			wantSub:     "writing_storytelling:presentation_publish:brunch_serial_publish",
		},
		{
			name:        "ai digital literacy life practice",
			sourceQuery: "AI digital literacy for daily life",
			goal:        "I want to use AI tools safely for daily tasks like email, search, and planning",
			wantSub:     "knowledge_hobby:habit_lifestyle:ai_digital_literacy_life_practice",
		},
		{
			name:        "creator tutorial video",
			sourceQuery: "create a beginner tutorial video",
			goal:        "I want to create a beginner tutorial video that explains my process clearly",
			wantSub:     "digital_creation:teaching_instruction:creator_tutorial_publish",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{
				SourceQuery:      tt.sourceQuery,
				LearningGoal:     &tt.goal,
				LearningLanguage: "en",
			}
			if got := inferGoalSubpatternKey(req); got != tt.wantSub {
				t.Fatalf("subpattern = %q, want %q", got, tt.wantSub)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishMOOCOERGuidance(t *testing.T) {
	goal := "Learn programming through problem sets and build a small Python program"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "MIT OCW 6.0001 Introduction to Computer Science and Programming in Python",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recurring goal subtype: mooc programming foundation") {
		t.Fatalf("prompt missing English MOOC subpattern guidance: %q", prompt)
	}
	if !strings.Contains(prompt, "Recommended flow for this subtype") {
		t.Fatalf("prompt missing English MOOC flow guidance: %q", prompt)
	}
	if strings.Contains(prompt, "영어 MOOC 프로그래밍 기초형") || strings.Contains(prompt, "리슨은 MOOC") {
		t.Fatalf("English prompt should not expose Korean operational notes: %q", prompt)
	}
}

func TestBuildRecommendationGoalMetadata_EnglishMOOCOER(t *testing.T) {
	goal := "Learn contour drawing, shading, and perspective to finish a small sketch"
	meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
		SourceQuery:      "Drawspace drawing lessons",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if meta.SubpatternKey != "visual_art:artifact_creation:structured_drawing_foundation" {
		t.Fatalf("subpattern = %q", meta.SubpatternKey)
	}
	if !strings.Contains(meta.QueryHint, "contour drawing") {
		t.Fatalf("query hint = %q, want contour drawing", meta.QueryHint)
	}
}

func TestEnglishProjectOERLearningPathRouting(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		wantPattern string
		wantSub     string
	}{
		{
			name:        "problem set to final project",
			sourceQuery: "Harvard CS50x syllabus",
			goal:        "Build computer science foundations through problem sets and prepare a final project",
			wantPattern: "digital_creation:foundation_build:artifact_creation",
			wantSub:     "digital_creation:foundation_build:problem_set_to_final_project",
		},
		{
			name:        "web dev project foundation",
			sourceQuery: "The Odin Project Foundations",
			goal:        "Create HTML CSS JavaScript projects for a beginner web development portfolio",
			wantPattern: "digital_creation:artifact_creation:foundation_build",
			wantSub:     "digital_creation:artifact_creation:web_dev_project_foundation",
		},
		{
			name:        "web platform core curriculum",
			sourceQuery: "MDN Learn Web Development",
			goal:        "Learn HTML CSS JavaScript accessibility and web tools as a foundation",
			wantPattern: "digital_creation:foundation_build:artifact_creation",
			wantSub:     "digital_creation:foundation_build:web_platform_core_curriculum",
		},
		{
			name:        "frontend competency map",
			sourceQuery: "MDN Curriculum front-end developer competency map",
			goal:        "Prepare front-end developer competencies and job ready best practices",
			wantPattern: "digital_creation:professional_transition:artifact_creation",
			wantSub:     "digital_creation:professional_transition:frontend_competency_map",
		},
		{
			name:        "creative coding visual project",
			sourceQuery: "Khan Academy Intro to JS Drawing and Animation",
			goal:        "Use JavaScript drawing and animation to finish a creative coding project",
			wantPattern: "digital_creation:artifact_creation:foundation_build",
			wantSub:     "digital_creation:artifact_creation:creative_coding_visual_project",
		},
		{
			name:        "game 3d guided pathway",
			sourceQuery: "Unity Learn Junior Programmer pathway",
			goal:        "Build a small real-time 3D game prototype with Unity scripting",
			wantPattern: "digital_creation:artifact_creation:foundation_build",
			wantSub:     "digital_creation:artifact_creation:game_3d_guided_pathway",
		},
		{
			name:        "open textbook chapter practice",
			sourceQuery: "OpenStax open textbook subjects",
			goal:        "Study core chapters and exercises from an open textbook",
			wantPattern: "knowledge_hobby:foundation_build:artifact_creation",
			wantSub:     "knowledge_hobby:foundation_build:open_textbook_chapter_practice",
		},
		{
			name:        "oer learning object cluster",
			sourceQuery: "LibreTexts learning objects and simulations",
			goal:        "Use OER worksheets simulations and visualizations to understand key concepts",
			wantPattern: "knowledge_hobby:foundation_build:artifact_creation",
			wantSub:     "knowledge_hobby:foundation_build:oer_learning_object_cluster",
		},
		{
			name:        "interactive tutor concept practice",
			sourceQuery: "CMU OLI Introduction to Biology open course",
			goal:        "Practice biology concepts with an interactive tutor and feedback",
			wantPattern: "knowledge_hobby:foundation_build:artifact_creation",
			wantSub:     "knowledge_hobby:foundation_build:interactive_tutor_concept_practice",
		},
		{
			name:        "cloud lab skill path",
			sourceQuery: "Google Cloud Skills Boost role based learning path",
			goal:        "Practice cloud task scenarios with hands-on labs and prepare a skill badge",
			wantPattern: "digital_creation:professional_transition:artifact_creation",
			wantSub:     "digital_creation:professional_transition:cloud_lab_skill_path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateCourseDraftRequest{
				SourceQuery:      tt.sourceQuery,
				LearningGoal:     &tt.goal,
				LearningLanguage: "en",
			}
			if got := formatGoalPatternKeyForLog(inferGoalPatternKey(req)); got != tt.wantPattern {
				t.Fatalf("pattern = %q, want %q", got, tt.wantPattern)
			}
			if got := inferGoalSubpatternKey(req); got != tt.wantSub {
				t.Fatalf("subpattern = %q, want %q", got, tt.wantSub)
			}
			if summary := analyzeGoalRouting(req); summary.RefinementMode != "subpattern_guidance_only" {
				t.Fatalf("refinement mode = %q, want subpattern_guidance_only", summary.RefinementMode)
			}
		})
	}
}

func TestBuildCurriculumGenerationPromptIncludesEnglishProjectOERGuidance(t *testing.T) {
	goal := "Create HTML CSS JavaScript projects for a beginner web development portfolio"
	prompt := buildCurriculumGenerationPrompt(CreateCourseDraftRequest{
		SourceQuery:      "The Odin Project Foundations",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if !strings.Contains(prompt, "Recurring goal subtype: web dev project foundation") {
		t.Fatalf("prompt missing English web project subpattern guidance: %q", prompt)
	}
	if !strings.Contains(prompt, "Recommended flow for this subtype") {
		t.Fatalf("prompt missing English web project flow guidance: %q", prompt)
	}
	if strings.Contains(prompt, "영어 웹 개발 프로젝트 기초형") || strings.Contains(prompt, "리슨은 개발 도구") {
		t.Fatalf("English prompt should not expose Korean operational notes: %q", prompt)
	}
}

func TestBuildRecommendationGoalMetadata_EnglishProjectOER(t *testing.T) {
	goal := "Practice cloud task scenarios with hands-on labs and prepare a skill badge"
	meta := BuildRecommendationGoalMetadata(CreateCourseDraftRequest{
		SourceQuery:      "Google Cloud Skills Boost role based learning path",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})

	if meta.SubpatternKey != "digital_creation:professional_transition:cloud_lab_skill_path" {
		t.Fatalf("subpattern = %q", meta.SubpatternKey)
	}
	if !strings.Contains(meta.QueryHint, "hands-on lab") {
		t.Fatalf("query hint = %q, want hands-on lab", meta.QueryHint)
	}
}
