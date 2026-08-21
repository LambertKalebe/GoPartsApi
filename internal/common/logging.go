package common

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type MyLogger struct {
	zerolog.Logger
}

var Logs = NewLogHub()
var Logger MyLogger

const maxLogHistory = 1000

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Grey   = "\033[90m"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func newConsoleWriter(out io.Writer, noColor bool) zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC822,
		NoColor:    noColor,
	}

	output.FormatLevel = func(i any) string {
		level := fmt.Sprintf("%s", i)

		switch level {
		case "panic":
			if noColor {
				return "PANIC |"
			}
			return Red + "PANIC |" + Reset

		case "fatal":
			if noColor {
				return "FATAL |"
			}
			return Red + "FATAL |" + Reset

		case "error":
			if noColor {
				return "ERROR |"
			}
			return Red + "ERROR |" + Reset

		case "warn":
			if noColor {
				return "WARN  |"
			}
			return Yellow + "WARN  |" + Reset

		case "info":
			if noColor {
				return "INFO  |"
			}
			return Green + "INFO  |" + Reset

		case "debug":
			if noColor {
				return "DEBUG |"
			}
			return Grey + "DEBUG |" + Reset

		case "trace":
			if noColor {
				return "TRACE |"
			}
			return Grey + "TRACE |" + Reset

		default:
			return fmt.Sprintf(
				"| %-6s|",
				strings.ToUpper(level),
			)
		}
	}

	output.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s:", i)
	}

	output.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	output.FormatErrFieldName = func(i any) string {
		return fmt.Sprintf("%s: ", i)
	}

	return output
}

func NewLogger() MyLogger {
	consoleOutput := newConsoleWriter(
		os.Stderr,
		false,
	)

	logger := zerolog.New(consoleOutput).
		With().
		Timestamp().
		Logger()

	Logger = MyLogger{
		Logger: logger,
	}

	return Logger
}

func (l *MyLogger) LogInfo() *zerolog.Event {
	return l.Logger.Info()
}

func (l *MyLogger) LogError() *zerolog.Event {
	return l.Logger.Error()
}

func (l *MyLogger) LogDebug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *MyLogger) LogWarn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *MyLogger) LogFatal() *zerolog.Event {
	return l.Logger.Fatal()
}

func NewLogHub() *logHub {
	return &logHub{
		clients: make(map[chan LogResponse]struct{}),
		history: make([]LogResponse, 0, maxLogHistory),
	}
}

func (h *logHub) Publish(entry LogResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.history = append(h.history, entry)

	if len(h.history) > maxLogHistory {
		h.history = h.history[len(h.history)-maxLogHistory:]
	}

	for client := range h.clients {
		select {
		case client <- entry:
		default:
			// Cliente lento não bloqueia o sistema.
		}
	}
}

func (h *logHub) Subscribe() chan LogResponse {
	client := make(chan LogResponse, 100)

	h.mu.Lock()
	h.clients[client] = struct{}{}
	h.mu.Unlock()

	return client
}

func (h *logHub) Unsubscribe(client chan LogResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client)
	}
}

func (h *logHub) History() []LogResponse {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]LogResponse, len(h.history))
	copy(result, h.history)

	return result
}
