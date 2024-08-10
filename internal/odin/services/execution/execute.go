package execution

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

//go:embed templates
var flakes embed.FS

type ExecutionService struct {
	queries   *db.Queries
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
}

// NewExecutionService returns a new instance of ExecutionService.
//
// Parameters:
// - queries: database queries
// - envConfig: environment configuration
// - logger: logger instance
//
// Returns:
// - *ExecutionService: a new instance of ExecutionService
func NewExecutionService(queries *db.Queries, envConfig *config.EnvConfig, logger *zerolog.Logger) *ExecutionService {
	return &ExecutionService{
		queries:   queries,
		envConfig: envConfig,
		logger:    logger,
	}
}

// prepareExecutionRequest Prepares an execution request based on the provided ExecutionRequest.
//
// Parameters:
// - req: The ExecutionRequest to prepare.
//
// Returns:
// - *models.ExecutionRequest: The prepared ExecutionRequest.
// - error: An error if the preparation fails.
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

// convertExecSpecToFlake converts an execution spec to a flake.
//
// Parameters:
// - execSpec: The execution request to convert.
// Returns:
// - string: The flake.
// - error: An error if the conversion fails.
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

// AddJob adds a job to the ExecutionService.
//
// Parameters:
// - ctx: The context.Context for the request.
// - req: The api.ExecutionRequest containing the job details.
//
// Returns:
// - int64: The ID of the inserted job.
// - error: An error if the job insertion fails.
func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := s.prepareExecutionRequest(req)
	if err != nil {
		return 0, err
	}
	job, err := s.queries.InsertJob(ctx, db.InsertJobParams{
		Script:     execReq.File.Content,
		Flake:      execReq.Environment,
		Language:   req.Language,
		ScriptPath: execReq.File.Name,
		Args:       pgtype.Text{String: execReq.Args, Valid: true},
	})
	if err != nil {
		return 0, err
	}
	return job.ID, nil
}
