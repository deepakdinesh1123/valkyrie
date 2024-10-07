package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *OdinServer) Execute(ctx context.Context, req *api.ExecutionRequest, params api.ExecuteParams) (api.ExecuteRes, error) {
	jobId, err := s.executionService.AddJob(ctx, req)
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
	return &api.ExecuteOK{JobId: jobId, Events: fmt.Sprintf("/executions/%d/events", jobId)}, nil
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

			s.logger.Debug().Str("status", job.CurrentState).Msg("CurrentState fetched")

			switch job.CurrentState {
			case "completed":
				res, err := s.queries.GetExecutionsForJob(ctx, db.GetExecutionsForJobParams{
					JobID:  pgtype.Int8{Int64: execId, Valid: true},
					Limit:  1,
					Offset: 0,
				})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution executions")
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
				s.logger.Warn().Str("status", job.CurrentState).Msg("Unknown status")
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

			s.logger.Debug().Str("status", job.CurrentState).Msg("CurrentState fetched")

			switch job.CurrentState {
			case "completed":
				c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("completed: %d\n\n", execId)))
				res, err := s.queries.GetExecutionsForJob(ctx, db.GetExecutionsForJobParams{
					JobID:  pgtype.Int8{Int64: execId, Valid: true},
					Limit:  1,
					Offset: 0,
				})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution executions")
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
				s.logger.Warn().Str("status", job.CurrentState).Msg("Unknown status")
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *OdinServer) GetAllExecutions(ctx context.Context, params api.GetAllExecutionsParams) (api.GetAllExecutionsRes, error) {
	execResDB, err := s.queries.GetAllExecutions(ctx, db.GetAllExecutionsParams{
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetAllExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get execution executions: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalExecutions(ctx)
	if err != nil {
		return &api.GetAllExecutionsInternalServerError{
			Message: fmt.Sprintf("Failed to get total execution executions: %v", err),
		}, nil
	}
	var executions []api.ExecutionResult
	for _, execution := range execResDB {
		executions = append(executions, api.ExecutionResult{
			JobId:      int64(execution.ExecRequestID.Int32),
			ExecLogs:   execution.ExecLogs,
			NixLogs:    api.OptString{Value: execution.NixLogs.String, Set: true},
			Script:     execution.Code,
			Flake:      execution.Flake,
			Args:       execution.Args.String,
			StartedAt:  execution.StartedAt.Time,
			FinishedAt: execution.FinishedAt.Time,
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

func (s *OdinServer) GetExecutionsForJob(ctx context.Context, params api.GetExecutionsForJobParams) (api.GetExecutionsForJobRes, error) {
	execRes, err := s.queries.GetExecutionsForJob(ctx, db.GetExecutionsForJobParams{
		JobID:  pgtype.Int8{Int64: params.JobId, Valid: true},
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetExecutionsForJobInternalServerError{
			Message: fmt.Sprintf("Failed to get execution executions: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalExecutionsForJob(ctx, pgtype.Int8{Int64: params.JobId, Valid: true})
	if err != nil {
		return &api.GetExecutionsForJobInternalServerError{
			Message: fmt.Sprintf("Failed to get total execution executions: %v", err),
		}, nil
	}
	var executions []api.ExecutionResult
	for _, execution := range execRes {
		executions = append(executions, api.ExecutionResult{
			JobId:      int64(execution.ExecRequestID.Int32),
			ExecId:     int64(execution.ID),
			ExecLogs:   execution.ExecLogs,
			NixLogs:    api.OptString{Value: execution.NixLogs.String, Set: true},
			Script:     execution.Code,
			Flake:      execution.Flake,
			Args:       execution.Args.String,
			StartedAt:  execution.StartedAt.Time,
			FinishedAt: execution.FinishedAt.Time,
		})
	}
	resp := api.GetExecutionsForJobOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
}

func (s *OdinServer) GetAllExecutionJobs(ctx context.Context, params api.GetAllExecutionJobsParams) (api.GetAllExecutionJobsRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.GetAllExecutionJobsForbidden{}, nil
	}
	executionsDB, err := s.queries.GetAllJobs(ctx, db.GetAllJobsParams{
		Limit:  params.PageSize.Value,
		Offset: params.Page.Value * params.PageSize.Value,
	})
	if err != nil {
		return &api.GetAllExecutionJobsInternalServerError{
			Message: fmt.Sprintf("Failed to get executions: %v", err),
		}, nil
	}
	total, err := s.queries.GetTotalJobs(ctx)
	if err != nil {
		return &api.GetAllExecutionJobsInternalServerError{
			Message: fmt.Sprintf("Failed to get total executions: %v", err),
		}, nil
	}
	var jobs []api.Job
	for _, job := range executionsDB {
		jobs = append(jobs, api.Job{
			JobId:     job.JobID,
			Script:    job.Code,
			Flake:     job.Flake,
			CreatedAt: job.UpdatedAt.Time,
			StartedAt: api.OptDateTime{Value: job.StartedAt.Time, Set: true},
			UpdatedAt: api.OptDateTime{Value: job.UpdatedAt.Time, Set: true},
		})
	}
	resp := api.GetAllExecutionJobsOK{
		Jobs: jobs,
		Pagination: api.PaginationResponse{
			Total: total,
		},
	}
	return &resp, nil
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

func (s *OdinServer) DeleteExecutionJob(ctx context.Context, params api.DeleteExecutionJobParams) (api.DeleteExecutionJobRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.DeleteExecutionJobForbidden{}, nil
	}
	state, err := s.queries.GetJobState(ctx, params.JobId)
	if err != nil {
		return &api.DeleteExecutionJobInternalServerError{
			Message: fmt.Sprintf("Failed to get job state: %v", err),
		}, nil
	}
	if state == "scheduled" {
		return &api.DeleteExecutionJobBadRequest{
			Message: fmt.Sprintf("Execution job with id %d is currently scheduled", params.JobId),
		}, nil
	}

	id, err := s.queries.DeleteJob(ctx, params.JobId)
	if id == 0 {
		return &api.DeleteExecutionJobBadRequest{
			Message: fmt.Sprintf("Execution job with id %d not found: %v", params.JobId, err),
		}, nil
	}
	if err != nil {
		return &api.DeleteExecutionJobInternalServerError{
			Message: fmt.Sprintf("Failed to delete execution: %v", err),
		}, nil
	}
	return &api.DeleteExecutionJobOK{}, nil
}

func (s *OdinServer) CancelExecutionJob(ctx context.Context, params api.CancelExecutionJobParams) (api.CancelExecutionJobRes, error) {
	err := s.queries.CancelJob(ctx, params.JobId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &api.CancelExecutionJobBadRequest{
				Message: fmt.Sprintf("Job not found: %v", err),
			}, nil
		}
		return &api.CancelExecutionJobInternalServerError{
			Message: fmt.Sprintf("Failed to cancel job: %v", err),
		}, nil
	}
	return &api.CancelExecutionJobOK{
		Message: fmt.Sprintf("Execution Job with id %d canceled", params.JobId),
	}, nil
}

func (s *OdinServer) GetExecutionResultById(ctx context.Context, params api.GetExecutionResultByIdParams) (api.GetExecutionResultByIdRes, error) {
	execution, err := s.queries.GetExecution(ctx, params.ExecId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &api.GetExecutionResultByIdBadRequest{
				Message: fmt.Sprintf("Execution execution with id %d not found: %v", params.ExecId, err),
			}, nil
		}
		return &api.GetExecutionResultByIdInternalServerError{
			Message: fmt.Sprintf("Failed to get execution execution: %v", err),
		}, nil
	}

	return &api.ExecutionResult{
		JobId:      int64(execution.ExecRequestID.Int32),
		Script:     execution.Code,
		Flake:      execution.Flake,
		CreatedAt:  execution.CreatedAt.Time,
		StartedAt:  execution.StartedAt.Time,
		FinishedAt: execution.FinishedAt.Time,
		ExecId:     int64(execution.ID),
		ExecLogs:   execution.ExecLogs,
		Args:       execution.Args.String,
	}, nil
}

func (s *OdinServer) GetExecutionJobById(ctx context.Context, params api.GetExecutionJobByIdParams) (api.GetExecutionJobByIdRes, error) {
	job, err := s.queries.GetJob(ctx, params.JobId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &api.GetExecutionJobByIdBadRequest{
				Message: fmt.Sprintf("Execution job with id %d not found: %v", params.JobId, err),
			}, nil
		}
		return &api.GetExecutionJobByIdInternalServerError{
			Message: fmt.Sprintf("Failed to get execution job: %v", err),
		}, nil
	}
	return &api.Job{
		JobId:     int64(job.ID),
		Script:    job.Code,
		Flake:     job.Flake,
		CreatedAt: job.CreatedAt.Time,
		StartedAt: api.OptDateTime{Value: job.StartedAt.Time, Set: true},
		UpdatedAt: api.OptDateTime{Value: job.UpdatedAt.Time, Set: true},
	}, nil
}

func (s *OdinServer) GetExecutionConfig(ctx context.Context, params api.GetExecutionConfigParams) (api.GetExecutionConfigRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.GetExecutionConfigForbidden{}, nil
	}

	return &api.ExecutionConfig{
		ODINWORKERPROVIDER:    s.envConfig.ODIN_CONTAINER_ENGINE,
		ODINWORKERCONCURRENCY: int32(s.envConfig.ODIN_WORKER_CONCURRENCY),
		ODINWORKERBUFFERSIZE:  int32(s.envConfig.ODIN_WORKER_BUFFER_SIZE),
		ODINWORKERTASKTIMEOUT: s.envConfig.ODIN_WORKER_TASK_TIMEOUT,
		ODINWORKERPOLLFREQ:    s.envConfig.ODIN_WORKER_POLL_FREQ,
		ODINWORKERRUNTIME:     s.envConfig.ODIN_WORKER_RUNTIME,
		ODINLOGLEVEL:          s.envConfig.ODIN_LOG_LEVEL,
	}, nil
}

func (s *OdinServer) GetAllLanguages(ctx context.Context, params api.GetAllLanguagesParams) (api.GetAllLanguagesRes, error) {
	var languages []string
	for lang := range config.Languages {
		languages = append(languages, lang)
	}
	return &api.GetAllLanguagesOK{
		Languages: languages,
	}, nil
}
