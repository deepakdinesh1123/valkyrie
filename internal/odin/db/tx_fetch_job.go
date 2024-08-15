package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type FetchJobTxParams struct {
	WorkerID int32
}

type FetchJobTxResult struct {
	Job Job
}

func (s *SQLStore) FetchJobTx(ctx context.Context, arg FetchJobTxParams) (FetchJobTxResult, error) {
	var fetchJobTxResult FetchJobTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		job, err := q.FetchJob(ctx, pgtype.Int4{Int32: arg.WorkerID, Valid: true})
		if err != nil {
			return err
		}
		fetchJobTxResult.Job = job
		return nil
	})
	if err != nil {
		return FetchJobTxResult{}, err
	}
	return fetchJobTxResult, nil
}
