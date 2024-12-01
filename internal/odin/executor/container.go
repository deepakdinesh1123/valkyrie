//go:build docker || podman

package executor

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/container"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func GetExecutor(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (Executor, error) {
	var Executor Executor
	var err error
	switch envConfig.ODIN_WORKER_EXECUTOR {
	case "container":
		Executor, err = container.NewContainerExecutor(ctx, envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			return nil, err
		}
	case "system":
		return nil, fmt.Errorf("System  executor is not supported")
	default:
		return nil, fmt.Errorf("invalid Executor")
	}
	return Executor, nil
}
