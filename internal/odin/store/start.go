//go:build all && !darwin

package store

import (
	"context"
	"fmt"
	"os"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func StartOdinStore(ctx context.Context, storeImage, storeContainerName, containerRuntime string, containerEngine string) error {

	switch containerEngine {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return fmt.Errorf("failed to create Docker client: %w", err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get user homeDir: %v", err)
		}

		storeSetupEnvFile := fmt.Sprintf("%s/.odin/store/setup/.env", homeDir)

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
							Source: fmt.Sprintf("%s/.odin/store/nix", homeDir),
							Target: "/nix",
						},
						{
							Type:   mount.TypeBind,
							Source: fmt.Sprintf("%s/.odin/store/setup", homeDir),
							Target: "/tmp/setup",
						},
					},
				}, nil, nil, storeContainerName)
				if err != nil {
					return fmt.Errorf("error creating odin store container: %w", err)
				}
				contInfo, err = cli.ContainerInspect(ctx, storeContainerName)
				if err != nil {
					return fmt.Errorf("error inspecting odin store container: %w", err)
				}
			} else {
				return fmt.Errorf("error inspecting odin store container: %w", err)
			}
		}

		// Ensure contInfo.State is not nil before accessing it
		if contInfo.State != nil && !contInfo.State.Running {
			// Start the container
			if err := cli.ContainerStart(ctx, storeContainerName, container.StartOptions{}); err != nil {
				return fmt.Errorf("error starting odin store container: %w", err)
			}
		}

		return nil
	case "podman":
	default:
		return fmt.Errorf("specified container engine %s not supported", containerEngine)
	}
	return nil
}
