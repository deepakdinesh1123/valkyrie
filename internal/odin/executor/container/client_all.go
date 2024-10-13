//go:build linux && !(docker || podman)

package container

import (
	"context"
	"fmt"
)

func GetContainerClient(ctx context.Context, cp *ContainerExecutor) (ContainerClient, error) {
	switch cp.envConfig.ODIN_CONTAINER_ENGINE {
	case "docker":
		containerClient, err := GetDockerProvider(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker containerClient")
		}
		return containerClient, nil
	case "podman":
		containerClient, err := GetPodmanClient(cp)
		if err != nil {
			return nil, fmt.Errorf("failed to create podman containerClient")
		}
		return containerClient, nil
	default:
		return nil, fmt.Errorf("invalid containerClient")
	}
}
