package curriculum

import "strings"

func inferGoalPatternKey(req CreateCourseDraftRequest) goalPatternKey {
	context := buildGoalInferenceContext(req)
	domain := inferDomainAxis(context)
	primary := inferGoalModePrimary(context, domain)
	secondary := inferGoalModeSecondary(context, domain, primary)
	return goalPatternKey{DomainAxis: domain, GoalModePrimary: primary, GoalModeSecondary: secondary}
}

func inferDomainAxis(context string) string {
	switch {
	case containsPublicLifelongOnlineCourseContext(context):
		return "knowledge_hobby"
	case containsBlendedLifelongParticipationContext(context):
		return "knowledge_hobby"
	case containsMidlifeTransitionLearningSupportContext(context):
		return "knowledge_hobby"
	case containsCompletionRefundOnlineCourseContext(context):
		return "knowledge_hobby"
	case containsFederatedLifelongOpenResourceContext(context):
		return "knowledge_hobby"
	case containsKoreanContentCreationPipelineContext(context):
		return "digital_creation"
	case containsKoreanDesignToolPracticalCertificationContext(context):
		return "digital_creation"
	case containsKoreanNCSWebPublisherTrainingContext(context):
		return "digital_creation"
	case containsRuralReturnAgricultureFoundationContext(context):
		return "knowledge_hobby"
	case containsUrbanAgriculturePracticeSeriesContext(context):
		return "knowledge_hobby"
	case containsAIDigitalLiteracyLifePracticeContext(context):
		return "knowledge_hobby"
	case containsLocalHobbyWorkshopSeriesContext(context):
		return "craft_making"
	case containsScienceSimulationInquiryPathContext(context):
		return "knowledge_hobby"
	case containsGuidedSimulationActivitySheetContext(context):
		return "knowledge_hobby"
	case containsFoodPreservationSafetyPathContext(context):
		return "cooking_baking"
	case containsCreatorTutorialPublishContext(context) && !containsAny(context, "코바늘", "뜨개", "뜨개질", "가죽공예", "도자기", "공예"):
		return "digital_creation"
	case containsWebDevProjectFoundationContext(context):
		return "digital_creation"
	case containsHorticultureDiagnosticFoundationContext(context):
		return "knowledge_hobby"
	case containsAIToolProductivityWorkflowContext(context):
		return "digital_creation"
	case containsWatercolorBeginnerLandscapeContext(context):
		return "visual_art"
	case containsRunning5KBeginnerRoutineContext(context):
		return "body_movement"
	case containsDailyEnglishConversationRoutineContext(context):
		return "language_communication"
	case containsArduinoSensorProjectContext(context):
		return "maker_technical_hobby"
	case containsKoreanHomeCookingBasicsContext(context):
		return "cooking_baking"
	case containsComputationalMusicAnalysisProjectContext(context):
		return "digital_creation"
	case containsDesignThinkingProblemFramingPathContext(context):
		return "maker_technical_hobby"
	case containsNoBuildPrototypeTestContext(context):
		return "maker_technical_hobby"
	case containsClimateDataGraphingExplanationContext(context):
		return "knowledge_hobby"
	case containsCitizenScienceObservationRecordContext(context):
		return "knowledge_hobby"
	case containsPrimarySourceInquiryNoteContext(context):
		return "knowledge_hobby"
	case containsSewingProjectFoundationContext(context):
		return "craft_making"
	case containsWoodworkingSmallProjectContext(context):
		return "craft_making"
	case containsKnittingBasicWearableContext(context):
		return "craft_making"
	case containsCrochetSmallDollProjectContext(context):
		return "craft_making"
	case containsPotteryHandbuildingCupProjectContext(context):
		return "craft_making"
	case containsCardboardCircuitInventionPathContext(context):
		return "maker_technical_hobby"
	case containsMusicProductionFullPipelineContext(context):
		return "digital_creation"
	case containsPhotographyFundamentalsProjectContext(context):
		return "visual_art"
	case containsMealPrepWeekdayLunchContext(context):
		return "cooking_baking"
	case containsHouseplantRepottingCareContext(context):
		return "knowledge_hobby"
	case containsPersonalFinanceBudgetRoutineContext(context):
		return "knowledge_hobby"
	case containsDogBasicTrainingRoutineContext(context):
		return "knowledge_hobby"
	case containsNotionStudyDashboardContext(context):
		return "digital_creation"
	case containsGardenDesignMaintenancePathContext(context):
		return "knowledge_hobby"
	case containsHomeCafeLatteFoundationContext(context):
		return "cooking_baking"
	case containsCookingBasicsFoundationContext(context):
		return "cooking_baking"
	case containsCriticalReadingToWritingDraftContext(context):
		return "writing_storytelling"
	case containsPhotographyStudioProjectContext(context):
		return "visual_art"
	case containsDesignLabPrototypeProjectContext(context):
		return "maker_technical_hobby"
	case containsOrganicGrowingCyclePlanContext(context):
		return "knowledge_hobby"
	case containsAdobeCreativeLearningPathContext(context):
		return "digital_creation"
	case containsTemplateDesignQuickOutputContext(context):
		return "digital_creation"
	case containsSequentialCreativeClassPathContext(context):
		return "visual_art"
	case containsAnimationFundamentalsShotProgressionContext(context):
		return "digital_creation"
	case containsInteractiveMusicProductionFoundationContext(context):
		return "digital_creation"
	case containsSongDrivenInstrumentPathContext(context):
		return "instrument_performance"
	case containsMakerProjectClassContext(context):
		return "maker_technical_hobby"
	case containsShortDiplomaAssessmentPathContext(context):
		return "knowledge_hobby"
	case containsShortCourseWeeklyDiscussionContext(context):
		return "knowledge_hobby"
	case containsAuditCoursePracticePathContext(context):
		return "knowledge_hobby"
	case containsProblemSetFinalProjectContext(context):
		return "digital_creation"
	case containsWebDevProjectFoundationContext(context):
		return "digital_creation"
	case containsWebPlatformCoreCurriculumContext(context):
		return "digital_creation"
	case containsFrontendCompetencyMapContext(context):
		return "digital_creation"
	case containsCreativeCodingVisualProjectContext(context):
		return "digital_creation"
	case containsGame3DGuidedPathwayContext(context):
		return "digital_creation"
	case containsOpenTextbookChapterPracticeContext(context):
		return "knowledge_hobby"
	case containsOERLearningObjectClusterContext(context):
		return "knowledge_hobby"
	case containsInteractiveTutorConceptPracticeContext(context):
		return "knowledge_hobby"
	case containsCloudLabSkillPathContext(context):
		return "digital_creation"
	case containsEnglishMOOCProgrammingContext(context):
		return "digital_creation"
	case containsKoreanOpenLectureWeeklySurveyContext(context):
		return "knowledge_hobby"
	case containsEnglishMOOCHumanitiesSurveyContext(context):
		return "knowledge_hobby"
	case containsDigitalLiteracyBadgedCourseContext(context):
		return "knowledge_hobby"
	case containsRoleSkillLearningPathContext(context):
		return "digital_creation"
	case containsWatercolorBeginnerLandscapeContext(context):
		return "visual_art"
	case containsStructuredDrawingFoundationContext(context):
		return "visual_art"
	case containsVideoEditingShortformProjectContext(context):
		return "digital_creation"
	case containsCreativeToolBeginnerOutputContext(context):
		return "digital_creation"
	case containsKMOOCOnlineCourseContext(context):
		return "knowledge_hobby"
	case containsUniversitySurveyContext(context):
		return "knowledge_hobby"
	case containsCreditBankStandardTheoryContext(context):
		return "knowledge_hobby"
	case containsCreditBankInstructionalDesignContext(context):
		return "knowledge_hobby"
	case containsCreditBankPracticumContext(context):
		return "knowledge_hobby"
	case containsUniversityStudioWritingContext(context):
		return "writing_storytelling"
	case containsUniversityStudioVisualContext(context):
		return "visual_art"
	case containsUniversityStudioCreativeTechContext(context):
		return "digital_creation"
	case containsVibeCodingContext(context):
		return "digital_creation"
	case containsMIDIProductionContext(context):
		return "digital_creation"
	case containsVocalSongPracticeContext(context):
		return "instrument_performance"
	case containsInstrumentPerformanceContext(context):
		return "instrument_performance"
	case containsAny(context, "메이크업", "makeup"):
		return "visual_art"
	case containsBeautyServicePracticalContext(context):
		return "craft_making"
	case containsAny(context, "그림", "수채화", "드로잉", "스케치", "일러스트", "페인팅", "캘리그라피", "손글씨", "붓펜", "레터링", "calligraphy", "lettering", "brush pen", "watercolor", "landscape postcard"):
		return "visual_art"
	case containsCandleSoapResinProjectContext(context):
		return "craft_making"
	case containsAny(context, "코바늘", "뜨개", "뜨개질", "가죽공예", "도자기", "공예",
		"leathercraft", "leather craft", "leather wallet", "leatherworking", "leather working",
		"woodworking", "wood work", "woodcraft", "birdhouse", "bird house", "crochet", "granny square", "pottery", "ceramic", "glazing"):
		return "craft_making"
	case containsAny(context, "글쓰기", "에세이", "브런치", "연재", "출간", "소설", "카피",
		"writing", "journal", "journaling", "reflection", "reflections", "write every day",
		"essay", "essay series", "publish a short essay", "weekly essay"):
		return "writing_storytelling"
	case containsAny(context, "영어", "일본어", "회화", "대화", "jlpt", "toeic", "toefl", "opic", "언어",
		"travel english", "travel conversation", "conversation", "speaking", "airport", "hotel", "restaurant"):
		return "language_communication"
	case containsAny(context, "영상편집", "프리미어", "다빈치", "캔바", "유튜브 편집", "youtube", "디지털",
		"포토샵", "photoshop", "gtq", "그래픽기술자격", "웹디자인기능사", "웹퍼블리셔", "프론트엔드",
		"미디", "midi", "daw", "bandlab", "밴드랩", "lmms", "tunepad", "soundtrap", "음악제작", "음악 제작", "뮤직 프로덕션", "music production", "beat making", "비트메이킹", "비트 만들기",
		"바이브 코딩", "바이브코딩", "vibe coding", "ai coding", "ai 코딩", "cursor", "lovable", "replit", "claude code", "github copilot", "copilot", "supabase", "vercel",
		"video editing", "portfolio", "publish a video", "short video", "vertical video", "podcast", "first podcast", "audio episode", "web app", "full-stack", "full stack", "authentication", "database", "notion", "study dashboard",
		"aws", "cloud practitioner", "solutions architect", "ncp", "naver cloud", "google cloud", "gcp", "azure", "az-900", "az 900", "az-104", "az 104",
		"associate cloud engineer", "digital leader", "정보처리기사", "컴퓨터활용능력", "컴활",
		"office spreadsheet certification", "spreadsheet certification", "excel certification", "office tool certification"):
		return "digital_creation"
	case containsDanceCoverRoutineContext(context):
		return "body_movement"
	case containsHomeStrengthRoutineContext(context):
		return "body_movement"
	case containsYogaDailyRoutineContext(context):
		return "body_movement"
	case containsAny(context, "요가", "필라테스", "댄스", "춤", "살사", "운동", "yoga", "dance", "stretching", "flexibility", "workout", "exercise", "fitness"):
		return "body_movement"
	case containsAny(context, "베이킹", "제과", "제빵", "케이크", "디저트", "요리", "파티세리", "파티셰", "한식조리기능사", "조리기능사",
		"baking", "cake", "dessert", "sponge cake", "whipped cream", "cook", "cooking", "meal prep", "healthy lunch"):
		return "cooking_baking"
	case containsAny(context, "와인", "천문", "관찰", "별보기", "바둑", "직업상담사", "대기환경기사",
		"garden", "gardening", "herb", "herbs", "balcony garden", "container garden",
		"houseplant", "repotting", "personal finance", "monthly budget", "spending categories", "dog training", "positive reinforcement",
		"climate data", "data graph", "graph explanation", "data trend",
		"counseling certification", "counselor certification", "environmental engineering certification",
		"environmental engineer exam", "biodiversity observation", "inaturalist", "primary source analysis"):
		return "knowledge_hobby"
	case containsAny(context, "아두이노", "3d 프린", "3d프린", "3d printing", "메이커", "전자회로",
		"드론", "초경량비행장치", "지게차", "건축기사", "건축산업기사", "토목기사", "전기기사", "전기산업기사", "전기기능사", "전기공사산업기사",
		"용접", "피복아크", "자동차정비", "정비기능사",
		"arduino", "sensor demo", "maker workshop", "drone", "drone license", "drone operator",
		"cardboard circuit", "circuit invention", "led", "leds", "switches", "design lab", "prototype for user feedback"):
		return "maker_technical_hobby"
	default:
		return ""
	}
}

func inferGoalSubpatternKey(req CreateCourseDraftRequest) string {
	context := buildGoalInferenceContext(req)
	pattern := inferGoalPatternKey(req)

	switch {
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsPublicLifelongOnlineCourseContext(context):
		return "knowledge_hobby:habit_lifestyle:public_lifelong_online_course"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "participation_service" && containsBlendedLifelongParticipationContext(context):
		return "knowledge_hobby:participation_service:blended_lifelong_participation"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsMidlifeTransitionLearningSupportContext(context):
		return "knowledge_hobby:professional_transition:midlife_transition_learning_support"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsCompletionRefundOnlineCourseContext(context):
		return "knowledge_hobby:professional_transition:completion_refund_online_course"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsFederatedLifelongOpenResourceContext(context):
		return "knowledge_hobby:foundation_build:federated_lifelong_open_resource"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "teaching_instruction" && containsCreatorTutorialPublishContext(context):
		return "digital_creation:teaching_instruction:creator_tutorial_publish"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsKoreanContentCreationPipelineContext(context):
		return "digital_creation:artifact_creation:korean_content_creation_pipeline"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsKoreanDesignToolPracticalCertificationContext(context):
		return "digital_creation:certification_assessment:korean_design_tool_practical_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "professional_transition" && containsKoreanNCSWebPublisherTrainingContext(context):
		return "digital_creation:professional_transition:korean_ncs_web_publisher_training"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsRuralReturnAgricultureFoundationContext(context):
		return "knowledge_hobby:professional_transition:rural_return_agriculture_foundation"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsUrbanAgriculturePracticeSeriesContext(context):
		return "knowledge_hobby:habit_lifestyle:urban_agriculture_practice_series"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsAIDigitalLiteracyLifePracticeContext(context):
		return "knowledge_hobby:habit_lifestyle:ai_digital_literacy_life_practice"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsLocalHobbyWorkshopSeriesContext(context):
		return "craft_making:artifact_creation:local_hobby_workshop_series"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "concept_mastery" && containsScienceSimulationInquiryPathContext(context):
		return "knowledge_hobby:concept_mastery:science_simulation_inquiry_path"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "concept_mastery" && containsGuidedSimulationActivitySheetContext(context):
		return "knowledge_hobby:concept_mastery:guided_simulation_activity_sheet"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "habit_lifestyle" && containsFoodPreservationSafetyPathContext(context):
		return "cooking_baking:habit_lifestyle:food_preservation_safety_path"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "habit_lifestyle" && containsMealPrepWeekdayLunchContext(context):
		return "cooking_baking:habit_lifestyle:meal_prep_weekday_lunch"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "concept_mastery" && containsHorticultureDiagnosticFoundationContext(context):
		return "knowledge_hobby:concept_mastery:horticulture_diagnostic_foundation"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsWatercolorBeginnerLandscapeContext(context):
		return "visual_art:artifact_creation:watercolor_beginner_landscape"
	case pattern.DomainAxis == "body_movement" && pattern.GoalModePrimary == "habit_lifestyle" && containsRunning5KBeginnerRoutineContext(context):
		return "body_movement:habit_lifestyle:running_5k_beginner_routine"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "performance_execution" && containsDailyEnglishConversationRoutineContext(context):
		return "language_communication:performance_execution:daily_english_conversation_routine"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "artifact_creation" && containsKoreanHomeCookingBasicsContext(context):
		return "cooking_baking:artifact_creation:korean_home_cooking_basics"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsHouseplantRepottingCareContext(context):
		return "knowledge_hobby:habit_lifestyle:houseplant_repotting_care"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsPersonalFinanceBudgetRoutineContext(context):
		return "knowledge_hobby:habit_lifestyle:personal_finance_budget_routine"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsDogBasicTrainingRoutineContext(context):
		return "knowledge_hobby:habit_lifestyle:dog_basic_training_routine"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsNotionStudyDashboardContext(context):
		return "digital_creation:artifact_creation:notion_study_dashboard"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "technical_skill" && containsAIToolProductivityWorkflowContext(context):
		return "digital_creation:technical_skill:ai_tool_productivity_workflow"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "technical_skill" && containsComputationalMusicAnalysisProjectContext(context):
		return "digital_creation:technical_skill:computational_music_analysis_project"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsDesignThinkingProblemFramingPathContext(context):
		return "maker_technical_hobby:artifact_creation:design_thinking_problem_framing_path"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsNoBuildPrototypeTestContext(context):
		return "maker_technical_hobby:artifact_creation:no_build_prototype_test"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "concept_mastery" && containsClimateDataGraphingExplanationContext(context):
		return "knowledge_hobby:concept_mastery:climate_data_graphing_explanation"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsCitizenScienceObservationRecordContext(context):
		return "knowledge_hobby:habit_lifestyle:citizen_science_observation_record"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "concept_mastery" && containsPrimarySourceInquiryNoteContext(context):
		return "knowledge_hobby:concept_mastery:primary_source_inquiry_note"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsSewingProjectFoundationContext(context):
		return "craft_making:artifact_creation:sewing_project_foundation"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsArduinoSensorProjectContext(context):
		return "maker_technical_hobby:artifact_creation:arduino_sensor_project"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsCardboardCircuitInventionPathContext(context):
		return "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsMusicProductionFullPipelineContext(context):
		return "digital_creation:artifact_creation:music_production_full_pipeline"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsPhotographyFundamentalsProjectContext(context):
		return "visual_art:artifact_creation:photography_fundamentals_project"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsGardenDesignMaintenancePathContext(context):
		return "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "artifact_creation" && containsCookingBasicsFoundationContext(context):
		return "cooking_baking:artifact_creation:cooking_basics_foundation"
	case pattern.DomainAxis == "writing_storytelling" && pattern.GoalModePrimary == "presentation_publish" && containsCriticalReadingToWritingDraftContext(context):
		return "writing_storytelling:presentation_publish:critical_reading_to_writing_draft"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsPhotographyStudioProjectContext(context):
		return "visual_art:artifact_creation:photography_studio_project"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsDesignLabPrototypeProjectContext(context):
		return "maker_technical_hobby:artifact_creation:design_lab_prototype_project"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsOrganicGrowingCyclePlanContext(context):
		return "knowledge_hobby:habit_lifestyle:organic_growing_cycle_plan"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsAdobeCreativeLearningPathContext(context):
		return "digital_creation:artifact_creation:adobe_creative_learning_path"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsTemplateDesignQuickOutputContext(context):
		return "digital_creation:artifact_creation:template_design_quick_output"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsSequentialCreativeClassPathContext(context):
		return "visual_art:artifact_creation:sequential_creative_class_path"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsAnimationFundamentalsShotProgressionContext(context):
		return "digital_creation:artifact_creation:animation_fundamentals_shot_progression"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsInteractiveMusicProductionFoundationContext(context):
		return "digital_creation:artifact_creation:interactive_music_production_foundation"
	case pattern.DomainAxis == "instrument_performance" && pattern.GoalModePrimary == "performance_execution" && containsSongDrivenInstrumentPathContext(context):
		return "instrument_performance:performance_execution:song_driven_instrument_path"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "artifact_creation" && containsMakerProjectClassContext(context):
		return "maker_technical_hobby:artifact_creation:maker_project_class"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsShortDiplomaAssessmentPathContext(context):
		return "knowledge_hobby:professional_transition:short_diploma_assessment_path"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsShortCourseWeeklyDiscussionContext(context):
		return "knowledge_hobby:foundation_build:short_course_weekly_discussion"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsAuditCoursePracticePathContext(context):
		return "knowledge_hobby:foundation_build:audit_course_practice_path"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "foundation_build" && containsProblemSetFinalProjectContext(context):
		return "digital_creation:foundation_build:problem_set_to_final_project"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsWebDevProjectFoundationContext(context):
		return "digital_creation:artifact_creation:web_dev_project_foundation"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "foundation_build" && containsWebPlatformCoreCurriculumContext(context):
		return "digital_creation:foundation_build:web_platform_core_curriculum"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "professional_transition" && containsFrontendCompetencyMapContext(context):
		return "digital_creation:professional_transition:frontend_competency_map"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsCreativeCodingVisualProjectContext(context):
		return "digital_creation:artifact_creation:creative_coding_visual_project"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsGame3DGuidedPathwayContext(context):
		return "digital_creation:artifact_creation:game_3d_guided_pathway"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsOpenTextbookChapterPracticeContext(context):
		return "knowledge_hobby:foundation_build:open_textbook_chapter_practice"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsOERLearningObjectClusterContext(context):
		return "knowledge_hobby:foundation_build:oer_learning_object_cluster"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsInteractiveTutorConceptPracticeContext(context):
		return "knowledge_hobby:foundation_build:interactive_tutor_concept_practice"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "professional_transition" && containsCloudLabSkillPathContext(context):
		return "digital_creation:professional_transition:cloud_lab_skill_path"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "foundation_build" && containsEnglishMOOCProgrammingContext(context):
		return "digital_creation:foundation_build:mooc_programming_foundation"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsKoreanOpenLectureWeeklySurveyContext(context):
		return "knowledge_hobby:foundation_build:korean_open_lecture_weekly_survey"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "foundation_build" && containsEnglishMOOCHumanitiesSurveyContext(context):
		return "knowledge_hobby:foundation_build:mooc_humanities_survey"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsDigitalLiteracyBadgedCourseContext(context):
		return "knowledge_hobby:habit_lifestyle:digital_literacy_badged_course"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "professional_transition" && containsRoleSkillLearningPathContext(context):
		return "digital_creation:professional_transition:learning_path_role_skill"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsStructuredDrawingFoundationContext(context):
		return "visual_art:artifact_creation:structured_drawing_foundation"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsVideoEditingShortformProjectContext(context):
		return "digital_creation:artifact_creation:video_editing_shortform_project"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsCreativeToolBeginnerOutputContext(context):
		return "digital_creation:artifact_creation:creative_tool_beginner_output"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsKMOOCOnlineCourseContext(context):
		return "knowledge_hobby:habit_lifestyle:kmooc_online_course_survey"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "habit_lifestyle" && containsUniversitySurveyContext(context):
		return "knowledge_hobby:habit_lifestyle:university_survey"
	case pattern.DomainAxis == "writing_storytelling" && pattern.GoalModePrimary == "presentation_publish" && containsUniversityStudioWritingContext(context):
		return "writing_storytelling:presentation_publish:university_studio"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsUniversityStudioVisualContext(context):
		return "visual_art:artifact_creation:university_studio"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsUniversityStudioCreativeTechContext(context):
		return "digital_creation:artifact_creation:university_studio"
	case pattern.DomainAxis == "instrument_performance" && pattern.GoalModePrimary == "participation_service" && containsAny(context, "찬양단", "워십", "교회", "예배", "봉사"):
		return "instrument_performance:participation_service:worship_team_support"
	case pattern.DomainAxis == "instrument_performance" && pattern.GoalModePrimary == "performance_execution" && containsVocalSongPracticeContext(context):
		return "instrument_performance:performance_execution:vocal_song_practice"
	case pattern.DomainAxis == "instrument_performance" && pattern.GoalModePrimary == "performance_execution" && containsInstrumentSongCompletionContext(context):
		return "instrument_performance:performance_execution:song_completion"
	case pattern.DomainAxis == "body_movement" && pattern.GoalModePrimary == "performance_execution" && containsDanceCoverRoutineContext(context):
		return "body_movement:performance_execution:dance_cover_song_routine"
	case pattern.DomainAxis == "body_movement" && pattern.GoalModePrimary == "habit_lifestyle" && containsRunning5KBeginnerRoutineContext(context):
		return "body_movement:habit_lifestyle:running_5k_beginner_routine"
	case pattern.DomainAxis == "body_movement" && pattern.GoalModePrimary == "habit_lifestyle" && containsHomeStrengthRoutineContext(context):
		return "body_movement:habit_lifestyle:home_strength_routine"
	case pattern.DomainAxis == "body_movement" && pattern.GoalModePrimary == "habit_lifestyle" && containsYogaDailyRoutineContext(context):
		return "body_movement:habit_lifestyle:yoga_daily_routine"
	case pattern.DomainAxis == "writing_storytelling" && pattern.GoalModePrimary == "presentation_publish" && containsAny(context, "브런치", "연재", "발행",
		"essay series", "weekly essay", "publish a short essay", "online publication"):
		return "writing_storytelling:presentation_publish:brunch_serial_publish"
	case pattern.DomainAxis == "writing_storytelling" && pattern.GoalModePrimary == "habit_lifestyle" && containsAny(context,
		"매일", "루틴", "습관", "일기", "기록",
		"daily", "habit", "journal", "reflection", "reflections", "write every day", "weekly review"):
		return "writing_storytelling:habit_lifestyle:daily_writing_habit"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "teaching_instruction" && containsAny(context, "유튜브", "youtube", "강의", "시연"):
		return "craft_making:teaching_instruction:craft_teach_youtube"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsCandleSoapResinProjectContext(context):
		return "craft_making:artifact_creation:candle_soap_resin_project"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsWoodworkingSmallProjectContext(context):
		return "craft_making:artifact_creation:woodworking_small_project"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsKnittingBasicWearableContext(context):
		return "craft_making:artifact_creation:knitting_basic_wearable"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsCrochetSmallDollProjectContext(context):
		return "craft_making:artifact_creation:crochet_small_doll_project"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsPotteryHandbuildingCupProjectContext(context):
		return "craft_making:artifact_creation:pottery_handbuilding_cup_project"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsLeathercraftContext(context) && containsAny(context, "지갑", "카드지갑", "반지갑", "동전지갑", "머니클립", "wallet"):
		return "craft_making:artifact_creation:leathercraft_wallet_project"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsLeathercraftContext(context) && containsAny(context, "가방", "핸드백", "미니백", "클러치", "파우치", "입체", "bag", "clutch", "pouch"):
		return "craft_making:artifact_creation:leathercraft_dimensional_bag"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "artifact_creation" && containsLeathercraftContext(context):
		return "craft_making:artifact_creation:leathercraft_basic_accessory"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "presentation_publish" && containsCalligraphyContext(context):
		return "visual_art:presentation_publish:calligraphy_publish_showcase"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsCalligraphyContext(context) && containsAny(context, "엽서", "카드", "명언", "문구", "문장", "액자", "선물", "작품", "quote", "card", "postcard", "composition"):
		return "visual_art:artifact_creation:calligraphy_quote_art_project"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "artifact_creation" && containsCalligraphyContext(context):
		return "visual_art:artifact_creation:calligraphy_basic_lettering"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "certification_assessment" && containsEnglishScoreCertificationContext(context):
		return "language_communication:certification_assessment:english_score_exam_certification"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "certification_assessment" && containsEnglishSpeakingInterviewCertificationContext(context):
		return "language_communication:certification_assessment:english_speaking_interview_certification"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "certification_assessment" && containsEnglishIntegratedFourSkillsCertificationContext(context):
		return "language_communication:certification_assessment:english_integrated_four_skills_certification"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "jlpt", "toeic", "toefl", "opic", "시험", "자격증", "취득"):
		return "language_communication:certification_assessment:exam_targeted_language"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "performance_execution" && containsDailyEnglishConversationRoutineContext(context):
		return "language_communication:performance_execution:daily_english_conversation_routine"
	case pattern.DomainAxis == "language_communication" && pattern.GoalModePrimary == "performance_execution" && containsAny(context,
		"여행", "공항", "식당", "해외", "길 묻기",
		"travel", "airport", "hotel", "restaurant", "directions", "conversation", "travel english"):
		return "language_communication:performance_execution:travel_conversation"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "presentation_publish" && containsAny(context, "전시", "인스타", "instagram", "업로드", "공개"):
		return "visual_art:presentation_publish:art_publish_showcase"
	case pattern.DomainAxis == "visual_art" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "메이크업", "미용사"):
		return "visual_art:certification_assessment:beauty_service_practical_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "presentation_publish" && containsAny(context,
		"포트폴리오", "업로드", "게시", "릴스", "유튜브", "공개",
		"portfolio", "publish", "upload", "showcase", "short video", "vertical video", "youtube", "reels",
		"podcast", "episode", "audio"):
		return "digital_creation:presentation_publish:portfolio_publish"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "teaching_instruction" && containsAny(context, "튜토리얼", "강의", "유튜브", "설명", "시연",
		"tutorial video", "how-to video", "explain my process", "beginner tutorial"):
		return "digital_creation:teaching_instruction:creator_tutorial_publish"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsMIDIProductionContext(context) && containsAny(context, "개념", "기초", "원리", "컨트롤러", "controller", "송수신", "sending", "receiving"):
		return "digital_creation:artifact_creation:midi_foundation_workflow"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsMIDIProductionContext(context) && containsAny(context, "전체", "완곡", "완성곡", "트랙", "track", "song", "곡", "export", "arrangement", "편곡", "automation", "오토메이션"):
		return "digital_creation:artifact_creation:freeware_midi_full_track"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsMIDIProductionContext(context):
		return "digital_creation:artifact_creation:freeware_midi_beatmaking"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsVibeCodingContext(context) && containsAny(context, "claude code", "클로드 코드", "mcp", "hook", "hooks", "agent", "multi-agent", "멀티 에이전트", "agentic", "자동화"):
		return "digital_creation:artifact_creation:claude_code_agentic_workflow"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsVibeCodingContext(context) && containsAny(context, "full-stack", "fullstack", "풀스택", "supabase", "vercel", "auth", "인증", "database", "db", "배포", "deploy", "shipping", "보안", "security"):
		return "digital_creation:artifact_creation:vibe_coding_fullstack_ship"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "artifact_creation" && containsVibeCodingContext(context):
		return "digital_creation:artifact_creation:vibe_coding_mvp_app"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "aws", "ncp", "naver cloud", "google cloud", "gcp", "azure", "az-900", "az 900", "cloud practitioner", "digital leader"):
		return "digital_creation:certification_assessment:cloud_foundation_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "solutions architect", "associate cloud engineer", "az-104", "az 104", "professional"):
		return "digital_creation:certification_assessment:cloud_operator_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "정보처리기사"):
		return "digital_creation:certification_assessment:software_engineering_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context,
		"컴퓨터활용능력", "컴활", "office spreadsheet certification", "spreadsheet certification", "excel certification", "office tool certification"):
		return "digital_creation:certification_assessment:office_tool_practical_certification"
	case pattern.DomainAxis == "digital_creation" && pattern.GoalModePrimary == "certification_assessment" && containsKoreanDesignToolPracticalCertificationContext(context):
		return "digital_creation:certification_assessment:korean_design_tool_practical_certification"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "teaching_instruction" && containsAny(context, "워크숍", "튜토리얼", "강의", "시연", "설명",
		"workshop", "demo", "demonstrate", "explain it", "sensor demo"):
		return "maker_technical_hobby:teaching_instruction:maker_workshop_demo"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "드론", "4종", "4 종", "초경량비행장치",
		"drone", "drone operator", "drone license", "basic drone"):
		return "maker_technical_hobby:certification_assessment:drone_operator_basic_4class"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "드론", "1종", "1 종", "2종", "2 종", "3종", "3 종"):
		return "maker_technical_hobby:certification_assessment:drone_operator_practical_1to3class"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "지게차", "용접", "피복아크", "자동차정비", "정비기능사"):
		return "maker_technical_hobby:certification_assessment:equipment_operation_practical_certification"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "건축기사", "건축산업기사", "토목기사"):
		return "maker_technical_hobby:certification_assessment:construction_engineering_certification"
	case pattern.DomainAxis == "maker_technical_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "전기기사", "전기산업기사", "전기기능사", "전기공사산업기사"):
		return "maker_technical_hobby:certification_assessment:electrical_engineering_certification"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "teaching_instruction" && containsAny(context,
		"클래스", "강의", "레시피", "시연", "설명", "class", "recipe demo", "recipe demonstration", "demonstrate a recipe", "baking class", "explain a recipe"):
		return "cooking_baking:teaching_instruction:baking_class_demo"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "artifact_creation" && containsKoreanHomeCookingBasicsContext(context):
		return "cooking_baking:artifact_creation:korean_home_cooking_basics"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "artifact_creation" && containsHomeCafeLatteFoundationContext(context):
		return "cooking_baking:artifact_creation:home_cafe_latte_foundation"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "제과기능사", "제빵기능사"):
		return "cooking_baking:certification_assessment:baking_written_practical_certification"
	case pattern.DomainAxis == "cooking_baking" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context, "한식조리기능사", "조리기능사"):
		return "cooking_baking:certification_assessment:korean_cooking_practical_certification"
	case pattern.DomainAxis == "craft_making" && pattern.GoalModePrimary == "certification_assessment" && containsBeautyServicePracticalContext(context):
		return "craft_making:certification_assessment:beauty_service_practical_certification"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context,
		"직업상담사", "counseling certification", "counselor certification", "counselling certification", "career counselor exam", "counseling exam"):
		return "knowledge_hobby:certification_assessment:counseling_certification"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "certification_assessment" && containsAny(context,
		"대기환경기사", "environmental engineering certification", "environmental engineer exam", "environmental engineering exam", "air pollution engineer"):
		return "knowledge_hobby:certification_assessment:environmental_engineering_certification"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsCreditBankStandardTheoryContext(context):
		return "knowledge_hobby:professional_transition:credit_bank_standard_theory"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "teaching_instruction" && containsCreditBankInstructionalDesignContext(context):
		return "knowledge_hobby:teaching_instruction:credit_bank_instructional_design"
	case pattern.DomainAxis == "knowledge_hobby" && pattern.GoalModePrimary == "professional_transition" && containsCreditBankPracticumContext(context):
		return "knowledge_hobby:professional_transition:credit_bank_practicum"
	default:
		return ""
	}
}

func inferGoalModePrimary(context, domain string) string {
	switch domain {
	case "instrument_performance":
		if containsAny(context, "찬양단", "워십", "교회", "예배", "봉사") {
			return "participation_service"
		}
		if containsAny(context, "합주", "밴드") &&
			!containsInstrumentSongCompletionContext(context) &&
			!containsAny(context, "리듬", "루트음", "곡", "연주") {
			return "participation_service"
		}
		return "performance_execution"
	case "language_communication":
		if containsLanguageCertificationContext(context) {
			return "certification_assessment"
		}
		return "performance_execution"
	case "craft_making":
		if containsPracticalCertificationContext(context) {
			return "certification_assessment"
		}
		if containsAny(context, "강의", "유튜브", "가르치", "튜토리얼", "시연") {
			return "teaching_instruction"
		}
		if containsAny(context, "루틴", "웰빙", "힐링", "일상") {
			return "habit_lifestyle"
		}
		return "artifact_creation"
	case "visual_art":
		if containsPracticalCertificationContext(context) {
			return "certification_assessment"
		}
		if containsPhotographyFundamentalsProjectContext(context) {
			return "artifact_creation"
		}
		if containsAny(context, "판매", "etsy", "수익", "클래스 오픈", "포트폴리오 판매") {
			return "professional_transition"
		}
		if containsAny(context, "전시", "업로드", "게시", "공개", "발행", "인스타", "instagram", "올려", "공유해", "공유하고") {
			return "presentation_publish"
		}
		if containsAny(context, "루틴", "매일", "일상", "습관") {
			return "habit_lifestyle"
		}
		return "artifact_creation"
	case "writing_storytelling":
		if containsUniversityStudioWritingContext(context) {
			return "presentation_publish"
		}
		if containsAny(context, "강의", "워크숍", "가르치") {
			return "teaching_instruction"
		}
		if containsAny(context, "내 생각 정리", "생각을 정리", "마음 정리", "기록 정리", "일기", "기록", "정리하고 싶",
			"daily writing", "writing habit", "journal habit", "short reflections", "write short reflections", "write every day") {
			return "habit_lifestyle"
		}
		if containsAny(context, "브런치", "연재", "출간", "공개", "발행",
			"essay series", "weekly essay", "publish a short essay", "online publication") {
			return "presentation_publish"
		}
		if containsAny(context, "루틴", "습관", "매일", "꾸준히", "routine", "habit", "daily", "weekly review") {
			return "habit_lifestyle"
		}
		return "presentation_publish"
	case "body_movement":
		if containsAny(context, "참여", "공연", "소셜댄스", "합류", "팀") {
			return "participation_service"
		}
		if containsDanceCoverRoutineContext(context) {
			return "performance_execution"
		}
		return "habit_lifestyle"
	case "digital_creation":
		if containsAIToolProductivityWorkflowContext(context) {
			return "technical_skill"
		}
		if containsKoreanContentCreationPipelineContext(context) {
			return "artifact_creation"
		}
		if containsCreatorTutorialPublishContext(context) || containsAny(context, "강의 영상", "튜토리얼 영상", "시연 영상") {
			return "teaching_instruction"
		}
		if containsKoreanNCSWebPublisherTrainingContext(context) {
			return "professional_transition"
		}
		if containsAny(context, "수익", "판매", "프리랜", "의뢰", "클라이언트", "side hustle", "부업", "부수입", "사이드잡", "돈 벌") {
			return "professional_transition"
		}
		if containsComputationalMusicAnalysisProjectContext(context) {
			return "technical_skill"
		}
		if containsVideoEditingShortformProjectContext(context) {
			if containsAny(context, "업로드", "게시", "공개", "발행", "유튜브", "인스타", "instagram", "publish", "upload", "portfolio", "showcase", "release") {
				return "presentation_publish"
			}
			return "artifact_creation"
		}
		if containsFrontendCompetencyMapContext(context) || containsCloudLabSkillPathContext(context) {
			return "professional_transition"
		}
		if containsWebDevProjectFoundationContext(context) || containsCreativeCodingVisualProjectContext(context) || containsGame3DGuidedPathwayContext(context) ||
			containsAdobeCreativeLearningPathContext(context) || containsTemplateDesignQuickOutputContext(context) ||
			containsNotionStudyDashboardContext(context) ||
			containsCreativeToolBeginnerOutputContext(context) ||
			containsAnimationFundamentalsShotProgressionContext(context) || containsInteractiveMusicProductionFoundationContext(context) ||
			containsMusicProductionFullPipelineContext(context) {
			return "artifact_creation"
		}
		if containsRoleSkillLearningPathContext(context) && !containsAny(context, "자격증", "자격", "취득", "시험", "certification", "certified", "exam", "associate", "practitioner") {
			return "professional_transition"
		}
		if containsAny(context, "자격증", "자격", "취득", "시험", "certification", "associate", "practitioner", "기사", "기능사", "산업기사") ||
			containsCloudCertificationContext(context) || containsKoreanDesignToolPracticalCertificationContext(context) || containsAny(context, "정보처리기사", "컴퓨터활용능력", "컴활") {
			return "certification_assessment"
		}
		if containsProblemSetFinalProjectContext(context) || containsWebPlatformCoreCurriculumContext(context) || containsEnglishMOOCProgrammingContext(context) {
			return "foundation_build"
		}
		if containsRoleSkillLearningPathContext(context) || containsAny(context, "수익", "판매", "프리랜", "의뢰", "클라이언트", "side hustle", "부업", "부수입", "사이드잡", "돈 벌") {
			return "professional_transition"
		}
		if containsAny(context, "업로드", "게시", "연재", "공개", "발행", "유튜브",
			"publish", "upload", "portfolio", "showcase", "release", "youtube", "short video") {
			return "presentation_publish"
		}
		if containsAny(context, "강의", "튜토리얼", "가르치", "시연", "설명") {
			return "teaching_instruction"
		}
		return "artifact_creation"
	case "knowledge_hobby":
		if containsBlendedLifelongParticipationContext(context) {
			return "participation_service"
		}
		if containsMidlifeTransitionLearningSupportContext(context) || containsCompletionRefundOnlineCourseContext(context) ||
			containsRuralReturnAgricultureFoundationContext(context) {
			return "professional_transition"
		}
		if containsFederatedLifelongOpenResourceContext(context) {
			return "foundation_build"
		}
		if containsPublicLifelongOnlineCourseContext(context) || containsUrbanAgriculturePracticeSeriesContext(context) ||
			containsAIDigitalLiteracyLifePracticeContext(context) {
			return "habit_lifestyle"
		}
		if containsScienceSimulationInquiryPathContext(context) || containsGuidedSimulationActivitySheetContext(context) ||
			containsHorticultureDiagnosticFoundationContext(context) || containsClimateDataGraphingExplanationContext(context) ||
			containsPrimarySourceInquiryNoteContext(context) {
			return "concept_mastery"
		}
		if containsCitizenScienceObservationRecordContext(context) {
			return "habit_lifestyle"
		}
		if containsHouseplantRepottingCareContext(context) || containsPersonalFinanceBudgetRoutineContext(context) ||
			containsDogBasicTrainingRoutineContext(context) {
			return "habit_lifestyle"
		}
		if containsGardenDesignMaintenancePathContext(context) || containsOrganicGrowingCyclePlanContext(context) {
			return "habit_lifestyle"
		}
		if containsShortDiplomaAssessmentPathContext(context) {
			return "professional_transition"
		}
		if containsShortCourseWeeklyDiscussionContext(context) || containsAuditCoursePracticePathContext(context) {
			return "foundation_build"
		}
		if containsOpenTextbookChapterPracticeContext(context) || containsOERLearningObjectClusterContext(context) || containsInteractiveTutorConceptPracticeContext(context) {
			return "foundation_build"
		}
		if containsDigitalLiteracyBadgedCourseContext(context) {
			return "habit_lifestyle"
		}
		if containsKoreanOpenLectureWeeklySurveyContext(context) {
			return "foundation_build"
		}
		if containsEnglishMOOCHumanitiesSurveyContext(context) {
			return "foundation_build"
		}
		if containsCreditBankInstructionalDesignContext(context) {
			return "teaching_instruction"
		}
		if containsCreditBankPracticumContext(context) {
			return "professional_transition"
		}
		if containsCreditBankStandardTheoryContext(context) {
			return "professional_transition"
		}
		if containsTheoryCertificationContext(context) {
			return "certification_assessment"
		}
		return "habit_lifestyle"
	case "maker_technical_hobby":
		if containsTheoryCertificationContext(context) || containsPracticalCertificationContext(context) {
			return "certification_assessment"
		}
		if containsCardboardCircuitInventionPathContext(context) || containsDesignLabPrototypeProjectContext(context) ||
			containsDesignThinkingProblemFramingPathContext(context) || containsNoBuildPrototypeTestContext(context) {
			return "artifact_creation"
		}
		if containsAny(context, "강의", "튜토리얼", "워크숍", "가르치", "시연", "설명",
			"workshop", "demo", "demonstrate", "explain it", "sensor demo") {
			return "teaching_instruction"
		}
		if containsAny(context, "수익", "판매", "의뢰", "상품화", "사이드 허슬", "side hustle") {
			return "professional_transition"
		}
		return "artifact_creation"
	case "cooking_baking":
		if containsPracticalCertificationContext(context) {
			return "certification_assessment"
		}
		if containsMealPrepWeekdayLunchContext(context) {
			return "habit_lifestyle"
		}
		if containsFoodPreservationSafetyPathContext(context) {
			return "habit_lifestyle"
		}
		if containsAny(context, "강의", "가르치", "클래스", "시연",
			"class", "recipe demo", "recipe demonstration", "demonstrate a recipe", "baking class", "explain a recipe") {
			return "teaching_instruction"
		}
		if containsCookingBasicsFoundationContext(context) {
			return "artifact_creation"
		}
		if containsAny(context, "수익", "판매", "주문", "브랜드", "창업") {
			return "professional_transition"
		}
		if containsKoreanHomeCookingBasicsContext(context) {
			return "artifact_creation"
		}
		if containsAny(context, "습관", "일상", "건강식", "루틴", "건강한 식사", "식단", "집밥", "가족에게", "해주고 싶") {
			return "habit_lifestyle"
		}
		return "artifact_creation"
	default:
		return ""
	}
}

func inferGoalModeSecondary(context, domain, primary string) string {
	switch domain {
	case "instrument_performance":
		if primary == "participation_service" {
			return "performance_execution"
		}
	case "language_communication":
		return "foundation_build"
	case "craft_making":
		if containsLocalHobbyWorkshopSeriesContext(context) {
			return "participation_service"
		}
		if primary == "teaching_instruction" {
			return "artifact_creation"
		}
		return "foundation_build"
	case "visual_art":
		if primary == "professional_transition" {
			return "artifact_creation"
		}
		if primary == "habit_lifestyle" {
			return "foundation_build"
		}
		return "artifact_creation"
	case "writing_storytelling":
		if primary == "teaching_instruction" {
			return "presentation_publish"
		}
		if primary == "presentation_publish" {
			return "artifact_creation"
		}
		if primary == "habit_lifestyle" {
			return "foundation_build"
		}
	case "body_movement", "knowledge_hobby", "maker_technical_hobby", "cooking_baking", "digital_creation":
		if domain == "knowledge_hobby" && containsBlendedLifelongParticipationContext(context) {
			return "habit_lifestyle"
		}
		if domain == "knowledge_hobby" && containsMidlifeTransitionLearningSupportContext(context) {
			return "foundation_build"
		}
		if domain == "knowledge_hobby" && containsCompletionRefundOnlineCourseContext(context) {
			return "habit_lifestyle"
		}
		if domain == "knowledge_hobby" && containsFederatedLifelongOpenResourceContext(context) {
			return "habit_lifestyle"
		}
		if domain == "knowledge_hobby" && containsKoreanOpenLectureWeeklySurveyContext(context) {
			return "habit_lifestyle"
		}
		if domain == "digital_creation" && containsKoreanDesignToolPracticalCertificationContext(context) {
			return "artifact_creation"
		}
		if domain == "knowledge_hobby" && containsRuralReturnAgricultureFoundationContext(context) {
			return "habit_lifestyle"
		}
		if domain == "knowledge_hobby" && containsUrbanAgriculturePracticeSeriesContext(context) {
			return "artifact_creation"
		}
		if domain == "knowledge_hobby" && containsAIDigitalLiteracyLifePracticeContext(context) {
			return "digital_literacy"
		}
		if domain == "craft_making" && containsLocalHobbyWorkshopSeriesContext(context) {
			return "participation_service"
		}
		if domain == "knowledge_hobby" && containsCreditBankInstructionalDesignContext(context) {
			return "professional_transition"
		}
		if domain == "knowledge_hobby" && containsCreditBankPracticumContext(context) {
			return "participation_service"
		}
		if domain == "knowledge_hobby" && containsCreditBankStandardTheoryContext(context) {
			return "foundation_build"
		}
		if primary == "teaching_instruction" {
			return "artifact_creation"
		}
		if primary == "professional_transition" {
			return "artifact_creation"
		}
		if primary == "presentation_publish" {
			return "professional_transition"
		}
		if primary == "foundation_build" {
			return "artifact_creation"
		}
		return "foundation_build"
	}
	return ""
}

func analyzeGoalRouting(req CreateCourseDraftRequest) goalRoutingSummary {
	patternKey := inferGoalPatternKey(req)
	subpatternKey := inferGoalSubpatternKey(req)

	switch {
	case subpatternKey != "":
		if refinement, ok := lookupGoalSubpatternRefinement(subpatternKey); ok && refinement.Refine != nil {
			return goalRoutingSummary{
				PatternKey:     patternKey,
				SubpatternKey:  subpatternKey,
				RefinementMode: "subpattern_strong",
			}
		}
		if _, ok := goalSubpatternSpecs[subpatternKey]; ok {
			return goalRoutingSummary{
				PatternKey:     patternKey,
				SubpatternKey:  subpatternKey,
				RefinementMode: "subpattern_guidance_only",
			}
		}
	case patternKey.DomainAxis == "" && patternKey.GoalModePrimary == "":
		return goalRoutingSummary{RefinementMode: "unclassified"}
	}

	if refinement, ok := lookupGoalPatternRefinement(patternKey); ok && refinement.Refine != nil {
		return goalRoutingSummary{
			PatternKey:     patternKey,
			SubpatternKey:  subpatternKey,
			RefinementMode: "pattern_strong",
		}
	}
	if refinement, ok := lookupGoalPatternSoftRefinement(patternKey); ok && refinement.Apply != nil {
		return goalRoutingSummary{
			PatternKey:     patternKey,
			SubpatternKey:  subpatternKey,
			RefinementMode: "pattern_soft",
		}
	}
	if _, ok := lookupGoalPatternSpec(patternKey); ok {
		return goalRoutingSummary{
			PatternKey:     patternKey,
			SubpatternKey:  subpatternKey,
			RefinementMode: "pattern_guidance_only",
		}
	}

	return goalRoutingSummary{
		PatternKey:     patternKey,
		SubpatternKey:  subpatternKey,
		RefinementMode: "unclassified",
	}
}

func formatGoalPatternKeyForLog(key goalPatternKey) string {
	parts := make([]string, 0, 3)
	if strings.TrimSpace(key.DomainAxis) != "" {
		parts = append(parts, strings.TrimSpace(key.DomainAxis))
	}
	if strings.TrimSpace(key.GoalModePrimary) != "" {
		parts = append(parts, strings.TrimSpace(key.GoalModePrimary))
	}
	if strings.TrimSpace(key.GoalModeSecondary) != "" {
		parts = append(parts, strings.TrimSpace(key.GoalModeSecondary))
	}
	if len(parts) == 0 {
		return "unclassified"
	}
	return strings.Join(parts, ":")
}
