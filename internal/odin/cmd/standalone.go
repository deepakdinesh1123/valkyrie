package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
)

var (
	cleanDB bool
)

var StandaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Odin standalone mode",
	Long:  `Start Odin in standalone mode`,
	RunE:  standaloneExec,
}

func standaloneExec(cmd *cobra.Command, args []string) error {
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		log.Err(err).Msg("Failed to get environment config")
		return err
	}

	logLevel := cmd.Flag("log-level").Value.String()
	config := logs.NewLogConfig(logs.WithLevel(logLevel), logs.WithExport(envConfig.ODIN_EXPORT_LOGS))
	logger := logs.GetLogger(config)
	logger.Info().Msg("Starting Odin in standalone mode")

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	if cleanDB {
		os.RemoveAll(envConfig.POSTGRES_STANDALONE_PATH)
	}

	srv, err := server.NewServer(ctx, envConfig, true, true, logger)
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
	logger.Info().Msg("Odin server and worker stopped")
	return nil
}

func init() {
	StandaloneCmd.Flags().BoolVarP(&cleanDB, "clean-db", "c", false, "Delete existing DB")
}
