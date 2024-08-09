// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	DeleteJob(ctx context.Context, id int64) error
	FetchJob(ctx context.Context, workerID pgtype.Int4) (Jobqueue, error)
	GetAllJobs(ctx context.Context, arg GetAllJobsParams) ([]GetAllJobsRow, error)
	GetAllWorkers(ctx context.Context, arg GetAllWorkersParams) ([]Worker, error)
	GetJob(ctx context.Context, id int64) (Jobqueue, error)
	GetResultUsingExecutionID(ctx context.Context, id int64) (Jobqueue, error)
	GetTotalJobs(ctx context.Context) (int64, error)
	GetTotalWorkers(ctx context.Context) (int64, error)
	GetWorker(ctx context.Context, name pgtype.Text) (Worker, error)
	InsertJob(ctx context.Context, arg InsertJobParams) (Jobqueue, error)
	InsertWorker(ctx context.Context, name pgtype.Text) (Worker, error)
	UpdateJob(ctx context.Context, arg UpdateJobParams) (Jobqueue, error)
}

var _ Querier = (*Queries)(nil)
