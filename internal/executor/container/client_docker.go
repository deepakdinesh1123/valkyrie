//go:build docker || darwin

package container

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/executor/container/docker"
)

func GetContainerClient(ctx context.Context, ce *ContainerExecutor) (ContainerClient, error) {
	switch ce.EnvConfig.RUNTIME {
	case "docker":
		containerClient, err := docker.GetDockerProvider(
			ce.EnvConfig,
			ce.Queries,
			ce.WorkerId,
			ce.Tp, ce.Mp,
			ce.Logger,
			ce.Pool,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker containerClient")
		}
		return containerClient, nil
	default:
		return nil, fmt.Errorf("engine not supported")
	}
}
