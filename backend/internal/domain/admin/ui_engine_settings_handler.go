package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const dashboardUIEngineSettingsKey = "dashboard"

type dashboardUIEngineSettings struct {
	CompactBreakpoint         int  `json:"compact_breakpoint"`
	PhoneBreakpoint           int  `json:"phone_breakpoint"`
	ShortViewportHeight       int  `json:"short_viewport_height"`
	CTAMinLaunchDurationMs    int  `json:"cta_min_launch_duration_ms"`
	SelectionCTALaunchDelayMs int  `json:"selection_cta_launch_delay_ms"`
	LumiQuickActionLimit      int  `json:"lumi_quick_action_limit"`
	LumiDesktopPanelEnabled   bool `json:"lumi_desktop_panel_enabled"`
	LumiMobileSheetEnabled    bool `json:"lumi_mobile_sheet_enabled"`
}

type dashboardUIEngineSettingsResponse struct {
	Settings  dashboardUIEngineSettings `json:"settings"`
	UpdatedAt string                    `json:"updated_at,omitempty"`
}

func defaultDashboardUIEngineSettings() dashboardUIEngineSettings {
	return dashboardUIEngineSettings{
		CompactBreakpoint:         1180,
		PhoneBreakpoint:           520,
		ShortViewportHeight:       460,
		CTAMinLaunchDurationMs:    920,
		SelectionCTALaunchDelayMs: 720,
		LumiQuickActionLimit:      2,
		LumiDesktopPanelEnabled:   true,
		LumiMobileSheetEnabled:    true,
	}
}

func normalizeDashboardUIEngineSettings(input dashboardUIEngineSettings) dashboardUIEngineSettings {
	settings := defaultDashboardUIEngineSettings()

	if input.CompactBreakpoint >= 720 && input.CompactBreakpoint <= 1600 {
		settings.CompactBreakpoint = input.CompactBreakpoint
	}
	if input.PhoneBreakpoint >= 320 && input.PhoneBreakpoint <= 960 {
		settings.PhoneBreakpoint = input.PhoneBreakpoint
	}
	if settings.PhoneBreakpoint > settings.CompactBreakpoint {
		settings.PhoneBreakpoint = settings.CompactBreakpoint
	}
	if input.ShortViewportHeight >= 280 && input.ShortViewportHeight <= 900 {
		settings.ShortViewportHeight = input.ShortViewportHeight
	}
	if input.CTAMinLaunchDurationMs >= 0 && input.CTAMinLaunchDurationMs <= 5000 {
		settings.CTAMinLaunchDurationMs = input.CTAMinLaunchDurationMs
	}
	if input.SelectionCTALaunchDelayMs >= 0 && input.SelectionCTALaunchDelayMs <= 5000 {
		settings.SelectionCTALaunchDelayMs = input.SelectionCTALaunchDelayMs
	}
	if input.LumiQuickActionLimit >= 0 && input.LumiQuickActionLimit <= 4 {
		settings.LumiQuickActionLimit = input.LumiQuickActionLimit
	}

	settings.LumiDesktopPanelEnabled = input.LumiDesktopPanelEnabled
	settings.LumiMobileSheetEnabled = input.LumiMobileSheetEnabled

	return settings
}

func (h *AdminHandler) getPublicDashboardUIEngineSettingsHandler(c *gin.Context) {
	response, err := h.readDashboardUIEngineSettingsResponse(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read ui engine settings"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminHandler) getDashboardUIEngineSettingsHandler(c *gin.Context) {
	response, err := h.readDashboardUIEngineSettingsResponse(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read ui engine settings"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminHandler) updateDashboardUIEngineSettingsHandler(c *gin.Context) {
	var req struct {
		Settings dashboardUIEngineSettings `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	normalized := normalizeDashboardUIEngineSettings(req.Settings)
	payload, err := json.Marshal(normalized)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode ui engine settings"})
		return
	}

	if _, err := h.db.Exec(c.Request.Context(), `
		INSERT INTO ui_engine_settings (key, value, updated_at)
		VALUES ($1, $2::jsonb, NOW())
		ON CONFLICT (key)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, dashboardUIEngineSettingsKey, string(payload)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist ui engine settings"})
		return
	}

	response, err := h.readDashboardUIEngineSettingsResponse(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload ui engine settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "ui engine settings saved",
		"settings":   response.Settings,
		"updated_at": response.UpdatedAt,
	})
}

func (h *AdminHandler) readDashboardUIEngineSettingsResponse(c *gin.Context) (dashboardUIEngineSettingsResponse, error) {
	settings := defaultDashboardUIEngineSettings()

	var raw []byte
	var updatedAt *time.Time
	err := h.db.QueryRow(c.Request.Context(), `
		SELECT value, updated_at
		FROM ui_engine_settings
		WHERE key = $1
	`, dashboardUIEngineSettingsKey).Scan(&raw, &updatedAt)
	if err != nil {
		if err != pgx.ErrNoRows {
			return dashboardUIEngineSettingsResponse{}, err
		}
		return dashboardUIEngineSettingsResponse{
			Settings: settings,
		}, nil
	}

	var stored dashboardUIEngineSettings
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &stored); err != nil {
			return dashboardUIEngineSettingsResponse{}, err
		}
		settings = normalizeDashboardUIEngineSettings(stored)
	}

	response := dashboardUIEngineSettingsResponse{
		Settings: settings,
	}
	if updatedAt != nil {
		response.UpdatedAt = updatedAt.Format(time.RFC3339)
	}

	return response, nil
}
