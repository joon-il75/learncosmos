package curriculum

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) getCourseAggregateByID(ctx context.Context, courseID uuid.UUID) (*CourseAggregate, error) {
	var course Course
	err := r.pool.QueryRow(ctx, `
		SELECT c.id, c.user_id, c.source_draft_id, c.source_query, c.learning_goal, c.current_level,
		       c.duration_weeks, c.study_hours_per_week, c.preferred_format,
		       c.title, c.description, COALESCE(c.completion_criteria, '{}'::text[]), c.status, c.planet_type_id,
		       pt.name, pt.asset_path, c.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, c.created_at, c.updated_at
		FROM courses c
		LEFT JOIN planet_types pt ON pt.id = c.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = c.planet_texture_map_id
		WHERE c.id = $1
	`, courseID).Scan(
		&course.ID,
		&course.UserID,
		&course.SourceDraftID,
		&course.SourceQuery,
		&course.LearningGoal,
		&course.CurrentLevel,
		&course.DurationWeeks,
		&course.StudyHoursPerWeek,
		&course.PreferredFormat,
		&course.Title,
		&course.Description,
		&course.CompletionCriteria,
		&course.Status,
		&course.PlanetTypeID,
		&course.PlanetTypeName,
		&course.PlanetTypeAsset,
		&course.PlanetTextureMapID,
		&course.PlanetTextureMapName,
		&course.PlanetTextureMapAsset,
		&course.PlanetTextureMapRotationDurationSeconds,
		&course.PlanetTextureMapRotationDirection,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}

	lessonRows, err := r.pool.Query(ctx, `
		SELECT id, course_id, parent_lesson_id, title, objective, summary, difficulty_level,
		       lesson_role, source_type, order_index, created_at, updated_at
		FROM course_lessons
		WHERE course_id = $1
		ORDER BY order_index ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course lessons: %w", err)
	}
	defer lessonRows.Close()

	mainTrees := []CourseLessonTree{}
	mainIndex := map[uuid.UUID]int{}
	subIndex := map[uuid.UUID][2]int{}
	for lessonRows.Next() {
		var lesson CourseLesson
		if err := lessonRows.Scan(
			&lesson.ID, &lesson.CourseID, &lesson.ParentLessonID, &lesson.Title, &lesson.Objective, &lesson.Summary,
			&lesson.DifficultyLevel, &lesson.LessonRole, &lesson.SourceType, &lesson.OrderIndex,
			&lesson.CreatedAt, &lesson.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course lesson: %w", err)
		}
		if lesson.ParentLessonID == nil {
			mainIndex[lesson.ID] = len(mainTrees)
			mainTrees = append(mainTrees, CourseLessonTree{
				Lesson:     lesson,
				Points:     []CoursePointAggregate{},
				SubLessons: []CourseLessonTree{},
			})
		} else {
			mainPos, ok := mainIndex[*lesson.ParentLessonID]
			if !ok {
				continue
			}
			subPos := len(mainTrees[mainPos].SubLessons)
			mainTrees[mainPos].SubLessons = append(mainTrees[mainPos].SubLessons, CourseLessonTree{
				Lesson:     lesson,
				Points:     []CoursePointAggregate{},
				SubLessons: []CourseLessonTree{},
			})
			subIndex[lesson.ID] = [2]int{mainPos, subPos}
		}
	}

	coursePointRows, err := r.pool.Query(ctx, `
		SELECT id, course_lesson_id, point_type, content_id,
		       external_url, title, description, point_goal, point_category, thumbnail_url, price_type, rank_score,
		       template_type, status, order_index, created_at, updated_at
		FROM course_points
		WHERE course_id = $1
		ORDER BY order_index ASC, created_at ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course points: %w", err)
	}
	defer coursePointRows.Close()

	for coursePointRows.Next() {
		var p CoursePoint
		var priceTypeStr *string
		var templateType *ResearchNodeTemplateType
		p.CourseID = courseID
		if err := coursePointRows.Scan(
			&p.ID, &p.CourseLessonID, &p.PointType, &p.ContentID,
			&p.ExternalURL, &p.Title, &p.Description, &p.PointGoal, &p.PointCategory, &p.ThumbnailURL, &priceTypeStr, &p.RankScore,
			&templateType, &p.Status, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course point: %w", err)
		}
		if priceTypeStr != nil && *priceTypeStr != "" {
			pt := PriceType(*priceTypeStr)
			p.PriceType = &pt
		}
		p.TemplateType = templateType

		if pos, ok := subIndex[p.CourseLessonID]; ok {
			mainTrees[pos[0]].SubLessons[pos[1]].Points = append(mainTrees[pos[0]].SubLessons[pos[1]].Points, CoursePointAggregate{
				Point:     p,
				Questions: []CoursePointQuestion{},
			})
			continue
		}
		if mainPos, ok := mainIndex[p.CourseLessonID]; ok {
			mainTrees[mainPos].Points = append(mainTrees[mainPos].Points, CoursePointAggregate{
				Point:     p,
				Questions: []CoursePointQuestion{},
			})
			continue
		}
	}
	if err := coursePointRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course points: %w", err)
	}

	pointBlockRows, err := r.pool.Query(ctx, `
		SELECT cpb.id, cpb.course_point_id, cpb.user_id, cpb.block_type, cpb.content, cpb.order_index, cpb.created_at, cpb.updated_at
		FROM course_point_blocks cpb
		JOIN course_points cp ON cp.id = cpb.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpb.order_index ASC, cpb.created_at ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point blocks: %w", err)
	}
	defer pointBlockRows.Close()

	for pointBlockRows.Next() {
		var block CoursePointBlock
		if err := pointBlockRows.Scan(
			&block.ID,
			&block.CoursePointID,
			&block.UserID,
			&block.BlockType,
			&block.Content,
			&block.OrderIndex,
			&block.CreatedAt,
			&block.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course point block: %w", err)
		}
		attachCoursePointEntry(mainTrees, block.CoursePointID, func(point *CoursePointAggregate) {
			point.Blocks = append(point.Blocks, block)
		})
	}
	if err := pointBlockRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course point blocks: %w", err)
	}

	pointAISummaryRows, err := r.pool.Query(ctx, `
		SELECT cpas.course_point_id, cpas.source_title, cpas.source_description, cpas.summary,
		       COALESCE(cpas.learning_language, 'ko'), cpas.updated_at
		FROM course_point_ai_summaries cpas
		JOIN course_points cp ON cp.id = cpas.course_point_id
		WHERE cp.course_id = $1
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point ai summaries: %w", err)
	}
	defer pointAISummaryRows.Close()

	for pointAISummaryRows.Next() {
		var pointID uuid.UUID
		var sourceTitle string
		var sourceDescription string
		var summary string
		var learningLanguage string
		var updatedAt time.Time
		if err := pointAISummaryRows.Scan(&pointID, &sourceTitle, &sourceDescription, &summary, &learningLanguage, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan course point ai summary: %w", err)
		}
		attachCoursePointEntry(mainTrees, pointID, func(point *CoursePointAggregate) {
			point.AISummary = &PointAISummaryEntry{
				SourceTitle:       sourceTitle,
				SourceDescription: sourceDescription,
				Summary:           summary,
				LearningLanguage:  normalizeLearningLanguage(learningLanguage),
				UpdatedAt:         updatedAt,
			}
		})
	}
	if err := pointAISummaryRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course point ai summaries: %w", err)
	}

	pointJournalRows, err := r.pool.Query(ctx, `
		SELECT cpje.course_point_id, cpje.observation, cpje.reflection, cpje.next_step,
		       COALESCE(cpje.core_concept, ''), COALESCE(cpje.my_explanation, ''),
		       COALESCE(cpje.examples, ''), COALESCE(cpje.confused_parts, ''),
		       COALESCE(cpje.reference_links, ''), cpje.updated_at
		FROM course_point_journal_entries cpje
		JOIN course_points cp ON cp.id = cpje.course_point_id
		WHERE cp.course_id = $1
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point journal entries: %w", err)
	}
	defer pointJournalRows.Close()

	for pointJournalRows.Next() {
		var pointID uuid.UUID
		var observation string
		var reflection string
		var nextStep string
		var coreConcept string
		var myExplanation string
		var examples string
		var confusedParts string
		var referenceLinks string
		var updatedAt time.Time
		if err := pointJournalRows.Scan(&pointID, &observation, &reflection, &nextStep, &coreConcept, &myExplanation, &examples, &confusedParts, &referenceLinks, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan course point journal entry: %w", err)
		}
		attachCoursePointEntry(mainTrees, pointID, func(point *CoursePointAggregate) {
			point.JournalEntry = &DraftJournalEntry{
				Observation:    observation,
				Reflection:     reflection,
				NextStep:       nextStep,
				CoreConcept:    coreConcept,
				MyExplanation:  myExplanation,
				Examples:       examples,
				ConfusedParts:  confusedParts,
				ReferenceLinks: referenceLinks,
				UpdatedAt:      updatedAt,
			}
		})
	}

	pointObservationNoteRows, err := r.pool.Query(ctx, `
		SELECT cpon.id, cpon.course_point_id, cpon.user_id, cpon.note_type, cpon.content,
		       cpon.order_index, cpon.created_at, cpon.updated_at
		FROM course_point_observation_notes cpon
		JOIN course_points cp ON cp.id = cpon.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpon.order_index ASC, cpon.created_at ASC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_observation_notes") {
			return nil, fmt.Errorf("get course point observation notes: %w", err)
		}
	} else {
		defer pointObservationNoteRows.Close()
		for pointObservationNoteRows.Next() {
			note, err := scanCoursePointObservationNote(pointObservationNoteRows)
			if err != nil {
				return nil, fmt.Errorf("scan course point observation note: %w", err)
			}
			attachCoursePointEntry(mainTrees, note.CoursePointID, func(point *CoursePointAggregate) {
				point.ObservationNotes = append(point.ObservationNotes, note)
			})
		}
		if err := pointObservationNoteRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point observation notes: %w", err)
		}
	}

	pointRecordRows, err := r.pool.Query(ctx, `
		SELECT cpre.course_point_id, cpre.study_minutes, cpre.practice_count, cpre.confidence_level, cpre.application_note, cpre.updated_at
		FROM course_point_record_entries cpre
		JOIN course_points cp ON cp.id = cpre.course_point_id
		WHERE cp.course_id = $1
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point record entries: %w", err)
	}
	defer pointRecordRows.Close()

	for pointRecordRows.Next() {
		var pointID uuid.UUID
		var studyMinutes int
		var practiceCount int
		var confidenceLevel int
		var applicationNote string
		var updatedAt time.Time
		if err := pointRecordRows.Scan(&pointID, &studyMinutes, &practiceCount, &confidenceLevel, &applicationNote, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan course point record entry: %w", err)
		}
		attachCoursePointEntry(mainTrees, pointID, func(point *CoursePointAggregate) {
			point.RecordEntry = &DraftRecordEntry{
				StudyMinutes:    studyMinutes,
				PracticeCount:   practiceCount,
				ConfidenceLevel: confidenceLevel,
				ApplicationNote: applicationNote,
				UpdatedAt:       updatedAt,
			}
		})
	}

	pointArtifactRows, err := r.pool.Query(ctx, `
		SELECT cpa.id, cpa.course_point_id, cpa.artifact_type, cpa.title, cpa.url, cpa.description,
		       COALESCE(cpa.point_category, ''), COALESCE(cpa.production_process, ''),
		       COALESCE(cpa.learned_points, ''), COALESCE(cpa.difficult_points, ''),
		       cpa.visibility, cpa.order_index, cpa.created_at, cpa.updated_at
		FROM course_point_artifacts cpa
		JOIN course_points cp ON cp.id = cpa.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpa.order_index ASC, cpa.created_at ASC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_artifacts") {
			return nil, fmt.Errorf("get course point artifacts: %w", err)
		}
		legacyArtifactRows, legacyErr := r.pool.Query(ctx, `
			SELECT cpae.course_point_id, cpae.artifact_type, cpae.title, cpae.url, cpae.description, cpae.updated_at
			FROM course_point_artifact_entries cpae
			JOIN course_points cp ON cp.id = cpae.course_point_id
			WHERE cp.course_id = $1
		`, courseID)
		if legacyErr != nil {
			return nil, fmt.Errorf("get legacy course point artifact entries: %w", legacyErr)
		}
		defer legacyArtifactRows.Close()
		for legacyArtifactRows.Next() {
			var pointID uuid.UUID
			var artifactType string
			var title string
			var url string
			var description string
			var updatedAt time.Time
			if err := legacyArtifactRows.Scan(&pointID, &artifactType, &title, &url, &description, &updatedAt); err != nil {
				return nil, fmt.Errorf("scan legacy course point artifact entry: %w", err)
			}
			attachCoursePointEntry(mainTrees, pointID, func(point *CoursePointAggregate) {
				point.ArtifactEntry = &DraftArtifactEntry{
					ArtifactType: artifactType,
					Title:        title,
					URL:          url,
					Description:  description,
					UpdatedAt:    updatedAt,
				}
			})
		}
		if err := legacyArtifactRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate legacy course point artifact entries: %w", err)
		}
	} else {
		defer pointArtifactRows.Close()
		for pointArtifactRows.Next() {
			var artifact CoursePointArtifact
			if err := pointArtifactRows.Scan(
				&artifact.ID,
				&artifact.CoursePointID,
				&artifact.ArtifactType,
				&artifact.Title,
				&artifact.URL,
				&artifact.Description,
				&artifact.PointCategory,
				&artifact.ProductionProcess,
				&artifact.LearnedPoints,
				&artifact.DifficultPoints,
				&artifact.Visibility,
				&artifact.OrderIndex,
				&artifact.CreatedAt,
				&artifact.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("scan course point artifact: %w", err)
			}
			attachCoursePointEntry(mainTrees, artifact.CoursePointID, func(point *CoursePointAggregate) {
				point.Artifacts = append(point.Artifacts, artifact)
				if point.ArtifactEntry == nil {
					point.ArtifactEntry = &DraftArtifactEntry{
						ArtifactType: artifact.ArtifactType,
						Title:        artifact.Title,
						URL:          artifact.URL,
						Description:  artifact.Description,
						UpdatedAt:    artifact.UpdatedAt,
					}
				}
			})
		}
		if err := pointArtifactRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point artifacts: %w", err)
		}
	}

	pointAttachmentRows, err := r.pool.Query(ctx, `
		SELECT cpa.id, cpa.course_point_id, cpa.artifact_id, cpa.user_id, cpa.provider, cpa.attachment_type,
		       COALESCE(cpa.source_context, 'work_attachment'),
		       COALESCE(cpa.title, ''), COALESCE(cpa.url, ''), COALESCE(cpa.file_path, ''),
		       cpa.file_size, COALESCE(cpa.mime_type, ''), cpa.created_at, cpa.updated_at
		FROM course_point_attachments cpa
		JOIN course_points cp ON cp.id = cpa.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpa.created_at ASC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_attachments") {
			return nil, fmt.Errorf("get course point attachments: %w", err)
		}
	} else {
		defer pointAttachmentRows.Close()
		for pointAttachmentRows.Next() {
			attachment, err := scanCoursePointAttachment(pointAttachmentRows)
			if err != nil {
				return nil, fmt.Errorf("scan course point attachment: %w", err)
			}
			attachCoursePointEntry(mainTrees, attachment.CoursePointID, func(point *CoursePointAggregate) {
				point.Attachments = append(point.Attachments, attachment)
			})
		}
		if err := pointAttachmentRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point attachments: %w", err)
		}
	}

	pointMaterialReportRows, err := r.pool.Query(ctx, `
		SELECT cpmr.id, cpmr.course_point_id, cpmr.user_id, cpmr.attachment_id,
		       cpmr.target_type, cpmr.report_type, cpmr.message, cpmr.status,
		       cpmr.target_content_id, COALESCE(cpmr.target_url, ''), COALESCE(cpmr.target_title, ''),
		       cpmr.replacement_content_id, COALESCE(cpmr.replacement_url, ''), COALESCE(cpmr.replacement_title, ''), cpmr.replaced_at,
		       cpmr.created_at, cpmr.updated_at
		FROM course_point_material_reports cpmr
		JOIN course_points cp ON cp.id = cpmr.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpmr.created_at DESC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_material_reports") {
			return nil, fmt.Errorf("get course point material reports: %w", err)
		}
	} else {
		defer pointMaterialReportRows.Close()
		for pointMaterialReportRows.Next() {
			var report CoursePointMaterialReport
			if err := pointMaterialReportRows.Scan(
				&report.ID,
				&report.CoursePointID,
				&report.UserID,
				&report.AttachmentID,
				&report.TargetType,
				&report.ReportType,
				&report.Message,
				&report.Status,
				&report.TargetContentID,
				&report.TargetURL,
				&report.TargetTitle,
				&report.ReplacementContentID,
				&report.ReplacementURL,
				&report.ReplacementTitle,
				&report.ReplacedAt,
				&report.CreatedAt,
				&report.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("scan course point material report: %w", err)
			}
			attachCoursePointEntry(mainTrees, report.CoursePointID, func(point *CoursePointAggregate) {
				point.MaterialReports = append(point.MaterialReports, report)
			})
		}
		if err := pointMaterialReportRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point material reports: %w", err)
		}
	}

	pointPracticeLogRows, err := r.pool.Query(ctx, `
		SELECT cppl.id, cppl.course_point_id, cppl.user_id, cppl.title, COALESCE(cppl.activity_name, cppl.title), cppl.attempt_count,
		       cppl.success_count, cppl.failure_count, cppl.duration_minutes, cppl.blocked_part,
		       cppl.changed_method, cppl.achievement_note, COALESCE(cppl.achievement, cppl.achievement_note),
		       cppl.next_plan, COALESCE(cppl.next_practice, cppl.next_plan), cppl.created_at, cppl.updated_at
		FROM course_point_practice_logs cppl
		JOIN course_points cp ON cp.id = cppl.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cppl.created_at ASC, cppl.updated_at ASC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_practice_logs") {
			return nil, fmt.Errorf("get course point practice logs: %w", err)
		}
	} else {
		defer pointPracticeLogRows.Close()
		for pointPracticeLogRows.Next() {
			var log CoursePointPracticeLog
			if err := pointPracticeLogRows.Scan(
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
			); err != nil {
				return nil, fmt.Errorf("scan course point practice log: %w", err)
			}
			attachCoursePointEntry(mainTrees, log.CoursePointID, func(point *CoursePointAggregate) {
				point.PracticeLogs = append(point.PracticeLogs, log)
			})
		}
		if err := pointPracticeLogRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point practice logs: %w", err)
		}
	}

	pointEventRows, err := r.pool.Query(ctx, `
		WITH ranked_events AS (
			SELECT cpe.id, cpe.course_point_id, cpe.user_id, cpe.event_type, cpe.event_payload, cpe.created_at,
			       ROW_NUMBER() OVER (PARTITION BY cpe.course_point_id ORDER BY cpe.created_at DESC) AS rn
			FROM course_point_events cpe
			JOIN course_points cp ON cp.id = cpe.course_point_id
			WHERE cp.course_id = $1
		)
		SELECT id, course_point_id, user_id, event_type, event_payload, created_at
		FROM ranked_events
		WHERE rn <= 50
		ORDER BY course_point_id ASC, created_at ASC
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_events") {
			return nil, fmt.Errorf("get course point events: %w", err)
		}
	} else {
		defer pointEventRows.Close()
		for pointEventRows.Next() {
			var event CoursePointEvent
			if err := pointEventRows.Scan(
				&event.ID,
				&event.CoursePointID,
				&event.UserID,
				&event.EventType,
				&event.EventPayload,
				&event.CreatedAt,
			); err != nil {
				return nil, fmt.Errorf("scan course point event: %w", err)
			}
			attachCoursePointEntry(mainTrees, event.CoursePointID, func(point *CoursePointAggregate) {
				point.Events = append(point.Events, event)
			})
		}
		if err := pointEventRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate course point events: %w", err)
		}
	}

	researchMaterialStatusRows, err := r.pool.Query(ctx, `
		WITH material_state AS (
			SELECT cp.id AS point_id
			FROM course_points cp
			WHERE cp.course_id = $1
			  AND cp.point_type = 'research'
		),
		block_state AS (
			SELECT
				course_point_id,
				COUNT(*) AS block_count
			FROM course_point_blocks
			GROUP BY course_point_id
		),
		attachment_state AS (
			SELECT
				course_point_id,
				COUNT(*) AS attachment_count
			FROM course_point_attachments
			WHERE source_context = 'research_material'
			GROUP BY course_point_id
		),
		combined_material_state AS (
			SELECT
				material_state.point_id,
				COALESCE(block_state.block_count, 0) + COALESCE(attachment_state.attachment_count, 0) AS material_count
			FROM material_state
			LEFT JOIN block_state ON block_state.course_point_id = material_state.point_id
			LEFT JOIN attachment_state ON attachment_state.course_point_id = material_state.point_id
		),
		confirmation_events AS (
			SELECT
				course_point_id,
				event_type,
				created_at,
				ROW_NUMBER() OVER (PARTITION BY course_point_id ORDER BY created_at DESC) AS rn
			FROM course_point_events
			WHERE event_type IN ('research_material_confirmed', 'research_material_unconfirmed')
		),
		confirmation_state AS (
			SELECT course_point_id, event_type, created_at AS confirmed_at
			FROM confirmation_events
			WHERE rn = 1
		)
		SELECT
			combined_material_state.point_id,
			COALESCE(
				combined_material_state.material_count > 0
				AND confirmation_state.event_type = 'research_material_confirmed',
				false
			) AS confirmed,
			CASE WHEN confirmation_state.event_type = 'research_material_confirmed' THEN confirmation_state.confirmed_at ELSE NULL END AS confirmed_at
		FROM combined_material_state
		LEFT JOIN confirmation_state ON confirmation_state.course_point_id = combined_material_state.point_id
	`, courseID)
	if err != nil {
		if !isUndefinedTableError(err, "course_point_events") {
			return nil, fmt.Errorf("get research material confirmation states: %w", err)
		}
	} else {
		defer researchMaterialStatusRows.Close()
		for researchMaterialStatusRows.Next() {
			var pointID uuid.UUID
			var confirmed bool
			var confirmedAt *time.Time
			if err := researchMaterialStatusRows.Scan(&pointID, &confirmed, &confirmedAt); err != nil {
				return nil, fmt.Errorf("scan research material confirmation state: %w", err)
			}
			attachCoursePointEntry(mainTrees, pointID, func(point *CoursePointAggregate) {
				point.ResearchMaterialConfirmed = confirmed
				if confirmed {
					point.ResearchMaterialConfirmedAt = confirmedAt
				} else {
					point.ResearchMaterialConfirmedAt = nil
				}
			})
		}
		if err := researchMaterialStatusRows.Err(); err != nil {
			return nil, fmt.Errorf("iterate research material confirmation states: %w", err)
		}
	}

	pointQuestionRows, err := r.pool.Query(ctx, `
		SELECT cpq.id, cpq.course_point_id, cpq.goal_profile_version, cpq.title, cpq.question, cpq.question_type,
		       cpq.answer_method, cpq.answer, cpq.status, cpq.created_by, cpq.created_at, cpq.updated_at
		FROM course_point_questions cpq
		JOIN course_points cp ON cp.id = cpq.course_point_id
		WHERE cp.course_id = $1
		ORDER BY cpq.created_at ASC, cpq.updated_at ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point questions: %w", err)
	}
	defer pointQuestionRows.Close()

	for pointQuestionRows.Next() {
		var question CoursePointQuestion
		if err := pointQuestionRows.Scan(
			&question.ID,
			&question.CoursePointID,
			&question.GoalProfileVersion,
			&question.Title,
			&question.Question,
			&question.QuestionType,
			&question.AnswerMethod,
			&question.Answer,
			&question.Status,
			&question.CreatedBy,
			&question.CreatedAt,
			&question.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course point question: %w", err)
		}
		attachCoursePointEntry(mainTrees, question.CoursePointID, func(point *CoursePointAggregate) {
			point.Questions = append(point.Questions, question)
		})
	}
	if err := pointQuestionRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course point questions: %w", err)
	}

	pointQuestionFeedbackRows, err := r.pool.Query(ctx, `
		SELECT cpqf.course_point_question_id, cpqf.feedback
		FROM course_point_question_ai_feedbacks cpqf
		JOIN course_point_questions cpq ON cpq.id = cpqf.course_point_question_id
		JOIN course_points cp ON cp.id = cpq.course_point_id
		WHERE cp.course_id = $1
	`, courseID)
	if err != nil {
		if isUndefinedTableError(err, "course_point_question_ai_feedbacks") {
			goto selfEvaluations
		}
		return nil, fmt.Errorf("get course point question ai feedbacks: %w", err)
	}
	defer pointQuestionFeedbackRows.Close()

	for pointQuestionFeedbackRows.Next() {
		var questionID uuid.UUID
		var feedback string
		if err := pointQuestionFeedbackRows.Scan(&questionID, &feedback); err != nil {
			return nil, fmt.Errorf("scan course point question ai feedback: %w", err)
		}
		for mainIdx := range mainTrees {
			for pointIdx := range mainTrees[mainIdx].Points {
				for questionIdx := range mainTrees[mainIdx].Points[pointIdx].Questions {
					if mainTrees[mainIdx].Points[pointIdx].Questions[questionIdx].ID == questionID {
						mainTrees[mainIdx].Points[pointIdx].Questions[questionIdx].AIFeedback = &feedback
					}
				}
			}
			for subIdx := range mainTrees[mainIdx].SubLessons {
				for pointIdx := range mainTrees[mainIdx].SubLessons[subIdx].Points {
					for questionIdx := range mainTrees[mainIdx].SubLessons[subIdx].Points[pointIdx].Questions {
						if mainTrees[mainIdx].SubLessons[subIdx].Points[pointIdx].Questions[questionIdx].ID == questionID {
							mainTrees[mainIdx].SubLessons[subIdx].Points[pointIdx].Questions[questionIdx].AIFeedback = &feedback
						}
					}
				}
			}
		}
	}
	if err := pointQuestionFeedbackRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course point question ai feedbacks: %w", err)
	}

selfEvaluations:
	pointSelfEvalRows, err := r.pool.Query(ctx, `
		SELECT cpse.course_point_id, cpse.goal_profile_version, cpse.understanding, cpse.application_note,
		       cpse.proficiency, cpse.understanding_score, COALESCE(cpse.understanding_reason, ''),
		       cpse.application_score, COALESCE(cpse.application_reason, ''),
		       cpse.proficiency_score, COALESCE(cpse.proficiency_reason, ''),
		       cpse.problem_solving_score, COALESCE(cpse.problem_solving_reason, ''),
		       cpse.expression_score, COALESCE(cpse.expression_reason, ''),
		       cpse.goal_alignment_note, cpse.final_score, cpse.updated_at
		FROM course_point_self_evaluations cpse
		JOIN course_points cp ON cp.id = cpse.course_point_id
		WHERE cp.course_id = $1
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course point self evaluations: %w", err)
	}
	defer pointSelfEvalRows.Close()

	for pointSelfEvalRows.Next() {
		var selfEval CoursePointSelfEvaluation
		if err := pointSelfEvalRows.Scan(
			&selfEval.CoursePointID,
			&selfEval.GoalProfileVersion,
			&selfEval.Understanding,
			&selfEval.ApplicationNote,
			&selfEval.Proficiency,
			&selfEval.UnderstandingScore,
			&selfEval.UnderstandingReason,
			&selfEval.ApplicationScore,
			&selfEval.ApplicationReason,
			&selfEval.ProficiencyScore,
			&selfEval.ProficiencyReason,
			&selfEval.ProblemSolvingScore,
			&selfEval.ProblemSolvingReason,
			&selfEval.ExpressionScore,
			&selfEval.ExpressionReason,
			&selfEval.GoalAlignmentNote,
			&selfEval.FinalScore,
			&selfEval.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course point self evaluation: %w", err)
		}
		attachCoursePointEntry(mainTrees, selfEval.CoursePointID, func(point *CoursePointAggregate) {
			point.SelfEvaluation = &selfEval
		})
	}
	if err := pointSelfEvalRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate course point self evaluations: %w", err)
	}

	runtimeRows, err := r.pool.Query(ctx, `
		SELECT id, user_id, course_id, course_lesson_id, status, started_at, completed_at, created_at, updated_at
		FROM course_lesson_runtime_entries
		WHERE user_id = $1
		  AND course_id = $2
	`, course.UserID, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course lesson runtime entries: %w", err)
	}
	defer runtimeRows.Close()

	for runtimeRows.Next() {
		var entry CourseLessonRuntimeEntry
		if err := runtimeRows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.CourseID,
			&entry.CourseLessonID,
			&entry.Status,
			&entry.StartedAt,
			&entry.CompletedAt,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan course lesson runtime entry: %w", err)
		}
		pos, ok := subIndex[entry.CourseLessonID]
		if !ok {
			continue
		}
		mainTrees[pos[0]].SubLessons[pos[1]].RuntimeEntry = &entry
	}

	applyResearchMaterialStatusToLessonTrees(mainTrees)

	return &CourseAggregate{
		Course:  course,
		Lessons: mainTrees,
	}, nil
}
