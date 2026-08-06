package system

import "github.com/labstack/echo/v5"

func Routes(g *echo.Group) {
	g.GET("/health", health)
}
