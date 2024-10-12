package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/puddle/v2"
)

type PodmanClient struct {
	connection context.Context
	*ContainerExecutor
}

var getPodmanConnectionOnce sync.Once
var podmanConnection context.Context

func GetPodmanClient(cp *ContainerExecutor) (*PodmanClient, error) {
	connection := getPodmanConnection()
	if connection == nil {
		return nil, fmt.Errorf("could not get podman connection")
	}
	return &PodmanClient{
		connection:        connection,
		ContainerExecutor: cp,
	}, nil
}

func getPodmanConnection() context.Context {
	getPodmanConnectionOnce.Do(func() {
		// sock_dir := os.Getenv("XDG_RUNTIME_DIR")
		// socket := "unix:" + sock_dir + "/podman/podman.sock"
		pc, err := bindings.NewConnection(context.Background(), "unix:/run/podman/podman.sock")
		if err != nil {
			return
		}
		podmanConnection = pc
	})
	return podmanConnection
}

func (p *PodmanClient) WriteFiles(ctx context.Context, containerID string, prepDir string, job db.Job) error {
	execReq, err := p.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		return err
	}
	files := map[string]string{
		"flake.nix":  execReq.Flake,
		execReq.Path: execReq.Code,
	}

	tarFilePath, err := CreateTarArchive(files, prepDir)
	if err != nil {
		return err
	}

	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return fmt.Errorf("Failed to open tar file")
	}
	defer tarFile.Close()
	defer os.Remove(tarFilePath)

	copyF, err := containers.CopyFromArchive(p.connection, containerID, "/odin", tarFile)
	if err != nil {
		return fmt.Errorf("Failed to copy files to container")
	}
	err = copyF()
	if err != nil {
		return fmt.Errorf("Failed to copy files to container")
	}
	return nil
}

func (p *PodmanClient) GetContainer(ctx context.Context) (*puddle.Resource[Container], error) {
	cont, err := p.ContainerExecutor.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	err = p.pool.CreateResource(ctx)
	if err != nil {
		p.logger.Debug().Msg("Container pool might be full")
	}
	return cont, nil
}

func (p *PodmanClient) Execute(ctx context.Context, containerID string, command []string) (bool, string, error) {
	done := make(chan bool, 1)
	success := make(chan bool, 1)
	errch := make(chan error)

	var output bytes.Buffer

	go func() {
		execId, err := containers.ExecCreate(p.connection, containerID, &handlers.ExecCreateConfig{
			ExecConfig: container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				WorkingDir:   "/odin",
				Cmd:          []string{"nix", "run"},
			},
		})
		if err != nil {
			done <- true
			errch <- err
			success <- false
			return
		}

		r, w, err := os.Pipe()
		if err != nil {
			done <- true
			errch <- err
			success <- false
			return
		}
		defer r.Close()
		execOpts := new(containers.ExecStartAndAttachOptions).WithOutputStream(w).WithAttachOutput(true).WithErrorStream(w)
		err = containers.ExecStartAndAttach(p.connection, execId, execOpts)
		if err != nil {
			done <- true
			errch <- err
			success <- false
			return
		}

		// copying output in a separate goroutine, so that printing doesn't remain blocked forever
		go func() {
			for {
				if r == nil {
					return
				}
				_, err = io.Copy(&output, r)
				if err != nil {
					return
				}
			}
		}()
		w.Close()
		r.Close()
		done <- true
	}()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				containerStopTimeout := uint(0)
				containers.Stop(ctx, containerID, &containers.StopOptions{
					Timeout: &containerStopTimeout,
				})
				return false, output.String(), nil
			case context.Canceled:
				return false, "", nil
			}
		case <-done:
			return <-success, output.String(), nil
		}
	}
}
