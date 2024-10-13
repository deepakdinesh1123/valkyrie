package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		sock_dir := os.Getenv("XDG_RUNTIME_DIR")
		socket := "unix:" + sock_dir + "/podman/podman.sock"
		pc, err := bindings.NewConnection(context.Background(), socket)
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
		"exec.sh":    execReq.NixScript,
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

	copyF, err := containers.CopyFromArchive(
		p.connection,
		containerID,
		filepath.Join("/home", p.user, "odin"),
		tarFile,
	)
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
		p.logger.Err(err).Msg("Error when acquiring container")
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

	var output string

	go func(ctx context.Context) {
		defer func() {
			done <- true
		}()

		select {
		case <-ctx.Done():
			p.logger.Info().Msg("Context cancelled before execution")
			success <- false
			return
		default:
			p.logger.Info().Msg("Exec created")
			execId, err := containers.ExecCreate(p.connection, containerID, &handlers.ExecCreateConfig{
				ExecConfig: container.ExecOptions{
					AttachStderr: true,
					AttachStdout: true,
					Cmd:          command,
				},
			})
			if err != nil {
				errch <- err
				success <- false
				return
			}

			r, w, err := os.Pipe()
			if err != nil {
				errch <- err
				success <- false
				return
			}
			defer r.Close()
			defer w.Close() // Ensure the writer is closed

			p.logger.Info().Msg("Starting execution")
			execOpts := new(containers.ExecStartAndAttachOptions).WithOutputStream(w).WithAttachOutput(true).WithErrorStream(w)
			err = containers.ExecStartAndAttach(p.connection, execId, execOpts)
			if err != nil {
				errch <- err
				success <- false
				return
			}

			outputC := make(chan string)

			// Copying output in a separate goroutine, so that printing doesn't remain blocked forever
			go func() {
				select {
				case <-ctx.Done():
					p.logger.Info().Msg("Context cancelled during output copying")
					return
				default:
					var output bytes.Buffer
					p.logger.Info().Msg("Copying output")
					_, _ = io.Copy(&output, r)
					outputC <- output.String()
					p.logger.Info().Msg("Output copied")
				}
			}()

			// Wait for output or context cancellation
			select {
			case out := <-outputC:
				p.logger.Info().Msg("Output received")
				success <- true
				output = out
			case <-ctx.Done():
				p.logger.Info().Msg("Context cancelled while waiting for output")
				success <- false
				return
			}
		}
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			p.logger.Info().Msg("Timeout killing process")
			stderr := true
			stdout := true
			stdoutChan := make(chan string)
			stderrChan := make(chan string)
			containers.Logs(p.connection, containerID, &containers.LogOptions{
				Stderr: &stderr,
				Stdout: &stdout,
			}, stdoutChan, stderrChan)
			p.logger.Info().Str("output", <-stdoutChan).Msg("Container logs")
			return false, "", nil
		case context.Canceled:
			return false, "", nil
		}
	case <-done:
		return <-success, string(output), nil
	}
	return false, "", nil
}
