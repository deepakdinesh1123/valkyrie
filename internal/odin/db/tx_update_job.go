package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateJobResultTxParams struct {
	StartTime time.Time
	Job       Job
	WorkerId  int32
	ExecLogs  string
	NixLogs   string
	Success   bool
	Retry     bool
}

type UpdateJobTxResult struct {
	execution Execution
}

func (s *SQLStore) UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error) {
	var updateJobTxResult UpdateJobTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		if arg.Success {
			err := q.UpdateJobCompleted(ctx, arg.Job.JobID)
			if err != nil {
				return err
			}
		} else {
			if !arg.Retry {
				err := q.CancelJob(ctx, arg.Job.JobID)
				if err != nil {
					return err
				}
			} else {
				err := q.RetryJob(ctx, arg.Job.JobID)
				if err != nil {
					return err
				}
			}
		}
		execution, err := q.InsertExecution(ctx, InsertExecutionParams{
			JobID:         pgtype.Int8{Int64: int64(arg.Job.JobID), Valid: true},
			WorkerID:      pgtype.Int4{Int32: arg.WorkerId, Valid: true},
			StartedAt:     pgtype.Timestamptz{Time: arg.StartTime, Valid: true},
			FinishedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ExecRequestID: arg.Job.ExecRequestID,
			ExecLogs:      arg.ExecLogs,
			NixLogs:       pgtype.Text{String: arg.NixLogs, Valid: true},
			Success:       pgtype.Bool{Bool: arg.Success, Valid: true},
		})
		if err != nil {
			return err
		}
		updateJobTxResult.execution = execution
		return nil
	})
	if err != nil {
		return UpdateJobTxResult{}, err
	}
	return updateJobTxResult, nil
}
