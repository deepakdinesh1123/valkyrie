package system

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *SystemProvider) Execute(ctx context.Context, wg *sync.WaitGroup, execReq db.Jobqueue, cancel context.CancelFunc) {
	defer wg.Done()
	defer cancel()
	dir := fmt.Sprintf("%s/%s", s.envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, execReq.CreatedAt.Time.Format("20060102150405"))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
		}
	}
	err := os.WriteFile(fmt.Sprintf("%s/%s", dir, "flake.nix"), []byte(execReq.Flake.String), os.ModePerm)
	if err != nil {
		s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s", dir, execReq.ScriptPath.String), []byte(execReq.Script.String), os.ModePerm)
	if err != nil {
		s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
	}
	execCmd := exec.Command("nix", "run")
	execCmd.Dir = dir
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
	}
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
	}
	err = execCmd.Start()
	if err != nil {
		s.queries.UpdateJob(ctx, db.UpdateJobParams{ID: execReq.ID, Concat: pgtype.Text{String: err.Error(), Valid: true}})
	}
	s.logger.Info().Msg("Task started")
	s.logger.Info().Msg(execCmd.String())

	done := make(chan error, 1)
	go func() {
		done <- execCmd.Wait()
		close(done)
	}()

	// Create a scanner to read from stdout
	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			s.queries.UpdateJob(ctx, db.UpdateJobParams{
				ID:     execReq.ID,
				Concat: pgtype.Text{String: string(stdoutScanner.Bytes()), Valid: true},
			})
		}
	}()

	// Create a scanner to read from stderr
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			s.queries.UpdateJob(ctx, db.UpdateJobParams{
				ID:     execReq.ID,
				Concat: pgtype.Text{String: string(stderrScanner.Bytes()), Valid: true},
			})
		}
	}()

	select {
	case <-ctx.Done():
		if execCmd.Process != nil {
			_ = execCmd.Process.Kill()
		}
		s.logger.Info().Msg("Task cancelled")
		s.queries.UpdateJob(ctx, db.UpdateJobParams{
			ID:     execReq.ID,
			Concat: pgtype.Text{String: "cancelled", Valid: true},
		})
	case <-done:
		s.logger.Info().Msg("Task completed")
		s.queries.UpdateJob(ctx, db.UpdateJobParams{
			ID:     execReq.ID,
			Concat: pgtype.Text{String: "completed", Valid: true},
		})
	}
}
