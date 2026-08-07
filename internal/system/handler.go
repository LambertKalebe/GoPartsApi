package system

import (
	"github.com/labstack/echo/v5"
)

// Health
// @Summary Check
// @Description Verifica se a API está online
// @Tags Health
// @Produce json
// @Success 200 {object} healthResponse
// @Router /health [get]
func health(c *echo.Context) error {
	res := serviceHealth()
	return c.JSON(200, res)
}
