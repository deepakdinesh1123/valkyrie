// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Executionresult struct {
	ID             int32
	Result         pgtype.Text
	Cmdlineargs    pgtype.Text
	Environment    pgtype.Text
	Flake          pgtype.Text
	DependencyFile pgtype.Text
	ExecutedAt     pgtype.Timestamp
	Sandbox        pgtype.UUID
}
