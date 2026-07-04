package curriculum

func applyDomainSpecificLessonRefinement(req CreateCourseDraftRequest, doc generatedDraftDocument) generatedDraftDocument {
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		return doc
	}

	subpatternKey := inferGoalSubpatternKey(req)
	if refinement, ok := lookupGoalSubpatternRefinement(subpatternKey); ok && refinement.Refine != nil {
		doc.MainLessons = refinement.Refine(doc.MainLessons)
		doc.CompletionCriteria = ensureCompletionCriterion(doc.CompletionCriteria, refinement.CompletionCaption)
		doc.CompletionCriteria = removeCompletionCriteriaDuplicatedByLessons(doc.CompletionCriteria, doc.MainLessons)
		return doc
	}

	key := inferGoalPatternKey(req)
	refinement, ok := lookupGoalPatternRefinement(key)
	if ok && refinement.Refine != nil {
		doc.MainLessons = refinement.Refine(doc.MainLessons)
		doc.CompletionCriteria = ensureCompletionCriterion(doc.CompletionCriteria, refinement.CompletionCaption)
		doc.CompletionCriteria = removeCompletionCriteriaDuplicatedByLessons(doc.CompletionCriteria, doc.MainLessons)
		return doc
	}
	softRefinement, ok := lookupGoalPatternSoftRefinement(key)
	if ok && softRefinement.Apply != nil {
		doc = softRefinement.Apply(doc)
		doc.CompletionCriteria = ensureCompletionCriterion(doc.CompletionCriteria, softRefinement.CompletionCaption)
		doc.CompletionCriteria = removeCompletionCriteriaDuplicatedByLessons(doc.CompletionCriteria, doc.MainLessons)
	}
	return doc
}

func lookupGoalPatternRefinement(key goalPatternKey) (goalPatternRefinement, bool) {
	if refinement, ok := goalPatternRefinements[key]; ok {
		return refinement, true
	}
	key.GoalModeSecondary = ""
	refinement, ok := goalPatternRefinements[key]
	return refinement, ok
}

func lookupGoalPatternSoftRefinement(key goalPatternKey) (goalPatternSoftRefinement, bool) {
	if refinement, ok := goalPatternSoftRefinements[key]; ok {
		return refinement, true
	}
	key.GoalModeSecondary = ""
	refinement, ok := goalPatternSoftRefinements[key]
	return refinement, ok
}

func lookupGoalSubpatternRefinement(key string) (goalSubpatternRefinement, bool) {
	refinement, ok := goalSubpatternRefinements[key]
	return refinement, ok
}
