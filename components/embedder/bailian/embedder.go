package bailian

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/cloudwego/eino/components/embedding"
)

type EmbeddingConfig struct {
	APIKey         string
	Model          string
	BaseURL        string
	Dimension      *int
	TextType       *string
	OutputType     *string
	EncodingFormat *string
	Instruct       *string
}

type Embedder struct {
	cfg        *EmbeddingConfig
	httpClient *http.Client
	endpoint   string
}

func NewEmbedder(ctx context.Context, cfg *EmbeddingConfig) (embedding.Embedder, error) {
	if cfg == nil || cfg.Model == "" || cfg.APIKey == "" {
		return nil, errors.New("invalid embedding config")
	}
	e := &Embedder{
		cfg:        cfg,
		httpClient: http.DefaultClient,
		endpoint:   normalizeBaseURL(cfg.BaseURL),
	}
	return e, nil
}

func normalizeBaseURL(u string) string {
	if u == "" || strings.Contains(u, "/compatible-mode/") {
		if strings.Contains(u, "dashscope-intl.aliyuncs.com") {
			return "https://dashscope-intl.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
		}
		return "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	}
	if strings.HasSuffix(u, "#") {
		return strings.TrimSuffix(u, "#")
	}
	return u
}

type apiRequest struct {
	Model      string         `json:"model"`
	Input      apiInput       `json:"input"`
	Parameters *apiParameters `json:"parameters,omitempty"`
}

type apiInput struct {
	Texts []string `json:"texts"`
}

type apiParameters struct {
	Dimension      *int   `json:"dimension,omitempty"`
	TextType       string `json:"text_type,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
	OutputType     string `json:"output_type,omitempty"`
	Instruct       string `json:"instruct,omitempty"`
}

type apiResponse struct {
	Output struct {
		Embeddings []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
}

type apiResponseFull struct {
	Output struct {
		Embeddings []struct {
			SparseEmbedding []struct {
				Index int     `json:"index"`
				Value float64 `json:"value"`
				Token string  `json:"token"`
			} `json:"sparse_embedding,omitempty"`
			Embedding []float64 `json:"embedding"`
			TextIndex int       `json:"text_index"`
		} `json:"embeddings"`
	} `json:"output"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (e *Embedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, errors.New("texts is empty")
	}

	if strings.Contains(e.cfg.Model, "text-embedding-v3") || strings.Contains(e.cfg.Model, "text-embedding-v4") {
		if len(texts) > 10 {
			return nil, errors.New("too many input texts for v3/v4 (<=10)")
		}
	} else if len(texts) > 25 {
		return nil, errors.New("too many input texts (<=25)")
	}

	textType := "document"
	if e.cfg.TextType != nil && *e.cfg.TextType != "" {
		textType = *e.cfg.TextType
	}
	encoding := "float"
	if e.cfg.EncodingFormat != nil && *e.cfg.EncodingFormat != "" {
		encoding = *e.cfg.EncodingFormat
	}
	var outputType, instruct string
	if e.cfg.OutputType != nil {
		outputType = *e.cfg.OutputType
	}
	if e.cfg.Instruct != nil {
		instruct = *e.cfg.Instruct
	}

	reqBody := apiRequest{
		Model: e.cfg.Model,
		Input: apiInput{Texts: texts},
		Parameters: &apiParameters{
			Dimension:      e.cfg.Dimension,
			TextType:       textType,
			EncodingFormat: encoding,
			OutputType:     outputType,
			Instruct:       instruct,
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)

	resp, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ar apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, err
	}

	if len(ar.Output.Embeddings) == 0 {
		return nil, errors.New("empty embeddings")
	}

	out := make([][]float64, 0, len(ar.Output.Embeddings))
	for _, item := range ar.Output.Embeddings {
		out = append(out, item.Embedding)
	}
	return out, nil
}

func (e *Embedder) EmbedStringsExt(ctx context.Context, texts []string, opts ...embedding.Option) (*domain.EmbeddingsResponse, error) {
	if len(texts) == 0 {
		return nil, errors.New("texts is empty")
	}

	if strings.Contains(e.cfg.Model, "text-embedding-v3") || strings.Contains(e.cfg.Model, "text-embedding-v4") {
		if len(texts) > 10 {
			return nil, errors.New("too many input texts for v3/v4 (<=10)")
		}
	} else if len(texts) > 25 {
		return nil, errors.New("too many input texts (<=25)")
	}

	textType := "document"
	if e.cfg.TextType != nil && *e.cfg.TextType != "" {
		textType = *e.cfg.TextType
	}
	encoding := "float"
	if e.cfg.EncodingFormat != nil && *e.cfg.EncodingFormat != "" {
		encoding = *e.cfg.EncodingFormat
	}
	var outputType, instruct string
	if e.cfg.OutputType != nil {
		outputType = *e.cfg.OutputType
	}
	if e.cfg.Instruct != nil {
		instruct = *e.cfg.Instruct
	}

	reqBody := apiRequest{
		Model: e.cfg.Model,
		Input: apiInput{Texts: texts},
		Parameters: &apiParameters{
			Dimension:      e.cfg.Dimension,
			TextType:       textType,
			EncodingFormat: encoding,
			OutputType:     outputType,
			Instruct:       instruct,
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)

	resp, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var ar apiResponseFull
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, err
	}

	if len(ar.Output.Embeddings) == 0 {
		return nil, errors.New("empty embeddings")
	}

	out := domain.EmbeddingsResponse{
		Embeddings: make([]domain.EmbeddingItem, 0, len(ar.Output.Embeddings)),
		Usage: domain.EmbeddingUsage{
			TotalTokens: ar.Usage.TotalTokens,
		},
	}

	for _, item := range ar.Output.Embeddings {
		se := make([]domain.SparseEmbedding, 0, len(item.SparseEmbedding))
		for _, s := range item.SparseEmbedding {
			se = append(se, domain.SparseEmbedding{Index: s.Index, Value: s.Value, Token: s.Token})
		}
		out.Embeddings = append(out.Embeddings, domain.EmbeddingItem{
			SparseEmbedding: se,
			Embedding:       item.Embedding,
			TextIndex:       item.TextIndex,
		})
	}

	return &out, nil
}
