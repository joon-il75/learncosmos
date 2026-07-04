package admin

import (
	"os"
	"strings"
	"testing"
)

func TestSuperAdminSettingsReadErrorsDoNotExposeRawInternalErrors(t *testing.T) {
	for _, file := range []string{"settings_handler.go", "settings_ads_affiliate.go", "settings_points.go"} {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		text := string(source)
		if strings.Contains(text, `StatusInternalServerError, gin.H{"error": err.Error()}`) {
			t.Fatalf("%s exposes raw internal errors from super-admin read APIs", file)
		}
	}
}

func TestRankerArtifactStatusDoesNotExposeInternalPaths(t *testing.T) {
	status := recommendationRankerArtifactStatus{
		ModelPath:    "/srv/learnweaver/models/ranker/model.txt",
		ManifestPath: "/srv/learnweaver/models/ranker/model.txt.manifest.json",
		CheckedPaths: []string{
			"/srv/learnweaver/models/ranker/model.txt",
			"models/ranker/fallback.lgb",
		},
	}
	sanitizeRecommendationRankerArtifactStatus(&status)
	if status.ModelPath != "model.txt" {
		t.Fatalf("ModelPath = %q, want model.txt", status.ModelPath)
	}
	if status.ManifestPath != "model.txt.manifest.json" {
		t.Fatalf("ManifestPath = %q, want model.txt.manifest.json", status.ManifestPath)
	}
	if len(status.CheckedPaths) != 2 || status.CheckedPaths[0] != "model.txt" || status.CheckedPaths[1] != "fallback.lgb" {
		t.Fatalf("CheckedPaths = %#v, want file labels", status.CheckedPaths)
	}
}

func TestRankerArtifactStatusUsesSafeReasons(t *testing.T) {
	source, err := os.ReadFile("recommendation_debug_ranker_helpers.go")
	if err != nil {
		t.Fatalf("read recommendation_debug_ranker_helpers.go: %v", err)
	}
	text := string(source)
	for _, forbidden := range []string{
		`result["reason"] = err.Error()`,
		`" + err.Error()`,
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("ranker artifact status must not expose raw internal errors via %s", forbidden)
		}
	}
	if !strings.Contains(text, "defer sanitizeRecommendationRankerArtifactStatus(&status)") {
		t.Fatal("ranker artifact status must sanitize internal paths before returning")
	}
}
