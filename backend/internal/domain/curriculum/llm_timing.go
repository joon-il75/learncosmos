package curriculum

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/learnweaver/backend/internal/pkg/llm"
)

type timedLLMClient struct {
	base      llm.Client
	elapsedMS atomic.Int64
}

func newTimedLLMClient(base llm.Client) *timedLLMClient {
	return &timedLLMClient{base: base}
}

func (c *timedLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	startedAt := time.Now()
	raw, err := c.base.Complete(ctx, prompt)
	c.elapsedMS.Add(time.Since(startedAt).Milliseconds())
	return raw, err
}

func (c *timedLLMClient) ElapsedMS() int64 {
	if c == nil {
		return 0
	}
	return c.elapsedMS.Load()
}
