package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type encryptedEnvFile struct {
	Version    int    `json:"version"`
	Type       string `json:"type"`
	Provider   string `json:"provider"`
	CreatedAt  string `json:"created_at"`
	Ciphertext string `json:"ciphertext"`
}

type kmsClient struct {
	baseURL     string
	keyRef      string
	accessKey   string
	secretKey   string
	httpClient  *http.Client
	encryptPath string
	decryptPath string
}

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ncp-kms-env failed: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: ncp_kms_env encrypt|decrypt")
	}
	switch args[0] {
	case "encrypt":
		return runEncrypt(ctx, args[1:])
	case "decrypt":
		return runDecrypt(ctx, args[1:])
	default:
		return errors.New("usage: ncp_kms_env encrypt|decrypt")
	}
}

func runEncrypt(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("encrypt", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inPath := fs.String("in", "", "input plaintext env path")
	outPath := fs.String("out", "", "output encrypted json path")
	if err := fs.Parse(args); err != nil {
		return errors.New("invalid encrypt flags")
	}
	if strings.TrimSpace(*inPath) == "" || strings.TrimSpace(*outPath) == "" {
		return errors.New("encrypt requires -in and -out")
	}
	plaintext, err := os.ReadFile(*inPath)
	if err != nil {
		return errors.New("read plaintext env")
	}
	if err := validateEnvBytes(plaintext); err != nil {
		return err
	}
	client, err := newKMSClientFromEnv()
	if err != nil {
		return err
	}
	ciphertext, err := client.encrypt(ctx, string(plaintext))
	if err != nil {
		return err
	}
	sealed := encryptedEnvFile{
		Version:    1,
		Type:       "learnweaver-backend-env",
		Provider:   "ncp-kms",
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		Ciphertext: ciphertext,
	}
	encoded, err := json.MarshalIndent(sealed, "", "  ")
	if err != nil {
		return errors.New("marshal encrypted env")
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*outPath, encoded, 0600); err != nil {
		return errors.New("write encrypted env")
	}
	fmt.Println("backend env encrypted ok")
	return nil
}

func runDecrypt(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("decrypt", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inPath := fs.String("in", "", "input encrypted json path")
	outPath := fs.String("out", "", "output plaintext env path")
	if err := fs.Parse(args); err != nil {
		return errors.New("invalid decrypt flags")
	}
	if strings.TrimSpace(*inPath) == "" || strings.TrimSpace(*outPath) == "" {
		return errors.New("decrypt requires -in and -out")
	}
	sealedBytes, err := os.ReadFile(*inPath)
	if err != nil {
		return errors.New("read encrypted env")
	}
	var sealed encryptedEnvFile
	if err := json.Unmarshal(sealedBytes, &sealed); err != nil {
		return errors.New("parse encrypted env")
	}
	if sealed.Version != 1 || sealed.Type != "learnweaver-backend-env" || sealed.Provider != "ncp-kms" || strings.TrimSpace(sealed.Ciphertext) == "" {
		return errors.New("unsupported encrypted env format")
	}
	client, err := newKMSClientFromEnv()
	if err != nil {
		return err
	}
	plaintext, err := client.decrypt(ctx, sealed.Ciphertext)
	if err != nil {
		return err
	}
	plaintextBytes := []byte(plaintext)
	if err := validateEnvBytes(plaintextBytes); err != nil {
		return err
	}
	if err := os.WriteFile(*outPath, plaintextBytes, 0600); err != nil {
		return errors.New("write plaintext env")
	}
	fmt.Println("backend env decrypted ok")
	return nil
}

func newKMSClientFromEnv() (*kmsClient, error) {
	baseURL := strings.TrimRight(firstNonEmpty(os.Getenv("NCP_KMS_BASE_URL"), "https://kms.apigw.ntruss.com"), "/")
	keyRef := firstNonEmpty(os.Getenv("NCP_KMS_KEY_TAG"), os.Getenv("NCP_KMS_KEY_NAME"))
	if strings.TrimSpace(keyRef) == "" {
		return nil, errors.New("NCP_KMS_KEY_TAG or NCP_KMS_KEY_NAME is required")
	}
	accessKey, secretKey, err := kmsCredentialFromEnv()
	if err != nil {
		return nil, err
	}
	return &kmsClient{
		baseURL:     baseURL,
		keyRef:      strings.TrimSpace(keyRef),
		accessKey:   accessKey,
		secretKey:   secretKey,
		httpClient:  &http.Client{Timeout: 15 * time.Second},
		encryptPath: firstNonEmpty(os.Getenv("NCP_KMS_ENCRYPT_PATH_TEMPLATE"), "/kms/v1/keys/%s/encrypt"),
		decryptPath: firstNonEmpty(os.Getenv("NCP_KMS_DECRYPT_PATH_TEMPLATE"), "/kms/v1/keys/%s/decrypt"),
	}, nil
}

func (c *kmsClient) encrypt(ctx context.Context, plaintext string) (string, error) {
	payload := map[string]string{"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext))}
	body, err := c.do(ctx, http.MethodPost, c.path(c.encryptPath), payload)
	if err != nil {
		return "", fmt.Errorf("encrypt request: %w", err)
	}
	ciphertext := firstJSONValue(body, "data.ciphertext", "ciphertext", "result.ciphertext", "data")
	if strings.TrimSpace(ciphertext) == "" {
		return "", errors.New("encrypt response did not include ciphertext")
	}
	return ciphertext, nil
}

func (c *kmsClient) decrypt(ctx context.Context, ciphertext string) (string, error) {
	payload := map[string]string{"ciphertext": ciphertext}
	body, err := c.do(ctx, http.MethodPost, c.path(c.decryptPath), payload)
	if err != nil {
		return "", fmt.Errorf("decrypt request: %w", err)
	}
	plaintext := firstJSONValue(body, "data.plaintext", "plaintext", "result.plaintext", "data")
	if strings.TrimSpace(plaintext) == "" {
		return "", errors.New("decrypt response did not include plaintext")
	}
	if decoded, err := base64.StdEncoding.DecodeString(plaintext); err == nil {
		return string(decoded), nil
	}
	return plaintext, nil
}

func (c *kmsClient) do(ctx context.Context, method, requestPath string, payload map[string]string) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("marshal request")
	}
	endpoint, err := url.Parse(c.baseURL + requestPath)
	if err != nil {
		return nil, errors.New("invalid endpoint")
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("create request")
	}
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", c.accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", c.signature(method, endpoint.RequestURI(), timestamp))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("request failed")
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}
	return respBody, nil
}

func (c *kmsClient) path(template string) string {
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, url.PathEscape(c.keyRef))
	}
	return template
}

func (c *kmsClient) signature(method, requestURI, timestamp string) string {
	message := method + " " + requestURI + "\n" + timestamp + "\n" + c.accessKey
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func validateEnvBytes(data []byte) error {
	content := string(data)
	required := []string{"APP_ENV", "DATABASE_URL", "JWT_SECRET", "ENCRYPTION_KEY"}
	for _, key := range required {
		if !hasEnvKey(content, key) {
			return fmt.Errorf("required env key missing: %s", key)
		}
	}
	return nil
}

func hasEnvKey(content, key string) bool {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+"=") {
			return true
		}
	}
	return false
}

func firstJSONValue(body []byte, paths ...string) string {
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return ""
	}
	for _, path := range paths {
		if value, ok := jsonPathString(data, strings.Split(path, ".")); ok {
			return value
		}
	}
	return ""
}

func jsonPathString(value any, path []string) (string, bool) {
	if len(path) == 0 {
		v, ok := value.(string)
		return v, ok
	}
	object, ok := value.(map[string]any)
	if !ok {
		return "", false
	}
	next, ok := object[path[0]]
	if !ok {
		return "", false
	}
	return jsonPathString(next, path[1:])
}

func kmsCredentialFromEnv() (string, string, error) {
	mode := strings.TrimSpace(strings.ToLower(os.Getenv("NCP_KMS_CREDENTIAL_MODE")))
	if mode == "dedicated" {
		accessKey := strings.TrimSpace(os.Getenv("NCP_KMS_ACCESS_KEY"))
		secretKey := strings.TrimSpace(os.Getenv("NCP_KMS_SECRET_KEY"))
		if accessKey == "" || secretKey == "" {
			return "", "", errors.New("NCP_KMS_ACCESS_KEY/NCP_KMS_SECRET_KEY is required when NCP_KMS_CREDENTIAL_MODE=dedicated")
		}
		return accessKey, secretKey, nil
	}

	accessKey := firstNonEmpty(os.Getenv("NCP_ACCESS_KEY"), os.Getenv("NCP_OBJECT_STORAGE_ACCESS_KEY"))
	secretKey := firstNonEmpty(os.Getenv("NCP_SECRET_KEY"), os.Getenv("NCP_OBJECT_STORAGE_SECRET_KEY"))
	if strings.TrimSpace(accessKey) == "" || strings.TrimSpace(secretKey) == "" {
		return "", "", errors.New("NCP_ACCESS_KEY/NCP_SECRET_KEY or NCP_OBJECT_STORAGE_ACCESS_KEY/NCP_OBJECT_STORAGE_SECRET_KEY is required")
	}
	return strings.TrimSpace(accessKey), strings.TrimSpace(secretKey), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
