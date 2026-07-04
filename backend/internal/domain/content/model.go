package content

import (
	"time"

	"github.com/google/uuid"
)

// ContentType 허용 타입
type ContentType string

const (
	ContentTypeYouTube  ContentType = "youtube"
	ContentTypeBlog     ContentType = "blog"
	ContentTypeArticle  ContentType = "article"
	ContentTypeInternal ContentType = "internal"
)

type ContentStatus string

const (
	ContentStatusActive     ContentStatus = "active"
	ContentStatusRedirect   ContentStatus = "redirect"
	ContentStatusRestricted ContentStatus = "restricted"
	ContentStatusBroken     ContentStatus = "broken"
	ContentStatusUnknown    ContentStatus = "unknown"
)

// Content — contents 테이블 도메인 모델
type Content struct {
	ID                uuid.UUID     `json:"id"`
	UserID            uuid.UUID     `json:"user_id"`
	ContentType       ContentType   `json:"content_type"`
	ExternalSource    *string       `json:"external_source"`
	ExternalContentID *string       `json:"external_content_id"`
	CanonicalURL      *string       `json:"canonical_url"`
	URL               *string       `json:"url"`
	Title             string        `json:"title"`
	Description       *string       `json:"description"`
	ThumbnailURL      *string       `json:"thumbnail_url"`
	DurationSeconds   *int          `json:"duration_seconds"`
	Author            *string       `json:"author"`
	Language          string        `json:"language"`
	IsPublic          bool          `json:"is_public"`
	QualityScore      *float64      `json:"quality_score"`
	ContentStatus     ContentStatus `json:"content_status"`
	HealthScore       float64       `json:"health_score"`
	LastCheckedAt     *time.Time    `json:"last_checked_at"`
	HTTPStatus        *int          `json:"http_status"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	// embedding은 API 응답에서 제외 (용량 문제)
}

// CreateContentRequest — POST /api/v1/contents 요청 바디
type CreateContentRequest struct {
	ContentType       ContentType `json:"content_type" binding:"required"`
	ExternalSource    *string     `json:"external_source"`
	ExternalContentID *string     `json:"external_content_id"`
	URL               *string     `json:"url"`
	Title             string      `json:"title" binding:"required,max=200"`
	Description       *string     `json:"description"`
	ThumbnailURL      *string     `json:"thumbnail_url"`
	DurationSeconds   *int        `json:"duration_seconds"`
	Author            *string     `json:"author"`
	Language          string      `json:"language"`
}

// UpdateContentRequest — PATCH /api/v1/contents/:id 요청 바디
// 모든 필드 optional (nil = 변경 없음)
type UpdateContentRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	ThumbnailURL *string `json:"thumbnail_url"`
	Author       *string `json:"author"`
	Language     *string `json:"language"`
	URL          *string `json:"url"`
}

// ListContentsQuery — GET /api/v1/contents 쿼리 파라미터
type ListContentsQuery struct {
	ContentType string `form:"type"`  // 필터: youtube/blog/article/internal
	Page        int    `form:"page"`  // 기본값 1
	Limit       int    `form:"limit"` // 기본값 20, 최대 100
}

// ListContentsResponse — 목록 응답 (페이지네이션 포함)
type ListContentsResponse struct {
	Contents []Content `json:"contents"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}
