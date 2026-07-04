package curriculum

import (
	"sync"

	"github.com/google/uuid"
)

type courseGenerationGoalLock struct {
	mu       sync.Mutex
	inFlight map[string]struct{}
}

func newCourseGenerationGoalLock() *courseGenerationGoalLock {
	return &courseGenerationGoalLock{inFlight: make(map[string]struct{})}
}

func (l *courseGenerationGoalLock) TryAcquire(userID, goalID uuid.UUID) (func(), bool) {
	if l == nil || userID == uuid.Nil || goalID == uuid.Nil {
		return func() {}, true
	}
	key := userID.String() + ":" + goalID.String()
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.inFlight[key]; exists {
		return nil, false
	}
	l.inFlight[key] = struct{}{}
	return func() {
		l.mu.Lock()
		delete(l.inFlight, key)
		l.mu.Unlock()
	}, true
}

func formatOptionalUUID(id *uuid.UUID) string {
	if id == nil {
		return "none"
	}
	return id.String()
}
