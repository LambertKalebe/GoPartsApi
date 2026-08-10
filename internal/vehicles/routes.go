package vehicles

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.GET("", vehiclesHandler, middleware.JWT)
	g.GET("/:id", vehicleByIdHandler, middleware.JWT)
	g.GET("/:id/details", vehicleDetailsByIdHandler, middleware.JWT)
}
