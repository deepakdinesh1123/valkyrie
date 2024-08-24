package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adhocore/gronx"
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
	CronExpression      string
	MaxRetries          int
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
			log.Printf("GetExecRequestByHash error: %v", err)
			switch err {
			case pgx.ErrNoRows:
				log.Println("InsertExecRequest")
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
		jobParams.LastScheduledAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
		jobParams.MaxRetries = pgtype.Int4{Int32: int32(arg.MaxRetries), Valid: true}
		if arg.CronExpression == "" {
			jobParams.NextRunAt = pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}
		} else {
			gron := gronx.New()
			if !gron.IsValid(arg.CronExpression) {
				return fmt.Errorf("invalid cron expression: %s", arg.CronExpression)
			}
			jobParams.CronExpression = pgtype.Text{String: arg.CronExpression, Valid: true}
			nextRunAt, err := gronx.NextTick(arg.CronExpression, true)
			if err != nil {
				log.Printf("NextTick error: %v", err)
				return err
			}
			jobParams.NextRunAt = pgtype.Timestamptz{
				Time:  nextRunAt,
				Valid: true,
			}
		}
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
