package domain

import "context"

type Reranker interface {
	Rerank(ctx context.Context, req RerankRequest) (RerankResponse, error)
}

type RerankRequest struct {
	N               *int     `json:"n,omitempty"`
	Documents       []string `json:"documents"`
	Query           string   `json:"query"`
	ReturnDocuments bool     `json:"return_documents"`
}

type RerankResponse struct {
	Results []Result `json:"results"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Result struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       string  `json:"document,omitempty"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
