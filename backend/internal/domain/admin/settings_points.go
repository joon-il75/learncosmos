package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── 포인트 정책 ───

// GET /api/v1/super-admin/settings/point-policy
func (h *AdminHandler) GetPointPolicy(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `SELECT key, value FROM point_settings`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	policy := map[string]int{}
	for rows.Next() {
		var key string
		var value int
		rows.Scan(&key, &value)
		policy[key] = value
	}
	c.JSON(http.StatusOK, gin.H{"policy": policy})
}

// PUT /api/v1/super-admin/settings/point-policy
func (h *AdminHandler) UpdatePointPolicy(c *gin.Context) {
	var req struct {
		WelcomePoints    int `json:"welcome_points"`
		CourseGenCost    int `json:"course_gen_cost"`
		LessonRecCost    int `json:"lesson_rec_cost"`
		AdminMaxGrant    int `json:"admin_max_grant"`
		ProMonthlyPoints int `json:"pro_monthly_points"`
		TutorCost        int `json:"tutor_cost"`
		VisionCost       int `json:"vision_cost"`
		QuizCost         int `json:"quiz_cost"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]int{
		"welcome_points":     req.WelcomePoints,
		"course_gen_cost":    req.CourseGenCost,
		"lesson_rec_cost":    req.LessonRecCost,
		"admin_max_grant":    req.AdminMaxGrant,
		"pro_monthly_points": req.ProMonthlyPoints,
		"tutor_cost":         req.TutorCost,
		"vision_cost":        req.VisionCost,
		"quiz_cost":          req.QuizCost,
	}
	for key, val := range updates {
		h.db.Exec(c.Request.Context(), `UPDATE point_settings SET value=$1, updated_at=NOW() WHERE key=$2`, val, key)
	}
	c.JSON(http.StatusOK, gin.H{"message": "포인트 정책이 저장되었습니다"})
}

// ─── 구매 상품 ───

// GET /api/v1/super-admin/settings/products
func (h *AdminHandler) GetPointProducts(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, name, points, price_krw, is_active, sort_order,
		       product_type, billing_period, monthly_points, ad_free, premium_access, description
		FROM point_products ORDER BY sort_order
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var id, name string
		var points, price, sort int
		var isActive bool
		var productType string
		var billingPeriod, description *string
		var monthlyPoints int
		var adFree, premiumAccess bool
		rows.Scan(&id, &name, &points, &price, &isActive, &sort, &productType, &billingPeriod, &monthlyPoints, &adFree, &premiumAccess, &description)
		products = append(products, map[string]interface{}{
			"id": id, "name": name, "points": points,
			"price_krw": price, "is_active": isActive, "sort_order": sort,
			"product_type": productType, "billing_period": billingPeriod,
			"monthly_points": monthlyPoints, "ad_free": adFree,
			"premium_access": premiumAccess, "description": description,
		})
	}
	if products == nil {
		products = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// PUT /api/v1/super-admin/settings/products
func (h *AdminHandler) SavePointProducts(c *gin.Context) {
	var req struct {
		Products []struct {
			ID            *string `json:"id"`
			Name          string  `json:"name"`
			Points        int     `json:"points"`
			PriceKRW      int     `json:"price_krw"`
			IsActive      bool    `json:"is_active"`
			SortOrder     int     `json:"sort_order"`
			ProductType   string  `json:"product_type"`
			BillingPeriod *string `json:"billing_period"`
			MonthlyPoints int     `json:"monthly_points"`
			AdFree        bool    `json:"ad_free"`
			PremiumAccess bool    `json:"premium_access"`
			Description   *string `json:"description"`
		} `json:"products"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, p := range req.Products {
		if p.ProductType == "" {
			p.ProductType = "one_time"
		}
		if p.ID == nil {
			h.db.Exec(c.Request.Context(), `
				INSERT INTO point_products (
					name, points, price_krw, is_active, sort_order,
					product_type, billing_period, monthly_points, ad_free, premium_access, description
				)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			`, p.Name, p.Points, p.PriceKRW, p.IsActive, i+1, p.ProductType, p.BillingPeriod, p.MonthlyPoints, p.AdFree, p.PremiumAccess, p.Description)
		} else {
			h.db.Exec(c.Request.Context(), `
				UPDATE point_products
				SET name=$1, points=$2, price_krw=$3, is_active=$4, sort_order=$5,
				    product_type=$6, billing_period=$7, monthly_points=$8,
				    ad_free=$9, premium_access=$10, description=$11
				WHERE id=$12
			`, p.Name, p.Points, p.PriceKRW, p.IsActive, i+1, p.ProductType, p.BillingPeriod, p.MonthlyPoints, p.AdFree, p.PremiumAccess, p.Description, *p.ID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "상품 목록이 저장되었습니다"})
}
