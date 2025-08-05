package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/internal/usecase"
)

type ModelKit struct {
	usecase domain.ModelUsecase
}

func NewModelKit(
	echo *echo.Echo,
) *ModelKit {
	m := &ModelKit{usecase: usecase.NewModelUsecase()}

	g := echo.Group("/api/v1/model/modelkit")

	g.POST("/check", m.CheckModel)
	g.GET("/models", m.ListModel)

	return m
}

// Check 检查模型
//
//	@Tags			ModelKitModel
//	@Summary		检查模型
//	@Description	检查模型
//	@ID				check-model
//	@Accept			json
//	@Produce		json
//	@Param			model	body		domain.CheckModelReq	true	"模型"
//	@Success		200		{object}	domain.Resp{data=domain.Model}
//	@Router			/api/v1/model/modelkit/check [post]
func (h *ModelKit) CheckModel(c echo.Context) error {
	var req domain.CheckModelReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	m, err := h.usecase.CheckModel(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Resp{
			Code:    1,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, domain.Resp{
		Code:    0,
		Message: "success",
		Data:    m,
	})
}

// List 获取模型列表
//
//	@Tags			ModelKitModel
//	@Summary		获取模型列表
//	@Description	根据筛选条件获取模型列表，支持分页，按创建时间降序排列
//	@ID				list-model
//	@Accept			json
//	@Produce		json
//	@Param			req	query		domain.ListModelReq	true	"模型"
//	@Success		200	{object}	domain.Resp{data=[]domain.Model}
//	@Router			/api/v1/model/modelkit/models [get]
func (h *ModelKit) ListModel(c echo.Context) error {
	var req domain.ListModelReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Resp{
			Code:    1,
			Message: err.Error(),
		})
	}
	models, err := h.usecase.ListModel(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Resp{
			Code:    1,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, domain.Resp{
		Code:    0,
		Message: "success",
		Data:    models,
	})
}
