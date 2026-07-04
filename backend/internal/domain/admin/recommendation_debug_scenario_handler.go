package admin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

func (h *AdminHandler) ListRecommendationDebugScenarios(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			s.id::text,
			COALESCE(s.created_by::text, ''),
			s.course_title,
			s.initial_user_intent,
			s.status,
			s.notes,
			COALESCE(r.run_count, 0),
			COALESCE(l.label_count, 0),
			s.created_at,
			s.updated_at
		FROM recommendation_debug_scenarios s
		LEFT JOIN (
			SELECT scenario_id, COUNT(*) AS run_count
			FROM recommendation_debug_runs
			GROUP BY scenario_id
		) r ON r.scenario_id = s.id
		LEFT JOIN (
			SELECT scenario_id, COUNT(*) AS label_count
			FROM recommendation_debug_labels
			GROUP BY scenario_id
		) l ON l.scenario_id = s.id
		WHERE s.status <> 'archived'
		ORDER BY s.updated_at DESC
		LIMIT 50
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list scenarios"})
		return
	}
	defer rows.Close()

	items := make([]recommendationDebugScenarioListRow, 0)
	for rows.Next() {
		var item recommendationDebugScenarioListRow
		if err := rows.Scan(
			&item.ID,
			&item.CreatedBy,
			&item.CourseTitle,
			&item.InitialUserIntent,
			&item.Status,
			&item.Notes,
			&item.RunCount,
			&item.LabelCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan scenarios"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list scenarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"scenarios": items})
}

func (h *AdminHandler) CreateRecommendationDebugScenario(c *gin.Context) {
	var req recommendationDebugScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.CourseTitle = strings.TrimSpace(req.CourseTitle)
	req.InitialUserIntent = strings.TrimSpace(req.InitialUserIntent)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.CourseTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_title is required"})
		return
	}
	if len([]rune(req.CourseTitle)) > 160 || len([]rune(req.InitialUserIntent)) > 2000 || len([]rune(req.Notes)) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scenario input too long"})
		return
	}

	adminID, ok := auth.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var createdBy *uuid.UUID
	if parsedAdminID, err := uuid.Parse(adminID); err == nil {
		createdBy = &parsedAdminID
	}

	var item recommendationDebugScenarioDetailRow
	if err := h.db.QueryRow(c.Request.Context(), `
		INSERT INTO recommendation_debug_scenarios (
			created_by,
			course_title,
			initial_user_intent,
			status,
			notes
		)
		VALUES (
			CASE WHEN $1::uuid IS NOT NULL AND EXISTS (SELECT 1 FROM users WHERE id = $1::uuid) THEN $1::uuid ELSE NULL END,
			$2,
			$3,
			'draft',
			$4
		)
		RETURNING
			id::text,
			COALESCE(created_by::text, ''),
			course_title,
			initial_user_intent,
			status,
			COALESCE(goal_profile_snapshot::text, ''),
			COALESCE(generated_lessons_snapshot::text, ''),
			notes,
			created_at,
			updated_at
	`, createdBy, req.CourseTitle, req.InitialUserIntent, req.Notes).Scan(
		&item.ID,
		&item.CreatedBy,
		&item.CourseTitle,
		&item.InitialUserIntent,
		&item.Status,
		&item.GoalProfileSnapshot,
		&item.GeneratedLessonsSnapshot,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		log.Printf("[recommendation-debug] failed to create scenario: admin_id=%s course_title_summary=%s err=%s", adminID, logsafe.Summary(req.CourseTitle), logsafe.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create scenario"})
		return
	}

	item.Runs = []recommendationDebugRunRow{}
	item.Labels = []recommendationDebugLabelRow{}
	c.JSON(http.StatusCreated, gin.H{"scenario": item})
}

func (h *AdminHandler) GetRecommendationDebugScenario(c *gin.Context) {
	scenarioID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scenario id"})
		return
	}

	var item recommendationDebugScenarioDetailRow
	if err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			id::text,
			COALESCE(created_by::text, ''),
			course_title,
			initial_user_intent,
			status,
			COALESCE(goal_profile_snapshot::text, ''),
			COALESCE(generated_lessons_snapshot::text, ''),
			notes,
			created_at,
			updated_at
		FROM recommendation_debug_scenarios
		WHERE id = $1
	`, scenarioID).Scan(
		&item.ID,
		&item.CreatedBy,
		&item.CourseTitle,
		&item.InitialUserIntent,
		&item.Status,
		&item.GoalProfileSnapshot,
		&item.GeneratedLessonsSnapshot,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "scenario not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load scenario"})
		return
	}

	runs, err := h.listRecommendationDebugRuns(c.Request.Context(), scenarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load scenario runs"})
		return
	}
	labels, err := h.listRecommendationDebugLabels(c.Request.Context(), scenarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load scenario labels"})
		return
	}
	item.Runs = runs
	item.Labels = labels

	c.JSON(http.StatusOK, gin.H{"scenario": item})
}
