// Package logger provides structured logging for the application
package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Init initializes the logger with the specified level and output
func Init(level string, pretty bool) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	var output io.Writer = os.Stdout
	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	zerolog.SetGlobalLevel(logLevel)
	log = zerolog.New(output).With().Timestamp().Logger()
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Error logs an error message
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	log.Info().Msg(fmt.Sprintf(format, v...))
}

// Errorf logs a formatted error message
func Errorf(err error, format string, v ...interface{}) {
	log.Error().Err(err).Msg(fmt.Sprintf(format, v...))
}
