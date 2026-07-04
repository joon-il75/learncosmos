package curriculum

import (
	"fmt"
	"strings"
)

type fallbackDraftPlan struct {
	Description        string
	MainLessons        []generatedMainLesson
	CompletionCriteria []string
}

func buildFallbackDraftPlan(req CreateCourseDraftRequest) fallbackDraftPlan {
	sourceQuery := strings.TrimSpace(req.SourceQuery)
	learningGoal := strings.TrimSpace(derefString(req.LearningGoal))
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		return buildEnglishFallbackDraftPlan(req, sourceQuery, learningGoal)
	}

	description := fmt.Sprintf("'%s' 목표를 달성하기 위해 '%s'를 단계적으로 익히는 초기 커리큘럼 초안입니다.", learningGoal, sourceQuery)
	if usage := trimmedGoalField(req.GoalUsageContext); usage != "" {
		description = fmt.Sprintf("'%s' 목표를 '%s' 맥락에서 달성하기 위해 '%s'를 단계적으로 익히는 초기 커리큘럼 초안입니다.", learningGoal, usage, sourceQuery)
	}
	lowerContext := buildGoalInferenceContext(req)

	switch {
	case containsEnglishScoreCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 구조와 파트 구성을 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 파트 구조와 점수형 시험 운영 방식을 파악한다.", learningGoal)},
				{Title: "어휘·문법·독해·청해 핵심을 묶어 익힌다", Objective: "빈출 어휘와 문법, 독해, 청해의 핵심 패턴을 파트별 문제 맥락으로 정리한다."},
				{Title: "파트별 문제풀이와 시간 대응을 연습한다", Objective: "각 파트의 문제 흐름을 익히고 제한 시간 안에 풀어내는 연습을 한다."},
				{Title: "시험 직전 점수형 운영 흐름을 정리한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 파트별 약점과 시간 배분을 점검한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"파트별 핵심 문제 유형에 맞는 풀이 흐름을 적용할 수 있다.",
				"어휘·문법·독해·청해 약점을 스스로 구분해 보완할 수 있다.",
				"모의 문제를 실제 시험 시간 안에 운영할 수 있다.",
			},
		}
	case containsEnglishSpeakingInterviewCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "인터뷰 구조와 자기소개 흐름을 익힌다", Objective: fmt.Sprintf("'%s' 목표에 필요한 인터뷰 구성과 기본 자기소개 흐름을 정리한다.", learningGoal)},
				{Title: "주제별 답변 패턴을 만든다", Objective: "일상, 경험, 취미, 장소 같은 자주 나오는 주제에 대해 답변 구조를 만든다."},
				{Title: "돌발 질문과 확장 답변을 연결한다", Objective: "비교, 문제 해결, 과거 경험 같은 확장형 질문에 즉흥적으로 답하는 연습을 한다."},
				{Title: "시험 직전 인터뷰 실전 흐름을 점검한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 답변 길이, 자연스러움, 주제 전환 흐름을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"주요 주제에 대해 짧고 일관된 답변 흐름을 만들 수 있다.",
				"돌발 질문에도 핵심 내용을 유지하며 답변을 이어갈 수 있다.",
				"실전 인터뷰 길이에 맞춰 답변 속도와 분량을 조절할 수 있다.",
			},
		}
	case containsEnglishIntegratedFourSkillsCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 유형과 4영역 구조를 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 시험 유형과 Reading, Listening, Writing, Speaking 구조를 파악한다.", learningGoal)},
				{Title: "Reading·Listening 기반을 다진다", Objective: "지문 이해, 강의·대화 청취, 핵심 정보 파악을 시험 task 기준으로 정리한다."},
				{Title: "Writing·Speaking task를 연결해 연습한다", Objective: "쓰기와 말하기 task를 단독으로가 아니라 읽기·듣기 입력과 연결해 수행하는 연습을 한다."},
				{Title: "시험 직전 4영역 운영 흐름을 점검한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 4영역 task 전환과 시간 배분을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"4영역 task의 요구사항 차이를 구분해 대응할 수 있다.",
				"읽기·듣기 입력을 바탕으로 말하기·쓰기 응답을 구성할 수 있다.",
				"시험 직전 4영역 전체 운영 흐름을 스스로 점검할 수 있다.",
			},
		}
	case containsLanguageCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 범위와 N4 핵심 표현 범주를 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 시험 구성과 핵심 어휘·표현 범위를 파악한다.", learningGoal)},
				{Title: "핵심 문법과 필수 어휘를 문제형으로 익힌다", Objective: "자주 나오는 문법, 어휘, 문자 문제를 문장 단위와 문제풀이 흐름으로 익힌다."},
				{Title: "독해와 청해의 핵심 정보를 잡아낸다", Objective: "짧은 독해와 청해에서 주제, 의도, 핵심 정보를 안정적으로 찾는 연습을 한다."},
				{Title: "유형별 문제풀이 흐름을 정리한다", Objective: "문자·어휘·문법·독해·청해 문제를 시간 안에 푸는 기본 흐름을 만든다."},
				{Title: "시험 직전 약점을 보완하고 실전 감각을 맞춘다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 약점 영역과 시간 배분을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"핵심 문법과 빈출 어휘를 문제풀이에 적용할 수 있다.",
				"독해와 청해에서 핵심 정보를 안정적으로 찾을 수 있다.",
				"모의 문제를 실제 시험 흐름에 가깝게 풀 수 있다.",
			},
		}
	case containsCloudCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 구조와 핵심 서비스 범주를 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 시험 구성과 주요 서비스 범주를 파악한다.", learningGoal)},
				{Title: "핵심 개념과 대표 서비스를 연결한다", Objective: "컴퓨트, 스토리지, 네트워크, 보안 같은 핵심 개념과 대표 서비스를 사례와 함께 이해한다."},
				{Title: "대표 시나리오로 서비스 선택 흐름을 익힌다", Objective: "요구사항에 맞는 서비스 선택과 기본 아키텍처 판단 흐름을 연습한다."},
				{Title: "유형별 문제풀이와 운영 판단을 연습한다", Objective: "입문형 또는 실무형 시험에서 자주 나오는 시나리오 문제를 풀며 약점을 정리한다."},
				{Title: "시험 직전 핵심 서비스와 시나리오를 점검한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 핵심 서비스, 시나리오, 시간 배분을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"핵심 서비스 범주와 역할을 구분해 설명할 수 있다.",
				"대표 시나리오에서 적절한 서비스 선택 흐름을 적용할 수 있다.",
				"모의 문제를 풀며 약점 영역을 스스로 점검할 수 있다.",
			},
		}
	case containsPracticalCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 구조와 기본 안전·위생 기준을 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 시험 구성과 기본 안전·위생 기준을 파악한다.", learningGoal)},
				{Title: "핵심 작업 패턴과 기본 조작을 익힌다", Objective: "대표 실기 과제에 필요한 기본 조작과 핵심 작업 패턴을 반복 연습한다."},
				{Title: "대표 과제를 순서대로 완성한다", Objective: "실기 공개과제나 대표 작업을 시작부터 끝까지 순서대로 수행하는 흐름을 익힌다."},
				{Title: "시간 안에 수행하는 실전 흐름을 만든다", Objective: "제한 시간 안에 작업 순서, 정확도, 마무리를 유지하는 연습을 한다."},
				{Title: "시험 직전 체크포인트를 점검한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 안전·위생·작업 순서·시간 배분을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"대표 실기 과제를 순서대로 수행할 수 있다.",
				"안전·위생·기본 기준을 지키며 작업을 이어갈 수 있다.",
				"시험 직전 체크포인트와 시간 배분 흐름을 스스로 점검할 수 있다.",
			},
		}
	case containsTheoryCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "시험 구조와 핵심 과목 범위를 정리한다", Objective: fmt.Sprintf("'%s' 목표에 필요한 시험 구조와 핵심 과목 범위를 파악한다.", learningGoal)},
				{Title: "핵심 개념과 빈출 영역을 묶어 익힌다", Objective: "주요 과목과 빈출 개념을 대표 문제 맥락과 함께 정리한다."},
				{Title: "유형별 문제풀이 흐름을 만든다", Objective: "계산형, 개념형, 응용형 등 대표 문제 유형을 구분해 풀이 흐름을 익힌다."},
				{Title: "약점 과목과 빈출 문제를 보완한다", Objective: "취약한 과목과 자주 틀리는 문제 유형을 다시 묶어 정리한다."},
				{Title: "시험 직전 과목별 풀이 흐름을 점검한다", Objective: fmt.Sprintf("'%s' 목표와 연결되는 시험 직전 수준까지 과목별 핵심과 시간 배분을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"핵심 과목과 빈출 개념을 문제풀이에 연결할 수 있다.",
				"대표 문제 유형별 풀이 흐름을 스스로 적용할 수 있다.",
				"시험 직전 약한 과목을 중심으로 복습 흐름을 정리할 수 있다.",
			},
		}
	case strings.Contains(lowerContext, "기타") && strings.Contains(lowerContext, "찬양"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "기초 셋팅과 첫 소리 내기", Objective: "기타 구조와 기본 자세를 익히고, 다운/업 스트로크로 한 줄씩 소리를 낼 수 있게 한다."},
				{Title: "찬양 반주 오픈 코드 익히기", Objective: "찬양곡에서 자주 쓰는 오픈 코드(Am, Em, G, C, D)를 짚어 분명한 소리를 낼 수 있게 한다."},
				{Title: "코드 전환과 스트로크 패턴 안정화하기", Objective: "코드 사이를 끊기지 않고 전환하면서 기본 스트로크 패턴을 유지할 수 있게 한다."},
				{Title: "쉬운 찬양곡 한 곡 완주하기", Objective: "쉬운 찬양곡 한 곡을 처음부터 끝까지 안정적으로 반주할 수 있게 연습한다."},
				{Title: "찬양단 합주 흐름에 맞춰 준비하기", Objective: "곡 순서, 템포, 시작과 마침을 맞추며 찬양단 합주에 참여할 준비를 한다."},
			},
			CompletionCriteria: []string{
				"찬양곡 1~2곡을 기본 코드와 스트로크로 끊기지 않고 반주할 수 있다.",
				"자주 쓰는 오픈 코드 전환을 템포를 크게 잃지 않고 이어갈 수 있다.",
				"찬양단 연습에서 시작, 진행, 마무리 타이밍을 따라가며 반주에 참여할 수 있다.",
			},
		}
	case strings.Contains(lowerContext, "기타"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "기초 셋팅과 첫 소리 내기", Objective: "기타 구조와 기본 자세를 파악하고, 다운/업 스트로크와 한 줄씩 계이름 누르기로 처음 소리를 낸다."},
				{Title: "오픈 코드 소리 내기", Objective: fmt.Sprintf("'%s' 목표에 필요한 오픈 코드(Am, Em, G, C, D)를 짚어 분명한 소리를 낼 수 있게 한다.", learningGoal)},
				{Title: "코드 전환과 스트로크 패턴 익히기", Objective: "코드 간 전환을 끊기지 않게 이어가고, 목표 곡 스타일에 맞는 스트로크 패턴을 반복 연습한다."},
				{Title: "바레코드(F코드) 통과하기", Objective: "F코드를 포함한 바레코드를 짚어 소리를 내고, 앞 코드와 전환할 수 있게 한다."},
				{Title: "목표 곡 한 곡 완주하기", Objective: fmt.Sprintf("'%s' 목표와 연결되는 곡을 처음부터 끝까지 안정적으로 연주할 수 있게 한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"오픈 코드와 바레코드를 포함한 코드 전환을 끊기지 않고 이어갈 수 있다.",
				"목표 곡 한 곡을 처음부터 끝까지 일정한 리듬으로 연주할 수 있다.",
				"배운 코드와 스트로크 패턴을 새로운 쉬운 곡에 스스로 적용할 수 있다.",
			},
		}
	case strings.Contains(lowerContext, "피아노") || strings.Contains(lowerContext, "드럼") || strings.Contains(lowerContext, "연주"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "기초 셋팅과 첫 소리 내기", Objective: fmt.Sprintf("'%s' 목표에 필요한 악기 구조와 기본 자세를 익히고, 처음 소리를 낸다.", learningGoal)},
				{Title: "핵심 기술과 기본 패턴 익히기", Objective: fmt.Sprintf("'%s'에 필요한 핵심 기술과 기본 패턴을 반복 연습한다.", learningGoal)},
				{Title: "연결과 리듬 안정화하기", Objective: "배운 기술을 자연스럽게 연결하고 리듬과 속도를 안정적으로 유지한다."},
				{Title: "짧은 곡이나 예제로 적용하기", Objective: "배운 기술을 짧은 곡이나 예제에 적용해 실제 연주 흐름을 익힌다."},
				{Title: "목표 곡 완주 준비하기", Objective: fmt.Sprintf("'%s' 목표와 연결되는 연주를 처음부터 끝까지 이어갈 수 있게 준비한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"핵심 기술과 기본 패턴을 끊기지 않게 이어서 수행할 수 있다.",
				"짧은 예제나 곡을 처음부터 끝까지 안정적으로 연주할 수 있다.",
				"실전 상황 직전까지 필요한 준비 흐름을 스스로 반복할 수 있다.",
			},
		}
	case strings.Contains(lowerContext, "회화") || strings.Contains(lowerContext, "영어") || strings.Contains(lowerContext, "일본어") || strings.Contains(lowerContext, "언어"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "핵심 표현과 상황 정리하기", Objective: fmt.Sprintf("'%s' 목표에 필요한 핵심 상황과 표현 범위를 정리한다.", learningGoal)},
				{Title: "기본 문장 패턴 익히기", Objective: "자주 쓰는 기본 문장 패턴을 말할 수 있도록 반복 연습한다."},
				{Title: "짧은 응답과 질문 연결하기", Objective: "질문과 대답을 짧게 이어가며 실제 대화 흐름의 기초를 만든다."},
				{Title: "상황별 대화 연습하기", Objective: "목표 상황에서 필요한 표현을 묶어 실제 대화처럼 연습한다."},
				{Title: "실전 대화 직전 점검하기", Objective: fmt.Sprintf("'%s' 목표와 연결되는 실제 대화 직전 수준까지 흐름을 정리한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"목표 상황에서 자주 쓰는 표현을 스스로 말할 수 있다.",
				"짧은 질문과 대답을 이어가며 대화를 몇 차례 유지할 수 있다.",
				"실전 상황 직전까지 필요한 표현 흐름을 큰 막힘 없이 사용할 수 있다.",
			},
		}
	case strings.Contains(lowerContext, "그림") || strings.Contains(lowerContext, "수채화") || strings.Contains(lowerContext, "드로잉") || strings.Contains(lowerContext, "스케치"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "재료와 기본 사용법 익히기", Objective: fmt.Sprintf("'%s' 목표에 필요한 재료와 기본 사용법을 익힌다.", learningGoal)},
				{Title: "기초 표현 연습하기", Objective: "선, 형태, 명암, 색 같은 기초 표현 요소를 반복 연습한다."},
				{Title: "작은 소재 완성하기", Objective: "짧고 작은 작업을 완성하며 표현 요소를 실제 결과물로 연결한다."},
				{Title: "목표 스타일로 적용하기", Objective: fmt.Sprintf("'%s' 목표에 맞는 스타일과 구성을 적용해본다.", learningGoal)},
				{Title: "최종 결과물 직전까지 다듬기", Objective: "구성, 완성도, 표현의 일관성을 다듬으며 최종 결과물 직전 단계까지 정리한다."},
			},
			CompletionCriteria: []string{
				"기본 표현 요소를 조합해 작은 작업을 완성할 수 있다.",
				"목표 스타일에 맞는 구성과 표현 방식을 스스로 선택할 수 있다.",
				"최종 결과물 직전 단계까지 큰 막힘 없이 작업을 이어갈 수 있다.",
			},
		}
	default:
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "목표에 필요한 요소 파악하기", Objective: fmt.Sprintf("'%s' 목표에 필요한 핵심 요소와 학습 범위를 정리한다.", learningGoal)},
				{Title: "기본 도구와 핵심 기술 익히기", Objective: fmt.Sprintf("'%s'를 위해 필요한 기본 도구와 핵심 기술을 익힌다.", learningGoal)},
				{Title: "반복 연습으로 흐름 안정화하기", Objective: "기본 기술을 반복해 익히고 자연스럽게 이어지는 작업 흐름을 만든다."},
				{Title: "작은 결과물로 적용하기", Objective: "작은 결과물이나 짧은 실습으로 배운 기술을 실제로 적용한다."},
				{Title: "목표 직전 단계까지 연결하기", Objective: fmt.Sprintf("'%s' 목표 직전 단계까지 필요한 준비와 적용 흐름을 완성한다.", learningGoal)},
			},
			CompletionCriteria: []string{
				"핵심 기술을 스스로 다시 수행할 수 있다.",
				"작은 결과물이나 실습 과제를 완성할 수 있다.",
				"최종 목표 직전 단계까지 필요한 흐름을 끊기지 않게 이어갈 수 있다.",
			},
		}
	}
}

func buildEnglishFallbackDraftPlan(req CreateCourseDraftRequest, sourceQuery, learningGoal string) fallbackDraftPlan {
	description := fmt.Sprintf("An initial staged curriculum draft for building '%s' toward the confirmed goal: '%s'.", sourceQuery, learningGoal)
	if usage := trimmedGoalField(req.GoalUsageContext); usage != "" {
		description = fmt.Sprintf("An initial staged curriculum draft for building '%s' in the context of '%s' toward the confirmed goal: '%s'.", sourceQuery, usage, learningGoal)
	}
	lowerContext := buildGoalInferenceContext(req)

	switch {
	case containsEnglishScoreCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Map the exam structure and scoring parts", Objective: fmt.Sprintf("Identify the exam format, part structure, and scoring flow needed for '%s'.", learningGoal)},
				{Title: "Build the core vocabulary and grammar base", Objective: "Organize frequent vocabulary, grammar, reading, and listening patterns by exam part."},
				{Title: "Practice part-based problem solving under time limits", Objective: "Apply part-specific solving routines and manage time across practice sets."},
				{Title: "Review the final test-taking routine", Objective: fmt.Sprintf("Check weak parts, timing, and final operating flow before attempting '%s'.", learningGoal)},
			},
			CompletionCriteria: []string{
				"Apply a solving routine to the main problem types in each exam part.",
				"Identify and improve weak vocabulary, grammar, reading, and listening areas.",
				"Run a timed practice set close to the real exam flow.",
			},
		}
	case containsEnglishSpeakingInterviewCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Understand the interview format and self-introduction flow", Objective: fmt.Sprintf("Organize the interview structure and basic speaking flow needed for '%s'.", learningGoal)},
				{Title: "Build answer patterns for common topics", Objective: "Create reusable answer structures for daily life, experiences, hobbies, and places."},
				{Title: "Connect unexpected questions to expanded answers", Objective: "Practice responding to comparison, problem-solving, and past-experience prompts."},
				{Title: "Run a final interview readiness check", Objective: "Check answer length, fluency, topic transitions, and final speaking rhythm before the test."},
			},
			CompletionCriteria: []string{
				"Give short, coherent answers for common interview topics.",
				"Keep the main point while responding to unexpected questions.",
				"Control answer speed and length for a realistic interview session.",
			},
		}
	case containsEnglishIntegratedFourSkillsCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Map the test format and four-skill structure", Objective: fmt.Sprintf("Understand the reading, listening, writing, and speaking structure needed for '%s'.", learningGoal)},
				{Title: "Build the reading and listening input base", Objective: "Practice finding key information in passages, lectures, and conversations."},
				{Title: "Connect writing and speaking tasks to input", Objective: "Use reading and listening input to produce structured writing and speaking responses."},
				{Title: "Run a final four-skill timing check", Objective: "Review task transitions and time allocation across all four skills before the test."},
			},
			CompletionCriteria: []string{
				"Distinguish the requirements of the four major task types.",
				"Use reading and listening input to support speaking and writing responses.",
				"Run a realistic final practice flow across all four skills.",
			},
		}
	case containsCloudCertificationContext(lowerContext) || containsTheoryCertificationContext(lowerContext) || containsPracticalCertificationContext(lowerContext) || containsLanguageCertificationContext(lowerContext):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Map the exam scope and core requirements", Objective: fmt.Sprintf("Identify the exam structure, required topics, and readiness criteria for '%s'.", learningGoal)},
				{Title: "Build the core concepts and task patterns", Objective: "Organize the main concepts, service areas, or task patterns into practical study units."},
				{Title: "Practice representative question or task types", Objective: "Apply the core material to representative problems, scenarios, or practical tasks."},
				{Title: "Strengthen weak areas through targeted review", Objective: "Find weak areas from practice results and rebuild the most important foundations."},
				{Title: "Run a final readiness check before the exam", Objective: "Check timing, accuracy, procedure, and final review priorities before the actual exam."},
			},
			CompletionCriteria: []string{
				"Explain the core scope and requirements of the target exam.",
				"Apply representative solving or task routines without major gaps.",
				"Review weak areas and complete a final exam-readiness routine.",
			},
		}
	case containsAny(lowerContext, "guitar", "piano", "drum", "instrument", "performance", "song"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Set up the basics and produce the first sound", Objective: fmt.Sprintf("Learn the basic setup, posture, and first playable sounds needed for '%s'.", learningGoal)},
				{Title: "Build the core technique patterns", Objective: "Practice the core finger, rhythm, chord, or movement patterns required by the goal."},
				{Title: "Connect patterns into a stable flow", Objective: "Link the learned patterns without stopping and stabilize rhythm, timing, and transitions."},
				{Title: "Apply the skills to a short target piece", Objective: "Use the learned skills in a short piece, section, or performance example."},
				{Title: "Prepare the final performance flow", Objective: "Review weak sections and prepare the full flow right before the target performance."},
			},
			CompletionCriteria: []string{
				"Perform the core patterns with stable timing and fewer interruptions.",
				"Apply the learned skills to a short target piece or example.",
				"Prepare the final performance flow up to the pre-performance stage.",
			},
		}
	case containsAny(lowerContext, "conversation", "english", "japanese", "language", "travel"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Organize the key situations and expressions", Objective: fmt.Sprintf("Identify the situations and expression range needed for '%s'.", learningGoal)},
				{Title: "Build reusable sentence patterns", Objective: "Practice the basic sentence patterns needed to speak or respond in the target situation."},
				{Title: "Connect short questions and answers", Objective: "Create a basic conversation flow by linking short questions and answers."},
				{Title: "Practice situation-based conversations", Objective: "Apply the expressions in realistic scenes related to the target goal."},
				{Title: "Run a final real-situation readiness check", Objective: "Review the speaking flow and key expressions before using them in the real situation."},
			},
			CompletionCriteria: []string{
				"Use common expressions for the target situation without major hesitation.",
				"Maintain several turns of short questions and answers.",
				"Prepare the expression flow needed right before the real situation.",
			},
		}
	case containsAny(lowerContext, "drawing", "painting", "sketch", "illustration", "calligraphy", "lettering", "design", "art"):
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Set up the materials and basic technique", Objective: fmt.Sprintf("Learn the materials, tools, and basic handling needed for '%s'.", learningGoal)},
				{Title: "Practice the core visual elements", Objective: "Practice lines, shapes, values, colors, composition, or lettering elements."},
				{Title: "Complete a small focused piece", Objective: "Turn the practiced elements into a small finished study or sample."},
				{Title: "Apply the target style or composition", Objective: "Use the target style, theme, or composition in a more complete work."},
				{Title: "Refine the work before final presentation", Objective: "Review composition, consistency, and finish before the final output stage."},
			},
			CompletionCriteria: []string{
				"Combine basic visual elements into a small completed work.",
				"Choose a style or composition that matches the target goal.",
				"Bring the final work close to a presentable stage.",
			},
		}
	default:
		return fallbackDraftPlan{
			Description: description,
			MainLessons: []generatedMainLesson{
				{Title: "Identify the core requirements for the goal", Objective: fmt.Sprintf("Clarify the key elements, scope, and first priorities needed for '%s'.", learningGoal)},
				{Title: "Build the essential tools and core skills", Objective: fmt.Sprintf("Practice the basic tools, concepts, and skills needed to make progress in '%s'.", sourceQuery)},
				{Title: "Stabilize the workflow through repeated practice", Objective: "Repeat the core actions until the workflow becomes smoother and more reliable."},
				{Title: "Apply the skills in a small practical output", Objective: "Use the learned skills in a small project, example, or practical task."},
				{Title: "Prepare the final pre-goal integration", Objective: "Connect the weak parts, review the full flow, and prepare the step right before the final goal."},
			},
			CompletionCriteria: []string{
				"Explain the main requirements and scope of the confirmed goal.",
				"Apply the core skills in a small practical output or task.",
				"Prepare the integrated workflow needed right before the final goal.",
			},
		}
	}
}
