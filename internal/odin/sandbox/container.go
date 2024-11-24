package sandbox

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/pool"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/sandbox/container"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetSandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (SandboxHandler, error) {
	containerPool, err := pool.NewSandboxPool(ctx, int32(envConfig.ODIN_HOT_CONTAINER), int32(envConfig.ODIN_WORKER_CONCURRENCY), envConfig.ODIN_CONTAINER_ENGINE)
	if err != nil {
		return nil, err
	}
	switch envConfig.ODIN_WORKER_SANDBOX_HANDLER {
	case "docker":
		return container.NewDockerSandboxHandler(ctx, queries, workerId, tp, mp, envConfig, containerPool, logger)
	default:
		return nil, fmt.Errorf("invalid sandbox handler")
	}
}
