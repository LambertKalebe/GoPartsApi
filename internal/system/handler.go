package system

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Health
// @Summary Check
// @Description Verifica se a API está online
// @Tags Health
// @Produce json
// @Success 200 {object} healthResponse
// @Router /go-api/system/health [get]
func health(c *echo.Context) error {
	res := serviceHealth()
	if res.Healthy == false {
		return c.JSON(http.StatusServiceUnavailable, res)
	}
	return c.JSON(http.StatusOK, res)
}
