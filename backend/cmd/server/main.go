package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/config"
	"github.com/learnweaver/backend/internal/domain/admin"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/domain/demoaccount"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/user"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
	"github.com/learnweaver/backend/internal/jobs"
	"github.com/learnweaver/backend/internal/pkg/cache"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/router"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal("invalid config: ", err)
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer pool.Close()

	redisClient, err := cache.NewClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}

	userRepo := user.NewPostgresRepository(pool)

	oauthCfg := auth.OAuthConfig{
		Google: auth.ProviderConfig{
			ClientID:     cfg.OAuth.Google.ClientID,
			ClientSecret: cfg.OAuth.Google.ClientSecret,
			RedirectURL:  cfg.OAuth.Google.RedirectURL,
			AuthURL:      cfg.OAuth.Google.AuthURL,
			TokenURL:     cfg.OAuth.Google.TokenURL,
			UserInfoURL:  cfg.OAuth.Google.UserInfoURL,
		},
		Kakao: auth.ProviderConfig{
			ClientID:     cfg.OAuth.Kakao.ClientID,
			ClientSecret: cfg.OAuth.Kakao.ClientSecret,
			RedirectURL:  cfg.OAuth.Kakao.RedirectURL,
			AuthURL:      cfg.OAuth.Kakao.AuthURL,
			TokenURL:     cfg.OAuth.Kakao.TokenURL,
			UserInfoURL:  cfg.OAuth.Kakao.UserInfoURL,
		},
		Naver: auth.ProviderConfig{
			ClientID:     cfg.OAuth.Naver.ClientID,
			ClientSecret: cfg.OAuth.Naver.ClientSecret,
			RedirectURL:  cfg.OAuth.Naver.RedirectURL,
			AuthURL:      cfg.OAuth.Naver.AuthURL,
			TokenURL:     cfg.OAuth.Naver.TokenURL,
			UserInfoURL:  cfg.OAuth.Naver.UserInfoURL,
		},
	}

	authSvc := auth.NewService(userRepo, redisClient, cfg.JWTSecret, oauthCfg)
	demoAccountRepo := demoaccount.NewRepository(pool)
	authHandler := auth.NewHandler(authSvc, demoAccountRepo, auth.DemoLoginConfig{
		Enabled:      cfg.DemoLoginEnabled,
		Secret:       cfg.DemoLoginSecret,
		AllowedHosts: cfg.DemoLoginAllowedHosts,
	})

	adminHandler := admin.NewAdminHandler(pool, cfg, authSvc)

	if cfg.ContentHealthCheckEnabled {
		contentRepo := content.NewRepository(pool)
		healthChecker := jobs.NewContentHealthChecker(contentRepo, log.Default())
		healthChecker.Start(ctx, 24*time.Hour, 200)
		log.Println("content health checker enabled")
	}

	if cfg.LessonSearchPrerunSchedulerEnabled {
		prerunScheduler := jobs.NewLessonSearchPrerunScheduler(jobs.LessonSearchPrerunSchedulerConfig{
			Enabled:        cfg.LessonSearchPrerunSchedulerEnabled,
			Interval:       time.Duration(cfg.LessonSearchPrerunIntervalHours) * time.Hour,
			StartDelay:     time.Duration(cfg.LessonSearchPrerunStartDelayMinutes) * time.Minute,
			Timeout:        time.Duration(cfg.LessonSearchPrerunTimeoutSeconds) * time.Second,
			Provider:       cfg.LessonSearchPrerunProvider,
			TopN:           cfg.LessonSearchPrerunTopN,
			RiskTopN:       cfg.LessonSearchPrerunRiskTopN,
			RiskFixtureIDs: cfg.LessonSearchPrerunRiskFixtureIDs,
			Retention:      time.Duration(cfg.LessonSearchPrerunReportRetentionDays) * 24 * time.Hour,
			YouTubeAPIKey:  cfg.YouTubeAPIKey,
			NaverClientID:  cfg.OAuth.Naver.ClientID,
			NaverSecret:    cfg.OAuth.Naver.ClientSecret,
		}, lessonsearchprerun.NewReportRepository(pool), jobs.NewPostgresAdvisoryLocker(pool, 0), log.Default())
		prerunScheduler.Start(ctx)
		log.Println("lesson search prerun scheduler enabled")
	}

	if cfg.LLMWorkerEnabled {
		llmJobRepo := llmjobs.NewRepository(pool)
		llmQueue := llmjobs.NewRedisStreamQueue(redisClient, llmjobs.DefaultStreamKey, llmjobs.DefaultConsumerGroup, llmjobs.DefaultDLQKey)
		if err := llmQueue.EnsureGroup(ctx); err != nil {
			log.Fatal("failed to initialize llm worker queue:", err)
		}
		llmGateway := llmjobs.NewGateway(llmJobRepo, llmQueue)
		adminHandler.SetLLMGateway(llmGateway)
		runtime := llmjobs.NewRuntime(llmJobRepo, llmQueue, llmjobs.RuntimeOptions{
			WorkerCount: cfg.LLMWorkerCount,
			Logger:      log.Default(),
		})
		runtime.RegisterProcessor(llmjobs.FeatureDummy, llmjobs.NewDummyProcessor())
		curriculumRepo := curriculum.NewRepository(pool)
		curriculumGoalRepo := goal.NewRepository(pool)
		curriculumHandler := curriculum.NewHandlerWithRedis(curriculumRepo, curriculumGoalRepo, curriculum.NewService(), redisClient)
		runtime.RegisterProcessor(llmjobs.FeatureCourseDraftCreate, curriculum.NewCourseDraftCreateProcessor(curriculumHandler))
		runtime.RegisterProcessor(llmjobs.FeaturePointAISummary, curriculum.NewPointAISummaryProcessor(curriculumHandler))
		runtime.RegisterProcessor(llmjobs.FeaturePointQuestionGenerate, curriculum.NewPointQuestionGenerateProcessor(curriculumHandler))
		runtime.RegisterProcessor(llmjobs.FeaturePointFeedbackGenerate, curriculum.NewPointFeedbackGenerateProcessor(curriculumHandler))
		runtime.RegisterProcessor(llmjobs.FeaturePointSelfEvaluationDraft, curriculum.NewPointSelfEvaluationDraftProcessor(curriculumHandler))
		runtime.RegisterProcessor(llmjobs.FeatureAdminRecommendationDebugGoalTurn, admin.NewRecommendationDebugGoalTurnProcessor(adminHandler))
		runtime.RegisterProcessor(llmjobs.FeatureAdminRecommendationDebugLessonsGenerate, admin.NewRecommendationDebugLessonsGenerateProcessor(adminHandler))
		goalRepo := goal.NewRepository(pool)
		goalHandler := goal.NewHandler(goalRepo, goal.NewService())
		runtime.RegisterProcessor(llmjobs.FeatureGoalInterviewTurn, goal.NewGoalInterviewTurnProcessor(goalHandler))
		runtime.Start(ctx)
	}

	r := router.NewRouter(
		authHandler,
		adminHandler,
		pool,
		redisClient,
		authSvc.JWTMiddleware(),
		authSvc.JWTMiddlewareHeaderOnly(),
		authSvc.RequireRequiredConsents(),
	)

	log.Printf("LearnWeaver backend starting on :%s", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
