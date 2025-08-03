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

	g.POST("/check", web.BindHandler(m.CheckModel))
	g.PUT("", web.BindHandler(m.UpdateModel))
	g.GET("", web.BindHandler(m.GetModel))
	g.GET("/list", web.BindHandler(m.ListModel))

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

// Get 获取模型
//
//	@Tags			Model
//	@Summary		获取模型
//	@Description	根据模型名称和提供商获取模型信息
//	@ID				get-model
//	@Accept			json
//	@Produce		json
//	@Param			req	query	domain.GetModelReq	true	"模型"
//	@Success		200			{object}	web.Resp{data=domain.Model}
//	@Router			/api/v1/model [get]
func (h *ModelHandler) GetModel(c *web.Context, req domain.GetModelReq) error {
	m, err := h.usecase.GetModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.Success(m)
}

// List 获取模型列表
//
//	@Tags			Model
//	@Summary		获取模型列表
//	@Description	根据筛选条件获取模型列表，支持分页，按创建时间降序排列
//	@ID				list-model
//	@Accept			json
//	@Produce		json
//	@Param			req	query	domain.ListModelReq	true	"模型"
//	@Success		200			{object}	web.Resp{data=[]domain.Model}
//	@Router			/api/v1/model/list [get]
func (h *ModelHandler) ListModel(c *web.Context, req domain.ListModelReq) error {
	models, err := h.usecase.ListModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.Success(models)
}
