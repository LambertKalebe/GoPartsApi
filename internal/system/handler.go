package system

import (
	"encoding/json"
	"fmt"
	"g0/internal/common"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// Health
// @Summary Check
// @Description Verifica se a API está online
// @Tags System
// @Produce json
// @Success 200 {object} healthResponse
// @Router /api/system/health [get]
func health(c *echo.Context) error {
	res := serviceHealth()

	if !res.Healthy {
		return c.JSON(http.StatusServiceUnavailable, res)
	}

	return c.JSON(http.StatusOK, res)
}

// Stop
// @Summary Stop
// @Description Finaliza o servidor
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/system/stop [post]
func stopServer(c *echo.Context) error {
	go func() {
		time.Sleep(500 * time.Millisecond)

		common.Logger.LogWarn().
			Msg("Desligando servidor...")

		common.ExitCode.Store(0)

		if common.ServerCancel != nil {
			common.ServerCancel()
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Servidor desligando",
	})
}

// Restart
// @Summary Restart
// @Description Reinicia o servidor
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/system/restart [post]
func restartServer(c *echo.Context) error {
	go func() {
		time.Sleep(500 * time.Millisecond)

		common.Logger.LogWarn().
			Msg("Reiniciando servidor...")

		common.ExitCode.Store(1)

		if common.ServerCancel != nil {
			common.ServerCancel()
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Servidor reiniciando",
	})
}

// LogsStream
// @Summary Stream de logs
// @Description Transmite os logs do sistema em tempo real usando Server-Sent Events (SSE)
// @Tags System
// @Produce text/event-stream
// @Success 200 {string} string "Stream SSE de logs"
// @Failure 500 {object} httpcustom.ErrorResponse "Streaming não suportado"
// @Router /api/system/logstream [get]
func logsStream(c *echo.Context) error {
	res := c.Response()

	flusher, ok := res.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Streaming não suportado",
		)
	}

	res.Header().Set(
		"Content-Type",
		"text/event-stream",
	)

	res.Header().Set(
		"Cache-Control",
		"no-cache",
	)

	res.Header().Set(
		"Connection",
		"keep-alive",
	)

	client := common.Logs.Subscribe()
	defer common.Logs.Unsubscribe(client)

	send := func(entry common.LogResponse) error {
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(
			res,
			"data: %s\n\n",
			data,
		)

		if err != nil {
			return err
		}

		flusher.Flush()

		return nil
	}

	// Histórico
	for _, entry := range common.Logs.History() {
		if err := send(entry); err != nil {
			return nil
		}
	}

	// Novos logs
	for {
		select {
		case <-c.Request().Context().Done():
			return nil

		case entry, ok := <-client:
			if !ok {
				return nil
			}

			if err := send(entry); err != nil {
				return nil
			}
		}
	}
}
