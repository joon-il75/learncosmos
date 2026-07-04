package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/pkg/abuseguard"
)

var superAdminLoginRateLimiter = abuseguard.NewFixedWindowLimiter(5, time.Minute)

func superAdminLoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, retryAfter := superAdminLoginRateLimiter.Allow("super-admin-login:" + c.ClientIP())
		if !ok {
			seconds := int(retryAfter.Seconds())
			if retryAfter > 0 && retryAfter%time.Second != 0 {
				seconds++
			}
			if seconds < 1 {
				seconds = 1
			}
			c.Header("Retry-After", strconv.Itoa(seconds))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests", "error_code": "rate_limited"})
			return
		}
		c.Next()
	}
}
