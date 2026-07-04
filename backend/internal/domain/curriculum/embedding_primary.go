package curriculum

import (
	"os"
	"strings"

	"github.com/learnweaver/backend/internal/pkg/embedder"
)

func newPrimaryEmbeddingClient() (embedder.Embedder, error) {
	return embedder.NewProviderClient(primaryEmbeddingProvider(), strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")))
}
