package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AddJobTxParams struct {
	Code                string
	Flake               string
	Hash                string
	Args                string
	Path                string
	ProgrammingLanguage string
	MaxRetries          int
	Timeout             int32
}

type AddJobTxResult struct {
	JobID int64
}

func (s *SQLStore) AddJobTx(ctx context.Context, arg AddJobTxParams) (AddJobTxResult, error) {
	var addJobTxResult AddJobTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		execReq, err := s.GetExecRequestByHash(ctx, arg.Hash)
		var execId int32
		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				execId, err = s.InsertExecRequest(ctx, InsertExecRequestParams{
					Code:                arg.Code,
					Flake:               arg.Flake,
					Hash:                arg.Hash,
					Args:                pgtype.Text{String: arg.Args, Valid: true},
					ProgrammingLanguage: pgtype.Text{String: arg.ProgrammingLanguage, Valid: true},
					Path:                arg.Path,
				})
				if err != nil {
					return err
				}
			default:
				log.Println("GetExecRequestByHash error default: ", err)
				return err
			}
		} else {
			execId = execReq.ID
		}

		var jobParams InsertJobParams
		jobParams.ExecRequestID = pgtype.Int4{Int32: execId, Valid: true}
		jobParams.MaxRetries = pgtype.Int4{Int32: int32(arg.MaxRetries), Valid: true}
		jobParams.TimeOut = pgtype.Int4{Int32: arg.Timeout, Valid: true}
		job, err := s.InsertJob(ctx, jobParams)
		if err != nil {
			log.Printf("InsertJob error: %v", err)
			return err
		}
		addJobTxResult.JobID = job.ID
		return nil
	})
	if err != nil {
		log.Printf("execTx error: %v", err)
		return AddJobTxResult{}, err
	}
	return addJobTxResult, nil
}
