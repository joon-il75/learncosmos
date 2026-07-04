package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── LLM/API key settings ───

func (h *AdminHandler) GetLLMSettings(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT feature, provider, model, endpoint_url
		FROM system_llm_settings
		ORDER BY feature
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	var settings []map[string]interface{}
	for rows.Next() {
		var feature, provider, model string
		var endpointURL *string
		if err := rows.Scan(&feature, &provider, &model, &endpointURL); err != nil {
			continue
		}
		settings = append(settings, map[string]interface{}{
			"feature":      feature,
			"provider":     provider,
			"model":        model,
			"endpoint_url": endpointURL,
		})
	}
	if settings == nil {
		settings = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// PUT /api/v1/super-admin/settings/llm
func (h *AdminHandler) UpdateLLMSettings(c *gin.Context) {
	var req struct {
		Settings []struct {
			Feature     string  `json:"feature"`
			Provider    string  `json:"provider"`
			Model       string  `json:"model"`
			EndpointURL *string `json:"endpoint_url"`
		} `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, s := range req.Settings {
		h.db.Exec(c.Request.Context(), `
			INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (feature)
			DO UPDATE SET provider=EXCLUDED.provider, model=EXCLUDED.model, endpoint_url=EXCLUDED.endpoint_url, updated_at=NOW()
		`, s.Feature, s.Provider, s.Model, s.EndpointURL)
	}
	c.JSON(http.StatusOK, gin.H{"message": "저장되었습니다"})
}

// ─── API 키 ───

// GET /api/v1/super-admin/settings/api-keys
func (h *AdminHandler) GetAPIKeys(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT provider,
		       CASE WHEN api_key_encrypted IS NOT NULL AND api_key_encrypted != '' THEN 'SET' ELSE '' END,
		       CASE WHEN api_key_secondary_encrypted IS NOT NULL AND api_key_secondary_encrypted != '' THEN 'SET' ELSE '' END,
		       endpoint_url,
		       updated_at
		FROM system_api_keys ORDER BY provider
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "설정 조회 실패"})
		return
	}
	defer rows.Close()

	var keys []map[string]interface{}
	for rows.Next() {
		var provider, keyStatus, key2Status string
		var endpointURL *string
		var updatedAt *time.Time
		rows.Scan(&provider, &keyStatus, &key2Status, &endpointURL, &updatedAt)
		var updatedAtStr *string
		if updatedAt != nil {
			formatted := updatedAt.Format("2006-01-02 15:04:05")
			updatedAtStr = &formatted
		}
		keys = append(keys, map[string]interface{}{
			"provider":     provider,
			"key_status":   keyStatus,
			"key2_status":  key2Status,
			"endpoint_url": endpointURL,
			"has_key":      keyStatus == "SET",
			"has_key2":     key2Status == "SET",
			"updated_at":   updatedAtStr,
		})
	}
	if keys == nil {
		keys = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// PUT /api/v1/super-admin/settings/api-keys
func (h *AdminHandler) UpdateAPIKey(c *gin.Context) {
	var req struct {
		Provider    string  `json:"provider"`
		APIKey      string  `json:"api_key"`
		APIKey2     string  `json:"api_key2"`
		EndpointURL *string `json:"endpoint_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Provider = strings.TrimSpace(req.Provider)
	if req.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	if err := h.settingsStore.UpsertProviderCredentials(
		c.Request.Context(),
		req.Provider,
		req.APIKey,
		req.APIKey2,
		req.EndpointURL,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API 키 저장에 실패했습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "API 키가 저장되었습니다"})
}
