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
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
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

func (s *ExecutionService) prepareExecutionRequest(req *api.ExecutionRequest) (*models.ExecutionRequest, error) {
	scriptName := fmt.Sprintf("main.%s", config.LANGUAGE_EXTENSION[req.Language])
	if req.Environment.Value.Type == "Flake" {
		return &models.ExecutionRequest{
			Environment: string(req.Environment.Value.Flake),
			File: models.File{
				Name:    scriptName,
				Content: req.Code,
			},
		}, nil
	} else if req.Environment.Value.Type == "ExecutionEnvironmentSpec" {
		flake, err := s.convertExecSpecToFlake(req)
		if err != nil {
			return nil, &ExecutionServiceError{
				Type:    "flake",
				Message: err.Error(),
			}
		}
		return &models.ExecutionRequest{
			Environment: flake,
			File: models.File{
				Name:    scriptName,
				Content: req.Code,
			},
			Args: req.Environment.Value.ExecutionEnvironmentSpec.Args.Value,
		}, nil
	} else if !req.Environment.Set {
		flake, err := s.convertExecSpecToFlake(req)
		if err != nil {
			return nil, &ExecutionServiceError{
				Type:    "flake",
				Message: err.Error(),
			}
		}
		return &models.ExecutionRequest{
			Environment: flake,
			File: models.File{
				Name:    scriptName,
				Content: req.Code,
			},
			Args: req.Environment.Value.ExecutionEnvironmentSpec.Args.Value,
		}, nil
	}
	return nil, &ExecutionServiceError{
		Type:    "environment",
		Message: "invalid environment type",
	}
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec *api.ExecutionRequest) (string, error) {
	tmplF, err := flakes.ReadFile(fmt.Sprintf("templates/%s.tmpl", execSpec.Language))
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

	err = tmpl.Execute(&res, execSpec.Environment.Value.ExecutionEnvironmentSpec)
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
	if req.Environment.Value.Type == "" {
		return 0, &ExecutionServiceError{
			Type:    "environment",
			Message: "Specify flake or language",
		}
	}
	execReq, err := s.prepareExecutionRequest(req)
	if err != nil {
		return 0, err
	}
	var jobParams db.AddJobTxParams
	jobParams.Code = execReq.File.Content
	jobParams.Flake = execReq.Environment
	jobParams.ProgrammingLanguage = req.Language
	jobParams.MaxRetries = req.MaxRetries.Value
	jobParams.Path = execReq.File.Name
	if req.CronExpression.Set {
		jobParams.CronExpression = req.CronExpression.Value
	}

	hash := calculateHash(jobParams.Code, jobParams.ProgrammingLanguage, jobParams.Flake, jobParams.Path)
	jobParams.Hash = hash

	job, err := s.queries.AddJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job")
		return 0, err
	}
	return int64(job.JobID), nil
}

func calculateHash(code string, language string, flake string, path string) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(code))
	hashInstance.Write([]byte(language))
	hashInstance.Write([]byte(flake))
	hashInstance.Write([]byte(path))
	return hex.EncodeToString(hashInstance.Sum(nil))
}
