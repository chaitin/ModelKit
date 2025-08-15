package main

import (
	v1 "github.com/chaitin/ModelKit/handler/http/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title			ModelKit API
// @version		1.0
// @description	ModelKit API server for model management
// @host			localhost:8080
// @BasePath		/
func main() {
	echo := echo.New()

	// 添加CORS中间件
	echo.Use(middleware.CORS())

	v1.NewModelKit(echo, nil, false)

	echo.Start(":8080")
}
