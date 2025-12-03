package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func baaiTest() {
	ctx := context.Background()
	modelkit := usecase.NewModelKit(nil)
	reranker, err := modelkit.GetReranker(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiZhiCloud,
		ModelName: "bge-m3",
		BaseURL:   "https://model-square.app.baizhi.cloud/v1/rerank",
		APIKey:    os.Getenv("baizhiapikey"),
	})
	if err != nil {
		log.Fatalf("NewEmbedder failed: %v", err)
	}
	topN := 3
	req := domain.RerankRequest{
		N:               &topN,
		Documents:       []string{"火鸡面", "火鸡", "面", "鸡", "火"},
		Query:           "食物",
		ReturnDocuments: false,
	}
	rerankedDocs, err := reranker.Rerank(ctx, req)
	if err != nil {
		log.Fatalf("EmbedStrings failed: %v", err)
	}
	for _, doc := range rerankedDocs.Results {
		fmt.Printf("Index: %d, Score: %.4f\n", doc.Index, doc.RelevanceScore)
		fmt.Println(doc.Document)
	}

}

func bailianTest() {
	ctx := context.Background()
	modelkit := usecase.NewModelKit(nil)
	reranker, err := modelkit.GetReranker(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiLian,
		ModelName: "qwen3-rerank",
		BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#",
		APIKey:    os.Getenv("bailianapikey"),
	})
	if err != nil {
		log.Fatalf("NewEmbedder failed: %v", err)
	}
	topN := 99
	req := domain.RerankRequest{
		N:               &topN,
		Documents:       []string{"火鸡面", "火鸡", "面", "鸡", "火"},
		Query:           "动物",
		ReturnDocuments: false,
	}
	rerankedDocs, err := reranker.Rerank(ctx, req)
	if err != nil {
		log.Fatalf("EmbedStrings failed: %v", err)
	}
	for _, doc := range rerankedDocs.Results {
		fmt.Printf("Index: %d, Score: %.4f\n", doc.Index, doc.RelevanceScore)
		fmt.Println(doc.Document)
	}

}

func main() {
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
	baaiTest()
	bailianTest()
}
