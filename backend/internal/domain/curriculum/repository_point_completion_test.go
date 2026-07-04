package curriculum

import "testing"

func TestLearningPointCompletionReady(t *testing.T) {
	t.Parallel()

	answered := "현재 목표에 맞게 적용할 수 있다"
	score := 3
	completeSelfEvaluation := &CoursePointSelfEvaluation{
		UnderstandingScore:   &score,
		UnderstandingReason:  "핵심 개념을 내 말로 설명할 수 있다.",
		ApplicationScore:     &score,
		ApplicationReason:    "다음 학습 작업에 바로 적용할 수 있다.",
		ProficiencyScore:     &score,
		ProficiencyReason:    "반복 연습이 더 필요하지만 흐름은 이해했다.",
		ProblemSolvingScore:  &score,
		ProblemSolvingReason: "막힌 부분을 확인했고 다음 해결 방법을 적었다.",
		ExpressionScore:      &score,
		ExpressionReason:     "정리한 내용을 다른 사람에게 설명할 수 있다.",
		GoalAlignmentNote:    "현재 목표 달성의 핵심 기초다.",
	}
	tests := []struct {
		name           string
		questions      []CoursePointQuestion
		selfEvaluation *CoursePointSelfEvaluation
		want           bool
	}{
		{
			name:           "missing answered question",
			questions:      []CoursePointQuestion{{Question: "무엇을 배웠나?", Status: PointQuestionStatusPending}},
			selfEvaluation: completeSelfEvaluation,
			want:           false,
		},
		{
			name: "missing self evaluation",
			questions: []CoursePointQuestion{{
				Question: "무엇을 배웠나?",
				Answer:   &answered,
				Status:   PointQuestionStatusAnswered,
			}},
			selfEvaluation: nil,
			want:           false,
		},
		{
			name: "missing goal connection note",
			questions: []CoursePointQuestion{{
				Question: "무엇을 배웠나?",
				Answer:   &answered,
				Status:   PointQuestionStatusAnswered,
			}},
			selfEvaluation: &CoursePointSelfEvaluation{
				UnderstandingScore:   &score,
				UnderstandingReason:  "핵심 개념을 이해했다.",
				ApplicationScore:     &score,
				ApplicationReason:    "적용할 수 있다.",
				ProficiencyScore:     &score,
				ProficiencyReason:    "반복할 수 있다.",
				ProblemSolvingScore:  &score,
				ProblemSolvingReason: "막힌 부분을 풀었다.",
				ExpressionScore:      &score,
				ExpressionReason:     "설명할 수 있다.",
			},
			want: false,
		},
		{
			name: "legacy application note alone does not satisfy completion",
			questions: []CoursePointQuestion{{
				Question: "어디에 적용할까?",
				Answer:   &answered,
				Status:   PointQuestionStatusAnswered,
			}},
			selfEvaluation: &CoursePointSelfEvaluation{
				Understanding:   4,
				Proficiency:     4,
				ApplicationNote: "다음 연습 세션에 바로 써본다",
			},
			want: false,
		},
		{
			name: "goal alignment note alone does not satisfy completion",
			questions: []CoursePointQuestion{{
				Question: "목표와 어떻게 연결되나?",
				Answer:   &answered,
				Status:   PointQuestionStatusAnswered,
			}},
			selfEvaluation: &CoursePointSelfEvaluation{
				Understanding:     4,
				Proficiency:       4,
				GoalAlignmentNote: "현재 목표 달성의 핵심 기초다",
			},
			want: false,
		},
		{
			name: "all six self evaluation items satisfy completion",
			questions: []CoursePointQuestion{{
				Question: "목표와 어떻게 연결되나?",
				Answer:   &answered,
				Status:   PointQuestionStatusAnswered,
			}},
			selfEvaluation: completeSelfEvaluation,
			want:           true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := learningPointCompletionReady(tt.questions, tt.selfEvaluation); got != tt.want {
				t.Fatalf("learningPointCompletionReady() = %v, want %v", got, tt.want)
			}
		})
	}
}
