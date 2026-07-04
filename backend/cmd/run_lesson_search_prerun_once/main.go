package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/config"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
	"github.com/learnweaver/backend/internal/jobs"
	"github.com/learnweaver/backend/internal/pkg/db"
)

func main() {
	envFile := strings.TrimSpace(os.Getenv("LEARNWEAVER_BACKEND_ENV_FILE"))
	if envFile == "" {
		envFile = "/run/learnweaver/backend.env"
	}
	_ = godotenv.Load(envFile)
	_ = godotenv.Load()

	provider := flag.String("provider", "", "provider filter: all, youtube, naver_blog")
	topN := flag.Int("top-n", 0, "main fixture top-N override")
	riskTopN := flag.Int("risk-top-n", -1, "risk fixture top-N override")
	riskFixtureIDs := flag.String("risk-fixture-ids", "", "space or comma separated risk fixture ids override")
	timeoutSeconds := flag.Int("timeout-seconds", 0, "run timeout seconds override")
	retentionDays := flag.Int("retention-days", 0, "report retention days override")
	runType := flag.String("run-type", "manual_once", "report run_type")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if *provider != "" {
		cfg.LessonSearchPrerunProvider = *provider
	}
	if *topN > 0 {
		cfg.LessonSearchPrerunTopN = *topN
	}
	if *riskTopN >= 0 {
		cfg.LessonSearchPrerunRiskTopN = *riskTopN
	}
	if *riskFixtureIDs != "" {
		cfg.LessonSearchPrerunRiskFixtureIDs = splitFixtureIDs(*riskFixtureIDs)
	}
	if *timeoutSeconds > 0 {
		cfg.LessonSearchPrerunTimeoutSeconds = *timeoutSeconds
	}
	if *retentionDays > 0 {
		cfg.LessonSearchPrerunReportRetentionDays = *retentionDays
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.LessonSearchPrerunTimeoutSeconds+10)*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer pool.Close()

	scheduler := jobs.NewLessonSearchPrerunScheduler(jobs.LessonSearchPrerunSchedulerConfig{
		Enabled:        true,
		Interval:       time.Duration(cfg.LessonSearchPrerunIntervalHours) * time.Hour,
		StartDelay:     0,
		Timeout:        time.Duration(cfg.LessonSearchPrerunTimeoutSeconds) * time.Second,
		Provider:       cfg.LessonSearchPrerunProvider,
		TopN:           cfg.LessonSearchPrerunTopN,
		RiskTopN:       cfg.LessonSearchPrerunRiskTopN,
		RiskFixtureIDs: cfg.LessonSearchPrerunRiskFixtureIDs,
		Retention:      time.Duration(cfg.LessonSearchPrerunReportRetentionDays) * 24 * time.Hour,
		RunType:        strings.TrimSpace(*runType),
		YouTubeAPIKey:  cfg.YouTubeAPIKey,
		NaverClientID:  cfg.OAuth.Naver.ClientID,
		NaverSecret:    cfg.OAuth.Naver.ClientSecret,
	}, lessonsearchprerun.NewReportRepository(pool), jobs.NewPostgresAdvisoryLocker(pool, 0), log.Default())

	if err := scheduler.RunOnce(ctx); err != nil {
		log.Fatal("lesson search prerun run once failed: ", err)
	}

	reports, err := lessonsearchprerun.NewReportRepository(pool).ListRecentReports(ctx, 1)
	if err != nil {
		log.Fatal("failed to load latest lesson search prerun report: ", err)
	}
	if len(reports) == 0 {
		log.Fatal("lesson search prerun completed but no report was stored")
	}
	latest := reports[0]
	if err := lessonsearchprerun.EnsureNoSecretMarkers(latest.Summary, latest.NoiseTaxonomySummary, latest.ErrorCode, latest.ErrorMessage); err != nil {
		log.Fatal("latest lesson search prerun report failed safety check: ", err)
	}

	fmt.Printf("lesson search prerun report stored id=%s status=%s provider=%s top_n=%d risk_top_n=%d created_at=%s\n",
		latest.ID.String(), latest.Status, latest.ProviderFilter, latest.TopN, latest.RiskTopN, latest.CreatedAt.Format(time.RFC3339))
}

func splitFixtureIDs(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t'
	})
	out := make([]string, 0, len(fields))
	seen := map[string]bool{}
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" || seen[field] {
			continue
		}
		seen[field] = true
		out = append(out, field)
	}
	return out
}
