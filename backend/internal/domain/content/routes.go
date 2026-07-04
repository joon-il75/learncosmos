package content

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/safety"
)

// RegisterRoutes — router.go에서 호출
func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, authMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc, alphaAccessMiddleware gin.HandlerFunc) {
	repo := NewRepository(pool)
	handler := NewHandler(repo)
	handler.SetSafetyService(safety.NewService(safety.NewRepository(pool)))

	contents := rg.Group("/contents", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		contents.POST("/parse", handler.Parse) // /:id 앞에 위치해야 함
		contents.POST("", handler.Create)
		contents.GET("", handler.List)
		contents.GET("/:id", handler.GetByID)
		contents.PATCH("/:id", handler.Update)
		contents.DELETE("/:id", handler.Delete)
	}
}
