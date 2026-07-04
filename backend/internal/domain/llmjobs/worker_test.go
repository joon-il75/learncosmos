package llmjobs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeRuntimeRepo struct {
	jobs map[uuid.UUID]*Job
}

func newFakeRuntimeRepo(job *Job) *fakeRuntimeRepo {
	return &fakeRuntimeRepo{jobs: map[uuid.UUID]*Job{job.ID: job}}
}

func (r *fakeRuntimeRepo) MarkRunning(ctx context.Context, id uuid.UUID, workerID string, providerInFlight, providerQueueDepth *int) (*Job, error) {
	job := r.jobs[id]
	if job == nil {
		return nil, ErrJobNotFound
	}
	if job.Status != StatusQueued {
		return nil, ErrInvalidTransition
	}
	job.Status = StatusRunning
	job.WorkerID = &workerID
	job.Attempts++
	queueWait := 12
	job.QueueWaitMS = &queueWait
	now := time.Now()
	job.StartedAt = &now
	job.UpdatedAt = now
	return job, nil
}

func (r *fakeRuntimeRepo) MarkSucceeded(ctx context.Context, id uuid.UUID, input CompleteJobInput) (*Job, error) {
	job := r.jobs[id]
	if job == nil {
		return nil, ErrJobNotFound
	}
	if job.Status != StatusRunning {
		return nil, ErrInvalidTransition
	}
	payload, err := json.Marshal(input.ResultRef)
	if err != nil {
		return nil, err
	}
	job.Status = StatusSucceeded
	job.ResultRef = payload
	job.ProviderWaitMS = input.ProviderWaitMS
	job.LLMGenerationElapsedMS = input.LLMGenerationElapsedMS
	job.ParseValidateElapsedMS = input.ParseValidateElapsedMS
	job.PersistElapsedMS = input.PersistElapsedMS
	job.TotalJobElapsedMS = input.TotalJobElapsedMS
	now := time.Now()
	job.FinishedAt = &now
	job.UpdatedAt = now
	return job, nil
}

func (r *fakeRuntimeRepo) MarkFailed(ctx context.Context, id uuid.UUID, input FailJobInput) (*Job, error) {
	job := r.jobs[id]
	if job == nil {
		return nil, ErrJobNotFound
	}
	if job.Status != StatusQueued && job.Status != StatusRunning {
		return nil, ErrInvalidTransition
	}
	job.Status = StatusFailed
	code := input.ErrorCode
	message := input.ErrorMessage
	job.ErrorCode = &code
	job.ErrorMessage = &message
	job.Retryable = input.Retryable
	job.ProviderWaitMS = input.ProviderWaitMS
	job.LLMGenerationElapsedMS = input.LLMGenerationElapsedMS
	job.ParseValidateElapsedMS = input.ParseValidateElapsedMS
	job.PersistElapsedMS = input.PersistElapsedMS
	job.TotalJobElapsedMS = input.TotalJobElapsedMS
	now := time.Now()
	job.FinishedAt = &now
	job.UpdatedAt = now
	return job, nil
}

type fakeRuntimeQueue struct {
	acked []string
	dlq   []string
}

func (q *fakeRuntimeQueue) Enqueue(ctx context.Context, jobID uuid.UUID) error { return nil }
func (q *fakeRuntimeQueue) Read(ctx context.Context, consumer string, block time.Duration) (*QueueMessage, error) {
	return nil, nil
}
func (q *fakeRuntimeQueue) Ack(ctx context.Context, messageID string) error {
	q.acked = append(q.acked, messageID)
	return nil
}
func (q *fakeRuntimeQueue) MoveToDLQ(ctx context.Context, message QueueMessage, reason string) error {
	q.dlq = append(q.dlq, reason)
	return nil
}

func TestRuntimeHandleMessageRecordsSuccessMetricsAndAck(t *testing.T) {
	job := newRuntimeTestJob(FeatureDummy)
	repo := newFakeRuntimeRepo(job)
	queue := &fakeRuntimeQueue{}
	runtime := newTestRuntime(repo, queue)
	providerWait := 7
	generation := 13
	parse := 3
	persist := 5
	runtime.RegisterProcessor(FeatureDummy, ProcessorFunc(func(ctx context.Context, job *Job) (CompleteJobInput, error) {
		return CompleteJobInput{
			ResultRef:              map[string]any{"ok": true},
			ProviderWaitMS:         &providerWait,
			LLMGenerationElapsedMS: &generation,
			ParseValidateElapsedMS: &parse,
			PersistElapsedMS:       &persist,
		}, nil
	}))

	runtime.handleMessage(context.Background(), "worker-test", QueueMessage{MessageID: "1-0", JobID: job.ID})

	if job.Status != StatusSucceeded {
		t.Fatalf("status = %s, want %s", job.Status, StatusSucceeded)
	}
	if job.WorkerID == nil || *job.WorkerID != "worker-test" {
		t.Fatalf("worker_id = %v, want worker-test", job.WorkerID)
	}
	if job.Attempts != 1 {
		t.Fatalf("attempts = %d, want 1", job.Attempts)
	}
	assertIntPtr(t, "queue_wait_ms", job.QueueWaitMS, 12)
	assertIntPtr(t, "provider_wait_ms", job.ProviderWaitMS, providerWait)
	assertIntPtr(t, "llm_generation_elapsed_ms", job.LLMGenerationElapsedMS, generation)
	assertIntPtr(t, "parse_validate_elapsed_ms", job.ParseValidateElapsedMS, parse)
	assertIntPtr(t, "persist_elapsed_ms", job.PersistElapsedMS, persist)
	if job.TotalJobElapsedMS == nil || *job.TotalJobElapsedMS < 0 {
		t.Fatalf("total_job_elapsed_ms = %v, want non-negative", job.TotalJobElapsedMS)
	}
	var result map[string]bool
	if err := json.Unmarshal(job.ResultRef, &result); err != nil {
		t.Fatalf("result_ref unmarshal: %v", err)
	}
	if !result["ok"] {
		t.Fatalf("result_ref = %s, want ok=true", string(job.ResultRef))
	}
	assertQueueAcked(t, queue, "1-0")
	if len(queue.dlq) != 0 {
		t.Fatalf("dlq = %v, want empty", queue.dlq)
	}
}

func TestRuntimeHandleMessageRecordsProcessorFailureAndAck(t *testing.T) {
	job := newRuntimeTestJob(FeaturePointAISummary)
	repo := newFakeRuntimeRepo(job)
	queue := &fakeRuntimeQueue{}
	runtime := newTestRuntime(repo, queue)
	runtime.RegisterProcessor(FeaturePointAISummary, ProcessorFunc(func(ctx context.Context, job *Job) (CompleteJobInput, error) {
		return CompleteJobInput{}, NewProcessorError("point_ai_summary_failed", "summary failed", true)
	}))

	runtime.handleMessage(context.Background(), "worker-test", QueueMessage{MessageID: "2-0", JobID: job.ID})

	if job.Status != StatusFailed {
		t.Fatalf("status = %s, want %s", job.Status, StatusFailed)
	}
	if job.ErrorCode == nil || *job.ErrorCode != "point_ai_summary_failed" {
		t.Fatalf("error_code = %v, want point_ai_summary_failed", job.ErrorCode)
	}
	if !job.Retryable {
		t.Fatal("retryable = false, want true")
	}
	if job.TotalJobElapsedMS == nil || *job.TotalJobElapsedMS < 0 {
		t.Fatalf("total_job_elapsed_ms = %v, want non-negative", job.TotalJobElapsedMS)
	}
	assertQueueAcked(t, queue, "2-0")
	if len(queue.dlq) != 0 {
		t.Fatalf("dlq = %v, want empty", queue.dlq)
	}
}

func TestRuntimeHandleMessageMovesUnsupportedFeatureToDLQ(t *testing.T) {
	job := newRuntimeTestJob("missing_feature")
	repo := newFakeRuntimeRepo(job)
	queue := &fakeRuntimeQueue{}
	runtime := newTestRuntime(repo, queue)

	runtime.handleMessage(context.Background(), "worker-test", QueueMessage{MessageID: "3-0", JobID: job.ID})

	if job.Status != StatusFailed {
		t.Fatalf("status = %s, want %s", job.Status, StatusFailed)
	}
	if job.ErrorCode == nil || *job.ErrorCode != "llm_job_unsupported_feature" {
		t.Fatalf("error_code = %v, want llm_job_unsupported_feature", job.ErrorCode)
	}
	if job.Retryable {
		t.Fatal("retryable = true, want false")
	}
	assertQueueAcked(t, queue, "3-0")
	if len(queue.dlq) != 1 || queue.dlq[0] != "unsupported_feature" {
		t.Fatalf("dlq = %v, want unsupported_feature", queue.dlq)
	}
}

func newRuntimeTestJob(feature string) *Job {
	now := time.Now()
	return &Job{
		ID:             uuid.New(),
		Feature:        feature,
		Status:         StatusQueued,
		RequestRef:     json.RawMessage(`{}`),
		PromptInputRef: json.RawMessage(`{}`),
		ResultRef:      json.RawMessage(`{}`),
		CreatedAt:      now,
		QueuedAt:       &now,
		UpdatedAt:      now,
	}
}

func newTestRuntime(repo runtimeRepository, queue Queue) *Runtime {
	return &Runtime{
		repo:        repo,
		queue:       queue,
		processors:  map[string]Processor{},
		workerCount: 1,
		logger:      log.New(io.Discard, "", 0),
	}
}

func assertIntPtr(t *testing.T, name string, got *int, want int) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %d", name, got, want)
	}
}

func assertQueueAcked(t *testing.T, queue *fakeRuntimeQueue, messageID string) {
	t.Helper()
	if len(queue.acked) != 1 || queue.acked[0] != messageID {
		t.Fatalf("acked = %v, want [%s]", queue.acked, messageID)
	}
}

func TestRuntimeHandleMessageSanitizesPersistedProcessorFailure(t *testing.T) {
	job := newRuntimeTestJob(FeaturePointAISummary)
	repo := newFakeRuntimeRepo(job)
	queue := &fakeRuntimeQueue{}
	runtime := newTestRuntime(repo, queue)
	secretMessage := "provider failed https://api.example.com/v1/chat?api_key=secret-token user prompt=프라이빗 목표"
	runtime.RegisterProcessor(FeaturePointAISummary, ProcessorFunc(func(ctx context.Context, job *Job) (CompleteJobInput, error) {
		return CompleteJobInput{}, NewProcessorError("point_ai_summary_failed", secretMessage, true)
	}))

	runtime.handleMessage(context.Background(), "worker-test", QueueMessage{MessageID: "4-0", JobID: job.ID})

	if job.ErrorMessage == nil || *job.ErrorMessage == "" {
		t.Fatal("expected sanitized error_message")
	}
	for _, leaked := range []string{"secret-token", "api_key", "프라이빗", "user prompt", "api.example.com"} {
		if strings.Contains(*job.ErrorMessage, leaked) {
			t.Fatalf("error_message leaked %q in %q", leaked, *job.ErrorMessage)
		}
	}
	if !strings.HasPrefix(*job.ErrorMessage, "persisted_error ") {
		t.Fatalf("error_message = %q, want persisted_error summary", *job.ErrorMessage)
	}
}
