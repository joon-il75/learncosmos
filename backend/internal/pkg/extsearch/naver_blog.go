package extsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type NaverBlogProvider struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewNaverBlogProvider(clientID, clientSecret string) *NaverBlogProvider {
	return &NaverBlogProvider{
		clientID:     strings.TrimSpace(clientID),
		clientSecret: strings.TrimSpace(clientSecret),
		httpClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

type naverBlogSearchResponse struct {
	Items []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
		BloggerName string `json:"bloggername"`
		BloggerLink string `json:"bloggerlink"`
		PostDate    string `json:"postdate"`
	} `json:"items"`
}

var naverHTMLTagPattern = regexp.MustCompile(`<[^>]+>`)

func (p *NaverBlogProvider) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	resp, err := p.SearchPage(ctx, SearchPageRequest{Query: query, Limit: limit, Start: 1})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (p *NaverBlogProvider) SearchPage(ctx context.Context, pageReq SearchPageRequest) (SearchPageResponse, error) {
	if strings.TrimSpace(p.clientID) == "" || strings.TrimSpace(p.clientSecret) == "" {
		return SearchPageResponse{}, fmt.Errorf("naver blog api credentials are required")
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
	start := pageReq.Start
	if start <= 0 {
		start = 1
	}

	apiURL := fmt.Sprintf(
		"https://openapi.naver.com/v1/search/blog.json?query=%s&display=%d&start=%d&sort=sim",
		url.QueryEscape(query),
		limit,
		start,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return SearchPageResponse{}, err
	}
	req.Header.Set("X-Naver-Client-Id", p.clientID)
	req.Header.Set("X-Naver-Client-Secret", p.clientSecret)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return SearchPageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchPageResponse{}, fmt.Errorf("naver blog search api error: %d", resp.StatusCode)
	}

	var data naverBlogSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return SearchPageResponse{}, err
	}

	results := make([]SearchResult, 0, len(data.Items))
	for _, item := range data.Items {
		link := strings.TrimSpace(item.Link)
		if link == "" {
			continue
		}

		results = append(results, SearchResult{
			Source:            "naver_blog",
			ExternalContentID: naverBlogExternalID(link),
			CanonicalURL:      link,
			URL:               link,
			Title:             stripNaverHTML(item.Title),
			Description:       stripNaverHTML(item.Description),
			ThumbnailURL:      "",
			Author:            strings.TrimSpace(item.BloggerName),
			Language:          "ko",
		})
	}

	nextStart := 0
	if len(data.Items) >= limit && start+limit <= 1000 {
		nextStart = start + limit
	}
	return SearchPageResponse{
		Items:     results,
		NextStart: nextStart,
		RawCount:  len(data.Items),
	}, nil
}

func stripNaverHTML(value string) string {
	cleaned := naverHTMLTagPattern.ReplaceAllString(value, "")
	cleaned = html.UnescapeString(cleaned)
	return strings.Join(strings.Fields(strings.TrimSpace(cleaned)), " ")
}

func naverBlogExternalID(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return strings.TrimSpace(rawURL)
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Host))
	path := strings.Trim(strings.TrimSpace(parsed.Path), "/")
	if host == "" {
		return path
	}
	if path == "" {
		return host
	}
	return host + "/" + path
}
