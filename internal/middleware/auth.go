package middleware

import (
	"g0/internal/config"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

var JWT = echojwt.WithConfig(echojwt.Config{
	SigningKey:  config.JWTSecret,
	TokenLookup: "header:Authorization:Bearer,cookie:token",
	ErrorHandler: func(c *echo.Context, err error) error {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return err
	},
})
