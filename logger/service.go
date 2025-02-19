package logger

import (
	"os"

	"github.com/orkarstoft/dns-updater/config"
	"github.com/rs/zerolog"
)

func New(logType config.LogType, logLevel zerolog.Level) *zerolog.Logger {
	var logger zerolog.Logger
	switch logType {
	case config.LOGTYPE_JSON:
		logger = zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()
	case config.LOGTYPE_PRETTY:
		// logger = zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger().Pretty()
	case config.LOGTYPE_FILE:
		// logger = zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()
	}
	return &logger
}
