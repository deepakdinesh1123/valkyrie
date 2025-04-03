//go:build docker

package store

import (
	"context"
	"fmt"
	"os"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func StartValkyrieStore(ctx context.Context, storeImage, storeContainerName, containerRuntime string, containerEngine string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user homeDir: %v", err)
	}

	storeSetupEnvFile := fmt.Sprintf("%s/.valkyrie/store/setup/.env", homeDir)

	if _, err := os.Stat(storeSetupEnvFile); os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: %s\n", storeSetupEnvFile)
	}

	storeSetupEnv, err := config.LoadEnvFile(storeSetupEnvFile)
	if err != nil {
		return fmt.Errorf("error reading store setup env file")
	}

	contInfo, err := cli.ContainerInspect(ctx, storeContainerName)
	if err != nil {
		if client.IsErrNotFound(err) { // Container doesn't exist, create it
			_, err = cli.ContainerCreate(ctx, &container.Config{Image: storeImage, Env: storeSetupEnv}, &container.HostConfig{
				Runtime:     containerRuntime,
				NetworkMode: "bridge",
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/.valkyrie/store/nix", homeDir),
						Target: "/nix",
					},
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/.valkyrie/store/setup", homeDir),
						Target: "/tmp/setup",
					},
				},
			}, nil, nil, storeContainerName)
			if err != nil {
				return fmt.Errorf("error creating valkyrie store container: %w", err)
			}
			contInfo, err = cli.ContainerInspect(ctx, storeContainerName)
			if err != nil {
				return fmt.Errorf("error inspecting valkyrie store container: %w", err)
			}
		} else {
			return fmt.Errorf("error inspecting valkyrie store container: %w", err)
		}
	}

	// Ensure contInfo.State is not nil before accessing it
	if contInfo.State != nil && !contInfo.State.Running {
		// Start the container
		if err := cli.ContainerStart(ctx, storeContainerName, container.StartOptions{}); err != nil {
			return fmt.Errorf("error starting valkyrie store container: %w", err)
		}
	}

	return nil
}
