package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/server"
	"github.com/deepakdinesh1123/valkyrie/internal/worker"
)

var (
	cleanDB      bool
	migrateDB    bool
	initialiseDB bool
)

var StandaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "valkyrie standalone mode",
	Long:  `Start valkyrie in standalone mode`,
	RunE:  standaloneExec,
}

func standaloneExec(cmd *cobra.Command, args []string) error {
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		log.Err(err).Msg("Failed to get environment config")
		return err
	}

	logLevel := cmd.Flag("log-level").Value.String()
	config := logs.NewLogConfig(logs.WithLevel(logLevel), logs.WithExport(envConfig.EXPORT_LOGS))
	logger := logs.GetLogger(config)
	logger.Info().Msg("Starting valkyrie in standalone mode")

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	if cleanDB {
		os.RemoveAll(envConfig.POSTGRES_STANDALONE_PATH)
	}

	srv, err := server.NewServer(ctx, envConfig, true, true, initialiseDB, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to create server")
		return err
	}

	worker, err := worker.GetWorker(ctx, "", envConfig, true, true, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to create worker")
		return err
	}
	go func() {
		<-sigChan
		logger.Info().Msg("Shutting down worker")
		logger.Info().Msg("Removing lock")
		cancel()
	}()
	var wg sync.WaitGroup

	wg.Add(1)
	logger.Info().Msg("Starting worker")
	go worker.Run(ctx, &wg)

	wg.Add(1)
	logger.Info().Msg("Starting server")
	go srv.Start(ctx, &wg)

	wg.Wait()
	logger.Info().Msg("valkyrie server and worker stopped")
	return nil
}

func init() {
	StandaloneCmd.Flags().BoolVarP(&cleanDB, "clean-db", "c", false, "Delete existing DB")
	StandaloneCmd.Flags().BoolVarP(&migrateDB, "migrate", "m", true, "Migrate database")
	StandaloneCmd.Flags().BoolVarP(&initialiseDB, "initdb", "i", true, "Initialise database")
}
