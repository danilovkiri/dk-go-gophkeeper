// Package logger provides logging functionality.

package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// InitLog initializes a logger.
func InitLog(file *os.File) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	Logger := zerolog.New(multi).With().Timestamp().Logger()
	return &Logger
}
