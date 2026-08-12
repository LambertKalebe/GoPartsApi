package export

import "github.com/labstack/echo/v5"

func Routes(g *echo.Group) {
	g.GET("/images", imagesExportHandler)
}
