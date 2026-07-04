package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetPlanetByIDAndDraftStatuses(ctx context.Context, userID, planetID uuid.UUID, statuses []DraftStatus) (*PlanetAggregate, error) {
	draft, err := r.findDraftByPlanetIDAndStatuses(ctx, userID, planetID, statuses)
	if err != nil {
		return nil, err
	}
	planet, err := r.buildPlanetAggregateFromDraftState(ctx, userID, *draft)
	if err != nil {
		return nil, err
	}
	if planet.Planet.IsInactive {
		return nil, errInactivePlanetAccessDenied
	}
	return planet, nil
}

func (r *Repository) GetPlanetPointDetailByIDAndDraftStatuses(ctx context.Context, userID, planetID, pointID uuid.UUID, statuses []DraftStatus) (*PlanetPointDetail, error) {
	planet, err := r.GetPlanetByIDAndDraftStatuses(ctx, userID, planetID, statuses)
	if err != nil {
		return nil, err
	}

	detail, ok := resolvePlanetPointDetail(planet, pointID)
	if !ok {
		return nil, errLearningPointNotFound
	}
	return detail, nil
}

func (r *Repository) GetPlanetRecordsByIDAndDraftStatuses(ctx context.Context, userID, planetID uuid.UUID, statuses []DraftStatus) (*PlanetRecordAggregate, error) {
	planet, err := r.GetPlanetByIDAndDraftStatuses(ctx, userID, planetID, statuses)
	if err != nil {
		return nil, err
	}
	return buildPlanetRecordAggregate(planet), nil
}

func (r *Repository) GetPlanetRecordFeedByIDAndDraftStatuses(ctx context.Context, userID, planetID uuid.UUID, statuses []DraftStatus, routeKind string) (*PlanetRecordFeedResponse, error) {
	planet, err := r.GetPlanetByIDAndDraftStatuses(ctx, userID, planetID, statuses)
	if err != nil {
		return nil, err
	}
	return buildPlanetRecordFeed(planet, routeKind), nil
}

func (r *Repository) GetPlanetResultsByIDAndDraftStatuses(ctx context.Context, userID, planetID uuid.UUID, statuses []DraftStatus) (*PlanetResultAggregate, error) {
	planet, err := r.GetPlanetByIDAndDraftStatuses(ctx, userID, planetID, statuses)
	if err != nil {
		return nil, err
	}
	return buildPlanetResultAggregate(planet), nil
}

func (r *Repository) verifyDraftOwner(ctx context.Context, userID, draftID uuid.UUID) error {
	var exists bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM course_drafts WHERE id = $1 AND user_id = $2)`,
		draftID, userID,
	).Scan(&exists); err != nil {
		return fmt.Errorf("verify draft owner: %w", err)
	}
	if !exists {
		return errDraftResourceNotFound
	}
	return nil
}

func getEditableDraftStatuses(includeLearning bool) []DraftStatus {
	statuses := []DraftStatus{DraftStatusDraft, DraftStatusConfirmed}
	if includeLearning {
		statuses = append(statuses, DraftStatusLearning)
	}
	return statuses
}

func fetchDraftMainLessonContext(
	ctx context.Context,
	tx pgx.Tx,
	draftID, userID, lessonID uuid.UUID,
	allowLearning bool,
) (string, string, error) {
	var sourceQuery string
	var title string
	err := tx.QueryRow(ctx, `
		SELECT cd.source_query, cdl.title
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		WHERE cdl.id = $1
		  AND cdl.course_draft_id = $2
		  AND cdl.parent_lesson_id IS NULL
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		FOR UPDATE OF cdl
	`, lessonID, draftID, userID, getEditableDraftStatuses(allowLearning)).Scan(&sourceQuery, &title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", errDraftLevelNotFound
		}
		return "", "", fmt.Errorf("get draft main lesson context: %w", err)
	}
	return sourceQuery, title, nil
}

func (r *Repository) fetchDraftLessonNotes(
	ctx context.Context,
	draftID uuid.UUID,
) (map[uuid.UUID]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT course_draft_lesson_id, note
		FROM course_draft_detail_notes
		WHERE course_draft_id = $1
	`, draftID)
	if err != nil {
		return nil, fmt.Errorf("get draft lesson notes: %w", err)
	}
	defer rows.Close()

	notes := make(map[uuid.UUID]string)
	for rows.Next() {
		var lessonID uuid.UUID
		var note string
		if err := rows.Scan(&lessonID, &note); err != nil {
			return nil, fmt.Errorf("scan draft lesson note: %w", err)
		}
		notes[lessonID] = note
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate draft lesson notes: %w", err)
	}
	return notes, nil
}

func (r *Repository) fetchDraftResearchPointContext(
	ctx context.Context,
	userID, draftID, pointID uuid.UUID,
) (*CourseDraftPoint, error) {
	var point CourseDraftPoint
	err := r.pool.QueryRow(ctx, `
		SELECT p.id, p.course_draft_id, p.course_draft_lesson_id,
		       p.point_type, p.status, p.title, p.description, p.template_type,
		       p.order_index, p.created_at, p.updated_at
		FROM course_draft_points p
		JOIN course_drafts cd ON cd.id = p.course_draft_id
		WHERE p.id = $1
		  AND cd.id = $2
		  AND cd.user_id = $3
		  AND cd.status = ANY($4)
		  AND p.point_type = 'research'
		LIMIT 1
	`, pointID, draftID, userID, getEditableDraftStatuses(true)).Scan(
		&point.ID, &point.CourseDraftID, &point.CourseDraftLessonID,
		&point.PointType, &point.Status, &point.Title, &point.Description, &point.TemplateType,
		&point.OrderIndex, &point.CreatedAt, &point.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errDraftResourceNotFound
		}
		return nil, fmt.Errorf("get draft research point context: %w", err)
	}
	return &point, nil
}

// ── Research Node CRUD ──────────────────────────────────────────────────────

func nullableJSON(raw json.RawMessage) *string {
	if len(raw) == 0 {
		return nil
	}
	s := string(raw)
	return &s
}

func buildPlanetLessonTreesFromDraftLessons(draftLessons []DraftLessonTree, courseID uuid.UUID) []CourseLessonTree {
	trees := make([]CourseLessonTree, 0, len(draftLessons))
	for _, draftLesson := range draftLessons {
		lesson := CourseLesson{
			ID:                       draftLesson.Lesson.ID,
			CourseID:                 courseID,
			ParentLessonID:           draftLesson.Lesson.ParentLessonID,
			Title:                    draftLesson.Lesson.Title,
			Objective:                draftLesson.Lesson.Objective,
			Summary:                  draftLesson.Lesson.Summary,
			DifficultyLevel:          draftLesson.Lesson.DifficultyLevel,
			RecommendationSearchSpec: draftLesson.Lesson.RecommendationSearchSpec,
			LessonRole:               draftLesson.Lesson.LessonRole,
			SourceType:               draftLesson.Lesson.SourceType,
			OrderIndex:               draftLesson.Lesson.OrderIndex,
			CreatedAt:                draftLesson.Lesson.CreatedAt,
			UpdatedAt:                draftLesson.Lesson.UpdatedAt,
		}

		tree := CourseLessonTree{
			Lesson:     lesson,
			Points:     make([]CoursePointAggregate, 0, len(draftLesson.Points)),
			SubLessons: buildPlanetLessonTreesFromDraftLessons(draftLesson.SubLessons, courseID),
		}

		for _, draftPoint := range draftLesson.Points {
			point := CoursePointAggregate{
				Point: CoursePoint{
					ID:             draftPoint.Point.ID,
					CourseID:       courseID,
					CourseLessonID: draftPoint.Point.CourseDraftLessonID,
					PointType:      draftPoint.Point.PointType,
					Status:         draftPoint.Point.Status,
					Title:          draftPoint.Point.Title,
					Description:    draftPoint.Point.Description,
					TemplateType:   draftPoint.Point.TemplateType,
					ContentID:      draftPoint.Point.ContentID,
					ExternalURL:    draftPoint.Point.ExternalURL,
					ThumbnailURL:   draftPoint.Point.ThumbnailURL,
					OrderIndex:     draftPoint.Point.OrderIndex,
					CompletedAt:    draftPoint.Point.CompletedAt,
					CreatedAt:      draftPoint.Point.CreatedAt,
					UpdatedAt:      draftPoint.Point.UpdatedAt,
				},
				JournalEntry:  draftPoint.JournalEntry,
				RecordEntry:   draftPoint.RecordEntry,
				ArtifactEntry: draftPoint.ArtifactEntry,
				Questions:     []CoursePointQuestion{},
			}
			if draftPoint.Point.PriceType != nil {
				point.Point.PriceType = draftPoint.Point.PriceType
			}
			if draftPoint.Point.RankScore != nil {
				point.Point.RankScore = draftPoint.Point.RankScore
			}
			for _, block := range draftPoint.Blocks {
				point.Blocks = append(point.Blocks, CoursePointBlock{
					ID:            block.ID,
					CoursePointID: point.Point.ID,
					BlockType:     block.BlockType,
					Content:       block.Content,
					OrderIndex:    block.OrderIndex,
					CreatedAt:     block.CreatedAt,
					UpdatedAt:     block.UpdatedAt,
				})
			}
			tree.Points = append(tree.Points, point)
		}

		trees = append(trees, tree)
	}
	return trees
}

func (r *Repository) loadPlanetGoalContext(ctx context.Context, draft CourseDraft) (*PlanetGoalContext, error) {
	context := &PlanetGoalContext{
		LearningGoal:       draft.LearningGoal,
		GoalProfileID:      draft.GoalProfileID,
		GoalProfileVersion: draft.GoalProfileVersion,
	}

	if draft.GoalProfileID != nil {
		if err := r.pool.QueryRow(ctx, `
			SELECT confirmed_goal, usage_context, motivation
			FROM course_goal_profiles
			WHERE id = $1
			LIMIT 1
		`, *draft.GoalProfileID).Scan(
			&context.ConfirmedGoal,
			&context.UsageContext,
			&context.Motivation,
		); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("load goal context by id: %w", err)
		}
	} else {
		if err := r.pool.QueryRow(ctx, `
			SELECT id, version, confirmed_goal, usage_context, motivation
			FROM course_goal_profiles
			WHERE course_draft_id = $1
			ORDER BY version DESC, updated_at DESC
			LIMIT 1
		`, draft.ID).Scan(
			&context.GoalProfileID,
			&context.GoalProfileVersion,
			&context.ConfirmedGoal,
			&context.UsageContext,
			&context.Motivation,
		); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("load goal context by draft id: %w", err)
		}
	}

	if context.LearningGoal == nil && context.ConfirmedGoal != nil {
		context.LearningGoal = context.ConfirmedGoal
	}
	if context.LearningGoal == nil && context.GoalProfileID == nil && context.GoalProfileVersion == nil &&
		context.ConfirmedGoal == nil && context.UsageContext == nil && context.Motivation == nil {
		return nil, nil
	}
	return context, nil
}

func (r *Repository) listDraftsByStatuses(ctx context.Context, userID uuid.UUID, statuses []DraftStatus) ([]CourseDraft, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id, cd.planet_type_id,
		       pt.name, pt.asset_path, cd.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = cd.planet_texture_map_id
		WHERE cd.user_id = $1
		  AND cd.status = ANY($2)
		ORDER BY cd.last_accessed_at DESC NULLS LAST, cd.updated_at DESC, cd.created_at DESC
	`, userID, statuses)
	if err != nil {
		return nil, fmt.Errorf("list planet drafts: %w", err)
	}
	defer rows.Close()

	drafts := make([]CourseDraft, 0)
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
			return nil, fmt.Errorf("scan planet draft: %w", err)
		}
		drafts = append(drafts, draft)
	}
	return drafts, nil
}

func (r *Repository) findDraftByPlanetIDAndStatuses(ctx context.Context, userID, planetID uuid.UUID, statuses []DraftStatus) (*CourseDraft, error) {
	var draft CourseDraft
	err := r.pool.QueryRow(ctx, `
		SELECT cd.id, cd.user_id, cd.source_query, cd.learning_goal, cd.goal_profile_id, cd.goal_profile_version, cd.current_level,
		       cd.duration_weeks, cd.study_hours_per_week, cd.preferred_format,
		       cd.title, cd.description, COALESCE(cd.completion_criteria, '{}'::text[]), cd.status, cd.confirmed_course_id, cd.planet_type_id,
		       pt.name, pt.asset_path, cd.planet_texture_map_id, ptmp.name, ptmp.asset_path,
		       ptmp.rotation_duration_seconds, ptmp.rotation_direction, cd.last_accessed_at, cd.created_at, cd.updated_at
		FROM course_drafts cd
		LEFT JOIN planet_types pt ON pt.id = cd.planet_type_id
		LEFT JOIN planet_texture_maps ptmp ON ptmp.id = cd.planet_texture_map_id
		WHERE cd.user_id = $1
		  AND cd.status = ANY($2)
		  AND (cd.id = $3 OR cd.confirmed_course_id = $3)
		ORDER BY cd.last_accessed_at DESC NULLS LAST, cd.updated_at DESC
		LIMIT 1
	`, userID, statuses, planetID).Scan(
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
	)
	if err != nil {
		return nil, fmt.Errorf("get planet draft: %w", err)
	}
	return &draft, nil
}

func (r *Repository) buildPlanetAggregateFromDraftState(ctx context.Context, userID uuid.UUID, draft CourseDraft) (*PlanetAggregate, error) {
	if draft.ConfirmedCourseID != nil {
		if courseAggregate, err := r.getCourseAggregateByID(ctx, *draft.ConfirmedCourseID); err == nil {
			return r.buildPlanetAggregateFromCourse(ctx, draft, courseAggregate)
		}
	}

	draftAggregate, err := r.GetDraftByID(ctx, draft.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("load draft fallback for planet: %w", err)
	}
	return r.buildPlanetAggregateFromDraft(ctx, draftAggregate)
}

func (r *Repository) buildPlanetAggregateFromCourse(ctx context.Context, draft CourseDraft, courseAggregate *CourseAggregate) (*PlanetAggregate, error) {
	item := PlanetListItem{
		ID:                                      courseAggregate.Course.ID,
		DraftID:                                 draft.ID,
		Title:                                   courseAggregate.Course.Title,
		Status:                                  mapDraftStatusToPlanetStatus(draft.Status),
		IsInactive:                              courseAggregate.Course.Status == CourseStatusArchived,
		PlanetTypeID:                            courseAggregate.Course.PlanetTypeID,
		PlanetTypeName:                          courseAggregate.Course.PlanetTypeName,
		PlanetTypeAsset:                         courseAggregate.Course.PlanetTypeAsset,
		PlanetTextureMapID:                      courseAggregate.Course.PlanetTextureMapID,
		PlanetTextureMapName:                    courseAggregate.Course.PlanetTextureMapName,
		PlanetTextureMapAsset:                   courseAggregate.Course.PlanetTextureMapAsset,
		PlanetTextureMapRotationDurationSeconds: courseAggregate.Course.PlanetTextureMapRotationDurationSeconds,
		PlanetTextureMapRotationDirection:       courseAggregate.Course.PlanetTextureMapRotationDirection,
		LastAccessedAt:                          draft.LastAccessedAt,
		UpdatedAt:                               maxTime(courseAggregate.Course.UpdatedAt, draft.UpdatedAt),
		Destination:                             fmt.Sprintf("/dashboard/course-drafts/%s", draft.ID.String()),
	}
	progressSummary, err := r.loadExplorerProgressSummary(ctx, r.pool, draft.ID, draft.Status, draft.ConfirmedCourseID)
	if err != nil {
		return nil, fmt.Errorf("load course planet progress: %w", err)
	}
	if progressSummary != nil {
		applyPlanetProgressSummary(&item, *progressSummary)
	} else {
		totalRegions := 0
		completedRegions := 0
		for _, mainLesson := range courseAggregate.Lessons {
			regionTotal, regionCompleted := countCourseCompletableRegionsInMainLesson(mainLesson)
			totalRegions += regionTotal
			completedRegions += regionCompleted
		}
		applyPlanetProgressSummary(&item, buildPlanetProgressSummary(draft.Status, totalRegions, completedRegions, 0, 0))
	}
	if draft.Status == DraftStatusArchived {
		item.ShareCount = 0
	}

	goalContext, err := r.loadPlanetGoalContext(ctx, draft)
	if err != nil {
		return nil, err
	}
	if goalContext != nil && draft.ConfirmedCourseID != nil {
		goalContext.GoalReadiness = r.loadGoalReadiness(ctx, *draft.ConfirmedCourseID)
	}

	return &PlanetAggregate{
		Planet:      item,
		GoalContext: goalContext,
		Lessons:     courseAggregate.Lessons,
	}, nil
}

func (r *Repository) loadGoalReadiness(ctx context.Context, courseID uuid.UUID) *float64 {
	var totalCompleted, withNote int64
	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE cp.status = 'completed'),
			COUNT(*) FILTER (
				WHERE cp.status = 'completed'
				AND se.goal_alignment_note IS NOT NULL
				AND btrim(se.goal_alignment_note) != ''
			)
		FROM course_points cp
		LEFT JOIN course_point_self_evaluations se ON se.course_point_id = cp.id
		WHERE cp.course_id = $1
	`, courseID).Scan(&totalCompleted, &withNote)
	if err != nil || totalCompleted == 0 {
		return nil
	}
	value := float64(withNote) / float64(totalCompleted) * 100
	return &value
}

func (r *Repository) buildPlanetAggregateFromDraft(ctx context.Context, draftAggregate *DraftAggregate) (*PlanetAggregate, error) {
	item := PlanetListItem{
		ID:                                      draftAggregate.Draft.ID,
		DraftID:                                 draftAggregate.Draft.ID,
		Title:                                   draftAggregate.Draft.Title,
		Status:                                  mapDraftStatusToPlanetStatus(draftAggregate.Draft.Status),
		IsInactive:                              draftAggregate.Draft.IsInactive,
		PlanetTypeID:                            draftAggregate.Draft.PlanetTypeID,
		PlanetTypeName:                          draftAggregate.Draft.PlanetTypeName,
		PlanetTypeAsset:                         draftAggregate.Draft.PlanetTypeAsset,
		PlanetTextureMapID:                      draftAggregate.Draft.PlanetTextureMapID,
		PlanetTextureMapName:                    draftAggregate.Draft.PlanetTextureMapName,
		PlanetTextureMapAsset:                   draftAggregate.Draft.PlanetTextureMapAsset,
		PlanetTextureMapRotationDurationSeconds: draftAggregate.Draft.PlanetTextureMapRotationDurationSeconds,
		PlanetTextureMapRotationDirection:       draftAggregate.Draft.PlanetTextureMapRotationDirection,
		LastAccessedAt:                          draftAggregate.Draft.LastAccessedAt,
		UpdatedAt:                               draftAggregate.Draft.UpdatedAt,
		Destination:                             fmt.Sprintf("/dashboard/course-drafts/%s", draftAggregate.Draft.ID.String()),
	}
	progressSummary, err := r.loadExplorerProgressSummary(ctx, r.pool, draftAggregate.Draft.ID, draftAggregate.Draft.Status, draftAggregate.Draft.ConfirmedCourseID)
	if err != nil {
		return nil, fmt.Errorf("load draft planet progress: %w", err)
	}
	if progressSummary != nil {
		applyPlanetProgressSummary(&item, *progressSummary)
	} else {
		totalRegions := 0
		completedRegions := 0
		for _, mainLesson := range draftAggregate.Lessons {
			regionTotal, regionCompleted := countDraftCompletableRegionsInMainLesson(mainLesson)
			totalRegions += regionTotal
			completedRegions += regionCompleted
		}
		applyPlanetProgressSummary(&item, buildPlanetProgressSummary(draftAggregate.Draft.Status, totalRegions, completedRegions, 0, 0))
	}
	if draftAggregate.Draft.Status == DraftStatusArchived {
		item.ShareCount = 0
	}

	goalContext, err := r.loadPlanetGoalContext(ctx, draftAggregate.Draft)
	if err != nil {
		return nil, err
	}

	return &PlanetAggregate{
		Planet:      item,
		GoalContext: goalContext,
		Lessons:     buildPlanetLessonTreesFromDraftLessons(draftAggregate.Lessons, draftAggregate.Draft.ID),
	}, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
