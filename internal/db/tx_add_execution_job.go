package db

import (
	"context"
	"log"

	"github.com/deepakdinesh1123/valkyrie/internal/db/jsonschema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AddJobTxParams struct {
	Hash                 string
	Code                 string
	Flake                string
	LanguageDependencies []string
	SystemDependencies   []string
	CmdLineArgs          string
	CompilerArgs         string
	Files                []byte
	Input                string
	Command              string
	Setup                string
	MaxRetries           int
	Timeout              int32
	LangVersion          int64
	SystemSetup          string
	PkgIndex             string
	Extension            string
}

type AddJobTxResult struct {
	JobID int64
}

func (s *SQLStore) AddExecJobTx(ctx context.Context, arg AddJobTxParams) (AddJobTxResult, error) {
	var addJobTxResult AddJobTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		execReq, err := s.GetExecRequestByHash(ctx, arg.Hash)
		var execId int32
		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				execId, err = s.InsertExecRequest(ctx, InsertExecRequestParams{
					Hash:                 arg.Hash,
					Code:                 pgtype.Text{String: arg.Code, Valid: true},
					LanguageDependencies: arg.LanguageDependencies,
					SystemDependencies:   arg.SystemDependencies,
					Flake:                arg.Flake,
					CmdLineArgs:          pgtype.Text{String: arg.CmdLineArgs, Valid: true},
					CompileArgs:          pgtype.Text{String: arg.CompilerArgs, Valid: true},
					Command:              pgtype.Text{String: arg.Command, Valid: true},
					LanguageVersion:      arg.LangVersion,
					Setup:                pgtype.Text{String: arg.Setup, Valid: true},
					Files:                arg.Files,
					Input:                pgtype.Text{String: arg.Input, Valid: true},
					SystemSetup:          pgtype.Text{String: arg.SystemSetup, Valid: true},
					PkgIndex:             pgtype.Text{String: arg.PkgIndex, Valid: true},
					Extension:            pgtype.Text{String: arg.Extension, Valid: true},
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
		jobParams.Arguments = jsonschema.JobArguments{
			ExecConfig: jsonschema.ExecConfig{
				ExecReqId: execId,
			},
		}
		jobParams.MaxRetries = pgtype.Int4{Int32: int32(arg.MaxRetries), Valid: true}
		jobParams.TimeOut = pgtype.Int4{Int32: arg.Timeout, Valid: true}
		jobParams.JobType = "execution"
		job, err := s.InsertJob(ctx, jobParams)
		if err != nil {
			log.Printf("InsertJob error: %v", err)
			return err
		}
		addJobTxResult.JobID = job.JobID
		return nil
	})
	if err != nil {
		log.Printf("execTx error: %v", err)
		return AddJobTxResult{}, err
	}
	return addJobTxResult, nil
}
