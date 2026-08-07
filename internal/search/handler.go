package search

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Products
// @Summary Products
// @Description Realiza uma pesquisa de produtos
// @Tags Search
// @Produce json
// @Param search query string false "Pesquisa"
// @Param limit query int false "limite de resultados"
// @Router /search [get]
// @Success 200 {object} productSearchResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func productSearchHandler(c *echo.Context) error {
	req := new(productSearchRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}
	fmt.Println("Handler", req.Limit)
	fmt.Println("Handler", req)
	resp, err := serviceSearchProducts(req.Search, req.Limit)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}
	err = c.JSON(200, resp)
	return nil
}
