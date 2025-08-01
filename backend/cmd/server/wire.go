//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"

	"github.com/google/wire"

	"github.com/GoYoko/web"

	"github.com/chaitin/ModelKit/backend/config"
	"github.com/chaitin/ModelKit/backend/db"
	v1 "github.com/chaitin/ModelKit/backend/internal/model/handler/http/v1"
	"github.com/chaitin/ModelKit/backend/pkg/version"
)

type Server struct {
	config        *config.Config
	web           *web.Web
	ent           *db.Client
	logger        *slog.Logger
	modelV1       *v1.ModelHandler
	version       *version.VersionInfo
}

func newServer() (*Server, error) {
	wire.Build(
		wire.Struct(new(Server), "*"),
		appSet,
	)
	return &Server{}, nil
}
