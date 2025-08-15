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
)

func ParseModelProvider(s string) (ModelProvider, error) {
	// 转换为小写进行不区分大小写的比较
	switch strings.ToLower(s) {
	case "siliconflow":
		return ModelProviderSiliconFlow, nil
	case "openai":
		return ModelProviderOpenAI, nil
	case "ollama":
		return ModelProviderOllama, nil
	case "deepseek":
		return ModelProviderDeepSeek, nil
	case "moonshot":
		return ModelProviderMoonshot, nil
	case "azureopenai":
		return ModelProviderAzureOpenAI, nil
	case "baizhicloud", "baizhiyun":
		return ModelProviderBaiZhiCloud, nil
	case "hunyuan":
		return ModelProviderHunyuan, nil
	case "bailian":
		return ModelProviderBaiLian, nil
	case "volcengine":
		return ModelProviderVolcengine, nil
	case "gemini":
		return ModelProviderGemini, nil
	case "zhipu":
		return ModelProviderZhiPu, nil
	default:
		return "", errors.New("invalid model provider")
	}
}
