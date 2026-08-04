package register

import (
	"g0/database"
	"net/http"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func Route(g *echo.Group) {
	g.POST("/register", UserRegister)
}

type User struct {
	Username string `json:"username" example:"User"`
	Password string `json:"password" example:"123"`
}

type ResponseRegister struct {
	Message string `json:"message" example:"User created"`
}

const InsertUserQuery = `
	INSERT INTO users (username, pass_hash, role)
	VALUES (?, ?, ?)
`

// UserRegister
// @Summary Registrar usuário
// @Description Cria um novo usuário com a senha criptografada utilizando bcrypt.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param request body User true "Dados do usuário"
// @Success 201 {object} ResponseRegister
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /auth/register [post]
func UserRegister(c *echo.Context) error {
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

	return c.JSON(http.StatusCreated, ResponseRegister{
		Message: "User created",
	})
}
