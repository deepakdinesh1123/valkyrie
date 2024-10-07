//go:build podman && !darwin

package container

import (
	"context"
	"fmt"
)

func GetContainerClient(ctx context.Context, cp *ContainerProvider) (ContainerClient, error) {
	switch cp.envConfig.ODIN_CONTAINER_ENGINE {
	case "docker":
		return nil, fmt.Errorf("docker engine not supported")
	case " podman":
		client, err := GetPodmanClient(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to create podman client")
		}
		return client, nil
	case "default":
		return nil, fmt.Errorf("invalid client")
	}
	return nil, nil
}
