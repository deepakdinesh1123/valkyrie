// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Job struct {
	ID         int64              `db:"id" json:"id"`
	InsertedAt pgtype.Timestamptz `db:"inserted_at" json:"inserted_at"`
	WorkerID   pgtype.Int4        `db:"worker_id" json:"worker_id"`
	Script     string             `db:"script" json:"script"`
	ScriptPath string             `db:"script_path" json:"script_path"`
	Args       pgtype.Text        `db:"args" json:"args"`
	Flake      string             `db:"flake" json:"flake"`
	Language   string             `db:"language" json:"language"`
	Completed  bool               `db:"completed" json:"completed"`
	Running    bool               `db:"running" json:"running"`
}

type JobGroup struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type JobRun struct {
	ID         int64              `db:"id" json:"id"`
	JobID      int64              `db:"job_id" json:"job_id"`
	WorkerID   int32              `db:"worker_id" json:"worker_id"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"created_at"`
	StartedAt  pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	Script     string             `db:"script" json:"script"`
	Flake      string             `db:"flake" json:"flake"`
	Args       pgtype.Text        `db:"args" json:"args"`
	Logs       pgtype.Text        `db:"logs" json:"logs"`
}

type JobType struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type Worker struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}
