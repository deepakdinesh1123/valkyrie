package container

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db/jsonschema"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/pool"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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

func (d *DockerSH) StartContainerPool(ctx context.Context, envConfig *config.EnvConfig) error {
	containerPool, err := pool.NewSandboxPool(ctx, int32(envConfig.ODIN_HOT_CONTAINER), int32(envConfig.ODIN_WORKER_CONCURRENCY), envConfig.ODIN_CONTAINER_ENGINE)
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

	// CodeServerConfig, err := sandbox.GetCodeServerConfig()
	// if err != nil {
	// 	d.logger.Err(err).Msg("could not get sandbox config")
	// 	return
	// }
	// configYaml, err := yaml.Marshal(CodeServerConfig)
	// if err != nil {
	// 	d.logger.Err(err).Msg("could not convert sandbox config to yaml")
	// }

	// d.logger.Info().Str("Command", fmt.Sprintf("echo \"%s\" >> /home/valnix/.config/code-server/config.yaml", string(configYaml))).Msg("Adding config")
	// configExec, err := d.client.ContainerExecCreate(
	// 	ctx,
	// 	cont.Value().ID,
	// 	container.ExecOptions{
	// 		AttachStderr: true,
	// 		AttachStdout: true,
	// 		Cmd: []string{
	// 			"/bin/sh",
	// 			"-c",
	// 			fmt.Sprintf("echo \"%s\" >> /home/valnix/.config/code-server/config.yaml", string(configYaml)),
	// 		},
	// 	},
	// )
	// _, err = d.client.ContainerExecAttach(ctx, configExec.ID, container.ExecAttachOptions{})
	// if err != nil {
	// 	d.logger.Err(err).Msg("error adding code server config")
	// 	return
	// }

	// codeServerExec, err := d.client.ContainerExecCreate(
	// 	ctx,
	// 	cont.Value().ID,
	// 	container.ExecOptions{
	// 		AttachStderr: true,
	// 		AttachStdout: true,
	// 		Cmd: []string{
	// 			"/home/valnix/start.sh",
	// 		},
	// 	},
	// )
	// err = d.client.ContainerExecStart(ctx, codeServerExec.ID, container.ExecStartOptions{
	// 	Detach: true,
	// })
	// if err != nil {
	// 	d.logger.Err(err).Msg("error adding code server config")
	// 	return
	// }

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
					fmt.Sprintf("rm /home/valnix/work/flake.nix && echo \"%s\" >> /home/valnix/work/flake.nix", sandBox.Config.Flake),
				},
			},
		)

		_, err = d.client.ContainerExecAttach(ctx, flakeExec.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msgf("error adding flake to sandbox %s: %v", cont.Value().Name, err)
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
				WorkingDir: "/home/valnix",
			},
		)
		if err != nil {
			d.logger.Err(err).Msgf("error creating flake evaluation process in sandbox %s: %v", cont.Value().Name, err)
			return
		}

		err = d.client.ContainerExecStart(ctx, flakeEvalExec.ID, container.ExecStartOptions{Detach: true})
		if err != nil {
			d.logger.Err(err).Msgf("error evaluating flake in sandbox %s: %v", cont.Value().Name, err)
			return
		}
	}

	containerURL := fmt.Sprintf("http://%s", contInfo.NetworkSettings.Networks["bridge"].IPAddress)
	sandboxAgentUrl := fmt.Sprintf("ws://%s:1618/sandbox", contInfo.NetworkSettings.Networks["bridge"].IPAddress)
	if d.envConfig.ODIN_ENVIRONMENT == "prod" {
		// Handle prod url configuration here
	}

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

func (d *DockerSH) StartOdinStore(ctx context.Context, storeImage, storeContainerName, containerRuntime string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user homeDir: %v", err)
	}
	contInfo, err := d.client.ContainerInspect(ctx, storeContainerName)
	if err != nil {
		if client.IsErrNotFound(err) { // Container doesn't exist, create it
			_, err = d.client.ContainerCreate(ctx, &container.Config{Image: storeImage}, &container.HostConfig{
				Runtime:     containerRuntime,
				NetworkMode: "bridge",
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/.odin/store/nix", homeDir),
						Target: "/nix",
					},
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/.odin/store/setup", homeDir),
						Target: "/tmp/setup",
					},
				},
			}, nil, nil, storeContainerName)
			if err != nil {
				return fmt.Errorf("error creating odin store container: %w", err)
			}
			contInfo, err = d.client.ContainerInspect(ctx, storeContainerName)
			if err != nil {
				return fmt.Errorf("error inspecting odin store container: %w", err)
			}
		} else {
			return fmt.Errorf("error inspecting odin store container: %w", err)
		}
	}

	// Ensure contInfo.State is not nil before accessing it
	if contInfo.State != nil && !contInfo.State.Running {
		// Start the container
		if err := d.client.ContainerStart(ctx, storeContainerName, container.StartOptions{}); err != nil {
			return fmt.Errorf("error starting odin store container: %w", err)
		}
	}

	return nil
}
