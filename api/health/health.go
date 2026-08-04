package health

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func RegisterRoute(g *echo.Group) {
	g.GET("/health", Health)
}

// Health godoc
//
// @Summary Health Check
// @Description Verifica se a API está online
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func Health(c *echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"server_online": true,
		"message":       "O pai ta on",
	})
}
