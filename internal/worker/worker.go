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
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/executor"
	"github.com/deepakdinesh1123/valkyrie/internal/sandbox"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/gofrs/flock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
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
		DiskUsed float64
	}
}

type WorkerInfo struct {
	ID   int
	Name string
}

func GetWorker(ctx context.Context, name string, envConfig *config.EnvConfig, newWorker bool, standalone bool, logger *zerolog.Logger) (*Worker, error) {
	if newWorker {
		deleteWorkerInfo(envConfig.WORKER_INFO_FILE)
	}

	otelShutdown, tp, mp, _, err := telemetry.SetupOTelSDK(ctx, "Valkyrie Worker", envConfig)
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
	workerInfo, err := readWorkerInfo(envConfig.WORKER_INFO_FILE, logger)
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

	if envConfig.ENABLE_EXECUTION {
		exectr, err := executor.GetExecutor(ctx, queries, int32(wrkr.ID), tp, mp, envConfig, logger)
		if err != nil {
			return nil, err
		}
		wrkr.exectr = exectr
	}

	if envConfig.ENABLE_SANDBOX {
		sandboxHandler, err := sandbox.GetSandboxHandler(ctx, queries, int32(wrkr.ID), tp, mp, envConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("could not get sandbox handler: %s", err)
		}
		err = sandboxHandler.StartSandboxPool(ctx, envConfig)
		if err != nil {
			return nil, fmt.Errorf("error starting container pool: %v", err)
		}
		wrkr.sandboxHandler = sandboxHandler
	}
	logger.Info().Msgf("Starting worker %d with name %s", wrkr.ID, wrkr.Name)

	err = writeWorkerInfo(envConfig.WORKER_INFO_FILE, wrkr)
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
		defer w.exectr.Cleanup(ctx)
	}

	defer func() {
		var err error
		w.logger.Info().Msg("Shutting down opentelemetry")
		err = errors.Join(err, w.otelShutdown(context.Background()))
		if err != nil {
			w.logger.Err(err).Msg("Failed to shutdown OpenTelemetry")
		}
	}()

	tracer := w.tp.Tracer("valkyrie worker")

	// Create a root span for the entire worker run
	workerCtx, workerSpan := tracer.Start(ctx, "worker_run")
	defer workerSpan.End()

	// Acquire lock
	infLock := flock.New(w.envConfig.WORKER_INFO_FILE)
	locked, err := infLock.TryLock()
	if err != nil {
		w.logger.Err(err).Msg("Failed to acquire lock on worker info")
		workerSpan.RecordError(err)
		return err
	}
	if !locked {
		w.logger.Info().Msg("Worker: failed to acquire lock on worker info, another worker is running")
		lockErr := &WorkerError{Type: "Lock", Message: "Failed to acquire lock on worker info"}
		workerSpan.RecordError(lockErr)
		return lockErr
	}
	defer infLock.Unlock()

	var swg concurrency.SafeWaitGroup
	fetchJobTicker := time.NewTicker(time.Duration(w.envConfig.WORKER_POLL_FREQ) * time.Millisecond)
	heartBeatTicker := time.NewTicker(time.Duration(5) * time.Second)
	defer fetchJobTicker.Stop()
	defer heartBeatTicker.Stop()

	for {
		select {
		case <-workerCtx.Done():
			swg.Wait()
			err := workerCtx.Err()

			if w.sandboxHandler != nil {
				cleanupErr := w.sandboxHandler.Cleanup(context.TODO())
				if cleanupErr != nil {
					workerSpan.RecordError(cleanupErr)
					return fmt.Errorf("error cleaning up containers: %s", cleanupErr)
				}
			}
			w.queries.RequeueWorkerJobs(context.TODO(), pgtype.Int4{Valid: true, Int32: int32(w.ID)})

			switch err {
			case context.Canceled:
				w.logger.Info().Msg("Worker: context canceled")
				return nil
			default:
				w.logger.Err(err).Msg("Worker: context error")
				workerSpan.RecordError(err)
				return fmt.Errorf("context error: %s", err)
			}

		case <-fetchJobTicker.C:
			w.updateStats()

			// Check resource limits
			if w.WorkerStats.CPUUsage > w.envConfig.CPU_LIMIT {
				w.logger.Info().Float64("high CPU Usage", w.WorkerStats.CPUUsage).Msg("Worker: ")
				continue
			}
			if w.WorkerStats.MemUsage > w.envConfig.MEMORY_LIMIT {
				w.logger.Info().Float64("high memory usage", w.WorkerStats.MemUsage).Msg("Worker: ")
				continue
			}

			if w.WorkerStats.DiskUsed > w.envConfig.DISK_LIMIT {
				w.logger.Info().Float64("high memory usage", w.WorkerStats.MemUsage).Msg("Worker: pruning unsued images")
				if w.envConfig.ENABLE_EXECUTION {
					go w.exectr.PruneImages(ctx)
				}
			}

			// Handle execution jobs
			if w.envConfig.ENABLE_EXECUTION {
				if swg.Count() >= w.envConfig.WORKER_CONCURRENCY {
					w.logger.Info().Int("Tasks in progress", int(swg.Count())).Int32("Concurrency limit", w.envConfig.WORKER_CONCURRENCY).Msg("Worker: concurrency limit reached")
					continue
				}

				// Create a new span for each fetch job operation
				fetchCtx, fetchSpan := tracer.Start(workerCtx, "fetch_execution_job")

				res, err := w.queries.FetchJob(fetchCtx, db.FetchJobParams{
					Workerid: int32(w.ID),
					Jobtype:  "execution",
				})

				if err != nil {
					switch err {
					case pgx.ErrNoRows:
						// No jobs available, this is normal
						fetchSpan.End()
					case context.Canceled:
						w.logger.Info().Msg("Worker: context canceled")
						fetchSpan.RecordError(err)
						fetchSpan.End()
						swg.Wait()
						w.cleanup()
						return nil
					default:
						w.logger.Err(err).Msgf("Worker: failed to fetch job")
						fetchSpan.RecordError(err)
						fetchSpan.End()
						w.cleanup()
						return fmt.Errorf("failed to fetch job: %s", err)
					}
				} else {
					// Job found, execute it
					w.logger.Info().Msgf("Worker: fetched job %d", res.JobID)
					fetchSpan.SetAttributes(
						attribute.Int64("job.id", res.JobID),
						attribute.String("job.type", "execution"),
					)
					fetchSpan.End()

					// Create execution context with proper span propagation
					execCtx, execSpan := tracer.Start(workerCtx, "execute_job")
					execSpan.SetAttributes(attribute.Int64("job.id", res.JobID))

					swg.Add(1)
					go func(ctx context.Context, span trace.Span) {
						defer span.End()
						w.exectr.Execute(ctx, &swg, &res, w.logger.With().Int64("JOB_ID", res.JobID).Logger())
					}(execCtx, execSpan)
				}
			}

			// Handle sandbox jobs
			if w.envConfig.ENABLE_SANDBOX {
				sandboxCtx, sandboxSpan := tracer.Start(workerCtx, "fetch_sandbox_job")

				res, err := w.queries.FetchSandboxJobTx(sandboxCtx, db.FetchSandboxJobTxParams{WorkerID: int32(w.ID)})

				if err != nil {
					switch err {
					case pgx.ErrNoRows:
						// No sandbox jobs available, this is normal
						sandboxSpan.End()
					case context.Canceled:
						w.logger.Info().Msg("Worker: context canceled")
						sandboxSpan.RecordError(err)
						sandboxSpan.End()
						w.cleanup()
						w.logger.Info().Msg("cleanup complete")
						return nil
					default:
						w.logger.Err(err).Msgf("Worker: failed to fetch sandbox job")
						sandboxSpan.RecordError(err)
						sandboxSpan.End()
						w.cleanup()
						return &WorkerError{Type: "FetchSandboxJob", Message: err.Error()}
					}
				} else {
					// Sandbox job found
					w.logger.Info().Msgf("Worker: fetched sandbox job %d", res.Sandbox.SandboxID)
					sandboxSpan.SetAttributes(
						attribute.Int64("sandbox.id", res.Sandbox.SandboxID),
						attribute.String("job.type", "sandbox"),
					)
					sandboxSpan.End()

					swg.Add(1)
					go w.sandboxHandler.Create(workerCtx, &swg, res)
				}
			}

		case <-heartBeatTicker.C:
			// Create a span for heartbeat updates
			hbCtx, hbSpan := tracer.Start(workerCtx, "heartbeat_update")
			w.queries.UpdateHeartbeat(hbCtx, int32(w.ID))
			hbSpan.End()
		}
	}
}

func (w *Worker) cleanup() error {
	if w.envConfig.ENABLE_SANDBOX {
		if w.sandboxHandler != nil {
			w.sandboxHandler.Cleanup(context.TODO())
			err := w.queries.ClearSandboxes(context.TODO())
			if err != nil {
				return fmt.Errorf("error clearing sandboxes %v", err)
			}
		}
	}

	if w.envConfig.ENABLE_EXECUTION {
		w.exectr.Cleanup(context.Background())
	}
	return nil
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
