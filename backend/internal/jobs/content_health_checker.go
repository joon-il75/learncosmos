package jobs

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/pkg/safehttp"
)

type ContentHealthChecker struct {
	repo   *content.Repository
	client *http.Client
	logger *log.Logger
}

func NewContentHealthChecker(repo *content.Repository, logger *log.Logger) *ContentHealthChecker {
	return &ContentHealthChecker{
		repo:   repo,
		client: safehttp.NewClient(safehttp.ClientConfig{Timeout: 10 * time.Second}),
		logger: logger,
	}
}

func (c *ContentHealthChecker) Start(ctx context.Context, interval time.Duration, batchSize int) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 200
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.RunOnce(ctx, batchSize); err != nil && c.logger != nil {
					c.logger.Printf("content health checker failed: %v", err)
				}
			}
		}
	}()
}

func (c *ContentHealthChecker) RunOnce(ctx context.Context, batchSize int) error {
	checkedAt := time.Now()
	targets, err := c.repo.ListHealthCheckTargets(ctx, batchSize)
	if err != nil {
		return err
	}

	for _, target := range targets {
		status, healthScore, httpStatus := c.checkURL(ctx, target.URL)
		if err := c.repo.UpdateContentHealth(ctx, target.ID, status, healthScore, httpStatus, checkedAt); err != nil {
			return err
		}
	}

	return nil
}

func (c *ContentHealthChecker) checkURL(ctx context.Context, url string) (content.ContentStatus, float64, *int) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return content.ContentStatusBroken, 0.0, nil
	}

	resp, err := c.client.Do(req)
	if err != nil || resp == nil {
		return content.ContentStatusBroken, 0.0, nil
	}
	defer resp.Body.Close()

	return interpretHealthHTTPStatus(resp.StatusCode)
}

func interpretHealthHTTPStatus(code int) (content.ContentStatus, float64, *int) {
	statusCode := code
	switch {
	case code >= 200 && code < 300:
		return content.ContentStatusActive, 1.0, &statusCode
	case code == http.StatusMovedPermanently || code == http.StatusFound || code == http.StatusTemporaryRedirect || code == http.StatusPermanentRedirect:
		return content.ContentStatusRedirect, 0.8, &statusCode
	case code == http.StatusForbidden || code == http.StatusUnauthorized:
		return content.ContentStatusRestricted, 0.4, &statusCode
	case code == http.StatusNotFound || code >= 500:
		return content.ContentStatusBroken, 0.0, &statusCode
	default:
		return content.ContentStatusUnknown, 0.5, &statusCode
	}
}
