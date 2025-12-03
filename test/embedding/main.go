package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

func main() {
	ctx := context.Background()
	modelkit := usecase.NewModelKit(nil)
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderOpenAI,
		ModelName: "bge-m3",
		BaseURL:   "https://model-square.app.baizhi.cloud/v1",
		APIKey:    "wrPLDr96ojOkBJNFjy3L3iA1ybco54nwsAQVu63Av0rKAz4N",
	})
	if err != nil {
		log.Fatalf("NewEmbedder failed: %v", err)
	}
	embeddings, err := embedder.EmbedStrings(ctx, []string{
		"第一段文本",
		"第二段文本",
	})
	if err != nil {
		log.Fatalf("EmbedStrings failed: %v", err)
	}
	fmt.Println(embeddings)
}
