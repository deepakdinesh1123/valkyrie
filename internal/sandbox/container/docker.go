//go:build all || docker

package container

import (
	"context"
	"fmt"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/db/jsonschema"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

	containers []*puddle.Resource[pool.Container]
	mu         sync.Mutex
}

func NewDockerSandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (*DockerSH, error) {
	client := pool.GetDockerClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &DockerSH{
		client:    client,
		queries:   queries,
		envConfig: envConfig,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}

func (d *DockerSH) StartSandboxPool(ctx context.Context, envConfig *config.EnvConfig) error {
	containerPool, err := pool.NewSandboxPool(ctx, int32(envConfig.HOT_CONTAINER), int32(envConfig.WORKER_CONCURRENCY), envConfig.RUNTIME)
	if err != nil {
		return fmt.Errorf("error creating sandbox pool: %s", err)
	}
	d.containerPool = containerPool
	return nil
}

func (d *DockerSH) Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBoxJob db.FetchSandboxJobTxResult) {
	d.logger.Info().Msg("Creating sandbox")
	defer wg.Done()

	sandBox := sandBoxJob.Sandbox

	// Define error handler
	handleError := func(err error, message string) {
		if err != nil {
			d.logger.Err(err).Msg(message)
			// Mark sandbox as failed in database
			updateErr := d.queries.UpdateSandboxState(ctx, db.UpdateSandboxStateParams{
				SandboxID:    sandBox.SandboxID,
				CurrentState: "failed",
				Details:      jsonschema.SandboxDetails{Error: err.Error()},
			})
			if updateErr != nil {
				d.logger.Err(updateErr).Msg("failed to mark sandbox as failed")
			}
			return
		}
	}

	var cont *puddle.Resource[pool.Container]
	cont, err := d.containerPool.Acquire(ctx)
	if err != nil {
		handleError(err, "could not acquire container")
		return
	}
	go d.containerPool.CreateResource(ctx)

	var contInfo types.ContainerJSON
	contInfo, err = d.client.ContainerInspect(ctx, cont.Value().ID)
	if err != nil {
		if client.IsErrNotFound(err) {
			// Container not found, acquire a new container from the pool
			cont, err = d.containerPool.Acquire(ctx)
			if err != nil {
				handleError(err, "could not acquire container after not found")
				return
			}
			contInfo, err = d.client.ContainerInspect(ctx, cont.Value().ID)
			if err != nil {
				handleError(err, "error getting container config after re-acquiring")
				return
			}
		} else {
			handleError(err, "error getting container config")
			return
		}
	}

	if sandBox.Config.Flake != "" {
		flakeExec, err := d.client.ContainerExecCreate(
			ctx,
			cont.Value().ID,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd: []string{
					"/bin/sh",
					"-c",
					fmt.Sprintf("rm /home/valnix/work/flake.nix && echo \"%s\" >> /home/valnix/work/flake.nix && nix profile remove work", sandBox.Config.Flake),
				},
			},
		)

		// start execution
		if err := d.client.ContainerExecStart(ctx, flakeExec.ID, container.ExecStartOptions{}); err != nil {
			d.logger.Err(err).Msgf("error adding flake in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		// wait for execution to complete
		resp, err := d.client.ContainerExecAttach(ctx, flakeExec.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msgf("error attaching to flake addition process in sandbox %s: %v", cont.Value().Name, err)
			return
		}
		defer resp.Close() // Ensure stream is closed

		// Inspect the execution result
		execInspect, err := d.client.ContainerExecInspect(ctx, flakeExec.ID)
		if err != nil {
			d.logger.Err(err).Msgf("error inspecting flake addition process in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		if execInspect.ExitCode != 0 {
			d.logger.Error().Msgf("flake addition process failed in sandbox %s with exit code %d", cont.Value().Name, execInspect.ExitCode)
			return
		}

		flakeEvalExec, err := d.client.ContainerExecCreate(ctx, cont.Value().ID,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd: []string{
					"nix",
					"profile",
					"install",
					".",
				},
				WorkingDir: "/home/valnix/work",
			},
		)
		if err != nil {
			d.logger.Err(err).Msgf("error creating flake evaluation process in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		// start execution
		if err := d.client.ContainerExecStart(ctx, flakeEvalExec.ID, container.ExecStartOptions{}); err != nil {
			d.logger.Err(err).Msgf("error starting flake install process in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		// wait for execution to complete
		resp, err = d.client.ContainerExecAttach(ctx, flakeEvalExec.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msgf("error attaching to flake install process in sandbox %s: %v", cont.Value().Name, err)
			return
		}
		defer resp.Close() // Ensure stream is closed

		// Inspect the execution result
		execInspect, err = d.client.ContainerExecInspect(ctx, flakeEvalExec.ID)
		if err != nil {
			d.logger.Err(err).Msgf("error inspecting flake install process in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		if execInspect.ExitCode != 0 {
			d.logger.Error().Msgf("flake install process failed in sandbox %s with exit code %d", cont.Value().Name, execInspect.ExitCode)
			return
		}
	}

	containerURL := fmt.Sprintf("http:/%s-cs.%s", contInfo.Name, d.envConfig.SANDBOX_HOSTNAME)
	sandboxAgentUrl := fmt.Sprintf("ws:/%s-ag.%s/sandbox", contInfo.Name, d.envConfig.SANDBOX_HOSTNAME)

	d.logger.Info().Msg(containerURL)
	err = d.queries.MarkSandboxRunning(ctx, db.MarkSandboxRunningParams{
		SandboxID:       sandBox.SandboxID,
		SandboxUrl:      pgtype.Text{String: containerURL, Valid: true},
		SandboxAgentUrl: pgtype.Text{String: sandboxAgentUrl, Valid: true},
	})
	if err != nil {
		handleError(err, "error marking the container as running")
		return
	}

	d.mu.Lock()
	d.containers = append(d.containers, cont)
	d.mu.Unlock()

	err = d.queries.UpdateJobCompleted(ctx, sandBoxJob.Job.JobID)
	if err != nil {
		d.logger.Err(err).Msgf("error changing job state: %d", sandBoxJob.Job.JobID)
	}
}

func (d *DockerSH) Cleanup(ctx context.Context) error {
	for _, cont := range d.containers {
		cont.Destroy()
	}
	d.containerPool.Close()
	return nil
}
