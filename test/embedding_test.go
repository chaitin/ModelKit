package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func TestMain(m *testing.M) {
	f, err := os.Open(".env")
	if err == nil {
		defer func() { _ = f.Close() }()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if k, v, ok := strings.Cut(line, "="); ok {
				if err := os.Setenv(strings.TrimSpace(k), strings.TrimSpace(v)); err != nil {
					fmt.Fprintf(os.Stderr, "setenv error: %v\n", err)
				}
			}
		}
	}
	os.Exit(m.Run())
}

func TestEmbedderCombinations(t *testing.T) {
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	dims := []int{64, 512, 1536}
	textTypes := []string{"document", "query"}
	outputTypes := []string{"dense", "sparse"}
	encFormats := []string{"float"}
	instructs := []string{"", "检索：返回和食物相关的内容"}
	baseURLs := []string{
		"https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding",
		"https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
	}
	apiKey := os.Getenv("bailianapikey")
	if apiKey == "" {
		t.Fatalf("missing bailianapikey")
	}
	apiKeyBZ := strings.TrimSpace(os.Getenv("baizhiapikey"))
	if apiKeyBZ == "" {
		t.Fatalf("missing baizhiapikey")
	}
	bzBaseURLs := []string{
		"https://model-square.app.baizhi.cloud/v1",
		"https://model-square.app.baizhi.cloud/v1#",
	}
	cases := []struct {
		Provider  consts.ModelProvider
		ModelName string
		BaseURLs  []string
		APIKey    string
	}{
		{consts.ModelProviderBaiLian, "text-embedding-v4", baseURLs, apiKey},
		{consts.ModelProviderBaiZhiCloud, "bge-m3", bzBaseURLs, apiKeyBZ},
	}
	for _, c := range cases {
		for _, bu := range c.BaseURLs {
			for _, d := range dims {
				for _, tt := range textTypes {
					otList := outputTypes
					if c.Provider == consts.ModelProviderBaiZhiCloud {
						otList = []string{"dense"}
					}
					for _, ot := range otList {
						for _, ef := range encFormats {
							instrList := instructs
							if c.Provider == consts.ModelProviderBaiZhiCloud {
								instrList = []string{""}
							}
							for _, ins := range instrList {
								name := fmt.Sprintf("provider=%s model=%s base_url=%s dimension=%d text_type=%s output_type=%s encoding_format=%s instruct=%s", c.Provider, c.ModelName, bu, d, tt, ot, ef, ins)
								t.Run(name, func(t *testing.T) {
									dim := d
									ttv := tt
									otv := ot
									efv := ef
									insv := ins
									md := &domain.ModelMetadata{
										Provider:  c.Provider,
										ModelName: c.ModelName,
										BaseURL:   bu,
										APIKey:    c.APIKey,
									}
									if c.Provider == consts.ModelProviderBaiLian {
										md.EmbedderParam = domain.EmbedderParam{
											Dimension:      &dim,
											TextType:       &ttv,
											OutputType:     &otv,
											EncodingFormat: &efv,
											Instruct:       &insv,
										}
									}
									emb, err := mk.GetEmbedder(ctx, md)
									if err != nil {
										t.Logf("test case: %s", name)
										t.Fatalf("GetEmbedder: %v", err)
									}
									texts := []string{"x"}
									res, err := mk.UseEmbedder(ctx, emb, texts)
									if c.Provider == consts.ModelProviderBaiZhiCloud && strings.HasSuffix(bu, "#") {
										if err == nil {
											t.Logf("test case: %s", name)
											t.Fatalf("expected error for base_url ending with #: %s", bu)
										}
										t.Logf("test case: %s", name)
										t.Logf("expected error: %v", err)
										return
									}
									if err != nil {
										t.Logf("test case: %s", name)
										t.Fatalf("UseEmbedder: %v", err)
									}
									if len(res.Embeddings) != len(texts) {
										t.Logf("test case: %s", name)
										t.Fatalf("embeddings length mismatch: got %d want %d", len(res.Embeddings), len(texts))
									}
									if res.Embeddings[0].TextIndex != 0 {
										t.Logf("test case: %s", name)
										t.Fatalf("text_index mismatch: got %d want %d", res.Embeddings[0].TextIndex, 0)
									}
									if otv == "dense" {
										if len(res.Embeddings) == 0 {
											t.Logf("test case: %s", name)
											t.Fatalf("empty embeddings")
										}
										if c.Provider == consts.ModelProviderBaiLian {
											if len(res.Embeddings[0].Embedding) != dim {
												t.Logf("test case: %s", name)
												t.Fatalf("dense dim mismatch: got %d want %d", len(res.Embeddings[0].Embedding), dim)
											}
										}
										if len(res.Embeddings[0].SparseEmbedding) != 0 {
											t.Logf("test case: %s", name)
											t.Fatalf("sparse should be empty for dense output")
										}
									}
									if otv == "sparse" {
										if strings.ToLower(ttv) == "query" && insv != "" {
											if len(res.Embeddings[0].SparseEmbedding) != 0 {
												t.Logf("test case: %s", name)
												t.Fatalf("sparse should be empty for query+instruct")
											}
										} else {
											if len(res.Embeddings[0].SparseEmbedding) == 0 {
												t.Logf("test case: %s", name)
												t.Fatalf("sparse not available for this combination")
											}
										}
									}
									if res.Usage.TotalTokens < 0 {
										t.Logf("test case: %s", name)
										t.Fatalf("usage.total_tokens invalid: %d", res.Usage.TotalTokens)
									}
									t.Logf("test case: %s", name)
									for i := range res.Embeddings {
										t.Logf("embedding_index=%d dim=%d", i, len(res.Embeddings[i].Embedding))
										logFullFloat(t, res.Embeddings[i].Embedding)
										if len(res.Embeddings[i].SparseEmbedding) > 0 {
											logSparseEntriesFull(t, res.Embeddings[i].SparseEmbedding)
										}
									}
									t.Logf("usage_total_tokens=%d", res.Usage.TotalTokens)
								})
							}
						}
					}
				}
			}
		}
	}

	t.Run("baizhicloud_bge_m3", func(t *testing.T) {
		ctx := context.Background()
		mk := usecase.NewModelKit(nil)
		baseURLs := []string{
			"https://model-square.app.baizhi.cloud/v1",
			"https://model-square.app.baizhi.cloud/v1#",
		}
		apiKey := strings.TrimSpace(os.Getenv("baizhiapikey"))
		if apiKey == "" {
			t.Fatalf("missing baizhiapikey")
		}
		texts := []string{"x"}
		for _, bu := range baseURLs {
			name := fmt.Sprintf("provider=%s model=%s base_url=%s", consts.ModelProviderBaiZhiCloud, "bge-m3", bu)
			t.Run(name, func(t *testing.T) {
				emb, err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
					Provider:  consts.ModelProviderBaiZhiCloud,
					ModelName: "bge-m3",
					BaseURL:   bu,
					APIKey:    apiKey,
				})
				if err != nil {
					t.Logf("test case: %s", name)
					t.Fatalf("GetEmbedder: %v", err)
				}
				res, err := mk.UseEmbedder(ctx, emb, texts)
				if strings.HasSuffix(bu, "#") {
					if err == nil {
						t.Logf("test case: %s", name)
						t.Fatalf("expected error for base_url ending with #: %s", bu)
					}
					t.Logf("test case: %s", name)
					t.Logf("expected error: %v", err)
					return
				}
				if err != nil {
					t.Logf("test case: %s", name)
					t.Fatalf("UseEmbedder: %v", err)
				}
				if len(res.Embeddings) != len(texts) {
					t.Logf("test case: %s", name)
					t.Fatalf("embeddings length mismatch: got %d want %d", len(res.Embeddings), len(texts))
				}
				if res.Embeddings[0].TextIndex != 0 {
					t.Logf("test case: %s", name)
					t.Fatalf("text_index mismatch: got %d want %d", res.Embeddings[0].TextIndex, 0)
				}
				if len(res.Embeddings[0].Embedding) == 0 {
					t.Logf("test case: %s", name)
					t.Fatalf("empty embedding")
				}
				if res.Usage.TotalTokens < 0 {
					t.Logf("test case: %s", name)
					t.Fatalf("usage.total_tokens invalid: %d", res.Usage.TotalTokens)
				}
				t.Logf("test case: %s", name)
				for i := range res.Embeddings {
					t.Logf("embedding_index=%d dim=%d", i, len(res.Embeddings[i].Embedding))
					logFullFloat(t, res.Embeddings[i].Embedding)
					if len(res.Embeddings[i].SparseEmbedding) > 0 {
						logSparseEntriesFull(t, res.Embeddings[i].SparseEmbedding)
					}
				}
				t.Logf("usage_total_tokens=%d", res.Usage.TotalTokens)
			})
		}
	})
}

func logFullFloat(t *testing.T, v []float64) {
	s := make([]string, 0, len(v))
	for i := 0; i < len(v); i++ {
		s = append(s, fmt.Sprintf("%.6f", v[i]))
	}
	t.Logf("embedding: %s", strings.Join(s, " "))
}

func logSparseEntriesFull(t *testing.T, se []domain.SparseEmbedding) {
	parts := make([]string, 0, len(se))
	for i := 0; i < len(se); i++ {
		parts = append(parts, fmt.Sprintf("%d:%.6f:%s", se[i].Index, se[i].Value, se[i].Token))
	}
	t.Logf("sparse: %s", strings.Join(parts, " "))
}
