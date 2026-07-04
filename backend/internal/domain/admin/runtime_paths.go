package admin

import (
	"os"
	"path/filepath"
	"strings"
)

const defaultProjectRoot = "/home/weaver/learnweaver"

func projectRootPath() string {
	for _, key := range []string{"LEARNWEAVER_PROJECT_ROOT", "PROJECT_ROOT"} {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return filepath.Clean(value)
		}
	}
	return defaultProjectRoot
}

func backendStoragePath(name string) string {
	return filepath.Join(projectRootPath(), "backend", "storage", name)
}

func frontendPublicPath(parts ...string) string {
	items := append([]string{projectRootPath(), "frontend", "public"}, parts...)
	return filepath.Join(items...)
}

func frontendStandalonePublicPath(parts ...string) string {
	items := append([]string{projectRootPath(), "frontend", ".next", "standalone", "public"}, parts...)
	return filepath.Join(items...)
}

func lumiRuntimeConfigFilePath() string {
	if value := strings.TrimSpace(os.Getenv("LUMI_RUNTIME_CONFIG_PATH")); value != "" {
		return filepath.Clean(value)
	}
	return backendStoragePath("lumi-runtime-config.json")
}

func lumiAssetSourceFilePath() string {
	if value := strings.TrimSpace(os.Getenv("LUMI_ASSET_SOURCE_PATH")); value != "" {
		return filepath.Clean(value)
	}
	return frontendPublicPath("images", "lumi.webp")
}

func lumiAssetStandaloneFilePath() string {
	if value := strings.TrimSpace(os.Getenv("LUMI_ASSET_STANDALONE_PATH")); value != "" {
		return filepath.Clean(value)
	}
	return frontendStandalonePublicPath("images", "lumi.webp")
}

func planetTextureMapsSourceDirectory() string {
	if value := strings.TrimSpace(os.Getenv("PLANET_TEXTURE_MAP_DIR")); value != "" {
		return filepath.Clean(value)
	}
	return frontendPublicPath("textures", "planets", "maps")
}

func planetTextureMapsStandaloneDirectory() string {
	if value := strings.TrimSpace(os.Getenv("PLANET_TEXTURE_MAP_STANDALONE_DIR")); value != "" {
		return filepath.Clean(value)
	}
	return frontendStandalonePublicPath("textures", "planets", "maps")
}
