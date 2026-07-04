package curriculum

import (
	"fmt"
	"strings"
)

func moveMilestoneLessonsToCompletionCriteria(doc generatedDraftDocument) generatedDraftDocument {
	filtered := make([]generatedMainLesson, 0, len(doc.MainLessons))
	for _, lesson := range doc.MainLessons {
		if isMilestoneLikeLesson(lesson) {
			doc.CompletionCriteria = append(doc.CompletionCriteria, buildCriterionFromLesson(lesson))
			continue
		}
		filtered = append(filtered, lesson)
	}
	doc.MainLessons = filtered
	return doc
}

func broadenAtomicLessons(req CreateCourseDraftRequest, doc generatedDraftDocument) generatedDraftDocument {
	goal := strings.TrimSpace(derefString(req.LearningGoal))
	for idx, lesson := range doc.MainLessons {
		if !isAtomicLesson(lesson) {
			continue
		}
		doc.MainLessons[idx] = generatedMainLesson{
			Title:                    buildIntroStageTitle(req.SourceQuery, goal, req.LearningLanguage),
			Objective:                buildIntroStageObjective(req.SourceQuery, goal, lesson.Objective, req.LearningLanguage),
			RecommendationSearchSpec: lesson.RecommendationSearchSpec,
		}
	}
	return doc
}

func mergeAtomicLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) <= 3 {
		return doc
	}

	lessons := append([]generatedMainLesson(nil), doc.MainLessons...)
	for idx := 0; idx < len(lessons); idx++ {
		if !isAtomicLesson(lessons[idx]) {
			continue
		}
		targetIdx := idx + 1
		if targetIdx >= len(lessons) {
			targetIdx = idx - 1
		}
		if targetIdx < 0 || targetIdx == idx {
			continue
		}
		lessons[targetIdx] = mergeLessonDetails(lessons[targetIdx], lessons[idx])
		lessons = append(lessons[:idx], lessons[idx+1:]...)
		idx--
		if len(lessons) <= 3 {
			break
		}
	}
	doc.MainLessons = lessons
	return doc
}

func dedupeNearLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) <= 1 {
		return doc
	}
	result := make([]generatedMainLesson, 0, len(doc.MainLessons))
	for _, lesson := range doc.MainLessons {
		if len(result) == 0 {
			result = append(result, lesson)
			continue
		}
		last := result[len(result)-1]
		if lessonsOverlap(last, lesson) {
			result[len(result)-1] = mergeLessonDetails(last, lesson)
			continue
		}
		result = append(result, lesson)
	}
	doc.MainLessons = result
	return doc
}

func ensurePreMilestoneLastLesson(req CreateCourseDraftRequest, doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) == 0 {
		return doc
	}
	lastIdx := len(doc.MainLessons) - 1
	last := doc.MainLessons[lastIdx]
	if isMilestoneLikeLesson(last) {
		doc.CompletionCriteria = append(doc.CompletionCriteria, buildCriterionFromLesson(last))
		doc.MainLessons = doc.MainLessons[:lastIdx]
		lastIdx--
	}
	if lastIdx < 0 {
		return doc
	}
	last = doc.MainLessons[lastIdx]
	if !isPreMilestoneLikeLesson(last) {
		goal := strings.TrimSpace(derefString(req.LearningGoal))
		doc.MainLessons[lastIdx] = generatedMainLesson{
			Title:     buildPreMilestoneTitle(req.SourceQuery, goal, req.LearningLanguage),
			Objective: buildPreMilestoneObjective(goal, req.LearningLanguage),
		}
	}
	return doc
}

func isAtomicLesson(lesson generatedMainLesson) bool {
	title := strings.ToLower(strings.TrimSpace(lesson.Title))
	objective := strings.ToLower(strings.TrimSpace(lesson.Objective))
	text := strings.TrimSpace(title + " " + objective)
	if text == "" {
		return false
	}
	if isMilestoneLikeLesson(lesson) || isPreMilestoneLikeLesson(lesson) {
		return false
	}

	// Atomic lesson은 "도구/재료/개념 하나만 소개하거나 이해하는 미시 단계"일 때만 잡는다.
	// 넓은 첫 단계(예: 재료+도구로 첫 결과물 만들기)까지 generic intro로 뭉개지지 않게 제한한다.
	microTitlePhrases := []string{
		"기본자세 익히기",
		"기본 자세 익히기",
		"자세 익히기",
		"튜닝하기",
		"튜닝 익히기",
		"재료 소개",
		"도구 소개",
		"장비 소개",
		"장비 이해",
		"기본 개념 이해하기",
		"기초 개념 이해하기",
		"이론 알아보기",
		"기본 잡기",
	}
	hasMicroTitle := false
	for _, phrase := range microTitlePhrases {
		if strings.Contains(title, phrase) {
			hasMicroTitle = true
			break
		}
	}
	if !hasMicroTitle {
		return false
	}

	broaderOutcomeTerms := []string{
		"완성", "만든다", "이어", "연결", "적용", "준비", "정리", "시작",
		"흐름", "결과물", "작품", "반주", "대화", "구도", "표현", "패턴", "형태",
	}
	for _, term := range broaderOutcomeTerms {
		if strings.Contains(text, term) {
			return false
		}
	}
	return true
}

func isMilestoneLikeLesson(lesson generatedMainLesson) bool {
	text := strings.ToLower(strings.TrimSpace(lesson.Title + " " + lesson.Objective))
	strongMilestonePhrases := []string{
		"실전 투입",
		"실제 참여",
		"공식 참여",
		"실제 연주",
		"공식 반주",
		"무대에 선",
		"무대에서 연주",
		"발표한다",
		"발표하기",
		"시험을 본",
		"시험 응시",
		"합격한다",
		"합격하기",
		"데뷔한다",
		"현업에 투입",
		"현업에서 바로",
		"찬양단에서 실제",
		"모임에서 바로 참여",
		"졸업식에서 발표",
	}
	for _, term := range strongMilestonePhrases {
		if strings.Contains(text, term) {
			return true
		}
	}
	return false
}

func isPreMilestoneLikeLesson(lesson generatedMainLesson) bool {
	text := strings.ToLower(strings.TrimSpace(lesson.Title + " " + lesson.Objective))
	for _, term := range []string{"준비", "적용", "리허설", "통합", "완주", "점검", "정리", "연습", "완성"} {
		if strings.Contains(text, term) {
			return true
		}
	}
	return false
}

func buildCriterionFromLesson(lesson generatedMainLesson) string {
	objective := strings.TrimSpace(lesson.Objective)
	if objective != "" {
		return objective
	}
	return strings.TrimSpace(lesson.Title)
}

func buildPreMilestoneTitle(sourceQuery, learningGoal, language string) string {
	context := strings.ToLower(strings.TrimSpace(sourceQuery + " " + learningGoal))
	if normalizeLearningLanguage(language) == "en" {
		switch {
		case containsAny(context, "sourdough", "starter bread", "simple loaf", "fermentation"):
			return "Bake and review the first sourdough loaf"
		case containsAny(context, "meal prep", "healthy lunch", "weekday meal"):
			return "Review the weekday meal prep routine"
		case containsAny(context, "indoor herb", "herb garden", "basil", "mint", "grow herbs"):
			return "Review the herb care and harvest routine"
		case containsAny(context, "houseplant", "repot", "repotting", "soil mix"):
			return "Monitor the repotted plant recovery"
		case containsAny(context, "chess", "opening", "middlegame"):
			return "Play and review the opening practice game"
		case containsAny(context, "personal finance", "monthly budget", "spending categories", "expenses"):
			return "Review the first weekly budget check"
		case containsAny(context, "mindfulness", "breathing routine", "focus"):
			return "Review the daily breathing practice log"
		case strings.Contains(context, "dog") && containsAny(context, "sit", "stay", "recall", "positive reinforcement"):
			return "Practice the commands in a short real session"
		case containsAny(context, "notion", "study dashboard", "weekly review"):
			return "Use the study dashboard for a weekly review"
		case containsEnglishSpeakingInterviewCertificationContext(context):
			return "Prepare final speaking interview drills"
		case containsLanguageCertificationContext(context) || containsCloudCertificationContext(context) || containsTheoryCertificationContext(context):
			return "Review the final exam practice flow"
		case containsAny(context, "drone", "drone operator", "drone license", "basic drone"):
			return "Review safe flight and license readiness"
		case containsPracticalCertificationContext(context):
			return "Check the practical exam workflow"
		case containsCreatorTutorialPublishContext(context):
			return "Prepare the tutorial for recording or publishing"
		case containsMIDIProductionContext(context):
			return "Prepare the final track export"
		case containsCalligraphyContext(context):
			return "Refine the quote card before completion"
		case containsLeathercraftContext(context):
			return "Finish the leather project workflow"
		case containsWoodworkingSmallProjectContext(context):
			if containsAny(context, "birdhouse", "bird house") {
				return "Inspect and finish the assembled birdhouse"
			}
			return "Finish and inspect the small wood project"
		case containsCrochetSmallDollProjectContext(context):
			if containsAny(context, "granny square", "coaster", "pouch") {
				return "Join and finish the crochet square project"
			}
			return "Finish and inspect the crochet project"
		case containsPotteryHandbuildingCupProjectContext(context):
			if containsAny(context, "glaze", "glazing", "bowl") {
				return "Review the glazed bowl before firing"
			}
			return "Finish and inspect the handbuilt pottery piece"
		case containsSewingProjectFoundationContext(context):
			return "Finish and inspect the sewn tote bag"
		case containsAny(context, "arduino", "sensor demo", "maker workshop", "prototype demo"):
			return "Prepare to demonstrate the working prototype"
		case containsCardboardCircuitInventionPathContext(context):
			return "Test and present the working circuit prototype"
		case containsDesignLabPrototypeProjectContext(context):
			return "Refine the prototype with user feedback"
		case containsGame3DGuidedPathwayContext(context):
			return "Playtest the guided 3D scene"
		case containsVibeCodingContext(context) && containsAny(context, "claude code", "mcp", "hook", "hooks", "agentic"):
			return "Run the Claude Code feature workflow end to end"
		case containsClimateDataGraphingExplanationContext(context):
			return "Prepare to explain the data trend clearly"
		case containsCitizenScienceObservationRecordContext(context):
			return "Review and share the biodiversity observation log"
		case containsPrimarySourceInquiryNoteContext(context):
			return "Complete the observe-reflect-question note"
		case containsShortCourseWeeklyDiscussionContext(context):
			return "Prepare the final weekly discussion contribution"
		case containsOrganicGrowingCyclePlanContext(context):
			return "Review the seed-to-harvest care cycle"
		case containsAIDigitalLiteracyLifePracticeContext(context):
			return "Finalize a safe daily AI workflow"
		case containsAny(context, "essay series", "weekly essay", "publish a short essay", "online publication"):
			return "Prepare the next essay for publication"
		case containsAny(context, "podcast", "episode", "audio"):
			return "Prepare the first podcast episode for publishing"
		case containsAny(context, "shorts", "short video", "vertical video", "captions", "thumbnail"):
			return "Prepare the short video for publishing"
		case containsAny(context, "workout", "exercise", "fitness", "home workout"):
			return "Review the weekly workout routine safely"
		case strings.Contains(context, "yoga") || strings.Contains(context, "flexibility") || strings.Contains(context, "breathing routine") || strings.Contains(context, "morning routine"):
			return "Finalize a repeatable short yoga routine"
		case strings.Contains(context, "dance") || strings.Contains(context, "choreography") || strings.Contains(context, "k-pop") || strings.Contains(context, "cover routine"):
			return "Prepare to record the target dance cover"
		case containsAny(context, "ukulele", "first song", "simple song"):
			return "Play the first song from start to finish"
		case containsAny(context, "harmonica", "blues riff"):
			return "Play the blues riff with clean notes and bends"
		case containsAny(context, "bass", "groove"):
			return "Play the target bass groove steadily"
		case strings.Contains(context, "guitar") || strings.Contains(context, "piano") || strings.Contains(context, "drum") || strings.Contains(context, "violin") || strings.Contains(context, "cello") || strings.Contains(context, "performance") || strings.Contains(context, "play "):
			return "Prepare to play through the target performance"
		case strings.Contains(context, "english") || strings.Contains(context, "japanese") || strings.Contains(context, "conversation") || strings.Contains(context, "language"):
			return "Prepare to sustain the target conversation"
		case containsAny(context, "watercolor", "watercolour", "landscape postcard"):
			return "Refine the watercolor postcard before completion"
		case strings.Contains(context, "drawing") || strings.Contains(context, "watercolor") || strings.Contains(context, "sketch") || strings.Contains(context, "painting") || strings.Contains(context, "photo"):
			return "Refine the target output before completion"
		default:
			return "Complete the pre-goal practice flow"
		}
	}
	switch {
	case containsLanguageCertificationContext(context) || containsCloudCertificationContext(context) || containsTheoryCertificationContext(context):
		return "시험 직전 문제 흐름을 정리한다"
	case containsPracticalCertificationContext(context):
		return "실기 시험 직전 작업 흐름을 점검한다"
	case strings.Contains(context, "기타") || strings.Contains(context, "피아노") || strings.Contains(context, "드럼") || strings.Contains(context, "연주"):
		return "목표 연주를 끝까지 이어갈 준비를 한다"
	case strings.Contains(context, "영어") || strings.Contains(context, "일본어") || strings.Contains(context, "회화") || strings.Contains(context, "언어"):
		return "목표 상황 대화를 이어갈 준비를 한다"
	case strings.Contains(context, "그림") || strings.Contains(context, "수채화") || strings.Contains(context, "드로잉") || strings.Contains(context, "스케치"):
		return "목표 결과물을 완성 직전까지 다듬는다"
	default:
		return "목표 직전 수행 흐름을 완성한다"
	}
}

func buildPreMilestoneObjective(learningGoal, language string) string {
	if normalizeLearningLanguage(language) == "en" {
		if strings.TrimSpace(learningGoal) == "" {
			return "Connect the preparation, application, and review steps needed before the final goal."
		}
		return fmt.Sprintf("Connect the preparation, application, and review steps needed before '%s'.", learningGoal)
	}
	if strings.TrimSpace(learningGoal) == "" {
		return "최종 목표 직전 단계까지 필요한 준비, 적용, 점검 흐름을 스스로 이어갈 수 있게 한다."
	}
	return fmt.Sprintf("'%s' 목표 직전 단계까지 필요한 준비, 적용, 점검 흐름을 스스로 이어갈 수 있게 한다.", learningGoal)
}
