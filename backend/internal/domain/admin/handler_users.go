package admin

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/auth"
)

type userRow struct {
	ID                 string  `json:"id"`
	Email              string  `json:"email"`
	Nickname           string  `json:"nickname"`
	Role               string  `json:"role"`
	PremiumAccess      bool    `json:"premium_access"`
	Provider           string  `json:"provider"`
	DisplayID          *string `json:"display_id"`
	CreatedAt          string  `json:"created_at"`
	TOTPResetRequested bool    `json:"totp_reset_requested"`
}

func (h *AdminHandler) listUsersHandler(c *gin.Context) {
	q := "%" + c.Query("q") + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	err := h.db.QueryRow(c.Request.Context(), `
		SELECT COUNT(*)
		FROM users u
		LEFT JOIN LATERAL (
			SELECT provider, provider_id
			FROM social_accounts
			WHERE user_id = u.id
			ORDER BY created_at DESC
			LIMIT 1
		) sa ON TRUE
		WHERE COALESCE(u.email,'') ILIKE $1
		   OR COALESCE(u.nickname,'') ILIKE $1
		   OR COALESCE(u.display_id,'') ILIKE $1
		   OR COALESCE(sa.provider,'') ILIKE $1
		   OR COALESCE(sa.provider_id,'') ILIKE $1
	`, q).Scan(&total)
	if err != nil {
		log.Printf("[super-admin/users] count query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	rows, err := h.db.Query(c.Request.Context(), `
		SELECT
			u.id,
			COALESCE(u.email, '') AS email,
			COALESCE(u.nickname, '') AS nickname,
			u.role,
			COALESCE(u.premium_access, false) AS premium_access,
			u.display_id,
			COALESCE(sa.provider, 'unknown') AS provider,
			u.created_at,
			COALESCE(u.totp_reset_requested, false) AS totp_reset_requested
		FROM users u
		LEFT JOIN LATERAL (
			SELECT provider, provider_id
			FROM social_accounts
			WHERE user_id = u.id
			ORDER BY created_at DESC
			LIMIT 1
		) sa ON TRUE
		WHERE COALESCE(u.email,'') ILIKE $1
		   OR COALESCE(u.nickname,'') ILIKE $1
		   OR COALESCE(u.display_id,'') ILIKE $1
		   OR COALESCE(sa.provider,'') ILIKE $1
		   OR COALESCE(sa.provider_id,'') ILIKE $1
		ORDER BY u.totp_reset_requested DESC, u.created_at DESC
		LIMIT $2 OFFSET $3
	`, q, limit, offset)
	if err != nil {
		log.Printf("[super-admin/users] list query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	users := []userRow{}
	for rows.Next() {
		var u userRow
		var createdAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Nickname, &u.Role, &u.PremiumAccess, &u.DisplayID, &u.Provider, &createdAt, &u.TOTPResetRequested); err != nil {
			log.Printf("[super-admin/users] row scan failed: %v", err)
			continue
		}
		u.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[super-admin/users] rows iteration failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *AdminHandler) updateRoleHandler(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if body.Role == string(auth.RoleSuperAdmin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot set super_admin role"})
		return
	}
	if body.Role != string(auth.RoleLearner) && body.Role != string(auth.RoleAdmin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	_, err := h.db.Exec(c.Request.Context(), `
		UPDATE users SET role=$1 WHERE id=$2
	`, body.Role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *AdminHandler) updatePremiumAccessHandler(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		PremiumAccess bool `json:"premium_access"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if _, err := h.db.Exec(c.Request.Context(), `
		UPDATE users SET premium_access = $1 WHERE id = $2
	`, body.PremiumAccess, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// GET /api/v1/admin/users — 학습자 목록 (포인트 잔액 포함)
func (h *AdminHandler) adminListUsersHandler(c *gin.Context) {
	q := "%" + c.Query("q") + "%"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	ctx := c.Request.Context()

	var total int
	h.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM users u
		LEFT JOIN LATERAL (
			SELECT provider, provider_id
			FROM social_accounts
			WHERE user_id = u.id
			ORDER BY created_at DESC
			LIMIT 1
		) sa ON TRUE
		WHERE COALESCE(u.email,'') ILIKE $1
		   OR COALESCE(u.nickname,'') ILIKE $1
		   OR COALESCE(u.display_id,'') ILIKE $1
		   OR COALESCE(sa.provider,'') ILIKE $1
		   OR COALESCE(sa.provider_id,'') ILIKE $1
	`, q).Scan(&total)

	rows, err := h.db.Query(ctx, `
		SELECT u.id,
		       COALESCE(u.email, '') AS email,
		       COALESCE(u.nickname, '') AS nickname,
		       u.role,
		       COALESCE(u.premium_access, false) AS premium_access,
		       COALESCE(sa.provider, 'unknown') AS provider,
		       u.display_id,
		       COALESCE(u.totp_reset_requested, false) AS totp_reset_requested,
		       u.created_at,
		       COALESCE(w.free_balance, 0) + COALESCE(w.paid_balance, 0) AS points
		FROM users u
		LEFT JOIN LATERAL (
			SELECT provider, provider_id
			FROM social_accounts
			WHERE user_id = u.id
			ORDER BY created_at DESC
			LIMIT 1
		) sa ON TRUE
		LEFT JOIN ai_point_wallets w ON w.user_id = u.id
		WHERE COALESCE(u.email,'') ILIKE $1
		   OR COALESCE(u.nickname,'') ILIKE $1
		   OR COALESCE(u.display_id,'') ILIKE $1
		   OR COALESCE(sa.provider,'') ILIKE $1
		   OR COALESCE(sa.provider_id,'') ILIKE $1
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3
	`, q, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	type adminUserRow struct {
		ID                 string  `json:"id"`
		Email              string  `json:"email"`
		Nickname           string  `json:"nickname"`
		Role               string  `json:"role"`
		PremiumAccess      bool    `json:"premium_access"`
		Provider           string  `json:"provider"`
		DisplayID          *string `json:"display_id"`
		TOTPResetRequested bool    `json:"totp_reset_requested"`
		Points             int     `json:"points"`
		CreatedAt          string  `json:"created_at"`
	}

	users := []adminUserRow{}
	for rows.Next() {
		var u adminUserRow
		var createdAt time.Time
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Nickname,
			&u.Role,
			&u.PremiumAccess,
			&u.Provider,
			&u.DisplayID,
			&u.TOTPResetRequested,
			&createdAt,
			&u.Points,
		); err != nil {
			continue
		}
		u.CreatedAt = createdAt.Format("2006-01-02")
		users = append(users, u)
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "total": total, "page": page, "limit": limit})
}
