package curriculum

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/content"
)

func (h *Handler) resolveSelectedExternalResources(ctx context.Context, userID uuid.UUID, draft *DraftAggregate) error {
	if draft == nil {
		return nil
	}
	contentLanguage := normalizeLearningLanguage(draft.Draft.GenerationLanguage)
	contentRepo := content.NewRepository(h.repo.pool)
	for mainIdx := range draft.Lessons {
		for subIdx := range draft.Lessons[mainIdx].SubLessons {
			for ptIdx := range draft.Lessons[mainIdx].SubLessons[subIdx].Points {
				point := &draft.Lessons[mainIdx].SubLessons[subIdx].Points[ptIdx]
				if point.Point.PointType != PointTypeExploration {
					continue
				}
				sel := point.Point.SelectionState
				if sel == nil || *sel != ResourceSelectionSelected || point.Point.ContentID != nil || point.Point.ExternalURL == nil {
					continue
				}
				createReq, ok := buildExternalContentCreateRequestFromPoint(point.Point, contentLanguage)
				if !ok {
					continue
				}
				saved, err := contentRepo.UpsertExternalContent(ctx, userID, createReq)
				if err != nil || saved == nil {
					return err
				}
				point.Point.ContentID = &saved.ID
				if point.Point.ExternalURL == nil && saved.URL != nil {
					point.Point.ExternalURL = saved.URL
				}
				if _, dbErr := h.repo.pool.Exec(ctx, `
					UPDATE course_draft_points
					SET content_id = $1, updated_at = NOW()
					WHERE id = $2
				`, saved.ID, point.Point.ID); dbErr != nil {
					return fmt.Errorf("persist external content ID: %w", dbErr)
				}
			}
		}
	}
	return nil
}

func buildExternalContentCreateRequestFromPoint(point CourseDraftPoint, language string) (content.CreateContentRequest, bool) {
	if point.ExternalURL == nil || strings.TrimSpace(*point.ExternalURL) == "" {
		return content.CreateContentRequest{}, false
	}
	urlValue := strings.TrimSpace(*point.ExternalURL)
	contentType, externalSource, externalContentID := inferExternalContentIdentity(urlValue)
	req := content.CreateContentRequest{
		ContentType:       contentType,
		ExternalSource:    externalSource,
		ExternalContentID: externalContentID,
		URL:               stringPointer(urlValue),
		Title:             strings.TrimSpace(point.Title),
		Description:       point.Description,
		ThumbnailURL:      point.ThumbnailURL,
		Language:          normalizeLearningLanguage(language),
	}
	return req, true
}

func inferExternalContentIdentity(rawURL string) (content.ContentType, *string, *string) {
	parsed, err := neturl.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return content.ContentTypeArticle, nil, nil
	}
	host := strings.ToLower(parsed.Host)
	switch {
	case strings.Contains(host, "youtube.com"), strings.Contains(host, "youtu.be"):
		videoID := strings.TrimSpace(parsed.Query().Get("v"))
		if videoID == "" && strings.Contains(host, "youtu.be") {
			videoID = strings.Trim(strings.TrimSpace(parsed.Path), "/")
		}
		source := "youtube"
		if videoID == "" {
			return content.ContentTypeYouTube, &source, nil
		}
		return content.ContentTypeYouTube, &source, &videoID
	case strings.Contains(host, "blog.naver.com"):
		source := "naver_blog"
		logNo := strings.Trim(strings.TrimSpace(parsed.Path), "/")
		if logNo == "" {
			return content.ContentTypeBlog, &source, nil
		}
		return content.ContentTypeBlog, &source, &logNo
	default:
		return content.ContentTypeArticle, nil, nil
	}
}

func (h *Handler) SelectDraftResource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	resourceID, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	if err := h.repo.SelectDraftLessonResource(c.Request.Context(), draftID, userID, resourceID); err != nil {
		switch {
		case errors.Is(err, errDraftResourceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		case errors.Is(err, errDraftResourceInvalidState):
			c.JSON(http.StatusBadRequest, gin.H{"error": "resource is not a selectable candidate"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select resource"})
		}
		return
	}

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resource selected but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": fetchedDraft})
}

func (h *Handler) AttachDraftLessonResource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req AttachDraftLessonResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.AttachDraftLessonResource(c.Request.Context(), draftID, userID, lessonID, req); err != nil {
		switch {
		case errors.Is(err, errDraftLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		case errors.Is(err, errDraftContentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		case errors.Is(err, errDraftResourceDuplicate):
			c.JSON(http.StatusConflict, gin.H{"error": "resource already attached to lesson"})
		case errors.Is(err, errDraftResourceInvalidState):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource selection state"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach resource"})
		}
		return
	}

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resource attached but failed to reload"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"draft": fetchedDraft})
}
