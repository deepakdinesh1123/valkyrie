package execution

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
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

	langVersion, err := s.queries.GetLanguageVersion(ctx, db.GetLanguageVersionParams{
		LanguageID: lang.ID,
		Version:    req.Version.Value,
	})

	if err != nil {
		return nil, err
	}

	execReq.LangNixPkg = langVersion.NixPackageName
	execReq.Language = req.Language.Value
	execReq.Code = req.Code.Value
	execReq.CmdLineArgs = req.CmdLineArgs.Value
	execReq.CompileArgs = req.CompileArgs.Value
	execReq.Input = req.Input.Value
	execReq.Command = req.Command.Value
	execReq.SystemDependencies = req.Environment.Value.SystemDependencies
	execReq.LanguageDependencies = req.Environment.Value.LanguageDependencies
	execReq.Setup = req.Environment.Value.Setup.Value
	execReq.ScriptName = fmt.Sprintf("main.%s", lang.Extension)
	execReq.Template = langVersion.Template
	execReq.LangVersion = langVersion.ID

	flake, err := s.convertExecSpecToFlake(execReq)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "flake",
			Message: err.Error(),
		}
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
	jobParams.CompileArgs = execReq.CompileArgs
	jobParams.Files = req.Files
	jobParams.Input = execReq.Input
	jobParams.Command = execReq.Command
	jobParams.LangVersion = execReq.LangVersion

	s.logger.Debug().Int64("Version", jobParams.LangVersion).Msg("Language")

	if !req.MaxRetries.Set {
		jobParams.MaxRetries = s.envConfig.ODIN_MAX_RETRIES
	} else {
		if req.MaxRetries.Value > s.envConfig.ODIN_MAX_RETRIES {
			return 0, fmt.Errorf("retries exceeds the maximum number of retries allowed")
		}
		jobParams.MaxRetries = req.MaxRetries.Value
	}

	if !req.Timeout.Set {
		jobParams.Timeout = int32(s.envConfig.ODIN_WORKER_TASK_TIMEOUT)
	} else {
		if req.Timeout.Value > int32(s.envConfig.ODIN_WORKER_TASK_TIMEOUT) {
			return 0, fmt.Errorf("retries exceeds the maximum number of retries allowed")
		}
		jobParams.Timeout = req.Timeout.Value
	}

	hash := calculateHash(jobParams.Code, jobParams.Flake, jobParams.Files)
	jobParams.Hash = hash

	job, err := s.queries.AddJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job")
		return 0, err
	}
	return int64(job.JobID), nil
}

func calculateHash(code string, flake string, files []byte) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(code))
	hashInstance.Write([]byte(flake))
	hashInstance.Write(files)
	return hex.EncodeToString(hashInstance.Sum(nil))
}
