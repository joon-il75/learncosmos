package admin

import "testing"

func TestRuntimePathsUseDefaultProjectRoot(t *testing.T) {
	t.Setenv("LEARNWEAVER_PROJECT_ROOT", "")
	t.Setenv("PROJECT_ROOT", "")
	t.Setenv("PLANET_TEXTURE_MAP_DIR", "")

	got := planetTextureMapsSourceDirectory()
	want := "/home/weaver/learnweaver/frontend/public/textures/planets/maps"
	if got != want {
		t.Fatalf("planetTextureMapsSourceDirectory() = %q, want %q", got, want)
	}
}

func TestRuntimePathsUseProjectRootOverride(t *testing.T) {
	t.Setenv("LEARNWEAVER_PROJECT_ROOT", "")
	t.Setenv("PROJECT_ROOT", "/home/cosmos/LearnCosmos")
	t.Setenv("LUMI_RUNTIME_CONFIG_PATH", "")
	t.Setenv("LUMI_ASSET_SOURCE_PATH", "")

	if got, want := lumiRuntimeConfigFilePath(), "/home/cosmos/LearnCosmos/backend/storage/lumi-runtime-config.json"; got != want {
		t.Fatalf("lumiRuntimeConfigFilePath() = %q, want %q", got, want)
	}
	if got, want := lumiAssetSourceFilePath(), "/home/cosmos/LearnCosmos/frontend/public/images/lumi.webp"; got != want {
		t.Fatalf("lumiAssetSourceFilePath() = %q, want %q", got, want)
	}
}

func TestRuntimePathsUseExplicitDirectoryOverride(t *testing.T) {
	t.Setenv("PROJECT_ROOT", "/home/cosmos/LearnCosmos")
	t.Setenv("PLANET_TEXTURE_MAP_DIR", "/srv/learncosmos/planet-maps")

	got := planetTextureMapsSourceDirectory()
	want := "/srv/learncosmos/planet-maps"
	if got != want {
		t.Fatalf("planetTextureMapsSourceDirectory() = %q, want %q", got, want)
	}
}
