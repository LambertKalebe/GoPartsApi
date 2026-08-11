package appbuilder

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Appbuilder
// @Summary Search
// @Description Realiza uma consulta de produtos por por um request de uma ou mais linhas
// @Tags AppBuilder
// @Produce json
// @Param search body appBuilderSearchRequest true "Linhas de pesquisa"
// @Accept json
// @Router /appbuilder [post]
// @Success 200 {object} appBuilderSearchResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func appBuilderHandler(c *echo.Context) error {
	req := new(appBuilderSearchRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	fmt.Printf("Handler Search: %#v\n", req.Search)
	fmt.Println("-----------------------------------")

	res, err := serviceAppBuilder(req.Search)
	if err != nil {
		fmt.Println("Handler", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
