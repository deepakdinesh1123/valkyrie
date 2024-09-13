// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: worker.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteWorker = `-- name: DeleteWorker :exec
delete from workers where id = $1
`

func (q *Queries) DeleteWorker(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteWorker, id)
	return err
}

const getAllWorkers = `-- name: GetAllWorkers :many
select id, name, created_at, last_heartbeat from workers
limit $1 offset $2
`

type GetAllWorkersParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

func (q *Queries) GetAllWorkers(ctx context.Context, arg GetAllWorkersParams) ([]Worker, error) {
	rows, err := q.db.Query(ctx, getAllWorkers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Worker
	for rows.Next() {
		var i Worker
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.LastHeartbeat,
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

const getStaleWorkers = `-- name: GetStaleWorkers :many
select id from workers
where last_heartbeat < now() - interval '1 minute'
`

func (q *Queries) GetStaleWorkers(ctx context.Context) ([]int32, error) {
	rows, err := q.db.Query(ctx, getStaleWorkers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTotalWorkers = `-- name: GetTotalWorkers :one
select count(*) from workers
`

func (q *Queries) GetTotalWorkers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalWorkers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getWorker = `-- name: GetWorker :one
select id, name, created_at, last_heartbeat from workers where name = $1
`

func (q *Queries) GetWorker(ctx context.Context, name string) (Worker, error) {
	row := q.db.QueryRow(ctx, getWorker, name)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.LastHeartbeat,
	)
	return i, err
}

const insertWorker = `-- name: InsertWorker :one
insert into workers
    (name)
values
    ($1)
returning id, name, created_at, last_heartbeat
`

func (q *Queries) InsertWorker(ctx context.Context, name string) (Worker, error) {
	row := q.db.QueryRow(ctx, insertWorker, name)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.LastHeartbeat,
	)
	return i, err
}

const updateHeartbeat = `-- name: UpdateHeartbeat :exec
update workers
set
    last_heartbeat = now()
where id = $1
`

func (q *Queries) UpdateHeartbeat(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, updateHeartbeat, id)
	return err
}

const workerTaskCount = `-- name: WorkerTaskCount :one
select count(*) from jobs
where status = 'scheduled' and worker_id = $1
`

func (q *Queries) WorkerTaskCount(ctx context.Context, workerID pgtype.Int4) (int64, error) {
	row := q.db.QueryRow(ctx, workerTaskCount, workerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
