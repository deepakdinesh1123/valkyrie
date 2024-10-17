package execution

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
)

//go:embed templates
var flakes embed.FS

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

func (s *ExecutionService) prepareExecutionRequest(req *api.ExecutionRequest) (*ExecutionRequest, error) {
	execReq := ExecutionRequest{}
	if req.Timeout.Value > 180 {
		return nil, &ExecutionServiceError{
			Type:    "timeout",
			Message: "Timeout cannot be more than 180 seconds",
		}
	}
	if config.Languages[req.Language] == nil {
		return nil, &ExecutionServiceError{
			Type:    "language",
			Message: "Language not supported",
		}
	}
	execReq.Language = req.Language
	execReq.LangNixPkg = config.Languages[req.Language]["nixPackageName"]
	scriptName := fmt.Sprintf("main.%s", config.Languages[req.Language]["extension"])
	execReq.File = File{
		Name:    scriptName,
		Content: req.Code,
	}
	execReq.Args = req.Environment.Value.Args.Value
	execReq.SystemDependencies = req.Environment.Value.SystemDependencies
	execReq.LanguageDependencies = req.Environment.Value.LanguageDependencies
	flake, err := s.convertExecSpecToFlake(execReq)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "flake",
			Message: err.Error(),
		}
	}

	nixScript, err := s.convertExecSpecToNixScript(execReq)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "nix script",
			Message: err.Error(),
		}
	}
	execReq.Flake = flake
	execReq.NixScript = nixScript
	return &execReq, nil
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec ExecutionRequest) (string, error) {
	tmplName := config.Languages[execSpec.Language]["flake_template"]
	tmplF, err := flakes.ReadFile(fmt.Sprintf("templates/%s", tmplName))
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to get template",
		}
	}
	var res bytes.Buffer
	tmpl, err := template.New(string("flake")).Parse(string(tmplF))
	if err != nil {
		s.logger.Err(err).Msg("failed to parse template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		s.logger.Err(err).Msg("failed to execute template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return res.String(), nil
}

func (s *ExecutionService) convertExecSpecToNixScript(execSpec ExecutionRequest) (string, error) {
	tmplName := config.Languages[execSpec.Language]["script_template"]
	tmplF, err := flakes.ReadFile(fmt.Sprintf("templates/%s", tmplName))
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to get template",
		}
	}
	var res bytes.Buffer
	tmpl, err := template.New(string("nix_script")).Parse(string(tmplF))
	if err != nil {
		s.logger.Err(err).Msg("failed to parse template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		s.logger.Err(err).Msg("failed to execute template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return res.String(), nil
}

func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := s.prepareExecutionRequest(req)
	if err != nil {
		return 0, err
	}
	var jobParams db.AddJobTxParams
	jobParams.Code = execReq.File.Content
	jobParams.Flake = execReq.Flake
	jobParams.NixScript = execReq.NixScript
	jobParams.ProgrammingLanguage = req.Language
	jobParams.MaxRetries = req.MaxRetries.Value
	jobParams.Path = execReq.File.Name
	jobParams.Timeout = req.Timeout.Value

	hash := calculateHash(jobParams.Code, jobParams.ProgrammingLanguage, jobParams.Flake, jobParams.NixScript, jobParams.Path)
	jobParams.Hash = hash

	job, err := s.queries.AddJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job")
		return 0, err
	}
	return int64(job.JobID), nil
}

func calculateHash(code string, language string, flake string, nix_script string, path string) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(code))
	hashInstance.Write([]byte(language))
	hashInstance.Write([]byte(flake))
	hashInstance.Write([]byte(path))
	hashInstance.Write([]byte(nix_script))
	return hex.EncodeToString(hashInstance.Sum(nil))
}
