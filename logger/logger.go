package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var ll *zerolog.Logger = nil

func SetupLog() *zerolog.Logger {

	if ll != nil {
		return ll
	}

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

	var ll zerolog.Logger

	if isEnvExist("TARDILOGFILE") { // env variable exists
		ll = SetupLogFile()
	} else { // env variable doesn't exist
		ll = zerolog.New(output).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return &ll
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

func SetupLogFile() zerolog.Logger {

	userDirName, _ := os.UserHomeDir()
	tardiLogFileName := userDirName + "/.tardigrade/tardilog.log"

	file, err := os.OpenFile(
		tardiLogFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	// defer file.Close()

	ll := zerolog.New(file).With().Timestamp().Logger()

	return ll
}

func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}
