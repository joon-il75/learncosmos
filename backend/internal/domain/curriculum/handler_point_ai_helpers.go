package curriculum

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func pointAILanguageInstruction(language string) string {
	if normalizeLearningLanguage(language) == "en" {
		return "Respond in English. Learner-facing summaries, questions, feedback, reasons, intents, and guidance must be in English. Keep JSON keys and enum values unchanged."
	}
	return "한국어로 응답하세요. 학습자에게 보이는 요약, 질문, 피드백, 이유, 의도, 안내 문구는 한국어로 작성하세요. JSON key와 enum 값은 변경하지 마세요."
}

func buildPointAISummaryPrompt(pointTitle, pointDescription, sourceTitle, sourceDescription, learningLanguage string) string {
	var sb strings.Builder
	sb.WriteString("당신은 학습자가 지금 바로 이해하고 다음 행동을 정할 수 있도록 돕는 학습 요약 도우미입니다.\n")
	sb.WriteString(pointAILanguageInstruction(learningLanguage))
	sb.WriteString("\n아래 탐험지점 정보와 원문 메타 정보를 바탕으로 3~4문장 요약만 작성하세요. 번호 목록, 제목, 마크다운은 쓰지 마세요.\n\n")
	sb.WriteString(fmt.Sprintf("탐험지점 제목: %s\n", strings.TrimSpace(pointTitle)))
	if strings.TrimSpace(pointDescription) != "" {
		sb.WriteString(fmt.Sprintf("탐험지점 설명: %s\n", strings.TrimSpace(pointDescription)))
	}
	if strings.TrimSpace(sourceTitle) != "" {
		sb.WriteString(fmt.Sprintf("원문 제목: %s\n", strings.TrimSpace(sourceTitle)))
	}
	if strings.TrimSpace(sourceDescription) != "" {
		sb.WriteString(fmt.Sprintf("원문 설명: %s\n", strings.TrimSpace(sourceDescription)))
	}
	sb.WriteString("\n요약 조건:\n")
	sb.WriteString("- 학습자가 이 지점에서 꼭 이해해야 할 핵심을 압축할 것\n")
	sb.WriteString("- 원문 설명을 그대로 길게 반복하지 말고, 학습 맥락에 맞춰 재구성할 것\n")
	sb.WriteString("- 마지막 문장은 바로 해볼 수 있는 작은 실천으로 마무리할 것\n")
	return sb.String()
}

func normalizePointAISummaryText(raw string) string {
	text := strings.TrimSpace(raw)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	if strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}") {
		var payload map[string]any
		if err := json.Unmarshal([]byte(text), &payload); err == nil {
			if summary := extractPointAISummaryString(payload); summary != "" {
				return summary
			}
			if len(payload) == 1 {
				for _, value := range payload {
					if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
						return strings.TrimSpace(str)
					}
				}
			}
		}
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(text, "{"), "}"))
		if inner != "" && !strings.Contains(inner, "\":") {
			return inner
		}
	}

	return text
}

func extractPointAISummaryString(payload map[string]any) string {
	for _, key := range []string{"summary", "요약", "text", "content", "answer", "result", "안내문구"} {
		if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	for _, value := range payload {
		if nested, ok := value.(map[string]any); ok {
			if summary := extractPointAISummaryString(nested); summary != "" {
				return summary
			}
		}
	}
	return ""
}

func findPointInPlanetLessons(lessons []CourseLessonTree, pointID uuid.UUID) (*CoursePointAggregate, string, bool) {
	for _, lesson := range lessons {
		for idx := range lesson.Points {
			if lesson.Points[idx].Point.ID == pointID {
				return &lesson.Points[idx], lesson.Lesson.Title, true
			}
		}
		if point, lessonTitle, ok := findPointInPlanetLessons(lesson.SubLessons, pointID); ok {
			return point, lessonTitle, true
		}
	}
	return nil, "", false
}

func buildPointAIQuestionPrompt(goalContext *PlanetGoalContext, point *CoursePointAggregate, lessonTitle, learningLanguage string) string {
	var sb strings.Builder
	sb.WriteString("당신은 학습자의 현재 목표와 연결되는 질문을 한 개만 제안하는 학습 코치입니다.\n")
	sb.WriteString(pointAILanguageInstruction(learningLanguage))
	sb.WriteString("\n")
	sb.WriteString("반드시 JSON 한 줄만 출력하세요. 형식: {\"question\":\"...\",\"question_type\":\"reflection|application|goal_alignment\"}\n")
	sb.WriteString("질문은 한 문장으로만 쓰고, 점수 평가나 정답 판정 표현은 금지합니다.\n")
	sb.WriteString("기존 질문과 최대한 겹치지 않게, 지금 지점의 이해 또는 적용을 앞으로 밀어주는 질문이어야 합니다.\n\n")
	sb.WriteString(fmt.Sprintf("지점 타입: %s\n", point.Point.PointType))
	sb.WriteString(fmt.Sprintf("지점 제목: %s\n", strings.TrimSpace(point.Point.Title)))
	if strings.TrimSpace(derefStr(point.Point.Description)) != "" {
		sb.WriteString(fmt.Sprintf("지점 설명: %s\n", strings.TrimSpace(derefStr(point.Point.Description))))
	}
	if strings.TrimSpace(lessonTitle) != "" {
		sb.WriteString(fmt.Sprintf("소속 리슨: %s\n", strings.TrimSpace(lessonTitle)))
	}
	if goalContext != nil {
		if strings.TrimSpace(derefStr(goalContext.LearningGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.LearningGoal))))
		} else if strings.TrimSpace(derefStr(goalContext.ConfirmedGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.ConfirmedGoal))))
		}
		if strings.TrimSpace(derefStr(goalContext.UsageContext)) != "" {
			sb.WriteString(fmt.Sprintf("활용 맥락: %s\n", strings.TrimSpace(derefStr(goalContext.UsageContext))))
		}
		if strings.TrimSpace(derefStr(goalContext.Motivation)) != "" {
			sb.WriteString(fmt.Sprintf("동기: %s\n", strings.TrimSpace(derefStr(goalContext.Motivation))))
		}
	}
	if point.AISummary != nil && strings.TrimSpace(point.AISummary.Summary) != "" {
		sb.WriteString(fmt.Sprintf("AI 요약: %s\n", strings.TrimSpace(point.AISummary.Summary)))
	}
	if len(point.Questions) > 0 {
		sb.WriteString("기존 질문:\n")
		for _, question := range point.Questions {
			sb.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(question.Question)))
		}
	}
	sb.WriteString("\n질문 생성 기준:\n")
	sb.WriteString("- exploration이면 실제 적용이나 관찰을 앞으로 밀어주는 질문을 우선한다.\n")
	sb.WriteString("- research이면 정리한 내용을 목표와 연결하거나 다음 정리 행동을 분명하게 만드는 질문을 우선한다.\n")
	sb.WriteString("- 질문 하나만 제안한다.\n")
	return sb.String()
}

func parsePointAIQuestionSuggestion(raw string, pointType PointType, learningLanguage string) pointAIQuestionSuggestion {
	defaultType := PointQuestionTypeReflection
	if pointType == PointTypeExploration {
		defaultType = PointQuestionTypeApplication
	} else if pointType == PointTypeResearch {
		defaultType = PointQuestionTypeGoalAlignment
	}

	suggestion := pointAIQuestionSuggestion{
		QuestionType: string(defaultType),
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		suggestion.Question = fallbackPointAIQuestionText(learningLanguage)
		return suggestion
	}

	var parsed pointAIQuestionSuggestion
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		suggestion.Question = strings.TrimSpace(parsed.Question)
		if normalizedType, ok := normalizePointQuestionType(PointQuestionType(strings.TrimSpace(parsed.QuestionType))); ok {
			suggestion.QuestionType = string(normalizedType)
		}
	}
	if suggestion.Question == "" {
		suggestion.Question = strings.Trim(trimmed, "\" \n\t")
	}
	if suggestion.Question == "" {
		suggestion.Question = fallbackPointAIQuestionText(learningLanguage)
	}
	return suggestion
}

func fallbackPointAIQuestionText(learningLanguage string) string {
	if normalizeLearningLanguage(learningLanguage) == "en" {
		return "How can you turn what you learned at this point into one next action for your current goal?"
	}
	return "이 지점의 내용을 현재 목표에 맞게 다음 행동으로 어떻게 옮길 수 있을까요?"
}

func buildPointAIFeedbackPrompt(goalContext *PlanetGoalContext, point *CoursePointAggregate, lessonTitle string, question CoursePointQuestion, learningLanguage string) string {
	var sb strings.Builder
	sb.WriteString("당신은 정답 판정자가 아니라 학습자가 답변을 작성하거나 보완하도록 돕는 답변 코치입니다.\n")
	sb.WriteString(pointAILanguageInstruction(learningLanguage))
	sb.WriteString("\n아래 질문과, 있을 경우 작성 중인 학습자 답변을 읽고 2~3문장만 작성하세요. 제목, 번호 목록, 마크다운은 쓰지 마세요.\n")
	sb.WriteString("중요: 점수화, 정답/오답 판정, 합격/불합격 표현은 금지합니다. 답변이 없으면 답변을 시작할 핵심 관점과 첫 문장 방향을 제안하세요. 답변이 있으면 잘못 달았거나 덜 완성됐을 가능성을 고려해 강점 하나, 빠진 부분 하나, 다음 보완 방향을 부드럽게 제안하세요.\n\n")
	sb.WriteString(fmt.Sprintf("지점 타입: %s\n", point.Point.PointType))
	sb.WriteString(fmt.Sprintf("지점 제목: %s\n", strings.TrimSpace(point.Point.Title)))
	if strings.TrimSpace(derefStr(point.Point.Description)) != "" {
		sb.WriteString(fmt.Sprintf("지점 설명: %s\n", strings.TrimSpace(derefStr(point.Point.Description))))
	}
	if strings.TrimSpace(lessonTitle) != "" {
		sb.WriteString(fmt.Sprintf("소속 리슨: %s\n", strings.TrimSpace(lessonTitle)))
	}
	if goalContext != nil {
		if strings.TrimSpace(derefStr(goalContext.LearningGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.LearningGoal))))
		} else if strings.TrimSpace(derefStr(goalContext.ConfirmedGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.ConfirmedGoal))))
		}
		if strings.TrimSpace(derefStr(goalContext.UsageContext)) != "" {
			sb.WriteString(fmt.Sprintf("활용 맥락: %s\n", strings.TrimSpace(derefStr(goalContext.UsageContext))))
		}
	}
	sb.WriteString(fmt.Sprintf("질문 유형: %s\n", question.QuestionType))
	sb.WriteString(fmt.Sprintf("질문: %s\n", strings.TrimSpace(question.Question)))
	if strings.TrimSpace(derefStr(question.Answer)) != "" {
		sb.WriteString(fmt.Sprintf("작성 중인 학습자 답변: %s\n", strings.TrimSpace(derefStr(question.Answer))))
	} else {
		sb.WriteString("작성 중인 학습자 답변: 아직 없음\n")
	}
	sb.WriteString("\n피드백 조건:\n")
	sb.WriteString("- 답변이 없으면 첫 문장에 질문을 풀어갈 핵심 관점을 제안하고, 마지막 문장에 학습자가 바로 쓸 수 있는 답변 시작 방향을 제안한다.\n")
	sb.WriteString("- 답변이 있으면 첫 문장에 현재 답변의 의도를 짚고, 마지막 문장에 목표와 연결되는 다음 보완 포인트 또는 고칠 방향을 제안한다.\n")
	sb.WriteString("- 단정적인 평가자 말투 대신 참고 제안 말투를 쓴다.\n")
	return sb.String()
}

func normalizePointAIFeedback(raw string, learningLanguage string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		if normalizeLearningLanguage(learningLanguage) == "en" {
			return "Your answer is starting to connect this point with your current goal. Next, make one situation more specific where you can use this idea right away."
		}
		return "지금 답변에는 현재 목표와 연결하려는 방향이 보입니다. 다음에는 이 지식을 어디에 바로 써볼지 한 가지 상황을 더 구체적으로 적어보면 좋습니다."
	}
	trimmed = strings.Trim(trimmed, "\" \n\t")
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var object map[string]any
		if err := json.Unmarshal([]byte(trimmed), &object); err == nil {
			for _, key := range []string{"feedback", "ai_coaching", "coaching", "content", "message"} {
				if value, ok := object[key].(string); ok && strings.TrimSpace(value) != "" {
					return strings.Trim(strings.TrimSpace(value), "\" \n\t")
				}
			}
		}
		trimmed = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{"), "}"))
	}
	return strings.Trim(trimmed, "\" \n\t")
}

func buildPointSelfEvaluationDraftPrompt(goalContext *PlanetGoalContext, point *CoursePointAggregate, lessonTitle string, applicationAnswers []LearningPointSelfEvaluationApplicationAnswer, learningLanguage string) string {
	var sb strings.Builder
	sb.WriteString("당신은 학습자의 자기평가를 대신 확정하지 않고, 학습자가 검토할 초안을 만드는 학습 코치입니다.\n")
	sb.WriteString(pointAILanguageInstruction(learningLanguage))
	sb.WriteString("\n")
	sb.WriteString("반드시 JSON 한 줄만 출력하세요. 마크다운, 설명, 코드블록은 금지합니다.\n")
	sb.WriteString("형식: {\"application_questions\":[{\"question\":\"...\",\"intent\":\"...\"},{\"question\":\"...\",\"intent\":\"...\"},{\"question\":\"...\",\"intent\":\"...\"}],\"understanding_score\":1-5,\"understanding_reason\":\"...\",\"application_score\":1-5,\"application_reason\":\"...\",\"proficiency_score\":1-5,\"proficiency_reason\":\"...\",\"problem_solving_score\":1-5,\"problem_solving_reason\":\"...\",\"expression_score\":1-5,\"expression_reason\":\"...\",\"goal_alignment_note\":\"...\"}\n")
	sb.WriteString("점수는 근거가 부족하면 3 이하로 보수적으로 제안하고, 이유는 한 문장으로 짧게 씁니다.\n\n")
	sb.WriteString(fmt.Sprintf("지점 타입: %s\n", point.Point.PointType))
	sb.WriteString(fmt.Sprintf("지점 제목: %s\n", strings.TrimSpace(point.Point.Title)))
	if strings.TrimSpace(derefStr(point.Point.Description)) != "" {
		sb.WriteString(fmt.Sprintf("지점 설명: %s\n", strings.TrimSpace(derefStr(point.Point.Description))))
	}
	if strings.TrimSpace(lessonTitle) != "" {
		sb.WriteString(fmt.Sprintf("소속 리슨: %s\n", strings.TrimSpace(lessonTitle)))
	}
	if goalContext != nil {
		if strings.TrimSpace(derefStr(goalContext.LearningGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.LearningGoal))))
		} else if strings.TrimSpace(derefStr(goalContext.ConfirmedGoal)) != "" {
			sb.WriteString(fmt.Sprintf("현재 목표: %s\n", strings.TrimSpace(derefStr(goalContext.ConfirmedGoal))))
		}
		if strings.TrimSpace(derefStr(goalContext.UsageContext)) != "" {
			sb.WriteString(fmt.Sprintf("활용 맥락: %s\n", strings.TrimSpace(derefStr(goalContext.UsageContext))))
		}
		if strings.TrimSpace(derefStr(goalContext.Motivation)) != "" {
			sb.WriteString(fmt.Sprintf("동기: %s\n", strings.TrimSpace(derefStr(goalContext.Motivation))))
		}
	}
	if point.JournalEntry != nil {
		sb.WriteString("내용정리:\n")
		writePromptLine(&sb, "핵심 개념", point.JournalEntry.CoreConcept)
		writePromptLine(&sb, "내 설명", point.JournalEntry.MyExplanation)
		writePromptLine(&sb, "예시", point.JournalEntry.Examples)
		writePromptLine(&sb, "헷갈린 부분", point.JournalEntry.ConfusedParts)
	}
	if point.RecordEntry != nil {
		sb.WriteString(fmt.Sprintf("탐험기록: 학습 %d분, 연습 %d회, 자신감 %d/5, 적용 메모: %s\n",
			point.RecordEntry.StudyMinutes,
			point.RecordEntry.PracticeCount,
			point.RecordEntry.ConfidenceLevel,
			strings.TrimSpace(point.RecordEntry.ApplicationNote),
		))
	}
	if len(point.Questions) > 0 {
		sb.WriteString("질문과 답변:\n")
		for _, question := range point.Questions {
			sb.WriteString(fmt.Sprintf("- Q: %s\n", strings.TrimSpace(question.Question)))
			if strings.TrimSpace(derefStr(question.Answer)) != "" {
				sb.WriteString(fmt.Sprintf("  A: %s\n", strings.TrimSpace(derefStr(question.Answer))))
			}
		}
	}
	if len(point.PracticeLogs) > 0 {
		sb.WriteString("연습/활동기록:\n")
		for _, log := range point.PracticeLogs {
			sb.WriteString(fmt.Sprintf("- %s / 성과: %s / 막힌 부분: %s / 다음: %s\n",
				strings.TrimSpace(log.Title),
				strings.TrimSpace(log.Achievement),
				strings.TrimSpace(log.BlockedPart),
				strings.TrimSpace(log.NextPractice),
			))
		}
	}
	if len(point.Artifacts) > 0 {
		sb.WriteString("결과물:\n")
		for _, artifact := range point.Artifacts {
			sb.WriteString(fmt.Sprintf("- %s / 배운 점: %s / 어려웠던 점: %s\n",
				strings.TrimSpace(artifact.Title),
				strings.TrimSpace(artifact.LearnedPoints),
				strings.TrimSpace(artifact.DifficultPoints),
			))
		}
	}
	if len(point.Blocks) > 0 {
		sb.WriteString(fmt.Sprintf("연구 블록 수: %d\n", len(point.Blocks)))
	}
	if len(applicationAnswers) > 0 {
		sb.WriteString("학습자 적용문제 답변:\n")
		for idx, answer := range applicationAnswers {
			sb.WriteString(fmt.Sprintf("%d. 문제: %s\n", idx+1, strings.TrimSpace(answer.Question)))
			sb.WriteString(fmt.Sprintf("   답변: %s\n", strings.TrimSpace(answer.Answer)))
		}
	}
	sb.WriteString("\n초안 작성 기준:\n")
	sb.WriteString("- application_questions는 현재 지점 내용을 실제 상황에 적용해 보는 짧은 문제 3개로 만든다.\n")
	sb.WriteString("- 학습자가 그대로 저장하기보다 검토하고 수정할 수 있는 초안으로 쓴다.\n")
	sb.WriteString("- 학습자 적용문제 답변이 있으면 해당 답변을 가장 중요한 근거로 평가 초안을 쓴다.\n")
	sb.WriteString("- goal_alignment_note는 현재 목표와 이 지점의 연결을 한두 문장으로 제안한다.\n")
	sb.WriteString("- evidence가 부족한 항목은 과장하지 않는다.\n")
	return sb.String()
}

func parsePointSelfEvaluationDraft(raw string, point *CoursePointAggregate, learningLanguage string) LearningPointSelfEvaluationAIDraft {
	fallback := fallbackPointSelfEvaluationDraft(point, learningLanguage)
	trimmed := extractJSONObject(raw)
	if trimmed == "" {
		return fallback
	}

	var parsed LearningPointSelfEvaluationAIDraft
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return fallback
	}
	parsed.ApplicationQuestions = normalizeSelfEvaluationApplicationQuestions(parsed.ApplicationQuestions, fallback.ApplicationQuestions)
	parsed.UnderstandingScore = clampSelfEvalScore(parsed.UnderstandingScore, fallback.UnderstandingScore)
	parsed.ApplicationScore = clampSelfEvalScore(parsed.ApplicationScore, fallback.ApplicationScore)
	parsed.ProficiencyScore = clampSelfEvalScore(parsed.ProficiencyScore, fallback.ProficiencyScore)
	parsed.ProblemSolvingScore = clampSelfEvalScore(parsed.ProblemSolvingScore, fallback.ProblemSolvingScore)
	parsed.ExpressionScore = clampSelfEvalScore(parsed.ExpressionScore, fallback.ExpressionScore)
	if strings.TrimSpace(parsed.UnderstandingReason) == "" {
		parsed.UnderstandingReason = fallback.UnderstandingReason
	}
	if strings.TrimSpace(parsed.ApplicationReason) == "" {
		parsed.ApplicationReason = fallback.ApplicationReason
	}
	if strings.TrimSpace(parsed.ProficiencyReason) == "" {
		parsed.ProficiencyReason = fallback.ProficiencyReason
	}
	if strings.TrimSpace(parsed.ProblemSolvingReason) == "" {
		parsed.ProblemSolvingReason = fallback.ProblemSolvingReason
	}
	if strings.TrimSpace(parsed.ExpressionReason) == "" {
		parsed.ExpressionReason = fallback.ExpressionReason
	}
	if strings.TrimSpace(parsed.GoalAlignmentNote) == "" {
		parsed.GoalAlignmentNote = fallback.GoalAlignmentNote
	}
	return parsed
}

func fallbackPointSelfEvaluationDraft(point *CoursePointAggregate, learningLanguage string) LearningPointSelfEvaluationAIDraft {
	score := 3
	if point.RecordEntry != nil && point.RecordEntry.ConfidenceLevel >= 1 && point.RecordEntry.ConfidenceLevel <= 5 {
		score = point.RecordEntry.ConfidenceLevel
	}
	answered := 0
	for _, question := range point.Questions {
		if strings.TrimSpace(derefStr(question.Answer)) != "" {
			answered++
		}
	}
	if normalizeLearningLanguage(learningLanguage) == "en" {
		return LearningPointSelfEvaluationAIDraft{
			ApplicationQuestions: []LearningPointSelfEvaluationApplicationQuestion{
				{Question: fmt.Sprintf("How can you apply the key idea from %s to one real situation?", strings.TrimSpace(point.Point.Title)), Intent: "Check practical application"},
				{Question: "If something blocked you during this point, what order would you use next time to solve it?", Intent: "Check problem-solving process"},
				{Question: "How would you explain this point in one sentence to someone learning it for the first time?", Intent: "Check understanding and expression"},
			},
			UnderstandingScore:   score,
			UnderstandingReason:  "This is a conservative estimate based on your saved notes and question answers.",
			ApplicationScore:     score,
			ApplicationReason:    "This draft uses your application notes and practice records as early evidence.",
			ProficiencyScore:     score,
			ProficiencyReason:    "A middle score is appropriate until you have more repeated practice evidence.",
			ProblemSolvingScore:  score,
			ProblemSolvingReason: "This estimate uses your blockers and next practice plan as the main clues.",
			ExpressionScore:      score,
			ExpressionReason:     "Your expression score can increase once your explanation and outputs become more concrete.",
			GoalAlignmentNote:    fmt.Sprintf("Review this point against your current goal. Use your %d answered question(s) and records as evidence, then adjust the final judgment yourself.", answered),
		}
	}
	return LearningPointSelfEvaluationAIDraft{
		ApplicationQuestions: []LearningPointSelfEvaluationApplicationQuestion{
			{Question: fmt.Sprintf("%s에서 배운 핵심을 실제 상황 하나에 어떻게 적용할 수 있을까요?", strings.TrimSpace(point.Point.Title)), Intent: "적용 가능성 확인"},
			{Question: "학습 중 막혔던 부분이 있었다면, 다음에는 어떤 순서로 해결해 볼 수 있을까요?", Intent: "문제해결 과정 확인"},
			{Question: "이 지점을 처음 배우는 사람에게 한 문장으로 설명한다면 어떻게 말할 수 있을까요?", Intent: "이해도와 표현 확인"},
		},
		UnderstandingScore:   score,
		UnderstandingReason:  "저장된 내용정리와 질문 답변을 기준으로 이해도를 보수적으로 추정했습니다.",
		ApplicationScore:     score,
		ApplicationReason:    "적용 메모와 연습/활동기록을 바탕으로 적용 가능성을 임시로 제안했습니다.",
		ProficiencyScore:     score,
		ProficiencyReason:    "반복 연습 기록이 충분히 쌓이기 전까지는 중간 수준으로 두는 것이 적절합니다.",
		ProblemSolvingScore:  score,
		ProblemSolvingReason: "막힌 부분과 다음 연습 계획을 기준으로 문제해결 준비도를 임시 평가했습니다.",
		ExpressionScore:      score,
		ExpressionReason:     "정리한 설명과 결과물이 더 구체화되면 표현 완성도를 높게 조정할 수 있습니다.",
		GoalAlignmentNote:    fmt.Sprintf("이 지점은 현재 목표와 연결해 다시 검토할 수 있습니다. 답변 있는 질문 %d개와 기록을 기준으로 마지막 판단은 직접 조정하세요.", answered),
	}
}

func normalizeSelfEvaluationApplicationQuestions(items []LearningPointSelfEvaluationApplicationQuestion, fallback []LearningPointSelfEvaluationApplicationQuestion) []LearningPointSelfEvaluationApplicationQuestion {
	normalized := make([]LearningPointSelfEvaluationApplicationQuestion, 0, 3)
	for _, item := range items {
		question := strings.TrimSpace(item.Question)
		if question == "" {
			continue
		}
		intent := strings.TrimSpace(item.Intent)
		if intent == "" {
			intent = "적용 확인"
		}
		normalized = append(normalized, LearningPointSelfEvaluationApplicationQuestion{Question: question, Intent: intent})
		if len(normalized) == 3 {
			break
		}
	}
	for _, item := range fallback {
		if len(normalized) == 3 {
			break
		}
		normalized = append(normalized, item)
	}
	return normalized
}

func clampSelfEvalScore(value int, fallback int) int {
	if value < 1 || value > 5 {
		return fallback
	}
	return value
}

func extractJSONObject(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start < 0 || end < start {
		return ""
	}
	return trimmed[start : end+1]
}

func writePromptLine(sb *strings.Builder, label string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	sb.WriteString(fmt.Sprintf("- %s: %s\n", label, value))
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
