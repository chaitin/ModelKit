package domain

import (
	"github.com/chaitin/ModelKit/v2/consts"
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
	Error  string          `json:"error"`
}

type ModelListItem struct {
	Model string `json:"model"`
}

type CheckModelReq struct {
	Provider   string      `json:"provider" query:"provider" validate:"required,oneof=OpenAI Ollama DeepSeek SiliconFlow Moonshot Other AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Gemini ZhiPu AiHubMix"`
	Model      string      `json:"model" query:"model_name" validate:"required"`
	BaseURL    string      `json:"base_url" query:"base_url" validate:"required"`
	APIKey     string      `json:"api_key" query:"api_key"`
	APIHeader  string      `json:"api_header" query:"api_header"`
	APIVersion string      `json:"api_version" query:"api_version"` // for azure openai
	Type       string      `json:"type" query:"model_type" validate:"required,oneof=chat embedding rerank llm"`
	Param      *ModelParam `json:"param" query:"param"`
}

type CheckModelResp struct {
	Error   string `json:"error"`
	Content string `json:"content"`
}

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

func From(ModelProvider ModelProvider) []ModelListItem {
	var result []ModelListItem
	for _, model := range ModelProvider.Models {
		result = append(result, ModelListItem{
			Model: model.ModelName,
		})
	}
	return result
}

// ModelResponseParser 定义模型响应解析器接口
type ModelResponseParser interface {
	ParseModels() []ModelListItem
}

type GithubModel struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Registry      string   `json:"registry"`
	Publisher     string   `json:"publisher"`
	Summary       string   `json:"summary"`
	RateLimitTier string   `json:"rate_limit_tier"`
	HTMLURL       string   `json:"html_url"`
	Version       string   `json:"version"`
	Capabilities  []string `json:"capabilities"`
	Limits        struct {
		MaxInputTokens  int `json:"max_input_tokens"`
		MaxOutputTokens int `json:"max_output_tokens"`
	} `json:"limits"`
	Tags                      []string `json:"tags"`
	SupportedInputModalities  []string `json:"supported_input_modalities"`
	SupportedOutputModalities []string `json:"supported_output_modalities"`
}

type GithubResp []GithubModel

// ParseModels 实现ModelResponseParser接口
func (g *GithubResp) ParseModels() []ModelListItem {
	var models []ModelListItem
	for _, item := range *g {
		models = append(models, ModelListItem{Model: item.ID})
	}
	return models
}
