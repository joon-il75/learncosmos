package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/extsearch"
)

const (
	externalSearchPageSize          = 5
	externalSearchMaxRawPerProvider = 10
	externalSearchTargetPerProvider = 5
)

type ExternalSearchProviderStats struct {
	Provider       string
	Query          string
	PagesAttempted int
	RawCount       int
	AcceptedCount  int
	FilteredCount  int
	ErrorReason    string
}

func CollectFilteredExternalSearchResults(
	ctx context.Context,
	providers []extsearch.SearchProvider,
	query string,
	accept func(extsearch.SearchResult) bool,
) ([]extsearch.SearchResult, []ExternalSearchProviderStats) {
	results := []extsearch.SearchResult{}
	stats := make([]ExternalSearchProviderStats, 0, len(providers))

	for _, provider := range providers {
		providerQuery := externalProviderSearchQuery(provider, query)
		stat := ExternalSearchProviderStats{
			Provider: externalProviderName(provider),
			Query:    providerQuery,
		}
		if strings.TrimSpace(providerQuery) == "" {
			stats = append(stats, stat)
			continue
		}

		pageToken := ""
		nextStart := 1
		for stat.RawCount < externalSearchMaxRawPerProvider && stat.AcceptedCount < externalSearchTargetPerProvider {
			remainingRaw := externalSearchMaxRawPerProvider - stat.RawCount
			limit := minInt(externalSearchPageSize, remainingRaw)
			pageReq := extsearch.SearchPageRequest{
				Query:     providerQuery,
				Limit:     limit,
				PageToken: pageToken,
				Start:     nextStart,
			}

			page, err := searchExternalProviderPage(ctx, provider, pageReq)
			stat.PagesAttempted++
			if err != nil {
				stat.ErrorReason = externalSearchErrorReason(err)
				break
			}

			rawCount := page.RawCount
			if rawCount <= 0 {
				rawCount = len(page.Items)
			}
			stat.RawCount += rawCount

			for _, item := range page.Items {
				if accept != nil && !accept(item) {
					stat.FilteredCount++
					continue
				}
				results = append(results, item)
				stat.AcceptedCount++
				if stat.AcceptedCount >= externalSearchTargetPerProvider {
					break
				}
			}

			if stat.AcceptedCount >= externalSearchTargetPerProvider || stat.RawCount >= externalSearchMaxRawPerProvider {
				break
			}
			if strings.TrimSpace(page.NextPageToken) != "" {
				pageToken = strings.TrimSpace(page.NextPageToken)
				nextStart = 0
				continue
			}
			if page.NextStart > 0 {
				pageToken = ""
				nextStart = page.NextStart
				continue
			}
			break
		}
		stats = append(stats, stat)
	}

	return results, stats
}

func IsInstructionalExternalSearchResult(result extsearch.SearchResult) bool {
	return isInstructionalExternalSearchResult(result)
}

func searchExternalProviderPage(ctx context.Context, provider extsearch.SearchProvider, req extsearch.SearchPageRequest) (extsearch.SearchPageResponse, error) {
	if paginated, ok := provider.(extsearch.PaginatedSearchProvider); ok {
		return paginated.SearchPage(ctx, req)
	}
	items, err := provider.Search(ctx, req.Query, req.Limit)
	if err != nil {
		return extsearch.SearchPageResponse{}, err
	}
	return extsearch.SearchPageResponse{
		Items:    items,
		RawCount: len(items),
	}, nil
}

func externalProviderName(provider extsearch.SearchProvider) string {
	switch provider.(type) {
	case *extsearch.YouTubeProvider:
		return "youtube"
	case *extsearch.NaverBlogProvider:
		return "naver_blog"
	default:
		return fmt.Sprintf("%T", provider)
	}
}

func externalSearchErrorReason(err error) string {
	if err == nil {
		return ""
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "youtube") && (strings.Contains(message, "quotaexceeded") || strings.Contains(message, "youtube.quota")):
		return "youtube_quota_exceeded"
	case strings.Contains(message, "youtube") && strings.Contains(message, "403"):
		return "youtube_search_api_error_403"
	case strings.Contains(message, "naver") && strings.Contains(message, "429"):
		return "naver_blog_rate_limited"
	case strings.Contains(message, "api key"), strings.Contains(message, "credentials"):
		return "external_provider_credentials_missing"
	default:
		return err.Error()
	}
}
