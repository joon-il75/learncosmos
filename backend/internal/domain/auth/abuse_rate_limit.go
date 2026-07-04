package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/abuseguard"
)

var (
	demoLoginRateLimiter    = abuseguard.NewFixedWindowLimiter(8, time.Minute)
	aiValidationRateLimiter = abuseguard.NewFixedWindowLimiter(6, 10*time.Minute)
)

func demoLoginRateLimitMiddleware() gin.HandlerFunc {
	return authRateLimitMiddleware(demoLoginRateLimiter, func(c *gin.Context) string {
		return "demo-login:" + normalizedHost(c.Request.Host) + ":" + c.ClientIP()
	})
}

func aiValidationRateLimitMiddleware() gin.HandlerFunc {
	return authRateLimitMiddleware(aiValidationRateLimiter, func(c *gin.Context) string {
		if userID, ok := GetCurrentUserID(c); ok && strings.TrimSpace(userID) != "" {
			return "ai-validation:user:" + userID
		}
		return "ai-validation:ip:" + c.ClientIP()
	})
}

func authRateLimitMiddleware(limiter *abuseguard.FixedWindowLimiter, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, retryAfter := limiter.Allow(keyFn(c))
		if !ok {
			seconds := int(retryAfter.Seconds())
			if retryAfter > 0 && retryAfter%time.Second != 0 {
				seconds++
			}
			if seconds < 1 {
				seconds = 1
			}
			c.Header("Retry-After", strconv.Itoa(seconds))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, apiError("rate_limited", "too many requests"))
			return
		}
		c.Next()
	}
}
