package container

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
	*ContainerProvider
}

func GetDockerClient(cp *ContainerProvider) (*DockerClient, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		client:            client,
		ContainerProvider: cp,
	}, nil
}

func newClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func (d *DockerClient) WriteFiles(ctx context.Context, containerID string, prepDir string, job db.Job) error {
	return nil
}

func (d *DockerClient) GetContainer(ctx context.Context, prepDir string) (Container, error) {
	return Container{}, nil
}

func (d *DockerClient) Execute(ctx context.Context, containerID string, command []string) error {
	return nil
}

func (d *DockerClient) ReadOutput(ctx context.Context, containerID string) (string, error) {
	return "", nil
}
