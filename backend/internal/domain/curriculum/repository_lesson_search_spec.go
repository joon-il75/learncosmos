package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LessonRecommendationSearchContext struct {
	Spec             LessonRecommendationSearchSpec
	SourceQuery      string
	LearningGoal     *string
	LearningLanguage string
	LessonTitle      string
	LessonObjective  *string
	OrderIndex       int
	LearningIntent   LearningIntentProfile
}

func (ctx LessonRecommendationSearchContext) FallbackSpec() LessonRecommendationSearchSpec {
	return BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(
		ctx.SourceQuery,
		ctx.LearningGoal,
		ctx.LearningLanguage,
		ctx.LessonTitle,
		ctx.LessonObjective,
		ctx.OrderIndex,
		ctx.LearningIntent,
	)
}

func (r *Repository) GetDraftLessonRecommendationSearchSpecByPointID(ctx context.Context, userID, pointID uuid.UUID) (LessonRecommendationSearchSpec, error) {
	var raw []byte
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(cdl.recommendation_search_spec, '{}'::jsonb)
		FROM course_draft_points cdp
		JOIN course_draft_lessons cdl ON cdl.id = cdp.course_draft_lesson_id
		JOIN course_drafts cd ON cd.id = cdp.course_draft_id
		WHERE cdp.id = $1
		  AND cd.user_id = $2
	`, pointID, userID).Scan(&raw)
	if err == nil {
		return ParseLessonRecommendationSearchSpec(raw), nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return LessonRecommendationSearchSpec{}, fmt.Errorf("get draft lesson recommendation search spec: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(
			NULLIF(cl.recommendation_search_spec, '{}'::jsonb),
			cdl.recommendation_search_spec,
			'{}'::jsonb
		)
		FROM course_points cp
		JOIN course_lessons cl ON cl.id = cp.course_lesson_id
		JOIN courses c ON c.id = cp.course_id
		LEFT JOIN course_draft_lessons cdl
		  ON cdl.course_draft_id = c.source_draft_id
		 AND cdl.order_index = cl.order_index
		 AND cdl.title = cl.title
		WHERE cp.id = $1
		  AND c.user_id = $2
	`, pointID, userID).Scan(&raw)
	if err != nil {
		return LessonRecommendationSearchSpec{}, fmt.Errorf("get lesson recommendation search spec: %w", err)
	}
	return ParseLessonRecommendationSearchSpec(raw), nil
}

func (r *Repository) GetLessonRecommendationSearchContextByPointID(ctx context.Context, userID, pointID uuid.UUID) (LessonRecommendationSearchContext, error) {
	var result LessonRecommendationSearchContext
	var raw []byte
	var rawProfile []byte
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(cdl.recommendation_search_spec, '{}'::jsonb),
		       cd.source_query, cd.learning_goal, COALESCE(cd.generation_language, 'ko'),
		       cdl.title, cdl.objective, cdl.order_index,
		       COALESCE(cgp.learning_intent_profile, '{}'::jsonb)
		FROM course_draft_points cdp
		JOIN course_draft_lessons cdl ON cdl.id = cdp.course_draft_lesson_id
		JOIN course_drafts cd ON cd.id = cdp.course_draft_id
		LEFT JOIN LATERAL (
			SELECT gp.learning_intent_profile
			FROM course_goal_profiles gp
			WHERE gp.user_id = cd.user_id
			  AND (gp.id = cd.goal_profile_id OR gp.course_draft_id = cd.id)
			ORDER BY CASE WHEN gp.id = cd.goal_profile_id THEN 0 ELSE 1 END, gp.is_active DESC, gp.updated_at DESC
			LIMIT 1
		) cgp ON true
		WHERE cdp.id = $1
		  AND cd.user_id = $2
	`, pointID, userID).Scan(
		&raw,
		&result.SourceQuery,
		&result.LearningGoal,
		&result.LearningLanguage,
		&result.LessonTitle,
		&result.LessonObjective,
		&result.OrderIndex,
		&rawProfile,
	)
	if err == nil {
		result.Spec = ParseLessonRecommendationSearchSpec(raw)
		result.LearningIntent = ParseLearningIntentProfile(rawProfile)
		return result, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return LessonRecommendationSearchContext{}, fmt.Errorf("get draft lesson recommendation search context: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(
			NULLIF(cl.recommendation_search_spec, '{}'::jsonb),
			cdl.recommendation_search_spec,
			'{}'::jsonb
		),
		       c.source_query, c.learning_goal, COALESCE(c.generation_language, 'ko'),
		       cl.title, cl.objective, cl.order_index,
		       COALESCE(cgp.learning_intent_profile, '{}'::jsonb)
		FROM course_points cp
		JOIN course_lessons cl ON cl.id = cp.course_lesson_id
		JOIN courses c ON c.id = cp.course_id
		LEFT JOIN course_draft_lessons cdl
		  ON cdl.course_draft_id = c.source_draft_id
		 AND cdl.order_index = cl.order_index
		 AND cdl.title = cl.title
		LEFT JOIN course_drafts source_cd ON source_cd.id = c.source_draft_id
		LEFT JOIN LATERAL (
			SELECT gp.learning_intent_profile
			FROM course_goal_profiles gp
			WHERE gp.user_id = c.user_id
			  AND (gp.id = source_cd.goal_profile_id OR gp.course_draft_id = c.source_draft_id)
			ORDER BY CASE WHEN gp.id = source_cd.goal_profile_id THEN 0 ELSE 1 END, gp.is_active DESC, gp.updated_at DESC
			LIMIT 1
		) cgp ON true
		WHERE cp.id = $1
		  AND c.user_id = $2
	`, pointID, userID).Scan(
		&raw,
		&result.SourceQuery,
		&result.LearningGoal,
		&result.LearningLanguage,
		&result.LessonTitle,
		&result.LessonObjective,
		&result.OrderIndex,
		&rawProfile,
	)
	if err != nil {
		return LessonRecommendationSearchContext{}, fmt.Errorf("get lesson recommendation search context: %w", err)
	}
	result.Spec = ParseLessonRecommendationSearchSpec(raw)
	result.LearningIntent = ParseLearningIntentProfile(rawProfile)
	return result, nil
}

func (r *Repository) ResolveDraftLessonIDByExplorerSubRegionID(ctx context.Context, userID, subRegionID uuid.UUID) (uuid.UUID, error) {
	var lessonID uuid.UUID
	err := r.pool.QueryRow(ctx, `
		SELECT es.course_draft_lesson_id
		FROM explorer_subregions es
		JOIN explorer_regions er ON er.id = es.region_id
		JOIN course_drafts cd ON cd.id = er.course_draft_id
		WHERE es.id = $1
		  AND cd.user_id = $2
		  AND es.course_draft_lesson_id IS NOT NULL
	`, subRegionID, userID).Scan(&lessonID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("resolve draft lesson id by explorer subregion: %w", err)
	}
	return lessonID, nil
}

func (r *Repository) ResolveDraftLessonIDByExplorerRegionID(ctx context.Context, userID, regionID uuid.UUID) (uuid.UUID, error) {
	var lessonID uuid.UUID
	err := r.pool.QueryRow(ctx, `
		SELECT er.course_draft_lesson_id
		FROM explorer_regions er
		JOIN course_drafts cd ON cd.id = er.course_draft_id
		WHERE er.id = $1
		  AND cd.user_id = $2
		  AND er.course_draft_lesson_id IS NOT NULL
	`, regionID, userID).Scan(&lessonID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("resolve draft lesson id by explorer region: %w", err)
	}
	return lessonID, nil
}

func (r *Repository) GetLessonRecommendationSearchContextByLessonID(ctx context.Context, userID, lessonID uuid.UUID) (LessonRecommendationSearchContext, error) {
	var result LessonRecommendationSearchContext
	var raw []byte
	var rawProfile []byte
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(cdl.recommendation_search_spec, '{}'::jsonb),
		       cd.source_query, cd.learning_goal, COALESCE(cd.generation_language, 'ko'),
		       cdl.title, cdl.objective, cdl.order_index,
		       COALESCE(cgp.learning_intent_profile, '{}'::jsonb)
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		LEFT JOIN LATERAL (
			SELECT gp.learning_intent_profile
			FROM course_goal_profiles gp
			WHERE gp.user_id = cd.user_id
			  AND (gp.id = cd.goal_profile_id OR gp.course_draft_id = cd.id)
			ORDER BY CASE WHEN gp.id = cd.goal_profile_id THEN 0 ELSE 1 END, gp.is_active DESC, gp.updated_at DESC
			LIMIT 1
		) cgp ON true
		WHERE cdl.id = $1
		  AND cd.user_id = $2
	`, lessonID, userID).Scan(
		&raw,
		&result.SourceQuery,
		&result.LearningGoal,
		&result.LearningLanguage,
		&result.LessonTitle,
		&result.LessonObjective,
		&result.OrderIndex,
		&rawProfile,
	)
	if err == nil {
		result.Spec = ParseLessonRecommendationSearchSpec(raw)
		result.LearningIntent = ParseLearningIntentProfile(rawProfile)
		return result, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return LessonRecommendationSearchContext{}, fmt.Errorf("get draft lesson recommendation search context by lesson: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(
			NULLIF(cl.recommendation_search_spec, '{}'::jsonb),
			cdl.recommendation_search_spec,
			'{}'::jsonb
		),
		       c.source_query, c.learning_goal, COALESCE(c.generation_language, 'ko'),
		       cl.title, cl.objective, cl.order_index,
		       COALESCE(cgp.learning_intent_profile, '{}'::jsonb)
		FROM course_lessons cl
		JOIN courses c ON c.id = cl.course_id
		LEFT JOIN course_draft_lessons cdl
		  ON cdl.course_draft_id = c.source_draft_id
		 AND cdl.order_index = cl.order_index
		 AND cdl.title = cl.title
		LEFT JOIN course_drafts source_cd ON source_cd.id = c.source_draft_id
		LEFT JOIN LATERAL (
			SELECT gp.learning_intent_profile
			FROM course_goal_profiles gp
			WHERE gp.user_id = c.user_id
			  AND (gp.id = source_cd.goal_profile_id OR gp.course_draft_id = c.source_draft_id)
			ORDER BY CASE WHEN gp.id = source_cd.goal_profile_id THEN 0 ELSE 1 END, gp.is_active DESC, gp.updated_at DESC
			LIMIT 1
		) cgp ON true
		WHERE cl.id = $1
		  AND c.user_id = $2
	`, lessonID, userID).Scan(
		&raw,
		&result.SourceQuery,
		&result.LearningGoal,
		&result.LearningLanguage,
		&result.LessonTitle,
		&result.LessonObjective,
		&result.OrderIndex,
		&rawProfile,
	)
	if err != nil {
		return LessonRecommendationSearchContext{}, fmt.Errorf("get course lesson recommendation search context by lesson: %w", err)
	}
	result.Spec = ParseLessonRecommendationSearchSpec(raw)
	result.LearningIntent = ParseLearningIntentProfile(rawProfile)
	return result, nil
}

func ParseLearningIntentProfile(raw []byte) LearningIntentProfile {
	if len(raw) == 0 || string(raw) == "{}" {
		return LearningIntentProfile{}
	}
	var profile LearningIntentProfile
	if err := json.Unmarshal(raw, &profile); err != nil {
		return LearningIntentProfile{}
	}
	return NormalizeLearningIntentProfile(profile)
}

func (r *Repository) SaveFallbackLessonRecommendationSearchSpecByPointID(ctx context.Context, userID, pointID uuid.UUID, spec LessonRecommendationSearchSpec) error {
	searchSpecJSON, err := MarshalLessonRecommendationSearchSpec(spec)
	if err != nil {
		return fmt.Errorf("marshal fallback lesson recommendation search spec: %w", err)
	}
	if searchSpecJSON == "{}" {
		return nil
	}

	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_lessons cdl
		SET recommendation_search_spec = $3::jsonb,
		    updated_at = NOW()
		FROM course_draft_points cdp
		JOIN course_drafts cd ON cd.id = cdp.course_draft_id
		WHERE cdl.id = cdp.course_draft_lesson_id
		  AND cdp.id = $1
		  AND cd.user_id = $2
		  AND COALESCE(cdl.recommendation_search_spec, '{}'::jsonb) = '{}'::jsonb
	`, pointID, userID, searchSpecJSON)
	if err != nil {
		return fmt.Errorf("save draft fallback lesson recommendation search spec: %w", err)
	}
	if tag.RowsAffected() > 0 {
		return nil
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE course_lessons cl
		SET recommendation_search_spec = $3::jsonb,
		    updated_at = NOW()
		FROM course_points cp
		JOIN courses c ON c.id = cp.course_id
		WHERE cl.id = cp.course_lesson_id
		  AND cp.id = $1
		  AND c.user_id = $2
		  AND COALESCE(cl.recommendation_search_spec, '{}'::jsonb) = '{}'::jsonb
	`, pointID, userID, searchSpecJSON)
	if err != nil {
		return fmt.Errorf("save course fallback lesson recommendation search spec: %w", err)
	}
	return nil
}

func (r *Repository) SaveFallbackLessonRecommendationSearchSpecByLessonID(ctx context.Context, userID, lessonID uuid.UUID, spec LessonRecommendationSearchSpec) error {
	searchSpecJSON, err := MarshalLessonRecommendationSearchSpec(spec)
	if err != nil {
		return fmt.Errorf("marshal fallback lesson recommendation search spec: %w", err)
	}
	if searchSpecJSON == "{}" {
		return nil
	}

	tag, err := r.pool.Exec(ctx, `
		UPDATE course_draft_lessons cdl
		SET recommendation_search_spec = $3::jsonb,
		    updated_at = NOW()
		FROM course_drafts cd
		WHERE cd.id = cdl.course_draft_id
		  AND cdl.id = $1
		  AND cd.user_id = $2
		  AND COALESCE(cdl.recommendation_search_spec, '{}'::jsonb) = '{}'::jsonb
	`, lessonID, userID, searchSpecJSON)
	if err != nil {
		return fmt.Errorf("save draft fallback lesson recommendation search spec by lesson: %w", err)
	}
	if tag.RowsAffected() > 0 {
		return nil
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE course_lessons cl
		SET recommendation_search_spec = $3::jsonb,
		    updated_at = NOW()
		FROM courses c
		WHERE c.id = cl.course_id
		  AND cl.id = $1
		  AND c.user_id = $2
		  AND COALESCE(cl.recommendation_search_spec, '{}'::jsonb) = '{}'::jsonb
	`, lessonID, userID, searchSpecJSON)
	if err != nil {
		return fmt.Errorf("save course fallback lesson recommendation search spec by lesson: %w", err)
	}
	return nil
}

func (r *Repository) ResolveSubregionDraftLessonIDByID(ctx context.Context, userID, lessonOrSubregionID uuid.UUID) (uuid.UUID, error) {
	var resolved uuid.UUID
	err := r.pool.QueryRow(ctx, `
		SELECT cdl.id
		FROM explorer_subregions es
		JOIN explorer_regions er ON es.region_id = er.id
		JOIN course_draft_lessons cdl ON cdl.id = es.course_draft_lesson_id
		JOIN course_drafts cd ON cd.id = er.course_draft_id
		WHERE es.id = $1
		  AND cd.user_id = $2
	`, lessonOrSubregionID, userID).Scan(&resolved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("resolve subregion draft lesson id: %w", err)
	}
	return resolved, nil
}
