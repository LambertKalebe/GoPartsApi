package appbuilder

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Appbuilder
// @Summary Search
// @Description Realiza uma consulta de carros por um request de uma ou mais linhas
// @Tags AppBuilder
// @Produce json
// @Param search body appBuilderSearchRequest true "Linhas de pesquisa"
// @Accept json
// @Router /appbuilder [post]
// @Success 200 {object} appBuilderResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func appBuilderHandler(c *echo.Context) error {
	req := new(appBuilderSearchRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceAppBuilder(req.Search)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"Usuário ou senha inválidos",
		)
	}
	return c.JSON(http.StatusOK, res)
}
