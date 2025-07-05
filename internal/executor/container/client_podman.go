//go:build podman && !darwin

package container

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/executor/container/podman"
)

func GetContainerClient(ctx context.Context, ce *ContainerExecutor) (ContainerClient, error) {
	switch ce.EnvConfig.RUNTIME {
	case "podman":
		containerClient, err := podman.GetPodmanClient(
			ce.EnvConfig,
			ce.Queries,
			ce.WorkerId,
			ce.Tp, ce.Mp,
			ce.Logger,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker containerClient")
		}
		return containerClient, nil
	default:
		return nil, fmt.Errorf("engine not supported")
	}
}
