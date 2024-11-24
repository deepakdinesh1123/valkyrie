package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type AddSandboxTxParams struct {
	GitUrl string
}

type AddSandboxTxResult struct {
	SandboxId int
}

func (s *SQLStore) AddSandboxTx(ctx context.Context, arg AddSandboxTxParams) (AddSandboxTxResult, error) {
	var addSandboxTxResult AddSandboxTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		sandbox, err := q.InsertSandbox(ctx,
			pgtype.Text{String: arg.GitUrl, Valid: true},
		)
		if err != nil {
			return err
		}
		addSandboxTxResult.SandboxId = int(sandbox.SandboxID)
		return nil
	})
	if err != nil {
		return AddSandboxTxResult{}, err
	}
	return addSandboxTxResult, nil
}
