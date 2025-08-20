package domain

import "github.com/chaitin/ModelKit/consts"

type IModelProvider[T any] interface {
	ListModel(subType string, provider string) ([]T, error)
}

var ModelProviders map[consts.ModelProvider]ModelProvider
var TypeModelMap map[consts.ModelType][]ModelMetadata

type ModelProvider struct {
	OwnerName  consts.ModelProvider `json:"owner_name"`  // 提供商
	APIBase    string               `json:"api_base"`    // 接口地址 如：https://api.qwen.com
	APIKey     string               `json:"api_key"`     // 接口密钥 如：sk-xxxx
	APIVersion string               `json:"api_version"` // 接口版本 如：2023-05-15
	APIHeader  string               `json:"api_header"`  // 接口头 如：Authorization: Bearer sk-xxxx
	Models     []ModelMetadata      `json:"models"`      // 模型列表
}

func initModelProviders() {
	// 初始化模型提供商及其模型
	ModelProviders = map[consts.ModelProvider]ModelProvider{
		consts.ModelProviderBaiZhiCloud: {
			OwnerName: consts.ModelProviderBaiZhiCloud,
			Models:    getModelsByOwner(consts.ModelProviderBaiZhiCloud),
			APIBase:   "https://model-square.app.baizhi.cloud/v1",
		},
		consts.ModelProviderDeepSeek: {
			OwnerName: consts.ModelProviderDeepSeek,
			Models:    getModelsByOwner(consts.ModelProviderDeepSeek),
			APIBase:   "https://api.deepseek.com/v1",
		},
		consts.ModelProviderSiliconFlow: {
			OwnerName: consts.ModelProviderSiliconFlow,
			Models:    getModelsByOwner(consts.ModelProviderSiliconFlow),
			APIBase:   "https://api.siliconflow.cn/v1",
		},
		consts.ModelProviderOpenAI: {
			OwnerName: consts.ModelProviderOpenAI,
			Models:    getModelsByOwner(consts.ModelProviderOpenAI),
			APIBase:   "https://api.openai.com/v1",
		},
		consts.ModelProviderOllama: {
			OwnerName: consts.ModelProviderOllama,
			Models:    getModelsByOwner(consts.ModelProviderOllama),
			APIBase:   "http://localhost:11434",
		},
		consts.ModelProviderMoonshot: {
			OwnerName: consts.ModelProviderMoonshot,
			Models:    getModelsByOwner(consts.ModelProviderMoonshot),
			APIBase:   "https://api.moonshot.cn/v1",
		},
		consts.ModelProviderAzureOpenAI: {
			OwnerName: consts.ModelProviderAzureOpenAI,
			Models:    getModelsByOwner(consts.ModelProviderAzureOpenAI),
		},
		consts.ModelProviderHunyuan: {
			OwnerName: consts.ModelProviderHunyuan,
			Models:    getModelsByOwner(consts.ModelProviderHunyuan),
			APIBase:   "https://api.hunyuan.cloud.tencent.com/v1",
		},
		consts.ModelProviderBaiLian: {
			OwnerName: consts.ModelProviderBaiLian,
			Models:    getModelsByOwner(consts.ModelProviderBaiLian),
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
		},
		consts.ModelProviderVolcengine: {
			OwnerName: consts.ModelProviderVolcengine,
			Models:    getModelsByOwner(consts.ModelProviderVolcengine),
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
		},
		consts.ModelProviderGemini: {
			OwnerName: consts.ModelProviderGemini,
			Models:    getModelsByOwner(consts.ModelProviderGemini),
			APIBase:   "https://generativelanguage.googleapis.com/v1beta",
		},
		consts.ModelProviderZhiPu: {
			OwnerName: consts.ModelProviderZhiPu,
			Models:    getModelsByOwner(consts.ModelProviderZhiPu),
			APIBase:   "https://open.bigmodel.cn/api/paas/v4",
		},
		consts.ModelProviderAiHubMix: {
			OwnerName: consts.ModelProviderAiHubMix,
			Models:    getModelsByOwner(consts.ModelProviderAiHubMix),
			APIBase:   "https://aihubmix.com/v1",
		},
		consts.ModelProviderOcoolAI: {
			OwnerName: consts.ModelProviderOcoolAI,
			Models:    getModelsByOwner(consts.ModelProviderOcoolAI),
			APIBase:   "https://api.ocoolai.com/v1",
		},
		consts.ModelProviderPPIO: {
			OwnerName: consts.ModelProviderPPIO,
			Models:    getModelsByOwner(consts.ModelProviderPPIO),
			APIBase:   "https://api.ppinfra.com/v3/openai",
		},
		consts.ModelProviderAlayaNew: {
			OwnerName: consts.ModelProviderAlayaNew,
			Models:    getModelsByOwner(consts.ModelProviderAlayaNew),
			APIBase:   "https://deepseek.alayanew.com/v1",
		},
		consts.ModelProviderQiniu: {
			OwnerName: consts.ModelProviderQiniu,
			Models:    getModelsByOwner(consts.ModelProviderQiniu),
			APIBase:   "https://api.qnaigc.com/v1",
		},
		consts.ModelProviderDMXAPI: {
			OwnerName: consts.ModelProviderDMXAPI,
			Models:    getModelsByOwner(consts.ModelProviderDMXAPI),
			APIBase:   "https://www.dmxapi.cn/v1",
		},
		consts.ModelProviderBurnCloud: {
			OwnerName: consts.ModelProviderBurnCloud,
			Models:    getModelsByOwner(consts.ModelProviderBurnCloud),
			APIBase:   "https://ai.burncloud.com/v1",
		},
		consts.ModelProviderTokenFlux: {
			OwnerName: consts.ModelProviderTokenFlux,
			Models:    getModelsByOwner(consts.ModelProviderTokenFlux),
			APIBase:   "https://tokenflux.ai/v1",
		},
		consts.ModelProvider302AI: {
			OwnerName: consts.ModelProvider302AI,
			Models:    getModelsByOwner(consts.ModelProvider302AI),
			APIBase:   "https://api.302.ai/v1",
		},
		consts.ModelProviderCephalon: {
			OwnerName: consts.ModelProviderCephalon,
			Models:    getModelsByOwner(consts.ModelProviderCephalon),
			APIBase:   "https://cephalon.cloud/user-center/v1/model",
		},
		consts.ModelProviderLanyun: {
			OwnerName: consts.ModelProviderLanyun,
			Models:    getModelsByOwner(consts.ModelProviderLanyun),
			APIBase:   "https://maas-api.lanyun.net/v1",
		},
		consts.ModelProviderPH8: {
			OwnerName: consts.ModelProviderPH8,
			Models:    getModelsByOwner(consts.ModelProviderPH8),
			APIBase:   "https://ph8.co/v1",
		},
		consts.ModelProviderOpenRouter: {
			OwnerName: consts.ModelProviderOpenRouter,
			Models:    getModelsByOwner(consts.ModelProviderOpenRouter),
			APIBase:   "https://openrouter.ai/api/v1",
		},
		consts.ModelProviderNewAPI: {
			OwnerName: consts.ModelProviderNewAPI,
			Models:    getModelsByOwner(consts.ModelProviderNewAPI),
			APIBase:   "http://localhost:3000/v1",
		},
		consts.ModelProviderLMStudio: {
			OwnerName: consts.ModelProviderLMStudio,
			Models:    getModelsByOwner(consts.ModelProviderLMStudio),
			APIBase:   "http://localhost:1234/v1",
		},
		consts.ModelProviderAnthropic: {
			OwnerName: consts.ModelProviderAnthropic,
			Models:    getModelsByOwner(consts.ModelProviderAnthropic),
			APIBase:   "https://api.anthropic.com",
		},
		consts.ModelProviderVertexAI: {
			OwnerName: consts.ModelProviderVertexAI,
			Models:    getModelsByOwner(consts.ModelProviderVertexAI),
			APIBase:   "https://aiplatform.googleapis.com",
		},
		consts.ModelProviderGithub: {
			OwnerName: consts.ModelProviderGithub,
			Models:    getModelsByOwner(consts.ModelProviderGithub),
			APIBase:   "https://models.github.ai/inference/v1",
		},
		consts.ModelProviderCopilot: {
			OwnerName: consts.ModelProviderCopilot,
			Models:    getModelsByOwner(consts.ModelProviderCopilot),
			APIBase:   "https://api.githubcopilot.com/v1",
		},
		consts.ModelProviderYi: {
			OwnerName: consts.ModelProviderYi,
			Models:    getModelsByOwner(consts.ModelProviderYi),
			APIBase:   "https://api.lingyiwanwu.com/v1",
		},
		consts.ModelProviderBaichuan: {
			OwnerName: consts.ModelProviderBaichuan,
			Models:    getModelsByOwner(consts.ModelProviderBaichuan),
			APIBase:   "https://api.baichuan-ai.com/v1",
		},
		consts.ModelProviderStepFun: {
			OwnerName: consts.ModelProviderStepFun,
			Models:    getModelsByOwner(consts.ModelProviderStepFun),
			APIBase:   "https://api.stepfun.com/v1",
		},
		consts.ModelProviderInfini: {
			OwnerName: consts.ModelProviderInfini,
			Models:    getModelsByOwner(consts.ModelProviderInfini),
			APIBase:   "https://cloud.infini-ai.com/maas/v1",
		},
		consts.ModelProviderMiniMax: {
			OwnerName: consts.ModelProviderMiniMax,
			Models:    getModelsByOwner(consts.ModelProviderMiniMax),
			APIBase:   "https://api.minimax.chat/v1",
		},
		consts.ModelProviderGroq: {
			OwnerName: consts.ModelProviderGroq,
			Models:    getModelsByOwner(consts.ModelProviderGroq),
			APIBase:   "https://api.groq.com/openai/v1",
		},
		consts.ModelProviderTogether: {
			OwnerName: consts.ModelProviderTogether,
			Models:    getModelsByOwner(consts.ModelProviderTogether),
			APIBase:   "https://api.together.xyz/v1",
		},
		consts.ModelProviderFireworks: {
			OwnerName: consts.ModelProviderFireworks,
			Models:    getModelsByOwner(consts.ModelProviderFireworks),
			APIBase:   "https://api.fireworks.ai/inference/v1",
		},
		consts.ModelProviderNvidia: {
			OwnerName: consts.ModelProviderNvidia,
			Models:    getModelsByOwner(consts.ModelProviderNvidia),
			APIBase:   "https://integrate.api.nvidia.com/v1",
		},
		consts.ModelProviderGrok: {
			OwnerName: consts.ModelProviderGrok,
			Models:    getModelsByOwner(consts.ModelProviderGrok),
			APIBase:   "https://api.x.ai/v1",
		},
		consts.ModelProviderHyperbolic: {
			OwnerName: consts.ModelProviderHyperbolic,
			Models:    getModelsByOwner(consts.ModelProviderHyperbolic),
			APIBase:   "https://api.hyperbolic.xyz/v1",
		},
		consts.ModelProviderMistral: {
			OwnerName: consts.ModelProviderMistral,
			Models:    getModelsByOwner(consts.ModelProviderMistral),
			APIBase:   "https://api.mistral.ai/v1",
		},
		consts.ModelProviderJina: {
			OwnerName: consts.ModelProviderJina,
			Models:    getModelsByOwner(consts.ModelProviderJina),
			APIBase:   "https://api.jina.ai/v1",
		},
		consts.ModelProviderPerplexity: {
			OwnerName: consts.ModelProviderPerplexity,
			Models:    getModelsByOwner(consts.ModelProviderPerplexity),
			APIBase:   "https://api.perplexity.ai/v1",
		},
		consts.ModelProviderModelScope: {
			OwnerName: consts.ModelProviderModelScope,
			Models:    getModelsByOwner(consts.ModelProviderModelScope),
			APIBase:   "https://api-inference.modelscope.cn/v1",
		},
		consts.ModelProviderXirang: {
			OwnerName: consts.ModelProviderXirang,
			Models:    getModelsByOwner(consts.ModelProviderXirang),
			APIBase:   "https://wishub-x1.ctyun.cn/v1",
		},
		consts.ModelProviderTencentCloudTI: {
			OwnerName: consts.ModelProviderTencentCloudTI,
			Models:    getModelsByOwner(consts.ModelProviderTencentCloudTI),
			APIBase:   "https://api.lkeap.cloud.tencent.com/v1",
		},
		consts.ModelProviderBaiduCloud: {
			OwnerName: consts.ModelProviderBaiduCloud,
			Models:    getModelsByOwner(consts.ModelProviderBaiduCloud),
			APIBase:   "https://qianfan.baidubce.com/v2",
		},
		consts.ModelProviderGPUStack: {
			OwnerName: consts.ModelProviderGPUStack,
			Models:    getModelsByOwner(consts.ModelProviderGPUStack),
			APIBase:   "",
		},
		consts.ModelProviderVoyageAI: {
			OwnerName: consts.ModelProviderVoyageAI,
			Models:    getModelsByOwner(consts.ModelProviderVoyageAI),
			APIBase:   "https://api.voyageai.com/v1",
		},
		consts.ModelProviderAWSBedrock: {
			OwnerName: consts.ModelProviderAWSBedrock,
			Models:    getModelsByOwner(consts.ModelProviderAWSBedrock),
			APIBase:   "",
		},
		consts.ModelProviderPoe: {
			OwnerName: consts.ModelProviderPoe,
			Models:    getModelsByOwner(consts.ModelProviderPoe),
			APIBase:   "https://api.poe.com/v1",
		},
	}

	// 初始化按类型分组的模型映射
	TypeModelMap = make(map[consts.ModelType][]ModelMetadata)
	for i := range Models {
		model := Models[i]
		TypeModelMap[model.ModelType] = append(TypeModelMap[model.ModelType], model)
	}
}