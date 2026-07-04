package objectstorage

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLearningPointObjectKeyAccess(t *testing.T) {
	userID := "user-1"
	courseID := "course-1"
	pointID := "point-1"

	key := BuildLearningPointObjectKey(userID, courseID, pointID, "../Avatar PNG.JPG")
	if !IsLearningPointObjectKeyFor(userID, courseID, pointID, key) {
		t.Fatalf("expected generated key to be accessible: %s", key)
	}
	if !strings.HasPrefix(key, BuildLearningPointObjectPrefix(userID, courseID, pointID)) {
		t.Fatalf("expected generated key to use learning point prefix: %s", key)
	}
	if strings.Contains(key, "..") {
		t.Fatalf("expected generated key to remove path traversal fragments: %s", key)
	}

	cases := []struct {
		name     string
		owner    string
		course   string
		point    string
		key      string
		expected bool
	}{
		{name: "same owner", owner: userID, course: courseID, point: pointID, key: key, expected: true},
		{name: "leading slash is cleaned", owner: userID, course: courseID, point: pointID, key: "/" + key, expected: true},
		{name: "other owner", owner: "other-user", course: courseID, point: pointID, key: key, expected: false},
		{name: "other course", owner: userID, course: "other-course", point: pointID, key: key, expected: false},
		{name: "other point", owner: userID, course: courseID, point: "other-point", key: key, expected: false},
		{name: "prefix only", owner: userID, course: courseID, point: pointID, key: BuildLearningPointObjectPrefix(userID, courseID, pointID), expected: false},
		{name: "traversal key", owner: userID, course: courseID, point: pointID, key: "../" + key, expected: false},
		{name: "invalid owner part", owner: "../" + userID, course: courseID, point: pointID, key: key, expected: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsLearningPointObjectKeyFor(tc.owner, tc.course, tc.point, tc.key)
			if got != tc.expected {
				t.Fatalf("IsLearningPointObjectKeyFor() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestPresignGetObjectCapsTTL(t *testing.T) {
	client, err := NewClient(Config{
		Endpoint:   "https://example.test",
		Region:     "kr-standard",
		Bucket:     "bucket",
		AccessKey:  "access",
		SecretKey:  "secret",
		PresignTTL: 5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	rawURL, expiresAt, err := client.PresignGetObject("learners/user/courses/course/points/point/file.txt", time.Hour)
	if err != nil {
		t.Fatalf("PresignGetObject() error = %v", err)
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	if got := parsed.Query().Get("X-Amz-Expires"); got != "1800" {
		t.Fatalf("X-Amz-Expires = %q, want 1800", got)
	}
	if time.Until(expiresAt) > 31*time.Minute {
		t.Fatalf("expected expiresAt to be capped near 30 minutes, got %s", expiresAt)
	}
}
