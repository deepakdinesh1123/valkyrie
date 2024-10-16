package system

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

func (s *SystemExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job, logger zerolog.Logger) {
	tracer := s.tp.Tracer("Execute")
	_, span := tracer.Start(ctx, "Execute")
	defer span.End()

	span.AddEvent("Executing job")

	start := time.Now()
	defer wg.Done()
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_EXECUTOR_BASE_DIR, job.CreatedAt.Time.Format("20060102150405"))
	s.logger.Info().Str("dir", dir).Msg("Executing job")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Err(err).Msg("Failed to create directory")
		err := s.updateJob(ctx, &job, start, false)
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	if err := s.writeFiles(ctx, dir, job); err != nil {
		s.logger.Err(err).Msg("Failed to write files")
		err := s.updateJob(ctx, &job, start, false)
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	outFile, err := os.Create(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create output file")
		err := s.updateJob(ctx, &job, start, false)
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	defer outFile.Close()
	errFile, err := os.Create(filepath.Join(dir, "error.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create error file")
		err := s.updateJob(ctx, &job, start, false)
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	defer errFile.Close()

	var timeout int
	if job.TimeOut.Int32 > 0 { // By default, timeout is set to -1
		timeout = int(job.TimeOut.Int32)
	} else if job.TimeOut.Int32 == 0 {
		timeout = 0
	} else {
		timeout = s.envConfig.ODIN_WORKER_TASK_TIMEOUT
	}

	var tctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		tctx, cancel = context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(context.TODO())
	}
	defer cancel()

	execCmd := exec.CommandContext(tctx, "nix", "run")
	execCmd.Cancel = func() error {
		s.logger.Info().Msg("Task timed out. Terminating execution")
		syscall.Kill(-execCmd.Process.Pid, syscall.SIGKILL)
		return nil
	}
	execCmd.Dir = dir
	execCmd.Stdout = outFile
	execCmd.Stderr = errFile
	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	done := make(chan bool, 1)

	go func() {
		s.logger.Info().Msg("Executing nix run command")
		if err := execCmd.Run(); err != nil {
			if tctx.Err() != nil {
				switch tctx.Err() {
				case context.DeadlineExceeded:
					done <- true
					return
				}
			}
			s.logger.Err(err).Msg("Failed to execute command")
			err := s.updateJob(ctx, &job, start, false)
			if err != nil {
				s.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		done <- true
	}()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				s.logger.Info().Msg("context canceled wating for process to exit")
				<-done
				s.updateJob(context.TODO(), &job, start, true)
				return
			default:
				s.logger.Info().Msg("context error killing process")
				err := execCmd.Process.Kill()
				if err != nil {
					s.logger.Err(err).Msg("Failed to send kill signal")
				}
				s.updateJob(context.TODO(), &job, start, false)
				return
			}
		case <-done:
			s.updateJob(context.TODO(), &job, start, true)
			return
		}
	}
}

func (s *SystemExecutor) writeFiles(ctx context.Context, dir string, job db.Job) error {
	execReq, err := s.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		return err
	}
	files := map[string]string{
		"flake.nix":  execReq.Flake,
		execReq.Path: execReq.Code,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *SystemExecutor) updateJob(ctx context.Context, job *db.Job, startTime time.Time, success bool) error {
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_EXECUTOR_BASE_DIR, job.CreatedAt.Time.Format("20060102150405"))
	out, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		return err
	}
	stderr, err := os.ReadFile(filepath.Join(dir, "error.txt"))
	if err != nil {
		return err
	}
	retry := true
	if job.Retries.Int32+1 >= job.MaxRetries.Int32 || success {
		retry = false
	}
	if _, err := s.queries.UpdateJobResultTx(ctx, db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       *job,
		Success:   success,
		WorkerId:  s.workerId,
		Retry:     retry,
		ExecLogs:  string(out),
		NixLogs:   string(stderr),
	}); err != nil {
		return err
	}
	return nil
}

func (s *SystemExecutor) Cleanup() {}
