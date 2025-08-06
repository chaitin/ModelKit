package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/usecase"
)

type ModelKitHandler struct {
	modelProvider domain.IModelProvider[domain.Model]
}

func NewModelKitHandler(
	echo *echo.Echo,
) *ModelKitHandler {
	m := &ModelKitHandler{
		modelProvider: usecase.NewOpenAI(),
	}

	// g := echo.Group("/api/v1/model/modelkit")

	// g.POST("/check", m.CheckModel)
	// g.GET("/models", m.ListModel)
	// g.GET("/panda/models", m.PandaModelList)

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
// func (h *ModelKitHandler) CheckModel(c echo.Context) error {
// 	var req domain.CheckModelReq
// 	if err := c.Bind(&req); err != nil {
// 		return err
// 	}
// 	// m, err := h.modelProvider.CheckModel(c.Request().Context(), &req)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.Resp{
// 			Code:    1,
// 			Message: err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusOK, domain.Resp{
// 		Code:    0,
// 		Message: "success",
// 		Data:    m,
// 	})
// }

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
// func (h *ModelKitHandler) ListModel(c echo.Context) error {
// 	var req domain.ListModelReq
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, domain.Resp{
// 			Code:    1,
// 			Message: err.Error(),
// 		})
// 	}
// 	models, err := h.modelProvider.ListModel(string(req.SubType), string(req.OwnedBy))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.Resp{
// 			Code:    1,
// 			Message: err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusOK, domain.Resp{
// 		Code:    0,
// 		Message: "success",
// 		Data:    models,
// 	})
// }

// PandaModelList 获取模型列表
//
//	@Tags			ModelKitModel
//	@Summary		获取模型列表
//	@Description	获取模型列表
//	@ID				panda-model-list
//	@Accept			json
//	@Produce		json
//	@Param			req	query		domain.GetProviderModelListReq	true	"模型"
//	@Success		200	{object}	domain.Resp{data=domain.GetProviderModelListResp}
//	@Router			/api/v1/model/modelkit/panda/models [get]
// func (h *ModelKitHandler) PandaModelList(c echo.Context) error {
// 	var req domain.GetProviderModelListReq
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, domain.Resp{
// 			Code:    1,
// 			Message: err.Error(),
// 		})
// 	}
// 	models, err := usecase.PandaModelList(c.Request().Context(), req.Provider, req.APIKey, req.BaseURL, req.APIHeader, req.Type, req.Type)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, domain.Resp{
// 			Code:    1,
// 			Message: err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusOK, domain.Resp{
// 		Code:    0,
// 		Message: "success",
// 		Data:    models,
// 	})
// }
