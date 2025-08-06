// pandawiki特供类型
package domain

import (
	"time"
)

type PandaModel struct {
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	APIKey     string `json:"api_key"`
	APIHeader  string `json:"api_header"`
	BaseURL    string `json:"base_url"`
	APIVersion string `json:"api_version"` // for azure openai
	Type       string `json:"type" gorm:"default:chat;uniqueIndex"`

	IsActive bool `json:"is_active" gorm:"default:false"`

	PromptTokens     uint64 `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens uint64 `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      uint64 `json:"total_tokens" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PandaGetProviderModelListReq struct {
	Provider  string `json:"provider" query:"provider" validate:"required,oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu"`
	BaseURL   string `json:"base_url" query:"base_url" validate:"required"`
	APIKey    string `json:"api_key" query:"api_key"`
	APIHeader string `json:"api_header" query:"api_header"`
	Type      string `json:"type" query:"type" validate:"required,oneof=chat embedding rerank"`
}
