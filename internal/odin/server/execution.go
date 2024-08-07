package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *Server) Execute(ctx context.Context, req *api.ExecutionRequest) (api.ExecuteRes, error) {
	execId, err := s.executionService.AddJob(ctx, req)
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
		ExecutionId: execResult.ID,
		Logs:        execResult.Logs.String,
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
			ExecutionId: exec.ID,
		})
	}
	return &executions, nil
}

func (s *Server) GetExecutionConfig(ctx context.Context) (api.GetExecutionConfigRes, error) {
	return &api.ExecutionConfig{
		ODINWORKERPROVIDER:    s.envConfig.ODIN_WORKER_PROVIDER,
		ODINWORKERCONCURRENCY: int32(s.envConfig.ODIN_WORKER_CONCURRENCY),
		ODINWORKERBUFFERSIZE:  int32(s.envConfig.ODIN_WORKER_BUFFER_SIZE),
		ODINWORKERTASKTIMEOUT: s.envConfig.ODIN_WORKER_TASK_TIMEOUT,
		ODINWORKERPOLLFREQ:    s.envConfig.ODIN_WORKER_POLL_FREQ,
		ODINWORKERRUNTIME:     s.envConfig.ODIN_WORKER_RUNTIME,
		ODINLOGLEVEL:          s.envConfig.ODIN_LOG_LEVEL,
	}, nil
}

func (s *Server) GetExecutionWorkers(ctx context.Context) (api.GetExecutionWorkersRes, error) {
	workersDB, err := s.queries.GetAllWorkers(ctx)
	if err != nil {
		return &api.GetExecutionWorkersInternalServerError{
			Message: fmt.Sprintf("Failed to get workers: %v", err),
		}, nil
	}
	var workers api.GetExecutionWorkersOKApplicationJSON
	for _, worker := range workersDB {
		workers = append(workers, api.ExecutionWorker{
			ID:        int64(worker.ID),
			Name:      worker.Name.String,
			CreatedAt: worker.CreatedAt.Time,
			Status:    "status",
		})
	}
	return &workers, nil
}

func (s *Server) DeleteExecution(ctx context.Context, params api.DeleteExecutionParams) (api.DeleteExecutionRes, error) {
	err := s.queries.DeleteJob(ctx, params.ExecutionId)
	if err != nil {
		return &api.DeleteExecutionBadRequest{
			Message: fmt.Sprintf("Failed to delete execution: %v", err),
		}, nil
	}
	return &api.DeleteExecutionOK{}, nil
}
