package container

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/pool"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/sandbox"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgtype"
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
	defer wg.Done()

	cont, err := d.containerPool.Acquire(ctx)
	if err != nil {
		d.logger.Err(err).Msg("could not acquire container")
		return
	}
	go d.containerPool.CreateResource(ctx)
	sandboxConfig, err := sandbox.GetSandboxConfig()
	if err != nil {
		d.logger.Err(err).Msg("could not get sandbox config")
		return
	}
	configYaml, err := yaml.Marshal(sandboxConfig)
	if err != nil {
		d.logger.Err(err).Msg("could not convert sandbox config to yaml")
	}

	d.logger.Debug().Str("Command", fmt.Sprintf("echo \"%s\" >> /home/valnix/.config/code-server/config.yaml", string(configYaml))).Msg("Adding config")
	configExec, err := d.client.ContainerExecCreate(
		ctx,
		cont.Value().ID,
		container.ExecOptions{
			AttachStderr: true,
			AttachStdout: true,
			Cmd: []string{
				"/bin/sh",
				"-c",
				fmt.Sprintf("echo \"%s\" >> /home/valnix/.config/code-server/config.yaml", string(configYaml)),
			},
		},
	)
	_, err = d.client.ContainerExecAttach(ctx, configExec.ID, container.ExecAttachOptions{})
	if err != nil {
		d.logger.Err(err).Msg("error adding code server config")
		return
	}

	codeServerExec, err := d.client.ContainerExecCreate(
		ctx,
		cont.Value().ID,
		container.ExecOptions{
			AttachStderr: true,
			AttachStdout: true,
			Cmd: []string{
				"/home/valnix/start.sh",
			},
		},
	)
	err = d.client.ContainerExecStart(ctx, codeServerExec.ID, container.ExecStartOptions{
		Detach: true,
	})
	if err != nil {
		d.logger.Err(err).Msg("error adding code server config")
		return
	}

	contInfo, err := d.client.ContainerInspect(ctx, cont.Value().ID)
	if err != nil {
		d.logger.Err(err).Msg("error getting container config")
	}

	containerURL := fmt.Sprintf("http://%s:9090", contInfo.NetworkSettings.IPAddress)

	err = d.queries.MarkSandboxRunning(ctx, db.MarkSandboxRunningParams{
		SandboxID:  sandBox.SandboxID,
		SandboxUrl: pgtype.Text{String: containerURL, Valid: true},
	})
	if err != nil {
		d.logger.Err(err).Msg("error marking the container as running")
	}
}

func (d *DockerSH) Cleanup(ctx context.Context) error {
	d.containerPool.Close()
	return nil
}
