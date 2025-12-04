package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestMain(m *testing.M) {
	f, err := os.Open(".env")
	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if k, v, ok := strings.Cut(line, "="); ok {
				os.Setenv(strings.TrimSpace(k), strings.TrimSpace(v))
			}
		}
		f.Close()
	}
	os.Exit(m.Run())
}

func TestBaiLian(t *testing.T) {
	if os.Getenv("bailianapikey") == "" {
		t.Skip("missing bailianapikey")
	}
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	t.Run("default", func(t *testing.T) {
		emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
			Provider:  consts.ModelProviderBaiLian,
			ModelName: "text-embedding-v4",
			BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
			APIKey:    os.Getenv("bailianapikey"),
		})
		if err != nil {
			t.Fatalf("GetEmbedder: %v", err)
		}
		texts := []string{"火鸡面", "测试文本", "向量模型"}
		res, err := mk.UseEmbedder(ctx, emb, texts)
		if err != nil {
			t.Skipf("UseEmbedder default: %v", err)
			return
		}
		if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) == 0 {
			t.Fatalf("empty embeddings")
		}
		t.Logf("[bailian default] texts=%d dim=%d", len(texts), len(res.Embeddings[0].Embedding))
		logHeadFloat(t, res.Embeddings[0].Embedding)
	})

	t.Run("dimension", func(t *testing.T) {
		for _, d := range []int{64, 128, 256, 512, 768, 1024, 1536, 2048} {
			dim := d
			t.Run(fmt.Sprintf("%d", dim), func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiLian,
					ModelName: "text-embedding-v4",
					BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
					APIKey:    os.Getenv("bailianapikey"),
					EmbedderParam: domain.EmbedderParam{
						Dimension: &dim,
					},
				})
				if err != nil {
					t.Skipf("GetEmbedder dimension=%d: %v", dim, err)
					return
				}
				res, err := mk.UseEmbedder(ctx, emb, []string{"火鸡面"})
				if err != nil {
					t.Skipf("UseEmbedder dimension=%d: %v", dim, err)
					return
				}
				if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) == 0 {
					t.Fatalf("empty embeddings")
				}
				t.Logf("[bailian dimension=%d] dim=%d", dim, len(res.Embeddings[0].Embedding))
				logHeadFloat(t, res.Embeddings[0].Embedding)
			})
		}
	})

	t.Run("text_type", func(t *testing.T) {
		for _, tt := range []string{"document", "query"} {
			t.Run(tt, func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiLian,
					ModelName: "text-embedding-v4",
					BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
					APIKey:    os.Getenv("bailianapikey"),
					EmbedderParam: domain.EmbedderParam{
						TextType: &tt,
					},
				})
				if err != nil {
					t.Skipf("GetEmbedder text_type=%s: %v", tt, err)
					return
				}
				res, err := mk.UseEmbedder(ctx, emb, []string{"火鸡面"})
				if err != nil {
					t.Skipf("UseEmbedder text_type=%s: %v", tt, err)
					return
				}
				if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) == 0 {
					t.Fatalf("empty embeddings")
				}
				t.Logf("[bailian text_type=%s] dim=%d", tt, len(res.Embeddings[0].Embedding))
				logHeadFloat(t, res.Embeddings[0].Embedding)
			})
		}
	})

	t.Run("output_type", func(t *testing.T) {
		for _, ot := range []string{"dense", "sparse", "dense&sparse"} {
			t.Run(ot, func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiLian,
					ModelName: "text-embedding-v4",
					BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
					APIKey:    os.Getenv("bailianapikey"),
					EmbedderParam: domain.EmbedderParam{
						OutputType: &ot,
					},
				})
				if err != nil {
					t.Skipf("GetEmbedder output_type=%s: %v", ot, err)
					return
				}
				res, err := mk.UseEmbedder(ctx, emb, []string{"火鸡面"})
				if err != nil {
					t.Skipf("UseEmbedder output_type=%s: %v", ot, err)
					return
				}
				t.Logf("[bailian output_type=%s] dim=%d", ot, len(res.Embeddings[0].Embedding))
				logHeadFloat(t, res.Embeddings[0].Embedding)
				if ot != "dense" {
					se := res.Embeddings[0].SparseEmbedding
					t.Logf("[bailian output_type=%s] sparse_nnz=%d", ot, len(se))
					logSparseEntriesHead(t, se)
				}
			})
		}
	})

	t.Run("encoding_format", func(t *testing.T) {
		for _, ef := range []string{"base64", "float"} {
			t.Run(ef, func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiLian,
					ModelName: "text-embedding-v4",
					BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
					APIKey:    os.Getenv("bailianapikey"),
					EmbedderParam: domain.EmbedderParam{
						EncodingFormat: &ef,
					},
				})
				if err != nil {
					t.Skipf("GetEmbedder encoding_format=%s: %v", ef, err)
					return
				}
				res, err := mk.UseEmbedder(ctx, emb, []string{"火鸡面"})
				if err != nil {
					t.Skipf("UseEmbedder encoding_format=%s: %v", ef, err)
					return
				}
				dim := len(res.Embeddings[0].Embedding)
				t.Logf("[bailian encoding_format=%s] dim=%d", ef, dim)
				logHeadFloat(t, res.Embeddings[0].Embedding)
			})
		}
	})

	t.Run("instruct", func(t *testing.T) {
		for _, instr := range []string{"检索：返回和食物相关的内容", "检索：返回和动物相关的内容"} {
			tt := "query"
			t.Run(instr, func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiLian,
					ModelName: "text-embedding-v4",
					BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
					APIKey:    os.Getenv("bailianapikey"),
					EmbedderParam: domain.EmbedderParam{
						TextType: &tt,
						Instruct: &instr,
					},
				})
				if err != nil {
					t.Skipf("GetEmbedder instruct=%s: %v", instr, err)
					return
				}
				res, err := mk.UseEmbedder(ctx, emb, []string{"火鸡面"})
				if err != nil {
					t.Skipf("UseEmbedder instruct=%s: %v", instr, err)
					return
				}
				if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) == 0 {
					t.Fatalf("empty embeddings")
				}
				t.Logf("[bailian instruct=%s text_type=%s] dim=%d", instr, tt, len(res.Embeddings[0].Embedding))
				logHeadFloat(t, res.Embeddings[0].Embedding)
			})
		}
	})
}

func TestOpenAI_Default(t *testing.T) {
	if os.Getenv("baizhiapikey") == "" {
		t.Skip("missing baizhiapikey")
	}
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)
	emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderOpenAI,
		ModelName: "bge-m3",
		BaseURL:   "https://model-square.app.baizhi.cloud/v1",
		APIKey:    os.Getenv("baizhiapikey"),
	})
	if err != nil {
		t.Fatalf("GetEmbedder: %v", err)
	}
	texts := []string{"风哀", "渚回"}
	res, err := mk.UseEmbedder(ctx, emb, texts)
	if err != nil {
		t.Fatalf("UseEmbedder: %v", err)
	}
	if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) == 0 {
		t.Fatalf("empty embeddings")
	}
	t.Logf("[openai default] texts=%d dim=%d", len(texts), len(res.Embeddings[0].Embedding))
	logHeadFloat(t, res.Embeddings[0].Embedding)
}

func TestEmbedderCombinations(t *testing.T) {
	type blReq struct {
		Model string `json:"model"`
		Input struct {
			Texts []string `json:"texts"`
		} `json:"input"`
		Parameters struct {
			Dimension      *int   `json:"dimension"`
			TextType       string `json:"text_type"`
			EncodingFormat string `json:"encoding_format"`
			OutputType     string `json:"output_type"`
			Instruct       string `json:"instruct"`
		} `json:"parameters"`
	}
	type blSparse struct {
		Index int     `json:"index"`
		Value float64 `json:"value"`
		Token string  `json:"token"`
	}
	type blRespEmb struct {
		SparseEmbedding []blSparse `json:"sparse_embedding,omitempty"`
		Embedding       []float64  `json:"embedding"`
		TextIndex       int        `json:"text_index"`
	}
	type blResp struct {
		StatusCode int    `json:"status_code"`
		RequestID  string `json:"request_id"`
		Code       string `json:"code"`
		Message    string `json:"message"`
		Output     struct {
			Embeddings []blRespEmb `json:"embeddings"`
		} `json:"output"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}
	type mockTransport struct {
		lastReq  *http.Request
		lastBody []byte
		build    func(*http.Request, []byte) ([]byte, int)
	}
	mt := &mockTransport{}
	mt.build = func(r *http.Request, b []byte) ([]byte, int) {
		var req blReq
		_ = json.Unmarshal(b, &req)
		dim := 1536
		if req.Parameters.Dimension != nil {
			dim = *req.Parameters.Dimension
		}
		dense := make([]float64, 0, dim)
		if req.Parameters.OutputType == "dense" || req.Parameters.OutputType == "dense&sparse" || req.Parameters.OutputType == "" {
			dense = make([]float64, dim)
			for i := 0; i < dim; i++ {
				dense[i] = float64(i%7) * 0.1
			}
		}
		var sparse []blSparse
		if req.Parameters.OutputType == "sparse" || req.Parameters.OutputType == "dense&sparse" {
			sparse = make([]blSparse, 8)
			for i := 0; i < 8; i++ {
				sparse[i] = blSparse{Index: i + 1, Value: float64(i) * 0.01, Token: fmt.Sprintf("t%d", i)}
			}
		}
		resp := blResp{
			StatusCode: http.StatusOK,
			RequestID:  "mock",
			Code:       "",
			Message:    "",
		}
		resp.Output.Embeddings = []blRespEmb{{
			SparseEmbedding: sparse,
			Embedding:       dense,
			TextIndex:       0,
		}}
		resp.Usage.TotalTokens = 1
		out, _ := json.Marshal(resp)
		return out, http.StatusOK
	}
	mtPtr := mt
	var tr http.RoundTripper = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		b, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(b))
		mtPtr.lastReq = req
		mtPtr.lastBody = b
		body, code := mtPtr.build(req, b)
		return &http.Response{
			StatusCode: code,
			Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	prev := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: tr}
	defer func() { http.DefaultClient = prev }()

	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	dims := []int{64, 128, 256, 512, 768, 1024, 1536, 2048}
	textTypes := []string{"document", "query"}
	outputTypes := []string{"dense", "sparse", "dense&sparse"}
	encFormats := []string{"float", "base64"}
	instructs := []string{"", "检索：返回和食物相关的内容"}
	baseURLs := []string{
		"",
		"https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
		"https://dashscope-intl.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
		"https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
	}
	for _, bu := range baseURLs {
		for _, d := range dims {
			for _, tt := range textTypes {
				for _, ot := range outputTypes {
					for _, ef := range encFormats {
						for _, ins := range instructs {
							name := fmt.Sprintf("bailian/%d/%s/%s/%s", d, tt, ot, ef)
							t.Run(name, func(t *testing.T) {
								dim := d
								ttv := tt
								otv := ot
								efv := ef
								insv := ins
								emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
									Provider:  consts.ModelProviderBaiLian,
									ModelName: "text-embedding-v4",
									BaseURL:   bu,
									APIKey:    "test",
									EmbedderParam: domain.EmbedderParam{
										Dimension:      &dim,
										TextType:       &ttv,
										OutputType:     &otv,
										EncodingFormat: &efv,
										Instruct:       &insv,
									},
								})
								if err != nil {
									t.Fatalf("GetEmbedder: %v", err)
								}
								texts := []string{"x"}
								res, err := mk.UseEmbedder(ctx, emb, texts)
								if err != nil {
									t.Fatalf("UseEmbedder: %v", err)
								}
								var req blReq
								_ = json.Unmarshal(mtPtr.lastBody, &req)
								wantEndpoint := bu
								if wantEndpoint == "" || strings.Contains(wantEndpoint, "/compatible-mode/") {
									if strings.Contains(wantEndpoint, "dashscope-intl.aliyuncs.com") {
										wantEndpoint = "https://dashscope-intl.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
									} else {
										wantEndpoint = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
									}
								}
								if strings.HasSuffix(wantEndpoint, "#") {
									wantEndpoint = strings.TrimSuffix(wantEndpoint, "#")
								}
								if mtPtr.lastReq.URL.String() != wantEndpoint {
									t.Fatalf("endpoint mismatch: got %s want %s", mtPtr.lastReq.URL.String(), wantEndpoint)
								}
								if req.Parameters.Dimension == nil || *req.Parameters.Dimension != dim {
									t.Fatalf("dimension mismatch: got %v want %d", req.Parameters.Dimension, dim)
								}
								if req.Parameters.TextType != ttv {
									t.Fatalf("text_type mismatch: got %s want %s", req.Parameters.TextType, ttv)
								}
								if req.Parameters.EncodingFormat != efv {
									t.Fatalf("encoding_format mismatch: got %s want %s", req.Parameters.EncodingFormat, efv)
								}
								if req.Parameters.OutputType != otv {
									t.Fatalf("output_type mismatch: got %s want %s", req.Parameters.OutputType, otv)
								}
								if req.Parameters.Instruct != insv {
									t.Fatalf("instruct mismatch: got %s want %s", req.Parameters.Instruct, insv)
								}
								if len(res.Embeddings) != len(texts) {
									t.Fatalf("embeddings length mismatch: got %d want %d", len(res.Embeddings), len(texts))
								}
								if res.Embeddings[0].TextIndex != 0 {
									t.Fatalf("text_index mismatch: got %d want %d", res.Embeddings[0].TextIndex, 0)
								}
								if otv == "dense" {
									if len(res.Embeddings) == 0 || len(res.Embeddings[0].Embedding) != dim {
										t.Fatalf("dense dim mismatch: got %d want %d", len(res.Embeddings[0].Embedding), dim)
									}
									if len(res.Embeddings[0].SparseEmbedding) != 0 {
										t.Fatalf("sparse should be empty for dense output")
									}
								}
								if otv == "sparse" {
									if len(res.Embeddings[0].SparseEmbedding) == 0 {
										t.Fatalf("sparse empty")
									}
									if len(res.Embeddings[0].Embedding) != 0 {
										t.Fatalf("dense should be empty for sparse output: got %d", len(res.Embeddings[0].Embedding))
									}
								}
								if otv == "dense&sparse" {
									if len(res.Embeddings[0].SparseEmbedding) == 0 {
										t.Fatalf("sparse empty")
									}
									if len(res.Embeddings[0].Embedding) != dim {
										t.Fatalf("dense dim mismatch: got %d want %d", len(res.Embeddings[0].Embedding), dim)
									}
								}
								if res.Usage.TotalTokens != 1 {
									t.Fatalf("usage.total_tokens mismatch: got %d want %d", res.Usage.TotalTokens, 1)
								}
							})
						}
					}
				}
			}
		}
	}

	t.Run("provider_matrix", func(t *testing.T) {
		cases := []struct {
			Provider  consts.ModelProvider
			ModelName string
			BaseURL   string
			APIKey    string
			WantErr   bool
		}{
			{consts.ModelProviderGemini, "embedding-001", "https://generativelanguage.googleapis.com", "", true},
			{consts.ModelProviderOllama, "nomic-embed-text", "http://localhost:11434", "", false},
			{consts.ModelProviderOllama, "nomic-embed-text", "http://localhost:11434/v1", "", false},
			{consts.ModelProviderOpenAI, "text-embedding-3-small", "https://api.openai.com/v1", "", false},
			{consts.ModelProviderAzureOpenAI, "text-embedding-3-small", "https://example.openai.azure.com", "", false},
			{consts.ModelProviderVolcengine, "bge-m3", "https://ark.example.com/v3", "", false},
		}
		for _, c := range cases {
			name := fmt.Sprintf("%s/%s", c.Provider, c.ModelName)
			t.Run(name, func(t *testing.T) {
				ctx := context.Background()
				mk := usecase.NewModelKit(nil)
				_, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  c.Provider,
					ModelName: c.ModelName,
					BaseURL:   c.BaseURL,
					APIKey:    c.APIKey,
				})
				if c.WantErr && err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !c.WantErr && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	})
}

func logHeadFloat(t *testing.T, v []float64) {
	n := 8
	if len(v) < n {
		n = len(v)
	}
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, fmt.Sprintf("%.4f", v[i]))
	}
	t.Logf("head: %s", strings.Join(s, " "))
}

func logHeadStr(t *testing.T, v []string) {
	n := 8
	if len(v) < n {
		n = len(v)
	}
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		it := v[i]
		if len(it) > 16 {
			it = it[:16]
		}
		s = append(s, it)
	}
	t.Logf("head(base64): %s", strings.Join(s, " "))
}

func logSparseEntriesHead(t *testing.T, se []domain.SparseEmbedding) {
	n := 8
	if len(se) < n {
		n = len(se)
	}
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, fmt.Sprintf("%d:%.4f:%s", se[i].Index, se[i].Value, se[i].Token))
	}
	t.Logf("sparse_head: %s", strings.Join(parts, " "))
}
