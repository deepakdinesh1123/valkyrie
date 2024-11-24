// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ExecRequest struct {
	ID                   int32       `db:"id" json:"id"`
	Hash                 string      `db:"hash" json:"hash"`
	Code                 pgtype.Text `db:"code" json:"code"`
	Flake                string      `db:"flake" json:"flake"`
	LanguageDependencies []string    `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string    `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text `db:"compile_args" json:"compile_args"`
	Files                []byte      `db:"files" json:"files"`
	Input                pgtype.Text `db:"input" json:"input"`
	Command              pgtype.Text `db:"command" json:"command"`
	Setup                pgtype.Text `db:"setup" json:"setup"`
	LanguageVersion      int64       `db:"language_version" json:"language_version"`
}

type Execution struct {
	ExecID        int64              `db:"exec_id" json:"exec_id"`
	JobID         pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID      pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt     pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt    pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"created_at"`
	ExecRequestID pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs      string             `db:"exec_logs" json:"exec_logs"`
	NixLogs       pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success       pgtype.Bool        `db:"success" json:"success"`
}

type Job struct {
	JobID         int64              `db:"job_id" json:"job_id"`
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

type JobType struct {
	ID        int32              `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"created_at"`
}

type Language struct {
	ID             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Extension      string `db:"extension" json:"extension"`
	MonacoLanguage string `db:"monaco_language" json:"monaco_language"`
	Template       string `db:"template" json:"template"`
	DefaultCode    string `db:"default_code" json:"default_code"`
}

type LanguageVersion struct {
	ID             int64       `db:"id" json:"id"`
	LanguageID     int64       `db:"language_id" json:"language_id"`
	Version        string      `db:"version" json:"version"`
	NixPackageName string      `db:"nix_package_name" json:"nix_package_name"`
	Template       pgtype.Text `db:"template" json:"template"`
	SearchQuery    string      `db:"search_query" json:"search_query"`
	DefaultVersion bool        `db:"default_version" json:"default_version"`
}

type Package struct {
	PackageID int64       `db:"package_id" json:"package_id"`
	Name      string      `db:"name" json:"name"`
	Version   string      `db:"version" json:"version"`
	Pkgtype   string      `db:"pkgtype" json:"pkgtype"`
	Language  pgtype.Text `db:"language" json:"language"`
	StorePath pgtype.Text `db:"store_path" json:"store_path"`
	TsvSearch interface{} `db:"tsv_search" json:"tsv_search"`
}

type Sandbox struct {
	SandboxID    int64              `db:"sandbox_id" json:"sandbox_id"`
	WorkerID     pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt    pgtype.Timestamptz `db:"started_at" json:"started_at"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
	GitUrl       pgtype.Text        `db:"git_url" json:"git_url"`
	SandboxUrl   pgtype.Text        `db:"sandbox_url" json:"sandbox_url"`
	CurrentState string             `db:"current_state" json:"current_state"`
}

type Worker struct {
	ID            int32              `db:"id" json:"id"`
	Name          string             `db:"name" json:"name"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"created_at"`
	LastHeartbeat pgtype.Timestamptz `db:"last_heartbeat" json:"last_heartbeat"`
	CurrentState  string             `db:"current_state" json:"current_state"`
}
