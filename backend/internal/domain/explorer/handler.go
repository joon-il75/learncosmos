package explorer

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/systemsettings"
)

type Handler struct {
	repo           *Repository
	service        *Service
	curriculumRepo *curriculum.Repository
	contentRepo    *content.Repository
	settingsStore  *systemsettings.Store
	openAIAPIKey   string
	youtubeAPIKey  string
}

func NewHandler(
	repo *Repository,
	service *Service,
	curriculumRepo *curriculum.Repository,
	contentRepo *content.Repository,
	settingsStore *systemsettings.Store,
	openAIAPIKey string,
	youtubeAPIKey string,
) *Handler {
	return &Handler{
		repo:           repo,
		service:        service,
		curriculumRepo: curriculumRepo,
		contentRepo:    contentRepo,
		settingsStore:  settingsStore,
		openAIAPIKey:   openAIAPIKey,
		youtubeAPIKey:  youtubeAPIKey,
	}
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists || raw == nil {
		return uuid.Nil, false
	}

	switch v := raw.(type) {
	case uuid.UUID:
		return v, true
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, false
		}
		return id, true
	default:
		return uuid.Nil, false
	}
}

func (h *Handler) verifyCourseOwner(c *gin.Context, courseID, userID uuid.UUID) bool {
	ownerID, err := h.repo.GetCourseOwner(c.Request.Context(), courseID)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify course owner"})
		return false
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}

	return true
}

func (h *Handler) isPreLearningCourse(c *gin.Context, courseID uuid.UUID) (bool, bool) {
	status, err := h.repo.GetCourseStatus(c.Request.Context(), courseID)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
			return false, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check course status"})
		return false, false
	}
	return status == "draft" || status == "confirmed", true
}

func (h *Handler) resolveParentCourseID(c *gin.Context, parentKind ParentKind, parentID uuid.UUID) (uuid.UUID, error) {
	return h.resolveParentCourseIDFromBackground(c.Request.Context(), parentKind, parentID)
}

func (h *Handler) resolveParentCourseIDFromBackground(ctx context.Context, parentKind ParentKind, parentID uuid.UUID) (uuid.UUID, error) {
	switch parentKind {
	case ParentKindRegion:
		region, err := h.repo.GetRegionByID(ctx, parentID)
		if err != nil {
			return uuid.Nil, ErrNotFound
		}
		return region.CourseDraftID, nil
	case ParentKindSubRegion:
		subregion, err := h.repo.GetSubRegionByID(ctx, parentID)
		if err != nil {
			return uuid.Nil, ErrNotFound
		}
		region, err := h.repo.GetRegionByID(ctx, subregion.RegionID)
		if err != nil {
			return uuid.Nil, ErrNotFound
		}
		return region.CourseDraftID, nil
	default:
		return uuid.Nil, ErrNotFound
	}
}

func (h *Handler) writeResolveParentError(c *gin.Context, err error) {
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "parent not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve parent"})
}
