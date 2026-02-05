# ModelKit Chat 介绍
- 支持 `OpenAI API` 
- 支持的供应商：
  `兼容OpenAI API 的所有供应商`、`AzureOpenAI`、`Ollama` 、 `DeepSeek`、`Gemini`、`BaiLian`
# 创建chat

```go
package main

import (
    "context"
    "os"
    "log/slog"
    "github.com/chaitin/ModelKit/v2/consts"
    "github.com/chaitin/ModelKit/v2/domain"
    "github.com/chaitin/ModelKit/v2/usecase"
	"github.com/cloudwego/eino/schema"
)

func main() {
    ctx := context.Background()
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
        Level:     slog.LevelInfo,
    }))
    mk := usecase.NewModelKit(logger)

    chatModel, err := mk.GetChatModel(ctx, &domain.ModelMetadata{
        Provider:  consts.ModelProviderOpenAI,
        ModelName: "gpt-4o-mini",
        BaseURL:   "https://api.openai.com/v1",
        APIKey:    "sk-xxxxxx",
    })

    if err != nil {
        logger.Error("get chat model failed", "err", err)
    }
}
```

字段说明（ModelMetadata）：

- `provider`：模型提供商，取值如 `OpenAI`、`AzureOpenAI`、`Ollama`、`DeepSeek`、`Gemini`、`BaiLian`。
- `model_name`：对话模型 ID，例如 `gpt-4o-mini`、`deepseek-chat`。
- `base_url`：OpenAI 兼容客户端会自动调用 `/chat/completions`；不要在 `base_url` 中包含该路径。`Ollama` 若以 `/v1` 结尾走兼容模式，否则走原生。
- `api_key`：鉴权密钥，作为 `Authorization: Bearer <API_KEY>` 使用。
- `api_version`：仅 `AzureOpenAI` 需要，未设置将默认 `2024-10-21`。
- `api_header`：可选的自定义请求头（`key=value` 按行拼接）。
高级参数: 
- `max_tokens`：最大生成长度。
- `temperature`：采样温度。
- `top_p`：核采样。
- `stop`：停止序列。
- `presence_penalty`：存在惩罚。
- `frequency_penalty`：频率惩罚。
- `response_format`：结构化响应格式（OpenAI 兼容）。
- `seed`：确定性采样。
- `logit_bias`：Logit 偏置。

# 使用chat
## 非流式生成

```go
t := float32(0.2)
chatModel, _ := mk.GetChatModel(ctx, &domain.ModelMetadata{
    Provider:  consts.ModelProviderDeepSeek,
    ModelName: "deepseek-chat",
    BaseURL:   "https://api.deepseek.com",
    APIKey:    "sk-xxxxxx",
    Temperature: &t,
})
msgs := []*schema.Message{ schema.UserMessage("讲个笑话") }
res, _ := chatModel.Generate(ctx, msgs)
```

示例结果

```json
{
  "content": "当然，这里有一个轻松的笑话……"
}
```

## 流式生成

```go
stream, _ := chatModel.Stream(ctx, msgs)
for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        break
    }
    fmt.Print(chunk.Content)
}
```

示例结果

```text
当 然 ， 这 里 有 一 个 轻 松 的 笑 话 …
```