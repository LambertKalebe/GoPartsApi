package login

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"g0/database"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func Route(g *echo.Group) {
	g.POST("/login", Login)
}

type RequestLogin struct {
	Username string `json:"username" example:"User"`
	Password string `json:"password" example:"123"`
}

type ResponseLogin struct {
	Message string `json:"message" example:"Login realizado com sucesso"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UserDB struct {
	ID       int
	Username string
	PassHash string
	Role     string
}

const CheckUserQuery = `
	SELECT id, username, pass_hash, role
	FROM users
	WHERE username = ?;
`

// Login
// @Summary Login
// @Description Realiza autenticação do usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param request body RequestLogin true "Credenciais"
// @Success 200 {object} ResponseLogin
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /auth/login [post]
func Login(c *echo.Context) error {
	req := new(RequestLogin)

	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	var user UserDB

	err := database.DB.QueryRow(
		CheckUserQuery,
		req.Username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PassHash,
		&user.Role,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return c.String(http.StatusUnauthorized, "Usuário ou senha inválidos")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PassHash),
		[]byte(req.Password),
	)

	if err != nil {
		return c.String(http.StatusUnauthorized, "Usuário ou senha inválidos")
	}

	token, err := GenerateToken(
		user.ID,
		user.Username,
		user.Role,
	)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, ResponseLogin{
		Message: "Login realizado com sucesso",
		Token:   token,
	})
}
