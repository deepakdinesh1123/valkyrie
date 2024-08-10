create table workers (
    id int primary key,
    name text not null unique,
    created_at timestamptz not null default now()
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