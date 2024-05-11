// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: execution_result.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getAllExecutionResults = `-- name: GetAllExecutionResults :many
SELECT id, result, cmdlineargs, environment, flake, dependency_file, executed_at, sandbox FROM ExecutionResult
ORDER BY executed_at
`

func (q *Queries) GetAllExecutionResults(ctx context.Context) ([]Executionresult, error) {
	rows, err := q.db.Query(ctx, getAllExecutionResults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Executionresult
	for rows.Next() {
		var i Executionresult
		if err := rows.Scan(
			&i.ID,
			&i.Result,
			&i.Cmdlineargs,
			&i.Environment,
			&i.Flake,
			&i.DependencyFile,
			&i.ExecutedAt,
			&i.Sandbox,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getResultUsingSandboxID = `-- name: GetResultUsingSandboxID :one
SELECT result
FROM ExecutionResult
WHERE sandbox = $1 LIMIT 1
`

func (q *Queries) GetResultUsingSandboxID(ctx context.Context, sandbox pgtype.UUID) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, getResultUsingSandboxID, sandbox)
	var result pgtype.Text
	err := row.Scan(&result)
	return result, err
}

const insertExecutionResult = `-- name: InsertExecutionResult :one
INSERT INTO ExecutionResult (
    result,
    cmdLineArgs,
    environment,
    flake,
    dependency_file, 
    sandbox
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, result, cmdlineargs, environment, flake, dependency_file, executed_at, sandbox
`

type InsertExecutionResultParams struct {
	Result         pgtype.Text
	Cmdlineargs    pgtype.Text
	Environment    pgtype.Text
	Flake          pgtype.Text
	DependencyFile pgtype.Text
	Sandbox        pgtype.UUID
}

func (q *Queries) InsertExecutionResult(ctx context.Context, arg InsertExecutionResultParams) (Executionresult, error) {
	row := q.db.QueryRow(ctx, insertExecutionResult,
		arg.Result,
		arg.Cmdlineargs,
		arg.Environment,
		arg.Flake,
		arg.DependencyFile,
		arg.Sandbox,
	)
	var i Executionresult
	err := row.Scan(
		&i.ID,
		&i.Result,
		&i.Cmdlineargs,
		&i.Environment,
		&i.Flake,
		&i.DependencyFile,
		&i.ExecutedAt,
		&i.Sandbox,
	)
	return i, err
}
