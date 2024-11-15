package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5"
)

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
			s.logger.Error().Stack().Err(err).Msg("Failed to execute")
			return &api.ExecuteInternalServerError{
				Message: fmt.Sprintf("Failed to execute: %v", err),
			}, nil
		}
	}
	return &api.ExecuteOK{ExecutionId: execId}, nil
}

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
			ExecutionId: int64(execRes.ExecRequestID.Int32),
			Logs:        execRes.Logs,
			Script:      execRes.Code,
			Flake:       execRes.Flake,
			Args:        execRes.Args.String,
			StartedAt:   execRes.StartedAt.Time,
			FinishedAt:  execRes.FinishedAt.Time,
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
			ExecutionId: int64(execRes.ExecRequestID.Int32),
			Logs:        execRes.Logs,
			Script:      execRes.Code,
			Flake:       execRes.Flake,
			Args:        execRes.Args.String,
			StartedAt:   execRes.StartedAt.Time,
			FinishedAt:  execRes.FinishedAt.Time,
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
			Script:      exec.Code,
			Flake:       exec.Flake,
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

func (s *OdinServer) DeleteJob(ctx context.Context, params api.DeleteJobParams) (api.DeleteJobRes, error) {
	err := s.queries.DeleteJob(ctx, params.JobId)
	if err != nil {
		return &api.DeleteJobBadRequest{
			Message: fmt.Sprintf("Failed to delete execution: %v", err),
		}, nil
	}
	return &api.DeleteJobOK{}, nil
}

func (s *OdinServer) CancelJob(ctx context.Context, params api.CancelJobParams) (api.CancelJobRes, error) {
	err := s.queries.CancelJob(ctx, params.JobId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &api.CancelJobBadRequest{
				Message: fmt.Sprintf("Job not found: %v", err),
			}, nil
		}
		return &api.CancelJobInternalServerError{
			Message: fmt.Sprintf("Failed to cancel job: %v", err),
		}, nil
	}
	return &api.CancelJobOK{}, nil
}
