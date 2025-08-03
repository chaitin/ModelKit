package pkg

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/text/language"

	"github.com/GoYoko/web"
	"github.com/GoYoko/web/locale"

	"github.com/chaitin/ModelKit/backend/config"
	"github.com/chaitin/ModelKit/backend/errcode"
	"github.com/chaitin/ModelKit/backend/pkg/logger"
	"github.com/chaitin/ModelKit/backend/pkg/store"
)

var Provider = wire.NewSet(
	NewWeb,
	logger.NewLogger,
	store.NewEntDB,
)

func NewWeb(cfg *config.Config) *web.Web {
	w := web.New()
	l := locale.NewLocalizerWithFile(language.Chinese, errcode.LocalFS, []string{"locale.zh.toml"})
	w.SetLocale(l)
	if cfg.Debug {
		w.Use(middleware.Logger())
	}
	return w
}
