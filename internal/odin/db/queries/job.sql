-- name: GetFlake :one
SELECT flake 
FROM exec_request 
WHERE id = (
    SELECT CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) 
    FROM jobs 
    WHERE job_id = $1
);

-- name: FetchJob :one
with cte as (
    select job_id
    from jobs
    where 
        current_state = 'pending'
        and job_type = @JobType::text
        and retries < max_retries
    order by job_id asc
    for update skip locked
    limit 1
)
update jobs
set current_state = 'scheduled', 
    started_at = now(), 
    worker_id = @WorkerId::int, 
    updated_at = now()
where job_id = (select job_id from cte)
returning *;

-- name: InsertJob :one
insert into jobs
    (arguments, max_retries, time_out, job_type)
values
    ($1, $2, $3, $4)
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

-- name: GetAllExecutionJobs :many
SELECT * 
FROM jobs
INNER JOIN exec_request 
ON CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) = exec_request.id
WHERE job_id >= $1
ORDER BY jobs.job_id
LIMIT $2;

-- name: GetExecutionJob :one
select * from jobs inner join exec_request on CAST(arguments->'ExecConfig'->>'exec_req_id' AS INT) = exec_request.id where jobs.job_id = $1;

-- name: GetExecution :one
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.exec_id = $1;

-- name: GetLatestExecution :one
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.job_id = $1
order by finished_at desc
limit 1;

-- name: GetJobState :one
select current_state from jobs where job_id = $1;

-- name: DeleteJob :one
delete from jobs where job_id = $1 and current_state in ('pending', 'cancelled', 'failed') returning job_id;

-- name: GetTotalJobs :one
SELECT count(*) FROM jobs;

-- name: GetExecutionsForJob :many
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where executions.job_id = $1 and exec_id >= $2
order by finished_at desc
limit $3;

-- name: GetAllExecutions :many
select * from executions
inner join exec_request on executions.exec_request_id = exec_request.id
where exec_id >= $1
order by started_at desc
limit $2;

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



