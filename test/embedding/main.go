package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
)

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
	// bailianTest()
	openaiTest()
}

func bailianTest() {
	ctx := context.Background()
	modelkit := usecase.NewModelKit(nil)

	bailianDefaultTest(ctx, modelkit)
	bailianDimensionTest(ctx, modelkit)
	bailianTextTypeTest(ctx, modelkit)
	bailianOutputTypeTest(ctx, modelkit)
	bailianEncodingFormatTest(ctx, modelkit)
	bailianInstructTest(ctx, modelkit)
}

func bailianDefaultTest(ctx context.Context, modelkit *usecase.ModelKit) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiLian,
		ModelName: "text-embedding-v4",
		BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:    os.Getenv("bailianapikey"),
	})
	if err != nil {
		log.Fatalf("NewEmbedder failed: %v", err)
	}
	texts := []string{"火鸡面", "测试文本", "向量模型"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Fatalf("UseEmbedder failed: %v", err)
	}
	fmt.Printf("[bailian default] texts=%d dim=%d\n", len(texts), len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func bailianDimensionTest(ctx context.Context, modelkit *usecase.ModelKit) {
	for _, d := range []int{64, 128, 256, 512, 768, 1024, 1536, 2048} {
		bailianDimensionOption(ctx, modelkit, d)
	}
}

func bailianTextTypeTest(ctx context.Context, modelkit *usecase.ModelKit) {
	bailianTextTypeOption(ctx, modelkit, "document")
	bailianTextTypeOption(ctx, modelkit, "query")
}

func bailianOutputTypeTest(ctx context.Context, modelkit *usecase.ModelKit) {
	bailianOutputTypeOption(ctx, modelkit, "dense")
	bailianOutputTypeOption(ctx, modelkit, "sparse")
	bailianOutputTypeOption(ctx, modelkit, "dense&sparse")
}

func bailianEncodingFormatTest(ctx context.Context, modelkit *usecase.ModelKit) {
	bailianEncodingFormatOption(ctx, modelkit, "float")
	bailianEncodingFormatOption(ctx, modelkit, "base64")
}

func bailianInstructTest(ctx context.Context, modelkit *usecase.ModelKit) {
	bailianInstructOption(ctx, modelkit, "query", "检索：返回和食物相关的内容")
	bailianInstructOption(ctx, modelkit, "query", "检索：返回和动物相关的内容")
}

func openaiTest() {
	ctx := context.Background()
	modelkit := usecase.NewModelKit(nil)
	openaiDefaultTest(ctx, modelkit)
}

func openaiDefaultTest(ctx context.Context, modelkit *usecase.ModelKit) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderOpenAI,
		ModelName: "bge-m3",
		BaseURL:   "https://model-square.app.baizhi.cloud/v1",
		APIKey:    os.Getenv("baizhiapikey"),
	})
	if err != nil {
		log.Fatalf("NewEmbedder failed: %v", err)
	}
	texts := []string{"风哀", "渚回"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Fatalf("UseEmbedder failed: %v", err)
	}
	fmt.Printf("[openai default] texts=%d dim=%d\n", len(texts), len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func bailianDimensionOption(ctx context.Context, modelkit *usecase.ModelKit, dim int) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiLian,
		ModelName: "text-embedding-v4",
		BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:    os.Getenv("bailianapikey"),
		Dimension: &dim,
	})
	if err != nil {
		log.Printf("NewEmbedder failed: %v", err)
		return
	}
	texts := []string{"火鸡面"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Printf("UseEmbedder failed: %v", err)
		return
	}
	fmt.Printf("[bailian dimension=%d] dim=%d\n", dim, len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func bailianTextTypeOption(ctx context.Context, modelkit *usecase.ModelKit, tt string) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiLian,
		ModelName: "text-embedding-v4",
		BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:    os.Getenv("bailianapikey"),
		TextType:  &tt,
	})
	if err != nil {
		log.Printf("NewEmbedder failed: %v", err)
		return
	}
	texts := []string{"火鸡面"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Printf("UseEmbedder failed: %v", err)
		return
	}
	fmt.Printf("[bailian text_type=%s] dim=%d\n", tt, len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func bailianOutputTypeOption(ctx context.Context, modelkit *usecase.ModelKit, ot string) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:   consts.ModelProviderBaiLian,
		ModelName:  "text-embedding-v4",
		BaseURL:    "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:     os.Getenv("bailianapikey"),
		OutputType: &ot,
	})
	if err != nil {
		log.Printf("NewEmbedder failed: %v", err)
		return
	}
	texts := []string{"火鸡面"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Printf("UseEmbedder failed: %v", err)
		return
	}
	fmt.Printf("[bailian output_type=%s] dim=%d\n", ot, len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
	if ot != "dense" {
		se := res.Output.Embeddings[0].SparseEmbedding
		fmt.Printf("[bailian output_type=%s] sparse_nnz=%d\n", ot, len(se))
		printSparseEntriesHead(se)
	}
}

func bailianEncodingFormatOption(ctx context.Context, modelkit *usecase.ModelKit, ef string) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:       consts.ModelProviderBaiLian,
		ModelName:      "text-embedding-v4",
		BaseURL:        "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:         os.Getenv("bailianapikey"),
		EncodingFormat: &ef,
	})
	if err != nil {
		log.Printf("NewEmbedder failed: %v", err)
		return
	}
	texts := []string{"火鸡面"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Printf("UseEmbedder failed: %v", err)
		return
	}
	fmt.Printf("[bailian encoding_format=%s] dim=%d\n", ef, len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func bailianInstructOption(ctx context.Context, modelkit *usecase.ModelKit, tt, instr string) {
	embedder, err := modelkit.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  consts.ModelProviderBaiLian,
		ModelName: "text-embedding-v4",
		BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
		APIKey:    os.Getenv("bailianapikey"),
		TextType:  &tt,
		Instruct:  &instr,
	})
	if err != nil {
		log.Printf("NewEmbedder failed: %v", err)
		return
	}
	texts := []string{"火鸡面"}
	res, err := modelkit.UseEmbedder(ctx, embedder, texts)
	if err != nil {
		log.Printf("UseEmbedder failed: %v", err)
		return
	}
	fmt.Printf("[bailian instruct=%s text_type=%s] dim=%d\n", instr, tt, len(res.Output.Embeddings[0].Embedding))
	printHead(res.Output.Embeddings[0].Embedding)
}

func printHead(v []float64) {
	n := 8
	if len(v) < n {
		n = len(v)
	}
	fmt.Printf("head: ")
	for i := 0; i < n; i++ {
		fmt.Printf("%.4f ", v[i])
	}
	fmt.Println()
}

func printSparseEntriesHead(se []domain.SparseEntry) {
	n := 8
	if len(se) < n {
		n = len(se)
	}
	fmt.Printf("sparse_head: ")
	for i := 0; i < n; i++ {
		fmt.Printf("%d:%.4f:%s ", se[i].Index, se[i].Value, se[i].Token)
	}
	fmt.Println()
}

func bailianFetchSparse(ctx context.Context, ot string) ([]int, []float64, error) {
	body := map[string]any{
		"model": "text-embedding-v4",
		"input": map[string]any{
			"texts": []string{"火鸡面"},
		},
		"parameters": map[string]any{
			"output_type": ot,
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}
	url := "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("bailianapikey"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf(resp.Status)
	}
	var ar struct {
		Output struct {
			Embeddings []struct {
				Embedding       []float64       `json:"embedding"`
				SparseEmbedding json.RawMessage `json:"sparse_embedding"`
			} `json:"embeddings"`
		} `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, nil, err
	}
	if len(ar.Output.Embeddings) == 0 {
		return nil, nil, fmt.Errorf("empty embeddings")
	}
	var obj struct {
		Indices []int     `json:"indices"`
		Values  []float64 `json:"values"`
	}
	if err := json.Unmarshal(ar.Output.Embeddings[0].SparseEmbedding, &obj); err == nil {
		return obj.Indices, obj.Values, nil
	}
	var arr1 []struct {
		Index int     `json:"index"`
		Value float64 `json:"value"`
	}
	if err := json.Unmarshal(ar.Output.Embeddings[0].SparseEmbedding, &arr1); err == nil {
		idx := make([]int, 0, len(arr1))
		vals := make([]float64, 0, len(arr1))
		for _, it := range arr1 {
			idx = append(idx, it.Index)
			vals = append(vals, it.Value)
		}
		return idx, vals, nil
	}
	var arr2 [][]float64
	if err := json.Unmarshal(ar.Output.Embeddings[0].SparseEmbedding, &arr2); err == nil {
		idx := make([]int, 0, len(arr2))
		vals := make([]float64, 0, len(arr2))
		for _, it := range arr2 {
			if len(it) >= 2 {
				idx = append(idx, int(it[0]))
				vals = append(vals, it[1])
			}
		}
		return idx, vals, nil
	}
	return nil, nil, fmt.Errorf("unsupported sparse format")
}

func printSparseHead(idx []int, vals []float64) {
	n := 8
	if len(vals) < n {
		n = len(vals)
	}
	fmt.Printf("sparse_head: ")
	for i := 0; i < n; i++ {
		fmt.Printf("%d:%.4f ", idx[i], vals[i])
	}
	fmt.Println()
}
