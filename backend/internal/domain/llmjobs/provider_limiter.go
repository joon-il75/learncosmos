package llmjobs

import (
	"context"
	"sync"
	"time"
)

type ProviderLimiter struct {
	sem chan struct{}
	mu  sync.Mutex
}

type ProviderLimiterSnapshot struct {
	WaitMS        int
	InFlightCount int
	QueueDepth    int
}

func NewProviderLimiter(maxConcurrency int) *ProviderLimiter {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}
	return &ProviderLimiter{sem: make(chan struct{}, maxConcurrency)}
}

func (l *ProviderLimiter) Acquire(ctx context.Context) (func(), ProviderLimiterSnapshot, error) {
	start := time.Now()
	queueDepth := len(l.sem)
	select {
	case l.sem <- struct{}{}:
		snapshot := ProviderLimiterSnapshot{
			WaitMS:        int(time.Since(start).Milliseconds()),
			InFlightCount: len(l.sem),
			QueueDepth:    queueDepth,
		}
		return func() {
			<-l.sem
		}, snapshot, nil
	case <-ctx.Done():
		return nil, ProviderLimiterSnapshot{
			WaitMS:        int(time.Since(start).Milliseconds()),
			InFlightCount: len(l.sem),
			QueueDepth:    queueDepth,
		}, ctx.Err()
	}
}
