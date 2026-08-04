package health

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Route(g *echo.Group) {
	g.GET("/health", Health)
}

// Health
//
// @Summary Health Check
// @Description Verifica se a API está online
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /system/health [get]
func Health(c *echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"server_online": true,
		"message":       "Server is running",
	})
}
