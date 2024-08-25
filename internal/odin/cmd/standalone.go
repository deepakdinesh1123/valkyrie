package cmd

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
)

var (
	newWorker bool
)

var StandaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Odin standalone mode",
	Long:  `Start Odin in standalone mode`,
	RunE:  standaloneExec,
}

func standaloneExec(cmd *cobra.Command, args []string) error {
	logLevel := cmd.Flag("log-level").Value.String()
	logger := logs.GetLogger(logLevel)
	logger.Info().Msg("Starting Odin in standalone mode")

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

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

	srv, err := server.NewServer(ctx, envConfig, true, applyMigrations, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to create server")
		return err
	}
	name, err := cmd.Flags().GetString("worker-name")
	if err != nil {
		logger.Err(err).Msg("Failed to get worker-name flag")
		name = ""
	}

	worker, err := worker.GetWorker(ctx, name, envConfig, newWorker, true, logger)
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
	StandaloneCmd.Flags().Bool("migrate", true, "Migrate database")
	StandaloneCmd.Flags().Bool("clean-db", false, "Clean database")
	StandaloneCmd.Flags().String("worker-name", "", "Name of the worker")
	StandaloneCmd.Flags().BoolVarP(&newWorker, "new", "n", false, "Create new worker(Deletes existing worker info)")
}
