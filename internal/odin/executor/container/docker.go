//go:build docker

package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jackc/puddle/v2"
)

type DockerProvider struct {
	client *client.Client
	*ContainerExecutor
}

var getDockerClientOnce sync.Once
var dockerclient *client.Client

func GetDockerProvider(ce *ContainerExecutor) (*DockerProvider, error) {
	client := getDockerClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &DockerProvider{
		client:            client,
		ContainerExecutor: ce,
	}, nil
}

func getDockerClient() *client.Client {
	getDockerClientOnce.Do(
		func() {
			c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return
			}
			dockerclient = c
		},
	)
	return dockerclient
}

func (d *DockerProvider) WriteFiles(ctx context.Context, containerID string, prepDir string, job db.Job) error {
	execReq, err := d.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		return err
	}
	files := map[string]string{
		"exec.sh":    execReq.NixScript,
		execReq.Path: execReq.Code,
	}

	tarFilePath, err := CreateTarArchive(files, prepDir)
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
		filepath.Join("/home/valnix/odin"),
		tarFile,
		container.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: true},
	)
	if err != nil {
		d.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	return nil
}

func (d *DockerProvider) GetContainer(ctx context.Context) (*puddle.Resource[Container], error) {
	cont, err := d.ContainerExecutor.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	err = d.pool.CreateResource(ctx)
	if err != nil {
		d.logger.Debug().Msg("Container pool might be full")
	}
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
				d.logger.Info().Msg("Timelimit exced")
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
			Cmd:          []string{"sh", "-c", "cat ~/odin/output.txt"},
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
	return string(out), nil
}
