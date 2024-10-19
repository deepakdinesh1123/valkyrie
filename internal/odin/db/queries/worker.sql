create table workers (
    id int primary key,
    name text not null unique,
    created_at timestamptz not null default now(),
    last_heartbeat timestamptz
);

-- name: InsertWorker :one
insert into workers
    (name)
values
    ($1)
returning *;

-- name: GetWorker :one
select * from workers where name = $1;

-- name: GetAllWorkers :many
select * from workers
limit $1 offset $2;

-- name: GetTotalWorkers :one
select count(*) from workers;

-- name: UpdateHeartbeat :exec
update workers
set
    last_heartbeat = now()
where id = $1;

-- name: GetStaleWorkers :many
select id from workers
where last_heartbeat < now() - interval '1 minute';

-- name: DeleteWorker :exec
delete from workers where id = $1;

-- name: WorkerTaskCount :one
select count(*) from jobs
where status = 'scheduled' and worker_id = $1;