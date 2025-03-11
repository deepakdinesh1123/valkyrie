//go:build docker

package store

import (
	"context"
	"fmt"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var RealiseCmd = &cobra.Command{
	Use:   "realise",
	Short: "realise odin store",
	Long:  "realise Odin store",
	RunE:  realisePackages,
}

func realisePackages(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return fmt.Errorf("error fetching environment config %s", err)
	}

	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(false),
		db.IsStandalone(false),
		db.IsWorker(false),
	)

	config := logs.NewLogConfig()
	logger := logs.GetLogger(config)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		return fmt.Errorf("error getting DB connection: %s", err)
	}

	packages, err := queries.GetAllPackages(ctx)
	if err != nil {
		logger.Err(err).Msg("could not fetch packages")
	}

	switch envConfig.ODIN_RUNTIME {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return fmt.Errorf("failed to create Docker client: %w", err)
		}

		contInfo, err := cli.ContainerInspect(cmd.Context(), "odin-store")
		for _, pkg := range packages {
			err = realiseDockerStore(pkg, contInfo.ID, cli)
			if err != nil {
				logger.Err(err).Msgf("could not realise package %s", pkg.Name)
			}
		}
	case "podman":
		return fmt.Errorf("podman provider not supported")
	}

	return nil
}

func realiseDockerStore(pkg db.GetAllPackagesRow, containerID string, cli *client.Client) error {
	var pkgName string
	if pkg.Pkgtype == "system" {
		pkgName = pkg.Name
	} else {
		pkgName = fmt.Sprintf("%s.%s", pkg.Language.String, pkg.Name)
	}
	fmt.Printf("Realising package %s in container %s\n", pkgName, containerID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	execConfig := container.ExecOptions{
		Cmd:          []string{"nix-shell", "-p", pkgName, "--run", "exit 0"},
		AttachStdout: true,
		AttachStderr: true,
	}

	respID, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("failed to create exec configuration: %w", err)
	}

	resp, err := cli.ContainerExecAttach(ctx, respID.ID, container.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()

	// Inspect the exec instance to get the exit code
	execInspect, err := cli.ContainerExecInspect(ctx, respID.ID)
	if err != nil {
		return fmt.Errorf("failed to inspect exec instance: %w", err)
	}

	if execInspect.ExitCode != 0 {
		return fmt.Errorf("command exited with non-zero status: %d", execInspect.ExitCode)
	}

	return nil
}
