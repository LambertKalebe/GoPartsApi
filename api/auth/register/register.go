package register

import (
	"g0/database"
	"net/http"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoute(g *echo.Group) {
	g.POST("/register", RegisterUser)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const InsertUserQuery = `
	INSERT INTO users (username, pass_hash, role)
	VALUES (?, ?, ?)
`

func RegisterUser(c *echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	EncryptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao gerar a senha: "+err.Error())
	}

	_, err = database.DB.Exec(
		InsertUserQuery,
		user.Username,
		string(EncryptedPass),
		"user",
	)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created",
	})
}
