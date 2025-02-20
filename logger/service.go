package logger

import (
	"fmt"
	"os"

	"github.com/orkarstoft/dns-updater/config"
	"github.com/rs/zerolog"
)

type LoggerError struct {
	Operation string
	Err       error
}

func (e *LoggerError) Error() string {
	return fmt.Sprintf("logger error: %s: %v", e.Operation, e.Err)
}

func (e *LoggerError) Unwrap() error {
	return e.Err
}

func New(logType config.LogType, logLevel zerolog.Level) (*zerolog.Logger, error) {
	var logger zerolog.Logger

	switch logType {
	case config.LOGTYPE_JSON:
		logger = zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Timestamp().
			Logger()

	case config.LOGTYPE_PRETTY:
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
		logger = zerolog.New(output).
			Level(logLevel).
			With().
			Timestamp().
			Logger()

	case config.LOGTYPE_FILE:
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, &LoggerError{
				Operation: "open_log_file",
				Err:       err,
			}
		}
		logger = zerolog.New(file).
			Level(logLevel).
			With().
			Timestamp().
			Logger()

	default:
		return nil, &LoggerError{
			Operation: "invalid_log_type",
			Err:       fmt.Errorf("unsupported log type: %s", logType),
		}
	}

	return &logger, nil
}
