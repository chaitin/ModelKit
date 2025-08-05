package domain

import (
	"context"

	"github.com/chaitin/ModelKit/backend/consts"
)

type ModelUsecase interface {
	CheckModel(ctx context.Context, req *CheckModelReq) (*Model, error)
	ListModel(ctx context.Context, req *ListModelReq) ([]*Model, error)
	PandaModelList(ctx context.Context, req *GetProviderModelListReq) (*GetProviderModelListResp, error)
}

type CheckModelReq struct {
	Owner   consts.ModelOwner `json:"owner" validate:"required"`      // 提供商
	ModelID string            `json:"model_name" validate:"required"` // 模型名称
	APIKey  string            `json:"api_key" validate:"required"`    // 接口密钥
	SubType consts.ModelType  `json:"sub_type" validate:"required"`   // 模型类型
}

type ListModelReq struct {
	OwnedBy consts.ModelOwner `json:"owned_by" query:"owned_by"` // 提供商
	SubType consts.ModelType  `json:"sub_type" query:"sub_type"` // 模型类型
}

type ModelOwner struct {
	OwnerName  consts.ModelOwner `json:"owner_name"`  // 提供商
	APIBase    string            `json:"api_base"`    // 接口地址 如：https://api.qwen.com
	APIKey     string            `json:"api_key"`     // 接口密钥 如：sk-xxxx
	APIVersion string            `json:"api_version"` // 接口版本 如：2023-05-15
	APIHeader  string            `json:"api_header"`  // 接口头 如：Authorization: Bearer sk-xxxx
	Models     []*Model          `json:"models"`      // 模型列表
}

type CheckModelResp struct {
	Error   string `json:"error"`
	Content string `json:"content"`
}

type Model struct {
	ID        string            `json:"id"`         // 模型的名字
	Object    string            `json:"object"`     // 总是model
	Created   int               `json:"created"`    // 创建时间
	OwnedBy   consts.ModelOwner `json:"owned_by"`   // 提供商
	ModelType consts.ModelType  `json:"model_type"` // 模型类型
}

var Models []Model
var ModelOwners map[consts.ModelOwner]*ModelOwner
var TypeModelMap map[consts.ModelType][]*Model

func getModelsByOwner(owner consts.ModelOwner) []*Model {
	var ms []*Model
	for i := range Models {
		if Models[i].OwnedBy == owner {
			ms = append(ms, &Models[i])
		}
	}
	return ms
}

func init() {
	initModels()
	initModelOwners()
}

// getBaiZhiCloudModels 返回百智云模型列表
func getBaiZhiCloudModels() []Model {
	return []Model{
		{ID: "qwen-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-72b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-14b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-coder-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-r1", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "kimi-k2-0711-preview", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-coder-480b-a35b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-v3", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen-turbo", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-32b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiZhiCloud, ModelType: consts.ModelTypeChat},
	}
}

// getDeepSeekModels 返回Deepseek模型列表
func getDeepSeekModels() []Model {
	return []Model{
		{ID: "deepseek-chat", Object: "model", OwnedBy: consts.ModelOwnerDeepSeek, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-reasoner", Object: "model", OwnedBy: consts.ModelOwnerDeepSeek, ModelType: consts.ModelTypeChat},
	}
}

// getHunyuanModels 返回腾讯混元模型列表
func getHunyuanModels() []Model {
	return []Model{
		{ID: "hunyuan-pro", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-vision", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-lite", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-standard-32K", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-standard", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-standard-256k", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-functioncall", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeFunctionCall},
		{ID: "hunyuan-role", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-code", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeCoder},
		{ID: "hunyuan-turbo-vision", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-turbo-latest", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbo", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-large", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-large-longcontext", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-latest", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-20250226", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-20250313", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-t1-latest", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-t1-20250321", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-t1-vision", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-turbos-20250515", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-large-vision", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-t1-20250529", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-20250604", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-vision-20250619", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-t1-vision-20250619", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
		{ID: "hunyuan-a13b", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-t1-20250711", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-turbos-20250716", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeChat},
		{ID: "hunyuan-vision-7b-20250720", Object: "model", OwnedBy: consts.ModelOwnerHunyuan, ModelType: consts.ModelTypeVision},
	}
}

// getBaiLianModels 返回阿里百炼模型列表
func getBaiLianModels() []Model {
	return []Model{
		{ID: "qwen3-coder-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen3-coder-plus-2025-07-22", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-plus-2025-07-14", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-coder-480b-a35b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-mt-turbo", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-mt-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-tts-2025-05-22", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qvq-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qvq-plus-2025-05-15", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qvq-max-2025-05-15", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-4b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-32b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-30b-a3b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-235b-a22b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-14b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-1.7b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-0.6b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-8b", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-vl-max-2025-04-02", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeVision},
		{ID: "qwen-vl-ocr-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeVision},
		{ID: "qwen-vl-ocr", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeVision},
		{ID: "qwen-coder-plus-1106", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-coder-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-coder-plus-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen2.5-coder-3b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen2.5-coder-0.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen2.5-coder-14b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen2.5-coder-32b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-math-turbo", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-3b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-math-1.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-1.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-0.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-32b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-72b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-coder-7b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen2.5-math-7b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-coder-turbo", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-coder-turbo-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-math-plus-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-math-plus-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-math-turbo-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-plus-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-plus-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-turbo-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-math-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-14b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-7b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-math-72b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-coder-turbo-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen-math-turbo-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-latest", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-turbo-0919", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2-1.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2-72b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2-7b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2-0.5b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen2-57b-a14b-instruct", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-long", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-vl-max", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeVision},
		{ID: "qwen-vl-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeVision},
		{ID: "qwen-max-0428", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen1.5-110b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen1.5-0.5b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-1.8b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-1.8b-longcontext-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-7b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-14b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-72b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "codeqwen1.5-7b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeCoder},
		{ID: "qwen1.5-32b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen1.5-72b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-longcontext", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-1201", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen1.5-1.8b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen1.5-14b-chat", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-turbo", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-plus", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-0403", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max-0107", Object: "model", OwnedBy: consts.ModelOwnerBaiLian, ModelType: consts.ModelTypeChat},
	}
}

// getVolcengineModels 返回火山引擎模型列表
func getVolcengineModels() []Model {
	return []Model{
		{ID: "doubao-seed-1.6-250615", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "doubao-seed-1.6-flash-250615", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "doubao-seed-1.6-flash-250715", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "doubao-seed-1.6-thinking-250615", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "doubao-seed-1.6-thinking-250715", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "doubao-1.5-thinking-vision-pro-250428", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeVision},
		{ID: "Doubao-1.5-thinking-pro-250415", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-r1-250528", Object: "model", OwnedBy: consts.ModelOwnerVolcengine, ModelType: consts.ModelTypeChat},
	}
}

// getOpenAIModels 返回OpenAI模型列表
func getOpenAIModels() []Model {
	return []Model{
		{ID: "text-embedding-ada-002", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ID: "whisper-1", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-3.5-turbo", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "tts-1", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-3.5-turbo-16k", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "davinci-002", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "babbage-002", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-3.5-turbo-instruct", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-3.5-turbo-instruct-0914", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "dall-e-3", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeVision},
		{ID: "dall-e-2", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeVision},
		{ID: "gpt-3.5-turbo-1106", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "tts-1-hd", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "tts-1-1106", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "tts-1-hd-1106", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "text-embedding-3-small", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ID: "text-embedding-3-large", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeEmbedding},
		{ID: "gpt-3.5-turbo-0125", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-2024-05-13", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-mini-2024-07-18", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-mini", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-2024-08-06", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o1-mini-2024-09-12", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o1-mini", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-audio-preview-2024-10-01", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-audio-preview", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "omni-moderation-latest", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "omni-moderation-2024-09-26", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-audio-preview-2024-12-17", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-audio-preview-2024-12-17", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-audio-preview", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-2024-11-20", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-search-preview-2025-03-11", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-search-preview", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-search-preview-2025-03-11", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-search-preview", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-transcribe", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-transcribe", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4o-mini-tts", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
		{ID: "gpt-4.1-2025-04-14", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-mini-2025-04-14", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-mini", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-nano-2025-04-14", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-nano", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-image-1", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeVision},
		{ID: "gpt-4o-audio-preview-2025-06-03", Object: "model", OwnedBy: consts.ModelOwnerOpenAI, ModelType: consts.ModelTypeFunctionCall},
	}
}

// getSiliconFlowModels 返回硅基流动模型列表
func getSiliconFlowModels() []Model {
	return []Model{
		{ID: "stabilityai/stable-diffusion-xl-base-1.0", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "THUDM/glm-4-9b-chat", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "internlm/internlm2_5-7b-chat", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "BAAI/bge-large-en-v1.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "BAAI/bge-large-zh-v1.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "Pro/Qwen/Qwen2-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/THUDM/glm-4-9b-chat", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "black-forest-labs/FLUX.1-schnell", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "black-forest-labs/FLUX.1-dev", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "FunAudioLLM/SenseVoiceSmall", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "netease-youdao/bce-embedding-base_v1", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "BAAI/bge-m3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "netease-youdao/bce-reranker-base_v1", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "BAAI/bge-reranker-v2-m3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "deepseek-ai/DeepSeek-V2.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-72B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-14B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-32B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/black-forest-labs/FLUX.1-schnell", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-Coder-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "Pro/Qwen/Qwen2.5-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-72B-Instruct-128K", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2-VL-72B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Pro/BAAI/bge-m3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "stabilityai/stable-diffusion-3-5-large", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "LoRA/Qwen/Qwen2.5-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "fishaudio/fish-speech-1.4", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "Pro/Qwen/Qwen2.5-Coder-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "LoRA/Qwen/Qwen2.5-72B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-Coder-32B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "Pro/BAAI/bge-reranker-v2-m3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "RVC-Boss/GPT-SoVITS", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "fishaudio/fish-speech-1.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "black-forest-labs/FLUX.1-pro", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "LoRA/Qwen/Qwen2.5-14B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "LoRA/Qwen/Qwen2.5-32B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "FunAudioLLM/CosyVoice2-0.5B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "deepseek-ai/deepseek-vl2", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Qwen/QVQ-72B-Preview", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "LoRA/black-forest-labs/FLUX.1-dev", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen2.5-VL-72B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Pro/Qwen/Qwen2.5-VL-7B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "deepseek-ai/DeepSeek-V3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-ai/DeepSeek-R1", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/deepseek-ai/DeepSeek-R1-Distill-Qwen-7B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-ai/DeepSeek-R1-Distill-Qwen-14B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-ai/DeepSeek-R1-Distill-Qwen-32B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-ai/DeepSeek-R1-Distill-Qwen-7B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "Pro/deepseek-ai/DeepSeek-R1", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/deepseek-ai/DeepSeek-V3", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Kwai-Kolors/Kolors", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Qwen/QwQ-32B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Wan-AI/Wan2.1-T2V-14B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Wan-AI/Wan2.1-T2V-14B-Turbo", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Wan-AI/Wan2.1-I2V-14B-720P", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Wan-AI/Wan2.1-I2V-14B-720P-Turbo", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "Qwen/Qwen2.5-VL-32B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeVision},
		{ID: "THUDM/GLM-Z1-32B-0414", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "THUDM/GLM-4-32B-0414", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "THUDM/GLM-Z1-9B-0414", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "THUDM/GLM-4-9B-0414", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "THUDM/GLM-Z1-Rumination-32B-0414", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-8B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-14B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-32B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-30B-A3B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-235B-A22B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Tongyi-Zhiwen/QwenLong-L1-32B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-Embedding-8B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "Qwen/Qwen3-Embedding-4B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "Qwen/Qwen3-Embedding-0.6B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeEmbedding},
		{ID: "MiniMaxAI/MiniMax-M1-80k", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-Reranker-0.6B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "Qwen/Qwen3-Reranker-4B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "Qwen/Qwen3-Reranker-8B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeReranker},
		{ID: "fnlp/MOSS-TTSD-v0.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeFunctionCall},
		{ID: "moonshotai/Kimi-Dev-72B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "tencent/Hunyuan-A13B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "baidu/ERNIE-4.5-300B-A47B", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "THUDM/GLM-4.1V-9B-Thinking", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/THUDM/GLM-4.1V-9B-Thinking", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "ascend-tribe/pangu-pro-moe", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "moonshotai/Kimi-K2-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Pro/moonshotai/Kimi-K2-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-235B-A22B-Instruct-2507", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-Coder-480B-A35B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
		{ID: "Qwen/Qwen3-235B-A22B-Thinking-2507", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "zai-org/GLM-4.5", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "zai-org/GLM-4.5-Air", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-30B-A3B-Instruct-2507", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-30B-A3B-Thinking-2507", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeChat},
		{ID: "Qwen/Qwen3-Coder-30B-A3B-Instruct", Object: "model", OwnedBy: consts.ModelOwnerSiliconFlow, ModelType: consts.ModelTypeCoder},
	}
}

func getMoonshotModels() []Model {
	return []Model{
		{ID: "moonshot-v1-auto", Object: "model", OwnedBy: consts.ModelOwnerMoonshot, ModelType: consts.ModelTypeChat},
		{ID: "moonshot-v1-8k", Object: "model", OwnedBy: consts.ModelOwnerMoonshot, ModelType: consts.ModelTypeChat},
		{ID: "moonshot-v1-32k", Object: "model", OwnedBy: consts.ModelOwnerMoonshot, ModelType: consts.ModelTypeChat},
		{ID: "moonshot-v1-128k", Object: "model", OwnedBy: consts.ModelOwnerMoonshot, ModelType: consts.ModelTypeChat},
	}

}

func getAzureOpenAIModels() []Model {
	return []Model{
		{ID: "gpt-4", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-mini", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4o-nano", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-mini", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "gpt-4.1-nano", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o1", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o1-mini", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o3", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o3-mini", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
		{ID: "o4-mini", Object: "model", OwnedBy: consts.ModelOwnerAzureOpenAI, ModelType: consts.ModelTypeChat},
	}
}

// getGeminiModels 返回Google Gemini模型列表
func getGeminiModels() []Model {
	return []Model{
		{ID: "gemini-2.5-pro", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.5-flash", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.5-flash-lite-preview-06-17", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.5-flash-preview-tts", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.5-pro-preview-tts", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.0-flash", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-2.0-flash-lite", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-1.5-flash", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-1.5-flash-8b", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-1.5-pro", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeChat},
		{ID: "gemini-embedding-001", Object: "model", OwnedBy: consts.ModelOwnerGemini, ModelType: consts.ModelTypeEmbedding},
	}
}

// getZhiPuModels 返回智谱模型列表
func getZhiPuModels() []Model {
	return []Model{
		{ID: "glm-4.5", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4.5-x", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4.5-air", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4.5-airx", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4.5-flash", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-plus", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-air-250414", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-airx", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-long", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-flashx-250414", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4-flash-250414", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-z1-air", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-z1-airx", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-z1-flashx", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-z1-flash", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeChat},
		{ID: "glm-4v-plus-0111", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeVision},
		{ID: "glm-4v-flash", Object: "model", OwnedBy: consts.ModelOwnerZhiPu, ModelType: consts.ModelTypeVision},
	}
}

func initModels() {
	Models = make([]Model, 0, 200)

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

func initModelOwners() {
	// 初始化模型提供商及其模型
	ModelOwners = map[consts.ModelOwner]*ModelOwner{
		consts.ModelOwnerBaiZhiCloud: {
			OwnerName: consts.ModelOwnerBaiZhiCloud,
			Models:    getModelsByOwner(consts.ModelOwnerBaiZhiCloud),
			APIBase:   "https://model-square.app.baizhi.cloud/v1",
		},
		consts.ModelOwnerDeepSeek: {
			OwnerName: consts.ModelOwnerDeepSeek,
			Models:    getModelsByOwner(consts.ModelOwnerDeepSeek),
			APIBase:   "https://api.deepseek.com/v1",
		},
		consts.ModelOwnerSiliconFlow: {
			OwnerName: consts.ModelOwnerSiliconFlow,
			Models:    getModelsByOwner(consts.ModelOwnerSiliconFlow),
			APIBase:   "https://api.siliconflow.cn/v1",
		},
		consts.ModelOwnerOpenAI: {
			OwnerName: consts.ModelOwnerOpenAI,
			Models:    getModelsByOwner(consts.ModelOwnerOpenAI),
			APIBase:   "https://api.openai.com/v1",
		},
		consts.ModelOwnerOllama: {
			OwnerName: consts.ModelOwnerOllama,
			Models:    getModelsByOwner(consts.ModelOwnerOllama),
			APIBase:   "http://localhost:11434",
		},
		consts.ModelOwnerMoonshot: {
			OwnerName: consts.ModelOwnerMoonshot,
			Models:    getModelsByOwner(consts.ModelOwnerMoonshot),
			APIBase:   "https://api.moonshot.cn/v1",
		},
		consts.ModelOwnerAzureOpenAI: {
			OwnerName: consts.ModelOwnerAzureOpenAI,
			Models:    getModelsByOwner(consts.ModelOwnerAzureOpenAI),
		},
		consts.ModelOwnerHunyuan: {
			OwnerName: consts.ModelOwnerHunyuan,
			Models:    getModelsByOwner(consts.ModelOwnerHunyuan),
			APIBase:   "https://api.hunyuan.cloud.tencent.com/v1",
		},
		consts.ModelOwnerBaiLian: {
			OwnerName: consts.ModelOwnerBaiLian,
			Models:    getModelsByOwner(consts.ModelOwnerBaiLian),
			APIBase:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
		},
		consts.ModelOwnerVolcengine: {
			OwnerName: consts.ModelOwnerVolcengine,
			Models:    getModelsByOwner(consts.ModelOwnerVolcengine),
			APIBase:   "https://ark.cn-beijing.volces.com/api/v3",
		},
	}

	// 初始化按类型分组的模型映射
	TypeModelMap = make(map[consts.ModelType][]*Model)
	for i := range Models {
		model := &Models[i]
		TypeModelMap[model.ModelType] = append(TypeModelMap[model.ModelType], model)
	}
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func From(modelOwner *ModelOwner) []ProviderModelListItem {
	var result []ProviderModelListItem
	for _, model := range modelOwner.Models {
		result = append(result, ProviderModelListItem{
			Model: model.ID,
		})
	}
	return result
}
