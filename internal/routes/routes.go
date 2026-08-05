package routes

import (
	"g0/api/products"
	productIdInfo "g0/api/products/id"
	"g0/internal/auth"
	"g0/internal/health"
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Route(e *echo.Echo) {
	api := e.Group("")

	// Rotas públicas
	authGroup := api.Group("auth")
	systemGroup := api.Group("system")

	health.Routes(systemGroup)
	auth.Routes(authGroup)

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

	// Rotas privadas
	private := api.Group("")
	private.Use(middleware.JWT)

	productsPrivateGroup := private.Group("/products")

	products.Route(productsPrivateGroup)
	productIdInfo.Route(productsPrivateGroup)
}
