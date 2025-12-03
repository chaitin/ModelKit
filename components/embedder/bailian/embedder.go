package bailian

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
