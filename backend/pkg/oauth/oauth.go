package oauth

import (
	"fmt"

	"github.com/chaitin/ModelKit/backend/consts"
	"github.com/chaitin/ModelKit/backend/domain"
)

func NewOAuther(config domain.OAuthConfig) (domain.OAuther, error) {
	switch config.Platform {
	case consts.UserPlatformDingTalk:
		return NewDingTalk(config), nil
	case consts.UserPlatformCustom:
		return NewCustomOAuth(config), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", config.Platform)
	}
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}
