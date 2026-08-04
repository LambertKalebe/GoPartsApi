package middleware

import echojwt "github.com/labstack/echo-jwt/v5"

var JWTSecret = []byte("secret")

var JWT = echojwt.WithConfig(echojwt.Config{
	SigningKey:  JWTSecret,
	TokenLookup: "header:Authorization:Bearer ,cookie:Auth",
})
