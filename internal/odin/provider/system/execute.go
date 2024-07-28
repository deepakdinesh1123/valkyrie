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

func (s *SystemProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Jobqueue) {
	defer wg.Done()
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, execReq.CreatedAt.Time.Format("20060102150405"))
	s.logger.Info().Str("dir", dir).Msg("Executing job")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Err(err).Msg("Failed to create directory")
		s.updateJob(ctx, execReq.ID, err.Error())
		return
	}
	if err := s.writeFiles(dir, execReq); err != nil {
		s.logger.Err(err).Msg("Failed to write files")
		s.updateJob(ctx, execReq.ID, err.Error())
		return
	}
	outFile, err := os.Create(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create output file")
		s.updateJob(ctx, execReq.ID, err.Error())
		return
	}
	defer outFile.Close()
	errFile, err := os.Create(filepath.Join(dir, "error.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create error file")
		s.updateJob(ctx, execReq.ID, err.Error())
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

	done := make(chan struct{})

	go func() {
		if err := execCmd.Run(); err != nil {
			if tctx.Err() != nil {
				switch tctx.Err() {
				case context.DeadlineExceeded:
					close(done)
					return
				}
			}
			s.logger.Err(err).Msg("Failed to execute command")
			s.updateJob(ctx, execReq.ID, err.Error())
			close(done)
			return
		}
		close(done)
	}()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				s.logger.Info().Msg("context canceled wating for process to exit")
				<-done
				s.updateExecutionDetails(context.TODO(), dir, execReq)
				return
			default:
				s.logger.Info().Msg("context error killing process")
				err := execCmd.Process.Kill()
				if err != nil {
					s.logger.Err(err).Msg("Failed to send kill signal")
				}
				s.updateExecutionDetails(context.TODO(), dir, execReq)
				return
			}
		case <-done:
			s.updateExecutionDetails(context.TODO(), dir, execReq)
			return
		}
	}
}

func (s *SystemProvider) updateExecutionDetails(ctx context.Context, dir string, execReq db.Jobqueue) {
	out, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to read output file")
		return
	}
	s.updateJob(ctx, execReq.ID, string(out))
}

func (s *SystemProvider) writeFiles(dir string, execReq db.Jobqueue) error {
	files := map[string]string{
		"flake.nix":               execReq.Flake.String,
		execReq.ScriptPath.String: execReq.Script.String,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *SystemProvider) updateJob(ctx context.Context, jobID int64, message string) {
	if _, err := s.queries.UpdateJob(ctx, db.UpdateJobParams{
		ID:   jobID,
		Logs: pgtype.Text{String: message, Valid: true},
	}); err != nil {
		s.logger.Error().Err(err).Msg("Failed to update job")
	}
}
