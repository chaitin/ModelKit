package consts

type ModelStatus string

const (
	ModelStatusActive   ModelStatus = "active"
	ModelStatusInactive ModelStatus = "inactive"
)

type ModelType string

const (
	ModelTypeChat      ModelType = "chat"
	ModelTypeCoder     ModelType = "coder"
	ModelTypeEmbedding ModelType = "embedding"
	ModelTypeReranker  ModelType = "reranker"
)

type ModelOwner string

const (
	ModelOwnerSiliconFlow ModelOwner = "SiliconFlow"
	ModelOwnerOpenAI      ModelOwner = "OpenAI"
	ModelOwnerOllama      ModelOwner = "Ollama"
	ModelOwnerDeepSeek    ModelOwner = "DeepSeek"
	ModelOwnerMoonshot    ModelOwner = "Moonshot"
	ModelOwnerAzureOpenAI ModelOwner = "AzureOpenAI"
	ModelOwnerBaiZhiCloud ModelOwner = "BaiZhiCloud"
	ModelOwnerHunyuan     ModelOwner = "Hunyuan"
	ModelOwnerBaiLian     ModelOwner = "BaiLian"
	ModelOwnerVolcengine  ModelOwner = "Volcengine"
)
