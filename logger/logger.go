package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func SetupLog() *zerolog.Logger {
	var output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("|%-4s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf(", %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return (fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return &log
}

func SetLogLevelInfo() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func SetLogLevelDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func SwapLogLevel() {
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if zerolog.GlobalLevel() == zerolog.InfoLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
