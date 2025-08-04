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
	Type       consts.ModelType     `json:"type" validate:"required,oneof=llm coder embedding rerank"`
	Provider   consts.ModelProvider `json:"provider" validate:"required"`   // 提供商
	ModelID    string               `json:"model_name" validate:"required"` // 模型名称
	APIBase    string               `json:"api_base" validate:"required"`   // 接口地址
	APIKey     string               `json:"api_key" validate:"required"`    // 接口密钥
	APIVersion string               `json:"api_version"`
	APIHeader  string               `json:"api_header"`
}

type ListModelReq struct {
	Owner string           `json:"owner" query:"owner"` // 提供商
	Type  consts.ModelType `json:"type" query:"type"`   // 模型类型
}

type ModelProvider struct {
	Provider   consts.ModelProvider `json:"provider"`    // 提供商
	APIBase    string               `json:"api_base"`    // 接口地址 如：https://api.qwen.com
	APIKey     string               `json:"api_key"`     // 接口密钥 如：sk-xxxx
	APIVersion string               `json:"api_version"` // 接口版本 如：2023-05-15
	APIHeader  string               `json:"api_header"`  // 接口头 如：Authorization: Bearer sk-xxxx
	Models     []*Model             `json:"models"`      // 模型列表
}

type CheckModelResp struct {
	Error   string `json:"error"`
	Content string `json:"content"`
}

type Model struct {
	ID        string               `json:"id"`         // 模型的名字
	Object    string               `json:"object"`     // 总是model
	Created   int                  `json:"created"`    // 创建时间
	OwnedBy   consts.ModelProvider `json:"owned_by"`   // 提供商
	ModelType consts.ModelType     `json:"model_type"` // 模型类型
}

var Models []Model
var OwnerModelMap map[consts.ModelProvider][]Model
var TypeModelMap map[consts.ModelType][]Model

func init() {
	initModels()
}

func initModels() {
	Models = []Model{
		// 百智云模型
		{ID: "qwen-plus", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-72b-instruct", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-14b-instruct", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen-max", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-coder-plus", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-r1", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "kimi-k2-0711-preview", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen3-coder-480b-a35b-instruct", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeCoder},
		{ID: "deepseek-v3", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen-turbo", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		{ID: "qwen2.5-32b-instruct", Object: "model", OwnedBy: consts.ModelProviderBaiZhiCloud, ModelType: consts.ModelTypeChat},
		// Deepseek
		{ID: "deepseek-chat", Object: "model", OwnedBy: consts.ModelProviderDeepSeek, ModelType: consts.ModelTypeChat},
		{ID: "deepseek-reasoner", Object: "model", OwnedBy: consts.ModelProviderDeepSeek, ModelType: consts.ModelTypeChat},
	}

	// 初始化按提供商分组的模型映射
	OwnerModelMap = make(map[consts.ModelProvider][]Model)
	for _, model := range Models {
		OwnerModelMap[model.OwnedBy] = append(OwnerModelMap[model.OwnedBy], model)
	}

	// 初始化按类型分组的模型映射
	TypeModelMap = make(map[consts.ModelType][]Model)
	for _, model := range Models {
		TypeModelMap[model.ModelType] = append(TypeModelMap[model.ModelType], model)
	}
}
