package curriculum

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateDraft(ctx context.Context, aggregate *DraftAggregate, chargeAmount int) error {
	if aggregate == nil {
		return fmt.Errorf("draft aggregate is nil")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.insertDraftAggregateTx(ctx, tx, aggregate, chargeAmount); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit draft create: %w", err)
	}
	return nil
}

func (r *Repository) CreateDraftAndStartLearning(ctx context.Context, draft *DraftAggregate, confirmed *CourseAggregate, chargeAmount int) error {
	if draft == nil {
		return fmt.Errorf("draft aggregate is nil")
	}
	if confirmed == nil {
		return fmt.Errorf("confirmed course aggregate is nil")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin create-and-start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.insertDraftAggregateTx(ctx, tx, draft, chargeAmount); err != nil {
		return err
	}
	if err := attachGoalToDraftTx(ctx, tx, draft.Draft.UserID, draft.Draft.ID, draft.Draft.GoalProfileID); err != nil {
		return err
	}
	if draft.Draft.PlanetTypeID != nil {
		confirmed.Course.PlanetTypeID = draft.Draft.PlanetTypeID
	}
	if draft.Draft.PlanetTextureMapID != nil {
		confirmed.Course.PlanetTextureMapID = draft.Draft.PlanetTextureMapID
	}
	if err := insertConfirmedCourseTx(ctx, tx, confirmed, draft.Draft.ID, draft.Draft.UserID); err != nil {
		return err
	}
	if err := setDraftLearningStateTx(ctx, tx, draft.Draft.UserID, draft.Draft.ID, confirmed.Course.ID, draft.Draft.SourceQuery); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit create-and-start tx: %w", err)
	}
	return nil
}

func (r *Repository) insertDraftAggregateTx(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate, chargeAmount int) error {
	if err := chargePointsTxWithDetail(ctx, tx, aggregate.Draft.UserID, chargeAmount, "use_course_gen", &PointTransactionDetail{
		Feature:       "course_draft_create",
		ReferenceType: "course_draft",
		ReferenceID:   &aggregate.Draft.ID,
		Description:   strings.TrimSpace(aggregate.Draft.Title),
	}); err != nil {
		return err
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO course_drafts (
			id, user_id, source_query, learning_goal, goal_profile_id, goal_profile_version, current_level,
			duration_weeks, study_hours_per_week, preferred_format,
			generation_language, title, description, completion_criteria, status, planet_type_id, planet_texture_map_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17
		)
	`,
		aggregate.Draft.ID,
		aggregate.Draft.UserID,
		aggregate.Draft.SourceQuery,
		aggregate.Draft.LearningGoal,
		aggregate.Draft.GoalProfileID,
		aggregate.Draft.GoalProfileVersion,
		aggregate.Draft.CurrentLevel,
		aggregate.Draft.DurationWeeks,
		aggregate.Draft.StudyHoursPerWeek,
		aggregate.Draft.PreferredFormat,
		normalizeLearningLanguage(aggregate.Draft.GenerationLanguage),
		aggregate.Draft.Title,
		aggregate.Draft.Description,
		aggregate.Draft.CompletionCriteria,
		aggregate.Draft.Status,
		aggregate.Draft.PlanetTypeID,
		aggregate.Draft.PlanetTextureMapID,
	)
	if err != nil {
		return fmt.Errorf("insert course draft: %w", err)
	}

	insertLesson := func(lesson CourseDraftLesson) error {
		searchSpecJSON, marshalErr := MarshalLessonRecommendationSearchSpec(lesson.RecommendationSearchSpec)
		if marshalErr != nil {
			return fmt.Errorf("marshal lesson recommendation search spec: %w", marshalErr)
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO course_draft_lessons (
				id, course_draft_id, parent_lesson_id, title, objective, summary,
				difficulty_level, lesson_role, source_type, order_index, recommendation_search_spec
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		`,
			lesson.ID,
			aggregate.Draft.ID,
			lesson.ParentLessonID,
			lesson.Title,
			lesson.Objective,
			lesson.Summary,
			lesson.DifficultyLevel,
			lesson.LessonRole,
			lesson.SourceType,
			lesson.OrderIndex,
			searchSpecJSON,
		)
		return err
	}

	for _, mainTree := range aggregate.Lessons {
		if err := insertLesson(mainTree.Lesson); err != nil {
			return fmt.Errorf("insert main lesson: %w", err)
		}
		for _, subTree := range mainTree.SubLessons {
			if err := insertLesson(subTree.Lesson); err != nil {
				return fmt.Errorf("insert sub lesson: %w", err)
			}
		}
	}

	if err := seedExplorerPlanFromDraftLessons(ctx, tx, aggregate); err != nil {
		return fmt.Errorf("seed explorer plan: %w", err)
	}

	eventPayload, err := json.Marshal(map[string]any{
		"title":                aggregate.Draft.Title,
		"main_lesson_count":    len(aggregate.Lessons),
		"goal_profile_id":      aggregate.Draft.GoalProfileID,
		"goal_profile_version": aggregate.Draft.GoalProfileVersion,
		"generation_language":  normalizeLearningLanguage(aggregate.Draft.GenerationLanguage),
	})
	if err != nil {
		return fmt.Errorf("marshal recommendation event payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO recommendation_events (
			user_id, course_draft_id, event_type, source_query, payload
		) VALUES ($1, $2, $3, $4, $5::jsonb)
	`,
		aggregate.Draft.UserID,
		aggregate.Draft.ID,
		EventCurriculumGenerated,
		aggregate.Draft.SourceQuery,
		string(eventPayload),
	)
	if err != nil {
		return fmt.Errorf("insert recommendation event: %w", err)
	}
	return nil
}

func attachGoalToDraftTx(ctx context.Context, tx pgx.Tx, userID, draftID uuid.UUID, goalProfileID *uuid.UUID) error {
	if goalProfileID == nil || *goalProfileID == uuid.Nil {
		return nil
	}

	tag, err := tx.Exec(ctx, `
		UPDATE course_goal_profiles
		SET course_draft_id = $3,
		    updated_at = now()
		WHERE id = $1
		  AND user_id = $2
		  AND is_active = true
		  AND (course_draft_id IS NULL OR course_draft_id = $3)
	`, *goalProfileID, userID, draftID)
	if err != nil {
		return fmt.Errorf("attach goal to draft: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: goal profile is not active or already attached", ErrActiveGoalAlreadyAttached)
	}
	return nil
}

type explorerRegionSeed struct {
	ID          uuid.UUID
	Name        string
	Description *string
	OrderIndex  int
	Nodes       []explorerNodeSeed
	SubRegions  []explorerSubRegionSeed
}

type explorerSubRegionSeed struct {
	ID          uuid.UUID
	Name        string
	Description *string
	OrderIndex  int
	Nodes       []explorerNodeSeed
}

type explorerNodeSeed struct {
	ID           uuid.UUID
	DraftPointID uuid.UUID
	ParentKind   string
	ParentID     uuid.UUID
	NodeType     string
	Title        string
	OrderIndex   int
	SourceType   *string
	SourceURL    *string
	ContentID    *uuid.UUID
	Summary      *string
	ResearchType *string
	LayoutType   *string
	BlockCount   int
}

func seedExplorerPlanFromDraftLessons(ctx context.Context, tx pgx.Tx, aggregate *DraftAggregate) error {
	if aggregate == nil {
		return nil
	}

	for regionIndex, mainTree := range aggregate.Lessons {
		regionSeed := buildExplorerRegionSeed(mainTree, regionIndex)
		if _, err := tx.Exec(ctx, `
			INSERT INTO explorer_regions (id, course_draft_id, name, description, order_index)
			VALUES ($1, $2, $3, $4, $5)
		`,
			regionSeed.ID,
			aggregate.Draft.ID,
			regionSeed.Name,
			regionSeed.Description,
			regionSeed.OrderIndex,
		); err != nil {
			return err
		}

		for _, nodeSeed := range regionSeed.Nodes {
			if _, err := tx.Exec(ctx, `
				INSERT INTO explorer_nodes (
					id, parent_kind, parent_id, node_type, draft_point_id, title, order_index,
					source_type, source_url, content_id, summary,
					research_type, layout_type, block_count
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			`,
				nodeSeed.ID,
				nodeSeed.ParentKind,
				nodeSeed.ParentID,
				nodeSeed.NodeType,
				nodeSeed.DraftPointID,
				nodeSeed.Title,
				nodeSeed.OrderIndex,
				nodeSeed.SourceType,
				nodeSeed.SourceURL,
				nodeSeed.ContentID,
				nodeSeed.Summary,
				nodeSeed.ResearchType,
				nodeSeed.LayoutType,
				nodeSeed.BlockCount,
			); err != nil {
				return err
			}
		}

		for _, subSeed := range regionSeed.SubRegions {
			if _, err := tx.Exec(ctx, `
				INSERT INTO explorer_subregions (id, region_id, name, description, order_index)
				VALUES ($1, $2, $3, $4, $5)
			`,
				subSeed.ID,
				regionSeed.ID,
				subSeed.Name,
				subSeed.Description,
				subSeed.OrderIndex,
			); err != nil {
				return err
			}

			for _, nodeSeed := range subSeed.Nodes {
				if _, err := tx.Exec(ctx, `
					INSERT INTO explorer_nodes (
						id, parent_kind, parent_id, node_type, draft_point_id, title, order_index,
						source_type, source_url, content_id, summary,
						research_type, layout_type, block_count
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
				`,
					nodeSeed.ID,
					nodeSeed.ParentKind,
					nodeSeed.ParentID,
					nodeSeed.NodeType,
					nodeSeed.DraftPointID,
					nodeSeed.Title,
					nodeSeed.OrderIndex,
					nodeSeed.SourceType,
					nodeSeed.SourceURL,
					nodeSeed.ContentID,
					nodeSeed.Summary,
					nodeSeed.ResearchType,
					nodeSeed.LayoutType,
					nodeSeed.BlockCount,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func buildExplorerRegionSeed(mainTree DraftLessonTree, regionIndex int) explorerRegionSeed {
	seed := explorerRegionSeed{
		ID:          mainTree.Lesson.ID,
		Name:        mainTree.Lesson.Title,
		Description: firstNonEmptyString(mainTree.Lesson.Summary, mainTree.Lesson.Objective),
		OrderIndex:  regionIndex,
		Nodes:       buildExplorerNodeSeeds("region", mainTree.Lesson.ID, mainTree.Points),
		SubRegions:  make([]explorerSubRegionSeed, 0, len(mainTree.SubLessons)),
	}
	for subIndex, subTree := range mainTree.SubLessons {
		seed.SubRegions = append(seed.SubRegions, explorerSubRegionSeed{
			ID:          subTree.Lesson.ID,
			Name:        subTree.Lesson.Title,
			Description: firstNonEmptyString(subTree.Lesson.Summary, subTree.Lesson.Objective),
			OrderIndex:  subIndex,
			Nodes:       buildExplorerNodeSeeds("subregion", subTree.Lesson.ID, subTree.Points),
		})
	}
	return seed
}

func buildExplorerNodeSeeds(parentKind string, parentID uuid.UUID, points []DraftPointAggregate) []explorerNodeSeed {
	seeds := make([]explorerNodeSeed, 0, len(points))
	for _, point := range points {
		node := explorerNodeSeed{
			ID:           point.Point.ID,
			DraftPointID: point.Point.ID,
			ParentKind:   parentKind,
			ParentID:     parentID,
			NodeType:     string(point.Point.PointType),
			Title:        point.Point.Title,
			OrderIndex:   point.Point.OrderIndex,
			ContentID:    point.Point.ContentID,
			Summary:      point.Point.Description,
			BlockCount:   len(point.Blocks),
		}
		if point.Point.PointType == PointTypeExploration {
			node.SourceType = inferExplorerSeedSourceType(point.Point)
			node.SourceURL = point.Point.ExternalURL
		} else {
			node.ResearchType = inferExplorerSeedResearchType(point.Point.TemplateType)
			layout := "basic"
			node.LayoutType = &layout
		}
		seeds = append(seeds, node)
	}
	return seeds
}

func inferExplorerSeedSourceType(point CourseDraftPoint) *string {
	if point.ExternalURL != nil && strings.TrimSpace(*point.ExternalURL) != "" {
		sourceType := "web"
		url := strings.ToLower(strings.TrimSpace(*point.ExternalURL))
		if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
			sourceType = "youtube"
		}
		return &sourceType
	}
	if point.ContentID != nil {
		sourceType := "internal"
		return &sourceType
	}
	return nil
}

func inferExplorerSeedResearchType(templateType *ResearchNodeTemplateType) *string {
	if templateType == nil {
		researchType := "free"
		return &researchType
	}
	var researchType string
	switch *templateType {
	case ResearchNodeTemplateConceptSummary:
		researchType = "concept"
	case ResearchNodeTemplatePracticeStrategy:
		researchType = "practice"
	case ResearchNodeTemplateProblemSolving:
		researchType = "problem"
	default:
		researchType = "free"
	}
	return &researchType
}

func firstNonEmptyString(values ...*string) *string {
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}
		return &trimmed
	}
	return nil
}
