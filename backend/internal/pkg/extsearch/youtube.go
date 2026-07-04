package extsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type YouTubeProvider struct {
	apiKey            string
	relevanceLanguage string
	httpClient        *http.Client
}

func NewYouTubeProvider(apiKey string) *YouTubeProvider {
	return NewYouTubeProviderWithLanguage(apiKey, "ko")
}

func NewYouTubeProviderWithLanguage(apiKey, relevanceLanguage string) *YouTubeProvider {
	return &YouTubeProvider{
		apiKey:            strings.TrimSpace(apiKey),
		relevanceLanguage: normalizeYouTubeRelevanceLanguage(relevanceLanguage),
		httpClient:        &http.Client{Timeout: 15 * time.Second},
	}
}

func normalizeYouTubeRelevanceLanguage(language string) string {
	switch strings.TrimSpace(strings.ToLower(language)) {
	case "en":
		return "en"
	default:
		return "ko"
	}
}

type youTubeSearchResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			Thumbnails   struct {
				High *struct {
					URL string `json:"url"`
				} `json:"high"`
				Medium *struct {
					URL string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

func (p *YouTubeProvider) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	resp, err := p.SearchPage(ctx, SearchPageRequest{Query: query, Limit: limit})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (p *YouTubeProvider) SearchPage(ctx context.Context, pageReq SearchPageRequest) (SearchPageResponse, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return SearchPageResponse{}, fmt.Errorf("youtube api key is required")
	}
	query := strings.TrimSpace(pageReq.Query)
	if query == "" {
		return SearchPageResponse{Items: []SearchResult{}}, nil
	}
	limit := pageReq.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 10 {
		limit = 10
	}

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("type", "video")
	params.Set("maxResults", fmt.Sprintf("%d", limit))
	params.Set("q", query)
	params.Set("key", p.apiKey)
	if p.relevanceLanguage != "" {
		params.Set("relevanceLanguage", p.relevanceLanguage)
	}
	if pageToken := strings.TrimSpace(pageReq.PageToken); pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	apiURL := "https://www.googleapis.com/youtube/v3/search?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return SearchPageResponse{}, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return SearchPageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return SearchPageResponse{}, fmt.Errorf("youtube search api error: %d", resp.StatusCode)
	}

	var data youTubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return SearchPageResponse{}, err
	}

	results := make([]SearchResult, 0, len(data.Items))
	for _, item := range data.Items {
		videoID := strings.TrimSpace(item.ID.VideoID)
		if videoID == "" {
			continue
		}
		canonicalURL := "https://www.youtube.com/watch?v=" + videoID
		thumbnailURL := ""
		if item.Snippet.Thumbnails.High != nil {
			thumbnailURL = item.Snippet.Thumbnails.High.URL
		} else if item.Snippet.Thumbnails.Medium != nil {
			thumbnailURL = item.Snippet.Thumbnails.Medium.URL
		}
		results = append(results, SearchResult{
			Source:            "youtube",
			ExternalContentID: videoID,
			CanonicalURL:      canonicalURL,
			URL:               canonicalURL,
			Title:             item.Snippet.Title,
			Description:       item.Snippet.Description,
			ThumbnailURL:      thumbnailURL,
			Author:            item.Snippet.ChannelTitle,
			Language:          p.relevanceLanguage,
		})
	}
	return SearchPageResponse{
		Items:         results,
		NextPageToken: strings.TrimSpace(data.NextPageToken),
		RawCount:      len(data.Items),
	}, nil
}
