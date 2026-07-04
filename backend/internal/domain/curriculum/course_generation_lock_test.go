package curriculum

import (
	"testing"

	"github.com/google/uuid"
)

func TestCourseGenerationGoalLockBlocksDuplicateGoal(t *testing.T) {
	lock := newCourseGenerationGoalLock()
	userID := uuid.New()
	goalID := uuid.New()

	release, ok := lock.TryAcquire(userID, goalID)
	if !ok {
		t.Fatal("first acquire failed")
	}
	if _, ok := lock.TryAcquire(userID, goalID); ok {
		t.Fatal("duplicate acquire for same user and goal should be blocked")
	}
	release()
	if releaseAgain, ok := lock.TryAcquire(userID, goalID); !ok {
		t.Fatal("acquire after release failed")
	} else {
		releaseAgain()
	}
}

func TestCourseGenerationGoalLockAllowsDifferentGoalsAndUsers(t *testing.T) {
	lock := newCourseGenerationGoalLock()
	userID := uuid.New()
	goalID := uuid.New()

	release, ok := lock.TryAcquire(userID, goalID)
	if !ok {
		t.Fatal("first acquire failed")
	}
	defer release()

	if releaseOtherGoal, ok := lock.TryAcquire(userID, uuid.New()); !ok {
		t.Fatal("different goal for same user should not share lock")
	} else {
		releaseOtherGoal()
	}
	if releaseOtherUser, ok := lock.TryAcquire(uuid.New(), goalID); !ok {
		t.Fatal("same goal id for different user should not share lock")
	} else {
		releaseOtherUser()
	}
}
