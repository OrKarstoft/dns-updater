package logger

import (
	"os"
	"time"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(cfg config.Log) {
	zerolog.TimeFieldFormat = time.RFC3339
	loglevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		loglevel = zerolog.InfoLevel
		log.Warn().Str("level", cfg.Level).Msg("invalid log level, defaulting to INFO")
	}
	zerolog.SetGlobalLevel(loglevel)

	var logger zerolog.Logger

	switch cfg.Type {
	case config.LOGTYPE_PRETTY:
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
		logger.Info().Msg("using pretty console logger")
	case config.LOGTYPE_JSON:
		fallthrough
	default:
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	log.Logger = logger
}
