package llm

import (
	"context"
	"fmt"
)

// Client is the common interface for all LLM providers.
type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type CompletionUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

type UsageReporter interface {
	LastUsage() (CompletionUsage, bool)
	ModelName() string
}

// NewClient creates an LLM client for the given provider.
// endpointURL is optional and currently used only by the llama provider.
func NewClient(provider, apiKey, modelName string, endpointURL ...string) (Client, error) {
	ep := ""
	if len(endpointURL) > 0 {
		ep = endpointURL[0]
	}
	switch provider {
	case "openai":
		return NewOpenAIClient(apiKey, modelName), nil
	case "google":
		return NewGoogleClient(apiKey, modelName), nil
	case "grok":
		return NewGrokClient(apiKey, modelName), nil
	case "anthropic":
		return NewAnthropicClient(apiKey, modelName), nil
	case "solar":
		return NewSolarClient(apiKey, modelName), nil
	case "hyperclova":
		return NewHyperCLOVAClient(apiKey, modelName), nil
	case "llama":
		return NewLlamaClient(apiKey, modelName, ep), nil
	case "exaone":
		return NewExaoneClient(apiKey, modelName), nil
	default:
		return nil, fmt.Errorf("llm: unknown provider %q", provider)
	}
}
