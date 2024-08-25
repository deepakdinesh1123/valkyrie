package server

import (
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
)

var ServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Odin server",
	Long:  `Start Odin server`,
	RunE:  serverExec,
}

func serverExec(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	logLevel := cmd.Flag("log-level").Value.String()
	logger := logs.GetLogger(logLevel)
	logger.Info().Msg("Starting Odin in standalone mode")

	pg_user := cmd.Flag("pg-user").Value.String()
	pg_password := cmd.Flag("pg-password").Value.String()
	pg_port := cmd.Flag("pg-port").Value.String()
	pg_port_int, err := strconv.ParseUint(pg_port, 10, 32)
	if err != nil {
		logger.Err(err).Msg("Failed to parse pg-port")
		return err
	}
	pg_db := cmd.Flag("pg-db").Value.String()

	envConfig, err := config.GetEnvConfig(
		config.WithPostgresDB(pg_db),
		config.WithPostgresUser(pg_user),
		config.WithPostgresPassword(pg_password),
		config.WithPostgresPort(uint32(pg_port_int)),
	)
	if err != nil {
		logger.Err(err).Msg("Failed to get environment config")
		return err
	}

	applyMigrations, err := cmd.Flags().GetBool("migrate")
	if err != nil {
		logger.Err(err).Msg("Failed to get migrate flag")
		return err
	}

	srv, err := server.NewServer(ctx, envConfig, false, applyMigrations, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to create server")
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	logger.Info().Msg("Starting server")
	go srv.Start(ctx, &wg)
	wg.Wait()
	return nil
}

func init() {
	ServerStartCmd.Flags().Bool("migrate", true, "Migrate database")
}
