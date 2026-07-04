package curriculum

import (
	"strings"
	"testing"
)

func TestNormalizePointAISummaryText(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "plain",
			raw:  "핵심을 짧게 정리합니다.",
			want: "핵심을 짧게 정리합니다.",
		},
		{
			name: "braced text",
			raw:  "{핵심을 짧게 정리합니다.}",
			want: "핵심을 짧게 정리합니다.",
		},
		{
			name: "json summary",
			raw:  `{"summary":"핵심을 짧게 정리합니다."}`,
			want: "핵심을 짧게 정리합니다.",
		},
		{
			name: "korean json summary",
			raw:  `{"요약":"일렉기타 장비를 이해하고 소리를 실험해보세요.","질문":"앰프와 오디오 인터페이스의 차이는 무엇인가요?","안내문구":"직접 장비를 사용해 보세요."}`,
			want: "일렉기타 장비를 이해하고 소리를 실험해보세요.",
		},
		{
			name: "nested korean json summary",
			raw:  `{"학습요약":{"요약":"카드지갑 만들기는 재단, 바느질, 엣지코트 마감을 연결해 완성하는 과정입니다.","질문":"엣지코트는 언제 바르나요?"}}`,
			want: "카드지갑 만들기는 재단, 바느질, 엣지코트 마감을 연결해 완성하는 과정입니다.",
		},
		{
			name: "json fenced",
			raw:  "```json\n{\"text\":\"핵심을 짧게 정리합니다.\"}\n```",
			want: "핵심을 짧게 정리합니다.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePointAISummaryText(tt.raw); got != tt.want {
				t.Fatalf("normalizePointAISummaryText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPointAIPromptsUseLearningLanguage(t *testing.T) {
	summaryPrompt := buildPointAISummaryPrompt("Node", "Practice", "Source", "Description", "en")
	if !strings.Contains(summaryPrompt, "Respond in English") {
		t.Fatalf("expected english summary prompt instruction, got %q", summaryPrompt)
	}

	questionPrompt := buildPointAIQuestionPrompt(nil, &CoursePointAggregate{Point: CoursePoint{PointType: PointTypeExploration, Title: "Node"}}, "Lesson", "en")
	if !strings.Contains(questionPrompt, "Respond in English") {
		t.Fatalf("expected english question prompt instruction, got %q", questionPrompt)
	}

	draft := parsePointSelfEvaluationDraft("", &CoursePointAggregate{Point: CoursePoint{Title: "Node"}}, "en")
	if len(draft.ApplicationQuestions) == 0 || strings.Contains(draft.GoalAlignmentNote, "현재 목표") {
		t.Fatalf("expected english self evaluation fallback, got %+v", draft)
	}
}
