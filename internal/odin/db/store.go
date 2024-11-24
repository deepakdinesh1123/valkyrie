package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error)
	FetchJobTx(ctx context.Context, arg FetchJobTxParams) (FetchJobTxResult, error)
	AddJobTx(ctx context.Context, arg AddJobTxParams) (AddJobTxResult, error)
	AddSandboxTx(ctx context.Context, arg AddSandboxTxParams) (AddSandboxTxResult, error)
	FetchSandboxJobTx(ctx context.Context, arg FetchSandboxJobTxParams) (FetchSandboxJobTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
