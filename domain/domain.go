package domain

import (
	"github.com/chaitin/ModelKit/consts"
)

type ModelKit interface {
	// CheckModel(ctx context.Context, req *CheckModelReq) (*Model, error)
}

type ModelListReq struct {
	Provider  string `json:"provider" query:"provider" validate:"required"`
	BaseURL   string `json:"base_url" query:"base_url" validate:"required"`
	APIKey    string `json:"api_key" query:"api_key"`
	APIHeader string `json:"api_header" query:"api_header"`
	Type      string `json:"type" query:"type" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

type ModelListResp struct {
	Models []ModelListItem `json:"models"`
}

type ModelListItem struct {
	Model string `json:"model"`
}

type CheckModelReq struct {
	Provider   string `json:"provider" validate:"required,oneof=OpenAI Ollama DeepSeek SiliconFlow Moonshot Other AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu"`
	Model      string `json:"model" validate:"required"`
	BaseURL    string `json:"base_url" validate:"required"`
	APIKey     string `json:"api_key"`
	APIHeader  string `json:"api_header"`
	APIVersion string `json:"api_version"` // for azure openai
	Type       string `json:"type" validate:"required,oneof=chat embedding rerank"`
}

type CheckModelResp struct {
	Error   string `json:"error"`
	Content string `json:"content"`
}

type ModelMetadata struct {
	ModelName string               `json:"id"`         // 模型的名字
	Object    string               `json:"object"`     // 总是model
	Created   int                  `json:"created"`    // 创建时间
	Provider  consts.ModelProvider `json:"provider"`   // 提供商
	ModelType consts.ModelType     `json:"model_type"` // 模型类型

	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	APIHeader  string `json:"api_header"`
	APIVersion string `json:"api_version"` // for azure openai
}

var Models []ModelMetadata

func getModelsByOwner(owner consts.ModelProvider) []ModelMetadata {
	var ms []ModelMetadata
	for i := range Models {
		if Models[i].Provider == owner {
			ms = append(ms, Models[i])
		}
	}
	return ms
}

func init() {
	initModels()
	initModelProviders()
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func From(ModelProvider ModelProvider) []ModelListItem {
	var result []ModelListItem
	for _, model := range ModelProvider.Models {
		result = append(result, ModelListItem{
			Model: model.ModelName,
		})
	}
	return result
}
