package v1

import (
	"net/http"

	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/pkg/log"
	"github.com/chaitin/ModelKit/usecase"
	"github.com/labstack/echo/v4"
)

type ModelKit struct{}

func NewModelKit(
	echo *echo.Echo,
	logger *log.Logger,
	isApmEnabled bool,
) *ModelKit {
	m := &ModelKit{}

	// 注册路由
	g := echo.Group("/api/v1/modelkit")
	g.GET("/models", m.GetModelList)
	g.POST("/check", m.CheckModel)

	return m
}

// GetModelList 获取模型列表
//
//	@Tags			Panda
//	@Summary		获取模型列表
//	@Description	获取指定提供商的模型列表
//	@ID				panda-get-model-list
//	@Accept			json
//	@Produce		json
//	@Param			provider	query		string	true	"模型提供商"
//	@Param			base_url	query		string	true	"Base URL"
//	@Param			api_key		query		string	false	"API Key"
//	@Param			api_header	query		string	false	"API Header"
//	@Param			type		query		string	true	"模型类型"
//	@Success		200			{object}	domain.Response{data=domain.GetProviderModelListResp}
//	@Router			/api/v1/panda/models [get]
func (p *ModelKit) GetModelList(c echo.Context) error {
	var req domain.ModelListReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
		})
	}

	// 验证参数
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数验证失败: " + err.Error(),
		})
	}

	resp, err := usecase.ModelList(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Message: "获取模型列表成功",
		Data:    resp,
	})
}

// CheckModel 检查模型
//
//	@Tags			Panda
//	@Summary		检查模型
//	@Description	检查模型配置是否正确
//	@ID				panda-check-model
//	@Accept			json
//	@Produce		json
//	@Param			model	body		domain.CheckModelReq	true	"模型配置"
//	@Success		200		{object}	domain.Response{data=domain.CheckModelResp}
//	@Router			/api/v1/panda/check [post]
func (p *ModelKit) CheckModel(c echo.Context) error {
	var req domain.CheckModelReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
		})
	}

	resp, err := usecase.CheckModel(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// 如果检查过程中有错误，返回错误响应
	if resp.Error != "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "模型检查失败",
			Data:    resp,
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Message: "模型检查成功",
		Data:    resp,
	})
}
