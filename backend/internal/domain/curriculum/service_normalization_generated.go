package curriculum

import (
	"fmt"
	"strings"
)

func normalizeGeneratedMainLessons(lessons []generatedMainLesson) []generatedMainLesson {
	result := make([]generatedMainLesson, 0, len(lessons))
	seen := make(map[string]struct{}, len(lessons))
	for _, lesson := range lessons {
		title := strings.TrimSpace(lesson.Title)
		objective := strings.TrimSpace(lesson.Objective)
		if title == "" {
			continue
		}
		key := strings.ToLower(title)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, generatedMainLesson{
			Title:                    title,
			Objective:                objective,
			RecommendationSearchSpec: lesson.RecommendationSearchSpec,
		})
	}
	return result
}

func normalizeGeneratedDraftDocument(doc generatedDraftDocument) generatedDraftDocument {
	doc.Title = strings.TrimSpace(doc.Title)
	doc.Description = strings.TrimSpace(doc.Description)
	if len(doc.MainLessons) == 0 {
		doc.MainLessons = doc.Lessons
	}
	doc.MainLessons = normalizeGeneratedMainLessons(doc.MainLessons)
	doc.Lessons = nil
	doc.CompletionCriteria = normalizeCompletionCriteria(doc.CompletionCriteria)
	return doc
}

func refineGeneratedDraftDocument(req CreateCourseDraftRequest, doc generatedDraftDocument) generatedDraftDocument {
	doc = normalizeGeneratedDraftDocument(doc)
	doc = ensureGeneratedDraftLanguage(req, doc)
	if len(doc.MainLessons) == 0 {
		return doc
	}

	doc = broadenAtomicLessons(req, doc)
	doc = moveMilestoneLessonsToCompletionCriteria(doc)
	doc = mergeAtomicLessons(doc)
	doc = dedupeNearLessons(doc)
	doc = ensurePreMilestoneLastLesson(req, doc)
	doc = applyDomainSpecificLessonRefinement(req, doc)
	doc.MainLessons = clampLessons(doc.MainLessons)
	doc.MainLessons = cleanupLessonObjectives(doc.MainLessons)
	doc.CompletionCriteria = normalizeCompletionCriteria(doc.CompletionCriteria)
	doc = normalizeGeneratedDraftDocument(doc)
	return ensureGeneratedDraftLanguage(req, doc)
}

func ensureGeneratedDraftLanguage(req CreateCourseDraftRequest, doc generatedDraftDocument) generatedDraftDocument {
	if normalizeLearningLanguage(req.LearningLanguage) != "en" {
		return doc
	}
	fallbackSubject := strings.TrimSpace(req.SourceQuery)
	if fallbackSubject == "" {
		fallbackSubject = strings.TrimSpace(derefString(req.LearningGoal))
	}
	if fallbackSubject == "" {
		fallbackSubject = "Learning Goal"
	}
	if containsHangul(doc.Title) {
		doc.Title = strings.TrimSpace(fallbackSubject + " Learning Course")
	}
	if containsHangul(doc.Description) {
		doc.Description = fmt.Sprintf("A structured learning path for %s.", fallbackSubject)
	}
	return doc
}

func applyCurriculumReview(base generatedDraftDocument, review reviewedDraftDocument) generatedDraftDocument {
	revisedLessons := normalizeGeneratedMainLessons(review.RevisedMainLessons)
	if len(revisedLessons) >= 3 && len(revisedLessons) <= 8 {
		base.MainLessons = revisedLessons
	}
	revisedCriteria := normalizeCompletionCriteria(review.RevisedCompletionCriteria)
	if len(revisedCriteria) > 0 {
		base.CompletionCriteria = revisedCriteria
	}
	return normalizeGeneratedDraftDocument(base)
}

func clampLessons(lessons []generatedMainLesson) []generatedMainLesson {
	if len(lessons) <= 8 {
		return lessons
	}
	return append([]generatedMainLesson(nil), lessons[:8]...)
}

func cleanupLessonObjectives(lessons []generatedMainLesson) []generatedMainLesson {
	result := make([]generatedMainLesson, 0, len(lessons))
	for _, lesson := range lessons {
		lesson.Objective = dedupeObjectiveSentences(lesson.Objective)
		result = append(result, lesson)
	}
	return result
}

func mergeLessonDetails(target, source generatedMainLesson) generatedMainLesson {
	target.Title = strings.TrimSpace(target.Title)
	target.Objective = strings.TrimSpace(target.Objective)
	source.Objective = strings.TrimSpace(source.Objective)
	if target.Objective == "" {
		target.Objective = source.Objective
		return target
	}
	if source.Objective == "" || strings.Contains(target.Objective, source.Objective) {
		return target
	}
	target.Objective = target.Objective + " 필요 시 " + source.Objective
	return target
}

func lessonsOverlap(a, b generatedMainLesson) bool {
	aTitle := strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(a.Title)), ""))
	bTitle := strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(b.Title)), ""))
	if aTitle == "" || bTitle == "" {
		return false
	}
	return aTitle == bTitle || strings.Contains(aTitle, bTitle) || strings.Contains(bTitle, aTitle)
}

func ensureCompletionCriterion(criteria []string, criterion string) []string {
	normalized := normalizeCompletionCriteria(criteria)
	target := strings.TrimSpace(criterion)
	if target == "" {
		return normalized
	}
	for _, item := range normalized {
		if strings.EqualFold(strings.TrimSpace(item), target) {
			return normalized
		}
	}
	return append(normalized, target)
}

func removeCompletionCriteriaDuplicatedByLessons(criteria []string, lessons []generatedMainLesson) []string {
	normalized := normalizeCompletionCriteria(criteria)
	if len(normalized) == 0 || len(lessons) == 0 {
		return normalized
	}

	lessonTexts := make(map[string]struct{}, len(lessons)*2)
	for _, lesson := range lessons {
		title := strings.TrimSpace(lesson.Title)
		if title != "" {
			lessonTexts[strings.ToLower(title)] = struct{}{}
		}
		objective := strings.TrimSpace(lesson.Objective)
		if objective != "" {
			lessonTexts[strings.ToLower(objective)] = struct{}{}
		}
	}

	filtered := make([]string, 0, len(normalized))
	for _, item := range normalized {
		if _, exists := lessonTexts[strings.ToLower(strings.TrimSpace(item))]; exists {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
