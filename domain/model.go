package domain

import (
	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/cloudwego/eino-ext/libs/acl/openai"
)

type ModelParam struct {
	ContextWindow      int      `json:"context_window"`
	MaxTokens          int      `json:"max_tokens"`
	R1Enabled          bool     `json:"r1_enabled"`
	SupportComputerUse bool     `json:"support_computer_use"`
	SupportImages      bool     `json:"support_images"`
	SupportPromptCache bool     `json:"support_prompt_cache"`
	Temperature        *float32 `json:"temperature"`
}

type ModelMetadata struct {
	// 基础参数
	ModelName string               `json:"id"`         // 模型的名字
	Object    string               `json:"object"`     // 总是model
	Created   int                  `json:"created"`    // 创建时间
	Provider  consts.ModelProvider `json:"provider"`   // 提供商
	ModelType consts.ModelType     `json:"model_type"` // 模型类型
	// api 调用相关参数
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	APIHeader  string `json:"api_header"`
	APIVersion string `json:"api_version"` // for azure openai
	// 高级参数
	// 限制生成的最大token数量,可选,默认为模型最大值, Ollama不支持
	MaxTokens *int `json:"max_tokens"`
	// 采样温度参数,建议与TopP二选一,范围0-2,值越大输出越随机,可选,默认1.0
	Temperature *float32 `json:"temperature"`
	// 控制采样多样性,建议与Temperature二选一,范围0-1,值越小输出越聚焦,可选,默认1.0
	TopP *float32 `json:"top_p"`
	// API停止生成的序列标记,可选,例如:[]string{"\n", "User:"}
	Stop []string `json:"stop"`
	// 基于存在惩罚重复,范围-2到2,正值增加新主题可能性,可选,默认0, Gemini不支持
	PresencePenalty *float32 `json:"presence_penalty"`
	// 指定模型响应的格式,可选,用于结构化输出, DS,Gemini,Ollama不支持
	ResponseFormat *openai.ChatCompletionResponseFormat `json:"response_format"`
	// 启用确定性采样以获得一致输出,可选,用于可重现结果,  DS,Gemini不支持
	Seed *int `json:"seed"`
	// 基于频率惩罚重复,范围-2到2,正值降低重复可能性,可选,默认0, Gemini不支持
	FrequencyPenalty *float32 `json:"frequency_penalty"`
	// 修改特定token在补全中出现的可能性,可选,token ID到偏置值(-100到100)的映射, DS,Gemini,Ollama不支持
	LogitBias map[string]int `json:"logit_bias"`
	// Embeddng高级参数
	EmbedderParam EmbedderParam `json:"embedder_param"`
}

type EmbedderParam struct {
	// 向量维度，可选：2048(仅v4)、1536(仅v4)、1024、768、512、256、128、64
	Dimension *int `json:"dimension"`
	// 文本类型，检索任务建议区分：query 与 document，默认 document
	TextType *string `json:"text_type"`
	// 输出类型，仅 v3/v4 支持：dense、sparse、dense&sparse，默认 dense
	OutputType *string `json:"output_type"`
	// 输出编码格式，默认 float
	EncodingFormat *string `json:"encoding_format"`
	// 检索指令，仅在 text-embedding-v4 且 TextType=query 时生效
	Instruct *string `json:"instruct"`
}

var Models []ModelMetadata

// getBaiZhiCloudModels 返回百智云模型列表
func getBaiZhiCloudModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "qwen-plus", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-72b-instruct", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-14b-instruct", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-coder-plus", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ModelName: "deepseek-r1", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "kimi-k2-0711-preview", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-coder-480b-a35b-instruct", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ModelName: "deepseek-v3", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-turbo", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-32b-instruct", Object: "model", Provider: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
	}
}

// getDeepSeekModels 返回Deepseek模型列表
func getDeepSeekModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "deepseek-chat", Object: "model", Provider: consts.ModelProviderDeepSeek, ModelType: consts.ModelTypeChat},
		{ModelName: "deepseek-reasoner", Object: "model", Provider: consts.ModelProviderDeepSeek, ModelType: consts.ModelTypeChat},
	}
}

// getHunyuanModels 返回腾讯混元模型列表
func getHunyuanModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "hunyuan-pro", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-vision", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-lite", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-standard-32K", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-standard", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-standard-256k", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-functioncall", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "hunyuan-role", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-code", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeCoder},
		{ModelName: "hunyuan-turbo-vision", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-turbo-latest", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbo", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-large", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-large-longcontext", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-latest", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-20250226", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-20250313", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-t1-latest", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-t1-20250321", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-t1-vision", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-turbos-20250515", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-large-vision", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-t1-20250529", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-20250604", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-vision-20250619", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-t1-vision-20250619", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
		{ModelName: "hunyuan-a13b", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-t1-20250711", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-turbos-20250716", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeChat},
		{ModelName: "hunyuan-vision-7b-20250720", Object: "model", Provider: consts.ModelProviderHunyuan, ModelType: consts.ModelTypeVision},
	}
}

// getBaiLianModels 返回阿里百炼模型列表
func getBaiLianModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "text-embedding-v1", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeEmbedding, BaseURL: "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"},
		{ModelName: "text-embedding-v2", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeEmbedding, BaseURL: "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"},
		{ModelName: "text-embedding-v3", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeEmbedding, BaseURL: "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"},
		{ModelName: "gte-rerank", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeRerank, BaseURL: "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#"},
		{ModelName: "qwen3-rerank", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeRerank, BaseURL: "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank#"},
		{ModelName: "qwen3-coder-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen3-coder-plus-2025-07-22", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-plus-2025-07-14", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-coder-480b-a35b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-mt-turbo", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-mt-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-tts-2025-05-22", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qvq-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qvq-plus-2025-05-15", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qvq-max-2025-05-15", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-4b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-32b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-30b-a3b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-235b-a22b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-14b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-1.7b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-0.6b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen3-8b", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-vl-max-2025-04-02", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeVision},
		{ModelName: "qwen-vl-ocr-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeVision},
		{ModelName: "qwen-vl-ocr", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeVision},
		{ModelName: "qwen-coder-plus-1106", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-coder-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-coder-plus-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen2.5-coder-3b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen2.5-coder-0.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen2.5-coder-14b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen2.5-coder-32b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-math-turbo", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-3b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-math-1.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-1.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-0.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-32b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-72b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-coder-7b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen2.5-math-7b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-coder-turbo", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-coder-turbo-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-math-plus-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-math-plus-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-math-turbo-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-plus-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-plus-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-turbo-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-math-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-14b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-7b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2.5-math-72b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-coder-turbo-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen-math-turbo-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-latest", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-turbo-0919", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2-1.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2-72b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2-7b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2-0.5b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen2-57b-a14b-instruct", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-long", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-vl-max", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeVision},
		{ModelName: "qwen-vl-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeVision},
		{ModelName: "qwen-max-0428", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen1.5-110b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen1.5-0.5b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-1.8b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-1.8b-longcontext-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-7b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-14b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-72b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "codeqwen1.5-7b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeCoder},
		{ModelName: "qwen1.5-32b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen1.5-72b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-longcontext", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-1201", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen1.5-1.8b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen1.5-14b-chat", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-turbo", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-plus", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-0403", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
		{ModelName: "qwen-max-0107", Object: "model", Provider: consts.ModelProviderBaiLian, ModelType: consts.ModelTypeChat},
	}
}

// getVolcengineModels 返回火山引擎模型列表
func getVolcengineModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "doubao-seed-1.6-250615", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		{ModelName: "doubao-seed-1.6-flash-250615", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		// {ModelName: "doubao-seed-1.6-flash-250715", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		{ModelName: "doubao-seed-1.6-thinking-250615", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		// {ModelName: "doubao-seed-1.6-thinking-250715", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		{ModelName: "doubao-1.5-thinking-vision-pro-250428", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeVision},
		// {ModelName: "Doubao-1.5-thinking-pro-250415", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
		{ModelName: "deepseek-r1-250528", Object: "model", Provider: consts.ModelProviderVolcengine, ModelType: consts.ModelTypeChat},
	}
}

// getOpenAIModels 返回OpenAI模型列表
func getOpenAIModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "text-embedding-ada-002", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "whisper-1", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-3.5-turbo", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "tts-1", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-3.5-turbo-16k", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "davinci-002", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "babbage-002", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-3.5-turbo-instruct", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-3.5-turbo-instruct-0914", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "dall-e-3", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeVision},
		{ModelName: "dall-e-2", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeVision},
		{ModelName: "gpt-3.5-turbo-1106", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "tts-1-hd", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "tts-1-1106", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "tts-1-hd-1106", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "text-embedding-3-small", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "text-embedding-3-large", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "gpt-3.5-turbo-0125", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-2024-05-13", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-mini-2024-07-18", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-mini", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-2024-08-06", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o1-mini-2024-09-12", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o1-mini", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-audio-preview-2024-10-01", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-audio-preview", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "omni-moderation-latest", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "omni-moderation-2024-09-26", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-audio-preview-2024-12-17", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-audio-preview-2024-12-17", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-audio-preview", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-2024-11-20", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-search-preview-2025-03-11", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-search-preview", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-search-preview-2025-03-11", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-search-preview", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-transcribe", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-transcribe", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4o-mini-tts", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "gpt-4.1-2025-04-14", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-mini-2025-04-14", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-mini", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-nano-2025-04-14", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-nano", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-image-1", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeVision},
		{ModelName: "gpt-4o-audio-preview-2025-06-03", Object: "model", Provider: consts.ModelProviderOpenAI, ModelType: consts.ModelTypeFunctionCall},
	}
}

// getSiliconFlowModels 返回硅基流动模型列表
func getSiliconFlowModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "stabilityai/stable-diffusion-xl-base-1.0", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "THUDM/glm-4-9b-chat", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "internlm/internlm2_5-7b-chat", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "BAAI/bge-large-en-v1.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "BAAI/bge-large-zh-v1.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "Pro/Qwen/Qwen2-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/THUDM/glm-4-9b-chat", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "black-forest-labs/FLUX.1-schnell", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "black-forest-labs/FLUX.1-dev", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "FunAudioLLM/SenseVoiceSmall", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "netease-youdao/bce-embedding-base_v1", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "BAAI/bge-m3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "netease-youdao/bce-reranker-base_v1", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "BAAI/bge-reranker-v2-m3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "deepseek-ai/DeepSeek-V2.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-72B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-14B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-32B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/black-forest-labs/FLUX.1-schnell", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-Coder-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "Pro/Qwen/Qwen2.5-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-72B-Instruct-128K", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2-VL-72B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Pro/BAAI/bge-m3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "stabilityai/stable-diffusion-3-5-large", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "LoRA/Qwen/Qwen2.5-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "fishaudio/fish-speech-1.4", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "Pro/Qwen/Qwen2.5-Coder-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "LoRA/Qwen/Qwen2.5-72B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-Coder-32B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "Pro/BAAI/bge-reranker-v2-m3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "RVC-Boss/GPT-SoVITS", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "fishaudio/fish-speech-1.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "black-forest-labs/FLUX.1-pro", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "LoRA/Qwen/Qwen2.5-14B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "LoRA/Qwen/Qwen2.5-32B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "FunAudioLLM/CosyVoice2-0.5B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "deepseek-ai/deepseek-vl2", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Qwen/QVQ-72B-Preview", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "LoRA/black-forest-labs/FLUX.1-dev", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen2.5-VL-72B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Pro/Qwen/Qwen2.5-VL-7B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "deepseek-ai/DeepSeek-V3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "deepseek-ai/DeepSeek-R1", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/deepseek-ai/DeepSeek-R1-Distill-Qwen-7B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "deepseek-ai/DeepSeek-R1-Distill-Qwen-14B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "deepseek-ai/DeepSeek-R1-Distill-Qwen-32B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "deepseek-ai/DeepSeek-R1-Distill-Qwen-7B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "Pro/deepseek-ai/DeepSeek-R1", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/deepseek-ai/DeepSeek-V3", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Kwai-Kolors/Kolors", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Qwen/QwQ-32B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Wan-AI/Wan2.1-T2V-14B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Wan-AI/Wan2.1-T2V-14B-Turbo", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Wan-AI/Wan2.1-I2V-14B-720P", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Wan-AI/Wan2.1-I2V-14B-720P-Turbo", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "Qwen/Qwen2.5-VL-32B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeVision},
		{ModelName: "THUDM/GLM-Z1-32B-0414", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "THUDM/GLM-4-32B-0414", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "THUDM/GLM-Z1-9B-0414", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "THUDM/GLM-4-9B-0414", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "THUDM/GLM-Z1-Rumination-32B-0414", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-8B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-14B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-32B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-30B-A3B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-235B-A22B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Tongyi-Zhiwen/QwenLong-L1-32B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-Embedding-8B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "Qwen/Qwen3-Embedding-4B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "Qwen/Qwen3-Embedding-0.6B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ModelName: "MiniMaxAI/MiniMax-M1-80k", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-Reranker-0.6B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "Qwen/Qwen3-Reranker-4B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "Qwen/Qwen3-Reranker-8B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeRerank},
		{ModelName: "fnlp/MOSS-TTSD-v0.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ModelName: "moonshotai/Kimi-Dev-72B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "tencent/Hunyuan-A13B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "baidu/ERNIE-4.5-300B-A47B", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "THUDM/GLM-4.1V-9B-Thinking", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/THUDM/GLM-4.1V-9B-Thinking", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "ascend-tribe/pangu-pro-moe", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "moonshotai/Kimi-K2-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Pro/moonshotai/Kimi-K2-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-235B-A22B-Instruct-2507", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-Coder-480B-A35B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ModelName: "Qwen/Qwen3-235B-A22B-Thinking-2507", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "zai-org/GLM-4.5", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "zai-org/GLM-4.5-Air", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-30B-A3B-Instruct-2507", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-30B-A3B-Thinking-2507", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeChat},
		{ModelName: "Qwen/Qwen3-Coder-30B-A3B-Instruct", Object: "model", Provider: consts.ModelProviderSiliconFlow, ModelType: consts.ModelTypeCoder},
	}
}

// 月之暗面模型
func getMoonshotModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "moonshot-v1-auto", Object: "model", Provider: consts.ModelProviderMoonshot, ModelType: consts.ModelTypeChat},
		{ModelName: "moonshot-v1-8k", Object: "model", Provider: consts.ModelProviderMoonshot, ModelType: consts.ModelTypeChat},
		{ModelName: "moonshot-v1-32k", Object: "model", Provider: consts.ModelProviderMoonshot, ModelType: consts.ModelTypeChat},
		{ModelName: "moonshot-v1-128k", Object: "model", Provider: consts.ModelProviderMoonshot, ModelType: consts.ModelTypeChat},
	}

}

// getAzureOpenAIModels 返回Azure OpenAI模型列表
func getAzureOpenAIModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "gpt-4", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-mini", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4o-nano", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-mini", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "gpt-4.1-nano", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o1", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o1-mini", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o3", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o3-mini", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ModelName: "o4-mini", Object: "model", Provider: consts.ModelProviderAzureOpenAI, ModelType: consts.ModelTypeChat},
	}
}

// getGeminiModels 返回Google Gemini模型列表
func getGeminiModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "gemini-2.5-pro", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.5-flash", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.5-flash-lite-preview-06-17", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.5-flash-preview-tts", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.5-pro-preview-tts", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.0-flash", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-2.0-flash-lite", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-1.5-flash", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-1.5-flash-8b", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-1.5-pro", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeChat},
		{ModelName: "gemini-embedding-001", Object: "model", Provider: consts.ModelProviderGemini, ModelType: consts.ModelTypeEmbedding},
	}
}

// getZhiPuModels 返回智谱模型列表
func getZhiPuModels() []ModelMetadata {
	return []ModelMetadata{
		{ModelName: "glm-4.5", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4.5-x", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		{ModelName: "glm-4.5-air", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4.5-airx", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		{ModelName: "glm-4.5-flash", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-plus", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-air-250414", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-airx", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-long", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-flashx-250414", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4-flash-250414", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-z1-air", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-z1-airx", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-z1-flashx", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-z1-flash", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeChat},
		// {ModelName: "glm-4v-plus-0111", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeVision},
		// {ModelName: "glm-4v-flash", Object: "model", Provider: consts.ModelProviderZhiPu, ModelType: consts.ModelTypeVision},
	}
}

func initModels() {
	Models = make([]ModelMetadata, 0, 200)

	Models = append(Models, getBaiZhiCloudModels()...)
	Models = append(Models, getDeepSeekModels()...)
	Models = append(Models, getHunyuanModels()...)
	Models = append(Models, getBaiLianModels()...)
	Models = append(Models, getVolcengineModels()...)
	Models = append(Models, getOpenAIModels()...)
	Models = append(Models, getSiliconFlowModels()...)
	Models = append(Models, getMoonshotModels()...)
	Models = append(Models, getAzureOpenAIModels()...)
	Models = append(Models, getZhiPuModels()...)
	Models = append(Models, getGeminiModels()...)
}
