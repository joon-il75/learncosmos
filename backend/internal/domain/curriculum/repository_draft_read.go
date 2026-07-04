package curriculum

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) ListDrafts(ctx context.Context, userID uuid.UUID) ([]CourseDraft, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       COALESCE(cd.generation_language, 'ko') AS generation_language,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id, cd.planet_type_id,
		       pt.name, pt.asset_path, cd.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = cd.planet_texture_map_id
		WHERE cd.user_id = $1
		ORDER BY cd.updated_at DESC, cd.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list drafts: %w", err)
	}
	defer rows.Close()

	drafts := []CourseDraft{}
	for rows.Next() {
		var draft CourseDraft
		if err := rows.Scan(
			&draft.ID,
			&draft.UserID,
			&draft.SourceQuery,
			&draft.LearningGoal,
			&draft.GoalProfileID,
			&draft.GoalProfileVersion,
			&draft.CurrentLevel,
			&draft.DurationWeeks,
			&draft.StudyHoursPerWeek,
			&draft.PreferredFormat,
			&draft.GenerationLanguage,
			&draft.Title,
			&draft.Description,
			&draft.CompletionCriteria,
			&draft.Status,
			&draft.ConfirmedCourseID,
			&draft.PlanetTypeID,
			&draft.PlanetTypeName,
			&draft.PlanetTypeAsset,
			&draft.PlanetTextureMapID,
			&draft.PlanetTextureMapName,
			&draft.PlanetTextureMapAsset,
			&draft.PlanetTextureMapRotationDurationSeconds,
			&draft.PlanetTextureMapRotationDirection,
			&draft.LastAccessedAt,
			&draft.CreatedAt,
			&draft.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan draft: %w", err)
		}
		drafts = append(drafts, draft)
	}
	return drafts, nil
}

func (r *Repository) ListDraftsForAdmin(ctx context.Context, limit int) ([]CourseDraft, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}

	rows, err := r.pool.Query(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       COALESCE(cd.generation_language, 'ko') AS generation_language,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id, cd.planet_type_id,
		       pt.name, pt.asset_path, cd.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = cd.planet_texture_map_id
		ORDER BY cd.updated_at DESC, cd.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list drafts for admin: %w", err)
	}
	defer rows.Close()

	drafts := []CourseDraft{}
	for rows.Next() {
		var draft CourseDraft
		if err := rows.Scan(
			&draft.ID,
			&draft.UserID,
			&draft.SourceQuery,
			&draft.LearningGoal,
			&draft.GoalProfileID,
			&draft.GoalProfileVersion,
			&draft.CurrentLevel,
			&draft.DurationWeeks,
			&draft.StudyHoursPerWeek,
			&draft.PreferredFormat,
			&draft.GenerationLanguage,
			&draft.Title,
			&draft.Description,
			&draft.CompletionCriteria,
			&draft.Status,
			&draft.ConfirmedCourseID,
			&draft.PlanetTypeID,
			&draft.PlanetTypeName,
			&draft.PlanetTypeAsset,
			&draft.PlanetTextureMapID,
			&draft.PlanetTextureMapName,
			&draft.PlanetTextureMapAsset,
			&draft.PlanetTextureMapRotationDurationSeconds,
			&draft.PlanetTextureMapRotationDirection,
			&draft.LastAccessedAt,
			&draft.CreatedAt,
			&draft.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan admin draft: %w", err)
		}
		drafts = append(drafts, draft)
	}
	return drafts, nil
}

func (r *Repository) loadDraftLessons(ctx context.Context, draftID uuid.UUID) ([]DraftLessonTree, error) {
	lessonNotes, err := r.fetchDraftLessonNotes(ctx, draftID)
	if err != nil {
		return nil, err
	}

	lessonRows, err := r.pool.Query(ctx, `
		SELECT cdl.id, cdl.course_draft_id, cdl.parent_lesson_id,
		       cdl.title, cdl.objective, cdl.summary, cdl.difficulty_level,
		       cdje.observation, cdje.reflection, cdje.next_step, cdje.updated_at,
		       cdre.study_minutes, cdre.practice_count, cdre.confidence_level, cdre.application_note, cdre.updated_at,
		       cdae.artifact_type, cdae.title, cdae.url, cdae.description, cdae.updated_at,
		       cdl.lesson_role, cdl.source_type, cdl.order_index, COALESCE(cdl.recommendation_search_spec, '{}'::jsonb), cdl.created_at, cdl.updated_at
		FROM course_draft_lessons cdl
		LEFT JOIN course_draft_journal_entries cdje ON cdje.course_draft_lesson_id = cdl.id
		LEFT JOIN course_draft_record_entries cdre ON cdre.course_draft_lesson_id = cdl.id
		LEFT JOIN course_draft_artifact_entries cdae ON cdae.course_draft_lesson_id = cdl.id
		WHERE cdl.course_draft_id = $1
		ORDER BY cdl.order_index ASC
	`, draftID)
	if err != nil {
		return nil, fmt.Errorf("get draft lessons: %w", err)
	}
	defer lessonRows.Close()

	mainTrees := []DraftLessonTree{}
	mainIndex := map[uuid.UUID]int{}
	subIndex := map[uuid.UUID][2]int{}

	for lessonRows.Next() {
		var lesson CourseDraftLesson
		var journalObservation *string
		var journalReflection *string
		var journalNextStep *string
		var journalUpdatedAt *time.Time
		var recordStudyMinutes *int
		var recordPracticeCount *int
		var recordConfidenceLevel *int
		var recordApplicationNote *string
		var recordUpdatedAt *time.Time
		var artifactType *string
		var artifactTitle *string
		var artifactURL *string
		var artifactDescription *string
		var artifactUpdatedAt *time.Time
		var recommendationSearchSpecJSON []byte
		if err := lessonRows.Scan(
			&lesson.ID, &lesson.CourseDraftID, &lesson.ParentLessonID,
			&lesson.Title, &lesson.Objective, &lesson.Summary, &lesson.DifficultyLevel,
			&journalObservation, &journalReflection, &journalNextStep, &journalUpdatedAt,
			&recordStudyMinutes, &recordPracticeCount, &recordConfidenceLevel, &recordApplicationNote, &recordUpdatedAt,
			&artifactType, &artifactTitle, &artifactURL, &artifactDescription, &artifactUpdatedAt,
			&lesson.LessonRole, &lesson.SourceType, &lesson.OrderIndex,
			&recommendationSearchSpecJSON,
			&lesson.CreatedAt, &lesson.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan draft lesson: %w", err)
		}
		lesson.RecommendationSearchSpec = ParseLessonRecommendationSearchSpec(recommendationSearchSpecJSON)
		if note, ok := lessonNotes[lesson.ID]; ok {
			lesson.OperationNote = &note
		}
		if journalUpdatedAt != nil {
			lesson.JournalEntry = &DraftJournalEntry{
				Observation: strings.TrimSpace(valueOrEmpty(journalObservation)),
				Reflection:  strings.TrimSpace(valueOrEmpty(journalReflection)),
				NextStep:    strings.TrimSpace(valueOrEmpty(journalNextStep)),
				UpdatedAt:   *journalUpdatedAt,
			}
		}
		if recordUpdatedAt != nil {
			lesson.RecordEntry = &DraftRecordEntry{
				StudyMinutes:    valueOrZero(recordStudyMinutes),
				PracticeCount:   valueOrZero(recordPracticeCount),
				ConfidenceLevel: valueOrDefault(recordConfidenceLevel, 3),
				ApplicationNote: strings.TrimSpace(valueOrEmpty(recordApplicationNote)),
				UpdatedAt:       *recordUpdatedAt,
			}
		}
		if artifactUpdatedAt != nil {
			lesson.ArtifactEntry = &DraftArtifactEntry{
				ArtifactType: strings.TrimSpace(valueOrEmpty(artifactType)),
				Title:        strings.TrimSpace(valueOrEmpty(artifactTitle)),
				URL:          strings.TrimSpace(valueOrEmpty(artifactURL)),
				Description:  strings.TrimSpace(valueOrEmpty(artifactDescription)),
				UpdatedAt:    *artifactUpdatedAt,
			}
		}
		if lesson.ParentLessonID == nil {
			mainIndex[lesson.ID] = len(mainTrees)
			mainTrees = append(mainTrees, DraftLessonTree{
				Lesson:     lesson,
				Points:     []DraftPointAggregate{},
				SubLessons: []DraftLessonTree{},
			})
		} else {
			mainPos, ok := mainIndex[*lesson.ParentLessonID]
			if !ok {
				continue
			}
			subPos := len(mainTrees[mainPos].SubLessons)
			mainTrees[mainPos].SubLessons = append(mainTrees[mainPos].SubLessons, DraftLessonTree{
				Lesson:     lesson,
				Points:     []DraftPointAggregate{},
				SubLessons: []DraftLessonTree{},
			})
			subIndex[lesson.ID] = [2]int{mainPos, subPos}
		}
	}
	if err := lessonRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate draft lessons: %w", err)
	}

	pointRows, err := r.pool.Query(ctx, `
		SELECT id, course_draft_lesson_id, point_type, selection_state, content_id,
		       external_url, title, description, thumbnail_url, price_type, rank_score,
		       template_type, status, order_index, created_at, updated_at
		FROM course_draft_points
		WHERE course_draft_id = $1
		ORDER BY order_index ASC, created_at ASC
	`, draftID)
	if err != nil {
		return nil, fmt.Errorf("get draft points: %w", err)
	}
	defer pointRows.Close()

	for pointRows.Next() {
		var p CourseDraftPoint
		var selectionState *ResourceSelectionState
		var priceTypeStr *string
		var templateType *ResearchNodeTemplateType
		p.CourseDraftID = draftID
		if err := pointRows.Scan(
			&p.ID, &p.CourseDraftLessonID, &p.PointType, &selectionState, &p.ContentID,
			&p.ExternalURL, &p.Title, &p.Description, &p.ThumbnailURL, &priceTypeStr, &p.RankScore,
			&templateType, &p.Status, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan draft point: %w", err)
		}
		p.SelectionState = selectionState
		if priceTypeStr != nil && *priceTypeStr != "" {
			pt := PriceType(*priceTypeStr)
			p.PriceType = &pt
		}
		p.TemplateType = templateType

		if p.PointType == PointTypeExploration {
			pos, ok := subIndex[p.CourseDraftLessonID]
			if !ok {
				continue
			}
			mainTrees[pos[0]].SubLessons[pos[1]].Points = append(mainTrees[pos[0]].SubLessons[pos[1]].Points, DraftPointAggregate{Point: p})
		} else {
			mainPos, ok := mainIndex[p.CourseDraftLessonID]
			if !ok {
				continue
			}
			mainTrees[mainPos].Points = append(mainTrees[mainPos].Points, DraftPointAggregate{Point: p})
		}
	}
	if err := pointRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate draft points: %w", err)
	}

	return mainTrees, nil
}

func (r *Repository) GetDraftByID(ctx context.Context, draftID, userID uuid.UUID) (*DraftAggregate, error) {
	var draft CourseDraft
	err := r.pool.QueryRow(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       COALESCE(cd.generation_language, 'ko') AS generation_language,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id,
		       COALESCE(
		           c.status = 'archived',
		           cd.status = 'archived'
		           AND cd.confirmed_course_id IS NULL
		           AND NOT EXISTS (
		               SELECT 1
		               FROM recommendation_events re
		               WHERE re.course_draft_id = cd.id
		                 AND re.event_type = 'curriculum_completed'
		           )
		       ), cd.planet_type_id,
		       pt.name, pt.asset_path, cd.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN courses c ON c.id = cd.confirmed_course_id AND c.user_id = cd.user_id
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = cd.planet_texture_map_id
		WHERE cd.id = $1 AND cd.user_id = $2
	`, draftID, userID).Scan(
		&draft.ID,
		&draft.UserID,
		&draft.SourceQuery,
		&draft.LearningGoal,
		&draft.GoalProfileID,
		&draft.GoalProfileVersion,
		&draft.CurrentLevel,
		&draft.DurationWeeks,
		&draft.StudyHoursPerWeek,
		&draft.PreferredFormat,
		&draft.GenerationLanguage,
		&draft.Title,
		&draft.Description,
		&draft.CompletionCriteria,
		&draft.Status,
		&draft.ConfirmedCourseID,
		&draft.IsInactive,
		&draft.PlanetTypeID,
		&draft.PlanetTypeName,
		&draft.PlanetTypeAsset,
		&draft.PlanetTextureMapID,
		&draft.PlanetTextureMapName,
		&draft.PlanetTextureMapAsset,
		&draft.PlanetTextureMapRotationDurationSeconds,
		&draft.PlanetTextureMapRotationDirection,
		&draft.LastAccessedAt,
		&draft.CreatedAt,
		&draft.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get draft: %w", err)
	}

	lessons, err := r.loadDraftLessons(ctx, draftID)
	if err != nil {
		return nil, err
	}

	return &DraftAggregate{
		Draft:   draft,
		Lessons: lessons,
	}, nil
}

func (r *Repository) GetDraftByIDForAdmin(ctx context.Context, draftID uuid.UUID) (*DraftAggregate, error) {
	var draft CourseDraft
	err := r.pool.QueryRow(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       COALESCE(cd.generation_language, 'ko') AS generation_language,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id, cd.planet_type_id,
		       pt.name, pt.asset_path, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		WHERE cd.id = $1
	`, draftID).Scan(
		&draft.ID,
		&draft.UserID,
		&draft.SourceQuery,
		&draft.LearningGoal,
		&draft.GoalProfileID,
		&draft.GoalProfileVersion,
		&draft.CurrentLevel,
		&draft.DurationWeeks,
		&draft.StudyHoursPerWeek,
		&draft.PreferredFormat,
		&draft.GenerationLanguage,
		&draft.Title,
		&draft.Description,
		&draft.CompletionCriteria,
		&draft.Status,
		&draft.ConfirmedCourseID,
		&draft.PlanetTypeID,
		&draft.PlanetTypeName,
		&draft.PlanetTypeAsset,
		&draft.LastAccessedAt,
		&draft.CreatedAt,
		&draft.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get draft for admin: %w", err)
	}

	lessons, err := r.loadDraftLessons(ctx, draftID)
	if err != nil {
		return nil, err
	}

	return &DraftAggregate{
		Draft:   draft,
		Lessons: lessons,
	}, nil
}
