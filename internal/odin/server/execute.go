package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/mq"

	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/database"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/google/uuid"
)

func (s *Server) Execute(ctx context.Context, req *api.ExecutionRequest) (api.ExecuteRes, error) {
	execRequest := execution.ExecutionRequest{
		ExecutionID: uuid.New().String(),
		Devenv:      req.Devenv,
		File: execution.File{
			Name:    req.File.Name.Value,
			Content: req.File.Content.Value,
		},
	}
	execRequestJSON, err := json.Marshal(execRequest)
	if err != nil {
		return &api.ExecuteBadRequest{
			Message: fmt.Sprintf("Failed to marshal request: %v", err),
		}, nil
	}
	err = mq.Publish("execute", execRequestJSON)
	if err != nil {
		return &api.ExecuteInternalServerError{
			Message: fmt.Sprintf("Failed to publish message: %v", err),
		}, nil
	}
	return &api.ExecutionResult{
		ExecutionID: execRequest.ExecutionID,
	}, nil
}

func (s *Server) GetExecutionResult(ctx context.Context, params api.GetExecutionResultParams) (api.GetExecutionResultRes, error) {
	execId, err := uuid.Parse(params.ExecutionID)
	if err != nil {
		return &api.GetExecutionResultBadRequest{
			Message: fmt.Sprintf("Invalid execution id: %s", params.ExecutionID),
		}, nil
	}
	execResult, err := database.Queries.GetResultUsingExecutionID(ctx, execId)
	if err != nil {
		return &api.GetExecutionResultNotFound{
			Message: fmt.Sprintf("Execution result not found for execution id: %s", params.ExecutionID),
		}, nil
	}
	return &api.ExecutionResult{
		ExecutionID: execResult.ExecutionID.String(),
		Result:      api.NewOptString(execResult.Result.String),
		Status:      api.NewOptString("Executed"),
	}, nil
}
