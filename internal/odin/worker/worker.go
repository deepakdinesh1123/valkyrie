package worker

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/gofrs/flock"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Worker struct {
	ID           int
	Name         string
	queries      db.Store
	envConfig    *config.EnvConfig
	exectr       executor.Executor
	logger       *zerolog.Logger
	tp           trace.TracerProvider
	mp           metric.MeterProvider
	otelShutdown func(context.Context) error

	WorkerStats struct {
		CPUUsage  float64
		MemAvail  uint64
		MemTotal  uint64
		MemUsed   uint64
		Timestamp time.Time
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
	exectr, err := executor.GetExecutor(ctx, queries, int32(wrkr.ID), tp, mp, envConfig, logger)
	if err != nil {
		return nil, err
	}
	wrkr.exectr = exectr
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
	defer func() {
		var err error

		w.logger.Info().Msg("Shutting down opentelemetry")
		err = errors.Join(err, w.otelShutdown(context.Background()))
		if err != nil {
			w.logger.Err(err).Msg("Failed to shutdown OpenTelemetry")
		}
	}()

	tracer := w.tp.Tracer("worker")
	tracerCtx, span := tracer.Start(ctx, "Run")
	defer span.End()

	span.AddEvent("Acquiring lock on worker info")
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
	fetchJobTicker := time.NewTicker(time.Duration(w.envConfig.ODIN_WORKER_POLL_FREQ) * time.Second)
	heartBeatTicker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Int32("Tasks in progress", swg.Count()).Msg("Worker: context done")
			swg.Wait()
			err := ctx.Err()
			fetchJobTicker.Stop()
			w.exectr.Cleanup()
			switch err {
			case context.Canceled:
				w.logger.Info().Msg("Worker: context canceled")
				return nil
			default:
				w.logger.Err(err).Msg("Worker: context error")
				return &WorkerError{Type: "Context", Message: err.Error()}
			}
		case <-fetchJobTicker.C:
			w.updateStats()
			// if w.WorkerStats.CPUUsage > 75 {
			// 	w.logger.Info().Float64("CPU Usage", w.WorkerStats.CPUUsage).Msg("Worker: high usage")
			// 	continue
			// }
			if swg.Count() >= w.envConfig.ODIN_WORKER_CONCURRENCY {
				w.logger.Info().Int("Tasks in progress", int(swg.Count())).Int32("Concurrency limit", w.envConfig.ODIN_WORKER_CONCURRENCY).Msg("Worker: concurrency limit reached")
				continue
			}
			res, err := w.queries.FetchJobTx(ctx, db.FetchJobTxParams{WorkerID: int32(w.ID)})
			if err != nil {
				switch err {
				case pgx.ErrNoRows:
					continue
				case context.Canceled:
					w.logger.Info().Msg("Worker: context canceled")
					w.exectr.Cleanup()
					return nil
				default:
					w.logger.Err(err).Msgf("Worker: failed to fetch job")
					return &WorkerError{Type: "FetchJob", Message: err.Error()}
				}
			}
			w.logger.Info().Msgf("Worker: fetched job %d", res.Job.JobID)
			swg.Add(1)
			span.AddEvent("Executing job")
			go w.exectr.Execute(tracerCtx, &swg, res.Job, w.logger.With().Int64("JOB_ID", res.Job.JobID).Logger())
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
