package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── 광고 슬롯 ───

// GET /api/v1/super-admin/settings/ad-slots
func (h *AdminHandler) GetAdSlots(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, slot_type, title, link_url, category_slug, is_active, sort_order
		FROM ad_slots ORDER BY sort_order, created_at
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	var slots []map[string]interface{}
	for rows.Next() {
		var id, slotType, title, linkURL string
		var categorySlug *string
		var isActive bool
		var sortOrder int
		rows.Scan(&id, &slotType, &title, &linkURL, &categorySlug, &isActive, &sortOrder)
		slots = append(slots, map[string]interface{}{
			"id": id, "slot_type": slotType, "title": title,
			"link_url": linkURL, "category_slug": categorySlug,
			"is_active": isActive, "sort_order": sortOrder,
		})
	}
	if slots == nil {
		slots = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"slots": slots})
}

// PUT /api/v1/super-admin/settings/ad-slots
func (h *AdminHandler) SaveAdSlots(c *gin.Context) {
	var req struct {
		Slots []struct {
			ID           *string `json:"id"`
			SlotType     string  `json:"slot_type"`
			Title        string  `json:"title"`
			LinkURL      string  `json:"link_url"`
			CategorySlug *string `json:"category_slug"`
			IsActive     bool    `json:"is_active"`
			SortOrder    int     `json:"sort_order"`
		} `json:"slots"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, s := range req.Slots {
		if s.ID == nil {
			h.db.Exec(c.Request.Context(), `
				INSERT INTO ad_slots (slot_type, title, link_url, category_slug, is_active, sort_order)
				VALUES ($1,$2,$3,$4,$5,$6)
			`, s.SlotType, s.Title, s.LinkURL, s.CategorySlug, s.IsActive, i+1)
		} else {
			h.db.Exec(c.Request.Context(), `
				UPDATE ad_slots
				SET slot_type=$1, title=$2, link_url=$3, category_slug=$4, is_active=$5, sort_order=$6
				WHERE id=$7
			`, s.SlotType, s.Title, s.LinkURL, s.CategorySlug, s.IsActive, i+1, *s.ID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "광고 슬롯이 저장되었습니다"})
}

// DELETE /api/v1/super-admin/settings/ad-slots/:id
func (h *AdminHandler) DeleteAdSlot(c *gin.Context) {
	id := c.Param("id")
	h.db.Exec(c.Request.Context(), `DELETE FROM ad_slots WHERE id=$1`, id)
	c.JSON(http.StatusOK, gin.H{"message": "삭제되었습니다"})
}

// ─── 제휴 설정 ───

// GET /api/v1/super-admin/settings/affiliate
func (h *AdminHandler) GetAffiliateSettings(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `SELECT provider, tracking_id FROM affiliate_settings ORDER BY provider`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	settings := map[string]string{}
	for rows.Next() {
		var provider string
		var trackingID *string
		rows.Scan(&provider, &trackingID)
		if trackingID != nil {
			settings[provider] = *trackingID
		} else {
			settings[provider] = ""
		}
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// PUT /api/v1/super-admin/settings/affiliate
func (h *AdminHandler) UpdateAffiliateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for provider, trackingID := range req {
		h.db.Exec(c.Request.Context(), `
			UPDATE affiliate_settings SET tracking_id=$1, updated_at=NOW() WHERE provider=$2
		`, trackingID, provider)
	}
	c.JSON(http.StatusOK, gin.H{"message": "제휴 설정이 저장되었습니다"})
}
