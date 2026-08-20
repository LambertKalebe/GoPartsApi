// Retornar retornar alguns dados como o diametro do tambor de freio (para o filtro do appbuilder funcionar)
// Adicionar configuração no FTS5 (sedan, hatch, etc)

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
// @Router /go-api/appbuilder [post]
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
	if res.Total == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound, "Nenhum resultado encontrado",
		)
	}
	return c.JSON(http.StatusOK, res)
}
