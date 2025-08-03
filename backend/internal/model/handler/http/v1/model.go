package v1

import (
	"context"
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

	g.GET("/provider/supported", web.BindHandler(m.GetProviderModelList))
	g.GET("", web.BaseHandler(m.ListModel))
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

// List 获取模型列表
//
//	@Tags			Model
//	@Summary		获取模型列表
//	@Description	获取模型列表
//	@ID				list-model
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{data=domain.AllModelResp}
//	@Router			/api/v1/model [get]
func (h *ModelHandler) ListModel(c *web.Context) error {
	models, err := h.usecase.ListModel(c.Request().Context())
	if err != nil {
		return err
	}
	return c.Success(models)
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

// GetProviderModelList 获取供应商支持的模型列表
//
//	@Tags			Model
//	@Summary		获取供应商支持的模型列表
//	@Description	获取供应商支持的模型列表
//	@ID				get-provider-model-list
//	@Accept			json
//	@Produce		json
//	@Param			param	query		domain.GetProviderModelListReq	true	"模型类型"
//	@Success		200		{object}	web.Resp{data=domain.GetProviderModelListResp}
//	@Router			/api/v1/model/provider/supported [get]
func (h *ModelHandler) GetProviderModelList(c *web.Context, req domain.GetProviderModelListReq) error {
	resp, err := h.usecase.GetProviderModelList(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.Success(resp)
}

func (h *ModelHandler) InitModel() error {
	return h.usecase.InitModel(context.Background())
}
