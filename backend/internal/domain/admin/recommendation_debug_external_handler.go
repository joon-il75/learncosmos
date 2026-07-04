package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

func (h *AdminHandler) SaveRecommendationDebugExternalCandidate(c *gin.Context) {
	scenario, adminID, ok := h.loadRecommendationDebugScenarioForGoal(c)
	if !ok {
		return
	}

	var req recommendationDebugExternalCandidateSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.LessonID = strings.TrimSpace(req.LessonID)
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.URL = strings.TrimSpace(req.URL)
	req.Source = strings.TrimSpace(req.Source)
	req.ExternalContentID = strings.TrimSpace(req.ExternalContentID)
	req.ThumbnailURL = strings.TrimSpace(req.ThumbnailURL)
	req.Author = strings.TrimSpace(req.Author)
	req.Language = strings.TrimSpace(req.Language)
	if req.Title == "" || req.URL == "" || req.LessonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_id, title and url are required"})
		return
	}
	if len([]rune(req.Title)) > 200 || len([]rune(req.Description)) > 2000 || len([]rune(req.URL)) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "external candidate input too long"})
		return
	}

	lessonSnapshot, err := recommendationDebugLessonsFromSnapshot(scenario.GeneratedLessonsSnapshot)
	if err != nil || !recommendationDebugLessonExists(lessonSnapshot, req.LessonID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson snapshot not found"})
		return
	}

	ownerUserID, err := h.resolveRecommendationDebugContentOwner(c.Request.Context(), adminID, scenario)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content owner user unavailable"})
		return
	}

	urlValue := req.URL
	language := req.Language
	if language == "" {
		language = "ko"
	}
	createReq := content.CreateContentRequest{
		ContentType:       recommendationDebugExternalContentType(req.Source),
		ExternalSource:    recommendationDebugStringPtr(req.Source),
		ExternalContentID: recommendationDebugStringPtr(req.ExternalContentID),
		URL:               &urlValue,
		Title:             req.Title,
		Description:       recommendationDebugStringPtr(req.Description),
		ThumbnailURL:      recommendationDebugStringPtr(req.ThumbnailURL),
		Author:            recommendationDebugStringPtr(req.Author),
		Language:          language,
	}

	contentRepo := content.NewRepository(h.db)
	saved, err := contentRepo.UpsertExternalContent(c.Request.Context(), ownerUserID, createReq)
	if err != nil {
		log.Printf("[recommendation-debug] external candidate save failed scenario_id=%s lesson_id=%s url=%q err=%s", scenario.ID, req.LessonID, logsafe.URL(req.URL), logsafe.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save external candidate"})
		return
	}
	h.enqueueRecommendationDebugContentEmbedding(contentRepo, saved)

	c.JSON(http.StatusOK, gin.H{
		"content":          saved,
		"embedding_status": "queued",
	})
}

func recommendationDebugLessonExists(snapshot *recommendationDebugGeneratedLessonsSnapshot, lessonID string) bool {
	if snapshot == nil {
		return false
	}
	for _, lesson := range snapshot.Lessons {
		if lesson.LessonID == lessonID {
			return true
		}
	}
	return false
}

func recommendationDebugLessonsFromSnapshot(raw string) (*recommendationDebugGeneratedLessonsSnapshot, error) {
	var snapshot recommendationDebugGeneratedLessonsSnapshot
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &snapshot); err != nil {
		return nil, err
	}
	if len(snapshot.Lessons) == 0 {
		return nil, fmt.Errorf("generated lessons are empty")
	}
	return &snapshot, nil
}

func recommendationDebugStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func recommendationDebugExternalContentType(source string) content.ContentType {
	switch strings.TrimSpace(source) {
	case "youtube":
		return content.ContentTypeYouTube
	case "naver_blog":
		return content.ContentTypeBlog
	default:
		return content.ContentTypeArticle
	}
}

func (h *AdminHandler) resolveRecommendationDebugContentOwner(ctx context.Context, adminID uuid.UUID, scenario recommendationDebugScenarioDetailRow) (uuid.UUID, error) {
	if recommendationDebugUserExists(ctx, h.db, adminID) {
		return adminID, nil
	}
	if scenario.CreatedBy != "" {
		if createdBy, err := uuid.Parse(scenario.CreatedBy); err == nil && recommendationDebugUserExists(ctx, h.db, createdBy) {
			return createdBy, nil
		}
	}
	var fallback uuid.UUID
	if err := h.db.QueryRow(ctx, `SELECT id FROM users ORDER BY created_at ASC LIMIT 1`).Scan(&fallback); err != nil {
		return uuid.Nil, err
	}
	return fallback, nil
}

func recommendationDebugUserExists(ctx context.Context, db interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, userID uuid.UUID) bool {
	if userID == uuid.Nil {
		return false
	}
	var exists bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&exists); err != nil {
		return false
	}
	return exists
}

func (h *AdminHandler) enqueueRecommendationDebugContentEmbedding(repo *content.Repository, saved *content.Content) {
	if saved == nil {
		return
	}
	go func(item content.Content) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		provider := "embedding_gemma"
		model := recommendationDebugEmbeddingGemmaModel()
		text := embedder.BuildText(item.Title, derefString(item.Description))
		client, err := embedder.NewProviderClient(provider, strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")))
		if err != nil {
			log.Printf("[recommendation-debug] external embedding client failed content_id=%s err=%s", item.ID, logsafe.Error(err))
			_ = repo.SaveEmbeddingFailure(ctx, item.ID, provider, model, recommendationDebugEmbeddingGemmaDimension(), text, err)
			return
		}
		vector, err := client.Embed(ctx, text)
		if err != nil {
			log.Printf("[recommendation-debug] external embedding failed content_id=%s err=%s", item.ID, logsafe.Error(err))
			_ = repo.SaveEmbeddingFailure(ctx, item.ID, provider, model, recommendationDebugEmbeddingGemmaDimension(), text, err)
			return
		}
		if err := repo.SaveEmbeddingReady(ctx, item.ID, provider, model, text, vector); err != nil {
			log.Printf("[recommendation-debug] external embedding save failed content_id=%s err=%s", item.ID, logsafe.Error(err))
		}
	}(*saved)
}

func recommendationDebugEmbeddingGemmaModel() string {
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		return "embedding-gemma"
	}
	return model
}

func recommendationDebugEmbeddingGemmaDimension() int {
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
