package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

// Execute Adds a job to the execution queue and returns the execution ID.
//
// Parameters:
// - ctx: the context of the execution request
// - req: the execution request
//
// Returns:
// - api.ExecuteRes: the result of the execution
// - error: any error that occurred during execution
func (s *OdinServer) Execute(ctx context.Context, req *api.ExecutionRequest) (api.ExecuteRes, error) {
	execId, err := s.executionService.AddJob(ctx, req)
	if err != nil {
		switch err.(type) {
		case *execution.ExecutionServiceError:
			return &api.ExecuteInternalServerError{
				Message: fmt.Sprintf("Execution Service: %v", err),
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

// GetAllExecutionResults Retrieves all execution results.
//
// Parameters:
// - ctx: The context of the request.
// - params: The parameters of the request, including pagination options.
//
// Returns:
// - api.GetAllExecutionResultsRes: The response containing the execution results.
// - error: Any error that occurred during the request.
func (s *OdinServer) GetAllExecutionResults(ctx context.Context, params api.GetAllExecutionResultsParams) (api.GetAllExecutionResultsRes, error) {
	execResDB, err := s.queries.GetAllExecutionResults(ctx, db.GetAllExecutionResultsParams{
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetAllExecutionResultsInternalServerError{
			Message: fmt.Sprintf("Failed to get execution results: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalExecutions(ctx)
	if err != nil {
		return &api.GetAllExecutionResultsInternalServerError{
			Message: fmt.Sprintf("Failed to get total execution results: %v", err),
		}, nil
	}
	var executions []api.ExecutionResult
	for _, execRes := range execResDB {
		executions = append(executions, api.ExecutionResult{
			ExecutionId: execRes.ID,
			Logs:        execRes.Logs.String,
			Script:      execRes.Script,
			Flake:       execRes.Flake,
			Args:        execRes.Args.String,
			StartedAt:   execRes.StartedAt.Time,
			FinishedAt:  execRes.FinishedAt.Time,
			CreatedAt:   execRes.CreatedAt.Time,
		})
	}
	resp := api.GetAllExecutionResultsOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
}

// GetExecutionResultsById returns the execution results for a specific job ID.
//
// Parameters:
// - ctx: The context of the request.
// - params: The parameters of the request, including the job ID, page size, and page value.
//
// Returns:
// - api.GetExecutionResultsByIdRes: The response containing the execution results.
// - error: Any error that occurred during the request.
func (s *OdinServer) GetExecutionResultsById(ctx context.Context, params api.GetExecutionResultsByIdParams) (api.GetExecutionResultsByIdRes, error) {
	execRes, err := s.queries.GetExecutionResultsByID(ctx, db.GetExecutionResultsByIDParams{
		JobID:  params.JobId,
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetExecutionResultsByIdInternalServerError{
			Message: fmt.Sprintf("Failed to get execution results: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalExecutionsForJob(ctx, params.JobId)
	if err != nil {
		return &api.GetExecutionResultsByIdInternalServerError{
			Message: fmt.Sprintf("Failed to get total execution results: %v", err),
		}, nil
	}
	var executions []api.ExecutionResult
	for _, execRes := range execRes {
		executions = append(executions, api.ExecutionResult{
			ExecutionId: execRes.ID,
			Logs:        execRes.Logs.String,
			Script:      execRes.Script,
			Flake:       execRes.Flake,
			Args:        execRes.Args.String,
			StartedAt:   execRes.StartedAt.Time,
			FinishedAt:  execRes.FinishedAt.Time,
			CreatedAt:   execRes.CreatedAt.Time,
		})
	}
	resp := api.GetExecutionResultsByIdOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
}

// GetAllExecutions Retrieves a list of all executions.
//
// Parameters:
// - ctx: The context for the request.
// - params: The parameters for the request, including pagination options.
//
// Returns:
// - api.GetAllExecutionsRes: The response containing the list of executions and pagination metadata.
// - error: Any error that occurred during the request.
func (s *OdinServer) GetAllExecutions(ctx context.Context, params api.GetAllExecutionsParams) (api.GetAllExecutionsRes, error) {
	executionsDB, err := s.queries.GetAllJobs(ctx, db.GetAllJobsParams{
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetAllExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get executions: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalJobs(ctx)
	if err != nil {
		return &api.GetAllExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get total executions: %v", err),
		}, nil
	}
	var executions []api.Execution
	for _, exec := range executionsDB {
		executions = append(executions, api.Execution{
			ExecutionId: exec.ID,
			Script:      exec.Script,
			Flake:       exec.Flake,
			CreatedAt:   exec.InsertedAt.Time,
		})
	}
	resp := api.GetAllExecutionsOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
}

// GetExecutionConfig returns the execution configuration.
//
// Parameters:
// - ctx (context.Context): the context for the request.
//
// Returns:
// - api.GetExecutionConfigRes: the execution configuration response.
// - error: any error that occurred while retrieving the configuration.
func (s *OdinServer) GetExecutionConfig(ctx context.Context) (api.GetExecutionConfigRes, error) {
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

// GetExecutionWorkers returns a list of execution workers based on the provided pagination parameters.
//
// Parameters:
// - ctx: The context for the request.
// - params: The pagination parameters for the request.
// Returns:
// - api.GetExecutionWorkersRes: The response containing the list of execution workers.
// - error: Any error that occurred during the request.
func (s *OdinServer) GetExecutionWorkers(ctx context.Context, params api.GetExecutionWorkersParams) (api.GetExecutionWorkersRes, error) {
	workersDB, err := s.queries.GetAllWorkers(ctx, db.GetAllWorkersParams{
		Limit:  params.PageSize.Value,
		Offset: params.PageSize.Value * params.Page.Value,
	})
	if err != nil {
		return &api.GetExecutionWorkersInternalServerError{
			Message: fmt.Sprintf("Failed to get workers: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalWorkers(ctx)
	if err != nil {
		return &api.GetExecutionWorkersInternalServerError{
			Message: fmt.Sprintf("Failed to get total workers: %v", err),
		}, nil
	}
	var workers []api.ExecutionWorker
	for _, worker := range workersDB {
		workers = append(workers, api.ExecutionWorker{
			ID:        worker.ID,
			Name:      worker.Name,
			CreatedAt: worker.CreatedAt.Time,
			Status:    "status",
		})
	}
	resp := api.GetExecutionWorkersOK{
		Workers: workers,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
}

// DeleteJob deletes a job by its ID.
//
// Parameters:
// - ctx: context for the request.
// - params: DeleteJobParams containing the JobId to delete.
// Returns:
// - api.DeleteJobRes: result of the delete operation.
// - error: error if the delete operation fails.
func (s *OdinServer) DeleteJob(ctx context.Context, params api.DeleteJobParams) (api.DeleteJobRes, error) {
	err := s.queries.DeleteJob(ctx, params.JobId)
	if err != nil {
		return &api.DeleteJobBadRequest{
			Message: fmt.Sprintf("Failed to delete execution: %v", err),
		}, nil
	}
	return &api.DeleteJobOK{}, nil
}
