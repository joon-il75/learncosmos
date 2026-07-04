package curriculum

import "testing"

func TestNormalizePointAIFeedback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "trim quoted text",
			raw:  "\"현재 답변은 방향이 분명합니다. 다음에는 적용 상황을 더 구체적으로 적어보세요.\"",
			want: "현재 답변은 방향이 분명합니다. 다음에는 적용 상황을 더 구체적으로 적어보세요.",
		},
		{
			name: "fallback on empty",
			raw:  "   ",
			want: "지금 답변에는 현재 목표와 연결하려는 방향이 보입니다. 다음에는 이 지식을 어디에 바로 써볼지 한 가지 상황을 더 구체적으로 적어보면 좋습니다.",
		},
		{
			name: "extract feedback from json object",
			raw:  `{"feedback":"질문을 풀 때 모집단과 표본을 먼저 구분해 보세요."}`,
			want: "질문을 풀 때 모집단과 표본을 먼저 구분해 보세요.",
		},
		{
			name: "trim outer braces from plain text",
			raw:  "{질문을 풀 때 모집단과 표본을 먼저 구분해 보세요.}",
			want: "질문을 풀 때 모집단과 표본을 먼저 구분해 보세요.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizePointAIFeedback(tt.raw, "ko"); got != tt.want {
				t.Fatalf("normalizePointAIFeedback() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizePointAIFeedbackEnglishFallback(t *testing.T) {
	t.Parallel()
	got := normalizePointAIFeedback(" ", "en")
	if got == "" || got == "지금 답변에는 현재 목표와 연결하려는 방향이 보입니다. 다음에는 이 지식을 어디에 바로 써볼지 한 가지 상황을 더 구체적으로 적어보면 좋습니다." {
		t.Fatalf("expected english fallback, got %q", got)
	}
}
