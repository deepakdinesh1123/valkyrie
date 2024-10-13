package container

import (
	"bytes"
	"context"
	"fmt"
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

	var execId string
	var err error

	go func() {
		defer func() {
			done <- true
		}()

		execId, err = containers.ExecCreate(p.connection, containerID, &handlers.ExecCreateConfig{
			ExecConfig: container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          command,
			},
		})
		p.logger.Info().Msg("Exec created")
		if err != nil {
			return
		}

		p.logger.Info().Msg("Starting execution")
		err = containers.ExecStart(p.connection, execId, nil)
		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				p.logger.Info().Msg("Timelimit exceed")
				return
			default:
				execInfo, err := containers.ExecInspect(p.connection, execId, nil)
				if err != nil {
					p.logger.Err(err).Msg("Could not inspect exec")
					return
				}
				if !execInfo.Running {
					p.logger.Info().Int("Exit Code", execInfo.ExitCode).Msg("Execution process exit")
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			p.logger.Info().Msg("Killing process")
			if execId != "" {
				err := containers.ExecRemove(
					p.connection,
					execId,
					(&containers.ExecRemoveOptions{}).WithForce(true),
				)
				if err != nil {
					return false, "", fmt.Errorf("error stopping container: %s", err)
				}
			}
			out, err := p.ReadExecLogs(containerID)
			if err != nil {
				return false, "", fmt.Errorf("error reading output: %s", err)
			}
			return true, out, nil
		case context.Canceled:
			return false, "", nil
		}
	case <-done:
		out, err := p.ReadExecLogs(containerID)
		return true, out, err
	}
	return false, "", nil
}

func (p *PodmanClient) ReadExecLogs(containerID string) (string, error) {
	var execLogs bytes.Buffer

	execId, err := containers.ExecCreate(p.connection, containerID, &handlers.ExecCreateConfig{
		ExecConfig: container.ExecOptions{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{"bash", "-c", "cat ~/odin/output.txt"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("Could not create exec: %s", err)
	}
	err = containers.ExecStartAndAttach(
		p.connection,
		execId,
		(&containers.ExecStartAndAttachOptions{}).WithAttachOutput(true).WithOutputStream(&execLogs),
	)
	if err != nil {
		return "", fmt.Errorf("Could not attach to exec")
	}
	p.logger.Info().Msg(execLogs.String())
	return execLogs.String(), nil
}
