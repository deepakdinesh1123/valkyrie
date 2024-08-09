// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// DeleteExecution implements deleteExecution operation.
//
// Delete execution.
//
// DELETE /executions/{executionId}/
func (UnimplementedHandler) DeleteExecution(ctx context.Context, params DeleteExecutionParams) (r DeleteExecutionRes, _ error) {
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

// GetExecutionConfig implements getExecutionConfig operation.
//
// Get execution config.
//
// GET /execution/config/
func (UnimplementedHandler) GetExecutionConfig(ctx context.Context) (r GetExecutionConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionResult implements getExecutionResult operation.
//
// Get execution result.
//
// GET /executions/{executionId}/
func (UnimplementedHandler) GetExecutionResult(ctx context.Context, params GetExecutionResultParams) (r GetExecutionResultRes, _ error) {
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

// GetExecutions implements getExecutions operation.
//
// Get all executions.
//
// GET /executions/
func (UnimplementedHandler) GetExecutions(ctx context.Context, params GetExecutionsParams) (r GetExecutionsRes, _ error) {
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
