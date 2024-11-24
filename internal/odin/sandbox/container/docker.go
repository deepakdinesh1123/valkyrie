package container

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/pool"
	"github.com/docker/docker/client"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DockerSH struct {
	client        *client.Client
	queries       db.Store
	envConfig     *config.EnvConfig
	workerId      int32
	logger        *zerolog.Logger
	tp            trace.TracerProvider
	mp            metric.MeterProvider
	containerPool *puddle.Pool[pool.Container]
}

func NewDockerSandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, containerPool *puddle.Pool[pool.Container], logger *zerolog.Logger) (*DockerSH, error) {
	client := pool.GetDockerClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &DockerSH{
		client:        client,
		queries:       queries,
		envConfig:     envConfig,
		workerId:      workerId,
		logger:        logger,
		tp:            tp,
		mp:            mp,
		containerPool: containerPool,
	}, nil
}

func (d *DockerSH) Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBox db.Sandbox) {
	cont := d.containerPool.Acquire(ctx)
}
