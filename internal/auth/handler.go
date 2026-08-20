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
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
// @Router /api/auth/login [post]
func login(c *echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceLogin(*req)
	if errors.Is(err, errInvalidCredentials) {
		return echo.NewHTTPError(
			http.StatusUnauthorized, "Usuário ou senha inválidos")
	}
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
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
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
// @Router /api/auth/register [post]
func register(c *echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceRegister(*req)
	if errors.Is(err, errUserAlreadyExists) {
		return echo.NewHTTPError(
			http.StatusConflict, "Usuário já existe")
	}
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, res)

}

// Logout
// @Summary Logout
// @Description Apaga os cookies do usuário.
// @Tags Auth
// @Produce json
// @Success 200 {object} logoutResponse
// @Router /api/auth/logout [post]
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
// @Success 200 {object} aboutMeResponse
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Router /api/auth/aboutme [get]
func aboutMe(c *echo.Context) error {
	if c.Get("user") == nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized, "Usuário não logado")
	}
	token := c.Get("user").(*jwt.Token)
	if token == nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized, "Usuário não autenticado")
	}
	res, err := serviceAboutMe(token)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
