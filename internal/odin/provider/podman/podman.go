//go:build linux

package podman

import (
	"context"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type PodmanProvider struct {
	queries db.Store
	conn      context.Context
	envConfig *config.EnvConfig
	workerId int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
}

func NewPodmanProvider(env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*PodmanProvider, error) {
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
    socket := "unix:" + sock_dir + "/podman/podman.sock"
	logger.Info().Str("Socket", socket).Msg("Connecting to podman")
	conn, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to podman")
		return nil, err
	}
	logger.Info().Msg("Successfully connected to podman")
	return &PodmanProvider{
		queries:   queries,
		envConfig: env,
		logger:    logger,
		conn:      conn,
		workerId:  workerId,
		tp:        tp,
		mp:        mp,
	}, nil
}
