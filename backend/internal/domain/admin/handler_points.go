package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/auth"
)

// ─── 포인트 지급 ───

type GrantPointRequest struct {
	Amount int    `json:"amount" binding:"required"`
	Memo   string `json:"memo"`
}

// POST /api/v1/super-admin/users/:id/grant-points
func (h *AdminHandler) GrantUserPoints(c *gin.Context) {
	targetUserID := c.Param("id")

	var req GrantPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청: " + err.Error()})
		return
	}
	if req.Amount == 0 || req.Amount < -50 || req.Amount > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "차감/지급량은 -50~50pt 범위이며 0은 불가합니다."})
		return
	}

	ctx := c.Request.Context()

	// 대상 사용자 존재 확인 + 식별자 조회
	var targetIdentifier string
	err := h.db.QueryRow(ctx,
		`SELECT COALESCE(NULLIF(email,''), display_id::text, id::text) FROM users WHERE id = $1`,
		targetUserID,
	).Scan(&targetIdentifier)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다."})
		return
	}

	// 트랜잭션으로 포인트 지급 + 로그 기록
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "트랜잭션 시작 실패"})
		return
	}
	defer tx.Rollback(ctx)

	if err := h.applyPaidPointDelta(ctx, tx, targetUserID, req.Amount, req.Memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ② point_grant_logs 기록
	adminID, _ := auth.GetCurrentUserID(c)
	adminEmail := c.GetString("user_email")

	_, err = tx.Exec(ctx,
		`INSERT INTO point_grant_logs
		   (admin_id, admin_email, target_user_id, target_identifier, amount, memo)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		adminID, adminEmail, targetUserID, targetIdentifier, req.Amount, req.Memo,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "로그 기록 실패: " + err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "커밋 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "포인트가 지급되었습니다.",
		"amount":  req.Amount,
		"user_id": targetUserID,
	})
}

// GET /api/v1/super-admin/audit/point-grants
func (h *AdminHandler) GetPointGrantLogs(c *gin.Context) {
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
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM point_grant_logs`).Scan(&total)

	rows, err := h.db.Query(ctx,
		`SELECT pgl.id, pgl.admin_id,
		        COALESCE(
		            NULLIF(u.nickname,''),
		            NULLIF(u.email,''),
		            pgl.admin_email,
		            pgl.admin_id::text
		        ) AS admin_display,
		        pgl.target_user_id, pgl.target_identifier,
		        pgl.amount, COALESCE(pgl.memo,''), pgl.created_at
		 FROM   point_grant_logs pgl
		 LEFT JOIN users u ON u.id::text = pgl.admin_id
		 ORDER  BY pgl.created_at DESC
		 LIMIT  $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
		return
	}
	defer rows.Close()

	type LogRow struct {
		ID               string    `json:"id"`
		AdminID          string    `json:"admin_id"`
		AdminDisplay     string    `json:"admin_display"`
		TargetUserID     string    `json:"target_user_id"`
		TargetIdentifier string    `json:"target_identifier"`
		Amount           int       `json:"amount"`
		Memo             string    `json:"memo"`
		CreatedAt        time.Time `json:"created_at"`
	}

	logs := []LogRow{}
	for rows.Next() {
		var l LogRow
		if err := rows.Scan(
			&l.ID, &l.AdminID, &l.AdminDisplay,
			&l.TargetUserID, &l.TargetIdentifier,
			&l.Amount, &l.Memo, &l.CreatedAt,
		); err != nil {
			continue
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GET /api/v1/admin/point-policy — 관리자용 포인트 정책 조회 (admin_max_grant 포함)
func (h *AdminHandler) adminGetPointPolicyHandler(c *gin.Context) {
	var adminMaxGrant int
	err := h.db.QueryRow(c.Request.Context(),
		`SELECT value FROM point_settings WHERE key = 'admin_max_grant'`,
	).Scan(&adminMaxGrant)
	if err != nil {
		adminMaxGrant = 50 // 기본값
	}
	c.JSON(http.StatusOK, gin.H{"admin_max_grant": adminMaxGrant})
}

// POST /api/v1/admin/users/:id/grant-points — 관리자가 포인트 지급 (정책 상한 적용)
func (h *AdminHandler) adminGrantPointsHandler(c *gin.Context) {
	targetUserID := c.Param("id")

	var req struct {
		Amount int    `json:"amount" binding:"required"`
		Memo   string `json:"memo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청: " + err.Error()})
		return
	}
	if req.Amount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "0pt는 처리할 수 없습니다."})
		return
	}

	ctx := c.Request.Context()

	// 정책 상한 조회
	var adminMaxGrant int
	if err := h.db.QueryRow(ctx, `SELECT value FROM point_settings WHERE key = 'admin_max_grant'`).Scan(&adminMaxGrant); err != nil {
		adminMaxGrant = 50
	}
	if req.Amount > adminMaxGrant || req.Amount < -adminMaxGrant {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("차감/지급량은 -%d~%dpt 범위여야 합니다.", adminMaxGrant, adminMaxGrant)})
		return
	}

	var targetIdentifier string
	err := h.db.QueryRow(ctx,
		`SELECT COALESCE(NULLIF(email,''), display_id::text, id::text) FROM users WHERE id = $1`,
		targetUserID,
	).Scan(&targetIdentifier)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다."})
		return
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "트랜잭션 시작 실패"})
		return
	}
	defer tx.Rollback(ctx)

	if err := h.applyPaidPointDelta(ctx, tx, targetUserID, req.Amount, req.Memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, _ := auth.GetCurrentUserID(c)
	adminEmail := c.GetString("user_email")

	_, err = tx.Exec(ctx,
		`INSERT INTO point_grant_logs (admin_id, admin_email, target_user_id, target_identifier, amount, memo)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		adminID, adminEmail, targetUserID, targetIdentifier, req.Amount, req.Memo,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "로그 기록 실패"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "커밋 실패"})
		return
	}

	message := "포인트가 지급되었습니다."
	if req.Amount < 0 {
		message = "포인트가 차감되었습니다."
	}
	c.JSON(http.StatusOK, gin.H{"message": message, "amount": req.Amount})
}

func (h *AdminHandler) applyPaidPointDelta(ctx context.Context, tx pgx.Tx, targetUserID string, amount int, memo string) error {
	var freeBalance, paidBalance int
	err := tx.QueryRow(ctx, `
		SELECT COALESCE(free_balance, 0), COALESCE(paid_balance, 0)
		FROM ai_point_wallets
		WHERE user_id = $1
	`, targetUserID).Scan(&freeBalance, &paidBalance)
	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("포인트 조회 실패")
	}

	if err == pgx.ErrNoRows {
		if amount < 0 {
			return fmt.Errorf("차감할 지급 포인트가 없습니다.")
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO ai_point_wallets (id, user_id, free_balance, paid_balance, updated_at)
			VALUES ($1, $2, 0, $3, NOW())
		`, uuid.New().String(), targetUserID, amount)
		if err != nil {
			return fmt.Errorf("포인트 업데이트 실패")
		}
		return h.insertAdminPointTransaction(ctx, tx, targetUserID, amount, memo)
	}

	newPaidBalance := paidBalance + amount
	if newPaidBalance < 0 {
		return fmt.Errorf("지급 포인트보다 많이 차감할 수 없습니다.")
	}

	_, err = tx.Exec(ctx, `
		UPDATE ai_point_wallets
		SET paid_balance = $1, updated_at = NOW()
		WHERE user_id = $2
	`, newPaidBalance, targetUserID)
	if err != nil {
		return fmt.Errorf("포인트 업데이트 실패")
	}
	return h.insertAdminPointTransaction(ctx, tx, targetUserID, amount, memo)
}

func (h *AdminHandler) insertAdminPointTransaction(ctx context.Context, tx pgx.Tx, targetUserID string, amount int, memo string) error {
	transactionType := "purchase"
	description := "관리자 포인트 지급"
	if amount < 0 {
		description = "관리자 포인트 차감"
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO ai_point_transactions (id, user_id, type, amount, feature, reference_type, reference_id, description, metadata, created_at)
		VALUES (
			$1, $2, $3, $4, 'admin_point_adjustment', 'admin_action', NULL, $5,
			CASE WHEN $6 = '' THEN '{}'::jsonb ELSE jsonb_build_object('memo', $6) END,
			NOW()
		)
	`, uuid.New().String(), targetUserID, transactionType, amount, description, strings.TrimSpace(memo))
	if err != nil {
		return fmt.Errorf("포인트 거래 기록 실패")
	}
	return nil
}

// GET /api/v1/public/point-policy — 인증 불필요, 공개 포인트 정책 조회
func (h *AdminHandler) publicPointPolicyHandler(c *gin.Context) {
	ctx := c.Request.Context()
	getValue := func(key string, def int) int {
		var v int
		if err := h.db.QueryRow(ctx, `SELECT value FROM point_settings WHERE key = $1`, key).Scan(&v); err != nil {
			return def
		}
		return v
	}
	c.JSON(http.StatusOK, gin.H{
		"welcome_points":  getValue("welcome_points", 30),
		"course_gen_cost": getValue("course_gen_cost", 5),
		"lesson_rec_cost": getValue("lesson_rec_cost", 1),
	})
}
