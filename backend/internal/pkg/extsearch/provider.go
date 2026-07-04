package extsearch

import "context"

type SearchResult struct {
	Source            string
	ExternalContentID string
	CanonicalURL      string
	URL               string
	Title             string
	Description       string
	ThumbnailURL      string
	Author            string
	Language          string
}

type SearchProvider interface {
	Search(ctx context.Context, query string, limit int) ([]SearchResult, error)
}

type SearchPageRequest struct {
	Query     string
	Limit     int
	PageToken string
	Start     int
}

type SearchPageResponse struct {
	Items         []SearchResult
	NextPageToken string
	NextStart     int
	RawCount      int
}

type PaginatedSearchProvider interface {
	SearchProvider
	SearchPage(ctx context.Context, req SearchPageRequest) (SearchPageResponse, error)
}
