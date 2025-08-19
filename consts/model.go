package consts

import (
	"errors"
	"strings"
)

type ModelStatus string

const (
	ModelStatusActive   ModelStatus = "active"
	ModelStatusInactive ModelStatus = "inactive"
)

type ModelType string

const (
	ModelTypeChat         ModelType = "llm"
	ModelTypeCoder        ModelType = "coder"
	ModelTypeEmbedding    ModelType = "embedding"
	ModelTypeRerank       ModelType = "reranker"
	ModelTypeVision       ModelType = "vision"
	ModelTypeFunctionCall ModelType = "function_call"
)

func ParseModelType(s string) (ModelType, error) {
	switch s {
	case "llm", "chat":
		return ModelTypeChat, nil
	case "coder", "code":
		return ModelTypeCoder, nil
	case "embedding":
		return ModelTypeEmbedding, nil
	case "reranker", "rerank":
		return ModelTypeRerank, nil
	case "vision":
		return ModelTypeVision, nil
	case "function_call":
		return ModelTypeFunctionCall, nil
	default:
		return "", errors.New("invalid model type")
	}
}

type ModelProvider string

const (
	ModelProviderSiliconFlow ModelProvider = "SiliconFlow"
	ModelProviderOpenAI      ModelProvider = "OpenAI"
	ModelProviderOllama      ModelProvider = "Ollama"
	ModelProviderDeepSeek    ModelProvider = "DeepSeek"
	ModelProviderMoonshot    ModelProvider = "Moonshot"
	ModelProviderAzureOpenAI ModelProvider = "AzureOpenAI"
	ModelProviderBaiZhiCloud ModelProvider = "BaiZhiCloud"
	ModelProviderHunyuan     ModelProvider = "Hunyuan"
	ModelProviderBaiLian     ModelProvider = "BaiLian"
	ModelProviderVolcengine  ModelProvider = "Volcengine"
	ModelProviderGemini      ModelProvider = "Gemini"
	ModelProviderZhiPu       ModelProvider = "ZhiPu"
	ModelProviderOther       ModelProvider = "Other"
)

func ParseModelProvider(s string) ModelProvider {
	// 转换为小写进行不区分大小写的比较
	switch strings.ToLower(s) {
	case "siliconflow":
		return ModelProviderSiliconFlow
	case "openai":
		return ModelProviderOpenAI
	case "ollama":
		return ModelProviderOllama
	case "deepseek":
		return ModelProviderDeepSeek
	case "moonshot":
		return ModelProviderMoonshot
	case "azureopenai":
		return ModelProviderAzureOpenAI
	case "baizhicloud", "baizhiyun":
		return ModelProviderBaiZhiCloud
	case "hunyuan":
		return ModelProviderHunyuan
	case "bailian":
		return ModelProviderBaiLian
	case "volcengine":
		return ModelProviderVolcengine
	case "gemini":
		return ModelProviderGemini
	case "zhipu":
		return ModelProviderZhiPu
	default:
		return ModelProviderOther
	}
}
