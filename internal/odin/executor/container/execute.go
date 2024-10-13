package container

import (
	"context"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

func (ce *ContainerExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job) {
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
		tctx, cancel = context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(context.TODO())
	}
	defer cancel()

	cc, err := GetContainerClient(tctx, ce)
	if err != nil {
		ce.logger.Err(err).Msg("could not get container client")
		ce.checkFailed(ce.queries.UpdateJobResultTx(ctx, jobRes))
		return
	}
	cont, err := cc.GetContainer(tctx)
	if err != nil {
		ce.logger.Err(err).Msg("could not get container")
		ce.checkFailed(ce.queries.UpdateJobResultTx(ctx, jobRes))
		return
	}
	defer cont.Destroy()
	contInfo := cont.Value()
	err = cc.WriteFiles(tctx, contInfo.ID, contInfo.HostPrepDir, job)
	if err != nil {
		ce.logger.Err(err).Msg("could not write files")
		ce.checkFailed(ce.queries.UpdateJobResultTx(ctx, jobRes))
		return
	}
	success, output, err := cc.Execute(tctx, contInfo.ID, []string{"bash", "nix_run.sh"})
	if err != nil {
		ce.logger.Err(err).Msg(err.Error())
		ce.checkFailed(ce.queries.UpdateJobResultTx(ctx, jobRes))
		return
	}
	ce.logger.Info().Bool("Success", success).Msg(output)
	ce.logger.Info().Msg("Destroying container")

	jobRes.Success = true
	jobRes.Retry = false
	jobRes.ExecLogs = output
	ce.checkFailed(ce.queries.UpdateJobResultTx(ctx, jobRes))
}

func (ce *ContainerExecutor) checkFailed(res db.UpdateJobTxResult, err error) {
	if err != nil {
		ce.logger.Err(err).Msg("Failed to update job")
	}
}

func (ce *ContainerExecutor) Cleanup() {
	ce.logger.Debug().Msg("Cleaning up")
	ce.pool.Close()
}
