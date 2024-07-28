package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
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
	logger := logs.GetLogger()
	logger.Info().Msg("Starting Odin in standalone mode")

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get environment config")
		return err
	}

	applyMigrations, err := cmd.Flags().GetBool("migrate")
	if err != nil {
		logger.Err(err).Msg("Failed to get migrate flag")
		return err
	}

	queries, err := db.GetDBConnection(ctx, false, envConfig, applyMigrations, nil, nil, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to get database connection")
		return err
	}
	srv := server.NewServer(ctx, envConfig, queries, logger)

	logger.Info().Msg("Starting Odin server")
	addr := fmt.Sprintf("%s:%s", envConfig.ODIN_SERVER_HOST, envConfig.ODIN_SERVER_PORT)
	err = http.ListenAndServe(addr, srv)
	if err != nil {
		logger.Err(err).Msg("Failed to start server")
		return err
	}
	return nil
}

func init() {
	ServerStartCmd.Flags().Bool("migrate", true, "Migrate database")
}
