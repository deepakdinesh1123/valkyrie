//go:build all && !darwin

package store

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
		sock_dir := os.Getenv("XDG_RUNTIME_DIR")
		socket := "unix:" + sock_dir + "/podman/podman.sock"
		pc, err := bindings.NewConnection(context.Background(), socket)

		contInfo, err := containers.Inspect(pc, "odin-store", nil)

		if err != nil {
			return fmt.Errorf("error inspecting odin-store container")
		}

		for _, pkg := range packages {
			err = realisePodmanStore(pkg, contInfo.ID, pc)
			if err != nil {
				logger.Err(err).Msgf("could not realise package %s", pkg.Name)
			}
		}
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
		Cmd:          []string{"/home/valnix/.nix-profile/bin/nix-shell", "-p", pkgName, "--run", "exit 0"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execInfo, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("failed to create exec configuration: %w", err)
	}

	resp, err := cli.ContainerExecAttach(ctx, execInfo.ID, container.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start package: %v", err)
	}
	defer resp.Close()

	// Read the output while the command is running
	outputBuf := new(bytes.Buffer)
	_, err = stdcopy.StdCopy(outputBuf, outputBuf, resp.Reader)
	if err != nil {
		return fmt.Errorf("error reading command output: %w", err)
	}

	// Inspect the exec instance to get the exit code
	execInspect, err := cli.ContainerExecInspect(ctx, execInfo.ID)
	if err != nil {
		return fmt.Errorf("failed to inspect exec instance: %w", err)
	}

	if execInspect.ExitCode != 0 {
		return fmt.Errorf("command exited with non-zero status: %d\nOutput: %s",
			execInspect.ExitCode, outputBuf.String())
	}

	return nil
}

func realisePodmanStore(pkg db.GetAllPackagesRow, containerID string, cli context.Context) error {
	var pkgName string
	if pkg.Pkgtype == "system" {
		pkgName = pkg.Name
	} else {
		pkgName = fmt.Sprintf("%s.%s", pkg.Language.String, pkg.Name)
	}
	fmt.Printf("Realising package %s in container %s\n", pkgName, containerID)

	execId, err := containers.ExecCreate(cli, containerID, &handlers.ExecCreateConfig{
		ExecConfig: container.ExecOptions{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{"/home/valnix/.nix-profile/bin/nix-shell", "-p", pkgName, "--run", "exit 0"},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create exec instance: %w", err)
	}

	err = containers.ExecStartAndAttach(cli, execId, nil)
	if err != nil {
		return fmt.Errorf("failed to start and attach exec instance: %w", err)
	}

	execDetails, err := containers.ExecInspect(cli, execId, nil)
	if err != nil {
		return fmt.Errorf("failed to inspect exec instance: %w", err)
	}

	if execDetails.ExitCode != 0 {
		return fmt.Errorf("exec process failed with exit code %d", execDetails.ExitCode)
	}

	return nil
}
