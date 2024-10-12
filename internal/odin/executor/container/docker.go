package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
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
		filepath.Join("/home", d.user, "/odin"),
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
	errch := make(chan error)
	success := make(chan bool)

	var out []byte

	go func() {
		resp, err := d.client.ContainerExecCreate(
			ctx,
			containerID,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          command,
				User:         d.user,
			},
		)
		if err != nil {
			done <- true
			errch <- err
			success <- false
		}
		hijResp, err := d.client.ContainerExecAttach(ctx, resp.ID, container.ExecAttachOptions{})
		if err != nil {
			done <- true
			errch <- err
			success <- false
			return
		}
		if hijResp.Reader != nil {
			out, err = io.ReadAll(hijResp.Reader)
			if err != nil {
				done <- true
				return
			}
		}
		done <- true
		success <- true
	}()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				containerStopTimeout := 0
				d.client.ContainerStop(ctx, containerID, container.StopOptions{
					Timeout: &containerStopTimeout,
				})
				reader, err := d.client.ContainerLogs(context.TODO(), containerID, container.LogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Follow:     false,
					Tail:       "all",
				})
				if err != nil {
					return false, "", err
				}
				out, err = io.ReadAll(reader)
				if err != nil {
					return false, "", err
				}
				return false, string(out), nil
			case context.Canceled:
				return false, "", nil
			}
		case <-done:
			return <-success, string(out), nil
		}
	}
}
