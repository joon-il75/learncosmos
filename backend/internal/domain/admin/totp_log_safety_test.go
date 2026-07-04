package admin

import (
	"os"
	"strings"
	"testing"
)

func TestTOTPValidationDoesNotLogCodesOrSecretMetadata(t *testing.T) {
	source, err := os.ReadFile("totp.go")
	if err != nil {
		t.Fatalf("read totp.go: %v", err)
	}
	if !strings.Contains(string(source), "secretmanager.Decrypt(encryptedSecret)") {
		t.Fatal("ValidateTOTP must use secretmanager.Decrypt for v1/v2 read compatibility")
	}

	for _, forbidden := range []string{
		"[TOTP] Attempt",
		"Code:",
		"SecretLen",
	} {
		if strings.Contains(string(source), forbidden) {
			t.Fatalf("totp.go contains forbidden log-sensitive token %q", forbidden)
		}
	}
}

func TestRecommendationDebugLogsDoNotWriteRawScenarioInputs(t *testing.T) {
	checks := map[string][]string{
		"recommendation_debug_scenario_handler.go": {"course_title=%q", "req.CourseTitle, err"},
		"recommendation_debug_external_handler.go": {"req.URL, err"},
	}
	for file, forbiddenTokens := range checks {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		text := string(source)
		for _, forbidden := range forbiddenTokens {
			if strings.Contains(text, forbidden) {
				t.Fatalf("%s contains forbidden log-sensitive token %q", file, forbidden)
			}
		}
	}
}
