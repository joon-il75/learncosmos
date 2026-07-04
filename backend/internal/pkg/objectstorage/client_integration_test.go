package objectstorage

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestIntegrationNaverObjectStorage(t *testing.T) {
	if os.Getenv("OBJECT_STORAGE_INTEGRATION") != "1" {
		t.Skip("set OBJECT_STORAGE_INTEGRATION=1 to run")
	}
	_ = godotenv.Load("../../../.env")

	client, err := NewClient(LoadConfigFromEnv())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	key := "tmp/integration-probes/objectstorage-" + time.Now().UTC().Format("20060102T150405") + ".txt"
	body := []byte("learnweaver object storage probe\n")
	if err := client.PutObject(context.Background(), key, body, "text/plain; charset=utf-8"); err != nil {
		t.Fatalf("put object: %v", err)
	}

	rawURL, _, err := client.PresignGetObject(key, time.Minute)
	if err != nil {
		t.Fatalf("presign get object: %v", err)
	}
	resp, err := http.Get(rawURL)
	if err != nil {
		t.Fatalf("get presigned object: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		t.Fatalf("presigned get status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read presigned object: %v", err)
	}
	if string(got) != string(body) {
		t.Fatalf("presigned object mismatch: got %q want %q", string(got), string(body))
	}
}
