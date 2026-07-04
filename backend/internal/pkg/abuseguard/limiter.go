package abuseguard

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit  int
	window time.Duration
	now    func() time.Time
	mu     sync.Mutex
	items  map[string]fixedWindowItem
}

type fixedWindowItem struct {
	count      int
	windowEnds time.Time
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &FixedWindowLimiter{
		limit:  limit,
		window: window,
		now:    time.Now,
		items:  make(map[string]fixedWindowItem),
	}
}

func (l *FixedWindowLimiter) Allow(key string) (bool, time.Duration) {
	if key == "" {
		key = "unknown"
	}
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	item := l.items[key]
	if item.windowEnds.IsZero() || !now.Before(item.windowEnds) {
		l.items[key] = fixedWindowItem{count: 1, windowEnds: now.Add(l.window)}
		l.sweepExpiredLocked(now)
		return true, 0
	}

	if item.count >= l.limit {
		retryAfter := item.windowEnds.Sub(now)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	item.count++
	l.items[key] = item
	return true, 0
}

func (l *FixedWindowLimiter) sweepExpiredLocked(now time.Time) {
	if len(l.items) < 1024 {
		return
	}
	for key, item := range l.items {
		if !now.Before(item.windowEnds) {
			delete(l.items, key)
		}
	}
}
