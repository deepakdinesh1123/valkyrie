package db

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db/jsonschema"
	"github.com/jackc/pgx/v5/pgtype"
)

type AddSandboxTxParams struct {
	SandboxConfig jsonschema.SandboxConfig
	MaxRetries    int
}

type AddSandboxTxResult struct {
	SandboxId int64
}

func (s *SQLStore) AddSandboxJobTx(ctx context.Context, arg AddSandboxTxParams) (AddSandboxTxResult, error) {
	var addSandboxTxResult AddSandboxTxResult
	err := s.execTx(ctx, func(q *Queries) error {

		sandbox, err := q.InsertSandbox(ctx, InsertSandboxParams{
			Config: arg.SandboxConfig,
			Details: jsonschema.SandboxDetails{
				Message: "Adding sandbox job",
			},
		})
		if err != nil {
			return err
		}
		addSandboxTxResult.SandboxId = sandbox.SandboxID
		arg.SandboxConfig.SandboxId = sandbox.SandboxID
		_, err = q.InsertJob(ctx, InsertJobParams{
			JobType: "sandbox",
			Arguments: jsonschema.JobArguments{
				SandboxConfig: arg.SandboxConfig,
			},
			MaxRetries: pgtype.Int4{Int32: int32(arg.MaxRetries), Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return AddSandboxTxResult{}, err
	}
	return addSandboxTxResult, nil
}
