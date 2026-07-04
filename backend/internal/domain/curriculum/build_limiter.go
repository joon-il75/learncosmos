package curriculum

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
)

var errDraftBuildBusy = errors.New("draft_build_busy")

type draftBuildLimiter struct {
	slots               chan struct{}
	inFlight            atomic.Int32
	waiting             atomic.Int32
	skipReviewThreshold int32
	maxQueueDepth       int32
}

func newDraftBuildLimiter(maxConcurrency int, skipReviewRatio float64, maxQueueDepth int) *draftBuildLimiter {
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}
	if skipReviewRatio <= 0 || skipReviewRatio > 1 {
		skipReviewRatio = 0.75
	}
	if maxQueueDepth < 0 {
		maxQueueDepth = 0
	}

	threshold := int32(math.Ceil(float64(maxConcurrency) * skipReviewRatio))
	if threshold < 1 {
		threshold = 1
	}

	return &draftBuildLimiter{
		slots:               make(chan struct{}, maxConcurrency),
		skipReviewThreshold: threshold,
		maxQueueDepth:       int32(maxQueueDepth),
	}
}

func (l *draftBuildLimiter) Acquire(ctx context.Context) error {
	if l == nil {
		return nil
	}
	if l.maxQueueDepth >= 0 && l.inFlight.Load() >= int32(cap(l.slots)) && l.waiting.Load() >= l.maxQueueDepth {
		return fmt.Errorf("%w: queue saturated", errDraftBuildBusy)
	}

	l.waiting.Add(1)
	defer l.waiting.Add(-1)

	select {
	case l.slots <- struct{}{}:
		l.inFlight.Add(1)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: %v", errDraftBuildBusy, ctx.Err())
	}
}

func (l *draftBuildLimiter) Release() {
	select {
	case <-l.slots:
		l.inFlight.Add(-1)
	default:
	}
}

func (l *draftBuildLimiter) ShouldSkipReview() bool {
	return l.inFlight.Load() >= l.skipReviewThreshold
}

type draftBuildLimiterSnapshot struct {
	InFlight   int32
	QueueDepth int32
	Capacity   int
}

func (l *draftBuildLimiter) Snapshot() draftBuildLimiterSnapshot {
	if l == nil {
		return draftBuildLimiterSnapshot{}
	}
	return draftBuildLimiterSnapshot{
		InFlight:   l.inFlight.Load(),
		QueueDepth: l.waiting.Load(),
		Capacity:   cap(l.slots),
	}
}
