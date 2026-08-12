package appbuilder

import "github.com/labstack/echo/v5"

func Routes(g *echo.Group) {
	g.POST("", appBuilderHandler)
}
