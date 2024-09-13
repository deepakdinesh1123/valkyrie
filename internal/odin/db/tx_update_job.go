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
	Message   string
	Success   bool
	Retry     bool
}

type UpdateJobTxResult struct {
	JobRun JobRun
}

func (s *SQLStore) UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error) {
	var updateJobTxResult UpdateJobTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		if arg.Success {
			err := q.UpdateJobCompleted(ctx, arg.Job.ID)
			if err != nil {
				return err
			}
		} else {
			if !arg.Retry {
				err := q.CancelJob(ctx, arg.Job.ID)
				if err != nil {
					return err
				}
			} else {
				err := q.RetryJob(ctx, arg.Job.ID)
				if err != nil {
					return err
				}
			}
		}
		jobRun, err := q.InsertJobRun(ctx, InsertJobRunParams{
			JobID:         arg.Job.ID,
			WorkerID:      arg.WorkerId,
			StartedAt:     pgtype.Timestamptz{Time: arg.StartTime, Valid: true},
			FinishedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ExecRequestID: arg.Job.ExecRequestID,
			ExecLogs:      arg.Message,
		})
		if err != nil {
			return err
		}
		updateJobTxResult.JobRun = jobRun
		return nil
	})
	if err != nil {
		return UpdateJobTxResult{}, err
	}
	return updateJobTxResult, nil
}
