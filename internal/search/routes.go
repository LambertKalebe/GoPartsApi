package search

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.GET("/products", productSearchHandler, middleware.JWT)
	g.GET("/cars", carSearchHandler, middleware.JWT)
}
