package explorer

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/systemsettings"
)

func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, authMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc, alphaAccessMiddleware gin.HandlerFunc) {
	repo := NewRepository(pool)
	service := NewService()
	contentRepo := content.NewRepository(pool)
	curriculumRepo := curriculum.NewRepository(pool)
	settingsStore := systemsettings.NewStore(pool)
	handler := NewHandler(
		repo,
		service,
		curriculumRepo,
		contentRepo,
		settingsStore,
		os.Getenv("OPENAI_API_KEY"),
		os.Getenv("YOUTUBE_API_KEY"),
	)

	ex := rg.Group("/explorer", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		ex.GET("/course/:courseId", handler.GetCourse)

		ex.POST("/region", handler.CreateRegion)
		ex.PATCH("/region/:regionId", handler.UpdateRegion)
		ex.DELETE("/region/:regionId", handler.DeleteRegion)
		ex.PATCH("/region/:regionId/status", handler.UpdateRegionStatus)

		ex.POST("/subregion", handler.CreateSubRegion)
		ex.PATCH("/subregion/:subRegionId", handler.UpdateSubRegion)
		ex.DELETE("/subregion/:subRegionId", handler.DeleteSubRegion)
		ex.PATCH("/subregion/:subRegionId/status", handler.UpdateSubRegionStatus)

		ex.POST("/exploration-node", handler.CreateExplorationNode)
		ex.POST("/research-node", handler.CreateResearchNode)
		ex.GET("/node/:nodeId", handler.GetNode)
		ex.PATCH("/node/:nodeId", handler.UpdateNode)
		ex.DELETE("/node/:nodeId", handler.DeleteNode)
		ex.PATCH("/node/:nodeId/status", handler.UpdateNodeStatus)

		ex.PATCH("/research-node/:nodeId/type", handler.SetResearchNodeType)
		ex.POST("/exploration-node/link-check", handler.CheckExplorationNodeLink)

		ex.POST("/save", handler.SavePlan)
	}
}
