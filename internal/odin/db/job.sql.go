// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: job.sql

package db

import (
	"context"

	jsonschema "github.com/deepakdinesh1123/valkyrie/internal/odin/db/jsonschema"
	"github.com/jackc/pgx/v5/pgtype"
)

const cancelJob = `-- name: CancelJob :exec
update jobs set current_state = 'cancelled', updated_at = now(), worker_id = null where job_id = $1
`

func (q *Queries) CancelJob(ctx context.Context, jobID int64) error {
	_, err := q.db.Exec(ctx, cancelJob, jobID)
	return err
}

const deleteJob = `-- name: DeleteJob :one
delete from jobs where job_id = $1 and current_state in ('pending', 'cancelled', 'failed') returning job_id
`

func (q *Queries) DeleteJob(ctx context.Context, jobID int64) (int64, error) {
	row := q.db.QueryRow(ctx, deleteJob, jobID)
	var job_id int64
	err := row.Scan(&job_id)
	return job_id, err
}

const fetchJob = `-- name: FetchJob :one
with cte as (
    select job_id
    from jobs
    where 
        current_state = 'pending'
        and job_type = $2::text
        and retries < max_retries
    order by job_id asc
    for update skip locked
    limit 1
)
update jobs
set current_state = 'scheduled', 
    started_at = now(), 
    worker_id = $1::int, 
    updated_at = now()
where job_id = (select job_id from cte)
returning job_id, created_at, updated_at, time_out, started_at, arguments, current_state, retries, max_retries, worker_id, job_type
`

type FetchJobParams struct {
	Workerid int32  `db:"workerid" json:"workerid"`
	Jobtype  string `db:"jobtype" json:"jobtype"`
}

func (q *Queries) FetchJob(ctx context.Context, arg FetchJobParams) (Job, error) {
	row := q.db.QueryRow(ctx, fetchJob, arg.Workerid, arg.Jobtype)
	var i Job
	err := row.Scan(
		&i.JobID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TimeOut,
		&i.StartedAt,
		&i.Arguments,
		&i.CurrentState,
		&i.Retries,
		&i.MaxRetries,
		&i.WorkerID,
		&i.JobType,
	)
	return i, err
}

const getAllExecutionJobs = `-- name: GetAllExecutionJobs :many
SELECT job_id, created_at, updated_at, time_out, started_at, arguments, current_state, retries, max_retries, worker_id, job_type, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version 
FROM jobs
INNER JOIN exec_request 
ON CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) = exec_request.id
WHERE job_id >= $1
ORDER BY jobs.job_id
LIMIT $2
`

type GetAllExecutionJobsParams struct {
	JobID int64 `db:"job_id" json:"job_id"`
	Limit int64 `db:"limit" json:"limit"`
}

type GetAllExecutionJobsRow struct {
	JobID                int64                   `db:"job_id" json:"job_id"`
	CreatedAt            pgtype.Timestamptz      `db:"created_at" json:"created_at"`
	UpdatedAt            pgtype.Timestamptz      `db:"updated_at" json:"updated_at"`
	TimeOut              pgtype.Int4             `db:"time_out" json:"time_out"`
	StartedAt            pgtype.Timestamptz      `db:"started_at" json:"started_at"`
	Arguments            jsonschema.JobArguments `db:"arguments" json:"arguments"`
	CurrentState         string                  `db:"current_state" json:"current_state"`
	Retries              pgtype.Int4             `db:"retries" json:"retries"`
	MaxRetries           pgtype.Int4             `db:"max_retries" json:"max_retries"`
	WorkerID             pgtype.Int4             `db:"worker_id" json:"worker_id"`
	JobType              string                  `db:"job_type" json:"job_type"`
	ID                   int32                   `db:"id" json:"id"`
	Hash                 string                  `db:"hash" json:"hash"`
	Code                 pgtype.Text             `db:"code" json:"code"`
	Flake                string                  `db:"flake" json:"flake"`
	LanguageDependencies []string                `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string                `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text             `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text             `db:"compile_args" json:"compile_args"`
	Files                []byte                  `db:"files" json:"files"`
	Input                pgtype.Text             `db:"input" json:"input"`
	Command              pgtype.Text             `db:"command" json:"command"`
	Setup                pgtype.Text             `db:"setup" json:"setup"`
	LanguageVersion      int64                   `db:"language_version" json:"language_version"`
}

func (q *Queries) GetAllExecutionJobs(ctx context.Context, arg GetAllExecutionJobsParams) ([]GetAllExecutionJobsRow, error) {
	rows, err := q.db.Query(ctx, getAllExecutionJobs, arg.JobID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllExecutionJobsRow
	for rows.Next() {
		var i GetAllExecutionJobsRow
		if err := rows.Scan(
			&i.JobID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TimeOut,
			&i.StartedAt,
			&i.Arguments,
			&i.CurrentState,
			&i.Retries,
			&i.MaxRetries,
			&i.WorkerID,
			&i.JobType,
			&i.ID,
			&i.Hash,
			&i.Code,
			&i.Flake,
			&i.LanguageDependencies,
			&i.SystemDependencies,
			&i.CmdLineArgs,
			&i.CompileArgs,
			&i.Files,
			&i.Input,
			&i.Command,
			&i.Setup,
			&i.LanguageVersion,
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

const getAllExecutions = `-- name: GetAllExecutions :many
select exec_id, job_id, worker_id, started_at, finished_at, created_at, exec_request_id, exec_logs, nix_logs, success, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where exec_id >= $1
order by started_at desc
limit $2
`

type GetAllExecutionsParams struct {
	ExecID int64 `db:"exec_id" json:"exec_id"`
	Limit  int64 `db:"limit" json:"limit"`
}

type GetAllExecutionsRow struct {
	ExecID               int64              `db:"exec_id" json:"exec_id"`
	JobID                pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID             pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt            pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt           pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	CreatedAt            pgtype.Timestamptz `db:"created_at" json:"created_at"`
	ExecRequestID        pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs             string             `db:"exec_logs" json:"exec_logs"`
	NixLogs              pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success              pgtype.Bool        `db:"success" json:"success"`
	ID                   int32              `db:"id" json:"id"`
	Hash                 string             `db:"hash" json:"hash"`
	Code                 pgtype.Text        `db:"code" json:"code"`
	Flake                string             `db:"flake" json:"flake"`
	LanguageDependencies []string           `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string           `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text        `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text        `db:"compile_args" json:"compile_args"`
	Files                []byte             `db:"files" json:"files"`
	Input                pgtype.Text        `db:"input" json:"input"`
	Command              pgtype.Text        `db:"command" json:"command"`
	Setup                pgtype.Text        `db:"setup" json:"setup"`
	LanguageVersion      int64              `db:"language_version" json:"language_version"`
}

func (q *Queries) GetAllExecutions(ctx context.Context, arg GetAllExecutionsParams) ([]GetAllExecutionsRow, error) {
	rows, err := q.db.Query(ctx, getAllExecutions, arg.ExecID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllExecutionsRow
	for rows.Next() {
		var i GetAllExecutionsRow
		if err := rows.Scan(
			&i.ExecID,
			&i.JobID,
			&i.WorkerID,
			&i.StartedAt,
			&i.FinishedAt,
			&i.CreatedAt,
			&i.ExecRequestID,
			&i.ExecLogs,
			&i.NixLogs,
			&i.Success,
			&i.ID,
			&i.Hash,
			&i.Code,
			&i.Flake,
			&i.LanguageDependencies,
			&i.SystemDependencies,
			&i.CmdLineArgs,
			&i.CompileArgs,
			&i.Files,
			&i.Input,
			&i.Command,
			&i.Setup,
			&i.LanguageVersion,
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

const getExecution = `-- name: GetExecution :one
select exec_id, job_id, worker_id, started_at, finished_at, created_at, exec_request_id, exec_logs, nix_logs, success, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.exec_id = $1
`

type GetExecutionRow struct {
	ExecID               int64              `db:"exec_id" json:"exec_id"`
	JobID                pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID             pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt            pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt           pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	CreatedAt            pgtype.Timestamptz `db:"created_at" json:"created_at"`
	ExecRequestID        pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs             string             `db:"exec_logs" json:"exec_logs"`
	NixLogs              pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success              pgtype.Bool        `db:"success" json:"success"`
	ID                   int32              `db:"id" json:"id"`
	Hash                 string             `db:"hash" json:"hash"`
	Code                 pgtype.Text        `db:"code" json:"code"`
	Flake                string             `db:"flake" json:"flake"`
	LanguageDependencies []string           `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string           `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text        `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text        `db:"compile_args" json:"compile_args"`
	Files                []byte             `db:"files" json:"files"`
	Input                pgtype.Text        `db:"input" json:"input"`
	Command              pgtype.Text        `db:"command" json:"command"`
	Setup                pgtype.Text        `db:"setup" json:"setup"`
	LanguageVersion      int64              `db:"language_version" json:"language_version"`
}

func (q *Queries) GetExecution(ctx context.Context, execID int64) (GetExecutionRow, error) {
	row := q.db.QueryRow(ctx, getExecution, execID)
	var i GetExecutionRow
	err := row.Scan(
		&i.ExecID,
		&i.JobID,
		&i.WorkerID,
		&i.StartedAt,
		&i.FinishedAt,
		&i.CreatedAt,
		&i.ExecRequestID,
		&i.ExecLogs,
		&i.NixLogs,
		&i.Success,
		&i.ID,
		&i.Hash,
		&i.Code,
		&i.Flake,
		&i.LanguageDependencies,
		&i.SystemDependencies,
		&i.CmdLineArgs,
		&i.CompileArgs,
		&i.Files,
		&i.Input,
		&i.Command,
		&i.Setup,
		&i.LanguageVersion,
	)
	return i, err
}

const getExecutionJob = `-- name: GetExecutionJob :one
select job_id, created_at, updated_at, time_out, started_at, arguments, current_state, retries, max_retries, worker_id, job_type, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version from jobs inner join exec_request on CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) = exec_request.id where jobs.job_id = $1
`

type GetExecutionJobRow struct {
	JobID                int64                   `db:"job_id" json:"job_id"`
	CreatedAt            pgtype.Timestamptz      `db:"created_at" json:"created_at"`
	UpdatedAt            pgtype.Timestamptz      `db:"updated_at" json:"updated_at"`
	TimeOut              pgtype.Int4             `db:"time_out" json:"time_out"`
	StartedAt            pgtype.Timestamptz      `db:"started_at" json:"started_at"`
	Arguments            jsonschema.JobArguments `db:"arguments" json:"arguments"`
	CurrentState         string                  `db:"current_state" json:"current_state"`
	Retries              pgtype.Int4             `db:"retries" json:"retries"`
	MaxRetries           pgtype.Int4             `db:"max_retries" json:"max_retries"`
	WorkerID             pgtype.Int4             `db:"worker_id" json:"worker_id"`
	JobType              string                  `db:"job_type" json:"job_type"`
	ID                   int32                   `db:"id" json:"id"`
	Hash                 string                  `db:"hash" json:"hash"`
	Code                 pgtype.Text             `db:"code" json:"code"`
	Flake                string                  `db:"flake" json:"flake"`
	LanguageDependencies []string                `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string                `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text             `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text             `db:"compile_args" json:"compile_args"`
	Files                []byte                  `db:"files" json:"files"`
	Input                pgtype.Text             `db:"input" json:"input"`
	Command              pgtype.Text             `db:"command" json:"command"`
	Setup                pgtype.Text             `db:"setup" json:"setup"`
	LanguageVersion      int64                   `db:"language_version" json:"language_version"`
}

func (q *Queries) GetExecutionJob(ctx context.Context, jobID int64) (GetExecutionJobRow, error) {
	row := q.db.QueryRow(ctx, getExecutionJob, jobID)
	var i GetExecutionJobRow
	err := row.Scan(
		&i.JobID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TimeOut,
		&i.StartedAt,
		&i.Arguments,
		&i.CurrentState,
		&i.Retries,
		&i.MaxRetries,
		&i.WorkerID,
		&i.JobType,
		&i.ID,
		&i.Hash,
		&i.Code,
		&i.Flake,
		&i.LanguageDependencies,
		&i.SystemDependencies,
		&i.CmdLineArgs,
		&i.CompileArgs,
		&i.Files,
		&i.Input,
		&i.Command,
		&i.Setup,
		&i.LanguageVersion,
	)
	return i, err
}

const getExecutionsForJob = `-- name: GetExecutionsForJob :many
select exec_id, job_id, worker_id, started_at, finished_at, created_at, exec_request_id, exec_logs, nix_logs, success, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.job_id = $1 and exec_id >= $2
order by finished_at desc
limit $3
`

type GetExecutionsForJobParams struct {
	JobID  pgtype.Int8 `db:"job_id" json:"job_id"`
	ExecID int64       `db:"exec_id" json:"exec_id"`
	Limit  int64       `db:"limit" json:"limit"`
}

type GetExecutionsForJobRow struct {
	ExecID               int64              `db:"exec_id" json:"exec_id"`
	JobID                pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID             pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt            pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt           pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	CreatedAt            pgtype.Timestamptz `db:"created_at" json:"created_at"`
	ExecRequestID        pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs             string             `db:"exec_logs" json:"exec_logs"`
	NixLogs              pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success              pgtype.Bool        `db:"success" json:"success"`
	ID                   int32              `db:"id" json:"id"`
	Hash                 string             `db:"hash" json:"hash"`
	Code                 pgtype.Text        `db:"code" json:"code"`
	Flake                string             `db:"flake" json:"flake"`
	LanguageDependencies []string           `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string           `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text        `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text        `db:"compile_args" json:"compile_args"`
	Files                []byte             `db:"files" json:"files"`
	Input                pgtype.Text        `db:"input" json:"input"`
	Command              pgtype.Text        `db:"command" json:"command"`
	Setup                pgtype.Text        `db:"setup" json:"setup"`
	LanguageVersion      int64              `db:"language_version" json:"language_version"`
}

func (q *Queries) GetExecutionsForJob(ctx context.Context, arg GetExecutionsForJobParams) ([]GetExecutionsForJobRow, error) {
	rows, err := q.db.Query(ctx, getExecutionsForJob, arg.JobID, arg.ExecID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExecutionsForJobRow
	for rows.Next() {
		var i GetExecutionsForJobRow
		if err := rows.Scan(
			&i.ExecID,
			&i.JobID,
			&i.WorkerID,
			&i.StartedAt,
			&i.FinishedAt,
			&i.CreatedAt,
			&i.ExecRequestID,
			&i.ExecLogs,
			&i.NixLogs,
			&i.Success,
			&i.ID,
			&i.Hash,
			&i.Code,
			&i.Flake,
			&i.LanguageDependencies,
			&i.SystemDependencies,
			&i.CmdLineArgs,
			&i.CompileArgs,
			&i.Files,
			&i.Input,
			&i.Command,
			&i.Setup,
			&i.LanguageVersion,
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

const getFlake = `-- name: GetFlake :one
SELECT flake 
FROM exec_request 
WHERE id = (
    SELECT CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) 
    FROM jobs 
    WHERE job_id = $1
)
`

func (q *Queries) GetFlake(ctx context.Context, jobID int64) (string, error) {
	row := q.db.QueryRow(ctx, getFlake, jobID)
	var flake string
	err := row.Scan(&flake)
	return flake, err
}

const getJobState = `-- name: GetJobState :one
select current_state from jobs where job_id = $1
`

func (q *Queries) GetJobState(ctx context.Context, jobID int64) (string, error) {
	row := q.db.QueryRow(ctx, getJobState, jobID)
	var current_state string
	err := row.Scan(&current_state)
	return current_state, err
}

const getLatestExecution = `-- name: GetLatestExecution :one
select exec_id, job_id, worker_id, started_at, finished_at, created_at, exec_request_id, exec_logs, nix_logs, success, id, hash, code, flake, language_dependencies, system_dependencies, cmd_line_args, compile_args, files, input, command, setup, language_version from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.job_id = $1
order by finished_at desc
limit 1
`

type GetLatestExecutionRow struct {
	ExecID               int64              `db:"exec_id" json:"exec_id"`
	JobID                pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID             pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt            pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt           pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	CreatedAt            pgtype.Timestamptz `db:"created_at" json:"created_at"`
	ExecRequestID        pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs             string             `db:"exec_logs" json:"exec_logs"`
	NixLogs              pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success              pgtype.Bool        `db:"success" json:"success"`
	ID                   int32              `db:"id" json:"id"`
	Hash                 string             `db:"hash" json:"hash"`
	Code                 pgtype.Text        `db:"code" json:"code"`
	Flake                string             `db:"flake" json:"flake"`
	LanguageDependencies []string           `db:"language_dependencies" json:"language_dependencies"`
	SystemDependencies   []string           `db:"system_dependencies" json:"system_dependencies"`
	CmdLineArgs          pgtype.Text        `db:"cmd_line_args" json:"cmd_line_args"`
	CompileArgs          pgtype.Text        `db:"compile_args" json:"compile_args"`
	Files                []byte             `db:"files" json:"files"`
	Input                pgtype.Text        `db:"input" json:"input"`
	Command              pgtype.Text        `db:"command" json:"command"`
	Setup                pgtype.Text        `db:"setup" json:"setup"`
	LanguageVersion      int64              `db:"language_version" json:"language_version"`
}

func (q *Queries) GetLatestExecution(ctx context.Context, jobID pgtype.Int8) (GetLatestExecutionRow, error) {
	row := q.db.QueryRow(ctx, getLatestExecution, jobID)
	var i GetLatestExecutionRow
	err := row.Scan(
		&i.ExecID,
		&i.JobID,
		&i.WorkerID,
		&i.StartedAt,
		&i.FinishedAt,
		&i.CreatedAt,
		&i.ExecRequestID,
		&i.ExecLogs,
		&i.NixLogs,
		&i.Success,
		&i.ID,
		&i.Hash,
		&i.Code,
		&i.Flake,
		&i.LanguageDependencies,
		&i.SystemDependencies,
		&i.CmdLineArgs,
		&i.CompileArgs,
		&i.Files,
		&i.Input,
		&i.Command,
		&i.Setup,
		&i.LanguageVersion,
	)
	return i, err
}

const getTotalExecutions = `-- name: GetTotalExecutions :one
select count(*) from executions
`

func (q *Queries) GetTotalExecutions(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalExecutions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalExecutionsForJob = `-- name: GetTotalExecutionsForJob :one
select count(*) from executions where job_id = $1
`

func (q *Queries) GetTotalExecutionsForJob(ctx context.Context, jobID pgtype.Int8) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalExecutionsForJob, jobID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalJobs = `-- name: GetTotalJobs :one
SELECT count(*) FROM jobs
`

func (q *Queries) GetTotalJobs(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalJobs)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertExecution = `-- name: InsertExecution :one
insert into executions
    (job_id, worker_id, started_at, finished_at, exec_request_id, exec_logs, nix_logs, success)
values
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning exec_id, job_id, worker_id, started_at, finished_at, created_at, exec_request_id, exec_logs, nix_logs, success
`

type InsertExecutionParams struct {
	JobID         pgtype.Int8        `db:"job_id" json:"job_id"`
	WorkerID      pgtype.Int4        `db:"worker_id" json:"worker_id"`
	StartedAt     pgtype.Timestamptz `db:"started_at" json:"started_at"`
	FinishedAt    pgtype.Timestamptz `db:"finished_at" json:"finished_at"`
	ExecRequestID pgtype.Int4        `db:"exec_request_id" json:"exec_request_id"`
	ExecLogs      string             `db:"exec_logs" json:"exec_logs"`
	NixLogs       pgtype.Text        `db:"nix_logs" json:"nix_logs"`
	Success       pgtype.Bool        `db:"success" json:"success"`
}

func (q *Queries) InsertExecution(ctx context.Context, arg InsertExecutionParams) (Execution, error) {
	row := q.db.QueryRow(ctx, insertExecution,
		arg.JobID,
		arg.WorkerID,
		arg.StartedAt,
		arg.FinishedAt,
		arg.ExecRequestID,
		arg.ExecLogs,
		arg.NixLogs,
		arg.Success,
	)
	var i Execution
	err := row.Scan(
		&i.ExecID,
		&i.JobID,
		&i.WorkerID,
		&i.StartedAt,
		&i.FinishedAt,
		&i.CreatedAt,
		&i.ExecRequestID,
		&i.ExecLogs,
		&i.NixLogs,
		&i.Success,
	)
	return i, err
}

const insertJob = `-- name: InsertJob :one
insert into jobs
    (arguments, max_retries, time_out, job_type)
values
    ($1, $2, $3, $4)
returning job_id, created_at, updated_at, time_out, started_at, arguments, current_state, retries, max_retries, worker_id, job_type
`

type InsertJobParams struct {
	Arguments  jsonschema.JobArguments `db:"arguments" json:"arguments"`
	MaxRetries pgtype.Int4             `db:"max_retries" json:"max_retries"`
	TimeOut    pgtype.Int4             `db:"time_out" json:"time_out"`
	JobType    string                  `db:"job_type" json:"job_type"`
}

func (q *Queries) InsertJob(ctx context.Context, arg InsertJobParams) (Job, error) {
	row := q.db.QueryRow(ctx, insertJob,
		arg.Arguments,
		arg.MaxRetries,
		arg.TimeOut,
		arg.JobType,
	)
	var i Job
	err := row.Scan(
		&i.JobID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TimeOut,
		&i.StartedAt,
		&i.Arguments,
		&i.CurrentState,
		&i.Retries,
		&i.MaxRetries,
		&i.WorkerID,
		&i.JobType,
	)
	return i, err
}

const pruneCompletedJobs = `-- name: PruneCompletedJobs :exec
delete from jobs where current_state = 'completed'
`

func (q *Queries) PruneCompletedJobs(ctx context.Context) error {
	_, err := q.db.Exec(ctx, pruneCompletedJobs)
	return err
}

const requeueLTJobs = `-- name: RequeueLTJobs :exec
update jobs
set
    current_state = 'pending',
    updated_at = now(),
    started_at = null,
    worker_id = null,
    retries = retries::integer + 1
where current_state = 'scheduled' 
  and started_at + time_out * INTERVAL '1 second' < now() and time_out > 0
`

func (q *Queries) RequeueLTJobs(ctx context.Context) error {
	_, err := q.db.Exec(ctx, requeueLTJobs)
	return err
}

const requeueWorkerJobs = `-- name: RequeueWorkerJobs :exec
update jobs
set
    current_state = 'pending',
    worker_id = null,
    started_at = null,
    retries = retries::integer + 1,
    updated_at = now()
where current_state = 'scheduled' 
  and worker_id = $1
`

func (q *Queries) RequeueWorkerJobs(ctx context.Context, workerID pgtype.Int4) error {
	_, err := q.db.Exec(ctx, requeueWorkerJobs, workerID)
	return err
}

const retryJob = `-- name: RetryJob :exec
update jobs
set
    current_state = 'pending',
    retries = retries::integer + 1,
    started_at = null,
    updated_at = now(),
    worker_id = null
where job_id = $1 AND current_state = 'scheduled'
`

func (q *Queries) RetryJob(ctx context.Context, jobID int64) error {
	_, err := q.db.Exec(ctx, retryJob, jobID)
	return err
}

const stopJob = `-- name: StopJob :exec
update jobs set current_state = 'pending', updated_at = now(), worker_id = null where job_id = $1
`

func (q *Queries) StopJob(ctx context.Context, jobID int64) error {
	_, err := q.db.Exec(ctx, stopJob, jobID)
	return err
}

const updateJobCompleted = `-- name: UpdateJobCompleted :exec
update jobs
set
    current_state = 'completed',
    updated_at = now()
where job_id = $1 AND current_state = 'scheduled'
`

func (q *Queries) UpdateJobCompleted(ctx context.Context, jobID int64) error {
	_, err := q.db.Exec(ctx, updateJobCompleted, jobID)
	return err
}

const updateJobFailed = `-- name: updateJobFailed :exec
update jobs
set
    current_state = 'failed',
    updated_at = now(),
    retries = retries::integer + 1
where job_id = $1 AND current_state = 'scheduled'
`

func (q *Queries) updateJobFailed(ctx context.Context, jobID int64) error {
	_, err := q.db.Exec(ctx, updateJobFailed, jobID)
	return err
}
