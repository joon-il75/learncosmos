package content

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	"github.com/learnweaver/backend/internal/pkg/metaparser"
)

type Handler struct {
	repo          *Repository
	youtubeAPIKey string
	safetyService *safety.Service
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo:          repo,
		youtubeAPIKey: os.Getenv("YOUTUBE_API_KEY"),
	}
}

func (h *Handler) SetSafetyService(service *safety.Service) {
	if h == nil {
		return
	}
	h.safetyService = service
}

// getUserID — JWT 미들웨어에서 세팅된 user_id (string) 추출 후 UUID 파싱
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	idStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func embeddingGemmaModel() string {
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		return "embedding-gemma"
	}
	return model
}

func embeddingGemmaDimension() int {
	raw := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_DIMENSION"))
	if raw == "" {
		return 768
	}
	dimension, err := strconv.Atoi(raw)
	if err != nil || dimension <= 0 {
		return 768
	}
	return dimension
}

func (h *Handler) enforceSafety(c *gin.Context, userID uuid.UUID, targetID *uuid.UUID, text string, locale string, metadata map[string]any) bool {
	if h == nil || h.safetyService == nil {
		return false
	}
	result, err := h.safetyService.Enforce(c.Request.Context(), safety.ModerateInput{
		UserID:     &userID,
		TargetType: "content_metadata",
		TargetID:   targetID,
		Text:       text,
		Locale:     locale,
		Route:      c.FullPath(),
		Metadata:   metadata,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to moderate input", "error_code": "safety_moderation_failed"})
		return true
	}
	if safety.IsBlocked(result) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "input blocked by safety policy",
			"error_code": safety.ErrorCodeBlocked,
			"safety":     result,
		})
		return true
	}
	return false
}

func createContentModerationText(req CreateContentRequest) (string, []string) {
	parts := []string{}
	fields := []string{}
	appendText := func(field string, value *string) {
		if value == nil {
			return
		}
		text := strings.TrimSpace(*value)
		if text == "" {
			return
		}
		fields = append(fields, field)
		parts = append(parts, text)
	}
	appendText("title", &req.Title)
	appendText("description", req.Description)
	appendText("author", req.Author)
	appendText("url", req.URL)
	return strings.Join(parts, "\n"), fields
}

func updateContentModerationText(req UpdateContentRequest) (string, []string) {
	parts := []string{}
	fields := []string{}
	appendText := func(field string, value *string) {
		if value == nil {
			return
		}
		text := strings.TrimSpace(*value)
		if text == "" {
			return
		}
		fields = append(fields, field)
		parts = append(parts, text)
	}
	appendText("title", req.Title)
	appendText("description", req.Description)
	appendText("author", req.Author)
	appendText("url", req.URL)
	return strings.Join(parts, "\n"), fields
}

// POST /api/v1/contents
func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// youtube/blog/article은 URL 필수
	if req.ContentType != ContentTypeInternal && (req.URL == nil || strings.TrimSpace(*req.URL) == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required for " + string(req.ContentType)})
		return
	}
	if err := normalizeCreateContentRequestURLs(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	if text, fields := createContentModerationText(req); text != "" {
		if h.enforceSafety(c, userID, nil, text, req.Language, map[string]any{
			"field":        "content_metadata",
			"fields":       fields,
			"content_type": req.ContentType,
		}) {
			return
		}
	}

	content, err := h.repo.UpsertExternalContent(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create content"})
		return
	}

	// 비동기 임베딩 (응답 반환 후 백그라운드 처리)
	go func(contentID uuid.UUID, title, description string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		provider := "embedding_gemma"
		model := embeddingGemmaModel()
		text := embedder.BuildText(title, description)
		client, err := embedder.NewProviderClient(provider, strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")))
		if err != nil {
			log.Printf("[embedder] client 생성 실패 content_id=%s err=%v", contentID, err)
			_ = h.repo.SaveEmbeddingFailure(ctx, contentID, provider, model, embeddingGemmaDimension(), text, err)
			return
		}
		vector, err := client.Embed(ctx, text)
		if err != nil {
			log.Printf("[embedder] 실패 content_id=%s err=%v", contentID, err)
			_ = h.repo.SaveEmbeddingFailure(ctx, contentID, provider, model, embeddingGemmaDimension(), text, err)
			return
		}
		if err := h.repo.SaveEmbeddingReady(ctx, contentID, provider, model, text, vector); err != nil {
			log.Printf("[embedder] DB 저장 실패 content_id=%s err=%v", contentID, err)
			return
		}
		log.Printf("[embedder] 완료 provider=%s model=%s content_id=%s dim=%d", provider, model, contentID, len(vector))
	}(content.ID, content.Title, func() string {
		if content.Description != nil {
			return *content.Description
		}
		return ""
	}())

	c.JSON(http.StatusCreated, content)
}

// GET /api/v1/contents
func (h *Handler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var q ListContentsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.repo.List(c.Request.Context(), userID, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list contents"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /api/v1/contents/:id
func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	content, err := h.repo.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get content"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// PATCH /api/v1/contents/:id
func (h *Handler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := normalizeUpdateContentRequestURLs(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	if text, fields := updateContentModerationText(req); text != "" {
		locale := ""
		if req.Language != nil {
			locale = *req.Language
		}
		if h.enforceSafety(c, userID, &id, text, locale, map[string]any{
			"field":      "content_metadata",
			"fields":     fields,
			"content_id": id.String(),
		}) {
			return
		}
	}

	content, err := h.repo.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update content"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// ParseURLRequest — POST /api/v1/contents/parse 요청 바디
type ParseURLRequest struct {
	URL string `json:"url" binding:"required"`
}

// Parse — URL 메타 파싱 (저장 안 함, 미리보기용)
// POST /api/v1/contents/parse
func (h *Handler) Parse(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ParseURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := metaparser.Parse(c.Request.Context(), req.URL, h.youtubeAPIKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"title":            "",
			"description":      "",
			"thumbnail_url":    "",
			"author":           "",
			"duration_seconds": 0,
			"content_type":     metaparser.DetectContentType(req.URL),
			"source_url":       req.URL,
			"parse_error":      logsafe.Error(err),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// DELETE /api/v1/contents/:id
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id, userID); err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete content"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
