package search

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Products
// @Summary Products
// @Description Realiza uma pesquisa de produtos
// @Tags Search
// @Produce json
// @Param search query string true "Pesquisa"
// @Param limit query int false "limite de resultados (padrão: 100)"
// @Router /api/search/products [get]
// @Success 200 {object} productSearchResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func productSearchHandler(c *echo.Context) error {
	req := new(productSearchRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}
	resp, err := serviceSearchProducts(req.Search, req.Limit)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	if len(resp.Products) == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound, "Nenhum carro encontrado")
	}
	err = c.JSON(200, resp)
	return nil
}

// Cars
// @Summary Cars
// @Description Realiza uma pesquisa de carros
// @Tags Search
// @Produce json
// @Param search query string true "Pesquisa"
// @Param limit query int false "limite de resultados (padrão: 100)"
// @Router /api/search/cars [get]
// @Success 200 {object} carSearchResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func carSearchHandler(c *echo.Context) error {
	req := new(carSearchRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}
	resp, err := serviceSearchCars(req.Search, req.Limit)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	if len(resp.Cars) == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound, "Nenhum carro encontrado")
	}
	err = c.JSON(200, resp)
	return nil
}
