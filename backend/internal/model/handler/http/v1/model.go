package v1

import (
	"log/slog"

	"github.com/GoYoko/web"

	"github.com/chaitin/ModelKit/backend/domain"
)

type ModelHandler struct {
	usecase domain.ModelUsecase
	logger  *slog.Logger
}

func NewModelHandler(
	w *web.Web,
	usecase domain.ModelUsecase,
	logger *slog.Logger,
) *ModelHandler {
	m := &ModelHandler{usecase: usecase, logger: logger.With("handler", "model")}

	g := w.Group("/api/v1/model")

	// g.POST("", web.BindHandler(m.CreateModel))
	g.POST("/check", web.BindHandler(m.CheckModel))
	g.PUT("", web.BindHandler(m.UpdateModel))
	// g.DELETE("", web.BaseHandler(m.DeleteModel))

	return m
}

// Check 检查模型
//
//	@Tags			Model
//	@Summary		检查模型
//	@Description	检查模型
//	@ID				check-model
//	@Accept			json
//	@Produce		json
//	@Param			model	body		domain.CheckModelReq	true	"模型"
//	@Success		200		{object}	web.Resp{data=domain.Model}
//	@Router			/api/v1/model/check [post]
func (h *ModelHandler) CheckModel(c *web.Context, req domain.CheckModelReq) error {
	m, err := h.usecase.CheckModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.Success(m)
}


// Update 更新模型
//
//	@Tags			Model
//	@Summary		更新模型
//	@Description	更新模型
//	@ID				update-model
//	@Accept			json
//	@Produce		json
//	@Param			model	body		domain.UpdateModelReq	true	"模型"
//	@Success		200		{object}	web.Resp{data=domain.Model}
//	@Router			/api/v1/model [put]
func (h *ModelHandler) UpdateModel(c *web.Context, req domain.UpdateModelReq) error {
	m, err := h.usecase.UpdateModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.Success(m)
}
