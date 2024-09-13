create table jobs (
    id bigint primary key default nextval('jobs_id_seq'),
    created_at timestamptz not null  default now(),
    updated_at timestamptz,
    time_out int,
    started_at timestamptz,
    exec_request_id int references exec_request on delete set null,
    current_state TEXT NOT NULL CHECK (current_state IN ('pending', 'scheduled', 'completed', 'failed', 'cancelled')) DEFAULT 'pending',
    retries int default 0,
    max_retries int default 5,
    worker_id int references workers on delete set null
);

create table job_runs (
    id bigint primary key default nextval('job_runs_id_seq'),
    job_id bigint not null references jobs on delete set null,
    worker_id int not null references workers on delete set null,
    started_at timestamptz not null,
    finished_at timestamptz not null,
    exec_request_id int references exec_request on delete set null,
    exec_logs text not null,
    nix_logs text,
    success boolean
);

-- name: FetchJob :one
update jobs set current_state = 'scheduled', started_at = now(), worker_id = $1, updated_at = now()
where id = (
    select id from jobs
    where 
        current_state = 'pending'
        and retries < max_retries
    order by
        id asc
    for update skip locked
    limit 1
    )
returning *;

-- name: InsertJob :one
insert into jobs
    (exec_request_id, max_retries, time_out)
values
    ($1, $2, $3)
returning *;

-- name: UpdateJobCompleted :exec
update jobs
set
    current_state = 'completed',
    updated_at = now()
where id = $1 AND current_state = 'scheduled';

-- name: InsertJobRun :one
insert into job_runs
    (job_id, worker_id, started_at, finished_at, exec_request_id, exec_logs, nix_logs, success)
values
    ($1, $2, $3, $4, $5, $6, $7, $8)
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
order by finished_at desc
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
update jobs set current_state = 'pending', updated_at = now(), worker_id = null where id = $1;

-- name: CancelJob :exec
update jobs set current_state = 'cancelled', updated_at = now(), worker_id = null where id = $1;

-- name: updateJobFailed :exec
update jobs
set
    current_state = 'failed',
    updated_at = now()
where id = $1 AND current_state = 'scheduled';

-- name: RetryJob :exec
update jobs
set
    current_state = 'pending',
    retries = retries::integer + 1,
    updated_at = now(),
    worker_id = null
where id = $1 AND current_state = 'scheduled';

-- name: PruneCompletedJobs :exec
delete from jobs where current_state = 'completed';

-- name: RequeueLTJobs :exec
update jobs
set
    current_state = 'pending',
    updated_at = now(),
    worker_id = null
where current_state = 'scheduled' 
  and started_at + time_out * INTERVAL '1 second' < now() and time_out > 0;

-- name: RequeueWorkerJobs :exec
update jobs
set
    current_state = 'pending',
    worker_id = null,
    updated_at = now()
where current_state = 'scheduled' 
  and worker_id = $1;
