package middleware

import (
	"errors"
	"g0/internal/common"
	"g0/internal/config"
	"net/http"
	"time"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

var JWT = echojwt.WithConfig(echojwt.Config{
	SigningKey:  config.JWTSecret,
	TokenLookup: "header:Authorization:Bearer,cookie:token",
	ErrorHandler: func(c *echo.Context, err error) error {
		return echo.NewHTTPError(
			http.StatusUnauthorized, "Não autorizado")
	},
})

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start).String()
		status := http.StatusOK

		if err != nil {
			var he *echo.HTTPError
			if errors.As(err, &he) {
				status = he.Code
			} else {
				status = http.StatusInternalServerError
			}
		}

		fields := map[string]interface{}{
			"method":   c.Request().Method,
			"uri":      c.Request().URL.Path,
			"status":   status,
			"duration": duration,
		}

		if err != nil {
			fields["error"] = err.Error()

			common.Logger.LogError().
				Fields(fields).
				Msg("Request")

			return err
		}

		common.Logger.LogInfo().
			Fields(fields).
			Msg("Request")

		return nil
	}
}
