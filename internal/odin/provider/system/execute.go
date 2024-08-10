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
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *SystemProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Job) {
	start := time.Now()
	defer wg.Done()
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, execReq.InsertedAt.Time.Format("20060102150405"))
	s.logger.Info().Str("dir", dir).Msg("Executing job")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Err(err).Msg("Failed to create directory")
		err := s.updateJob(ctx, &execReq, start, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	if err := s.writeFiles(dir, execReq); err != nil {
		s.logger.Err(err).Msg("Failed to write files")
		err := s.updateJob(ctx, &execReq, start, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	outFile, err := os.Create(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create output file")
		err := s.updateJob(ctx, &execReq, start, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	defer outFile.Close()
	errFile, err := os.Create(filepath.Join(dir, "error.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create error file")
		err := s.updateJob(ctx, &execReq, start, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	defer errFile.Close()

	tctx, cancel := context.WithTimeout(context.TODO(), time.Duration(s.envConfig.ODIN_WORKER_TASK_TIMEOUT)*time.Second)
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
		if err := execCmd.Run(); err != nil {
			if tctx.Err() != nil {
				switch tctx.Err() {
				case context.DeadlineExceeded:
					done <- true
					return
				}
			}
			s.logger.Err(err).Msg("Failed to execute command")
			err := s.updateJob(ctx, &execReq, start, err.Error())
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
				s.updateExecutionDetails(context.TODO(), dir, start, execReq)
				return
			default:
				s.logger.Info().Msg("context error killing process")
				err := execCmd.Process.Kill()
				if err != nil {
					s.logger.Err(err).Msg("Failed to send kill signal")
				}
				s.updateExecutionDetails(context.TODO(), dir, start, execReq)
				return
			}
		case <-done:
			s.updateExecutionDetails(context.TODO(), dir, start, execReq)
			return
		}
	}
}

func (s *SystemProvider) updateExecutionDetails(ctx context.Context, dir string, startTime time.Time, execReq db.Job) {
	out, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to read output file")
		return
	}
	err = s.updateJob(ctx, &execReq, startTime, string(out))
	if err != nil {
		s.logger.Err(err).Msg("Failed to update job")
	}
}

func (s *SystemProvider) writeFiles(dir string, execReq db.Job) error {
	files := map[string]string{
		"flake.nix":        execReq.Flake,
		execReq.ScriptPath: execReq.Script,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (d *SystemProvider) updateJob(ctx context.Context, execReq *db.Job, startTime time.Time, message string) error {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := d.queries.UpdateJob(ctx, execReq.ID); err != nil {
		return err
	}
	if _, err := d.queries.InsertJobRun(ctx, db.InsertJobRunParams{
		JobID:      execReq.ID,
		WorkerID:   execReq.WorkerID.Int32,
		StartedAt:  pgtype.Timestamptz{Time: startTime, Valid: true},
		FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Script:     execReq.Script,
		Flake:      execReq.Flake,
		Args:       execReq.Args,
		Logs:       message,
		CreatedAt:  execReq.InsertedAt,
	}); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
