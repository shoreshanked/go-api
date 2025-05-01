package log_controller

import (
	"os"

	"github.com/rs/zerolog"
)

type LogService struct {
	logLevel zerolog.Level
	logger   zerolog.Logger
}

func (l *LogService) Logger() zerolog.Logger {
	return l.logger
}

func NewLogService(levelStr string) *LogService {
	level := parseLogLevel(levelStr)

	// Set global log level immediately
	zerolog.SetGlobalLevel(level)

	// Set up logger instance
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}).With().Timestamp().Logger()

	return &LogService{
		logLevel: level,
		logger:   logger,
	}
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		// Default fallback
		return zerolog.InfoLevel
	}
}
