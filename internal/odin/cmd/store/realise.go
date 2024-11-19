package store

import (
	"os/exec"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
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
		return logError("Error fetching environment config", err)
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
		return logError("Error getting DB connection", err)
	}

	packages, err := queries.GetAllPackages(ctx)
	if err != nil {
		logger.Err(err).Msg("could not fetch packages")
	}

	for _, pkg := range packages {
		err = realise(pkg)
		if err != nil {
			logger.Err(err).Msgf("could not realise package %s", pkg.Name)
		}
	}

	return nil
}

func realise(pkg db.GetAllPackagesRow) error {
	cmd := exec.Command("nix-shell", "-p", pkg.Name, "--run", "exit")
	err := cmd.Run()
	return err
}
