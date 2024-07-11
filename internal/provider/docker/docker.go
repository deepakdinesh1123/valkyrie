package docker

import (
	"github.com/docker/docker/client"
)

type DockerProvider struct {
	client *client.Client
}

func NewDockerProvider() (*DockerProvider, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return &DockerProvider{
		client: client,
	}, nil
}
