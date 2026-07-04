package content

import (
	"errors"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/urlsafe"
)

var errInvalidContentURL = errors.New("invalid content url")

func normalizeCreateContentRequestURLs(req *CreateContentRequest) error {
	required := req.ContentType != ContentTypeInternal
	urlValue, err := normalizeStoredContentURL(req.URL)
	if err != nil {
		return err
	}
	if required && urlValue == nil {
		return errInvalidContentURL
	}
	req.URL = urlValue
	return nil
}

func normalizeUpdateContentRequestURLs(req *UpdateContentRequest) error {
	if req.URL == nil {
		return nil
	}
	urlValue, err := normalizeStoredContentURL(req.URL)
	if err != nil {
		return err
	}
	if urlValue == nil {
		return errInvalidContentURL
	}
	req.URL = urlValue
	return nil
}

func normalizeStoredContentURL(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}
	normalized, ok := urlsafe.NormalizeHTTPURL(trimmed)
	if !ok {
		return nil, errInvalidContentURL
	}
	return &normalized, nil
}
