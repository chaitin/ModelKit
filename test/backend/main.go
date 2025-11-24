package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/usecase"
	"github.com/labstack/echo/v4"
)

type ModelKitHandler struct {
	logger   *slog.Logger
	modelkit *usecase.ModelKit
}

func NewModelKit(
	echo *echo.Echo,
	logger *slog.Logger,
	isApmEnabled bool,
	modelkit *usecase.ModelKit,
) *ModelKitHandler {
	m := &ModelKitHandler{
		logger:   logger,
		modelkit: modelkit,
	}

	// 注册路由
	g := echo.Group("/api/v1/modelkit")
	g.GET("/listmodel", m.GetModelList)
	g.GET("/checkmodel", m.CheckModel)

	return m
}

func (p *ModelKitHandler) GetModelList(c echo.Context) error {
	var req domain.ModelListReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusOK, domain.Response{
			Success: false,
			Message: "参数绑定失败: " + err.Error(),
			Data:    nil,
		})
	}

	resp, err := p.modelkit.ModelList(c.Request().Context(), &req)
	if err != nil {
		fmt.Println("err:", err)
		return c.JSON(http.StatusOK, domain.Response{
			Success: true,
			Message: err.Error(),
			Data:    resp,
		})
	}

	return c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Message: "获取模型列表成功",
		Data:    resp,
	})
}

func (p *ModelKitHandler) CheckModel(c echo.Context) error {
	// 手动绑定参数，避免 param 字段的 JSON 解析问题
	var req domain.CheckModelReq

	// 绑定基本字段
	req.Provider = c.QueryParam("provider")
	req.Model = c.QueryParam("model")
	req.BaseURL = c.QueryParam("base_url")
	req.APIKey = c.QueryParam("api_key")
	req.APIHeader = c.QueryParam("api_header")
	req.APIVersion = c.QueryParam("api_version")
	req.Type = c.QueryParam("type")

	p.logger.Info("CheckModel req:", req)

	// 单独处理 param 字段
	paramStr := c.QueryParam("param")
	if paramStr != "" {
		var param domain.ModelParam
		if err := json.Unmarshal([]byte(paramStr), &param); err != nil {
			return c.JSON(http.StatusBadRequest, domain.Response{
				Success: false,
				Message: "param 参数解析失败: " + err.Error(),
			})
		}
		req.Param = &param
	}

	// 验证必需参数
	if req.Provider == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "provider 参数是必需的",
		})
	}
	if req.Model == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "model_name 参数是必需的",
		})
	}
	if req.BaseURL == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "base_url 参数是必需的",
		})
	}
	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Message: "model_type 参数是必需的",
		})
	}

	resp, err := p.modelkit.CheckModel(c.Request().Context(), &req)
	if err != nil {
		fmt.Println("err:", err)
		return c.JSON(http.StatusOK, domain.Response{
			Success: true,
			Message: err.Error(),
			Data:    resp,
		})
	}

	// 如果检查过程中有错误，返回错误响应
	if resp.Error != "" {
		fmt.Println("resp.Error:", resp.Error)
		return c.JSON(http.StatusOK, domain.Response{
			Success: true,
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

	// 创建logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	// 添加CORS中间件
	echo.Use(middleware.CORS())

	// 创建ModelKit
	modelkit := usecase.NewModelKit(
		logger,
	)

	NewModelKit(echo, logger, false, modelkit)

	err := echo.Start(":8080")
	if err != nil {
		fmt.Println("err:", err)
	}
}
