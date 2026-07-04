package curriculum

import "testing"

func TestParsePointAIQuestionSuggestion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		raw       string
		pointType PointType
		wantQ     string
		wantType  string
	}{
		{
			name:      "json response",
			raw:       `{"question":"이 지식을 이번 목표에 맞게 어디에 바로 써볼 수 있을까요?","question_type":"application"}`,
			pointType: PointTypeExploration,
			wantQ:     "이 지식을 이번 목표에 맞게 어디에 바로 써볼 수 있을까요?",
			wantType:  "application",
		},
		{
			name:      "plain text fallback",
			raw:       "이 지점의 내용을 현재 목표와 어떻게 연결할까요?",
			pointType: PointTypeResearch,
			wantQ:     "이 지점의 내용을 현재 목표와 어떻게 연결할까요?",
			wantType:  "goal_alignment",
		},
		{
			name:      "empty fallback",
			raw:       "",
			pointType: PointTypeExploration,
			wantQ:     "이 지점의 내용을 현재 목표에 맞게 다음 행동으로 어떻게 옮길 수 있을까요?",
			wantType:  "application",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parsePointAIQuestionSuggestion(tt.raw, tt.pointType, "ko")
			if got.Question != tt.wantQ {
				t.Fatalf("question = %q, want %q", got.Question, tt.wantQ)
			}
			if got.QuestionType != tt.wantType {
				t.Fatalf("questionType = %q, want %q", got.QuestionType, tt.wantType)
			}
		})
	}
}

func TestParsePointAIQuestionSuggestionEnglishFallback(t *testing.T) {
	t.Parallel()
	got := parsePointAIQuestionSuggestion("", PointTypeExploration, "en")
	if got.Question == "" || got.Question == "이 지점의 내용을 현재 목표에 맞게 다음 행동으로 어떻게 옮길 수 있을까요?" {
		t.Fatalf("expected english fallback question, got %q", got.Question)
	}
	if got.QuestionType != "application" {
		t.Fatalf("questionType = %q, want application", got.QuestionType)
	}
}
