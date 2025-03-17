package db

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Store interface {
	Querier
	UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error)
	AddExecJobTx(ctx context.Context, arg AddJobTxParams) (AddJobTxResult, error)
	AddSandboxJobTx(ctx context.Context, arg AddSandboxTxParams) (AddSandboxTxResult, error)
	FetchSandboxJobTx(ctx context.Context, arg FetchSandboxJobTxParams) (FetchSandboxJobTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
	logger *zerolog.Logger
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
		logger:   logs.GetLogger(logs.NewLogConfig(logs.WithSource("DB"))),
	}
}
