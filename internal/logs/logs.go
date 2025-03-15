package logs

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogOptsFunc func(*LogConfig)

type LogConfig struct {
	level  string
	export string
	source string
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

func WithSource(source string) LogOptsFunc {
	return func(config *LogConfig) {
		config.source = source
	}
}

func GetLogger(config *LogConfig) *zerolog.Logger {
	var logLevel zerolog.Level
	switch config.level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var logOut io.Writer
	switch config.export {
	case "console":
		logOut = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	case "file":
		logOut = &lumberjack.Logger{
			Filename:   fmt.Sprintf("valkyrie_%s.log", config.source),
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
	}
	logger := zerolog.New(logOut).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	return &logger
}
