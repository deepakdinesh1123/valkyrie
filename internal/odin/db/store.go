package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error)
	FetchJobTx(ctx context.Context, arg FetchJobTxParams) (FetchJobTxResult, error)
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
