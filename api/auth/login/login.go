package login

import (
	"database/sql"
	"net/http"
	"time"

	"g0/database"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoute(g *echo.Group) {
	g.POST("/login", Login)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func Login(c *echo.Context) error {
	req := new(LoginRequest)

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

	if err == sql.ErrNoRows {
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
		Name:     "Auth",
		Value:    token,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login realizado com sucesso: " + token,
	})
}
