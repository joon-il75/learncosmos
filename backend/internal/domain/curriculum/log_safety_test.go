package curriculum

import (
	"os"
	"strings"
	"testing"
)

func TestCourseDraftLogsDoNotWriteRawSourceQuery(t *testing.T) {
	files := []string{
		"handler_draft_create.go",
		"handler_draft_build.go",
		"course_draft_worker.go",
	}
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		text := string(source)
		for _, forbidden := range []string{
			"source_query=%q",
			"source_query=%s",
		} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("%s contains raw source query log token %q", file, forbidden)
			}
		}
	}
}

func TestRecordBYOKValidationStoresProviderErrorCode(t *testing.T) {
	source, err := os.ReadFile("repository_config.go")
	if err != nil {
		t.Fatalf("read repository_config.go: %v", err)
	}
	text := string(source)
	if !strings.Contains(text, "trimmed := logsafe.ProviderErrorCode(message)") {
		t.Fatal("RecordBYOKValidation must store provider error codes instead of raw validation errors")
	}
	if strings.Contains(text, "trimmed := message") {
		t.Fatal("RecordBYOKValidation must not persist raw validation error messages")
	}
}
