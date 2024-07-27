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
)

//go:embed templates
var flakes embed.FS

type ExecutionService struct {
	queries   *db.Queries
	envConfig *config.EnvConfig
}

func NewExecutionService(queries *db.Queries, envConfig *config.EnvConfig) *ExecutionService {
	return &ExecutionService{
		queries:   queries,
		envConfig: envConfig,
	}
}

func PrepareExecutionRequest(req *api.ExecutionRequest) (*models.ExecutionRequest, error) {
	var scriptName string
	if !req.File.Name.Set {
		scriptName = fmt.Sprintf("main.%s", config.LANGUAGE_EXTENSION[req.Language.String])
	} else {
		scriptName = req.File.Name.Value
	}
	if req.Environment.Type == "ExecutionEnvironmentSpec" {
		flake, err := ConvertExecSpecToFlake(req.Environment.ExecutionEnvironmentSpec, req.Language.String)
		if err != nil {
			return nil, &ExecutionServiceError{
				Type:    "flake",
				Message: "failed to convert execution environment spec to flake",
			}
		}
		return &models.ExecutionRequest{
			Environment: flake,
			File: models.File{
				Name:    scriptName,
				Content: req.File.Content,
			},
		}, nil
	}
	return &models.ExecutionRequest{
		Environment: string(req.Environment.Flake),
		File: models.File{
			Name:    scriptName,
			Content: req.File.Content,
		},
		Language:   req.Language.String,
		ScriptName: scriptName,
	}, nil
}

func ConvertExecSpecToFlake(execSpec api.ExecutionEnvironmentSpec, language string) (string, error) {
	tmpl, err := flakes.ReadFile(fmt.Sprintf("templates/%s", language))
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}
	buffer := new(bytes.Buffer)
	err = template.Must(template.New("tmpl").Parse(string(tmpl))).Execute(buffer, execSpec)
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return buffer.String(), nil
}

func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := PrepareExecutionRequest(req)
	if err != nil {
		return 0, err
	}
	job, err := s.queries.InsertJob(ctx, db.InsertJobParams{
		Script: pgtype.Text{
			String: execReq.File.Content,
			Valid:  true},
		Flake:      pgtype.Text{String: execReq.Environment, Valid: true},
		Language:   pgtype.Text{String: execReq.Language, Valid: true},
		ScriptPath: pgtype.Text{String: execReq.ScriptName, Valid: true},
	})
	if err != nil {
		return 0, err
	}
	return job.ID, nil
}
