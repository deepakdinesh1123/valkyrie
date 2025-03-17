//go:build podman && !darwin

package store

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func StartOdinStore(ctx context.Context, storeImage, storeContainerName string) error {
	sysCtx, err := bindings.NewConnection(ctx, "unix:/run/podman/podman.sock")
	if err != nil {
		return fmt.Errorf("failed to connect to Podman socket: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user homeDir: %v", err)
	}

	storeSetupEnvFile := fmt.Sprintf("%s/.odin/store/setup/.env", homeDir)

	if _, err := os.Stat(storeSetupEnvFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", storeSetupEnvFile)
	}

	// Check if container exists
	contInfo, err := containers.Inspect(sysCtx, storeContainerName, nil)
	if err != nil {
		// Container does not exist, create it
		s := specgen.NewSpecGenerator(storeImage, false)
		s.Name = storeContainerName
		s.Env = map[string]string{ /* Load env variables from file */ }
		s.Mounts = []specs.Mount{
			{
				Source:      fmt.Sprintf("%s/.odin/store/nix", homeDir),
				Destination: "/nix",
				Type:        "bind",
			},
			{
				Source:      fmt.Sprintf("%s/.odin/store/setup", homeDir),
				Destination: "/tmp/setup",
				Type:        "bind",
			},
		}
		_, err := containers.CreateWithSpec(sysCtx, s, nil)
		if err != nil {
			return fmt.Errorf("error creating odin store container: %w", err)
		}
		contInfo, err = containers.Inspect(sysCtx, storeContainerName, nil)
		if err != nil {
			return fmt.Errorf("error inspecting odin store container: %w", err)
		}
	}

	if contInfo.State != nil && contInfo.State.Running == false {
		// Start the container
		if err := containers.Start(sysCtx, storeContainerName, nil); err != nil {
			return fmt.Errorf("error starting odin store container: %w", err)
		}
	}

	return nil
}
