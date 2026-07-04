package llmjobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateJob(ctx context.Context, input CreateJobInput) (*Job, error) {
	feature := normalizeFeature(input.Feature)
	if feature == "" {
		return nil, ErrInvalidJobInput
	}

	requestRef, err := marshalRef(input.RequestRef)
	if err != nil {
		return nil, err
	}
	promptInputRef, err := marshalRef(input.PromptInputRef)
	if err != nil {
		return nil, err
	}

	if input.MaxActiveJobs > 0 {
		return r.createJobWithActiveLimit(ctx, input, feature, requestRef, promptInputRef)
	}
	return r.insertJob(ctx, r.pool, input, feature, requestRef, promptInputRef)
}

type jobInserter interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *Repository) createJobWithActiveLimit(ctx context.Context, input CreateJobInput, feature string, requestRef, promptInputRef []byte) (*Job, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, feature); err != nil {
		return nil, err
	}
	var activeCount int
	if err := tx.QueryRow(ctx, `
		SELECT count(*)::int
		FROM llm_jobs
		WHERE feature = $1
		  AND status IN ($2, $3)
	`, feature, StatusQueued, StatusRunning).Scan(&activeCount); err != nil {
		return nil, err
	}
	if activeCount >= input.MaxActiveJobs {
		return nil, ErrActiveJobLimitExceeded
	}

	job, err := r.insertJob(ctx, tx, input, feature, requestRef, promptInputRef)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return job, nil
}

func (r *Repository) insertJob(ctx context.Context, q jobInserter, input CreateJobInput, feature string, requestRef, promptInputRef []byte) (*Job, error) {
	id := uuid.New()
	row := q.QueryRow(ctx, `
		INSERT INTO llm_jobs (
			id, user_id, feature, idempotency_key, status, priority,
			request_ref, prompt_template_key, prompt_input_ref,
			provider, model, provider_mode, api_key_ref,
			point_cost, billing_status, queued_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7::jsonb, $8, $9::jsonb,
			$10, $11, $12, $13,
			$14, $15, now()
		)
		RETURNING `+jobColumns,
		id, input.UserID, feature, input.IdempotencyKey, StatusQueued, input.Priority,
		requestRef, input.PromptTemplateKey, promptInputRef,
		input.Provider, input.Model, input.ProviderMode, input.APIKeyRef,
		input.PointCost, input.BillingStatus,
	)
	job, err := scanJob(row)
	if isDuplicateJobError(err) {
		return nil, ErrDuplicateJob
	}
	return job, err
}

func (r *Repository) GetActiveJobByIdempotencyKey(ctx context.Context, feature, idempotencyKey string) (*Job, error) {
	feature = normalizeFeature(feature)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if feature == "" || idempotencyKey == "" {
		return nil, ErrInvalidJobInput
	}
	row := r.pool.QueryRow(ctx, `
		SELECT `+jobColumns+`
		FROM llm_jobs
		WHERE feature = $1
		  AND idempotency_key = $2
		  AND status IN ($3, $4, $5)
		ORDER BY created_at DESC
		LIMIT 1
	`, feature, idempotencyKey, StatusQueued, StatusRunning, StatusSucceeded)
	return scanJob(row)
}

func (r *Repository) CountActiveJobsByFeature(ctx context.Context, feature string) (int, error) {
	feature = normalizeFeature(feature)
	if feature == "" {
		return 0, ErrInvalidJobInput
	}
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT count(*)::int
		FROM llm_jobs
		WHERE feature = $1
		  AND status IN ($2, $3)
	`, feature, StatusQueued, StatusRunning).Scan(&count)
	return count, err
}

func (r *Repository) GetActiveJobEstimate(ctx context.Context, feature string, jobID uuid.UUID) (ActiveJobEstimate, error) {
	feature = normalizeFeature(feature)
	if feature == "" {
		return ActiveJobEstimate{}, ErrInvalidJobInput
	}

	var estimate ActiveJobEstimate
	var status Status
	var createdAt time.Time
	if err := r.pool.QueryRow(ctx, `
		SELECT status, created_at
		FROM llm_jobs
		WHERE id = $1 AND feature = $2
	`, jobID, feature).Scan(&status, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ActiveJobEstimate{}, ErrJobNotFound
		}
		return ActiveJobEstimate{}, err
	}

	if err := r.pool.QueryRow(ctx, `
		SELECT count(*)::int
		FROM llm_jobs
		WHERE feature = $1
		  AND status IN ($2, $3)
	`, feature, StatusQueued, StatusRunning).Scan(&estimate.ActiveCount); err != nil {
		return ActiveJobEstimate{}, err
	}

	switch status {
	case StatusQueued:
		var position int
		if err := r.pool.QueryRow(ctx, `
			SELECT count(*)::int
			FROM llm_jobs
			WHERE feature = $1
			  AND status IN ($2, $3)
			  AND (created_at < $4 OR (created_at = $4 AND id::text <= $5))
		`, feature, StatusQueued, StatusRunning, createdAt, jobID.String()).Scan(&position); err != nil {
			return ActiveJobEstimate{}, err
		}
		estimate.QueuePosition = &position
	case StatusRunning:
		position := 0
		estimate.QueuePosition = &position
	}

	return estimate, nil
}

func isDuplicateJobError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "idx_llm_jobs_idempotency_key"
}

func (r *Repository) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+jobColumns+` FROM llm_jobs WHERE id = $1`, id)
	return scanJob(row)
}

func (r *Repository) GetJobForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Job, error) {
	job, err := r.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}
	if job.UserID == nil || *job.UserID != userID {
		return nil, ErrJobAccessDenied
	}
	return job, nil
}

func (r *Repository) MarkRunning(ctx context.Context, id uuid.UUID, workerID string, providerInFlight, providerQueueDepth *int) (*Job, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE llm_jobs
		SET status = $2,
		    started_at = COALESCE(started_at, now()),
		    updated_at = now(),
		    queue_wait_ms = CASE
		      WHEN queued_at IS NOT NULL THEN GREATEST(0, floor(extract(epoch from (now() - queued_at)) * 1000)::int)
		      ELSE queue_wait_ms
		    END,
		    worker_id = $3,
		    attempts = attempts + 1,
		    provider_inflight_count = COALESCE($4, provider_inflight_count),
		    provider_queue_depth = COALESCE($5, provider_queue_depth)
		WHERE id = $1 AND status = $6
		RETURNING `+jobColumns,
		id, StatusRunning, workerID, providerInFlight, providerQueueDepth, StatusQueued,
	)
	return scanJob(row)
}

func (r *Repository) MarkSucceeded(ctx context.Context, id uuid.UUID, input CompleteJobInput) (*Job, error) {
	resultRef, err := marshalRef(input.ResultRef)
	if err != nil {
		return nil, err
	}
	row := r.pool.QueryRow(ctx, `
		UPDATE llm_jobs
		SET status = $2,
		    phase = NULL,
		    result_ref = $3::jsonb,
		    error_code = NULL,
		    error_message = NULL,
		    retryable = false,
		    provider_wait_ms = COALESCE($4, provider_wait_ms),
		    llm_generation_elapsed_ms = COALESCE($5, llm_generation_elapsed_ms),
		    parse_validate_elapsed_ms = COALESCE($6, parse_validate_elapsed_ms),
		    persist_elapsed_ms = COALESCE($7, persist_elapsed_ms),
		    total_job_elapsed_ms = COALESCE(
		      $8,
		      CASE WHEN started_at IS NOT NULL THEN GREATEST(0, floor(extract(epoch from (now() - started_at)) * 1000)::int) ELSE total_job_elapsed_ms END
		    ),
		    finished_at = now(),
		    updated_at = now()
		WHERE id = $1 AND status = $9
		RETURNING `+jobColumns,
		id, StatusSucceeded, resultRef,
		input.ProviderWaitMS, input.LLMGenerationElapsedMS, input.ParseValidateElapsedMS, input.PersistElapsedMS, input.TotalJobElapsedMS,
		StatusRunning,
	)
	return scanJob(row)
}

func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, input FailJobInput) (*Job, error) {
	errorCode := strings.TrimSpace(input.ErrorCode)
	if errorCode == "" {
		errorCode = "llm_generation_failed"
	}
	errorMessage := logsafe.PersistedError(input.ErrorMessage)
	row := r.pool.QueryRow(ctx, `
		UPDATE llm_jobs
		SET status = $2,
		    phase = NULL,
		    error_code = $3,
		    error_message = NULLIF($4, ''),
		    retryable = $5,
		    provider_wait_ms = COALESCE($6, provider_wait_ms),
		    llm_generation_elapsed_ms = COALESCE($7, llm_generation_elapsed_ms),
		    parse_validate_elapsed_ms = COALESCE($8, parse_validate_elapsed_ms),
		    persist_elapsed_ms = COALESCE($9, persist_elapsed_ms),
		    total_job_elapsed_ms = COALESCE(
		      $10,
		      CASE WHEN started_at IS NOT NULL THEN GREATEST(0, floor(extract(epoch from (now() - started_at)) * 1000)::int) ELSE total_job_elapsed_ms END
		    ),
		    finished_at = now(),
		    updated_at = now()
		WHERE id = $1 AND status IN ($11, $12)
		RETURNING `+jobColumns,
		id, StatusFailed, errorCode, errorMessage, input.Retryable,
		input.ProviderWaitMS, input.LLMGenerationElapsedMS, input.ParseValidateElapsedMS, input.PersistElapsedMS, input.TotalJobElapsedMS,
		StatusQueued, StatusRunning,
	)
	return scanJob(row)
}

func (r *Repository) SetPhase(ctx context.Context, id uuid.UUID, phase string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE llm_jobs
		SET phase = NULLIF($2, ''), updated_at = now()
		WHERE id = $1 AND status = $3
	`, id, strings.TrimSpace(phase), StatusRunning)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInvalidTransition
	}
	return nil
}

func marshalRef(value map[string]any) ([]byte, error) {
	if value == nil {
		value = map[string]any{}
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid json ref", ErrInvalidJobInput)
	}
	return payload, nil
}

const jobColumns = `
	id, user_id, feature, idempotency_key, status, priority, phase,
	request_ref, prompt_template_key, prompt_input_ref,
	provider, model, provider_mode, api_key_ref, result_ref,
	error_code, error_message, retryable, point_cost, billing_status,
	queue_wait_ms, provider_wait_ms, llm_generation_elapsed_ms,
	parse_validate_elapsed_ms, persist_elapsed_ms, total_job_elapsed_ms,
	provider_inflight_count, provider_queue_depth, worker_id, attempts,
	created_at, queued_at, started_at, finished_at, updated_at
`

func scanJob(row pgx.Row) (*Job, error) {
	var job Job
	err := row.Scan(
		&job.ID, &job.UserID, &job.Feature, &job.IdempotencyKey, &job.Status, &job.Priority, &job.Phase,
		&job.RequestRef, &job.PromptTemplateKey, &job.PromptInputRef,
		&job.Provider, &job.Model, &job.ProviderMode, &job.APIKeyRef, &job.ResultRef,
		&job.ErrorCode, &job.ErrorMessage, &job.Retryable, &job.PointCost, &job.BillingStatus,
		&job.QueueWaitMS, &job.ProviderWaitMS, &job.LLMGenerationElapsedMS,
		&job.ParseValidateElapsedMS, &job.PersistElapsedMS, &job.TotalJobElapsedMS,
		&job.ProviderInFlightCount, &job.ProviderQueueDepth, &job.WorkerID, &job.Attempts,
		&job.CreatedAt, &job.QueuedAt, &job.StartedAt, &job.FinishedAt, &job.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrJobNotFound
	}
	if job.RequestRef == nil {
		job.RequestRef = json.RawMessage(`{}`)
	}
	if job.PromptInputRef == nil {
		job.PromptInputRef = json.RawMessage(`{}`)
	}
	if job.ResultRef == nil {
		job.ResultRef = json.RawMessage(`{}`)
	}
	return &job, err
}

func elapsedMS(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
