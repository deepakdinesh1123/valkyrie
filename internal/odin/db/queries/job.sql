create table jobs (
    id bigint primary key default nextval('jobs_id_seq'),
    created_at timestamptz not null  default now(),
    updated_at timestamptz,
    exec_request_id int references exec_request on delete set null,
    status TEXT NOT NULL CHECK (status IN ('pending', 'scheduled', 'completed', 'failed', 'cancelled')) DEFAULT 'pending',
    retries int default 0,
    max_retries int default 5
);

create table job_runs (
    id bigint primary key default nextval('job_runs_id_seq'),
    job_id bigint not null references jobs on delete set null,
    worker_id int not null references workers on delete set null,
    started_at timestamptz not null ,
    finished_at timestamptz not null ,
    exec_request_id int references exec_request on delete set null,
    logs text not null
);

-- name: FetchJob :one
update jobs set status = 'scheduled'
where id = (
    select id from jobs
    where 
        status = 'pending'
        and retries < max_retries
    order by
        id asc
    for update skip locked
    limit 1
    )
returning *;

-- name: InsertJob :one
insert into jobs
    (exec_request_id, max_retries)
values
    ($1, $2)
returning *;

-- name: UpdateJobCompleted :exec
update jobs
set
    status = 'completed',
    updated_at = now()
where id = $1 AND status = 'scheduled';

-- name: InsertJobRun :one
insert into job_runs
    (job_id, worker_id, started_at, finished_at, exec_request_id, logs)
values
    ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetAllJobs :many
select * from jobs
inner join exec_request on jobs.exec_request_id = exec_request.id
order by jobs.id
limit $1 offset $2;

-- name: GetJob :one
select * from jobs inner join exec_request on jobs.exec_request_id = exec_request.id where jobs.id = $1;

-- name: DeleteJob :exec
delete from jobs where id = $1 and completed = false;

-- name: GetTotalJobs :one
SELECT count(*) FROM jobs;

-- name: GetExecutionResultsByID :many
select * from job_runs
inner join exec_request on job_runs.exec_request_id = exec_request.id
where job_runs.job_id = $1
limit $2 offset $3;

-- name: GetAllExecutionResults :many
select * from job_runs
inner join exec_request on job_runs.exec_request_id = exec_request.id
order by started_at desc
limit $1 offset $2;

-- name: GetTotalExecutionsForJob :one
select count(*) from job_runs where job_id = $1;

-- name: GetTotalExecutions :one
select count(*) from job_runs;

-- name: StopJob :exec
update jobs set status = 'pending' where id = $1;

-- name: CancelJob :exec
update jobs set status = 'cancelled' where id = $1;

-- name: updateJobFailed :exec
update jobs
set
    status = 'failed',
    updated_at = now()
where id = $1 AND status = 'scheduled';

-- name: RetryJob :exec
update jobs
set
    status = 'pending',
    retries = retries + 1,
    updated_at = now()
where id = $1 AND status = 'scheduled';

-- name: PruneCompletedJobs :exec
delete from jobs where status = 'completed';