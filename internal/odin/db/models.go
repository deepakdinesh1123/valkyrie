// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ExecRequest struct {
	ID                  int32       `db:"id" json:"id"`
	Hash                string      `db:"hash" json:"hash"`
	Code                string      `db:"code" json:"code"`
	Path                string      `db:"path" json:"path"`
	Flake               string      `db:"flake" json:"flake"`
	Args                pgtype.Text `db:"args" json:"args"`
	ProgrammingLanguage pgtype.Text `db:"programming_language" json:"programming_language"`
}

type Job struct {
	ID            int64              `db:"id" json:"id"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	TimeOut       pgtype.Int4        `db:"time_out" json:"time_out"`
	StartedAt     pgtype.Timestamptz `db:"started_at" json:"started_at"`
	ExecRequestID pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	CurrentState  string             `db:"current_state" json:"current_state"`
	Retries       pgtype.Int4        `db:"retries" json:"retries"`
	MaxRetries    pgtype.Int4        `db:"max_retries" json:"max_retries"`
	WorkerID      pgtype.Int4        `db:"worker_id" json:"worker_id"`
}

type JobGroup struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type JobRun struct {
	ID            int64              `db:"id" json:"id"`
	JobID         pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID      pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt     pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt    pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	ExecRequestID pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs      string             `db:"exec_logs" json:"exec_logs"`
	NixLogs       pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success       pgtype.Bool        `db:"success" json:"success"`
}

type JobType struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type Worker struct {
	ID            int32              `db:"id" json:"id"`
	Name          string             `db:"name" json:"name"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"created_at"`
	LastHeartbeat pgtype.Timestamptz `db:"last_heartbeat" json:"last_heartbeat"`
}
