package logs

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

func init() {
	if Logger != nil {
		return
	}
	buildInfo, _ := debug.ReadBuildInfo()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Str("go_version", buildInfo.GoVersion).
		Logger()
	Logger = &logger
}
