-- name: CreateWorker :one
insert into workers
    (name)
values
    ($1)
returning *;

-- name: InsertWorker :one
insert into workers
    (id, name)
values
    ($1, $2)
returning *;

-- name: GetWorker :one
update workers
    set current_state = 'active'
where name = $1
returning *;

-- name: GetAllWorkers :many
select * from workers;

-- name: GetTotalWorkers :one
select count(*) from workers;

-- name: UpdateHeartbeat :exec
update workers
set
    last_heartbeat = now()
where id = $1;

-- name: GetStaleWorkers :many
update workers
    set current_state='stale'
where last_heartbeat < now() - interval '20 seconds' and current_state != 'stale'
returning id;

-- name: DeleteWorker :exec
delete from workers where id = $1;

-- name: WorkerTaskCount :one
select count(*) from jobs
where current_state = 'scheduled' and worker_id = $1;
