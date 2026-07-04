package admin

import (
	"encoding/binary"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	planetTextureMapMaxBytes  = 15 << 20
	planetTextureMapMaxWidth  = 4096
	planetTextureMapMaxHeight = 3072
	planetTextureMapMaxPixels = planetTextureMapMaxWidth * planetTextureMapMaxHeight
	planetTextureMapColumns   = 4
	planetTextureMapRows      = 3
	defaultPlanetTextureMapID = "8f7a8f12-4c7f-4e0f-9f6c-2e9b1f5d0c31"
)

type publicPlanetTextureMapAssetResponse struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	AssetPath               string `json:"asset_path"`
	PublicURL               string `json:"public_url"`
	Width                   int    `json:"width"`
	Height                  int    `json:"height"`
	Columns                 int    `json:"columns"`
	Rows                    int    `json:"rows"`
	CellWidth               int    `json:"cell_width"`
	CellHeight              int    `json:"cell_height"`
	RotationDurationSeconds int    `json:"rotation_duration_seconds"`
	RotationDirection       string `json:"rotation_direction"`
	Version                 int64  `json:"version,omitempty"`
}

type planetTextureMapAssetResponse struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	AssetPath               string `json:"asset_path"`
	PublicURL               string `json:"public_url"`
	Width                   int    `json:"width"`
	Height                  int    `json:"height"`
	Columns                 int    `json:"columns"`
	Rows                    int    `json:"rows"`
	CellWidth               int    `json:"cell_width"`
	CellHeight              int    `json:"cell_height"`
	RotationDurationSeconds int    `json:"rotation_duration_seconds"`
	RotationDirection       string `json:"rotation_direction"`
	IsActive                bool   `json:"is_active"`
	IsBuiltin               bool   `json:"is_builtin"`
	Exists                  bool   `json:"exists"`
	SizeBytes               int64  `json:"size_bytes"`
	UpdatedAt               string `json:"updated_at,omitempty"`
	Version                 int64  `json:"version,omitempty"`
	CreatedAt               string `json:"created_at,omitempty"`
}

type planetTextureMapRecord struct {
	ID                      uuid.UUID
	Name                    string
	Description             string
	AssetPath               string
	Width                   int
	Height                  int
	Columns                 int
	Rows                    int
	CellWidth               int
	CellHeight              int
	RotationDurationSeconds int
	RotationDirection       string
	IsActive                bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type updatePlanetTextureMapRequest struct {
	Name                    *string `json:"name"`
	Description             *string `json:"description"`
	IsActive                *bool   `json:"is_active"`
	RotationDurationSeconds *int    `json:"rotation_duration_seconds"`
	RotationDirection       *string `json:"rotation_direction"`
}

func buildPublicPlanetTextureMapAssetResponse(resp planetTextureMapAssetResponse) publicPlanetTextureMapAssetResponse {
	return publicPlanetTextureMapAssetResponse{
		ID:                      resp.ID,
		Name:                    resp.Name,
		Description:             resp.Description,
		AssetPath:               resp.AssetPath,
		PublicURL:               resp.PublicURL,
		Width:                   resp.Width,
		Height:                  resp.Height,
		Columns:                 resp.Columns,
		Rows:                    resp.Rows,
		CellWidth:               resp.CellWidth,
		CellHeight:              resp.CellHeight,
		RotationDurationSeconds: resp.RotationDurationSeconds,
		RotationDirection:       resp.RotationDirection,
		Version:                 resp.Version,
	}
}

func validatePlanetTextureMapDimensions(width, height int) error {
	if width <= 0 || height <= 0 {
		return errors.New("invalid texture map dimensions")
	}
	if width > planetTextureMapMaxWidth || height > planetTextureMapMaxHeight || width*height > planetTextureMapMaxPixels {
		return errors.New("texture map dimensions are too large")
	}
	if width%planetTextureMapColumns != 0 || height%planetTextureMapRows != 0 {
		return errors.New("texture map size must be divisible by 4 columns and 3 rows")
	}
	return nil
}

func buildPlanetTextureMapAssetResponse(record planetTextureMapRecord) (planetTextureMapAssetResponse, error) {
	targetPath := filepath.Join(planetTextureMapsSourceDirectory(), filepath.Base(record.AssetPath))
	info, err := os.Stat(targetPath)
	exists := err == nil
	if err != nil && !os.IsNotExist(err) {
		return planetTextureMapAssetResponse{}, err
	}

	resp := planetTextureMapAssetResponse{
		ID:                      record.ID.String(),
		Name:                    record.Name,
		Description:             record.Description,
		AssetPath:               record.AssetPath,
		PublicURL:               record.AssetPath,
		Width:                   record.Width,
		Height:                  record.Height,
		Columns:                 record.Columns,
		Rows:                    record.Rows,
		CellWidth:               record.CellWidth,
		CellHeight:              record.CellHeight,
		RotationDurationSeconds: record.RotationDurationSeconds,
		RotationDirection:       record.RotationDirection,
		IsActive:                record.IsActive,
		IsBuiltin:               record.ID.String() == defaultPlanetTextureMapID,
		Exists:                  exists,
		CreatedAt:               record.CreatedAt.Format(time.RFC3339),
	}
	if exists {
		resp.SizeBytes = info.Size()
		resp.UpdatedAt = info.ModTime().Format(time.RFC3339)
		resp.Version = info.ModTime().Unix()
		return resp, nil
	}
	resp.UpdatedAt = record.UpdatedAt.Format(time.RFC3339)
	resp.Version = record.UpdatedAt.Unix()
	return resp, nil
}

func parsePlanetTextureRotationDuration(raw string) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 36
	}
	return clampPlanetTextureRotationDuration(value)
}

func clampPlanetTextureRotationDuration(value int) int {
	if value < 18 {
		return 18
	}
	if value > 90 {
		return 90
	}
	return value
}

func normalizePlanetTextureRotationDirection(raw string) string {
	if strings.EqualFold(strings.TrimSpace(raw), "right") {
		return "right"
	}
	return "left"
}

func parsePlanetTypeActive(raw string) bool {
	return !strings.EqualFold(strings.TrimSpace(raw), "false")
}

func nullableString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func writePlanetTypeAssetFile(path string, payload []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}

func isSupportedPlanetTextureUploadExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	default:
		return false
	}
}

func convertPlanetTextureUploadToWebP(c *gin.Context, payload []byte, ext string) ([]byte, error) {
	if ext == ".webp" {
		return payload, nil
	}

	tmpDir, err := os.MkdirTemp("", "learnweaver-planet-texture-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input"+ext)
	outputPath := filepath.Join(tmpDir, "output.webp")
	if err := os.WriteFile(inputPath, payload, 0o600); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(
		c.Request.Context(),
		"ffmpeg",
		"-y",
		"-hide_banner",
		"-loglevel", "error",
		"-i", inputPath,
		"-frames:v", "1",
		"-c:v", "libwebp",
		"-lossless", "1",
		"-compression_level", "6",
		outputPath,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, errors.New(strings.TrimSpace(string(output)))
	}

	return os.ReadFile(outputPath)
}

func readWebPDimensions(payload []byte) (int, int, error) {
	if len(payload) < 20 || string(payload[0:4]) != "RIFF" || string(payload[8:12]) != "WEBP" {
		return 0, 0, errors.New("invalid webp header")
	}

	for offset := 12; offset+8 <= len(payload); {
		chunkType := string(payload[offset : offset+4])
		chunkSize := int(binary.LittleEndian.Uint32(payload[offset+4 : offset+8]))
		chunkStart := offset + 8
		chunkEnd := chunkStart + chunkSize
		if chunkSize < 0 || chunkEnd > len(payload) {
			return 0, 0, errors.New("invalid webp chunk")
		}

		switch chunkType {
		case "VP8X":
			if chunkSize < 10 {
				return 0, 0, errors.New("invalid vp8x chunk")
			}
			width := 1 + int(payload[chunkStart+4]) + int(payload[chunkStart+5])<<8 + int(payload[chunkStart+6])<<16
			height := 1 + int(payload[chunkStart+7]) + int(payload[chunkStart+8])<<8 + int(payload[chunkStart+9])<<16
			return width, height, nil
		case "VP8L":
			if chunkSize < 5 || payload[chunkStart] != 0x2f {
				return 0, 0, errors.New("invalid vp8l chunk")
			}
			bits := uint32(payload[chunkStart+1]) | uint32(payload[chunkStart+2])<<8 | uint32(payload[chunkStart+3])<<16 | uint32(payload[chunkStart+4])<<24
			width := int(bits&0x3fff) + 1
			height := int((bits>>14)&0x3fff) + 1
			return width, height, nil
		case "VP8 ":
			if chunkSize < 10 || payload[chunkStart+3] != 0x9d || payload[chunkStart+4] != 0x01 || payload[chunkStart+5] != 0x2a {
				return 0, 0, errors.New("invalid vp8 chunk")
			}
			width := int(binary.LittleEndian.Uint16(payload[chunkStart+6:chunkStart+8]) & 0x3fff)
			height := int(binary.LittleEndian.Uint16(payload[chunkStart+8:chunkStart+10]) & 0x3fff)
			return width, height, nil
		}

		offset = chunkEnd
		if chunkSize%2 == 1 {
			offset++
		}
	}

	return 0, 0, errors.New("webp dimensions not found")
}
