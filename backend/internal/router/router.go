package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/admin"
	"github.com/learnweaver/backend/internal/domain/alphaaccess"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/explorer"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/redis/go-redis/v9"
)

var localTrustedProxies = []string{"127.0.0.1", "::1"}

func configureTrustedProxies(r *gin.Engine) error {
	return r.SetTrustedProxies(localTrustedProxies)
}

func NewRouter(authHandler *auth.Handler, adminHandler *admin.AdminHandler, pool *pgxpool.Pool, redisClient *redis.Client, jwtMiddleware gin.HandlerFunc, headerOnlyJWT gin.HandlerFunc, consentMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	if err := configureTrustedProxies(r); err != nil {
		panic(err)
	}
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(unsafeRequestOriginGuard())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	alphaAccessHandler := alphaaccess.NewHandler(alphaaccess.NewRepository(pool))
	authHandler.RegisterRoutes(v1)
	adminHandler.RegisterRoutes(v1, jwtMiddleware, headerOnlyJWT)
	alphaAccessHandler.RegisterRoutes(v1, jwtMiddleware, consentMiddleware)
	alphaSuperAdmin := v1.Group("/super-admin")
	alphaSuperAdmin.Use(headerOnlyJWT, admin.TOTPMiddleware())
	alphaAccessHandler.RegisterSuperAdminRoutes(alphaSuperAdmin)
	alphaAccessMiddleware := alphaAccessHandler.RequireAlphaAccess()
	content.RegisterRoutes(v1, pool, jwtMiddleware, consentMiddleware, alphaAccessMiddleware)
	curriculum.RegisterRoutes(v1, pool, redisClient, jwtMiddleware, consentMiddleware, alphaAccessMiddleware)
	explorer.RegisterRoutes(v1, pool, jwtMiddleware, consentMiddleware, alphaAccessMiddleware)
	goal.RegisterRoutes(v1, pool, redisClient, jwtMiddleware, consentMiddleware, alphaAccessMiddleware)
	llmjobs.RegisterRoutes(v1, pool, redisClient, jwtMiddleware, headerOnlyJWT, consentMiddleware)
	safety.RegisterRoutes(v1, pool, jwtMiddleware, headerOnlyJWT, consentMiddleware, alphaAccessMiddleware, admin.TOTPMiddleware())

	return r
}
