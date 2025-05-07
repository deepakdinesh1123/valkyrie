package sandbox

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/sandbox/container"
	"github.com/deepakdinesh1123/valkyrie/internal/sandbox/k8s"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type SandboxHandler interface {
	Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBoxJob db.FetchSandboxJobTxResult)
	Cleanup(ctx context.Context) error
	StartSandboxPool(ctx context.Context, envConfig *config.EnvConfig) error
}

func GetSandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (SandboxHandler, error) {
	switch envConfig.RUNTIME {
	case "docker":
		return container.NewDockerSandboxHandler(ctx, queries, workerId, tp, mp, envConfig, logger)
	case "k8s":
		return k8s.NewK8SandboxHandler(ctx, queries, workerId, tp, mp, envConfig, logger)
	default:
		return nil, fmt.Errorf("invalid sandbox handler")
	}
}
