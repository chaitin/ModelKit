# ModelKit Reranker 介绍

支持BGE与Qwen重排序模型

# 创建reranker

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/chaitin/ModelKit/v2/consts"
    "github.com/chaitin/ModelKit/v2/domain"
    "github.com/chaitin/ModelKit/v2/usecase"
)

func main() {
    ctx := context.Background()
    // 创建logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
        Level:     slog.LevelInfo,
    }))
    mk := usecase.NewModelKit(logger)

    rk, err := mk.GetReranker(ctx, &domain.ModelMetadata{
        Provider:  consts.ModelProviderBaiLian,
        ModelName: "qwen3-rerank",
        BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#",
        APIKey:    "sk-xxxxxx",
    })

    if err != nil {
        logger.Error("get reranker failed", "err", err)
    }
}
```

字段说明（ModelMetadata）：

- `provider`：模型提供商，取值如 `BaiLian`、`BaiZhiCloud`、`Other`（OpenAI-style）。
- `model_name`：重排模型 ID，例如 `qwen3-rerank`、`bge-reranker-v2-m3`。
- `base_url`：通用模式下会自动添加 `/rerank` 路径；若 `base_url` 以 `#` 结尾，则强制使用输入的完整地址（适用于百炼的 `text-rerank` 路径）。
- `api_key`：鉴权密钥，作为 `Authorization: Bearer <API_KEY>` 使用。
- `api_header`：可选的自定义请求头（`key=value` 按行拼接），用于兼容某些平台的鉴权方式。

# 使用reranker

```go
docs := []string{"示例文档一", "示例文档二", "示例文档三"}
N := 2
res, _ := rk.Rerank(ctx, domain.RerankRequest{
    Documents:       docs,
    Query:           "查询词",
    ReturnDocuments: true,
    N:               &N,
})
```

示例结果：

```json
{
  "results": [
    { "index": 1, "relevance_score": 0.92, "document": "示例文档二" },
    { "index": 0, "relevance_score": 0.75, "document": "示例文档一" }
  ],
  "usage": {
    "prompt_tokens": 20,
    "input_tokens": 54,
    "output_tokens": 0,
    "total_tokens": 74
  }
}
```

参数说明（RerankRequest）：

- `documents`：待重排的文档字符串数组。
- `query`：查询词，用于计算与各文档的相关性。
- `return_documents`：是否在结果中返回文档内容文本，默认 false。
- `n`：返回前 N 条结果；未设置则返回全部；N≤0 视为 1；N 大于文档数量时按文档数量截断。

结果说明：

- `results`：按相关性从高到低排序的重排结果，包含原始 `index`、`relevance_score`，以及在 `return_documents=true` 时的 `document` 文本。
- `usage`：可选的计费/令牌使用信息，不同供应商字段可能有所不同。
