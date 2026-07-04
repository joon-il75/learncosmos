package metaparser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// extractVideoID — YouTube URL에서 video ID 추출
func extractVideoID(rawURL string) (string, error) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`),
		regexp.MustCompile(`youtube\.com/embed/([a-zA-Z0-9_-]{11})`),
		regexp.MustCompile(`youtube\.com/shorts/([a-zA-Z0-9_-]{11})`),
	}
	for _, p := range patterns {
		m := p.FindStringSubmatch(rawURL)
		if len(m) >= 2 {
			return m[1], nil
		}
	}
	return "", fmt.Errorf("YouTube video ID를 찾을 수 없습니다")
}

// parseDuration — ISO 8601 duration (PT1H2M3S) → seconds
func parseDuration(iso string) int {
	d, err := time.ParseDuration(
		strings.NewReplacer("PT", "", "H", "h", "M", "m", "S", "s").Replace(iso),
	)
	if err != nil {
		return 0
	}
	return int(d.Seconds())
}

type ytAPIResponse struct {
	Items []struct {
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			Thumbnails   struct {
				MaxRes *struct {
					URL string `json:"url"`
				} `json:"maxres"`
				High *struct {
					URL string `json:"url"`
				} `json:"high"`
				Medium *struct {
					URL string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// ParseYouTube — YouTube Data API v3로 메타 정보 추출
func ParseYouTube(ctx context.Context, rawURL string, apiKey string) (*MetaInfo, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("YouTube API 키 없음")
	}

	videoID, err := extractVideoID(rawURL)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet,contentDetails",
		url.QueryEscape(videoID), url.QueryEscape(apiKey),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API 오류: %d", resp.StatusCode)
	}

	var data ytAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Items) == 0 {
		return nil, fmt.Errorf("영상을 찾을 수 없습니다")
	}

	item := data.Items[0]
	s := item.Snippet

	// 썸네일: maxres → high → medium 순서로 선택
	thumbnailURL := ""
	if s.Thumbnails.MaxRes != nil {
		thumbnailURL = s.Thumbnails.MaxRes.URL
	} else if s.Thumbnails.High != nil {
		thumbnailURL = s.Thumbnails.High.URL
	} else if s.Thumbnails.Medium != nil {
		thumbnailURL = s.Thumbnails.Medium.URL
	}

	return &MetaInfo{
		Title:           s.Title,
		Description:     s.Description,
		ThumbnailURL:    thumbnailURL,
		Author:          s.ChannelTitle,
		DurationSeconds: parseDuration(item.ContentDetails.Duration),
		ContentType:     "youtube",
		SourceURL:       rawURL,
	}, nil
}
