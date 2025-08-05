package main

import (
	v1 "github.com/chaitin/ModelKit/handler/http/v1"
	"github.com/labstack/echo/v4"
)

// @title			ModelKit API
// @version		1.0
// @description	ModelKit API server for model management
// @host			localhost:8080
// @BasePath		/
func main() {
	echo := echo.New()

	v1.NewModelKitHandler(echo)

	echo.Start(":8080")
}
