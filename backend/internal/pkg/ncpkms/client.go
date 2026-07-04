package ncpkms

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://kms.apigw.ntruss.com"
	defaultEncryptPath  = "/kms/v1/keys/%s/encrypt"
	defaultDecryptPath  = "/kms/v1/keys/%s/decrypt"
	defaultPlainField   = "plaintext"
	defaultCipherField  = "ciphertext"
	defaultHTTPTimeout  = 15 * time.Second
	defaultCallDeadline = 20 * time.Second
)

// Client is a minimal Naver Cloud KMS client for runtime secret envelopes.
type Client struct {
	baseURL     string
	keyRef      string
	accessKey   string
	secretKey   string
	httpClient  *http.Client
	encryptPath string
	decryptPath string
	plainField  string
	cipherField string
}

// NewFromEnv builds a KMS client from the same environment contract used by the
// operational KMS scripts. By default it reuses the unified NCP key pair; set
// NCP_KMS_CREDENTIAL_MODE=dedicated to require NCP_KMS_ACCESS_KEY/SECRET_KEY.
func NewFromEnv() (*Client, error) {
	baseURL := strings.TrimRight(firstNonEmpty(os.Getenv("NCP_KMS_BASE_URL"), defaultBaseURL), "/")
	keyRef := firstNonEmpty(os.Getenv("NCP_KMS_KEY_TAG"), os.Getenv("NCP_KMS_KEY_NAME"))
	if strings.TrimSpace(keyRef) == "" {
		return nil, errors.New("NCP_KMS_KEY_TAG or NCP_KMS_KEY_NAME is required")
	}
	accessKey, secretKey, err := credentialFromEnv()
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:     baseURL,
		keyRef:      strings.TrimSpace(keyRef),
		accessKey:   accessKey,
		secretKey:   secretKey,
		httpClient:  &http.Client{Timeout: defaultHTTPTimeout},
		encryptPath: firstNonEmpty(os.Getenv("NCP_KMS_ENCRYPT_PATH_TEMPLATE"), defaultEncryptPath),
		decryptPath: firstNonEmpty(os.Getenv("NCP_KMS_DECRYPT_PATH_TEMPLATE"), defaultDecryptPath),
		plainField:  firstNonEmpty(os.Getenv("NCP_KMS_PLAINTEXT_FIELD"), defaultPlainField),
		cipherField: firstNonEmpty(os.Getenv("NCP_KMS_CIPHERTEXT_FIELD"), defaultCipherField),
	}, nil
}

// EncryptSecret encrypts plaintext and returns the ciphertext plus the key tag
// used by the current runtime configuration.
func (c *Client) EncryptSecret(plaintext string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultCallDeadline)
	defer cancel()
	ciphertext, err := c.encrypt(ctx, plaintext)
	if err != nil {
		return "", "", err
	}
	return ciphertext, c.keyRef, nil
}

// DecryptSecret decrypts ciphertext. keyTag is accepted for the secret envelope
// interface; NCP KMS routing is controlled by the runtime key reference.
func (c *Client) DecryptSecret(ciphertext string, _ string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultCallDeadline)
	defer cancel()
	return c.decrypt(ctx, ciphertext)
}

func (c *Client) encrypt(ctx context.Context, plaintext string) (string, error) {
	payload := map[string]string{c.plainField: base64.StdEncoding.EncodeToString([]byte(plaintext))}
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

func (c *Client) decrypt(ctx context.Context, ciphertext string) (string, error) {
	payload := map[string]string{c.cipherField: ciphertext}
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

func (c *Client) do(ctx context.Context, method, requestPath string, payload map[string]string) ([]byte, error) {
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

func (c *Client) path(template string) string {
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, url.PathEscape(c.keyRef))
	}
	return template
}

func (c *Client) signature(method, requestURI, timestamp string) string {
	message := method + " " + requestURI + "\n" + timestamp + "\n" + c.accessKey
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func credentialFromEnv() (string, string, error) {
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
		s, ok := value.(string)
		return s, ok
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
