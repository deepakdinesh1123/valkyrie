//go:build system || all

package native

import (
	"context"
	"os"
	"os/exec"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
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

func (s *SystemExecutor) GetExecCmd(ctx context.Context, outFile *os.File, errFile *os.File, dir string, execReq *db.ExecRequest) (*exec.Cmd, error) {
	execCmd := exec.CommandContext(ctx, "sh", "exec.sh")
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
	return execCmd, nil
}
