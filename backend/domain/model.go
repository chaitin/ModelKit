package domain

import (
	"context"

	"github.com/chaitin/ModelKit/backend/consts"
)

type ModelUsecase interface {
	CheckModel(ctx context.Context, req *CheckModelReq) (*Model, error)
	ListModel(ctx context.Context, req *ListModelReq) ([]*Model, error)
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
var ModelOwners map[consts.ModelOwner]ModelOwner
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
}

func initModels() {
	Models = []Model{
		// 百智云模型
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
		// Deepseek
		{ID: "deepseek-chat", Object: "model", OwnedBy: consts.ModelOwnerDeepSeek, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-reasoner", Object: "model", OwnedBy: consts.ModelOwnerDeepSeek, ModelType: consts.ModelTypeChat},
	}

	// 初始化模型提供商及其模型
	ModelOwners = map[consts.ModelOwner]ModelOwner{
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
