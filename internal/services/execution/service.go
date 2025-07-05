package execution

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/secret"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/rs/zerolog"
)

//go:embed templates
var ExecTemplates embed.FS

//go:embed scripts
var ExecScripts embed.FS

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
	lang, err := s.queries.GetLanguageByName(ctx, req.Language.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to get language %s: %w", req.Language.Value, err)
	}

	langVersion, err := s.getLanguageVersion(ctx, req, lang.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get language version: %w", err)
	}

	sort.Strings(req.Environment.Value.SystemDependencies)

	execReq := &ExecutionRequest{
		Language:             req.Language.Value,
		Code:                 req.Code.Value,
		CmdLineArgs:          req.CmdLineArgs.Value,
		CompilerArgs:         req.CompilerArgs.Value,
		Input:                req.Input.Value,
		Command:              req.Command.Value,
		SystemDependencies:   req.Environment.Value.SystemDependencies,
		LanguageDependencies: req.Environment.Value.LanguageDependencies,
		Setup:                req.Environment.Value.Setup.Value,
		LangVersion:          langVersion.ID,
		ScriptName:           s.buildScriptName(req, lang),
		Template:             s.getTemplate(langVersion, lang),
		NIXPKGS_REV:          s.envConfig.NIXERY_NIXPKGS_REV,
	}

	if req.Environment.Value.Secrets.Set {
		encodedSecrets, err := secret.EncodeSecrets(req.Environment.Value.Secrets.Value, s.envConfig.ENCKEY)
		if err != nil {
			return nil, fmt.Errorf("failed to encode secrets: %v", err)
		}
		execReq.Secrets = encodedSecrets
	}

	s.applyLanguageSpecificConfig(execReq, langVersion)

	if lang.Name != "generic" {
		err = s.addLanguageNixPackage(execReq, langVersion.NixPackageName.String)
		if err != nil {
			return nil, fmt.Errorf("failed to add : %w", err)
		}
	}

	flake, err := s.convertExecSpecToFlake(*execReq)
	if err != nil {
		return nil, fmt.Errorf("failed to convert exec spec to flake: %w", err)
	}
	execReq.Flake = flake

	return execReq, nil
}

func (s *ExecutionService) addLanguageNixPackage(execReq *ExecutionRequest, langNixPkg string) error {
	execReq.SystemDependencies = append(execReq.SystemDependencies, langNixPkg)
	return nil
}

func (s *ExecutionService) getLanguageVersion(ctx context.Context, req *api.ExecutionRequest, languageID int64) (db.LanguageVersion, error) {
	if req.Version.Value != "" {
		return s.queries.GetLanguageVersion(ctx, db.GetLanguageVersionParams{
			LanguageID: languageID,
			Version:    req.Version.Value,
		})
	}
	return s.queries.GetDefaultVersion(ctx, languageID)
}

func (s *ExecutionService) buildScriptName(req *api.ExecutionRequest, lang db.Language) string {
	extension := req.Extension.Value
	if extension == "" {
		extension = lang.Extension
	}
	return fmt.Sprintf("main.%s", extension)
}

func (s *ExecutionService) getTemplate(langVersion db.LanguageVersion, lang db.Language) string {
	if langVersion.Template.String != "" {
		return langVersion.Template.String
	}
	return lang.Template
}

func (s *ExecutionService) applyLanguageSpecificConfig(execReq *ExecutionRequest, langVersion db.LanguageVersion) {
	switch execReq.Language {
	case "python":
		execReq.SystemSetup = buildPythonSystemSetup(langVersion.Version)
		execReq.PkgIndex = s.envConfig.PY_INDEX
		execReq.SystemDependencies = append(execReq.SystemDependencies, "uv")
	}
}

func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := s.prepareExecutionRequest(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare execution request: %w", err)
	}

	jobParams, err := s.buildJobParams(req, execReq)
	if err != nil {
		return 0, fmt.Errorf("failed to build job parameters: %w", err)
	}

	job, err := s.queries.AddExecJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job to database")
		return 0, fmt.Errorf("failed to add job: %w", err)
	}

	return int64(job.JobID), nil
}

// buildJobParams constructs the database job parameters
func (s *ExecutionService) buildJobParams(req *api.ExecutionRequest, execReq *ExecutionRequest) (db.AddJobTxParams, error) {
	jobParams := db.AddJobTxParams{
		Code:                 execReq.Code,
		Flake:                execReq.Flake,
		LanguageDependencies: execReq.LanguageDependencies,
		SystemDependencies:   execReq.SystemDependencies,
		CmdLineArgs:          execReq.CmdLineArgs,
		CompilerArgs:         execReq.CompilerArgs,
		Files:                req.Files,
		Input:                execReq.Input,
		Command:              execReq.Command,
		LangVersion:          execReq.LangVersion,
		Setup:                execReq.Setup,
		SystemSetup:          execReq.SystemSetup,
		PkgIndex:             execReq.PkgIndex,
		Extension:            execReq.Extension,
		Secrets:              execReq.Secrets,
	}

	// Set max retries
	if err := s.setMaxRetries(req, &jobParams); err != nil {
		return jobParams, err
	}

	// Set timeout
	if err := s.setTimeout(req, &jobParams); err != nil {
		return jobParams, err
	}

	// Calculate and set hash
	jobParams.Hash = calculateHash(jobParams.Code, jobParams.Flake, jobParams.Files, jobParams.Input, jobParams.Secrets)

	return jobParams, nil
}

func (s *ExecutionService) setMaxRetries(req *api.ExecutionRequest, jobParams *db.AddJobTxParams) error {
	if !req.MaxRetries.Set {
		jobParams.MaxRetries = s.envConfig.MAX_RETRIES
		return nil
	}

	if req.MaxRetries.Value > s.envConfig.MAX_RETRIES {
		return fmt.Errorf("max retries (%d) exceeds the maximum allowed (%d)",
			req.MaxRetries.Value, s.envConfig.MAX_RETRIES)
	}

	jobParams.MaxRetries = req.MaxRetries.Value
	return nil
}

func (s *ExecutionService) setTimeout(req *api.ExecutionRequest, jobParams *db.AddJobTxParams) error {
	if !req.Timeout.Set {
		jobParams.Timeout = int32(s.envConfig.WORKER_MAX_TASK_TIMEOUT)
		return nil
	}

	if req.Timeout.Value > int32(s.envConfig.WORKER_MAX_TASK_TIMEOUT) {
		return fmt.Errorf("timeout (%d) exceeds the maximum allowed (%d)",
			req.Timeout.Value, s.envConfig.WORKER_MAX_TASK_TIMEOUT)
	}

	jobParams.Timeout = req.Timeout.Value
	return nil
}

func buildPythonSystemSetup(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		return fmt.Sprintf("export UV_PYTHON=$(which python%s.%s)", parts[0], parts[1])
	}
	return fmt.Sprintf("export UV_PYTHON=$(which python%s)", version)
}

func calculateHash(code, flake string, files []byte, input string, secrets []byte) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	hasher.Write([]byte(flake))
	hasher.Write(files)
	hasher.Write([]byte(input))
	hasher.Write(secrets)
	return hex.EncodeToString(hasher.Sum(nil))
}
