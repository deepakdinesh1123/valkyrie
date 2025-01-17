//go:build all || darwin

package executor

import (
	"context"
	"fmt"
	"os"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/container"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/system"
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
		if _, err := os.Stat(envConfig.ODIN_SYSTEM_EXECUTOR_BASE_DIR); os.IsNotExist(err) {
			err = os.Mkdir(envConfig.ODIN_SYSTEM_EXECUTOR_BASE_DIR, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("failed to create system Executor base directory: %s", err)
			}
		}
		Executor, err = system.NewSystemExecutor(ctx, envConfig, queries, workerId, tp, mp, logger)
		if err != nil {
			return nil, fmt.Errorf("error while creating Executor")
		}
	default:
		return nil, fmt.Errorf("invalid Executor")
	}
	return Executor, nil
}
