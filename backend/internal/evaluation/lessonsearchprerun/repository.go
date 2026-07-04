package lessonsearchprerun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ReportStatusRunning                   = "running"
	ReportStatusSucceeded                 = "succeeded"
	ReportStatusFailed                    = "failed"
	ReportStatusSkippedMissingCredentials = "skipped_missing_credentials"
	ReportStatusSkippedLockNotAcquired    = "skipped_lock_not_acquired"
)

var reportSecretMarkerPattern = regexp.MustCompile(`(?i)(YOUTUBE_API_KEY|NAVER_CLIENT_ID|NAVER_CLIENT_SECRET|OPENAI_API_KEY|ANTHROPIC_API_KEY|Authorization|Bearer|sk-[A-Za-z0-9])`)

type ReportRepository struct {
	pool *pgxpool.Pool
}

type SaveReportInput struct {
	RunType              string
	Status               string
	StartedAt            time.Time
	FinishedAt           *time.Time
	ProviderFilter       string
	TopN                 int
	RiskTopN             int
	Summary              any
	NoiseTaxonomySummary any
	Report               any
	ErrorCode            *string
	ErrorMessage         *string
}

type StoredReport struct {
	ID                   uuid.UUID       `json:"id"`
	RunType              string          `json:"run_type"`
	Status               string          `json:"status"`
	StartedAt            time.Time       `json:"started_at"`
	FinishedAt           *time.Time      `json:"finished_at,omitempty"`
	ProviderFilter       string          `json:"provider_filter"`
	TopN                 int             `json:"top_n"`
	RiskTopN             int             `json:"risk_top_n"`
	Summary              json.RawMessage `json:"summary"`
	NoiseTaxonomySummary json.RawMessage `json:"noise_taxonomy_summary"`
	Report               json.RawMessage `json:"report"`
	ErrorCode            *string         `json:"error_code,omitempty"`
	ErrorMessage         *string         `json:"error_message,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
}

func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

func (repo *ReportRepository) SaveReport(ctx context.Context, input SaveReportInput) (*StoredReport, error) {
	if repo == nil || repo.pool == nil {
		return nil, errors.New("lesson search prerun report repository is not configured")
	}
	normalized, err := normalizeSaveReportInput(input)
	if err != nil {
		return nil, err
	}
	summaryJSON, err := marshalReportJSON(normalized.Summary, "summary")
	if err != nil {
		return nil, err
	}
	noiseJSON, err := marshalReportJSON(normalized.NoiseTaxonomySummary, "noise taxonomy summary")
	if err != nil {
		return nil, err
	}
	reportJSON, err := marshalReportJSON(normalized.Report, "report")
	if err != nil {
		return nil, err
	}
	if err := EnsureNoSecretMarkers(summaryJSON, noiseJSON, reportJSON, normalized.ErrorCode, normalized.ErrorMessage); err != nil {
		return nil, err
	}

	row := repo.pool.QueryRow(ctx, `
		INSERT INTO lesson_search_prerun_reports (
			run_type, status, started_at, finished_at, provider_filter, top_n, risk_top_n,
			summary, noise_taxonomy_summary, report, error_code, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9::jsonb, $10::jsonb, $11, $12)
		RETURNING id, run_type, status, started_at, finished_at, provider_filter, top_n, risk_top_n,
		          summary, noise_taxonomy_summary, report, error_code, error_message, created_at
	`,
		normalized.RunType,
		normalized.Status,
		normalized.StartedAt,
		normalized.FinishedAt,
		normalized.ProviderFilter,
		normalized.TopN,
		normalized.RiskTopN,
		summaryJSON,
		noiseJSON,
		reportJSON,
		normalized.ErrorCode,
		normalized.ErrorMessage,
	)
	stored, err := scanStoredReport(row)
	if err != nil {
		return nil, fmt.Errorf("save lesson search prerun report: %w", err)
	}
	return stored, nil
}

func (repo *ReportRepository) ListRecentReports(ctx context.Context, limit int) ([]StoredReport, error) {
	if repo == nil || repo.pool == nil {
		return nil, errors.New("lesson search prerun report repository is not configured")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := repo.pool.Query(ctx, `
		SELECT id, run_type, status, started_at, finished_at, provider_filter, top_n, risk_top_n,
		       summary, noise_taxonomy_summary, report, error_code, error_message, created_at
		FROM lesson_search_prerun_reports
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list lesson search prerun reports: %w", err)
	}
	defer rows.Close()

	reports := []StoredReport{}
	for rows.Next() {
		report, err := scanStoredReport(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lesson search prerun reports: %w", err)
	}
	return reports, nil
}

func (repo *ReportRepository) DeleteReportsCreatedBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	if repo == nil || repo.pool == nil {
		return 0, errors.New("lesson search prerun report repository is not configured")
	}
	tag, err := repo.pool.Exec(ctx, `DELETE FROM lesson_search_prerun_reports WHERE created_at < $1`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete old lesson search prerun reports: %w", err)
	}
	return tag.RowsAffected(), nil
}

func EnsureNoSecretMarkers(values ...any) error {
	for _, value := range values {
		if value == nil {
			continue
		}
		var text string
		switch typed := value.(type) {
		case []byte:
			text = string(typed)
		case string:
			text = typed
		case *string:
			if typed == nil {
				continue
			}
			text = *typed
		default:
			raw, err := json.Marshal(typed)
			if err != nil {
				return fmt.Errorf("scan report secret markers: %w", err)
			}
			text = string(raw)
		}
		if reportSecretMarkerPattern.MatchString(text) {
			return errors.New("lesson search prerun report contains secret marker")
		}
	}
	return nil
}

func normalizeSaveReportInput(input SaveReportInput) (SaveReportInput, error) {
	input.RunType = strings.TrimSpace(input.RunType)
	if input.RunType == "" {
		input.RunType = "daily"
	}
	input.Status = strings.TrimSpace(input.Status)
	if input.Status == "" {
		input.Status = ReportStatusSucceeded
	}
	switch input.Status {
	case ReportStatusRunning, ReportStatusSucceeded, ReportStatusFailed, ReportStatusSkippedMissingCredentials, ReportStatusSkippedLockNotAcquired:
	default:
		return SaveReportInput{}, fmt.Errorf("unsupported lesson search prerun report status %q", input.Status)
	}
	if input.StartedAt.IsZero() {
		input.StartedAt = time.Now().UTC()
	}
	input.ProviderFilter = strings.ToLower(strings.TrimSpace(input.ProviderFilter))
	if input.ProviderFilter == "" {
		input.ProviderFilter = "all"
	}
	if input.TopN <= 0 {
		input.TopN = 5
	}
	if input.RiskTopN < 0 {
		return SaveReportInput{}, errors.New("risk_top_n must be greater than or equal to 0")
	}
	if input.Summary == nil {
		input.Summary = map[string]any{}
	}
	if input.NoiseTaxonomySummary == nil {
		input.NoiseTaxonomySummary = map[string]any{}
	}
	if input.Report == nil {
		input.Report = map[string]any{}
	}
	return input, nil
}

func marshalReportJSON(value any, label string) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal lesson search prerun %s: %w", label, err)
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, fmt.Errorf("lesson search prerun %s must be a JSON object: %w", label, err)
	}
	return raw, nil
}

type reportScanner interface {
	Scan(dest ...any) error
}

func scanStoredReport(row reportScanner) (*StoredReport, error) {
	var report StoredReport
	if err := row.Scan(
		&report.ID,
		&report.RunType,
		&report.Status,
		&report.StartedAt,
		&report.FinishedAt,
		&report.ProviderFilter,
		&report.TopN,
		&report.RiskTopN,
		&report.Summary,
		&report.NoiseTaxonomySummary,
		&report.Report,
		&report.ErrorCode,
		&report.ErrorMessage,
		&report.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("scan lesson search prerun report postgres error %s: %w", pgErr.Code, err)
		}
		return nil, fmt.Errorf("scan lesson search prerun report: %w", err)
	}
	return &report, nil
}
