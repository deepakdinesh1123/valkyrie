package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/token"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExecutionMessage struct {
	Status   string      `json:"status"`
	JobID    int64       `json:"jobId"`
	Logs     interface{} `json:"logs,omitempty"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
}

func (s *ValkyrieServer) Execute(ctx context.Context, req *api.ExecutionRequest, params api.ExecuteParams) (api.ExecuteRes, error) {
	if !s.envConfig.ENABLE_EXECUTION {
		return &api.ExecuteBadRequest{
			Message: "Execution is not enabled, please ask the admin to enable it",
		}, nil
	}
	if _, err := s.executionService.CheckExecRequest(ctx, req); err != nil {
		return &api.ExecuteBadRequest{
			Message: fmt.Sprintf("%s", err),
		}, nil
	}
	jobId, err := s.executionService.AddJob(ctx, req)
	if err != nil {
		return &api.ExecuteInternalServerError{
			Message: fmt.Sprintf("Error adding execution job: %s", err),
		}, nil
	}

	tkn, err := token.GenerateValkyrieToken(jwt.MapClaims{
		"jobId": strconv.Itoa(int(jobId)),
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	}, s.envConfig)
	if err != nil {
		return &api.ExecuteInternalServerError{
			Message: fmt.Sprintf("Failed to generate token: %v", err),
		}, nil
	}

	return &api.ExecuteOK{JobId: jobId, Events: fmt.Sprintf("/executions/%d/events", jobId), Websocket: fmt.Sprintf("/executions/%d/ws/?token=%s", jobId, tkn)}, nil
}

func (s *ValkyrieServer) ExecuteSSE(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	jobIDStr := chi.URLParam(req, "jobId")

	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to get executionId")
		http.Error(w, "Failed to get executionId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.logger.Error().Stack().Msg("Failed to get flusher")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-ctx.Done():
			job, err := s.queries.GetExecutionJob(context.TODO(), jobID)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to get job status")
				return
			}
			if job.CurrentState == "pending" {
				if _, err := s.queries.DeleteJob(context.TODO(), jobID); err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to delete pending job on disconnect")
				} else {
					s.logger.Info().Int64("executionId", jobID).Msg("Deleted pending job due to client disconnection")
				}
			}
			return
		default:
			job, err := s.queries.GetExecutionJob(ctx, jobID)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to get job")
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status:   "error",
					JobID:    jobID,
					ErrorMsg: fmt.Sprintf("Failed to get job: %s", err),
				})
				return
			}
			switch job.CurrentState {
			case "completed":
				res, err := s.queries.GetLatestExecution(ctx, pgtype.Int8{Int64: jobID, Valid: true})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution logs")
					sendExecutionMessage(w, flusher, ExecutionMessage{
						Status:   "error",
						JobID:    jobID,
						ErrorMsg: fmt.Sprintf("failed to get execution logs: %s", err),
					})
					return
				}
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status: "completed",
					JobID:  jobID,
					Logs:   res.ExecLogs,
				})
				return

			case "failed":
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status: "failed",
					JobID:  jobID,
				})
				return

			case "pending":
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status: "pending",
					JobID:  jobID,
				})
			case "scheduled":
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status: "scheduled",
					JobID:  jobID,
				})
			case "cancelled":
				sendExecutionMessage(w, flusher, ExecutionMessage{
					Status: "canceled",
					JobID:  jobID,
				})
			default:
				s.logger.Warn().Str("status", job.CurrentState).Msg("Unknown status")
			}

			time.Sleep(30 * time.Millisecond)
		}
	}
}

func (s *ValkyrieServer) ExecuteWebSocket(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	jobIDStr := chi.URLParam(req, "jobId")

	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to get executionId")
		http.Error(w, "Failed to get executionId", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("Failed to accept websocket connection")
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Connection closed")

	wsCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	disconnectChan := make(chan struct{})

	go func() {
		defer close(disconnectChan)
		for {
			_, _, err := conn.Read(wsCtx)
			if err != nil {
				s.logger.Info().Int64("jobID", jobID).Err(err).Msg("Client disconnected")
				cancel()
				return
			}
		}
	}()

	go func() {
		select {
		case <-wsCtx.Done():
		case <-disconnectChan:
		}

		job, err := s.queries.GetExecutionJob(context.Background(), jobID)
		if err != nil {
			s.logger.Error().Stack().Err(err).Msg("Failed to get job status for cleanup")
			return
		}

		if job.CurrentState == "pending" || job.CurrentState == "scheduled" {
			if _, err := s.queries.DeleteJob(context.Background(), jobID); err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to delete job on disconnect")
			} else {
				s.logger.Info().Int64("jobID", jobID).Str("state", job.CurrentState).
					Msg("Deleted job due to client disconnection")
			}
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-wsCtx.Done():
			s.logger.Info().Int64("jobID", jobID).Msg("WebSocket context cancelled")
			return

		case <-disconnectChan:
			s.logger.Info().Int64("jobID", jobID).Msg("Client disconnected via read error")
			return

		case <-ticker.C:
			job, err := s.queries.GetExecutionJob(wsCtx, jobID)
			if err != nil {
				s.logger.Error().Stack().Err(err).Msg("Failed to get job")
				msg := ExecutionMessage{
					Status:   "error",
					JobID:    jobID,
					ErrorMsg: fmt.Sprintf("Failed to get job: %s", err),
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send error message")
					return
				}
				conn.Close(websocket.StatusInternalError, "Failed to get job")
				return
			}

			switch job.CurrentState {
			case "completed":
				res, err := s.queries.GetLatestExecution(wsCtx, pgtype.Int8{Int64: jobID, Valid: true})
				if err != nil {
					s.logger.Error().Stack().Err(err).Msg("Failed to get execution logs")
					msg := ExecutionMessage{
						Status:   "error",
						JobID:    jobID,
						ErrorMsg: fmt.Sprintf("failed to get execution logs: %s", err),
					}
					if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
						s.logger.Error().Err(sendErr).Msg("Failed to send error message")
					}
					conn.Close(websocket.StatusInternalError, "Failed to get execution logs")
					return
				}
				msg := ExecutionMessage{
					Status: "completed",
					JobID:  jobID,
					Logs:   res.ExecLogs,
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send completion message")
				}
				conn.Close(websocket.StatusNormalClosure, "Job completed")
				return

			case "failed":
				msg := ExecutionMessage{
					Status: "failed",
					JobID:  jobID,
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send failure message")
				}
				conn.Close(websocket.StatusNormalClosure, "Job failed")
				return

			case "pending":
				msg := ExecutionMessage{
					Status: "pending",
					JobID:  jobID,
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send pending message")
					return
				}

			case "scheduled":
				msg := ExecutionMessage{
					Status: "scheduled",
					JobID:  jobID,
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send scheduled message")
					return
				}

			case "cancelled":
				msg := ExecutionMessage{
					Status: "canceled",
					JobID:  jobID,
				}
				if sendErr := sendWebSocketMessage(wsCtx, conn, msg); sendErr != nil {
					s.logger.Error().Err(sendErr).Msg("Failed to send cancelled message")
				}
				conn.Close(websocket.StatusNormalClosure, "Job canceled")
				return

			default:
				s.logger.Warn().Str("status", job.CurrentState).Msg("Unknown job status")
			}
		}
	}
}

// sendWebSocketMessage sends a message over the WebSocket connection
func sendWebSocketMessage(ctx context.Context, conn *websocket.Conn, message ExecutionMessage) error {

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		log.Printf("Failed to write WebSocket message: %v", err)
		return err
	}

	return nil
}

func (s *ValkyrieServer) GetAllExecutions(ctx context.Context, params api.GetAllExecutionsParams) (api.GetAllExecutionsRes, error) {

	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.GetAllExecutionsForbidden{}, nil
	}

	execResDB, err := s.queries.GetAllExecutions(ctx, db.GetAllExecutionsParams{
		Limit:  params.Limit.Value,
		ExecID: params.Cursor.Value,
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
			StartedAt:  execution.StartedAt.Time,
			FinishedAt: execution.FinishedAt.Time,
		})
	}
	var cursor int64
	if len(executions) > 0 {
		cursor = executions[len(executions)-1].JobId
	} else {
		cursor = 0
	}
	resp := api.GetAllExecutionsOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total:  total,
			Limit:  params.Limit.Value,
			Cursor: cursor,
		},
	}
	return &resp, nil
}

func (s *ValkyrieServer) GetExecutionsForJob(ctx context.Context, params api.GetExecutionsForJobParams) (api.GetExecutionsForJobRes, error) {
	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.GetExecutionsForJobForbidden{}, nil
	}

	execRes, err := s.queries.GetExecutionsForJob(ctx, db.GetExecutionsForJobParams{
		JobID:  pgtype.Int8{Int64: params.JobId, Valid: true},
		Limit:  params.Limit.Value,
		ExecID: params.Cursor.Value,
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
			StartedAt:  execution.StartedAt.Time,
			FinishedAt: execution.FinishedAt.Time,
		})
	}
	var cursor int64
	cursor = 0
	if len(executions) > 0 {
		cursor = executions[len(executions)-1].ExecId
	}
	resp := api.GetExecutionsForJobOK{
		Executions: executions,
		Pagination: api.PaginationResponse{
			Total:  total,
			Limit:  params.Limit.Value,
			Cursor: cursor,
		},
	}
	return &resp, nil
}

func (s *ValkyrieServer) GetAllExecutionJobs(ctx context.Context, params api.GetAllExecutionJobsParams) (api.GetAllExecutionJobsRes, error) {

	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.GetAllExecutionJobsForbidden{}, nil
	}

	executionsDB, err := s.queries.GetAllExecutionJobs(ctx, db.GetAllExecutionJobsParams{
		Limit: params.Limit.Value,
		JobID: params.Cursor.Value,
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
			Flake:     job.Flake,
			CreatedAt: job.UpdatedAt.Time,
			StartedAt: api.OptDateTime{Value: job.StartedAt.Time, Set: true},
			UpdatedAt: api.OptDateTime{Value: job.UpdatedAt.Time, Set: true},
		})
	}
	var cursor int64
	cursor = 0
	if len(jobs) > 0 {
		cursor = jobs[len(jobs)-1].JobId
	}
	resp := api.GetAllExecutionJobsOK{
		Jobs: jobs,
		Pagination: api.PaginationResponse{
			Total:  total,
			Limit:  params.Limit.Value,
			Cursor: cursor,
		},
	}
	return &resp, nil
}

func (s *ValkyrieServer) DeleteExecutionJob(ctx context.Context, params api.DeleteExecutionJobParams) (api.DeleteExecutionJobRes, error) {
	user := ctx.Value(config.RoleKey).(string)

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

func (s *ValkyrieServer) CancelExecutionJob(ctx context.Context, params api.CancelExecutionJobParams) (api.CancelExecutionJobRes, error) {
	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.CancelExecutionJobForbidden{}, nil
	}

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

func (s *ValkyrieServer) GetExecutionResultById(ctx context.Context, params api.GetExecutionResultByIdParams) (api.GetExecutionResultByIdRes, error) {
	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.GetExecutionResultByIdForbidden{}, nil
	}

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
		CreatedAt:  execution.CreatedAt.Time,
		StartedAt:  execution.StartedAt.Time,
		FinishedAt: execution.FinishedAt.Time,
		ExecId:     int64(execution.ID),
		ExecLogs:   execution.ExecLogs,
	}, nil
}

func (s *ValkyrieServer) GetExecutionJobById(ctx context.Context, params api.GetExecutionJobByIdParams) (api.GetExecutionJobByIdRes, error) {
	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.GetExecutionJobByIdForbidden{}, nil
	}
	job, err := s.queries.GetExecutionJob(ctx, params.JobId)
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
		CreatedAt: job.CreatedAt.Time,
		StartedAt: api.OptDateTime{Value: job.StartedAt.Time, Set: true},
		UpdatedAt: api.OptDateTime{Value: job.UpdatedAt.Time, Set: true},
	}, nil
}

func (s *ValkyrieServer) GetExecutionConfig(ctx context.Context) (api.GetExecutionConfigRes, error) {
	user := ctx.Value(config.RoleKey).(string)

	if user != "admin" {
		return &api.GetExecutionConfigForbidden{}, nil
	}

	return &api.ExecutionConfig{
		WORKERPROVIDER:    s.envConfig.RUNTIME,
		WORKERCONCURRENCY: int32(s.envConfig.WORKER_CONCURRENCY),
		WORKERPOLLFREQ:    s.envConfig.WORKER_POLL_FREQ,
		WORKERRUNTIME:     s.envConfig.CONTAINER_RUNTIME,
		LOGLEVEL:          s.envConfig.LOG_LEVEL,
	}, nil
}

func sendExecutionMessage(w http.ResponseWriter, flusher http.Flusher, message ExecutionMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintf(w, "data: {\"status\":\"error\",\"error\":\"Failed to marshal JSON\"}\n\n")
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}

// FlakeJobIdGet implements api.Handler.
func (s *ValkyrieServer) FetchFlake(ctx context.Context, params api.FetchFlakeParams) (api.FetchFlakeRes, error) {
	role := ctx.Value(config.RoleKey).(string)

	if role != "admin" {
		return &api.FetchFlakeForbidden{}, nil
	}
	flake, err := s.queries.GetFlake(ctx, params.JobId)
	if err != nil {
		return &api.FetchFlakeInternalServerError{}, err
	}
	return &api.FetchFlakeOK{
		Flake: flake,
	}, nil
}
