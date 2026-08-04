package middleware

import echojwt "github.com/labstack/echo-jwt/v5"

var JWTSecret = []byte("IHave10Dogs")

var JWT = echojwt.WithConfig(echojwt.Config{
	SigningKey:  JWTSecret,
	TokenLookup: "header:Authorization:Bearer ,cookie:token",
})
