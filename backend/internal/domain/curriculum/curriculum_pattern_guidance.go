package curriculum

import (
	"fmt"
	"strings"
)

func buildGoalTypeGuidance(sourceQuery, learningGoal string) string {
	lowerContext := strings.ToLower(strings.TrimSpace(sourceQuery + " " + learningGoal))
	switch {
	case containsEnglishScoreCertificationContext(lowerContext):
		return "- TOEIC/TEPS 같은 영어 점수형 시험은 회화형이 아니라 점수형 시험 progression으로 설계하세요.\n" +
			"- 최소한 시험 구조 파악, 어휘·문법·독해·청해 핵심 정리, 파트별 문제풀이, 시간 대응, 시험 직전 점검 단계를 포함하세요.\n" +
			"- 실제 시험 응시와 목표 점수 달성은 completion_criteria로 분리하세요.\n"
	case containsEnglishSpeakingInterviewCertificationContext(lowerContext):
		return "- OPIc 같은 영어 말하기 인터뷰형 시험은 문법 설명보다 인터뷰 답변 흐름과 즉흥 응답에 초점을 두세요.\n" +
			"- 최소한 인터뷰 구조 이해, 자기소개와 기본 주제 답변, 주제 확장과 경험 말하기, 돌발 질문 대응, 시험 직전 실전 인터뷰 점검 단계를 포함하세요.\n" +
			"- 실제 시험 응시와 목표 등급 획득은 completion_criteria로 분리하세요.\n"
	case containsEnglishIntegratedFourSkillsCertificationContext(lowerContext):
		return "- IELTS/TOEFL 같은 4영역 통합형 영어 시험은 읽기·듣기·쓰기·말하기를 분리하되 마지막에는 task 전환과 시간 운영이 함께 보이도록 설계하세요.\n" +
			"- 최소한 시험 유형 이해, Reading/Listening 기반 정리, Writing/Speaking task 적용, 4영역 통합 연습, 시험 직전 시간 배분 점검 단계를 포함하세요.\n" +
			"- 실제 시험 응시와 목표 점수 달성은 completion_criteria로 분리하세요.\n"
	case containsLanguageCertificationContext(lowerContext):
		return "- 언어 자격증 목표는 회화형이 아니라 시험형 progression으로 설계하세요.\n" +
			"- 최소한 시험 범위 파악, 핵심 어휘·문법·독해 기반 정리, 청해/응답 흐름 적응, 유형별 문제풀이, 시험 직전 점검 단계를 포함하세요.\n" +
			"- JLPT처럼 영역이 분명한 시험은 문자·어휘, 문법·독해, 청해 흐름이 드러나야 합니다.\n" +
			"- 실제 시험 응시와 합격은 completion_criteria로 분리하세요.\n"
	case containsCloudCertificationContext(lowerContext):
		return "- 클라우드 자격증 목표는 단순 툴 사용이 아니라 시험 범위 중심으로 설계하세요.\n" +
			"- 최소한 핵심 개념과 서비스 범주 정리, 대표 아키텍처/운영 시나리오 이해, 유형별 문제풀이, 시험 직전 정리 단계를 포함하세요.\n" +
			"- 입문형은 서비스 개념 분류와 사례 판단을, 실무형은 운영·배포·보안·네트워크 시나리오를 더 강하게 다루세요.\n" +
			"- 실제 시험 응시와 자격 취득은 completion_criteria로 분리하세요.\n"
	case containsPracticalCertificationContext(lowerContext):
		return "- 자격증 실기형 목표는 이론 설명보다 공개과제/실기 과제 수행 흐름 중심으로 설계하세요.\n" +
			"- 최소한 시험 구조와 안전·위생·기본 기준 파악, 핵심 작업 패턴 익히기, 대표 과제 완성, 시간 안에 수행하는 실전 흐름, 시험 직전 점검 단계를 포함하세요.\n" +
			"- 실제 시험 응시와 합격은 completion_criteria로 분리하세요.\n"
	case containsTheoryCertificationContext(lowerContext):
		return "- 자격증 시험형 목표는 시험 범위 파악, 핵심 과목/영역 정리, 유형별 문제풀이, 약점 보완, 시험 직전 점검 흐름으로 구성하세요.\n" +
			"- 개념 정의만 나열하지 말고, 실제 문제를 풀 수 있는 수준으로 lesson 경계를 잡으세요.\n" +
			"- 실제 시험 응시와 합격은 completion_criteria로 분리하세요.\n"
	case strings.Contains(lowerContext, "기타") && strings.Contains(lowerContext, "찬양"):
		return "- 기타·찬양 반주는 처음 진입 활동(기타 구조 파악, 기본 자세, 다운/업 스트로크, 한 줄 소리내기)을 첫 번째 리슨 하나로 묶으세요 (예: '기초 셋팅과 첫 소리 내기').\n" +
			"- 이후 리슨은 찬양 반주에 실제로 필요한 기술 관문 중심으로 잡으세요: 자주 쓰는 오픈 코드, 코드 전환 안정화, 스트로크 패턴, 찬양곡 한 곡 완주.\n" +
			"- '기본 자세'·'튜닝'·'기타 구조'는 단독 리슨이 되어서는 안 됩니다.\n" +
			"- 찬양단 실제 참여·합주는 completion_criteria로 분리하세요.\n"
	case strings.Contains(lowerContext, "기타"):
		return "- 기타 초보는 기타 구조 이해, 기본 자세, 다운/업 스트로크, 계이름 누르기, 한 줄 소리내기처럼 처음에 하는 진입 활동을 첫 번째 리슨 하나로 묶으세요 (예: '기초 셋팅과 첫 소리 내기').\n" +
			"- 이후 리슨은 학습자의 확정 목표에 맞는 실제 기술 관문으로 잡으세요. 통기타 초보의 대표 관문: 오픈 코드 소리 내기, 코드 전환, 바레코드(F코드) 통과, 스트로크/핑거피킹 패턴 선택, 목표 곡 완주.\n" +
			"- '기본 자세'·'튜닝'·'기타 구조'는 단독 리슨이 되어서는 안 됩니다.\n" +
			"- 학습자의 확정 목표(특정 곡, 장르, 상황)에 따라 관문 순서와 비중을 조정하세요.\n" +
			"- 공연·합주·공식 참여는 completion_criteria로 분리하세요.\n"
	case strings.Contains(lowerContext, "피아노") || strings.Contains(lowerContext, "드럼") || strings.Contains(lowerContext, "연주"):
		return "- 악기/연주 목표는 처음 진입 활동(악기 구조, 자세, 기본 소리 내기)을 첫 번째 리슨 하나로 묶으세요. '기본 자세'는 단독 리슨이 아닙니다.\n" +
			"- 이후 리슨은 학습자의 확정 목표에 맞는 실제 기술 관문 중심으로 잡으세요. 각 악기마다 관문이 다릅니다.\n" +
			"- 공연·합주·공식 참여는 completion_criteria로 분리하세요.\n"
	case strings.Contains(lowerContext, "회화") || strings.Contains(lowerContext, "영어") || strings.Contains(lowerContext, "일본어") || strings.Contains(lowerContext, "언어"):
		return "- 언어 목표는 최소한 핵심 표현, 기본 문장 패턴, 질문/응답 연결, 상황별 적용 단계를 포함하세요.\n"
	case strings.Contains(lowerContext, "그림") || strings.Contains(lowerContext, "수채화") || strings.Contains(lowerContext, "드로잉") || strings.Contains(lowerContext, "스케치"):
		return "- 시각 창작 목표는 최소한 재료/도구, 기초 표현, 작은 소재 완성, 목표 스타일 적용 단계를 포함하세요.\n"
	default:
		return "- 목표를 이루기 위한 핵심 요소 파악, 기본 기술 습득, 반복 연습, 작은 적용, 목표 직전 통합 단계가 자연스럽게 이어지게 구성하세요.\n"
	}
}

func lookupGoalPatternSpec(key goalPatternKey) (goalPatternSpec, bool) {
	if spec, ok := goalPatternSpecs[key]; ok {
		return spec, true
	}
	key.GoalModeSecondary = ""
	spec, ok := goalPatternSpecs[key]
	return spec, ok
}

func buildGoalPatternGuidance(req CreateCourseDraftRequest) string {
	key := inferGoalPatternKey(req)
	spec, ok := lookupGoalPatternSpec(key)
	if !ok {
		return ""
	}
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		return buildGoalPatternGuidanceEnglish(spec)
	}

	var b strings.Builder
	if spec.LessonCountMin > 0 && spec.LessonCountMax > 0 {
		b.WriteString(fmt.Sprintf("- 권장 lesson 수 범위: %d~%d개\n", spec.LessonCountMin, spec.LessonCountMax))
	}
	if len(spec.StepRoles) > 0 {
		b.WriteString("- 권장 단계 흐름: " + strings.Join(describeStepRoles(spec.StepRoles), " -> ") + "\n")
	}
	if strings.TrimSpace(spec.LastLessonRule) != "" {
		b.WriteString("- " + strings.TrimSpace(spec.LastLessonRule) + "\n")
	}
	for _, rule := range spec.CompletionCriteriaRules {
		if strings.TrimSpace(rule) != "" {
			b.WriteString("- " + strings.TrimSpace(rule) + "\n")
		}
	}
	if len(spec.AtomicDisallowExamples) > 0 {
		b.WriteString("- 독립 lesson 금지 예: " + strings.Join(spec.AtomicDisallowExamples, ", ") + "\n")
	}
	return b.String()
}

func buildGoalPatternGuidanceEnglish(spec goalPatternSpec) string {
	var b strings.Builder
	if spec.LessonCountMin > 0 && spec.LessonCountMax > 0 {
		b.WriteString(fmt.Sprintf("- Recommended lesson count range: %d-%d\n", spec.LessonCountMin, spec.LessonCountMax))
	}
	if len(spec.StepRoles) > 0 {
		b.WriteString("- Recommended stage flow: " + strings.Join(describeStepRolesEnglish(spec.StepRoles), " -> ") + "\n")
	}
	b.WriteString("- Each lesson must be a gateway stage toward the confirmed goal, not a tiny isolated topic.\n")
	b.WriteString("- Keep the final lesson as preparation, integration, rehearsal, or pre-launch work before the final goal event.\n")
	b.WriteString("- Put the final event itself, such as passing, publishing, performing, selling, or joining, in completion_criteria instead of as a lesson title.\n")
	if len(spec.AtomicDisallowExamples) > 0 {
		b.WriteString("- Do not make atomic setup details into standalone lessons; absorb them into a larger first-stage lesson.\n")
	}
	return b.String()
}

func buildGoalSubpatternGuidance(req CreateCourseDraftRequest) string {
	subpatternKey := inferGoalSubpatternKey(req)
	if subpatternKey == "" {
		return ""
	}
	spec, ok := goalSubpatternSpecs[subpatternKey]
	if !ok {
		return ""
	}
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		return buildGoalSubpatternGuidanceEnglish(subpatternKey, spec)
	}

	var b strings.Builder
	if strings.TrimSpace(spec.Name) != "" {
		b.WriteString("- 세부 목표형: " + strings.TrimSpace(spec.Name) + "\n")
	}
	if len(spec.StepRoles) > 0 {
		b.WriteString("- 이 세부 목표에서 권장하는 단계 흐름: " + strings.Join(describeStepRoles(spec.StepRoles), " -> ") + "\n")
	}
	for _, rule := range spec.StageRules {
		if strings.TrimSpace(rule) != "" {
			b.WriteString("- " + strings.TrimSpace(rule) + "\n")
		}
	}
	if strings.TrimSpace(spec.LastLessonRule) != "" {
		b.WriteString("- " + strings.TrimSpace(spec.LastLessonRule) + "\n")
	}
	for _, rule := range spec.CompletionCriteriaRules {
		if strings.TrimSpace(rule) != "" {
			b.WriteString("- " + strings.TrimSpace(rule) + "\n")
		}
	}
	if len(spec.AtomicDisallowExamples) > 0 {
		b.WriteString("- 이 세부 목표에서 특히 독립 lesson으로 만들지 말아야 할 예: " + strings.Join(spec.AtomicDisallowExamples, ", ") + "\n")
	}
	return b.String()
}

func buildGoalSubpatternGuidanceEnglish(subpatternKey string, spec goalSubpatternSpec) string {
	var b strings.Builder
	b.WriteString("- Recurring goal subtype: " + describeSubpatternEnglish(subpatternKey, spec.Name) + "\n")
	if len(spec.StepRoles) > 0 {
		b.WriteString("- Recommended flow for this subtype: " + strings.Join(describeStepRolesEnglish(spec.StepRoles), " -> ") + "\n")
	}
	b.WriteString("- Preserve the subtype-specific progression, but write all learner-facing titles, objectives, and completion criteria in English.\n")
	b.WriteString("- Avoid making tiny setup details into standalone lessons; combine them into a larger gateway stage.\n")
	b.WriteString("- Keep the final target event as completion criteria, not as the final lesson itself.\n")
	return b.String()
}

func describeSubpatternEnglish(subpatternKey, fallbackName string) string {
	parts := strings.Split(subpatternKey, ":")
	if len(parts) > 0 {
		key := strings.ReplaceAll(parts[len(parts)-1], "_", " ")
		if key != "" {
			return key
		}
	}
	if strings.TrimSpace(fallbackName) != "" {
		return strings.TrimSpace(fallbackName)
	}
	return "goal-specific subtype"
}

func describeStepRoles(stepRoles []string) []string {
	result := make([]string, 0, len(stepRoles))
	for _, role := range stepRoles {
		switch role {
		case "setup_intro":
			result = append(result, "기초 셋업과 진입")
		case "first_output":
			result = append(result, "첫 수행 또는 첫 결과물")
		case "core_pattern":
			result = append(result, "핵심 패턴 반복")
		case "structure_plan":
			result = append(result, "구조와 설계 정리")
		case "integration_practice":
			result = append(result, "통합 연습과 흐름 연결")
		case "artifact_finish":
			result = append(result, "결과물 완성")
		case "performance_prep":
			result = append(result, "실전 직전 준비")
		case "teaching_prep":
			result = append(result, "설명과 시연 준비")
		case "publish_prep":
			result = append(result, "공개와 게시 준비")
		case "prediction_prompt":
			result = append(result, "예측 질문")
		case "interactive_exploration":
			result = append(result, "상호작용 탐구")
		case "variable_experiment":
			result = append(result, "변수 실험")
		case "explanation_finish":
			result = append(result, "설명 정리")
		case "peer_discussion":
			result = append(result, "관찰 공유와 토론")
		case "reflection_finish":
			result = append(result, "성찰 정리")
		case "safety_intro":
			result = append(result, "안전 기준과 진입")
		case "procedure_plan":
			result = append(result, "절차 계획")
		case "diagnostic_practice":
			result = append(result, "진단 연습")
		case "analysis_project":
			result = append(result, "분석 프로젝트")
		case "presentation_prep":
			result = append(result, "발표 준비")
		case "empathy_setup":
			result = append(result, "사용자 공감")
		case "problem_define":
			result = append(result, "문제 정의")
		case "ideation_practice":
			result = append(result, "아이디어 확장")
		case "prototype_build":
			result = append(result, "프로토타입 제작")
		case "test_reflection":
			result = append(result, "테스트와 성찰")
		case "question_setup":
			result = append(result, "검증 질문 설정")
		case "prototype_plan":
			result = append(result, "프로토타입 계획")
		case "scrappy_build":
			result = append(result, "간단한 구성 제작")
		case "test_session":
			result = append(result, "반응 테스트")
		case "learning_reflection":
			result = append(result, "배운 점 정리")
		case "context_intro":
			result = append(result, "맥락 이해")
		case "data_reading":
			result = append(result, "데이터 읽기")
		case "model_build":
			result = append(result, "모델 만들기")
		case "trend_compare":
			result = append(result, "추세 비교")
		case "field_setup":
			result = append(result, "현장 준비")
		case "observation_capture":
			result = append(result, "관찰 수집")
		case "identification_practice":
			result = append(result, "식별 연습")
		case "community_record_finish":
			result = append(result, "관찰 기록 정리")
		case "source_setup":
			result = append(result, "자료 맥락 파악")
		case "observation_practice":
			result = append(result, "세부 관찰")
		case "reflection_practice":
			result = append(result, "추론과 반성")
		case "question_generation":
			result = append(result, "질문 생성")
		case "investigation_finish":
			result = append(result, "추가 조사 정리")
		default:
			result = append(result, role)
		}
	}
	return result
}

func describeStepRolesEnglish(stepRoles []string) []string {
	result := make([]string, 0, len(stepRoles))
	for _, role := range stepRoles {
		switch role {
		case "setup_intro":
			result = append(result, "setup and entry")
		case "first_output":
			result = append(result, "first performance or first output")
		case "core_pattern":
			result = append(result, "core pattern practice")
		case "structure_plan":
			result = append(result, "structure and plan")
		case "integration_practice":
			result = append(result, "integrated practice")
		case "artifact_finish":
			result = append(result, "finish the output")
		case "performance_prep":
			result = append(result, "pre-performance preparation")
		case "teaching_prep":
			result = append(result, "teaching and demo preparation")
		case "publish_prep":
			result = append(result, "publishing preparation")
		case "prediction_prompt":
			result = append(result, "prediction prompt")
		case "interactive_exploration":
			result = append(result, "interactive exploration")
		case "variable_experiment":
			result = append(result, "variable experiment")
		case "explanation_finish":
			result = append(result, "finish the explanation")
		case "peer_discussion":
			result = append(result, "peer discussion")
		case "reflection_finish":
			result = append(result, "finish the reflection")
		case "safety_intro":
			result = append(result, "safety setup")
		case "procedure_plan":
			result = append(result, "procedure plan")
		case "diagnostic_practice":
			result = append(result, "diagnostic practice")
		case "analysis_project":
			result = append(result, "analysis project")
		case "presentation_prep":
			result = append(result, "presentation preparation")
		case "empathy_setup":
			result = append(result, "empathy setup")
		case "problem_define":
			result = append(result, "problem definition")
		case "ideation_practice":
			result = append(result, "ideation practice")
		case "prototype_build":
			result = append(result, "prototype build")
		case "test_reflection":
			result = append(result, "test and reflection")
		case "question_setup":
			result = append(result, "test question setup")
		case "prototype_plan":
			result = append(result, "prototype plan")
		case "scrappy_build":
			result = append(result, "scrappy build")
		case "test_session":
			result = append(result, "test session")
		case "learning_reflection":
			result = append(result, "learning reflection")
		case "context_intro":
			result = append(result, "context introduction")
		case "data_reading":
			result = append(result, "data reading")
		case "model_build":
			result = append(result, "model build")
		case "trend_compare":
			result = append(result, "trend comparison")
		case "field_setup":
			result = append(result, "field setup")
		case "observation_capture":
			result = append(result, "observation capture")
		case "identification_practice":
			result = append(result, "identification practice")
		case "community_record_finish":
			result = append(result, "finish the community record")
		case "source_setup":
			result = append(result, "source setup")
		case "observation_practice":
			result = append(result, "observation practice")
		case "reflection_practice":
			result = append(result, "reflection practice")
		case "question_generation":
			result = append(result, "question generation")
		case "investigation_finish":
			result = append(result, "finish the investigation note")
		default:
			result = append(result, role)
		}
	}
	return result
}
