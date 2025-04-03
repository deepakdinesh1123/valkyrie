package container

import (
	"context"
	"os"
	"os/user"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type ContainerExecutor struct {
	Queries   db.Store
	EnvConfig *config.EnvConfig
	WorkerId  int32
	Logger    *zerolog.Logger
	Tp        trace.TracerProvider
	Mp        metric.MeterProvider
	User      string
	Pool      *puddle.Pool[pool.Container]
}

func NewContainerExecutor(ctx context.Context, env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*ContainerExecutor, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	pool, err := pool.NewContainerExecutionPool(ctx, int32(env.HOT_CONTAINER), env.WORKER_CONCURRENCY, env.RUNTIME)
	if err != nil {
		return nil, err
	}
	return &ContainerExecutor{
		EnvConfig: env,
		Logger:    logger,
		Queries:   queries,
		WorkerId:  workerId,
		Tp:        tp,
		Mp:        mp,
		User:      user.Username,
		Pool:      pool,
	}, nil
}

func (ce *ContainerExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job *db.Job, logger zerolog.Logger) {
	defer wg.Done()
	startTime := time.Now()

	jobRes := db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       job,
		Retry:     true,
		Success:   false,
		WorkerId:  ce.WorkerId,
	}
	if job.Retries.Int32+1 >= job.MaxRetries.Int32 {
		jobRes.Retry = false
	}

	var timeout int
	if job.TimeOut.Int32 > 0 { // By default, timeout is set to -1
		timeout = int(job.TimeOut.Int32)
	} else if job.TimeOut.Int32 == 0 {
		timeout = 0
	} else {
		timeout = ce.EnvConfig.WORKER_TASK_TIMEOUT
	}
	var tctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		tctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	logger.Debug().Msg("Getting container client")
	cc, err := GetContainerClient(tctx, ce)
	if err != nil {
		logger.Err(err).Msg("could not get container client")
		ce.checkFailed(ce.Queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Got container client")
	logger.Debug().Msg("Getting container")
	cont, err := cc.GetContainer(tctx)
	if err != nil {
		logger.Err(err).Msg("could not get container")
		ce.checkFailed(ce.Queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Got container")
	defer cont.Destroy()
	contInfo := cont.Value()
	logger.Debug().Msg("Writing files")
	err = cc.WriteFiles(tctx, contInfo.ID, os.TempDir(), job)
	if err != nil {
		logger.Err(err).Msg("could not write files")
		ce.checkFailed(ce.Queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Files written")
	success, output, err := cc.Execute(tctx, contInfo.ID, []string{"sh", "nix_run.sh"})
	if err != nil {
		logger.Err(err).Msg(err.Error())
		ce.checkFailed(ce.Queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Destroying container")

	jobRes.Success = success
	if success {
		jobRes.Retry = false
	}
	jobRes.ExecLogs = output
	ce.Logger.Debug().Str("output", output).Msg("Exec Logs")
	ce.checkFailed(ce.Queries.UpdateJobResultTx(context.TODO(), jobRes))
}

func (ce *ContainerExecutor) checkFailed(_ db.UpdateJobTxResult, err error) {
	if err != nil {
		ce.Logger.Error().Err(err).Stack().Msgf("An error occurred %s: ", err)
	}
}

func (ce *ContainerExecutor) Cleanup() {
	ce.Logger.Debug().Msg("Cleaning up")
	ce.Pool.Close()
}
