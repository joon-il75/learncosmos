package safety

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, authMiddleware gin.HandlerFunc, headerOnlyAuthMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc, alphaAccessMiddleware gin.HandlerFunc, superAdminMiddleware gin.HandlerFunc) {
	handler := NewHandler(NewService(NewRepository(pool)))

	userGroup := rg.Group("/safety", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		userGroup.POST("/moderate", handler.Moderate)
	}

	superAdminGroup := rg.Group("/super-admin/safety", headerOnlyAuthMiddleware, superAdminMiddleware)
	{
		superAdminGroup.GET("/moderation-logs", handler.ListLogs)
		superAdminGroup.GET("/ai-output-summary", handler.AIOutputSummary)
		superAdminGroup.GET("/ai-output-review-candidates", handler.AIOutputReviewCandidates)
		superAdminGroup.PATCH("/ai-output-review-candidates/review", handler.SaveAIOutputReviewDecision)
		superAdminGroup.GET("/moderation-rules", handler.ListRules)
		superAdminGroup.POST("/moderation-rules", handler.CreateRule)
		superAdminGroup.PATCH("/moderation-rules/:rule_id", handler.UpdateRule)
		superAdminGroup.POST("/moderation-rules/:rule_id/deactivate", handler.DeactivateRule)
		superAdminGroup.POST("/moderation-rules/:rule_id/reactivate", handler.ReactivateRule)
	}
}
