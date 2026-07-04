package curriculum

type goalPatternKey struct {
	DomainAxis        string
	GoalModePrimary   string
	GoalModeSecondary string
}

type goalPatternSpec struct {
	LessonCountMin          int
	LessonCountMax          int
	StepRoles               []string
	LastLessonRule          string
	CompletionCriteriaRules []string
	AtomicDisallowExamples  []string
}

type goalPatternRefinement struct {
	Refine            func([]generatedMainLesson) []generatedMainLesson
	CompletionCaption string
}

type goalPatternSoftRefinement struct {
	Apply             func(generatedDraftDocument) generatedDraftDocument
	CompletionCaption string
}

type goalSubpatternSpec struct {
	Name                    string
	StepRoles               []string
	StageRules              []string
	LastLessonRule          string
	CompletionCriteriaRules []string
	AtomicDisallowExamples  []string
}

type goalSubpatternRefinement struct {
	Refine            func([]generatedMainLesson) []generatedMainLesson
	CompletionCaption string
}

type goalRoutingSummary struct {
	PatternKey     goalPatternKey
	SubpatternKey  string
	RefinementMode string
}

var goalPatternSpecs = map[goalPatternKey]goalPatternSpec{
	{DomainAxis: "instrument_performance", GoalModePrimary: "performance_execution"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "first_output", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 목표 연주나 실제 연주 직전 준비 단계여야 합니다.",
		CompletionCriteriaRules: []string{"공연, 실전 무대, 공식 참여는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"튜닝", "기본 자세", "손 모양"},
	},
	{DomainAxis: "instrument_performance", GoalModePrimary: "participation_service"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 합주, 시작/마침, 템포 적응 같은 참여 직전 준비 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 합주 참여, 찬양단 봉사, 공식 반주는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"튜닝", "기본 자세", "손 모양"},
	},
	{DomainAxis: "language_communication", GoalModePrimary: "performance_execution"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 실제 대화 직전까지 표현 흐름을 연결하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 회화 수행은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"문법 설명", "표현 목록", "발음 설명"},
	},
	{DomainAxis: "language_communication", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 시험 직전까지 문제풀이와 응답 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"시험 응시와 인증 통과는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"문자표만 보기", "문법 정의 외우기"},
	},
	{DomainAxis: "craft_making", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 실기 시험 직전까지 작업 흐름과 위생 기준을 점검하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"도구 소개만 보기", "위생 규정만 읽기"},
	},
	{DomainAxis: "craft_making", GoalModePrimary: "artifact_creation"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "first_output", "core_pattern", "structure_plan", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 직접 만든 작품 한 점을 완성하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"판매 시작, 전시 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개", "첫 매듭", "기본 포인트 하나"},
	},
	{DomainAxis: "craft_making", GoalModePrimary: "teaching_instruction"}: {
		LessonCountMin:          5,
		LessonCountMax:          5,
		StepRoles:               []string{"first_output", "core_pattern", "structure_plan", "artifact_finish", "teaching_prep"},
		LastLessonRule:          "마지막 lesson은 제작 과정을 설명하고 시연할 준비 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 업로드와 강의 개시는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개", "기본 포인트 하나"},
	},
	{DomainAxis: "craft_making", GoalModePrimary: "habit_lifestyle"}: {
		LessonCountMin:          3,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "first_output", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 일상 루틴으로 이어질 수 있는 반복 흐름 정리 단계여야 합니다.",
		CompletionCriteriaRules: []string{"정기적인 루틴 유지 자체는 completion_criteria로 분리할 수 있습니다."},
		AtomicDisallowExamples:  []string{"재료 소개", "기본 포인트 하나"},
	},
	{DomainAxis: "visual_art", GoalModePrimary: "artifact_creation"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 완성도 있는 작품 한 점을 마무리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"전시와 판매는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개", "붓 설명", "기초 워시"},
	},
	{DomainAxis: "visual_art", GoalModePrimary: "habit_lifestyle"}: {
		LessonCountMin:          3,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "first_output", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 꾸준히 그릴 수 있는 개인 루틴 정리 단계여야 합니다.",
		CompletionCriteriaRules: []string{"매일 스케치 유지 자체는 completion_criteria로 분리할 수 있습니다."},
		AtomicDisallowExamples:  []string{"붓 설명", "도구 소개"},
	},
	{DomainAxis: "visual_art", GoalModePrimary: "professional_transition"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "structure_plan", "artifact_finish", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 첫 판매 직전까지 업로드와 운영 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 판매와 첫 고객 확보는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개", "도구 설명"},
	},
	{DomainAxis: "visual_art", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 실기 시험 직전까지 표현 순서와 결과물 완성 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"도구 소개만 보기", "색채 이론만 보기"},
	},
	{DomainAxis: "writing_storytelling", GoalModePrimary: "presentation_publish"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"first_output", "structure_plan", "integration_practice", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 공개 직전까지 글 묶음과 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"브런치 발행, 연재 시작, 출간은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"문장 이론 설명", "맞춤법만 보기"},
	},
	{DomainAxis: "writing_storytelling", GoalModePrimary: "teaching_instruction"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "first_output", "integration_practice", "structure_plan", "teaching_prep"},
		LastLessonRule:          "마지막 lesson은 글쓰기 과정을 설명할 수 있게 예시와 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 강의 오픈과 워크숍 진행은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"문장 이론 설명"},
	},
	{DomainAxis: "body_movement", GoalModePrimary: "participation_service"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 팀 참여나 파트너 흐름에 맞추는 실전 직전 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 공연 참여와 소셜댄스 참여는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"기본 자세 하나", "호흡 설명만 보기"},
	},
	{DomainAxis: "body_movement", GoalModePrimary: "habit_lifestyle"}: {
		LessonCountMin:          3,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 스스로 이어갈 수 있는 루틴 정리 단계여야 합니다.",
		CompletionCriteriaRules: []string{"주간 루틴 정착은 completion_criteria로 분리할 수 있습니다."},
		AtomicDisallowExamples:  []string{"호흡 설명", "정렬 설명"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "artifact_creation"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 결과물 한 편 또는 한 개를 완성하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 업로드와 포트폴리오 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"툴 인터페이스 설명만 보기"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "foundation_build"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 핵심 개념을 작은 산출물이나 문제해결 흐름으로 통합하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"공식 course 수료, certificate 발급, 별도 프로젝트 출시는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"설치만 하기", "문법 이름만 외우기", "퀴즈만 풀기"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "teaching_instruction"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "artifact_finish", "teaching_prep"},
		LastLessonRule:          "마지막 lesson은 제작 과정을 설명하고 시연할 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 강의 업로드와 튜토리얼 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"툴 소개만 보기", "단축키 설명만 보기"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "presentation_publish"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 게시 직전 품질 정리와 export 준비 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 게시와 공개 발행은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"툴 설명만 보기"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "professional_transition"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "structure_plan", "artifact_finish", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 포트폴리오, 채널, 판매 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 수주, 첫 판매, 프리랜서 시작은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"툴 설명만 보기"},
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 시험 직전까지 대표 시나리오와 문제풀이 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"서비스 목록만 보기", "단축키만 외우기"},
	},
	{DomainAxis: "knowledge_hobby", GoalModePrimary: "habit_lifestyle"}: {
		LessonCountMin:          3,
		LessonCountMax:          3,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 생활 속 반복 적용이나 관찰 루틴을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"정기적 실천 루틴은 completion_criteria로 분리할 수 있습니다."},
		AtomicDisallowExamples:  []string{"개념 설명만 보기"},
	},
	{DomainAxis: "knowledge_hobby", GoalModePrimary: "foundation_build"}: {
		LessonCountMin:          3,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 여러 개념, 사례, 텍스트를 비교·종합해 자기 언어로 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"공식 course 수료, paper 제출, 학점 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"용어 하나만 외우기", "작품명만 보기", "주차명 하나만 학습"},
	},
	{DomainAxis: "knowledge_hobby", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 시험 직전까지 주요 개념과 문제풀이 흐름을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"개념 정의만 보기"},
	},
	{DomainAxis: "maker_technical_hobby", GoalModePrimary: "artifact_creation"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "first_output", "core_pattern", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 작동하는 프로젝트나 제작물을 완성하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"판매와 외부 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"부품 소개만 보기"},
	},
	{DomainAxis: "maker_technical_hobby", GoalModePrimary: "teaching_instruction"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "first_output", "core_pattern", "artifact_finish", "teaching_prep"},
		LastLessonRule:          "마지막 lesson은 제작 원리와 시연 흐름을 설명할 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 워크숍 진행과 튜토리얼 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"부품 소개만 보기", "회로 설명만 보기"},
	},
	{DomainAxis: "maker_technical_hobby", GoalModePrimary: "professional_transition"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "structure_plan", "artifact_finish", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 상품화, 채널화, 판매 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"첫 판매와 첫 수주는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"부품 소개만 보기"},
	},
	{DomainAxis: "maker_technical_hobby", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "integration_practice", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 실기나 시험 직전까지 안전 기준과 작업 흐름을 점검하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"공구 소개만 보기", "법규 정의만 보기"},
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "teaching_instruction"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "artifact_finish", "teaching_prep"},
		LastLessonRule:          "마지막 lesson은 레시피와 시연 흐름을 설명할 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 클래스 오픈과 강의 공개는 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개만 보기", "도구 소개만 보기"},
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "artifact_creation"}: {
		LessonCountMin:          3,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "core_pattern", "artifact_finish"},
		LastLessonRule:          "마지막 lesson은 대표 디저트나 요리 한 점을 완성하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"판매와 클래스 시작은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개만 보기"},
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "habit_lifestyle"}: {
		LessonCountMin:          3,
		LessonCountMax:          3,
		StepRoles:               []string{"setup_intro", "first_output", "integration_practice"},
		LastLessonRule:          "마지막 lesson은 일상 식습관이나 조리 루틴을 정리하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"주간 식단 루틴과 건강식 습관 유지는 completion_criteria로 분리할 수 있습니다."},
		AtomicDisallowExamples:  []string{"재료 소개만 보기"},
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "professional_transition"}: {
		LessonCountMin:          4,
		LessonCountMax:          4,
		StepRoles:               []string{"setup_intro", "structure_plan", "artifact_finish", "publish_prep"},
		LastLessonRule:          "마지막 lesson은 브랜드, 메뉴, 판매 채널 준비를 마친 단계여야 합니다.",
		CompletionCriteriaRules: []string{"첫 주문과 첫 매출은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개만 보기"},
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "certification_assessment"}: {
		LessonCountMin:          4,
		LessonCountMax:          5,
		StepRoles:               []string{"setup_intro", "core_pattern", "artifact_finish", "performance_prep"},
		LastLessonRule:          "마지막 lesson은 실기 시험 직전까지 작업 순서와 시간 배분을 점검하는 단계여야 합니다.",
		CompletionCriteriaRules: []string{"실제 시험 응시와 자격 취득은 completion_criteria로 분리합니다."},
		AtomicDisallowExamples:  []string{"재료 소개만 보기", "위생 설명만 보기"},
	},
}

var goalPatternRefinements = map[goalPatternKey]goalPatternRefinement{
	{DomainAxis: "instrument_performance", GoalModePrimary: "participation_service"}: {
		Refine:            refineWorshipGuitarLessons,
		CompletionCaption: "찬양단 연습 흐름에 맞춰 반주를 준비하고 연결할 수 있다.",
	},
	{DomainAxis: "craft_making", GoalModePrimary: "teaching_instruction"}: {
		Refine:            refineCrochetTeachingLessons,
		CompletionCaption: "직접 만든 작품의 제작 과정을 설명하고 시연할 준비가 되어 있다.",
	},
	{DomainAxis: "language_communication", GoalModePrimary: "performance_execution"}: {
		Refine:            refineConversationPerformanceLessons,
		CompletionCaption: "목표 상황에서 짧은 질문과 응답을 이어가며 실제 대화를 유지할 수 있다.",
	},
	{DomainAxis: "writing_storytelling", GoalModePrimary: "presentation_publish"}: {
		Refine:            refineWritingPresentationLessons,
		CompletionCaption: "공개를 염두에 둔 글 묶음을 스스로 정리하고 다듬을 수 있다.",
	},
	{DomainAxis: "visual_art", GoalModePrimary: "artifact_creation"}: {
		Refine:            refineVisualArtArtifactLessons,
		CompletionCaption: "표현 요소를 조합해 완성도 있는 작품 한 점을 스스로 완성할 수 있다.",
	},
	{DomainAxis: "language_communication", GoalModePrimary: "certification_assessment"}: {
		Refine:            refineLanguageCertificationLessons,
		CompletionCaption: "목표 시험 유형에서 핵심 문제 흐름을 따라 답할 수 있다.",
	},
	{DomainAxis: "body_movement", GoalModePrimary: "habit_lifestyle"}: {
		Refine:            refineBodyMovementHabitLessons,
		CompletionCaption: "짧은 루틴을 스스로 이어가며 일상에 정착시킬 수 있다.",
	},
	{DomainAxis: "maker_technical_hobby", GoalModePrimary: "professional_transition"}: {
		Refine:            refineMakerProfessionalLessons,
		CompletionCaption: "작동하는 결과물과 소개 흐름을 바탕으로 판매나 의뢰를 준비할 수 있다.",
	},
	{DomainAxis: "visual_art", GoalModePrimary: "habit_lifestyle"}: {
		Refine:            refineVisualArtHabitLessons,
		CompletionCaption: "짧은 스케치나 드로잉 루틴을 일상에서 스스로 이어갈 수 있다.",
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "artifact_creation"}: {
		Refine:            refineDigitalCreationArtifactLessons,
		CompletionCaption: "기획한 디지털 결과물 한 편 또는 한 개를 스스로 완성할 수 있다.",
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "professional_transition"}: {
		Refine:            refineDigitalCreationProfessionalLessons,
		CompletionCaption: "완성한 디지털 결과물과 소개 자료를 바탕으로 판매나 의뢰를 준비할 수 있다.",
	},
	{DomainAxis: "visual_art", GoalModePrimary: "professional_transition"}: {
		Refine:            refineVisualArtProfessionalLessons,
		CompletionCaption: "완성한 작품과 소개 자료를 바탕으로 판매나 의뢰를 준비할 수 있다.",
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "artifact_creation"}: {
		Refine:            refineCookingArtifactLessons,
		CompletionCaption: "대표 요리나 디저트 결과물 한 점을 스스로 완성할 수 있다.",
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "professional_transition"}: {
		Refine:            refineCookingProfessionalLessons,
		CompletionCaption: "완성한 메뉴와 소개 흐름을 바탕으로 판매나 주문을 준비할 수 있다.",
	},
	{DomainAxis: "digital_creation", GoalModePrimary: "presentation_publish"}: {
		Refine:            refineDigitalCreationPublishLessons,
		CompletionCaption: "완성한 디지털 결과물을 공개 직전까지 정리하고 게시 준비를 마칠 수 있다.",
	},
	{DomainAxis: "writing_storytelling", GoalModePrimary: "teaching_instruction"}: {
		Refine:            refineWritingTeachingLessons,
		CompletionCaption: "자신의 글쓰기 과정을 예시와 함께 설명하고 강의 흐름을 정리할 수 있다.",
	},
	{DomainAxis: "craft_making", GoalModePrimary: "artifact_creation"}: {
		Refine:            refineCraftArtifactLessons,
		CompletionCaption: "작품 한 점을 기획하고 완성하는 흐름을 스스로 이어갈 수 있다.",
	},
	{DomainAxis: "body_movement", GoalModePrimary: "participation_service"}: {
		Refine:            refineBodyMovementParticipationLessons,
		CompletionCaption: "팀이나 파트너 흐름에 맞춰 실제 참여 직전까지 준비할 수 있다.",
	},
	{DomainAxis: "knowledge_hobby", GoalModePrimary: "habit_lifestyle"}: {
		Refine:            refineKnowledgeHabitLessons,
		CompletionCaption: "관찰과 기록, 선택 기준을 일상 루틴으로 이어갈 수 있다.",
	},
}

var goalPatternSoftRefinements = map[goalPatternKey]goalPatternSoftRefinement{
	{DomainAxis: "instrument_performance", GoalModePrimary: "performance_execution"}: {
		Apply:             softenInstrumentPerformanceExecution,
		CompletionCaption: "목표 곡이나 연주 과제를 실전 직전까지 스스로 준비할 수 있다.",
	},
	{DomainAxis: "writing_storytelling", GoalModePrimary: "habit_lifestyle"}: {
		Apply:             softenWritingHabitLessons,
		CompletionCaption: "짧은 글쓰기를 일상 루틴으로 이어갈 수 있다.",
	},
	{DomainAxis: "visual_art", GoalModePrimary: "presentation_publish"}: {
		Apply:             softenVisualArtPublishLessons,
		CompletionCaption: "공개나 전시를 염두에 둔 작품 묶음을 스스로 정리할 수 있다.",
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "teaching_instruction"}: {
		Apply:             softenCookingTeachingLessons,
		CompletionCaption: "레시피와 시연 흐름을 정리해 설명 준비를 할 수 있다.",
	},
	{DomainAxis: "cooking_baking", GoalModePrimary: "habit_lifestyle"}: {
		Apply:             softenCookingHabitLessons,
		CompletionCaption: "건강한 식사 준비를 반복 가능한 생활 루틴으로 이어갈 수 있다.",
	},
}
