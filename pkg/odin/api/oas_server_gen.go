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
	// PUT /executions/jobs/{JobId}
	CancelExecutionJob(ctx context.Context, params CancelExecutionJobParams) (CancelExecutionJobRes, error)
	// CreateLanguage implements createLanguage operation.
	//
	// Create a new language entry in the database.
	//
	// POST /languages/create
	CreateLanguage(ctx context.Context, req *Language, params CreateLanguageParams) (CreateLanguageRes, error)
	// CreateLanguageVersion implements createLanguageVersion operation.
	//
	// Create a new language version entry in the database.
	//
	// POST /language-versions/create
	CreateLanguageVersion(ctx context.Context, req *LanguageVersion, params CreateLanguageVersionParams) (CreateLanguageVersionRes, error)
	// DeleteExecutionJob implements deleteExecutionJob operation.
	//
	// Delete execution job.
	//
	// DELETE /executions/jobs/{JobId}
	DeleteExecutionJob(ctx context.Context, params DeleteExecutionJobParams) (DeleteExecutionJobRes, error)
	// DeleteLanguage implements deleteLanguage operation.
	//
	// Delete a specific language by its ID.
	//
	// DELETE /languages/{id}
	DeleteLanguage(ctx context.Context, params DeleteLanguageParams) (DeleteLanguageRes, error)
	// DeleteLanguageVersion implements deleteLanguageVersion operation.
	//
	// Delete a specific language version by its ID.
	//
	// DELETE /language-versions/{id}
	DeleteLanguageVersion(ctx context.Context, params DeleteLanguageVersionParams) (DeleteLanguageVersionRes, error)
	// Execute implements execute operation.
	//
	// Execute a script.
	//
	// POST /executions/execute
	Execute(ctx context.Context, req *ExecutionRequest, params ExecuteParams) (ExecuteRes, error)
	// FetchLanguagePackages implements FetchLanguagePackages operation.
	//
	// Initialize the search results content with a default set of language specific packages.
	//
	// GET /fetch/language
	FetchLanguagePackages(ctx context.Context, params FetchLanguagePackagesParams) (FetchLanguagePackagesRes, error)
	// FetchSystemPackages implements FetchSystemPackages operation.
	//
	// Initialize the search results content with a default set of system packages.
	//
	// GET /fetch/system
	FetchSystemPackages(ctx context.Context, params FetchSystemPackagesParams) (FetchSystemPackagesRes, error)
	// FlakeJobIdGet implements GET /flake/{jobId} operation.
	//
	// Fetches flake of a given job.
	//
	// GET /flake/{jobId}
	FlakeJobIdGet(ctx context.Context, params FlakeJobIdGetParams) (FlakeJobIdGetRes, error)
	// GetAllExecutionJobs implements getAllExecutionJobs operation.
	//
	// Get all execution jobs.
	//
	// GET /jobs/execution
	GetAllExecutionJobs(ctx context.Context, params GetAllExecutionJobsParams) (GetAllExecutionJobsRes, error)
	// GetAllExecutions implements getAllExecutions operation.
	//
	// Get all executions.
	//
	// GET /executions
	GetAllExecutions(ctx context.Context, params GetAllExecutionsParams) (GetAllExecutionsRes, error)
	// GetAllLanguageVersions implements getAllLanguageVersions operation.
	//
	// Retrieve a list of all language versions from the database.
	//
	// GET /language-versions
	GetAllLanguageVersions(ctx context.Context, params GetAllLanguageVersionsParams) (GetAllLanguageVersionsRes, error)
	// GetAllLanguages implements getAllLanguages operation.
	//
	// Retrieve a list of all languages from the database.
	//
	// GET /languages
	GetAllLanguages(ctx context.Context, params GetAllLanguagesParams) (GetAllLanguagesRes, error)
	// GetAllVersions implements getAllVersions operation.
	//
	// Retrieve a list of all language versions from the database.
	//
	// GET /languages/{id}/versions
	GetAllVersions(ctx context.Context, params GetAllVersionsParams) (GetAllVersionsRes, error)
	// GetExecutionConfig implements getExecutionConfig operation.
	//
	// Get execution config.
	//
	// GET /execution/config
	GetExecutionConfig(ctx context.Context, params GetExecutionConfigParams) (GetExecutionConfigRes, error)
	// GetExecutionJobById implements getExecutionJobById operation.
	//
	// Get execution job.
	//
	// GET /executions/jobs/{JobId}
	GetExecutionJobById(ctx context.Context, params GetExecutionJobByIdParams) (GetExecutionJobByIdRes, error)
	// GetExecutionResultById implements getExecutionResultById operation.
	//
	// Get execution result by id.
	//
	// GET /executions/{execId}
	GetExecutionResultById(ctx context.Context, params GetExecutionResultByIdParams) (GetExecutionResultByIdRes, error)
	// GetExecutionsForJob implements getExecutionsForJob operation.
	//
	// Get executions of given job.
	//
	// GET /jobs/{JobId}/executions
	GetExecutionsForJob(ctx context.Context, params GetExecutionsForJobParams) (GetExecutionsForJobRes, error)
	// GetLanguageById implements getLanguageById operation.
	//
	// Retrieve a language entry from the database using its ID.
	//
	// GET /languages/{id}
	GetLanguageById(ctx context.Context, params GetLanguageByIdParams) (GetLanguageByIdRes, error)
	// GetLanguageVersionById implements getLanguageVersionById operation.
	//
	// Retrieve a language version entry from the database using its ID.
	//
	// GET /language-versions/{id}
	GetLanguageVersionById(ctx context.Context, params GetLanguageVersionByIdParams) (GetLanguageVersionByIdRes, error)
	// GetVersion implements getVersion operation.
	//
	// Get version.
	//
	// GET /version
	GetVersion(ctx context.Context, params GetVersionParams) (GetVersionRes, error)
	// Health implements health operation.
	//
	// Health Check.
	//
	// GET /health
	Health(ctx context.Context) error
	// PackagesExist implements PackagesExist operation.
	//
	// Verify the package list is available for the language version while switching between language
	// versions.
	//
	// POST /packages/exist
	PackagesExist(ctx context.Context, req *PackageExistRequest, params PackagesExistParams) (PackagesExistRes, error)
	// SearchLanguagePackages implements SearchLanguagePackages operation.
	//
	// Search for language specific packages.
	//
	// GET /search/language
	SearchLanguagePackages(ctx context.Context, params SearchLanguagePackagesParams) (SearchLanguagePackagesRes, error)
	// SearchSystemPackages implements SearchSystemPackages operation.
	//
	// Search for system packages.
	//
	// GET /search/system
	SearchSystemPackages(ctx context.Context, params SearchSystemPackagesParams) (SearchSystemPackagesRes, error)
	// UpdateLanguage implements updateLanguage operation.
	//
	// Update the details of a specific language by its ID.
	//
	// PUT /languages/{id}
	UpdateLanguage(ctx context.Context, req *Language, params UpdateLanguageParams) (UpdateLanguageRes, error)
	// UpdateLanguageVersion implements updateLanguageVersion operation.
	//
	// Update the details of a specific language version by its ID.
	//
	// PUT /language-versions/{id}
	UpdateLanguageVersion(ctx context.Context, req *LanguageVersion, params UpdateLanguageVersionParams) (UpdateLanguageVersionRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
