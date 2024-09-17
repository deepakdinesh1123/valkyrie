package logs

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type LogOptsFunc func(*LogConfig)

type LogConfig struct {
	level  string
	export string
}

func NewLogConfig(opts ...LogOptsFunc) *LogConfig {
	logConfig := &LogConfig{
		level:  "info",
		export: "console",
	}

	for _, opt := range opts {
		opt(logConfig)
	}
	return logConfig
}

func WithLevel(level string) LogOptsFunc {
	return func(config *LogConfig) {
		config.level = level
	}
}

func WithExport(export string) LogOptsFunc {
	return func(config *LogConfig) {
		config.export = export
	}
}

func GetLogger(config *LogConfig) *zerolog.Logger {
	switch config.level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var logOut io.Writer
	switch config.export {
	case "console":
		logOut = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	case "file":
		file, err := os.OpenFile("valkyrie.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil
		}
		logOut = file
	}
	logger := zerolog.New(logOut).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	return &logger
}
