package cmd

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/docker"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/system"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
	"github.com/spf13/cobra"
)

var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start worker",
	Long:  `Start worker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		logger := logs.GetLogger()
		logger.Info().Msg("Starting worker")

		envConfig, err := config.GetEnvConfig()
		if err != nil {
			logger.Err(err).Msg("Failed to get environment config")
			return err
		}
		_, queries, err := db.GetDBConnection(ctx, false, envConfig, false, nil, nil, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to get database connection")
			return err
		}
		var provider provider.Provider
		switch envConfig.ODIN_WORKER_PROVIDER {
		case "docker":
			provider, err = docker.NewDockerProvider()
			if err != nil {
				logger.Err(err).Msg("Failed to create docker provider")
				return err
			}
		case "system":
			provider, err = system.NewSystemProvider()
			if err != nil {
				logger.Err(err).Msg("Failed to create system provider")
				return err
			}
		default:
			logger.Err(err).Msg("Invalid provider")
			return err
		}
		worker := worker.NewWorker(ctx, queries, envConfig, provider, logger)
		logger.Info().Msg("Starting worker")
		err = worker.Run(ctx)
		if err != nil {
			logger.Err(err).Msg("Failed to start worker")
			return err
		}
		return nil
	},
}
