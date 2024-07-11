package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExecutionService struct {
	queries *db.Queries
	// queue          *mq.MessageQueue
	valkyrieConfig *config.ValkyrieConfig
}

func NewExecutionService(queries *db.Queries, valkyrieConfig *config.ValkyrieConfig) *ExecutionService {
	return &ExecutionService{
		queries:        queries,
		valkyrieConfig: valkyrieConfig,
	}
}

func (s *ExecutionService) Execute(ctx context.Context, req *api.ExecutionRequest) (uuid.UUID, error) {
	executionId := uuid.New()
	execRequest, err := PrepareExecutionRequest(req)
	if err != nil {
		return uuid.Nil, err
	}
	execRequest.ExecutionID = executionId
	execRequestJSON, err := json.Marshal(execRequest)
	if err != nil {
		return uuid.Nil, &ExecutionServiceError{
			Type:    "json",
			Message: "failed to marshal execution request",
		}
	}
	err = mq.Publish("execute", execRequestJSON)
	if err != nil {
		return uuid.Nil, &ExecutionServiceError{
			Type:    "mq",
			Message: "failed to publish execution request",
		}
	}
	_, err = s.queries.InsertExecutionRequest(
		ctx,
		db.InsertExecutionRequestParams{
			ExecutionID: executionId,
			Environment: pgtype.Text{String: execRequest.Environment, Valid: true},
			Code:        pgtype.Text{String: req.File.Content.Value, Valid: true},
		},
	)
	if err != nil {
		return uuid.Nil, &ExecutionServiceError{
			Type:    "db",
			Message: "failed to insert execution request",
		}
	}
	return executionId, nil
}

func PrepareExecutionRequest(req *api.ExecutionRequest) (*execution.ExecutionRequest, error) {
	if req.Environment.Type == "ExecutionEnvironmentSpec" {
		flake, err := ConvertExecSpecToFlake(req.Environment.ExecutionEnvironmentSpec)
		if err != nil {
			return nil, &ExecutionServiceError{
				Type:    "flake",
				Message: "failed to convert execution environment spec to flake",
			}
		}
		return &execution.ExecutionRequest{
			Environment: flake,
			File: execution.File{
				Name:    req.File.Name.Value,
				Content: req.File.Content.Value,
			},
		}, nil
	}
	return &execution.ExecutionRequest{
		Environment: string(req.Environment.Flake),
		File: execution.File{
			Name:    req.File.Name.Value,
			Content: req.File.Content.Value,
		},
	}, nil
}

func ConvertExecSpecToFlake(execSpec api.ExecutionEnvironmentSpec) (string, error) {
	template_path := fmt.Sprintf("internal/odin/services/execution/templates/%s.tmpl", execSpec.Language.Name)
	flakeTmpl, err := template.ParseFiles(template_path)
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}
	buffer := new(bytes.Buffer)
	err = flakeTmpl.Execute(buffer, execSpec)
	if err != nil {
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return buffer.String(), nil
}
