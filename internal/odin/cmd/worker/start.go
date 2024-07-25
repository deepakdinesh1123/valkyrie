package worker

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/docker"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/system"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/spf13/cobra"
)

var (
	workerDir = "worker"
	infoFile  = "worker-info.json"
)

var WorkerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start worker",
	Long:  `Start worker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logs.GetLogger()
		logger.Info().Msg("Starting worker")
		ctx, cancel := context.WithCancel(cmd.Context())
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-sigChan
			cancel()
		}()

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
			if _, err := os.Stat(envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR); os.IsNotExist(err) {
				err = os.Mkdir(envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, os.ModePerm)
				if err != nil {
					logger.Err(err).Msg("Failed to create system provider base directory")
					return err
				}
			}
			provider, err = system.NewSystemProvider(envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, envConfig.ODIN_SYSTEM_PROVIDER_CLEAN_UP, queries, logger)
			if err != nil {
				logger.Err(err).Msg("Failed to create system provider")
				return err
			}
		default:
			logger.Err(err).Msg("Invalid provider")
			return err
		}

		var wrkr *worker.Worker
		workerInfo, err := readWorkerInfo()
		if err != nil {
			switch err.(type) {
			case *WorkerInfoNotFoundError:
				name, err := cmd.Flags().GetString("name")
				if err != nil {
					logger.Err(err).Msg("Failed to get name flag")
				}
				if name == "" {
					name = namesgenerator.GetRandomName(0)
				}
				wrkr, err = worker.GetWorker(ctx, name, queries, envConfig, provider, logger)
				if err != nil {
					logger.Err(err).Msg("Failed to create worker")
					return err
				}
			default:
				logger.Err(err).Msg("Failed to read worker info")
				return err
			}
		}
		if wrkr == nil && workerInfo != nil {
			wrkr, err = worker.GetWorker(ctx, workerInfo.Name, queries, envConfig, provider, logger)
			if err != nil {
				logger.Err(err).Msg("Failed to get worker")
				return err
			}
		}
		logger.Info().Msgf("Starting worker %d", wrkr.ID)

		err = writeWorkerInfo(wrkr)
		if err != nil {
			logger.Err(err).Msg("Failed to write worker info")
			return err
		}
		err = wrkr.Run(ctx)
		if err != nil {
			logger.Err(err).Msg("Failed to start worker")
			err := deleteWorkerInfo()
			if err != nil {
				logger.Err(err).Msg("Failed to delete worker info")
				return err
			}
			return err
		}
		return nil
	},
}

func init() {
	WorkerStartCmd.Flags().String("name", "", "Name of the worker")
	WorkerStartCmd.Flags().Bool("new", false, "Create new worker(Deletes info of any existing worker)")
}

func writeWorkerInfo(worker *worker.Worker) error {
	wrkrInfo := models.WorkerInfo{
		ID:   worker.ID,
		Name: worker.Name,
	}
	workerInfoBytes, err := json.Marshal(wrkrInfo)
	if err != nil {
		return err
	}
	err = os.WriteFile(infoFile, workerInfoBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func deleteWorkerInfo() error {
	err := os.Remove(infoFile)
	if err != nil {
		return err
	}
	return nil
}

func readWorkerInfo() (*models.WorkerInfo, error) {

	if _, err := os.Stat(infoFile); os.IsNotExist(err) {
		return nil, &WorkerInfoNotFoundError{}
	}
	workerInfoBytes, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var wrkrInfo models.WorkerInfo
	err = json.Unmarshal(workerInfoBytes, &wrkrInfo)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
