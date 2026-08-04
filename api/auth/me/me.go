package me

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func Route(g *echo.Group) {
	g.GET("/me", Validate)
}

type ResponseMe struct {
	Valid bool         `json:"valid" example:"true"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
	Role     string `json:"role" example:"admin"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"invalid or expired jwt"`
}

// Validate
//
//	@Summary		Me
//	@Description	Valida o JWT enviado no cookie e retorna os dados do usuário.
//	@Tags			Autenticação
//	@Security		CookieAuth
//	@securityDefinitions.apikey	CookieAuth
//	@in							cookie
//	@name						token
//	@Produce		json
//	@Success		200	{object}	ResponseMe
//	@Failure		401	{object}	ErrorResponse
//	@Router			/auth/me [get]
func Validate(c *echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)
	id := int64(claims["user_id"].(float64))

	return c.JSON(http.StatusOK, ResponseMe{
		Valid: true,
		User: UserResponse{
			ID:       id,
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
		},
	})
}
