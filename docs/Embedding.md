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
    mk := usecase.NewModelKit(nil)

    // 阿里云百炼（DashScope）
    emb1, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
        Provider:  consts.ModelProviderBaiLian,
        ModelName: "text-embedding-v4",
        BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
        APIKey:    os.Getenv("bailianapikey"),
    })

    // OpenAI 兼容（例如百智云 ModelSquare）
    emb2, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
        Provider:  consts.ModelProviderOpenAI,
        ModelName: "bge-m3",
        BaseURL:   "https://model-square.app.baizhi.cloud/v1",
        APIKey:    os.Getenv("baizhiapikey"),
    })

    _ = emb1
    _ = emb2
}
```

# 使用embedder

```go
ctx := context.Background()
mk := usecase.NewModelKit(nil)

emb, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
    Provider:  consts.ModelProviderBaiLian,
    ModelName: "text-embedding-v4",
    BaseURL:   "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
    APIKey:    os.Getenv("bailianapikey"),
})

texts := []string{"示例文本一", "示例文本二"}
res, _ := mk.UseEmbedder(ctx, emb, texts)

// res.Embeddings[i].Embedding 为 dense 向量
// res.Embeddings[i].SparseEmbedding（如返回）为 sparse 向量条目
```

## 支持的供应商

- `BaiLian`（阿里云百炼 DashScope）
- `OpenAI` 兼容（原生 OpenAI 及兼容平台）
- `AzureOpenAI`
- `Ollama`（本地或远端，支持 OpenAI 兼容与原生两种 URL）
- `Volcengine`（火山引擎方舟 Ark）

## 请求体与响应体

- 请求体：

```json
{
  "model": "text-embedding-v4",
  "input": {
    "texts": ["文本1", "文本2"]
  },
  "parameters": {
    "dimension": 1536,
    "text_type": "document",
    "encoding_format": "float",
    "output_type": "dense",
    "instruct": "检索指令"
  }
}
```

- 响应体：

```json
{
  "embeddings": [
    {
      "embedding": [0.01, 0.02],
      "sparse_embedding": [
        { "index": 12, "value": 0.45, "token": "示例" }
      ],
      "text_index": 0
    }
  ],
  "usage": { "total_tokens": 0 }
}
```

## 每个参数的使用方法

- `model_name`：向量模型 ID，例如 `text-embedding-v4`、`bge-m3`。
- `base_url`：服务地址。百炼默认可使用 `https://dashscope.aliyuncs.com/.../text-embedding#`。
- `api_key`：鉴权密钥，使用 `Authorization: Bearer <API_KEY>`。
- `dimension`：支持 `2048(仅v4)`、`1536(仅v4)`、`1024`、`768`、`512`、`256`、`128`、`64`。
- `text_type`：`document` 或 `query`。检索任务建议区分 query/document。
- `output_type`：`dense`、`sparse`、`dense&sparse`（v3/v4 支持）。
- `encoding_format`：`float` 或 `base64`。当前 SDK 解析为浮点数组。
- `instruct`：仅在 `text-embedding-v4` 且 `text_type=query` 时生效，用于检索指令对齐。

示例（设置参数）：

```go
emb, _ := mk.GetEmbedder(ctx, &domain.ModelMetadata{
    Provider:       consts.ModelProviderBaiLian,
    ModelName:      "text-embedding-v4",
    BaseURL:        "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding#",
    APIKey:         os.Getenv("bailianapikey"),
    Dimension:      ptr(1536),
    TextType:       ptrs("query"),
    OutputType:     ptrs("dense&sparse"),
    EncodingFormat: ptrs("float"),
    Instruct:       ptrs("检索：返回和食物相关的内容"),
})

texts := []string{"火鸡面"}
res, _ := mk.UseEmbedder(ctx, emb, texts)
_ = res
```

## 其它的注意事项

- v3/v4 模型单次 `texts` 最多 10 条；其它模型最多 25 条。
- `base_url` 允许以 `#` 结尾，SDK 会自动裁剪；如包含 `compatible-mode` 或为空，将回退至默认 DashScope 路径。
- 国际站点可使用 `dashscope-intl.aliyuncs.com`，会自动选择正确默认路径。
- `output_type=sparse` 的原始响应可能存在不同结构，SDK 已兼容多种格式并统一为条目列表。
- 使用 `AzureOpenAI` 时需指定 `api_version`，未指定将默认 `2024-10-21`。
- `Ollama` 当 `base_url` 以 `/v1` 结尾时走 OpenAI 兼容；否则使用原生 BaseURL（自动移除路径）。

```go
func ptr(v int) *int { return &v }
func ptrs(v string) *string { return &v }
```
