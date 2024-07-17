package execution

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/models"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/google/uuid"
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

func (s *ExecutionService) Execute(ctx context.Context, req *api.ExecutionRequest) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func PrepareExecutionRequest(req *api.ExecutionRequest) (*models.ExecutionRequest, error) {
	if req.Environment.Type == "ExecutionEnvironmentSpec" {
		flake, err := ConvertExecSpecToFlake(req.Environment.ExecutionEnvironmentSpec)
		if err != nil {
			return nil, &ExecutionServiceError{
				Type:    "flake",
				Message: "failed to convert execution environment spec to flake",
			}
		}
		return &models.ExecutionRequest{
			Environment: flake,
			File: models.File{
				Name:    req.File.Name.Value,
				Content: req.File.Content.Value,
			},
		}, nil
	}
	return &models.ExecutionRequest{
		Environment: string(req.Environment.Flake),
		File: models.File{
			Name:    req.File.Name.Value,
			Content: req.File.Content.Value,
		},
	}, nil
}

func ConvertExecSpecToFlake(execSpec api.ExecutionEnvironmentSpec) (string, error) {
	tmpl, err := flakes.ReadFile(fmt.Sprintf("templates/%s", execSpec.Language.Name))
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
