# ModelKit Embedder 介绍

技术

- 封装 CloudWeGo Eino 组件（`embedding/openai`、`embedding/ollama`），提供统一 Embedder 接口
- 内置 DashScope 支持（`components/embedder/bailian`），兼容 `text-embedding-v3/v4`
- 集成火山引擎 Ark Embedding（`eino-ext/components/embedding/ark`）

功能

- 支持 `OpenAI API` 与 `DashScope API`
- 支持生成 `稠密向量`、`稀疏向量`、`稠密+稀疏`
- 支持的供应商：
  - OpenAI API 兼容：`OpenAI`、`AzureOpenAI`（兼容模式）、`Ollama`（`/v1` 兼容模式）
  - 原生：`BaiLian (DashScope)`、`Ollama`（原生）、`Volcengine (Ark)`

# 创建embedder

```go
package main

import (
    "context"
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

    embedder,err := mk.GetEmbedder(ctx, &domain.ModelMetadata{
        Provider:  consts.ModelProviderBaiLian,
        ModelName: "text-embedding-v4",
        BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
        APIKey:    "sk-xxxxxx",
    })

    if err != nil {
        logger.Error("get embedder failed", "err", err)
    }


}

```

字段说明（ModelMetadata）：

- `provider`：模型提供商，取值如 `BaiLian`、`OpenAI`、`AzureOpenAI`、`Ollama`、`Volcengine`。
- `model_name`：向量模型 ID，例如 `text-embedding-v4`、`bge-m3`。
- `base_url`：Modelkit会自动添加 `/embeddings` 路径, 如果base_url以`#`结尾, 可以强制使用输入的base_url。
- `api_key`：鉴权密钥，作为 `Authorization: Bearer <API_KEY>` 使用。
- `api_version`：仅 `AzureOpenAI` 需要，未设置将默认 `2024-10-21`。
- `api_header`：可选的自定义请求头（`key=value` 按行拼接），用于兼容某些平台的鉴权方式。
- `dimension`：向量维度，可选值：`2048(仅v4)`、`1536(仅v4)`、`1024`、`768`、`512`、`256`、`128`、`64`。
- `text_type`：`document` 或 `query`；检索任务建议区分 `query/document`。
- `output_type`：`dense`、`sparse`、`dense&sparse`（v3/v4 支持）。
- `encoding_format`：`float` 
- `instruct`：检索指令，仅在 `text-embedding-v4` 且 `text_type=query` 时生效。

# 使用embedder

```go
texts := []string{"示例文本一", "示例文本二"}
res, _ := mk.UseEmbedder(ctx, embedder, texts)
```

## 生成稠密向量

```go
ot := "dense"
embedder, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
    Provider:  consts.ModelProviderBaiLian,
    ModelName: "text-embedding-v4",
    BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
    APIKey:    "sk-xxxxxx",
    EmbedderParam: domain.EmbedderParam{
        OutputType: &ot,
    },
})
texts := []string{"示例文本一", "示例文本二"}
res, _ := mk.UseEmbedder(ctx, embedder, texts)
```



示例结果

```json
{
  "embeddings": [
    {
      "text_index": 0,
      "embedding": [0.0123, -0.0456, 0.0789, 0.0042, -0.0178, 0.0321, -0.0567, 0.0890]
    },
    {
      "text_index": 1,
      "embedding": [0.0234, -0.0567, 0.0890, 0.0055, -0.0201, 0.0288, -0.0602, 0.0813]
    }
  ],
  "usage": { "total_tokens": 48 }
}
```

## 生成稀疏向量

```go
ot := "sparse"
embedder, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
    Provider:  consts.ModelProviderBaiLian,
    ModelName: "text-embedding-v4",
    BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
    APIKey:    "sk-xxxxxx",
    EmbedderParam: domain.EmbedderParam{
        OutputType: &ot,
    },
})
texts := []string{"示例文本一", "示例文本二"}
res, _ := mk.UseEmbedder(ctx, embedder, texts)
```

使用条件

- 仅 `text-embedding-v3`/`text-embedding-v4` 支持设置 `output_type` 生成稀疏向量
- 将 `output_type` 设为 `sparse` 或 `dense&sparse` 才会返回 `sparse_embedding`
- `text_type` 设 `query` , 同时设置 `instruct` 时 无法生成稀疏向量

示例结果

```json
{
  "embeddings": [
    {
      "text_index": 0,
      "sparse_embedding": [
        {"index": 123, "value": 0.0321, "token": "示例"},
        {"index": 456, "value": 0.0289, "token": "文本"},
        {"index": 789, "value": 0.0203, "token": "一"}
      ]
    },
    {
      "text_index": 1,
      "sparse_embedding": [
        {"index": 111, "value": 0.0400, "token": "示例"},
        {"index": 222, "value": 0.0312, "token": "文本"}
      ]
    }
  ],
  "usage": { "total_tokens": 52 }
}
```

## 生成稠密+稀疏向量

```go
ot := "dense&sparse"
embedder, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
    Provider:  consts.ModelProviderBaiLian,
    ModelName: "text-embedding-v4",
    BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
    APIKey:    "sk-xxxxxx",
    EmbedderParam: domain.EmbedderParam{
        OutputType: &ot,
    },
})
texts := []string{"示例文本一", "示例文本二"}
res, _ := mk.UseEmbedder(ctx, embedder, texts)
```

示例结果

```json
{
  "embeddings": [
    {
      "text_index": 0,
      "embedding": [0.0123, -0.0456, 0.0789, 0.0042, -0.0178, 0.0321, -0.0567, 0.0890],
      "sparse_embedding": [
        {"index": 120, "value": 0.0310, "token": "示例"},
        {"index": 451, "value": 0.0275, "token": "文本"}
      ]
    },
    {
      "text_index": 1,
      "embedding": [0.0234, -0.0567, 0.0890, 0.0055, -0.0201, 0.0288, -0.0602, 0.0813],
      "sparse_embedding": [
        {"index": 210, "value": 0.0399, "token": "示例"},
        {"index": 325, "value": 0.0305, "token": "文本"}
      ]
    }
  ],
  "usage": { "total_tokens": 50 }
}
```
