package container

import (
	"context"
	"os"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

func (ce *ContainerExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job, logger zerolog.Logger) {
	defer wg.Done()
	startTime := time.Now()

	jobRes := db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       job,
		Retry:     true,
		Success:   false,
		WorkerId:  ce.workerId,
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
		timeout = ce.envConfig.ODIN_WORKER_TASK_TIMEOUT
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
		ce.checkFailed(ce.queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Got container client")
	logger.Debug().Msg("Getting container")
	cont, err := cc.GetContainer(tctx)
	if err != nil {
		logger.Err(err).Msg("could not get container")
		ce.checkFailed(ce.queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Got container")
	// defer cont.Destroy()
	contInfo := cont.Value()
	logger.Debug().Msg("Writing files")
	err = cc.WriteFiles(tctx, contInfo.ID, os.TempDir(), job)
	if err != nil {
		logger.Err(err).Msg("could not write files")
		ce.checkFailed(ce.queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Files written")
	success, output, err := cc.Execute(tctx, contInfo.ID, []string{"sh", "nix_run.sh"})
	if err != nil {
		logger.Err(err).Msg(err.Error())
		ce.checkFailed(ce.queries.UpdateJobResultTx(context.TODO(), jobRes))
		return
	}
	logger.Debug().Msg("Destroying container")

	jobRes.Success = success
	if success {
		jobRes.Retry = false
	}
	jobRes.ExecLogs = output
	ce.logger.Debug().Str("output", output).Msg("Exec Logs")
	ce.checkFailed(ce.queries.UpdateJobResultTx(context.TODO(), jobRes))
}

func (ce *ContainerExecutor) checkFailed(_ db.UpdateJobTxResult, err error) {
	if err != nil {
		ce.logger.Error().Err(err).Stack().Msgf("An error occurred %s: ", err)
	}
}

func (ce *ContainerExecutor) Cleanup() {
	ce.logger.Debug().Msg("Cleaning up")
	ce.pool.Close()
}
