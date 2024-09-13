package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
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
	return &api.ExecuteOK{ExecutionId: execId, Events: fmt.Sprintf("/executions/%d/events", execId)}, nil
}

func (s *OdinServer) ExecuteSSE(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	execIdStr := req.PathValue("executionId")
	execId, err := strconv.ParseInt(execIdStr, 10, 64)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to get executionId")
		http.Error(w, "Failed to get executionId", http.StatusBadRequest)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	var flusher http.Flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		s.logger.Error().Stack().Msg("Failed to get flusher")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Int64("executionId", execId).Msg("Client disconnected")
			return
		default:
			job, err := s.queries.GetJob(ctx, execId)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to get job")
				fmt.Fprintf(w, "data: error\n\n")
				flusher.Flush()
				return
			}

			s.logger.Debug().Str("status", job.Status).Msg("Status fetched")

			switch job.Status {
			case "completed":
				res, err := s.queries.GetExecutionResultsByID(ctx, db.GetExecutionResultsByIDParams{
					JobID:  execId,
					Limit:  1,
					Offset: 0,
				})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution results")
					fmt.Fprintf(w, "data: error %s\n\n", err)
					flusher.Flush()
					return
				}
				fmt.Fprintf(w, "data: completed: %d\n\n", execId)
				flusher.Flush()
				fmt.Fprintf(w, "data: %s\n\n", res[0].ExecLogs)
				return
			case "failed":
				fmt.Fprintf(w, "data: failed: %d\n\n", execId)
				flusher.Flush()
				return
			case "pending":
				fmt.Fprintf(w, "data: pending: %d\n\n", execId)
				flusher.Flush()
			case "scheduled":
				fmt.Fprintf(w, "data: scheduled: %d\n\n", execId)
				flusher.Flush()
			default:
				s.logger.Warn().Str("status", job.Status).Msg("Unknown status")
			}

			time.Sleep(3 * time.Second)
		}
	}
}

func (s *OdinServer) ExecuteWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to accept websocket")
		http.Error(w, "Failed to accept websocket", http.StatusInternalServerError)
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "Closing connection")
	var execReq api.ExecutionRequest
	err = wsjson.Read(ctx, c, &execReq)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to read request")
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	execId, err := s.executionService.AddJob(ctx, &execReq)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to execute")
		http.Error(w, "Failed to execute", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Int64("executionId", execId).Msg("Client disconnected")
			return
		default:
			job, err := s.queries.GetJob(ctx, execId)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to get job")
				return
			}

			s.logger.Debug().Str("status", job.Status).Msg("Status fetched")

			switch job.Status {
			case "completed":
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("completed: %d\n\n", execId)))
				res, err := s.queries.GetExecutionResultsByID(ctx, db.GetExecutionResultsByIDParams{
					JobID:  execId,
					Limit:  1,
					Offset: 0,
				})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution results")
					c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("failed: %s\n\n", err)))
					return
				}
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("data: %s\n\n", res[0].ExecLogs)))
				return
			case "failed":
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("failed: %d\n\n", execId)))
				return
			case "pending":
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("pending: %d\n\n", execId)))
			case "scheduled":
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("scheduled: %d\n\n", execId)))
			default:
				s.logger.Warn().Str("status", job.Status).Msg("Unknown status")
			}
			time.Sleep(3 * time.Second)
		}
	}
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
			ExecLogs:    execRes.ExecLogs,
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
			ExecLogs:    execRes.ExecLogs,
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
