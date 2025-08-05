package domain

import "github.com/chaitin/ModelKit/consts"

type GetProviderModelListReq struct {
	Provider  string           `json:"provider" query:"provider" validate:"required,oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu"`
	BaseURL   string           `json:"base_url" query:"base_url" validate:"required"`
	APIKey    string           `json:"api_key" query:"api_key"`
	APIHeader string           `json:"api_header" query:"api_header"`
	Type      consts.ModelType `json:"type" query:"type" validate:"required,oneof=chat embedding rerank"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

type GetProviderModelListResp struct {
	Models []ProviderModelListItem `json:"models"`
}

type ProviderModelListItem struct {
	Model string `json:"model"`
}
