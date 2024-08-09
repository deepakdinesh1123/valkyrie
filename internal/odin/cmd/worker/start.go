package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
	"github.com/spf13/cobra"
)

var (
	newWorker bool
)

var WorkerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start worker",
	Long:  `Start worker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logLevel := cmd.Flag("log-level").Value.String()
		logger := logs.GetLogger(logLevel)
		logger.Info().Msg("Starting worker")
		envConfig, err := config.GetEnvConfig()
		if err != nil {
			logger.Err(err).Msg("Failed to get environment config")
			return err
		}
		ctx, cancel := context.WithCancel(cmd.Context())
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-sigChan
			logger.Info().Msg("Shutting down worker")
			cancel()
		}()

		name, err := cmd.Flags().GetString("worker-name")
		if err != nil {
			logger.Err(err).Msg("Failed to get worker-name flag")
		}
		worker, err := worker.GetWorker(ctx, name, envConfig, newWorker, false, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to create worker")
			return err
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go worker.Run(ctx, &wg)
		wg.Wait()
		return nil
	},
}

func init() {
	WorkerStartCmd.Flags().String("worker-name", "", "Name of the worker")
	WorkerStartCmd.Flags().BoolVarP(&newWorker, "new", "n", false, "Create new worker(Deletes existing worker info)")
}
