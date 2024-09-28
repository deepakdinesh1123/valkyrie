//go:build docker

package provider

import (
	"context"
	"errors"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/docker"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetProvider(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (Provider, error) {
	var provider Provider
	var err error
	switch envConfig.ODIN_WORKER_PROVIDER {
	case "docker":
		provider, err = docker.NewDockerProvider(envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to create docker provider")
			return nil, err
		}
	default:
		err := errors.New("Invalid provider only docker provider is supported")
		return nil, err
	}
	return provider, nil
}
