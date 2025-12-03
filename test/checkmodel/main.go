package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func loadEnv() {
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
}

func checkEmbedding(ctx context.Context, mk *usecase.ModelKit, provider, model, baseURL, apiKey string) {
	req := &domain.CheckModelReq{
		Provider: provider,
		Model:    model,
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Type:     string(consts.ModelTypeEmbedding),
	}
	resp, err := mk.CheckModel(ctx, req)
	if err != nil {
		fmt.Printf("[embedding] provider=%s model=%s error=%s\n", provider, model, err.Error())
		return
	}
	if resp.Error != "" {
		fmt.Printf("[embedding] provider=%s model=%s error=%s\n", provider, model, resp.Error)
		return
	}
	fmt.Printf("[embedding] provider=%s model=%s content=%s\n", provider, model, resp.Content)
}

func checkRerank(ctx context.Context, mk *usecase.ModelKit, provider, model, baseURL, apiKey string) {
	req := &domain.CheckModelReq{
		Provider: provider,
		Model:    model,
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Type:     string(consts.ModelTypeRerank),
	}
	resp, err := mk.CheckModel(ctx, req)
	if err != nil {
		fmt.Printf("[rerank] provider=%s model=%s error=%s\n", provider, model, err.Error())
		return
	}
	if resp.Error != "" {
		fmt.Printf("[rerank] provider=%s model=%s error=%s\n", provider, model, resp.Error)
		return
	}
	fmt.Printf("[rerank] provider=%s model=%s content=\n%s\n", provider, model, resp.Content)
}

func main() {
	loadEnv()
	ctx := context.Background()
	mk := usecase.NewModelKit(nil)

	checkEmbedding(ctx, mk,
		string(consts.ModelProviderOpenAI),
		"bge-m3",
		"https://model-square.app.baizhi.cloud/v1",
		os.Getenv("baizhiapikey"),
	)
	checkEmbedding(ctx, mk,
		string(consts.ModelProviderBaiLian),
		"text-embedding-v4",
		"https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		os.Getenv("bailianapikey"),
	)

	checkRerank(ctx, mk,
		string(consts.ModelProviderBaiZhiCloud),
		"bge-m3",
		"https://model-square.app.baizhi.cloud/v1/rerank",
		os.Getenv("baizhiapikey"),
	)
	checkRerank(ctx, mk,
		string(consts.ModelProviderBaiLian),
		"qwen3-rerank",
		"https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#",
		os.Getenv("bailianapikey"),
	)
}
