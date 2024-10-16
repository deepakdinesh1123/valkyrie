// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CancelExecutionJob implements cancelExecutionJob operation.
//
// Cancel Execution Job.
//
// PUT /executions/jobs/{JobId}
func (UnimplementedHandler) CancelExecutionJob(ctx context.Context, params CancelExecutionJobParams) (r CancelExecutionJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteExecutionJob implements deleteExecutionJob operation.
//
// Delete execution job.
//
// DELETE /executions/jobs/{JobId}
func (UnimplementedHandler) DeleteExecutionJob(ctx context.Context, params DeleteExecutionJobParams) (r DeleteExecutionJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteExecutionWorker implements deleteExecutionWorker operation.
//
// Delete execution worker.
//
// DELETE /executions/workers/{workerId}
func (UnimplementedHandler) DeleteExecutionWorker(ctx context.Context, params DeleteExecutionWorkerParams) (r DeleteExecutionWorkerRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Execute implements execute operation.
//
// Execute a script.
//
// POST /executions/execute
func (UnimplementedHandler) Execute(ctx context.Context, req *ExecutionRequest, params ExecuteParams) (r ExecuteRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetAllExecutionJobs implements getAllExecutionJobs operation.
//
// Get all execution jobs.
//
// GET /jobs/execution
func (UnimplementedHandler) GetAllExecutionJobs(ctx context.Context, params GetAllExecutionJobsParams) (r GetAllExecutionJobsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetAllExecutions implements getAllExecutions operation.
//
// Get all executions.
//
// GET /executions
func (UnimplementedHandler) GetAllExecutions(ctx context.Context, params GetAllExecutionsParams) (r GetAllExecutionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetAllLanguages implements getAllLanguages operation.
//
// Get all languages.
//
// GET /languages
func (UnimplementedHandler) GetAllLanguages(ctx context.Context, params GetAllLanguagesParams) (r GetAllLanguagesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionConfig implements getExecutionConfig operation.
//
// Get execution config.
//
// GET /execution/config
func (UnimplementedHandler) GetExecutionConfig(ctx context.Context, params GetExecutionConfigParams) (r GetExecutionConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionJobById implements getExecutionJobById operation.
//
// Get execution job.
//
// GET /executions/jobs/{JobId}
func (UnimplementedHandler) GetExecutionJobById(ctx context.Context, params GetExecutionJobByIdParams) (r GetExecutionJobByIdRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetExecutionResultById implements getExecutionResultById operation.
//
// Get execution result by id.
//
// GET /executions/{execId}
func (UnimplementedHandler) GetExecutionResultById(ctx context.Context, params GetExecutionResultByIdParams) (r GetExecutionResultByIdRes, _ error) {
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

// GetExecutionsForJob implements getExecutionsForJob operation.
//
// Get executions of given job.
//
// GET /jobs/{JobId}/executions
func (UnimplementedHandler) GetExecutionsForJob(ctx context.Context, params GetExecutionsForJobParams) (r GetExecutionsForJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetVersion implements getVersion operation.
//
// Get version.
//
// GET /version
func (UnimplementedHandler) GetVersion(ctx context.Context, params GetVersionParams) (r GetVersionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PackagesExist implements PackagesExist operation.
//
// Verify the package list is available for the language version while switching between language
// versions.
//
// POST /packages/exist/
func (UnimplementedHandler) PackagesExist(ctx context.Context, req *PackageExistRequest, params PackagesExistParams) (r PackagesExistRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SearchLanguagePackages implements SearchLanguagePackages operation.
//
// Search for language specific packages.
//
// GET /search/language
func (UnimplementedHandler) SearchLanguagePackages(ctx context.Context, params SearchLanguagePackagesParams) (r SearchLanguagePackagesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SearchSystemPackages implements SearchSystemPackages operation.
//
// Search for system packages.
//
// GET /search/system
func (UnimplementedHandler) SearchSystemPackages(ctx context.Context, params SearchSystemPackagesParams) (r SearchSystemPackagesRes, _ error) {
	return r, ht.ErrNotImplemented
}
