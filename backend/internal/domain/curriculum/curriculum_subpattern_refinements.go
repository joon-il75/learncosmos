package curriculum

var goalSubpatternRefinements = map[string]goalSubpatternRefinement{
	"instrument_performance:participation_service:worship_team_support": {
		Refine:            refineWorshipTeamSupportLessons,
		CompletionCaption: "찬양단 합주 흐름에 맞춰 반주를 준비하고 연결할 수 있다.",
	},
	"instrument_performance:performance_execution:song_completion": {
		Refine:            refineSongCompletionLessons,
		CompletionCaption: "목표 곡을 처음부터 끝까지 끊기지 않게 이어갈 수 있다.",
	},
	"writing_storytelling:presentation_publish:brunch_serial_publish": {
		Refine:            refineBrunchSerialPublishLessons,
		CompletionCaption: "브런치나 공개 플랫폼에 올릴 글 묶음을 스스로 정리하고 다듬을 수 있다.",
	},
	"writing_storytelling:habit_lifestyle:daily_writing_habit": {
		Refine:            refineDailyWritingHabitLessons,
		CompletionCaption: "짧은 글쓰기를 부담 없이 일상 루틴으로 이어갈 수 있다.",
	},
	"craft_making:teaching_instruction:craft_teach_youtube": {
		Refine:            refineCraftTeachYoutubeLessons,
		CompletionCaption: "직접 만든 작품의 제작 과정을 설명하고 시연할 준비가 되어 있다.",
	},
	"craft_making:artifact_creation:leathercraft_basic_accessory": {
		Refine:            refineLeathercraftBasicAccessoryLessons,
		CompletionCaption: "기본 공정을 따라 작은 가죽 소품 한 점을 완성하고 마감 상태를 점검할 수 있다.",
	},
	"craft_making:artifact_creation:leathercraft_wallet_project": {
		Refine:            refineLeathercraftWalletProjectLessons,
		CompletionCaption: "실제로 사용할 수 있는 가죽 지갑 한 점을 완성하고 수납과 마감 품질을 점검할 수 있다.",
	},
	"craft_making:artifact_creation:leathercraft_dimensional_bag": {
		Refine:            refineLeathercraftDimensionalBagLessons,
		CompletionCaption: "입체 구조가 유지되는 가죽 가방이나 파우치 한 점을 완성하고 조립 순서를 설명할 수 있다.",
	},
	"language_communication:performance_execution:travel_conversation": {
		Refine:            refineTravelConversationLessons,
		CompletionCaption: "여행 상황에서 짧은 질문과 응답을 이어가며 실제 대화를 시작할 수 있다.",
	},
	"visual_art:artifact_creation:calligraphy_basic_lettering": {
		Refine:            refineCalligraphyBasicLetteringLessons,
		CompletionCaption: "붓펜과 기본 획을 활용해 짧은 문장을 안정적인 캘리그라피 결과물로 완성할 수 있다.",
	},
	"visual_art:artifact_creation:calligraphy_quote_art_project": {
		Refine:            refineCalligraphyQuoteArtProjectLessons,
		CompletionCaption: "문구와 구도를 정해 선물하거나 보관할 수 있는 캘리그라피 작품 한 점을 완성할 수 있다.",
	},
	"visual_art:presentation_publish:calligraphy_publish_showcase": {
		Refine:            refineCalligraphyPublishShowcaseLessons,
		CompletionCaption: "공개할 캘리그라피 작품 묶음과 소개 흐름을 스스로 정리할 수 있다.",
	},
	"visual_art:presentation_publish:art_publish_showcase": {
		Refine:            refineArtPublishShowcaseLessons,
		CompletionCaption: "공개나 전시를 염두에 둔 작품 묶음을 스스로 정리할 수 있다.",
	},
	"digital_creation:presentation_publish:portfolio_publish": {
		Refine:            refinePortfolioPublishLessons,
		CompletionCaption: "완성한 디지털 결과물을 공개 직전까지 정리하고 게시 준비를 마칠 수 있다.",
	},
	"digital_creation:teaching_instruction:creator_tutorial_publish": {
		Refine:            refineCreatorTutorialPublishLessons,
		CompletionCaption: "완성한 디지털 결과물의 제작 과정을 설명하고 튜토리얼 공개를 준비할 수 있다.",
	},
	"digital_creation:artifact_creation:freeware_midi_beatmaking": {
		Refine:            refineFreewareMIDIBeatmakingLessons,
		CompletionCaption: "무료 도구에서 짧은 MIDI beat 한 개를 완성하고 재생·저장 상태를 점검할 수 있다.",
	},
	"digital_creation:artifact_creation:freeware_midi_full_track": {
		Refine:            refineFreewareMIDIFullTrackLessons,
		CompletionCaption: "무료 도구에서 짧은 곡 또는 전체 트랙을 export 직전 상태까지 정리할 수 있다.",
	},
	"digital_creation:artifact_creation:midi_foundation_workflow": {
		Refine:            refineMIDIFoundationWorkflowLessons,
		CompletionCaption: "MIDI 송수신과 note 입력 흐름을 이해하고 짧은 MIDI phrase를 만들어 재생할 수 있다.",
	},
	"digital_creation:artifact_creation:vibe_coding_mvp_app": {
		Refine:            refineVibeCodingMVPAppLessons,
		CompletionCaption: "AI 도구로 작은 앱 또는 MVP를 만들고 핵심 사용자 흐름을 직접 점검할 수 있다.",
	},
	"digital_creation:artifact_creation:vibe_coding_fullstack_ship": {
		Refine:            refineVibeCodingFullstackShipLessons,
		CompletionCaption: "AI로 만든 full-stack 앱의 데이터, 인증, 오류, 보안 상태를 배포 직전까지 점검할 수 있다.",
	},
	"digital_creation:artifact_creation:claude_code_agentic_workflow": {
		Refine:            refineClaudeCodeAgenticWorkflowLessons,
		CompletionCaption: "Claude Code 기반 반복 개발 workflow를 작은 프로젝트에 적용하고 테스트·정리까지 마칠 수 있다.",
	},
	"maker_technical_hobby:teaching_instruction:maker_workshop_demo": {
		Refine:            refineMakerWorkshopDemoLessons,
		CompletionCaption: "작동하는 프로젝트의 제작 과정과 시연 흐름을 설명할 준비를 할 수 있다.",
	},
	"cooking_baking:teaching_instruction:baking_class_demo": {
		Refine:            refineBakingClassDemoLessons,
		CompletionCaption: "완성한 메뉴의 레시피와 시연 흐름을 설명할 준비를 할 수 있다.",
	},
	"knowledge_hobby:professional_transition:credit_bank_practicum": {
		Refine:            refineCreditBankPracticumLessons,
		CompletionCaption: "기관실습과 사후정리를 마치고 실습 인정 직전 상태까지 준비할 수 있다.",
	},
	"knowledge_hobby:professional_transition:credit_bank_standard_theory": {
		Refine:            refineCreditBankStandardTheoryLessons,
		CompletionCaption: "핵심 이론을 사례와 실무 관점으로 연결해 현장 적용 직전 수준까지 정리할 수 있다.",
	},
	"knowledge_hobby:habit_lifestyle:university_survey": {
		Refine:            refineUniversitySurveyLessons,
		CompletionCaption: "핵심 텍스트와 사례를 비교·비평 관점으로 연결해 종합적으로 정리할 수 있다.",
	},
	"knowledge_hobby:habit_lifestyle:kmooc_online_course_survey": {
		Refine:            refineKMOOCOnlineCourseSurveyLessons,
		CompletionCaption: "공개강좌의 핵심 주제와 사례 적용 흐름을 평가 직전 수준으로 통합 정리할 수 있다.",
	},
	"writing_storytelling:presentation_publish:university_studio": {
		Refine:            refineUniversityStudioLessons,
		CompletionCaption: "작은 초안에서 출발해 제출 가능한 작업물 직전 수준까지 스스로 발전시킬 수 있다.",
	},
	"visual_art:artifact_creation:university_studio": {
		Refine:            refineUniversityStudioLessons,
		CompletionCaption: "작은 초안에서 출발해 제출 가능한 작업물 직전 수준까지 스스로 발전시킬 수 있다.",
	},
	"digital_creation:artifact_creation:university_studio": {
		Refine:            refineUniversityStudioLessons,
		CompletionCaption: "작은 초안에서 출발해 제출 가능한 작업물 직전 수준까지 스스로 발전시킬 수 있다.",
	},
	"knowledge_hobby:teaching_instruction:credit_bank_instructional_design": {
		Refine:            refineCreditBankInstructionalDesignLessons,
		CompletionCaption: "교안이나 평가 설계를 정리해 실제 수업 운영 직전 상태까지 준비할 수 있다.",
	},
}
