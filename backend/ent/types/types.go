package types

type DingtalkOAuth struct {
	Enable       bool   `json:"enable"`        // 钉钉OAuth开关
	ClientID     string `json:"client_id"`     // 钉钉客户端ID
	ClientSecret string `json:"client_secret"` // 钉钉客户端密钥
}

type BaseURL struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type CustomOAuth struct {
	Enable         bool     `json:"enable"`           // 自定义OAuth开关
	ClientID       string   `json:"client_id"`        // 自定义客户端ID
	ClientSecret   string   `json:"client_secret"`    // 自定义客户端密钥
	AuthorizeURL   string   `json:"authorize_url"`    // 自定义OAuth授权URL
	AccessTokenURL string   `json:"access_token_url"` // 自定义OAuth访问令牌URL
	UserInfoURL    string   `json:"userinfo_url"`     // 自定义OAuth用户信息URL
	Scopes         []string `json:"scopes"`           // 自定义OAuth Scope列表
	IDField        string   `json:"id_field"`         // 用户信息回包中的ID字段名
	NameField      string   `json:"name_field"`       // 用户信息回包中的用户名字段名`
	AvatarField    string   `json:"avatar_field"`     // 用户信息回包中的头像URL字段名`
	EmailField     string   `json:"email_field"`      // 用户信息回包中的邮箱字段名
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
