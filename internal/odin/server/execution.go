package server

import (
	"context"
	"fmt"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *Server) Execute(ctx context.Context, req *api.ExecutionRequest) (api.ExecuteRes, error) {
	execId, err := s.executionService.Execute(ctx, req)
	if err != nil {
		switch err.(type) {
		case *execution.ExecutionServiceError:
			return &api.ExecuteInternalServerError{
				Message: fmt.Sprintf("Failed to execute: %v", err),
			}, nil
		case *execution.TemplateError:
			return &api.ExecuteBadRequest{
				Message: fmt.Sprintf("Failed to execute: %v", err),
			}, nil
		default:
			return &api.ExecuteInternalServerError{
				Message: fmt.Sprintf("Failed to execute: %v", err),
			}, nil
		}
	}
	return &api.ExecuteOK{ExecutionId: execId}, nil
}

func (s *Server) GetExecutionResult(ctx context.Context, params api.GetExecutionResultParams) (api.GetExecutionResultRes, error) {
	execResult, err := s.queries.GetResultUsingExecutionID(ctx, params.ExecutionId)
	if err != nil {
		return &api.GetExecutionResultNotFound{}, nil
	}

	return &api.ExecutionResult{
		ExecutionId: execResult.ExecutionID,
		Result:      execResult.Result.String,
		Status:      execResult.ExecutionStatus.String,
	}, nil
}

func (s *Server) GetExecutions(ctx context.Context) (api.GetExecutionsRes, error) {
	executionsDB, err := s.queries.GetAllExecutions(ctx)
	if err != nil {
		return &api.GetExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get executions: %v", err),
		}, nil
	}
	var executions api.GetExecutionsOKApplicationJSON
	for _, exec := range executionsDB {
		executions = append(executions, api.Execution{
			ExecutionId:     exec.ExecutionID,
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
	execResultsDB, err := s.queries.GetAllExecutionResults(ctx)
	if err != nil {
		return &api.GetExecutionResultsInternalServerError{
			Message: fmt.Sprintf("Failed to get execution results: %v", err),
		}, nil
	}
	var execResults api.GetExecutionResultsOKApplicationJSON
	for _, execResult := range execResultsDB {
		execResults = append(execResults, api.ExecutionResult{
			ExecutionId: execResult.ExecutionID,
			Result:      execResult.Result.String,
			Status:      execResult.ExecutionStatus.String,
		})
	}
	return &execResults, nil
}
