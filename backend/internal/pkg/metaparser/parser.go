package metaparser

import (
	"context"
	"net/url"
	"strings"
)

// MetaInfo — URL에서 추출한 메타 정보
type MetaInfo struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	ThumbnailURL    string `json:"thumbnail_url"`
	Author          string `json:"author"`
	DurationSeconds int    `json:"duration_seconds"` // YouTube 전용, 나머지 0
	ContentType     string `json:"content_type"`     // 감지된 타입
	SourceURL       string `json:"source_url"`
}

// Parse — URL을 분석해 적절한 파서로 라우팅
func Parse(ctx context.Context, rawURL string, youtubeAPIKey string) (*MetaInfo, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	host := strings.ToLower(u.Host)

	// YouTube 감지
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		info, err := ParseYouTube(ctx, rawURL, youtubeAPIKey)
		if err == nil {
			return info, nil
		}
		// API 실패 시 OGP 폴백
	}

	// OGP 파싱 (블로그, 아티클 등)
	return ParseOGP(ctx, rawURL)
}

// DetectContentType — URL 패턴으로 content_type 감지
func DetectContentType(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "article"
	}
	host := strings.ToLower(u.Host)
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		return "youtube"
	}
	blogPlatforms := []string{"tistory.com", "velog.io", "medium.com", "brunch.co.kr", "naver.com/PostView"}
	for _, p := range blogPlatforms {
		if strings.Contains(rawURL, p) {
			return "blog"
		}
	}
	return "article"
}
