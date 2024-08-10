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

// Execute executes a job using the system provider.
//
// Parameters:
// - ctx: the context for the execution operation
// - wg: the wait group for the execution operation
// - execReq: the job execution request
func (s *SystemProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Job) {
	defer wg.Done()
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, execReq.InsertedAt.Time.Format("20060102150405"))
	s.logger.Info().Str("dir", dir).Msg("Executing job")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Err(err).Msg("Failed to create directory")
		err := s.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	if err := s.writeFiles(dir, execReq); err != nil {
		s.logger.Err(err).Msg("Failed to write files")
		err := s.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	outFile, err := os.Create(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create output file")
		err := s.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			s.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	defer outFile.Close()
	errFile, err := os.Create(filepath.Join(dir, "error.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to create error file")
		err := s.updateJob(ctx, &execReq, err.Error())
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
		syscall.Kill(-execCmd.Process.Pid, syscall.SIGKILL) // send SIGKILL to child process
		return nil
	}
	execCmd.Dir = dir
	execCmd.Stdout = outFile
	execCmd.Stderr = errFile
	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // prevent child process from receiving termination signal
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
			err := s.updateJob(ctx, &execReq, err.Error())
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

// updateExecutionDetails updates the execution details of a job by reading the output file and updating the job status.
//
// Parameters:
// - ctx: the context for the update operation
// - dir: the directory containing the output file
// - execReq: the job execution request
func (s *SystemProvider) updateExecutionDetails(ctx context.Context, dir string, execReq db.Job) {
	out, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to read output file")
		return
	}
	err = s.updateJob(ctx, &execReq, string(out))
	if err != nil {
		s.logger.Err(err).Msg("Failed to update job")
	}
}

// writeFiles writes the flake and script files to the odin execution directory.
//
// Parameters:
// - dir: the directory to write the files to
// - execReq: the job execution request containing the flake and script and other metadata
// Returns:
// - error: any error that occurs during file writing
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

// updateJob updates the job status to completed and inserts a new job run.
//
// Parameters:
// - ctx: the context for the update operation.
// - execReq: the job execution request.
// - message: the message to be logged.
// Returns:
// - error: an error if the update operation fails.
func (d *SystemProvider) updateJob(ctx context.Context, execReq *db.Job, message string) error {
	if err := d.queries.UpdateJob(ctx, execReq.ID); err != nil {
		return err
	}
	if _, err := d.queries.InsertJobRun(ctx, db.InsertJobRunParams{
		JobID:      execReq.ID,
		WorkerID:   execReq.WorkerID.Int32,
		StartedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Script:     message,
		Flake:      execReq.Flake,
		Args:       execReq.Args,
		Logs:       pgtype.Text{String: message, Valid: true},
		CreatedAt:  execReq.InsertedAt,
	}); err != nil {
		return err
	}
	return nil
}
