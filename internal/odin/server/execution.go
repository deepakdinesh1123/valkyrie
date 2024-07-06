package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/google/uuid"
)

func (s *Server) Execute(ctx context.Context, req *api.ExecutionRequest) (api.ExecuteRes, error) {
	executionId := uuid.New()
	execRequest := execution.ExecutionRequest{
		ExecutionID: executionId.String(),
		Environment: req.Environment,
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
	_, err = s.Queries.InsertExecutionRequest(
		ctx,
		db.InsertExecutionRequestParams{
			ExecutionID: executionId,
			Environment: pgtype.Text{String: req.Environment, Valid: true},
			Code:        pgtype.Text{String: req.File.Content.Value, Valid: true},
		},
	)
	if err != nil {
		return &api.ExecuteInternalServerError{
			Message: fmt.Sprintf("Failed to insert execution request: %v", err),
		}, nil
	}
	return &api.ExecuteOK{
		ExecutionID: executionId.String(),
	}, nil
}

func (s *Server) GetExecutionResult(ctx context.Context, params api.GetExecutionResultParams) (api.GetExecutionResultRes, error) {
	execId, err := uuid.Parse(params.ExecutionID)
	if err != nil {
		return &api.GetExecutionResultBadRequest{
			Message: fmt.Sprintf("Invalid execution id: %s", params.ExecutionID),
		}, nil
	}
	execResult, err := s.Queries.GetResultUsingExecutionID(ctx, execId)
	if err != nil {
		return &api.GetExecutionResultNotFound{
			Message: fmt.Sprintf("Execution result not found for execution id: %s", params.ExecutionID),
		}, nil
	}
	return &api.ExecutionResult{
		ExecutionID: execResult.ExecutionID.String(),
		Result:      execResult.Result.String,
		Status:      execResult.ExecutionStatus.String,
	}, nil
}

func (s *Server) GetExecutions(ctx context.Context) (api.GetExecutionsRes, error) {
	executionsDB, err := s.Queries.GetAllExecutions(ctx)
	if err != nil {
		return &api.GetExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get executions: %v", err),
		}, nil
	}
	var executions api.GetExecutionsOKApplicationJSON
	for _, exec := range executionsDB {
		executions = append(executions, api.Execution{
			ExecutionID:     exec.ExecutionID.String(),
			Code:            api.NewOptString(exec.Code.String),
			Environment:     exec.Environment.String,
			RequestedAt:     exec.RequestedAt.Time.Format(time.RFC3339),
			Result:          api.NewOptString(exec.Result.String),
			ExecutionStatus: api.NewOptString(exec.ExecutionStatus.String),
			ExecutedAt:      api.NewOptString(exec.ExecutedAt.Time.Format(time.RFC3339)),
		})
	}
	return &executions, nil
}

func (s *Server) GetExecutionResults(ctx context.Context) (api.GetExecutionResultsRes, error) {
	execResultsDB, err := s.Queries.GetAllExecutionResults(ctx)
	if err != nil {
		return &api.GetExecutionResultsInternalServerError{
			Message: fmt.Sprintf("Failed to get execution results: %v", err),
		}, nil
	}
	var execResults api.GetExecutionResultsOKApplicationJSON
	for _, execResult := range execResultsDB {
		execResults = append(execResults, api.ExecutionResult{
			ExecutionID: execResult.ExecutionID.String(),
			Result:      execResult.Result.String,
			Status:      execResult.ExecutionStatus.String,
		})
	}
	return &execResults, nil
}
