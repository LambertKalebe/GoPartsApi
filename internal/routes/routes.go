package routes

import (
	"g0/internal/appBuilder"
	"g0/internal/auth"
	"g0/internal/download"
	"g0/internal/products"
	"g0/internal/search"
	"g0/internal/system"
	"g0/internal/vehicles"

	"github.com/labstack/echo/v5"
)

func Route(e *echo.Echo) {
	// API
	system.Routes(e.Group("api/system"))
	auth.Routes(e.Group("api/auth"))
	products.Routes(e.Group("api/products"))
	vehicles.Routes(e.Group("api/vehicles"))
	search.Routes(e.Group("api/search"))
	appbuilder.Routes(e.Group("api/appbuilder"))
	download.Routes(e.Group("api/export"))

	// OpenAPI
	e.Static("/openapi", "./docs")

	// Scalar
	e.GET("/api/docs", func(c *echo.Context) error {
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

	// Frontend Vite
	e.Static("/assets", "./dist/assets")
	e.File("/", "./dist/index.html")
	e.GET("/*", func(c *echo.Context) error {
		return c.File("./dist/index.html")
	})
}
