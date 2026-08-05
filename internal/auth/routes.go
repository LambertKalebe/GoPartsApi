package auth

import (
	"g0/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Routes(g *echo.Group) {
	g.POST("/login", login)
	g.POST("/register", register)
	g.POST("/logout", logout)
	g.GET("/aboutme", aboutMe, middleware.JWT)
}
