package admin

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	lumiAssetMaxBytes  = 10 << 20
	lumiAssetMaxWidth  = 8192
	lumiAssetMaxHeight = 8192
	lumiAssetMaxPixels = 24 * 1024 * 1024
)

type lumiAssetMetaResponse struct {
	Path      string `json:"path"`
	PublicURL string `json:"public_url"`
	Exists    bool   `json:"exists"`
	SizeBytes int64  `json:"size_bytes"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Version   int64  `json:"version,omitempty"`
}

func (h *AdminHandler) getLumiAssetMetaHandler(c *gin.Context) {
	info, err := os.Stat(lumiAssetSourceFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, lumiAssetMetaResponse{
				Path:      lumiAssetSourceFilePath(),
				PublicURL: "/images/lumi.webp",
				Exists:    false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lumi asset metadata read failed"})
		return
	}

	c.JSON(http.StatusOK, lumiAssetMetaResponse{
		Path:      lumiAssetSourceFilePath(),
		PublicURL: "/images/lumi.webp",
		Exists:    true,
		SizeBytes: info.Size(),
		UpdatedAt: info.ModTime().Format(time.RFC3339),
		Version:   info.ModTime().Unix(),
	})
}

func (h *AdminHandler) uploadLumiAssetHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, lumiAssetMaxBytes+1024)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if fileHeader.Size > lumiAssetMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}

	if ext := strings.ToLower(filepath.Ext(fileHeader.Filename)); ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .webp sprite sheets are supported"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(io.LimitReader(file, lumiAssetMaxBytes+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}
	if int64(len(payload)) > lumiAssetMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large"})
		return
	}

	if err := validateLumiAssetUploadPayload(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targets := []string{lumiAssetSourceFilePath(), lumiAssetStandaloneFilePath()}
	for _, target := range targets {
		if err := writeLumiAssetFile(target, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write lumi asset"})
			return
		}
	}

	info, err := os.Stat(lumiAssetSourceFilePath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lumi asset metadata read failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "lumi asset uploaded",
		"asset": lumiAssetMetaResponse{
			Path:      lumiAssetSourceFilePath(),
			PublicURL: "/images/lumi.webp",
			Exists:    true,
			SizeBytes: info.Size(),
			UpdatedAt: info.ModTime().Format(time.RFC3339),
			Version:   info.ModTime().Unix(),
		},
	})
}

func validateLumiAssetUploadPayload(payload []byte) error {
	width, height, err := readWebPDimensions(payload)
	if err != nil {
		return errors.New("invalid webp sprite sheet")
	}
	if width <= 0 || height <= 0 || width > lumiAssetMaxWidth || height > lumiAssetMaxHeight || width*height > lumiAssetMaxPixels {
		return errors.New("webp sprite sheet dimensions are too large")
	}
	return nil
}

func writeLumiAssetFile(target string, payload []byte) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	tmpPath := target + ".uploading"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, target)
}
