package db

import (
	"context"
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
		job, err := q.FetchJob(ctx, FetchJobParams{
			Workerid: arg.WorkerID,
			Jobtype:  "sandbox",
		})
		if err != nil {
			return err
		}
		sandbox, err := q.GetSandbox(ctx, job.Arguments.SandboxId)
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
