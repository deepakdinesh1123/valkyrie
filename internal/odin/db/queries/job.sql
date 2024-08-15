create table jobs (
    id bigint primary key default nextval('jobs_id_seq'),
    inserted_at timestamptz not null default now(),
    -- group_id int not null, will be added later
    -- type_id int not null, will be added later
    worker_id int references workers on delete set null,
    script text not null,
    script_path text not null,
    args varchar(1024),
    flake text not null,
    language text not null,
    completed bool not null default false,
    running bool not null default false
);

create table job_runs (
    id bigint primary key default nextval('job_runs_id_seq'),
    job_id bigint not null,
    worker_id int not null references workers on delete set null,
    started_at timestamptz not null,
    finished_at timestamptz not null,
    created_at timestamptz not null,
    -- group_id int not null,
    -- type_id int not null,
    script text not null,
    flake text not null,
    args varchar(1024),
    logs text not null
);

-- name: FetchJob :one
update jobs set running = true, worker_id = $1
where id = (
    select id from jobs
    where (running = false and completed = false)
    order by
        id asc
    for update skip locked
    limit 1
    )
returning *;

-- name: InsertJob :one
insert into jobs
    (script, flake, language, script_path, args)
values
    ($1, $2, $3, $4, $5)
returning *;

-- name: UpdateJob :exec
update jobs
set
    completed = true
where id = $1 AND completed = false;

-- name: InsertJobRun :one
insert into job_runs
    (job_id, worker_id, started_at, finished_at, script, flake, args, logs, created_at)
values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: GetAllJobs :many
select * from jobs
limit $1 offset $2;

-- name: GetJob :one
select * from jobs where id = $1;

-- name: DeleteJob :exec
delete from jobs where id = $1 and completed = false;

-- name: GetTotalJobs :one
SELECT count(*) FROM jobs;

-- name: GetExecutionResultsByID :many
select * from job_runs where job_id = $1
limit $2 offset $3;

-- name: GetAllExecutionResults :many
select * from job_runs
order by started_at desc
limit $1 offset $2;

-- name: GetTotalExecutionsForJob :one
select count(*) from job_runs where job_id = $1;

-- name: GetTotalExecutions :one
select count(*) from job_runs;

-- name: StopJob :exec
update jobs set running = false, worker_id = null where id = $1;