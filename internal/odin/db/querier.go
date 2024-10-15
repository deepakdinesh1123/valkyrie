// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CancelJob(ctx context.Context, jobID int64) error
	CreateWorker(ctx context.Context, name string) (Worker, error)
	DeleteExecRequest(ctx context.Context, id int32) error
	DeleteJob(ctx context.Context, jobID int64) (int64, error)
	DeleteWorker(ctx context.Context, id int32) error
	FetchJob(ctx context.Context, workerID pgtype.Int4) (Job, error)
	GetAllExecutions(ctx context.Context, arg GetAllExecutionsParams) ([]GetAllExecutionsRow, error)
	GetAllJobs(ctx context.Context, arg GetAllJobsParams) ([]GetAllJobsRow, error)
	GetAllWorkers(ctx context.Context, arg GetAllWorkersParams) ([]Worker, error)
	GetExecRequest(ctx context.Context, id int32) (ExecRequest, error)
	GetExecRequestByHash(ctx context.Context, hash string) (ExecRequest, error)
	GetExecution(ctx context.Context, execID int64) (GetExecutionRow, error)
	GetExecutionsForJob(ctx context.Context, arg GetExecutionsForJobParams) ([]GetExecutionsForJobRow, error)
	GetJob(ctx context.Context, jobID int64) (GetJobRow, error)
	GetJobState(ctx context.Context, jobID int64) (string, error)
	GetStaleWorkers(ctx context.Context) ([]int32, error)
	GetTotalExecutions(ctx context.Context) (int64, error)
	GetTotalExecutionsForJob(ctx context.Context, jobID pgtype.Int8) (int64, error)
	GetTotalJobs(ctx context.Context) (int64, error)
	GetTotalWorkers(ctx context.Context) (int64, error)
	GetWorker(ctx context.Context, name string) (Worker, error)
	InsertExecRequest(ctx context.Context, arg InsertExecRequestParams) (int32, error)
	InsertExecution(ctx context.Context, arg InsertExecutionParams) (Execution, error)
	InsertJob(ctx context.Context, arg InsertJobParams) (Job, error)
	InsertWorker(ctx context.Context, arg InsertWorkerParams) (Worker, error)
	ListExecRequests(ctx context.Context, arg ListExecRequestsParams) ([]ExecRequest, error)
	PruneCompletedJobs(ctx context.Context) error
	RequeueLTJobs(ctx context.Context) error
	RequeueWorkerJobs(ctx context.Context, workerID pgtype.Int4) error
	RetryJob(ctx context.Context, jobID int64) error
	SearchLanguagePackages(ctx context.Context, arg SearchLanguagePackagesParams) ([]SearchLanguagePackagesRow, error)
	SearchSystemPackages(ctx context.Context, plaintoTsquery string) ([]SearchSystemPackagesRow, error)
	StopJob(ctx context.Context, jobID int64) error
	UpdateHeartbeat(ctx context.Context, id int32) error
	UpdateJobCompleted(ctx context.Context, jobID int64) error
	WorkerTaskCount(ctx context.Context, workerID pgtype.Int4) (int64, error)
	updateJobFailed(ctx context.Context, jobID int64) error
}

var _ Querier = (*Queries)(nil)
