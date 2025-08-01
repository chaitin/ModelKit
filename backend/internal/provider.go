package internal

import (
	"github.com/google/wire"

	modelv1 "github.com/chaitin/ModelKit/backend/internal/model/handler/http/v1"
	modelrepo "github.com/chaitin/ModelKit/backend/internal/model/repo"
	modelusecase "github.com/chaitin/ModelKit/backend/internal/model/usecase"
	"github.com/chaitin/ModelKit/backend/pkg/version"
)

var Provider = wire.NewSet(
	modelv1.NewModelHandler,
	modelusecase.NewModelUsecase,
	modelrepo.NewModelRepo,
	version.NewVersionInfo,
)
