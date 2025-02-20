package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/sandbox"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/gofrs/flock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Worker struct {
	ID             int
	Name           string
	queries        db.Store
	envConfig      *config.EnvConfig
	exectr         executor.Executor
	sandboxHandler sandbox.SandboxHandler
	logger         *zerolog.Logger
	tp             trace.TracerProvider
	mp             metric.MeterProvider
	otelShutdown   func(context.Context) error

	WorkerStats struct {
		CPUUsage float64
		MemUsage float64
	}
}

type WorkerInfo struct {
	ID   int
	Name string
}

func GetWorker(ctx context.Context, name string, envConfig *config.EnvConfig, newWorker bool, standalone bool, logger *zerolog.Logger) (*Worker, error) {
	if newWorker {
		deleteWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE)
	}

	otelShutdown, tp, mp, _, err := telemetry.SetupOTelSDK(ctx, "Odin Worker", envConfig)
	if err != nil {
		logger.Err(err).Msg("Failed to setup OpenTelemetry")
		return nil, err
	}

	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(false),
		db.IsStandalone(standalone),
		db.IsWorker(true),
		db.WithTracerProvider(tp),
	)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		logger.Err(err).Msg("Failed to get database connection")
		return nil, err
	}

	wrkr := &Worker{
		queries:      queries,
		envConfig:    envConfig,
		logger:       logger,
		tp:           tp, // trace provider
		mp:           mp, // metric provider
		otelShutdown: otelShutdown,
	}
	workerInfo, err := readWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE, logger)
	if err != nil {
		switch err.(type) {
		case *WorkerInfoNotFoundError:
			if name == "" {
				name = namesgenerator.GetRandomName(0)
			}
			wrkr.Name = name
			wrkr.ID, err = wrkr.upsertWorker(ctx, name, -1)
			if err != nil {
				logger.Err(err).Msg("Failed to create worker")
			}
		default:
			logger.Err(err).Msg("Failed to read worker info")
		}
	}
	if wrkr.ID == 0 && workerInfo != nil {
		logger.Info().Str("workerName", workerInfo.Name).Int("workerID", workerInfo.ID).Msgf("Found worker info")
		wrkr.Name = workerInfo.Name
		wrkr.ID, err = wrkr.upsertWorker(ctx, workerInfo.Name, workerInfo.ID)
		if err != nil {
			logger.Err(err).Msg("Failed to get worker")
		}
	}

	if envConfig.ODIN_ENABLE_EXECUTION {
		exectr, err := executor.GetExecutor(ctx, queries, int32(wrkr.ID), tp, mp, envConfig, logger)
		if err != nil {
			return nil, err
		}
		wrkr.exectr = exectr
	}

	if envConfig.ODIN_ENABLE_SANDBOX {
		sandboxHandler, err := sandbox.GetSandboxHandler(ctx, queries, int32(wrkr.ID), tp, mp, envConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("could not get sandbox handler: %s", err)
		}
		if !envConfig.ODIN_COMPOSE_ENV {
			err = sandboxHandler.StartOdinStore(ctx, envConfig.ODIN_STORE_IMAGE, envConfig.ODIN_STORE_CONTAINER, envConfig.ODIN_CONTAINER_RUNTIME)
			if err != nil {
				return nil, fmt.Errorf("could not start odin store: %v", err)
			}
		}
		err = sandboxHandler.StartContainerPool(ctx, envConfig)
		if err != nil {
			return nil, fmt.Errorf("error starting container pool: %v", err)
		}
		wrkr.sandboxHandler = sandboxHandler
	}
	logger.Info().Msgf("Starting worker %d with name %s", wrkr.ID, wrkr.Name)

	err = writeWorkerInfo(envConfig.ODIN_WORKER_INFO_FILE, wrkr)
	if err != nil {
		return nil, err
	}
	return wrkr, nil
}

func (w *Worker) upsertWorker(ctx context.Context, name string, id int) (int, error) {
	wrkr, err := w.queries.GetWorker(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			if id == -1 {
				wrkr, err = w.queries.CreateWorker(ctx, name)
				if err != nil {
					w.logger.Err(err).Msg("Worker: failed to create worker")
					return 0, err
				}
			} else {
				wrkr, err = w.queries.InsertWorker(ctx, db.InsertWorkerParams{
					ID:   int32(id),
					Name: name,
				})
				if err != nil {
					w.logger.Err(err).Msg("Worker: failed to insert worker")
					return 0, err
				}
			}
		} else {
			w.logger.Err(err).Msg("Worker: failed to get worker")
			return 0, err
		}
	}
	return int(wrkr.ID), nil
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) error {
	w.queries.UpdateHeartbeat(ctx, int32(w.ID))
	defer wg.Done()
	if w.exectr != nil {
		defer w.exectr.Cleanup()
	}

	defer func() {
		var err error

		w.logger.Info().Msg("Shutting down opentelemetry")
		err = errors.Join(err, w.otelShutdown(context.Background()))
		if err != nil {
			w.logger.Err(err).Msg("Failed to shutdown OpenTelemetry")
		}
	}()

	// tracer := w.tp.Tracer("worker")
	// tracerCtx, span := tracer.Start(ctx, "Run")
	// defer span.End()

	// span.AddEvent("Acquiring lock on worker info")
	infLock := flock.New(w.envConfig.ODIN_WORKER_INFO_FILE)
	locked, err := infLock.TryLock()
	if err != nil {
		w.logger.Err(err).Msg("Failed to acquire lock on worker info")
		return err
	}
	if !locked {
		w.logger.Info().Msg("Worker: failed to acquire lock on worker info, another worker is running")
		return &WorkerError{Type: "Lock", Message: "Failed to acquire lock on worker info"}
	}
	defer infLock.Unlock()
	var swg concurrency.SafeWaitGroup
	fetchJobTicker := time.NewTicker(time.Duration(w.envConfig.ODIN_WORKER_POLL_FREQ) * time.Millisecond)
	heartBeatTicker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-ctx.Done():
			// w.logger.Info().Int32("Tasks in progress", swg.Count()).Msg("Worker: context done")
			swg.Wait()
			err := ctx.Err()
			fetchJobTicker.Stop()
			if w.sandboxHandler != nil {
				err = w.sandboxHandler.Cleanup(context.TODO())
				if err != nil {
					return fmt.Errorf("error cleaning up containers: %s", err)
				}
			}
			w.queries.RequeueWorkerJobs(context.TODO(), pgtype.Int4{Valid: true, Int32: int32(w.ID)})
			switch err {
			case context.Canceled:
				w.logger.Info().Msg("Worker: context canceled")
				return nil
			default:
				w.logger.Err(err).Msg("Worker: context error")
				return fmt.Errorf("context error: %s", err)
			}
		case <-fetchJobTicker.C:
			w.updateStats()
			if w.WorkerStats.CPUUsage > w.envConfig.ODIN_CPU_LIMIT {
				w.logger.Info().Float64("high CPU Usage", w.WorkerStats.CPUUsage).Msg("Worker: ")
				continue
			}
			if w.WorkerStats.MemUsage > w.envConfig.ODIN_MEMORY_LIMIT {
				w.logger.Info().Float64("high memory usage", w.WorkerStats.MemUsage).Msg("Worker: ")
				continue
			}
			if w.envConfig.ODIN_ENABLE_EXECUTION {
				if swg.Count() >= w.envConfig.ODIN_WORKER_CONCURRENCY {
					w.logger.Info().Int("Tasks in progress", int(swg.Count())).Int32("Concurrency limit", w.envConfig.ODIN_WORKER_CONCURRENCY).Msg("Worker: concurrency limit reached")
					continue
				}
				res, err := w.queries.FetchJob(ctx, db.FetchJobParams{
					Workerid: int32(w.ID),
					Jobtype:  "execution",
				})
				if err != nil {
					switch err {
					case pgx.ErrNoRows:
						continue
					case context.Canceled:
						w.logger.Info().Msg("Worker: context canceled")
						swg.Wait()
						return nil
					default:
						w.logger.Err(err).Msgf("Worker: failed to fetch job")
						return fmt.Errorf("failed to fetch job: %s", err)
					}
				}
				w.logger.Info().Msgf("Worker: fetched job %d", res.JobID)
				swg.Add(1)
				// span.AddEvent("Executing job")
				go w.exectr.Execute(ctx, &swg, &res, w.logger.With().Int64("JOB_ID", res.JobID).Logger())
			}

			if w.envConfig.ODIN_ENABLE_SANDBOX {
				res, err := w.queries.FetchSandboxJobTx(ctx, db.FetchSandboxJobTxParams{WorkerID: int32(w.ID)})
				if err != nil {
					switch err {
					case pgx.ErrNoRows:
						continue
					case context.Canceled:
						w.logger.Info().Msg("Worker: context canceled")
						w.sandboxHandler.Cleanup(context.TODO())
						err = w.queries.ClearSandboxes(context.TODO())
						if err != nil {
							w.logger.Err(err).Msg("error clearing sandboxes")
						}
						w.logger.Info().Msg("cleanup complete")
						return nil
					default:
						w.logger.Err(err).Msgf("Worker: failed to fetch sandbox job")
						return &WorkerError{Type: "FetchSandboxJob", Message: err.Error()}
					}
				}
				w.logger.Info().Msgf("Worker: fetched sandbox job %d", res.Sandbox.SandboxID)
				swg.Add(1)
				go w.sandboxHandler.Create(ctx, &swg, res)
			}
		case <-heartBeatTicker.C:
			w.queries.UpdateHeartbeat(ctx, int32(w.ID))
		}
	}
}

func writeWorkerInfo(infoFile string, worker *Worker) error {
	wrkrInfo := WorkerInfo{
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

func readWorkerInfo(infoFile string, logger *zerolog.Logger) (*WorkerInfo, error) {
	if _, err := os.Stat(infoFile); err != nil {
		if os.IsNotExist(err) {
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
	var wrkrInfo WorkerInfo
	err = json.Unmarshal(workerInfoBytes, &wrkrInfo)
	if err != nil {
		return nil, err
	}
	return &wrkrInfo, nil
}
