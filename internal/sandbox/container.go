package sandbox

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/sandbox/container"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetSandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (SandboxHandler, error) {
	switch envConfig.RUNTIME {
	case "docker":
		return container.NewDockerSandboxHandler(ctx, queries, workerId, tp, mp, envConfig, logger)
	default:
		return nil, fmt.Errorf("invalid sandbox handler")
	}
}
