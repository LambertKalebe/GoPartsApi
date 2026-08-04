package routes

import (
	"g0/api/auth/login"
	"g0/api/auth/register"
	"g0/api/health"

	"github.com/labstack/echo/v5"
)

func RegisterRoute(e *echo.Echo) {
	api := e.Group("/api")

	// Rotas públicas
	public := api.Group("")

	health.RegisterRoute(public)
	register.RegisterRoute(public)
	login.RegisterRoute(public)

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
	// private := api.Group("")
	// private.Use(middleware.JWT)
}
