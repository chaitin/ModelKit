package domain

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
)

type SparseEntry struct {
	Index int     `json:"index"`
	Value float64 `json:"value"`
	Token string  `json:"token"`
}

type EmbeddingItem struct {
	SparseEmbedding []SparseEntry `json:"sparse_embedding,omitempty"`
	Embedding       []float64     `json:"embedding"`
	TextIndex       int           `json:"text_index"`
}

type EmbeddingsOutput struct {
	Embeddings []EmbeddingItem `json:"embeddings"`
}

type EmbeddingUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type EmbeddingsResponse struct {
	Output EmbeddingsOutput `json:"output"`
	Usage  EmbeddingUsage   `json:"usage"`
}

type EmbedderExt interface {
	EmbedStringsExt(ctx context.Context, texts []string, opts ...embedding.Option) (*EmbeddingsResponse, error)
}
