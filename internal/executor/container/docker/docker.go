//go:build docker || all || darwin

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/executor/container/common"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/deepakdinesh1123/valkyrie/internal/services/execution"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DockerProvider struct {
	client    *client.Client
	queries   db.Store
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	// user          string
	containerPool *puddle.Pool[pool.Container]
}

func GetDockerProvider(envConfig *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger, containerPool *puddle.Pool[pool.Container]) (*DockerProvider, error) {
	client := pool.GetDockerClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &DockerProvider{
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

func (d *DockerProvider) WriteFiles(ctx context.Context, containerID string, prepDir string, job *db.Job) error {

	execReq, err := d.queries.GetExecRequest(ctx, job.Arguments.ExecConfig.ExecReqId)
	if err != nil {
		return err
	}
	script, spec, err := execution.ConvertExecSpecToNixScript(ctx, &execReq, d.queries)
	if err != nil {
		return fmt.Errorf("error writing files: %s", err)
	}
	files := map[string]string{
		"exec.sh":       script,
		spec.ScriptName: execReq.Code.String,
		"input.txt":     execReq.Input.String,
	}

	tarFilePath, err := common.CreateTarArchive(files, execReq.Files, prepDir)
	if err != nil {
		return err
	}
	defer os.Remove(tarFilePath)

	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		d.logger.Err(err).Msg("Failed to open tar file")
		return err
	}
	defer tarFile.Close()
	err = d.client.CopyToContainer(
		ctx,
		containerID,
		filepath.Join("/home/valnix/valkyrie"),
		tarFile,
		container.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: true},
	)
	if err != nil {
		d.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	return nil
}

func (d *DockerProvider) GetContainer(ctx context.Context) (*puddle.Resource[pool.Container], error) {
	cont, err := d.containerPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	go d.containerPool.CreateResource(ctx)
	return cont, nil
}

func (d *DockerProvider) Execute(ctx context.Context, containerID string, command []string) (bool, string, error) {
	done := make(chan bool)

	var dexec types.IDResponse
	var err error

	go func() {
		defer func() {
			done <- true
		}()

		dexec, err = d.client.ContainerExecCreate(
			ctx,
			containerID,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          command,
			},
		)
		if err != nil {
			return
		}
		err = d.client.ContainerExecStart(ctx, dexec.ID, container.ExecAttachOptions{})
		if err != nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				d.logger.Info().Msg("Timelimit exceeded")
				return
			default:
				execInfo, err := d.client.ContainerExecInspect(ctx, dexec.ID)
				if err != nil {
					return
				}
				if !execInfo.Running {
					d.logger.Info().Int("Exit Code", execInfo.ExitCode).Msg("Execution process exit")
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			if dexec.ID != "" {
				stopExec, err := d.client.ContainerExecCreate(
					context.TODO(),
					containerID,
					container.ExecOptions{
						Cmd: []string{"sh", "nix_stop.sh"},
					},
				)
				if err != nil {
					return false, "", fmt.Errorf("could not create exec for nix stop: %s", err)
				}
				err = d.client.ContainerExecStart(context.TODO(), stopExec.ID, container.ExecStartOptions{})
				if err != nil {
					return false, "", fmt.Errorf("could not start the nix_stop script: %s", err)
				}
			}
			out, err := d.ReadExecLogs(context.TODO(), containerID)
			if err != nil {
				return false, "", fmt.Errorf("error reading output: %s", err)
			}
			return true, out, nil
		case context.Canceled:
			return false, "", fmt.Errorf("context canceled")
		}
	case <-done:
		out, err := d.ReadExecLogs(context.TODO(), containerID)
		if err != nil {
			return false, "", fmt.Errorf("error reading output: %s", err)
		}
		return true, out, nil
	}
	return false, "", nil
}

func (d *DockerProvider) ReadExecLogs(ctx context.Context, containerID string) (string, error) {
	var out []byte
	dexec, err := d.client.ContainerExecCreate(
		ctx,
		containerID,
		container.ExecOptions{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"sh", "-c", "cat ~/valkyrie/output.txt"},
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not create exec: %s", err)
	}
	resp, err := d.client.ContainerExecAttach(ctx, dexec.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("could not attach to container: %s", err)
	}
	if resp.Reader != nil {
		out, err = io.ReadAll(resp.Reader)
		if err != nil {
			return "", fmt.Errorf("could not read from hijacked response: %s", err)
		}
	}
	return stripCtlAndExtFromUTF8(string(out)), nil
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 || r == 10 {
			return r
		}
		return -1
	}, str)
}
