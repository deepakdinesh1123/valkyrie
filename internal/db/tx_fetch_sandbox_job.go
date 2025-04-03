package db

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/db/jsonschema"
)

type FetchSandboxJobTxParams struct {
	WorkerID int32
}

type FetchSandboxJobTxResult struct {
	Sandbox Sandbox
	Job     Job
}

func (s *SQLStore) FetchSandboxJobTx(ctx context.Context, arg FetchSandboxJobTxParams) (result FetchSandboxJobTxResult, err error) {
	err = s.execTx(ctx, func(q *Queries) error {
		job, err := q.FetchJob(ctx, FetchJobParams{
			Workerid: arg.WorkerID,
			Jobtype:  "sandbox",
		})
		if err != nil {
			return err
		}

		sandbox, err := q.GetSandbox(ctx, job.Arguments.SandboxConfig.SandboxId)
		if err != nil {
			return fmt.Errorf("failed to get sandbox: %w", err)
		}

		err = q.UpdateSandboxState(ctx, UpdateSandboxStateParams{
			SandboxID:    sandbox.SandboxID,
			CurrentState: "creating",
			Details: jsonschema.SandboxDetails{
				Message: "Creating sandbox",
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update sandbox state: %w", err)
		}

		result.Sandbox = sandbox
		result.Job = job
		return nil
	})
	if err != nil {
		return FetchSandboxJobTxResult{}, err
	}
	return result, nil
}
