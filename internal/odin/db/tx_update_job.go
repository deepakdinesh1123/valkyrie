package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateJobResultTxParams struct {
	StartTime time.Time
	Job       Job
	Message   string
	Success   bool
}

type UpdateJobTxResult struct {
	JobRun JobRun
}

func (s *SQLStore) UpdateJobResultTx(ctx context.Context, arg UpdateJobResultTxParams) (UpdateJobTxResult, error) {
	var updateJobTxResult UpdateJobTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		if arg.Success {
			err := q.UpdateJob(ctx, arg.Job.ID)
			if err != nil {
				return err
			}
		} else {
			err := q.StopJob(ctx, arg.Job.ID)
			if err != nil {
				return err
			}
		}
		jobRun, err := q.InsertJobRun(ctx, InsertJobRunParams{
			JobID:      arg.Job.ID,
			WorkerID:   arg.Job.WorkerID.Int32,
			StartedAt:  pgtype.Timestamptz{Time: arg.StartTime, Valid: true},
			FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			CreatedAt:  arg.Job.InsertedAt,
			Script:     arg.Job.Script,
			Flake:      arg.Job.Flake,
			Args:       arg.Job.Args,
			Logs:       arg.Message,
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
