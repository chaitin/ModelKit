package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func TestRerankerCombinations(t *testing.T) {
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	t.Run("provider_matrix", func(t *testing.T) {
		apiKeyBL := os.Getenv("bailianapikey")
		apiKeyBZ := os.Getenv("baizhiapikey")

		cases := []struct {
			Provider  consts.ModelProvider
			ModelName string
			BaseURLs  []string
			APIKey    string
			Docs      []string
		}{
			{consts.ModelProviderBaiLian, "qwen3-rerank", []string{"", "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank", "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#"}, apiKeyBL, []string{"a", "b", "c", "d", "e"}},
			{consts.ModelProviderBaiZhiCloud, "bge-m3", []string{"https://model-square.app.baizhi.cloud/v1", "https://model-square.app.baizhi.cloud/v1/rerank"}, apiKeyBZ, []string{"a", "b", "c", "d"}},
		}

		for _, c := range cases {
			if c.APIKey == "" {
				if c.Provider == consts.ModelProviderBaiLian {
					t.Skip("missing bailianapikey")
				} else {
					t.Skip("missing baizhiapikey")
				}
			}
			for _, bu := range c.BaseURLs {
				for _, n := range []int{0, 3, 9} {
					for _, rd := range []bool{false, true} {
						name := fmt.Sprintf("provider=%s model=%s base_url=%s n=%d return_documents=%t", c.Provider, c.ModelName, bu, n, rd)
						t.Run(name, func(t *testing.T) {
							nn := n
							rdv := rd
							rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
								Provider:  c.Provider,
								ModelName: c.ModelName,
								BaseURL:   bu,
								APIKey:    c.APIKey,
							})
							if err != nil {
								t.Logf("test case: %s", name)
								t.Fatalf("GetReranker: %v", err)
							}
							res, err := rk.Rerank(ctx, domain.RerankRequest{N: &nn, Documents: c.Docs, Query: "q", ReturnDocuments: rdv})
							if err != nil {
								t.Logf("test case: %s", name)
								t.Fatalf("Rerank: %v", err)
							}
							want := nn
							if want <= 0 {
								want = 1
							}
							if want > len(c.Docs) {
								want = len(c.Docs)
							}
							if len(res.Results) != want {
								t.Logf("test case: %s", name)
								t.Fatalf("results length mismatch: got %d want %d", len(res.Results), want)
							}
							if rdv {
								if res.Results[0].Document == "" {
									t.Logf("test case: %s", name)
									t.Fatalf("document empty when return_documents=true")
								}
							}
							t.Logf("test case: %s", name)
							for i := range res.Results {
								t.Logf("result_index=%d score=%.6f doc=%s", res.Results[i].Index, res.Results[i].RelevanceScore, res.Results[i].Document)
							}
							if res.Usage != nil {
								t.Logf("usage_prompt_tokens=%d total_tokens=%d input_tokens=%d output_tokens=%d", res.Usage.PromptTokens, res.Usage.TotalTokens, res.Usage.InputTokens, res.Usage.OutputTokens)
							}
						})
					}
				}
			}
		}
	})
}
