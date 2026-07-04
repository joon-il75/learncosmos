package llmjobs

import (
	"os"
	"strings"
	"testing"
)

func TestValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{name: "queued to running", from: StatusQueued, to: StatusRunning, want: true},
		{name: "running to succeeded", from: StatusRunning, to: StatusSucceeded, want: true},
		{name: "running to failed", from: StatusRunning, to: StatusFailed, want: true},
		{name: "succeeded terminal", from: StatusSucceeded, to: StatusRunning, want: false},
		{name: "failed terminal", from: StatusFailed, to: StatusRunning, want: false},
		{name: "unknown", from: Status("unknown"), to: StatusRunning, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStatusTransition(tt.from, tt.to); got != tt.want {
				t.Fatalf("ValidateStatusTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestContainsRedisBusyGroup(t *testing.T) {
	if !containsRedisBusyGroup("BUSYGROUP Consumer Group name already exists") {
		t.Fatal("expected BUSYGROUP error to match")
	}
	if containsRedisBusyGroup("ERR other") {
		t.Fatal("unexpected match")
	}
}

func TestRepositoryMarkFailedSanitizesPersistedErrorMessage(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatalf("read repository.go: %v", err)
	}
	text := string(source)
	if !strings.Contains(text, "errorMessage := logsafe.PersistedError(input.ErrorMessage)") {
		t.Fatal("MarkFailed must sanitize persisted llm job error messages")
	}
	if strings.Contains(text, "errorCode, input.ErrorMessage, input.Retryable") {
		t.Fatal("MarkFailed must not persist raw input.ErrorMessage")
	}
}
