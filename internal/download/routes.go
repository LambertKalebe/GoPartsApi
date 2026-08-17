package download

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.GET("/images", imagesExportHandler, middleware.JWT)
	g.POST("/apps", appExportHandler, middleware.JWT)
}
