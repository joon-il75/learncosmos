package curriculum

import (
	"testing"

	"github.com/google/uuid"
	goaldomain "github.com/learnweaver/backend/internal/domain/goal"
)

func TestBuildCreateDraftRequestFromGoalCarriesLearningIntent(t *testing.T) {
	confirmed := "가죽공예로 손 지갑을 완성한다"
	activeGoal := &goaldomain.GoalProfile{
		ID:             uuid.New(),
		ConfirmedGoal:  &confirmed,
		Language:       "ko",
		Version:        3,
		LearningIntent: goaldomain.LearningIntentProfile{LearnerLevel: "beginner", DesiredOutput: "손 지갑", PreferredActivities: []string{"프로젝트 실습"}, Source: goaldomain.LearningIntentSourceGoalChat},
	}

	req := buildCreateDraftRequestFromGoal(activeGoal, "ko")

	if req.LearningIntent.DesiredOutput != "손 지갑" {
		t.Fatalf("desired output = %q, want 손 지갑", req.LearningIntent.DesiredOutput)
	}
	if len(req.LearningIntent.PreferredActivities) != 1 || req.LearningIntent.PreferredActivities[0] != "프로젝트 실습" {
		t.Fatalf("preferred activities = %#v", req.LearningIntent.PreferredActivities)
	}
}
