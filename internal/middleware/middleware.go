package middleware

import (
	"errors"
	"fmt"
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

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()

		original := c.Response()

		writer := &responseWriter{
			ResponseWriter: original,
		}

		c.SetResponse(writer)

		err := next(c)

		status := writer.status

		statusColor := common.Green

		if status == 0 {
			status = http.StatusOK
		}

		duration := time.Since(start).Microseconds()

		var errorMessage string

		if err != nil {
			var he *echo.HTTPError

			if errors.As(err, &he) {
				status = he.Code
				errorMessage = fmt.Sprint(he.Message)
			} else {
				status = http.StatusInternalServerError
				errorMessage = err.Error()
			}
		}

		if status >= 400 {
			statusColor = common.Red
		}

		message := fmt.Sprintf(
			"%s%d%s - %s%s%s - %s%dμs%s",
			statusColor,
			status,
			common.Reset,

			common.Blue,
			c.Request().URL.Path,
			common.Reset,

			common.Yellow,
			duration,
			common.Reset,
		)

		if err != nil {
			message += fmt.Sprintf(" - " + errorMessage)
			common.Logger.LogError().Msg(message)
		} else {
			common.Logger.LogInfo().Msg(message)
		}

		return err
	}
}
