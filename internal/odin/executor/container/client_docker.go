//go:build docker

package container

import (
	"context"
	"fmt"
)

func GetContainerClient(ctx context.Context, cp *ContainerExecutor) (ContainerClient, error) {
	switch cp.envConfig.ODIN_CONTAINER_ENGINE {
	case "docker":
		client, err := GetDockerClient(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker client")
		}
		return client, nil
	case " podman":
		return nil, fmt.Errorf("podman engine not supported")
	case "default":
		return nil, fmt.Errorf("invalid client")
	}
	return nil, nil
}
