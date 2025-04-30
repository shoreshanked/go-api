package log_controller

import (
	"os"

	"github.com/rs/zerolog"
)

func SetupLogger() zerolog.Logger {
	// Use JSON formatting
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,             // Log to standard error
		TimeFormat: "2006-01-02 15:04:05", // Set the time format for the logs
		NoColor:    false,                 // Set to true if you don't want colors
	}).With().Timestamp().Logger()

	// Set the global logger (optional)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return logger
}
