package logs

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func GetLogger() *zerolog.Logger {
	logFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err) // Handle error properly in production
	}
	defer logFile.Close()

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(logFile).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	return &logger
}
