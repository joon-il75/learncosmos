package llmjobs

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultStreamKey     = "llm_jobs:v1"
	DefaultConsumerGroup = "llm-workers"
	DefaultDLQKey        = "llm_jobs:dead:v1"
)

type Queue interface {
	Enqueue(ctx context.Context, jobID uuid.UUID) error
	Read(ctx context.Context, consumer string, block time.Duration) (*QueueMessage, error)
	Ack(ctx context.Context, messageID string) error
	MoveToDLQ(ctx context.Context, message QueueMessage, reason string) error
}

type QueueMessage struct {
	MessageID string
	JobID     uuid.UUID
}

type RedisStreamQueue struct {
	client *redis.Client
	stream string
	group  string
	dlq    string
}

func NewRedisStreamQueue(client *redis.Client, stream, group, dlq string) *RedisStreamQueue {
	if stream == "" {
		stream = DefaultStreamKey
	}
	if group == "" {
		group = DefaultConsumerGroup
	}
	if dlq == "" {
		dlq = DefaultDLQKey
	}
	return &RedisStreamQueue{client: client, stream: stream, group: group, dlq: dlq}
}

func (q *RedisStreamQueue) EnsureGroup(ctx context.Context) error {
	if q == nil || q.client == nil {
		return ErrWorkerQueueDisabled
	}
	err := q.client.XGroupCreateMkStream(ctx, q.stream, q.group, "0").Err()
	if err != nil && !isBusyGroupErr(err) {
		return err
	}
	return nil
}

func (q *RedisStreamQueue) Enqueue(ctx context.Context, jobID uuid.UUID) error {
	if q == nil || q.client == nil {
		return ErrWorkerQueueDisabled
	}
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]any{"job_id": jobID.String()},
	}).Err()
}

func (q *RedisStreamQueue) Read(ctx context.Context, consumer string, block time.Duration) (*QueueMessage, error) {
	if q == nil || q.client == nil {
		return nil, ErrWorkerQueueDisabled
	}
	streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    q.group,
		Consumer: consumer,
		Streams:  []string{q.stream, ">"},
		Count:    1,
		Block:    block,
	}).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return nil, nil
	}
	msg := streams[0].Messages[0]
	raw, _ := msg.Values["job_id"].(string)
	jobID, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}
	return &QueueMessage{MessageID: msg.ID, JobID: jobID}, nil
}

func (q *RedisStreamQueue) Ack(ctx context.Context, messageID string) error {
	if q == nil || q.client == nil {
		return ErrWorkerQueueDisabled
	}
	return q.client.XAck(ctx, q.stream, q.group, messageID).Err()
}

func (q *RedisStreamQueue) MoveToDLQ(ctx context.Context, message QueueMessage, reason string) error {
	if q == nil || q.client == nil {
		return ErrWorkerQueueDisabled
	}
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.dlq,
		Values: map[string]any{
			"job_id":     message.JobID.String(),
			"message_id": message.MessageID,
			"reason":     reason,
		},
	}).Err()
}

func isBusyGroupErr(err error) bool {
	return err != nil && containsRedisBusyGroup(err.Error())
}

func containsRedisBusyGroup(value string) bool {
	return value == "BUSYGROUP Consumer Group name already exists" ||
		len(value) >= len("BUSYGROUP") && value[:len("BUSYGROUP")] == "BUSYGROUP"
}
