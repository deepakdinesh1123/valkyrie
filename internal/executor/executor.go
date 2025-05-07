package executor

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/executor/container"
	"github.com/deepakdinesh1123/valkyrie/internal/executor/k8s"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Executor interface {
	Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq *db.Job, logger zerolog.Logger)
	Cleanup()
}

func GetExecutor(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (Executor, error) {
	var Executor Executor
	var err error
	switch envConfig.RUNTIME {
	case "docker", "podman":
		Executor, err = container.NewContainerExecutor(ctx, envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			return nil, err
		}
	case "k8s":
		Executor, err = k8s.NewK8sExecutor(ctx, envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid Executor")
	}
	return Executor, nil
}
