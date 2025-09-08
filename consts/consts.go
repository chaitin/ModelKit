package consts

import (
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

func ParseModelType(s string) ModelType {
	switch s {
	case "llm", "chat", "analysis":
		return ModelTypeChat
	case "coder", "code":
		return ModelTypeCoder
	case "embedding":
		return ModelTypeEmbedding
	case "reranker", "rerank":
		return ModelTypeRerank
	case "vision":
		return ModelTypeVision
	case "function_call":
		return ModelTypeFunctionCall
	default:
		return ModelTypeChat
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

var ApiKeyBalanceKeyWords = []string{"quota", "billing", "balance", "payment required"}

type AddModelBaseURLErrType string

const (
	AddModelBaseURLErrTypeHost   AddModelBaseURLErrType = "host"    // 建议使用宿主机主机名
	AddModelBaseURLErrTypeV1Path AddModelBaseURLErrType = "v1_path" // 建议在API地址末尾添加/v1
	AddModelBaseURLErrTypeSlash  AddModelBaseURLErrType = "slash"   // 建议去掉末尾的/
	AddModelBaseURLErrTypeProtocol AddModelBaseURLErrType = "protocol" // 建议在url开头添加协议 http:// 或者 https://
	AddModelBaseURLErrTypeChatCompletions AddModelBaseURLErrType = "chat_completions" // 建议去掉/chat/completions路径
)

const LinuxHost = "172.17.0.1"
const MacWinHost = "host.docker.internal"
const LocalHost = "localhost"
const LocalIP = "127.0.0.1"
