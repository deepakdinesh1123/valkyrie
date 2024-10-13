package container

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	go func() {
		defer func() {
			done <- true
		}()

		execId, err := containers.ExecCreate(p.connection, containerID, &handlers.ExecCreateConfig{
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

		select {
		case <-ctx.Done():
			p.logger.Info().Msg("Timelimit exceed")
			return
		default:
			p.logger.Info().Msg("Inspecting")
			execInfo, err := containers.ExecInspect(p.connection, execId, nil)
			if err != nil {
				p.logger.Err(err).Msg("Could not inspect exec")
				return
			}
			if !execInfo.Running {
				p.logger.Info().Int("Exit Code", execInfo.ExitCode).Msg("Execution process exit")
				return
			}
			p.logger.Info().Msg("Exec session running")
		}
	}()

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			err := containers.Stop(
				context.TODO(),
				containerID,
				(&containers.StopOptions{}).WithTimeout(0),
			)
			if err != nil {
				return false, "", fmt.Errorf("Error stopping container")
			}
			out, err := p.ReadContainerLogs(containerID)
			if err != nil {
				return false, "", fmt.Errorf("Error reading output")
			}
			return false, out, nil
		case context.Canceled:
			return false, "", nil
		}
	case <-done:
		out, err := p.ReadContainerLogs(containerID)
		return true, out, err
	}
	return false, "", nil
}

func (p *PodmanClient) ReadContainerLogs(containerID string) (string, error) {

	var logs []string

	logsBuffer := 200

	done := make(chan bool)
	logout := make(chan string, logsBuffer)
	logerr := make(chan string, logsBuffer)

	logAppender := func() {
		for {
			select {
			case msg := <-logout:
				logs = append(logs, msg)
			case msg := <-logerr:
				logs = append(logs, msg)
			case <-done:
				return
			}
		}
	}

	go logAppender()

	options := new(containers.LogOptions).WithFollow(false)

	err := containers.Logs(p.connection, containerID, options, logout, logerr)
	if err != nil {
		return "", err
	}
	done <- true

	return strings.Join(logs, "\n"), nil
}
