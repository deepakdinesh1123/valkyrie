package server

import (
	"sync"

	"github.com/rs/zerolog/log"
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

var (
	migrateDB    bool
	initialiseDB bool
)

func serverExec(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		log.Err(err).Msg("Failed to get env config")
		return err
	}

	logLevel := cmd.Flag("log-level").Value.String()
	config := logs.NewLogConfig(
		logs.WithLevel(logLevel),
		logs.WithExport(envConfig.ODIN_EXPORT_LOGS),
		logs.WithSource("server"),
	)
	logger := logs.GetLogger(config)
	logger.Info().Msg("Starting Odin server")

	srv, err := server.NewServer(ctx, envConfig, false, migrateDB, initialiseDB, logger)
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
	ServerStartCmd.Flags().BoolVarP(&migrateDB, "migrate", "m", true, "Migrate database")
	ServerStartCmd.Flags().BoolVarP(&initialiseDB, "initdb", "i", true, "Initialise database")
}
