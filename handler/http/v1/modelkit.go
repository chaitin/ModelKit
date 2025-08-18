package v1

import (
	"net/http"
	"fmt"

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
	g.GET("/listmodel", m.GetModelList)
	g.GET("/checkmodel", m.CheckModel)

	return m
}

func (p *ModelKit) GetModelList(c echo.Context) error {
	var req domain.ModelListReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
		})
	}
	fmt.Println("req:", req)

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

func (p *ModelKit) CheckModel(c echo.Context) error {
// 打印原始请求信息
fmt.Printf("Request URL: %s\n", c.Request().URL)
fmt.Printf("Request Body: %v\n", c.Request().Body)

	var req domain.CheckModelReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
		})
	}
	fmt.Println("req:", req)

	resp, err := usecase.CheckModel(c.Request().Context(), &req)
	if err != nil {
		fmt.Println("err:", err)
		return c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	// 如果检查过程中有错误，返回错误响应
	if resp.Error != "" {
		fmt.Println("resp.Error:", resp.Error)
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
