package logger

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

func GetLogger() zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Str("go_version", buildInfo.GoVersion).
		Logger()
	return logger
}
