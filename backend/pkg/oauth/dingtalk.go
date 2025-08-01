package oauth

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/request"
)

type DingTalk struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	client       *request.Client
}

var _ domain.OAuther = &DingTalk{}

func NewDingTalk(config domain.OAuthConfig) *DingTalk {
	client := request.NewClient("https", "api.dingtalk.com", 30*time.Second)
	client.SetDebug(config.Debug)
	return &DingTalk{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURI:  config.RedirectURI,
		client:       client,
	}
}

// GetUserInfo implements domain.OAuther.
func (d *DingTalk) GetUserInfo(code string) (*domain.OAuthUserInfo, error) {
	accessToken, err := d.getAccessToken(code)
	if err != nil {
		return nil, err
	}
	return d.getUserInfo(accessToken)
}

type dingtalkAccessTokenResp struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int    `json:"expireIn"`
	RefreshToken string `json:"refreshToken"`
	CorpID       string `json:"corpId"`
}

type dingtalkAccessTokenReq struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrantType    string `json:"grantType"`
}

func (d *DingTalk) getAccessToken(code string) (string, error) {
	resp, err := request.Post[dingtalkAccessTokenResp](d.client, "/v1.0/oauth2/userAccessToken", dingtalkAccessTokenReq{
		ClientID:     d.ClientID,
		ClientSecret: d.ClientSecret,
		Code:         code,
		GrantType:    "authorization_code",
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

type dingtalkUserInfoResp struct {
	Nick      string `json:"nick"`
	AvatarURL string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

func (d *DingTalk) getUserInfo(accessToken string) (*domain.OAuthUserInfo, error) {
	resp, err := request.Get[dingtalkUserInfoResp](d.client, "/v1.0/contact/users/me", request.WithHeader(request.Header{
		"x-acs-dingtalk-access-token": accessToken,
	}))
	if err != nil {
		return nil, err
	}
	return &domain.OAuthUserInfo{
		ID:        resp.OpenID,
		UnionID:   resp.UnionID,
		Name:      resp.Nick,
		Email:     resp.Email,
		AvatarURL: resp.AvatarURL,
	}, nil
}

func (d *DingTalk) GetAuthorizeURL() (string, string) {
	state := uuid.NewString()
	url := fmt.Sprintf("https://login.dingtalk.com/oauth2/auth?response_type=code&scope=openid&client_id=%s&prompt=consent&state=%s&redirect_uri=%s", d.ClientID, state, d.RedirectURI)
	return state, url
}
