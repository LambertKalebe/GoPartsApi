package routes

import (
	"g0/api/auth/login"
	"g0/api/auth/me"
	"g0/api/auth/register"
	"g0/api/health"
	"g0/api/products"
	productIdInfo "g0/api/products/id"
	"g0/middleware"

	"github.com/labstack/echo/v5"
)

func Route(e *echo.Echo) {
	api := e.Group("")

	// Rotas públicas
	authGroup := api.Group("auth")
	systemGroup := api.Group("system")

	health.Route(systemGroup)
	register.Route(authGroup)
	login.Route(authGroup)

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

	authPrivateGroup := private.Group("/auth")
	productsPrivateGroup := private.Group("/products")

	me.Route(authPrivateGroup)
	products.Route(productsPrivateGroup)
	productIdInfo.Route(productsPrivateGroup)
}
