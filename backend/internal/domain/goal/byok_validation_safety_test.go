package goal

import (
	"os"
	"strings"
	"testing"
)

func TestRecordBYOKValidationStoresProviderErrorCode(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatalf("read repository.go: %v", err)
	}
	text := string(source)
	if !strings.Contains(text, "trimmed := logsafe.ProviderErrorCode(message)") {
		t.Fatal("RecordBYOKValidation must store provider error codes instead of raw validation errors")
	}
	if strings.Contains(text, "trimmed := message") {
		t.Fatal("RecordBYOKValidation must not persist raw validation error messages")
	}
}
