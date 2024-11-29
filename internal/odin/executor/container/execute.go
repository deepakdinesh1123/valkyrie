package container

import (
	"context"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

func (cp *ContainerExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job) {
	defer wg.Done()
	startTime := time.Now()

	var timeout int
	if job.TimeOut.Int32 > 0 { // By default, timeout is set to -1
		timeout = int(job.TimeOut.Int32)
	} else if job.TimeOut.Int32 == 0 {
		timeout = 0
	} else {
		timeout = cp.envConfig.ODIN_WORKER_TASK_TIMEOUT
	}
	var tctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		tctx, cancel = context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(context.TODO())
	}
	defer cancel()

	cc, err := GetContainerClient(tctx, cp)
	if err != nil {
		cp.logger.Err(err).Msg("could not container client")
	}
	cont, err := cc.GetContainer(tctx)
	if err != nil {
		cp.logger.Err(err).Msg("could not get container")
	}
	contInfo := cont.Value()
	err = cc.WriteFiles(tctx, contInfo.ID, contInfo.HostPrepDir, job)
	if err != nil {
		cp.logger.Err(err).Msg("could not write files")
		cont.Destroy()
	}
	success, output, err := cc.Execute(tctx, contInfo.ID, []string{"bash", "nix_run.sh"})
	if err != nil {
		cp.logger.Err(err).Msg("could not write files")
		cont.Destroy()
	}
	cp.logger.Info().Bool("Success", success).Msg(output)
	retry := true
	if job.Retries.Int32+1 >= job.MaxRetries.Int32 || success {
		retry = false
	}
	if _, err := cp.queries.UpdateJobResultTx(ctx, db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       job,
		Message:   "",
		Success:   success,
		Retry:     retry,
		WorkerId:  cp.workerId,
	}); err != nil {
		cp.logger.Err(err).Msg("could not update job")
		cont.Destroy()
	}
	cont.Destroy()
}

func (cp *ContainerExecutor) Cleanup() {
	cp.logger.Debug().Msg("Cleaning up")
	cp.pool.Close()
}
