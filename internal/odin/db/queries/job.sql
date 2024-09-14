create table jobs (
    job_id bigint primary key default nextval('jobs_id_seq'),
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

create table executions (
    exec_id bigint primary key default nextval('executions_id_seq'),
    job_id bigint not null references jobs on delete set null,
    worker_id int not null references workers on delete set null,
    created_at timestamptz not null default now(),
    started_at timestamptz not null,
    finished_at timestamptz not null,
    exec_request_id int references exec_request on delete set null,
    exec_logs text not null,
    nix_logs text,
    success boolean
);

-- name: FetchJob :one
update jobs set current_state = 'scheduled', started_at = now(), worker_id = $1, updated_at = now()
where job_id = (
    select job_id from jobs
    where 
        current_state = 'pending'
        and retries < max_retries
    order by
        job_id asc
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
where job_id = $1 AND current_state = 'scheduled';

-- name: InsertExecution :one
insert into executions
    (job_id, worker_id, started_at, finished_at, exec_request_id, exec_logs, nix_logs, success)
values
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetAllJobs :many
select * from jobs
inner join exec_request on jobs.exec_request_id = exec_request.id
order by jobs.job_id
limit $1 offset $2;

-- name: GetJob :one
select * from jobs inner join exec_request on jobs.exec_request_id = exec_request.id where jobs.job_id = $1;

-- name: GetExecution :one
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.exec_id = $1;

-- name: GetJobState :one
select current_state from jobs where job_id = $1;

-- name: DeleteJob :one
delete from jobs where job_id = $1 and completed = false and current_state in ('pending', 'cancelled', 'failed') returning job_id;

-- name: GetTotalJobs :one
SELECT count(*) FROM jobs;

-- name: GetExecutionsForJob :many
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.job_id = $1
order by finished_at desc
limit $2 offset $3;

-- name: GetAllExecutions :many
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
order by started_at desc
limit $1 offset $2;

-- name: GetTotalExecutionsForJob :one
select count(*) from executions where job_id = $1;

-- name: GetTotalExecutions :one
select count(*) from executions;

-- name: StopJob :exec
update jobs set current_state = 'pending', updated_at = now(), worker_id = null where job_id = $1;

-- name: CancelJob :exec
update jobs set current_state = 'cancelled', updated_at = now(), worker_id = null where job_id = $1;

-- name: updateJobFailed :exec
update jobs
set
    current_state = 'failed',
    updated_at = now(),
    retries = retries::integer + 1
where job_id = $1 AND current_state = 'scheduled';

-- name: RetryJob :exec
update jobs
set
    current_state = 'pending',
    retries = retries::integer + 1,
    started_at = null,
    updated_at = now(),
    worker_id = null
where job_id = $1 AND current_state = 'scheduled';

-- name: PruneCompletedJobs :exec
delete from jobs where current_state = 'completed';

-- name: RequeueLTJobs :exec
update jobs
set
    current_state = 'pending',
    updated_at = now(),
    started_at = null,
    worker_id = null,
    retries = retries::integer + 1
where current_state = 'scheduled' 
  and started_at + time_out * INTERVAL '1 second' < now() and time_out > 0;

-- name: RequeueWorkerJobs :exec
update jobs
set
    current_state = 'pending',
    worker_id = null,
    started_at = null,
    retries = retries::integer + 1,
    updated_at = now()
where current_state = 'scheduled' 
  and worker_id = $1;
