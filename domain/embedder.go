package domain

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
)

type SparseEmbedding struct {
	Index int     `json:"index"`
	Value float64 `json:"value"`
	Token string  `json:"token"`
}

type EmbeddingItem struct {
	SparseEmbedding []SparseEmbedding `json:"sparse_embedding,omitempty"`
	Embedding       []float64         `json:"embedding"`
	TextIndex       int               `json:"text_index"`
}

type EmbeddingUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type EmbeddingsResponse struct {
	Embeddings []EmbeddingItem `json:"embeddings"`
	Usage  EmbeddingUsage   `json:"usage"`
}

type EmbedderExt interface {
	EmbedStringsExt(ctx context.Context, texts []string, opts ...embedding.Option) (*EmbeddingsResponse, error)
}
