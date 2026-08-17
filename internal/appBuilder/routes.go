package appbuilder

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.POST("", appBuilderHandler, middleware.JWT)
}
