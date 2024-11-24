package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type FetchSandboxJobTxParams struct {
	WorkerID int32
}

type FetchSandboxJobTxResult struct {
	Sandbox Sandbox
}

func (s *SQLStore) FetchSandboxJobTx(ctx context.Context, arg FetchSandboxJobTxParams) (FetchSandboxJobTxResult, error) {
	var fetchSandboxTxResult FetchSandboxJobTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		sandbox, err := q.FetchSandboxJob(ctx, pgtype.Int4{Int32: arg.WorkerID, Valid: true})
		if err != nil {
			return err
		}
		fetchSandboxTxResult.Sandbox = sandbox
		return nil
	})
	if err != nil {
		return FetchSandboxJobTxResult{}, err
	}
	return fetchSandboxTxResult, nil
}
