package curriculum

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func buildPlanetRecordFeed(planet *PlanetAggregate, routeKind string) *PlanetRecordFeedResponse {
	if routeKind != "shared" {
		routeKind = "learning"
	}
	courseStatus := "learning"
	if routeKind == "shared" || planet.Planet.Status == PlanetStatusCompleted {
		courseStatus = "completed"
	}

	response := &PlanetRecordFeedResponse{
		RouteKind: routeKind,
		Course: PlanetRecordFeedCourse{
			ID:     planet.Planet.ID,
			Title:  planet.Planet.Title,
			Status: courseStatus,
		},
		Cards: []PlanetRecordCard{},
	}

	for _, lesson := range planet.Lessons {
		appendPlanetRecordFeedLesson(response, planet.Planet.ID, lesson, routeKind)
	}
	if courseStatus == "completed" && response.Course.CompletedPoints > 0 {
		response.Cards = append(response.Cards, buildCourseCompletedRecordCard(response.Course, planet.Planet.UpdatedAt, routeKind))
	}

	sort.SliceStable(response.Cards, func(i, j int) bool {
		return response.Cards[i].OccurredAt.After(response.Cards[j].OccurredAt)
	})
	response.Course.RecordCardCount = len(response.Cards)
	for _, card := range response.Cards {
		if card.Visibility == "share_candidate" {
			response.Course.ShareCandidateCount++
		}
	}
	response.Filters = buildPlanetRecordFeedFilters(response.Cards)
	return response
}

func appendPlanetRecordFeedLesson(response *PlanetRecordFeedResponse, courseID uuid.UUID, lesson CourseLessonTree, routeKind string) {
	totalPoints := countLessonPoints(lesson)
	if totalPoints > 0 {
		response.Course.TotalLessons++
	}
	completedPoints := countCompletedLessonPoints(lesson)
	response.Course.TotalPoints += totalPoints
	response.Course.CompletedPoints += completedPoints

	for _, point := range lesson.Points {
		response.Cards = append(response.Cards, buildPointRecordCards(courseID, lesson.Lesson, point, routeKind)...)
	}
	for _, sub := range lesson.SubLessons {
		appendPlanetRecordFeedLesson(response, courseID, sub, routeKind)
	}

	if totalPoints > 0 && completedPoints == totalPoints {
		response.Course.CompletedLessons++
		response.Cards = append(response.Cards, buildLessonCompletedRecordCard(courseID, lesson.Lesson, completedPoints, totalPoints, latestLessonCompletionTime(lesson), routeKind))
	}
}

func buildPointRecordCards(courseID uuid.UUID, lesson CourseLesson, point CoursePointAggregate, routeKind string) []PlanetRecordCard {
	pointID := point.Point.ID
	lessonID := lesson.ID
	actions := recordCardActions(routeKind)
	cards := []PlanetRecordCard{}

	if point.Point.Status == PointStatusCompleted {
		occurredAt := point.Point.UpdatedAt
		if point.Point.CompletedAt != nil {
			occurredAt = *point.Point.CompletedAt
		}
		eventID := latestPointEventID(point.Events, "point_completed")
		cardID := fmt.Sprintf("point_completed:%s:%s:%s", courseID, pointID, eventID)
		cards = append(cards, PlanetRecordCard{
			ID:                  cardID,
			CardType:            "point_completed",
			Title:               fmt.Sprintf("%s 학습콘텐츠를 완료했습니다", point.Point.Title),
			Summary:             recordFeedStringOrFallback(point.Point.Description, "학습콘텐츠 하나를 끝까지 탐험했습니다."),
			Badge:               "학습콘텐츠 완료",
			Category:            "completed",
			Visibility:          "share_candidate",
			CourseID:            courseID,
			LessonID:            &lessonID,
			PointID:             &pointID,
			LessonTitle:         lesson.Title,
			PointTitle:          point.Point.Title,
			PointType:           point.Point.PointType,
			OccurredAt:          occurredAt,
			ShareCandidateScore: shareCandidateScore(point, true, false),
			ShareText:           fmt.Sprintf("LearnWeaver에서 '%s' 학습콘텐츠를 완료했습니다.", point.Point.Title),
			Detail: PlanetRecordCardDetail{
				PrimaryText:           "학습콘텐츠를 완료하며 기록과 자기평가를 남겼습니다.",
				JournalExcerpt:        journalExcerpt(point.JournalEntry),
				QuestionExcerpt:       answeredQuestionExcerpt(point.Questions),
				SelfEvaluationSummary: selfEvaluationSummary(point.SelfEvaluation),
				PracticeSummary:       practiceSummary(point.PracticeLogs),
			},
			Actions: actions,
		})
	}

	for _, artifact := range point.Artifacts {
		cards = append(cards, PlanetRecordCard{
			ID:                  fmt.Sprintf("artifact_submitted:%s:%s:%s", courseID, pointID, artifact.ID),
			CardType:            "artifact_submitted",
			Title:               "공유할 수 있는 결과물을 만들었습니다",
			Summary:             firstRecordFeedText(artifact.Description, artifact.Title, "학습 결과물이 기록되었습니다."),
			Badge:               "결과물 생성",
			Category:            "artifact",
			Visibility:          "share_candidate",
			CourseID:            courseID,
			LessonID:            &lessonID,
			PointID:             &pointID,
			LessonTitle:         lesson.Title,
			PointTitle:          point.Point.Title,
			PointType:           point.Point.PointType,
			OccurredAt:          artifact.CreatedAt,
			ShareCandidateScore: shareCandidateScore(point, false, true),
			ShareText:           fmt.Sprintf("LearnWeaver에서 '%s' 결과물을 만들었습니다.", artifact.Title),
			Detail: PlanetRecordCardDetail{
				PrimaryText:   firstRecordFeedText(artifact.LearnedPoints, artifact.ProductionProcess, artifact.Description),
				ArtifactTitle: artifact.Title,
			},
			Actions: actions,
		})
	}

	for _, question := range point.Questions {
		if strings.TrimSpace(nilToString(question.Answer)) == "" {
			continue
		}
		cards = append(cards, PlanetRecordCard{
			ID:          fmt.Sprintf("question_answered:%s:%s:%s", courseID, pointID, question.ID),
			CardType:    "question_answered",
			Title:       "막혔던 질문에 답을 남겼습니다",
			Summary:     firstRecordFeedText(question.Question, "질문과 답변을 정리했습니다."),
			Badge:       "질문 해결",
			Category:    "work",
			Visibility:  "private",
			CourseID:    courseID,
			LessonID:    &lessonID,
			PointID:     &pointID,
			LessonTitle: lesson.Title,
			PointTitle:  point.Point.Title,
			PointType:   point.Point.PointType,
			OccurredAt:  question.UpdatedAt,
			ShareText:   fmt.Sprintf("LearnWeaver에서 '%s' 질문에 답을 남겼습니다.", point.Point.Title),
			Detail: PlanetRecordCardDetail{
				QuestionExcerpt: firstRecordFeedText(question.Question, nilToString(question.Answer)),
			},
			Actions: actions,
		})
	}

	if point.SelfEvaluation != nil {
		cards = append(cards, PlanetRecordCard{
			ID:                  fmt.Sprintf("self_evaluation_saved:%s:%s", courseID, pointID),
			CardType:            "self_evaluation_saved",
			Title:               "나의 이해도를 점검했습니다",
			Summary:             selfEvaluationSummary(point.SelfEvaluation),
			Badge:               "자기평가",
			Category:            "work",
			Visibility:          "private",
			CourseID:            courseID,
			LessonID:            &lessonID,
			PointID:             &pointID,
			LessonTitle:         lesson.Title,
			PointTitle:          point.Point.Title,
			PointType:           point.Point.PointType,
			OccurredAt:          point.SelfEvaluation.UpdatedAt,
			ShareCandidateScore: shareCandidateScore(point, false, false),
			ShareText:           fmt.Sprintf("LearnWeaver에서 '%s' 학습 후 자기평가를 남겼습니다.", point.Point.Title),
			Detail: PlanetRecordCardDetail{
				SelfEvaluationSummary: selfEvaluationSummary(point.SelfEvaluation),
			},
			Actions: actions,
		})
	}

	for _, log := range point.PracticeLogs {
		cards = append(cards, PlanetRecordCard{
			ID:          fmt.Sprintf("practice_log_added:%s:%s:%s", courseID, pointID, log.ID),
			CardType:    "practice_log_added",
			Title:       "연습의 흔적을 남겼습니다",
			Summary:     firstRecordFeedText(log.Achievement, log.AchievementNote, log.ActivityName, log.Title),
			Badge:       "연습 기록",
			Category:    "work",
			Visibility:  "private",
			CourseID:    courseID,
			LessonID:    &lessonID,
			PointID:     &pointID,
			LessonTitle: lesson.Title,
			PointTitle:  point.Point.Title,
			PointType:   point.Point.PointType,
			OccurredAt:  log.CreatedAt,
			ShareText:   fmt.Sprintf("LearnWeaver에서 '%s' 연습 기록을 남겼습니다.", point.Point.Title),
			Detail: PlanetRecordCardDetail{
				PracticeSummary: practiceLogSummary(log),
			},
			Actions: actions,
		})
	}

	if point.JournalEntry != nil && journalExcerpt(point.JournalEntry) != "" {
		cards = append(cards, PlanetRecordCard{
			ID:          fmt.Sprintf("journal_saved:%s:%s", courseID, pointID),
			CardType:    "journal_saved",
			Title:       "중요한 내용을 내 말로 정리했습니다",
			Summary:     journalExcerpt(point.JournalEntry),
			Badge:       "내용정리",
			Category:    "work",
			Visibility:  "private",
			CourseID:    courseID,
			LessonID:    &lessonID,
			PointID:     &pointID,
			LessonTitle: lesson.Title,
			PointTitle:  point.Point.Title,
			PointType:   point.Point.PointType,
			OccurredAt:  point.JournalEntry.UpdatedAt,
			ShareText:   fmt.Sprintf("LearnWeaver에서 '%s' 내용을 정리했습니다.", point.Point.Title),
			Detail: PlanetRecordCardDetail{
				JournalExcerpt: journalExcerpt(point.JournalEntry),
			},
			Actions: actions,
		})
	}

	return cards
}

func buildLessonCompletedRecordCard(courseID uuid.UUID, lesson CourseLesson, completedPoints, totalPoints int, occurredAt time.Time, routeKind string) PlanetRecordCard {
	lessonID := lesson.ID
	return PlanetRecordCard{
		ID:                  fmt.Sprintf("lesson_completed:%s:%s", courseID, lessonID),
		CardType:            "lesson_completed",
		Title:               fmt.Sprintf("%s 리슨의 학습콘텐츠를 모두 완료했습니다", lesson.Title),
		Summary:             fmt.Sprintf("학습콘텐츠 %d개를 완료해 하나의 리슨을 마쳤습니다.", completedPoints),
		Badge:               "리슨 완료",
		Category:            "completed",
		Visibility:          "share_candidate",
		CourseID:            courseID,
		LessonID:            &lessonID,
		LessonTitle:         lesson.Title,
		OccurredAt:          occurredAt,
		ShareCandidateScore: 2,
		ShareText:           fmt.Sprintf("LearnWeaver에서 '%s' 리슨의 학습콘텐츠를 모두 완료했습니다.", lesson.Title),
		Detail: PlanetRecordCardDetail{
			PrimaryText: "여러 학습콘텐츠 완료가 모여 만든 구간 성취입니다.",
			LessonProgress: &PlanetRecordLessonProgress{
				CompletedPoints: completedPoints,
				TotalPoints:     totalPoints,
				Percent:         percent(completedPoints, totalPoints),
			},
		},
		Actions: recordCardActions(routeKind),
	}
}

func buildCourseCompletedRecordCard(course PlanetRecordFeedCourse, occurredAt time.Time, routeKind string) PlanetRecordCard {
	return PlanetRecordCard{
		ID:                  fmt.Sprintf("course_completed:%s", course.ID),
		CardType:            "course_completed",
		Title:               "하나의 학습 행성을 완주했습니다",
		Summary:             fmt.Sprintf("%s에서 학습콘텐츠 %d개를 완료했습니다.", course.Title, course.CompletedPoints),
		Badge:               "행성 완주",
		Category:            "completed",
		Visibility:          "share_candidate",
		CourseID:            course.ID,
		OccurredAt:          occurredAt,
		ShareCandidateScore: 4,
		ShareText:           fmt.Sprintf("LearnWeaver에서 '%s' 학습 행성을 완주했습니다.", course.Title),
		Detail: PlanetRecordCardDetail{
			PrimaryText: fmt.Sprintf("완료 리슨 %d개, 완료 학습콘텐츠 %d개를 기록했습니다.", course.CompletedLessons, course.CompletedPoints),
		},
		Actions: recordCardActions(routeKind),
	}
}

func buildPlanetRecordFeedFilters(cards []PlanetRecordCard) []PlanetRecordFeedFilter {
	counts := map[string]int{"all": len(cards)}
	for _, card := range cards {
		counts[card.Category]++
		if card.Visibility == "share_candidate" {
			counts["share_candidate"]++
		}
	}
	return []PlanetRecordFeedFilter{
		{Key: "all", Label: "전체", Count: counts["all"]},
		{Key: "completed", Label: "완료", Count: counts["completed"]},
		{Key: "work", Label: "기록", Count: counts["work"]},
		{Key: "artifact", Label: "결과물", Count: counts["artifact"]},
		{Key: "share_candidate", Label: "공유 추천", Count: counts["share_candidate"]},
	}
}

func countLessonPoints(lesson CourseLessonTree) int {
	return len(lesson.Points)
}

func countCompletedLessonPoints(lesson CourseLessonTree) int {
	total := 0
	for _, point := range lesson.Points {
		if point.Point.Status == PointStatusCompleted {
			total++
		}
	}
	return total
}

func latestLessonCompletionTime(lesson CourseLessonTree) time.Time {
	latest := lesson.Lesson.UpdatedAt
	for _, point := range lesson.Points {
		if point.Point.CompletedAt != nil && point.Point.CompletedAt.After(latest) {
			latest = *point.Point.CompletedAt
		}
	}
	return latest
}

func latestPointEventID(events []CoursePointEvent, eventType string) string {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].EventType == eventType {
			return events[i].ID.String()
		}
	}
	return "current"
}

func shareCandidateScore(point CoursePointAggregate, completed, artifact bool) int {
	score := 0
	if completed {
		score += 2
	}
	if artifact || len(point.Artifacts) > 0 {
		score += 3
	}
	if point.SelfEvaluation != nil {
		score++
	}
	if answeredQuestionExcerpt(point.Questions) != "" {
		score++
	}
	if len(point.PracticeLogs) > 0 || point.RecordEntry != nil {
		score++
	}
	return score
}

func recordCardActions(routeKind string) PlanetRecordCardAction {
	return PlanetRecordCardAction{
		CanOpenLearningPage: routeKind == "learning",
		CanCopyShareText:    true,
		CanCreateShareCard:  true,
	}
}

func percent(completed, total int) int {
	if total <= 0 {
		return 0
	}
	return int(float64(completed)/float64(total)*100 + 0.5)
}

func recordFeedStringOrFallback(value *string, fallback string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return fallback
}

func firstRecordFeedText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func nilToString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func journalExcerpt(entry *DraftJournalEntry) string {
	if entry == nil {
		return ""
	}
	return firstRecordFeedText(entry.Reflection, entry.Observation, entry.NextStep)
}

func answeredQuestionExcerpt(questions []CoursePointQuestion) string {
	for _, question := range questions {
		if strings.TrimSpace(nilToString(question.Answer)) != "" {
			return firstRecordFeedText(question.Question, nilToString(question.Answer))
		}
	}
	return ""
}

func selfEvaluationSummary(evaluation *CoursePointSelfEvaluation) string {
	if evaluation == nil {
		return ""
	}
	if evaluation.FinalScore != nil {
		return fmt.Sprintf("최종 자기평가 %d점. %s", *evaluation.FinalScore, firstRecordFeedText(evaluation.GoalAlignmentNote, evaluation.ApplicationNote))
	}
	return fmt.Sprintf("이해도 %d/5, 적용 자신감 %d/5. %s", evaluation.Understanding, evaluation.Proficiency, firstRecordFeedText(evaluation.GoalAlignmentNote, evaluation.ApplicationNote))
}

func practiceSummary(logs []CoursePointPracticeLog) string {
	if len(logs) == 0 {
		return ""
	}
	return practiceLogSummary(logs[0])
}

func practiceLogSummary(log CoursePointPracticeLog) string {
	return firstRecordFeedText(log.Achievement, log.AchievementNote, log.ActivityName, log.Title, log.NextPractice)
}
