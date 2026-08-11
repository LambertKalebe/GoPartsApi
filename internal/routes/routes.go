package routes

import (
	appbuilder "g0/internal/AppBuilder"
	"g0/internal/auth"
	"g0/internal/products"
	"g0/internal/search"
	"g0/internal/system"
	"g0/internal/vehicles"

	"github.com/labstack/echo/v5"
)

func Route(e *echo.Echo) {
	system.Routes(e.Group("system"))
	auth.Routes(e.Group("auth"))
	products.Routes(e.Group("products"))
	vehicles.Routes(e.Group("vehicles"))
	search.Routes(e.Group("search"))
	appbuilder.Routes(e.Group("appbuilder"))

	// OpenAPI
	e.Static("/openapi", "./docs")

	// Scalar
	e.GET("/docs", func(c *echo.Context) error {
		return c.HTML(200, `
			<!doctype html>
			<html>
			<head>
				<title>G0 API Reference</title>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<div id="app"></div>

				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

				<script>
					Scalar.createApiReference('#app', {
						url: '/openapi/swagger.json'
					})
				</script>
			</body>
			</html>
			`)
	})

}
