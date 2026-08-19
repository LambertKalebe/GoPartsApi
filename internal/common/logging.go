package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type MyLogger struct {
	zerolog.Logger
}

var Logger MyLogger

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Grey   = "\033[90m"
)

func NewLogger() MyLogger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC822,
		NoColor:    false,
	}

	output.FormatLevel = func(i any) string {
		level := i.(string)

		switch level {
		case "panic":
			return "\033[31mPANIC |\033[0m"

		case "fatal":
			return "\033[31mFATAL |\033[0m"

		case "error":
			return "\033[31mERROR |\033[0m"

		case "warn":
			return "\033[33mWARN  |\033[0m"

		case "info":
			return "\033[32mINFO  |\033[0m"

		case "debug":
			return "\033[90mDEBUG |\033[0m"

		case "trace":
			return "\033[90mTRACE |\033[0m"

		default:
			return fmt.Sprintf("| %-6s|", strings.ToUpper(level))
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

	logger := zerolog.New(output).With().Timestamp().Logger()

	Logger = MyLogger{logger}
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
