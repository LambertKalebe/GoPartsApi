package httpcustom

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HTTPErrorHandler(c *echo.Context, err error) {
	code := http.StatusInternalServerError
	message := "Erro interno do servidor"

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code

		switch code {
		case http.StatusBadRequest:
			message = "Requisição inválida"

		case http.StatusUnauthorized:
			message = "Não autorizado"

		case http.StatusForbidden:
			message = "Acesso negado"

		case http.StatusNotFound:
			message = "Recurso não encontrado"

		case http.StatusMethodNotAllowed:
			message = "Método HTTP não permitido"

		case http.StatusConflict:
			message = "Conflito"

		case http.StatusUnprocessableEntity:
			message = "Dados inválidos"

		case http.StatusTooManyRequests:
			message = "Muitas requisições"

		default:
			message = he.Message
		}
	}

	_ = c.JSON(code, ErrorResponse{
		Status:  code,
		Message: message,
	})
}
