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
// @Param qnt query int false "Quantidade máxima de produtos por página"
// @Param publicOnly query bool false "Somente produtos públicos"
// @Param page query int false "Página desejada"
// @Router /products [get]
// @Success 200 {object} productsResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func productsHandler(c *echo.Context) error {
	req := new(productRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceProducts(req.Qnt, req.Page, req.PublicOnly)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}
	err = c.JSON(200, res)
	return nil
}

// ProductId
// @Summary ID
// @Description Realiza uma consulta basica de um produto com base no seu ID
// @Tags Products
// @Produce json
// @Param int query int false "ID do produto"
// @Router /products/{id} [get]
// @Success 200 {object} productByIdResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func productByIdHandler(c *echo.Context) error {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceProductById(id)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// ProductDetails
// @Summary Details
// @Description Realiza uma consulta completa de um produto com base no seu ID
// @Tags Products
// @Produce json
// @Param int query int false "ID do produto"
// @Router /products/{id}/details [get]
// @Success 200 {object} productDetailsByIdResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func productDetailsByIdHandler(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceProductDetailsById(id)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, res)
}
