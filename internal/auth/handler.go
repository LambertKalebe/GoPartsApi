package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// Login
// @Summary Login
// @Description Realiza um login de usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Credenciais"
// @Success 200 {object} loginResponse
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 401 {string} string "Usuário ou senha inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /auth/login [post]
func login(c *echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceLogin(*req)
	if errors.Is(err, errInvalidCredentials) {
		return c.String(http.StatusUnauthorized, "Usuário ou senha inválidos")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(loginCookie(res.Token))
	return c.JSON(http.StatusOK, res)
}

func loginCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
}

func logoutCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	}
}

// Register
// @Summary Register
// @Description Cria um novo usuário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Dados do usuário"
// @Success 201 {object} registerResponse
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /auth/register [post]
func register(c *echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceRegister(*req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)

}

// Logout
// @Summary Logout
// @Description Apaga os cookies do usuário.
// @Tags Auth
// @Produce json
// @Success 201 {object} logoutResponse
// @Router /auth/logout [post]
func logout(c *echo.Context) error {
	res := serviceLogout()
	c.SetCookie(logoutCookie())
	return c.JSON(http.StatusOK, res)
}

// aboutMe
// @Summary aboutMe
// @Description Retorna informações do usuário.
// @Tags Auth
// @Produce json
// @Success 201 {object} aboutMeResponse
// @Router /auth/aboutme [get]
func aboutMe(c *echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	res, err := serviceAboutMe(token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
