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
	ModelProviderSiliconFlow    ModelProvider = "SiliconFlow"
	ModelProviderOpenAI         ModelProvider = "OpenAI"
	ModelProviderOllama         ModelProvider = "Ollama"
	ModelProviderDeepSeek       ModelProvider = "DeepSeek"
	ModelProviderMoonshot       ModelProvider = "Moonshot"
	ModelProviderAzureOpenAI    ModelProvider = "AzureOpenAI"
	ModelProviderBaiZhiCloud    ModelProvider = "BaiZhiCloud"
	ModelProviderHunyuan        ModelProvider = "Hunyuan"
	ModelProviderBaiLian        ModelProvider = "BaiLian"
	ModelProviderVolcengine     ModelProvider = "Volcengine"
	ModelProviderGemini         ModelProvider = "Gemini"
	ModelProviderZhiPu          ModelProvider = "ZhiPu"
	ModelProviderAiHubMix       ModelProvider = "AiHubMix"
	ModelProviderOcoolAI        ModelProvider = "OcoolAI"
	ModelProviderPPIO           ModelProvider = "PPIO"
	ModelProviderAlayaNew       ModelProvider = "AlayaNew"
	ModelProviderQiniu          ModelProvider = "Qiniu"
	ModelProviderDMXAPI         ModelProvider = "DMXAPI"
	ModelProviderBurnCloud      ModelProvider = "BurnCloud"
	ModelProviderTokenFlux      ModelProvider = "TokenFlux"
	ModelProvider302AI          ModelProvider = "302AI"
	ModelProviderCephalon       ModelProvider = "Cephalon"
	ModelProviderLanyun         ModelProvider = "Lanyun"
	ModelProviderPH8            ModelProvider = "PH8"
	ModelProviderOpenRouter     ModelProvider = "OpenRouter"
	ModelProviderNewAPI         ModelProvider = "NewAPI"
	ModelProviderLMStudio       ModelProvider = "LMStudio"
	ModelProviderAnthropic      ModelProvider = "Anthropic"
	ModelProviderVertexAI       ModelProvider = "VertexAI"
	ModelProviderGithub         ModelProvider = "Github"
	ModelProviderCopilot        ModelProvider = "Copilot"
	ModelProviderYi             ModelProvider = "Yi"
	ModelProviderBaichuan       ModelProvider = "Baichuan"
	ModelProviderStepFun        ModelProvider = "StepFun"
	ModelProviderInfini         ModelProvider = "Infini"
	ModelProviderMiniMax        ModelProvider = "MiniMax"
	ModelProviderGroq           ModelProvider = "Groq"
	ModelProviderTogether       ModelProvider = "Together"
	ModelProviderFireworks      ModelProvider = "Fireworks"
	ModelProviderNvidia         ModelProvider = "Nvidia"
	ModelProviderGrok           ModelProvider = "Grok"
	ModelProviderHyperbolic     ModelProvider = "Hyperbolic"
	ModelProviderMistral        ModelProvider = "Mistral"
	ModelProviderJina           ModelProvider = "Jina"
	ModelProviderPerplexity     ModelProvider = "Perplexity"
	ModelProviderModelScope     ModelProvider = "ModelScope"
	ModelProviderXirang         ModelProvider = "Xirang"
	ModelProviderTencentCloudTI ModelProvider = "TencentCloudTI"
	ModelProviderBaiduCloud     ModelProvider = "BaiduCloud"
	ModelProviderGPUStack       ModelProvider = "GPUStack"
	ModelProviderVoyageAI       ModelProvider = "VoyageAI"
	ModelProviderAWSBedrock     ModelProvider = "AWSBedrock"
	ModelProviderPoe            ModelProvider = "Poe"
	ModelProviderOther          ModelProvider = "Other"
)

func ParseModelProvider(s string) ModelProvider {
	// 转换为小写进行不区分大小写的比较
	switch strings.ToLower(s) {
	case "siliconflow", "silicon":
		return ModelProviderSiliconFlow
	case "openai":
		return ModelProviderOpenAI
	case "ollama":
		return ModelProviderOllama
	case "deepseek":
		return ModelProviderDeepSeek
	case "moonshot":
		return ModelProviderMoonshot
	case "azureopenai", "azure-openai":
		return ModelProviderAzureOpenAI
	case "baizhicloud", "baizhiyun":
		return ModelProviderBaiZhiCloud
	case "hunyuan":
		return ModelProviderHunyuan
	case "bailian":
		return ModelProviderBaiLian
	case "volcengine", "doubao":
		return ModelProviderVolcengine
	case "gemini":
		return ModelProviderGemini
	case "zhipu":
		return ModelProviderZhiPu
	case "aihubmix":
		return ModelProviderAiHubMix
	case "ocoolai":
		return ModelProviderOcoolAI
	case "ppio":
		return ModelProviderPPIO
	case "alayanew":
		return ModelProviderAlayaNew
	case "qiniu":
		return ModelProviderQiniu
	case "dmxapi":
		return ModelProviderDMXAPI
	case "burncloud":
		return ModelProviderBurnCloud
	case "tokenflux":
		return ModelProviderTokenFlux
	case "302ai":
		return ModelProvider302AI
	case "cephalon":
		return ModelProviderCephalon
	case "lanyun":
		return ModelProviderLanyun
	case "ph8":
		return ModelProviderPH8
	case "openrouter":
		return ModelProviderOpenRouter
	case "new-api":
		return ModelProviderNewAPI
	case "lmstudio":
		return ModelProviderLMStudio
	case "anthropic":
		return ModelProviderAnthropic
	case "vertexai":
		return ModelProviderVertexAI
	case "github":
		return ModelProviderGithub
	case "copilot":
		return ModelProviderCopilot
	case "yi":
		return ModelProviderYi
	case "baichuan":
		return ModelProviderBaichuan
	case "stepfun":
		return ModelProviderStepFun
	case "infini":
		return ModelProviderInfini
	case "minimax":
		return ModelProviderMiniMax
	case "groq":
		return ModelProviderGroq
	case "together":
		return ModelProviderTogether
	case "fireworks":
		return ModelProviderFireworks
	case "nvidia":
		return ModelProviderNvidia
	case "grok":
		return ModelProviderGrok
	case "hyperbolic":
		return ModelProviderHyperbolic
	case "mistral":
		return ModelProviderMistral
	case "jina":
		return ModelProviderJina
	case "perplexity":
		return ModelProviderPerplexity
	case "modelscope":
		return ModelProviderModelScope
	case "xirang":
		return ModelProviderXirang
	case "tencent-cloud-ti":
		return ModelProviderTencentCloudTI
	case "baidu-cloud":
		return ModelProviderBaiduCloud
	case "gpustack":
		return ModelProviderGPUStack
	case "voyageai":
		return ModelProviderVoyageAI
	case "aws-bedrock":
		return ModelProviderAWSBedrock
	case "poe":
		return ModelProviderPoe
	default:
		return ModelProviderOther
	}
}
