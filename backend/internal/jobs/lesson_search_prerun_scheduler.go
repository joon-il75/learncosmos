package jobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

const lessonSearchPrerunAdvisoryLockKey int64 = 62026180618

type LessonSearchPrerunSchedulerConfig struct {
	Enabled        bool
	Interval       time.Duration
	StartDelay     time.Duration
	Timeout        time.Duration
	Provider       string
	TopN           int
	RiskTopN       int
	RiskFixtureIDs []string
	Retention      time.Duration
	FixturesPath   string
	RunType        string
	YouTubeAPIKey  string
	NaverClientID  string
	NaverSecret    string
}

type LessonSearchPrerunReportStore interface {
	SaveReport(ctx context.Context, input lessonsearchprerun.SaveReportInput) (*lessonsearchprerun.StoredReport, error)
	DeleteReportsCreatedBefore(ctx context.Context, cutoff time.Time) (int64, error)
}

type LessonSearchPrerunLocker interface {
	TryLock(ctx context.Context) (bool, func(context.Context) error, error)
}

type LessonSearchPrerunScheduler struct {
	cfg       LessonSearchPrerunSchedulerConfig
	reports   LessonSearchPrerunReportStore
	locker    LessonSearchPrerunLocker
	providers map[string]lessonsearchprerun.Provider
	logger    *log.Logger
	now       func() time.Time
}

type lessonSearchPrerunPersistedReport struct {
	Mode         string                         `json:"mode"`
	MainReport   lessonsearchprerun.Report      `json:"main_report"`
	RiskReports  []lessonsearchprerun.Report    `json:"risk_reports,omitempty"`
	QuotaSummary lessonSearchPrerunQuotaSummary `json:"quota_summary"`
}

type lessonSearchPrerunQuotaSummary struct {
	ProviderFilter            string   `json:"provider_filter"`
	TopN                      int      `json:"top_n"`
	RiskTopN                  int      `json:"risk_top_n"`
	RiskFixtureIDs            []string `json:"risk_fixture_ids,omitempty"`
	MainProviderRequests      int      `json:"main_provider_requests"`
	RiskProviderRequests      int      `json:"risk_provider_requests"`
	EstimatedProviderRequests int      `json:"estimated_provider_requests"`
	EstimatedCandidatesMax    int      `json:"estimated_candidates_max"`
}

type PostgresAdvisoryLocker struct {
	pool *pgxpool.Pool
	key  int64
}

func NewPostgresAdvisoryLocker(pool *pgxpool.Pool, key int64) *PostgresAdvisoryLocker {
	if key == 0 {
		key = lessonSearchPrerunAdvisoryLockKey
	}
	return &PostgresAdvisoryLocker{pool: pool, key: key}
}

func (locker *PostgresAdvisoryLocker) TryLock(ctx context.Context) (bool, func(context.Context) error, error) {
	if locker == nil || locker.pool == nil {
		return false, nil, errors.New("lesson search prerun advisory locker is not configured")
	}
	tx, err := locker.pool.Begin(ctx)
	if err != nil {
		return false, nil, fmt.Errorf("begin lesson search prerun lock transaction: %w", err)
	}
	var acquired bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, locker.key).Scan(&acquired); err != nil {
		_ = tx.Rollback(ctx)
		return false, nil, fmt.Errorf("acquire lesson search prerun advisory lock: %w", err)
	}
	if !acquired {
		_ = tx.Rollback(ctx)
		return false, nil, nil
	}
	return true, func(unlockCtx context.Context) error {
		if err := tx.Rollback(unlockCtx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return fmt.Errorf("release lesson search prerun advisory lock: %w", err)
		}
		return nil
	}, nil
}

func NewLessonSearchPrerunScheduler(cfg LessonSearchPrerunSchedulerConfig, reports LessonSearchPrerunReportStore, locker LessonSearchPrerunLocker, logger *log.Logger) *LessonSearchPrerunScheduler {
	cfg = normalizeLessonSearchPrerunSchedulerConfig(cfg)
	return &LessonSearchPrerunScheduler{
		cfg:       cfg,
		reports:   reports,
		locker:    locker,
		providers: lessonsearchprerun.BuildProvidersWithCredentials(cfg.YouTubeAPIKey, cfg.NaverClientID, cfg.NaverSecret),
		logger:    logger,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func (scheduler *LessonSearchPrerunScheduler) Start(ctx context.Context) {
	if scheduler == nil || !scheduler.cfg.Enabled {
		return
	}
	go scheduler.loop(ctx)
}

func (scheduler *LessonSearchPrerunScheduler) loop(ctx context.Context) {
	if scheduler.cfg.StartDelay > 0 {
		timer := time.NewTimer(scheduler.cfg.StartDelay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
		}
	}
	scheduler.runAndLog(ctx)

	ticker := time.NewTicker(scheduler.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			scheduler.runAndLog(ctx)
		}
	}
}

func (scheduler *LessonSearchPrerunScheduler) runAndLog(ctx context.Context) {
	if err := scheduler.RunOnce(ctx); err != nil && scheduler.logger != nil {
		scheduler.logger.Printf("[lesson-search-prerun] run failed err=%s", logsafe.Error(err))
	}
}

func (scheduler *LessonSearchPrerunScheduler) RunOnce(ctx context.Context) error {
	if scheduler == nil {
		return errors.New("lesson search prerun scheduler is nil")
	}
	if scheduler.reports == nil {
		return errors.New("lesson search prerun report store is not configured")
	}
	if scheduler.locker == nil {
		return errors.New("lesson search prerun locker is not configured")
	}
	startedAt := scheduler.now()
	acquired, unlock, err := scheduler.locker.TryLock(ctx)
	if err != nil {
		finishedAt := scheduler.now()
		_ = scheduler.saveReport(ctx, startedAt, &finishedAt, lessonsearchprerun.ReportStatusFailed, map[string]any{"lock_error": true}, map[string]any{}, map[string]any{"mode": "lock"}, ptrString("lock_failed"), ptrString("lock acquisition failed"))
		return err
	}
	if !acquired {
		finishedAt := scheduler.now()
		return scheduler.saveReport(ctx, startedAt, &finishedAt, lessonsearchprerun.ReportStatusSkippedLockNotAcquired, map[string]any{"skipped": "lock_not_acquired"}, map[string]any{}, map[string]any{"mode": "lock"}, nil, nil)
	}
	defer func() {
		if unlock != nil {
			if err := unlock(context.Background()); err != nil && scheduler.logger != nil {
				scheduler.logger.Printf("[lesson-search-prerun] unlock failed err=%s", logsafe.Error(err))
			}
		}
	}()

	if missing := scheduler.missingCredentialProviders(); len(missing) > 0 {
		finishedAt := scheduler.now()
		return scheduler.saveReport(ctx, startedAt, &finishedAt, lessonsearchprerun.ReportStatusSkippedMissingCredentials, map[string]any{"missing_credentials": missing}, map[string]any{}, map[string]any{"mode": lessonsearchprerun.DefaultMode, "provider_filter": scheduler.cfg.Provider}, nil, nil)
	}

	fixtures, err := lessonsearchprerun.LoadFixtures(scheduler.cfg.FixturesPath)
	if err != nil {
		finishedAt := scheduler.now()
		return scheduler.saveReport(ctx, startedAt, &finishedAt, lessonsearchprerun.ReportStatusFailed, map[string]any{"fixture_load_failed": true}, map[string]any{}, map[string]any{"mode": lessonsearchprerun.DefaultMode}, ptrString("fixture_load_failed"), ptrString("failed to load lesson search prerun fixtures"))
	}

	runCtx, cancel := context.WithTimeout(ctx, scheduler.cfg.Timeout)
	defer cancel()
	mainReport := lessonsearchprerun.BuildReport(runCtx, scheduler.runOptions("", scheduler.cfg.TopN), fixtures, scheduler.providers, scheduler.now())
	riskReports := scheduler.buildRiskReports(runCtx, fixtures)
	quotaSummary := scheduler.quotaSummary(mainReport, riskReports)
	persisted := lessonSearchPrerunPersistedReport{
		Mode:         lessonsearchprerun.DefaultMode,
		MainReport:   mainReport,
		RiskReports:  riskReports,
		QuotaSummary: quotaSummary,
	}
	finishedAt := scheduler.now()
	status := lessonsearchprerun.ReportStatusSucceeded
	var errorCode *string
	var errorMessage *string
	if mainReport.Summary.ProviderErrors+sumRiskProviderErrors(riskReports) > 0 {
		status = lessonsearchprerun.ReportStatusFailed
		errorCode = ptrString("provider_error")
		errorMessage = ptrString("one or more providers returned an error")
	}
	if err := scheduler.saveReport(ctx, startedAt, &finishedAt, status, scheduler.persistedSummary(mainReport, riskReports, quotaSummary), combinedNoiseTaxonomy(mainReport, riskReports), persisted, errorCode, errorMessage); err != nil {
		return err
	}
	if scheduler.cfg.Retention > 0 {
		_, err := scheduler.reports.DeleteReportsCreatedBefore(ctx, scheduler.now().Add(-scheduler.cfg.Retention))
		if err != nil {
			return err
		}
	}
	return nil
}

func (scheduler *LessonSearchPrerunScheduler) buildRiskReports(ctx context.Context, fixtures []lessonsearchprerun.Fixture) []lessonsearchprerun.Report {
	if scheduler.cfg.RiskTopN <= 0 || len(scheduler.cfg.RiskFixtureIDs) == 0 {
		return nil
	}
	reports := []lessonsearchprerun.Report{}
	for _, fixtureID := range scheduler.cfg.RiskFixtureIDs {
		fixtureID = strings.TrimSpace(fixtureID)
		if fixtureID == "" {
			continue
		}
		report := lessonsearchprerun.BuildReport(ctx, scheduler.runOptions(fixtureID, scheduler.cfg.RiskTopN), fixtures, scheduler.providers, scheduler.now())
		if report.Summary.FixturesProcessed == 0 {
			continue
		}
		reports = append(reports, report)
	}
	return reports
}

func (scheduler *LessonSearchPrerunScheduler) saveReport(ctx context.Context, startedAt time.Time, finishedAt *time.Time, status string, summary any, noise any, report any, errorCode *string, errorMessage *string) error {
	_, err := scheduler.reports.SaveReport(ctx, lessonsearchprerun.SaveReportInput{
		RunType:              scheduler.cfg.RunType,
		Status:               status,
		StartedAt:            startedAt,
		FinishedAt:           finishedAt,
		ProviderFilter:       scheduler.cfg.Provider,
		TopN:                 scheduler.cfg.TopN,
		RiskTopN:             scheduler.cfg.RiskTopN,
		Summary:              summary,
		NoiseTaxonomySummary: noise,
		Report:               report,
		ErrorCode:            errorCode,
		ErrorMessage:         errorMessage,
	})
	return err
}

func (scheduler *LessonSearchPrerunScheduler) runOptions(fixtureID string, topN int) lessonsearchprerun.Options {
	return lessonsearchprerun.Options{
		FixturesPath: scheduler.cfg.FixturesPath,
		Provider:     scheduler.cfg.Provider,
		FixtureID:    fixtureID,
		TopN:         topN,
		Mode:         lessonsearchprerun.DefaultMode,
		Timeout:      scheduler.cfg.Timeout,
	}
}

func (scheduler *LessonSearchPrerunScheduler) persistedSummary(mainReport lessonsearchprerun.Report, riskReports []lessonsearchprerun.Report, quota lessonSearchPrerunQuotaSummary) map[string]any {
	return map[string]any{
		"fixtures_total":               mainReport.Summary.FixturesTotal,
		"fixtures_processed":           mainReport.Summary.FixturesProcessed,
		"fixtures_skipped":             mainReport.Summary.FixturesSkipped,
		"provider_errors":              mainReport.Summary.ProviderErrors + sumRiskProviderErrors(riskReports),
		"noise_suspects":               mainReport.Summary.NoiseSuspects + sumRiskNoiseSuspects(riskReports),
		"relevant_like_candidates":     mainReport.Summary.RelevantLikeCandidates,
		"risk_reports":                 len(riskReports),
		"risk_provider_errors":         sumRiskProviderErrors(riskReports),
		"risk_noise_suspects":          sumRiskNoiseSuspects(riskReports),
		"quota":                        quota,
		"estimated_provider_requests":  quota.EstimatedProviderRequests,
		"estimated_candidates_max":     quota.EstimatedCandidatesMax,
		"main_provider_requests":       quota.MainProviderRequests,
		"risk_provider_requests":       quota.RiskProviderRequests,
		"risk_fixture_ids_configured":  append([]string(nil), scheduler.cfg.RiskFixtureIDs...),
		"risk_fixture_ids_processed":   processedRiskFixtureIDs(riskReports),
		"automatic_correction_applied": false,
	}
}

func (scheduler *LessonSearchPrerunScheduler) quotaSummary(mainReport lessonsearchprerun.Report, riskReports []lessonsearchprerun.Report) lessonSearchPrerunQuotaSummary {
	mainRequests := mainReport.Summary.FixturesProcessed
	riskRequests := 0
	for _, report := range riskReports {
		riskRequests += report.Summary.FixturesProcessed
	}
	return lessonSearchPrerunQuotaSummary{
		ProviderFilter:            scheduler.cfg.Provider,
		TopN:                      scheduler.cfg.TopN,
		RiskTopN:                  scheduler.cfg.RiskTopN,
		RiskFixtureIDs:            append([]string(nil), scheduler.cfg.RiskFixtureIDs...),
		MainProviderRequests:      mainRequests,
		RiskProviderRequests:      riskRequests,
		EstimatedProviderRequests: mainRequests + riskRequests,
		EstimatedCandidatesMax:    mainRequests*scheduler.cfg.TopN + riskRequests*scheduler.cfg.RiskTopN,
	}
}

func (scheduler *LessonSearchPrerunScheduler) missingCredentialProviders() []string {
	provider := strings.ToLower(strings.TrimSpace(scheduler.cfg.Provider))
	missing := []string{}
	if provider == "all" || provider == "youtube" {
		if strings.TrimSpace(scheduler.cfg.YouTubeAPIKey) == "" {
			missing = append(missing, "youtube")
		}
	}
	if provider == "all" || provider == "naver_blog" {
		if strings.TrimSpace(scheduler.cfg.NaverClientID) == "" || strings.TrimSpace(scheduler.cfg.NaverSecret) == "" {
			missing = append(missing, "naver_blog")
		}
	}
	return missing
}

func normalizeLessonSearchPrerunSchedulerConfig(cfg LessonSearchPrerunSchedulerConfig) LessonSearchPrerunSchedulerConfig {
	if cfg.Interval <= 0 {
		cfg.Interval = 24 * time.Hour
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	cfg.Provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	if cfg.Provider == "" {
		cfg.Provider = "all"
	}
	switch cfg.Provider {
	case "all", "youtube", "naver_blog":
	default:
		cfg.Provider = "all"
	}
	if cfg.TopN <= 0 {
		cfg.TopN = 5
	}
	if cfg.RiskTopN < 0 {
		cfg.RiskTopN = 0
	}
	cfg.RiskFixtureIDs = normalizeRiskFixtureIDs(cfg.RiskFixtureIDs)
	if cfg.FixturesPath == "" {
		cfg.FixturesPath = lessonsearchprerun.DefaultFixturePath
	}
	if strings.TrimSpace(cfg.RunType) == "" {
		cfg.RunType = "daily"
	}
	return cfg
}

func normalizeRiskFixtureIDs(values []string) []string {
	seen := map[string]struct{}{}
	result := []string{}
	for _, value := range values {
		for _, part := range strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' }) {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			if _, ok := seen[trimmed]; ok {
				continue
			}
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func combinedNoiseTaxonomy(mainReport lessonsearchprerun.Report, riskReports []lessonsearchprerun.Report) map[string]any {
	registered := map[string]int{}
	unregistered := map[string]int{}
	addNoiseCounts(registered, mainReport.NoiseTaxonomySummary.RegisteredAvoidHits)
	addNoiseCounts(unregistered, mainReport.NoiseTaxonomySummary.UnregisteredNoiseHints)
	for _, report := range riskReports {
		addNoiseCounts(registered, report.NoiseTaxonomySummary.RegisteredAvoidHits)
		addNoiseCounts(unregistered, report.NoiseTaxonomySummary.UnregisteredNoiseHints)
	}
	return map[string]any{
		"registered_avoid_hits":    tokenCountMap(registered),
		"unregistered_noise_hints": tokenCountMap(unregistered),
	}
}

func addNoiseCounts(target map[string]int, values []lessonsearchprerun.NoiseTokenCount) {
	for _, value := range values {
		if strings.TrimSpace(value.Token) == "" || value.Count <= 0 {
			continue
		}
		target[value.Token] += value.Count
	}
}

func tokenCountMap(values map[string]int) []map[string]any {
	result := []map[string]any{}
	for token, count := range values {
		result = append(result, map[string]any{"token": token, "count": count})
	}
	return result
}

func sumRiskProviderErrors(reports []lessonsearchprerun.Report) int {
	total := 0
	for _, report := range reports {
		total += report.Summary.ProviderErrors
	}
	return total
}

func sumRiskNoiseSuspects(reports []lessonsearchprerun.Report) int {
	total := 0
	for _, report := range reports {
		total += report.Summary.NoiseSuspects
	}
	return total
}

func processedRiskFixtureIDs(reports []lessonsearchprerun.Report) []string {
	ids := []string{}
	for _, report := range reports {
		for _, fixture := range report.Fixtures {
			if fixture.Skipped == "" || len(fixture.Candidates) > 0 || fixture.Error != "" {
				ids = append(ids, fixture.ID)
			}
		}
	}
	return ids
}

func ptrString(value string) *string {
	return &value
}
