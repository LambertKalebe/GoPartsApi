package system

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.GET("/health", health)
	g.GET("/logstream", logsStream)
	g.POST("/restart", restartServer, middleware.JWT)
	g.POST("/stop", stopServer, middleware.JWT)
}
