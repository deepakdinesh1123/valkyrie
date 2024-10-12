//go:build podman && linux

package provider

import (
	"context"
	"errors"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/podman"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetProvider(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (Provider, error) {
	var provider Provider
	var err error
	switch envConfig.ODIN_WORKER_PROVIDER {
	case "podman":
		provider, err = podman.NewPodmanProvider(envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to create podman provider")
			return nil, err
		}
		logger.Info().Msg("Using podman provider")
	default:
		err := errors.New("Invalid provider only podman provider is supported")
		return nil, err
	}
	return provider, nil
}
