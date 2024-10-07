package container

import (
	"context"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

type PodmanClient struct {
	connection context.Context
	*ContainerProvider
}

func GetPodmanClient(cp *ContainerProvider) (*PodmanClient, error) {
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	connection, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return nil, err
	}
	return &PodmanClient{
		connection:        connection,
		ContainerProvider: cp,
	}, nil
}

func (d *PodmanClient) WriteFiles(ctx context.Context, containerID string, prepDir string, job db.Job) error {
	return nil
}

func (d *PodmanClient) GetContainer(ctx context.Context, prepDir string) (Container, error) {
	return Container{}, nil
}

func (d *PodmanClient) Execute(ctx context.Context, containerID string, command []string) error {
	return nil
}

func (d *PodmanClient) ReadOutput(ctx context.Context, containerID string) (string, error) {
	return "", nil
}
