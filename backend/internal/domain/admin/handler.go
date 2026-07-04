package admin

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/config"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
	"github.com/learnweaver/backend/internal/pkg/systemsettings"
)

type AdminHandler struct {
	db                                       *pgxpool.Pool
	cfg                                      *config.Config
	authSvc                                  *auth.Service
	settingsStore                            *systemsettings.Store
	llmGateway                               *llmjobs.Gateway
	recommendationDebugGoalWorkerEnabled     bool
	recommendationDebugGoalWorkerSyncWait    time.Duration
	recommendationDebugLessonsWorkerEnabled  bool
	recommendationDebugLessonsWorkerSyncWait time.Duration
	lessonSearchPrerunReports                lessonSearchPrerunReportLister
}

func NewAdminHandler(db *pgxpool.Pool, cfg *config.Config, authSvc *auth.Service) *AdminHandler {
	return &AdminHandler{
		db:                                       db,
		cfg:                                      cfg,
		authSvc:                                  authSvc,
		settingsStore:                            systemsettings.NewStore(db),
		recommendationDebugGoalWorkerEnabled:     adminEnvBool("LLM_WORKER_FEATURE_ADMIN_RECOMMENDATION_DEBUG_GOAL_TURN", false),
		recommendationDebugGoalWorkerSyncWait:    adminEnvDuration("LLM_WORKER_ADMIN_RECOMMENDATION_DEBUG_GOAL_TURN_SYNC_WAIT", 30*time.Second),
		recommendationDebugLessonsWorkerEnabled:  adminEnvBool("LLM_WORKER_FEATURE_ADMIN_RECOMMENDATION_DEBUG_LESSONS_GENERATE", false),
		recommendationDebugLessonsWorkerSyncWait: adminEnvDuration("LLM_WORKER_ADMIN_RECOMMENDATION_DEBUG_LESSONS_GENERATE_SYNC_WAIT", 75*time.Second),
		lessonSearchPrerunReports:                lessonsearchprerun.NewReportRepository(db),
	}
}

func (h *AdminHandler) SetLLMGateway(gateway *llmjobs.Gateway) {
	h.llmGateway = gateway
}

func adminEnvBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func adminEnvDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil && d >= 0 {
		return d
	}
	if ms, err := strconv.Atoi(raw); err == nil && ms >= 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return fallback
}

func (h *AdminHandler) RegisterRoutes(r *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, headerOnlyJWT gin.HandlerFunc) {
	// Public: no auth required
	r.POST("/super-admin/login", superAdminLoginRateLimitMiddleware(), SuperAdminAuth(h.cfg, h.authSvc))
	r.GET("/public/point-policy", h.publicPointPolicyHandler)
	r.GET("/public/policies", h.GetPublicPolicies)
	r.GET("/public/policies/:type", h.GetPublicPolicyByType)
	r.GET("/public/lumi/runtime-config", h.getPublicLumiRuntimeConfigHandler)
	r.GET("/public/ui-engine-settings", h.getPublicDashboardUIEngineSettingsHandler)
	r.GET("/public/planet-texture-maps/active", h.publicActivePlanetTextureMapHandler)
	r.GET("/public/planet-texture-maps", h.publicListActivePlanetTextureMapsHandler)
	r.GET("/public/platform-notices", h.publicListPlatformNoticesHandler)
	r.GET("/public/platform-notices/:slug", h.publicGetPlatformNoticeHandler)

	// Protected: 헤더 전용 JWT (쿠키 무시) + super_admin role required
	protected := r.Group("/super-admin")
	protected.Use(headerOnlyJWT, TOTPMiddleware())
	{
		protected.GET("/users", h.listUsersHandler)
		protected.PUT("/users/:id/role", h.updateRoleHandler)
		protected.PATCH("/users/:id/role", h.patchRoleHandler)
		protected.PUT("/users/:id/premium-access", h.updatePremiumAccessHandler)

		// Settings
		protected.GET("/settings/llm", h.GetLLMSettings)
		protected.PUT("/settings/llm", h.UpdateLLMSettings)
		protected.GET("/settings/api-keys", h.GetAPIKeys)
		protected.PUT("/settings/api-keys", h.UpdateAPIKey)
		protected.GET("/settings/point-policy", h.GetPointPolicy)
		protected.PUT("/settings/point-policy", h.UpdatePointPolicy)
		protected.GET("/settings/products", h.GetPointProducts)
		protected.PUT("/settings/products", h.SavePointProducts)
		protected.GET("/settings/ad-slots", h.GetAdSlots)
		protected.PUT("/settings/ad-slots", h.SaveAdSlots)
		protected.DELETE("/settings/ad-slots/:id", h.DeleteAdSlot)
		protected.GET("/settings/affiliate", h.GetAffiliateSettings)
		protected.PUT("/settings/affiliate", h.UpdateAffiliateSettings)
		protected.GET("/settings/policies", h.GetPolicySettings)
		protected.PUT("/settings/policies", h.SavePolicyDrafts)
		protected.POST("/settings/policies/confirm", h.ConfirmPolicyDrafts)
		protected.GET("/settings/ui-engine", h.getDashboardUIEngineSettingsHandler)
		protected.PUT("/settings/ui-engine", h.updateDashboardUIEngineSettingsHandler)
		protected.GET("/assets/lumi/meta", h.getLumiAssetMetaHandler)
		protected.POST("/assets/lumi/upload", h.uploadLumiAssetHandler)
		protected.GET("/assets/planet-texture-maps", h.listPlanetTextureMapsHandler)
		protected.POST("/assets/planet-texture-maps", h.createPlanetTextureMapHandler)
		protected.PATCH("/assets/planet-texture-maps/:id", h.updatePlanetTextureMapHandler)
		protected.POST("/assets/planet-texture-maps/:id/file", h.replacePlanetTextureMapFileHandler)
		protected.DELETE("/assets/planet-texture-maps/:id", h.deletePlanetTextureMapHandler)
		protected.GET("/settings/lumi-runtime", h.getLumiRuntimeConfigHandler)
		protected.PUT("/settings/lumi-runtime", h.updateLumiRuntimeConfigHandler)
		protected.GET("/recommendation-debug/rollout-state", h.GetRecommendationRolloutState)
		protected.GET("/recommendation-debug/rollout-checkpoints", h.GetRecommendationRolloutCheckpoints)
		protected.GET("/recommendation-debug/rollout-events", h.ListRecommendationRolloutEvents)
		protected.GET("/recommendation-debug/rollout-metrics", h.ListRecommendationRolloutMetrics)
		protected.POST("/recommendation-debug/rollout-metrics/snapshot", h.CreateRecommendationRolloutMetricSnapshot)
		protected.GET("/recommendation-debug/rollout-routing-preview", h.PreviewRecommendationRolloutRouting)
		protected.PATCH("/recommendation-debug/rollout-ranker", h.UpdateRecommendationRolloutRanker)
		protected.GET("/recommendation-debug/ranker-artifact", h.GetRecommendationRankerArtifactStatus)
		protected.GET("/recommendation-debug/spec-metrics", h.GetRecommendationSpecMetrics)
		protected.GET("/recommendation-debug/lesson-search-prerun-reports", h.ListLessonSearchPrerunReports)
		protected.POST("/recommendation-debug/rollout-transition", h.TransitionRecommendationRollout)
		protected.GET("/recommendation-debug/shadow-status", h.GetRecommendationShadowStatus)
		protected.GET("/recommendation-debug/shadow-targets", h.ListRecommendationShadowBackfillTargets)
		protected.POST("/recommendation-debug/shadow-probe", h.ProbeRecommendationShadowEmbedding)
		protected.GET("/recommendation-debug/scenarios", h.ListRecommendationDebugScenarios)
		protected.POST("/recommendation-debug/scenarios", h.CreateRecommendationDebugScenario)
		protected.GET("/recommendation-debug/scenarios/:id", h.GetRecommendationDebugScenario)
		protected.POST("/recommendation-debug/scenarios/:id/goal/start", h.StartRecommendationDebugGoalInterview)
		protected.POST("/recommendation-debug/scenarios/:id/goal/message", h.SendRecommendationDebugGoalMessage)
		protected.POST("/recommendation-debug/scenarios/:id/goal/confirm", h.ConfirmRecommendationDebugGoal)
		protected.POST("/recommendation-debug/scenarios/:id/lessons/generate", h.GenerateRecommendationDebugLessons)
		protected.POST("/recommendation-debug/scenarios/:id/recommendations/compare", h.CompareRecommendationDebugScenarioRecommendations)
		protected.POST("/recommendation-debug/scenarios/:id/external-candidates/save", h.SaveRecommendationDebugExternalCandidate)
		protected.GET("/recommendation-debug/scenarios/:id/labels/summary", h.GetRecommendationDebugLabelSummary)
		protected.GET("/recommendation-debug/scenarios/:id/labels/export", h.ExportRecommendationDebugLabels)
		protected.POST("/recommendation-debug/scenarios/:id/labels", h.SaveRecommendationDebugLabel)
		protected.GET("/recommendation-debug/drafts", h.ListRecommendationDebugDrafts)
		protected.GET("/recommendation-debug/drafts/:id", h.GetRecommendationDebugDraft)
		protected.POST("/recommendation-debug/preview", h.PreviewRecommendationDebug)
		protected.POST("/recommendation-debug/preview-direct", h.PreviewRecommendationDebugDirect)
		protected.GET("/platform-notices", h.listPlatformNoticesHandler)
		protected.POST("/platform-notices", h.createPlatformNoticeHandler)
		protected.PATCH("/platform-notices/:id", h.updatePlatformNoticeHandler)
		protected.DELETE("/platform-notices/:id", h.deletePlatformNoticeHandler)
		protected.GET("/material-reports", h.listMaterialReportsHandler)
		protected.PATCH("/material-reports/:id", h.updateMaterialReportHandler)
		protected.GET("/demo-accounts", h.listDemoAccountsHandler)
		protected.POST("/demo-accounts", h.createDemoAccountsHandler)
		protected.PATCH("/demo-accounts/:id", h.updateDemoAccountHandler)
		protected.POST("/demo-accounts/:id/rotate-code", h.rotateDemoAccountCodeHandler)
		protected.GET("/demo-accounts/:id/purge-preview", h.previewDemoAccountPurgeHandler)
		protected.POST("/demo-accounts/:id/purge", h.purgeDemoAccountHandler)

		// 포인트 지급 & 감사 로그
		protected.POST("/users/:id/grant-points", h.GrantUserPoints)
		protected.GET("/audit/point-grants", h.GetPointGrantLogs)
	}

	// Admin user management: JWT + admin role
	adminMgmt := r.Group("/admin")
	adminMgmt.Use(jwtMiddleware, auth.RequireRole(auth.RoleAdmin), RequireAdminTOTPSetup(h.db), h.requireAdminTOTPSession())
	{
		adminMgmt.GET("/users", h.adminListUsersHandler)
		adminMgmt.POST("/users/:id/grant-points", h.adminGrantPointsHandler)
		adminMgmt.GET("/point-policy", h.adminGetPointPolicyHandler)
	}

	// Admin TOTP setup/verify/reset-request: JWT + admin role only
	totpSetup := r.Group("/admin/totp")
	totpSetup.Use(jwtMiddleware, auth.RequireRole(auth.RoleAdmin))
	{
		totpSetup.GET("/status", h.totpStatusHandler)
		totpSetup.POST("/setup", h.totpSetupHandler)
		totpSetup.POST("/verify", h.totpVerifyHandler)
		totpSetup.POST("/reset-request", h.totpResetRequestHandler)
	}

	// Admin TOTP validate: JWT only (로그인 흐름에서 호출)
	r.POST("/admin/totp/validate", jwtMiddleware, auth.RequireRole(auth.RoleAdmin), h.totpValidateHandler)

	// Super admin: TOTP 초기화
	protected.POST("/users/:id/totp-reset", h.totpResetHandler)

	// Admin protected routes: JWT + admin role + TOTP 설정 완료 + TOTP 코드 검증
	adminProtected := r.Group("/admin")
	adminProtected.Use(jwtMiddleware, RequireAdminTOTPSetup(h.db), RequireTOTPValidated(h.db))
	{
		// 추후 admin 전용 라우트 여기에 추가
	}
}
