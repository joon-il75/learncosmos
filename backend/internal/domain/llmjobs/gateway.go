package llmjobs

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

type Gateway struct {
	repo  *Repository
	queue Queue
}

func NewGateway(repo *Repository, queue Queue) *Gateway {
	return &Gateway{repo: repo, queue: queue}
}

func (g *Gateway) Submit(ctx context.Context, input CreateJobInput) (*Job, error) {
	if g == nil || g.repo == nil {
		return nil, ErrWorkerQueueDisabled
	}
	job, err := g.repo.CreateJob(ctx, input)
	if err != nil {
		return nil, err
	}
	if g.queue == nil {
		return job, ErrWorkerQueueDisabled
	}
	if err := g.queue.Enqueue(ctx, job.ID); err != nil {
		_, _ = g.repo.MarkFailed(ctx, job.ID, FailJobInput{
			ErrorCode:    "llm_job_queue_failed",
			ErrorMessage: logsafe.PersistedError(err.Error()),
			Retryable:    true,
		})
		return nil, err
	}
	return job, nil
}

func (g *Gateway) SubmitAndWait(ctx context.Context, input CreateJobInput, wait time.Duration) (*Job, error) {
	job, err := g.Submit(ctx, input)
	if err != nil {
		return nil, err
	}
	if wait <= 0 {
		return job, nil
	}
	deadline := time.NewTimer(wait)
	defer deadline.Stop()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return job, ctx.Err()
		case <-deadline.C:
			return job, nil
		case <-ticker.C:
			current, err := g.repo.GetJob(ctx, job.ID)
			if err != nil {
				return job, err
			}
			if terminalStatus(current.Status) {
				return current, nil
			}
		}
	}
}

func (g *Gateway) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	return g.repo.GetJob(ctx, id)
}

func (g *Gateway) GetActiveJobByIdempotencyKey(ctx context.Context, feature, idempotencyKey string) (*Job, error) {
	if g == nil || g.repo == nil {
		return nil, ErrWorkerQueueDisabled
	}
	return g.repo.GetActiveJobByIdempotencyKey(ctx, feature, idempotencyKey)
}

func (g *Gateway) CountActiveJobsByFeature(ctx context.Context, feature string) (int, error) {
	if g == nil || g.repo == nil {
		return 0, ErrWorkerQueueDisabled
	}
	return g.repo.CountActiveJobsByFeature(ctx, feature)
}

func (g *Gateway) GetActiveJobEstimate(ctx context.Context, feature string, jobID uuid.UUID) (ActiveJobEstimate, error) {
	if g == nil || g.repo == nil {
		return ActiveJobEstimate{}, ErrWorkerQueueDisabled
	}
	return g.repo.GetActiveJobEstimate(ctx, feature, jobID)
}
