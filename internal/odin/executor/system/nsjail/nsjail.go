//go:build system || all

package nsjail

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type NSJailExecutor struct {
	envConfig *config.EnvConfig
	queries   db.Store
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	workerId  int32
}

type JailSpec struct {
	FileSource string
}

type JailExecSpec struct {
	Setup       string
	Command     string
	CmdLineArgs string
	CompileArgs string
	Input       string
}

func NewNSJailExecutor(ctx context.Context, envConfig *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*NSJailExecutor, error) {
	return &NSJailExecutor{
		envConfig: envConfig,
		queries:   queries,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}

func (ns *NSJailExecutor) GetExecCmd(ctx context.Context, outFile *os.File, errFile *os.File, dir string, execReq *db.ExecRequest) (*exec.Cmd, error) {
	storePaths, err := ns.queries.GetPackageStorePaths(ctx, execReq.SystemDependencies)
	if err != nil {
		return nil, err
	}

	var jailConf bytes.Buffer
	tmpl, err := template.New(string("base.nsjail.tmpl")).ParseFS(
		execution.ExecTemplates,
		"templates/base.nsjail.tmpl",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %s", err)
	}

	err = tmpl.Execute(&jailConf, &JailSpec{
		FileSource: dir,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %s", err)
	}

	file, err := os.Create(filepath.Join(dir, "config.proto"))
	if err != nil {
		return nil, fmt.Errorf("error creating config.proto: %s", err)
	}
	defer file.Close() // Ensure the file is closed when done

	// Write content to the file
	_, err = file.WriteString(jailConf.String())
	if err != nil {
		return nil, fmt.Errorf("could not write config.proto file: %s", err)
	}

	langVersion, err := ns.queries.GetLanguageVersionByID(ctx, execReq.LanguageVersion)
	language, err := ns.queries.GetLanguageByID(ctx, langVersion.LanguageID)

	var jailExec bytes.Buffer
	tmplName := filepath.Join("templates", language.Name, langVersion.Template)

	execTmpl, err := template.New(string("base.nsexec.tmpl")).ParseFS(
		execution.ExecTemplates,
		"templates/base.nsexec.tmpl",
		tmplName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %s", err)
	}

	err = execTmpl.Execute(&jailExec, &JailExecSpec{
		Setup:       execReq.Setup.String,
		Command:     execReq.Command.String,
		CmdLineArgs: execReq.CmdLineArgs.String,
		CompileArgs: execReq.CompileArgs.String,
		Input:       execReq.Input.String,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %s", err)
	}

	execFile, err := os.Create(filepath.Join(dir, "nsexec.sh"))
	if err != nil {
		return nil, fmt.Errorf("error creating config.proto: %s", err)
	}
	defer file.Close() // Ensure the file is closed when done

	// Write content to the file
	_, err = execFile.WriteString(jailExec.String())
	if err != nil {
		return nil, fmt.Errorf("could not write config.proto file: %s", err)
	}

	var nsjailPath strings.Builder

	_, err = nsjailPath.WriteString("NIXPATHS=")
	if err != nil {
		return nil, fmt.Errorf("error while creating nsjail path")
	}

	for _, path := range storePaths {
		_, err := nsjailPath.WriteString(fmt.Sprintf("%s:", path.StorePath.String))
		if err != nil {
			return nil, fmt.Errorf("error while adding store paths: %s", err)
		}
	}

	languageStorePaths, err := ns.queries.GetPackageStorePaths(ctx, []string{langVersion.NixPackageName})
	if err != nil {
		return nil, fmt.Errorf("error while fetching language store paths: %s", err)
	}

	if len(languageStorePaths) == 0 {
		return nil, fmt.Errorf("language store path not found in DB")
	}

	nsjailPath.WriteString(languageStorePaths[0].StorePath.String)

	ns.logger.Debug().Msgf("Jail PATH is %s", nsjailPath.String())

	args := []string{"--config", fmt.Sprintf("%s/config.proto", dir), "-E", nsjailPath.String()}
	execCmd := exec.CommandContext(ctx, "nsjail", args...)
	execCmd.Cancel = func() error {
		ns.logger.Info().Msg("Task timed out. Terminating execution")
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
