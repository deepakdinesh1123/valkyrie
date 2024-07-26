package worker

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofrs/flock"
	"github.com/rs/zerolog"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/worker"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
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
		logger := logs.GetLogger()
		logger.Info().Msg("Starting worker")
		envConfig, err := config.GetEnvConfig()
		if err != nil {
			logger.Err(err).Msg("Failed to get environment config")
			return err
		}
		fileLock := flock.New(envConfig.ODIN_WORKER_INFO_FILE)
		ctx, cancel := context.WithCancel(cmd.Context())
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			<-sigChan
			logger.Info().Msg("Shutting down worker")
			logger.Info().Msg("Removing lock")
			fileLock.Unlock()
			cancel()
		}()
		if newWorker {
			deleteWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE)
		}
		_, queries, err := db.GetDBConnection(ctx, false, envConfig, false, nil, nil, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to get database connection")
			return err
		}
		prvdr, err := provider.GetProvider(ctx, queries, envConfig, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to get provider")
			return err
		}

		var wrkr *worker.Worker
		workerInfo, err := readWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to read worker info")
			switch err.(type) {
			case *WorkerInfoNotFoundError:
				logger.Info().Msgf("Creating new worker")
				name, err := cmd.Flags().GetString("name")
				if err != nil {
					logger.Err(err).Msg("Failed to get name flag")
				}
				if name == "" {
					name = namesgenerator.GetRandomName(0)
				}
				wrkr, err = worker.GetWorker(ctx, name, queries, envConfig, prvdr, logger)
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
			logger.Info().Msgf("Found worker info")
			wrkr, err = worker.GetWorker(ctx, workerInfo.Name, queries, envConfig, prvdr, logger)
			if err != nil {
				logger.Err(err).Msg("Failed to get worker")
				return err
			}
		}
		logger.Info().Msgf("Starting worker %d", wrkr.ID)

		err = writeWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE, wrkr)
		if err != nil {
			logger.Err(err).Msg("Failed to write worker info")
			return err
		}

		locked, err := fileLock.TryLock()
		if err != nil {
			logger.Err(err).Msg("Failed to lock worker info file")
			return err
		}
		if locked {
			err = wrkr.Run(ctx)
			if err != nil {
				logger.Err(err).Msg("Failed to start worker")
				err := deleteWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE)
				if err != nil {
					logger.Err(err).Msg("Failed to delete worker info")
					return err
				}
				return err
			}
		} else {
			logger.Info().Msg("Could not lock worker info file, worker is already running")
		}
		return nil
	},
}

func writeWorkerInfo(infoFile string, worker *worker.Worker) error {
	wrkrInfo := models.WorkerInfo{
		ID:   worker.ID,
		Name: worker.Name,
	}
	workerInfoBytes, err := json.Marshal(wrkrInfo)
	if err != nil {
		return err
	}
	f, err := os.Create(infoFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(workerInfoBytes)
	if err != nil {
		return err
	}
	return nil
}

func deleteWorkerInfo(infoFile string) error {
	err := os.Remove(infoFile)
	if err != nil {
		return err
	}
	return nil
}

func readWorkerInfo(infoFile string, logger *zerolog.Logger) (*models.WorkerInfo, error) {
	if _, err := os.Stat(infoFile); err != nil {
		if os.IsNotExist(err) {
			logger.Info().Msgf("Worker info not found at %s", infoFile)
			return nil, &WorkerInfoNotFoundError{}
		}
		if os.IsPermission(err) {
			logger.Err(err)
		}
		return nil, err
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
	return &wrkrInfo, nil
}

func init() {
	WorkerStartCmd.Flags().String("name", "", "Name of the worker")
	WorkerStartCmd.Flags().BoolVarP(&newWorker, "new", "n", false, "Create new worker(Deletes existing worker info)")
}
