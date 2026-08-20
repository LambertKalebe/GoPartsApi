package products

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// ProductsList
// @Summary List
// @Description Realiza uma consulta de produtos por paginação
// @Tags Products
// @Produce json
// @Param limit query int false "Quantidade máxima de produtos por página" default(100) minimum(1) maximum(500)
// @Param publicOnly query bool false "Somente produtos públicos"
// @Param page query int false "Página desejada"
// @Router /api/products [get]
// @Success 200 {object} productsResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func productsHandler(c *echo.Context) error {
	req := new(productRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceProducts(req.Limit, req.Page, req.PublicOnly)
	if err != nil {
		fmt.Println("Handler", err)
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	err = c.JSON(200, res)
	return nil
}

// ProductId
// @Summary ID
// @Description Realiza uma consulta basica de um produto com base no seu ID
// @Tags Products
// @Produce json
// @Router /api/products/{id} [get]
// @Param id path int true "ID do produto"
// @Success 200 {object} productByIdResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func productByIdHandler(c *echo.Context) error {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceProductById(id)
	if err != nil {
		fmt.Println("Handler", err)
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

// ProductDetails
// @Summary Details
// @Description Realiza uma consulta completa de um produto com base no seu ID
// @Tags Products
// @Produce json
// @Router /api/products/{id}/details [get]
// @Param id path int true "ID do produto"
// @Success 200 {object} productDetailsByIdResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func productDetailsByIdHandler(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceProductDetailsById(id)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
