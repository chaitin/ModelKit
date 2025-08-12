package consts

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
