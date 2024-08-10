// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// DeleteJob implements deleteJob operation.
//
// Delete job.
//
// DELETE /executions/{JobId}/
func (UnimplementedHandler) DeleteJob(ctx context.Context, params DeleteJobParams) (r DeleteJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Execute implements execute operation.
//
// Execute a script.
//
// POST /executions/execute/
func (UnimplementedHandler) Execute(ctx context.Context, req *ExecutionRequest) (r ExecuteRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetAllExecutionResults implements getAllExecutionResults operation.
//
// Get all execution results.
//
// GET /executions/results/
func (UnimplementedHandler) GetAllExecutionResults(ctx context.Context, params GetAllExecutionResultsParams) (r GetAllExecutionResultsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetAllExecutions implements getAllExecutions operation.
//
// Get all executions.
//
// GET /executions/
func (UnimplementedHandler) GetAllExecutions(ctx context.Context, params GetAllExecutionsParams) (r GetAllExecutionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionConfig implements getExecutionConfig operation.
//
// Get execution config.
//
// GET /execution/config/
func (UnimplementedHandler) GetExecutionConfig(ctx context.Context) (r GetExecutionConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionResultsById implements getExecutionResultsById operation.
//
// Get execution result.
//
// GET /executions/{JobId}/
func (UnimplementedHandler) GetExecutionResultsById(ctx context.Context, params GetExecutionResultsByIdParams) (r GetExecutionResultsByIdRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionWorkers implements getExecutionWorkers operation.
//
// Get all execution workers.
//
// GET /executions/workers
func (UnimplementedHandler) GetExecutionWorkers(ctx context.Context, params GetExecutionWorkersParams) (r GetExecutionWorkersRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetVersion implements getVersion operation.
//
// Get version.
//
// GET /version/
func (UnimplementedHandler) GetVersion(ctx context.Context) (r GetVersionRes, _ error) {
	return r, ht.ErrNotImplemented
}
