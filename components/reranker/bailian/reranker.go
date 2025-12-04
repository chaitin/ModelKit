package bailian

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

type RerankerConfig struct {
	APIKey  string
	Model   string
	BaseUrl string
}

type RerankRequest struct {
	Model      string           `json:"model"`
	Input      RerankInput      `json:"input"`
	Parameters *RerankParameter `json:"parameters,omitempty"`
}

type RerankInput struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type RerankParameter struct {
	ReturnDocuments bool   `json:"return_documents"`
	TopN            int    `json:"top_n,omitempty"`
	Instruct        string `json:"instruct,omitempty"`
}

type apiResponse struct {
	Output struct {
		Results []apiResult `json:"results"`
	} `json:"output"`
	Usage *apiUsage `json:"usage,omitempty"`
}

type apiResult struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       *apiDoc `json:"document,omitempty"`
}

type apiDoc struct {
	Text string `json:"text"`
}

type apiUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func NewReranker(ctx context.Context, config RerankerConfig) *Reranker {
	config.BaseUrl = normalizeBaseUrl(config.BaseUrl)
	return &Reranker{
		Ctx:    ctx,
		Config: config,
	}
}

func normalizeBaseUrl(u string) string {
	if u == "" {
		return "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
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
		Model: r.Config.Model,
		Input: RerankInput{
			Query:     req.Query,
			Documents: req.Documents,
		},
		Parameters: &RerankParameter{
			ReturnDocuments: req.ReturnDocuments,
			TopN:            *N,
		},
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
	defer func() { _ = rawResp.Body.Close() }()

	if rawResp.StatusCode != http.StatusOK {
		return domain.RerankResponse{}, errors.New("request failed, status code: " + strconv.Itoa(rawResp.StatusCode))
	}

	var resp apiResponse
	if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
		return domain.RerankResponse{}, err
	}
	if len(resp.Output.Results) == 0 {
		return domain.RerankResponse{}, errors.New("empty results")
	}

	if *N < len(resp.Output.Results) {
		resp.Output.Results = resp.Output.Results[:*N]
	}

	var rerankResp domain.RerankResponse
	rerankResp.Results = lo.Map(resp.Output.Results, func(item apiResult, _ int) domain.Result {
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

	if resp.Usage != nil {
		pt := resp.Usage.PromptTokens
		tt := resp.Usage.TotalTokens
		it := resp.Usage.InputTokens
		ot := resp.Usage.OutputTokens

		rerankResp.Usage = &domain.Usage{PromptTokens: pt, TotalTokens: tt, InputTokens: it, OutputTokens: ot}
	}

	return rerankResp, nil
}
