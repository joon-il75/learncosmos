package curriculum

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const tiptapInlineImageMaxUploadBytes = int64(5 * 1024 * 1024)

var tiptapInlineImageAllowedExtensions = map[string]struct{}{
	".gif":  {},
	".jpeg": {},
	".jpg":  {},
	".png":  {},
	".webp": {},
}

var tiptapInlineImageAllowedContentTypes = map[string]struct{}{
	"image/gif":  {},
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

var researchMaterialAllowedVideoContentTypes = map[string]struct{}{
	"video/mp4":       {},
	"video/quicktime": {},
	"video/webm":      {},
}

var researchMaterialAllowedSubtitleContentTypes = map[string]struct{}{
	"application/octet-stream": {},
	"application/x-subrip":     {},
	"text/plain":               {},
	"text/srt":                 {},
	"text/vtt":                 {},
}

func validateTiptapInlineImageUploadMetadata(fileHeader *multipart.FileHeader) error {
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if _, ok := tiptapInlineImageAllowedExtensions[extension]; !ok {
		return errors.New("only jpg, jpeg, png, webp, and gif images can be uploaded")
	}
	if fileHeader.Size > tiptapInlineImageMaxUploadBytes {
		return fmt.Errorf("image upload limit is %dMB", tiptapInlineImageMaxUploadBytes/1024/1024)
	}
	return nil
}

func validateResearchMaterialUploadMetadata(fileHeader *multipart.FileHeader, rawAttachmentType string) error {
	plannedAttachmentType := inferResearchMaterialAttachmentType(rawAttachmentType, fileHeader.Filename)
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	_, imageOK := researchMaterialAllowedImageExtensions[extension]
	_, docOK := researchMaterialAllowedDocumentExtensions[extension]
	_, videoOK := researchMaterialAllowedVideoExtensions[extension]
	_, subtitleOK := researchMaterialAllowedSubtitleExtensions[extension]
	if plannedAttachmentType == "thumbnail" {
		docOK = false
		videoOK = false
		subtitleOK = false
	}
	if plannedAttachmentType == "subtitle" {
		imageOK = false
		docOK = false
		videoOK = false
	}
	if !imageOK && !docOK && !videoOK && !subtitleOK {
		return errors.New("research material allows only document, zip, image, video, or subtitle files")
	}
	limit := researchMaterialDocumentMaxBytes
	if imageOK {
		limit = researchMaterialImageMaxBytes
	} else if videoOK {
		limit = researchMaterialVideoMaxBytes
	} else if subtitleOK {
		limit = researchMaterialSubtitleMaxBytes
	}
	if fileHeader.Size > limit {
		if imageOK {
			return fmt.Errorf("image upload limit is %dMB", researchMaterialImageMaxBytes/1024/1024)
		}
		if videoOK {
			return fmt.Errorf("video upload limit is %dMB", researchMaterialVideoMaxBytes/1024/1024)
		}
		if subtitleOK {
			return fmt.Errorf("subtitle upload limit is %dMB", researchMaterialSubtitleMaxBytes/1024/1024)
		}
		return fmt.Errorf("document upload limit is %dMB", researchMaterialDocumentMaxBytes/1024/1024)
	}
	return nil
}

func validateTiptapInlineImageUploadContent(contentType string) error {
	normalized := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if _, ok := tiptapInlineImageAllowedContentTypes[normalized]; !ok {
		return errors.New("only jpg, jpeg, png, webp, and gif images can be uploaded")
	}
	return nil
}

func validateResearchMaterialUploadContent(filename, rawAttachmentType, contentType string) error {
	plannedAttachmentType := inferResearchMaterialAttachmentType(rawAttachmentType, filename)
	extension := strings.ToLower(filepath.Ext(filename))
	normalized := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if plannedAttachmentType == "video" {
		if _, ok := researchMaterialAllowedVideoExtensions[extension]; !ok {
			return errors.New("only mp4, webm, and mov videos can be uploaded")
		}
		if _, ok := researchMaterialAllowedVideoContentTypes[normalized]; !ok {
			return errors.New("only mp4, webm, and mov videos can be uploaded")
		}
		return nil
	}
	if plannedAttachmentType == "thumbnail" {
		if _, ok := researchMaterialAllowedImageExtensions[extension]; !ok {
			return errors.New("only jpg, jpeg, png, webp, and gif thumbnails can be uploaded")
		}
		if _, ok := tiptapInlineImageAllowedContentTypes[normalized]; !ok {
			return errors.New("only jpg, jpeg, png, webp, and gif thumbnails can be uploaded")
		}
		return nil
	}
	if plannedAttachmentType == "image" {
		if _, ok := researchMaterialAllowedImageExtensions[extension]; !ok {
			return errors.New("only jpg, jpeg, png, webp, and gif images can be uploaded")
		}
		if _, ok := tiptapInlineImageAllowedContentTypes[normalized]; !ok {
			return errors.New("only jpg, jpeg, png, webp, and gif images can be uploaded")
		}
		return nil
	}
	if plannedAttachmentType == "subtitle" {
		if _, ok := researchMaterialAllowedSubtitleExtensions[extension]; !ok {
			return errors.New("only vtt and srt subtitles can be uploaded")
		}
		if _, ok := researchMaterialAllowedSubtitleContentTypes[normalized]; !ok {
			return errors.New("only vtt and srt subtitles can be uploaded")
		}
	}
	return nil
}

func researchMaterialUploadLimitBytes(rawAttachmentType, filename string) int64 {
	plannedAttachmentType := inferResearchMaterialAttachmentType(rawAttachmentType, filename)
	if plannedAttachmentType == "subtitle" {
		return researchMaterialSubtitleMaxBytes
	}
	if plannedAttachmentType == "thumbnail" {
		return researchMaterialImageMaxBytes
	}
	extension := strings.ToLower(filepath.Ext(filename))
	if _, ok := researchMaterialAllowedImageExtensions[extension]; ok {
		return researchMaterialImageMaxBytes
	}
	if _, ok := researchMaterialAllowedVideoExtensions[extension]; ok {
		return researchMaterialVideoMaxBytes
	}
	return researchMaterialDocumentMaxBytes
}

func detectUploadedAttachmentContentType(fileHeader *multipart.FileHeader) (string, error) {
	contentType, err := sniffUploadedAttachmentContentType(fileHeader)
	if err != nil {
		return "", err
	}
	if contentType != "application/octet-stream" {
		return contentType, nil
	}
	if headerContentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type")); headerContentType != "" {
		return headerContentType, nil
	}
	return contentType, nil
}

func sniffUploadedAttachmentContentType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	return http.DetectContentType(buffer[:n]), nil
}

func readUploadedAttachment(fileHeader *multipart.FileHeader, maxBytes int64) ([]byte, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	limited := io.LimitReader(file, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxBytes {
		return nil, "", http.ErrBodyReadAfterClose
	}

	contentType := http.DetectContentType(data)
	if contentType == "application/octet-stream" {
		if headerContentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type")); headerContentType != "" {
			contentType = headerContentType
		}
	}
	return data, contentType, nil
}

func normalizeUploadedAttachmentType(value string, contentType string) string {
	value = strings.TrimSpace(value)
	switch value {
	case "image", "file", "link", "code", "other", "video", "subtitle", "thumbnail":
		return value
	}
	normalizedContentType := strings.ToLower(contentType)
	if strings.HasPrefix(normalizedContentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(normalizedContentType, "video/") {
		return "video"
	}
	return "file"
}

func inferResearchMaterialAttachmentType(rawAttachmentType, filename string) string {
	rawAttachmentType = strings.TrimSpace(rawAttachmentType)
	if rawAttachmentType != "" {
		return rawAttachmentType
	}
	extension := strings.ToLower(filepath.Ext(filename))
	if _, ok := researchMaterialAllowedImageExtensions[extension]; ok {
		return "image"
	}
	if _, ok := researchMaterialAllowedVideoExtensions[extension]; ok {
		return "video"
	}
	if _, ok := researchMaterialAllowedSubtitleExtensions[extension]; ok {
		return "subtitle"
	}
	return "file"
}

func stringPtr(value string) *string {
	return &value
}
