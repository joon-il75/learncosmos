package llmjobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

type Processor interface {
	Process(ctx context.Context, job *Job) (CompleteJobInput, error)
}

type ProcessorFunc func(ctx context.Context, job *Job) (CompleteJobInput, error)

func (f ProcessorFunc) Process(ctx context.Context, job *Job) (CompleteJobInput, error) {
	return f(ctx, job)
}

type Runtime struct {
	repo        runtimeRepository
	queue       Queue
	processors  map[string]Processor
	workerCount int
	logger      *log.Logger
}

type runtimeRepository interface {
	MarkRunning(ctx context.Context, id uuid.UUID, workerID string, providerInFlight, providerQueueDepth *int) (*Job, error)
	MarkSucceeded(ctx context.Context, id uuid.UUID, input CompleteJobInput) (*Job, error)
	MarkFailed(ctx context.Context, id uuid.UUID, input FailJobInput) (*Job, error)
}

type RuntimeOptions struct {
	WorkerCount int
	Logger      *log.Logger
}

func NewRuntime(repo *Repository, queue Queue, opts RuntimeOptions) *Runtime {
	workerCount := opts.WorkerCount
	if workerCount <= 0 {
		workerCount = 6
	}
	logger := opts.Logger
	if logger == nil {
		logger = log.Default()
	}
	return &Runtime{
		repo:        repo,
		queue:       queue,
		processors:  map[string]Processor{},
		workerCount: workerCount,
		logger:      logger,
	}
}

func (r *Runtime) RegisterProcessor(feature string, processor Processor) {
	if r == nil || processor == nil {
		return
	}
	feature = normalizeFeature(feature)
	if feature == "" {
		return
	}
	r.processors[feature] = processor
}

func (r *Runtime) Start(ctx context.Context) {
	if r == nil || r.repo == nil || r.queue == nil {
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < r.workerCount; i++ {
		workerID := fmt.Sprintf("llm-worker-%02d", i+1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.runWorker(ctx, workerID)
		}()
	}
	go func() {
		<-ctx.Done()
		wg.Wait()
	}()
	r.logger.Printf("[llm-worker] started worker_count=%d", r.workerCount)
}

func (r *Runtime) runWorker(ctx context.Context, workerID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, err := r.queue.Read(ctx, workerID, 5*time.Second)
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			r.logger.Printf("[llm-worker] read error worker_id=%s err=%v", workerID, err)
			time.Sleep(time.Second)
			continue
		}
		if msg == nil {
			continue
		}
		r.handleMessage(ctx, workerID, *msg)
	}
}

func (r *Runtime) handleMessage(ctx context.Context, workerID string, msg QueueMessage) {
	start := time.Now()
	job, err := r.repo.MarkRunning(ctx, msg.JobID, workerID, nil, nil)
	if err != nil {
		r.logger.Printf("[llm-worker] mark running failed worker_id=%s job_id=%s err=%v", workerID, msg.JobID, err)
		_ = r.queue.Ack(ctx, msg.MessageID)
		return
	}

	processor := r.processors[job.Feature]
	if processor == nil {
		_, _ = r.repo.MarkFailed(ctx, job.ID, FailJobInput{
			ErrorCode:    "llm_job_unsupported_feature",
			ErrorMessage: ErrUnsupportedFeature.Error(),
			Retryable:    false,
		})
		_ = r.queue.MoveToDLQ(ctx, msg, "unsupported_feature")
		_ = r.queue.Ack(ctx, msg.MessageID)
		return
	}

	result, err := processor.Process(ctx, job)
	total := elapsedMS(start)
	result.TotalJobElapsedMS = &total
	if err != nil {
		errorCode := "llm_generation_failed"
		retryable := false
		var processorErr *ProcessorError
		if errors.As(err, &processorErr) {
			if processorErr.Code != "" {
				errorCode = processorErr.Code
			}
			retryable = processorErr.Retryable
		}
		_, _ = r.repo.MarkFailed(ctx, job.ID, FailJobInput{
			ErrorCode:         errorCode,
			ErrorMessage:      logsafe.PersistedError(err.Error()),
			Retryable:         retryable,
			TotalJobElapsedMS: &total,
		})
		_ = r.queue.Ack(ctx, msg.MessageID)
		r.logger.Printf("[llm-worker] job failed worker_id=%s job_id=%s feature=%s error_code=%s retryable=%t err=%s", workerID, job.ID, job.Feature, errorCode, retryable, logsafe.Error(err))
		return
	}
	if _, err := r.repo.MarkSucceeded(ctx, job.ID, result); err != nil {
		r.logger.Printf("[llm-worker] mark succeeded failed worker_id=%s job_id=%s err=%v", workerID, job.ID, err)
		_ = r.queue.MoveToDLQ(ctx, msg, "mark_succeeded_failed")
		_ = r.queue.Ack(ctx, msg.MessageID)
		return
	}
	_ = r.queue.Ack(ctx, msg.MessageID)
	r.logger.Printf("[llm-worker] job succeeded worker_id=%s job_id=%s feature=%s duration_ms=%d", workerID, job.ID, job.Feature, total)
}

func NewDummyProcessor() Processor {
	return ProcessorFunc(func(ctx context.Context, job *Job) (CompleteJobInput, error) {
		if job == nil {
			return CompleteJobInput{}, ErrJobNotFound
		}
		if err := ctx.Err(); err != nil {
			return CompleteJobInput{}, err
		}
		return CompleteJobInput{
			ResultRef: map[string]any{
				"dummy":  true,
				"job_id": job.ID.String(),
			},
		}, nil
	})
}
