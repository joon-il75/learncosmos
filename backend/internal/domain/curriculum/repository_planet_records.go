package curriculum

func buildPlanetRecordAggregate(planet *PlanetAggregate) *PlanetRecordAggregate {
	result := &PlanetRecordAggregate{
		Course: PlanetRecordCourse{
			ID:    planet.Planet.ID,
			Title: planet.Planet.Title,
		},
		Lessons: []PlanetRecordLesson{},
	}

	for _, lesson := range planet.Lessons {
		appendRecordLesson(result, lesson)
	}
	return result
}

func appendRecordLesson(result *PlanetRecordAggregate, lesson CourseLessonTree) {
	recordLesson := PlanetRecordLesson{
		LessonID:    lesson.Lesson.ID,
		LessonTitle: lesson.Lesson.Title,
		Points:      []PlanetRecordPoint{},
	}
	for _, point := range lesson.Points {
		result.Course.TotalPoints++
		recordPoint := PlanetRecordPoint{
			PointID:      point.Point.ID,
			PointTitle:   point.Point.Title,
			PointType:    point.Point.PointType,
			Blocks:       point.Blocks,
			JournalEntry: point.JournalEntry,
			Questions:    point.Questions,
			RecordEntry:  point.RecordEntry,
		}
		if point.SelfEvaluation != nil {
			recordPoint.SelfEvaluation = point.SelfEvaluation
			recordPoint.ApplicationNote = point.SelfEvaluation.ApplicationNote
			recordPoint.SelfEvaluationApplicationNote = point.SelfEvaluation.ApplicationNote
			recordPoint.GoalAlignmentNote = point.SelfEvaluation.GoalAlignmentNote
		}
		if isPlanetRecordPointRecorded(recordPoint) {
			result.Course.RecordedPoints++
		}
		recordLesson.Points = append(recordLesson.Points, recordPoint)
	}
	if len(recordLesson.Points) > 0 {
		result.Lessons = append(result.Lessons, recordLesson)
	}
	for _, sub := range lesson.SubLessons {
		appendRecordLesson(result, sub)
	}
}

func isPlanetRecordPointRecorded(point PlanetRecordPoint) bool {
	if point.JournalEntry != nil {
		return true
	}
	if point.RecordEntry != nil {
		return true
	}
	if len(point.Questions) > 0 {
		return true
	}
	if point.SelfEvaluation != nil {
		return true
	}
	if len(point.Blocks) > 0 {
		return true
	}
	return false
}

func buildPlanetResultAggregate(planet *PlanetAggregate) *PlanetResultAggregate {
	result := &PlanetResultAggregate{
		Course: PlanetResultCourse{
			ID:    planet.Planet.ID,
			Title: planet.Planet.Title,
		},
		Lessons: []PlanetResultLesson{},
	}

	for _, lesson := range planet.Lessons {
		appendResultLesson(result, lesson)
	}
	return result
}

func appendResultLesson(result *PlanetResultAggregate, lesson CourseLessonTree) {
	resultLesson := PlanetResultLesson{
		LessonID:    lesson.Lesson.ID,
		LessonTitle: lesson.Lesson.Title,
		Points:      []PlanetResultPoint{},
	}
	for _, point := range lesson.Points {
		result.Course.TotalPoints++
		artifacts := point.Artifacts
		if artifacts == nil {
			artifacts = []CoursePointArtifact{}
		}
		if len(artifacts) > 0 {
			result.Course.PointsWithArtifacts++
		}
		resultLesson.Points = append(resultLesson.Points, PlanetResultPoint{
			PointID:    point.Point.ID,
			PointTitle: point.Point.Title,
			PointType:  point.Point.PointType,
			Artifacts:  artifacts,
		})
	}
	if len(resultLesson.Points) > 0 {
		result.Lessons = append(result.Lessons, resultLesson)
	}
	for _, sub := range lesson.SubLessons {
		appendResultLesson(result, sub)
	}
}
