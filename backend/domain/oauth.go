package domain

import "github.com/chaitin/ModelKit/backend/consts"

type OAuther interface {
	GetAuthorizeURL() (state string, url string)
	GetUserInfo(code string) (*OAuthUserInfo, error)
}

type OAuthConfig struct {
	Debug        bool
	Platform     consts.UserPlatform
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	AuthorizeURL string
	TokenURL     string
	UserInfoURL  string
	IDField      string
	NameField    string
	AvatarField  string
	EmailField   string
}

type OAuthUserInfo struct {
	ID        string `json:"id"`
	UnionID   string `json:"union_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type OAuthSignUpOrInReq struct {
	Source      consts.LoginSource  `json:"source" query:"source" validate:"required" default:"plugin"` // 登录来源 plugin: 插件 browser: 浏览器; 默认为 plugin
	Platform    consts.UserPlatform `json:"platform" query:"platform" validate:"required"`              // 第三方平台 dingtalk
	SessionID   string              `json:"session_id" query:"session_id"`                              // 会话ID
	RedirectURL string              `json:"redirect_url" query:"redirect_url"`                          // 登录成功后跳转的 URL
	InviteCode  string              `json:"inviate_code" query:"inviate_code"`                          // 邀请码
	BaseURL     string              `json:"-"`
}

func (o OAuthSignUpOrInReq) OAuthKind() consts.OAuthKind {
	if o.InviteCode == "" {
		return consts.OAuthKindLogin
	}
	return consts.OAuthKindInvite
}

type OAuthCallbackReq struct {
	State   string `json:"state" query:"state" validate:"required"`
	Code    string `json:"code" query:"code" validate:"required"`
	IP      string `json:"-"`
	BaseURL string `json:"-"`
}

type OAuthURLResp struct {
	URL string `json:"url"`
}

type OAuthState struct {
	Source      consts.LoginSource  `json:"source"`                                        // 登录来源 plugin: 插件 browser: 浏览器; 默认为 plugin
	Kind        consts.OAuthKind    `json:"kind" query:"kind" validate:"required"`         // invite: 邀请登录 login: 登录
	SessionID   string              `json:"session_id"`                                    // 会话ID
	Platform    consts.UserPlatform `json:"platform" query:"platform" validate:"required"` // 第三方平台 dingtalk
	RedirectURL string              `json:"redirect_url" query:"redirect_url"`             // 登录成功后跳转的 URL
	InviteCode  string              `json:"inviate_code"`                                  // 邀请码
}

type OAuthAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type GetAccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURL  string `json:"redirect_uri"`
}
