package user

import (
	"os"
	"strings"
	"testing"
)

func TestNullableUserEmailMigrationContract(t *testing.T) {
	body, err := os.ReadFile("../../../migrations/117_allow_nullable_user_email.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	sql := string(body)

	required := []string{
		"ALTER COLUMN email DROP NOT NULL",
		"SET email = NULL",
		"DROP CONSTRAINT IF EXISTS users_email_key",
		"idx_users_email_unique_present",
		"WHERE email IS NOT NULL AND btrim(email) <> ''",
	}
	for _, needle := range required {
		if !strings.Contains(sql, needle) {
			t.Fatalf("migration missing %q", needle)
		}
	}
}
