package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4/middleware"

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
	fmt.Println("list model req:", req)

	resp, err := usecase.ModelList(c.Request().Context(), &req)
	if err != nil {
		fmt.Println("err:", err)
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
	var req domain.CheckModelReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
		})
	}
	fmt.Println("check model req:", req)

	resp, _ := usecase.CheckModel(c.Request().Context(), &req)

	// 如果检查过程中有错误，返回错误响应
	if resp.Error != "" {
		fmt.Println("resp.Error:", resp.Error)
		return c.JSON(http.StatusOK, domain.Response{
			Success: false,
			Message: resp.Error,
			Data:    resp,
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Message: "模型检查成功",
		Data:    resp,
	})
}

// @title			ModelKit API
// @version		1.0
// @description	ModelKit API server for model management
// @host			localhost:8080
// @BasePath		/
func main() {
	echo := echo.New()

	// 添加CORS中间件
	echo.Use(middleware.CORS())

	NewModelKit(echo, nil, false)

	err := echo.Start(":8080")
	if err != nil {
		fmt.Println("err:", err)
	}
}
