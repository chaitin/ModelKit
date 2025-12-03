package baai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/samber/lo"
)

type Reranker struct {
	Ctx    context.Context
	Config RerankerConfig
}

type RerankRequest struct {
	Model           string   `json:"model"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	ReturnDocuments bool     `json:"return_documents"`
}

type RerankResponse struct {
	Model   string         `json:"model"`
	Results []RerankResult `json:"results"`
	Usage   *RerankUsage   `json:"usage,omitempty"`
}

type RerankResult struct {
	Index          int             `json:"index"`
	RelevanceScore float64         `json:"relavance_score"`
	Document       *RerankDocument `json:"document,omitempty"`
}

type RerankDocument struct {
	Text string `json:"text"`
}

type RerankUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type RerankerConfig struct {
	APIKey  string
	Model   string
	BaseUrl string
}

func NewReranker(ctx context.Context, config RerankerConfig) *Reranker {
	config.BaseUrl = normalizeBaseUrl(config.BaseUrl)
	return &Reranker{
		Ctx:    ctx,
		Config: config,
	}
}

func normalizeBaseUrl(u string) string {
	if strings.HasSuffix(u, "/rerank") {
		return u
	}

	if !strings.HasSuffix(u, "/rerank") {
		return u + "/rerank"
	}

	if strings.HasSuffix(u, "#") {
		return strings.TrimSuffix(u, "#")
	}

	return u
}

func (r *Reranker) Rerank(ctx context.Context, req domain.RerankRequest) (domain.RerankResponse, error) {
	if len(req.Documents) == 0 || r.Config.Model == "" || r.Config.BaseUrl == "" {
		return domain.RerankResponse{}, errors.New("invalid params")
	}

	N := req.N
	if N == nil {
		lens := len(req.Documents)
		N = &lens
	} else {
		if *N <= 0 {
			*N = 1
		} else if *N > len(req.Documents) {
			*N = len(req.Documents)
		}
	}

	body := RerankRequest{
		Model:           r.Config.Model,
		Query:           req.Query,
		Documents:       req.Documents,
		ReturnDocuments: req.ReturnDocuments,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return domain.RerankResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.Config.BaseUrl, bytes.NewReader(raw))
	if err != nil {
		return domain.RerankResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+r.Config.APIKey)

	client := http.DefaultClient
	client.Transport = http.DefaultTransport

	rawResp, err := client.Do(httpReq)
	if err != nil {
		return domain.RerankResponse{}, err
	}
	defer rawResp.Body.Close()

	if rawResp.StatusCode != http.StatusOK {
		return domain.RerankResponse{}, errors.New("request failed, status code: " + strconv.Itoa(rawResp.StatusCode))
	}

	var resp RerankResponse
	if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
		return domain.RerankResponse{}, err
	}
	if len(resp.Results) == 0 {
		return domain.RerankResponse{}, errors.New("empty results")
	}

	if *N < len(resp.Results) {
		resp.Results = resp.Results[:*N]
	}
	var rerankResp domain.RerankResponse

	// 使用 lo 库将 resp.Results 转成 rerankResp
	rerankResp.Results = lo.Map(resp.Results, func(item RerankResult, _ int) domain.Result {
		var doc string
		if item.Document != nil {
			doc = item.Document.Text
		}
		return domain.Result{
			Index:          item.Index,
			RelevanceScore: item.RelevanceScore,
			Document:       doc,
		}
	})

	return rerankResp, nil
}
