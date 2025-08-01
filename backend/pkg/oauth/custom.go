package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/chaitin/ModelKit/backend/domain"
)

type CustomOAuth struct {
	cfg   domain.OAuthConfig
	oauth *oauth2.Config
}

func NewCustomOAuth(config domain.OAuthConfig) domain.OAuther {
	c := &CustomOAuth{
		cfg: config,
		oauth: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthorizeURL,
				TokenURL: config.TokenURL,
			},
			RedirectURL: config.RedirectURI,
			Scopes:      config.Scopes,
		},
	}

	return c
}

// GetAuthorizeURL implements domain.OAuther.
func (c *CustomOAuth) GetAuthorizeURL() (string, string) {
	state := uuid.NewString()
	url := c.oauth.AuthCodeURL(state)
	return state, url
}

// GetUserInfo implements domain.OAuther.
func (c *CustomOAuth) GetUserInfo(code string) (*domain.OAuthUserInfo, error) {
	info, err := c.getUserInfo(code)
	if err != nil {
		return nil, err
	}
	return &domain.OAuthUserInfo{
		ID:        fmt.Sprint(info[c.cfg.IDField]),
		AvatarURL: fmt.Sprint(info[c.cfg.AvatarField]),
		Name:      fmt.Sprint(info[c.cfg.NameField]),
		Email:     fmt.Sprint(info[c.cfg.EmailField]),
	}, nil
}

type UserInfo map[string]any

func (c *CustomOAuth) getUserInfo(code string) (UserInfo, error) {
	token, err := c.oauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	client := c.oauth.Client(context.Background(), token)
	res, err := client.Get(c.cfg.UserInfoURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	err = json.Unmarshal(buf, &userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
