package curriculum

import (
	"strings"
)

type RecommendationGoalMetadata struct {
	PatternKey    string
	SubpatternKey string
	QueryHint     string
	QueryTerms    []string
}

func BuildRecommendationGoalMetadata(req CreateCourseDraftRequest) RecommendationGoalMetadata {
	patternKey := inferGoalPatternKey(req)
	subpatternKey := inferGoalSubpatternKey(req)
	terms := buildRecommendationGoalTerms(patternKey, subpatternKey, buildGoalInferenceContext(req))
	return RecommendationGoalMetadata{
		PatternKey:    formatGoalPatternKeyForLog(patternKey),
		SubpatternKey: subpatternKey,
		QueryHint:     strings.Join(terms, " "),
		QueryTerms:    terms,
	}
}

func buildRecommendationGoalTerms(patternKey goalPatternKey, subpatternKey, context string) []string {
	var terms []string

	switch subpatternKey {
	case "knowledge_hobby:habit_lifestyle:public_lifelong_online_course":
		terms = []string{"평생학습", "온라인 강좌", "생활 실천", "핵심 개념", "수료 정리"}
	case "knowledge_hobby:participation_service:blended_lifelong_participation":
		terms = []string{"평생학습", "오프라인 강좌", "화상학습", "참여 준비", "학습이력"}
	case "knowledge_hobby:professional_transition:midlife_transition_learning_support":
		terms = []string{"중장년", "직업역량", "직업 전환", "취창업", "학습계획"}
	case "knowledge_hobby:professional_transition:completion_refund_online_course":
		terms = []string{"온라인 과정", "진도 관리", "수료", "환급", "적용 정리"}
	case "knowledge_hobby:foundation_build:federated_lifelong_open_resource":
		terms = []string{"늘배움", "평생학습", "공개강좌", "학습경로", "자기주도 학습"}
	case "digital_creation:artifact_creation:korean_content_creation_pipeline":
		terms = []string{"콘텐츠 제작", "방송영상", "게임", "애니메이션", "포트폴리오"}
	case "knowledge_hobby:professional_transition:rural_return_agriculture_foundation":
		terms = []string{"귀농귀촌", "농업 기초", "품목기술", "정착 준비", "교육 수료"}
	case "knowledge_hobby:habit_lifestyle:urban_agriculture_practice_series":
		terms = []string{"도시농업", "텃밭", "가정원예", "재배 실습", "관리 루틴"}
	case "knowledge_hobby:habit_lifestyle:ai_digital_literacy_life_practice":
		terms = []string{"디지털배움터", "스마트폰", "키오스크", "AI 활용", "보이스피싱 예방"}
	case "craft_making:artifact_creation:local_hobby_workshop_series":
		terms = []string{"생활취미", "평생학습센터", "워크숍", "첫 결과물", "피드백 반영"}
	case "knowledge_hobby:concept_mastery:science_simulation_inquiry_path":
		terms = []string{"phet simulation", "science simulation", "variable experiment", "inquiry activity", "concept explanation"}
	case "knowledge_hobby:concept_mastery:guided_simulation_activity_sheet":
		terms = []string{"phet activity sheet", "guided inquiry", "simulation worksheet", "small group activity", "reflection"}
	case "cooking_baking:habit_lifestyle:food_preservation_safety_path":
		terms = []string{"food preservation", "water bath canning", "high acid foods", "canning safety", "research tested recipe"}
	case "knowledge_hobby:concept_mastery:horticulture_diagnostic_foundation":
		terms = []string{"horticulture", "plant diagnostics", "integrated pest management", "plant pathology", "care plan"}
	case "digital_creation:technical_skill:computational_music_analysis_project":
		terms = []string{"computational music theory", "music21", "symbolic score", "corpus analysis", "algorithmic composition"}
	case "maker_technical_hobby:artifact_creation:design_thinking_problem_framing_path":
		terms = []string{"design thinking", "human centered design", "empathy", "prototype", "test reflection"}
	case "maker_technical_hobby:artifact_creation:no_build_prototype_test":
		terms = []string{"no-build prototype", "scrappy prototype", "prototype test", "user feedback", "design experiment"}
	case "knowledge_hobby:concept_mastery:climate_data_graphing_explanation":
		terms = []string{"climate data", "temperature anomaly", "global temperature trends", "graphing activity", "data explanation"}
	case "knowledge_hobby:habit_lifestyle:citizen_science_observation_record":
		terms = []string{"iNaturalist", "citizen science", "BioBlitz", "species identification", "biodiversity record"}
	case "knowledge_hobby:concept_mastery:primary_source_inquiry_note":
		terms = []string{"primary source analysis", "observe reflect question", "historical inquiry", "source analysis", "investigation note"}
	case "craft_making:artifact_creation:sewing_project_foundation":
		terms = []string{"sewing project", "sewing machine", "stitch", "pattern envelope", "hemming"}
	case "maker_technical_hobby:artifact_creation:cardboard_circuit_invention_path":
		terms = []string{"cardboard circuits", "electronics", "robotics", "prototype", "circuit playground"}
	case "digital_creation:artifact_creation:music_production_full_pipeline":
		terms = []string{"music production", "sound design", "arrangement", "mixing mastering", "daw"}
	case "visual_art:artifact_creation:photography_fundamentals_project":
		terms = []string{"photography fundamentals", "exposure", "composition", "lighting", "photo project"}
	case "knowledge_hobby:habit_lifestyle:garden_design_maintenance_path":
		terms = []string{"garden design", "gardening", "planting", "maintenance routine", "garden plan"}
	case "knowledge_hobby:habit_lifestyle:houseplant_repotting_care":
		terms = []string{"houseplant repotting", "soil mix", "plant recovery", "watering after repotting", "indoor plant care"}
	case "knowledge_hobby:habit_lifestyle:personal_finance_budget_routine":
		terms = []string{"personal finance", "monthly budget", "spending categories", "track spending", "expense review"}
	case "knowledge_hobby:habit_lifestyle:dog_basic_training_routine":
		terms = []string{"dog basic training", "sit stay come", "positive reinforcement", "recall training", "short training session"}
	case "cooking_baking:artifact_creation:cooking_basics_foundation":
		terms = []string{"cooking basics", "food hygiene", "knife skills", "sauce", "herbs spices"}
	case "cooking_baking:habit_lifestyle:meal_prep_weekday_lunch":
		terms = []string{"meal prep", "healthy lunch", "batch cooking", "meal storage", "weekly menu"}
	case "writing_storytelling:presentation_publish:critical_reading_to_writing_draft":
		terms = []string{"creative writing", "critical reading", "writing draft", "fiction", "poetry"}
	case "visual_art:artifact_creation:photography_studio_project":
		terms = []string{"photography studio", "digital imaging", "studio lighting", "critique", "exhibition"}
	case "maker_technical_hobby:artifact_creation:design_lab_prototype_project":
		terms = []string{"design lab", "prototype", "fabrication", "design process", "community partner"}
	case "knowledge_hobby:habit_lifestyle:organic_growing_cycle_plan":
		terms = []string{"organic gardening", "crop seasons", "sowing", "growing cycle", "care routine"}
	case "digital_creation:artifact_creation:adobe_creative_learning_path":
		terms = []string{"creative tool", "adobe", "photoshop", "illustrator", "portfolio"}
	case "digital_creation:artifact_creation:template_design_quick_output":
		terms = []string{"canva", "template design", "social media design", "brand kit", "quick output"}
	case "visual_art:artifact_creation:sequential_creative_class_path":
		terms = []string{"creative class", "illustration", "portfolio", "project", "gallery"}
	case "digital_creation:artifact_creation:animation_fundamentals_shot_progression":
		terms = []string{"blender animation", "animation fundamentals", "bouncing ball", "walk cycle", "character shot"}
	case "digital_creation:artifact_creation:interactive_music_production_foundation":
		terms = []string{"music production", "beat making", "song structure", "loop", "ableton"}
	case "instrument_performance:performance_execution:song_driven_instrument_path":
		terms = []string{"guitar lessons", "riff", "chords", "song practice", "instrument path"}
	case "maker_technical_hobby:artifact_creation:maker_project_class":
		terms = []string{"maker project", "diy", "electronics", "woodworking", "lesson plan"}
	case "knowledge_hobby:professional_transition:short_diploma_assessment_path":
		terms = []string{"diploma", "assessment", "career skills", "certificate", "self paced"}
	case "knowledge_hobby:foundation_build:short_course_weekly_discussion":
		terms = []string{"short course", "discussion", "case study", "weekly course", "professional development"}
	case "knowledge_hobby:foundation_build:audit_course_practice_path":
		terms = []string{"audit course", "lecture", "readings", "ungraded practice", "discussion"}
	case "digital_creation:foundation_build:problem_set_to_final_project":
		terms = []string{"computer science", "problem set", "final project", "programming", "cs50"}
	case "digital_creation:artifact_creation:web_dev_project_foundation":
		terms = []string{"web development", "html", "css", "javascript", "project"}
	case "digital_creation:artifact_creation:notion_study_dashboard":
		terms = []string{"Notion study dashboard", "tasks database", "notes page", "weekly review", "study planner"}
	case "digital_creation:foundation_build:web_platform_core_curriculum":
		terms = []string{"web platform", "html", "css", "javascript", "accessibility"}
	case "digital_creation:professional_transition:frontend_competency_map":
		terms = []string{"front-end developer", "competency", "accessibility", "responsive design", "job ready"}
	case "digital_creation:artifact_creation:creative_coding_visual_project":
		terms = []string{"creative coding", "javascript drawing", "animation", "interaction", "project"}
	case "digital_creation:artifact_creation:game_3d_guided_pathway":
		terms = []string{"unity", "game prototype", "real-time 3d", "scripting", "portfolio"}
	case "digital_creation:foundation_build:mooc_programming_foundation":
		terms = []string{"programming", "computer science", "problem set", "python", "project"}
	case "knowledge_hobby:foundation_build:mooc_humanities_survey":
		terms = []string{"humanities", "survey course", "readings", "discussion", "essay"}
	case "knowledge_hobby:habit_lifestyle:digital_literacy_badged_course":
		terms = []string{"digital literacy", "online safety", "digital skills", "information overload", "practice plan"}
	case "digital_creation:professional_transition:learning_path_role_skill":
		terms = []string{"learning path", "role based", "training module", "task scenario", "career path"}
	case "digital_creation:professional_transition:cloud_lab_skill_path":
		terms = []string{"cloud lab", "hands-on lab", "skill badge", "google cloud", "task scenario"}
	case "visual_art:artifact_creation:structured_drawing_foundation":
		terms = []string{"drawing foundation", "contour drawing", "shading", "perspective", "sketch"}
	case "digital_creation:artifact_creation:creative_tool_beginner_output":
		terms = []string{"creative tool", "beginner tutorial", "photoshop", "illustrator", "portfolio"}
	case "knowledge_hobby:foundation_build:open_textbook_chapter_practice":
		terms = []string{"open textbook", "chapter practice", "study guide", "concepts", "exercises"}
	case "knowledge_hobby:foundation_build:oer_learning_object_cluster":
		terms = []string{"oer", "learning objects", "simulation", "worksheet", "visualization"}
	case "knowledge_hobby:foundation_build:interactive_tutor_concept_practice":
		terms = []string{"interactive tutor", "concept practice", "biology", "feedback", "quiz"}
	case "language_communication:certification_assessment:english_score_exam_certification":
		terms = []string{"toeic", "teps", "문제풀이", "청해", "독해", "문법", "어휘", "시험"}
	case "language_communication:certification_assessment:english_speaking_interview_certification":
		terms = []string{"opic", "인터뷰", "말하기", "답변", "돌발질문", "주제"}
	case "language_communication:certification_assessment:english_integrated_four_skills_certification":
		terms = []string{"ielts", "toefl", "reading", "listening", "writing", "speaking", "task"}
	case "language_communication:certification_assessment:exam_targeted_language":
		terms = []string{"시험", "문제풀이", "청해", "독해", "문법"}
	case "digital_creation:certification_assessment:cloud_foundation_certification":
		terms = []string{"클라우드", "서비스", "핵심개념", "입문", "시험"}
	case "digital_creation:certification_assessment:cloud_operator_certification":
		terms = []string{"클라우드", "운영", "배포", "보안", "시나리오", "시험"}
	case "digital_creation:certification_assessment:software_engineering_certification":
		terms = []string{"정보처리기사", "기출", "문제풀이", "소프트웨어", "기사"}
	case "digital_creation:certification_assessment:office_tool_practical_certification":
		terms = []string{"컴활", "엑셀", "실기", "함수", "작업형"}
	case "maker_technical_hobby:certification_assessment:drone_operator_basic_4class":
		terms = []string{"드론", "4종", "법규", "안전", "이론"}
	case "maker_technical_hobby:certification_assessment:drone_operator_practical_1to3class":
		terms = []string{"드론", "실기", "비행", "조작", "절차"}
	case "maker_technical_hobby:certification_assessment:equipment_operation_practical_certification":
		terms = []string{"실기", "작업형", "안전", "조작", "장비"}
	case "maker_technical_hobby:certification_assessment:construction_engineering_certification":
		terms = []string{"기사", "구조", "시공", "법규", "문제풀이"}
	case "maker_technical_hobby:certification_assessment:electrical_engineering_certification":
		terms = []string{"전기", "계산", "회로", "설비", "문제풀이"}
	case "cooking_baking:certification_assessment:baking_written_practical_certification":
		terms = []string{"제과", "제빵", "실기", "반죽", "공정"}
	case "cooking_baking:certification_assessment:korean_cooking_practical_certification":
		terms = []string{"한식", "조리", "실기", "위생", "조리순서"}
	case "craft_making:certification_assessment:beauty_service_practical_certification",
		"visual_art:certification_assessment:beauty_service_practical_certification":
		terms = []string{"미용", "실기", "위생", "작업", "공개과제"}
	case "knowledge_hobby:certification_assessment:counseling_certification":
		terms = []string{"직업상담사", "이론", "기출", "문제풀이"}
	case "knowledge_hobby:certification_assessment:environmental_engineering_certification":
		terms = []string{"대기환경기사", "환경", "방지기술", "법규", "문제풀이"}
	}

	if len(terms) == 0 {
		switch {
		case patternKey.DomainAxis == "language_communication" && patternKey.GoalModePrimary == "certification_assessment":
			terms = []string{"시험", "문제풀이", "청해", "독해"}
		case patternKey.GoalModePrimary == "certification_assessment":
			terms = []string{"시험", "기출", "문제풀이"}
		}
	}

	if containsAny(context, "토플", "toefl") {
		terms = append(terms, "toefl")
	}
	if containsAny(context, "아이엘츠", "ielts") {
		terms = append(terms, "ielts")
	}
	if containsAny(context, "토익", "toeic") {
		terms = append(terms, "toeic")
	}
	if containsAny(context, "오픽", "opic") {
		terms = append(terms, "opic")
	}
	if containsAny(context, "텝스", "teps") {
		terms = append(terms, "teps")
	}

	return dedupeRecommendationTerms(terms)
}

func dedupeRecommendationTerms(terms []string) []string {
	result := make([]string, 0, len(terms))
	seen := make(map[string]struct{}, len(terms))
	for _, term := range terms {
		trimmed := strings.TrimSpace(term)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
