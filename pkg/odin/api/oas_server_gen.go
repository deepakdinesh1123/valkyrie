// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CancelExecutionJob implements cancelExecutionJob operation.
	//
	// Cancel Execution Job.
	//
	// PUT /executions/jobs/{JobId}/
	CancelExecutionJob(ctx context.Context, params CancelExecutionJobParams) (CancelExecutionJobRes, error)
	// DeleteExecutionJob implements deleteExecutionJob operation.
	//
	// Delete execution job.
	//
	// DELETE /executions/jobs/{JobId}/
	DeleteExecutionJob(ctx context.Context, params DeleteExecutionJobParams) (DeleteExecutionJobRes, error)
	// DeleteExecutionWorker implements deleteExecutionWorker operation.
	//
	// Delete execution worker.
	//
	// DELETE /executions/workers/{workerId}/
	DeleteExecutionWorker(ctx context.Context, params DeleteExecutionWorkerParams) (DeleteExecutionWorkerRes, error)
	// Execute implements execute operation.
	//
	// Execute a script.
	//
	// POST /executions/execute/
	Execute(ctx context.Context, req *ExecutionRequest) (ExecuteRes, error)
	// GenerateUserToken implements generateUserToken operation.
	//
	// Generate user token.
	//
	// GET /user/token/
	GenerateUserToken(ctx context.Context) (GenerateUserTokenRes, error)
	// GetAllExecutionJobs implements getAllExecutionJobs operation.
	//
	// Get all execution jobs.
	//
	// GET /jobs/execution/
	GetAllExecutionJobs(ctx context.Context, params GetAllExecutionJobsParams) (GetAllExecutionJobsRes, error)
	// GetAllExecutions implements getAllExecutions operation.
	//
	// Get all executions.
	//
	// GET /executions/
	GetAllExecutions(ctx context.Context, params GetAllExecutionsParams) (GetAllExecutionsRes, error)
	// GetExecutionConfig implements getExecutionConfig operation.
	//
	// Get execution config.
	//
	// GET /execution/config/
	GetExecutionConfig(ctx context.Context) (GetExecutionConfigRes, error)
	// GetExecutionJobById implements getExecutionJobById operation.
	//
	// Get execution job.
	//
	// GET /executions/jobs/{JobId}/
	GetExecutionJobById(ctx context.Context, params GetExecutionJobByIdParams) (GetExecutionJobByIdRes, error)
	// GetExecutionResultById implements getExecutionResultById operation.
	//
	// Get execution result by id.
	//
	// GET /executions/{execId}/
	GetExecutionResultById(ctx context.Context, params GetExecutionResultByIdParams) (GetExecutionResultByIdRes, error)
	// GetExecutionWorkers implements getExecutionWorkers operation.
	//
	// Get all execution workers.
	//
	// GET /executions/workers
	GetExecutionWorkers(ctx context.Context, params GetExecutionWorkersParams) (GetExecutionWorkersRes, error)
	// GetExecutionsForJob implements getExecutionsForJob operation.
	//
	// Get executions of given job.
	//
	// GET /jobs/{JobId}/executions/
	GetExecutionsForJob(ctx context.Context, params GetExecutionsForJobParams) (GetExecutionsForJobRes, error)
	// GetToken implements getToken operation.
	//
	// Get token.
	//
	// POST /admin/token/
	GetToken(ctx context.Context, req *GetTokenReq) (GetTokenRes, error)
	// GetVersion implements getVersion operation.
	//
	// Get version.
	//
	// GET /version/
	GetVersion(ctx context.Context) (GetVersionRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
