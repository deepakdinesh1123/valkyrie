// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Jobqueue struct {
	ID           int64
	CreatedBy    pgtype.Text
	CreatedAt    pgtype.Timestamp
	StartedAt    pgtype.Timestamp
	CompletedAt  pgtype.Timestamp
	Script       pgtype.Text
	Args         []byte
	Logs         pgtype.Text
	Flake        pgtype.Text
	Language     pgtype.Text
	MemPeak      pgtype.Int4
	Timeout      pgtype.Int4
	Priority     pgtype.Int4
	LeaseTimeout pgtype.Float8
	Queue        string
	JobType      string
}