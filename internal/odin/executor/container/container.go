package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Container struct {
	Name string
	ID   string
	PID  int

	HostPrepDir string
}

type ContainerExecutor struct {
	queries   db.Store
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	user      string
	pool      *puddle.Pool[Container]
}

func NewContainerExecutor(ctx context.Context, env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*ContainerExecutor, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	pool, err := NewContainerPool(ctx, int32(env.ODIN_HOT_CONTAINER), env.ODIN_WORKER_CONCURRENCY)
	if err != nil {
		return nil, err
	}
	return &ContainerExecutor{
		envConfig: env,
		logger:    logger,
		queries:   queries,
		workerId:  workerId,
		tp:        tp,
		mp:        mp,
		user:      user.Username,
		pool:      pool,
	}, nil
}

func KillContainer(pid int) error {
	_, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Container with given PID has already been killed")
	}

	cmd := exec.Command("kill", "-KILL", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill container: %w", err)
	}
	return nil
}
