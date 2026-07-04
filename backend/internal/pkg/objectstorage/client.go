package objectstorage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultEndpoint       = "https://kr.object.ncloudstorage.com"
	defaultRegion         = "kr-standard"
	defaultBucket         = "learnweaver-assets"
	defaultPresignTTL     = 5 * time.Minute
	defaultMaxUploadBytes = int64(20 * 1024 * 1024)
)

var ErrNotConfigured = errors.New("object storage is not configured")

type Config struct {
	Endpoint       string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
	PresignTTL     time.Duration
	MaxUploadBytes int64
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func LoadConfigFromEnv() Config {
	return Config{
		Endpoint:       getEnvString("NCP_OBJECT_STORAGE_ENDPOINT", defaultEndpoint),
		Region:         getEnvString("NCP_OBJECT_STORAGE_REGION", defaultRegion),
		Bucket:         getEnvString("NCP_OBJECT_STORAGE_BUCKET", defaultBucket),
		AccessKey:      strings.TrimSpace(os.Getenv("NCP_OBJECT_STORAGE_ACCESS_KEY")),
		SecretKey:      strings.TrimSpace(os.Getenv("NCP_OBJECT_STORAGE_SECRET_KEY")),
		PresignTTL:     getEnvDuration("NCP_OBJECT_STORAGE_PRESIGN_TTL_SECONDS", defaultPresignTTL),
		MaxUploadBytes: getEnvInt64("NCP_OBJECT_STORAGE_MAX_UPLOAD_BYTES", defaultMaxUploadBytes),
	}
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.AccessKey) == "" || strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, ErrNotConfigured
	}
	cfg.Endpoint = strings.TrimRight(strings.TrimSpace(cfg.Endpoint), "/")
	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.Bucket = strings.Trim(strings.TrimSpace(cfg.Bucket), "/")
	if cfg.Endpoint == "" || cfg.Region == "" || cfg.Bucket == "" {
		return nil, ErrNotConfigured
	}
	if cfg.PresignTTL <= 0 {
		cfg.PresignTTL = defaultPresignTTL
	}
	if cfg.MaxUploadBytes <= 0 {
		cfg.MaxUploadBytes = defaultMaxUploadBytes
	}
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}, nil
}

func (c *Client) Config() Config {
	return c.cfg
}

func (c *Client) PutObject(ctx context.Context, key string, body []byte, contentType string) error {
	if c == nil {
		return ErrNotConfigured
	}
	key = cleanObjectKey(key)
	if key == "" {
		return fmt.Errorf("object key is empty")
	}
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	sum := sha256.Sum256(body)
	payloadHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()
	objectURL := c.objectURL(key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, objectURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", amzDate(now))
	req.Header.Set("Authorization", c.authorizationHeader(req, payloadHash, now))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("object storage put failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (c *Client) PutObjectReader(ctx context.Context, key string, body io.ReadSeeker, size int64, contentType string) error {
	if c == nil {
		return ErrNotConfigured
	}
	key = cleanObjectKey(key)
	if key == "" {
		return fmt.Errorf("object key is empty")
	}
	if size < 0 {
		return fmt.Errorf("object size is invalid")
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return err
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, body); err != nil {
		return err
	}
	payloadHash := hex.EncodeToString(hash.Sum(nil))
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return err
	}

	now := time.Now().UTC()
	objectURL := c.objectURL(key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, objectURL, body)
	if err != nil {
		return err
	}
	req.ContentLength = size
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", amzDate(now))
	req.Header.Set("Authorization", c.authorizationHeader(req, payloadHash, now))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("object storage put failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (c *Client) DeleteObject(ctx context.Context, key string) error {
	if c == nil {
		return ErrNotConfigured
	}
	key = cleanObjectKey(key)
	if key == "" {
		return fmt.Errorf("object key is empty")
	}

	sum := sha256.Sum256(nil)
	payloadHash := hex.EncodeToString(sum[:])
	now := time.Now().UTC()
	objectURL := c.objectURL(key)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, objectURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", amzDate(now))
	req.Header.Set("Authorization", c.authorizationHeader(req, payloadHash, now))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("object storage delete failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (c *Client) GetObjectBytes(ctx context.Context, key string, maxBytes int64) ([]byte, string, error) {
	if c == nil {
		return nil, "", ErrNotConfigured
	}
	key = cleanObjectKey(key)
	if key == "" {
		return nil, "", fmt.Errorf("object key is empty")
	}
	if maxBytes <= 0 {
		maxBytes = 2 * 1024 * 1024
	}
	presignedURL, _, err := c.PresignGetObject(key, 0)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, presignedURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, "", fmt.Errorf("object storage get failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(body)) > maxBytes {
		return nil, "", fmt.Errorf("object exceeds max read size")
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) PresignGetObject(key string, ttl time.Duration) (string, time.Time, error) {
	if c == nil {
		return "", time.Time{}, ErrNotConfigured
	}
	key = cleanObjectKey(key)
	if key == "" {
		return "", time.Time{}, fmt.Errorf("object key is empty")
	}
	if ttl <= 0 {
		ttl = c.cfg.PresignTTL
	}
	if ttl > 30*time.Minute {
		ttl = 30 * time.Minute
	}

	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	parsed, err := url.Parse(c.objectURL(key))
	if err != nil {
		return "", time.Time{}, err
	}
	scope := c.credentialScope(now)
	query := parsed.Query()
	query.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	query.Set("X-Amz-Credential", c.cfg.AccessKey+"/"+scope)
	query.Set("X-Amz-Date", amzDate(now))
	query.Set("X-Amz-Expires", strconv.Itoa(int(ttl.Seconds())))
	query.Set("X-Amz-SignedHeaders", "host")
	parsed.RawQuery = canonicalQuery(query)

	canonicalRequest := strings.Join([]string{
		http.MethodGet,
		parsed.EscapedPath(),
		parsed.RawQuery,
		"host:" + parsed.Host + "\n",
		"host",
		"UNSIGNED-PAYLOAD",
	}, "\n")
	signature := c.signature(canonicalRequest, now)
	query.Set("X-Amz-Signature", signature)
	parsed.RawQuery = canonicalQuery(query)
	return parsed.String(), expiresAt, nil
}

func (c *Client) objectURL(key string) string {
	return c.cfg.Endpoint + "/" + escapePath(path.Join(c.cfg.Bucket, cleanObjectKey(key)))
}

func (c *Client) authorizationHeader(req *http.Request, payloadHash string, now time.Time) string {
	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := strings.Join([]string{
		"content-type:" + req.Header.Get("Content-Type"),
		"host:" + req.URL.Host,
		"x-amz-content-sha256:" + payloadHash,
		"x-amz-date:" + amzDate(now),
		"",
	}, "\n")
	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.EscapedPath(),
		"",
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")
	credential := c.cfg.AccessKey + "/" + c.credentialScope(now)
	return "AWS4-HMAC-SHA256 Credential=" + credential + ", SignedHeaders=" + signedHeaders + ", Signature=" + c.signature(canonicalRequest, now)
}

func (c *Client) signature(canonicalRequest string, now time.Time) string {
	hash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate(now),
		c.credentialScope(now),
		hex.EncodeToString(hash[:]),
	}, "\n")
	signingKey := deriveSigningKey(c.cfg.SecretKey, shortDate(now), c.cfg.Region, "s3")
	return hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
}

func (c *Client) credentialScope(now time.Time) string {
	return shortDate(now) + "/" + c.cfg.Region + "/s3/aws4_request"
}

func BuildLearningPointObjectKey(userID, courseID, pointID, originalName string) string {
	ext := strings.ToLower(path.Ext(originalName))
	if ext == "" {
		if exts, err := mime.ExtensionsByType(http.DetectContentType([]byte(originalName))); err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}
	base := strings.TrimSuffix(path.Base(originalName), path.Ext(originalName))
	base = sanitizeKeyPart(base)
	if base == "" {
		base = "attachment"
	}
	return BuildLearningPointObjectPrefix(userID, courseID, pointID) + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + base + ext
}

func BuildLearningPointObjectPrefix(userID, courseID, pointID string) string {
	return "learners/" + sanitizePathPart(userID) + "/courses/" + sanitizePathPart(courseID) + "/points/" + sanitizePathPart(pointID) + "/"
}

func IsLearningPointObjectKeyFor(userID, courseID, pointID, key string) bool {
	userID = sanitizePathPart(userID)
	courseID = sanitizePathPart(courseID)
	pointID = sanitizePathPart(pointID)
	if userID == "" || courseID == "" || pointID == "" {
		return false
	}
	key = cleanObjectKey(key)
	if key == "" {
		return false
	}
	prefix := BuildLearningPointObjectPrefix(userID, courseID, pointID)
	return strings.HasPrefix(key, prefix) && len(key) > len(prefix)
}

func cleanObjectKey(key string) string {
	key = strings.TrimSpace(strings.TrimPrefix(key, "/"))
	cleaned := path.Clean(key)
	if cleaned == "." || strings.HasPrefix(cleaned, "../") {
		return ""
	}
	return cleaned
}

func sanitizePathPart(value string) string {
	value = strings.Trim(strings.TrimSpace(value), "/")
	if value == "." || value == ".." || strings.Contains(value, "/") {
		return ""
	}
	return value
}

func sanitizeKeyPart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		keep := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if keep {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func escapePath(value string) string {
	parts := strings.Split(value, "/")
	for idx := range parts {
		parts[idx] = url.PathEscape(parts[idx])
	}
	return strings.Join(parts, "/")
}

func canonicalQuery(values url.Values) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0)
	for _, key := range keys {
		vals := append([]string(nil), values[key]...)
		sort.Strings(vals)
		for _, value := range vals {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}
	return strings.ReplaceAll(strings.Join(parts, "&"), "+", "%20")
}

func deriveSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func amzDate(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

func shortDate(t time.Time) string {
	return t.UTC().Format("20060102")
}

func getEnvString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

func getEnvInt64(key string, fallback int64) int64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
