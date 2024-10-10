package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func CreateLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05", NoColor: false}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return logger
}
