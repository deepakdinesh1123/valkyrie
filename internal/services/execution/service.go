package execution

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/rs/zerolog"
)

//go:embed templates
var ExecTemplates embed.FS

type ExecutionService struct {
	queries   db.Store
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
}

func NewExecutionService(queries db.Store, envConfig *config.EnvConfig, logger *zerolog.Logger) *ExecutionService {
	return &ExecutionService{
		queries:   queries,
		envConfig: envConfig,
		logger:    logger,
	}
}

func (s *ExecutionService) prepareExecutionRequest(ctx context.Context, req *api.ExecutionRequest) (*ExecutionRequest, error) {
	execReq := ExecutionRequest{}

	lang, err := s.queries.GetLanguageByName(ctx, req.Language.Value)
	if err != nil {
		return nil, err
	}

	var langVersion db.LanguageVersion

	if req.Version.Value != "" {
		langVersion, err = s.queries.GetLanguageVersion(ctx, db.GetLanguageVersionParams{
			LanguageID: lang.ID,
			Version:    req.Version.Value,
		})
	} else {
		langVersion, err = s.queries.GetDefaultVersion(ctx, lang.ID)
	}

	if err != nil {
		return nil, err
	}

	execReq.LangNixPkg = langVersion.NixPackageName
	execReq.Language = req.Language.Value
	execReq.Code = req.Code.Value
	execReq.CmdLineArgs = req.CmdLineArgs.Value
	execReq.CompilerArgs = req.CompilerArgs.Value
	execReq.Input = req.Input.Value
	execReq.Command = req.Command.Value
	execReq.SystemDependencies = req.Environment.Value.SystemDependencies
	execReq.LanguageDependencies = req.Environment.Value.LanguageDependencies
	execReq.Setup = req.Environment.Value.Setup.Value
	execReq.ScriptName = fmt.Sprintf("main.%s", lang.Extension)
	execReq.LangVersion = langVersion.ID

	if execReq.Language == "python" {
		execReq.SystemSetup = getPythonSystemSetup(langVersion.Version)
		execReq.PkgIndex = s.envConfig.PY_INDEX
	}

	if langVersion.Template.String != "" {
		execReq.Template = langVersion.Template.String
	} else {
		execReq.Template = lang.Template
	}

	flake, err := s.convertExecSpecToFlake(execReq)
	if err != nil {
		return nil, fmt.Errorf("error converting exec spec to flake: %s", err)
	}
	execReq.Flake = flake
	return &execReq, nil
}

func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := s.prepareExecutionRequest(ctx, req)
	if err != nil {
		return 0, err
	}
	var jobParams db.AddJobTxParams
	jobParams.Code = execReq.Code
	jobParams.Flake = execReq.Flake
	jobParams.LanguageDependencies = execReq.LanguageDependencies
	jobParams.SystemDependencies = execReq.SystemDependencies
	jobParams.CmdLineArgs = execReq.CmdLineArgs
	jobParams.CompilerArgs = execReq.CompilerArgs
	jobParams.Files = req.Files
	jobParams.Input = execReq.Input
	jobParams.Command = execReq.Command
	jobParams.LangVersion = execReq.LangVersion
	jobParams.Setup = execReq.Setup
	jobParams.SystemSetup = execReq.SystemSetup
	jobParams.PkgIndex = execReq.PkgIndex

	s.logger.Debug().Int64("Version", jobParams.LangVersion).Msg("Language")

	if !req.MaxRetries.Set {
		jobParams.MaxRetries = s.envConfig.MAX_RETRIES
	} else {
		if req.MaxRetries.Value > s.envConfig.MAX_RETRIES {
			return 0, fmt.Errorf("retries exceeds the maximum number of retries allowed")
		}
		jobParams.MaxRetries = req.MaxRetries.Value
	}

	if !req.Timeout.Set {
		jobParams.Timeout = int32(s.envConfig.WORKER_TASK_TIMEOUT)
	} else {
		if req.Timeout.Value > int32(s.envConfig.WORKER_TASK_TIMEOUT) {
			return 0, fmt.Errorf("retries exceeds the maximum number of retries allowed")
		}
		jobParams.Timeout = req.Timeout.Value
	}

	hash := calculateHash(jobParams.Code, jobParams.Flake, jobParams.Files, jobParams.Input)
	jobParams.Hash = hash

	job, err := s.queries.AddExecJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job")
		return 0, err
	}
	return int64(job.JobID), nil
}

func getPythonSystemSetup(version string) string {
	ver := strings.Split(version, ".")
	if len(ver) > 1 {
		return fmt.Sprintf(`export UV_PYTHON=$(which python%s)`, fmt.Sprintf("%s.%s", ver[0], ver[1]))
	}
	return fmt.Sprintf(`export UV_PYTHON=$(which python%s)`, version)
}

func calculateHash(code string, flake string, files []byte, input string) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(code))
	hashInstance.Write([]byte(flake))
	hashInstance.Write(files)
	hashInstance.Write([]byte(input))
	return hex.EncodeToString(hashInstance.Sum(nil))
}
