package domain

import (
	"context"

	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/chaitin/ModelKit/backend/db"
)

type ModelUsecase interface {
	UpdateModel(ctx context.Context, req *UpdateModelReq) (*Model, error)
	CheckModel(ctx context.Context, req *CheckModelReq) (*Model, error)
}

type ModelRepo interface {
	GetWithCache(ctx context.Context, modelType consts.ModelType) (*db.Model, error)
	UpdateModel(ctx context.Context, id string, fn func(tx *db.Tx, old *db.Model, up *db.ModelUpdateOne) error) (*db.Model, error)
}

type CheckModelReq struct {
	Type       consts.ModelType     `json:"type" validate:"required,oneof=llm coder embedding rerank"`
	Provider   consts.ModelProvider `json:"provider" validate:"required"`   // 提供商
	ModelName  string               `json:"model_name" validate:"required"` // 模型名称
	APIBase    string               `json:"api_base" validate:"required"`   // 接口地址
	APIKey     string               `json:"api_key" validate:"required"`    // 接口密钥
	APIVersion string               `json:"api_version"`
	APIHeader  string               `json:"api_header"`
}

type CreateModelReq struct {
	UserID     string               `json:"-"`
	ShowName   string               `json:"show_name"`                                                                                                                               // 模型显示名称
	ModelName  string               `json:"model_name" validate:"required"`                                                                                                          // 模型名称 如: deepseek-v3
	Provider   consts.ModelProvider `json:"provider" validate:"required,oneof=SiliconFlow OpenAI Ollama DeepSeek Moonshot AzureOpenAI BaiZhiCloud Hunyuan BaiLian Volcengine Other"` // 提供商
	APIBase    string               `json:"api_base" validate:"required"`                                                                                                            // 接口地址 如：https://api.qwen.com
	APIKey     string               `json:"api_key"`                                                                                                                                 // 接口密钥 如：sk-xxxx
	APIVersion string               `json:"api_version"`
	APIHeader  string               `json:"api_header"`
	ModelType  consts.ModelType     `json:"model_type"` // 模型类型 llm:对话模型 coder:代码模型
	Param      *ModelParam          `json:"param"`      // 高级参数
}

type ModelParam struct {
	R1Enabled          bool `json:"r1_enabled"`
	MaxTokens          int  `json:"max_tokens"`
	ContextWindow      int  `json:"context_window"`
	SupprtImages       bool `json:"support_images"`
	SupportComputerUse bool `json:"support_computer_use"`
	SupportPromptCache bool `json:"support_prompt_cache"`
}

func DefaultModelParam() *ModelParam {
	return &ModelParam{
		R1Enabled:          false,
		MaxTokens:          8192,
		ContextWindow:      64000,
		SupprtImages:       false,
		SupportComputerUse: false,
		SupportPromptCache: false,
	}
}

type UpdateModelReq struct {
	ID         string                `json:"id"`                                                                                                                                      // 模型ID
	ShowName   *string               `json:"show_name"`                                                                                                                               // 模型显示名称
	APIKey     *string               `json:"api_key"`                                                                                                                                 // 接口密钥 如：sk-xxxx
	APIVersion *string               `json:"api_version"`
	APIHeader  *string               `json:"api_header"`
	Status     *consts.ModelStatus   `json:"status"`          // 状态 active:启用 inactive:禁用
}

type Model struct {
	ID         string               `json:"id"`          // 模型ID
	ModelName  string               `json:"model_name"`  // 模型名称 如: deepseek-v3
	Provider   consts.ModelProvider `json:"provider"`    // 提供商
	APIBase    string               `json:"api_base"`    // 接口地址 如：https://api.qwen.com
	APIKey     string               `json:"api_key"`     // 接口密钥 如：sk-xxxx
	APIVersion string               `json:"api_version"` // 接口版本 如：2023-05-15
	APIHeader  string               `json:"api_header"`  // 接口头 如：Authorization: Bearer sk-xxxx
	ModelType  consts.ModelType     `json:"model_type"`  // 模型类型 llm:对话模型 coder:代码模型
	CreatedAt  int64                `json:"created_at"`  // 创建时间
	UpdatedAt  int64                `json:"updated_at"`  // 更新时间
}

func (m *Model) From(e *db.Model) *Model {
	if e == nil {
		return m
	}

	m.ID = e.ID.String()
	m.ModelName = e.ModelName
	m.Provider = e.Provider
	m.APIBase = e.APIBase
	m.APIKey = e.APIKey
	m.APIVersion = e.APIVersion
	m.APIHeader = e.APIHeader
	m.ModelType = e.ModelType
	m.CreatedAt = e.CreatedAt.Unix()
	m.UpdatedAt = e.UpdatedAt.Unix()

	return m
}

type CheckModelResp struct {
	Error   string `json:"error"`
	Content string `json:"content"`
}

type Setting struct {
	EnableSSO            bool   `json:"enable_sso"`             // 是否开启SSO
	ForceTwoFactorAuth   bool   `json:"force_two_factor_auth"`  // 是否强制两步验证
	DisablePasswordLogin bool   `json:"disable_password_login"` // 是否禁用密码登录
	EnableAutoLogin      bool   `json:"enable_auto_login"`      // 是否开启自动登录
	BaseURL              string `json:"base_url,omitempty"`     // base url 配置，为了支持前置代理
	CreatedAt            int64  `json:"created_at"`             // 创建时间
	UpdatedAt            int64  `json:"updated_at"`             // 更新时间
}
