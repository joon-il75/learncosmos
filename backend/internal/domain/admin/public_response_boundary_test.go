package admin

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPublicLumiRuntimeConfigResponseOmitsEditorMetadata(t *testing.T) {
	publicResponse := buildPublicLumiRuntimeConfigResponse(lumiRuntimeConfigResponse{
		Exists:    true,
		UpdatedAt: "2026-06-16T00:00:00Z",
		Version:   123,
		Rules: []lumiRuntimeRuleDraftPayload{{
			ID:          "rule-1",
			Label:       "Rule",
			Trigger:     "dashboard_loaded",
			Mode:        "dashboard",
			Page:        "dashboard",
			Action:      "none",
			Scene:       "home",
			Context:     "default",
			DockSlot:    "dashboard-floating",
			State:       "idle",
			MessageType: "info",
			Message:     "hello",
			Note:        "operator-only note",
		}},
		Messages: []lumiRuntimeMessageDraftPayload{{
			ID:          "message-1",
			Trigger:     "empty",
			Label:       "Message",
			Note:        "draft note",
			MessageType: "info",
			Message:     "hello",
			Preview:     "hi",
		}},
	})

	payload, err := json.Marshal(publicResponse)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	body := string(payload)
	for _, forbidden := range []string{"note", "exists", "updated_at", "version", "operator-only", "draft note"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("public Lumi runtime payload leaked %q: %s", forbidden, body)
		}
	}
	for _, required := range []string{"rules", "messages", "dockSlot", "messageType", "preview"} {
		if !strings.Contains(body, required) {
			t.Fatalf("public Lumi runtime payload missing %q: %s", required, body)
		}
	}
}

func TestPublicPlanetTextureMapResponseOmitsFileMetadata(t *testing.T) {
	publicResponse := buildPublicPlanetTextureMapAssetResponse(planetTextureMapAssetResponse{
		ID:                      "map-1",
		Name:                    "Map",
		Description:             "Active map",
		AssetPath:               "/textures/planets/maps/map.webp",
		PublicURL:               "/textures/planets/maps/map.webp",
		Width:                   2048,
		Height:                  768,
		Columns:                 4,
		Rows:                    3,
		CellWidth:               512,
		CellHeight:              256,
		RotationDurationSeconds: 36,
		RotationDirection:       "left",
		IsActive:                true,
		IsBuiltin:               true,
		Exists:                  true,
		SizeBytes:               123456,
		UpdatedAt:               "2026-06-16T00:00:00Z",
		Version:                 42,
		CreatedAt:               "2026-06-15T00:00:00Z",
	})

	payload, err := json.Marshal(publicResponse)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	body := string(payload)
	for _, forbidden := range []string{"is_active", "is_builtin", "exists", "size_bytes", "updated_at", "created_at"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("public planet texture payload leaked %q: %s", forbidden, body)
		}
	}
	for _, required := range []string{"public_url", "rotation_duration_seconds", "rotation_direction", "version"} {
		if !strings.Contains(body, required) {
			t.Fatalf("public planet texture payload missing %q: %s", required, body)
		}
	}
}
