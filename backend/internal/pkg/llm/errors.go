package llm

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
)

func safeLLMAPIError(provider string, statusCode int) error {
	return fmt.Errorf("%s: API error %d", strings.TrimSpace(provider), statusCode)
}

func safeLLMRequestError(provider string, err error) error {
	provider = strings.TrimSpace(provider)
	if err == nil {
		return fmt.Errorf("%s: request failed", provider)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%s: request canceled", provider)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: request timeout", provider)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("%s: request timeout", provider)
	}
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "timeout") || strings.Contains(lower, "deadline exceeded") || strings.Contains(lower, "i/o timeout") {
		return fmt.Errorf("%s: request timeout", provider)
	}
	return fmt.Errorf("%s: request failed", provider)
}
