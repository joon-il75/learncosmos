package curriculum

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/pkg/urlsafe"
)

const (
	maxLearningResearchBlockCount              = 3
	maxLearningResearchMaterialAttachmentCount = 10
	researchMaterialImageMaxBytes              = int64(5 * 1024 * 1024)
	researchMaterialDocumentMaxBytes           = int64(20 * 1024 * 1024)
	researchMaterialVideoMaxBytes              = int64(500 * 1024 * 1024)
	researchMaterialSubtitleMaxBytes           = int64(2 * 1024 * 1024)
)

var researchMaterialAllowedImageExtensions = map[string]struct{}{
	".gif":  {},
	".jpeg": {},
	".jpg":  {},
	".png":  {},
	".webp": {},
}

var researchMaterialAllowedDocumentExtensions = map[string]struct{}{
	".csv":  {},
	".doc":  {},
	".docx": {},
	".hwp":  {},
	".hwpx": {},
	".md":   {},
	".odt":  {},
	".pdf":  {},
	".ppt":  {},
	".pptx": {},
	".rtf":  {},
	".txt":  {},
	".xls":  {},
	".xlsx": {},
	".zip":  {},
}

var researchMaterialAllowedVideoExtensions = map[string]struct{}{
	".mov":  {},
	".mp4":  {},
	".webm": {},
}

var researchMaterialAllowedSubtitleExtensions = map[string]struct{}{
	".srt": {},
	".vtt": {},
}

func attachCoursePointEntry(mainTrees []CourseLessonTree, pointID uuid.UUID, attach func(point *CoursePointAggregate)) {
	for mainIdx := range mainTrees {
		for pointIdx := range mainTrees[mainIdx].Points {
			if mainTrees[mainIdx].Points[pointIdx].Point.ID == pointID {
				attach(&mainTrees[mainIdx].Points[pointIdx])
				return
			}
		}
		for subIdx := range mainTrees[mainIdx].SubLessons {
			for pointIdx := range mainTrees[mainIdx].SubLessons[subIdx].Points {
				if mainTrees[mainIdx].SubLessons[subIdx].Points[pointIdx].Point.ID == pointID {
					attach(&mainTrees[mainIdx].SubLessons[subIdx].Points[pointIdx])
					return
				}
			}
		}
	}
}

func flattenPlanetPointDetails(lessons []CourseLessonTree, levelTitle string) []PlanetPointDetail {
	details := make([]PlanetPointDetail, 0)
	for _, lesson := range lessons {
		for _, point := range lesson.Points {
			details = append(details, PlanetPointDetail{
				LevelTitle:  levelTitle,
				LessonTitle: lesson.Lesson.Title,
				LessonID:    lesson.Lesson.ID,
				Point:       point,
			})
		}
		if len(lesson.SubLessons) > 0 {
			details = append(details, flattenPlanetPointDetails(lesson.SubLessons, levelTitle)...)
		}
	}
	return details
}

func resolvePlanetPointDetail(planet *PlanetAggregate, pointID uuid.UUID) (*PlanetPointDetail, bool) {
	sequence := make([]PlanetPointDetail, 0)
	for _, lesson := range planet.Lessons {
		sequence = append(sequence, flattenPlanetPointDetails([]CourseLessonTree{lesson}, lesson.Lesson.Title)...)
	}

	for idx := range sequence {
		if sequence[idx].Point.Point.ID != pointID {
			continue
		}

		detail := sequence[idx]
		detail.Planet = planet.Planet
		detail.GoalContext = planet.GoalContext
		if idx > 0 {
			prevID := sequence[idx-1].Point.Point.ID
			detail.PreviousPointID = &prevID
		}
		if idx < len(sequence)-1 {
			nextID := sequence[idx+1].Point.Point.ID
			detail.NextPointID = &nextID
		}
		return &detail, true
	}

	return nil, false
}

func applyResearchMaterialStatusToLessonTrees(lessons []CourseLessonTree) {
	for lessonIdx := range lessons {
		for pointIdx := range lessons[lessonIdx].Points {
			applyResearchMaterialStatus(&lessons[lessonIdx].Points[pointIdx])
		}
		applyResearchMaterialStatusToLessonTrees(lessons[lessonIdx].SubLessons)
	}
}

func applyResearchMaterialStatus(point *CoursePointAggregate) {
	if point == nil || point.Point.PointType != PointTypeResearch {
		return
	}

	materialCount := len(point.Blocks)
	for _, attachment := range point.Attachments {
		if strings.TrimSpace(attachment.SourceContext) != "research_material" {
			continue
		}
		materialCount++
	}
	if materialCount == 0 {
		return
	}

	var latestConfirmationEvent *CoursePointEvent
	for idx := range point.Events {
		eventType := strings.TrimSpace(point.Events[idx].EventType)
		if eventType != "research_material_confirmed" && eventType != "research_material_unconfirmed" {
			continue
		}
		if latestConfirmationEvent == nil || point.Events[idx].CreatedAt.After(latestConfirmationEvent.CreatedAt) {
			latestConfirmationEvent = &point.Events[idx]
		}
	}

	if latestConfirmationEvent != nil && strings.TrimSpace(latestConfirmationEvent.EventType) == "research_material_confirmed" {
		confirmedAt := latestConfirmationEvent.CreatedAt
		point.ResearchMaterialConfirmed = true
		point.ResearchMaterialConfirmedAt = &confirmedAt
	}
}

func normalizeResearchNodeBlockType(blockType ResearchNodeBlockType) (ResearchNodeBlockType, bool) {
	switch blockType {
	case ResearchNodeBlockTypeText, ResearchNodeBlockTypeImage, ResearchNodeBlockTypeLink:
		return blockType, true
	default:
		return "", false
	}
}

func normalizePointQuestionType(questionType PointQuestionType) (PointQuestionType, bool) {
	switch questionType {
	case "", PointQuestionTypeReflection:
		return PointQuestionTypeReflection, true
	case PointQuestionTypeApplication, PointQuestionTypeGoalAlignment:
		return questionType, true
	default:
		return "", false
	}
}

func normalizePointQuestionStatus(status PointQuestionStatus) (PointQuestionStatus, bool) {
	switch status {
	case PointQuestionStatusPending, PointQuestionStatusAnswered:
		return status, true
	default:
		return "", false
	}
}

func normalizePointArtifactInput(artifactType, title, url, description string) (string, string, string, string, error) {
	artifactType = strings.TrimSpace(artifactType)
	title = strings.TrimSpace(title)
	url = strings.TrimSpace(url)
	description = sanitizePointArtifactDescription(description)
	if title == "" {
		return "", "", "", "", errLearningPointArtifactInvalidInput
	}
	if artifactType == "" {
		artifactType = "note"
	}
	if url != "" {
		normalizedURL, ok := urlsafe.NormalizeHTTPURL(url)
		if !ok {
			return "", "", "", "", errLearningPointArtifactInvalidInput
		}
		url = normalizedURL
	}
	return artifactType, title, url, description, nil
}

func normalizePointArtifactExtraInput(pointCategoryValue, productionProcessValue, learnedPointsValue, difficultPointsValue, visibilityValue *string) (string, string, string, string, string, error) {
	pointCategory := trimOptionalString(pointCategoryValue)
	productionProcess := trimOptionalString(productionProcessValue)
	learnedPoints := trimOptionalString(learnedPointsValue)
	difficultPoints := trimOptionalString(difficultPointsValue)
	visibility := trimOptionalString(visibilityValue)
	if visibility == "" {
		visibility = "private"
	}
	switch visibility {
	case "private", "community", "portfolio":
	default:
		return "", "", "", "", "", errLearningPointArtifactInvalidInput
	}
	return pointCategory, productionProcess, learnedPoints, difficultPoints, visibility, nil
}

func normalizePointAttachmentInput(providerValue, attachmentTypeValue, sourceContextValue, titleValue, urlValue, filePathValue, mimeTypeValue *string, fileSize *int64) (string, string, string, string, string, string, *int64, string, error) {
	provider := trimOptionalString(providerValue)
	if provider == "" {
		provider = "learner"
	}
	switch provider {
	case "creator", "platform", "learner":
	default:
		return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
	}

	attachmentType := trimOptionalString(attachmentTypeValue)
	if attachmentType == "" {
		attachmentType = "link"
	}
	switch attachmentType {
	case "image", "file", "link", "code", "other", "video", "subtitle", "thumbnail":
	default:
		return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
	}

	sourceContext := trimOptionalString(sourceContextValue)
	if sourceContext == "" {
		sourceContext = "work_attachment"
	}
	switch sourceContext {
	case "research_material", "work_attachment", "artifact":
	default:
		return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
	}

	title := trimOptionalString(titleValue)
	url := trimOptionalString(urlValue)
	filePath := trimOptionalString(filePathValue)
	mimeType := trimOptionalString(mimeTypeValue)
	if title == "" || (url == "" && filePath == "") {
		return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
	}
	if url != "" {
		normalizedURL, ok := urlsafe.NormalizeHTTPURL(url)
		if !ok {
			return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
		}
		url = normalizedURL
	}
	if fileSize != nil && *fileSize < 0 {
		return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
	}
	if sourceContext == "research_material" || sourceContext == "artifact" {
		if provider != "learner" {
			return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
		}
		if attachmentType != "image" && attachmentType != "file" && attachmentType != "link" && attachmentType != "video" && attachmentType != "subtitle" && attachmentType != "thumbnail" {
			return "", "", "", "", "", "", nil, "", errLearningPointArtifactInvalidInput
		}
		if err := validateResearchMaterialAttachmentReference(url, filePath, fileSize); err != nil {
			return "", "", "", "", "", "", nil, "", err
		}
	}
	return provider, attachmentType, sourceContext, title, url, filePath, fileSize, mimeType, nil
}

func validateResearchMaterialAttachmentReference(url, filePath string, fileSize *int64) error {
	reference := strings.TrimSpace(filePath)
	if reference == "" {
		reference = strings.TrimSpace(url)
	}
	extension := strings.ToLower(filepath.Ext(strings.Split(strings.Split(reference, "?")[0], "#")[0]))
	if extension == "" {
		return errLearningPointArtifactInvalidInput
	}
	_, imageOK := researchMaterialAllowedImageExtensions[extension]
	_, docOK := researchMaterialAllowedDocumentExtensions[extension]
	_, videoOK := researchMaterialAllowedVideoExtensions[extension]
	_, subtitleOK := researchMaterialAllowedSubtitleExtensions[extension]
	if !imageOK && !docOK && !videoOK && !subtitleOK {
		return errLearningPointArtifactInvalidInput
	}
	if fileSize == nil {
		return nil
	}
	limit := researchMaterialDocumentMaxBytes
	if imageOK {
		limit = researchMaterialImageMaxBytes
	} else if videoOK {
		limit = researchMaterialVideoMaxBytes
	} else if subtitleOK {
		limit = researchMaterialSubtitleMaxBytes
	}
	if *fileSize > limit {
		return errLearningPointArtifactInvalidInput
	}
	return nil
}

func normalizePointMaterialReportInput(targetTypeValue, reportTypeValue, messageValue *string) (string, string, string, error) {
	targetType := trimOptionalString(targetTypeValue)
	if targetType == "" {
		targetType = "source"
	}
	switch targetType {
	case "source", "ai_summary", "attachment", "other":
	default:
		return "", "", "", errLearningPointArtifactInvalidInput
	}

	reportType := trimOptionalString(reportTypeValue)
	if reportType == "" {
		reportType = "other"
	}
	switch reportType {
	case "broken_link", "wrong_content", "unsafe_content", "copyright", "low_quality", "other":
	default:
		return "", "", "", errLearningPointArtifactInvalidInput
	}

	message := trimOptionalString(messageValue)
	if len([]rune(message)) < 5 || len([]rune(message)) > 1000 {
		return "", "", "", errLearningPointArtifactInvalidInput
	}
	return targetType, reportType, message, nil
}

func normalizePracticeLogInput(
	titleValue *string,
	activityNameValue *string,
	attemptCount *int,
	successCount *int,
	failureCount *int,
	durationMinutes *int,
	blockedPartValue *string,
	changedMethodValue *string,
	achievementNoteValue *string,
	achievementValue *string,
	nextPlanValue *string,
	nextPracticeValue *string,
) (string, string, *int, *int, *int, *int, string, string, string, string, string, string, error) {
	if titleValue == nil {
		return "", "", nil, nil, nil, nil, "", "", "", "", "", "", errLearningPointPracticeLogInvalidInput
	}
	title := strings.TrimSpace(*titleValue)
	activityName := trimOptionalString(activityNameValue)
	if activityName == "" {
		activityName = title
	}
	blockedPart := ""
	if blockedPartValue != nil {
		blockedPart = strings.TrimSpace(*blockedPartValue)
	}
	changedMethod := ""
	if changedMethodValue != nil {
		changedMethod = strings.TrimSpace(*changedMethodValue)
	}
	achievementNote := ""
	if achievementNoteValue != nil {
		achievementNote = strings.TrimSpace(*achievementNoteValue)
	}
	achievement := trimOptionalString(achievementValue)
	if achievement == "" {
		achievement = achievementNote
	}
	nextPlan := ""
	if nextPlanValue != nil {
		nextPlan = strings.TrimSpace(*nextPlanValue)
	}
	nextPractice := trimOptionalString(nextPracticeValue)
	if nextPractice == "" {
		nextPractice = nextPlan
	}
	if title == "" || negativeIntPtr(attemptCount) || negativeIntPtr(successCount) || negativeIntPtr(failureCount) || negativeIntPtr(durationMinutes) {
		return "", "", nil, nil, nil, nil, "", "", "", "", "", "", errLearningPointPracticeLogInvalidInput
	}
	return title, activityName, attemptCount, successCount, failureCount, durationMinutes, blockedPart, changedMethod, achievementNote, achievement, nextPlan, nextPractice, nil
}

func trimOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func nullableTrimmedString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizePointQuestionAnswerMethod(value *string) *string {
	method := trimOptionalString(value)
	if method == "" {
		method = "self"
	}
	switch method {
	case "ai", "community", "self", "creator":
		return &method
	default:
		fallback := "self"
		return &fallback
	}
}

func negativeIntPtr(value *int) bool {
	return value != nil && *value < 0
}

func validScorePtr(value *int) bool {
	return value == nil || (*value >= 1 && *value <= 5)
}

func scanCoursePointPracticeLog(row pgx.Row) (CoursePointPracticeLog, error) {
	var log CoursePointPracticeLog
	err := row.Scan(
		&log.ID,
		&log.CoursePointID,
		&log.UserID,
		&log.Title,
		&log.ActivityName,
		&log.AttemptCount,
		&log.SuccessCount,
		&log.FailureCount,
		&log.DurationMinutes,
		&log.BlockedPart,
		&log.ChangedMethod,
		&log.AchievementNote,
		&log.Achievement,
		&log.NextPlan,
		&log.NextPractice,
		&log.CreatedAt,
		&log.UpdatedAt,
	)
	return log, err
}

func sanitizePointArtifactDescription(value string) string {
	return sanitizeTiptapHTML(value)
}
