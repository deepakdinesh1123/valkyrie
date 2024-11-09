//go:build system || all

package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type SystemExecutor struct {
	envConfig *config.EnvConfig
	queries   db.Store
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	workerId  int32
}

func NewSystemExecutor(ctx context.Context, envConfig *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*SystemExecutor, error) {
	return &SystemExecutor{
		envConfig: envConfig,
		queries:   queries,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}

type SystemExecutionClient interface {
	GetExecCmd(ctx context.Context, outFile *os.File, errFile *os.File, dir string, execReq *db.ExecRequest) (*exec.Cmd, error)
}

func (s *SystemExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job *db.Job, logger zerolog.Logger) {
	startTime := time.Now()
	defer wg.Done()
	dir := filepath.Join(s.envConfig.ODIN_SYSTEM_EXECUTOR_BASE_DIR, job.CreatedAt.Time.Format("20060102150405"))
	s.logger.Info().Str("dir", dir).Msg("Executing execReq")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		s.logger.Err(err)
		return
	}

	s.logger.Info().Msg("getting exec request")
	res, err := s.queries.GetTotalJobs(ctx)
	s.logger.Info().Int64("id", res).Msg("Log")
	execReq, err := s.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		s.logger.Err(err)
		return
	}

	s.logger.Info().Msg("writing files")
	if err := s.writeFiles(ctx, dir, &execReq); err != nil {
		s.logger.Err(err).Msg("Failed to update execReq")
		return
	}
	outFile, err := os.Create(filepath.Join(dir, "output.txt"))
	if err != nil {
		s.logger.Err(err)
		return
	}
	defer outFile.Close()
	// errFile, err := os.Create(filepath.Join(dir, "error.txt"))
	// if err != nil {
	// 	s.logger.Err(err)
	// 	return
	// }
	// defer errFile.Close()

	jobRes := db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       job,
		Retry:     true,
		Success:   false,
		WorkerId:  s.workerId,
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
		timeout = s.envConfig.ODIN_WORKER_TASK_TIMEOUT
	}

	var tctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		tctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	sc, err := getSystemExecutorClient(ctx, s)
	if err != nil {
		s.logger.Err(err).Msg("could not get client")
		return
	}

	done := make(chan bool, 1)

	execCmd, err := sc.GetExecCmd(tctx, outFile, outFile, dir, &execReq)
	if err != nil {
		s.logger.Err(err).Msg("could not get exec command")
		return
	}

	go func() {
		if err := execCmd.Run(); err != nil {
			if tctx.Err() != nil {
				switch tctx.Err() {
				case context.DeadlineExceeded:
					done <- true
					return
				}
			}
			s.checkFailed(s.queries.UpdateJobResultTx(ctx, jobRes))
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
				err := execCmd.Process.Kill()
				if err != nil {
					s.logger.Err(err).Msg("error killing process")
				}
				out, err := s.ReadOutput(dir)
				if err != nil {
					s.logger.Err(err).Msg("error reading output")
					return
				}
				jobRes.ExecLogs = string(out)
				s.logger.Debug().Bytes("output", out).Msg("Exec Result")
				s.checkFailed(s.queries.UpdateJobResultTx(ctx, jobRes))
				return
			default:
				err := execCmd.Process.Kill()
				if err != nil {
					s.logger.Err(err).Msg("Failed to send kill signal")
				}
				out, err := s.ReadOutput(dir)
				if err != nil {
					s.logger.Err(err).Msg("Error reading output")
				}
				jobRes.ExecLogs = string(out)
				s.logger.Debug().Bytes("output", out).Msg("Exec Result")
				s.checkFailed(s.queries.UpdateJobResultTx(ctx, jobRes))
				return
			}
		case <-done:
			out, err := s.ReadOutput(dir)
			if err != nil {
				s.logger.Err(err).Msg("Error reading output")
			}
			jobRes.ExecLogs = string(out)
			jobRes.Retry = false
			jobRes.Success = true
			s.logger.Debug().Bytes("output", out).Msg("Exec Result")
			s.checkFailed(s.queries.UpdateJobResultTx(ctx, jobRes))
			return
		}
	}
}

func (s *SystemExecutor) checkFailed(_ db.UpdateJobTxResult, err error) {
	if err != nil {
		s.logger.Error().Err(err).Stack().Msgf("An error occurred %s: ", err)
	}
}

func (s *SystemExecutor) writeFiles(ctx context.Context, dir string, execReq *db.ExecRequest) error {
	script, spec, err := execution.ConvertExecSpecToNixScript(ctx, execReq, s.queries)
	if err != nil {
		return fmt.Errorf("error writing files: %s", err)
	}
	files := map[string]string{
		"exec.sh":       script,
		spec.ScriptName: execReq.Code.String,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
func (s *SystemExecutor) Cleanup() {}

func (s *SystemExecutor) ReadOutput(dir string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, "output.txt"))
}
