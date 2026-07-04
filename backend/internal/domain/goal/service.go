package goal

type Service struct{}

func NewService() *Service { return &Service{} }

var genericInterviewFallbackMessage = "계속 말씀해주세요. 더 잘 도와드릴 수 있을 것 같아요."

func (s *Service) RefreshLearningIntentProfile(profile *GoalProfile) {
	if profile == nil {
		return
	}
	profile.LearningIntent = BuildLearningIntentProfileWithSource(profile, LearningIntentSourceGoalChat)
}

func BeginGoalRevision(profile *GoalProfile, message string) {
	if profile == nil {
		return
	}
	if profile.InterviewState == StateConfirmed || profile.InterviewState == StateAwaitingRebuildDecision {
		if profile.RevisionSnapshot == nil {
			profile.RevisionSnapshot = snapshotGoalRevisionState(profile)
		}
		profile.Version++
		profile.RebuildDecision = nil
	}
	profile.InterviewState = StateRevisingGoal
	profile.Messages = append(profile.Messages, InterviewMessage{Role: "user", Content: message})
}

func snapshotGoalRevisionState(profile *GoalProfile) *GoalRevisionState {
	if profile == nil {
		return nil
	}
	messages := make([]InterviewMessage, len(profile.Messages))
	copy(messages, profile.Messages)
	return &GoalRevisionState{
		UserIntent:        profile.UserIntent,
		Motivation:        profile.Motivation,
		UsageContext:      profile.UsageContext,
		ConfirmedGoal:     profile.ConfirmedGoal,
		GoalType:          profile.GoalType,
		OutputType:        profile.OutputType,
		DifficultyLevel:   profile.DifficultyLevel,
		TimeHorizon:       profile.TimeHorizon,
		SummarizedContext: profile.SummarizedContext,
		LearningIntent:    NormalizeLearningIntentProfile(profile.LearningIntent),
		InterviewState:    profile.InterviewState,
		Messages:          messages,
		Version:           profile.Version,
	}
}
