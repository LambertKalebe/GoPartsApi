package products

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.GET("", productsHandler, middleware.JWT)
	g.GET("/:id", productByIdHandler, middleware.JWT)
	g.GET("/:id/details", productDetailsByIdHandler, middleware.JWT)
}
